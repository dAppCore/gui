// pkg/window/window.go
package window

import (
	"fmt"
	"sync"
)

// Window is CoreGUI's own window descriptor — NOT a Wails type alias.
type Window struct {
	Name              string
	Title             string
	URL               string
	Width, Height     int
	X, Y              int
	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	Frameless         bool
	Hidden            bool
	AlwaysOnTop       bool
	BackgroundColour  [4]uint8
	DisableResize     bool
	EnableFileDrop    bool
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
	platform Platform
	state    *StateManager
	layout   *LayoutManager
	windows  map[string]PlatformWindow
	mu       sync.RWMutex
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
		w.Width = 1280
	}
	if w.Height == 0 {
		w.Height = 800
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
