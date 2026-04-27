package application

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// ButtonState controls window button appearance.
type ButtonState int

const (
	ButtonEnabled  ButtonState = 0
	ButtonDisabled ButtonState = 1
	ButtonHidden   ButtonState = 2
)

// LRTB represents Left, Right, Top, Bottom border sizes.
type LRTB struct {
	Left, Right, Top, Bottom int
}

// ContextMenuData carries context menu trigger details.
type ContextMenuData struct {
	Name string
	X, Y int
	Data any
}

// BrowserWindow represents a browser client connection in server mode.
// Implements the Window interface — most methods are no-ops since browser
// clients are controlled via WebSocket, not native APIs.
//
//	browserWindow := application.NewBrowserWindow(1, "client-abc123")
type BrowserWindow struct {
	mu                sync.RWMutex
	id                uint
	name              string
	clientID          string
	visible           bool
	focused           bool
	maximised         bool
	minimised         bool
	fullscreen        bool
	resizable         bool
	ignoreMouseEvents bool
	alwaysOnTop       bool
	frameless         bool
	title             string
	url               string
	html              string
	x                 int
	y                 int
	width             int
	height            int
	minWidth          int
	minHeight         int
	maxWidth          int
	maxHeight         int
	zoom              float64
	backgroundColour  RGBA
}

// NewBrowserWindow creates a browser window with the given ID and client ID.
//
//	browserWindow := application.NewBrowserWindow(1, "nanoid-abc123")
func NewBrowserWindow(id uint, clientID string) *BrowserWindow {
	return &BrowserWindow{
		id:       id,
		name:     fmt.Sprintf("browser-%d", id),
		clientID: clientID,
		visible:  true,
		zoom:     1.0,
	}
}

func (browserWindow *BrowserWindow) ID() uint {
	if browserWindow == nil {
		return 0
	}
	return browserWindow.id
}

func (browserWindow *BrowserWindow) Name() string {
	if browserWindow == nil {
		return ""
	}
	return browserWindow.name
}

func (browserWindow *BrowserWindow) ClientID() string {
	if browserWindow == nil {
		return ""
	}
	return browserWindow.clientID
}

func (browserWindow *BrowserWindow) DispatchWailsEvent(event *CustomEvent) {}
func (browserWindow *BrowserWindow) EmitEvent(name string, data ...any) bool {
	return true
}

func (browserWindow *BrowserWindow) Error(message string, arguments ...any) {}
func (browserWindow *BrowserWindow) Info(message string, arguments ...any)  {}

// No-op methods — browser windows are controlled via WebSocket, not native APIs.
func (browserWindow *BrowserWindow) Center() {}
func (browserWindow *BrowserWindow) Close() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.visible = false
	browserWindow.focused = false
	browserWindow.minimised = false
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) DisableSizeConstraints()  {}
func (browserWindow *BrowserWindow) EnableSizeConstraints()   {}
func (browserWindow *BrowserWindow) ExecJS(javascript string) {}
func (browserWindow *BrowserWindow) Focus() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.focused = true
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) ForceReload() {}
func (browserWindow *BrowserWindow) Fullscreen() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.fullscreen = true
	browserWindow.maximised = false
	browserWindow.minimised = false
	browserWindow.visible = true
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) GetBorderSizes() *LRTB { return &LRTB{} }
func (browserWindow *BrowserWindow) GetScreen() (*Screen, error) {
	return &Screen{}, nil
}
func (browserWindow *BrowserWindow) GetZoom() float64 {
	if browserWindow == nil {
		return 1.0
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	if browserWindow.zoom == 0 {
		return 1.0
	}
	return browserWindow.zoom
}
func (browserWindow *BrowserWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {
}
func (browserWindow *BrowserWindow) HandleMessage(message string)      {}
func (browserWindow *BrowserWindow) HandleWindowEvent(identifier uint) {}
func (browserWindow *BrowserWindow) Height() int {
	if browserWindow == nil {
		return 0
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.height
}
func (browserWindow *BrowserWindow) Hide() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.visible = false
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) HideMenuBar() {}
func (browserWindow *BrowserWindow) IsFocused() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.focused
}
func (browserWindow *BrowserWindow) IsFullscreen() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.fullscreen
}
func (browserWindow *BrowserWindow) IsIgnoreMouseEvents() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.ignoreMouseEvents
}
func (browserWindow *BrowserWindow) IsMaximised() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.maximised
}
func (browserWindow *BrowserWindow) IsMinimised() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.minimised
}
func (browserWindow *BrowserWindow) HandleKeyEvent(accelerator string) {}
func (browserWindow *BrowserWindow) Maximise() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.maximised = true
	browserWindow.fullscreen = false
	browserWindow.minimised = false
	browserWindow.visible = true
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) Minimise() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.minimised = true
	browserWindow.maximised = false
	browserWindow.fullscreen = false
	browserWindow.visible = false
	browserWindow.focused = false
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}
func (browserWindow *BrowserWindow) OpenContextMenu(data *ContextMenuData) {}
func (browserWindow *BrowserWindow) Position() (int, int) {
	if browserWindow == nil {
		return 0, 0
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.x, browserWindow.y
}
func (browserWindow *BrowserWindow) RelativePosition() (int, int) {
	return browserWindow.Position()
}
func (browserWindow *BrowserWindow) Reload() {}
func (browserWindow *BrowserWindow) Resizable() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.resizable
}
func (browserWindow *BrowserWindow) Restore() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.maximised = false
	browserWindow.minimised = false
	browserWindow.fullscreen = false
	browserWindow.visible = true
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) Run() {}
func (browserWindow *BrowserWindow) SetPosition(x, y int) {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.x = x
	browserWindow.y = y
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) SetAlwaysOnTop(alwaysOnTop bool) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.alwaysOnTop = alwaysOnTop
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetBackgroundColour(colour RGBA) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.backgroundColour = colour
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetFrameless(frameless bool) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.frameless = frameless
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetHTML(html string) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.html = html
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMinimiseButtonState(state ButtonState) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMaximiseButtonState(state ButtonState) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetCloseButtonState(state ButtonState) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMaxSize(maxWidth, maxHeight int) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.maxWidth = maxWidth
	browserWindow.maxHeight = maxHeight
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMinSize(minWidth, minHeight int) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.minWidth = minWidth
	browserWindow.minHeight = minHeight
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetRelativePosition(x, y int) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.SetPosition(x, y)
	return browserWindow
}
func (browserWindow *BrowserWindow) SetResizable(resizable bool) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.resizable = resizable
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetIgnoreMouseEvents(ignore bool) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.ignoreMouseEvents = ignore
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetSize(width, height int) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.width = width
	browserWindow.height = height
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetTitle(title string) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.title = title
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetURL(url string) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.url = url
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetZoom(magnification float64) Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.zoom = magnification
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) Show() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.visible = true
	browserWindow.minimised = false
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) ShowMenuBar() {}
func (browserWindow *BrowserWindow) Size() (int, int) {
	if browserWindow == nil {
		return 0, 0
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.width, browserWindow.height
}
func (browserWindow *BrowserWindow) OpenDevTools() {}
func (browserWindow *BrowserWindow) ToggleFullscreen() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.fullscreen = !browserWindow.fullscreen
	if browserWindow.fullscreen {
		browserWindow.maximised = false
		browserWindow.minimised = false
		browserWindow.visible = true
	}
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) ToggleMaximise() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.maximised = !browserWindow.maximised
	if browserWindow.maximised {
		browserWindow.fullscreen = false
		browserWindow.minimised = false
		browserWindow.visible = true
	}
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) ToggleMenuBar()   {}
func (browserWindow *BrowserWindow) ToggleFrameless() {}
func (browserWindow *BrowserWindow) UnFullscreen() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.fullscreen = false
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) UnMaximise() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.maximised = false
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) UnMinimise() {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.minimised = false
	browserWindow.visible = true
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) Width() int {
	if browserWindow == nil {
		return 0
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.width
}
func (browserWindow *BrowserWindow) IsVisible() bool {
	if browserWindow == nil {
		return false
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return browserWindow.visible
}
func (browserWindow *BrowserWindow) Bounds() Rect {
	if browserWindow == nil {
		return Rect{}
	}
	browserWindow.mu.RLock()
	defer browserWindow.mu.RUnlock()
	return Rect{X: browserWindow.x, Y: browserWindow.y, Width: browserWindow.width, Height: browserWindow.height}
}
func (browserWindow *BrowserWindow) SetBounds(bounds Rect) {
	if browserWindow == nil {
		return
	}
	browserWindow.mu.Lock()
	browserWindow.x = bounds.X
	browserWindow.y = bounds.Y
	browserWindow.width = bounds.Width
	browserWindow.height = bounds.Height
	browserWindow.mu.Unlock()
}
func (browserWindow *BrowserWindow) Zoom()    {}
func (browserWindow *BrowserWindow) ZoomIn()  {}
func (browserWindow *BrowserWindow) ZoomOut() {}
func (browserWindow *BrowserWindow) ZoomReset() Window {
	if browserWindow == nil {
		return nil
	}
	browserWindow.mu.Lock()
	browserWindow.zoom = 1.0
	browserWindow.mu.Unlock()
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMenu(menu *Menu) {}
func (browserWindow *BrowserWindow) SnapAssist()        {}
func (browserWindow *BrowserWindow) SetContentProtection(protection bool) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetEnabled(enabled bool) {}
func (browserWindow *BrowserWindow) Flash(enabled bool)      {}
func (browserWindow *BrowserWindow) Print() error            { return nil }
func (browserWindow *BrowserWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}

// Internal platform hooks — no-ops for browser windows.

func (browserWindow *BrowserWindow) InitiateFrontendDropProcessing(filenames []string, x int, y int) {
}
func (browserWindow *BrowserWindow) shouldUnconditionallyClose() bool { return true }
func (browserWindow *BrowserWindow) cut()                             {}
func (browserWindow *BrowserWindow) copy()                            {}
func (browserWindow *BrowserWindow) paste()                           {}
func (browserWindow *BrowserWindow) undo()                            {}
func (browserWindow *BrowserWindow) redo()                            {}
func (browserWindow *BrowserWindow) delete()                          {}
func (browserWindow *BrowserWindow) selectAll()                       {}

// NativeWindow returns nil — browser windows have no native handle.
//
//	ptr := w.NativeWindow()
func (browserWindow *BrowserWindow) NativeWindow() unsafe.Pointer { return nil }

// AttachModal registers a modal window that blocks this window.
//
//	w.AttachModal(confirmDialog)
func (browserWindow *BrowserWindow) AttachModal(modalWindow Window) {}
