package display

import (
	"context"
	"runtime"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"forge.lthn.ai/core/config"

	"forge.lthn.ai/core/gui/pkg/chat"
	"forge.lthn.ai/core/gui/pkg/contextmenu"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/dock"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/events"
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

type Options struct{}

// WindowInfo is an alias for window.WindowInfo (backward compatibility).
type WindowInfo = window.WindowInfo

// Service orchestrates window, systray, and menu sub-services via IPC.
// Bridges IPC actions to WebSocket events for TypeScript apps.
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp       *application.App
	app            App
	configData     map[string]map[string]any
	configFile     *config.Config // config instance for file persistence
	events         *WSEventManager
	schemeHandlers map[string]SchemeHandler
}

// New returns a display Service with empty config sections.
// s, _ := display.New(); s.loadConfigFrom("/path/to/config.yaml")
func New() (*Service, error) {
	return &Service{
		configData: map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
		},
		schemeHandlers: make(map[string]SchemeHandler),
	}, nil
}

// Register binds the display service to a Core instance.
// core.WithService(display.Register(app))      // production (Wails app)
// core.WithService(display.Register(nil))      // tests (no Wails runtime)
func Register(wailsApp *application.App) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		s, err := New()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
		s.wailsApp = wailsApp
		return core.Result{Value: s, OK: true}
	}
}

// OnStartup loads config and registers handlers before sub-services start.
// Config handlers are registered first — sub-services query them during their own OnStartup.
func (s *Service) OnStartup(_ context.Context) core.Result {
	s.loadConfig()

	// Register config query handler — available NOW for sub-services
	s.Core().RegisterQuery(s.handleConfigQuery)

	// Register config save actions
	s.Core().Action("display.saveWindowConfig", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(window.TaskSaveConfig)
		s.configData["window"] = t.Config
		s.persistSection("window", t.Config)
		return core.Result{OK: true}
	})
	s.Core().Action("display.saveSystrayConfig", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(systray.TaskSaveConfig)
		s.configData["systray"] = t.Config
		s.persistSection("systray", t.Config)
		return core.Result{OK: true}
	})
	s.Core().Action("display.saveMenuConfig", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(menu.TaskSaveConfig)
		s.configData["menu"] = t.Config
		s.persistSection("menu", t.Config)
		return core.Result{OK: true}
	})
	s.Core().Action("display.buildPreload", func(_ context.Context, opts core.Options) core.Result {
		script, err := s.BuildPreloadScript(opts.String("url"))
		return core.Result{}.New(script, err)
	})
	s.Core().Action("display.resolveScheme", func(ctx context.Context, opts core.Options) core.Result {
		return s.ResolveScheme(ctx, opts.String("url"))
	})
	s.registerDefaultSchemes()

	// Initialise Wails wrappers if app is available (nil in tests)
	if s.wailsApp != nil {
		s.app = newWailsApp(s.wailsApp)
		s.events = NewWSEventManager()
	}

	return core.Result{OK: true}
}

// HandleIPCEvents bridges IPC actions from sub-services to WebSocket events for TS apps.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) core.Result {
	switch m := msg.(type) {
	case core.ActionServiceStartup:
		// All services have completed OnStartup — safe to call sub-services
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
	case ActionIDECommand:
		if s.events != nil {
			s.events.Emit(Event{Type: EventIDECommand,
				Data: map[string]any{"command": m.Command}})
		}
	case events.ActionEventFired:
		if s.events != nil {
			s.events.Emit(Event{Type: EventCustomEvent,
				Data: map[string]any{"name": m.Event.Name, "data": m.Event.Data}})
		}
	case dock.ActionProgressChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventDockProgress,
				Data: map[string]any{"progress": m.Progress}})
		}
	case dock.ActionBounceStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventDockBounce,
				Data: map[string]any{"requestId": m.RequestID, "bounceType": m.BounceType}})
		}
	case notification.ActionNotificationActionTriggered:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationAction,
				Data: map[string]any{"notificationId": m.NotificationID, "actionId": m.ActionID}})
		}
	case notification.ActionNotificationDismissed:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationDismiss,
				Data: map[string]any{"id": m.ID}})
		}
	case chat.ActionConversationCreated:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatConversation,
				Data: map[string]any{"action": "created", "conversation": m.Conversation}})
		}
	case chat.ActionConversationUpdated:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatConversation,
				Data: map[string]any{"action": "updated", "conversation": m.Conversation}})
		}
	case chat.ActionConversationDeleted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatConversation,
				Data: map[string]any{"action": "deleted", "conversationId": m.ConversationID}})
		}
	case chat.ActionConversationCleared:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatConversation,
				Data: map[string]any{"action": "cleared", "conversationId": m.ConversationID}})
		}
	case chat.ActionMessageAdded:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatMessage,
				Data: map[string]any{"conversationId": m.ConversationID, "message": m.Message}})
		}
	case chat.ActionStreamStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatMessage,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"streamId":       m.StreamID,
					"state":          "started",
				}})
		}
	case chat.ActionTokenAppended:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatToken,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"content":        m.Content,
				}})
		}
	case chat.ActionStreamFinished:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatMessage,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"state":          "finished",
					"finishReason":   m.FinishReason,
				}})
		}
	case chat.ActionThinkingStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatThinkingStart,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"startedAt":      m.StartedAt,
				}})
		}
	case chat.ActionThinkingAppended:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatThinkingAppend,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"content":        m.Content,
				}})
		}
	case chat.ActionThinkingEnded:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatThinkingEnd,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"durationMs":     m.DurationMS,
				}})
		}
	case chat.ActionToolCallStarted:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatToolCall,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"call":           m.Call,
				}})
		}
	case chat.ActionToolResultReady:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatToolResult,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"messageId":      m.MessageID,
					"result":         m.Result,
				}})
		}
	case chat.ActionImageQueued:
		if s.events != nil {
			s.events.Emit(Event{Type: EventChatImageQueued,
				Data: map[string]any{
					"conversationId": m.ConversationID,
					"attachment":     m.Attachment,
				}})
		}
	}
	return core.Result{OK: true}
}

// WSMessage represents a command received from a WebSocket client.
type WSMessage struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data,omitempty"`
}

// requireStringField extracts a string field from WebSocket data and fails when it is missing.
func requireStringField(data map[string]any, key string) (string, error) {
	v, _ := data[key].(string)
	if v == "" {
		return "", coreerr.E("display.requireStringField", "missing required field \""+key+"\"", nil)
	}
	return v, nil
}

// wsRequire is kept for backward compatibility inside the display package.
func wsRequire(data map[string]any, key string) (string, error) {
	return requireStringField(data, key)
}

func optionsFromMap(data map[string]any) core.Options {
	items := make([]core.Option, 0, len(data))
	for key, value := range data {
		items = append(items, core.Option{Key: key, Value: value})
	}
	return core.NewOptions(items...)
}

// wsOptions is kept for backward compatibility inside the display package.
func wsOptions(data map[string]any) core.Options {
	return optionsFromMap(data)
}

// handleWSMessage bridges WebSocket commands to IPC calls.
func (s *Service) handleWSMessage(msg WSMessage) core.Result {
	ctx := context.Background()
	c := s.Core()

	switch msg.Action {
	case "chat:send":
		return c.Action("gui.chat.send").Run(ctx, wsOptions(msg.Data))
	case "chat:clear":
		return c.Action("gui.chat.clear").Run(ctx, wsOptions(msg.Data))
	case "chat:history":
		return c.Action("gui.chat.history").Run(ctx, wsOptions(msg.Data))
	case "chat:models":
		return c.Action("gui.chat.models").Run(ctx, wsOptions(msg.Data))
	case "chat:select-model":
		return c.Action("gui.chat.selectModel").Run(ctx, wsOptions(msg.Data))
	case "chat:settings:save":
		return c.Action("gui.chat.settings.save").Run(ctx, wsOptions(msg.Data))
	case "chat:settings:load":
		return c.Action("gui.chat.settings.load").Run(ctx, wsOptions(msg.Data))
	case "chat:settings:reset":
		return c.Action("gui.chat.settings.reset").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:list":
		return c.Action("gui.chat.conversations.list").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:get":
		return c.Action("gui.chat.conversations.get").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:delete":
		return c.Action("gui.chat.conversations.delete").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:search":
		return c.Action("gui.chat.conversations.search").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:new":
		return c.Action("gui.chat.conversations.new").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:rename":
		return c.Action("gui.chat.conversations.rename").Run(ctx, wsOptions(msg.Data))
	case "chat:conversations:export":
		return c.Action("gui.chat.conversations.export").Run(ctx, wsOptions(msg.Data))
	case "chat:attach-image":
		return c.Action("gui.chat.attachImage").Run(ctx, wsOptions(msg.Data))
	case "keybinding:add":
		accelerator, _ := msg.Data["accelerator"].(string)
		description, _ := msg.Data["description"].(string)
		return c.Action("keybinding.add").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: keybinding.TaskAdd{Accelerator: accelerator, Description: description}},
		))
	case "keybinding:remove":
		accelerator, _ := msg.Data["accelerator"].(string)
		return c.Action("keybinding.remove").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: keybinding.TaskRemove{Accelerator: accelerator}},
		))
	case "keybinding:list":
		return c.QUERY(keybinding.QueryList{})
	case "browser:open-url":
		url, _ := msg.Data["url"].(string)
		return c.Action("browser.openURL").Run(ctx, core.NewOptions(
			core.Option{Key: "url", Value: url},
		))
	case "browser:open-file":
		path, _ := msg.Data["path"].(string)
		return c.Action("browser.openFile").Run(ctx, core.NewOptions(
			core.Option{Key: "path", Value: path},
		))
	case "dock:show":
		return c.Action("dock.showIcon").Run(ctx, core.NewOptions())
	case "dock:hide":
		return c.Action("dock.hideIcon").Run(ctx, core.NewOptions())
	case "dock:badge":
		label, _ := msg.Data["label"].(string)
		return c.Action("dock.setBadge").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: dock.TaskSetBadge{Label: label}},
		))
	case "dock:badge-remove":
		return c.Action("dock.removeBadge").Run(ctx, core.NewOptions())
	case "dock:visible":
		return c.QUERY(dock.QueryVisible{})
	case "contextmenu:add":
		name, _ := msg.Data["name"].(string)
		marshalResult := core.JSONMarshal(msg.Data["menu"])
		var menuDef contextmenu.ContextMenuDef
		if marshalResult.OK {
			menuJSON, _ := marshalResult.Value.([]byte)
			core.JSONUnmarshal(menuJSON, &menuDef)
		}
		return c.Action("contextmenu.add").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: contextmenu.TaskAdd{Name: name, Menu: menuDef}},
		))
	case "contextmenu:remove":
		name, _ := msg.Data["name"].(string)
		return c.Action("contextmenu.remove").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: contextmenu.TaskRemove{Name: name}},
		))
	case "contextmenu:get":
		name, _ := msg.Data["name"].(string)
		return c.QUERY(contextmenu.QueryGet{Name: name})
	case "contextmenu:list":
		return c.QUERY(contextmenu.QueryList{})
	case "webview:eval":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		script, _ := msg.Data["script"].(string)
		return c.Action("webview.evaluate").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskEvaluate{Window: w, Script: script}},
		))
	case "webview:click":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.Action("webview.click").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskClick{Window: w, Selector: sel}},
		))
	case "webview:type":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		text, _ := msg.Data["text"].(string)
		return c.Action("webview.type").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskType{Window: w, Selector: sel, Text: text}},
		))
	case "webview:navigate":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		url, e := wsRequire(msg.Data, "url")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.Action("webview.navigate").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskNavigate{Window: w, URL: url}},
		))
	case "webview:screenshot":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.Action("webview.screenshot").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskScreenshot{Window: w}},
		))
	case "webview:scroll":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		x, _ := msg.Data["x"].(float64)
		y, _ := msg.Data["y"].(float64)
		return c.Action("webview.scroll").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskScroll{Window: w, X: int(x), Y: int(y)}},
		))
	case "webview:hover":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.Action("webview.hover").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskHover{Window: w, Selector: sel}},
		))
	case "webview:select":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		val, _ := msg.Data["value"].(string)
		return c.Action("webview.select").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskSelect{Window: w, Selector: sel, Value: val}},
		))
	case "webview:check":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		checked, _ := msg.Data["checked"].(bool)
		return c.Action("webview.check").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskCheck{Window: w, Selector: sel, Checked: checked}},
		))
	case "webview:upload":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		pathsRaw, _ := msg.Data["paths"].([]any)
		var paths []string
		for _, p := range pathsRaw {
			if ps, ok := p.(string); ok {
				paths = append(paths, ps)
			}
		}
		return c.Action("webview.uploadFile").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskUploadFile{Window: w, Selector: sel, Paths: paths}},
		))
	case "webview:viewport":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		width, _ := msg.Data["width"].(float64)
		height, _ := msg.Data["height"].(float64)
		return c.Action("webview.setViewport").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskSetViewport{Window: w, Width: int(width), Height: int(height)}},
		))
	case "webview:clear-console":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.Action("webview.clearConsole").Run(ctx, core.NewOptions(
			core.Option{Key: "task", Value: webview.TaskClearConsole{Window: w}},
		))
	case "webview:console":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		level, _ := msg.Data["level"].(string)
		limit := 100
		if l, ok := msg.Data["limit"].(float64); ok {
			limit = int(l)
		}
		return c.QUERY(webview.QueryConsole{Window: w, Level: level, Limit: limit})
	case "webview:query":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.QUERY(webview.QuerySelector{Window: w, Selector: sel})
	case "webview:query-all":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, e := wsRequire(msg.Data, "selector")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.QUERY(webview.QuerySelectorAll{Window: w, Selector: sel})
	case "webview:dom-tree":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		sel, _ := msg.Data["selector"].(string) // selector optional for dom-tree (defaults to root)
		return c.QUERY(webview.QueryDOMTree{Window: w, Selector: sel})
	case "webview:url":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.QUERY(webview.QueryURL{Window: w})
	case "webview:title":
		w, e := wsRequire(msg.Data, "window")
		if e != nil {
			return core.Result{Value: e, OK: false}
		}
		return c.QUERY(webview.QueryTitle{Window: w})
	default:
		return core.Result{}
	}
}

// handleTrayAction processes tray menu item clicks.
func (s *Service) handleTrayAction(actionID string) {
	ctx := context.Background()
	c := s.Core()
	switch actionID {
	case "open-desktop":
		// Show all windows
		infos := s.ListWindowInfos()
		for _, info := range infos {
			_ = c.Action("window.focus").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: window.TaskFocus{Name: info.Name}},
			))
		}
	case "close-desktop":
		// Hide all windows — future: add TaskHideWindow
	case "env-info":
		// Query environment info via IPC and show as dialog
		r := c.QUERY(environment.QueryInfo{})
		if r.OK {
			info, _ := r.Value.(environment.EnvironmentInfo)
			details := "OS: " + info.OS + "\nArch: " + info.Arch + "\nPlatform: " +
				info.Platform.Name + " " + info.Platform.Version
			_ = c.Action("dialog.message").Run(ctx, core.NewOptions(
				core.Option{Key: "task", Value: dialog.TaskMessageDialog{
					Options: dialog.MessageDialogOptions{
						Type: dialog.DialogInfo, Title: "Environment",
						Message: details, Buttons: []string{"OK"},
					},
				}},
			))
		}
	case "quit":
		if s.app != nil {
			s.app.Quit()
		}
	}
}

func guiConfigPath() string {
	home := core.Env("HOME")
	if home == "" {
		return core.JoinPath(".core", "gui", "config.yaml")
	}
	return core.JoinPath(home, ".core", "gui", "config.yaml")
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
}

func (s *Service) handleConfigQuery(_ *core.Core, q core.Query) core.Result {
	switch q.(type) {
	case window.QueryConfig:
		return core.Result{Value: s.configData["window"], OK: true}
	case systray.QueryConfig:
		return core.Result{Value: s.configData["systray"], OK: true}
	case menu.QueryConfig:
		return core.Result{Value: s.configData["menu"], OK: true}
	case events.QueryServerInfo:
		if s.events == nil {
			return core.Result{Value: events.ServerInfo{}, OK: true}
		}
		return core.Result{Value: s.events.Info(), OK: true}
	default:
		return core.Result{}
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
	svc, ok := core.ServiceFor[*window.Service](s.Core(), "window")
	if !ok {
		return nil
	}
	return svc
}

// --- Window Management (delegates via IPC) ---

// OpenWindow creates a new window via IPC.
func (s *Service) OpenWindow(options ...window.WindowOption) error {
	spec, err := window.ApplyOptions(options...)
	if err != nil {
		return err
	}
	r := s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{Window: spec}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
		return coreerr.E("display.OpenWindow", "window.open action failed", nil)
	}
	return nil
}

// GetWindowInfo returns information about a window via IPC.
func (s *Service) GetWindowInfo(name string) (*window.WindowInfo, error) {
	r := s.Core().QUERY(window.QueryWindowByName{Name: name})
	if !r.OK {
		return nil, coreerr.E("display.GetWindowInfo", "window service not available", nil)
	}
	info, _ := r.Value.(*window.WindowInfo)
	return info, nil
}

// ListWindowInfos returns information about all tracked windows via IPC.
func (s *Service) ListWindowInfos() []window.WindowInfo {
	r := s.Core().QUERY(window.QueryWindowList{})
	if !r.OK {
		return nil
	}
	list, _ := r.Value.([]window.WindowInfo)
	return list
}

// SetWindowPosition moves a window via IPC.
func (s *Service) SetWindowPosition(name string, x, y int) error {
	r := s.Core().Action("window.setPosition").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetPosition{Name: name, X: x, Y: y}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowSize resizes a window via IPC.
func (s *Service) SetWindowSize(name string, width, height int) error {
	r := s.Core().Action("window.setSize").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetSize{Name: name, Width: width, Height: height}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowBounds sets both position and size of a window via IPC.
func (s *Service) SetWindowBounds(name string, x, y, width, height int) error {
	if err := s.SetWindowPosition(name, x, y); err != nil {
		return err
	}
	return s.SetWindowSize(name, width, height)
}

// MaximizeWindow maximizes a window via IPC.
func (s *Service) MaximizeWindow(name string) error {
	r := s.Core().Action("window.maximise").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskMaximise{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// MinimizeWindow minimizes a window via IPC.
func (s *Service) MinimizeWindow(name string) error {
	r := s.Core().Action("window.minimise").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskMinimise{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// FocusWindow brings a window to the front via IPC.
func (s *Service) FocusWindow(name string) error {
	r := s.Core().Action("window.focus").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFocus{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// CloseWindow closes a window via IPC.
func (s *Service) CloseWindow(name string) error {
	r := s.Core().Action("window.close").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskCloseWindow{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// RestoreWindow restores a maximized/minimized window.
func (s *Service) RestoreWindow(name string) error {
	r := s.Core().Action("window.restore").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRestore{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowVisibility shows or hides a window.
func (s *Service) SetWindowVisibility(name string, visible bool) error {
	r := s.Core().Action("window.setVisibility").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetVisibility{Name: name, Visible: visible}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowAlwaysOnTop sets whether a window stays on top.
func (s *Service) SetWindowAlwaysOnTop(name string, alwaysOnTop bool) error {
	r := s.Core().Action("window.setAlwaysOnTop").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetAlwaysOnTop{Name: name, AlwaysOnTop: alwaysOnTop}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowTitle changes a window's title.
func (s *Service) SetWindowTitle(name string, title string) error {
	r := s.Core().Action("window.setTitle").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetTitle{Name: name, Title: title}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowFullscreen sets a window to fullscreen mode.
func (s *Service) SetWindowFullscreen(name string, fullscreen bool) error {
	r := s.Core().Action("window.fullscreen").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskFullscreen{Name: name, Fullscreen: fullscreen}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SetWindowBackgroundColour sets the background colour of a window.
func (s *Service) SetWindowBackgroundColour(name string, r, g, b, a uint8) error {
	result := s.Core().Action("window.setBackgroundColour").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetBackgroundColour{
			Name: name, Red: r, Green: g, Blue: b, Alpha: a,
		}},
	))
	if !result.OK {
		if e, ok := result.Value.(error); ok {
			return e
		}
	}
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
		return "", coreerr.E("display.GetWindowTitle", "window not found: "+name, nil)
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

// CreateWindowOptions specifies the initial state for a new named window.
// svc.CreateWindow(display.CreateWindowOptions{Name: "settings", URL: "/settings", Width: 800, Height: 600})
type CreateWindowOptions struct {
	Name   string `json:"name"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

func (s *Service) CreateWindow(options CreateWindowOptions) (*window.WindowInfo, error) {
	if options.Name == "" {
		return nil, coreerr.E("display.CreateWindow", "window name is required", nil)
	}
	r := s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   options.Name,
				Title:  options.Title,
				URL:    options.URL,
				Width:  options.Width,
				Height: options.Height,
				X:      options.X,
				Y:      options.Y,
			},
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, e
		}
		return nil, coreerr.E("display.CreateWindow", "window.open action failed", nil)
	}
	info, _ := r.Value.(window.WindowInfo)
	return &info, nil
}

// --- Layout delegation ---

// SaveLayout saves the current window arrangement as a named layout.
func (s *Service) SaveLayout(name string) error {
	r := s.Core().Action("window.saveLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSaveLayout{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// RestoreLayout applies a saved layout.
func (s *Service) RestoreLayout(name string) error {
	r := s.Core().Action("window.restoreLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRestoreLayout{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// ListLayouts returns all saved layout names with metadata.
func (s *Service) ListLayouts() []window.LayoutInfo {
	r := s.Core().QUERY(window.QueryLayoutList{})
	if !r.OK {
		return nil
	}
	layouts, _ := r.Value.([]window.LayoutInfo)
	return layouts
}

// DeleteLayout removes a saved layout by name.
func (s *Service) DeleteLayout(name string) error {
	r := s.Core().Action("window.deleteLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskDeleteLayout{Name: name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// GetLayout returns a specific layout by name.
func (s *Service) GetLayout(name string) *window.Layout {
	r := s.Core().QUERY(window.QueryLayoutGet{Name: name})
	if !r.OK {
		return nil
	}
	layout, _ := r.Value.(*window.Layout)
	return layout
}

// --- Tiling/snapping delegation ---

// TileWindows arranges windows in a tiled layout.
func (s *Service) TileWindows(mode window.TileMode, windowNames []string) error {
	r := s.Core().Action("window.tileWindows").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskTileWindows{Mode: mode.String(), Windows: windowNames}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// SnapWindow snaps a window to a screen edge or corner.
func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	r := s.Core().Action("window.snapWindow").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSnapWindow{Name: name, Position: position.String()}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// StackWindows arranges windows in a cascade pattern.
func (s *Service) StackWindows(windowNames []string, offsetX, offsetY int) error {
	r := s.Core().Action("window.stackWindows").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskStackWindows{Windows: windowNames, OffsetX: offsetX, OffsetY: offsetY}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
	}
	return nil
}

// ApplyWorkflowLayout applies a predefined layout for a specific workflow.
func (s *Service) ApplyWorkflowLayout(workflow window.WorkflowLayout) error {
	r := s.Core().Action("window.applyWorkflow").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskApplyWorkflow{Workflow: workflow.String()}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return e
		}
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
		{Role: pointerTo(menu.RoleAppMenu)},
		{Role: pointerTo(menu.RoleFileMenu)},
		{Role: pointerTo(menu.RoleViewMenu)},
		{Role: pointerTo(menu.RoleEditMenu)},
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
		{Role: pointerTo(menu.RoleWindowMenu)},
		{Role: pointerTo(menu.RoleHelpMenu)},
	}

	// On non-macOS, remove the AppMenu role
	if runtime.GOOS != "darwin" {
		items = items[1:] // skip AppMenu
	}

	_ = s.Core().Action("menu.setAppMenu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: menu.TaskSetAppMenu{Items: items}},
	))
}

// pointerTo returns a pointer to value.
func pointerTo[T any](value T) *T { return &value }

// --- Menu handler methods ---

func (s *Service) handleNewWorkspace() {
	_ = s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   "workspace-new",
				Title:  "New Workspace",
				URL:    "/workspace/new",
				Width:  500,
				Height: 400,
			},
		}},
	))
}

func (s *Service) handleListWorkspaces() {
	r := s.Core().Service("workspace")
	if !r.OK || r.Value == nil {
		return
	}
	lister, ok := r.Value.(interface{ ListWorkspaces() []string })
	if !ok {
		return
	}
	_ = lister.ListWorkspaces()
}

func (s *Service) handleNewFile() {
	_ = s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   "editor",
				Title:  "New File - Editor",
				URL:    "/#/developer/editor?new=true",
				Width:  1200,
				Height: 800,
			},
		}},
	))
}

func (s *Service) handleOpenFile() {
	r := s.Core().Action("dialog.openFile").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dialog.TaskOpenFile{
			Options: dialog.OpenFileOptions{
				Title:         "Open File",
				AllowMultiple: false,
			},
		}},
	))
	if !r.OK {
		return
	}
	paths, ok := r.Value.([]string)
	if !ok || len(paths) == 0 {
		return
	}
	_ = s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   "editor",
				Title:  paths[0] + " - Editor",
				URL:    "/#/developer/editor?file=" + paths[0],
				Width:  1200,
				Height: 800,
			},
		}},
	))
}

func (s *Service) handleSaveFile() { _ = s.Core().ACTION(ActionIDECommand{Command: "save"}) }
func (s *Service) handleOpenEditor() {
	_ = s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   "editor",
				Title:  "Editor",
				URL:    "/#/developer/editor",
				Width:  1200,
				Height: 800,
			},
		}},
	))
}

func (s *Service) handleOpenTerminal() {
	_ = s.Core().Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Window: &window.Window{
				Name:   "terminal",
				Title:  "Terminal",
				URL:    "/#/developer/terminal",
				Width:  800,
				Height: 500,
			},
		}},
	))
}
func (s *Service) handleRun()   { _ = s.Core().ACTION(ActionIDECommand{Command: "run"}) }
func (s *Service) handleBuild() { _ = s.Core().ACTION(ActionIDECommand{Command: "build"}) }

// --- Tray (setup delegated via IPC) ---

func (s *Service) setupTray() {
	_ = s.Core().Action("systray.setMenu").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayMenu{Items: []systray.TrayMenuItem{
			{Label: "Open Desktop", ActionID: "open-desktop"},
			{Label: "Close Desktop", ActionID: "close-desktop"},
			{Type: "separator"},
			{Label: "Environment Info", ActionID: "env-info"},
			{Type: "separator"},
			{Label: "Quit", ActionID: "quit"},
		}}},
	))
}
