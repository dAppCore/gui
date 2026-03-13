package window

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the window service.
type Options struct{}

// Service is a core.Service managing window lifecycle via IPC.
// It embeds ServiceRuntime for Core access and composes Manager for platform operations.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
}

// OnStartup queries config from the display orchestrator and registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	// Query config — display registers its handler before us (registration order guarantee).
	// If display is not registered, handled=false and we skip config.
	cfg, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if wCfg, ok := cfg.(map[string]any); ok {
			s.applyConfig(wCfg)
		}
	}

	// Register QUERY and TASK handlers manually.
	// ACTION handler (HandleIPCEvents) is auto-registered by WithService —
	// do NOT call RegisterAction here or actions will double-fire.
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(cfg map[string]any) {
	if w, ok := cfg["default_width"]; ok {
		if _, ok := w.(int); ok {
			// TODO: s.manager.SetDefaultWidth(width) — add when Manager API is extended
		}
	}
	if h, ok := cfg["default_height"]; ok {
		if _, ok := h.(int); ok {
			// TODO: s.manager.SetDefaultHeight(height) — add when Manager API is extended
		}
	}
	if sf, ok := cfg["state_file"]; ok {
		if _, ok := sf.(string); ok {
			// TODO: s.manager.State().SetPath(stateFile) — add when StateManager API is extended
		}
	}
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Query Handlers ---

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryWindowList:
		return s.queryWindowList(), true, nil
	case QueryWindowByName:
		return s.queryWindowByName(q.Name), true, nil
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
				Name: name, X: x, Y: y, Width: w, Height: h,
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
		Name: name, X: x, Y: y, Width: w, Height: h,
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
		return nil, true, s.taskSetSize(t.Name, t.W, t.H)
	case TaskMaximise:
		return nil, true, s.taskMaximise(t.Name)
	case TaskMinimise:
		return nil, true, s.taskMinimise(t.Name)
	case TaskFocus:
		return nil, true, s.taskFocus(t.Name)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskOpenWindow(t TaskOpenWindow) (any, bool, error) {
	pw, err := s.manager.Open(t.Opts...)
	if err != nil {
		return nil, true, err
	}
	x, y := pw.Position()
	w, h := pw.Size()
	info := WindowInfo{Name: pw.Name(), X: x, Y: y, Width: w, Height: h}

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
				_ = s.Core().ACTION(ActionWindowResized{Name: e.Name, W: w, H: h})
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
		return fmt.Errorf("window not found: %s", name)
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
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetPosition(x, y)
	s.manager.State().UpdatePosition(name, x, y)
	return nil
}

func (s *Service) taskSetSize(name string, w, h int) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetSize(w, h)
	s.manager.State().UpdateSize(name, w, h)
	return nil
}

func (s *Service) taskMaximise(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Maximise()
	s.manager.State().UpdateMaximized(name, true)
	return nil
}

func (s *Service) taskMinimise(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Minimise()
	return nil
}

func (s *Service) taskFocus(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Focus()
	return nil
}

// Manager returns the underlying window Manager for direct access.
func (s *Service) Manager() *Manager {
	return s.manager
}
