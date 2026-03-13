# CoreGUI Spec B: New Input

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Add keybinding and browser packages as core.Services, enhance window service with file drop events

## Context

Spec A extracted 5 packages from display. This spec adds input capabilities not yet present in core/gui: keyboard shortcuts, browser delegation, and file drop events. Each follows the three-layer pattern (IPC Bus → Service → Platform Interface).

## Architecture

### New Packages

| Package | Platform methods | IPC pattern |
|---------|-----------------|-------------|
| `pkg/keybinding` | `Add()`, `Remove()`, `GetAll()` | Task add/remove, Query list, Action triggered |
| `pkg/browser` | `OpenURL()`, `OpenFile()` | Tasks (side-effects) |

### Window Service Enhancement

File drop is not a separate package — it extends `pkg/window` with an `OnFileDrop` callback on `PlatformWindow` and a new `ActionFilesDropped` broadcast.

---

## Package Designs

### 1. pkg/keybinding

**Platform interface:**

```go
type Platform interface {
    Add(accelerator string, handler func()) error
    Remove(accelerator string) error
    GetAll() []string
}
```

Platform-aware accelerator syntax: `Cmd+S` (macOS), `Ctrl+S` (Windows/Linux). Special keys: `F1-F12`, `Escape`, `Enter`, `Space`, `Tab`, `Backspace`, `Delete`, arrow keys.

**Our own types:**

```go
type BindingInfo struct {
    Accelerator string `json:"accelerator"`
    Description string `json:"description"`
}
```

**IPC messages:**

```go
// Queries
type QueryList struct{}  // → []BindingInfo

// Tasks
type TaskAdd struct {
    Accelerator string `json:"accelerator"`
    Description string `json:"description"`
}  // → error

type TaskRemove struct {
    Accelerator string `json:"accelerator"`
}  // → error

// Actions
type ActionTriggered struct {
    Accelerator string `json:"accelerator"`
}
```

**Service logic:** The Service maintains a `map[string]BindingInfo` registry. When `TaskAdd` is received, it returns `ErrAlreadyRegistered` if the accelerator exists (callers must `TaskRemove` first to rebind). Otherwise, the Service calls `platform.Add(accelerator, callback)` where the callback broadcasts `ActionTriggered` via `s.Core().ACTION()`. `QueryList` reads from the in-memory registry (not `platform.GetAll()`) — `platform.GetAll()` is for adapter-level reconciliation only. The display orchestrator bridges `ActionTriggered` → `Event{Type: "keybinding.triggered"}` for WS clients.

**Wails adapter:** Wraps `app.KeyBinding.Add()`, `app.KeyBinding.Remove()`, `app.KeyBinding.GetAll()`.

**Config:** None — bindings are registered programmatically.

**WS bridge events:** `keybinding.triggered` (broadcast). TS apps call `keybinding:add`, `keybinding:remove`, `keybinding:list` via WS→IPC.

---

### 2. pkg/browser

**Platform interface:**

```go
type Platform interface {
    OpenURL(url string) error
    OpenFile(path string) error
}
```

**IPC messages (all Tasks):**

```go
type TaskOpenURL struct{ URL string `json:"url"` }   // → error
type TaskOpenFile struct{ Path string `json:"path"` } // → error
```

**Wails adapter:** Wraps `app.Browser.OpenURL()` and `app.Browser.OpenFile()`.

**Config:** None — stateless.

**WS bridge:** No actions. TS apps call `browser:open-url`, `browser:open-file` via WS→IPC Tasks.

---

### 3. Window Service Enhancement — File Drop

**PlatformWindow interface addition:**

```go
// Added to existing PlatformWindow interface in pkg/window/platform.go
OnFileDrop(handler func(paths []string, targetID string))
```

**New Action in pkg/window/messages.go:**

```go
type ActionFilesDropped struct {
    Name     string   `json:"name"`     // window name
    Paths    []string `json:"paths"`
    TargetID string   `json:"targetId,omitempty"`
}
```

**Service wiring:** The existing `trackWindow()` method in `pkg/window/service.go` gains a call to `pw.OnFileDrop()` that broadcasts `ActionFilesDropped`. Adding `OnFileDrop` to `PlatformWindow` is a breaking interface change — `MockWindow` in `mock_platform.go` and the Wails adapter must both gain no-op stubs.

**WS bridge:** Display orchestrator adds a case for `window.ActionFilesDropped` → `Event{Type: "window.filedrop"}`.

**Note:** File drop is opt-in per window via the existing `EnableFileDrop` field in `PlatformWindowOptions`. The HTML target element uses `data-file-drop-target` attribute (Wails v3 convention). HTML5 internal drag-and-drop is purely frontend (JS/CSS) — no Go package needed.

---

## Display Orchestrator Changes

### New HandleIPCEvents Cases

```go
case keybinding.ActionTriggered:
    if s.events != nil {
        s.events.Emit(Event{Type: "keybinding.triggered",
            Data: map[string]any{"accelerator": m.Accelerator}})
    }
case window.ActionFilesDropped:
    if s.events != nil {
        s.events.Emit(Event{Type: "window.filedrop", Window: m.Name,
            Data: map[string]any{"paths": m.Paths, "targetId": m.TargetID}})
    }
```

### New WS→IPC Cases

```go
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
```

## Dependency Direction

```
pkg/display (orchestrator)
├── imports pkg/keybinding (message types only)
├── imports pkg/browser (message types only)
└── ... existing imports ...

pkg/keybinding (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface

pkg/browser (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface
```

No circular dependencies.

## Testing Strategy

```go
func TestKeybinding_AddAndTrigger_Good(t *testing.T) {
    mock := &mockPlatform{}
    c, _ := core.New(
        core.WithService(keybinding.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, err := c.PERFORM(keybinding.TaskAdd{
        Accelerator: "Ctrl+S", Description: "Save",
    })
    require.NoError(t, err)
    assert.True(t, handled)

    // Simulate shortcut trigger via mock
    mock.trigger("Ctrl+S")

    // Verify action was broadcast (captured via registered handler)
}

func TestBrowser_OpenURL_Good(t *testing.T) {
    mock := &mockPlatform{}
    c, _ := core.New(
        core.WithService(browser.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, err := c.PERFORM(browser.TaskOpenURL{URL: "https://example.com"})
    require.NoError(t, err)
    assert.True(t, handled)
    assert.Equal(t, "https://example.com", mock.lastURL)
}
```

## Deferred Work

- **Per-window keybindings**: Currently global only. Per-window scoping requires handler-receives-window pattern.
- **Browser fallback**: Copy URL to clipboard when browser open fails.
- **HTML5 drag-drop helpers**: TS SDK utilities for internal drag-drop — purely frontend.

## Licence

EUPL-1.2
