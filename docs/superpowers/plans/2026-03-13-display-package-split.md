# CoreGUI Display Package Split — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split `pkg/display/` monolith (3,910 LOC, 15 files) into 4 focused packages: `pkg/window`, `pkg/systray`, `pkg/menu`, and a slimmed `pkg/display` orchestrator.

**Architecture:** Each sub-package defines a `Platform` interface insulating Wails v3. The orchestrator (`pkg/display`) composes `window.Manager`, `systray.Manager`, `menu.Manager` and owns the WebSocket event bridge, dialogs, clipboard, notifications, and theme. No circular dependencies — sub-packages are peers.

**Tech Stack:** Go 1.26, Wails v3 alpha-74, gorilla/websocket, testify, core/go DI

**Spec:** `docs/superpowers/specs/2026-03-13-display-package-split-design.md`

---

## File Map

### New files to create

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/window` | `platform.go` | `Platform`, `PlatformWindow`, `PlatformWindowOptions`, `WindowEvent` interfaces |
| `pkg/window` | `window.go` | `Window` struct (own type), `Manager` struct, CRUD |
| `pkg/window` | `options.go` | `WindowOption` functional options against `Window` |
| `pkg/window` | `state.go` | `WindowStateManager` — JSON persistence |
| `pkg/window` | `layout.go` | `LayoutManager` — named arrangements |
| `pkg/window` | `tiling.go` | `TileMode`, `SnapPosition`, tiling/snapping/stacking/workflows |
| `pkg/window` | `wails.go` | Wails adapter implementing `Platform` + `PlatformWindow` |
| `pkg/window` | `mock_test.go` | `mockPlatform` + `mockWindow` |
| `pkg/window` | `window_test.go` | All window tests |
| `pkg/systray` | `platform.go` | `Platform`, `PlatformTray`, `PlatformMenu` interfaces |
| `pkg/systray` | `types.go` | `TrayMenuItem` struct |
| `pkg/systray` | `tray.go` | `Manager` struct, lifecycle, icon, tooltip, label |
| `pkg/systray` | `menu.go` | Dynamic menu builder, callback registry |
| `pkg/systray` | `wails.go` | Wails adapter |
| `pkg/systray` | `mock_test.go` | `mockPlatform` + `mockTray` |
| `pkg/systray` | `tray_test.go` | All systray tests |
| `pkg/menu` | `platform.go` | `Platform`, `PlatformMenu`, `PlatformMenuItem` interfaces |
| `pkg/menu` | `menu.go` | `Manager` struct, `MenuItem`, builder (structure only) |
| `pkg/menu` | `wails.go` | Wails adapter |
| `pkg/menu` | `mock_test.go` | `mockPlatform` + `mockMenu` |
| `pkg/menu` | `menu_test.go` | All menu tests |

### Files to modify

| File | Change |
|------|--------|
| `pkg/display/display.go` | Remove window CRUD/tiling/snapping/workflows (~800 LOC). Compose sub-managers. |
| `pkg/display/interfaces.go` | Remove migrated interfaces. Keep `DialogManager`, `EnvManager`, `EventManager`, `Logger`. |
| `pkg/display/window.go` | DELETE — replaced by `pkg/window/options.go` |
| `pkg/display/window_state.go` | DELETE — replaced by `pkg/window/state.go` |
| `pkg/display/layout.go` | DELETE — replaced by `pkg/window/layout.go` |
| `pkg/display/tray.go` | DELETE — replaced by `pkg/systray/tray.go` |
| `pkg/display/menu.go` | DELETE — handlers stay in display.go, structure moves to `pkg/menu` |
| `pkg/display/actions.go` | Update `ActionOpenWindow` to use `window.Window` not Wails type |
| `pkg/display/events.go` | `AttachWindowListeners` accepts `window.PlatformWindow`, add `EventSource` interface |
| `pkg/display/display_test.go` | Update for new imports, split window tests to `pkg/window` |
| `pkg/display/mocks_test.go` | Remove migrated mocks, keep display-level mocks |

### Files to move

| From | To |
|------|-----|
| `pkg/display/ui/` (entire dir) | `ui/` (top-level) |
| `pkg/display/assets/apptray.png` | `pkg/systray/assets/apptray.png` |

### New shared types file

| File | Types |
|------|-------|
| `pkg/display/types.go` | `WindowHandle` interface, `ScreenInfo`, `WorkArea` structs |

---

## Chunk 1: pkg/window

### Task 1: Platform interfaces and mock

**Files:**
- Create: `pkg/window/platform.go`
- Create: `pkg/window/mock_test.go`

- [ ] **Step 1: Write platform.go**

```go
// pkg/window/platform.go
package window

// Platform abstracts the windowing backend (Wails v3).
type Platform interface {
	CreateWindow(opts PlatformWindowOptions) PlatformWindow
	GetWindows() []PlatformWindow
}

// PlatformWindowOptions are the backend-specific options passed to CreateWindow.
type PlatformWindowOptions struct {
	Name                  string
	Title                 string
	URL                   string
	Width, Height         int
	X, Y                  int
	MinWidth, MinHeight   int
	MaxWidth, MaxHeight   int
	Frameless             bool
	Hidden                bool
	AlwaysOnTop           bool
	BackgroundColour      [4]uint8 // RGBA
	DisableResize         bool
	EnableDragAndDrop     bool
	Centered              bool
}

// PlatformWindow is a live window handle from the backend.
type PlatformWindow interface {
	// Identity
	Name() string

	// Queries
	Position() (int, int)
	Size() (int, int)
	IsMaximised() bool
	IsFocused() bool

	// Mutations
	SetTitle(title string)
	SetPosition(x, y int)
	SetSize(width, height int)
	SetBackgroundColour(r, g, b, a uint8)
	SetVisibility(visible bool)
	SetAlwaysOnTop(alwaysOnTop bool)

	// Window state
	Maximise()
	Restore()
	Minimise()
	Focus()
	Close()
	Show()
	Hide()
	Fullscreen()
	UnFullscreen()

	// Events
	OnWindowEvent(handler func(event WindowEvent))
}

// WindowEvent is emitted by the backend for window state changes.
type WindowEvent struct {
	Type string // "focus", "blur", "move", "resize", "close"
	Name string // window name
	Data map[string]any
}
```

- [ ] **Step 2: Write mock_test.go**

```go
// pkg/window/mock_test.go
package window

type mockPlatform struct {
	windows []*mockWindow
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{}
}

func (m *mockPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	w := &mockWindow{
		name: opts.Name, title: opts.Title, url: opts.URL,
		width: opts.Width, height: opts.Height,
		x: opts.X, y: opts.Y,
	}
	m.windows = append(m.windows, w)
	return w
}

func (m *mockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.windows))
	for i, w := range m.windows {
		out[i] = w
	}
	return out
}

type mockWindow struct {
	name, title, url      string
	width, height, x, y   int
	maximised, focused     bool
	visible, alwaysOnTop   bool
	closed                 bool
	eventHandlers          []func(WindowEvent)
}

func (w *mockWindow) Name() string                              { return w.name }
func (w *mockWindow) Position() (int, int)                      { return w.x, w.y }
func (w *mockWindow) Size() (int, int)                          { return w.width, w.height }
func (w *mockWindow) IsMaximised() bool                         { return w.maximised }
func (w *mockWindow) IsFocused() bool                           { return w.focused }
func (w *mockWindow) SetTitle(title string)                     { w.title = title }
func (w *mockWindow) SetPosition(x, y int)                      { w.x = x; w.y = y }
func (w *mockWindow) SetSize(width, height int)                 { w.width = width; w.height = height }
func (w *mockWindow) SetBackgroundColour(r, g, b, a uint8)      {}
func (w *mockWindow) SetVisibility(visible bool)                { w.visible = visible }
func (w *mockWindow) SetAlwaysOnTop(alwaysOnTop bool)           { w.alwaysOnTop = alwaysOnTop }
func (w *mockWindow) Maximise()                                 { w.maximised = true }
func (w *mockWindow) Restore()                                  { w.maximised = false }
func (w *mockWindow) Minimise()                                 {}
func (w *mockWindow) Focus()                                    { w.focused = true }
func (w *mockWindow) Close()                                    { w.closed = true }
func (w *mockWindow) Show()                                     { w.visible = true }
func (w *mockWindow) Hide()                                     { w.visible = false }
func (w *mockWindow) Fullscreen()                               {}
func (w *mockWindow) UnFullscreen()                             {}
func (w *mockWindow) OnWindowEvent(handler func(WindowEvent))   { w.eventHandlers = append(w.eventHandlers, handler) }

// emit fires a test event to all registered handlers.
func (w *mockWindow) emit(e WindowEvent) {
	for _, h := range w.eventHandlers {
		h(e)
	}
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd /Users/snider/Code/core/gui && go build ./pkg/window/...`
Expected: SUCCESS (no test binary, just compile check)

- [ ] **Step 4: Commit**

```bash
git add pkg/window/platform.go pkg/window/mock_test.go
git commit -m "feat(window): add Platform and PlatformWindow interfaces"
```

---

### Task 2: Window struct, options, and Manager

**Files:**
- Create: `pkg/window/window.go`
- Create: `pkg/window/options.go`
- Create: `pkg/window/window_test.go`

- [ ] **Step 1: Write window_test.go — Window struct and option tests**

```go
// pkg/window/window_test.go
package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowDefaults(t *testing.T) {
	w := &Window{}
	assert.Equal(t, "", w.Name)
	assert.Equal(t, 0, w.Width)
}

func TestWindowOption_Name_Good(t *testing.T) {
	w := &Window{}
	err := WithName("main")(w)
	require.NoError(t, err)
	assert.Equal(t, "main", w.Name)
}

func TestWindowOption_Title_Good(t *testing.T) {
	w := &Window{}
	err := WithTitle("My App")(w)
	require.NoError(t, err)
	assert.Equal(t, "My App", w.Title)
}

func TestWindowOption_URL_Good(t *testing.T) {
	w := &Window{}
	err := WithURL("/dashboard")(w)
	require.NoError(t, err)
	assert.Equal(t, "/dashboard", w.URL)
}

func TestWindowOption_Size_Good(t *testing.T) {
	w := &Window{}
	err := WithSize(1280, 720)(w)
	require.NoError(t, err)
	assert.Equal(t, 1280, w.Width)
	assert.Equal(t, 720, w.Height)
}

func TestWindowOption_Position_Good(t *testing.T) {
	w := &Window{}
	err := WithPosition(100, 200)(w)
	require.NoError(t, err)
	assert.Equal(t, 100, w.X)
	assert.Equal(t, 200, w.Y)
}

func TestApplyOptions_Good(t *testing.T) {
	w, err := ApplyOptions(
		WithName("test"),
		WithTitle("Test Window"),
		WithURL("/test"),
		WithSize(800, 600),
	)
	require.NoError(t, err)
	assert.Equal(t, "test", w.Name)
	assert.Equal(t, "Test Window", w.Title)
	assert.Equal(t, "/test", w.URL)
	assert.Equal(t, 800, w.Width)
	assert.Equal(t, 600, w.Height)
}

func TestApplyOptions_Bad(t *testing.T) {
	_, err := ApplyOptions(func(w *Window) error {
		return assert.AnError
	})
	assert.Error(t, err)
}

func TestApplyOptions_Empty_Good(t *testing.T) {
	w, err := ApplyOptions()
	require.NoError(t, err)
	assert.NotNil(t, w)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -v`
Expected: FAIL — `Window`, `WithName`, `ApplyOptions` etc. undefined

- [ ] **Step 3: Write window.go — Window struct and Manager**

```go
// pkg/window/window.go
package window

import (
	"fmt"
	"sync"
)

// Window is CoreGUI's own window descriptor — NOT a Wails type alias.
type Window struct {
	Name                  string
	Title                 string
	URL                   string
	Width, Height         int
	X, Y                  int
	MinWidth, MinHeight   int
	MaxWidth, MaxHeight   int
	Frameless             bool
	Hidden                bool
	AlwaysOnTop           bool
	BackgroundColour      [4]uint8
	DisableResize         bool
	EnableDragAndDrop     bool
	Centered              bool
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
		DisableResize: w.DisableResize, EnableDragAndDrop: w.EnableDragAndDrop,
		Centered: w.Centered,
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
```

- [ ] **Step 4: Write options.go — WindowOption functional options**

```go
// pkg/window/options.go
package window

// WindowOption is a functional option applied to a Window descriptor.
type WindowOption func(*Window) error

// ApplyOptions creates a Window and applies all options in order.
func ApplyOptions(opts ...WindowOption) (*Window, error) {
	w := &Window{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func WithName(name string) WindowOption {
	return func(w *Window) error { w.Name = name; return nil }
}

func WithTitle(title string) WindowOption {
	return func(w *Window) error { w.Title = title; return nil }
}

func WithURL(url string) WindowOption {
	return func(w *Window) error { w.URL = url; return nil }
}

func WithSize(width, height int) WindowOption {
	return func(w *Window) error { w.Width = width; w.Height = height; return nil }
}

func WithPosition(x, y int) WindowOption {
	return func(w *Window) error { w.X = x; w.Y = y; return nil }
}

func WithMinSize(width, height int) WindowOption {
	return func(w *Window) error { w.MinWidth = width; w.MinHeight = height; return nil }
}

func WithMaxSize(width, height int) WindowOption {
	return func(w *Window) error { w.MaxWidth = width; w.MaxHeight = height; return nil }
}

func WithFrameless(frameless bool) WindowOption {
	return func(w *Window) error { w.Frameless = frameless; return nil }
}

func WithHidden(hidden bool) WindowOption {
	return func(w *Window) error { w.Hidden = hidden; return nil }
}

func WithAlwaysOnTop(alwaysOnTop bool) WindowOption {
	return func(w *Window) error { w.AlwaysOnTop = alwaysOnTop; return nil }
}

func WithBackgroundColour(r, g, b, a uint8) WindowOption {
	return func(w *Window) error { w.BackgroundColour = [4]uint8{r, g, b, a}; return nil }
}

func WithCentered(centered bool) WindowOption {
	return func(w *Window) error { w.Centered = centered; return nil }
}
```

- [ ] **Step 5: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -v`
Expected: PASS (some tests will fail on StateManager — that's Task 3)

Note: `NewStateManager` and `NewLayoutManager` are referenced but not yet created. Add stubs:

```go
// Temporary stubs in window.go (remove after Task 3)
// func NewStateManager() *StateManager { return &StateManager{} }
// func NewLayoutManager() *LayoutManager { return &LayoutManager{} }
```

Actually — write minimal stubs in state.go and layout.go so tests pass:

- [ ] **Step 5a: Write minimal state.go stub**

```go
// pkg/window/state.go
package window

// StateManager persists window positions to disk.
// Full implementation in Task 3.
type StateManager struct{}

func NewStateManager() *StateManager { return &StateManager{} }

// ApplyState restores saved position/size to a Window descriptor.
func (sm *StateManager) ApplyState(w *Window) {}
```

- [ ] **Step 5b: Write minimal layout.go stub**

```go
// pkg/window/layout.go
package window

// LayoutManager persists named window arrangements.
// Full implementation in Task 3.
type LayoutManager struct{}

func NewLayoutManager() *LayoutManager { return &LayoutManager{} }
```

- [ ] **Step 6: Add Manager tests to window_test.go**

Append to `pkg/window/window_test.go`:

```go
// newTestManager creates a Manager with a mock platform for testing.
func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Open_Good(t *testing.T) {
	m, p := newTestManager()
	pw, err := m.Open(WithName("test"), WithTitle("Test"), WithURL("/test"), WithSize(800, 600))
	require.NoError(t, err)
	assert.NotNil(t, pw)
	assert.Equal(t, "test", pw.Name())
	assert.Len(t, p.windows, 1)
}

func TestManager_Open_Defaults_Good(t *testing.T) {
	m, _ := newTestManager()
	pw, err := m.Open()
	require.NoError(t, err)
	assert.Equal(t, "main", pw.Name())
	w, h := pw.Size()
	assert.Equal(t, 1280, w)
	assert.Equal(t, 800, h)
}

func TestManager_Open_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, err := m.Open(func(w *Window) error { return assert.AnError })
	assert.Error(t, err)
}

func TestManager_Get_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("findme"))
	pw, ok := m.Get("findme")
	assert.True(t, ok)
	assert.Equal(t, "findme", pw.Name())
}

func TestManager_Get_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, ok := m.Get("nonexistent")
	assert.False(t, ok)
}

func TestManager_List_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"))
	_, _ = m.Open(WithName("b"))
	names := m.List()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "a")
	assert.Contains(t, names, "b")
}

func TestManager_Remove_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("temp"))
	m.Remove("temp")
	_, ok := m.Get("temp")
	assert.False(t, ok)
}
```

- [ ] **Step 7: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -v -count=1`
Expected: ALL PASS

- [ ] **Step 8: Commit**

```bash
git add pkg/window/
git commit -m "feat(window): add Window struct, options, and Manager with CRUD"
```

---

### Task 3: State persistence

**Files:**
- Modify: `pkg/window/state.go` (replace stub)
- Test: `pkg/window/window_test.go` (append state tests)

- [ ] **Step 1: Write state tests**

Append to `pkg/window/window_test.go`:

```go
func TestStateManager_SetGet_Good(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	state := WindowState{X: 100, Y: 200, Width: 800, Height: 600, Maximized: false}
	sm.SetState("main", state)
	got, ok := sm.GetState("main")
	assert.True(t, ok)
	assert.Equal(t, 100, got.X)
	assert.Equal(t, 800, got.Width)
}

func TestStateManager_SetGet_Bad(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	_, ok := sm.GetState("nonexistent")
	assert.False(t, ok)
}

func TestStateManager_CaptureState_Good(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	w := &mockWindow{name: "cap", x: 50, y: 60, width: 1024, height: 768, maximised: true}
	sm.CaptureState(w)
	got, ok := sm.GetState("cap")
	assert.True(t, ok)
	assert.Equal(t, 50, got.X)
	assert.Equal(t, 1024, got.Width)
	assert.True(t, got.Maximized)
}

func TestStateManager_ApplyState_Good(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	sm.SetState("win", WindowState{X: 10, Y: 20, Width: 640, Height: 480})
	w := &Window{Name: "win", Width: 1280, Height: 800}
	sm.ApplyState(w)
	assert.Equal(t, 10, w.X)
	assert.Equal(t, 20, w.Y)
	assert.Equal(t, 640, w.Width)
	assert.Equal(t, 480, w.Height)
}

func TestStateManager_ListStates_Good(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	sm.SetState("a", WindowState{Width: 100})
	sm.SetState("b", WindowState{Width: 200})
	names := sm.ListStates()
	assert.Len(t, names, 2)
}

func TestStateManager_Clear_Good(t *testing.T) {
	sm := NewStateManager()
	sm.configDir = t.TempDir()
	sm.SetState("a", WindowState{Width: 100})
	sm.Clear()
	names := sm.ListStates()
	assert.Empty(t, names)
}

func TestStateManager_Persistence_Good(t *testing.T) {
	dir := t.TempDir()
	sm1 := NewStateManager()
	sm1.configDir = dir
	sm1.SetState("persist", WindowState{X: 42, Y: 84, Width: 500, Height: 300})
	sm1.ForceSync()

	sm2 := NewStateManager()
	sm2.configDir = dir
	sm2.load()
	got, ok := sm2.GetState("persist")
	assert.True(t, ok)
	assert.Equal(t, 42, got.X)
	assert.Equal(t, 500, got.Width)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -run TestStateManager -v`
Expected: FAIL — stub has no real implementation

- [ ] **Step 3: Replace state.go with full implementation**

Migrate from `pkg/display/window_state.go` (262 LOC), changing:
- Accept `PlatformWindow` in `CaptureState` (not `*application.WebviewWindow`)
- Return/modify `*Window` in `ApplyState` (not `*application.WebviewWindowOptions`)
- Export `configDir` field for test injection

```go
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
type StateManager struct {
	configDir string
	states    map[string]WindowState
	mu        sync.RWMutex
	saveTimer *time.Timer
}

// NewStateManager creates a StateManager loading from the default config directory.
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

func (sm *StateManager) filePath() string {
	return filepath.Join(sm.configDir, "window_state.json")
}

func (sm *StateManager) load() {
	if sm.configDir == "" {
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
	if sm.configDir == "" {
		return
	}
	sm.mu.RLock()
	data, err := json.MarshalIndent(sm.states, "", "  ")
	sm.mu.RUnlock()
	if err != nil {
		return
	}
	_ = os.MkdirAll(sm.configDir, 0o755)
	_ = os.WriteFile(sm.filePath(), data, 0o644)
}

func (sm *StateManager) scheduleSave() {
	if sm.saveTimer != nil {
		sm.saveTimer.Stop()
	}
	sm.saveTimer = time.AfterFunc(500*time.Millisecond, sm.save)
}

// GetState returns the saved state for a window name.
func (sm *StateManager) GetState(name string) (WindowState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.states[name]
	return s, ok
}

// SetState saves state for a window name (debounced disk write).
func (sm *StateManager) SetState(name string, state WindowState) {
	state.UpdatedAt = time.Now().UnixMilli()
	sm.mu.Lock()
	sm.states[name] = state
	sm.mu.Unlock()
	sm.scheduleSave()
}

// UpdatePosition updates only the position fields.
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
func (sm *StateManager) CaptureState(pw PlatformWindow) {
	x, y := pw.Position()
	w, h := pw.Size()
	sm.SetState(pw.Name(), WindowState{
		X: x, Y: y, Width: w, Height: h,
		Maximized: pw.IsMaximised(),
	})
}

// ApplyState restores saved position/size to a Window descriptor.
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
func (sm *StateManager) Clear() {
	sm.mu.Lock()
	sm.states = make(map[string]WindowState)
	sm.mu.Unlock()
	sm.scheduleSave()
}

// ForceSync writes state to disk immediately.
func (sm *StateManager) ForceSync() {
	if sm.saveTimer != nil {
		sm.saveTimer.Stop()
	}
	sm.save()
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -run TestStateManager -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/window/state.go pkg/window/window_test.go
git commit -m "feat(window): add StateManager with JSON persistence"
```

---

### Task 4: Layout management

**Files:**
- Modify: `pkg/window/layout.go` (replace stub)
- Test: `pkg/window/window_test.go` (append layout tests)

- [ ] **Step 1: Write layout tests**

Append to `pkg/window/window_test.go`:

```go
func TestLayoutManager_SaveGet_Good(t *testing.T) {
	lm := NewLayoutManager()
	lm.configDir = t.TempDir()
	states := map[string]WindowState{
		"editor": {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 1080},
	}
	err := lm.SaveLayout("coding", states)
	require.NoError(t, err)

	layout, ok := lm.GetLayout("coding")
	assert.True(t, ok)
	assert.Equal(t, "coding", layout.Name)
	assert.Len(t, layout.Windows, 2)
}

func TestLayoutManager_GetLayout_Bad(t *testing.T) {
	lm := NewLayoutManager()
	lm.configDir = t.TempDir()
	_, ok := lm.GetLayout("nonexistent")
	assert.False(t, ok)
}

func TestLayoutManager_ListLayouts_Good(t *testing.T) {
	lm := NewLayoutManager()
	lm.configDir = t.TempDir()
	_ = lm.SaveLayout("a", map[string]WindowState{})
	_ = lm.SaveLayout("b", map[string]WindowState{})
	layouts := lm.ListLayouts()
	assert.Len(t, layouts, 2)
}

func TestLayoutManager_DeleteLayout_Good(t *testing.T) {
	lm := NewLayoutManager()
	lm.configDir = t.TempDir()
	_ = lm.SaveLayout("temp", map[string]WindowState{})
	lm.DeleteLayout("temp")
	_, ok := lm.GetLayout("temp")
	assert.False(t, ok)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -run TestLayoutManager -v`
Expected: FAIL

- [ ] **Step 3: Replace layout.go with full implementation**

Migrate from `pkg/display/layout.go` (150 LOC):

```go
// pkg/window/layout.go
package window

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Layout is a named window arrangement.
type Layout struct {
	Name      string                  `json:"name"`
	Windows   map[string]WindowState  `json:"windows"`
	CreatedAt int64                   `json:"createdAt"`
	UpdatedAt int64                   `json:"updatedAt"`
}

// LayoutInfo is a summary of a layout.
type LayoutInfo struct {
	Name        string `json:"name"`
	WindowCount int    `json:"windowCount"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// LayoutManager persists named window arrangements to ~/.config/Core/layouts.json.
type LayoutManager struct {
	configDir string
	layouts   map[string]Layout
	mu        sync.RWMutex
}

// NewLayoutManager creates a LayoutManager loading from the default config directory.
func NewLayoutManager() *LayoutManager {
	lm := &LayoutManager{
		layouts: make(map[string]Layout),
	}
	configDir, err := os.UserConfigDir()
	if err == nil {
		lm.configDir = filepath.Join(configDir, "Core")
	}
	lm.load()
	return lm
}

func (lm *LayoutManager) filePath() string {
	return filepath.Join(lm.configDir, "layouts.json")
}

func (lm *LayoutManager) load() {
	if lm.configDir == "" {
		return
	}
	data, err := os.ReadFile(lm.filePath())
	if err != nil {
		return
	}
	lm.mu.Lock()
	defer lm.mu.Unlock()
	_ = json.Unmarshal(data, &lm.layouts)
}

func (lm *LayoutManager) save() {
	if lm.configDir == "" {
		return
	}
	lm.mu.RLock()
	data, err := json.MarshalIndent(lm.layouts, "", "  ")
	lm.mu.RUnlock()
	if err != nil {
		return
	}
	_ = os.MkdirAll(lm.configDir, 0o755)
	_ = os.WriteFile(lm.filePath(), data, 0o644)
}

// SaveLayout creates or updates a named layout.
func (lm *LayoutManager) SaveLayout(name string, windowStates map[string]WindowState) error {
	if name == "" {
		return fmt.Errorf("layout name cannot be empty")
	}
	now := time.Now().UnixMilli()
	lm.mu.Lock()
	existing, exists := lm.layouts[name]
	layout := Layout{
		Name:      name,
		Windows:   windowStates,
		UpdatedAt: now,
	}
	if exists {
		layout.CreatedAt = existing.CreatedAt
	} else {
		layout.CreatedAt = now
	}
	lm.layouts[name] = layout
	lm.mu.Unlock()
	lm.save()
	return nil
}

// GetLayout returns a layout by name.
func (lm *LayoutManager) GetLayout(name string) (Layout, bool) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	l, ok := lm.layouts[name]
	return l, ok
}

// ListLayouts returns info summaries for all layouts.
func (lm *LayoutManager) ListLayouts() []LayoutInfo {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	infos := make([]LayoutInfo, 0, len(lm.layouts))
	for _, l := range lm.layouts {
		infos = append(infos, LayoutInfo{
			Name: l.Name, WindowCount: len(l.Windows),
			CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
		})
	}
	return infos
}

// DeleteLayout removes a layout by name.
func (lm *LayoutManager) DeleteLayout(name string) {
	lm.mu.Lock()
	delete(lm.layouts, name)
	lm.mu.Unlock()
	lm.save()
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -run TestLayoutManager -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/window/layout.go pkg/window/window_test.go
git commit -m "feat(window): add LayoutManager with JSON persistence"
```

---

### Task 5: Tiling, snapping, and workflows

**Files:**
- Create: `pkg/window/tiling.go`
- Test: `pkg/window/window_test.go` (append tiling tests)

- [ ] **Step 1: Write tiling tests**

Append to `pkg/window/window_test.go`:

```go
func TestTileMode_String_Good(t *testing.T) {
	assert.Equal(t, "left-half", TileModeLeftHalf.String())
	assert.Equal(t, "grid", TileModeGrid.String())
}

func TestManager_TileWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"), WithSize(800, 600))
	_, _ = m.Open(WithName("b"), WithSize(800, 600))
	err := m.TileWindows(TileModeLeftRight, []string{"a", "b"}, 1920, 1080)
	require.NoError(t, err)
	a, _ := m.Get("a")
	b, _ := m.Get("b")
	aw, _ := a.Size()
	bw, _ := b.Size()
	assert.Equal(t, 960, aw)
	assert.Equal(t, 960, bw)
}

func TestManager_TileWindows_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.TileWindows(TileModeLeftRight, []string{"nonexistent"}, 1920, 1080)
	assert.Error(t, err)
}

func TestManager_SnapWindow_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("snap"), WithSize(800, 600))
	err := m.SnapWindow("snap", SnapLeft, 1920, 1080)
	require.NoError(t, err)
	w, _ := m.Get("snap")
	x, _ := w.Position()
	assert.Equal(t, 0, x)
	sw, _ := w.Size()
	assert.Equal(t, 960, sw)
}

func TestManager_StackWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("s1"), WithSize(800, 600))
	_, _ = m.Open(WithName("s2"), WithSize(800, 600))
	err := m.StackWindows([]string{"s1", "s2"}, 30, 30)
	require.NoError(t, err)
	s2, _ := m.Get("s2")
	x, y := s2.Position()
	assert.Equal(t, 30, x)
	assert.Equal(t, 30, y)
}

func TestWorkflowLayout_Good(t *testing.T) {
	assert.Equal(t, "coding", WorkflowCoding.String())
	assert.Equal(t, "debugging", WorkflowDebugging.String())
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -run "TestTile|TestSnap|TestStack|TestWorkflow" -v`
Expected: FAIL

- [ ] **Step 3: Write tiling.go**

Migrate tiling/snapping/stacking/workflow code from `pkg/display/display.go` lines 859-1293:

```go
// pkg/window/tiling.go
package window

import "fmt"

// TileMode defines how windows are arranged.
type TileMode int

const (
	TileModeLeftHalf TileMode = iota
	TileModeRightHalf
	TileModeTopHalf
	TileModeBottomHalf
	TileModeTopLeft
	TileModeTopRight
	TileModeBottomLeft
	TileModeBottomRight
	TileModeLeftRight
	TileModeGrid
)

var tileModeNames = map[TileMode]string{
	TileModeLeftHalf: "left-half", TileModeRightHalf: "right-half",
	TileModeTopHalf: "top-half", TileModeBottomHalf: "bottom-half",
	TileModeTopLeft: "top-left", TileModeTopRight: "top-right",
	TileModeBottomLeft: "bottom-left", TileModeBottomRight: "bottom-right",
	TileModeLeftRight: "left-right", TileModeGrid: "grid",
}

func (m TileMode) String() string { return tileModeNames[m] }

// SnapPosition defines where a window snaps to.
type SnapPosition int

const (
	SnapLeft SnapPosition = iota
	SnapRight
	SnapTop
	SnapBottom
	SnapTopLeft
	SnapTopRight
	SnapBottomLeft
	SnapBottomRight
	SnapCenter
)

// WorkflowLayout is a predefined arrangement for common tasks.
type WorkflowLayout int

const (
	WorkflowCoding    WorkflowLayout = iota // 70/30 split
	WorkflowDebugging                       // 60/40 split
	WorkflowPresenting                      // maximised
	WorkflowSideBySide                      // 50/50 split
)

var workflowNames = map[WorkflowLayout]string{
	WorkflowCoding: "coding", WorkflowDebugging: "debugging",
	WorkflowPresenting: "presenting", WorkflowSideBySide: "side-by-side",
}

func (w WorkflowLayout) String() string { return workflowNames[w] }

// TileWindows arranges the named windows in the given mode across the screen area.
func (m *Manager) TileWindows(mode TileMode, names []string, screenW, screenH int) error {
	windows := make([]PlatformWindow, 0, len(names))
	for _, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return fmt.Errorf("window %q not found", name)
		}
		windows = append(windows, pw)
	}
	if len(windows) == 0 {
		return fmt.Errorf("no windows to tile")
	}

	halfW, halfH := screenW/2, screenH/2

	switch mode {
	case TileModeLeftRight:
		w := screenW / len(windows)
		for i, pw := range windows {
			pw.SetPosition(i*w, 0)
			pw.SetSize(w, screenH)
		}
	case TileModeGrid:
		cols := 2
		if len(windows) > 4 {
			cols = 3
		}
		cellW := screenW / cols
		for i, pw := range windows {
			row := i / cols
			col := i % cols
			rows := (len(windows) + cols - 1) / cols
			cellH := screenH / rows
			pw.SetPosition(col*cellW, row*cellH)
			pw.SetSize(cellW, cellH)
		}
	case TileModeLeftHalf:
		for _, pw := range windows {
			pw.SetPosition(0, 0)
			pw.SetSize(halfW, screenH)
		}
	case TileModeRightHalf:
		for _, pw := range windows {
			pw.SetPosition(halfW, 0)
			pw.SetSize(halfW, screenH)
		}
	case TileModeTopHalf:
		for _, pw := range windows {
			pw.SetPosition(0, 0)
			pw.SetSize(screenW, halfH)
		}
	case TileModeBottomHalf:
		for _, pw := range windows {
			pw.SetPosition(0, halfH)
			pw.SetSize(screenW, halfH)
		}
	case TileModeTopLeft:
		for _, pw := range windows {
			pw.SetPosition(0, 0)
			pw.SetSize(halfW, halfH)
		}
	case TileModeTopRight:
		for _, pw := range windows {
			pw.SetPosition(halfW, 0)
			pw.SetSize(halfW, halfH)
		}
	case TileModeBottomLeft:
		for _, pw := range windows {
			pw.SetPosition(0, halfH)
			pw.SetSize(halfW, halfH)
		}
	case TileModeBottomRight:
		for _, pw := range windows {
			pw.SetPosition(halfW, halfH)
			pw.SetSize(halfW, halfH)
		}
	}
	return nil
}

// SnapWindow snaps a window to a screen edge/corner/centre.
func (m *Manager) SnapWindow(name string, pos SnapPosition, screenW, screenH int) error {
	pw, ok := m.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}

	halfW, halfH := screenW/2, screenH/2

	switch pos {
	case SnapLeft:
		pw.SetPosition(0, 0)
		pw.SetSize(halfW, screenH)
	case SnapRight:
		pw.SetPosition(halfW, 0)
		pw.SetSize(halfW, screenH)
	case SnapTop:
		pw.SetPosition(0, 0)
		pw.SetSize(screenW, halfH)
	case SnapBottom:
		pw.SetPosition(0, halfH)
		pw.SetSize(screenW, halfH)
	case SnapTopLeft:
		pw.SetPosition(0, 0)
		pw.SetSize(halfW, halfH)
	case SnapTopRight:
		pw.SetPosition(halfW, 0)
		pw.SetSize(halfW, halfH)
	case SnapBottomLeft:
		pw.SetPosition(0, halfH)
		pw.SetSize(halfW, halfH)
	case SnapBottomRight:
		pw.SetPosition(halfW, halfH)
		pw.SetSize(halfW, halfH)
	case SnapCenter:
		cw, ch := pw.Size()
		pw.SetPosition((screenW-cw)/2, (screenH-ch)/2)
	}
	return nil
}

// StackWindows cascades windows with an offset.
func (m *Manager) StackWindows(names []string, offsetX, offsetY int) error {
	for i, name := range names {
		pw, ok := m.Get(name)
		if !ok {
			return fmt.Errorf("window %q not found", name)
		}
		pw.SetPosition(i*offsetX, i*offsetY)
	}
	return nil
}

// ApplyWorkflow arranges windows in a predefined workflow layout.
func (m *Manager) ApplyWorkflow(workflow WorkflowLayout, names []string, screenW, screenH int) error {
	if len(names) == 0 {
		return fmt.Errorf("no windows for workflow")
	}

	switch workflow {
	case WorkflowCoding:
		// 70/30 split — main editor + terminal
		mainW := screenW * 70 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(0, 0)
			pw.SetSize(mainW, screenH)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(mainW, 0)
				pw.SetSize(screenW-mainW, screenH)
			}
		}
	case WorkflowDebugging:
		// 60/40 split
		mainW := screenW * 60 / 100
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(0, 0)
			pw.SetSize(mainW, screenH)
		}
		if len(names) > 1 {
			if pw, ok := m.Get(names[1]); ok {
				pw.SetPosition(mainW, 0)
				pw.SetSize(screenW-mainW, screenH)
			}
		}
	case WorkflowPresenting:
		// Maximise first window
		if pw, ok := m.Get(names[0]); ok {
			pw.SetPosition(0, 0)
			pw.SetSize(screenW, screenH)
		}
	case WorkflowSideBySide:
		return m.TileWindows(TileModeLeftRight, names, screenW, screenH)
	}
	return nil
}
```

- [ ] **Step 4: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/... -v -count=1`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add pkg/window/tiling.go pkg/window/window_test.go
git commit -m "feat(window): add tiling, snapping, stacking, and workflow layouts"
```

---

### Task 6: Wails adapter for pkg/window

**Files:**
- Create: `pkg/window/wails.go`

- [ ] **Step 1: Write wails.go**

```go
// pkg/window/wails.go
package window

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

// NewWailsPlatform creates a Wails-backed Platform.
func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	wOpts := application.WebviewWindowOptions{
		Name:              opts.Name,
		Title:             opts.Title,
		URL:               opts.URL,
		Width:             opts.Width,
		Height:            opts.Height,
		X:                 opts.X,
		Y:                 opts.Y,
		MinWidth:          opts.MinWidth,
		MinHeight:         opts.MinHeight,
		MaxWidth:          opts.MaxWidth,
		MaxHeight:         opts.MaxHeight,
		Frameless:         opts.Frameless,
		Hidden:            opts.Hidden,
		AlwaysOnTop:       opts.AlwaysOnTop,
		DisableResize:     opts.DisableResize,
		EnableDragAndDrop: opts.EnableDragAndDrop,
		Centered:          opts.Centered,
		BackgroundColour:  application.NewRGBA(opts.BackgroundColour[0], opts.BackgroundColour[1], opts.BackgroundColour[2], opts.BackgroundColour[3]),
	}
	w := wp.app.NewWebviewWindowWithOptions(wOpts)
	return &wailsWindow{w: w}
}

func (wp *WailsPlatform) GetWindows() []PlatformWindow {
	all := wp.app.GetWindowByName // Wails doesn't expose GetAll directly
	// Use the app's internal window list — adapt based on Wails v3 API
	return nil // TODO: implement once Wails v3 exposes window enumeration
}

// wailsWindow wraps *application.WebviewWindow to implement PlatformWindow.
type wailsWindow struct {
	w *application.WebviewWindow
}

func (ww *wailsWindow) Name() string                          { return ww.w.Name() }
func (ww *wailsWindow) Position() (int, int)                  { return ww.w.Position() }
func (ww *wailsWindow) Size() (int, int)                      { return ww.w.Size() }
func (ww *wailsWindow) IsMaximised() bool                     { return ww.w.IsMaximised() }
func (ww *wailsWindow) IsFocused() bool                       { return ww.w.IsFocused() }
func (ww *wailsWindow) SetTitle(title string)                 { ww.w.SetTitle(title) }
func (ww *wailsWindow) SetPosition(x, y int)                  { ww.w.SetPosition(x, y) }
func (ww *wailsWindow) SetSize(width, height int)             { ww.w.SetSize(width, height) }
func (ww *wailsWindow) SetBackgroundColour(r, g, b, a uint8)  { ww.w.SetBackgroundColour(application.NewRGBA(r, g, b, a)) }
func (ww *wailsWindow) SetVisibility(visible bool)            { if visible { ww.w.Show() } else { ww.w.Hide() } }
func (ww *wailsWindow) SetAlwaysOnTop(alwaysOnTop bool)       { ww.w.SetAlwaysOnTop(alwaysOnTop) }
func (ww *wailsWindow) Maximise()                             { ww.w.Maximise() }
func (ww *wailsWindow) Restore()                              { ww.w.Restore() }
func (ww *wailsWindow) Minimise()                             { ww.w.Minimise() }
func (ww *wailsWindow) Focus()                                { ww.w.Focus() }
func (ww *wailsWindow) Close()                                { ww.w.Close() }
func (ww *wailsWindow) Show()                                 { ww.w.Show() }
func (ww *wailsWindow) Hide()                                 { ww.w.Hide() }
func (ww *wailsWindow) Fullscreen()                           { ww.w.Fullscreen() }
func (ww *wailsWindow) UnFullscreen()                         { ww.w.UnFullscreen() }

func (ww *wailsWindow) OnWindowEvent(handler func(event WindowEvent)) {
	name := ww.w.Name()
	ww.w.OnWindowEvent(func(e *application.WindowEvent) {
		handler(WindowEvent{
			Type: e.EventType.String(),
			Name: name,
		})
	})
}
```

Note: The `GetWindows()` and `OnWindowEvent` implementations may need adjusting based on exact Wails v3 API. The engineer should check `wails/v3/pkg/application` for the correct method signatures. The key contract is that the adapter wraps Wails and nothing outside this file touches Wails types.

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/snider/Code/core/gui && go build ./pkg/window/...`
Expected: SUCCESS (may need minor API adjustments — Wails v3 alpha)

- [ ] **Step 3: Commit**

```bash
git add pkg/window/wails.go
git commit -m "feat(window): add Wails v3 adapter"
```

---

## Chunk 2: pkg/systray + pkg/menu

### Task 7: pkg/systray — Platform, Manager, and menu builder

**Files:**
- Create: `pkg/systray/platform.go`
- Create: `pkg/systray/types.go`
- Create: `pkg/systray/tray.go`
- Create: `pkg/systray/menu.go`
- Create: `pkg/systray/mock_test.go`
- Create: `pkg/systray/tray_test.go`
- Create: `pkg/systray/wails.go`
- Move: `pkg/display/assets/apptray.png` → `pkg/systray/assets/apptray.png`

- [ ] **Step 1: Copy apptray.png asset FIRST (required for `//go:embed` in tray.go)**

```bash
mkdir -p /Users/snider/Code/core/gui/pkg/systray/assets
cp /Users/snider/Code/core/gui/pkg/display/assets/apptray.png /Users/snider/Code/core/gui/pkg/systray/assets/apptray.png
```

- [ ] **Step 2: Write platform.go**

```go
// pkg/systray/platform.go
package systray

// Platform abstracts the system tray backend.
type Platform interface {
	NewTray() PlatformTray
	NewMenu() PlatformMenu // Menu factory for building tray menus
}

// PlatformTray is a live tray handle from the backend.
type PlatformTray interface {
	SetIcon(data []byte)
	SetTemplateIcon(data []byte)
	SetTooltip(text string)
	SetLabel(text string)
	SetMenu(menu PlatformMenu)
	AttachWindow(w WindowHandle)
}

// PlatformMenu is a tray menu built by the backend.
type PlatformMenu interface {
	Add(label string) PlatformMenuItem
	AddSeparator()
}

// PlatformMenuItem is a single item in a tray menu.
type PlatformMenuItem interface {
	SetTooltip(text string)
	SetChecked(checked bool)
	SetEnabled(enabled bool)
	OnClick(fn func())
	AddSubmenu() PlatformMenu
}

// WindowHandle is a cross-package interface for window operations.
// Defined locally to avoid circular imports (display imports systray).
// pkg/window.PlatformWindow satisfies this implicitly.
type WindowHandle interface {
	Name() string
	Show()
	Hide()
	SetPosition(x, y int)
	SetSize(width, height int)
}
```

Note: `WindowHandle` is defined locally in `pkg/systray` — NOT imported from `pkg/display`. This avoids a circular dependency (`display` → `systray` → `display`). Go's implicit interface satisfaction means `window.PlatformWindow` satisfies this without any coupling.

- [ ] **Step 3: Write types.go**

```go
// pkg/systray/types.go
package systray

// TrayMenuItem describes a menu item for dynamic tray menus.
type TrayMenuItem struct {
	Label    string         `json:"label"`
	Type     string         `json:"type"` // "normal", "separator", "checkbox", "radio"
	Checked  bool           `json:"checked,omitempty"`
	Disabled bool           `json:"disabled,omitempty"`
	Tooltip  string         `json:"tooltip,omitempty"`
	Submenu  []TrayMenuItem `json:"submenu,omitempty"`
	ActionID string         `json:"action_id,omitempty"`
}
```

- [ ] **Step 4: Write tray.go — Manager struct**

```go
// pkg/systray/tray.go
package systray

import (
	_ "embed"
	"fmt"
	"sync"
)

//go:embed assets/apptray.png
var defaultIcon []byte

// Manager manages the system tray lifecycle.
// State that was previously in package-level vars is now on the Manager.
type Manager struct {
	platform  Platform
	tray      PlatformTray
	callbacks map[string]func()
	mu        sync.RWMutex
}

// NewManager creates a systray Manager.
func NewManager(platform Platform) *Manager {
	return &Manager{
		platform:  platform,
		callbacks: make(map[string]func()),
	}
}

// Setup creates the system tray with default icon and tooltip.
func (m *Manager) Setup(tooltip, label string) error {
	m.tray = m.platform.NewTray()
	if m.tray == nil {
		return fmt.Errorf("platform returned nil tray")
	}
	m.tray.SetTemplateIcon(defaultIcon)
	m.tray.SetTooltip(tooltip)
	m.tray.SetLabel(label)
	return nil
}

// SetIcon sets the tray icon.
func (m *Manager) SetIcon(data []byte) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetIcon(data)
	return nil
}

// SetTemplateIcon sets the template icon (macOS).
func (m *Manager) SetTemplateIcon(data []byte) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetTemplateIcon(data)
	return nil
}

// SetTooltip sets the tray tooltip.
func (m *Manager) SetTooltip(text string) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetTooltip(text)
	return nil
}

// SetLabel sets the tray label.
func (m *Manager) SetLabel(text string) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.SetLabel(text)
	return nil
}

// AttachWindow attaches a panel window to the tray.
func (m *Manager) AttachWindow(w WindowHandle) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	m.tray.AttachWindow(w)
	return nil
}

// Tray returns the underlying platform tray for direct access.
func (m *Manager) Tray() PlatformTray {
	return m.tray
}

// IsActive returns whether a tray has been created.
func (m *Manager) IsActive() bool {
	return m.tray != nil
}
```

- [ ] **Step 5: Write menu.go — Dynamic menu builder and callback registry**

```go
// pkg/systray/menu.go
package systray

import "fmt"

// SetMenu sets a dynamic menu on the tray from TrayMenuItem descriptors.
func (m *Manager) SetMenu(items []TrayMenuItem) error {
	if m.tray == nil {
		return fmt.Errorf("tray not initialised")
	}
	menu := m.buildMenu(items)
	m.tray.SetMenu(menu)
	return nil
}

// buildMenu recursively builds a PlatformMenu from TrayMenuItem descriptors.
func (m *Manager) buildMenu(items []TrayMenuItem) PlatformMenu {
	menu := m.platform.NewMenu()
	for _, item := range items {
		if item.Type == "separator" {
			menu.AddSeparator()
			continue
		}
		if len(item.Submenu) > 0 {
			sub := m.buildMenu(item.Submenu)
			mi := menu.Add(item.Label)
			_ = mi.AddSubmenu()
			_ = sub // TODO: wire sub into parent via platform
			continue
		}
		mi := menu.Add(item.Label)
		if item.Tooltip != "" {
			mi.SetTooltip(item.Tooltip)
		}
		if item.Disabled {
			mi.SetEnabled(false)
		}
		if item.Checked {
			mi.SetChecked(true)
		}
		if item.ActionID != "" {
			actionID := item.ActionID
			mi.OnClick(func() {
				if cb, ok := m.GetCallback(actionID); ok {
					cb()
				}
			})
		}
	}
	return menu
}

// RegisterCallback registers a callback for a menu action ID.
func (m *Manager) RegisterCallback(actionID string, callback func()) {
	m.mu.Lock()
	m.callbacks[actionID] = callback
	m.mu.Unlock()
}

// UnregisterCallback removes a callback.
func (m *Manager) UnregisterCallback(actionID string) {
	m.mu.Lock()
	delete(m.callbacks, actionID)
	m.mu.Unlock()
}

// GetCallback returns the callback for an action ID.
func (m *Manager) GetCallback(actionID string) (func(), bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cb, ok := m.callbacks[actionID]
	return cb, ok
}

// GetInfo returns tray status information.
func (m *Manager) GetInfo() map[string]any {
	return map[string]any{
		"active": m.IsActive(),
	}
}
```

- [ ] **Step 6: Write mock_test.go**

```go
// pkg/systray/mock_test.go
package systray

type mockPlatform struct {
	trays []*mockTray
	menus []*mockTrayMenu
}

func newMockPlatform() *mockPlatform { return &mockPlatform{} }

func (p *mockPlatform) NewTray() PlatformTray {
	t := &mockTray{}
	p.trays = append(p.trays, t)
	return t
}

func (p *mockPlatform) NewMenu() PlatformMenu {
	m := &mockTrayMenu{}
	p.menus = append(p.menus, m)
	return m
}

type mockTrayMenu struct {
	items []string
}

func (m *mockTrayMenu) Add(label string) PlatformMenuItem { m.items = append(m.items, label); return &mockTrayMenuItem{} }
func (m *mockTrayMenu) AddSeparator()                     { m.items = append(m.items, "---") }

type mockTrayMenuItem struct{}

func (mi *mockTrayMenuItem) SetTooltip(text string)   {}
func (mi *mockTrayMenuItem) SetChecked(checked bool)   {}
func (mi *mockTrayMenuItem) SetEnabled(enabled bool)   {}
func (mi *mockTrayMenuItem) OnClick(fn func())         {}
func (mi *mockTrayMenuItem) AddSubmenu() PlatformMenu  { return &mockTrayMenu{} }

type mockTray struct {
	icon, templateIcon []byte
	tooltip, label     string
	menu               PlatformMenu
	attachedWindow     WindowHandle
}

func (t *mockTray) SetIcon(data []byte)                { t.icon = data }
func (t *mockTray) SetTemplateIcon(data []byte)         { t.templateIcon = data }
func (t *mockTray) SetTooltip(text string)              { t.tooltip = text }
func (t *mockTray) SetLabel(text string)                { t.label = text }
func (t *mockTray) SetMenu(menu PlatformMenu)           { t.menu = menu }
func (t *mockTray) AttachWindow(w WindowHandle)         { t.attachedWindow = w }
```

- [ ] **Step 7: Write tray_test.go**

```go
// pkg/systray/tray_test.go
package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Setup_Good(t *testing.T) {
	m, p := newTestManager()
	err := m.Setup("Core", "Core")
	require.NoError(t, err)
	assert.True(t, m.IsActive())
	assert.Len(t, p.trays, 1)
	assert.Equal(t, "Core", p.trays[0].tooltip)
	assert.Equal(t, "Core", p.trays[0].label)
	assert.NotEmpty(t, p.trays[0].templateIcon) // default icon embedded
}

func TestManager_SetIcon_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	err := m.SetIcon([]byte{1, 2, 3})
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, p.trays[0].icon)
}

func TestManager_SetIcon_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.SetIcon([]byte{1})
	assert.Error(t, err) // tray not initialised
}

func TestManager_SetTooltip_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetTooltip("New Tooltip")
	assert.Equal(t, "New Tooltip", p.trays[0].tooltip)
}

func TestManager_SetLabel_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetLabel("New Label")
	assert.Equal(t, "New Label", p.trays[0].label)
}

func TestManager_RegisterCallback_Good(t *testing.T) {
	m, _ := newTestManager()
	called := false
	m.RegisterCallback("test-action", func() { called = true })
	cb, ok := m.GetCallback("test-action")
	assert.True(t, ok)
	cb()
	assert.True(t, called)
}

func TestManager_RegisterCallback_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, ok := m.GetCallback("nonexistent")
	assert.False(t, ok)
}

func TestManager_UnregisterCallback_Good(t *testing.T) {
	m, _ := newTestManager()
	m.RegisterCallback("remove-me", func() {})
	m.UnregisterCallback("remove-me")
	_, ok := m.GetCallback("remove-me")
	assert.False(t, ok)
}

func TestManager_GetInfo_Good(t *testing.T) {
	m, _ := newTestManager()
	info := m.GetInfo()
	assert.False(t, info["active"].(bool))
	_ = m.Setup("Core", "Core")
	info = m.GetInfo()
	assert.True(t, info["active"].(bool))
}
```

- [ ] **Step 8: Write wails.go adapter**

```go
// pkg/systray/wails.go
package systray

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) NewTray() PlatformTray {
	return &wailsTray{tray: wp.app.NewSystemTray(), app: wp.app}
}

type wailsTray struct {
	tray *application.SystemTray
	app  *application.App
}

func (wt *wailsTray) SetIcon(data []byte)        { wt.tray.SetIcon(data) }
func (wt *wailsTray) SetTemplateIcon(data []byte) { wt.tray.SetTemplateIcon(data) }
func (wt *wailsTray) SetTooltip(text string)      { wt.tray.SetTooltip(text) }
func (wt *wailsTray) SetLabel(text string)        { wt.tray.SetLabel(text) }

func (wt *wailsTray) SetMenu(menu PlatformMenu) {
	// Menu constructed via Wails application.Menu — adapt as needed
}

func (wt *wailsTray) AttachWindow(w WindowHandle) {
	// Wails systray can attach a window — adapt based on v3 API
}
```

- [ ] **Step 9: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/systray/... -v -count=1`
Expected: ALL PASS (no `pkg/display` import needed — `WindowHandle` is local)

- [ ] **Step 10: Commit**

```bash
git add pkg/systray/ pkg/display/types.go
git commit -m "feat(systray): add Manager with platform abstraction and callback registry"
```

---

### Task 8: pkg/menu — Platform and builder

**Files:**
- Create: `pkg/menu/platform.go`
- Create: `pkg/menu/menu.go`
- Create: `pkg/menu/mock_test.go`
- Create: `pkg/menu/menu_test.go`
- Create: `pkg/menu/wails.go`

- [ ] **Step 1: Write platform.go**

```go
// pkg/menu/platform.go
package menu

// Platform abstracts the menu backend.
type Platform interface {
	NewMenu() PlatformMenu
	SetApplicationMenu(menu PlatformMenu)
}

// PlatformMenu is a live menu handle.
type PlatformMenu interface {
	Add(label string) PlatformMenuItem
	AddSeparator()
	AddSubmenu(label string) PlatformMenu
	// Roles — macOS menu roles
	AddRole(role MenuRole)
}

// PlatformMenuItem is a single menu item.
type PlatformMenuItem interface {
	SetAccelerator(accel string) PlatformMenuItem
	SetTooltip(text string) PlatformMenuItem
	SetChecked(checked bool) PlatformMenuItem
	SetEnabled(enabled bool) PlatformMenuItem
	OnClick(fn func()) PlatformMenuItem
}

// MenuRole is a predefined platform menu role.
type MenuRole int

const (
	RoleAppMenu MenuRole = iota
	RoleFileMenu
	RoleEditMenu
	RoleViewMenu
	RoleWindowMenu
	RoleHelpMenu
)
```

- [ ] **Step 2: Write menu.go — Manager + MenuItem (structure only)**

```go
// pkg/menu/menu.go
package menu

// MenuItem describes a menu item for construction (structure only — no handlers).
type MenuItem struct {
	Label       string
	Accelerator string
	Type        string // "normal", "separator", "checkbox", "radio", "submenu"
	Checked     bool
	Disabled    bool
	Tooltip     string
	Children    []MenuItem
	Role        *MenuRole
	OnClick     func() // Injected by orchestrator, not by menu package consumer
}

// Manager builds application menus via a Platform backend.
type Manager struct {
	platform Platform
}

// NewManager creates a menu Manager.
func NewManager(platform Platform) *Manager {
	return &Manager{platform: platform}
}

// Build constructs a PlatformMenu from a tree of MenuItems.
func (m *Manager) Build(items []MenuItem) PlatformMenu {
	menu := m.platform.NewMenu()
	m.buildItems(menu, items)
	return menu
}

func (m *Manager) buildItems(menu PlatformMenu, items []MenuItem) {
	for _, item := range items {
		if item.Role != nil {
			menu.AddRole(*item.Role)
			continue
		}
		if item.Type == "separator" {
			menu.AddSeparator()
			continue
		}
		if len(item.Children) > 0 {
			sub := menu.AddSubmenu(item.Label)
			m.buildItems(sub, item.Children)
			continue
		}
		mi := menu.Add(item.Label)
		if item.Accelerator != "" {
			mi.SetAccelerator(item.Accelerator)
		}
		if item.Tooltip != "" {
			mi.SetTooltip(item.Tooltip)
		}
		if item.OnClick != nil {
			mi.OnClick(item.OnClick)
		}
	}
}

// SetApplicationMenu builds and sets the application menu.
func (m *Manager) SetApplicationMenu(items []MenuItem) {
	menu := m.Build(items)
	m.platform.SetApplicationMenu(menu)
}

// Platform returns the underlying platform.
func (m *Manager) Platform() Platform {
	return m.platform
}
```

- [ ] **Step 3: Write mock_test.go**

```go
// pkg/menu/mock_test.go
package menu

type mockPlatform struct {
	menus    []*mockMenu
	appMenu  PlatformMenu
}

func newMockPlatform() *mockPlatform { return &mockPlatform{} }

func (p *mockPlatform) NewMenu() PlatformMenu {
	m := &mockMenu{}
	p.menus = append(p.menus, m)
	return m
}

func (p *mockPlatform) SetApplicationMenu(menu PlatformMenu) { p.appMenu = menu }

type mockMenu struct {
	items []*mockMenuItem
	subs  []*mockMenu
	roles []MenuRole
}

func (m *mockMenu) Add(label string) PlatformMenuItem {
	mi := &mockMenuItem{label: label}
	m.items = append(m.items, mi)
	return mi
}

func (m *mockMenu) AddSeparator() {
	m.items = append(m.items, &mockMenuItem{label: "---"})
}

func (m *mockMenu) AddSubmenu(label string) PlatformMenu {
	sub := &mockMenu{}
	m.subs = append(m.subs, sub)
	m.items = append(m.items, &mockMenuItem{label: label, isSubmenu: true})
	return sub
}

func (m *mockMenu) AddRole(role MenuRole) { m.roles = append(m.roles, role) }

type mockMenuItem struct {
	label, accel, tooltip string
	checked, enabled      bool
	isSubmenu             bool
	onClick               func()
}

func (mi *mockMenuItem) SetAccelerator(accel string) PlatformMenuItem { mi.accel = accel; return mi }
func (mi *mockMenuItem) SetTooltip(text string) PlatformMenuItem      { mi.tooltip = text; return mi }
func (mi *mockMenuItem) SetChecked(checked bool) PlatformMenuItem     { mi.checked = checked; return mi }
func (mi *mockMenuItem) SetEnabled(enabled bool) PlatformMenuItem     { mi.enabled = enabled; return mi }
func (mi *mockMenuItem) OnClick(fn func()) PlatformMenuItem           { mi.onClick = fn; return mi }
```

- [ ] **Step 4: Write menu_test.go**

```go
// pkg/menu/menu_test.go
package menu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Build_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "File"},
		{Label: "Edit"},
	}
	menu := m.Build(items)
	assert.NotNil(t, menu)
	assert.Len(t, p.menus, 1)
	assert.Len(t, p.menus[0].items, 2)
	assert.Equal(t, "File", p.menus[0].items[0].label)
}

func TestManager_Build_Separator_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Above"},
		{Type: "separator"},
		{Label: "Below"},
	}
	m.Build(items)
	assert.Len(t, p.menus[0].items, 3)
	assert.Equal(t, "---", p.menus[0].items[1].label)
}

func TestManager_Build_Submenu_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Parent", Children: []MenuItem{
			{Label: "Child 1"},
			{Label: "Child 2"},
		}},
	}
	m.Build(items)
	assert.Len(t, p.menus[0].subs, 1)
	assert.Len(t, p.menus[0].subs[0].items, 2)
}

func TestManager_Build_Accelerator_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{
		{Label: "Save", Accelerator: "CmdOrCtrl+S"},
	}
	m.Build(items)
	assert.Equal(t, "CmdOrCtrl+S", p.menus[0].items[0].accel)
}

func TestManager_Build_OnClick_Good(t *testing.T) {
	m, p := newTestManager()
	called := false
	items := []MenuItem{
		{Label: "Action", OnClick: func() { called = true }},
	}
	m.Build(items)
	p.menus[0].items[0].onClick()
	assert.True(t, called)
}

func TestManager_Build_Role_Good(t *testing.T) {
	m, p := newTestManager()
	appMenu := RoleAppMenu
	items := []MenuItem{
		{Role: &appMenu},
	}
	m.Build(items)
	assert.Contains(t, p.menus[0].roles, RoleAppMenu)
}

func TestManager_SetApplicationMenu_Good(t *testing.T) {
	m, p := newTestManager()
	items := []MenuItem{{Label: "Test"}}
	m.SetApplicationMenu(items)
	assert.NotNil(t, p.appMenu)
}

func TestManager_Build_Empty_Good(t *testing.T) {
	m, _ := newTestManager()
	menu := m.Build(nil)
	assert.NotNil(t, menu)
}
```

- [ ] **Step 5: Write wails.go adapter**

```go
// pkg/menu/wails.go
package menu

import "github.com/wailsapp/wails/v3/pkg/application"

// WailsPlatform implements Platform using Wails v3.
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

func (wp *WailsPlatform) NewMenu() PlatformMenu {
	return &wailsMenu{menu: application.NewMenu()}
}

func (wp *WailsPlatform) SetApplicationMenu(menu PlatformMenu) {
	if wm, ok := menu.(*wailsMenu); ok {
		wp.app.SetMenu(wm.menu)
	}
}

type wailsMenu struct {
	menu *application.Menu
}

func (wm *wailsMenu) Add(label string) PlatformMenuItem {
	return &wailsMenuItem{item: wm.menu.Add(label)}
}

func (wm *wailsMenu) AddSeparator() {
	wm.menu.AddSeparator()
}

func (wm *wailsMenu) AddSubmenu(label string) PlatformMenu {
	sub := wm.menu.AddSubmenu(label)
	return &wailsMenu{menu: sub}
}

func (wm *wailsMenu) AddRole(role MenuRole) {
	switch role {
	case RoleAppMenu:
		wm.menu.AddRole(application.AppMenu)
	case RoleFileMenu:
		wm.menu.AddRole(application.FileMenu)
	case RoleEditMenu:
		wm.menu.AddRole(application.EditMenu)
	case RoleViewMenu:
		wm.menu.AddRole(application.ViewMenu)
	case RoleWindowMenu:
		wm.menu.AddRole(application.WindowMenu)
	case RoleHelpMenu:
		wm.menu.AddRole(application.HelpMenu)
	}
}

type wailsMenuItem struct {
	item *application.MenuItem
}

func (mi *wailsMenuItem) SetAccelerator(accel string) PlatformMenuItem {
	mi.item.SetAccelerator(accel)
	return mi
}

func (mi *wailsMenuItem) SetTooltip(text string) PlatformMenuItem {
	mi.item.SetTooltip(text)
	return mi
}

func (mi *wailsMenuItem) SetChecked(checked bool) PlatformMenuItem {
	mi.item.SetChecked(checked)
	return mi
}

func (mi *wailsMenuItem) SetEnabled(enabled bool) PlatformMenuItem {
	if enabled {
		mi.item.SetEnabled(true)
	} else {
		mi.item.SetEnabled(false)
	}
	return mi
}

func (mi *wailsMenuItem) OnClick(fn func()) PlatformMenuItem {
	mi.item.OnClick(func(*application.Context) { fn() })
	return mi
}
```

- [ ] **Step 6: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/menu/... -v -count=1`
Expected: ALL PASS

- [ ] **Step 7: Commit**

```bash
git add pkg/menu/
git commit -m "feat(menu): add Manager with platform abstraction and builder"
```

---

## Chunk 3: pkg/display refactor

### Task 9: Shared types in pkg/display

**Files:**
- Modify: `pkg/display/types.go` (expand stub from Task 7)

- [ ] **Step 1: Expand types.go**

```go
// pkg/display/types.go
package display

// WindowHandle provides a cross-package interface for window operations.
// Both pkg/window.PlatformWindow and pkg/systray use this — no peer imports.
type WindowHandle interface {
	Name() string
	Show()
	Hide()
	SetPosition(x, y int)
	SetSize(width, height int)
}

// ScreenInfo describes a display screen.
type ScreenInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	IsPrimary bool   `json:"isPrimary"`
}

// WorkArea describes the usable area of a screen (excluding dock/menubar).
type WorkArea struct {
	ScreenID string `json:"screenId"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// EventSource abstracts the application event system (Wails insulation for WSEventManager).
// WSEventManager receives this instead of calling application.Get() directly.
type EventSource interface {
	OnThemeChange(handler func(isDark bool)) func()
	Emit(name string, data ...any) bool
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/display/types.go
git commit -m "feat(display): add shared types (WindowHandle, ScreenInfo, WorkArea, EventSource)"
```

---

### Task 10: Refactor pkg/display — Orchestrator

This is the largest task. The orchestrator composes the three sub-managers and delegates.

**Files:**
- Modify: `pkg/display/display.go` — Remove ~800 LOC of window CRUD/tiling/snapping, replace with delegation
- Modify: `pkg/display/interfaces.go` — Remove migrated interfaces
- Modify: `pkg/display/events.go` — Accept `window.PlatformWindow` instead of Wails types
- Modify: `pkg/display/actions.go` — Use `window.Window` instead of Wails type
- Delete: `pkg/display/window.go` — Replaced by `pkg/window/options.go`
- Delete: `pkg/display/window_state.go` — Replaced by `pkg/window/state.go`
- Delete: `pkg/display/layout.go` — Replaced by `pkg/window/layout.go`
- Delete: `pkg/display/tray.go` — Replaced by `pkg/systray/tray.go`
- Delete: `pkg/display/menu.go` — Handlers stay in display.go, structure in `pkg/menu`
- Modify: `pkg/display/display_test.go` — Update for new imports
- Modify: `pkg/display/mocks_test.go` — Remove migrated mocks

**Strategy:** This task is large but mechanical. The engineer should:
1. Delete the old files first
2. Update `display.go` to compose sub-managers
3. Update imports and types
4. Fix tests

- [ ] **Step 1: Delete replaced files**

```bash
cd /Users/snider/Code/core/gui
rm pkg/display/window.go
rm pkg/display/window_state.go
rm pkg/display/layout.go
rm pkg/display/tray.go
rm pkg/display/menu.go
```

- [ ] **Step 2: Update actions.go**

Replace the `ActionOpenWindow` struct to use `window.Window`:

```go
// pkg/display/actions.go
package display

import "forge.lthn.ai/core/gui/pkg/window"

// ActionOpenWindow is an IPC message type requesting a new window.
type ActionOpenWindow struct {
	window.Window
}
```

- [ ] **Step 3: Update interfaces.go — Keep only display-level interfaces**

Remove `WindowManager`, `MenuManager`, `SystemTrayManager` and their Wails adapters. Keep `DialogManager`, `EnvManager`, `EventManager`, `Logger`, and the `App` interface reduced to what display still needs directly:

```go
// pkg/display/interfaces.go
package display

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// App abstracts the Wails application for the orchestrator.
type App interface {
	Dialog() DialogManager
	Env() EnvManager
	Event() EventManager
	Logger() Logger
	Quit()
}

// DialogManager wraps Wails dialog operations.
type DialogManager interface {
	Info() *application.MessageDialog
	Warning() *application.MessageDialog
	OpenFile() *application.OpenFileDialogStruct
}

// EnvManager wraps Wails environment queries.
type EnvManager interface {
	Info() application.EnvironmentInfo
	IsDarkMode() bool
}

// EventManager wraps Wails application events.
type EventManager interface {
	OnApplicationEvent(eventType events.ApplicationEventType, handler func(*application.ApplicationEvent)) func()
	Emit(name string, data ...any) bool
}

// Logger wraps Wails logging.
type Logger interface {
	Info(message string, args ...any)
}

// wailsApp wraps *application.App for the App interface.
type wailsApp struct {
	app *application.App
}

func newWailsApp(app *application.App) *wailsApp {
	return &wailsApp{app: app}
}

func (w *wailsApp) Dialog() DialogManager   { return &wailsDialogManager{app: w.app} }
func (w *wailsApp) Env() EnvManager         { return &wailsEnvManager{app: w.app} }
func (w *wailsApp) Event() EventManager     { return &wailsEventManager{app: w.app} }
func (w *wailsApp) Logger() Logger           { return &wailsLogger{app: w.app} }
func (w *wailsApp) Quit()                   { w.app.Quit() }

type wailsDialogManager struct{ app *application.App }
func (d *wailsDialogManager) Info() *application.MessageDialog       { return d.app.InfoDialog() }
func (d *wailsDialogManager) Warning() *application.MessageDialog    { return d.app.WarningDialog() }
func (d *wailsDialogManager) OpenFile() *application.OpenFileDialogStruct { return d.app.OpenFileDialogWithOptions() }

type wailsEnvManager struct{ app *application.App }
func (e *wailsEnvManager) Info() application.EnvironmentInfo { return e.app.Info() }
func (e *wailsEnvManager) IsDarkMode() bool                 { return e.app.IsDarkMode() }

type wailsEventManager struct{ app *application.App }
func (ev *wailsEventManager) OnApplicationEvent(eventType events.ApplicationEventType, handler func(*application.ApplicationEvent)) func() {
	return ev.app.OnApplicationEvent(eventType, handler)
}
func (ev *wailsEventManager) Emit(name string, data ...any) bool { return ev.app.EmitEvent(name, data...) }

type wailsLogger struct{ app *application.App }
func (l *wailsLogger) Info(message string, args ...any) { l.app.Logger.Info(message, args...) }
```

- [ ] **Step 4: Refactor display.go — Compose sub-managers**

The `Service` struct changes from managing windows directly to delegating to `window.Manager`, `systray.Manager`, and `menu.Manager`.

```go
// pkg/display/display.go — updated Service struct and key methods

import (
	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
)

type Service struct {
	*core.ServiceRuntime[Options]
	app         App
	windows     *window.Manager
	tray        *systray.Manager
	menus       *menu.Manager
	events      *WSEventManager
	eventSource EventSource
}

// New creates an unregistered Service.
func New() (*Service, error) {
	return &Service{}, nil
}

// Register creates a Service bound to Core DI.
func Register(c *core.Core) (any, error) {
	s := &Service{}
	s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
	return s, nil
}

// ServiceStartup initialises sub-managers with the Wails app.
func (s *Service) ServiceStartup(app any) {
	// Cast to *application.App, create platform adapters
	// wailsApp := app.(*application.App)
	// s.windows = window.NewManager(window.NewWailsPlatform(wailsApp))
	// s.tray = systray.NewManager(systray.NewWailsPlatform(wailsApp))
	// s.menus = menu.NewManager(menu.NewWailsPlatform(wailsApp))
	// s.app = newWailsApp(wailsApp)
	// s.events = NewWSEventManager(s.eventSource)
	// s.buildMenu()
	// s.setupTray()
}

// --- Public API delegates to window.Manager ---

func (s *Service) OpenWindow(opts ...window.WindowOption) error {
	pw, err := s.windows.Open(opts...)
	if err != nil {
		return err
	}
	s.trackWindow(pw)
	return nil
}

func (s *Service) SetWindowPosition(name string, x, y int) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}
	pw.SetPosition(x, y)
	s.windows.State().UpdatePosition(name, x, y)
	return nil
}

func (s *Service) SetWindowSize(name string, width, height int) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}
	pw.SetSize(width, height)
	s.windows.State().UpdateSize(name, width, height)
	return nil
}

func (s *Service) MaximizeWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}
	pw.Maximise()
	s.windows.State().UpdateMaximized(name, true)
	return nil
}

func (s *Service) FocusWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}
	pw.Focus()
	return nil
}

func (s *Service) CloseWindow(name string) error {
	pw, ok := s.windows.Get(name)
	if !ok {
		return fmt.Errorf("window %q not found", name)
	}
	s.windows.State().CaptureState(pw)
	pw.Close()
	s.windows.Remove(name)
	return nil
}

// --- Layout delegation ---

func (s *Service) SaveLayout(name string) error {
	states := make(map[string]window.WindowState)
	for _, n := range s.windows.List() {
		if pw, ok := s.windows.Get(n); ok {
			x, y := pw.Position()
			w, h := pw.Size()
			states[n] = window.WindowState{X: x, Y: y, Width: w, Height: h, Maximized: pw.IsMaximised()}
		}
	}
	return s.windows.Layout().SaveLayout(name, states)
}

func (s *Service) RestoreLayout(name string) error {
	layout, ok := s.windows.Layout().GetLayout(name)
	if !ok {
		return fmt.Errorf("layout %q not found", name)
	}
	for wName, state := range layout.Windows {
		if pw, ok := s.windows.Get(wName); ok {
			pw.SetPosition(state.X, state.Y)
			pw.SetSize(state.Width, state.Height)
			if state.Maximized {
				pw.Maximise()
			}
		}
	}
	return nil
}

// --- Tiling/snapping delegation ---

func (s *Service) TileWindows(mode window.TileMode, names []string) error {
	// Use primary screen dimensions — screen queries remain in display
	return s.windows.TileWindows(mode, names, 1920, 1080) // TODO: use actual screen size
}

func (s *Service) SnapWindow(name string, position window.SnapPosition) error {
	return s.windows.SnapWindow(name, position, 1920, 1080) // TODO: use actual screen size
}

// --- trackWindow attaches event listeners for state persistence ---

func (s *Service) trackWindow(pw window.PlatformWindow) {
	s.events.AttachWindowListeners(pw)
	s.events.EmitWindowEvent(EventWindowCreate, pw.Name(), nil)
}

// --- setupTray delegates to systray.Manager ---

func (s *Service) setupTray() {
	_ = s.tray.Setup("Core", "Core")
	s.tray.RegisterCallback("open-desktop", func() {
		for _, name := range s.windows.List() {
			if pw, ok := s.windows.Get(name); ok {
				pw.Show()
			}
		}
	})
	s.tray.RegisterCallback("close-desktop", func() {
		for _, name := range s.windows.List() {
			if pw, ok := s.windows.Get(name); ok {
				pw.Hide()
			}
		}
	})
	s.tray.RegisterCallback("env-info", func() { s.ShowEnvironmentDialog() })
	s.tray.RegisterCallback("quit", func() { s.app.Quit() })
	_ = s.tray.SetMenu([]systray.TrayMenuItem{
		{Label: "Open Desktop", ActionID: "open-desktop"},
		{Label: "Close Desktop", ActionID: "close-desktop"},
		{Type: "separator"},
		{Label: "Environment Info", ActionID: "env-info"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	})
}

// --- Handler methods (stay in display — use s.windows.Open) ---

func (s *Service) handleNewWorkspace() {
	_ = s.OpenWindow(window.WithName("workspace-new"), window.WithTitle("New Workspace"),
		window.WithURL("/workspace/new"), window.WithSize(500, 400))
}

func (s *Service) handleNewFile() {
	_ = s.OpenWindow(window.WithName("editor-new"), window.WithTitle("New File"),
		window.WithURL("/developer/editor?new=true"), window.WithSize(1200, 800))
}

func (s *Service) handleOpenFile() {
	// File dialog → open editor with file path
	// Uses s.app.Dialog().OpenFile() which stays in display
}

func (s *Service) handleSaveFile()     { s.app.Event().Emit("ide:save") }
func (s *Service) handleOpenEditor()   {
	_ = s.OpenWindow(window.WithName("editor"), window.WithTitle("Editor"),
		window.WithURL("/developer/editor"), window.WithSize(1200, 800))
}
func (s *Service) handleOpenTerminal() {
	_ = s.OpenWindow(window.WithName("terminal"), window.WithTitle("Terminal"),
		window.WithURL("/developer/terminal"), window.WithSize(800, 500))
}
func (s *Service) handleRun()   { s.app.Event().Emit("ide:run") }
func (s *Service) handleBuild() { s.app.Event().Emit("ide:build") }

func ptr[T any](v T) *T { return &v }
```

**Note on screen queries:** Methods like `GetScreens()`, `GetWorkAreas()`, `GetPrimaryScreen()` currently call `application.Get()` directly. These stay in `pkg/display` and will be insulated in a follow-up task via a `ScreenProvider` interface. For now they remain as direct Wails calls — the priority is splitting window/systray/menu cleanly.

**Key pattern for menu handlers staying in display:**

```go
func (s *Service) buildMenu() {
	items := []menu.MenuItem{
		{Role: ptr(menu.RoleAppMenu)},
		{Role: ptr(menu.RoleFileMenu)},
		{Label: "Workspace", Children: []menu.MenuItem{
			{Label: "New...", OnClick: s.handleNewWorkspace},
			{Label: "List", OnClick: s.handleListWorkspaces},
		}},
		{Label: "Developer", Children: []menu.MenuItem{
			{Label: "New File", Accelerator: "CmdOrCtrl+N", OnClick: s.handleNewFile},
			{Label: "Open File...", Accelerator: "CmdOrCtrl+O", OnClick: s.handleOpenFile},
			{Label: "Save", Accelerator: "CmdOrCtrl+S", OnClick: s.handleSaveFile},
			{Type: "separator"},
			{Label: "Editor", OnClick: s.handleOpenEditor},
			{Label: "Terminal", OnClick: s.handleOpenTerminal},
			{Type: "separator"},
			{Label: "Run", Accelerator: "CmdOrCtrl+R", OnClick: s.handleRun},
			{Label: "Build", Accelerator: "CmdOrCtrl+B", OnClick: s.handleBuild},
		}},
		{Role: ptr(menu.RoleEditMenu)},
		{Role: ptr(menu.RoleViewMenu)},
		{Role: ptr(menu.RoleWindowMenu)},
		{Role: ptr(menu.RoleHelpMenu)},
	}
	s.menus.SetApplicationMenu(items)
}

func ptr[T any](v T) *T { return &v }
```

Handler methods (`handleNewWorkspace`, `handleOpenFile`, etc.) stay in `display.go` — they use `s.windows.Open(...)` to create windows and `s.app.Event().Emit(...)` for IDE events.

- [ ] **Step 5: Update events.go — PlatformWindow + EventSource insulation**

1. `AttachWindowListeners` accepts `window.PlatformWindow` instead of Wails concrete type
2. `SetupWindowEventListeners` uses `EventSource` instead of `application.Get()` directly
3. `NewWSEventManager` accepts `EventSource`

```go
// events.go — key changes (keep existing WebSocket logic intact)

import "forge.lthn.ai/core/gui/pkg/window"

// NewWSEventManager now accepts an EventSource for theme change events.
func NewWSEventManager(es EventSource) *WSEventManager {
	em := &WSEventManager{
		eventSource: es,
		// ... existing fields ...
	}
	return em
}

// AttachWindowListeners accepts PlatformWindow (not *application.WebviewWindow).
func (em *WSEventManager) AttachWindowListeners(pw window.PlatformWindow) {
	pw.OnWindowEvent(func(e window.WindowEvent) {
		em.EmitWindowEvent(EventType(e.Type), e.Name, e.Data)
	})
}

// SetupWindowEventListeners uses EventSource (not application.Get()).
func (em *WSEventManager) SetupWindowEventListeners() {
	if em.eventSource != nil {
		em.eventSource.OnThemeChange(func(isDark bool) {
			theme := "light"
			if isDark {
				theme = "dark"
			}
			em.EmitWindowEvent(EventThemeChange, "", map[string]any{"theme": theme})
		})
	}
}
```

- [ ] **Step 6: Update display_test.go and mocks_test.go**

Remove tests that moved to sub-packages (window option tests, tray tests). Keep orchestrator-level tests. Update mocks to compose sub-package mocks:

```go
// mocks_test.go — simplified
type mockApp struct {
	dialogManager *mockDialogManager
	envManager    *mockEnvManager
	eventManager  *mockEventManager
	logger        *mockLogger
	quitCalled    bool
}
// ... only dialog/env/event/logger mocks remain
```

- [ ] **Step 7: Run all tests**

Run: `cd /Users/snider/Code/core/gui && go test ./... -v -count=1`
Expected: ALL PASS across all 4 packages

- [ ] **Step 8: Commit**

```bash
git add pkg/display/display.go pkg/display/interfaces.go pkg/display/events.go pkg/display/actions.go pkg/display/types.go pkg/display/display_test.go pkg/display/mocks_test.go
git rm pkg/display/window.go pkg/display/window_state.go pkg/display/layout.go pkg/display/tray.go pkg/display/menu.go
git commit -m "refactor(display): compose window/systray/menu sub-packages into orchestrator"
```

---

## Chunk 4: ui/ move and final verification

### Task 11: Move ui/ to top level

**Files:**
- Move: `pkg/display/ui/` → `ui/`
- Modify: Any `go:embed` directives referencing `ui/`

- [ ] **Step 1: Move ui/ directory**

```bash
cd /Users/snider/Code/core/gui
mv pkg/display/ui ui
```

- [ ] **Step 2: Check for go:embed references**

Search for any `go:embed` directives in pkg/display/ that reference `ui/`:

```bash
grep -r "go:embed" pkg/display/
```

If found, these likely embed the Angular build output. Update paths from `ui/dist` to `../../ui/dist` or move the embed directive to a top-level file.

- [ ] **Step 3: Verify Angular project is intact**

```bash
ls ui/package.json ui/angular.json ui/src/main.ts
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor: move ui/ demo to top level"
```

---

### Task 12: Final verification

- [ ] **Step 1: Build all packages**

```bash
cd /Users/snider/Code/core/gui && go build ./...
```

- [ ] **Step 2: Run all tests**

```bash
cd /Users/snider/Code/core/gui && go test ./... -v -count=1
```

- [ ] **Step 3: Run go vet**

```bash
cd /Users/snider/Code/core/gui && go vet ./...
```

- [ ] **Step 4: Verify no circular dependencies**

```bash
cd /Users/snider/Code/core/gui && go list -f '{{.ImportPath}}: {{join .Imports "\n  "}}' ./pkg/window/ ./pkg/systray/ ./pkg/menu/ ./pkg/display/
```

Verify:
- `pkg/window` does NOT import `pkg/display`, `pkg/systray`, or `pkg/menu`
- `pkg/systray` imports `pkg/display` (for WindowHandle) but NOT `pkg/window` or `pkg/menu`
- `pkg/menu` does NOT import `pkg/display`, `pkg/window`, or `pkg/systray`
- `pkg/display` imports `pkg/window`, `pkg/systray`, `pkg/menu`

- [ ] **Step 5: Verify workspace builds**

```bash
cd /Users/snider/Code && go build ./...
```

- [ ] **Step 6: Commit and push**

```bash
cd /Users/snider/Code/core/gui
git add -A
git commit -m "chore: final verification after display package split"
git push origin main
```

---

## Breaking API Changes

This split changes the public API. Downstream consumers (LEM, Mining, IDE) will need updates:

| Old (pkg/display) | New | Notes |
|---|---|---|
| `type Window = application.WebviewWindowOptions` | `window.Window` (own struct) | No longer a Wails alias |
| `WindowOption func(*application.WebviewWindowOptions) error` | `window.WindowOption func(*window.Window) error` | Rewritten against CoreGUI's Window |
| `WindowName("x")` | `window.WithName("x")` | Renamed to `With*` prefix |
| `display.TileMode` (string) | `window.TileMode` (int iota) | Type changed |
| `display.SnapPosition` (string) | `window.SnapPosition` (int iota) | Type changed |
| `SetTrayMenu(items)` | `systray.Manager.SetMenu(items)` | Now on Manager |
| `RegisterTrayMenuCallback(id, fn)` | `systray.Manager.RegisterCallback(id, fn)` | Now on Manager |

Screen query methods (`GetScreens`, `GetWorkAreas`, etc.) remain on `pkg/display.Service` unchanged. Dialog, clipboard, notification, and theme APIs are unchanged.

## Deferred Work

- **Screen insulation:** `GetScreens()`, `GetWorkAreas()`, etc. still call `application.Get()` directly. A future `ScreenProvider` interface will complete the insulation.
- **Existing layouts.json / window_state.json:** JSON field naming is preserved (camelCase) for backward compatibility with existing persisted files.
- **Clipboard image/HTML:** Clipboard remains text-only; parsed types exist but aren't used.

## Key References

| File | Role |
|------|------|
| `docs/superpowers/specs/2026-03-13-display-package-split-design.md` | Approved design spec |
| `pkg/display/display.go` | Current monolith (1,294 LOC) |
| `pkg/display/interfaces.go` | Current Wails abstraction layer |
| `pkg/display/tray.go` | Current tray with package-level globals |
| `pkg/display/menu.go` | Current menu with embedded click handlers |
| `pkg/display/window_state.go` | Current state persistence |
| `pkg/display/display_test.go` | Existing 63 test cases |
| `pkg/display/mocks_test.go` | Existing mock infrastructure |
