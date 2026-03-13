# CoreGUI Spec B: New Input — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add keybinding and browser packages as core.Services, enhance the window service with file drop events, and wire all new IPC messages through the display orchestrator.

**Architecture:** keybinding and browser follow the three-layer pattern (IPC Bus -> Service -> Platform Interface) established by window/systray/menu. File drop extends the existing window package. The display orchestrator gains HandleIPCEvents cases for new Action types and WS->IPC bridge cases for new Task/Query types.

**Tech Stack:** Go, core/go DI framework, Wails v3 (behind Platform interfaces), gorilla/websocket (WSEventManager)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-new-input-design.md`

---

## File Structure

### New files (2 packages x 5 files = 10 files)

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/keybinding/` | `platform.go` | Platform interface (3 methods: Add, Remove, GetAll) |
| | `messages.go` | IPC message types (QueryList, TaskAdd, TaskRemove, ActionTriggered) + BindingInfo |
| | `register.go` | Register() factory closure |
| | `service.go` | Service struct, OnStartup(), handlers, in-memory registry |
| | `service_test.go` | Tests with mock platform |
| `pkg/browser/` | `platform.go` | Platform interface (2 methods: OpenURL, OpenFile) |
| | `messages.go` | IPC message types (TaskOpenURL, TaskOpenFile) |
| | `register.go` | Register() factory closure |
| | `service.go` | Service struct, OnStartup(), handlers |
| | `service_test.go` | Tests with mock platform |

### Modified files

| File | Change |
|------|--------|
| `pkg/window/platform.go` | Add `OnFileDrop(handler func(paths []string, targetID string))` to PlatformWindow interface |
| `pkg/window/messages.go` | Add `ActionFilesDropped` action type |
| `pkg/window/service.go` | Add `pw.OnFileDrop()` call in `trackWindow()` |
| `pkg/window/mock_platform.go` | Add `OnFileDrop` to exported MockWindow |
| `pkg/window/mock_test.go` | Add `OnFileDrop` to unexported mockWindow |
| `pkg/window/service_test.go` | Add file drop test |
| `pkg/display/events.go` | Add `EventKeybindingTriggered`, `EventWindowFileDrop` constants |
| `pkg/display/display.go` | Add HandleIPCEvents cases + WS->IPC cases in handleWSMessage + new imports |

---

## Task 1: Create pkg/keybinding

**Files:**
- Create: `pkg/keybinding/platform.go`
- Create: `pkg/keybinding/messages.go`
- Create: `pkg/keybinding/register.go`
- Create: `pkg/keybinding/service.go`
- Create: `pkg/keybinding/service_test.go`

- [ ] **Step 1: Create platform.go**

```go
// pkg/keybinding/platform.go
package keybinding

// Platform abstracts the keyboard shortcut backend (Wails v3).
type Platform interface {
	// Add registers a global keyboard shortcut with the given accelerator string.
	// The handler is called when the shortcut is triggered.
	// Accelerator syntax is platform-aware: "Cmd+S" (macOS), "Ctrl+S" (Windows/Linux).
	// Special keys: F1-F12, Escape, Enter, Space, Tab, Backspace, Delete, arrow keys.
	Add(accelerator string, handler func()) error

	// Remove unregisters a previously registered keyboard shortcut.
	Remove(accelerator string) error

	// GetAll returns all currently registered accelerator strings.
	// Used for adapter-level reconciliation only — not read by QueryList.
	GetAll() []string
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/keybinding/messages.go
package keybinding

import "errors"

// ErrAlreadyRegistered is returned when attempting to add a binding
// that already exists. Callers must TaskRemove first to rebind.
var ErrAlreadyRegistered = errors.New("keybinding: accelerator already registered")

// BindingInfo describes a registered keyboard shortcut.
type BindingInfo struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// --- Queries ---

// QueryList returns all registered bindings. Result: []BindingInfo
type QueryList struct{}

// --- Tasks ---

// TaskAdd registers a new keyboard shortcut. Result: nil
// Returns ErrAlreadyRegistered if the accelerator is already bound.
type TaskAdd struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// TaskRemove unregisters a keyboard shortcut. Result: nil
type TaskRemove struct {
	Accelerator string `json:"accelerator"`
}

// --- Actions ---

// ActionTriggered is broadcast when a registered shortcut is activated.
type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
```

- [ ] **Step 3: Create register.go**

```go
// pkg/keybinding/register.go
package keybinding

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			bindings:       make(map[string]BindingInfo),
		}, nil
	}
}
```

- [ ] **Step 4: Write failing test**

```go
// pkg/keybinding/service_test.go
package keybinding

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatform records Add/Remove calls and allows triggering shortcuts.
type mockPlatform struct {
	mu       sync.Mutex
	handlers map[string]func()
	removed  []string
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{handlers: make(map[string]func())}
}

func (m *mockPlatform) Add(accelerator string, handler func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[accelerator] = handler
	return nil
}

func (m *mockPlatform) Remove(accelerator string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, accelerator)
	m.removed = append(m.removed, accelerator)
	return nil
}

func (m *mockPlatform) GetAll() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, 0, len(m.handlers))
	for k := range m.handlers {
		out = append(out, k)
	}
	return out
}

// trigger simulates a shortcut keypress by calling the registered handler.
func (m *mockPlatform) trigger(accelerator string) {
	m.mu.Lock()
	h, ok := m.handlers[accelerator]
	m.mu.Unlock()
	if ok {
		h()
	}
}

func newTestKeybindingService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "keybinding")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	mp := newMockPlatform()
	svc, _ := newTestKeybindingService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskAdd_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, handled, err := c.PERFORM(TaskAdd{
		Accelerator: "Ctrl+S", Description: "Save",
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify binding registered on platform
	assert.Contains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskAdd_Bad_Duplicate(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Second add with same accelerator should fail
	_, handled, err := c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save Again"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, handled, err := c.PERFORM(TaskRemove{Accelerator: "Ctrl+S"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify removed from platform
	assert.NotContains(t, mp.GetAll(), "Ctrl+S")
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, handled, err := c.PERFORM(TaskRemove{Accelerator: "Ctrl+X"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+Z", Description: "Undo"})

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]BindingInfo)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]BindingInfo)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_TriggerBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	// Capture broadcast actions
	var triggered ActionTriggered
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionTriggered); ok {
			mu.Lock()
			triggered = a
			mu.Unlock()
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})

	// Simulate shortcut trigger via mock
	mp.trigger("Ctrl+S")

	mu.Lock()
	assert.Equal(t, "Ctrl+S", triggered.Accelerator)
	mu.Unlock()
}

func TestTaskAdd_Good_RebindAfterRemove(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestKeybindingService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save"})
	_, _, _ = c.PERFORM(TaskRemove{Accelerator: "Ctrl+S"})

	// Should succeed after remove
	_, handled, err := c.PERFORM(TaskAdd{Accelerator: "Ctrl+S", Description: "Save v2"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify new description
	result, _, _ := c.QUERY(QueryList{})
	list := result.([]BindingInfo)
	assert.Len(t, list, 1)
	assert.Equal(t, "Save v2", list[0].Description)
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryList{})
	assert.False(t, handled)
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/keybinding/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 6: Create service.go**

```go
// pkg/keybinding/service.go
package keybinding

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the keybinding service.
type Options struct{}

// Service is a core.Service managing keyboard shortcuts via IPC.
// It maintains an in-memory registry of bindings and delegates
// platform-level registration to the Platform interface.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	bindings map[string]BindingInfo
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
	case QueryList:
		return s.queryList(), true, nil
	default:
		return nil, false, nil
	}
}

// queryList reads from the in-memory registry (not platform.GetAll()).
func (s *Service) queryList() []BindingInfo {
	result := make([]BindingInfo, 0, len(s.bindings))
	for _, info := range s.bindings {
		result = append(result, info)
	}
	return result
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskAdd:
		return nil, true, s.taskAdd(t)
	case TaskRemove:
		return nil, true, s.taskRemove(t)
	default:
		return nil, false, nil
	}
}

func (s *Service) taskAdd(t TaskAdd) error {
	if _, exists := s.bindings[t.Accelerator]; exists {
		return ErrAlreadyRegistered
	}

	// Register on platform with a callback that broadcasts ActionTriggered
	err := s.platform.Add(t.Accelerator, func() {
		_ = s.Core().ACTION(ActionTriggered{Accelerator: t.Accelerator})
	})
	if err != nil {
		return fmt.Errorf("keybinding: platform add failed: %w", err)
	}

	s.bindings[t.Accelerator] = BindingInfo{
		Accelerator: t.Accelerator,
		Description: t.Description,
	}
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.bindings[t.Accelerator]; !exists {
		return fmt.Errorf("keybinding: not registered: %s", t.Accelerator)
	}

	err := s.platform.Remove(t.Accelerator)
	if err != nil {
		return fmt.Errorf("keybinding: platform remove failed: %w", err)
	}

	delete(s.bindings, t.Accelerator)
	return nil
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/keybinding/ -v`
Expected: PASS (10 tests)

- [ ] **Step 8: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/keybinding/
git commit -m "feat(keybinding): add keybinding core.Service with Platform interface and IPC

Implements pkg/keybinding with three-layer pattern: IPC Bus -> Service -> Platform.
Service maintains in-memory registry, ErrAlreadyRegistered on duplicates.
QueryList reads from service registry, not platform.GetAll().
ActionTriggered broadcast on shortcut trigger via platform callback.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Create pkg/browser

**Files:**
- Create: `pkg/browser/platform.go`
- Create: `pkg/browser/messages.go`
- Create: `pkg/browser/register.go`
- Create: `pkg/browser/service.go`
- Create: `pkg/browser/service_test.go`

- [ ] **Step 1: Create platform.go**

```go
// pkg/browser/platform.go
package browser

// Platform abstracts the system browser/file-opener backend.
type Platform interface {
	// OpenURL opens the given URL in the default system browser.
	OpenURL(url string) error

	// OpenFile opens the given file path with the system default application.
	OpenFile(path string) error
}
```

- [ ] **Step 2: Create messages.go**

```go
// pkg/browser/messages.go
package browser

// --- Tasks (all side-effects, no queries or actions) ---

// TaskOpenURL opens a URL in the default system browser. Result: nil
type TaskOpenURL struct {
	URL string `json:"url"`
}

// TaskOpenFile opens a file with the system default application. Result: nil
type TaskOpenFile struct {
	Path string `json:"path"`
}
```

- [ ] **Step 3: Create register.go**

```go
// pkg/browser/register.go
package browser

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
// pkg/browser/service_test.go
package browser

import (
	"context"
	"errors"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	lastURL  string
	lastPath string
	urlErr   error
	fileErr  error
}

func (m *mockPlatform) OpenURL(url string) error {
	m.lastURL = url
	return m.urlErr
}

func (m *mockPlatform) OpenFile(path string) error {
	m.lastPath = path
	return m.fileErr
}

func newTestBrowserService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "browser")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	mp := &mockPlatform{}
	svc, _ := newTestBrowserService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskOpenURL_Good(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "https://example.com", mp.lastURL)
}

func TestTaskOpenURL_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{urlErr: errors.New("browser not found")}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOpenFile_Good(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenFile{Path: "/tmp/readme.txt"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/readme.txt", mp.lastPath)
}

func TestTaskOpenFile_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{fileErr: errors.New("file not found")}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenFile{Path: "/nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOpenURL_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	assert.False(t, handled)
}
```

- [ ] **Step 5: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/browser/ -v`
Expected: FAIL — `Service` type not defined

- [ ] **Step 6: Create service.go**

```go
// pkg/browser/service.go
package browser

import (
	"context"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the browser service.
type Options struct{}

// Service is a core.Service that delegates browser/file-open operations
// to the platform. It is stateless — no queries, no actions.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(ctx context.Context) error {
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	return nil
}

// --- Task Handlers ---

func (s *Service) handleTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskOpenURL:
		return nil, true, s.platform.OpenURL(t.URL)
	case TaskOpenFile:
		return nil, true, s.platform.OpenFile(t.Path)
	default:
		return nil, false, nil
	}
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/browser/ -v`
Expected: PASS (6 tests)

- [ ] **Step 8: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/browser/
git commit -m "feat(browser): add browser core.Service with Platform interface and IPC

Implements pkg/browser with three-layer pattern: IPC Bus -> Service -> Platform.
Stateless service — delegates OpenURL and OpenFile to platform adapter.
No queries or actions, tasks only.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Window file drop enhancement

**Files:**
- Modify: `pkg/window/platform.go`
- Modify: `pkg/window/messages.go`
- Modify: `pkg/window/service.go`
- Modify: `pkg/window/mock_platform.go`
- Modify: `pkg/window/mock_test.go`
- Modify: `pkg/window/service_test.go`

- [ ] **Step 1: Add OnFileDrop to PlatformWindow interface in platform.go**

In `pkg/window/platform.go`, add to the `PlatformWindow` interface after the `OnWindowEvent` method:

```go
	// File drop
	OnFileDrop(handler func(paths []string, targetID string))
```

The full `PlatformWindow` interface becomes:

```go
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

	// File drop
	OnFileDrop(handler func(paths []string, targetID string))
}
```

- [ ] **Step 2: Add ActionFilesDropped to messages.go**

In `pkg/window/messages.go`, add after the existing action types:

```go
type ActionFilesDropped struct {
	Name     string   `json:"name"`     // window name
	Paths    []string `json:"paths"`
	TargetID string   `json:"targetId,omitempty"`
}
```

- [ ] **Step 3: Add OnFileDrop call in trackWindow() in service.go**

In `pkg/window/service.go`, add at the end of the `trackWindow()` method (after the `pw.OnWindowEvent` block):

```go
	pw.OnFileDrop(func(paths []string, targetID string) {
		_ = s.Core().ACTION(ActionFilesDropped{
			Name:     pw.Name(),
			Paths:    paths,
			TargetID: targetID,
		})
	})
```

The full `trackWindow` method becomes:

```go
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
```

- [ ] **Step 4: Add OnFileDrop to exported MockWindow in mock_platform.go**

In `pkg/window/mock_platform.go`, add a field and method to `MockWindow`:

Add field to MockWindow struct:

```go
	fileDropHandlers []func(paths []string, targetID string)
```

Add method:

```go
func (w *MockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}
```

The full `MockWindow` struct becomes:

```go
type MockWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, focused   bool
	visible, alwaysOnTop bool
	closed               bool
	eventHandlers        []func(WindowEvent)
	fileDropHandlers     []func(paths []string, targetID string)
}
```

And the full method set gains:

```go
func (w *MockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}
```

- [ ] **Step 5: Add OnFileDrop to unexported mockWindow in mock_test.go**

In `pkg/window/mock_test.go`, add a field and method to `mockWindow`:

Add field to mockWindow struct:

```go
	fileDropHandlers []func(paths []string, targetID string)
```

Add method:

```go
func (w *mockWindow) OnFileDrop(handler func(paths []string, targetID string)) {
	w.fileDropHandlers = append(w.fileDropHandlers, handler)
}
```

Add helper to simulate file drops:

```go
// emitFileDrop simulates a file drop on the window.
func (w *mockWindow) emitFileDrop(paths []string, targetID string) {
	for _, h := range w.fileDropHandlers {
		h(paths, targetID)
	}
}
```

The full `mockWindow` struct becomes:

```go
type mockWindow struct {
	name, title, url     string
	width, height, x, y  int
	maximised, focused   bool
	visible, alwaysOnTop bool
	closed               bool
	eventHandlers        []func(WindowEvent)
	fileDropHandlers     []func(paths []string, targetID string)
}
```

- [ ] **Step 6: Write failing file drop test in service_test.go**

Add to `pkg/window/service_test.go`:

```go
func TestFileDrop_Good(t *testing.T) {
	_, c := newTestWindowService(t)

	// Open a window
	result, _, _ := c.PERFORM(TaskOpenWindow{
		Opts: []WindowOption{WithName("drop-test")},
	})
	info := result.(WindowInfo)
	assert.Equal(t, "drop-test", info.Name)

	// Capture broadcast actions
	var dropped ActionFilesDropped
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionFilesDropped); ok {
			mu.Lock()
			dropped = a
			mu.Unlock()
		}
		return nil
	})

	// Get the mock window and simulate file drop
	svc := core.MustServiceFor[*Service](c, "window")
	pw, ok := svc.Manager().Get("drop-test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	mw.emitFileDrop([]string{"/tmp/file1.txt", "/tmp/file2.txt"}, "upload-zone")

	mu.Lock()
	assert.Equal(t, "drop-test", dropped.Name)
	assert.Equal(t, []string{"/tmp/file1.txt", "/tmp/file2.txt"}, dropped.Paths)
	assert.Equal(t, "upload-zone", dropped.TargetID)
	mu.Unlock()
}
```

Note: Add `"sync"` to the import block in `service_test.go` if not already present.

- [ ] **Step 7: Run test to verify it fails**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/ -v -run TestFileDrop`
Expected: FAIL — `OnFileDrop` method missing from interface / `ActionFilesDropped` not defined

- [ ] **Step 8: Apply all modifications from Steps 1-5**

Apply the changes described in Steps 1-5 (platform.go, messages.go, service.go, mock_platform.go, mock_test.go).

- [ ] **Step 9: Run all window tests to verify they pass**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/window/ -v`
Expected: PASS (all existing tests + new TestFileDrop_Good)

- [ ] **Step 10: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/window/
git commit -m "feat(window): add file drop support to PlatformWindow interface

Adds OnFileDrop(handler func(paths []string, targetID string)) to PlatformWindow.
trackWindow() now wires file drop callbacks to ActionFilesDropped broadcasts.
Updates both exported MockWindow and unexported mockWindow with the new method.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Display orchestrator updates

**Files:**
- Modify: `pkg/display/events.go`
- Modify: `pkg/display/display.go`

- [ ] **Step 1: Add new EventType constants in events.go**

In `pkg/display/events.go`, add two new constants after the existing ones:

```go
	EventKeybindingTriggered EventType = "keybinding.triggered"
	EventWindowFileDrop      EventType = "window.filedrop"
```

The full constants block becomes:

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
	EventKeybindingTriggered EventType = "keybinding.triggered"
	EventWindowFileDrop      EventType = "window.filedrop"
)
```

- [ ] **Step 2: Add new imports in display.go**

In `pkg/display/display.go`, add two new imports:

```go
	"forge.lthn.ai/core/gui/pkg/browser"
	"forge.lthn.ai/core/gui/pkg/keybinding"
```

The full import block becomes:

```go
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"forge.lthn.ai/core/go-config"
	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/browser"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/keybinding"
	"forge.lthn.ai/core/gui/pkg/menu"
	"forge.lthn.ai/core/gui/pkg/notification"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/systray"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/wailsapp/wails/v3/pkg/application"
)
```

- [ ] **Step 3: Add HandleIPCEvents cases in display.go**

In `pkg/display/display.go`, add two new cases to `HandleIPCEvents` before the closing `}` of the switch statement (after the `screen.ActionScreensChanged` case):

```go
	case keybinding.ActionTriggered:
		if s.events != nil {
			s.events.Emit(Event{Type: EventKeybindingTriggered,
				Data: map[string]any{"accelerator": m.Accelerator}})
		}
	case window.ActionFilesDropped:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowFileDrop, Window: m.Name,
				Data: map[string]any{"paths": m.Paths, "targetId": m.TargetID}})
		}
```

The full HandleIPCEvents method becomes:

```go
// HandleIPCEvents is auto-discovered and registered by core.WithService.
// It bridges sub-service IPC actions to WebSocket events for TS apps.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	switch m := msg.(type) {
	case core.ActionServiceStartup:
		// All services have completed OnStartup — safe to PERFORM on sub-services
		s.buildMenu()
		s.setupTray()
	case window.ActionWindowOpened:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowCreate, Window: m.Name,
				Data: map[string]any{"name": m.Name}})
		}
	case window.ActionWindowClosed:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowClose, Window: m.Name,
				Data: map[string]any{"name": m.Name}})
		}
	case window.ActionWindowMoved:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowMove, Window: m.Name,
				Data: map[string]any{"x": m.X, "y": m.Y}})
		}
	case window.ActionWindowResized:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowResize, Window: m.Name,
				Data: map[string]any{"w": m.W, "h": m.H}})
		}
	case window.ActionWindowFocused:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowFocus, Window: m.Name})
		}
	case window.ActionWindowBlurred:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowBlur, Window: m.Name})
		}
	case systray.ActionTrayClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventTrayClick})
		}
	case systray.ActionTrayMenuItemClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventTrayMenuItemClick,
				Data: map[string]any{"actionId": m.ActionID}})
		}
		s.handleTrayAction(m.ActionID)
	case environment.ActionThemeChanged:
		if s.events != nil {
			theme := "light"
			if m.IsDark {
				theme = "dark"
			}
			s.events.Emit(Event{Type: EventThemeChange,
				Data: map[string]any{"isDark": m.IsDark, "theme": theme}})
		}
	case notification.ActionNotificationClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventNotificationClick,
				Data: map[string]any{"id": m.ID}})
		}
	case screen.ActionScreensChanged:
		if s.events != nil {
			s.events.Emit(Event{Type: EventScreenChange,
				Data: map[string]any{"screens": m.Screens}})
		}
	case keybinding.ActionTriggered:
		if s.events != nil {
			s.events.Emit(Event{Type: EventKeybindingTriggered,
				Data: map[string]any{"accelerator": m.Accelerator}})
		}
	case window.ActionFilesDropped:
		if s.events != nil {
			s.events.Emit(Event{Type: EventWindowFileDrop, Window: m.Name,
				Data: map[string]any{"paths": m.Paths, "targetId": m.TargetID}})
		}
	}
	return nil
}
```

- [ ] **Step 4: Add handleWSMessage method to Service (if not already present)**

The display orchestrator needs a `handleWSMessage` method that bridges WebSocket commands to IPC. If this method does not yet exist, add it. If it already exists, add new cases to the switch.

Add to `pkg/display/display.go`:

```go
// WSMessage represents a command received from a WebSocket client.
type WSMessage struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data,omitempty"`
}

// handleWSMessage bridges WebSocket commands to IPC calls.
func (s *Service) handleWSMessage(msg WSMessage) (any, bool, error) {
	var result any
	var handled bool
	var err error

	switch msg.Action {
	case "keybinding:add":
		accelerator, _ := msg.Data["accelerator"].(string)
		description, _ := msg.Data["description"].(string)
		result, handled, err = s.Core().PERFORM(keybinding.TaskAdd{
			Accelerator: accelerator, Description: description,
		})
	case "keybinding:remove":
		accelerator, _ := msg.Data["accelerator"].(string)
		result, handled, err = s.Core().PERFORM(keybinding.TaskRemove{
			Accelerator: accelerator,
		})
	case "keybinding:list":
		result, handled, err = s.Core().QUERY(keybinding.QueryList{})
	case "browser:open-url":
		url, _ := msg.Data["url"].(string)
		result, handled, err = s.Core().PERFORM(browser.TaskOpenURL{URL: url})
	case "browser:open-file":
		path, _ := msg.Data["path"].(string)
		result, handled, err = s.Core().PERFORM(browser.TaskOpenFile{Path: path})
	default:
		return nil, false, nil
	}

	return result, handled, err
}
```

Note: The `WSMessage` type and `handleWSMessage` method may need to be integrated into an existing WS message dispatch mechanism. If the display package already has a WS command router, add the new cases to its switch statement instead of creating a new method. Use safe comma-ok type assertions throughout (never bare `msg.Data["key"].(string)`, always `msg.Data["key"].(string)` with the ok value discarded via `_`).

- [ ] **Step 5: Run all display tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/display/ -v`
Expected: PASS (verify the new imports compile and existing tests still pass)

- [ ] **Step 6: Run full test suite**

Run: `cd /Users/snider/Code/core/gui && go test ./... -v`
Expected: PASS (all packages)

- [ ] **Step 7: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/
git commit -m "feat(display): wire keybinding, browser, and file drop through orchestrator

Adds EventKeybindingTriggered and EventWindowFileDrop EventType constants.
HandleIPCEvents bridges keybinding.ActionTriggered and window.ActionFilesDropped
to WS events. handleWSMessage bridges WS commands to IPC for keybinding:add,
keybinding:remove, keybinding:list, browser:open-url, browser:open-file.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 5: Final verification and commit

- [ ] **Step 1: Run full test suite**

```bash
cd /Users/snider/Code/core/gui && go test ./... -v
```

Expected: All tests pass across all packages.

- [ ] **Step 2: Run linting**

```bash
cd /Users/snider/Code/core/gui && go vet ./...
```

Expected: No warnings.

- [ ] **Step 3: Verify no import cycles**

```bash
cd /Users/snider/Code/core/gui && go build ./...
```

Expected: Clean build. Dependency direction:
- `pkg/keybinding` and `pkg/browser` are independent (import only `core/go`)
- `pkg/display` imports `pkg/keybinding` and `pkg/browser` (message types only)
- No circular dependencies

- [ ] **Step 4: Review file inventory**

New files created (10):
- `pkg/keybinding/platform.go`
- `pkg/keybinding/messages.go`
- `pkg/keybinding/register.go`
- `pkg/keybinding/service.go`
- `pkg/keybinding/service_test.go`
- `pkg/browser/platform.go`
- `pkg/browser/messages.go`
- `pkg/browser/register.go`
- `pkg/browser/service.go`
- `pkg/browser/service_test.go`

Modified files (8):
- `pkg/window/platform.go` — `OnFileDrop` added to `PlatformWindow`
- `pkg/window/messages.go` — `ActionFilesDropped` added
- `pkg/window/service.go` — `trackWindow()` wires file drop
- `pkg/window/mock_platform.go` — `MockWindow.OnFileDrop` added
- `pkg/window/mock_test.go` — `mockWindow.OnFileDrop` + `emitFileDrop` added
- `pkg/window/service_test.go` — `TestFileDrop_Good` added
- `pkg/display/events.go` — 2 new `EventType` constants
- `pkg/display/display.go` — 2 new `HandleIPCEvents` cases + WS->IPC bridge + 2 new imports

- [ ] **Step 5: Final commit (if not already committed per-task)**

If individual task commits were made, no final commit is needed. If working in a single branch, squash or leave as individual commits per preference.

---

## Summary

| Task | Package | Tests | Key design decisions |
|------|---------|-------|---------------------|
| 1 | `pkg/keybinding` | 10 | In-memory `map[string]BindingInfo` registry; `ErrAlreadyRegistered` on duplicate; `QueryList` reads registry not `platform.GetAll()`; callback broadcasts `ActionTriggered` |
| 2 | `pkg/browser` | 6 | Stateless; tasks only, no queries/actions; direct platform delegation |
| 3 | `pkg/window` (enhance) | 1 new | `OnFileDrop` on `PlatformWindow` interface; `trackWindow()` wires to `ActionFilesDropped`; both mocks updated |
| 4 | `pkg/display` (wire) | 0 new (compile check) | 2 new `EventType` constants; 2 new `HandleIPCEvents` cases; 5 new WS->IPC cases with safe comma-ok assertions |
| 5 | Verification | full suite | `go test ./...`, `go vet ./...`, `go build ./...` |

## Licence

EUPL-1.2
