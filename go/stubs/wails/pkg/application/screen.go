package application

import "sync"

type Alignment int
type OffsetReference int

const (
	TOP Alignment = iota
	RIGHT
	BOTTOM
	LEFT
)

const (
	BEGIN OffsetReference = iota
	END
)

// Screen describes a physical or logical display.
//
//	primary := manager.GetPrimary()
//	if primary != nil { useSize(primary.Size) }
type Screen struct {
	ID               string  // A unique identifier for the display
	Name             string  // The name of the display
	ScaleFactor      float32 // The scale factor of the display (DPI/96)
	X                int     // The x-coordinate of the top-left corner of the display
	Y                int     // The y-coordinate of the top-left corner of the display
	Size             Size    // The logical size of the display
	Bounds           Rect    // The logical bounds of the display
	PhysicalBounds   Rect    // The physical bounds before scaling
	WorkArea         Rect    // The work area (excluding taskbars etc.)
	PhysicalWorkArea Rect    // The physical work area before scaling
	IsPrimary        bool    // Whether this is the primary display
	Rotation         float32 // The rotation of the display in degrees
}

// Rect is an axis-aligned rectangle in either logical or physical pixels.
//
//	if rect.Contains(Point{X: 100, Y: 200}) { ... }
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Point is a two-dimensional coordinate.
//
//	centre := Point{X: bounds.X + bounds.Width/2, Y: bounds.Y + bounds.Height/2}
type Point struct {
	X int
	Y int
}

// Size holds the width and height dimensions of a display or region.
//
//	if size.Width > 1920 { useHiResLayout() }
type Size struct {
	Width  int
	Height int
}

type ScreenPlacement struct {
	Screen          *Screen
	Parent          *Screen
	Alignment       Alignment
	Offset          int
	OffsetReference OffsetReference
}

// Origin returns the top-left corner of the rectangle.
func (r Rect) Origin() Point {
	return Point{X: r.X, Y: r.Y}
}

// Corner returns the exclusive bottom-right corner (X+Width, Y+Height).
func (r Rect) Corner() Point {
	return Point{X: r.X + r.Width, Y: r.Y + r.Height}
}

// IsEmpty reports whether the rectangle has non-positive area.
func (r Rect) IsEmpty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// Contains reports whether point pt lies within the rectangle.
//
//	if bounds.Contains(cursorPoint) { highlightWindow() }
func (r Rect) Contains(pt Point) bool {
	return pt.X >= r.X && pt.X < r.X+r.Width && pt.Y >= r.Y && pt.Y < r.Y+r.Height
}

// RectSize returns the dimensions of the rectangle as a Size value.
func (r Rect) RectSize() Size {
	return Size{Width: r.Width, Height: r.Height}
}

func (r Rect) Size() Size {
	return r.RectSize()
}

func (s Screen) Origin() Point {
	return Point{X: s.X, Y: s.Y}
}

func (p ScreenPlacement) Apply() {
	if p.Screen == nil || p.Parent == nil {
		return
	}

	x := p.Parent.X
	y := p.Parent.Y

	switch p.Alignment {
	case TOP:
		x += p.Offset
		y -= p.Screen.Size.Height
	case RIGHT:
		x += p.Parent.Size.Width
		y += p.Offset
	case BOTTOM:
		x += p.Offset
		y += p.Parent.Size.Height
	case LEFT:
		x -= p.Screen.Size.Width
		y += p.Offset
	}

	p.Screen.X = x
	p.Screen.Y = y
	p.Screen.Bounds.X = x
	p.Screen.Bounds.Y = y
}

// ScreenManager tracks connected screens and the active screen.
//
//	manager.SetScreens(detectedScreens)
//	primary := manager.GetPrimary()
type ScreenManager struct {
	mu      sync.RWMutex
	screens []*Screen
	current *Screen
	primary *Screen
}

func newScreenManager() *ScreenManager {
	return &ScreenManager{}
}

// SetScreens replaces the full list of known screens and recomputes primary.
//
//	manager.SetScreens(platformDetectedScreens)
func (m *ScreenManager) SetScreens(screens []*Screen) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.screens = append([]*Screen(nil), screens...)
	m.primary = nil
	for _, screen := range m.screens {
		if screen.IsPrimary {
			m.primary = screen
			break
		}
	}
	if m.current == nil && m.primary != nil {
		m.current = m.primary
	}
}

// SetCurrent marks the given screen as the currently active one.
//
//	manager.SetCurrent(screenUnderPointer)
func (m *ScreenManager) SetCurrent(screen *Screen) {
	if m == nil {
		return
	}
	m.mu.Lock()
	m.current = screen
	m.mu.Unlock()
}

// GetAll returns all registered screens.
//
//	for _, s := range manager.GetAll() { renderMonitorPreview(s) }
func (m *ScreenManager) GetAll() []*Screen {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Screen, len(m.screens))
	copy(out, m.screens)
	return out
}

// GetPrimary returns the primary screen, or nil if none has been registered.
//
//	primary := manager.GetPrimary()
func (m *ScreenManager) GetPrimary() *Screen {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.primary
}

// GetCurrent returns the most recently active screen, or the primary if unset.
//
//	current := manager.GetCurrent()
func (m *ScreenManager) GetCurrent() *Screen {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.current != nil {
		return m.current
	}
	return m.primary
}

func (m *ScreenManager) LayoutScreens(screens []*Screen) error {
	m.SetScreens(screens)
	return nil
}

func (m *ScreenManager) All() []*Screen {
	return m.GetAll()
}

func (m *ScreenManager) Primary() *Screen {
	return m.GetPrimary()
}

func (m *ScreenManager) Current() *Screen {
	return m.GetCurrent()
}

func (m *ScreenManager) DipToPhysicalPoint(dipPoint Point) Point {
	return dipPoint
}

func (m *ScreenManager) PhysicalToDipPoint(physicalPoint Point) Point {
	return physicalPoint
}

func (m *ScreenManager) DipToPhysicalRect(dipRect Rect) Rect {
	return dipRect
}

func (m *ScreenManager) PhysicalToDipRect(physicalRect Rect) Rect {
	return physicalRect
}

func (m *ScreenManager) ScreenNearestPhysicalPoint(physicalPoint Point) *Screen {
	return m.screenNearestPoint(physicalPoint)
}

func (m *ScreenManager) ScreenNearestDipPoint(dipPoint Point) *Screen {
	return m.screenNearestPoint(dipPoint)
}

func (m *ScreenManager) ScreenNearestPhysicalRect(physicalRect Rect) *Screen {
	return m.screenNearestRect(physicalRect)
}

func (m *ScreenManager) ScreenNearestDipRect(dipRect Rect) *Screen {
	return m.screenNearestRect(dipRect)
}

func (m *ScreenManager) screenNearestPoint(point Point) *Screen {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, screen := range m.screens {
		if screen != nil && screen.Bounds.Contains(point) {
			return screen
		}
	}

	if m.current != nil {
		return m.current
	}
	if m.primary != nil {
		return m.primary
	}
	if len(m.screens) > 0 {
		return m.screens[0]
	}
	return nil
}

func (m *ScreenManager) screenNearestRect(rect Rect) *Screen {
	if rect.IsEmpty() {
		return m.screenNearestPoint(Point{})
	}
	return m.screenNearestPoint(rect.Origin())
}
