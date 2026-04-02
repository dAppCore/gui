package window

import (
	"context"
	"fmt"
	"strings"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/screen"
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
	case QueryLayoutList:
		return s.manager.Layout().ListLayouts(), true, nil
	case QueryLayoutGet:
		l, ok := s.manager.Layout().GetLayout(q.Name)
		if !ok {
			return (*Layout)(nil), true, nil
		}
		return &l, true, nil
	case QueryFindSpace:
		screenW, screenH := s.primaryScreenSize()
		return s.manager.FindSpace(screenW, screenH, q.Width, q.Height), true, nil
	case QueryLayoutSuggestion:
		screenW, screenH := s.primaryScreenSize()
		return s.manager.SuggestLayout(screenW, screenH, q.WindowCount), true, nil
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
		return nil, true, s.taskSetSize(t.Name, t.W, t.H)
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
	case TaskArrangePair:
		return nil, true, s.taskArrangePair(t.First, t.Second)
	case TaskBesideEditor:
		return nil, true, s.taskBesideEditor(t.Editor, t.Window)
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

func (s *Service) taskRestore(name string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.Restore()
	s.manager.State().UpdateMaximized(name, false)
	return nil
}

func (s *Service) taskSetTitle(name, title string) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetTitle(title)
	return nil
}

func (s *Service) taskSetVisibility(name string, visible bool) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
	}
	pw.SetVisibility(visible)
	return nil
}

func (s *Service) taskFullscreen(name string, fullscreen bool) error {
	pw, ok := s.manager.Get(name)
	if !ok {
		return fmt.Errorf("window not found: %s", name)
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
		return fmt.Errorf("layout not found: %s", name)
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
		return fmt.Errorf("unknown tile mode: %s", mode)
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
		return fmt.Errorf("unknown snap position: %s", position)
	}
	screenW, screenH := s.primaryScreenSize()
	return s.manager.SnapWindow(name, pos, screenW, screenH)
}

func (s *Service) taskArrangePair(first, second string) error {
	screenW, screenH := s.primaryScreenSize()
	return s.manager.ArrangePair(first, second, screenW, screenH)
}

func (s *Service) taskBesideEditor(editorName, windowName string) error {
	screenW, screenH := s.primaryScreenSize()
	if editorName == "" {
		editorName = s.detectEditorWindow()
	}
	if editorName == "" {
		return fmt.Errorf("editor window not found")
	}
	if windowName == "" {
		windowName = s.detectCompanionWindow(editorName)
	}
	if windowName == "" {
		return fmt.Errorf("companion window not found")
	}
	return s.manager.BesideEditor(editorName, windowName, screenW, screenH)
}

func (s *Service) detectEditorWindow() string {
	for _, info := range s.queryWindowList() {
		if looksLikeEditor(info.Name, info.Title) {
			return info.Name
		}
	}
	return ""
}

func (s *Service) detectCompanionWindow(editorName string) string {
	for _, info := range s.queryWindowList() {
		if info.Name == editorName {
			continue
		}
		if !looksLikeEditor(info.Name, info.Title) {
			return info.Name
		}
	}
	return ""
}

func looksLikeEditor(name, title string) bool {
	return containsAny(name, "editor", "ide", "code", "workspace") || containsAny(title, "editor", "ide", "code")
}

func containsAny(value string, needles ...string) bool {
	lower := strings.ToLower(value)
	for _, needle := range needles {
		if strings.Contains(lower, needle) {
			return true
		}
	}
	return false
}

func (s *Service) primaryScreenSize() (int, int) {
	result, handled, err := s.Core().QUERY(screen.QueryPrimary{})
	if err == nil && handled {
		if scr, ok := result.(*screen.Screen); ok && scr != nil {
			if scr.WorkArea.Width > 0 && scr.WorkArea.Height > 0 {
				return scr.WorkArea.Width, scr.WorkArea.Height
			}
			if scr.Bounds.Width > 0 && scr.Bounds.Height > 0 {
				return scr.Bounds.Width, scr.Bounds.Height
			}
			if scr.Size.Width > 0 && scr.Size.Height > 0 {
				return scr.Size.Width, scr.Size.Height
			}
		}
	}
	return 1920, 1080
}

// Manager returns the underlying window Manager for direct access.
func (s *Service) Manager() *Manager {
	return s.manager
}
