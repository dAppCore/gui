package application

import (
	"context"
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
	if state, ok := appCompatStates.Load(app); ok {
		return state.(*appCompatState)
	}
	state := &appCompatState{ctx: context.Background()}
	actual, _ := appCompatStates.LoadOrStore(app, state)
	if app.KeyBinding.bindings == nil {
		app.KeyBinding = *NewKeyBindingManager()
	}
	return actual.(*appCompatState)
}

func newStubApp(options Options) *App {
	app := &App{}
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
		return &globalApplication.Screen
	}
	return defaultScreenManager
}
