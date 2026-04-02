package application

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// RGBA represents a colour.
type RGBA struct {
	R, G, B, A uint8
}

// NewRGBA creates a colour value.
func NewRGBA(r, g, b, a uint8) RGBA { return RGBA{R: r, G: g, B: b, A: a} }

// Logger is a minimal logger used by the repo.
type Logger struct{}

func (l *Logger) Info(message string, args ...any) {}

// Context carries event data.
type Context struct {
	droppedFiles   []string
	dropTargetData *DropTargetDetails
}

func (c *Context) DroppedFiles() []string {
	if c == nil {
		return nil
	}
	out := make([]string, len(c.droppedFiles))
	copy(out, c.droppedFiles)
	return out
}

func (c *Context) DropTargetDetails() *DropTargetDetails {
	if c == nil || c.dropTargetData == nil {
		return nil
	}
	d := *c.dropTargetData
	return &d
}

// DropTargetDetails describes the drop target.
type DropTargetDetails struct {
	ElementID string
}

// WindowEvent wraps window event context.
type WindowEvent struct {
	ctx *Context
}

func (e *WindowEvent) Context() *Context {
	if e == nil {
		return nil
	}
	if e.ctx == nil {
		e.ctx = &Context{}
	}
	return e.ctx
}

// WebviewWindowOptions configures a new window.
type WebviewWindowOptions struct {
	Name             string
	Title            string
	URL              string
	Width            int
	Height           int
	X                int
	Y                int
	MinWidth         int
	MinHeight        int
	MaxWidth         int
	MaxHeight        int
	Frameless        bool
	Hidden           bool
	AlwaysOnTop      bool
	DisableResize    bool
	EnableFileDrop   bool
	BackgroundColour RGBA
}

// WebviewWindow is a lightweight in-memory window handle.
type WebviewWindow struct {
	opts          WebviewWindowOptions
	title         string
	x, y          int
	width, height int
	minimised     bool
	maximised     bool
	focused       bool
	visible       bool
	alwaysOnTop   bool
	fullscreen    bool
	devtoolsOpen  bool
	eventHandlers map[events.WindowEventType][]func(*WindowEvent)
	mu            sync.Mutex
}

func newWebviewWindow(opts WebviewWindowOptions) *WebviewWindow {
	return &WebviewWindow{
		opts:          opts,
		title:         opts.Title,
		x:             opts.X,
		y:             opts.Y,
		width:         opts.Width,
		height:        opts.Height,
		visible:       !opts.Hidden,
		alwaysOnTop:   opts.AlwaysOnTop,
		eventHandlers: make(map[events.WindowEventType][]func(*WindowEvent)),
	}
}

func (w *WebviewWindow) Name() string          { return w.opts.Name }
func (w *WebviewWindow) Position() (int, int)  { return w.x, w.y }
func (w *WebviewWindow) Size() (int, int)      { return w.width, w.height }
func (w *WebviewWindow) IsVisible() bool       { return w.visible }
func (w *WebviewWindow) IsMinimised() bool     { return w.minimised }
func (w *WebviewWindow) IsMaximised() bool     { return w.maximised }
func (w *WebviewWindow) IsFocused() bool       { return w.focused }
func (w *WebviewWindow) SetTitle(title string) { w.title = title }
func (w *WebviewWindow) SetPosition(x, y int)  { w.x, w.y = x, y }
func (w *WebviewWindow) SetSize(width, height int) {
	w.width, w.height = width, height
}
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) {}
func (w *WebviewWindow) SetOpacity(opacity float32)      {}
func (w *WebviewWindow) SetVisibility(visible bool)      { w.visible = visible }
func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) { w.alwaysOnTop = alwaysOnTop }
func (w *WebviewWindow) Maximise()                       { w.maximised = true; w.minimised = false; w.visible = true }
func (w *WebviewWindow) Restore()                        { w.maximised = false; w.minimised = false; w.visible = true }
func (w *WebviewWindow) Minimise()                       { w.minimised = true; w.maximised = false; w.visible = false }
func (w *WebviewWindow) Focus()                          { w.focused = true }
func (w *WebviewWindow) Close()                          {}
func (w *WebviewWindow) Show()                           { w.visible = true }
func (w *WebviewWindow) Hide()                           { w.visible = false }
func (w *WebviewWindow) Fullscreen()                     { w.fullscreen = true }
func (w *WebviewWindow) UnFullscreen()                   { w.fullscreen = false }
func (w *WebviewWindow) OpenDevTools()                   { w.devtoolsOpen = true }
func (w *WebviewWindow) CloseDevTools()                  { w.devtoolsOpen = false }
func (w *WebviewWindow) DevToolsOpen() bool              { return w.devtoolsOpen }

func (w *WebviewWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], callback)
	return func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		handlers := w.eventHandlers[eventType]
		if len(handlers) == 0 {
			return
		}
		w.eventHandlers[eventType] = handlers[:len(handlers)-1]
	}
}

func (w *WebviewWindow) trigger(eventType events.WindowEventType, event *WindowEvent) {
	w.mu.Lock()
	handlers := append([]func(*WindowEvent){}, w.eventHandlers[eventType]...)
	w.mu.Unlock()
	for _, handler := range handlers {
		handler(event)
	}
}

// WindowManager manages in-memory windows.
type WindowManager struct {
	windows []*WebviewWindow
}

func (wm *WindowManager) NewWithOptions(opts WebviewWindowOptions) *WebviewWindow {
	w := newWebviewWindow(opts)
	wm.windows = append(wm.windows, w)
	return w
}

func (wm *WindowManager) GetAll() []any {
	out := make([]any, len(wm.windows))
	for i, w := range wm.windows {
		out[i] = w
	}
	return out
}

// Menu role constants.
type MenuRole int

const (
	AppMenu MenuRole = iota
	FileMenu
	EditMenu
	ViewMenu
	WindowMenu
	HelpMenu
)

// Menu is a lightweight in-memory menu.
type Menu struct {
	items []*MenuItem
}

func NewMenu() *Menu { return &Menu{} }

func (m *Menu) Add(label string) *MenuItem {
	item := &MenuItem{label: label}
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AddSeparator() {}

func (m *Menu) AddSubmenu(label string) *Menu {
	return &Menu{}
}

func (m *Menu) AddRole(role MenuRole) {}

func (m *Menu) SetApplicationMenu(menu *Menu) {}

// MenuItem is a lightweight menu item.
type MenuItem struct {
	label       string
	accelerator string
	tooltip     string
	checked     bool
	enabled     bool
	onClick     func(*Context)
}

func (mi *MenuItem) SetAccelerator(accel string) *MenuItem {
	mi.accelerator = accel
	return mi
}

func (mi *MenuItem) SetTooltip(text string) *MenuItem {
	mi.tooltip = text
	return mi
}

func (mi *MenuItem) SetChecked(checked bool) *MenuItem {
	mi.checked = checked
	return mi
}

func (mi *MenuItem) SetEnabled(enabled bool) *MenuItem {
	mi.enabled = enabled
	return mi
}

func (mi *MenuItem) OnClick(fn func(ctx *Context)) *MenuItem {
	mi.onClick = fn
	return mi
}

// SystemTray models a tray icon.
type SystemTray struct {
	icon           []byte
	templateIcon   []byte
	tooltip        string
	label          string
	menu           *Menu
	attachedWindow interface {
		Show()
		Hide()
		Focus()
		IsVisible() bool
	}
	onClick func()
}

func (st *SystemTray) SetIcon(icon []byte) *SystemTray {
	st.icon = append([]byte(nil), icon...)
	return st
}

func (st *SystemTray) SetTemplateIcon(icon []byte) *SystemTray {
	st.templateIcon = append([]byte(nil), icon...)
	return st
}

func (st *SystemTray) SetTooltip(tooltip string) {
	st.tooltip = tooltip
}

func (st *SystemTray) SetLabel(label string) {
	st.label = label
}

func (st *SystemTray) SetMenu(menu *Menu) *SystemTray {
	st.menu = menu
	return st
}

func (st *SystemTray) Show() {}
func (st *SystemTray) Hide() {}
func (st *SystemTray) OnClick(callback func()) *SystemTray {
	st.onClick = callback
	return st
}

func (st *SystemTray) AttachWindow(window interface {
	Show()
	Hide()
	Focus()
	IsVisible() bool
}) *SystemTray {
	st.attachedWindow = window
	st.OnClick(func() {
		if st.attachedWindow == nil {
			return
		}
		if st.attachedWindow.IsVisible() {
			st.attachedWindow.Hide()
			return
		}
		st.attachedWindow.Show()
		st.attachedWindow.Focus()
	})
	return st
}

func (st *SystemTray) Click() {
	if st.onClick != nil {
		st.onClick()
	}
}

// SystemTrayManager creates trays.
type SystemTrayManager struct {
	app *App
}

func (stm *SystemTrayManager) New() *SystemTray { return &SystemTray{} }

// MenuManager manages application menus.
type MenuManager struct {
	appMenu *Menu
}

func (mm *MenuManager) SetApplicationMenu(menu *Menu) { mm.appMenu = menu }

// App is the top-level application container.
type App struct {
	Window     *WindowManager
	Menu       *MenuManager
	SystemTray *SystemTrayManager
	Logger     *Logger
	quit       bool
}

func NewApp() *App {
	app := &App{}
	app.Window = &WindowManager{}
	app.Menu = &MenuManager{}
	app.SystemTray = &SystemTrayManager{app: app}
	app.Logger = &Logger{}
	return app
}

func (a *App) Quit() { a.quit = true }

func (a *App) NewMenu() *Menu { return NewMenu() }
