package application

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"
)

type appCompatState struct {
	mu       sync.RWMutex
	options  Options
	ctx      context.Context
	services []Service
	shutdown []func()
	running  bool
	hidden   bool
	icon     []byte
}

type webviewCompatState struct {
	mu               sync.RWMutex
	physicalBounds   Rect
	physicalBoundsOK bool
	devtoolsOpen     bool
	keyBindings      map[string]func(Window)
}

type trayCompatState struct {
	mu                 sync.RWMutex
	darkModeIcon       []byte
	iconPosition       IconPosition
	visible            bool
	windowOffset       int
	windowDebounce     time.Duration
	rightClickHandler  func()
	doubleClickHandler func()
	rightDoubleClick   func()
	mouseEnterHandler  func()
	mouseLeaveHandler  func()
	onMenuOpen         func()
	onMenuClose        func()
}

var (
	appCompatStates     sync.Map
	webviewCompatStates sync.Map
	trayCompatStates    sync.Map

	defaultScreenManager = NewScreenManager()
)

func ensureAppCompatState(app *App) *appCompatState {
	if app == nil {
		return &appCompatState{ctx: context.Background()}
	}
	initialiseAppManagers(app, Options{})
	if state, ok := appCompatStates.Load(app); ok {
		return state.(*appCompatState)
	}
	state := &appCompatState{ctx: context.Background()}
	actual, _ := appCompatStates.LoadOrStore(app, state)
	return actual.(*appCompatState)
}

func newStubApp(options Options) *App {
	app := &App{}
	initialiseAppManagers(app, options)
	state := ensureAppCompatState(app)
	state.mu.Lock()
	state.options = options
	state.services = append([]Service(nil), options.Services...)
	state.mu.Unlock()
	for accelerator, callback := range options.KeyBindings {
		app.KeyBinding.Add(accelerator, callback)
	}
	return app
}

func initialiseAppManagers(app *App, options Options) {
	if app == nil {
		return
	}
	if app.Logger == nil {
		if options.Logger != nil {
			app.Logger = options.Logger
		} else {
			handlerOptions := &slog.HandlerOptions{Level: options.LogLevel}
			app.Logger = slog.New(slog.NewTextHandler(io.Discard, handlerOptions))
		}
	}
	if app.Window == nil {
		app.Window = &WindowManager{}
	}
	if app.Menu == nil {
		app.Menu = &MenuManager{}
	}
	if app.SystemTray == nil {
		app.SystemTray = &SystemTrayManager{}
	}
	if app.Dialog == nil {
		app.Dialog = &DialogManager{}
	}
	if app.Event == nil {
		app.Event = &EventManager{}
	}
	if app.Browser == nil {
		app.Browser = &BrowserManager{}
	}
	if app.Clipboard == nil {
		app.Clipboard = &ClipboardManager{}
	}
	if app.ContextMenu == nil {
		app.ContextMenu = &ContextMenuManager{}
	}
	if app.Env == nil && app.Environment != nil {
		app.Env = app.Environment
	}
	if app.Env == nil {
		app.Env = &EnvironmentManager{}
	}
	if app.Environment == nil {
		app.Environment = app.Env
	}
	if app.Screen == nil {
		app.Screen = NewScreenManager()
	}
	if app.KeyBinding == nil {
		app.KeyBinding = NewKeyBindingManager()
	}
}

func ensureWebviewCompatState(window *WebviewWindow) *webviewCompatState {
	if state, ok := webviewCompatStates.Load(window); ok {
		return state.(*webviewCompatState)
	}
	state := &webviewCompatState{
		physicalBounds: Rect{X: window.posX, Y: window.posY, Width: window.sizeW, Height: window.sizeH},
		keyBindings:    make(map[string]func(Window)),
	}
	for accelerator, callback := range window.opts.KeyBindings {
		state.keyBindings[accelerator] = callback
	}
	actual, _ := webviewCompatStates.LoadOrStore(window, state)
	return actual.(*webviewCompatState)
}

func ensureTrayCompatState(tray *SystemTray) *trayCompatState {
	if state, ok := trayCompatStates.Load(tray); ok {
		return state.(*trayCompatState)
	}
	state := &trayCompatState{
		iconPosition:   NSImageLeading,
		visible:        true,
		windowDebounce: 200 * time.Millisecond,
	}
	actual, _ := trayCompatStates.LoadOrStore(tray, state)
	return actual.(*trayCompatState)
}

func appScreens() *ScreenManager {
	if globalApplication != nil {
		return globalApplication.Screen
	}
	return defaultScreenManager
}
