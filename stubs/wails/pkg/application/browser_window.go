package application

import (
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// ButtonState represents the visual state of a window control button.
//
//	window.SetMinimiseButtonState(ButtonHidden)
type ButtonState int

const (
	ButtonEnabled  ButtonState = 0
	ButtonDisabled ButtonState = 1
	ButtonHidden   ButtonState = 2
)

// LRTB holds left/right/top/bottom border sizes in pixels.
//
//	sizes := window.GetBorderSizes()
type LRTB struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

// ContextMenuData carries context-menu position and metadata from the frontend.
//
//	window.OpenContextMenu(&ContextMenuData{Id: "file-menu"})
type ContextMenuData struct {
	Id   string `json:"id"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Data string `json:"data"`
}

func (c *ContextMenuData) clone() *ContextMenuData {
	if c == nil {
		return nil
	}
	copy := *c
	return &copy
}

// CustomEvent carries a named event with arbitrary data from the frontend.
//
//	window.DispatchWailsEvent(&CustomEvent{Name: "ready", Data: nil})
type CustomEvent struct {
	Name      string `json:"name"`
	Data      any    `json:"data"`
	Sender    string `json:"sender,omitempty"`
	cancelled bool
}

// Cancel prevents the event from reaching further listeners.
func (e *CustomEvent) Cancel() { e.cancelled = true }

// IsCancelled reports whether Cancel has been called.
func (e *CustomEvent) IsCancelled() bool { return e.cancelled }

// BrowserWindow represents a browser client connection in server mode.
// It satisfies the Window interface so browser clients are treated
// uniformly with native windows throughout the codebase.
//
//	bw := NewBrowserWindow(1, "nano-abc123")
//	bw.Focus() // no-op in browser mode
type BrowserWindow struct {
	mu       sync.RWMutex
	id       uint
	name     string
	clientID string
}

// NewBrowserWindow constructs a BrowserWindow with the given numeric ID and client nano-ID.
//
//	bw := NewBrowserWindow(1, "nano-abc123")
func NewBrowserWindow(id uint, clientID string) *BrowserWindow {
	return &BrowserWindow{
		id:       id,
		name:     "browser-window",
		clientID: clientID,
	}
}

// ID returns the numeric window identifier.
func (b *BrowserWindow) ID() uint { return b.id }

// Name returns the window name.
func (b *BrowserWindow) Name() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.name
}

// ClientID returns the runtime nano-ID for this client.
func (b *BrowserWindow) ClientID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.clientID
}

// No-op windowing methods — browser windows have no native chrome.

func (b *BrowserWindow) Center()                  {}
func (b *BrowserWindow) Close()                   {}
func (b *BrowserWindow) DisableSizeConstraints()  {}
func (b *BrowserWindow) EnableSizeConstraints()   {}
func (b *BrowserWindow) ExecJS(_ string)          {}
func (b *BrowserWindow) Focus()                   {}
func (b *BrowserWindow) ForceReload()             {}
func (b *BrowserWindow) HideMenuBar()             {}
func (b *BrowserWindow) OpenDevTools()            {}
func (b *BrowserWindow) Reload()                  {}
func (b *BrowserWindow) Restore()                 {}
func (b *BrowserWindow) Run()                     {}
func (b *BrowserWindow) SetPosition(_ int, _ int) {}
func (b *BrowserWindow) ShowMenuBar()             {}
func (b *BrowserWindow) SnapAssist()              {}
func (b *BrowserWindow) ToggleFrameless()         {}
func (b *BrowserWindow) ToggleFullscreen()        {}
func (b *BrowserWindow) ToggleMaximise()          {}
func (b *BrowserWindow) ToggleMenuBar()           {}
func (b *BrowserWindow) UnFullscreen()            {}
func (b *BrowserWindow) UnMaximise()              {}
func (b *BrowserWindow) UnMinimise()              {}
func (b *BrowserWindow) SetEnabled(_ bool)        {}
func (b *BrowserWindow) Flash(_ bool)             {}
func (b *BrowserWindow) SetMenu(_ *Menu)          {}
func (b *BrowserWindow) SetBounds(_ Rect)         {}
func (b *BrowserWindow) Zoom()                    {}
func (b *BrowserWindow) ZoomIn()                  {}
func (b *BrowserWindow) ZoomOut()                 {}
func (b *BrowserWindow) OpenContextMenu(_ *ContextMenuData) {}
func (b *BrowserWindow) HandleMessage(_ string)             {}
func (b *BrowserWindow) HandleWindowEvent(_ uint)           {}
func (b *BrowserWindow) HandleKeyEvent(_ string)            {}
func (b *BrowserWindow) AttachModal(_ Window)               {}

// Internal editing stubs.
func (b *BrowserWindow) cut()       {}
func (b *BrowserWindow) copy()      {}
func (b *BrowserWindow) paste()     {}
func (b *BrowserWindow) undo()      {}
func (b *BrowserWindow) redo()      {}
func (b *BrowserWindow) delete()    {}
func (b *BrowserWindow) selectAll() {}

// shouldUnconditionallyClose always returns false for browser windows.
func (b *BrowserWindow) shouldUnconditionallyClose() bool { return false }

// Methods returning Window for chaining — always no-op for browser windows.

func (b *BrowserWindow) Fullscreen() Window                           { return b }
func (b *BrowserWindow) Hide() Window                                 { return b }
func (b *BrowserWindow) Maximise() Window                             { return b }
func (b *BrowserWindow) Minimise() Window                             { return b }
func (b *BrowserWindow) Show() Window                                 { return b }
func (b *BrowserWindow) SetAlwaysOnTop(_ bool) Window                 { return b }
func (b *BrowserWindow) SetBackgroundColour(_ RGBA) Window            { return b }
func (b *BrowserWindow) SetFrameless(_ bool) Window                   { return b }
func (b *BrowserWindow) SetHTML(_ string) Window                      { return b }
func (b *BrowserWindow) SetMinimiseButtonState(_ ButtonState) Window  { return b }
func (b *BrowserWindow) SetMaximiseButtonState(_ ButtonState) Window  { return b }
func (b *BrowserWindow) SetCloseButtonState(_ ButtonState) Window     { return b }
func (b *BrowserWindow) SetMaxSize(_ int, _ int) Window               { return b }
func (b *BrowserWindow) SetMinSize(_ int, _ int) Window               { return b }
func (b *BrowserWindow) SetRelativePosition(_ int, _ int) Window      { return b }
func (b *BrowserWindow) SetResizable(_ bool) Window                   { return b }
func (b *BrowserWindow) SetIgnoreMouseEvents(_ bool) Window           { return b }
func (b *BrowserWindow) SetSize(_ int, _ int) Window                  { return b }
func (b *BrowserWindow) SetTitle(_ string) Window                     { return b }
func (b *BrowserWindow) SetURL(_ string) Window                       { return b }
func (b *BrowserWindow) SetZoom(_ float64) Window                     { return b }
func (b *BrowserWindow) SetContentProtection(_ bool) Window           { return b }
func (b *BrowserWindow) ZoomReset() Window                            { return b }

// Methods returning simple zero values.

func (b *BrowserWindow) GetBorderSizes() *LRTB               { return nil }
func (b *BrowserWindow) GetScreen() (*Screen, error)         { return nil, nil }
func (b *BrowserWindow) GetZoom() float64                    { return 1.0 }
func (b *BrowserWindow) Height() int                         { return 0 }
func (b *BrowserWindow) Width() int                          { return 0 }
func (b *BrowserWindow) IsFocused() bool                     { return false }
func (b *BrowserWindow) IsFullscreen() bool                  { return false }
func (b *BrowserWindow) IsIgnoreMouseEvents() bool           { return false }
func (b *BrowserWindow) IsMaximised() bool                   { return false }
func (b *BrowserWindow) IsMinimised() bool                   { return false }
func (b *BrowserWindow) IsVisible() bool                     { return true }
func (b *BrowserWindow) Resizable() bool                     { return false }
func (b *BrowserWindow) Position() (int, int)                { return 0, 0 }
func (b *BrowserWindow) RelativePosition() (int, int)        { return 0, 0 }
func (b *BrowserWindow) Size() (int, int)                    { return 0, 0 }
func (b *BrowserWindow) Bounds() Rect                        { return Rect{} }
func (b *BrowserWindow) NativeWindow() unsafe.Pointer        { return nil }
func (b *BrowserWindow) Print() error                        { return nil }

// DispatchWailsEvent is a no-op for browser windows; events are broadcast via WebSocket.
func (b *BrowserWindow) DispatchWailsEvent(_ *CustomEvent) {}

// EmitEvent broadcasts a named event; always returns false in the stub.
func (b *BrowserWindow) EmitEvent(_ string, _ ...any) bool { return false }

// Error logs an error message (no-op in the stub).
func (b *BrowserWindow) Error(_ string, _ ...any) {}

// Info logs an info message (no-op in the stub).
func (b *BrowserWindow) Info(_ string, _ ...any) {}

// OnWindowEvent registers a callback for a window event type; returns an unsubscribe func.
//
//	unsubscribe := bw.OnWindowEvent(events.Common.WindowClosing, fn)
func (b *BrowserWindow) OnWindowEvent(_ events.WindowEventType, _ func(*WindowEvent)) func() {
	return func() {}
}

// RegisterHook registers a lifecycle hook; returns an unsubscribe func.
func (b *BrowserWindow) RegisterHook(_ events.WindowEventType, _ func(*WindowEvent)) func() {
	return func() {}
}

// handleDragAndDropMessage is a no-op for browser windows.
func (b *BrowserWindow) handleDragAndDropMessage(_ []string, _ *DropTargetDetails) {}

// InitiateFrontendDropProcessing is a no-op for browser windows.
func (b *BrowserWindow) InitiateFrontendDropProcessing(_ []string, _ int, _ int) {}
