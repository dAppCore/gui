package display

import (
	"context"
	"fmt"
	"runtime"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"
)

// Options holds configuration for the display service.
type Options struct{}

// Service manages windowing, dialogs, and other visual elements.
// It orchestrates sub-services (window, systray, menu) via IPC and bridges
// IPC actions to WebSocket events for TypeScript apps.
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	config     Options
	configData map[string]map[string]any
	windows    *window.Manager
	tray       *systray.Manager
	menus      *menu.Manager
	notifier   *notifications.NotificationService
	events     *WSEventManager
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
		s.events = NewWSEventManager(newWailsEventSource(s.wailsApp))
		s.events.SetupWindowEventListeners()
	}

	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
// It bridges sub-service IPC actions to WebSocket events for TS apps.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	if s.events == nil && s.wailsApp != nil {
		return nil // No WS event manager (testing without Wails)
	}

	switch m := msg.(type) {
	case core.ActionServiceStartup:
		// All services have completed OnStartup — safe to PERFORM on sub-services
		if s.menus != nil {
			s.buildMenu()
		}
		if s.tray != nil {
			s.setupTray()
		}
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
	}
	return nil
}

// handleTrayAction processes tray menu item clicks.
func (s *Service) handleTrayAction(actionID string) {
	switch actionID {
	case "open-desktop":
		// Show all windows
		if s.windows != nil {
			for _, name := range s.windows.List() {
				if pw, ok := s.windows.Get(name); ok {
					pw.Show()
				}
			}
		}
	case "close-desktop":
		// Hide all windows — future: add TaskHideWindow
		if s.windows != nil {
			for _, name := range s.windows.List() {
				if pw, ok := s.windows.Get(name); ok {
					pw.Hide()
				}
			}
		}
	case "env-info":
		if s.app != nil {
			s.ShowEnvironmentDialog()
		}
	case "quit":
		if s.app != nil {
			s.app.Quit()
		}
	}
}

func (s *Service) loadConfig() {
	// In-memory defaults. go-config integration is deferred work.
	if s.configData == nil {
		s.configData = map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
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
		return nil, true, nil
	case systray.TaskSaveConfig:
		s.configData["systray"] = t.Value
		return nil, true, nil
	case menu.TaskSaveConfig:
		s.configData["menu"] = t.Value
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

// --- Window Management (delegates to window.Manager) ---

// OpenWindow creates a new window with the given options.
func (s *Service) OpenWindow(opts ...window.WindowOption) error {
	pw, err := s.windows.Open(opts...)
	if err != nil {
		return err
	}
	s.trackWindow(pw)
	return nil
}

// trackWindow attaches event listeners for state persistence and WebSocket events.
func (s *Service) trackWindow(pw window.PlatformWindow) {
	if s.events != nil {
		s.events.EmitWindowEvent(EventWindowCreate, pw.Name(), map[string]any{
			"name": pw.Name(),
		})
		s.events.AttachWindowListeners(pw)
	}
}

// GetWindowInfo returns information about a window by name.
func (s *Service) GetWindowInfo(name string) (*WindowInfo, error) {
	pw, ok := s.windows.Get(name)
	if !ok {
		return nil, fmt.Errorf("window not found: %s", name)
	}
	x, y := pw.Position()
	w, h := pw.Size()
	return &WindowInfo{
		Name:      name,
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		Maximized: pw.IsMaximised(),
	}, nil
}

// ListWindowInfos returns information about all tracked windows.
func (s *Service) ListWindowInfos() []WindowInfo {
	names := s.windows.List()
	result := make([]WindowInfo, 0, len(names))
	for _, name := range names {
		if pw, ok := s.windows.Get(name); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			result = append(result, WindowInfo{
				Name:      name,
				X:         x,
				Y:         y,
				Width:     w,
				Height:    h,
				Maximized: pw.IsMaximised(),
			})
		}
	}
	return result
}

// SetWindowPosition moves a window to the specified position.
func (s *Service) SetWindowPosition(name string, x, y int) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetPosition(x, y)
	s.windows.State().UpdatePosition(name, x, y)
	return nil
}

// SetWindowSize resizes a window.
func (s *Service) SetWindowSize(name string, width, height int) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetSize(width, height)
	s.windows.State().UpdateSize(name, width, height)
	return nil
}

// SetWindowBounds sets both position and size of a window.
func (s *Service) SetWindowBounds(name string, x, y, width, height int) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetPosition(x, y)
	pw.SetSize(width, height)
	return nil
}

// MaximizeWindow maximizes a window.
func (s *Service) MaximizeWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Maximise()
	s.windows.State().UpdateMaximized(name, true)
	return nil
}

// RestoreWindow restores a maximized/minimized window.
func (s *Service) RestoreWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Restore()
	return nil
}

// MinimizeWindow minimizes a window.
func (s *Service) MinimizeWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Minimise()
	return nil
}

// FocusWindow brings a window to the front.
func (s *Service) FocusWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Focus()
	return nil
}

// CloseWindow closes a window by name.
func (s *Service) CloseWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	s.windows.State().CaptureState(pw)
	pw.Close()
	s.windows.Remove(name)
	return nil
}

// SetWindowVisibility shows or hides a window.
func (s *Service) SetWindowVisibility(name string, visible bool) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetVisibility(visible)
	return nil
}

// SetWindowAlwaysOnTop sets whether a window stays on top.
func (s *Service) SetWindowAlwaysOnTop(name string, alwaysOnTop bool) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetAlwaysOnTop(alwaysOnTop)
	return nil
}

// SetWindowTitle changes a window's title.
func (s *Service) SetWindowTitle(name string, title string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetTitle(title)
	return nil
}

// SetWindowFullscreen sets a window to fullscreen mode.
func (s *Service) SetWindowFullscreen(name string, fullscreen bool) error {
	pw, ok := s.windows.Get(name)
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
func (s *Service) SetWindowBackgroundColour(name string, r, g, b, a uint8) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetBackgroundColour(r, g, b, a)
	return nil
}

// GetFocusedWindow returns the name of the currently focused window.
func (s *Service) GetFocusedWindow() string {
	for _, name := range s.windows.List() {
		if pw, ok := s.windows.Get(name); ok {
			if pw.IsFocused() {
				return name
			}
		}
	}
	return ""
}

// GetWindowTitle returns the title of a window by name.
func (s *Service) GetWindowTitle(name string) (string, error) {
	_, ok := s.windows.Get(name)
	if !ok {
		return "", fmt.Errorf("window not found: %s", name)
	}
	return name, nil // Wails v3 doesn't expose a title getter
}

// ResetWindowState clears saved window positions.
func (s *Service) ResetWindowState() error {
	if s.windows != nil {
		s.windows.State().Clear()
	}
	return nil
}

// GetSavedWindowStates returns all saved window states.
func (s *Service) GetSavedWindowStates() map[string]window.WindowState {
	if s.windows == nil {
		return nil
	}
	result := make(map[string]window.WindowState)
	for _, name := range s.windows.State().ListStates() {
		if state, ok := s.windows.State().GetState(name); ok {
			result[name] = state
		}
	}
	return result
}

// WindowInfo contains information about a window for MCP.
type WindowInfo struct {
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Maximized bool   `json:"maximized"`
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
func (s *Service) CreateWindow(opts CreateWindowOptions) (*WindowInfo, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("window name is required")
	}
	err := s.OpenWindow(
		window.WithName(opts.Name),
		window.WithTitle(opts.Title),
		window.WithURL(opts.URL),
		window.WithSize(opts.Width, opts.Height),
		window.WithPosition(opts.X, opts.Y),
	)
	if err != nil {
		return nil, err
	}
	return &WindowInfo{
		Name:   opts.Name,
		X:      opts.X,
		Y:      opts.Y,
		Width:  opts.Width,
		Height: opts.Height,
	}, nil
}

// --- Layout delegation ---

// SaveLayout saves the current window arrangement as a named layout.
func (s *Service) SaveLayout(name string) error {
	if s.windows == nil {
		return fmt.Errorf("window manager not initialized")
	}
	states := make(map[string]window.WindowState)
	for _, n := range s.windows.List() {
		if pw, ok := s.windows.Get(n); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			states[n] = window.WindowState{X: x, Y: y, Width: w, Height: h, Maximized: pw.IsMaximised()}
		}
	}
	return s.windows.Layout().SaveLayout(name, states)
}

// RestoreLayout applies a saved layout.
func (s *Service) RestoreLayout(name string) error {
	if s.windows == nil {
		return fmt.Errorf("window manager not initialized")
	}
	layout, ok := s.windows.Layout().GetLayout(name)
	if !ok {
		return fmt.Errorf("layout not found: %s", name)
	}
	for wName, state := range layout.Windows {
		if pw, ok := s.windows.Get(wName); ok {
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
	if s.windows == nil {
		return nil
	}
	return s.windows.Layout().ListLayouts()
}

// DeleteLayout removes a saved layout by name.
func (s *Service) DeleteLayout(name string) error {
	if s.windows == nil {
		return fmt.Errorf("window manager not initialized")
	}
	s.windows.Layout().DeleteLayout(name)
	return nil
}

// GetLayout returns a specific layout by name.
func (s *Service) GetLayout(name string) *window.Layout {
	if s.windows == nil {
		return nil
	}
	layout, ok := s.windows.Layout().GetLayout(name)
	if !ok {
		return nil
	}
	return &layout
}

// --- Tiling/snapping delegation ---

// TileWindows arranges windows in a tiled layout.
func (s *Service) TileWindows(mode window.TileMode, windowNames []string) error {
	return s.windows.TileWindows(mode, windowNames, 1920, 1080) // TODO: use actual screen size
}

// SnapWindow snaps a window to a screen edge or corner.
func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	return s.windows.SnapWindow(name, position, 1920, 1080) // TODO: use actual screen size
}

// StackWindows arranges windows in a cascade pattern.
func (s *Service) StackWindows(windowNames []string, offsetX, offsetY int) error {
	return s.windows.StackWindows(windowNames, offsetX, offsetY)
}

// ApplyWorkflowLayout applies a predefined layout for a specific workflow.
func (s *Service) ApplyWorkflowLayout(workflow window.WorkflowLayout) error {
	return s.windows.ApplyWorkflow(workflow, s.windows.List(), 1920, 1080)
}

// --- Screen queries (remain in display — use application.Get() directly) ---

// ScreenInfo contains information about a display screen.
type ScreenInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Primary bool   `json:"primary"`
}

// WorkArea represents usable screen space.
type WorkArea struct {
	ScreenID string `json:"screenId"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// GetScreens returns information about all available screens.
func (s *Service) GetScreens() []ScreenInfo {
	app := application.Get()
	if app == nil || app.Screen == nil {
		return nil
	}
	screens := app.Screen.GetAll()
	if screens == nil {
		return nil
	}
	result := make([]ScreenInfo, 0, len(screens))
	for _, screen := range screens {
		result = append(result, ScreenInfo{
			ID:      screen.ID,
			Name:    screen.Name,
			X:       screen.Bounds.X,
			Y:       screen.Bounds.Y,
			Width:   screen.Bounds.Width,
			Height:  screen.Bounds.Height,
			Primary: screen.IsPrimary,
		})
	}
	return result
}

// GetWorkAreas returns the usable work area for all screens.
func (s *Service) GetWorkAreas() []WorkArea {
	app := application.Get()
	if app == nil || app.Screen == nil {
		return nil
	}
	screens := app.Screen.GetAll()
	if screens == nil {
		return nil
	}
	result := make([]WorkArea, 0, len(screens))
	for _, screen := range screens {
		result = append(result, WorkArea{
			ScreenID: screen.ID,
			X:        screen.WorkArea.X,
			Y:        screen.WorkArea.Y,
			Width:    screen.WorkArea.Width,
			Height:   screen.WorkArea.Height,
		})
	}
	return result
}

// GetPrimaryScreen returns information about the primary screen.
func (s *Service) GetPrimaryScreen() (*ScreenInfo, error) {
	app := application.Get()
	if app == nil || app.Screen == nil {
		return nil, fmt.Errorf("screen service not available")
	}
	screens := app.Screen.GetAll()
	for _, screen := range screens {
		if screen.IsPrimary {
			return &ScreenInfo{
				ID: screen.ID, Name: screen.Name,
				X: screen.Bounds.X, Y: screen.Bounds.Y,
				Width: screen.Bounds.Width, Height: screen.Bounds.Height,
				Primary: true,
			}, nil
		}
	}
	return nil, fmt.Errorf("no primary screen found")
}

// GetScreen returns information about a specific screen by ID.
func (s *Service) GetScreen(id string) (*ScreenInfo, error) {
	app := application.Get()
	if app == nil || app.Screen == nil {
		return nil, fmt.Errorf("screen service not available")
	}
	screens := app.Screen.GetAll()
	for _, screen := range screens {
		if screen.ID == id {
			return &ScreenInfo{
				ID: screen.ID, Name: screen.Name,
				X: screen.Bounds.X, Y: screen.Bounds.Y,
				Width: screen.Bounds.Width, Height: screen.Bounds.Height,
				Primary: screen.IsPrimary,
			}, nil
		}
	}
	return nil, fmt.Errorf("screen not found: %s", id)
}

// GetScreenAtPoint returns the screen containing a specific point.
func (s *Service) GetScreenAtPoint(x, y int) (*ScreenInfo, error) {
	app := application.Get()
	if app == nil || app.Screen == nil {
		return nil, fmt.Errorf("screen service not available")
	}
	screens := app.Screen.GetAll()
	for _, screen := range screens {
		bounds := screen.Bounds
		if x >= bounds.X && x < bounds.X+bounds.Width &&
			y >= bounds.Y && y < bounds.Y+bounds.Height {
			return &ScreenInfo{
				ID: screen.ID, Name: screen.Name,
				X: bounds.X, Y: bounds.Y,
				Width: bounds.Width, Height: bounds.Height,
				Primary: screen.IsPrimary,
			}, nil
		}
	}
	return nil, fmt.Errorf("no screen found at point (%d, %d)", x, y)
}

// GetScreenForWindow returns the screen containing a specific window.
func (s *Service) GetScreenForWindow(name string) (*ScreenInfo, error) {
	info, err := s.GetWindowInfo(name)
	if err != nil {
		return nil, err
	}
	centerX := info.X + info.Width/2
	centerY := info.Y + info.Height/2
	return s.GetScreenAtPoint(centerX, centerY)
}

// ShowEnvironmentDialog displays environment information.
func (s *Service) ShowEnvironmentDialog() {
	envInfo := s.app.Env().Info()
	details := "Environment Information:\n\n"
	details += fmt.Sprintf("Operating System: %s\n", envInfo.OS)
	details += fmt.Sprintf("Architecture: %s\n", envInfo.Arch)
	details += fmt.Sprintf("Debug Mode: %t\n\n", envInfo.Debug)
	details += fmt.Sprintf("Dark Mode: %t\n\n", s.app.Env().IsDarkMode())
	details += "Platform Information:"
	for key, value := range envInfo.PlatformInfo {
		details += fmt.Sprintf("\n%s: %v", key, value)
	}
	if envInfo.OSInfo != nil {
		details += fmt.Sprintf("\n\nOS Details:\nName: %s\nVersion: %s",
			envInfo.OSInfo.Name, envInfo.OSInfo.Version)
	}
	dialog := s.app.Dialog().Info()
	dialog.SetTitle("Environment Information")
	dialog.SetMessage(details)
	dialog.Show()
}

// GetEventManager returns the event manager for WebSocket event subscriptions.
func (s *Service) GetEventManager() *WSEventManager {
	return s.events
}

// --- Menu (handlers stay in display, structure delegated to menu.Manager) ---

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

	s.menus.SetApplicationMenu(items)
}

func ptr[T any](v T) *T { return &v }

// --- Menu handler methods ---

func (s *Service) handleNewWorkspace() {
	_ = s.OpenWindow(window.WithName("workspace-new"), window.WithTitle("New Workspace"),
		window.WithURL("/workspace/new"), window.WithSize(500, 400))
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
	_ = s.OpenWindow(window.WithName("editor"), window.WithTitle("New File - Editor"),
		window.WithURL("/#/developer/editor?new=true"), window.WithSize(1200, 800))
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
	_ = s.OpenWindow(window.WithName("editor"), window.WithTitle(result+" - Editor"),
		window.WithURL("/#/developer/editor?file="+result), window.WithSize(1200, 800))
}

func (s *Service) handleSaveFile()   { s.app.Event().Emit("ide:save") }
func (s *Service) handleOpenEditor() {
	_ = s.OpenWindow(window.WithName("editor"), window.WithTitle("Editor"),
		window.WithURL("/#/developer/editor"), window.WithSize(1200, 800))
}
func (s *Service) handleOpenTerminal() {
	_ = s.OpenWindow(window.WithName("terminal"), window.WithTitle("Terminal"),
		window.WithURL("/#/developer/terminal"), window.WithSize(800, 500))
}
func (s *Service) handleRun()   { s.app.Event().Emit("ide:run") }
func (s *Service) handleBuild() { s.app.Event().Emit("ide:build") }

// --- Tray (setup delegated to systray.Manager) ---

func (s *Service) setupTray() {
	_ = s.tray.Setup("Core", "Core")
	s.tray.RegisterCallback("open-desktop", func() {
		for _, name := range s.windows.List() {
			if pw, ok := s.windows.Get(name); ok {
				pw.Show()
			}
		}
	})
	s.tray.RegisterCallback("close-desktop", func() {
		for _, name := range s.windows.List() {
			if pw, ok := s.windows.Get(name); ok {
				pw.Hide()
			}
		}
	})
	s.tray.RegisterCallback("env-info", func() { s.ShowEnvironmentDialog() })
	s.tray.RegisterCallback("quit", func() { s.app.Quit() })
	_ = s.tray.SetMenu([]systray.TrayMenuItem{
		{Label: "Open Desktop", ActionID: "open-desktop"},
		{Label: "Close Desktop", ActionID: "close-desktop"},
		{Type: "separator"},
		{Label: "Environment Info", ActionID: "env-info"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	})
}
