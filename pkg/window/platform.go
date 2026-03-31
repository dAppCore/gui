// pkg/window/platform.go
package window

// Platform abstracts the windowing backend (Wails v3).
type Platform interface {
	CreateWindow(options PlatformWindowOptions) PlatformWindow
	GetWindows() []PlatformWindow
}

// PlatformWindowOptions are the backend-specific options passed to CreateWindow.
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
type PlatformWindow interface {
	// Identity
	Name() string
	Title() string

	// Queries
	Position() (int, int)
	Size() (int, int)
	IsMaximised() bool
	IsFocused() bool

	// Mutations
	SetTitle(title string)
	SetPosition(x, y int)
	SetSize(width, height int)
	SetBackgroundColour(r, g, b, a uint8)
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

	// Zoom
	GetZoom() float64
	SetZoom(factor float64)
	ZoomIn()
	ZoomOut()

	// Content
	SetURL(url string)
	SetHTML(html string)
	ExecJS(js string)

	// Bounds
	GetBounds() Bounds
	SetBounds(bounds Bounds)

	// Extras
	ToggleFullscreen()
	Print() error
	Flash(enabled bool)

	// Events
	OnWindowEvent(handler func(event WindowEvent))

	// File drop
	OnFileDrop(handler func(paths []string, targetID string))
}

// Bounds holds the position and dimensions of a window.
type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WindowEvent is emitted by the backend for window state changes.
type WindowEvent struct {
	Type string // "focus", "blur", "move", "resize", "close"
	Name string // window name
	Data map[string]any
}
