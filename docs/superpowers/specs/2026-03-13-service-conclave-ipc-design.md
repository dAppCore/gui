# CoreGUI Service Conclave — IPC Integration

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Wire core/go IPC into core/gui's split packages, enabling cross-package communication via ACTION/QUERY/PERFORM

## Context

CoreGUI (`forge.lthn.ai/core/gui`) has just been split from a monolithic `pkg/display/` into four packages: `pkg/window`, `pkg/systray`, `pkg/menu`, and a slimmed `pkg/display` orchestrator. Each sub-package defines a `Platform` interface insulating Wails v3.

Today, the orchestrator calls sub-package methods directly. This design replaces direct calls with core/go's IPC bus, making each sub-package a full `core.Service` that communicates via typed messages. This enables:

- Cross-package communication without import coupling
- Declarative window config in `.core/gui/` (foundation for multi-window apps)
- TS apps talking to individual services through the existing WebSocket bridge
- Independent testability with mock platforms and mock core

## Architecture

### Three-Layer Stack

```
IPC Bus (core/go ACTION/QUERY/PERFORM)
    ↓
Manager (pkg/window.Service, pkg/systray.Service, pkg/menu.Service)
    ↓
Platform Interface (Wails v3 adapters)
```

Each sub-package is a `core.Service` that registers its own IPC handlers during `OnStartup`. The display orchestrator is also a `core.Service` — it owns the Wails `*application.App`, wraps it in Platform adapters, and manages config.

No sub-service ever sees Wails types. The orchestrator creates Platform adapters and passes them to sub-service factories.

### Service Registration

Each sub-package exposes a `Register(platform)` factory that returns `func(*core.Core) (any, error)` — the signature `core.WithService` requires. The factory captures the Platform adapter in a closure:

```go
// pkg/window/register.go
func Register(p Platform) func(*core.Core) (any, error) {
    return func(c *core.Core) (any, error) {
        return &Service{
            ServiceRuntime: core.NewServiceRuntime[Options](c, Options{}),
            platform:       p,
        }, nil
    }
}
```

App-level wiring:

```go
wailsApp := application.New(application.Options{...})

windowPlatform := window.NewWailsPlatform(wailsApp)
trayPlatform := systray.NewWailsPlatform(wailsApp)
menuPlatform := menu.NewWailsPlatform(wailsApp)

core.New(
    core.WithService(display.Register(wailsApp)),
    core.WithService(window.Register(windowPlatform)),
    core.WithService(systray.Register(trayPlatform)),
    core.WithService(menu.Register(menuPlatform)),
    core.WithServiceLock(),
)
```

**Startup order**: Display registers first (owns Wails and config), then window/systray/menu. `WithServiceLock()` freezes the registry. `ServiceStartup` calls `OnStartup` sequentially in registration order.

**Critical constraint**: The display orchestrator's `OnStartup` MUST register its `QueryConfig` handler synchronously before returning. Sub-services depend on this query being available when their own `OnStartup` runs.

```go
// pkg/display/service.go
func (s *Service) OnStartup(ctx context.Context) error {
    s.loadConfig()  // Load .core/gui/config.yaml via go-config
    s.Core().RegisterQuery(s.handleQuery)  // QueryConfig available NOW
    s.Core().RegisterTask(s.handleTask)    // TaskSaveConfig available NOW
    // ... remaining setup
    return nil
}
```

**Shutdown order**: Reverse — sub-services save state before display tears down Wails.

## IPC Message Types

Each sub-service defines its own message types. The pattern: Queries return data, Tasks mutate state and may return results, Actions are fire-and-forget broadcasts.

### Window Messages (`pkg/window/messages.go`)

```go
// Queries (read-only)
type QueryWindowList struct{}                    // → []WindowInfo
type QueryWindowByName struct{ Name string }     // → *WindowInfo (nil if not found)

// Tasks (side-effects)
type TaskOpenWindow struct{ Opts []WindowOption }  // → WindowInfo
type TaskCloseWindow struct{ Name string }         // handler persists state BEFORE emitting ActionWindowClosed
type TaskSetPosition struct{ Name string; X, Y int }
type TaskSetSize struct{ Name string; W, H int }
type TaskMaximise struct{ Name string }
type TaskMinimise struct{ Name string }
type TaskFocus struct{ Name string }

// Actions (broadcasts)
type ActionWindowOpened struct{ Name string }
type ActionWindowClosed struct{ Name string }
type ActionWindowMoved struct{ Name string; X, Y int }
type ActionWindowResized struct{ Name string; W, H int }
type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }
```

### Systray Messages (`pkg/systray/messages.go`)

```go
type TaskSetTrayIcon struct{ Data []byte }
type TaskSetTrayMenu struct{ Items []TrayMenuItem }
type TaskShowPanel struct{}
type TaskHidePanel struct{}
type ActionTrayClicked struct{}
```

### Menu Messages (`pkg/menu/messages.go`)

```go
type TaskSetAppMenu struct{ Items []MenuItem }
type QueryGetAppMenu struct{}  // → []MenuItem
```

### Config Messages (display orchestrator)

```go
type QueryConfig struct{ Key string }                          // → map[string]any
type TaskSaveConfig struct{ Key string; Value map[string]any }  // serialisable to YAML
```

## Config via IPC

The display orchestrator owns `.core/gui/config.yaml` (loaded via go-config). Sub-services QUERY for their config section during `OnStartup`. Config saves route through the orchestrator so there is one writer to disk.

```yaml
# .core/gui/config.yaml
window:
  state_file: window_state.json
  default_width: 1024
  default_height: 768

systray:
  icon: apptray.png
  tooltip: "Core GUI"

menu:
  show_dev_tools: true
```

Sub-service startup pattern:

```go
func (s *Service) OnStartup(ctx context.Context) error {
    // Query config from the display orchestrator (registered before us)
    cfg, _, _ := s.Core().QUERY(display.QueryConfig{Key: "window"})
    if wCfg, ok := cfg.(map[string]any); ok {
        s.applyConfig(wCfg)
    }
    // Register QUERY and TASK handlers manually.
    // ACTION handler (HandleIPCEvents) is auto-registered by WithService —
    // do NOT call RegisterAction here or actions will double-fire.
    s.Core().RegisterQuery(s.handleQuery)
    s.Core().RegisterTask(s.handleTask)
    return nil
}

// HandleIPCEvents is auto-discovered and registered by core.WithService.
// Used for filtering broadcast actions relevant to this service.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
    switch msg.(type) {
    case core.ActionServiceStartup:
        // post-startup work if needed
    }
    return nil
}
```

## WSEventManager Bridge

The WSEventManager (WebSocket pub/sub for TS apps) stays in `pkg/display`. It bridges IPC actions to WebSocket events and vice versa.

**IPC → WebSocket** (Go to TS):

The display orchestrator's `HandleIPCEvents` converts IPC actions to `Event` structs and calls `WSEventManager.Emit()`:

```go
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
    switch m := msg.(type) {
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
    }
    return nil
}
```

**WebSocket → IPC** (TS to Go):

Inbound WS messages include a `RequestID` for response correlation. The orchestrator PERFORMs the task and writes back the result or error:

```go
func (s *Service) handleWSMessage(conn *websocket.Conn, msg WSMessage) {
    var result any
    var handled bool
    var err error

    switch msg.Type {
    case "window:open":
        result, handled, err = s.Core().PERFORM(window.TaskOpenWindow{
            Opts: []window.WindowOption{
                window.WithName(msg.Data["name"].(string)),
                window.WithURL(msg.Data["url"].(string)),
            },
        })
    case "window:close":
        result, handled, err = s.Core().PERFORM(
            window.TaskCloseWindow{Name: msg.Data["name"].(string)})
    }

    s.writeResponse(conn, msg.RequestID, result, handled, err)
}
```

TS apps talk WebSocket, the orchestrator translates to/from IPC. The TS app never knows about core/go's bus.

## Dependency Direction

```
pkg/display (orchestrator, core.Service)
├── imports pkg/window (message types only)
├── imports pkg/systray (message types only)
├── imports pkg/menu (message types only)
├── imports core/go (DI, IPC)
└── imports core/go-config (config loading)

pkg/window (core.Service)
├── imports core/go (DI, IPC)
└── uses Platform interface (Wails adapter injected)

pkg/systray (core.Service)
├── imports core/go (DI, IPC)
└── uses Platform interface (Wails adapter injected)

pkg/menu (core.Service)
├── imports core/go (DI, IPC)
└── uses Platform interface (Wails adapter injected)
```

No circular dependencies. Sub-packages do not import each other or the orchestrator. The orchestrator imports sub-package message types for the WSEventManager bridge.

## Testing Strategy

Each sub-service is independently testable with mock platforms and a real core.Core:

```go
func newTestWindowService(t *testing.T) (*Service, *core.Core) {
    c, err := core.New(
        core.WithService(Register(&mockPlatform{})),
        core.WithServiceLock(),
    )
    require.NoError(t, err)
    c.ServiceStartup(context.Background(), nil)
    svc := core.MustServiceFor[*Service](c, "window")
    return svc, c
}

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
    c, _ := core.New(core.WithServiceLock())
    _, handled, _ := c.PERFORM(TaskOpenWindow{})
    assert.False(t, handled)
}
```

Integration test with the full conclave:

```go
func TestServiceConclave_Good(t *testing.T) {
    c, _ := core.New(
        core.WithService(display.Register(mockWailsApp)),
        core.WithService(window.Register(&mockWindowPlatform{})),
        core.WithService(systray.Register(&mockTrayPlatform{})),
        core.WithService(menu.Register(&mockMenuPlatform{})),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, _ := c.PERFORM(window.TaskOpenWindow{
        Opts: []window.WindowOption{window.WithName("main")},
    })
    assert.True(t, handled)

    val, handled, _ := c.QUERY(display.QueryConfig{Key: "window"})
    assert.True(t, handled)
    assert.NotNil(t, val)
}

func TestServiceConclave_Bad(t *testing.T) {
    // Sub-service starts without display — config QUERY returns handled=false
    c, _ := core.New(
        core.WithService(window.Register(&mockWindowPlatform{})),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, _ := c.QUERY(display.QueryConfig{Key: "window"})
    assert.False(t, handled, "no display service means no config handler")
}
```

No Wails runtime needed for any test.

### Service Name Convention

`core.WithService` auto-derives the service name from the type's package path (last segment, lowercased). Canonical names: `"display"`, `"window"`, `"systray"`, `"menu"`. Any existing `ServiceName()` methods returning the full module path should be removed to avoid lookup mismatches.

## Deferred Work

- **Declarative window config**: `.core/gui/windows.yaml` defining named windows with size/position/URL, restored on launch. Foundation is here (config QUERY pattern), but the declarative layer is a follow-up.
- **Screen insulation**: `GetScreens()`, `GetWorkAreas()` still call `application.Get()` directly. Will be wrapped in a `ScreenProvider` interface.
- **go-config dependency**: core/gui currently has no go-config dependency. Adding it is part of implementation.
- **WS response envelope**: Full request/response protocol for WS→IPC (RequestID, error codes, retry semantics). This spec adds the foundation (`writeResponse`), the full envelope schema is a follow-up.

## Licence

EUPL-1.2
