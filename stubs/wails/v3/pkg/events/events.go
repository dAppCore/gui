package events

// WindowEventType identifies a window event.
type WindowEventType string

// Common exposes the event names used by the repo.
var Common = struct {
	WindowFocus        WindowEventType
	WindowLostFocus    WindowEventType
	WindowDidMove      WindowEventType
	WindowDidResize    WindowEventType
	WindowClosing      WindowEventType
	WindowFilesDropped WindowEventType
}{
	WindowFocus:        "focus",
	WindowLostFocus:    "blur",
	WindowDidMove:      "move",
	WindowDidResize:    "resize",
	WindowClosing:      "close",
	WindowFilesDropped: "files-dropped",
}
