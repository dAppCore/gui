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
	mu       sync.RWMutex
	id       uint
	name     string
	clientID string
}

// NewBrowserWindow creates a browser window with the given ID and client ID.
//
//	browserWindow := application.NewBrowserWindow(1, "nanoid-abc123")
func NewBrowserWindow(id uint, clientID string) *BrowserWindow {
	return &BrowserWindow{
		id:       id,
		name:     fmt.Sprintf("browser-%d", id),
		clientID: clientID,
	}
}

func (browserWindow *BrowserWindow) ID() uint         { return browserWindow.id }
func (browserWindow *BrowserWindow) Name() string     { return browserWindow.name }
func (browserWindow *BrowserWindow) ClientID() string { return browserWindow.clientID }

func (browserWindow *BrowserWindow) DispatchWailsEvent(event *CustomEvent) {}
func (browserWindow *BrowserWindow) EmitEvent(name string, data ...any) bool {
	return true
}

func (browserWindow *BrowserWindow) Error(message string, arguments ...any) {}
func (browserWindow *BrowserWindow) Info(message string, arguments ...any)  {}

// No-op methods — browser windows are controlled via WebSocket, not native APIs.
func (browserWindow *BrowserWindow) Center()                  {}
func (browserWindow *BrowserWindow) Close()                   {}
func (browserWindow *BrowserWindow) DisableSizeConstraints()  {}
func (browserWindow *BrowserWindow) EnableSizeConstraints()   {}
func (browserWindow *BrowserWindow) ExecJS(javascript string) {}
func (browserWindow *BrowserWindow) Focus()                   {}
func (browserWindow *BrowserWindow) ForceReload()             {}
func (browserWindow *BrowserWindow) Fullscreen() Window       { return browserWindow }
func (browserWindow *BrowserWindow) GetBorderSizes() *LRTB    { return nil }
func (browserWindow *BrowserWindow) GetScreen() (*Screen, error) {
	return nil, nil
}
func (browserWindow *BrowserWindow) GetZoom() float64 { return 1.0 }
func (browserWindow *BrowserWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {
}
func (browserWindow *BrowserWindow) HandleMessage(message string)      {}
func (browserWindow *BrowserWindow) HandleWindowEvent(identifier uint) {}
func (browserWindow *BrowserWindow) Height() int                       { return 0 }
func (browserWindow *BrowserWindow) Hide() Window                      { return browserWindow }
func (browserWindow *BrowserWindow) HideMenuBar()                      {}
func (browserWindow *BrowserWindow) IsFocused() bool                   { return false }
func (browserWindow *BrowserWindow) IsFullscreen() bool                { return false }
func (browserWindow *BrowserWindow) IsIgnoreMouseEvents() bool         { return false }
func (browserWindow *BrowserWindow) IsMaximised() bool                 { return false }
func (browserWindow *BrowserWindow) IsMinimised() bool                 { return false }
func (browserWindow *BrowserWindow) HandleKeyEvent(accelerator string) {}
func (browserWindow *BrowserWindow) Maximise() Window                  { return browserWindow }
func (browserWindow *BrowserWindow) Minimise() Window                  { return browserWindow }
func (browserWindow *BrowserWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}
func (browserWindow *BrowserWindow) OpenContextMenu(data *ContextMenuData)  {}
func (browserWindow *BrowserWindow) Position() (int, int)                   { return 0, 0 }
func (browserWindow *BrowserWindow) RelativePosition() (int, int)           { return 0, 0 }
func (browserWindow *BrowserWindow) Reload()                                {}
func (browserWindow *BrowserWindow) Resizable() bool                        { return false }
func (browserWindow *BrowserWindow) Restore()                               {}
func (browserWindow *BrowserWindow) Run()                                   {}
func (browserWindow *BrowserWindow) SetPosition(x, y int)                   {}
func (browserWindow *BrowserWindow) SetAlwaysOnTop(alwaysOnTop bool) Window { return browserWindow }
func (browserWindow *BrowserWindow) SetBackgroundColour(colour RGBA) Window { return browserWindow }
func (browserWindow *BrowserWindow) SetFrameless(frameless bool) Window     { return browserWindow }
func (browserWindow *BrowserWindow) SetHTML(html string) Window             { return browserWindow }
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
	return browserWindow
}
func (browserWindow *BrowserWindow) SetMinSize(minWidth, minHeight int) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetRelativePosition(x, y int) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetResizable(resizable bool) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetIgnoreMouseEvents(ignore bool) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) SetSize(width, height int) Window { return browserWindow }
func (browserWindow *BrowserWindow) SetTitle(title string) Window     { return browserWindow }
func (browserWindow *BrowserWindow) SetURL(url string) Window         { return browserWindow }
func (browserWindow *BrowserWindow) SetZoom(magnification float64) Window {
	return browserWindow
}
func (browserWindow *BrowserWindow) Show() Window          { return browserWindow }
func (browserWindow *BrowserWindow) ShowMenuBar()          {}
func (browserWindow *BrowserWindow) Size() (int, int)      { return 0, 0 }
func (browserWindow *BrowserWindow) OpenDevTools()         {}
func (browserWindow *BrowserWindow) ToggleFullscreen()     {}
func (browserWindow *BrowserWindow) ToggleMaximise()       {}
func (browserWindow *BrowserWindow) ToggleMenuBar()        {}
func (browserWindow *BrowserWindow) ToggleFrameless()      {}
func (browserWindow *BrowserWindow) UnFullscreen()         {}
func (browserWindow *BrowserWindow) UnMaximise()           {}
func (browserWindow *BrowserWindow) UnMinimise()           {}
func (browserWindow *BrowserWindow) Width() int            { return 0 }
func (browserWindow *BrowserWindow) IsVisible() bool       { return true }
func (browserWindow *BrowserWindow) Bounds() Rect          { return Rect{} }
func (browserWindow *BrowserWindow) SetBounds(bounds Rect) {}
func (browserWindow *BrowserWindow) Zoom()                 {}
func (browserWindow *BrowserWindow) ZoomIn()               {}
func (browserWindow *BrowserWindow) ZoomOut()              {}
func (browserWindow *BrowserWindow) ZoomReset() Window     { return browserWindow }
func (browserWindow *BrowserWindow) SetMenu(menu *Menu)    {}
func (browserWindow *BrowserWindow) SnapAssist()           {}
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
