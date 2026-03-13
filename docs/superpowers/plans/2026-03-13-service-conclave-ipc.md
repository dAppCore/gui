# Service Conclave IPC Integration — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Wire core/go's IPC bus (ACTION/QUERY/PERFORM) into core/gui's split packages, making each sub-package a full `core.Service` that communicates via typed messages.

**Architecture:** Three-layer stack — IPC Bus → Service (embeds `ServiceRuntime`, composes Manager) → Platform Interface (Wails adapters). Each sub-package exposes `Register(platform)` returning a factory closure for `WithService`. The display orchestrator owns config and bridges IPC actions to WebSocket events for TS apps.

**Tech Stack:** Go 1.26, `forge.lthn.ai/core/go` v0.2.2 (DI/IPC), `forge.lthn.ai/core/gui` (Wails v3 abstraction), testify (assert/require)

**Spec:** `docs/superpowers/specs/2026-03-13-service-conclave-ipc-design.md`

---

## File Structure

### New Files

| File | Responsibility |
|------|---------------|
| `pkg/window/messages.go` | Window IPC message types (queries, tasks, actions) + WindowInfo |
| `pkg/window/service.go` | Window Service: embeds ServiceRuntime, composes Manager, IPC handlers |
| `pkg/window/register.go` | `Register(Platform)` factory returning closure for `WithService` |
| `pkg/window/service_test.go` | Window Service IPC tests |
| `pkg/systray/messages.go` | Systray IPC message types |
| `pkg/systray/service.go` | Systray Service: embeds ServiceRuntime, composes Manager |
| `pkg/systray/register.go` | `Register(Platform)` factory |
| `pkg/systray/service_test.go` | Systray Service IPC tests |
| `pkg/menu/messages.go` | Menu IPC message types |
| `pkg/menu/service.go` | Menu Service: embeds ServiceRuntime, composes Manager |
| `pkg/menu/register.go` | `Register(Platform)` factory |
| `pkg/menu/service_test.go` | Menu Service IPC tests |

### Modified Files

| File | Changes |
|------|---------|
| `pkg/display/display.go` | Refactor `Register()` to closure pattern accepting `*application.App`, add `OnStartup` with config query/task handlers registered synchronously, add `HandleIPCEvents` for IPC→WS bridge, add `configData` map, convert delegation methods to IPC, remove `ServiceName()` |
| `pkg/display/actions.go` | Remove file (replaced by sub-package message types) |
| `pkg/display/events.go` | Add `EventTrayClick` constant |
| `pkg/display/display_test.go` | Add IPC conclave integration tests |

### Design Note: Circular Import Avoidance

The spec shows `display.QueryConfig{Key: "window"}` used by sub-services during OnStartup. This creates a circular import (window→display→window). Solution: each sub-package defines its own `QueryConfig` type. The display orchestrator's `handleQuery` switches on all sub-package query types (which it already imports per the dependency direction). Sub-services use their local `QueryConfig{}` — no display import needed.

---

## Chunk 1: Window + Systray IPC Layers

### Task 1: Window IPC Layer

**Files:**
- Create: `pkg/window/messages.go`
- Create: `pkg/window/service.go`
- Create: `pkg/window/register.go`
- Create: `pkg/window/service_test.go`

- [ ] **Step 1: Write failing test for Register + Service creation**

Create `pkg/window/service_test.go`:

```go
package window

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestWindowService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/window/ -v -run TestRegister_Good`
Expected: FAIL — `Register`, `Service` not defined

- [ ] **Step 3: Create messages.go with window message types + WindowInfo**

Create `pkg/window/messages.go`:

```go
package window

// WindowInfo contains information about a window.
type WindowInfo struct {
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Maximized bool   `json:"maximized"`
	Focused   bool   `json:"focused"`
}

// --- Queries (read-only) ---

// QueryWindowList returns all tracked windows. Result: []WindowInfo
type QueryWindowList struct{}

// QueryWindowByName returns a single window by name. Result: *WindowInfo (nil if not found)
type QueryWindowByName struct{ Name string }

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// --- Tasks (side-effects) ---

// TaskOpenWindow creates a new window. Result: WindowInfo
type TaskOpenWindow struct{ Opts []WindowOption }

// TaskCloseWindow closes a window. Handler persists state BEFORE emitting ActionWindowClosed.
type TaskCloseWindow struct{ Name string }

// TaskSetPosition moves a window.
type TaskSetPosition struct {
	Name string
	X, Y int
}

// TaskSetSize resizes a window.
type TaskSetSize struct {
	Name string
	W, H int
}

// TaskMaximise maximises a window.
type TaskMaximise struct{ Name string }

// TaskMinimise minimises a window.
type TaskMinimise struct{ Name string }

// TaskFocus brings a window to the front.
type TaskFocus struct{ Name string }

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions (broadcasts) ---

type ActionWindowOpened struct{ Name string }
type ActionWindowClosed struct{ Name string }

type ActionWindowMoved struct {
	Name string
	X, Y int
}

type ActionWindowResized struct {
	Name string
	W, H int
}

type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }
```

- [ ] **Step 4: Create register.go with factory closure**

Create `pkg/window/register.go`:

```go
package window

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
```

- [ ] **Step 5: Create service.go with Service struct and OnStartup**

Create `pkg/window/service.go`:

```go
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
	// Apply config to manager defaults — future expansion.
	// e.g., default_width, default_height, state_file path.
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
```

- [ ] **Step 6: Run test to verify it passes**

Run: `go test ./pkg/window/ -v -run TestRegister_Good`
Expected: PASS

- [ ] **Step 7: Write failing tests for TaskOpenWindow + QueryWindowList + QueryWindowByName**

Append to `pkg/window/service_test.go`:

```go
func TestTaskOpenWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{
		Opts: []WindowOption{WithName("test"), WithURL("/")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestTaskOpenWindow_Bad(t *testing.T) {
	// No window service registered — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskOpenWindow{})
	assert.False(t, handled)
}

func TestQueryWindowList_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("a")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("b")}})

	result, handled, err := c.QUERY(QueryWindowList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]WindowInfo)
	assert.Len(t, list, 2)
}

func TestQueryWindowByName_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	result, handled, err := c.QUERY(QueryWindowByName{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(*WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestQueryWindowByName_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.QUERY(QueryWindowByName{Name: "nonexistent"})
	require.NoError(t, err)
	assert.True(t, handled) // handled=true, result is nil (not found)
	assert.Nil(t, result)
}
```

- [ ] **Step 8: Run tests to verify they pass**

Run: `go test ./pkg/window/ -v -run "TestTask|TestQuery"`
Expected: PASS (all 5 tests)

- [ ] **Step 9: Write failing tests for TaskCloseWindow, TaskSetPosition, TaskSetSize**

Append to `pkg/window/service_test.go`:

```go
func TestTaskCloseWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify window is removed
	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	assert.Nil(t, result)
}

func TestTaskCloseWindow_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetPosition_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetPosition{Name: "test", X: 100, Y: 200})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
}

func TestTaskSetSize_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetSize{Name: "test", W: 800, H: 600})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestTaskMaximise_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskMaximise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.True(t, info.Maximized)
}
```

- [ ] **Step 10: Run tests to verify they pass**

Run: `go test ./pkg/window/ -v -run "TestTaskClose|TestTaskSet|TestTaskMax"`
Expected: PASS

- [ ] **Step 11: Commit**

```bash
git add pkg/window/messages.go pkg/window/service.go pkg/window/register.go pkg/window/service_test.go
git commit -m "feat(window): add IPC layer — Service, Register factory, message types

Window package is now a full core.Service with typed IPC messages.
Register(Platform) factory closure captures platform adapter for WithService.
OnStartup queries config and registers query/task handlers.
Platform events converted to IPC actions via trackWindow.

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

### Task 2: Systray IPC Layer

**Files:**
- Create: `pkg/systray/messages.go`
- Create: `pkg/systray/service.go`
- Create: `pkg/systray/register.go`
- Create: `pkg/systray/service_test.go`

- [ ] **Step 1: Write failing test for Register + Service creation**

Create `pkg/systray/service_test.go`:

```go
package systray

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSystrayService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "systray")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestSystrayService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/systray/ -v -run TestRegister_Good`
Expected: FAIL — `Register`, `Service` not defined

- [ ] **Step 3: Create messages.go with systray message types**

Create `pkg/systray/messages.go`:

```go
package systray

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// --- Tasks ---

// TaskSetTrayIcon sets the tray icon.
type TaskSetTrayIcon struct{ Data []byte }

// TaskSetTrayMenu sets the tray menu items.
type TaskSetTrayMenu struct{ Items []TrayMenuItem }

// TaskShowPanel shows the tray panel window.
type TaskShowPanel struct{}

// TaskHidePanel hides the tray panel window.
type TaskHidePanel struct{}

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions ---

// ActionTrayClicked is broadcast when the tray icon is clicked.
type ActionTrayClicked struct{}

// ActionTrayMenuItemClicked is broadcast when a tray menu item is clicked.
type ActionTrayMenuItemClicked struct{ ActionID string }
```

- [ ] **Step 4: Create register.go and service.go**

Create `pkg/systray/register.go`:

```go
package systray

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
```

Create `pkg/systray/service.go`:

```go
package systray

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the systray service.
type Options struct{}

// Service is a core.Service managing the system tray via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
}

// OnStartup queries config and registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	cfg, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if tCfg, ok := cfg.(map[string]any); ok {
			s.applyConfig(tCfg)
		}
	}
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(cfg map[string]any) {
	// Apply config — tooltip, icon path, etc.
	tooltip, _ := cfg["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	_ = s.manager.Setup(tooltip, tooltip)
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetTrayIcon:
		return nil, true, s.manager.SetIcon(t.Data)
	case TaskSetTrayMenu:
		return nil, true, s.taskSetTrayMenu(t)
	case TaskShowPanel:
		// Panel show — deferred (requires WindowHandle integration)
		return nil, true, nil
	case TaskHidePanel:
		// Panel hide — deferred (requires WindowHandle integration)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) taskSetTrayMenu(t TaskSetTrayMenu) error {
	// Register IPC-emitting callbacks for each menu item
	for _, item := range t.Items {
		if item.ActionID != "" {
			actionID := item.ActionID
			s.manager.RegisterCallback(actionID, func() {
				_ = s.Core().ACTION(ActionTrayMenuItemClicked{ActionID: actionID})
			})
		}
	}
	return s.manager.SetMenu(t.Items)
}

// Manager returns the underlying systray Manager.
func (s *Service) Manager() *Manager {
	return s.manager
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./pkg/systray/ -v -run TestRegister_Good`
Expected: PASS

- [ ] **Step 6: Write tests for TaskSetTrayIcon + TaskSetTrayMenu**

Append to `pkg/systray/service_test.go`:

```go
func TestTaskSetTrayIcon_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)

	// Setup tray first (normally done via config in OnStartup)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	icon := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	_, handled, err := c.PERFORM(TaskSetTrayIcon{Data: icon})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetTrayMenu_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)

	require.NoError(t, svc.manager.Setup("Test", "Test"))

	items := []TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}
	_, handled, err := c.PERFORM(TaskSetTrayMenu{Items: items})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetTrayIcon_Bad(t *testing.T) {
	// No systray service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetTrayIcon{Data: nil})
	assert.False(t, handled)
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./pkg/systray/ -v -run "TestTask"`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add pkg/systray/messages.go pkg/systray/service.go pkg/systray/register.go pkg/systray/service_test.go
git commit -m "feat(systray): add IPC layer — Service, Register factory, message types

Systray package is now a full core.Service with typed IPC messages.
Menu item clicks emit ActionTrayMenuItemClicked via IPC.

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

## Chunk 2: Menu IPC Layer + Display Config Refactor

### Task 3: Menu IPC Layer

**Files:**
- Create: `pkg/menu/messages.go`
- Create: `pkg/menu/service.go`
- Create: `pkg/menu/register.go`
- Create: `pkg/menu/service_test.go`

- [ ] **Step 1: Write failing test for Register + Service creation**

Create `pkg/menu/service_test.go`:

```go
package menu

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMenuService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "menu")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestMenuService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/menu/ -v -run TestRegister_Good`
Expected: FAIL

- [ ] **Step 3: Create messages.go with menu message types**

Create `pkg/menu/messages.go`:

```go
package menu

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// QueryGetAppMenu returns the current app menu item descriptors.
// Result: []MenuItem
type QueryGetAppMenu struct{}

// TaskSetAppMenu sets the application menu. OnClick closures work because
// core/go IPC is in-process (no serialisation boundary).
type TaskSetAppMenu struct{ Items []MenuItem }

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }
```

- [ ] **Step 4: Create register.go and service.go**

Create `pkg/menu/register.go`:

```go
package menu

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			manager:        NewManager(p),
		}, nil
	}
}
```

Create `pkg/menu/service.go`:

```go
package menu

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the menu service.
type Options struct{}

// Service is a core.Service managing application menus via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	manager  *Manager
	platform Platform
	items    []MenuItem // last-set menu items for QueryGetAppMenu
}

// OnStartup queries config and registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	cfg, handled, _ := s.Core().QUERY(QueryConfig{})
	if handled {
		if mCfg, ok := cfg.(map[string]any); ok {
			s.applyConfig(mCfg)
		}
	}
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

func (s *Service) applyConfig(cfg map[string]any) {
	// Apply config — e.g., show_dev_tools
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryGetAppMenu:
		return s.items, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskSetAppMenu:
		s.items = t.Items
		s.manager.SetApplicationMenu(t.Items)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

// Manager returns the underlying menu Manager.
func (s *Service) Manager() *Manager {
	return s.manager
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./pkg/menu/ -v -run TestRegister_Good`
Expected: PASS

- [ ] **Step 6: Write tests for TaskSetAppMenu + QueryGetAppMenu**

Append to `pkg/menu/service_test.go`:

```go
func TestTaskSetAppMenu_Good(t *testing.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{
		{Label: "File", Children: []MenuItem{
			{Label: "New"},
			{Type: "separator"},
			{Label: "Quit"},
		}},
	}
	_, handled, err := c.PERFORM(TaskSetAppMenu{Items: items})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestQueryGetAppMenu_Good(t *testing.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{{Label: "File"}, {Label: "Edit"}}
	_, _, _ = c.PERFORM(TaskSetAppMenu{Items: items})

	result, handled, err := c.QUERY(QueryGetAppMenu{})
	require.NoError(t, err)
	assert.True(t, handled)
	menu := result.([]MenuItem)
	assert.Len(t, menu, 2)
	assert.Equal(t, "File", menu[0].Label)
}

func TestTaskSetAppMenu_Bad(t *testing.T) {
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetAppMenu{})
	assert.False(t, handled)
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./pkg/menu/ -v -run "TestTask|TestQuery"`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add pkg/menu/messages.go pkg/menu/service.go pkg/menu/register.go pkg/menu/service_test.go
git commit -m "feat(menu): add IPC layer — Service, Register factory, message types

Menu package is now a full core.Service with typed IPC messages.
TaskSetAppMenu carries MenuItems with OnClick closures (in-process IPC).

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

### Task 4: Display Config Handlers & Register Refactor

**Files:**
- Modify: `pkg/display/display.go`
- Modify: `pkg/display/events.go`
- Remove: `pkg/display/actions.go`
- Modify: `pkg/display/display_test.go`

This task refactors the display service to:
1. Accept `*application.App` in Register (closure pattern)
2. Add `configData` for in-memory config storage
3. Register config query/task handlers synchronously in OnStartup
4. Remove `ServiceName()` (auto-derived by WithService as `"display"`)
5. Add `EventTrayClick` constant to events.go
6. Remove `ServiceStartup(ctx, options)` and `Startup(ctx)` methods — replaced by `OnStartup(ctx)` (the `Startable` interface). Keeping both risks double-initialisation.
7. Remove existing tests that break due to signature changes: `TestRegister` (old `Register(c)` signature), `TestServiceName` (method removed), `TestActionOpenWindow_Good` (uses deleted `ActionOpenWindow` type from `actions.go`)
8. Call `buildMenu()` and `setupTray()` from `HandleIPCEvents` on `ActionServiceStartup` — NOT in `OnStartup`. Display's OnStartup runs first (registration order), before sub-services register their task handlers. `ActionServiceStartup` fires after ALL services complete `OnStartup`, so IPC PERFORM calls to menu/systray will succeed.

- [ ] **Step 1: Write failing test for new Register closure pattern**

Add to `pkg/display/display_test.go`:

```go
func TestRegisterClosure_Good(t *testing.T) {
	factory := Register(nil) // nil wailsApp for testing
	assert.NotNil(t, factory)

	c, err := core.New(
		core.WithService(factory),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	svc := core.MustServiceFor[*Service](c, "display")
	assert.NotNil(t, svc)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/display/ -v -run TestRegisterClosure_Good`
Expected: FAIL — `Register` signature doesn't match

- [ ] **Step 3: Refactor Register to closure pattern**

In `pkg/display/display.go`, replace the current `Register` function:

```go
// Register creates a factory closure that captures the Wails app.
// Pass nil for testing without a Wails runtime.
func Register(wailsApp *application.App) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		s, err := New()
		if err != nil {
			return nil, err
		}
		s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
		s.wailsApp = wailsApp
		return s, nil
	}
}
```

Add `wailsApp` field and `configData` to the Service struct:

```go
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	config     Options
	configData map[string]map[string]any
	windows    *window.Manager
	tray       *systray.Manager
	menus      *menu.Manager
	notifier   *notifications.NotificationService
	events     *WSEventManager
}
```

Update `New()` to initialise configData:

```go
func New() (*Service, error) {
	return &Service{
		configData: map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
		},
	}, nil
}
```

Remove `newDisplayService()` (no longer needed — inlined into `New()`).

Remove `ServiceName()` method (auto-derived as `"display"` by WithService).

- [ ] **Step 4: Add OnStartup with config handlers registered synchronously**

Add to `pkg/display/display.go`:

```go
// OnStartup loads config and registers IPC handlers synchronously.
// CRITICAL: config handlers MUST be registered before returning —
// sub-services depend on them during their own OnStartup.
func (s *Service) OnStartup(ctx context.Context) error {
	s.loadConfig()

	// Register config query/task handlers — available NOW for sub-services
	s.Core().RegisterQuery(s.handleConfigQuery)
	s.Core().RegisterTask(s.handleConfigTask)

	// Initialise Wails wrappers if app is available (nil in tests)
	if s.wailsApp != nil {
		s.app = newWailsApp(s.wailsApp)
		s.events = NewWSEventManager(newWailsEventSource(s.wailsApp))
		s.events.SetupWindowEventListeners()
	}

	return nil
}

// Remove the old ServiceStartup(ctx, options) and Startup(ctx) methods.
// OnStartup replaces both — the Startable interface is called by Core.ServiceStartup.

func (s *Service) loadConfig() {
	// In-memory defaults. go-config integration is deferred work.
	if s.configData == nil {
		s.configData = map[string]map[string]any{
			"window":  {},
			"systray": {},
			"menu":    {},
		}
	}
}

func (s *Service) handleConfigQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case window.QueryConfig:
		return s.configData["window"], true, nil
	case systray.QueryConfig:
		return s.configData["systray"], true, nil
	case menu.QueryConfig:
		return s.configData["menu"], true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleConfigTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case window.TaskSaveConfig:
		s.configData["window"] = t.Value
		return nil, true, nil
	case systray.TaskSaveConfig:
		s.configData["systray"] = t.Value
		return nil, true, nil
	case menu.TaskSaveConfig:
		s.configData["menu"] = t.Value
		return nil, true, nil
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 5: Run test to verify Register works**

Run: `go test ./pkg/display/ -v -run TestRegisterClosure_Good`
Expected: PASS

- [ ] **Step 6: Write tests for config query/task handlers**

Add to `pkg/display/display_test.go`:

```go
func newTestDisplayService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "display")
	return svc, c
}

func TestConfigQuery_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)

	// Set window config
	svc.configData["window"] = map[string]any{
		"default_width": 1024,
	}

	result, handled, err := c.QUERY(window.QueryConfig{})
	require.NoError(t, err)
	assert.True(t, handled)
	cfg := result.(map[string]any)
	assert.Equal(t, 1024, cfg["default_width"])
}

func TestConfigQuery_Bad(t *testing.T) {
	// No display service — window config query returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(window.QueryConfig{})
	assert.False(t, handled)
}

func TestConfigTask_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	newCfg := map[string]any{"default_width": 800}
	_, handled, err := c.PERFORM(window.TaskSaveConfig{Value: newCfg})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify config was saved
	result, _, _ := c.QUERY(window.QueryConfig{})
	cfg := result.(map[string]any)
	assert.Equal(t, 800, cfg["default_width"])
}
```

- [ ] **Step 7: Run tests to verify config handlers work**

Run: `go test ./pkg/display/ -v -run "TestConfig"`
Expected: PASS

- [ ] **Step 8: Add EventTrayClick to events.go**

In `pkg/display/events.go`, add to the const block:

```go
EventTrayClick         EventType = "tray.click"
EventTrayMenuItemClick EventType = "tray.menuitem.click"
```

- [ ] **Step 9: Remove actions.go**

Delete `pkg/display/actions.go` — the `ActionOpenWindow` type is replaced by `window.TaskOpenWindow`.

- [ ] **Step 10: Run all tests to verify nothing is broken**

Run: `go test ./pkg/... -v`
Expected: PASS (all packages)

- [ ] **Step 11: Commit**

```bash
git rm pkg/display/actions.go
git add pkg/display/display.go pkg/display/events.go pkg/display/display_test.go
git commit -m "feat(display): refactor to closure Register pattern, add config IPC handlers

Register(wailsApp) returns factory closure for WithService.
OnStartup registers config query/task handlers synchronously.
Handles window.QueryConfig, systray.QueryConfig, menu.QueryConfig.
Remove ServiceName() — auto-derived as 'display' by WithService.
Remove actions.go — replaced by sub-package message types.

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

## Chunk 3: IPC Bridge + Integration Tests

### Task 5: Display IPC→WS Bridge

**Files:**
- Modify: `pkg/display/display.go`
- Modify: `pkg/display/display_test.go`

This task adds `HandleIPCEvents` to the display service, converting sub-service IPC actions to WebSocket events for TS apps.

- [ ] **Step 1: Write failing test for HandleIPCEvents bridge**

Add to `pkg/display/display_test.go`:

```go
func TestHandleIPCEvents_WindowOpened_Good(t *testing.T) {
	c, err := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	// Open a window — this should trigger ActionWindowOpened
	// which HandleIPCEvents should convert to a WS event
	result, handled, err := c.PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{window.WithName("test")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(window.WindowInfo)
	assert.Equal(t, "test", info.Name)
}
```

Note: This test verifies the full IPC flow (PERFORM → window.Service → ACTION → display.HandleIPCEvents). A full WS integration test is deferred — this validates the IPC wiring.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/display/ -v -run TestHandleIPCEvents_WindowOpened_Good`
Expected: FAIL — `HandleIPCEvents` not defined, or `window.NewMockPlatform` not exported

- [ ] **Step 3: Export mock platform from window package for integration tests**

In `pkg/window/`, create or update a test-helper file that exports mock constructors. Since mocks are in `mock_test.go` (test-only), we need an exported mock for cross-package tests.

Create `pkg/window/mock_platform.go`:

```go
package window

// MockPlatform is an exported mock for cross-package integration tests.
// For internal tests, use the unexported mockPlatform in mock_test.go.
type MockPlatform struct {
	Windows []*MockWindow
}

func NewMockPlatform() *MockPlatform {
	return &MockPlatform{}
}

func (m *MockPlatform) CreateWindow(opts PlatformWindowOptions) PlatformWindow {
	w := &MockWindow{
		name: opts.Name, title: opts.Title, url: opts.URL,
		width: opts.Width, height: opts.Height,
		x: opts.X, y: opts.Y,
	}
	m.Windows = append(m.Windows, w)
	return w
}

func (m *MockPlatform) GetWindows() []PlatformWindow {
	out := make([]PlatformWindow, len(m.Windows))
	for i, w := range m.Windows {
		out[i] = w
	}
	return out
}

type MockWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, focused   bool
	visible, alwaysOnTop bool
	closed               bool
	eventHandlers        []func(WindowEvent)
}

func (w *MockWindow) Name() string                            { return w.name }
func (w *MockWindow) Position() (int, int)                    { return w.x, w.y }
func (w *MockWindow) Size() (int, int)                        { return w.width, w.height }
func (w *MockWindow) IsMaximised() bool                       { return w.maximised }
func (w *MockWindow) IsFocused() bool                         { return w.focused }
func (w *MockWindow) SetTitle(title string)                   { w.title = title }
func (w *MockWindow) SetPosition(x, y int)                    { w.x = x; w.y = y }
func (w *MockWindow) SetSize(width, height int)               { w.width = width; w.height = height }
func (w *MockWindow) SetBackgroundColour(r, g, b, a uint8)    {}
func (w *MockWindow) SetVisibility(visible bool)              { w.visible = visible }
func (w *MockWindow) SetAlwaysOnTop(alwaysOnTop bool)         { w.alwaysOnTop = alwaysOnTop }
func (w *MockWindow) Maximise()                               { w.maximised = true }
func (w *MockWindow) Restore()                                { w.maximised = false }
func (w *MockWindow) Minimise()                               {}
func (w *MockWindow) Focus()                                  { w.focused = true }
func (w *MockWindow) Close()                                  { w.closed = true }
func (w *MockWindow) Show()                                   { w.visible = true }
func (w *MockWindow) Hide()                                   { w.visible = false }
func (w *MockWindow) Fullscreen()                             {}
func (w *MockWindow) UnFullscreen()                           {}
func (w *MockWindow) OnWindowEvent(handler func(WindowEvent)) { w.eventHandlers = append(w.eventHandlers, handler) }
```

Do the same for systray — create `pkg/systray/mock_platform.go`:

```go
package systray

// MockPlatform is an exported mock for cross-package integration tests.
type MockPlatform struct{}

func NewMockPlatform() *MockPlatform { return &MockPlatform{} }

func (m *MockPlatform) NewTray() PlatformTray { return &mockTray{} }
func (m *MockPlatform) NewMenu() PlatformMenu { return &mockMenu{} }

type mockTray struct {
	icon, templateIcon []byte
	tooltip, label     string
}

func (t *mockTray) SetIcon(data []byte)           { t.icon = data }
func (t *mockTray) SetTemplateIcon(data []byte)    { t.templateIcon = data }
func (t *mockTray) SetTooltip(text string)         { t.tooltip = text }
func (t *mockTray) SetLabel(text string)           { t.label = text }
func (t *mockTray) SetMenu(menu PlatformMenu)      {}
func (t *mockTray) AttachWindow(w WindowHandle)     {}

type mockMenu struct{ items []mockMenuItem }

func (m *mockMenu) Add(label string) PlatformMenuItem {
	mi := &mockMenuItem{label: label}
	m.items = append(m.items, *mi)
	return mi
}
func (m *mockMenu) AddSeparator() {}

type mockMenuItem struct {
	label, tooltip string
	checked, enabled bool
	onClick          func()
}

func (mi *mockMenuItem) SetTooltip(tip string)         { mi.tooltip = tip }
func (mi *mockMenuItem) SetChecked(checked bool)        { mi.checked = checked }
func (mi *mockMenuItem) SetEnabled(enabled bool)        { mi.enabled = enabled }
func (mi *mockMenuItem) OnClick(fn func())              { mi.onClick = fn }
func (mi *mockMenuItem) AddSubmenu() PlatformMenu       { return &mockMenu{} }
```

And for menu — create `pkg/menu/mock_platform.go`:

```go
package menu

// MockPlatform is an exported mock for cross-package integration tests.
type MockPlatform struct{}

func NewMockPlatform() *MockPlatform { return &MockPlatform{} }

func (m *MockPlatform) NewMenu() PlatformMenu             { return &mockPlatformMenu{} }
func (m *MockPlatform) SetApplicationMenu(menu PlatformMenu) {}

type mockPlatformMenu struct{}

func (m *mockPlatformMenu) Add(label string) PlatformMenuItem { return &mockPlatformMenuItem{} }
func (m *mockPlatformMenu) AddSeparator()                     {}
func (m *mockPlatformMenu) AddSubmenu(label string) PlatformMenu { return &mockPlatformMenu{} }
func (m *mockPlatformMenu) AddRole(role MenuRole)             {}

type mockPlatformMenuItem struct{}

func (mi *mockPlatformMenuItem) SetAccelerator(acc string) PlatformMenuItem { return mi }
func (mi *mockPlatformMenuItem) SetTooltip(tip string) PlatformMenuItem     { return mi }
func (mi *mockPlatformMenuItem) SetChecked(checked bool) PlatformMenuItem   { return mi }
func (mi *mockPlatformMenuItem) SetEnabled(enabled bool) PlatformMenuItem   { return mi }
func (mi *mockPlatformMenuItem) OnClick(fn func()) PlatformMenuItem         { return mi }
```

Note: Mock interfaces verified against `pkg/menu/platform.go` (PlatformMenu has `AddSubmenu`+`AddRole`, PlatformMenuItem has fluent returns) and `pkg/systray/platform.go` (PlatformMenu has `Add`+`AddSeparator` only, PlatformMenuItem has void returns).

- [ ] **Step 4: Add HandleIPCEvents to display.Service**

In `pkg/display/display.go`, add:

```go
// HandleIPCEvents is auto-discovered and registered by core.WithService.
// It bridges sub-service IPC actions to WebSocket events for TS apps.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	if s.events == nil {
		return nil // No WS event manager (testing without Wails)
	}

	switch m := msg.(type) {
	case core.ActionServiceStartup:
		// All services have completed OnStartup — safe to PERFORM on sub-services
		s.buildMenu()
		s.setupTray()
	case window.ActionWindowOpened:
		s.events.Emit(Event{Type: EventWindowCreate, Window: m.Name,
			Data: map[string]any{"name": m.Name}})
	case window.ActionWindowClosed:
		s.events.Emit(Event{Type: EventWindowClose, Window: m.Name,
			Data: map[string]any{"name": m.Name}})
	case window.ActionWindowMoved:
		s.events.Emit(Event{Type: EventWindowMove, Window: m.Name,
			Data: map[string]any{"x": m.X, "y": m.Y}})
	case window.ActionWindowResized:
		s.events.Emit(Event{Type: EventWindowResize, Window: m.Name,
			Data: map[string]any{"w": m.W, "h": m.H}})
	case window.ActionWindowFocused:
		s.events.Emit(Event{Type: EventWindowFocus, Window: m.Name})
	case window.ActionWindowBlurred:
		s.events.Emit(Event{Type: EventWindowBlur, Window: m.Name})
	case systray.ActionTrayClicked:
		s.events.Emit(Event{Type: EventTrayClick})
	case systray.ActionTrayMenuItemClicked:
		s.events.Emit(Event{Type: EventTrayMenuItemClick,
			Data: map[string]any{"actionId": m.ActionID}})
		s.handleTrayAction(m.ActionID)
	}
	return nil
}
```

- [ ] **Step 5: Run test to verify IPC bridge works**

Run: `go test ./pkg/display/ -v -run TestHandleIPCEvents_WindowOpened_Good`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add pkg/window/mock_platform.go pkg/systray/mock_platform.go pkg/menu/mock_platform.go
git add pkg/display/display.go pkg/display/display_test.go
git commit -m "feat(display): add HandleIPCEvents IPC→WS bridge

Display HandleIPCEvents converts sub-service actions to WS events.
Export mock platforms from each sub-package for integration tests.

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

### Task 6: Integration Tests, Delegation Conversion & Cleanup

**Files:**
- Modify: `pkg/display/display.go`
- Modify: `pkg/display/display_test.go`

This task:
1. Adds full conclave integration tests (all 4 services)
2. Converts display delegation methods to use IPC instead of direct Manager calls
3. Updates menu/tray setup to use IPC

- [ ] **Step 1: Write full conclave integration test**

Add to `pkg/display/display_test.go`:

```go
func newTestConclave(t *testing.T) *core.Core {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithService(systray.Register(systray.NewMockPlatform())),
		core.WithService(menu.Register(menu.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return c
}

func TestServiceConclave_Good(t *testing.T) {
	c := newTestConclave(t)

	// Open a window via IPC
	result, handled, err := c.PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{window.WithName("main")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(window.WindowInfo)
	assert.Equal(t, "main", info.Name)

	// Query window config from display
	val, handled, err := c.QUERY(window.QueryConfig{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.NotNil(t, val)

	// Set app menu via IPC
	_, handled, err = c.PERFORM(menu.TaskSetAppMenu{Items: []menu.MenuItem{
		{Label: "File"},
	}})
	require.NoError(t, err)
	assert.True(t, handled)

	// Query app menu via IPC
	menuResult, handled, _ := c.QUERY(menu.QueryGetAppMenu{})
	assert.True(t, handled)
	items := menuResult.([]menu.MenuItem)
	assert.Len(t, items, 1)
}

func TestServiceConclave_Bad(t *testing.T) {
	// Sub-service starts without display — config QUERY returns handled=false
	c, err := core.New(
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, _ := c.QUERY(window.QueryConfig{})
	assert.False(t, handled, "no display service means no config handler")
}
```

- [ ] **Step 2: Run integration tests**

Run: `go test ./pkg/display/ -v -run "TestServiceConclave"`
Expected: PASS

- [ ] **Step 3: Convert display delegation methods to use IPC**

In `pkg/display/display.go`, update the window management methods to route through IPC. This decouples display from direct Manager access:

```go
// OpenWindow creates a new window via IPC.
func (s *Service) OpenWindow(opts ...window.WindowOption) error {
	_, _, err := s.Core().PERFORM(window.TaskOpenWindow{Opts: opts})
	return err
}

// GetWindowInfo returns information about a window via IPC.
func (s *Service) GetWindowInfo(name string) (*window.WindowInfo, error) {
	result, handled, err := s.Core().QUERY(window.QueryWindowByName{Name: name})
	if err != nil {
		return nil, err
	}
	if !handled {
		return nil, fmt.Errorf("window service not available")
	}
	info, _ := result.(*window.WindowInfo)
	return info, nil
}

// ListWindowInfos returns information about all tracked windows via IPC.
func (s *Service) ListWindowInfos() []window.WindowInfo {
	result, handled, _ := s.Core().QUERY(window.QueryWindowList{})
	if !handled {
		return nil
	}
	list, _ := result.([]window.WindowInfo)
	return list
}

// SetWindowPosition moves a window via IPC.
func (s *Service) SetWindowPosition(name string, x, y int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetPosition{Name: name, X: x, Y: y})
	return err
}

// SetWindowSize resizes a window via IPC.
func (s *Service) SetWindowSize(name string, width, height int) error {
	_, _, err := s.Core().PERFORM(window.TaskSetSize{Name: name, W: width, H: height})
	return err
}

// MaximizeWindow maximizes a window via IPC.
func (s *Service) MaximizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMaximise{Name: name})
	return err
}

// MinimizeWindow minimizes a window via IPC.
func (s *Service) MinimizeWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskMinimise{Name: name})
	return err
}

// FocusWindow brings a window to the front via IPC.
func (s *Service) FocusWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskFocus{Name: name})
	return err
}

// CloseWindow closes a window via IPC.
func (s *Service) CloseWindow(name string) error {
	_, _, err := s.Core().PERFORM(window.TaskCloseWindow{Name: name})
	return err
}
```

Update menu click handlers to use IPC:

```go
func (s *Service) handleNewWorkspace() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("workspace-new"),
			window.WithTitle("New Workspace"),
			window.WithURL("/workspace/new"),
			window.WithSize(500, 400),
		},
	})
}

func (s *Service) handleNewFile() {
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle("New File - Editor"),
			window.WithURL("/#/developer/editor?new=true"),
			window.WithSize(1200, 800),
		},
	})
}
```

Update `buildMenu` to use IPC:

```go
func (s *Service) buildMenu() {
	items := []menu.MenuItem{
		// ... same item definitions as before ...
	}
	if runtime.GOOS != "darwin" {
		items = items[1:]
	}
	_, _, _ = s.Core().PERFORM(menu.TaskSetAppMenu{Items: items})
}
```

Update `setupTray` to use IPC:

```go
func (s *Service) setupTray() {
	_, _, _ = s.Core().PERFORM(systray.TaskSetTrayMenu{Items: []systray.TrayMenuItem{
		{Label: "Open Desktop", ActionID: "open-desktop"},
		{Label: "Close Desktop", ActionID: "close-desktop"},
		{Type: "separator"},
		{Label: "Environment Info", ActionID: "env-info"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}})
}
```

The `handleTrayAction` call is already wired in Task 5's `HandleIPCEvents` (combined with WS emission). Add the handler method:

```go
func (s *Service) handleTrayAction(actionID string) {
	switch actionID {
	case "open-desktop":
		// Show all windows
		infos := s.ListWindowInfos()
		for _, info := range infos {
			_, _, _ = s.Core().PERFORM(window.TaskFocus{Name: info.Name})
		}
	case "close-desktop":
		// Hide all windows — future: add TaskHideWindow
	case "env-info":
		s.ShowEnvironmentDialog()
	case "quit":
		if s.app != nil {
			s.app.Quit()
		}
	}
}
```

- [ ] **Step 4: Remove direct Manager fields from Service struct**

Remove `windows`, `tray`, `menus` fields from the Service struct. Remove `trackWindow` method (now handled by window.Service). Update remaining methods that need direct access (screen queries still use `application.Get()` — deferred work).

The Service struct becomes:

```go
type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	config     Options
	configData map[string]map[string]any
	notifier   *notifications.NotificationService
	events     *WSEventManager
}
```

Methods that referenced `s.windows`, `s.tray`, `s.menus` now go through IPC. Methods that need `s.windows.State()` or `s.windows.Layout()` for state persistence and layout management need to either:
- Query the window service directly (via `core.ServiceFor`)
- Or expose state/layout queries via new IPC messages

For now, remove the layout/state delegation methods from display (they can be accessed directly via the window service). Keep `SaveLayout`, `RestoreLayout`, `ListLayouts`, `DeleteLayout`, `GetLayout` using `core.ServiceFor`:

```go
func (s *Service) windowService() *window.Service {
	svc, err := core.ServiceFor[*window.Service](s.Core(), "window")
	if err != nil {
		return nil
	}
	return svc
}
```

Callers of `windowService()` must nil-check the result before use. If `nil`, the window service is not registered — layout/state operations should no-op gracefully. Example pattern for all layout/state delegation methods (`SaveLayout`, `RestoreLayout`, `ListLayouts`, `DeleteLayout`, `GetLayout`):

```go
func (s *Service) SaveLayout(name string) error {
	ws := s.windowService()
	if ws == nil {
		return fmt.Errorf("window service not available")
	}
	return ws.Manager().Layout().Save(name, ws.Manager())
}
```

Apply the same nil-guard pattern to all methods that call `windowService()`.

- [ ] **Step 5: Remove WindowInfo from display (use window.WindowInfo)**

In `pkg/display/display.go`:
- Remove `WindowInfo` struct definition
- Add type alias: `type WindowInfo = window.WindowInfo` (backward compat)
- Update `CreateWindowOptions.CreateWindow` return type to `*window.WindowInfo`

- [ ] **Step 6: Run all tests**

Run: `go test ./pkg/... -v`
Expected: PASS

- [ ] **Step 7: Verify build**

Run: `go build ./...`
Expected: SUCCESS

- [ ] **Step 8: Commit**

```bash
git add pkg/display/display.go pkg/display/display_test.go
git commit -m "feat(display): convert delegation to IPC, full conclave integration

Display methods now route through IPC bus instead of direct Manager calls.
Menu/tray setup uses PERFORM. Tray click actions handled via HandleIPCEvents.
WindowInfo aliased from window package. Direct Manager refs removed.
Integration tests verify full 4-service conclave startup and communication.

Co-Authored-By: Virgil <virgil@lethean.io>"
```

---

## Key References

| File | Role |
|------|------|
| `docs/superpowers/specs/2026-03-13-service-conclave-ipc-design.md` | Approved spec |
| `docs/superpowers/plans/2026-03-13-display-package-split.md` | Prerequisite plan (completed) |
| `/Users/snider/Code/host-uk/core/pkg/core/core.go` | Core DI — WithService, ACTION, QUERY, PERFORM |
| `/Users/snider/Code/host-uk/core/pkg/core/interfaces.go` | Message, Query, Task, QueryHandler, TaskHandler, Startable |
| `/Users/snider/Code/host-uk/core/pkg/core/message_bus.go` | Bus internals — action, query, perform dispatch |
| `/Users/snider/Code/host-uk/core/pkg/core/runtime_pkg.go` | ServiceRuntime[T] generic helper |
| `docs/ref/wails-v3/` | Wails v3 alpha-74 reference docs |
