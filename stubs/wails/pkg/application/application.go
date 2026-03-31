package application

import (
	"sync"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Context mirrors the callback context type exposed by Wails.
type Context struct{}

// Logger is a minimal logger surface used by the GUI packages.
type Logger struct{}

func (l Logger) Info(message string, args ...any) {}

// RGBA stores a colour with alpha.
type RGBA struct {
	Red, Green, Blue, Alpha uint8
}

// NewRGBA constructs an RGBA value.
func NewRGBA(red, green, blue, alpha uint8) RGBA {
	return RGBA{Red: red, Green: green, Blue: blue, Alpha: alpha}
}

// MenuRole identifies a platform menu role.
type MenuRole int

const (
	AppMenu MenuRole = iota
	FileMenu
	EditMenu
	ViewMenu
	WindowMenu
	HelpMenu
)

// MenuItem is a minimal menu item implementation.
type MenuItem struct {
	Label       string
	Accelerator string
	Tooltip     string
	Checked     bool
	Enabled     bool
	onClick     func(*Context)
}

func (mi *MenuItem) SetAccelerator(accel string) { mi.Accelerator = accel }
func (mi *MenuItem) SetTooltip(text string)      { mi.Tooltip = text }
func (mi *MenuItem) SetChecked(checked bool)     { mi.Checked = checked }
func (mi *MenuItem) SetEnabled(enabled bool)     { mi.Enabled = enabled }
func (mi *MenuItem) OnClick(fn func(*Context))   { mi.onClick = fn }

// Menu is a minimal menu tree used by the GUI wrappers.
type Menu struct {
	Items []*MenuItem
}

func NewMenu() *Menu { return &Menu{} }

func (m *Menu) Add(label string) *MenuItem {
	item := &MenuItem{Label: label, Enabled: true}
	m.Items = append(m.Items, item)
	return item
}

func (m *Menu) AddSeparator() {
	m.Items = append(m.Items, &MenuItem{Label: "---"})
}

func (m *Menu) AddSubmenu(label string) *Menu {
	submenu := &Menu{}
	m.Items = append(m.Items, &MenuItem{Label: label})
	return submenu
}

func (m *Menu) AddRole(role MenuRole) {
	m.Items = append(m.Items, &MenuItem{Label: role.String(), Enabled: true})
}

func (role MenuRole) String() string {
	switch role {
	case AppMenu:
		return "app"
	case FileMenu:
		return "file"
	case EditMenu:
		return "edit"
	case ViewMenu:
		return "view"
	case WindowMenu:
		return "window"
	case HelpMenu:
		return "help"
	default:
		return "unknown"
	}
}

// MenuManager owns the application menu.
type MenuManager struct {
	applicationMenu *Menu
}

func (m *MenuManager) SetApplicationMenu(menu *Menu) { m.applicationMenu = menu }

// SystemTray represents a tray instance.
type SystemTray struct {
	icon           []byte
	templateIcon   []byte
	tooltip        string
	label          string
	menu           *Menu
	attachedWindow *WebviewWindow
}

func (t *SystemTray) SetIcon(data []byte)         { t.icon = append([]byte(nil), data...) }
func (t *SystemTray) SetTemplateIcon(data []byte) { t.templateIcon = append([]byte(nil), data...) }
func (t *SystemTray) SetTooltip(text string)      { t.tooltip = text }
func (t *SystemTray) SetLabel(text string)        { t.label = text }
func (t *SystemTray) SetMenu(menu *Menu)          { t.menu = menu }
func (t *SystemTray) AttachWindow(w *WebviewWindow) {
	t.attachedWindow = w
}

// SystemTrayManager creates tray instances.
type SystemTrayManager struct{}

func (m *SystemTrayManager) New() *SystemTray { return &SystemTray{} }

// WindowEventContext carries drag-and-drop details for a window event.
type WindowEventContext struct {
	droppedFiles []string
	dropDetails  *DropTargetDetails
}

func (c *WindowEventContext) DroppedFiles() []string {
	return append([]string(nil), c.droppedFiles...)
}

func (c *WindowEventContext) DropTargetDetails() *DropTargetDetails {
	if c.dropDetails == nil {
		return nil
	}
	details := *c.dropDetails
	return &details
}

// DropTargetDetails mirrors the fields consumed by the GUI wrappers.
type DropTargetDetails struct {
	ElementID string
}

// WindowEvent mirrors the event object passed to window callbacks.
type WindowEvent struct {
	ctx *WindowEventContext
}

func (e *WindowEvent) Context() *WindowEventContext {
	if e.ctx == nil {
		e.ctx = &WindowEventContext{}
	}
	return e.ctx
}

// WebviewWindowOptions configures a window instance.
type WebviewWindowOptions struct {
	Name                string
	Title               string
	URL                 string
	Width, Height       int
	X, Y                int
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	Frameless           bool
	Hidden              bool
	AlwaysOnTop         bool
	DisableResize       bool
	EnableFileDrop      bool
	BackgroundColour    RGBA
}

// WebviewWindow is a lightweight, in-memory window implementation.
type WebviewWindow struct {
	mu            sync.RWMutex
	opts          WebviewWindowOptions
	title         string
	x, y          int
	width, height int
	maximised     bool
	focused       bool
	visible       bool
	alwaysOnTop   bool
	fullscreen    bool
	closed        bool
	eventHandlers map[events.WindowEventType][]func(*WindowEvent)
}

func newWebviewWindow(options WebviewWindowOptions) *WebviewWindow {
	return &WebviewWindow{
		opts:          options,
		title:         options.Title,
		x:             options.X,
		y:             options.Y,
		width:         options.Width,
		height:        options.Height,
		visible:       !options.Hidden,
		alwaysOnTop:   options.AlwaysOnTop,
		eventHandlers: make(map[events.WindowEventType][]func(*WindowEvent)),
	}
}

func (w *WebviewWindow) Name() string { return w.opts.Name }
func (w *WebviewWindow) Title() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.title
}
func (w *WebviewWindow) Position() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.x, w.y
}
func (w *WebviewWindow) Size() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.width, w.height
}
func (w *WebviewWindow) IsMaximised() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.maximised
}
func (w *WebviewWindow) IsFocused() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.focused
}

func (w *WebviewWindow) SetTitle(title string) {
	w.mu.Lock()
	w.title = title
	w.mu.Unlock()
}

func (w *WebviewWindow) SetPosition(x, y int) {
	w.mu.Lock()
	w.x = x
	w.y = y
	w.mu.Unlock()
}

func (w *WebviewWindow) SetSize(width, height int) {
	w.mu.Lock()
	w.width = width
	w.height = height
	w.mu.Unlock()
}

func (w *WebviewWindow) SetBackgroundColour(colour RGBA) {}

func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) {
	w.mu.Lock()
	w.alwaysOnTop = alwaysOnTop
	w.mu.Unlock()
}

func (w *WebviewWindow) Maximise() {
	w.mu.Lock()
	w.maximised = true
	w.mu.Unlock()
}

func (w *WebviewWindow) Restore() {
	w.mu.Lock()
	w.maximised = false
	w.fullscreen = false
	w.mu.Unlock()
}

func (w *WebviewWindow) Minimise() {}

func (w *WebviewWindow) Focus() {
	w.mu.Lock()
	w.focused = true
	w.mu.Unlock()
}

func (w *WebviewWindow) Close() {
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()
}

func (w *WebviewWindow) Show() {
	w.mu.Lock()
	w.visible = true
	w.mu.Unlock()
}

func (w *WebviewWindow) Hide() {
	w.mu.Lock()
	w.visible = false
	w.mu.Unlock()
}

func (w *WebviewWindow) Fullscreen() {
	w.mu.Lock()
	w.fullscreen = true
	w.mu.Unlock()
}

func (w *WebviewWindow) UnFullscreen() {
	w.mu.Lock()
	w.fullscreen = false
	w.mu.Unlock()
}

func (w *WebviewWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	w.mu.Lock()
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], callback)
	w.mu.Unlock()
	return func() {}
}

// WindowManager manages in-memory windows.
type WindowManager struct {
	mu      sync.RWMutex
	windows []*WebviewWindow
}

func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow {
	window := newWebviewWindow(options)
	wm.mu.Lock()
	wm.windows = append(wm.windows, window)
	wm.mu.Unlock()
	return window
}

func (wm *WindowManager) GetAll() []any {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	out := make([]any, 0, len(wm.windows))
	for _, window := range wm.windows {
		out = append(out, window)
	}
	return out
}

// App is the top-level application object used by the GUI packages.
type App struct {
	Logger     Logger
	Window     WindowManager
	Menu       MenuManager
	SystemTray SystemTrayManager
}

func (a *App) Quit() {}

func (a *App) NewMenu() *Menu {
	return NewMenu()
}
