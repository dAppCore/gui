// pkg/window/platform.go
package window

// Platform abstracts the windowing backend (Wails v3).
// Use: var p window.Platform
type Platform interface {
	CreateWindow(opts PlatformWindowOptions) PlatformWindow
	GetWindows() []PlatformWindow
}

// PlatformWindowOptions are the backend-specific options passed to CreateWindow.
// Use: opts := window.PlatformWindowOptions{Name: "editor"}
type PlatformWindowOptions struct {
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
	BackgroundColour    [4]uint8 // RGBA
	DisableResize       bool
	EnableFileDrop      bool
}

// PlatformWindow is a live window handle from the backend.
// Use: var w window.PlatformWindow
type PlatformWindow interface {
	// Identity
	Name() string
	Title() string

	// Queries
	Position() (int, int)
	Size() (int, int)
	IsVisible() bool
	IsMinimised() bool
	IsMaximised() bool
	IsFocused() bool

	// Mutations
	SetTitle(title string)
	SetPosition(x, y int)
	SetSize(width, height int)
	SetBackgroundColour(r, g, b, a uint8)
	SetOpacity(opacity float32)
	SetVisibility(visible bool)
	SetAlwaysOnTop(alwaysOnTop bool)

	// Window state
	Maximise()
	Restore()
	Minimise()
	Focus()
	Close()
	Show()
	Hide()
	Fullscreen()
	UnFullscreen()
	OpenDevTools()
	CloseDevTools()

	// Events
	OnWindowEvent(handler func(event WindowEvent))

	// File drop
	OnFileDrop(handler func(paths []string, targetID string))
}

// WindowEvent is emitted by the backend for window state changes.
// Use: evt := window.WindowEvent{Type: "focus", Name: "editor"}
type WindowEvent struct {
	Type string // "focus", "blur", "move", "resize", "close"
	Name string // window name
	Data map[string]any
}
