// pkg/window/window.go
package window

import (
	"sync"

	coreerr "forge.lthn.ai/core/go-log"
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
// window.NewManager(window.NewWailsPlatform(app))
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManager(),
		layout:   NewLayoutManager(),
		windows:  make(map[string]PlatformWindow),
	}
}

// NewManagerWithDir creates a window Manager with a custom config directory for state/layout persistence.
// window.NewManagerWithDir(window.NewMockPlatform(), t.TempDir())
func NewManagerWithDir(platform Platform, configDir string) *Manager {
	return &Manager{
		platform: platform,
		state:    NewStateManagerWithDir(configDir),
		layout:   NewLayoutManagerWithDir(configDir),
		windows:  make(map[string]PlatformWindow),
	}
}

func (m *Manager) SetDefaultWidth(width int) {
	if width > 0 {
		m.defaultWidth = width
	}
}

func (m *Manager) SetDefaultHeight(height int) {
	if height > 0 {
		m.defaultHeight = height
	}
}

// Deprecated: use CreateWindow(Window{Name: "settings", URL: "/settings", Width: 800, Height: 600}).
func (m *Manager) Open(options ...WindowOption) (PlatformWindow, error) {
	w, err := ApplyOptions(options...)
	if err != nil {
		return nil, coreerr.E("window.Manager.Open", "failed to apply options", err)
	}
	return m.CreateWindow(*w)
}

// CreateWindow creates a window from a Window descriptor.
// window.NewManager(window.NewWailsPlatform(app)).CreateWindow(window.Window{Name: "settings", URL: "/settings", Width: 800, Height: 600})
func (m *Manager) CreateWindow(spec Window) (PlatformWindow, error) {
	_, pw, err := m.createWindow(spec)
	return pw, err
}

// Deprecated: use CreateWindow(Window{Name: "settings", URL: "/settings", Width: 800, Height: 600}).
func (m *Manager) Create(w *Window) (PlatformWindow, error) {
	if w == nil {
		return nil, coreerr.E("window.Manager.Create", "window descriptor is required", nil)
	}
	spec, pw, err := m.createWindow(*w)
	if err != nil {
		return nil, err
	}
	*w = spec
	return pw, nil
}

func (m *Manager) createWindow(spec Window) (Window, PlatformWindow, error) {
	if spec.Name == "" {
		spec.Name = "main"
	}
	if spec.Title == "" {
		spec.Title = "Core"
	}
	if spec.Width == 0 {
		if m.defaultWidth > 0 {
			spec.Width = m.defaultWidth
		} else {
			spec.Width = 1280
		}
	}
	if spec.Height == 0 {
		if m.defaultHeight > 0 {
			spec.Height = m.defaultHeight
		} else {
			spec.Height = 800
		}
	}
	if spec.URL == "" {
		spec.URL = "/"
	}

	// Apply saved state if available.
	m.state.ApplyState(&spec)

	pw := m.platform.CreateWindow(spec.ToPlatformOptions())

	m.mu.Lock()
	m.windows[spec.Name] = pw
	m.mu.Unlock()

	return spec, pw, nil
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
