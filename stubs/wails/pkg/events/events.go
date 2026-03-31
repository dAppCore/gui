package events

// WindowEventType identifies a window event emitted by the application layer.
type WindowEventType int

const (
	WindowFocus WindowEventType = iota
	WindowLostFocus
	WindowDidMove
	WindowDidResize
	WindowClosing
	WindowFilesDropped
)

// Common matches the event namespace used by the real Wails package.
var Common = struct {
	WindowFocus        WindowEventType
	WindowLostFocus    WindowEventType
	WindowDidMove      WindowEventType
	WindowDidResize    WindowEventType
	WindowClosing      WindowEventType
	WindowFilesDropped WindowEventType
}{
	WindowFocus:        WindowFocus,
	WindowLostFocus:    WindowLostFocus,
	WindowDidMove:      WindowDidMove,
	WindowDidResize:    WindowDidResize,
	WindowClosing:      WindowClosing,
	WindowFilesDropped: WindowFilesDropped,
}
