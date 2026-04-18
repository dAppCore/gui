package application

import (
	"sync"
	"unsafe"

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
	submenu     *Menu
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
	item := NewSubMenuItem(label)
	m.Items = append(m.Items, item)
	return item.GetSubmenu()
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

func (m *MenuManager) SetApplicationMenu(menu *Menu) {
	if m == nil {
		return
	}
	m.applicationMenu = menu
}

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
	if c == nil {
		return nil
	}
	return append([]string(nil), c.droppedFiles...)
}

func (c *WindowEventContext) DropTargetDetails() *DropTargetDetails {
	if c == nil {
		return nil
	}
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
	if e == nil {
		return nil
	}
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
	HTML                string
	JS                  string
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
	id            uint
	opts          WebviewWindowOptions
	title         string
	url           string
	html          string
	x, y          int
	width, height int
	maximised     bool
	focused       bool
	visible       bool
	alwaysOnTop   bool
	opacity       float64
	fullscreen    bool
	closed        bool
	execJSCalls   []string
	eventHandlers map[events.WindowEventType][]func(*WindowEvent)
}

func newWebviewWindow(id uint, options WebviewWindowOptions) *WebviewWindow {
	window := &WebviewWindow{
		opts:          options,
		id:            id,
		title:         options.Title,
		url:           options.URL,
		html:          options.HTML,
		x:             options.X,
		y:             options.Y,
		width:         options.Width,
		height:        options.Height,
		visible:       !options.Hidden,
		alwaysOnTop:   options.AlwaysOnTop,
		opacity:       1.0,
		eventHandlers: make(map[events.WindowEventType][]func(*WindowEvent)),
	}
	if options.JS != "" {
		window.execJSCalls = append(window.execJSCalls, options.JS)
	}
	return window
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

func (w *WebviewWindow) SetTitle(title string) Window {
	w.mu.Lock()
	w.title = title
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) SetPosition(x, y int) {
	w.mu.Lock()
	w.x = x
	w.y = y
	w.mu.Unlock()
}

func (w *WebviewWindow) SetSize(width, height int) Window {
	w.mu.Lock()
	w.width = width
	w.height = height
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window { return w }

func (w *WebviewWindow) SetAlwaysOnTop(alwaysOnTop bool) Window {
	w.mu.Lock()
	w.alwaysOnTop = alwaysOnTop
	w.mu.Unlock()
	return w
}

// SetOpacity sets the whole-window opacity.
//
//	w.SetOpacity(0.85)
func (w *WebviewWindow) SetOpacity(opacity float64) Window {
	w.mu.Lock()
	w.opacity = opacity
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) Maximise() Window {
	w.mu.Lock()
	w.maximised = true
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) Restore() {
	w.mu.Lock()
	w.maximised = false
	w.fullscreen = false
	w.mu.Unlock()
}

func (w *WebviewWindow) Minimise() Window { return w }

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

func (w *WebviewWindow) Show() Window {
	w.mu.Lock()
	w.visible = true
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) Hide() Window {
	w.mu.Lock()
	w.visible = false
	w.mu.Unlock()
	return w
}

func (w *WebviewWindow) Fullscreen() Window {
	w.mu.Lock()
	w.fullscreen = true
	w.mu.Unlock()
	return w
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

// ID returns a stable numeric identifier for this window.
//
//	id := w.ID()
func (w *WebviewWindow) ID() uint { return w.id }

// ClientID returns the client identifier (empty for native windows).
//
//	cid := w.ClientID()
func (w *WebviewWindow) ClientID() string { return "" }

// Width returns the current window width in logical pixels.
//
//	px := w.Width()
func (w *WebviewWindow) Width() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.width
}

// Height returns the current window height in logical pixels.
//
//	px := w.Height()
func (w *WebviewWindow) Height() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.height
}

// IsVisible reports whether the window is currently shown.
//
//	if w.IsVisible() { ... }
func (w *WebviewWindow) IsVisible() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.visible
}

// IsFullscreen reports whether the window is in fullscreen mode.
//
//	if w.IsFullscreen() { ... }
func (w *WebviewWindow) IsFullscreen() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.fullscreen
}

// IsMinimised reports whether the window is minimised.
//
//	if w.IsMinimised() { ... }
func (w *WebviewWindow) IsMinimised() bool { return false }

// IsIgnoreMouseEvents reports whether mouse events are being suppressed.
//
//	if w.IsIgnoreMouseEvents() { ... }
func (w *WebviewWindow) IsIgnoreMouseEvents() bool { return false }

// Resizable reports whether the window can be resized by the user.
//
//	if w.Resizable() { ... }
func (w *WebviewWindow) Resizable() bool { return true }

// Bounds returns the current position and size as a Rect.
//
//	r := w.Bounds()
func (w *WebviewWindow) Bounds() Rect {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return Rect{X: w.x, Y: w.y, Width: w.width, Height: w.height}
}

// SetBounds sets position and size simultaneously.
//
//	w.SetBounds(Rect{X: 100, Y: 100, Width: 1280, Height: 800})
func (w *WebviewWindow) SetBounds(bounds Rect) {
	w.mu.Lock()
	w.x, w.y, w.width, w.height = bounds.X, bounds.Y, bounds.Width, bounds.Height
	w.mu.Unlock()
}

// RelativePosition returns the position relative to the screen origin.
//
//	rx, ry := w.RelativePosition()
func (w *WebviewWindow) RelativePosition() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.x, w.y
}

// SetRelativePosition sets the position relative to the screen.
//
//	w.SetRelativePosition(0, 0)
func (w *WebviewWindow) SetRelativePosition(x, y int) Window {
	w.mu.Lock()
	w.x = x
	w.y = y
	w.mu.Unlock()
	return w
}

// SetMinSize sets the minimum window dimensions.
//
//	w.SetMinSize(640, 480)
func (w *WebviewWindow) SetMinSize(minWidth, minHeight int) Window { return w }

// SetMaxSize sets the maximum window dimensions.
//
//	w.SetMaxSize(3840, 2160)
func (w *WebviewWindow) SetMaxSize(maxWidth, maxHeight int) Window { return w }

// Center positions the window at the centre of the screen.
//
//	w.Center()
func (w *WebviewWindow) Center() {}

// SetURL navigates the webview to the given URL.
//
//	w.SetURL("https://example.com")
func (w *WebviewWindow) SetURL(url string) Window {
	w.mu.Lock()
	w.url = url
	w.mu.Unlock()
	return w
}

// SetHTML replaces the webview content with the given HTML string.
//
//	w.SetHTML("<h1>Hello</h1>")
func (w *WebviewWindow) SetHTML(html string) Window {
	w.mu.Lock()
	w.html = html
	w.mu.Unlock()
	return w
}

// SetFrameless toggles the window frame.
//
//	w.SetFrameless(true)
func (w *WebviewWindow) SetFrameless(frameless bool) Window { return w }

// SetResizable controls whether the user can resize the window.
//
//	w.SetResizable(false)
func (w *WebviewWindow) SetResizable(b bool) Window { return w }

// SetIgnoreMouseEvents suppresses or restores mouse event delivery.
//
//	w.SetIgnoreMouseEvents(true)
func (w *WebviewWindow) SetIgnoreMouseEvents(ignore bool) Window { return w }

// SetMinimiseButtonState controls the minimise button appearance.
//
//	w.SetMinimiseButtonState(ButtonHidden)
func (w *WebviewWindow) SetMinimiseButtonState(state ButtonState) Window { return w }

// SetMaximiseButtonState controls the maximise button appearance.
//
//	w.SetMaximiseButtonState(ButtonDisabled)
func (w *WebviewWindow) SetMaximiseButtonState(state ButtonState) Window { return w }

// SetCloseButtonState controls the close button appearance.
//
//	w.SetCloseButtonState(ButtonEnabled)
func (w *WebviewWindow) SetCloseButtonState(state ButtonState) Window { return w }

// SetEnabled enables or disables user interaction with the window.
//
//	w.SetEnabled(false)
func (w *WebviewWindow) SetEnabled(enabled bool) {}

// SetContentProtection prevents the window contents from being captured.
//
//	w.SetContentProtection(true)
func (w *WebviewWindow) SetContentProtection(protection bool) Window { return w }

// SetMenu attaches a menu to the window.
//
//	w.SetMenu(myMenu)
func (w *WebviewWindow) SetMenu(menu *Menu) {}

// ShowMenuBar makes the menu bar visible.
//
//	w.ShowMenuBar()
func (w *WebviewWindow) ShowMenuBar() {}

// HideMenuBar hides the menu bar.
//
//	w.HideMenuBar()
func (w *WebviewWindow) HideMenuBar() {}

// ToggleMenuBar toggles menu bar visibility.
//
//	w.ToggleMenuBar()
func (w *WebviewWindow) ToggleMenuBar() {}

// ToggleFrameless toggles the window frame.
//
//	w.ToggleFrameless()
func (w *WebviewWindow) ToggleFrameless() {}

// ExecJS executes a JavaScript string in the webview.
//
//	w.ExecJS("document.title = 'Ready'")
func (w *WebviewWindow) ExecJS(js string) {
	w.mu.Lock()
	w.execJSCalls = append(w.execJSCalls, js)
	w.mu.Unlock()
}

// Reload reloads the current page.
//
//	w.Reload()
func (w *WebviewWindow) Reload() {}

// ForceReload bypasses the cache and reloads.
//
//	w.ForceReload()
func (w *WebviewWindow) ForceReload() {}

// OpenDevTools opens the browser developer tools panel.
//
//	w.OpenDevTools()
func (w *WebviewWindow) OpenDevTools() {}

// OpenContextMenu triggers a named context menu at the given position.
//
//	w.OpenContextMenu(&ContextMenuData{Name: "edit", X: 100, Y: 200})
func (w *WebviewWindow) OpenContextMenu(data *ContextMenuData) {}

// Zoom applies the default zoom level.
//
//	w.Zoom()
func (w *WebviewWindow) Zoom() {}

// ZoomIn increases the zoom level by one step.
//
//	w.ZoomIn()
func (w *WebviewWindow) ZoomIn() {}

// ZoomOut decreases the zoom level by one step.
//
//	w.ZoomOut()
func (w *WebviewWindow) ZoomOut() {}

// ZoomReset returns the zoom level to 1.0.
//
//	w.ZoomReset()
func (w *WebviewWindow) ZoomReset() Window { return w }

// GetZoom returns the current zoom magnification factor.
//
//	z := w.GetZoom()
func (w *WebviewWindow) GetZoom() float64 { return 1.0 }

// GetOpacity returns the current window opacity.
//
//	alpha := w.GetOpacity()
func (w *WebviewWindow) GetOpacity() float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.opacity
}

// SetZoom sets the zoom magnification factor.
//
//	w.SetZoom(1.5)
func (w *WebviewWindow) SetZoom(magnification float64) Window { return w }

// RegisterHook registers a pre-event hook for the given window event type.
//
//	cancel := w.RegisterHook(events.Common.WindowClose, func(e *WindowEvent) { saveState() })
//	defer cancel()
func (w *WebviewWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}

// EmitEvent fires a named event from this window.
//
//	w.EmitEvent("user:login", payload)
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool { return false }

// DispatchWailsEvent sends a custom event through the Wails event bus.
//
//	w.DispatchWailsEvent(&CustomEvent{Name: "init"})
func (w *WebviewWindow) DispatchWailsEvent(event *CustomEvent) {}

// GetScreen returns the screen on which this window is currently displayed.
//
//	screen, err := w.GetScreen()
func (w *WebviewWindow) GetScreen() (*Screen, error) { return &Screen{}, nil }

// GetBorderSizes returns the platform-specific window border dimensions.
//
//	borders := w.GetBorderSizes()
func (w *WebviewWindow) GetBorderSizes() *LRTB { return &LRTB{} }

// EnableSizeConstraints activates the min/max size limits.
//
//	w.EnableSizeConstraints()
func (w *WebviewWindow) EnableSizeConstraints() {}

// DisableSizeConstraints removes the min/max size limits.
//
//	w.DisableSizeConstraints()
func (w *WebviewWindow) DisableSizeConstraints() {}

// AttachModal registers a modal window that blocks this window.
//
//	w.AttachModal(confirmDialog)
func (w *WebviewWindow) AttachModal(modalWindow Window) {}

// Flash requests the window manager to flash or bounce this window.
//
//	w.Flash(true)
func (w *WebviewWindow) Flash(enabled bool) {}

// Print opens the system print dialog for the webview contents.
//
//	err := w.Print()
func (w *WebviewWindow) Print() error { return nil }

// Error logs an error-level message on behalf of this window.
//
//	w.Error("load failed: %s", err)
func (w *WebviewWindow) Error(message string, args ...any) {}

// Info logs an info-level message on behalf of this window.
//
//	w.Info("window ready")
func (w *WebviewWindow) Info(message string, args ...any) {}

// NativeWindow returns the platform-specific window handle (nil in stub).
//
//	ptr := w.NativeWindow()
func (w *WebviewWindow) NativeWindow() unsafe.Pointer { return nil }

// Run starts the window event loop.
//
//	w.Run()
func (w *WebviewWindow) Run() {}

// UnMaximise restores the window from maximised state.
//
//	w.UnMaximise()
func (w *WebviewWindow) UnMaximise() {
	w.mu.Lock()
	w.maximised = false
	w.mu.Unlock()
}

// UnMinimise restores the window from minimised state.
//
//	w.UnMinimise()
func (w *WebviewWindow) UnMinimise() {}

// ToggleFullscreen switches between fullscreen and windowed mode.
//
//	w.ToggleFullscreen()
func (w *WebviewWindow) ToggleFullscreen() {
	w.mu.Lock()
	w.fullscreen = !w.fullscreen
	w.mu.Unlock()
}

// ToggleMaximise switches between maximised and restored state.
//
//	w.ToggleMaximise()
func (w *WebviewWindow) ToggleMaximise() {
	w.mu.Lock()
	w.maximised = !w.maximised
	w.mu.Unlock()
}

// SnapAssist triggers the platform snap-assist feature.
//
//	w.SnapAssist()
func (w *WebviewWindow) SnapAssist() {}

// Internal platform hooks — no-ops in the stub.

func (w *WebviewWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {}
func (w *WebviewWindow) InitiateFrontendDropProcessing(filenames []string, x int, y int)            {}
func (w *WebviewWindow) HandleMessage(message string)                                               {}
func (w *WebviewWindow) HandleWindowEvent(id uint)                                                  {}
func (w *WebviewWindow) HandleKeyEvent(acceleratorString string)                                    {}
func (w *WebviewWindow) shouldUnconditionallyClose() bool                                           { return false }
func (w *WebviewWindow) cut()                                                                       {}
func (w *WebviewWindow) copy()                                                                      {}
func (w *WebviewWindow) paste()                                                                     {}
func (w *WebviewWindow) undo()                                                                      {}
func (w *WebviewWindow) redo()                                                                      {}
func (w *WebviewWindow) delete()                                                                    {}
func (w *WebviewWindow) selectAll()                                                                 {}

// WindowManager manages in-memory windows.
type WindowManager struct {
	mu      sync.RWMutex
	nextID  uint
	windows []*WebviewWindow
}

func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow {
	if wm == nil {
		return nil
	}
	wm.mu.Lock()
	wm.nextID++
	window := newWebviewWindow(wm.nextID, options)
	wm.windows = append(wm.windows, window)
	wm.mu.Unlock()
	return window
}

// GetAll returns all windows managed by this manager.
//
//	for _, w := range wm.GetAll() { w.Show() }
func (wm *WindowManager) GetAll() []Window {
	if wm == nil {
		return nil
	}
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	out := make([]Window, 0, len(wm.windows))
	for _, window := range wm.windows {
		out = append(out, window)
	}
	return out
}

// App is the top-level application object used by the GUI packages.
//
//	app := &application.App{}
//	app.Dialog.Info().SetTitle("Done").SetMessage("Saved.").Show()
//	app.Event.Emit("user:login", payload)
type App struct {
	Logger      Logger
	Window      WindowManager
	Menu        MenuManager
	SystemTray  SystemTrayManager
	Dialog      DialogManager
	Event       EventManager
	Browser     BrowserManager
	Clipboard   ClipboardManager
	ContextMenu ContextMenuManager
	Environment EnvironmentManager
	Screen      ScreenManager
	KeyBinding  KeyBindingManager
}

func (a *App) Quit() {}

func (a *App) NewMenu() *Menu {
	return NewMenu()
}
