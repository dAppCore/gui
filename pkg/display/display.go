// pkg/display/display.go
package display

import (
	"context"
	"runtime"

	corego "dappco.re/go/core"
	coreutil "dappco.re/go/core"
	"forge.lthn.ai/core/config"
	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"

	"dappco.re/go/core/gui/pkg/browser"
	"dappco.re/go/core/gui/pkg/contextmenu"
	"dappco.re/go/core/gui/pkg/dialog"
	"dappco.re/go/core/gui/pkg/dock"
	"dappco.re/go/core/gui/pkg/environment"
	"dappco.re/go/core/gui/pkg/events"
	"dappco.re/go/core/gui/pkg/keybinding"
	"dappco.re/go/core/gui/pkg/lifecycle"
	"dappco.re/go/core/gui/pkg/menu"
	"dappco.re/go/core/gui/pkg/notification"
	"dappco.re/go/core/gui/pkg/screen"
	"dappco.re/go/core/gui/pkg/systray"
	"dappco.re/go/core/gui/pkg/webview"
	"dappco.re/go/core/gui/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Options struct{}

// WindowInfo is an alias for window.WindowInfo (backward compatibility).
type WindowInfo = window.WindowInfo

// Service orchestrates window, systray, and menu sub-services via IPC.
// Bridges IPC actions to WebSocket events for TypeScript apps.
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	configData map[string]map[string]any
	configFile *config.Config // config instance for file persistence
	events     *WSEventManager
	chat       *ChatStore
	schemes    map[string]SchemeHandler
}

// NewService returns a display Service with empty config sections.
// svc, _ := display.NewService(); _, _ = svc.CreateWindow(display.CreateWindowOptions{Name: "settings", URL: "/settings", Width: 800, Height: 600})
func NewService() (*Service, error) {
	return &Service{
		configData: map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
		},
		chat:    NewChatStore(),
		schemes: make(map[string]SchemeHandler),
	}, nil
}

// Deprecated: use NewService().
func New() (*Service, error) {
	return NewService()
}

// Register binds the display service to a Core instance.
// core.WithService(display.Register(app))      // production (Wails app)
// core.WithService(display.Register(nil))      // tests (no Wails runtime)
func Register(wailsApp *application.App) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		s, err := NewService()
		if err != nil {
			return nil, err
		}
		s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
		s.wailsApp = wailsApp
		return s, nil
	}
}

// OnStartup loads config and registers handlers before sub-services start.
// Config handlers are registered first — sub-services query them during their own OnStartup.
func (s *Service) OnStartup(ctx context.Context) error {
	s.loadConfig()

	// Register config query/task handlers — available NOW for sub-services
	s.Core().RegisterQuery(s.handleConfigQuery)
	s.Core().RegisterTask(s.handleConfigTask)
	s.Core().RegisterQuery(s.handleChatQuery)
	s.Core().RegisterTask(s.handleChatTask)

	s.registerBuiltinSchemes()

	// Initialise Wails wrappers if app is available (nil in tests)
	if s.wailsApp != nil {
		s.app = newWailsApp(s.wailsApp)
		s.events = NewWSEventManager()
	}

	return nil
}

// HandleIPCEvents bridges IPC actions from sub-services to WebSocket events for TS apps.
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
				Data: map[string]any{"w": m.Width, "h": m.Height}})
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
	case events.ActionEventFired:
		if s.events != nil {
			s.events.Emit(Event{Type: EventCustomEvent,
				Data: map[string]any{"name": m.Name, "data": m.Data}})
		}
	case dock.ActionProgressChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventDockProgress,
				Data: map[string]any{"value": m.Value}})
		}
	case dock.ActionBounceStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventDockBounce,
				Data: map[string]any{"bounceId": m.BounceID, "type": m.Type}})
		}
	case notification.ActionNotificationActionTriggered:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationAction,
				Data: map[string]any{"notificationId": m.NotificationID, "actionId": m.ActionID}})
		}
	case notification.ActionNotificationDismissed:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationDismiss,
				Data: map[string]any{"notificationId": m.NotificationID}})
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
		return "", coreerr.E("display.wsRequire", "missing required field \""+key+"\"", nil)
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
	case "window:list":
		result, handled, err = s.Core().QUERY(window.QueryWindowList{})
	case "window:get":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = s.Core().QUERY(window.QueryWindowByName{Name: name})
	case "window:focused":
		result, handled, err = s.GetFocusedWindow(), true, nil
	case "window:create":
		var opts CreateWindowOptions
		encodedR := corego.JSONMarshal(msg.Data)
		if encodedR.OK {
			_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &opts)
		}
		return nil, false, corego.Wrap(err, "display.ws", "ws: invalid window create options")
		info, createErr := s.CreateWindow(opts)
		if createErr != nil {
			return nil, false, createErr
		}
		result, handled, err = info, true, nil
	case "window:close":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.CloseWindow(name)
	case "window:position":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		x, _ := msg.Data["x"].(float64)
		y, _ := msg.Data["y"].(float64)
		result, handled, err = nil, true, s.SetWindowPosition(name, int(x), int(y))
	case "window:size":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		width, _ := msg.Data["width"].(float64)
		height, _ := msg.Data["height"].(float64)
		result, handled, err = nil, true, s.SetWindowSize(name, int(width), int(height))
	case "window:bounds":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		x, _ := msg.Data["x"].(float64)
		y, _ := msg.Data["y"].(float64)
		width, _ := msg.Data["width"].(float64)
		height, _ := msg.Data["height"].(float64)
		result, handled, err = nil, true, s.SetWindowBounds(name, int(x), int(y), int(width), int(height))
	case "window:maximize":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.MaximizeWindow(name)
	case "window:minimize":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.MinimizeWindow(name)
	case "window:restore":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.RestoreWindow(name)
	case "window:focus":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.FocusWindow(name)
	case "focus:set":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.FocusSet(name)
	case "window:visibility":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		visible, _ := msg.Data["visible"].(bool)
		result, handled, err = nil, true, s.SetWindowVisibility(name, visible)
	case "window:title-set":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		title, e := wsRequire(msg.Data, "title")
		if e != nil {
			return nil, false, e
		}
		result, handled, err = nil, true, s.SetWindowTitle(name, title)
	case "window:title-get":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		title, titleErr := s.GetWindowTitle(name)
		if titleErr != nil {
			return nil, false, titleErr
		}
		result, handled, err = title, true, nil
	case "window:always-on-top":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		alwaysOnTop, _ := msg.Data["alwaysOnTop"].(bool)
		result, handled, err = nil, true, s.SetWindowAlwaysOnTop(name, alwaysOnTop)
	case "window:background-colour":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		red, _ := msg.Data["red"].(float64)
		green, _ := msg.Data["green"].(float64)
		blue, _ := msg.Data["blue"].(float64)
		alpha, _ := msg.Data["alpha"].(float64)
		result, handled, err = nil, true, s.SetWindowBackgroundColour(name, uint8(red), uint8(green), uint8(blue), uint8(alpha))
	case "window:opacity":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		opacity, _ := msg.Data["opacity"].(float64)
		result, handled, err = nil, true, s.SetWindowOpacity(name, float32(opacity))
	case "window:fullscreen":
		name, e := wsRequire(msg.Data, "name")
		if e != nil {
			return nil, false, e
		}
		fullscreen, _ := msg.Data["fullscreen"].(bool)
		result, handled, err = nil, true, s.SetWindowFullscreen(name, fullscreen)
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
		menuJSON := coreutil.JSONMarshalString(msg.Data["menu"])
		var menuDef contextmenu.ContextMenuDef
		_ = coreutil.JSONUnmarshalString(menuJSON, &menuDef)
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
	case "webview:errors":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return nil, false, e
		}
		limit := 0
		if l, ok := msg.Data["limit"].(float64); ok {
			limit = int(l)
		}
		result, handled, err = s.Core().QUERY(webview.QueryExceptions{Window: w, Limit: limit})
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
			return nil, false, corego.NewError(corego.Sprintf("ws: unknown workflow %q", workflowName))
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
		screenWidth, screenHeight := s.primaryScreenSize()
		result, handled, err = s.Core().QUERY(window.QueryFindSpace{
			Width:        width,
			Height:       height,
			ScreenWidth:  screenWidth,
			ScreenHeight: screenHeight,
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
			return nil, false, corego.NewError(corego.Sprintf("ws: missing required field %q", "data"))
		}
		decoded, decodeErr := base64.StdEncoding.DecodeString(data)
		if decodeErr != nil {
			return nil, false, corego.Wrap(decodeErr, "display.ws", "ws: invalid base64 image data")
		}
		result, handled, err = s.Core().PERFORM(clipboard.TaskSetImage{Data: decoded})
	case "notification:show":
		var opts notification.NotificationOptions
		encodedR := corego.JSONMarshal(msg.Data)
		if encodedR.OK {
			_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &opts)
		}
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
			encodedR := corego.JSONMarshal(raw)
			if encodedR.OK {
				_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &actions)
			}
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
			return nil, false, corego.Wrap(decodeErr, "display.ws", "ws: invalid base64 tray icon data")
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetTrayIcon{Data: decoded})
	case "tray:set-menu":
		raw, ok := msg.Data["items"]
		if !ok {
			return nil, false, corego.NewError(corego.Sprintf("ws: missing required field %q", "items"))
		}
		encodedItemsR := corego.JSONMarshal(raw)
		var items []systray.TrayMenuItem
		if encodedItemsR.OK {
			if r := corego.JSONUnmarshal(encodedItemsR.Value.([]byte), &items); !r.OK {
				return nil, false, corego.E("display.ws", "invalid tray menu items", nil)
			}
		}
		result, handled, err = s.Core().PERFORM(systray.TaskSetTrayMenu{Items: items})
	case "tray:info":
		result, handled, err = s.GetTrayInfo(), true, nil
	case "event:info":
		result, handled, err = s.GetEventInfo(), true, nil
	case "theme:get":
		result, handled, err = s.GetTheme(), true, nil
	case "theme:system":
		result, handled, err = s.GetSystemTheme(), true, nil
	case "theme:set":
		if theme, ok := msg.Data["theme"].(string); ok && theme != "" {
			result, handled, err = nil, true, s.SetThemeMode(theme)
			break
		}
		isDark, _ := msg.Data["isDark"].(bool)
		result, handled, err = nil, true, s.SetTheme(isDark)
	case "dialog:open-file":
		var opts dialog.OpenFileOptions
		encodedR := corego.JSONMarshal(msg.Data)
		if encodedR.OK {
			_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &opts)
		}
		return nil, false, corego.Wrap(err, "display.ws", "ws: invalid open file options")
		paths, openErr := s.OpenFileDialog(opts)
		if openErr != nil {
			return nil, false, openErr
		}
		result, handled, err = paths, true, nil
	case "dialog:save-file":
		var opts dialog.SaveFileOptions
		encodedR := corego.JSONMarshal(msg.Data)
		if encodedR.OK {
			_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &opts)
		}
		return nil, false, corego.Wrap(err, "display.ws", "ws: invalid save file options")
		path, saveErr := s.SaveFileDialog(opts)
		if saveErr != nil {
			return nil, false, saveErr
		}
		result, handled, err = path, true, nil
	case "dialog:open-directory":
		var opts dialog.OpenDirectoryOptions
		encodedR := corego.JSONMarshal(msg.Data)
		if encodedR.OK {
			_ = corego.JSONUnmarshal(encodedR.Value.([]byte), &opts)
		}
		return nil, false, corego.Wrap(err, "display.ws", "ws: invalid open directory options")
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
			details := "OS: " + info.OS + "\nArch: " + info.Arch + "\nPlatform: " +
				info.Platform.Name + " " + info.Platform.Version
			_, _, _ = s.Core().PERFORM(dialog.TaskMessageDialog{
				Options: dialog.MessageDialogOptions{
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
	home := coreutil.Env("DIR_HOME")
	if home == "" {
		return coreutil.JoinPath(".core", "gui", "config.yaml")
	}
	return coreutil.JoinPath(home, ".core", "gui", "config.yaml")
}

func (s *Service) loadConfig() {
	if s.configFile != nil {
		return // Already loaded (e.g., via loadConfigFrom in tests)
	}
	s.loadConfigFrom(guiConfigPath())
}

func (s *Service) loadConfigFrom(path string) {
	configFile, err := config.New(config.WithPath(path))
	if err != nil {
		// Non-critical — continue with empty configData
		return
	}
	s.configFile = configFile

	for _, section := range []string{"window", "systray", "menu"} {
		var data map[string]any
		if err := configFile.Get(section, &data); err == nil && data != nil {
			s.configData[section] = data
		}
	}

	s.chat.Load(configFile)
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
		s.configData["window"] = t.Config
		s.persistSection("window", t.Config)
		return nil, true, nil
	case systray.TaskSaveConfig:
		s.configData["systray"] = t.Config
		s.persistSection("systray", t.Config)
		return nil, true, nil
	case menu.TaskSaveConfig:
		s.configData["menu"] = t.Config
		s.persistSection("menu", t.Config)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) persistSection(key string, value map[string]any) {
	if s.configFile == nil {
		return
	}
	_ = s.configFile.Set(key, value)
	_ = s.configFile.Commit()
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

// Deprecated: use CreateWindow(display.CreateWindowOptions{Name: "settings", URL: "/settings", Width: 800, Height: 600}).
func (s *Service) OpenWindow(options ...window.WindowOption) error {
	spec, err := window.ApplyOptions(options...)
	if err != nil {
		return err
	}
	_, _, err = s.Core().PERFORM(window.TaskOpenWindow{Window: spec})
	return err
}

// GetWindowInfo returns information about a window via IPC.
// Use: info, err := svc.GetWindowInfo("editor")
func (s *Service) GetWindowInfo(name string) (*window.WindowInfo, error) {
	result, handled, err := s.Core().QUERY(window.QueryWindowByName{Name: name})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, coreerr.E("display.GetWindowInfo", "window service not available", nil)
	}
	info, _ := result.(*window.WindowInfo)
	return info, nil
}

// ListWindowInfos returns information about all tracked windows via IPC.
// Use: infos := svc.ListWindowInfos()
func (s *Service) ListWindowInfos() []window.WindowInfo {
	result, handled, _ := s.Core().QUERY(window.QueryWindowList{})
	if !handled {
		return nil
	}
	list, _ := result.([]window.WindowInfo)
	return list
}

// SetWindowPosition moves a window via IPC.
// Use: _ = svc.SetWindowPosition("editor", 160, 120)
func (s *Service) SetWindowPosition(name string, x, y int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetPosition{Name: name, X: x, Y: y})
	return err
}

// SetWindowSize resizes a window via IPC.
// Use: _ = svc.SetWindowSize("editor", 1280, 800)
func (s *Service) SetWindowSize(name string, width, height int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetSize{Name: name, Width: width, Height: height})
	return err
}

// SetWindowBounds sets both position and size of a window via IPC.
// Use: _ = svc.SetWindowBounds("editor", 160, 120, 1280, 800)
func (s *Service) SetWindowBounds(name string, x, y, width, height int) error {
	if _, _, err := s.Core().PERFORM(window.TaskSetPosition{Name: name, X: x, Y: y}); err != nil {
		return err
	}
	_, _, err := s.Core().PERFORM(window.TaskSetSize{Name: name, Width: width, Height: height})
	return err
}

// MaximizeWindow maximizes a window via IPC.
// Use: _ = svc.MaximizeWindow("editor")
func (s *Service) MaximizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMaximise{Name: name})
	return err
}

// MinimizeWindow minimizes a window via IPC.
// Use: _ = svc.MinimizeWindow("editor")
func (s *Service) MinimizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMinimise{Name: name})
	return err
}

// FocusWindow brings a window to the front via IPC.
// Use: _ = svc.FocusWindow("editor")
func (s *Service) FocusWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskFocus{Name: name})
	return err
}

// FocusSet is a compatibility alias for FocusWindow.
// Use: _ = svc.FocusSet("editor")
func (s *Service) FocusSet(name string) error {
	return s.FocusWindow(name)
}

// CloseWindow closes a window via IPC.
// Use: _ = svc.CloseWindow("editor")
func (s *Service) CloseWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskCloseWindow{Name: name})
	return err
}

// RestoreWindow restores a maximized/minimized window.
func (s *Service) RestoreWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskRestore{Name: name})
	return err
}

// SetWindowVisibility shows or hides a window.
func (s *Service) SetWindowVisibility(name string, visible bool) error {
	_, _, err := s.Core().PERFORM(window.TaskSetVisibility{Name: name, Visible: visible})
	return err
}

// SetWindowAlwaysOnTop sets whether a window stays on top.
func (s *Service) SetWindowAlwaysOnTop(name string, alwaysOnTop bool) error {
	_, _, err := s.Core().PERFORM(window.TaskSetAlwaysOnTop{Name: name, AlwaysOnTop: alwaysOnTop})
	return err
}

// SetWindowTitle changes a window's title.
func (s *Service) SetWindowTitle(name string, title string) error {
	_, _, err := s.Core().PERFORM(window.TaskSetTitle{Name: name, Title: title})
	return err
}

// SetWindowFullscreen sets a window to fullscreen mode.
func (s *Service) SetWindowFullscreen(name string, fullscreen bool) error {
	_, _, err := s.Core().PERFORM(window.TaskFullscreen{Name: name, Fullscreen: fullscreen})
	return err
}

// SetWindowBackgroundColour sets the background colour of a window.
func (s *Service) SetWindowBackgroundColour(name string, r, g, b, a uint8) error {
	_, _, err := s.Core().PERFORM(window.TaskSetBackgroundColour{
		Name: name, Red: r, Green: g, Blue: b, Alpha: a,
	})
	return err
}

// GetFocusedWindow returns the name of the currently focused window.
// Use: focused := svc.GetFocusedWindow()
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
// Use: title, err := svc.GetWindowTitle("editor")
func (s *Service) GetWindowTitle(name string) (string, error) {
	info, err := s.GetWindowInfo(name)
	if err != nil {
		return "", err
	}
	if info == nil {
		return "", coreerr.E("display.GetWindowTitle", "window not found: "+name, nil)
	}
	return info.Title, nil
}

// ResetWindowState clears saved window positions.
// Use: _ = svc.ResetWindowState()
func (s *Service) ResetWindowState() error {
	ws := s.windowService()
	if ws != nil {
		ws.Manager().State().Clear()
	}
	return nil
}

// GetSavedWindowStates returns all saved window states.
// Use: states := svc.GetSavedWindowStates()
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

// CreateWindowOptions specifies the initial state for a new named window.
// svc.CreateWindow(display.CreateWindowOptions{Name: "settings", URL: "/settings", Width: 800, Height: 600})
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

func (s *Service) CreateWindow(options CreateWindowOptions) (*window.WindowInfo, error) {
	if options.Name == "" {
		return nil, coreerr.E("display.CreateWindow", "window name is required", nil)
	}
	result, _, err := s.Core().PERFORM(window.TaskOpenWindow{
		Window: &window.Window{
			Name:   options.Name,
			Title:  options.Title,
			URL:    options.URL,
			Width:  options.Width,
			Height: options.Height,
			X:      options.X,
			Y:      options.Y,
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
// Use: _ = svc.SaveLayout("coding")
func (s *Service) SaveLayout(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskSaveLayout{Name: name})
	return err
}

// RestoreLayout applies a saved layout.
// Use: _ = svc.RestoreLayout("coding")
func (s *Service) RestoreLayout(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskRestoreLayout{Name: name})
	return err
}

// ListLayouts returns all saved layout names with metadata.
// Use: layouts := svc.ListLayouts()
func (s *Service) ListLayouts() []window.LayoutInfo {
	result, handled, _ := s.Core().QUERY(window.QueryLayoutList{})
	if !handled {
		return nil
	}
	layouts, _ := result.([]window.LayoutInfo)
	return layouts
}

// DeleteLayout removes a saved layout by name.
// Use: _ = svc.DeleteLayout("coding")
func (s *Service) DeleteLayout(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskDeleteLayout{Name: name})
	return err
}

// GetLayout returns a specific layout by name.
// Use: layout := svc.GetLayout("coding")
func (s *Service) GetLayout(name string) *window.Layout {
	result, handled, _ := s.Core().QUERY(window.QueryLayoutGet{Name: name})
	if !handled {
		return nil
	}
	layout, _ := result.(*window.Layout)
	return layout
}

// --- Tiling/snapping delegation ---

// TileWindows arranges windows in a tiled layout.
// Use: _ = svc.TileWindows(window.TileModeLeftRight, []string{"left", "right"})
func (s *Service) TileWindows(mode window.TileMode, windowNames []string) error {
	_, _, err := s.Core().PERFORM(window.TaskTileWindows{Mode: mode.String(), Windows: windowNames})
	return err
}

// SnapWindow snaps a window to a screen edge or corner.
// Use: _ = svc.SnapWindow("editor", window.SnapRight)
func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	_, _, err := s.Core().PERFORM(window.TaskSnapWindow{Name: name, Position: position.String()})
	return err
}

// StackWindows arranges windows in a cascade pattern.
// Use: _ = svc.StackWindows([]string{"editor", "terminal"}, 24, 24)
func (s *Service) StackWindows(windowNames []string, offsetX, offsetY int) error {
	_, _, err := s.Core().PERFORM(window.TaskStackWindows{Windows: windowNames, OffsetX: offsetX, OffsetY: offsetY})
	return err
}

// ApplyWorkflowLayout applies a predefined layout for a specific workflow.
// Use: _ = svc.ApplyWorkflowLayout(window.WorkflowCoding)
func (s *Service) ApplyWorkflowLayout(workflow window.WorkflowLayout) error {
	_, _, err := s.Core().PERFORM(window.TaskApplyWorkflow{
		Workflow: workflow.String(),
	})
	return err
}

// GetEventManager returns the event manager for WebSocket event subscriptions.
// Use: events := svc.GetEventManager()
func (s *Service) GetEventManager() *WSEventManager {
	return s.events
}

// GetEventInfo returns a summary of the live WebSocket event server state.
// Use: info := svc.GetEventInfo()
func (s *Service) GetEventInfo() EventServerInfo {
	if s.events == nil {
		return EventServerInfo{}
	}
	return s.events.Info()
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
		Window: &window.Window{
			Name:   "workspace-new",
			Title:  "New Workspace",
			URL:    "/workspace/new",
			Width:  500,
			Height: 400,
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
		Window: &window.Window{
			Name:   "editor",
			Title:  "New File - Editor",
			URL:    "/#/developer/editor?new=true",
			Width:  1200,
			Height: 800,
		},
	})
}

func (s *Service) handleOpenFile() {
	result, handled, err := s.Core().PERFORM(dialog.TaskOpenFile{
		Options: dialog.OpenFileOptions{
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
		Window: &window.Window{
			Name:   "editor",
			Title:  paths[0] + " - Editor",
			URL:    "/#/developer/editor?file=" + paths[0],
			Width:  1200,
			Height: 800,
		},
	})
}

func (s *Service) handleSaveFile() { _ = s.Core().ACTION(ActionIDECommand{Command: "save"}) }
func (s *Service) handleOpenEditor() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Window: &window.Window{
			Name:   "editor",
			Title:  "Editor",
			URL:    "/#/developer/editor",
			Width:  1200,
			Height: 800,
		},
	})
}
func (s *Service) handleOpenTerminal() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Window: &window.Window{
			Name:   "terminal",
			Title:  "Terminal",
			URL:    "/#/developer/terminal",
			Width:  800,
			Height: 500,
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
