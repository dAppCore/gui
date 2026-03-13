# CoreGUI Spec D: Context Menus & Final Cleanup — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `pkg/contextmenu` as the final Wails v3 Manager API wrapper (completing full coverage), then remove all stale Wails wrapper types from `pkg/display/interfaces.go`, migrating the 4 remaining direct `s.app.*` calls to IPC. The `App` interface shrinks to `Quit()` + `Logger()` only.

**Architecture:** `pkg/contextmenu` follows the identical three-layer pattern (IPC Bus -> Service -> Platform Interface) established by `pkg/keybinding`. The display orchestrator gains HandleIPCEvents + WS bridge cases for context menu actions, a new `ActionIDECommand` message type, and loses ~75 lines of stale wrapper code.

**Tech Stack:** Go, core/go DI framework, Wails v3 (behind Platform interfaces)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-final-cleanup-design.md`

---

## File Structure

### New files

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/contextmenu/` | `platform.go` | Platform interface (4 methods) + `ContextMenuDef`, `MenuItemDef` types |
| | `messages.go` | IPC message types (QueryGet, QueryList, TaskAdd, TaskRemove, ActionItemClicked) |
| | `register.go` | `Register(Platform)` factory closure |
| | `service.go` | Service struct, OnStartup(), query/task handlers, in-memory registry |
| | `service_test.go` | Tests with mock platform (`_Good/_Bad/_Ugly` naming) |

### Modified files

| File | Change |
|------|--------|
| `pkg/display/display.go` | Add `contextmenu.ActionItemClicked` + `ActionIDECommand` HandleIPCEvents cases, add `contextmenu:*` WS->IPC cases, migrate `handleOpenFile`/`handleSaveFile`/`handleRun`/`handleBuild` from direct `s.app.*` to IPC |
| `pkg/display/events.go` | Add `EventContextMenuClick` + `EventIDECommand` constants |
| `pkg/display/interfaces.go` | Remove `DialogManager`, `EnvManager`, `EventManager` interfaces + `wailsDialogManager`, `wailsEnvManager`, `wailsEventManager` structs + `App.Dialog()`, `App.Env()`, `App.Event()` methods. Shrink `App` to `Quit()` + `Logger()` |

### New file in display

| File | Responsibility |
|------|---------------|
| `pkg/display/messages.go` | `ActionIDECommand` message type |

---

## Task 1: Create pkg/contextmenu

**Files:**
- Create: `pkg/contextmenu/platform.go`
- Create: `pkg/contextmenu/messages.go`
- Create: `pkg/contextmenu/register.go`
- Create: `pkg/contextmenu/service.go`
- Test: `pkg/contextmenu/service_test.go`

### Step 1: Create platform.go with types

- [ ] **Create `pkg/contextmenu/platform.go`**

```go
// pkg/contextmenu/platform.go
package contextmenu

// Platform abstracts the context menu backend (Wails v3).
// The Add callback must broadcast ActionItemClicked via s.Core().ACTION()
// when a menu item is clicked — the adapter translates MenuItemDef.ActionID
// to a callback that does this.
type Platform interface {
	// Add registers a context menu by name.
	// The onItemClick callback is called with (menuName, actionID, data)
	// when any item in the menu is clicked. The adapter creates per-item
	// OnClick handlers that call this with the appropriate ActionID.
	Add(name string, menu ContextMenuDef, onItemClick func(menuName, actionID, data string)) error

	// Remove unregisters a context menu by name.
	Remove(name string) error

	// Get returns a context menu definition by name, or false if not found.
	Get(name string) (*ContextMenuDef, bool)

	// GetAll returns all registered context menu definitions.
	GetAll() map[string]ContextMenuDef
}

// ContextMenuDef describes a context menu and its items.
type ContextMenuDef struct {
	Name  string        `json:"name"`
	Items []MenuItemDef `json:"items"`
}

// MenuItemDef describes a single item in a context menu.
// Items may be nested (submenu children via Items field).
type MenuItemDef struct {
	Label       string        `json:"label"`
	Type        string        `json:"type,omitempty"`        // "" (normal), "separator", "checkbox", "radio", "submenu"
	Accelerator string        `json:"accelerator,omitempty"`
	Enabled     *bool         `json:"enabled,omitempty"`     // nil = true (default)
	Checked     bool          `json:"checked,omitempty"`
	ActionID    string        `json:"actionId,omitempty"`    // identifies which item was clicked
	Items       []MenuItemDef `json:"items,omitempty"`       // submenu children (recursive)
}
```

### Step 2: Create messages.go

- [ ] **Create `pkg/contextmenu/messages.go`**

```go
// pkg/contextmenu/messages.go
package contextmenu

import "errors"

// ErrMenuNotFound is returned when attempting to remove or get a menu
// that does not exist in the registry.
var ErrMenuNotFound = errors.New("contextmenu: menu not found")

// --- Queries ---

// QueryGet returns a single context menu by name. Result: *ContextMenuDef (nil if not found)
type QueryGet struct {
	Name string `json:"name"`
}

// QueryList returns all registered context menus. Result: map[string]ContextMenuDef
type QueryList struct{}

// --- Tasks ---

// TaskAdd registers a context menu. Result: nil
// If a menu with the same name already exists it is replaced (remove + re-add).
type TaskAdd struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskRemove unregisters a context menu. Result: nil
// Returns ErrMenuNotFound if the menu does not exist.
type TaskRemove struct {
	Name string `json:"name"`
}

// --- Actions ---

// ActionItemClicked is broadcast when a context menu item is clicked.
// The Data field is populated from the CSS --custom-contextmenu-data property
// on the element that triggered the context menu.
type ActionItemClicked struct {
	MenuName string `json:"menuName"`
	ActionID string `json:"actionId"`
	Data     string `json:"data,omitempty"`
}
```

### Step 3: Create register.go

- [ ] **Create `pkg/contextmenu/register.go`**

```go
// pkg/contextmenu/register.go
package contextmenu

import "forge.lthn.ai/core/go/pkg/core"

// Register creates a factory closure that captures the Platform adapter.
// The returned function has the signature WithService requires: func(*Core) (any, error).
func Register(p Platform) func(*core.Core) (any, error) {
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
			platform:       p,
			menus:          make(map[string]ContextMenuDef),
		}, nil
	}
}
```

### Step 4: Write failing test

- [ ] **Create `pkg/contextmenu/service_test.go`**

```go
// pkg/contextmenu/service_test.go
package contextmenu

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPlatform records Add/Remove calls and allows simulating clicks.
type mockPlatform struct {
	mu            sync.Mutex
	menus         map[string]ContextMenuDef
	clickHandlers map[string]func(menuName, actionID, data string)
	removed       []string
	addErr        error
	removeErr     error
}

func newMockPlatform() *mockPlatform {
	return &mockPlatform{
		menus:         make(map[string]ContextMenuDef),
		clickHandlers: make(map[string]func(menuName, actionID, data string)),
	}
}

func (m *mockPlatform) Add(name string, menu ContextMenuDef, onItemClick func(string, string, string)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.addErr != nil {
		return m.addErr
	}
	m.menus[name] = menu
	m.clickHandlers[name] = onItemClick
	return nil
}

func (m *mockPlatform) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.removeErr != nil {
		return m.removeErr
	}
	delete(m.menus, name)
	delete(m.clickHandlers, name)
	m.removed = append(m.removed, name)
	return nil
}

func (m *mockPlatform) Get(name string) (*ContextMenuDef, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	menu, ok := m.menus[name]
	if !ok {
		return nil, false
	}
	return &menu, true
}

func (m *mockPlatform) GetAll() map[string]ContextMenuDef {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]ContextMenuDef, len(m.menus))
	for k, v := range m.menus {
		out[k] = v
	}
	return out
}

// simulateClick simulates a context menu item click by calling the registered handler.
func (m *mockPlatform) simulateClick(menuName, actionID, data string) {
	m.mu.Lock()
	h, ok := m.clickHandlers[menuName]
	m.mu.Unlock()
	if ok {
		h(menuName, actionID, data)
	}
}

func newTestContextMenuService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "contextmenu")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	mp := newMockPlatform()
	svc, _ := newTestContextMenuService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskAdd_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{
			Name: "file-menu",
			Items: []MenuItemDef{
				{Label: "Open", ActionID: "open"},
				{Label: "Delete", ActionID: "delete"},
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify menu registered on platform
	_, ok := mp.Get("file-menu")
	assert.True(t, ok)
}

func TestTaskAdd_Good_ReplaceExisting(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add initial menu
	_, _, _ = c.PERFORM(TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "A", ActionID: "a"}}},
	})

	// Replace with new menu
	_, handled, err := c.PERFORM(TaskAdd{
		Name: "ctx",
		Menu: ContextMenuDef{Name: "ctx", Items: []MenuItemDef{{Label: "B", ActionID: "b"}}},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify registry has new menu
	result, _, _ := c.QUERY(QueryGet{Name: "ctx"})
	def := result.(*ContextMenuDef)
	require.Len(t, def.Items, 1)
	assert.Equal(t, "B", def.Items[0].Label)
}

func TestTaskRemove_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Add then remove
	_, _, _ = c.PERFORM(TaskAdd{
		Name: "test",
		Menu: ContextMenuDef{Name: "test"},
	})
	_, handled, err := c.PERFORM(TaskRemove{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify removed from registry
	result, _, _ := c.QUERY(QueryGet{Name: "test"})
	assert.Nil(t, result)
}

func TestTaskRemove_Bad_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskRemove{Name: "nonexistent"})
	assert.True(t, handled)
	assert.ErrorIs(t, err, ErrMenuNotFound)
}

func TestQueryGet_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "my-menu",
		Menu: ContextMenuDef{
			Name:  "my-menu",
			Items: []MenuItemDef{{Label: "Edit", ActionID: "edit"}},
		},
	})

	result, handled, err := c.QUERY(QueryGet{Name: "my-menu"})
	require.NoError(t, err)
	assert.True(t, handled)
	def := result.(*ContextMenuDef)
	assert.Equal(t, "my-menu", def.Name)
	assert.Len(t, def.Items, 1)
}

func TestQueryGet_Good_NotFound(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	result, handled, err := c.QUERY(QueryGet{Name: "missing"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryList_Good(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, _, _ = c.PERFORM(TaskAdd{Name: "a", Menu: ContextMenuDef{Name: "a"}})
	_, _, _ = c.PERFORM(TaskAdd{Name: "b", Menu: ContextMenuDef{Name: "b"}})

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.(map[string]ContextMenuDef)
	assert.Len(t, list, 2)
}

func TestQueryList_Good_Empty(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	result, handled, err := c.QUERY(QueryList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.(map[string]ContextMenuDef)
	assert.Len(t, list, 0)
}

func TestTaskAdd_Good_ClickBroadcast(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	// Capture broadcast actions
	var clicked ActionItemClicked
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionItemClicked); ok {
			mu.Lock()
			clicked = a
			mu.Unlock()
		}
		return nil
	})

	_, _, _ = c.PERFORM(TaskAdd{
		Name: "file-menu",
		Menu: ContextMenuDef{
			Name: "file-menu",
			Items: []MenuItemDef{
				{Label: "Open", ActionID: "open"},
			},
		},
	})

	// Simulate click via mock
	mp.simulateClick("file-menu", "open", "file-123")

	mu.Lock()
	assert.Equal(t, "file-menu", clicked.MenuName)
	assert.Equal(t, "open", clicked.ActionID)
	assert.Equal(t, "file-123", clicked.Data)
	mu.Unlock()
}

func TestTaskAdd_Good_SubmenuItems(t *testing.T) {
	mp := newMockPlatform()
	_, c := newTestContextMenuService(t, mp)

	_, handled, err := c.PERFORM(TaskAdd{
		Name: "nested",
		Menu: ContextMenuDef{
			Name: "nested",
			Items: []MenuItemDef{
				{Label: "File", Type: "submenu", Items: []MenuItemDef{
					{Label: "New", ActionID: "new"},
					{Label: "Open", ActionID: "open"},
				}},
				{Type: "separator"},
				{Label: "Quit", ActionID: "quit"},
			},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryGet{Name: "nested"})
	def := result.(*ContextMenuDef)
	assert.Len(t, def.Items, 3)
	assert.Len(t, def.Items[0].Items, 2) // submenu children
}

func TestQueryList_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryList{})
	assert.False(t, handled)
}
```

- [ ] **Run test to verify it fails**

```bash
cd /Users/snider/Code/core/gui && go test ./pkg/contextmenu/ -v
```

Expected: FAIL -- `Service` type not defined, `Options` type not defined

### Step 5: Create service.go

- [ ] **Create `pkg/contextmenu/service.go`**

```go
// pkg/contextmenu/service.go
package contextmenu

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
)

// Options holds configuration for the context menu service.
type Options struct{}

// Service is a core.Service managing context menus via IPC.
// It maintains an in-memory registry of menus (map[string]ContextMenuDef)
// and delegates platform-level registration to the Platform interface.
type Service struct {
	*core.ServiceRuntime[Options]
	platform Platform
	menus    map[string]ContextMenuDef
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
	switch q := q.(type) {
	case QueryGet:
		return s.queryGet(q), true, nil
	case QueryList:
		return s.queryList(), true, nil
	default:
		return nil, false, nil
	}
}

// queryGet returns a single menu definition by name, or nil if not found.
func (s *Service) queryGet(q QueryGet) *ContextMenuDef {
	menu, ok := s.menus[q.Name]
	if !ok {
		return nil
	}
	return &menu
}

// queryList returns a copy of all registered menus.
func (s *Service) queryList() map[string]ContextMenuDef {
	result := make(map[string]ContextMenuDef, len(s.menus))
	for k, v := range s.menus {
		result[k] = v
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
	// If menu already exists, remove it first (replace semantics)
	if _, exists := s.menus[t.Name]; exists {
		_ = s.platform.Remove(t.Name)
		delete(s.menus, t.Name)
	}

	// Register on platform with a callback that broadcasts ActionItemClicked
	err := s.platform.Add(t.Name, t.Menu, func(menuName, actionID, data string) {
		_ = s.Core().ACTION(ActionItemClicked{
			MenuName: menuName,
			ActionID: actionID,
			Data:     data,
		})
	})
	if err != nil {
		return fmt.Errorf("contextmenu: platform add failed: %w", err)
	}

	s.menus[t.Name] = t.Menu
	return nil
}

func (s *Service) taskRemove(t TaskRemove) error {
	if _, exists := s.menus[t.Name]; !exists {
		return ErrMenuNotFound
	}

	err := s.platform.Remove(t.Name)
	if err != nil {
		return fmt.Errorf("contextmenu: platform remove failed: %w", err)
	}

	delete(s.menus, t.Name)
	return nil
}
```

- [ ] **Run tests to verify they pass**

```bash
cd /Users/snider/Code/core/gui && go test ./pkg/contextmenu/ -v
```

Expected: PASS (13 tests)

- [ ] **Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/contextmenu/
git commit -m "feat(contextmenu): add contextmenu core.Service with Platform interface and IPC

Completes full Wails v3 Manager API coverage through the IPC bus.
Service maintains in-memory registry, delegates to Platform for native
context menu operations, broadcasts ActionItemClicked on menu item clicks.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Display orchestrator — add contextmenu integration + IDE command message

**Files:**
- Create: `pkg/display/messages.go`
- Modify: `pkg/display/events.go`
- Modify: `pkg/display/display.go`

### Step 1: Create messages.go with ActionIDECommand

- [ ] **Create `pkg/display/messages.go`**

```go
// pkg/display/messages.go
package display

// ActionIDECommand is broadcast when a menu handler triggers an IDE command
// (save, run, build). Replaces direct s.app.Event().Emit("ide:*") calls.
// Listeners (e.g. editor windows) handle this via HandleIPCEvents.
type ActionIDECommand struct {
	Command string `json:"command"` // "save", "run", "build"
}

// EventIDECommand is the WS event type for IDE commands.
const EventIDECommand EventType = "ide.command"
```

### Step 2: Add EventContextMenuClick to events.go

- [ ] **Add constant to `pkg/display/events.go`**

In the `const` block, add after `EventSystemResume`:

```go
EventContextMenuClick   EventType = "contextmenu.item-clicked"
```

The full addition is a single line in the existing `const ( ... )` block, placed after `EventSystemResume EventType = "system.resume"`:

```go
EventContextMenuClick   EventType = "contextmenu.item-clicked"
```

### Step 3: Add HandleIPCEvents cases for contextmenu + IDE command

- [ ] **Modify `pkg/display/display.go` — add import for contextmenu**

Add to the import block:

```go
"forge.lthn.ai/core/gui/pkg/contextmenu"
```

- [ ] **Modify `pkg/display/display.go` — add HandleIPCEvents cases**

In the `HandleIPCEvents` method, add these two cases inside the `switch m := msg.(type)` block, before the closing `}`:

```go
	case contextmenu.ActionItemClicked:
		if s.events != nil {
			s.events.Emit(Event{Type: EventContextMenuClick,
				Data: map[string]any{
					"menuName": m.MenuName,
					"actionId": m.ActionID,
					"data":     m.Data,
				}})
		}
	case ActionIDECommand:
		if s.events != nil {
			s.events.Emit(Event{Type: EventIDECommand,
				Data: map[string]any{"command": m.Command}})
		}
```

### Step 4: Add WS->IPC cases for contextmenu

- [ ] **Modify `pkg/display/display.go` — add contextmenu WS cases in `handleWSMessage`**

Add to the import block (if not already present):

```go
"encoding/json"
```

In the `handleWSMessage` method's `switch msg.Action` block, add before the `default:` case:

```go
	case "contextmenu:add":
		name, _ := msg.Data["name"].(string)
		menuJSON, _ := json.Marshal(msg.Data["menu"])
		var menuDef contextmenu.ContextMenuDef
		_ = json.Unmarshal(menuJSON, &menuDef)
		result, handled, err = s.Core().PERFORM(contextmenu.TaskAdd{
			Name: name, Menu: menuDef,
		})
	case "contextmenu:remove":
		name, _ := msg.Data["name"].(string)
		result, handled, err = s.Core().PERFORM(contextmenu.TaskRemove{Name: name})
	case "contextmenu:get":
		name, _ := msg.Data["name"].(string)
		result, handled, err = s.Core().QUERY(contextmenu.QueryGet{Name: name})
	case "contextmenu:list":
		result, handled, err = s.Core().QUERY(contextmenu.QueryList{})
```

**Note on `encoding/json` import:** The `display.go` file does not currently import `encoding/json`. However, the `events.go` file (same package) already uses `encoding/json`. Since this is the same package, the import is available. But since `display.go` is a separate file, you MUST add `"encoding/json"` to its import block. Check whether it is already present before adding.

- [ ] **Run `go vet` to verify no compilation errors**

```bash
cd /Users/snider/Code/core/gui && go vet ./pkg/display/
```

Expected: PASS (no errors)

- [ ] **Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/messages.go pkg/display/events.go pkg/display/display.go
git commit -m "feat(display): add contextmenu integration and ActionIDECommand to orchestrator

Add HandleIPCEvents cases for contextmenu.ActionItemClicked and
ActionIDECommand, WS->IPC bridge cases for contextmenu:add/remove/get/list,
EventContextMenuClick and EventIDECommand constants.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Display cleanup — migrate stale calls, remove wrappers

**Files:**
- Modify: `pkg/display/display.go`
- Modify: `pkg/display/interfaces.go`

### Step 1: Migrate handleSaveFile, handleRun, handleBuild to IPC

- [ ] **Modify `pkg/display/display.go` — replace 3 direct `s.app.Event().Emit()` calls**

Replace these three method bodies:

**Before (line ~829):**
```go
func (s *Service) handleSaveFile()   { s.app.Event().Emit("ide:save") }
```

**After:**
```go
func (s *Service) handleSaveFile() { _ = s.Core().ACTION(ActionIDECommand{Command: "save"}) }
```

**Before (line ~850):**
```go
func (s *Service) handleRun()   { s.app.Event().Emit("ide:run") }
```

**After:**
```go
func (s *Service) handleRun() { _ = s.Core().ACTION(ActionIDECommand{Command: "run"}) }
```

**Before (line ~851):**
```go
func (s *Service) handleBuild() { s.app.Event().Emit("ide:build") }
```

**After:**
```go
func (s *Service) handleBuild() { _ = s.Core().ACTION(ActionIDECommand{Command: "build"}) }
```

### Step 2: Migrate handleOpenFile to use dialog.TaskOpenFile IPC

- [ ] **Modify `pkg/display/display.go` — replace handleOpenFile**

**Before (lines ~810-827):**
```go
func (s *Service) handleOpenFile() {
	dialog := s.app.Dialog().OpenFile()
	dialog.SetTitle("Open File")
	dialog.CanChooseFiles(true)
	dialog.CanChooseDirectories(false)
	result, err := dialog.PromptForSingleSelection()
	if err != nil || result == "" {
		return
	}
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle(result + " - Editor"),
			window.WithURL("/#/developer/editor?file=" + result),
			window.WithSize(1200, 800),
		},
	})
}
```

**After:**
```go
func (s *Service) handleOpenFile() {
	result, handled, err := s.Core().PERFORM(dialog.TaskOpenFile{
		Opts: dialog.OpenFileOptions{
			Title:         "Open File",
			AllowMultiple: false,
		},
	})
	if err != nil || !handled {
		return
	}
	paths, ok := result.([]string)
	if !ok || len(paths) == 0 {
		return
	}
	_, _, _ = s.Core().PERFORM(window.TaskOpenWindow{
		Opts: []window.WindowOption{
			window.WithName("editor"),
			window.WithTitle(paths[0] + " - Editor"),
			window.WithURL("/#/developer/editor?file=" + paths[0]),
			window.WithSize(1200, 800),
		},
	})
}
```

**Verify import:** `"forge.lthn.ai/core/gui/pkg/dialog"` must be in the import block. It is already imported on the current `display.go` (used in `handleTrayAction` for `dialog.TaskMessageDialog`).

### Step 3: Remove stale wrappers from interfaces.go

- [ ] **Rewrite `pkg/display/interfaces.go`**

Replace the entire file contents with:

```go
// pkg/display/interfaces.go
package display

import "github.com/wailsapp/wails/v3/pkg/application"

// App abstracts the Wails application for the orchestrator.
// After Spec D cleanup, only Quit() and Logger() remain —
// all other Wails Manager APIs are accessed via IPC.
type App interface {
	Logger() Logger
	Quit()
}

// Logger wraps Wails logging.
type Logger interface {
	Info(message string, args ...any)
}

// wailsApp wraps *application.App for the App interface.
type wailsApp struct {
	app *application.App
}

func newWailsApp(app *application.App) App {
	return &wailsApp{app: app}
}

func (w *wailsApp) Logger() Logger { return w.app.Logger }
func (w *wailsApp) Quit()          { w.app.Quit() }
```

**Removed types:**
- `DialogManager` interface (3 methods)
- `EnvManager` interface (2 methods)
- `EventManager` interface (2 methods)
- `wailsDialogManager` struct + 3 methods
- `wailsEnvManager` struct + 2 methods
- `wailsEventManager` struct + 2 methods
- `App.Dialog()`, `App.Env()`, `App.Event()` method requirements

**Removed imports:**
- `"github.com/wailsapp/wails/v3/pkg/events"` (was only used by `wailsEventManager.OnApplicationEvent`)

### Step 4: Verify no remaining references to removed methods

- [ ] **Check for any remaining `s.app.Dialog()`, `s.app.Env()`, `s.app.Event()` calls**

```bash
cd /Users/snider/Code/core/gui && grep -rn 's\.app\.Dialog()\|s\.app\.Env()\|s\.app\.Event()' pkg/display/
```

Expected: No matches. If there are matches, they indicate missed migrations and must be fixed before proceeding.

- [ ] **Run `go vet` to verify compilation**

```bash
cd /Users/snider/Code/core/gui && go vet ./pkg/display/
```

Expected: PASS. If it fails due to unused imports (e.g. the `events` package), remove them. If it fails due to `s.app.Dialog()` or similar calls elsewhere in display.go, those are migrations missed in Steps 1-2 and must be addressed.

- [ ] **Run all tests**

```bash
cd /Users/snider/Code/core/gui && go test ./pkg/... -v
```

Expected: All tests PASS across all packages.

- [ ] **Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/display.go pkg/display/interfaces.go
git commit -m "refactor(display): migrate stale Wails calls to IPC, remove wrapper types

Migrate handleOpenFile to dialog.TaskOpenFile IPC, handleSaveFile/handleRun/
handleBuild to ActionIDECommand IPC. Remove DialogManager, EnvManager,
EventManager interfaces and wailsDialogManager, wailsEnvManager,
wailsEventManager adapter structs. App interface now has Quit() + Logger() only.

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Final verification and commit

**Files:**
- No new files — verification only

### Step 1: Full test suite

- [ ] **Run full test suite**

```bash
cd /Users/snider/Code/core/gui && go test ./pkg/... -v -count=1
```

Expected: All tests PASS.

### Step 2: Vet and lint

- [ ] **Run go vet on entire module**

```bash
cd /Users/snider/Code/core/gui && go vet ./...
```

Expected: PASS.

### Step 3: Verify no remaining Wails imports leak through display interfaces

- [ ] **Check that `pkg/display/interfaces.go` only imports `application` (not `events`)**

```bash
cd /Users/snider/Code/core/gui && grep -n 'events' pkg/display/interfaces.go
```

Expected: No matches.

### Step 4: Verify App interface surface

- [ ] **Confirm App interface is exactly `Quit() + Logger()`**

```bash
cd /Users/snider/Code/core/gui && grep -A5 'type App interface' pkg/display/interfaces.go
```

Expected output:
```
type App interface {
    Logger() Logger
    Quit()
}
```

### Step 5: Count Wails Manager API coverage

- [ ] **Verify all 11 Wails v3 Manager APIs have core.Service wrappers**

| # | Wails Manager | Package | Spec |
|---|--------------|---------|------|
| 1 | Window | `pkg/window` | A (display split) |
| 2 | Systray | `pkg/systray` | A (display split) |
| 3 | Menu | `pkg/menu` | A (display split) |
| 4 | Clipboard | `pkg/clipboard` | B (extract & insulate) |
| 5 | Dialog | `pkg/dialog` | B (extract & insulate) |
| 6 | Notification | `pkg/notification` | B (extract & insulate) |
| 7 | Environment | `pkg/environment` | B (extract & insulate) |
| 8 | Screen | `pkg/screen` | B (extract & insulate) |
| 9 | Keybinding | `pkg/keybinding` | C (new input services) |
| 10 | Browser | `pkg/browser` | C (new input services) |
| 11 | ContextMenu | `pkg/contextmenu` | D (this spec) |

All 11 covered.

---

## Summary of Changes

### Lines added (approx)
| File | Lines |
|------|-------|
| `pkg/contextmenu/platform.go` | ~40 |
| `pkg/contextmenu/messages.go` | ~35 |
| `pkg/contextmenu/register.go` | ~15 |
| `pkg/contextmenu/service.go` | ~95 |
| `pkg/contextmenu/service_test.go` | ~230 |
| `pkg/display/messages.go` | ~10 |
| **Total new** | **~425** |

### Lines modified
| File | Change |
|------|--------|
| `pkg/display/events.go` | +1 line (EventContextMenuClick) |
| `pkg/display/display.go` | +30 lines (HandleIPCEvents cases, WS cases), -20 lines (migrated handlers) |

### Lines removed (approx)
| File | Lines removed |
|------|--------------|
| `pkg/display/interfaces.go` | ~55 lines (reduced from ~78 to ~23) |

### Net effect
- **+425 lines** new contextmenu package (5 files)
- **+10 lines** new display/messages.go
- **+11 lines** net in display.go (new cases minus migrated code)
- **-55 lines** removed from interfaces.go
- **Net: ~+391 lines**, all stale Wails wrappers eliminated, full Manager API coverage achieved
