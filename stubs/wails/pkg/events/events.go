package events

// ApplicationEventType identifies a platform-level application event.
// Matches the type used by the real Wails v3 package.
//
//	em.OnApplicationEvent(events.Common.ThemeChanged, handler)
type ApplicationEventType uint

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
	ApplicationOpenedWithFile  ApplicationEventType
	ApplicationStarted         ApplicationEventType
	ApplicationLaunchedWithUrl ApplicationEventType
	ThemeChanged               ApplicationEventType
	WindowFocus                WindowEventType
	WindowLostFocus            WindowEventType
	WindowDidMove              WindowEventType
	WindowDidResize            WindowEventType
	WindowClosing              WindowEventType
	WindowFilesDropped         WindowEventType
}{
	ApplicationOpenedWithFile:  1024,
	ApplicationStarted:         1025,
	ApplicationLaunchedWithUrl: 1026,
	ThemeChanged:               1027,
	WindowFocus:                WindowFocus,
	WindowLostFocus:            WindowLostFocus,
	WindowDidMove:              WindowDidMove,
	WindowDidResize:            WindowDidResize,
	WindowClosing:              WindowClosing,
	WindowFilesDropped:         WindowFilesDropped,
}
