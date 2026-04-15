package application

import "math"

// Alignment describes how a screen is positioned relative to another screen.
type Alignment int

// OffsetReference determines whether an offset is measured from the beginning or end edge.
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

// ScreenPlacement describes how one screen is placed relative to another.
type ScreenPlacement struct {
	screen          *Screen
	parent          *Screen
	alignment       Alignment
	offset          int
	offsetReference OffsetReference
}

// Size returns the dimensions of the rect.
func (r Rect) Size() Size {
	return Size{Width: r.Width, Height: r.Height}
}

// LayoutScreens replaces the tracked screen list.
func (m *ScreenManager) LayoutScreens(screens []*Screen) error {
	if len(screens) > 0 {
		hasPrimary := false
		for _, screen := range screens {
			if screen != nil && screen.IsPrimary {
				hasPrimary = true
				break
			}
		}
		if !hasPrimary && screens[0] != nil {
			screens[0].IsPrimary = true
		}
	}
	m.SetScreens(screens)
	return nil
}

// DipToPhysicalRect converts a DIP rectangle to physical coordinates.
func (m *ScreenManager) DipToPhysicalRect(dipRect Rect) Rect {
	origin := m.DipToPhysicalPoint(dipRect.Origin())
	corner := m.DipToPhysicalPoint(dipRect.Corner())
	return Rect{
		X:      origin.X,
		Y:      origin.Y,
		Width:  corner.X - origin.X,
		Height: corner.Y - origin.Y,
	}
}

// PhysicalToDipRect converts a physical rectangle to DIP coordinates.
func (m *ScreenManager) PhysicalToDipRect(physicalRect Rect) Rect {
	origin := m.PhysicalToDipPoint(physicalRect.Origin())
	corner := m.PhysicalToDipPoint(physicalRect.Corner())
	return Rect{
		X:      origin.X,
		Y:      origin.Y,
		Width:  corner.X - origin.X,
		Height: corner.Y - origin.Y,
	}
}

// ScreenNearestPhysicalPoint returns the nearest screen for physical coordinates.
func (m *ScreenManager) ScreenNearestPhysicalPoint(physicalPoint Point) *Screen {
	return nearestScreenPoint(m.GetAll(), physicalPoint, true)
}

// ScreenNearestDipPoint returns the nearest screen for DIP coordinates.
func (m *ScreenManager) ScreenNearestDipPoint(dipPoint Point) *Screen {
	return nearestScreenPoint(m.GetAll(), dipPoint, false)
}

// ScreenNearestPhysicalRect returns the nearest screen for a physical rectangle.
func (m *ScreenManager) ScreenNearestPhysicalRect(physicalRect Rect) *Screen {
	return nearestScreenRect(m.GetAll(), physicalRect, true)
}

// ScreenNearestDipRect returns the nearest screen for a DIP rectangle.
func (m *ScreenManager) ScreenNearestDipRect(dipRect Rect) *Screen {
	return nearestScreenRect(m.GetAll(), dipRect, false)
}

// DipToPhysicalPoint converts a DIP point using the active screen manager.
func DipToPhysicalPoint(dipPoint Point) Point {
	return appScreens().DipToPhysicalPoint(dipPoint)
}

// PhysicalToDipPoint converts a physical point using the active screen manager.
func PhysicalToDipPoint(physicalPoint Point) Point {
	return appScreens().PhysicalToDipPoint(physicalPoint)
}

// DipToPhysicalRect converts a DIP rect using the active screen manager.
func DipToPhysicalRect(dipRect Rect) Rect {
	return appScreens().DipToPhysicalRect(dipRect)
}

// PhysicalToDipRect converts a physical rect using the active screen manager.
func PhysicalToDipRect(physicalRect Rect) Rect {
	return appScreens().PhysicalToDipRect(physicalRect)
}

// ScreenNearestPhysicalPoint returns the nearest screen for the given physical point.
func ScreenNearestPhysicalPoint(physicalPoint Point) *Screen {
	return appScreens().ScreenNearestPhysicalPoint(physicalPoint)
}

// ScreenNearestDipPoint returns the nearest screen for the given DIP point.
func ScreenNearestDipPoint(dipPoint Point) *Screen {
	return appScreens().ScreenNearestDipPoint(dipPoint)
}

// ScreenNearestPhysicalRect returns the nearest screen for the given physical rect.
func ScreenNearestPhysicalRect(physicalRect Rect) *Screen {
	return appScreens().ScreenNearestPhysicalRect(physicalRect)
}

// ScreenNearestDipRect returns the nearest screen for the given DIP rect.
func ScreenNearestDipRect(dipRect Rect) *Screen {
	return appScreens().ScreenNearestDipRect(dipRect)
}

func nearestScreenPoint(screens []*Screen, point Point, physical bool) *Screen {
	if len(screens) == 0 {
		return nil
	}
	var best *Screen
	bestDistance := math.MaxFloat64
	for _, screen := range screens {
		if screen == nil {
			continue
		}
		bounds := screen.Bounds
		if physical {
			bounds = screen.PhysicalBounds
		}
		if bounds.Contains(point) {
			return screen
		}
		dx := float64(0)
		if point.X < bounds.X {
			dx = float64(bounds.X - point.X)
		} else if point.X > bounds.X+bounds.Width {
			dx = float64(point.X - (bounds.X + bounds.Width))
		}
		dy := float64(0)
		if point.Y < bounds.Y {
			dy = float64(bounds.Y - point.Y)
		} else if point.Y > bounds.Y+bounds.Height {
			dy = float64(point.Y - (bounds.Y + bounds.Height))
		}
		distance := dx*dx + dy*dy
		if distance < bestDistance {
			best = screen
			bestDistance = distance
		}
	}
	return best
}

func nearestScreenRect(screens []*Screen, rect Rect, physical bool) *Screen {
	if len(screens) == 0 {
		return nil
	}
	center := Point{X: rect.X + rect.Width/2, Y: rect.Y + rect.Height/2}
	if screen := nearestScreenPoint(screens, center, physical); screen != nil {
		return screen
	}
	return screens[0]
}
