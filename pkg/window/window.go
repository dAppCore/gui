// pkg/window/window.go
package window

import (
	"fmt"
	"math"
	"sync"
)

// Window is CoreGUI's own window descriptor — NOT a Wails type alias.
type Window struct {
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
	BackgroundColour    [4]uint8
	DisableResize       bool
	EnableFileDrop      bool
}

// ToPlatformOptions converts a Window to PlatformWindowOptions for the backend.
func (w *Window) ToPlatformOptions() PlatformWindowOptions {
	return PlatformWindowOptions{
		Name: w.Name, Title: w.Title, URL: w.URL,
		Width: w.Width, Height: w.Height, X: w.X, Y: w.Y,
		MinWidth: w.MinWidth, MinHeight: w.MinHeight,
		MaxWidth: w.MaxWidth, MaxHeight: w.MaxHeight,
		Frameless: w.Frameless, Hidden: w.Hidden,
		AlwaysOnTop: w.AlwaysOnTop, BackgroundColour: w.BackgroundColour,
		DisableResize: w.DisableResize, EnableFileDrop: w.EnableFileDrop,
	}
}

// Manager manages window lifecycle through a Platform backend.
type Manager struct {
	platform      Platform
	state         *StateManager
	layout        *LayoutManager
	windows       map[string]PlatformWindow
	defaultWidth  int
	defaultHeight int
	mu            sync.RWMutex
}

// NewManager creates a window Manager with the given platform backend.
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManager(),
		layout:   NewLayoutManager(),
		windows:  make(map[string]PlatformWindow),
	}
}

// NewManagerWithDir creates a window Manager with a custom config directory for state/layout persistence.
// Useful for testing or when the default config directory is not appropriate.
func NewManagerWithDir(platform Platform, configDir string) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManagerWithDir(configDir),
		layout:   NewLayoutManagerWithDir(configDir),
		windows:  make(map[string]PlatformWindow),
	}
}

// SetDefaultWidth overrides the fallback width used when a window is created without one.
func (m *Manager) SetDefaultWidth(width int) {
	if width > 0 {
		m.defaultWidth = width
	}
}

// SetDefaultHeight overrides the fallback height used when a window is created without one.
func (m *Manager) SetDefaultHeight(height int) {
	if height > 0 {
		m.defaultHeight = height
	}
}

// Open creates a window using functional options, applies saved state, and tracks it.
func (m *Manager) Open(opts ...WindowOption) (PlatformWindow, error) {
	w, err := ApplyOptions(opts...)
	if err != nil {
		return nil, fmt.Errorf("window.Manager.Open: %w", err)
	}
	return m.Create(w)
}

// Create creates a window from a Window descriptor.
func (m *Manager) Create(w *Window) (PlatformWindow, error) {
	if w.Name == "" {
		w.Name = "main"
	}
	if w.Title == "" {
		w.Title = "Core"
	}
	if w.Width == 0 {
		if m.defaultWidth > 0 {
			w.Width = m.defaultWidth
		} else {
			w.Width = 1280
		}
	}
	if w.Height == 0 {
		if m.defaultHeight > 0 {
			w.Height = m.defaultHeight
		} else {
			w.Height = 800
		}
	}
	if w.URL == "" {
		w.URL = "/"
	}

	// Apply saved state if available
	m.state.ApplyState(w)

	pw := m.platform.CreateWindow(w.ToPlatformOptions())

	m.mu.Lock()
	m.windows[w.Name] = pw
	m.mu.Unlock()

	return pw, nil
}

// Get returns a tracked window by name.
func (m *Manager) Get(name string) (PlatformWindow, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pw, ok := m.windows[name]
	return pw, ok
}

// List returns all tracked window names.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.windows))
	for name := range m.windows {
		names = append(names, name)
	}
	return names
}

// Remove stops tracking a window by name.
func (m *Manager) Remove(name string) {
	m.mu.Lock()
	delete(m.windows, name)
	m.mu.Unlock()
}

// Platform returns the underlying platform for direct access.
func (m *Manager) Platform() Platform {
	return m.platform
}

// State returns the state manager for window persistence.
func (m *Manager) State() *StateManager {
	return m.state
}

// Layout returns the layout manager.
func (m *Manager) Layout() *LayoutManager {
	return m.layout
}

// SuggestLayout returns a simple layout recommendation for the given screen.
func (m *Manager) SuggestLayout(screenW, screenH, windowCount int) LayoutSuggestion {
	if windowCount <= 1 {
		return LayoutSuggestion{
			Mode:           "single",
			Columns:        1,
			Rows:           1,
			PrimaryWidth:   screenW,
			SecondaryWidth: 0,
			Description:    "Focus the primary window and keep the screen uncluttered.",
		}
	}

	if windowCount == 2 {
		return LayoutSuggestion{
			Mode:           "side-by-side",
			Columns:        2,
			Rows:           1,
			PrimaryWidth:   screenW / 2,
			SecondaryWidth: screenW - (screenW / 2),
			Description:    "Split the screen into two equal panes.",
		}
	}

	if windowCount <= 4 {
		return LayoutSuggestion{
			Mode:           "quadrants",
			Columns:        2,
			Rows:           2,
			PrimaryWidth:   screenW / 2,
			SecondaryWidth: screenW / 2,
			Description:    "Use a 2x2 grid for the active windows.",
		}
	}

	cols := 3
	rows := int(math.Ceil(float64(windowCount) / float64(cols)))
	return LayoutSuggestion{
		Mode:           "grid",
		Columns:        cols,
		Rows:           rows,
		PrimaryWidth:   screenW / cols,
		SecondaryWidth: screenW / cols,
		Description:    "Use a dense grid to keep every window visible.",
	}
}

// FindSpace returns a free placement suggestion for a new window.
func (m *Manager) FindSpace(screenW, screenH, width, height int) SpaceInfo {
	if width <= 0 {
		width = screenW / 2
	}
	if height <= 0 {
		height = screenH / 2
	}

	occupied := make([]struct {
		x, y, w, h int
	}, 0)
	for _, name := range m.List() {
		pw, ok := m.Get(name)
		if !ok {
			continue
		}
		x, y := pw.Position()
		w, h := pw.Size()
		occupied = append(occupied, struct {
			x, y, w, h int
		}{x: x, y: y, w: w, h: h})
	}

	step := int(math.Max(40, math.Min(float64(width), float64(height))/6))
	if step < 40 {
		step = 40
	}

	for y := 0; y+height <= screenH; y += step {
		for x := 0; x+width <= screenW; x += step {
			if !intersectsAny(x, y, width, height, occupied) {
				return SpaceInfo{
					X: x, Y: y, Width: width, Height: height,
					ScreenWidth: screenW, ScreenHeight: screenH,
					Reason: "first available gap",
				}
			}
		}
	}

	return SpaceInfo{
		X: (screenW - width) / 2, Y: (screenH - height) / 2,
		Width: width, Height: height,
		ScreenWidth: screenW, ScreenHeight: screenH,
		Reason: "center fallback",
	}
}

// ArrangePair places two windows side-by-side with a balanced split.
func (m *Manager) ArrangePair(first, second string, screenW, screenH int) error {
	left, ok := m.Get(first)
	if !ok {
		return fmt.Errorf("window %q not found", first)
	}
	right, ok := m.Get(second)
	if !ok {
		return fmt.Errorf("window %q not found", second)
	}

	leftW := screenW / 2
	rightW := screenW - leftW
	left.SetPosition(0, 0)
	left.SetSize(leftW, screenH)
	right.SetPosition(leftW, 0)
	right.SetSize(rightW, screenH)
	return nil
}

// BesideEditor places a target window beside an editor window, using a 70/30 split.
func (m *Manager) BesideEditor(editorName, windowName string, screenW, screenH int) error {
	editor, ok := m.Get(editorName)
	if !ok {
		return fmt.Errorf("window %q not found", editorName)
	}
	target, ok := m.Get(windowName)
	if !ok {
		return fmt.Errorf("window %q not found", windowName)
	}

	editorW := screenW * 70 / 100
	if editorW <= 0 {
		editorW = screenW / 2
	}
	targetW := screenW - editorW
	editor.SetPosition(0, 0)
	editor.SetSize(editorW, screenH)
	target.SetPosition(editorW, 0)
	target.SetSize(targetW, screenH)
	return nil
}

func intersectsAny(x, y, w, h int, occupied []struct{ x, y, w, h int }) bool {
	for _, r := range occupied {
		if x < r.x+r.w && x+w > r.x && y < r.y+r.h && y+h > r.y {
			return true
		}
	}
	return false
}
