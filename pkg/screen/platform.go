// pkg/screen/platform.go
package screen

// Platform abstracts the screen/display backend.
// Use: var p screen.Platform
type Platform interface {
	GetAll() []Screen
	GetPrimary() *Screen
}

// Screen describes a display/monitor.
// Use: scr := screen.Screen{ID: "display-1"}
type Screen struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ScaleFactor float64 `json:"scaleFactor"`
	Size        Size    `json:"size"`
	Bounds      Rect    `json:"bounds"`
	WorkArea    Rect    `json:"workArea"`
	IsPrimary   bool    `json:"isPrimary"`
	Rotation    float64 `json:"rotation"`
}

// Rect represents a rectangle with position and dimensions.
// Use: rect := screen.Rect{X: 0, Y: 0, Width: 1920, Height: 1080}
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Size represents dimensions.
// Use: size := screen.Size{Width: 1920, Height: 1080}
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
