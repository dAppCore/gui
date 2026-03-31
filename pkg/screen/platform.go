// pkg/screen/platform.go
package screen

// Platform abstracts the screen/display backend.
type Platform interface {
	GetAll() []Screen
	GetPrimary() *Screen
	GetCurrent() *Screen
}

// Screen describes a display/monitor.
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
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Contains reports whether the point (x, y) lies within the rectangle.
//
//	if rect.Contains(mouseX, mouseY) { handleClick() }
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

// Overlaps reports whether the rectangle r overlaps with other.
//
//	if bounds.Overlaps(workArea) { show() }
func (r Rect) Overlaps(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Center returns the centre point of the rectangle.
//
//	cx, cy := rect.Center()
func (r Rect) Center() (x, y int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// Size represents dimensions.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ScreenPlacement describes a desired window position relative to a screen.
//
//	p := screen.ScreenPlacement{ScreenID: "1", X: 100, Y: 200, Width: 800, Height: 600}
//	p.Apply(platformWindow)
type ScreenPlacement struct {
	ScreenID string `json:"screenId"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// Placer is implemented by platform windows that can be repositioned.
type Placer interface {
	SetPosition(x, y int)
	SetSize(width, height int)
}

// Apply positions and sizes the given Placer according to the placement.
//
//	placement.Apply(pw)
func (p ScreenPlacement) Apply(target Placer) {
	if p.Width > 0 && p.Height > 0 {
		target.SetSize(p.Width, p.Height)
	}
	target.SetPosition(p.X, p.Y)
}
