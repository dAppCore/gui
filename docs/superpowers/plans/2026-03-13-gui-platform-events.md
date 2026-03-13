# GUI Platform & Events Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add dock/badge and application lifecycle event packages as independent core.Services with Platform interface insulation, IPC messages, and display orchestrator integration.

**Architecture:** Each package follows the three-layer pattern (IPC Bus -> Service -> Platform Interface) established by window/systray/menu. The dock service handles taskbar icon visibility and badge labels. The lifecycle service registers platform callbacks during OnStartup and broadcasts Actions for application and system events. The display orchestrator gains 9 new EventType constants, HandleIPCEvents cases for both packages, and WS->IPC cases for dock commands.

**Tech Stack:** Go, core/go DI framework, Wails v3 (behind Platform interfaces), gorilla/websocket (WSEventManager)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-platform-events-design.md`

---

## File Structure

### New files (2 packages x 5 files = 10 files)

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/dock/` | `platform.go` | Platform interface (5 methods: ShowIcon, HideIcon, SetBadge, RemoveBadge, IsVisible) |
| | `messages.go` | IPC message types (1 Query, 4 Tasks, 1 Action) |
| | `register.go` | Register() factory closure |
| | `service.go` | Service struct, OnStartup(), handlers, visibility broadcast |
| | `service_test.go` | Tests with mock platform |
| `pkg/lifecycle/` | `platform.go` | Platform interface (2 methods) + EventType enum (7 values) |
| | `messages.go` | IPC message types (8 Actions) |
| | `register.go` | Register() factory closure |
| | `service.go` | Service struct, OnStartup() callback registration, OnShutdown() cancel |
| | `service_test.go` | Tests with mock platform + event simulation |

### Modified files

| File | Change |
|------|--------|
| `pkg/display/events.go` | Add 9 new EventType constants |
| `pkg/display/display.go` | Add imports for dock + lifecycle, HandleIPCEvents cases (10 new), WS->IPC cases (5 new dock commands) |

---

## Task 1: Create pkg/dock

**Files:**
- Create: `pkg/dock/platform.go`
- Create: `pkg/dock/messages.go`
- Create: `pkg/dock/register.go`
- Create: `pkg/dock/service.go`
- Create: `pkg/dock/service_test.go`

- [ ] **Step 1: Create platform.go**

```go
// pkg/dock/platform.go
package dock

// Platform abstracts the dock/taskbar backend (Wails v3).
// macOS: dock icon show/hide + badge.
// Windows: taskbar badge only (show/hide not supported).
// Linux: not supported — adapter returns nil for all operations.
type Platform interface {
	ShowIcon() error
	HideIcon() error
	SetBadge(label string) error
	RemoveBadge() error
	IsVisible() bool
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/dock/messages.go
package dock

// --- Queries (read-only) ---

// QueryVisible returns whether the dock icon is visible. Result: bool
type QueryVisible struct{}

// --- Tasks (side-effects) ---

// TaskShowIcon shows the dock/taskbar icon. Result: nil
type TaskShowIcon struct{}

// TaskHideIcon hides the dock/taskbar icon. Result: nil
type TaskHideIcon struct{}

// TaskSetBadge sets the dock/taskbar badge label.
// Empty string "" shows the default system badge indicator.
// Numeric "3", "99" shows unread count. Text "New", "Paused" shows brief status.
// Result: nil
type TaskSetBadge struct{ Label string }

// TaskRemoveBadge removes the dock/taskbar badge. Result: nil
type TaskRemoveBadge struct{}

// --- Actions (broadcasts) ---

// ActionVisibilityChanged is broadcast after a successful TaskShowIcon or TaskHideIcon.
type ActionVisibilityChanged struct{ Visible bool }
```

- [ ] **Step 3: Create register.go**

```go
// pkg/dock/register.go
package dock

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}
```

- [ ] **Step 4: Write failing test**

```go
// pkg/dock/service_test.go
package dock

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	visible   bool
	badge     string
	hasBadge  bool
	showErr   error
	hideErr   error
	badgeErr  error
	removeErr error
}

func (m *mockPlatform) ShowIcon() error {
	if m.showErr != nil {
		return m.showErr
	}
	m.visible = true
	return nil
}

func (m *mockPlatform) HideIcon() error {
	if m.hideErr != nil {
		return m.hideErr
	}
	m.visible = false
	return nil
}

func (m *mockPlatform) SetBadge(label string) error {
	if m.badgeErr != nil {
		return m.badgeErr
	}
	m.badge = label
	m.hasBadge = true
	return nil
}

func (m *mockPlatform) RemoveBadge() error {
	if m.removeErr != nil {
		return m.removeErr
	}
	m.badge = ""
	m.hasBadge = false
	return nil
}

func (m *mockPlatform) IsVisible() bool { return m.visible }

// --- Test helpers ---

func newTestDockService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := &mockPlatform{visible: true}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "dock")
	return svc, c, mock
}

// --- Tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestDockService(t)
	assert.NotNil(t, svc)
}

func TestQueryVisible_Good(t *testing.T) {
	_, c, _ := newTestDockService(t)
	result, handled, err := c.QUERY(QueryVisible{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)
}

func TestQueryVisible_Bad(t *testing.T) {
	// No dock service registered — QUERY returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(QueryVisible{})
	assert.False(t, handled)
}

func TestTaskShowIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = false // Start hidden

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskShowIcon{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.visible)
	require.NotNil(t, received)
	assert.True(t, received.Visible)
}

func TestTaskHideIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = true // Start visible

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskHideIcon{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, mock.visible)
	require.NotNil(t, received)
	assert.False(t, received.Visible)
}

func TestTaskSetBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	_, handled, err := c.PERFORM(TaskSetBadge{Label: "3"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "3", mock.badge)
	assert.True(t, mock.hasBadge)
}

func TestTaskSetBadge_EmptyLabel_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	_, handled, err := c.PERFORM(TaskSetBadge{Label: ""})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", mock.badge)
	assert.True(t, mock.hasBadge) // Empty string = default system badge indicator
}

func TestTaskRemoveBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	// Set a badge first
	_, _, _ = c.PERFORM(TaskSetBadge{Label: "5"})

	_, handled, err := c.PERFORM(TaskRemoveBadge{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", mock.badge)
	assert.False(t, mock.hasBadge)
}

func TestTaskShowIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.showErr = assert.AnError

	_, handled, err := c.PERFORM(TaskShowIcon{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskHideIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.hideErr = assert.AnError

	_, handled, err := c.PERFORM(TaskHideIcon{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetBadge_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.badgeErr = assert.AnError

	_, handled, err := c.PERFORM(TaskSetBadge{Label: "3"})
	assert.True(t, handled)
	assert.Error(t, err)
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/dock/ -v`
Expected: FAIL — `Service` type not defined, `Options` type not defined

- [ ] **Step 6: Create service.go**

```go
// pkg/dock/service.go
package dock

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the dock service.
type Options struct{}

// Service is a core.Service managing dock/taskbar operations via IPC.
// It embeds ServiceRuntime for Core access and delegates to Platform.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Query Handlers ---

func (s *Service) handleQuery(c *core.Core, q core.Query) (any, bool, error) {
	switch q.(type) {
	case QueryVisible:
		return s.platform.IsVisible(), true, nil
	default:
		return nil, false, nil
	}
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskShowIcon:
		if err := s.platform.ShowIcon(); err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionVisibilityChanged{Visible: true})
		return nil, true, nil
	case TaskHideIcon:
		if err := s.platform.HideIcon(); err != nil {
			return nil, true, err
		}
		_ = s.Core().ACTION(ActionVisibilityChanged{Visible: false})
		return nil, true, nil
	case TaskSetBadge:
		if err := s.platform.SetBadge(t.Label); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	case TaskRemoveBadge:
		if err := s.platform.RemoveBadge(); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/dock/ -v`
Expected: PASS (12 tests)

- [ ] **Step 8: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/dock/
git commit -m "feat(dock): add dock/badge core.Service with Platform interface and IPC"
```

---

## Task 2: Create pkg/lifecycle

**Files:**
- Create: `pkg/lifecycle/platform.go`
- Create: `pkg/lifecycle/messages.go`
- Create: `pkg/lifecycle/register.go`
- Create: `pkg/lifecycle/service.go`
- Create: `pkg/lifecycle/service_test.go`

- [ ] **Step 1: Create platform.go with EventType enum**

```go
// pkg/lifecycle/platform.go
package lifecycle

// EventType identifies application and system lifecycle events.
type EventType int

const (
	EventApplicationStarted EventType = iota
	EventWillTerminate                        // macOS only
	EventDidBecomeActive                      // macOS only
	EventDidResignActive                      // macOS only
	EventPowerStatusChanged                   // Windows only (APMPowerStatusChange)
	EventSystemSuspend                        // Windows only (APMSuspend)
	EventSystemResume                         // Windows only (APMResume)
)

// Platform abstracts the application lifecycle backend (Wails v3).
// OnApplicationEvent registers a handler for a fire-and-forget event type.
// OnOpenedWithFile registers a handler for file-open events (carries path data).
// Both return a cancel function that deregisters the handler.
// Platform-specific events no-op silently on unsupported OS (adapter registers nothing).
type Platform interface {
	OnApplicationEvent(eventType EventType, handler func()) func()
	OnOpenedWithFile(handler func(path string)) func()
}
```

- [ ] **Step 2: Create messages.go with 8 Actions**

```go
// pkg/lifecycle/messages.go
package lifecycle

// All lifecycle events are broadcasts (Actions). There are no Queries or Tasks.

// ActionApplicationStarted fires when the platform application starts.
// Distinct from core.ActionServiceStartup — this is platform-level readiness.
type ActionApplicationStarted struct{}

// ActionOpenedWithFile fires when the application is opened with a file argument.
type ActionOpenedWithFile struct{ Path string }

// ActionWillTerminate fires when the application is about to terminate (macOS only).
type ActionWillTerminate struct{}

// ActionDidBecomeActive fires when the application becomes the active app (macOS only).
type ActionDidBecomeActive struct{}

// ActionDidResignActive fires when the application resigns active status (macOS only).
type ActionDidResignActive struct{}

// ActionPowerStatusChanged fires on power status changes (Windows only: APMPowerStatusChange).
type ActionPowerStatusChanged struct{}

// ActionSystemSuspend fires when the system is about to suspend (Windows only: APMSuspend).
type ActionSystemSuspend struct{}

// ActionSystemResume fires when the system resumes from suspend (Windows only: APMResume).
type ActionSystemResume struct{}
```

- [ ] **Step 3: Create register.go**

```go
// pkg/lifecycle/register.go
package lifecycle

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
		}, nil
	}
}
```

- [ ] **Step 4: Write failing test**

```go
// pkg/lifecycle/service_test.go
package lifecycle

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	mu       sync.Mutex
	handlers map[EventType][]func()
	fileHandlers []func(string)
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		handlers: make(map[EventType][]func()),
	}
}

func (m *mockPlatform) OnApplicationEvent(eventType EventType, handler func()) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[eventType] = append(m.handlers[eventType], handler)
	idx := len(m.handlers[eventType]) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.handlers[eventType]) {
			m.handlers[eventType] = append(m.handlers[eventType][:idx], m.handlers[eventType][idx+1:]...)
		}
	}
}

func (m *mockPlatform) OnOpenedWithFile(handler func(string)) func() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fileHandlers = append(m.fileHandlers, handler)
	idx := len(m.fileHandlers) - 1
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if idx < len(m.fileHandlers) {
			m.fileHandlers = append(m.fileHandlers[:idx], m.fileHandlers[idx+1:]...)
		}
	}
}

// simulateEvent fires all registered handlers for the given event type.
func (m *mockPlatform) simulateEvent(eventType EventType) {
	m.mu.Lock()
	handlers := make([]func(), len(m.handlers[eventType]))
	copy(handlers, m.handlers[eventType])
	m.mu.Unlock()
	for _, h := range handlers {
		h()
	}
}

// simulateFileOpen fires all registered file-open handlers.
func (m *mockPlatform) simulateFileOpen(path string) {
	m.mu.Lock()
	handlers := make([]func(string), len(m.fileHandlers))
	copy(handlers, m.fileHandlers)
	m.mu.Unlock()
	for _, h := range handlers {
		h(path)
	}
}

// handlerCount returns the number of registered handlers for event-based + file-based.
func (m *mockPlatform) handlerCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := len(m.fileHandlers)
	for _, handlers := range m.handlers {
		count += len(handlers)
	}
	return count
}

// --- Test helpers ---

func newTestLifecycleService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := newMockPlatform()
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "lifecycle")
	return svc, c, mock
}

// --- Tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestLifecycleService(t)
	assert.NotNil(t, svc)
}

func TestApplicationStarted_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventApplicationStarted)
	assert.True(t, received)
}

func TestDidBecomeActive_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionDidBecomeActive); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventDidBecomeActive)
	assert.True(t, received)
}

func TestDidResignActive_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionDidResignActive); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventDidResignActive)
	assert.True(t, received)
}

func TestWillTerminate_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionWillTerminate); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventWillTerminate)
	assert.True(t, received)
}

func TestPowerStatusChanged_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionPowerStatusChanged); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventPowerStatusChanged)
	assert.True(t, received)
}

func TestSystemSuspend_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionSystemSuspend); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventSystemSuspend)
	assert.True(t, received)
}

func TestSystemResume_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionSystemResume); ok {
			received = true
		}
		return nil
	})

	mock.simulateEvent(EventSystemResume)
	assert.True(t, received)
}

func TestOpenedWithFile_Good(t *testing.T) {
	_, c, mock := newTestLifecycleService(t)

	var receivedPath string
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionOpenedWithFile); ok {
			receivedPath = a.Path
		}
		return nil
	})

	mock.simulateFileOpen("/Users/snider/Documents/test.txt")
	assert.Equal(t, "/Users/snider/Documents/test.txt", receivedPath)
}

func TestOnShutdown_CancelsAll_Good(t *testing.T) {
	svc, _, mock := newTestLifecycleService(t)

	// Verify handlers were registered during OnStartup
	assert.Greater(t, mock.handlerCount(), 0, "handlers should be registered after OnStartup")

	// Shutdown should cancel all registrations
	err := svc.OnShutdown(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 0, mock.handlerCount(), "all handlers should be cancelled after OnShutdown")
}

func TestRegister_Bad(t *testing.T) {
	// No lifecycle service registered — actions are not received
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)

	var received bool
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if _, ok := msg.(ActionApplicationStarted); ok {
			received = true
		}
		return nil
	})

	// No way to trigger events without the service
	assert.False(t, received)
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/lifecycle/ -v`
Expected: FAIL — `Service` type not defined, `Options` type not defined

- [ ] **Step 6: Create service.go with OnStartup/OnShutdown**

```go
// pkg/lifecycle/service.go
package lifecycle

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the lifecycle service.
type Options struct{}

// Service is a core.Service that registers platform lifecycle callbacks
// and broadcasts corresponding IPC Actions. It implements both Startable
// and Stoppable: OnStartup registers all callbacks, OnShutdown cancels them.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	cancels  []func()
}

// OnStartup registers a platform callback for each EventType and for file-open.
// Each callback broadcasts the corresponding Action via s.Core().ACTION().
func (s *Service) OnStartup(ctx context.Context) error {
	// Register fire-and-forget event callbacks
	eventActions := map[EventType]func(){
		EventApplicationStarted: func() { _ = s.Core().ACTION(ActionApplicationStarted{}) },
		EventWillTerminate:      func() { _ = s.Core().ACTION(ActionWillTerminate{}) },
		EventDidBecomeActive:    func() { _ = s.Core().ACTION(ActionDidBecomeActive{}) },
		EventDidResignActive:    func() { _ = s.Core().ACTION(ActionDidResignActive{}) },
		EventPowerStatusChanged: func() { _ = s.Core().ACTION(ActionPowerStatusChanged{}) },
		EventSystemSuspend:      func() { _ = s.Core().ACTION(ActionSystemSuspend{}) },
		EventSystemResume:       func() { _ = s.Core().ACTION(ActionSystemResume{}) },
	}

	for eventType, handler := range eventActions {
		cancel := s.platform.OnApplicationEvent(eventType, handler)
		s.cancels = append(s.cancels, cancel)
	}

	// Register file-open callback (carries data)
	cancel := s.platform.OnOpenedWithFile(func(path string) {
		_ = s.Core().ACTION(ActionOpenedWithFile{Path: path})
	})
	s.cancels = append(s.cancels, cancel)

	return nil
}

// OnShutdown cancels all registered platform callbacks.
func (s *Service) OnShutdown(ctx context.Context) error {
	for _, cancel := range s.cancels {
		cancel()
	}
	s.cancels = nil
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
// Lifecycle events are all outbound (platform -> IPC) so there is nothing to handle here.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/lifecycle/ -v`
Expected: PASS (12 tests)

- [ ] **Step 8: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/lifecycle/
git commit -m "feat(lifecycle): add application lifecycle core.Service with Platform interface and IPC"
```

---

## Task 3: Display orchestrator updates

**Files:**
- Modify: `pkg/display/events.go`
- Modify: `pkg/display/display.go`

- [ ] **Step 1: Add 9 new EventType constants to events.go**

Add after the existing constants in `pkg/display/events.go`:

```go
// Dock events
EventDockVisibility    EventType = "dock.visibility-changed"

// Application lifecycle events
EventAppStarted        EventType = "app.started"
EventAppOpenedWithFile EventType = "app.opened-with-file"
EventAppWillTerminate  EventType = "app.will-terminate"
EventAppActive         EventType = "app.active"
EventAppInactive       EventType = "app.inactive"

// System events
EventSystemPowerChange EventType = "system.power-change"
EventSystemSuspend     EventType = "system.suspend"
EventSystemResume      EventType = "system.resume"
```

The full const block in `pkg/display/events.go` should become:

```go
const (
	EventWindowFocus  EventType = "window.focus"
	EventWindowBlur   EventType = "window.blur"
	EventWindowMove   EventType = "window.move"
	EventWindowResize EventType = "window.resize"
	EventWindowClose  EventType = "window.close"
	EventWindowCreate EventType = "window.create"
	EventThemeChange         EventType = "theme.change"
	EventScreenChange        EventType = "screen.change"
	EventNotificationClick   EventType = "notification.click"
	EventTrayClick           EventType = "tray.click"
	EventTrayMenuItemClick   EventType = "tray.menuitem.click"
	EventDockVisibility      EventType = "dock.visibility-changed"
	EventAppStarted          EventType = "app.started"
	EventAppOpenedWithFile   EventType = "app.opened-with-file"
	EventAppWillTerminate    EventType = "app.will-terminate"
	EventAppActive           EventType = "app.active"
	EventAppInactive         EventType = "app.inactive"
	EventSystemPowerChange   EventType = "system.power-change"
	EventSystemSuspend       EventType = "system.suspend"
	EventSystemResume        EventType = "system.resume"
)
```

- [ ] **Step 2: Add dock and lifecycle imports to display.go**

Add to the import block in `pkg/display/display.go`:

```go
"forge.lthn.ai/core/gui/pkg/dock"
"forge.lthn.ai/core/gui/pkg/lifecycle"
```

- [ ] **Step 3: Add HandleIPCEvents cases for dock actions**

Add after the existing `screen.ActionScreensChanged` case in `HandleIPCEvents`:

```go
case dock.ActionVisibilityChanged:
	if s.events != nil {
		s.events.Emit(Event{Type: EventDockVisibility,
			Data: map[string]any{"visible": m.Visible}})
	}
```

- [ ] **Step 4: Add HandleIPCEvents cases for lifecycle actions (8 cases)**

Add after the dock case:

```go
case lifecycle.ActionApplicationStarted:
	if s.events != nil {
		s.events.Emit(Event{Type: EventAppStarted})
	}
case lifecycle.ActionOpenedWithFile:
	if s.events != nil {
		s.events.Emit(Event{Type: EventAppOpenedWithFile,
			Data: map[string]any{"path": m.Path}})
	}
case lifecycle.ActionWillTerminate:
	if s.events != nil {
		s.events.Emit(Event{Type: EventAppWillTerminate})
	}
case lifecycle.ActionDidBecomeActive:
	if s.events != nil {
		s.events.Emit(Event{Type: EventAppActive})
	}
case lifecycle.ActionDidResignActive:
	if s.events != nil {
		s.events.Emit(Event{Type: EventAppInactive})
	}
case lifecycle.ActionPowerStatusChanged:
	if s.events != nil {
		s.events.Emit(Event{Type: EventSystemPowerChange})
	}
case lifecycle.ActionSystemSuspend:
	if s.events != nil {
		s.events.Emit(Event{Type: EventSystemSuspend})
	}
case lifecycle.ActionSystemResume:
	if s.events != nil {
		s.events.Emit(Event{Type: EventSystemResume})
	}
```

- [ ] **Step 5: Add WS->IPC cases for dock commands**

The display orchestrator needs a `handleWSMessage` method (or extend the existing WS message handling). Add 5 dock WS->IPC cases. These go in the WS message handler (the method that dispatches `msg.Type` strings to IPC calls). If no such method exists yet, add it as part of the `handleMessages` flow in the WSEventManager or as a new method on Service:

```go
case "dock:show":
	result, handled, err = s.Core().PERFORM(dock.TaskShowIcon{})
case "dock:hide":
	result, handled, err = s.Core().PERFORM(dock.TaskHideIcon{})
case "dock:badge":
	label, _ := msg.Data["label"].(string)
	result, handled, err = s.Core().PERFORM(dock.TaskSetBadge{Label: label})
case "dock:badge-remove":
	result, handled, err = s.Core().PERFORM(dock.TaskRemoveBadge{})
case "dock:visible":
	result, handled, err = s.Core().QUERY(dock.QueryVisible{})
```

Key detail: `dock:badge` uses the safe comma-ok type assertion `label, _ := msg.Data["label"].(string)` so a missing or non-string `label` defaults to `""` (which is a valid badge value per the spec — shows default system badge indicator).

Lifecycle events are outbound-only (Actions -> WS). No inbound WS->IPC cases needed for lifecycle.

- [ ] **Step 6: Run all tests to verify nothing is broken**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/display/ ./pkg/dock/ ./pkg/lifecycle/ -v`
Expected: PASS (all existing display tests + new dock + lifecycle tests)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/events.go pkg/display/display.go
git commit -m "feat(display): add dock + lifecycle HandleIPCEvents and WS bridge integration"
```

---

## Task 4: Final verification and commit

- [ ] **Step 1: Run the full test suite**

Run: `cd /Users/snider/Code/core/gui && go test ./... -v`
Expected: PASS — all packages compile and all tests pass.

- [ ] **Step 2: Run lint and vet**

Run: `cd /Users/snider/Code/core/gui && go vet ./pkg/dock/ ./pkg/lifecycle/ ./pkg/display/`
Expected: No issues.

- [ ] **Step 3: Verify no circular dependencies**

Confirm dependency direction:
- `pkg/dock` imports `forge.lthn.ai/core/go/pkg/core` only
- `pkg/lifecycle` imports `forge.lthn.ai/core/go/pkg/core` only
- `pkg/display` imports `pkg/dock` and `pkg/lifecycle` (message types only)
- Neither `pkg/dock` nor `pkg/lifecycle` import `pkg/display`

Run: `cd /Users/snider/Code/core/gui && go vet ./...`
Expected: No import cycle errors.

- [ ] **Step 4: Final commit (if any formatting or vet fixes were needed)**

```bash
cd /Users/snider/Code/core/gui
git add -A
git commit -m "chore(gui): platform & events final cleanup and verification"
```

---

## Summary

| Task | Files | Tests | Description |
|------|-------|-------|-------------|
| 1 | 5 created | 12 | pkg/dock: Platform interface, IPC messages, Service with visibility broadcast |
| 2 | 5 created | 12 | pkg/lifecycle: Platform interface, EventType enum, 8 Action types, OnStartup/OnShutdown |
| 3 | 2 modified | 0 new (existing pass) | Display orchestrator: 9 EventType constants, 10 HandleIPCEvents cases, 5 WS->IPC cases |
| 4 | 0 | full suite | Verification: tests, vet, lint, dependency direction |

**Total:** 10 new files, 2 modified files, ~24 tests, 4 commits.

## Platform Mapping Reference

| Our Type | Wails Event | Platforms |
|----------|-------------|-----------|
| `EventApplicationStarted` | `events.Common.ApplicationStarted` | all |
| `EventWillTerminate` | `events.Mac.ApplicationWillTerminate` | macOS only |
| `EventDidBecomeActive` | `events.Mac.ApplicationDidBecomeActive` | macOS only |
| `EventDidResignActive` | `events.Mac.ApplicationDidResignActive` | macOS only |
| `EventPowerStatusChanged` | `events.Windows.APMPowerStatusChange` | Windows only |
| `EventSystemSuspend` | `events.Windows.APMSuspend` | Windows only |
| `EventSystemResume` | `events.Windows.APMResume` | Windows only |

## Licence

EUPL-1.2
