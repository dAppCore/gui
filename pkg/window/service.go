package window

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

type Options struct{}

type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
}

func (s *Service) OnStartup(ctx context.Context) error {
	// Query config — display registers its handler before us (registration order guarantee).
	// If display is not registered, handled=false and we skip config.
	configValue, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if windowConfig, ok := configValue.(map[string]any); ok {
			s.applyConfig(windowConfig)
		}
	}

	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(configData map[string]any) {
	if width, ok := configData["default_width"]; ok {
		if width, ok := width.(int); ok {
			s.manager.SetDefaultWidth(width)
		}
	}
	if height, ok := configData["default_height"]; ok {
		if height, ok := height.(int); ok {
			s.manager.SetDefaultHeight(height)
		}
	}
	if stateFile, ok := configData["state_file"]; ok {
		if stateFile, ok := stateFile.(string); ok {
			s.manager.State().SetPath(stateFile)
		}
	}
}

func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryWindowList:
		return s.queryWindowList(), true, nil
	case QueryWindowByName:
		return s.queryWindowByName(q.Name), true, nil
	case QueryLayoutList:
		return s.manager.Layout().ListLayouts(), true, nil
	case QueryLayoutGet:
		l, ok := s.manager.Layout().GetLayout(q.Name)
		if !ok {
			return (*Layout)(nil), true, nil
		}
		return &l, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) queryWindowList() []WindowInfo {
	names := s.manager.List()
	result := make([]WindowInfo, 0, len(names))
	for _, name := range names {
		if pw, ok := s.manager.Get(name); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			result = append(result, WindowInfo{
				Name: name, Title: pw.Title(), X: x, Y: y, Width: w, Height: h,
				Maximized: pw.IsMaximised(),
				Focused:   pw.IsFocused(),
			})
		}
	}
	return result
}

func (s *Service) queryWindowByName(name string) *WindowInfo {
	pw, ok := s.manager.Get(name)
	if !ok {
		return nil
	}
	x, y := pw.Position()
	w, h := pw.Size()
	return &WindowInfo{
		Name: name, Title: pw.Title(), X: x, Y: y, Width: w, Height: h,
		Maximized: pw.IsMaximised(),
		Focused:   pw.IsFocused(),
	}
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenWindow:
		return s.taskOpenWindow(t)
	case TaskCloseWindow:
		return nil, true, s.taskCloseWindow(t.Name)
	case TaskSetPosition:
		return nil, true, s.taskSetPosition(t.Name, t.X, t.Y)
	case TaskSetSize:
		return nil, true, s.taskSetSize(t.Name, t.Width, t.Height)
	case TaskMaximise:
		return nil, true, s.taskMaximise(t.Name)
	case TaskMinimise:
		return nil, true, s.taskMinimise(t.Name)
	case TaskFocus:
		return nil, true, s.taskFocus(t.Name)
	case TaskRestore:
		return nil, true, s.taskRestore(t.Name)
	case TaskSetTitle:
		return nil, true, s.taskSetTitle(t.Name, t.Title)
	case TaskSetVisibility:
		return nil, true, s.taskSetVisibility(t.Name, t.Visible)
	case TaskFullscreen:
		return nil, true, s.taskFullscreen(t.Name, t.Fullscreen)
	case TaskSaveLayout:
		return nil, true, s.taskSaveLayout(t.Name)
	case TaskRestoreLayout:
		return nil, true, s.taskRestoreLayout(t.Name)
	case TaskDeleteLayout:
		s.manager.Layout().DeleteLayout(t.Name)
		return nil, true, nil
	case TaskTileWindows:
		return nil, true, s.taskTileWindows(t.Mode, t.Windows)
	case TaskSnapWindow:
		return nil, true, s.taskSnapWindow(t.Name, t.Position)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskOpenWindow(t TaskOpenWindow) (any, bool, error) {
	pw, err := s.manager.Open(t.Options...)
	if err != nil {
		return nil, true, err
	}
	x, y := pw.Position()
	w, h := pw.Size()
	info := WindowInfo{Name: pw.Name(), Title: pw.Title(), X: x, Y: y, Width: w, Height: h}

	// Attach platform event listeners that convert to IPC actions
	s.trackWindow(pw)

	// Broadcast to all listeners
	_ = s.Core().ACTION(ActionWindowOpened{Name: pw.Name()})
	return info, true, nil
}

// trackWindow attaches platform event listeners that emit IPC actions.
func (s *Service) trackWindow(pw PlatformWindow) {
	pw.OnWindowEvent(func(e WindowEvent) {
		switch e.Type {
		case "focus":
			_ = s.Core().ACTION(ActionWindowFocused{Name: e.Name})
		case "blur":
			_ = s.Core().ACTION(ActionWindowBlurred{Name: e.Name})
		case "move":
			if data := e.Data; data != nil {
				x, _ := data["x"].(int)
				y, _ := data["y"].(int)
				_ = s.Core().ACTION(ActionWindowMoved{Name: e.Name, X: x, Y: y})
			}
		case "resize":
			if data := e.Data; data != nil {
				w, _ := data["w"].(int)
				h, _ := data["h"].(int)
				_ = s.Core().ACTION(ActionWindowResized{Name: e.Name, Width: w, Height: h})
			}
		case "close":
			_ = s.Core().ACTION(ActionWindowClosed{Name: e.Name})
		}
	})
	pw.OnFileDrop(func(paths []string, targetID string) {
		_ = s.Core().ACTION(ActionFilesDropped{
			Name:     pw.Name(),
			Paths:    paths,
			TargetID: targetID,
		})
	})
}

func (s *Service) taskCloseWindow(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskClose", "window not found: "+name, nil)
	}
	// Persist state BEFORE closing (spec requirement)
	s.manager.State().CaptureState(pw)
	pw.Close()
	s.manager.Remove(name)
	_ = s.Core().ACTION(ActionWindowClosed{Name: name})
	return nil
}

func (s *Service) taskSetPosition(name string, x, y int) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskSetPosition", "window not found: "+name, nil)
	}
	pw.SetPosition(x, y)
	s.manager.State().UpdatePosition(name, x, y)
	return nil
}

func (s *Service) taskSetSize(name string, width, height int) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskSetSize", "window not found: "+name, nil)
	}
	pw.SetSize(width, height)
	s.manager.State().UpdateSize(name, width, height)
	return nil
}

func (s *Service) taskMaximise(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskMaximise", "window not found: "+name, nil)
	}
	pw.Maximise()
	s.manager.State().UpdateMaximized(name, true)
	return nil
}

func (s *Service) taskMinimise(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskMinimise", "window not found: "+name, nil)
	}
	pw.Minimise()
	return nil
}

func (s *Service) taskFocus(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskFocus", "window not found: "+name, nil)
	}
	pw.Focus()
	return nil
}

func (s *Service) taskRestore(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskRestore", "window not found: "+name, nil)
	}
	pw.Restore()
	s.manager.State().UpdateMaximized(name, false)
	return nil
}

func (s *Service) taskSetTitle(name, title string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskSetTitle", "window not found: "+name, nil)
	}
	pw.SetTitle(title)
	return nil
}

func (s *Service) taskSetVisibility(name string, visible bool) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskSetVisibility", "window not found: "+name, nil)
	}
	pw.SetVisibility(visible)
	return nil
}

func (s *Service) taskFullscreen(name string, fullscreen bool) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return coreerr.E("window.taskFullscreen", "window not found: "+name, nil)
	}
	if fullscreen {
		pw.Fullscreen()
	} else {
		pw.UnFullscreen()
	}
	return nil
}

func (s *Service) taskSaveLayout(name string) error {
	windows := s.queryWindowList()
	states := make(map[string]WindowState, len(windows))
	for _, w := range windows {
		states[w.Name] = WindowState{
			X: w.X, Y: w.Y, Width: w.Width, Height: w.Height,
			Maximized: w.Maximized,
		}
	}
	return s.manager.Layout().SaveLayout(name, states)
}

func (s *Service) taskRestoreLayout(name string) error {
	layout, ok := s.manager.Layout().GetLayout(name)
	if !ok {
		return coreerr.E("window.taskRestoreLayout", "layout not found: "+name, nil)
	}
	for winName, state := range layout.Windows {
		pw, found := s.manager.Get(winName)
		if !found {
			continue
		}
		pw.SetPosition(state.X, state.Y)
		pw.SetSize(state.Width, state.Height)
		if state.Maximized {
			pw.Maximise()
		}
	}
	return nil
}

var tileModeMap = map[string]TileMode{
	"left-half": TileModeLeftHalf, "right-half": TileModeRightHalf,
	"top-half": TileModeTopHalf, "bottom-half": TileModeBottomHalf,
	"top-left": TileModeTopLeft, "top-right": TileModeTopRight,
	"bottom-left": TileModeBottomLeft, "bottom-right": TileModeBottomRight,
	"left-right": TileModeLeftRight, "grid": TileModeGrid,
}

func (s *Service) taskTileWindows(mode string, names []string) error {
	tm, ok := tileModeMap[mode]
	if !ok {
		return coreerr.E("window.taskTileWindows", "unknown tile mode: "+mode, nil)
	}
	if len(names) == 0 {
		names = s.manager.List()
	}
	// Default screen size — callers can query screen_primary for actual values.
	return s.manager.TileWindows(tm, names, 1920, 1080)
}

var snapPosMap = map[string]SnapPosition{
	"left": SnapLeft, "right": SnapRight,
	"top": SnapTop, "bottom": SnapBottom,
	"top-left": SnapTopLeft, "top-right": SnapTopRight,
	"bottom-left": SnapBottomLeft, "bottom-right": SnapBottomRight,
	"center": SnapCenter, "centre": SnapCenter,
}

func (s *Service) taskSnapWindow(name, position string) error {
	pos, ok := snapPosMap[position]
	if !ok {
		return coreerr.E("window.taskSnapWindow", "unknown snap position: "+position, nil)
	}
	return s.manager.SnapWindow(name, pos, 1920, 1080)
}

// Manager returns the underlying window Manager for direct access.
func (s *Service) Manager() *Manager {
	return s.manager
}
