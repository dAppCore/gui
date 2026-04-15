package application

import (
	"log/slog"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Context mirrors the callback context type exposed by Wails.
//
//	item.OnClick(func(ctx *Context) { openPrefs() })
type Context struct {
	clickedMenuItem *MenuItem
	contextMenuData *ContextMenuData
	checked         bool
}

func newContext() *Context { return &Context{} }

func (ctx *Context) withClickedMenuItem(item *MenuItem) *Context {
	ctx.clickedMenuItem = item
	return ctx
}

func (ctx *Context) withContextMenuData(data *ContextMenuData) *Context {
	ctx.contextMenuData = data
	return ctx
}

func (ctx *Context) withChecked(checked bool) *Context {
	ctx.checked = checked
	return ctx
}

// Logger is a minimal logger surface used by the GUI packages.
type Logger struct{}

func (l Logger) Info(message string, args ...any) {}

// RGBA stores a colour with alpha.
//
//	colour := NewRGBA(255, 128, 0, 255) // opaque orange
type RGBA struct {
	Red, Green, Blue, Alpha uint8
}

// NewRGBA constructs an RGBA value.
//
//	colour := NewRGBA(255, 128, 0, 255) // opaque orange
func NewRGBA(red, green, blue, alpha uint8) RGBA {
	return RGBA{Red: red, Green: green, Blue: blue, Alpha: alpha}
}

// Menu is a menu tree used by the GUI wrappers.
//
//	menu := NewMenu()
//	menu.Add("Save").SetAccelerator("CmdOrCtrl+S").OnClick(func(ctx *Context) { save() })
type Menu struct {
	label string
	Items []*MenuItem
}

// NewMenu creates a new, empty Menu.
//
//	menu := NewMenu()
func NewMenu() *Menu { return &Menu{} }

// Add appends a new text item with the given label.
func (m *Menu) Add(label string) *MenuItem {
	item := NewMenuItem(label)
	item.disabled = false
	m.Items = append(m.Items, item)
	return item
}

// AddSeparator appends a separator item.
func (m *Menu) AddSeparator() {
	m.Items = append(m.Items, NewMenuItemSeparator())
}

// AddSubmenu appends a submenu item and returns the child Menu.
func (m *Menu) AddSubmenu(label string) *Menu {
	item := NewSubMenuItem(label)
	m.Items = append(m.Items, item)
	return item.submenu
}

// AddRole appends a platform-role item.
func (m *Menu) AddRole(role Role) *Menu {
	m.Items = append(m.Items, NewRole(role))
	return m
}

// AppendItem appends an already-constructed MenuItem.
func (m *Menu) AppendItem(item *MenuItem) {
	m.Items = append(m.Items, item)
}

// Clone returns a deep copy of the menu tree.
func (m *Menu) Clone() *Menu {
	cloned := &Menu{label: m.label}
	for _, item := range m.Items {
		cloned.Items = append(cloned.Items, item.Clone())
	}
	return cloned
}

// Destroy frees all items in the menu.
func (m *Menu) Destroy() {
	for _, item := range m.Items {
		item.Destroy()
	}
	m.Items = nil
}

func (m *Menu) setContextData(data *ContextMenuData) {
	for _, item := range m.Items {
		item.contextMenuData = data
		if item.submenu != nil {
			item.submenu.setContextData(data)
		}
	}
}

// MenuManager owns the application menu.
//
//	app.Menu.SetApplicationMenu(menu)
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
	attachedWindow Window
	onClick        func()
}

func (t *SystemTray) SetIcon(data []byte) *SystemTray {
	t.icon = append([]byte(nil), data...)
	return t
}

func (t *SystemTray) SetTemplateIcon(data []byte) *SystemTray {
	t.templateIcon = append([]byte(nil), data...)
	return t
}

func (t *SystemTray) SetTooltip(text string) { t.tooltip = text }
func (t *SystemTray) SetLabel(text string)   { t.label = text }
func (t *SystemTray) SetMenu(menu *Menu) *SystemTray {
	t.menu = menu
	return t
}

func (t *SystemTray) OnClick(callback func()) *SystemTray {
	t.onClick = callback
	return t
}

// AttachWindow associates a window with the tray icon (shown on click).
func (t *SystemTray) AttachWindow(w Window) *SystemTray {
	t.attachedWindow = w
	t.OnClick(func() {
		if t.attachedWindow == nil {
			return
		}
		if t.attachedWindow.IsVisible() {
			t.attachedWindow.Hide()
			return
		}
		t.attachedWindow.Show()
		t.attachedWindow.Focus()
	})
	return t
}

// Click triggers the registered tray click handler.
func (t *SystemTray) Click() {
	if t.onClick != nil {
		t.onClick()
	}
}

// SystemTrayManager creates tray instances.
type SystemTrayManager struct{}

func (m *SystemTrayManager) New() *SystemTray {
	tray := &SystemTray{}
	ensureTrayCompatState(tray)
	return tray
}

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

// DropTargetDetails mirrors the Wails drop-target payload.
type DropTargetDetails struct {
	X          int               `json:"x"`
	Y          int               `json:"y"`
	ElementID  string            `json:"id"`
	ClassList  []string          `json:"classList"`
	Attributes map[string]string `json:"attributes,omitempty"`
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

// WebviewWindow is a lightweight, in-memory window implementation
// that satisfies the Window interface.
type WebviewWindow struct {
	mu                sync.RWMutex
	opts              WebviewWindowOptions
	windowID          uint
	title             string
	posX, posY        int
	sizeW, sizeH      int
	maximised         bool
	focused           bool
	visible           bool
	alwaysOnTop       bool
	isFullscreen      bool
	closed            bool
	zoom              float64
	resizable         bool
	ignoreMouseEvents bool
	enabled           bool
	eventHandlers     map[events.WindowEventType][]func(*WindowEvent)
}

var globalWindowID uint
var globalWindowIDMu sync.Mutex

func nextWindowID() uint {
	globalWindowIDMu.Lock()
	defer globalWindowIDMu.Unlock()
	globalWindowID++
	return globalWindowID
}

func newWebviewWindow(options WebviewWindowOptions) *WebviewWindow {
	return &WebviewWindow{
		opts:          options,
		windowID:      nextWindowID(),
		title:         options.Title,
		posX:          options.X,
		posY:          options.Y,
		sizeW:         options.Width,
		sizeH:         options.Height,
		visible:       !options.Hidden,
		alwaysOnTop:   options.AlwaysOnTop,
		zoom:          options.Zoom,
		resizable:     !options.DisableResize,
		enabled:       true,
		eventHandlers: make(map[events.WindowEventType][]func(*WindowEvent)),
	}
}

// ID returns the unique numeric identifier for the window.
func (w *WebviewWindow) ID() uint { return w.windowID }

// Name returns the window name set in WebviewWindowOptions.
func (w *WebviewWindow) Name() string { return w.opts.Name }

// Show makes the window visible and returns it for chaining.
func (w *WebviewWindow) Show() Window {
	w.mu.Lock()
	w.visible = true
	w.mu.Unlock()
	return w
}

// Hide makes the window invisible and returns it for chaining.
func (w *WebviewWindow) Hide() Window {
	w.mu.Lock()
	w.visible = false
	w.mu.Unlock()
	return w
}

// IsVisible reports whether the window is currently visible.
func (w *WebviewWindow) IsVisible() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.visible
}

// Close marks the window as closed.
func (w *WebviewWindow) Close() {
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()
}

// Focus marks the window as focused.
func (w *WebviewWindow) Focus() {
	w.mu.Lock()
	w.focused = true
	w.mu.Unlock()
}

// Run is a no-op in the stub (the real implementation enters the run loop).
func (w *WebviewWindow) Run() {}

// Center is a no-op in the stub.
func (w *WebviewWindow) Center() {}

// Position returns the current x/y position.
func (w *WebviewWindow) Position() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.posX, w.posY
}

// RelativePosition returns the position relative to the screen.
func (w *WebviewWindow) RelativePosition() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.posX, w.posY
}

// Size returns the current width and height.
func (w *WebviewWindow) Size() (int, int) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.sizeW, w.sizeH
}

// Width returns the current window width.
func (w *WebviewWindow) Width() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.sizeW
}

// Height returns the current window height.
func (w *WebviewWindow) Height() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.sizeH
}

// Bounds returns the window's position and size as a Rect.
func (w *WebviewWindow) Bounds() Rect {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return Rect{X: w.posX, Y: w.posY, Width: w.sizeW, Height: w.sizeH}
}

// SetPosition sets the top-left corner position.
func (w *WebviewWindow) SetPosition(x, y int) {
	w.mu.Lock()
	w.posX, w.posY = x, y
	w.mu.Unlock()
}

// SetRelativePosition sets position relative to the screen and returns the window.
func (w *WebviewWindow) SetRelativePosition(x, y int) Window {
	w.SetPosition(x, y)
	return w
}

// SetSize sets the window dimensions and returns the window.
func (w *WebviewWindow) SetSize(width, height int) Window {
	w.mu.Lock()
	w.sizeW, w.sizeH = width, height
	w.mu.Unlock()
	return w
}

// SetBounds sets position and size in one call.
func (w *WebviewWindow) SetBounds(bounds Rect) {
	w.mu.Lock()
	w.posX, w.posY = bounds.X, bounds.Y
	w.sizeW, w.sizeH = bounds.Width, bounds.Height
	w.mu.Unlock()
}

// SetMaxSize is a no-op in the stub.
func (w *WebviewWindow) SetMaxSize(maxWidth, maxHeight int) Window { return w }

// SetMinSize is a no-op in the stub.
func (w *WebviewWindow) SetMinSize(minWidth, minHeight int) Window { return w }

// EnableSizeConstraints is a no-op in the stub.
func (w *WebviewWindow) EnableSizeConstraints() {}

// DisableSizeConstraints is a no-op in the stub.
func (w *WebviewWindow) DisableSizeConstraints() {}

// Resizable reports whether the user can resize the window.
func (w *WebviewWindow) Resizable() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.resizable
}

// SetResizable enables or disables user resizing and returns the window.
func (w *WebviewWindow) SetResizable(b bool) Window {
	w.mu.Lock()
	w.resizable = b
	w.mu.Unlock()
	return w
}

// Maximise maximises the window and returns it.
func (w *WebviewWindow) Maximise() Window {
	w.mu.Lock()
	w.maximised = true
	w.mu.Unlock()
	return w
}

// UnMaximise restores from maximised state.
func (w *WebviewWindow) UnMaximise() {
	w.mu.Lock()
	w.maximised = false
	w.mu.Unlock()
}

// ToggleMaximise toggles between maximised and normal.
func (w *WebviewWindow) ToggleMaximise() {
	w.mu.Lock()
	w.maximised = !w.maximised
	w.mu.Unlock()
}

// IsMaximised reports whether the window is maximised.
func (w *WebviewWindow) IsMaximised() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.maximised
}

// Minimise minimises the window and returns it.
func (w *WebviewWindow) Minimise() Window {
	w.mu.Lock()
	w.visible = false
	w.mu.Unlock()
	return w
}

// UnMinimise restores from minimised state.
func (w *WebviewWindow) UnMinimise() {
	w.mu.Lock()
	w.visible = true
	w.mu.Unlock()
}

// IsMinimised always returns false in the stub.
func (w *WebviewWindow) IsMinimised() bool { return false }

// Fullscreen enters fullscreen and returns the window.
func (w *WebviewWindow) Fullscreen() Window {
	w.mu.Lock()
	w.isFullscreen = true
	w.mu.Unlock()
	return w
}

// UnFullscreen exits fullscreen.
func (w *WebviewWindow) UnFullscreen() {
	w.mu.Lock()
	w.isFullscreen = false
	w.mu.Unlock()
}

// ToggleFullscreen toggles between fullscreen and normal.
func (w *WebviewWindow) ToggleFullscreen() {
	w.mu.Lock()
	w.isFullscreen = !w.isFullscreen
	w.mu.Unlock()
}

// IsFullscreen reports whether the window is in fullscreen mode.
func (w *WebviewWindow) IsFullscreen() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isFullscreen
}

// Restore exits both fullscreen and maximised states.
func (w *WebviewWindow) Restore() {
	w.mu.Lock()
	w.maximised = false
	w.isFullscreen = false
	w.mu.Unlock()
}

// SnapAssist is a no-op in the stub.
func (w *WebviewWindow) SnapAssist() {}

// SetTitle updates the window title and returns the window.
func (w *WebviewWindow) SetTitle(title string) Window {
	w.mu.Lock()
	w.title = title
	w.mu.Unlock()
	return w
}

// SetURL is a no-op in the stub.
func (w *WebviewWindow) SetURL(s string) Window { return w }

// SetHTML is a no-op in the stub.
func (w *WebviewWindow) SetHTML(html string) Window { return w }

// SetMinimiseButtonState is a no-op in the stub.
func (w *WebviewWindow) SetMinimiseButtonState(state ButtonState) Window { return w }

// SetMaximiseButtonState is a no-op in the stub.
func (w *WebviewWindow) SetMaximiseButtonState(state ButtonState) Window { return w }

// SetCloseButtonState is a no-op in the stub.
func (w *WebviewWindow) SetCloseButtonState(state ButtonState) Window { return w }

// SetMenu is a no-op in the stub.
func (w *WebviewWindow) SetMenu(menu *Menu) {}

// ShowMenuBar is a no-op in the stub.
func (w *WebviewWindow) ShowMenuBar() {}

// HideMenuBar is a no-op in the stub.
func (w *WebviewWindow) HideMenuBar() {}

// ToggleMenuBar is a no-op in the stub.
func (w *WebviewWindow) ToggleMenuBar() {}

// SetBackgroundColour is a no-op in the stub.
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window { return w }

// SetAlwaysOnTop sets the always-on-top flag and returns the window.
func (w *WebviewWindow) SetAlwaysOnTop(b bool) Window {
	w.mu.Lock()
	w.alwaysOnTop = b
	w.mu.Unlock()
	return w
}

// SetFrameless is a no-op in the stub.
func (w *WebviewWindow) SetFrameless(frameless bool) Window { return w }

// ToggleFrameless is a no-op in the stub.
func (w *WebviewWindow) ToggleFrameless() {}

// SetIgnoreMouseEvents sets the mouse-event passthrough flag and returns the window.
func (w *WebviewWindow) SetIgnoreMouseEvents(ignore bool) Window {
	w.mu.Lock()
	w.ignoreMouseEvents = ignore
	w.mu.Unlock()
	return w
}

// IsIgnoreMouseEvents reports whether mouse events are being ignored.
func (w *WebviewWindow) IsIgnoreMouseEvents() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.ignoreMouseEvents
}

// SetContentProtection is a no-op in the stub.
func (w *WebviewWindow) SetContentProtection(protection bool) Window { return w }

// GetZoom returns the current zoom magnification.
func (w *WebviewWindow) GetZoom() float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.zoom
}

// SetZoom sets the zoom magnification and returns the window.
func (w *WebviewWindow) SetZoom(magnification float64) Window {
	w.mu.Lock()
	w.zoom = magnification
	w.mu.Unlock()
	return w
}

// Zoom is a no-op in the stub.
func (w *WebviewWindow) Zoom() {}

// ZoomIn is a no-op in the stub.
func (w *WebviewWindow) ZoomIn() {}

// ZoomOut is a no-op in the stub.
func (w *WebviewWindow) ZoomOut() {}

// ZoomReset resets zoom to 1.0 and returns the window.
func (w *WebviewWindow) ZoomReset() Window {
	w.mu.Lock()
	w.zoom = 1.0
	w.mu.Unlock()
	return w
}

// GetBorderSizes returns zero insets in the stub.
func (w *WebviewWindow) GetBorderSizes() *LRTB { return &LRTB{} }

// GetScreen returns nil in the stub.
func (w *WebviewWindow) GetScreen() (*Screen, error) { return nil, nil }

// ExecJS is a no-op in the stub.
func (w *WebviewWindow) ExecJS(js string) {}

// EmitEvent always returns false in the stub.
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool { return false }

// DispatchWailsEvent is a no-op in the stub.
func (w *WebviewWindow) DispatchWailsEvent(event *CustomEvent) {}

// OnWindowEvent registers an event callback and returns an unsubscribe function.
func (w *WebviewWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	w.mu.Lock()
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], callback)
	w.mu.Unlock()
	return func() {}
}

// RegisterHook is an alias for OnWindowEvent.
func (w *WebviewWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return w.OnWindowEvent(eventType, callback)
}

// handleDragAndDropMessage is a no-op in the stub.
func (w *WebviewWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {}

// InitiateFrontendDropProcessing is a no-op in the stub.
func (w *WebviewWindow) InitiateFrontendDropProcessing(filenames []string, x int, y int) {}

// HandleMessage is a no-op in the stub.
func (w *WebviewWindow) HandleMessage(message string) {}

// HandleWindowEvent is a no-op in the stub.
func (w *WebviewWindow) HandleWindowEvent(id uint) {}

// HandleKeyEvent is a no-op in the stub.
func (w *WebviewWindow) HandleKeyEvent(acceleratorString string) {
	state := ensureWebviewCompatState(w)
	state.mu.RLock()
	callback := state.keyBindings[acceleratorString]
	state.mu.RUnlock()
	if callback != nil {
		callback(w)
	}
}

// OpenContextMenu is a no-op in the stub.
func (w *WebviewWindow) OpenContextMenu(data *ContextMenuData) {}

// AttachModal is a no-op in the stub.
func (w *WebviewWindow) AttachModal(modalWindow Window) {}

// OpenDevTools marks devtools as open in the stub.
func (w *WebviewWindow) OpenDevTools() {
	state := ensureWebviewCompatState(w)
	state.mu.Lock()
	state.devtoolsOpen = true
	state.mu.Unlock()
}

// Print always returns nil in the stub.
func (w *WebviewWindow) Print() error { return nil }

// Flash is a no-op in the stub.
func (w *WebviewWindow) Flash(enabled bool) {}

// IsFocused reports whether the window has focus.
func (w *WebviewWindow) IsFocused() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.focused
}

// NativeWindow returns nil in the stub.
func (w *WebviewWindow) NativeWindow() unsafe.Pointer { return nil }

// SetEnabled sets the window's enabled state.
func (w *WebviewWindow) SetEnabled(enabled bool) {
	w.mu.Lock()
	w.enabled = enabled
	w.mu.Unlock()
}

// Reload is a no-op in the stub.
func (w *WebviewWindow) Reload() {}

// ForceReload is a no-op in the stub.
func (w *WebviewWindow) ForceReload() {}

// Info is a no-op in the stub.
func (w *WebviewWindow) Info(message string, args ...any) {}

// Error is a no-op in the stub.
func (w *WebviewWindow) Error(message string, args ...any) {}

// shouldUnconditionallyClose always returns false in the stub.
func (w *WebviewWindow) shouldUnconditionallyClose() bool { return false }

// Internal editing stubs — satisfy the Window interface.
func (w *WebviewWindow) cut()       {}
func (w *WebviewWindow) copy()      {}
func (w *WebviewWindow) paste()     {}
func (w *WebviewWindow) undo()      {}
func (w *WebviewWindow) redo()      {}
func (w *WebviewWindow) delete()    {}
func (w *WebviewWindow) selectAll() {}

// Title returns the current window title.
func (w *WebviewWindow) Title() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.title
}

// WindowManager manages in-memory windows.
//
//	win := app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Main"})
type WindowManager struct {
	mu        sync.RWMutex
	windows   map[uint]Window
	callbacks []func(Window)
	current   uint
}

func (wm *WindowManager) init() {
	if wm.windows == nil {
		wm.windows = make(map[uint]Window)
	}
}

// NewWithOptions creates and registers a new window.
func (wm *WindowManager) NewWithOptions(options WebviewWindowOptions) *WebviewWindow {
	window := newWebviewWindow(options)
	wm.mu.Lock()
	wm.init()
	wm.windows[window.windowID] = window
	wm.current = window.windowID
	callbacks := append([]func(Window){}, wm.callbacks...)
	wm.mu.Unlock()
	for _, callback := range callbacks {
		if callback != nil {
			callback(window)
		}
	}
	return window
}

// New creates a window with default options.
func (wm *WindowManager) New() *WebviewWindow {
	return wm.NewWithOptions(WebviewWindowOptions{})
}

// GetAll returns all managed windows as the Window interface slice.
func (wm *WindowManager) GetAll() []Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	out := make([]Window, 0, len(wm.windows))
	for _, window := range wm.windows {
		out = append(out, window)
	}
	return out
}

// GetByName finds a window by name, returning it and whether it was found.
func (wm *WindowManager) GetByName(name string) (Window, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	for _, window := range wm.windows {
		if window != nil && window.Name() == name {
			return window, true
		}
	}
	return nil, false
}

// GetByID finds a window by its numeric ID.
func (wm *WindowManager) GetByID(id uint) (Window, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	window, exists := wm.windows[id]
	return window, exists
}

// Remove unregisters a window by ID.
func (wm *WindowManager) Remove(windowID uint) {
	wm.mu.Lock()
	wm.init()
	delete(wm.windows, windowID)
	if wm.current == windowID {
		wm.current = 0
	}
	wm.mu.Unlock()
}

// Get aliases GetByName.
func (wm *WindowManager) Get(name string) (Window, bool) {
	return wm.GetByName(name)
}

// OnCreate registers a callback for future window creation.
func (wm *WindowManager) OnCreate(callback func(Window)) {
	if callback == nil {
		return
	}
	wm.mu.Lock()
	wm.callbacks = append(wm.callbacks, callback)
	wm.mu.Unlock()
}

// Current returns the last-added window.
func (wm *WindowManager) Current() Window {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.windows[wm.current]
}

// Add registers an externally created window.
func (wm *WindowManager) Add(window Window) {
	if window == nil {
		return
	}
	wm.mu.Lock()
	wm.init()
	wm.windows[window.ID()] = window
	wm.current = window.ID()
	callbacks := append([]func(Window){}, wm.callbacks...)
	wm.mu.Unlock()
	for _, callback := range callbacks {
		if callback != nil {
			callback(window)
		}
	}
}

// RemoveByName removes a window by name.
func (wm *WindowManager) RemoveByName(name string) bool {
	window, ok := wm.GetByName(name)
	if !ok {
		return false
	}
	wm.Remove(window.ID())
	return true
}

// App is the top-level application object used by the GUI packages.
//
//	app := &application.App{}
//	win := app.Window.NewWithOptions(application.WebviewWindowOptions{Title: "Main"})
type App struct {
	Logger      *slog.Logger
	Window      *WindowManager
	Menu        *MenuManager
	SystemTray  *SystemTrayManager
	Dialog      *DialogManager
	Event       *EventManager
	Browser     *BrowserManager
	Clipboard   *ClipboardManager
	ContextMenu *ContextMenuManager
	Env         *EnvironmentManager
	Environment *EnvironmentManager
	Screen      *ScreenManager
	KeyBinding  *KeyBindingManager
}

// NewApp creates a zero-config in-memory application stub.
func NewApp() *App { return newStubApp(Options{}) }

// Quit executes registered shutdown callbacks.
func (a *App) Quit() {
	state := ensureAppCompatState(a)
	state.mu.Lock()
	callbacks := append([]func(){}, state.shutdown...)
	state.running = false
	state.mu.Unlock()
	for _, callback := range callbacks {
		if callback != nil {
			callback()
		}
	}
}

// NewMenu creates a new empty Menu.
func (a *App) NewMenu() *Menu {
	return NewMenu()
}
