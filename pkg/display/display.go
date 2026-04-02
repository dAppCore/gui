package display

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"encoding/json"
	"forge.lthn.ai/core/config"
	"forge.lthn.ai/core/go/pkg/core"

	"forge.lthn.ai/core/gui/pkg/browser"
	"forge.lthn.ai/core/gui/pkg/clipboard"
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
	"forge.lthn.ai/core/gui/pkg/webview"
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
	cfg        *config.Config // config instance for file persistence
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
	case webview.ActionConsoleMessage:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWebviewConsole, Window: m.Window,
				Data: map[string]any{"message": m.Message}})
		}
	case webview.ActionException:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWebviewException, Window: m.Window,
				Data: map[string]any{"exception": m.Exception}})
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

// wsRequire extracts a string field from WS data and returns an error if it is empty.
func wsRequire(data map[string]any, key string) (string, error) {
	v, _ := data[key].(string)
	if v == "" {
		return "", fmt.Errorf("ws: missing required field %q", key)
	}
	return v, nil
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
	case "webview:eval":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		script, _ := msg.Data["script"].(string)
		result, handled, err = s.Core().PERFORM(webview.TaskEvaluate{Window: w, Script: script})
	case "webview:click":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskClick{Window: w, Selector: sel})
	case "webview:type":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		text, _ := msg.Data["text"].(string)
		result, handled, err = s.Core().PERFORM(webview.TaskType{Window: w, Selector: sel, Text: text})
	case "webview:navigate":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		url, e := wsRequire(msg.Data, "url")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskNavigate{Window: w, URL: url})
	case "webview:screenshot":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskScreenshot{Window: w})
	case "webview:screenshot-element":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskScreenshotElement{Window: w, Selector: sel})
	case "webview:scroll":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		x, _ := msg.Data["x"].(float64)
		y, _ := msg.Data["y"].(float64)
		result, handled, err = s.Core().PERFORM(webview.TaskScroll{Window: w, X: int(x), Y: int(y)})
	case "webview:hover":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskHover{Window: w, Selector: sel})
	case "webview:select":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		val, _ := msg.Data["value"].(string)
		result, handled, err = s.Core().PERFORM(webview.TaskSelect{Window: w, Selector: sel, Value: val})
	case "webview:check":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		checked, _ := msg.Data["checked"].(bool)
		result, handled, err = s.Core().PERFORM(webview.TaskCheck{Window: w, Selector: sel, Checked: checked})
	case "webview:upload":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		pathsRaw, _ := msg.Data["paths"].([]any)
		var paths []string
		for _, p := range pathsRaw {
			if ps, ok := p.(string); ok {
				paths = append(paths, ps)
			}
		}
		result, handled, err = s.Core().PERFORM(webview.TaskUploadFile{Window: w, Selector: sel, Paths: paths})
	case "webview:viewport":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		width, _ := msg.Data["width"].(float64)
		height, _ := msg.Data["height"].(float64)
		result, handled, err = s.Core().PERFORM(webview.TaskSetViewport{Window: w, Width: int(width), Height: int(height)})
	case "webview:clear-console":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskClearConsole{Window: w})
	case "webview:highlight":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		colour, _ := msg.Data["colour"].(string)
		result, handled, err = s.Core().PERFORM(webview.TaskHighlight{Window: w, Selector: sel, Colour: colour})
	case "webview:computed-style":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryComputedStyle{Window: w, Selector: sel})
	case "webview:performance":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryPerformance{Window: w})
	case "webview:resources":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryResources{Window: w})
	case "webview:network":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		limit := 0
		if l, ok := msg.Data["limit"].(float64); ok {
			limit = int(l)
		}
		result, handled, err = s.Core().QUERY(webview.QueryNetwork{Window: w, Limit: limit})
	case "webview:network-inject":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskInjectNetworkLogging{Window: w})
	case "webview:network-clear":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskClearNetworkLog{Window: w})
	case "webview:print":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskPrint{Window: w})
	case "webview:pdf":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskExportPDF{Window: w})
	case "webview:console":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		level, _ := msg.Data["level"].(string)
		limit := 100
		if l, ok := msg.Data["limit"].(float64); ok {
			limit = int(l)
		}
		result, handled, err = s.Core().QUERY(webview.QueryConsole{Window: w, Level: level, Limit: limit})
	case "webview:query":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QuerySelector{Window: w, Selector: sel})
	case "webview:query-all":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QuerySelectorAll{Window: w, Selector: sel})
	case "webview:dom-tree":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, _ := msg.Data["selector"].(string) // selector optional for dom-tree (defaults to root)
		result, handled, err = s.Core().QUERY(webview.QueryDOMTree{Window: w, Selector: sel})
	case "webview:source":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryDOMTree{Window: w})
	case "webview:element-info":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QuerySelector{Window: w, Selector: sel})
	case "webview:url":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryURL{Window: w})
	case "webview:title":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(webview.QueryTitle{Window: w})
	case "webview:devtools-open":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskOpenDevTools{Window: w})
	case "webview:devtools-close":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(webview.TaskCloseDevTools{Window: w})
	case "layout:beside-editor":
		editor, _ := msg.Data["editor"].(string)
		windowName, _ := msg.Data["window"].(string)
		result, handled, err = s.Core().PERFORM(window.TaskBesideEditor{
			Editor: editor,
			Window: windowName,
		})
	case "layout:stack":
		offsetX, _ := msg.Data["offsetX"].(float64)
		offsetY, _ := msg.Data["offsetY"].(float64)
		var names []string
		if raw, ok := msg.Data["windows"].([]any); ok {
			for _, v := range raw {
				if name, ok := v.(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
		result, handled, err = s.Core().PERFORM(window.TaskStackWindows{
			Windows: names,
			OffsetX: int(offsetX),
			OffsetY: int(offsetY),
		})
	case "layout:workflow":
		workflowName, e := wsRequire(msg.Data, "workflow")
		if e != nil {
			return nil, false, e
		}
		workflow, ok := window.ParseWorkflowLayout(workflowName)
		if !ok {
			return nil, false, fmt.Errorf("ws: unknown workflow %q", workflowName)
		}
		var names []string
		if raw, ok := msg.Data["windows"].([]any); ok {
			for _, v := range raw {
				if name, ok := v.(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
		result, handled, err = s.Core().PERFORM(window.TaskApplyWorkflow{
			Workflow: workflow,
			Windows:  names,
		})
	case "window:arrange-pair":
		first, e := wsRequire(msg.Data, "first")
		if e != nil {
			return nil, false, e
		}
		second, e := wsRequire(msg.Data, "second")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(window.TaskArrangePair{
			First:  first,
			Second: second,
		})
	case "layout:suggest":
		windowCount := 0
		if count, ok := msg.Data["windowCount"].(float64); ok {
			windowCount = int(count)
		}
		screenWidth := 0
		if width, ok := msg.Data["screenWidth"].(float64); ok {
			screenWidth = int(width)
		}
		screenHeight := 0
		if height, ok := msg.Data["screenHeight"].(float64); ok {
			screenHeight = int(height)
		}
		if windowCount <= 0 {
			windowCount = len(s.ListWindowInfos())
		}
		if screenWidth <= 0 || screenHeight <= 0 {
			screenWidth, screenHeight = s.primaryScreenSize()
		}
		result, handled, err = s.Core().QUERY(window.QueryLayoutSuggestion{
			WindowCount:  windowCount,
			ScreenWidth:  screenWidth,
			ScreenHeight: screenHeight,
		})
	case "screen:find-space":
		width := 0
		if w, ok := msg.Data["width"].(float64); ok {
			width = int(w)
		}
		height := 0
		if h, ok := msg.Data["height"].(float64); ok {
			height = int(h)
		}
		result, handled, err = s.Core().QUERY(window.QueryFindSpace{
			Width:  width,
			Height: height,
		})
	case "screen:list":
		result, handled, err = s.Core().QUERY(screen.QueryAll{})
	case "screen:get":
		id, e := wsRequire(msg.Data, "id")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(screen.QueryByID{ID: id})
	case "screen:primary":
		result, handled, err = s.Core().QUERY(screen.QueryPrimary{})
	case "screen:at-point":
		x, _ := msg.Data["x"].(float64)
		y, _ := msg.Data["y"].(float64)
		result, handled, err = s.Core().QUERY(screen.QueryAtPoint{X: int(x), Y: int(y)})
	case "screen:work-areas":
		result, handled, err = s.Core().QUERY(screen.QueryWorkAreas{})
	case "screen:for-window":
		name, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		screenInfo, screenErr := s.GetScreenForWindow(name)
		if screenErr != nil {
			return nil, false, screenErr
		}
		result, handled, err = screenInfo, true, nil
	case "clipboard:read":
		result, handled, err = s.Core().QUERY(clipboard.QueryText{})
	case "clipboard:write":
		text, e := wsRequire(msg.Data, "text")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(clipboard.TaskSetText{Text: text})
	case "clipboard:has":
		textResult, textHandled, textErr := s.Core().QUERY(clipboard.QueryText{})
		if textErr != nil {
			return nil, false, textErr
		}
		hasContent := false
		if textHandled {
			if content, ok := textResult.(clipboard.ClipboardContent); ok {
				hasContent = content.HasContent
			}
		}
		if !hasContent {
			imageResult, imageHandled, imageErr := s.Core().QUERY(clipboard.QueryImage{})
			if imageErr != nil {
				return nil, false, imageErr
			}
			if imageHandled {
				if content, ok := imageResult.(clipboard.ClipboardImageContent); ok {
					hasContent = content.HasContent
				}
			}
		}
		result, handled, err = hasContent, true, nil
	case "clipboard:clear":
		result, handled, err = s.Core().PERFORM(clipboard.TaskClear{})
	case "clipboard:read-image":
		result, handled, err = s.Core().QUERY(clipboard.QueryImage{})
	case "clipboard:write-image":
		data, ok := msg.Data["data"].(string)
		if !ok || data == "" {
			return nil, false, fmt.Errorf("ws: missing required field %q", "data")
		}
		decoded, decodeErr := base64.StdEncoding.DecodeString(data)
		if decodeErr != nil {
			return nil, false, fmt.Errorf("ws: invalid base64 image data: %w", decodeErr)
		}
		result, handled, err = s.Core().PERFORM(clipboard.TaskSetImage{Data: decoded})
	case "notification:show":
		var opts notification.NotificationOptions
		encoded, _ := json.Marshal(msg.Data)
		_ = json.Unmarshal(encoded, &opts)
		result, handled, err = s.Core().PERFORM(notification.TaskSend{Opts: opts})
	case "notification:info":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(notification.TaskSend{Opts: notification.NotificationOptions{
			Title:    title,
			Message:  message,
			Severity: notification.SeverityInfo,
		}})
	case "notification:warning":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(notification.TaskSend{Opts: notification.NotificationOptions{
			Title:    title,
			Message:  message,
			Severity: notification.SeverityWarning,
		}})
	case "notification:error":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(notification.TaskSend{Opts: notification.NotificationOptions{
			Title:    title,
			Message:  message,
			Severity: notification.SeverityError,
		}})
	case "notification:with-actions":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		subtitle, _ := msg.Data["subtitle"].(string)
		actions := make([]notification.NotificationAction, 0)
		if raw, ok := msg.Data["actions"]; ok {
			encoded, _ := json.Marshal(raw)
			_ = json.Unmarshal(encoded, &actions)
		}
		result, handled, err = s.Core().PERFORM(notification.TaskSend{
			Opts: notification.NotificationOptions{
				Title:    title,
				Message:  message,
				Subtitle: subtitle,
				Actions:  actions,
			},
		})
	case "notification:clear":
		result, handled, err = s.Core().PERFORM(notification.TaskClear{})
	case "notification:permission-request":
		result, handled, err = s.Core().PERFORM(notification.TaskRequestPermission{})
	case "notification:permission-check":
		result, handled, err = s.Core().QUERY(notification.QueryPermission{})
	case "tray:show-message":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(systray.TaskShowMessage{
			Title:   title,
			Message: message,
		})
	case "tray:set-tooltip":
		tooltip, e := wsRequire(msg.Data, "tooltip")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetTooltip{Tooltip: tooltip})
	case "tray:set-label":
		label, e := wsRequire(msg.Data, "label")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetLabel{Label: label})
	case "tray:set-icon":
		data, e := wsRequire(msg.Data, "data")
		if e != nil {
			return nil, false, e
		}
		decoded, decodeErr := base64.StdEncoding.DecodeString(data)
		if decodeErr != nil {
			return nil, false, fmt.Errorf("ws: invalid base64 tray icon data: %w", decodeErr)
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetTrayIcon{Data: decoded})
	case "tray:set-menu":
		raw, ok := msg.Data["items"]
		if !ok {
			return nil, false, fmt.Errorf("ws: missing required field %q", "items")
		}
		encoded, _ := json.Marshal(raw)
		var items []systray.TrayMenuItem
		if err := json.Unmarshal(encoded, &items); err != nil {
			return nil, false, fmt.Errorf("ws: invalid tray menu items: %w", err)
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetTrayMenu{Items: items})
	case "tray:info":
		result, handled, err = s.GetTrayInfo(), true, nil
	case "theme:get":
		result, handled, err = s.GetTheme(), true, nil
	case "theme:system":
		result, handled, err = s.GetSystemTheme(), true, nil
	case "theme:set":
		isDark, _ := msg.Data["isDark"].(bool)
		result, handled, err = nil, true, s.SetTheme(isDark)
	case "dialog:open-file":
		var opts dialog.OpenFileOptions
		encoded, _ := json.Marshal(msg.Data)
		if err := json.Unmarshal(encoded, &opts); err != nil {
			return nil, false, fmt.Errorf("ws: invalid open file options: %w", err)
		}
		paths, openErr := s.OpenFileDialog(opts)
		if openErr != nil {
			return nil, false, openErr
		}
		result, handled, err = paths, true, nil
	case "dialog:save-file":
		var opts dialog.SaveFileOptions
		encoded, _ := json.Marshal(msg.Data)
		if err := json.Unmarshal(encoded, &opts); err != nil {
			return nil, false, fmt.Errorf("ws: invalid save file options: %w", err)
		}
		path, saveErr := s.SaveFileDialog(opts)
		if saveErr != nil {
			return nil, false, saveErr
		}
		result, handled, err = path, true, nil
	case "dialog:open-directory":
		var opts dialog.OpenDirectoryOptions
		encoded, _ := json.Marshal(msg.Data)
		if err := json.Unmarshal(encoded, &opts); err != nil {
			return nil, false, fmt.Errorf("ws: invalid open directory options: %w", err)
		}
		path, dirErr := s.OpenDirectoryDialog(opts)
		if dirErr != nil {
			return nil, false, dirErr
		}
		result, handled, err = path, true, nil
	case "dialog:confirm":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		confirmed, confirmErr := s.ConfirmDialog(title, message)
		if confirmErr != nil {
			return nil, false, confirmErr
		}
		result, handled, err = confirmed, true, nil
	case "dialog:prompt":
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		message, e := wsRequire(msg.Data, "message")
		if e != nil {
			return nil, false, e
		}
		button, accepted, promptErr := s.PromptDialog(title, message)
		if promptErr != nil {
			return nil, false, promptErr
		}
		_ = accepted
		result, handled, err = button, true, nil
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
		// Hide all tracked windows using the existing visibility task.
		infos := s.ListWindowInfos()
		for _, info := range infos {
			_, _, _ = s.Core().PERFORM(window.TaskSetVisibility{
				Name:    info.Name,
				Visible: false,
			})
		}
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

// FocusSet is a compatibility alias for FocusWindow.
func (s *Service) FocusSet(name string) error {
	return s.FocusWindow(name)
}

// CloseWindow closes a window via IPC.
func (s *Service) CloseWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskCloseWindow{Name: name})
	return err
}

// RestoreWindow restores a maximized/minimized window via IPC.
func (s *Service) RestoreWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskRestore{Name: name})
	return err
}

// SetWindowVisibility shows or hides a window via IPC.
func (s *Service) SetWindowVisibility(name string, visible bool) error {
	_, _, err := s.Core().PERFORM(window.TaskSetVisibility{Name: name, Visible: visible})
	return err
}

// SetWindowAlwaysOnTop sets whether a window stays on top via IPC.
func (s *Service) SetWindowAlwaysOnTop(name string, alwaysOnTop bool) error {
	_, _, err := s.Core().PERFORM(window.TaskSetAlwaysOnTop{Name: name, AlwaysOnTop: alwaysOnTop})
	return err
}

// SetWindowTitle changes a window's title via IPC.
func (s *Service) SetWindowTitle(name string, title string) error {
	_, _, err := s.Core().PERFORM(window.TaskSetTitle{Name: name, Title: title})
	return err
}

// SetWindowFullscreen sets a window to fullscreen mode via IPC.
func (s *Service) SetWindowFullscreen(name string, fullscreen bool) error {
	_, _, err := s.Core().PERFORM(window.TaskFullscreen{Name: name, Fullscreen: fullscreen})
	return err
}

// SetWindowBackgroundColour sets the background colour of a window via IPC.
func (s *Service) SetWindowBackgroundColour(name string, r, g, b, a uint8) error {
	_, _, err := s.Core().PERFORM(window.TaskSetBackgroundColour{
		Name:  name,
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a,
	})
	return err
}

// SetWindowOpacity updates a window's opacity via IPC.
func (s *Service) SetWindowOpacity(name string, opacity float32) error {
	_, _, err := s.Core().PERFORM(window.TaskSetOpacity{
		Name:    name,
		Opacity: opacity,
	})
	return err
}

// ClearWebviewConsole clears the captured console buffer for a window.
func (s *Service) ClearWebviewConsole(name string) error {
	_, _, err := s.Core().PERFORM(webview.TaskClearConsole{Window: name})
	return err
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
	return info.Title, nil
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
	Name             string   `json:"name"`
	Title            string   `json:"title,omitempty"`
	URL              string   `json:"url,omitempty"`
	X                int      `json:"x,omitempty"`
	Y                int      `json:"y,omitempty"`
	Width            int      `json:"width,omitempty"`
	Height           int      `json:"height,omitempty"`
	AlwaysOnTop      bool     `json:"alwaysOnTop,omitempty"`
	BackgroundColour [4]uint8 `json:"backgroundColour,omitempty"`
}

// CreateWindow creates a new window with the specified options.
func (s *Service) CreateWindow(opts CreateWindowOptions) (*window.WindowInfo, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf("window name is required")
	}
	result, _, err := s.Core().PERFORM(window.TaskOpenWindow{
		Window: &window.Window{
			Name:             opts.Name,
			Title:            opts.Title,
			URL:              opts.URL,
			X:                opts.X,
			Y:                opts.Y,
			Width:            opts.Width,
			Height:           opts.Height,
			AlwaysOnTop:      opts.AlwaysOnTop,
			BackgroundColour: opts.BackgroundColour,
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
	screenWidth, screenHeight := s.primaryScreenSize()
	return ws.Manager().TileWindows(mode, windowNames, screenWidth, screenHeight)
}

// SnapWindow snaps a window to a screen edge or corner.
func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	screenWidth, screenHeight := s.primaryScreenSize()
	return ws.Manager().SnapWindow(name, position, screenWidth, screenHeight)
}

func (s *Service) primaryScreenSize() (int, int) {
	const fallbackWidth = 1920
	const fallbackHeight = 1080

	result, handled, err := s.Core().QUERY(screen.QueryPrimary{})
	if err != nil || !handled {
		return fallbackWidth, fallbackHeight
	}

	primary, ok := result.(*screen.Screen)
	if !ok || primary == nil {
		return fallbackWidth, fallbackHeight
	}

	width := primary.WorkArea.Width
	height := primary.WorkArea.Height
	if width <= 0 || height <= 0 {
		width = primary.Bounds.Width
		height = primary.Bounds.Height
	}
	if width <= 0 || height <= 0 {
		return fallbackWidth, fallbackHeight
	}

	return width, height
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
	screenWidth, screenHeight := s.primaryScreenSize()
	return ws.Manager().ApplyWorkflow(workflow, ws.Manager().List(), screenWidth, screenHeight)
}

// ArrangeWindowPair places two windows side by side using the window manager's balanced split.
func (s *Service) ArrangeWindowPair(first, second string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	screenWidth, screenHeight := s.primaryScreenSize()
	return ws.Manager().ArrangePair(first, second, screenWidth, screenHeight)
}

// FindSpace returns a free placement suggestion for a new window.
func (s *Service) FindSpace(width, height int) (window.SpaceInfo, error) {
	ws := s.windowService()
	if ws == nil {
		return window.SpaceInfo{}, fmt.Errorf("window service not available")
	}
	screenWidth, screenHeight := s.primaryScreenSize()
	if width <= 0 {
		width = screenWidth / 2
	}
	if height <= 0 {
		height = screenHeight / 2
	}
	return ws.Manager().FindSpace(screenWidth, screenHeight, width, height), nil
}

// --- Screen management ---

// GetScreens returns all known screens.
func (s *Service) GetScreens() []screen.Screen {
	result, handled, _ := s.Core().QUERY(screen.QueryAll{})
	if !handled {
		return nil
	}
	screens, _ := result.([]screen.Screen)
	return screens
}

// GetScreen returns a screen by ID.
func (s *Service) GetScreen(id string) (*screen.Screen, error) {
	result, handled, err := s.Core().QUERY(screen.QueryByID{ID: id})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("screen service not available")
	}
	scr, _ := result.(*screen.Screen)
	return scr, nil
}

// GetPrimaryScreen returns the primary screen.
func (s *Service) GetPrimaryScreen() (*screen.Screen, error) {
	result, handled, err := s.Core().QUERY(screen.QueryPrimary{})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("screen service not available")
	}
	scr, _ := result.(*screen.Screen)
	return scr, nil
}

// GetScreenAtPoint returns the screen containing the specified point.
func (s *Service) GetScreenAtPoint(x, y int) (*screen.Screen, error) {
	result, handled, err := s.Core().QUERY(screen.QueryAtPoint{X: x, Y: y})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("screen service not available")
	}
	scr, _ := result.(*screen.Screen)
	return scr, nil
}

// GetScreenForWindow returns the screen containing the named window.
func (s *Service) GetScreenForWindow(name string) (*screen.Screen, error) {
	info, err := s.GetWindowInfo(name)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	x := info.X
	y := info.Y
	if info.Width > 0 && info.Height > 0 {
		x += info.Width / 2
		y += info.Height / 2
	}
	return s.GetScreenAtPoint(x, y)
}

// GetWorkAreas returns the usable area of every screen.
func (s *Service) GetWorkAreas() []screen.Rect {
	result, handled, _ := s.Core().QUERY(screen.QueryWorkAreas{})
	if !handled {
		return nil
	}
	areas, _ := result.([]screen.Rect)
	return areas
}

// --- Clipboard ---

// ReadClipboard returns the current clipboard text content.
func (s *Service) ReadClipboard() (string, error) {
	result, handled, err := s.Core().QUERY(clipboard.QueryText{})
	if err != nil {
		return "", err
	}
	if !handled {
		return "", core.E("display.ReadClipboard", "clipboard service not available", nil)
	}
	content, _ := result.(clipboard.ClipboardContent)
	return content.Text, nil
}

// WriteClipboard writes text to the clipboard.
func (s *Service) WriteClipboard(text string) error {
	result, handled, err := s.Core().PERFORM(clipboard.TaskSetText{Text: text})
	if err != nil {
		return err
	}
	if !handled {
		return core.E("display.WriteClipboard", "clipboard service not available", nil)
	}
	if ok, _ := result.(bool); !ok {
		return core.E("display.WriteClipboard", "clipboard write failed", nil)
	}
	return nil
}

// HasClipboard reports whether the clipboard has text or image content.
func (s *Service) HasClipboard() bool {
	textResult, textHandled, _ := s.Core().QUERY(clipboard.QueryText{})
	if textHandled {
		if content, ok := textResult.(clipboard.ClipboardContent); ok && content.HasContent {
			return true
		}
	}
	imageResult, imageHandled, _ := s.Core().QUERY(clipboard.QueryImage{})
	if imageHandled {
		if content, ok := imageResult.(clipboard.ClipboardImageContent); ok && content.HasContent {
			return true
		}
	}
	return false
}

// ClearClipboard clears clipboard text and any image data when supported.
func (s *Service) ClearClipboard() error {
	result, handled, err := s.Core().PERFORM(clipboard.TaskClear{})
	if err != nil {
		return err
	}
	if !handled {
		return core.E("display.ClearClipboard", "clipboard service not available", nil)
	}
	if ok, _ := result.(bool); !ok {
		return core.E("display.ClearClipboard", "clipboard clear failed", nil)
	}
	return nil
}

// ReadClipboardImage returns the clipboard image content.
func (s *Service) ReadClipboardImage() (clipboard.ClipboardImageContent, error) {
	result, handled, err := s.Core().QUERY(clipboard.QueryImage{})
	if err != nil {
		return clipboard.ClipboardImageContent{}, err
	}
	if !handled {
		return clipboard.ClipboardImageContent{}, core.E("display.ReadClipboardImage", "clipboard service not available", nil)
	}
	content, _ := result.(clipboard.ClipboardImageContent)
	return content, nil
}

// WriteClipboardImage writes raw image data to the clipboard.
func (s *Service) WriteClipboardImage(data []byte) error {
	result, handled, err := s.Core().PERFORM(clipboard.TaskSetImage{Data: data})
	if err != nil {
		return err
	}
	if !handled {
		return core.E("display.WriteClipboardImage", "clipboard service not available", nil)
	}
	if ok, _ := result.(bool); !ok {
		return core.E("display.WriteClipboardImage", "clipboard image write failed", nil)
	}
	return nil
}

// --- Notifications ---

// ShowNotification sends a native notification.
func (s *Service) ShowNotification(opts notification.NotificationOptions) error {
	_, handled, err := s.Core().PERFORM(notification.TaskSend{Opts: opts})
	if err != nil {
		return err
	}
	if !handled {
		return core.E("display.ShowNotification", "notification service not available", nil)
	}
	return nil
}

// ShowInfoNotification sends an informational notification.
func (s *Service) ShowInfoNotification(title, message string) error {
	return s.ShowNotification(notification.NotificationOptions{
		Title:    title,
		Message:  message,
		Severity: notification.SeverityInfo,
	})
}

// ShowWarningNotification sends a warning notification.
func (s *Service) ShowWarningNotification(title, message string) error {
	return s.ShowNotification(notification.NotificationOptions{
		Title:    title,
		Message:  message,
		Severity: notification.SeverityWarning,
	})
}

// ShowErrorNotification sends an error notification.
func (s *Service) ShowErrorNotification(title, message string) error {
	return s.ShowNotification(notification.NotificationOptions{
		Title:    title,
		Message:  message,
		Severity: notification.SeverityError,
	})
}

// RequestNotificationPermission requests notification permission.
func (s *Service) RequestNotificationPermission() (bool, error) {
	result, handled, err := s.Core().PERFORM(notification.TaskRequestPermission{})
	if err != nil {
		return false, err
	}
	if !handled {
		return false, fmt.Errorf("notification service not available")
	}
	granted, _ := result.(bool)
	return granted, nil
}

// CheckNotificationPermission checks notification permission.
func (s *Service) CheckNotificationPermission() (bool, error) {
	result, handled, err := s.Core().QUERY(notification.QueryPermission{})
	if err != nil {
		return false, err
	}
	if !handled {
		return false, fmt.Errorf("notification service not available")
	}
	status, _ := result.(notification.PermissionStatus)
	return status.Granted, nil
}

// ClearNotifications clears notifications when supported.
func (s *Service) ClearNotifications() error {
	_, handled, err := s.Core().PERFORM(notification.TaskClear{})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("notification service not available")
	}
	return nil
}

// --- Dialogs ---

// OpenFileDialog opens a file picker and returns all selected paths.
func (s *Service) OpenFileDialog(opts dialog.OpenFileOptions) ([]string, error) {
	result, handled, err := s.Core().PERFORM(dialog.TaskOpenFile{Opts: opts})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("dialog service not available")
	}
	paths, _ := result.([]string)
	return paths, nil
}

// OpenSingleFileDialog opens a file picker and returns the first selected path.
func (s *Service) OpenSingleFileDialog(opts dialog.OpenFileOptions) (string, error) {
	paths, err := s.OpenFileDialog(opts)
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", nil
	}
	return paths[0], nil
}

// SaveFileDialog opens a save dialog and returns the selected path.
func (s *Service) SaveFileDialog(opts dialog.SaveFileOptions) (string, error) {
	result, handled, err := s.Core().PERFORM(dialog.TaskSaveFile{Opts: opts})
	if err != nil {
		return "", err
	}
	if !handled {
		return "", fmt.Errorf("dialog service not available")
	}
	path, _ := result.(string)
	return path, nil
}

// OpenDirectoryDialog opens a directory picker and returns the selected path.
func (s *Service) OpenDirectoryDialog(opts dialog.OpenDirectoryOptions) (string, error) {
	result, handled, err := s.Core().PERFORM(dialog.TaskOpenDirectory{Opts: opts})
	if err != nil {
		return "", err
	}
	if !handled {
		return "", fmt.Errorf("dialog service not available")
	}
	path, _ := result.(string)
	return path, nil
}

// ConfirmDialog shows a confirmation prompt.
func (s *Service) ConfirmDialog(title, message string) (bool, error) {
	result, handled, err := s.Core().PERFORM(dialog.TaskMessageDialog{
		Opts: dialog.MessageDialogOptions{
			Type:    dialog.DialogQuestion,
			Title:   title,
			Message: message,
			Buttons: []string{"Yes", "No"},
		},
	})
	if err != nil {
		return false, err
	}
	if !handled {
		return false, fmt.Errorf("dialog service not available")
	}
	button, _ := result.(string)
	return button == "Yes" || button == "OK", nil
}

// PromptDialog shows a prompt-style dialog and returns the selected button.
func (s *Service) PromptDialog(title, message string) (string, bool, error) {
	result, handled, err := s.Core().PERFORM(dialog.TaskMessageDialog{
		Opts: dialog.MessageDialogOptions{
			Type:    dialog.DialogInfo,
			Title:   title,
			Message: message,
			Buttons: []string{"OK", "Cancel"},
		},
	})
	if err != nil {
		return "", false, err
	}
	if !handled {
		return "", false, fmt.Errorf("dialog service not available")
	}
	button, _ := result.(string)
	return button, button == "OK", nil
}

// DialogMessage shows an informational, warning, or error message via the notification pipeline.
func (s *Service) DialogMessage(kind, title, message string) error {
	var severity notification.NotificationSeverity
	switch kind {
	case "warning":
		severity = notification.SeverityWarning
	case "error":
		severity = notification.SeverityError
	default:
		severity = notification.SeverityInfo
	}
	_, _, err := s.Core().PERFORM(notification.TaskSend{
		Opts: notification.NotificationOptions{
			Title:    title,
			Message:  message,
			Severity: severity,
		},
	})
	return err
}

// --- Theme ---

// GetTheme returns the current theme state.
func (s *Service) GetTheme() *Theme {
	result, handled, err := s.Core().QUERY(environment.QueryTheme{})
	if err != nil || !handled {
		return nil
	}
	theme, ok := result.(environment.ThemeInfo)
	if !ok {
		return nil
	}
	return &Theme{IsDark: theme.IsDark}
}

// GetSystemTheme returns the current system theme preference.
func (s *Service) GetSystemTheme() string {
	result, handled, err := s.Core().QUERY(environment.QueryTheme{})
	if err != nil || !handled {
		return ""
	}
	theme, ok := result.(environment.ThemeInfo)
	if !ok {
		return ""
	}
	if theme.IsDark {
		return "dark"
	}
	return "light"
}

// SetTheme overrides the application theme.
func (s *Service) SetTheme(isDark bool) error {
	_, handled, err := s.Core().PERFORM(environment.TaskSetTheme{IsDark: isDark})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("environment service not available")
	}
	return nil
}

// --- Tray ---

// SetTrayIcon sets the tray icon image.
func (s *Service) SetTrayIcon(data []byte) error {
	_, handled, err := s.Core().PERFORM(systray.TaskSetTrayIcon{Data: data})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("systray service not available")
	}
	return nil
}

// SetTrayTooltip updates the tray tooltip.
func (s *Service) SetTrayTooltip(tooltip string) error {
	_, handled, err := s.Core().PERFORM(systray.TaskSetTooltip{Tooltip: tooltip})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("systray service not available")
	}
	return nil
}

// SetTrayLabel updates the tray label.
func (s *Service) SetTrayLabel(label string) error {
	_, handled, err := s.Core().PERFORM(systray.TaskSetLabel{Label: label})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("systray service not available")
	}
	return nil
}

// SetTrayMenu replaces the tray menu items.
func (s *Service) SetTrayMenu(items []systray.TrayMenuItem) error {
	_, handled, err := s.Core().PERFORM(systray.TaskSetTrayMenu{Items: items})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("systray service not available")
	}
	return nil
}

// GetTrayInfo returns current tray state information.
func (s *Service) GetTrayInfo() map[string]any {
	svc, err := core.ServiceFor[*systray.Service](s.Core(), "systray")
	if err != nil || svc == nil || svc.Manager() == nil {
		return nil
	}
	return svc.Manager().GetInfo()
}

// ShowTrayMessage shows a tray message or notification.
func (s *Service) ShowTrayMessage(title, message string) error {
	_, handled, err := s.Core().PERFORM(systray.TaskShowMessage{Title: title, Message: message})
	if err != nil {
		return err
	}
	if !handled {
		return fmt.Errorf("systray service not available")
	}
	return nil
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
	result, handled, err := s.Core().PERFORM(dialog.TaskOpenFile{
		Opts: dialog.OpenFileOptions{
			Title:         "Open File",
			AllowMultiple: false,
		},
	})
	if err != nil || !handled {
		return
	}
	paths, ok := result.([]string)
	if !ok || len(paths) == 0 {
		return
	}
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle(paths[0] + " - Editor"),
			window.WithURL("/#/developer/editor?file=" + paths[0]),
			window.WithSize(1200, 800),
		},
	})
}

func (s *Service) handleSaveFile() { _ = s.Core().ACTION(ActionIDECommand{Command: "save"}) }
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
func (s *Service) handleRun()   { _ = s.Core().ACTION(ActionIDECommand{Command: "run"}) }
func (s *Service) handleBuild() { _ = s.Core().ACTION(ActionIDECommand{Command: "build"}) }

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
