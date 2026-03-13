package display

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"forge.lthn.ai/core/go-config"
	"forge.lthn.ai/core/go/pkg/core"
	"encoding/json"

	"forge.lthn.ai/core/gui/pkg/browser"
	"forge.lthn.ai/core/gui/pkg/contextmenu"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/dock"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/keybinding"
	"forge.lthn.ai/core/gui/pkg/lifecycle"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/notification"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Options holds configuration for the display service.
type Options struct{}

// WindowInfo is an alias for window.WindowInfo (backward compatibility).
type WindowInfo = window.WindowInfo

// Service manages windowing, dialogs, and other visual elements.
// It orchestrates sub-services (window, systray, menu) via IPC and bridges
// IPC actions to WebSocket events for TypeScript apps.
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	config     Options
	configData map[string]map[string]any
	cfg        *config.Config // go-config instance for file persistence
	events *WSEventManager
}

// New is the constructor for the display service.
func New() (*Service, error) {
	return &Service{
		configData: map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
		},
	}, nil
}

// Register creates a factory closure that captures the Wails app.
// Pass nil for testing without a Wails runtime.
func Register(wailsApp *application.App) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		s, err := New()
		if err != nil {
			return nil, err
		}
		s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
		s.wailsApp = wailsApp
		return s, nil
	}
}

// OnStartup loads config and registers IPC handlers synchronously.
// CRITICAL: config handlers MUST be registered before returning —
// sub-services depend on them during their own OnStartup.
func (s *Service) OnStartup(ctx context.Context) error {
	s.loadConfig()

	// Register config query/task handlers — available NOW for sub-services
	s.Core().RegisterQuery(s.handleConfigQuery)
	s.Core().RegisterTask(s.handleConfigTask)

	// Initialise Wails wrappers if app is available (nil in tests)
	if s.wailsApp != nil {
		s.app = newWailsApp(s.wailsApp)
		s.events = NewWSEventManager()
	}

	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
// It bridges sub-service IPC actions to WebSocket events for TS apps.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	switch m := msg.(type) {
	case core.ActionServiceStartup:
		// All services have completed OnStartup — safe to PERFORM on sub-services
		s.buildMenu()
		s.setupTray()
	case window.ActionWindowOpened:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowCreate, Window: m.Name,
				Data: map[string]any{"name": m.Name}})
		}
	case window.ActionWindowClosed:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowClose, Window: m.Name,
				Data: map[string]any{"name": m.Name}})
		}
	case window.ActionWindowMoved:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowMove, Window: m.Name,
				Data: map[string]any{"x": m.X, "y": m.Y}})
		}
	case window.ActionWindowResized:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowResize, Window: m.Name,
				Data: map[string]any{"w": m.W, "h": m.H}})
		}
	case window.ActionWindowFocused:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowFocus, Window: m.Name})
		}
	case window.ActionWindowBlurred:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowBlur, Window: m.Name})
		}
	case systray.ActionTrayClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventTrayClick})
		}
	case systray.ActionTrayMenuItemClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventTrayMenuItemClick,
				Data: map[string]any{"actionId": m.ActionID}})
		}
		s.handleTrayAction(m.ActionID)
	case environment.ActionThemeChanged:
		if s.events != nil {
			theme := "light"
			if m.IsDark {
				theme = "dark"
			}
			s.events.Emit(Event{Type: EventThemeChange,
				Data: map[string]any{"isDark": m.IsDark, "theme": theme}})
		}
	case notification.ActionNotificationClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationClick,
				Data: map[string]any{"id": m.ID}})
		}
	case screen.ActionScreensChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventScreenChange,
				Data: map[string]any{"screens": m.Screens}})
		}
	case keybinding.ActionTriggered:
		if s.events != nil {
			s.events.Emit(Event{Type: EventKeybindingTriggered,
				Data: map[string]any{"accelerator": m.Accelerator}})
		}
	case window.ActionFilesDropped:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowFileDrop, Window: m.Name,
				Data: map[string]any{"paths": m.Paths, "targetId": m.TargetID}})
		}
	case dock.ActionVisibilityChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventDockVisibility,
				Data: map[string]any{"visible": m.Visible}})
		}
	case lifecycle.ActionApplicationStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventAppStarted})
		}
	case lifecycle.ActionOpenedWithFile:
		if s.events != nil {
			s.events.Emit(Event{Type: EventAppOpenedWithFile,
				Data: map[string]any{"path": m.Path}})
		}
	case lifecycle.ActionWillTerminate:
		if s.events != nil {
			s.events.Emit(Event{Type: EventAppWillTerminate})
		}
	case lifecycle.ActionDidBecomeActive:
		if s.events != nil {
			s.events.Emit(Event{Type: EventAppActive})
		}
	case lifecycle.ActionDidResignActive:
		if s.events != nil {
			s.events.Emit(Event{Type: EventAppInactive})
		}
	case lifecycle.ActionPowerStatusChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventSystemPowerChange})
		}
	case lifecycle.ActionSystemSuspend:
		if s.events != nil {
			s.events.Emit(Event{Type: EventSystemSuspend})
		}
	case lifecycle.ActionSystemResume:
		if s.events != nil {
			s.events.Emit(Event{Type: EventSystemResume})
		}
	case contextmenu.ActionItemClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventContextMenuClick,
				Data: map[string]any{
					"menuName": m.MenuName,
					"actionId": m.ActionID,
					"data":     m.Data,
				}})
		}
	case ActionIDECommand:
		if s.events != nil {
			s.events.Emit(Event{Type: EventIDECommand,
				Data: map[string]any{"command": m.Command}})
		}
	}
	return nil
}

// WSMessage represents a command received from a WebSocket client.
type WSMessage struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data,omitempty"`
}

// handleWSMessage bridges WebSocket commands to IPC calls.
func (s *Service) handleWSMessage(msg WSMessage) (any, bool, error) {
	var result any
	var handled bool
	var err error

	switch msg.Action {
	case "keybinding:add":
		accelerator, _ := msg.Data["accelerator"].(string)
		description, _ := msg.Data["description"].(string)
		result, handled, err = s.Core().PERFORM(keybinding.TaskAdd{
			Accelerator: accelerator, Description: description,
		})
	case "keybinding:remove":
		accelerator, _ := msg.Data["accelerator"].(string)
		result, handled, err = s.Core().PERFORM(keybinding.TaskRemove{
			Accelerator: accelerator,
		})
	case "keybinding:list":
		result, handled, err = s.Core().QUERY(keybinding.QueryList{})
	case "browser:open-url":
		url, _ := msg.Data["url"].(string)
		result, handled, err = s.Core().PERFORM(browser.TaskOpenURL{URL: url})
	case "browser:open-file":
		path, _ := msg.Data["path"].(string)
		result, handled, err = s.Core().PERFORM(browser.TaskOpenFile{Path: path})
	case "dock:show":
		result, handled, err = s.Core().PERFORM(dock.TaskShowIcon{})
	case "dock:hide":
		result, handled, err = s.Core().PERFORM(dock.TaskHideIcon{})
	case "dock:badge":
		label, _ := msg.Data["label"].(string)
		result, handled, err = s.Core().PERFORM(dock.TaskSetBadge{Label: label})
	case "dock:badge-remove":
		result, handled, err = s.Core().PERFORM(dock.TaskRemoveBadge{})
	case "dock:visible":
		result, handled, err = s.Core().QUERY(dock.QueryVisible{})
	case "contextmenu:add":
		name, _ := msg.Data["name"].(string)
		menuJSON, _ := json.Marshal(msg.Data["menu"])
		var menuDef contextmenu.ContextMenuDef
		_ = json.Unmarshal(menuJSON, &menuDef)
		result, handled, err = s.Core().PERFORM(contextmenu.TaskAdd{
			Name: name, Menu: menuDef,
		})
	case "contextmenu:remove":
		name, _ := msg.Data["name"].(string)
		result, handled, err = s.Core().PERFORM(contextmenu.TaskRemove{Name: name})
	case "contextmenu:get":
		name, _ := msg.Data["name"].(string)
		result, handled, err = s.Core().QUERY(contextmenu.QueryGet{Name: name})
	case "contextmenu:list":
		result, handled, err = s.Core().QUERY(contextmenu.QueryList{})
	default:
		return nil, false, nil
	}

	return result, handled, err
}

// handleTrayAction processes tray menu item clicks.
func (s *Service) handleTrayAction(actionID string) {
	switch actionID {
	case "open-desktop":
		// Show all windows
		infos := s.ListWindowInfos()
		for _, info := range infos {
			_, _, _ = s.Core().PERFORM(window.TaskFocus{Name: info.Name})
		}
	case "close-desktop":
		// Hide all windows — future: add TaskHideWindow
	case "env-info":
		// Query environment info via IPC and show as dialog
		result, handled, _ := s.Core().QUERY(environment.QueryInfo{})
		if handled {
			info := result.(environment.EnvironmentInfo)
			details := fmt.Sprintf("OS: %s\nArch: %s\nPlatform: %s %s",
				info.OS, info.Arch, info.Platform.Name, info.Platform.Version)
			_, _, _ = s.Core().PERFORM(dialog.TaskMessageDialog{
				Opts: dialog.MessageDialogOptions{
					Type: dialog.DialogInfo, Title: "Environment",
					Message: details, Buttons: []string{"OK"},
				},
			})
		}
	case "quit":
		if s.app != nil {
			s.app.Quit()
		}
	}
}

func guiConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".core", "gui", "config.yaml")
	}
	return filepath.Join(home, ".core", "gui", "config.yaml")
}

func (s *Service) loadConfig() {
	if s.cfg != nil {
		return // Already loaded (e.g., via loadConfigFrom in tests)
	}
	s.loadConfigFrom(guiConfigPath())
}

func (s *Service) loadConfigFrom(path string) {
	cfg, err := config.New(config.WithPath(path))
	if err != nil {
		// Non-critical — continue with empty configData
		return
	}
	s.cfg = cfg

	for _, section := range []string{"window", "systray", "menu"} {
		var data map[string]any
		if err := cfg.Get(section, &data); err == nil && data != nil {
			s.configData[section] = data
		}
	}
}

func (s *Service) handleConfigQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case window.QueryConfig:
		return s.configData["window"], true, nil
	case systray.QueryConfig:
		return s.configData["systray"], true, nil
	case menu.QueryConfig:
		return s.configData["menu"], true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleConfigTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case window.TaskSaveConfig:
		s.configData["window"] = t.Value
		s.persistSection("window", t.Value)
		return nil, true, nil
	case systray.TaskSaveConfig:
		s.configData["systray"] = t.Value
		s.persistSection("systray", t.Value)
		return nil, true, nil
	case menu.TaskSaveConfig:
		s.configData["menu"] = t.Value
		s.persistSection("menu", t.Value)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) persistSection(key string, value map[string]any) {
	if s.cfg == nil {
		return
	}
	_ = s.cfg.Set(key, value)
	_ = s.cfg.Commit()
}

// --- Service accessors ---

// windowService returns the window service from Core, or nil if not registered.
func (s *Service) windowService() *window.Service {
	svc, err := core.ServiceFor[*window.Service](s.Core(), "window")
	if err != nil {
		return nil
	}
	return svc
}

// --- Window Management (delegates via IPC) ---

// OpenWindow creates a new window via IPC.
func (s *Service) OpenWindow(opts ...window.WindowOption) error {
	_, _, err := s.Core().PERFORM(window.TaskOpenWindow{Opts: opts})
	return err
}

// GetWindowInfo returns information about a window via IPC.
func (s *Service) GetWindowInfo(name string) (*window.WindowInfo, error) {
	result, handled, err := s.Core().QUERY(window.QueryWindowByName{Name: name})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("window service not available")
	}
	info, _ := result.(*window.WindowInfo)
	return info, nil
}

// ListWindowInfos returns information about all tracked windows via IPC.
func (s *Service) ListWindowInfos() []window.WindowInfo {
	result, handled, _ := s.Core().QUERY(window.QueryWindowList{})
	if !handled {
		return nil
	}
	list, _ := result.([]window.WindowInfo)
	return list
}

// SetWindowPosition moves a window via IPC.
func (s *Service) SetWindowPosition(name string, x, y int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetPosition{Name: name, X: x, Y: y})
	return err
}

// SetWindowSize resizes a window via IPC.
func (s *Service) SetWindowSize(name string, width, height int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetSize{Name: name, W: width, H: height})
	return err
}

// SetWindowBounds sets both position and size of a window via IPC.
func (s *Service) SetWindowBounds(name string, x, y, width, height int) error {
	if _, _, err := s.Core().PERFORM(window.TaskSetPosition{Name: name, X: x, Y: y}); err != nil {
		return err
	}
	_, _, err := s.Core().PERFORM(window.TaskSetSize{Name: name, W: width, H: height})
	return err
}

// MaximizeWindow maximizes a window via IPC.
func (s *Service) MaximizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMaximise{Name: name})
	return err
}

// MinimizeWindow minimizes a window via IPC.
func (s *Service) MinimizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMinimise{Name: name})
	return err
}

// FocusWindow brings a window to the front via IPC.
func (s *Service) FocusWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskFocus{Name: name})
	return err
}

// CloseWindow closes a window via IPC.
func (s *Service) CloseWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskCloseWindow{Name: name})
	return err
}

// RestoreWindow restores a maximized/minimized window.
// Uses direct Manager access (no IPC task for restore yet).
func (s *Service) RestoreWindow(name string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Restore()
	return nil
}

// SetWindowVisibility shows or hides a window.
// Uses direct Manager access (no IPC task for visibility yet).
func (s *Service) SetWindowVisibility(name string, visible bool) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetVisibility(visible)
	return nil
}

// SetWindowAlwaysOnTop sets whether a window stays on top.
// Uses direct Manager access (no IPC task for always-on-top yet).
func (s *Service) SetWindowAlwaysOnTop(name string, alwaysOnTop bool) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetAlwaysOnTop(alwaysOnTop)
	return nil
}

// SetWindowTitle changes a window's title.
// Uses direct Manager access (no IPC task for title yet).
func (s *Service) SetWindowTitle(name string, title string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetTitle(title)
	return nil
}

// SetWindowFullscreen sets a window to fullscreen mode.
// Uses direct Manager access (no IPC task for fullscreen yet).
func (s *Service) SetWindowFullscreen(name string, fullscreen bool) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	if fullscreen {
		pw.Fullscreen()
	} else {
		pw.UnFullscreen()
	}
	return nil
}

// SetWindowBackgroundColour sets the background colour of a window.
// Uses direct Manager access (no IPC task for background colour yet).
func (s *Service) SetWindowBackgroundColour(name string, r, g, b, a uint8) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	pw, ok := ws.Manager().Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetBackgroundColour(r, g, b, a)
	return nil
}

// GetFocusedWindow returns the name of the currently focused window.
func (s *Service) GetFocusedWindow() string {
	infos := s.ListWindowInfos()
	for _, info := range infos {
		if info.Focused {
			return info.Name
		}
	}
	return ""
}

// GetWindowTitle returns the title of a window by name.
func (s *Service) GetWindowTitle(name string) (string, error) {
	info, err := s.GetWindowInfo(name)
	if err != nil {
		return "", err
	}
	if info == nil {
		return "", fmt.Errorf("window not found: %s", name)
	}
	return info.Name, nil // Wails v3 doesn't expose a title getter
}

// ResetWindowState clears saved window positions.
func (s *Service) ResetWindowState() error {
	ws := s.windowService()
	if ws != nil {
		ws.Manager().State().Clear()
	}
	return nil
}

// GetSavedWindowStates returns all saved window states.
func (s *Service) GetSavedWindowStates() map[string]window.WindowState {
	ws := s.windowService()
	if ws == nil {
		return nil
	}
	result := make(map[string]window.WindowState)
	for _, name := range ws.Manager().State().ListStates() {
		if state, ok := ws.Manager().State().GetState(name); ok {
			result[name] = state
		}
	}
	return result
}

// CreateWindowOptions contains options for creating a new window.
type CreateWindowOptions struct {
	Name   string `json:"name"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// CreateWindow creates a new window with the specified options.
func (s *Service) CreateWindow(opts CreateWindowOptions) (*window.WindowInfo, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("window name is required")
	}
	result, _, err := s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName(opts.Name),
			window.WithTitle(opts.Title),
			window.WithURL(opts.URL),
			window.WithSize(opts.Width, opts.Height),
			window.WithPosition(opts.X, opts.Y),
		},
	})
	if err != nil {
		return nil, err
	}
	info := result.(window.WindowInfo)
	return &info, nil
}

// --- Layout delegation ---

// SaveLayout saves the current window arrangement as a named layout.
func (s *Service) SaveLayout(name string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	states := make(map[string]window.WindowState)
	for _, n := range ws.Manager().List() {
		if pw, ok := ws.Manager().Get(n); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			states[n] = window.WindowState{X: x, Y: y, Width: w, Height: h, Maximized: pw.IsMaximised()}
		}
	}
	return ws.Manager().Layout().SaveLayout(name, states)
}

// RestoreLayout applies a saved layout.
func (s *Service) RestoreLayout(name string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	layout, ok := ws.Manager().Layout().GetLayout(name)
	if !ok {
		return fmt.Errorf("layout not found: %s", name)
	}
	for wName, state := range layout.Windows {
		if pw, ok := ws.Manager().Get(wName); ok {
			pw.SetPosition(state.X, state.Y)
			pw.SetSize(state.Width, state.Height)
			if state.Maximized {
				pw.Maximise()
			} else {
				pw.Restore()
			}
		}
	}
	return nil
}

// ListLayouts returns all saved layout names with metadata.
func (s *Service) ListLayouts() []window.LayoutInfo {
	ws := s.windowService()
	if ws == nil {
		return nil
	}
	return ws.Manager().Layout().ListLayouts()
}

// DeleteLayout removes a saved layout by name.
func (s *Service) DeleteLayout(name string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	ws.Manager().Layout().DeleteLayout(name)
	return nil
}

// GetLayout returns a specific layout by name.
func (s *Service) GetLayout(name string) *window.Layout {
	ws := s.windowService()
	if ws == nil {
		return nil
	}
	layout, ok := ws.Manager().Layout().GetLayout(name)
	if !ok {
		return nil
	}
	return &layout
}

// --- Tiling/snapping delegation ---

// TileWindows arranges windows in a tiled layout.
func (s *Service) TileWindows(mode window.TileMode, windowNames []string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	return ws.Manager().TileWindows(mode, windowNames, 1920, 1080) // TODO: use actual screen size
}

// SnapWindow snaps a window to a screen edge or corner.
func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	return ws.Manager().SnapWindow(name, position, 1920, 1080) // TODO: use actual screen size
}

// StackWindows arranges windows in a cascade pattern.
func (s *Service) StackWindows(windowNames []string, offsetX, offsetY int) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	return ws.Manager().StackWindows(windowNames, offsetX, offsetY)
}

// ApplyWorkflowLayout applies a predefined layout for a specific workflow.
func (s *Service) ApplyWorkflowLayout(workflow window.WorkflowLayout) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	return ws.Manager().ApplyWorkflow(workflow, ws.Manager().List(), 1920, 1080)
}

// GetEventManager returns the event manager for WebSocket event subscriptions.
func (s *Service) GetEventManager() *WSEventManager {
	return s.events
}

// --- Menu (handlers stay in display, structure delegated via IPC) ---

func (s *Service) buildMenu() {
	items := []menu.MenuItem{
		{Role: ptr(menu.RoleAppMenu)},
		{Role: ptr(menu.RoleFileMenu)},
		{Role: ptr(menu.RoleViewMenu)},
		{Role: ptr(menu.RoleEditMenu)},
		{Label: "Workspace", Children: []menu.MenuItem{
			{Label: "New...", OnClick: s.handleNewWorkspace},
			{Label: "List", OnClick: s.handleListWorkspaces},
		}},
		{Label: "Developer", Children: []menu.MenuItem{
			{Label: "New File", Accelerator: "CmdOrCtrl+N", OnClick: s.handleNewFile},
			{Label: "Open File...", Accelerator: "CmdOrCtrl+O", OnClick: s.handleOpenFile},
			{Label: "Save", Accelerator: "CmdOrCtrl+S", OnClick: s.handleSaveFile},
			{Type: "separator"},
			{Label: "Editor", OnClick: s.handleOpenEditor},
			{Label: "Terminal", OnClick: s.handleOpenTerminal},
			{Type: "separator"},
			{Label: "Run", Accelerator: "CmdOrCtrl+R", OnClick: s.handleRun},
			{Label: "Build", Accelerator: "CmdOrCtrl+B", OnClick: s.handleBuild},
		}},
		{Role: ptr(menu.RoleWindowMenu)},
		{Role: ptr(menu.RoleHelpMenu)},
	}

	// On non-macOS, remove the AppMenu role
	if runtime.GOOS != "darwin" {
		items = items[1:] // skip AppMenu
	}

	_, _, _ = s.Core().PERFORM(menu.TaskSetAppMenu{Items: items})
}

func ptr[T any](v T) *T { return &v }

// --- Menu handler methods ---

func (s *Service) handleNewWorkspace() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("workspace-new"),
			window.WithTitle("New Workspace"),
			window.WithURL("/workspace/new"),
			window.WithSize(500, 400),
		},
	})
}

func (s *Service) handleListWorkspaces() {
	ws := s.Core().Service("workspace")
	if ws == nil {
		return
	}
	lister, ok := ws.(interface{ ListWorkspaces() []string })
	if !ok {
		return
	}
	_ = lister.ListWorkspaces()
}

func (s *Service) handleNewFile() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle("New File - Editor"),
			window.WithURL("/#/developer/editor?new=true"),
			window.WithSize(1200, 800),
		},
	})
}

func (s *Service) handleOpenFile() {
	dialog := s.app.Dialog().OpenFile()
	dialog.SetTitle("Open File")
	dialog.CanChooseFiles(true)
	dialog.CanChooseDirectories(false)
	result, err := dialog.PromptForSingleSelection()
	if err != nil || result == "" {
		return
	}
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle(result + " - Editor"),
			window.WithURL("/#/developer/editor?file=" + result),
			window.WithSize(1200, 800),
		},
	})
}

func (s *Service) handleSaveFile()   { s.app.Event().Emit("ide:save") }
func (s *Service) handleOpenEditor() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle("Editor"),
			window.WithURL("/#/developer/editor"),
			window.WithSize(1200, 800),
		},
	})
}
func (s *Service) handleOpenTerminal() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("terminal"),
			window.WithTitle("Terminal"),
			window.WithURL("/#/developer/terminal"),
			window.WithSize(800, 500),
		},
	})
}
func (s *Service) handleRun()   { s.app.Event().Emit("ide:run") }
func (s *Service) handleBuild() { s.app.Event().Emit("ide:build") }

// --- Tray (setup delegated via IPC) ---

func (s *Service) setupTray() {
	_, _, _ = s.Core().PERFORM(systray.TaskSetTrayMenu{Items: []systray.TrayMenuItem{
		{Label: "Open Desktop", ActionID: "open-desktop"},
		{Label: "Close Desktop", ActionID: "close-desktop"},
		{Type: "separator"},
		{Label: "Environment Info", ActionID: "env-info"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}})
}
