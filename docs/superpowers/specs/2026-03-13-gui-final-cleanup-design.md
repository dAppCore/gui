# CoreGUI Spec D: Context Menus & Final Cleanup

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Add context menu package as core.Service, remove stale Wails wrappers from display orchestrator

## Context

Specs A–C extracted 12 packages covering 10 of the 11 Wails v3 Manager APIs. Two gaps remain:

1. **`app.ContextMenu`** — the only Wails v3 manager without a core.Service wrapper
2. **Stale `interfaces.go`** — display orchestrator still holds `wailsDialogManager`, `wailsEnvManager`, `wailsEventManager` wrappers and 4 direct `s.app.*` calls that bypass IPC

This spec closes both gaps, achieving full Wails v3 Manager API coverage through the IPC bus.

## Architecture

### New Package

| Package | Platform methods | IPC pattern |
|---------|-----------------|-------------|
| `pkg/contextmenu` | `Add()`, `Remove()`, `Get()`, `GetAll()` | Tasks (register/remove), Query (get/list), Action (item clicked) |

### Display Cleanup

Remove the `App` interface's `Dialog()`, `Env()`, `Event()` methods and their Wails adapter implementations. Migrate the 4 remaining direct calls to use IPC. The `App` interface reduces to `Quit()` and `Logger()` only.

---

## Package Design

### 1. pkg/contextmenu

**Platform interface:**

```go
type Platform interface {
    Add(name string, menu ContextMenuDef) error
    Remove(name string) error
    Get(name string) (*ContextMenuDef, bool)
    GetAll() map[string]ContextMenuDef
}
```

**Our own types:**

```go
type ContextMenuDef struct {
    Name  string         `json:"name"`
    Items []MenuItemDef  `json:"items"`
}

type MenuItemDef struct {
    Label       string        `json:"label"`
    Type        string        `json:"type,omitempty"`        // "" (normal), "separator", "checkbox", "radio", "submenu"
    Accelerator string        `json:"accelerator,omitempty"`
    Enabled     *bool         `json:"enabled,omitempty"`     // nil = true (default)
    Checked     bool          `json:"checked,omitempty"`
    ActionID    string        `json:"actionId,omitempty"`    // identifies which item was clicked
    Items       []MenuItemDef `json:"items,omitempty"`       // submenu children
}
```

**IPC messages:**

```go
// Queries
type QueryGet struct{ Name string }      // → *ContextMenuDef (nil if not found)
type QueryList struct{}                   // → map[string]ContextMenuDef

// Tasks
type TaskAdd struct {
    Name string         `json:"name"`
    Menu ContextMenuDef `json:"menu"`
}  // → error

type TaskRemove struct {
    Name string `json:"name"`
}  // → error

// Actions
type ActionItemClicked struct {
    MenuName string `json:"menuName"`
    ActionID string `json:"actionId"`
    Data     string `json:"data,omitempty"`  // from --custom-contextmenu-data CSS
}
```

**Service logic:** The Service maintains a `map[string]ContextMenuDef` registry (mirroring the Wails `contextMenus` map). `TaskAdd` calls `platform.Add(name, menu)` — the Wails adapter translates `ContextMenuDef` to `*application.ContextMenu` with `OnClick` callbacks that broadcast `ActionItemClicked` via `s.Core().ACTION()`. `TaskRemove` calls `platform.Remove()` and deletes from registry. `QueryGet` and `QueryList` read from registry.

**Callback bridging:** The Wails adapter's `OnClick` handler receives `*application.Context` which provides `ContextMenuData()`. The adapter maps each `MenuItemDef.ActionID` to a callback that broadcasts `ActionItemClicked{MenuName, ActionID, ctx.ContextMenuData()}`.

**Wails adapter:** Wraps `app.ContextMenu.Add()`, `app.ContextMenu.Remove()`, `app.ContextMenu.Get()`, `app.ContextMenu.GetAll()`. Translates between our `ContextMenuDef`/`MenuItemDef` types and Wails `*ContextMenu`/`*MenuItem`.

**Config:** None — stateless.

**WS bridge events:** `contextmenu.item-clicked` (broadcast with menuName, actionId, data). TS apps call `contextmenu:add`, `contextmenu:remove`, `contextmenu:get`, `contextmenu:list` via WS→IPC.

---

## Display Orchestrator Changes

### New HandleIPCEvents Case

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
```

### New EventType Constant

```go
EventContextMenuClick EventType = "contextmenu.item-clicked"
```

### New WS→IPC Cases

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

### IDE Command Migration

Replace direct `s.app.Event().Emit()` calls with IPC Actions:

```go
// New Action in pkg/display/messages.go
type ActionIDECommand struct {
    Command string `json:"command"` // "save", "run", "build"
}
```

```go
// Before (stale):
func (s *Service) handleSaveFile() { s.app.Event().Emit("ide:save") }
func (s *Service) handleRun()      { s.app.Event().Emit("ide:run") }
func (s *Service) handleBuild()    { s.app.Event().Emit("ide:build") }

// After (IPC):
func (s *Service) handleSaveFile() { _ = s.Core().ACTION(ActionIDECommand{Command: "save"}) }
func (s *Service) handleRun()      { _ = s.Core().ACTION(ActionIDECommand{Command: "run"}) }
func (s *Service) handleBuild()    { _ = s.Core().ACTION(ActionIDECommand{Command: "build"}) }
```

The `HandleIPCEvents` bridges these to WS:

```go
case ActionIDECommand:
    if s.events != nil {
        s.events.Emit(Event{Type: EventIDECommand,
            Data: map[string]any{"command": m.Command}})
    }
```

### handleOpenFile Migration

```go
// Before (stale — uses s.app.Dialog() directly):
func (s *Service) handleOpenFile() {
    dialog := s.app.Dialog().OpenFile()
    dialog.SetTitle("Open File")
    // ...
}

// After (IPC):
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

### Removed Types

After migration, remove from `interfaces.go`:
- `DialogManager` interface
- `EnvManager` interface
- `EventManager` interface
- `wailsDialogManager` struct + methods
- `wailsEnvManager` struct + methods
- `wailsEventManager` struct + methods
- `App.Dialog()`, `App.Env()`, `App.Event()` methods from interface

The `App` interface reduces to:

```go
type App interface {
    Logger() Logger
    Quit()
}
```

## Dependency Direction

```
pkg/display (orchestrator)
├── imports pkg/contextmenu (message types only)
├── imports pkg/dialog (message types — for handleOpenFile migration)
├── ... existing imports ...

pkg/contextmenu (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface
```

No circular dependencies.

## Testing Strategy

```go
func TestContextMenu_AddAndClick_Good(t *testing.T) {
    mock := &mockPlatform{}
    c, _ := core.New(
        core.WithService(contextmenu.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, err := c.PERFORM(contextmenu.TaskAdd{
        Name: "file-menu",
        Menu: contextmenu.ContextMenuDef{
            Name: "file-menu",
            Items: []contextmenu.MenuItemDef{
                {Label: "Open", ActionID: "open"},
                {Label: "Delete", ActionID: "delete"},
            },
        },
    })
    require.NoError(t, err)
    assert.True(t, handled)

    // Simulate click via mock
    mock.simulateClick("file-menu", "open", "file-123")

    // Verify action was broadcast
}

func TestContextMenu_Remove_Good(t *testing.T) {
    mock := &mockPlatform{}
    c, _ := core.New(
        core.WithService(contextmenu.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    // Add then remove
    _, _, _ = c.PERFORM(contextmenu.TaskAdd{Name: "test", Menu: contextmenu.ContextMenuDef{Name: "test"}})
    _, handled, err := c.PERFORM(contextmenu.TaskRemove{Name: "test"})
    require.NoError(t, err)
    assert.True(t, handled)

    // Verify removed
    result, _, _ := c.QUERY(contextmenu.QueryGet{Name: "test"})
    assert.Nil(t, result)
}
```

## Deferred Work

- **Dynamic menu updates**: `TaskUpdate` to modify items on an existing menu (SetEnabled, SetLabel, SetChecked). Currently requires remove + re-add.
- **Per-window context menus**: Context menus are global — scoping to specific windows would require window ID in the registration.

## Licence

EUPL-1.2
