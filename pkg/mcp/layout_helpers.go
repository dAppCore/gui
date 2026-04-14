package mcp

import (
	"sort"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/window"
)

func (s *Subsystem) allWindows() ([]window.WindowInfo, error) {
	result, _, err := s.core.QUERY(window.QueryWindowList{})
	if err != nil {
		return nil, err
	}
	windows, _ := result.([]window.WindowInfo)
	return windows, nil
}

func (s *Subsystem) allScreens() ([]screen.Screen, error) {
	result, _, err := s.core.QUERY(screen.QueryAll{})
	if err != nil {
		return nil, err
	}
	screens, _ := result.([]screen.Screen)
	return screens, nil
}

func (s *Subsystem) primaryScreen() (*screen.Screen, error) {
	result, _, err := s.core.QUERY(screen.QueryPrimary{})
	if err != nil {
		return nil, err
	}
	scr, _ := result.(*screen.Screen)
	return scr, nil
}

func (s *Subsystem) screenByID(id string) (*screen.Screen, error) {
	if id == "" {
		return nil, nil
	}
	result, _, err := s.core.QUERY(screen.QueryByID{ID: id})
	if err != nil {
		return nil, err
	}
	scr, _ := result.(*screen.Screen)
	return scr, nil
}

func screenForWindowInfo(screens []screen.Screen, info window.WindowInfo) *screen.Screen {
	cx := info.X + info.Width/2
	cy := info.Y + info.Height/2
	for i := range screens {
		if screens[i].Bounds.Contains(cx, cy) {
			return &screens[i]
		}
	}
	return nil
}

func chooseScreenByIDOrPrimary(screens []screen.Screen, screenID string) *screen.Screen {
	if screenID != "" {
		for i := range screens {
			if screens[i].ID == screenID {
				return &screens[i]
			}
		}
	}
	for i := range screens {
		if screens[i].IsPrimary {
			return &screens[i]
		}
	}
	if len(screens) == 0 {
		return nil
	}
	return &screens[0]
}

func workAreaRect(scr *screen.Screen) screen.Rect {
	if scr == nil {
		return screen.Rect{}
	}
	if scr.WorkArea.Width > 0 && scr.WorkArea.Height > 0 {
		return scr.WorkArea
	}
	return scr.Bounds
}

func uniqueSorted(values []int) []int {
	sort.Ints(values)
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func clipRectToWorkArea(rect, workArea screen.Rect) (screen.Rect, bool) {
	x1 := max(rect.X, workArea.X)
	y1 := max(rect.Y, workArea.Y)
	x2 := min(rect.X+rect.Width, workArea.X+workArea.Width)
	y2 := min(rect.Y+rect.Height, workArea.Y+workArea.Height)
	if x2 <= x1 || y2 <= y1 {
		return screen.Rect{}, false
	}
	return screen.Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}, true
}

func findLargestFreeRect(workArea screen.Rect, occupied []screen.Rect, minWidth, minHeight int) (screen.Rect, bool) {
	xs := []int{workArea.X, workArea.X + workArea.Width}
	ys := []int{workArea.Y, workArea.Y + workArea.Height}

	for _, rect := range occupied {
		clipped, ok := clipRectToWorkArea(rect, workArea)
		if !ok {
			continue
		}
		xs = append(xs, clipped.X, clipped.X+clipped.Width)
		ys = append(ys, clipped.Y, clipped.Y+clipped.Height)
	}

	xs = uniqueSorted(xs)
	ys = uniqueSorted(ys)

	bestArea := -1
	best := screen.Rect{}

	for xi := 0; xi < len(xs)-1; xi++ {
		for xj := xi + 1; xj < len(xs); xj++ {
			width := xs[xj] - xs[xi]
			if width < minWidth {
				continue
			}
			for yi := 0; yi < len(ys)-1; yi++ {
				for yj := yi + 1; yj < len(ys); yj++ {
					height := ys[yj] - ys[yi]
					if height < minHeight {
						continue
					}
					candidate := screen.Rect{X: xs[xi], Y: ys[yi], Width: width, Height: height}
					if candidate.X < workArea.X || candidate.Y < workArea.Y ||
						candidate.X+candidate.Width > workArea.X+workArea.Width ||
						candidate.Y+candidate.Height > workArea.Y+workArea.Height {
						continue
					}
					overlaps := false
					for _, occ := range occupied {
						if candidate.Overlaps(occ) {
							overlaps = true
							break
						}
					}
					if overlaps {
						continue
					}
					area := candidate.Width * candidate.Height
					if area > bestArea || (area == bestArea && (candidate.Y < best.Y || (candidate.Y == best.Y && candidate.X < best.X))) {
						bestArea = area
						best = candidate
					}
				}
			}
		}
	}

	return best, bestArea >= 0
}

func applyRect(c *core.Core, windowName string, rect screen.Rect) error {
	if _, _, err := c.PERFORM(window.TaskSetPosition{Name: windowName, X: rect.X, Y: rect.Y}); err != nil {
		return err
	}
	_, _, err := c.PERFORM(window.TaskSetSize{Name: windowName, Width: rect.Width, Height: rect.Height})
	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
