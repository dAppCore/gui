# CoreGUI Spec C: Platform & Events

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Add dock/badge and application lifecycle event packages as core.Services

## Context

Spec A extracted display features, Spec B added input capabilities. This spec covers platform-specific features (dock icon, badge) and application lifecycle events — the remaining Wails v3 features needed for full IPC coverage.

Cursor management is excluded — CSS `cursor:` property handles it in the webview without Go involvement.

## Architecture

### New Packages

| Package | Platform methods | IPC pattern |
|---------|-----------------|-------------|
| `pkg/dock` | `ShowIcon()`, `HideIcon()`, `SetBadge()`, `RemoveBadge()` | Tasks (mutations), Query visibility |
| `pkg/lifecycle` | `OnApplicationEvent()` | Actions broadcast for each event type |

---

## Package Designs

### 1. pkg/dock

**Platform interface:**

```go
type Platform interface {
    ShowIcon() error
    HideIcon() error
    SetBadge(label string) error
    RemoveBadge() error
    IsVisible() bool
}
```

macOS: dock icon show/hide + badge. Windows: taskbar badge only (show/hide not supported). Linux: not supported — adapter returns nil. Platform adapter returns nil on unsupported OS.

**IPC messages:**

```go
// Queries
type QueryVisible struct{}  // → bool

// Tasks
type TaskShowIcon struct{}                 // → error
type TaskHideIcon struct{}                 // → error
type TaskSetBadge struct{ Label string }   // → error
type TaskRemoveBadge struct{}              // → error

// Actions
type ActionVisibilityChanged struct{ Visible bool }
```

**Badge conventions:**
- Empty string `""`: default system badge indicator
- Numeric `"3"`, `"99"`: unread count
- Text `"New"`, `"Paused"`: brief status labels

**Wails adapter:** Wraps `app.Dock.HideAppIcon()`, `app.Dock.ShowAppIcon()` (macOS only), and badge APIs.

**Service logic:** After a successful `TaskShowIcon`, the Service broadcasts `ActionVisibilityChanged{Visible: true}`. After `TaskHideIcon`, broadcasts `ActionVisibilityChanged{Visible: false}`. `QueryVisible` delegates to `platform.IsVisible()`.

**Config:** None — stateless. No config section required in display orchestrator.

**WS bridge events:** `dock.visibility-changed` (on show/hide). TS apps call `dock:show`, `dock:hide`, `dock:badge`, `dock:badge-remove`, `dock:visible` via WS→IPC.

---

### 2. pkg/lifecycle

**Platform interface:**

```go
type Platform interface {
    // OnApplicationEvent registers a handler for a fire-and-forget event type.
    // Events that carry data (e.g. file path) use dedicated methods.
    OnApplicationEvent(eventType EventType, handler func()) func()  // returns cancel
    OnOpenedWithFile(handler func(path string)) func()              // returns cancel
}
```

Separate method for file-open because it carries data (the file path).

**Our own types:**

```go
type EventType int

const (
    EventApplicationStarted EventType = iota
    EventWillTerminate
    EventDidBecomeActive
    EventDidResignActive
    EventPowerStatusChanged  // Windows: APMPowerStatusChange
    EventSystemSuspend       // Windows: APMSuspend
    EventSystemResume        // Windows: APMResume
)
```

**IPC messages (all Actions — lifecycle events are broadcasts):**

```go
type ActionApplicationStarted struct{}
type ActionOpenedWithFile struct{ Path string }
type ActionWillTerminate struct{}
type ActionDidBecomeActive struct{}
type ActionDidResignActive struct{}
type ActionPowerStatusChanged struct{}
type ActionSystemSuspend struct{}
type ActionSystemResume struct{}
```

**Service logic:** During `OnStartup`, the Service registers a platform callback for each `EventType` and for file-open. Each callback broadcasts the corresponding Action via `s.Core().ACTION()`. `OnShutdown` cancels all registrations.

**Wails adapter:** Wraps `app.Event.OnApplicationEvent()` for each Wails event type:

| Our EventType | Wails event | Platforms |
|---|---|---|
| `EventApplicationStarted` | `events.Common.ApplicationStarted` | all |
| `EventWillTerminate` | `events.Mac.ApplicationWillTerminate` | macOS only |
| `EventDidBecomeActive` | `events.Mac.ApplicationDidBecomeActive` | macOS only |
| `EventDidResignActive` | `events.Mac.ApplicationDidResignActive` | macOS only |
| `EventPowerStatusChanged` | `events.Windows.APMPowerStatusChange` | Windows only |
| `EventSystemSuspend` | `events.Windows.APMSuspend` | Windows only |
| `EventSystemResume` | `events.Windows.APMResume` | Windows only |

Platform-specific events no-op silently on unsupported OS (adapter registers nothing).

**Note:** `ActionApplicationStarted` maps to the Wails `ApplicationStarted` event, which fires when the platform application starts. This is distinct from `core.ActionServiceStartup`, which fires after all core.Services complete `OnStartup`. TS clients should use `app.started` for platform-level readiness and subscribe to services individually for service-level readiness.

**Config:** None — event-driven, no state. No config section required in display orchestrator.

**WS bridge events:**

| Action | WS Event Type |
|--------|--------------|
| `ActionApplicationStarted` | `app.started` |
| `ActionOpenedWithFile` | `app.opened-with-file` |
| `ActionWillTerminate` | `app.will-terminate` |
| `ActionDidBecomeActive` | `app.active` |
| `ActionDidResignActive` | `app.inactive` |
| `ActionPowerStatusChanged` | `system.power-change` |
| `ActionSystemSuspend` | `system.suspend` |
| `ActionSystemResume` | `system.resume` |

---

## Display Orchestrator Changes

### New EventType Constants (in `pkg/display/events.go`)

```go
EventDockVisibility    EventType = "dock.visibility-changed"
EventAppStarted        EventType = "app.started"
EventAppOpenedWithFile EventType = "app.opened-with-file"
EventAppWillTerminate  EventType = "app.will-terminate"
EventAppActive         EventType = "app.active"
EventAppInactive       EventType = "app.inactive"
EventSystemPowerChange EventType = "system.power-change"
EventSystemSuspend     EventType = "system.suspend"
EventSystemResume      EventType = "system.resume"
```

### New HandleIPCEvents Cases

```go
case dock.ActionVisibilityChanged:
    if s.events != nil {
        s.events.Emit(Event{Type: EventDockVisibility,
            Data: map[string]any{"visible": m.Visible}})
    }
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

### New WS→IPC Cases

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

Lifecycle events are outbound-only (Actions → WS). No inbound WS→IPC needed.

## Dependency Direction

```
pkg/display (orchestrator)
├── imports pkg/dock (message types only)
├── imports pkg/lifecycle (message types only)
└── ... existing imports ...

pkg/dock (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface

pkg/lifecycle (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface
```

No circular dependencies.

## Testing Strategy

```go
func TestDock_SetBadge_Good(t *testing.T) {
    mock := &mockPlatform{visible: true}
    c, _ := core.New(
        core.WithService(dock.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, err := c.PERFORM(dock.TaskSetBadge{Label: "3"})
    require.NoError(t, err)
    assert.True(t, handled)
    assert.Equal(t, "3", mock.lastBadge)
}

func TestLifecycle_BecomeActive_Good(t *testing.T) {
    mock := &mockPlatform{}
    c, _ := core.New(
        core.WithService(lifecycle.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    var received bool
    c.RegisterAction(func(_ *core.Core, msg core.Message) error {
        if _, ok := msg.(lifecycle.ActionDidBecomeActive); ok {
            received = true
        }
        return nil
    })

    mock.simulateEvent(lifecycle.EventDidBecomeActive)
    assert.True(t, received)
}
```

## Deferred Work

- **Window close hooks**: Cancellable `OnClose` (return false to prevent). Requires hook vs listener distinction in the window service — follow-up enhancement.
- **Custom badge styling**: Windows-specific font/colour options for `SetCustomBadge`. macOS uses system styling only.
- **Notification centre integration**: Deep integration with macOS/Windows notification centres beyond simple badge.

## Licence

EUPL-1.2
