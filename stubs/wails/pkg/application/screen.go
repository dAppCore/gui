package application

import (
	"math"
	"sync"
)

// Rect describes an axis-aligned rectangle.
// r := Rect{X: 0, Y: 0, Width: 1920, Height: 1080}
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Point is a 2D integer coordinate.
// p := Point{X: 100, Y: 200}
type Point struct {
	X int
	Y int
}

// Size holds a width/height pair.
// s := Size{Width: 1920, Height: 1080}
type Size struct {
	Width  int
	Height int
}

// Origin returns the top-left corner of the rect as a Point.
// origin := r.Origin()
func (r Rect) Origin() Point { return Point{X: r.X, Y: r.Y} }

// Corner returns the exclusive bottom-right corner of the rect.
// corner := r.Corner()
func (r Rect) Corner() Point { return Point{X: r.X + r.Width, Y: r.Y + r.Height} }

// InsideCorner returns the last pixel inside the rect.
// inside := r.InsideCorner()
func (r Rect) InsideCorner() Point { return Point{X: r.X + r.Width - 1, Y: r.Y + r.Height - 1} }

// RectSize returns the size of the rect.
// s := r.RectSize()
func (r Rect) RectSize() Size { return Size{Width: r.Width, Height: r.Height} }

// IsEmpty reports whether the rect has no area.
// if r.IsEmpty() { return }
func (r Rect) IsEmpty() bool { return r.Width <= 0 || r.Height <= 0 }

// Contains reports whether the point lies inside the rect.
// if r.Contains(Point{X: 100, Y: 200}) { ... }
func (r Rect) Contains(point Point) bool {
	return point.X >= r.X && point.X < r.X+r.Width &&
		point.Y >= r.Y && point.Y < r.Y+r.Height
}

// Intersect returns the overlapping rectangle of r and other.
// overlap := r.Intersect(otherRect)
func (r Rect) Intersect(other Rect) Rect {
	if r.IsEmpty() || other.IsEmpty() {
		return Rect{}
	}
	maxLeft := max(r.X, other.X)
	maxTop := max(r.Y, other.Y)
	minRight := min(r.X+r.Width, other.X+other.Width)
	minBottom := min(r.Y+r.Height, other.Y+other.Height)
	if minRight > maxLeft && minBottom > maxTop {
		return Rect{X: maxLeft, Y: maxTop, Width: minRight - maxLeft, Height: minBottom - maxTop}
	}
	return Rect{}
}

// Screen describes a physical display.
// primary := manager.GetPrimary()
type Screen struct {
	ID               string
	Name             string
	ScaleFactor      float32
	X                int
	Y                int
	Size             Size
	Bounds           Rect
	PhysicalBounds   Rect
	WorkArea         Rect
	PhysicalWorkArea Rect
	IsPrimary        bool
	Rotation         float32
}

// Origin returns the logical origin of the screen.
// origin := screen.Origin()
func (s Screen) Origin() Point { return Point{X: s.X, Y: s.Y} }

// scale converts a value between physical and DIP coordinates for this screen.
func (s Screen) scale(value int, toDIP bool) int {
	if toDIP {
		return int(math.Ceil(float64(value) / float64(s.ScaleFactor)))
	}
	return int(math.Floor(float64(value) * float64(s.ScaleFactor)))
}

// ScreenManager holds the list of known screens in memory.
// screens := manager.GetAll()
type ScreenManager struct {
	mu            sync.RWMutex
	screens       []*Screen
	primaryScreen *Screen
}

// NewScreenManager constructs an empty ScreenManager.
// manager := NewScreenManager()
func NewScreenManager() *ScreenManager {
	return &ScreenManager{}
}

// SetScreens replaces the tracked screen list; the first screen with IsPrimary == true becomes primary.
// manager.SetScreens([]*Screen{primary, secondary})
func (m *ScreenManager) SetScreens(screens []*Screen) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.screens = screens
	m.primaryScreen = nil
	for _, screen := range screens {
		if screen.IsPrimary {
			m.primaryScreen = screen
			break
		}
	}
}

// GetAll returns all tracked screens.
// for _, s := range manager.GetAll() { fmt.Println(s.Name) }
func (m *ScreenManager) GetAll() []*Screen {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Screen, len(m.screens))
	copy(result, m.screens)
	return result
}

// GetPrimary returns the primary screen, or nil if none is set.
// primary := manager.GetPrimary()
func (m *ScreenManager) GetPrimary() *Screen {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primaryScreen
}

// GetCurrent returns the screen nearest to the given DIP point.
// current := manager.GetCurrent(Point{X: 500, Y: 300})
func (m *ScreenManager) GetCurrent(dipPoint Point) *Screen {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, screen := range m.screens {
		if screen.Bounds.Contains(dipPoint) {
			return screen
		}
	}
	return m.primaryScreen
}

// DipToPhysicalPoint converts a DIP point to physical pixels using the nearest screen.
// physical := manager.DipToPhysicalPoint(Point{X: 100, Y: 200})
func (m *ScreenManager) DipToPhysicalPoint(dipPoint Point) Point {
	screen := m.GetCurrent(dipPoint)
	if screen == nil {
		return dipPoint
	}
	relativeX := dipPoint.X - screen.Bounds.X
	relativeY := dipPoint.Y - screen.Bounds.Y
	return Point{
		X: screen.PhysicalBounds.X + screen.scale(relativeX, false),
		Y: screen.PhysicalBounds.Y + screen.scale(relativeY, false),
	}
}

// PhysicalToDipPoint converts physical pixel coordinates to DIP using the nearest screen.
// dip := manager.PhysicalToDipPoint(Point{X: 200, Y: 400})
func (m *ScreenManager) PhysicalToDipPoint(physicalPoint Point) Point {
	m.mu.RLock()
	var nearest *Screen
	for _, screen := range m.screens {
		if screen.PhysicalBounds.Contains(physicalPoint) {
			nearest = screen
			break
		}
	}
	if nearest == nil {
		nearest = m.primaryScreen
	}
	m.mu.RUnlock()
	if nearest == nil {
		return physicalPoint
	}
	relativeX := physicalPoint.X - nearest.PhysicalBounds.X
	relativeY := physicalPoint.Y - nearest.PhysicalBounds.Y
	return Point{
		X: nearest.Bounds.X + nearest.scale(relativeX, true),
		Y: nearest.Bounds.Y + nearest.scale(relativeY, true),
	}
}
