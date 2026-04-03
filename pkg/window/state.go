// pkg/window/state.go
package window

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WindowState holds the persisted position/size of a window.
// JSON tags match existing window_state.json format for backward compat.
// Use: state := window.WindowState{X: 10, Y: 20, Width: 1280, Height: 800}
type WindowState struct {
	X         int    `json:"x,omitempty"`
	Y         int    `json:"y,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Maximized bool   `json:"maximized,omitempty"`
	Screen    string `json:"screen,omitempty"`
	URL       string `json:"url,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

// StateManager persists window positions to ~/.config/Core/window_state.json.
// Use: sm := window.NewStateManager()
type StateManager struct {
	configDir string
	statePath string
	states    map[string]WindowState
	mu        sync.RWMutex
	saveTimer *time.Timer
}

// NewStateManager creates a StateManager loading from the default config directory.
// Use: sm := window.NewStateManager()
func NewStateManager() *StateManager {
	sm := &StateManager{
		states: make(map[string]WindowState),
	}
	configDir, err := os.UserConfigDir()
	if err == nil {
		sm.configDir = filepath.Join(configDir, "Core")
	}
	sm.load()
	return sm
}

// NewStateManagerWithDir creates a StateManager loading from a custom config directory.
// Useful for testing or when the default config directory is not appropriate.
// Use: sm := window.NewStateManagerWithDir(t.TempDir())
func NewStateManagerWithDir(configDir string) *StateManager {
	sm := &StateManager{
		configDir: configDir,
		states:    make(map[string]WindowState),
	}
	sm.load()
	return sm
}

func (sm *StateManager) filePath() string {
	if sm.statePath != "" {
		return sm.statePath
	}
	return filepath.Join(sm.configDir, "window_state.json")
}

func (sm *StateManager) dataDir() string {
	if sm.statePath != "" {
		return filepath.Dir(sm.statePath)
	}
	return sm.configDir
}

// SetPath overrides the persisted state file path.
// Use: sm.SetPath(filepath.Join(t.TempDir(), "window_state.json"))
func (sm *StateManager) SetPath(path string) {
	if path == "" {
		return
	}
	sm.mu.Lock()
	if sm.saveTimer != nil {
		sm.saveTimer.Stop()
		sm.saveTimer = nil
	}
	sm.statePath = path
	sm.states = make(map[string]WindowState)
	sm.mu.Unlock()
	sm.load()
}

func (sm *StateManager) load() {
	if sm.configDir == "" && sm.statePath == "" {
		return
	}
	data, err := os.ReadFile(sm.filePath())
	if err != nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	_ = json.Unmarshal(data, &sm.states)
}

func (sm *StateManager) save() {
	if sm.configDir == "" && sm.statePath == "" {
		return
	}
	sm.mu.RLock()
	data, err := json.MarshalIndent(sm.states, "", "  ")
	sm.mu.RUnlock()
	if err != nil {
		return
	}
	_ = os.MkdirAll(sm.dataDir(), 0o755)
	_ = os.WriteFile(sm.filePath(), data, 0o644)
}

func (sm *StateManager) scheduleSave() {
	if sm.saveTimer != nil {
		sm.saveTimer.Stop()
	}
	sm.saveTimer = time.AfterFunc(500*time.Millisecond, sm.save)
}

// GetState returns the saved state for a window name.
// Use: state, ok := sm.GetState("editor")
func (sm *StateManager) GetState(name string) (WindowState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.states[name]
	return s, ok
}

// SetState saves state for a window name (debounced disk write).
// Use: sm.SetState("editor", window.WindowState{Width: 1280, Height: 800})
func (sm *StateManager) SetState(name string, state WindowState) {
	state.UpdatedAt = time.Now().UnixMilli()
	sm.mu.Lock()
	sm.states[name] = state
	sm.mu.Unlock()
	sm.scheduleSave()
}

// UpdatePosition updates only the position fields.
// Use: sm.UpdatePosition("editor", 160, 120)
func (sm *StateManager) UpdatePosition(name string, x, y int) {
	sm.mu.Lock()
	s := sm.states[name]
	s.X = x
	s.Y = y
	s.UpdatedAt = time.Now().UnixMilli()
	sm.states[name] = s
	sm.mu.Unlock()
	sm.scheduleSave()
}

// UpdateSize updates only the size fields.
// Use: sm.UpdateSize("editor", 1280, 800)
func (sm *StateManager) UpdateSize(name string, width, height int) {
	sm.mu.Lock()
	s := sm.states[name]
	s.Width = width
	s.Height = height
	s.UpdatedAt = time.Now().UnixMilli()
	sm.states[name] = s
	sm.mu.Unlock()
	sm.scheduleSave()
}

// UpdateMaximized updates the maximized flag.
// Use: sm.UpdateMaximized("editor", true)
func (sm *StateManager) UpdateMaximized(name string, maximized bool) {
	sm.mu.Lock()
	s := sm.states[name]
	s.Maximized = maximized
	s.UpdatedAt = time.Now().UnixMilli()
	sm.states[name] = s
	sm.mu.Unlock()
	sm.scheduleSave()
}

// CaptureState snapshots the current state from a PlatformWindow.
// Use: sm.CaptureState(pw)
func (sm *StateManager) CaptureState(pw PlatformWindow) {
	x, y := pw.Position()
	w, h := pw.Size()
	sm.SetState(pw.Name(), WindowState{
		X: x, Y: y, Width: w, Height: h,
		Maximized: pw.IsMaximised(),
	})
}

// ApplyState restores saved position/size to a Window descriptor.
// Use: sm.ApplyState(&window.Window{Name: "editor"})
func (sm *StateManager) ApplyState(w *Window) {
	s, ok := sm.GetState(w.Name)
	if !ok {
		return
	}
	if s.Width > 0 {
		w.Width = s.Width
	}
	if s.Height > 0 {
		w.Height = s.Height
	}
	w.X = s.X
	w.Y = s.Y
}

// ListStates returns all stored window names.
// Use: names := sm.ListStates()
func (sm *StateManager) ListStates() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	names := make([]string, 0, len(sm.states))
	for name := range sm.states {
		names = append(names, name)
	}
	return names
}

// Clear removes all stored states.
// Use: sm.Clear()
func (sm *StateManager) Clear() {
	sm.mu.Lock()
	sm.states = make(map[string]WindowState)
	sm.mu.Unlock()
	sm.scheduleSave()
}

// ForceSync writes state to disk immediately.
// Use: sm.ForceSync()
func (sm *StateManager) ForceSync() {
	if sm.saveTimer != nil {
		sm.saveTimer.Stop()
	}
	sm.save()
}
