# CoreGUI Display Package Split

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Restructure `pkg/display/` monolith into 4 focused packages

## Context

CoreGUI (`forge.lthn.ai/core/gui`) is an abstraction layer over Wails v3 — the "display server" that TypeScript applications bind to for window, tray, and menu management. Apps never import Wails directly; CoreGUI provides the stable API contract. If Wails breaks, it's fixed in one place.

This lineage traces back to `/Users/snider/Code/dappserver` — the client-side server for PWAs that formed a contract of APIs/polyfills inside the webview window to talk to the runtime.

Today, all 3,910 LOC live in a single `pkg/display/` package across 15 files. This design splits it into 4 packages with clear boundaries.

## Package Boundaries

### pkg/window — Window lifecycle and spatial management

Extracted from `display.go` (window CRUD, tiling, snapping, layouts), `window.go`, `window_state.go`, `layout.go`.

**Responsibilities:**
- `Window` struct (CoreGUI's own, NOT a type alias — replaces `type Window = application.WebviewWindowOptions`)
- `WindowOption` functional options rewritten against CoreGUI's `Window` struct: `func(*Window) error`
- `WindowStateManager` — JSON persistence to `~/.config/Core/window_state.json`
- `LayoutManager` — named window arrangements to `~/.config/Core/layouts.json`
- Tiling (9 modes), snapping (9 positions), stacking
- Workflow presets (coding, debugging, presenting, side-by-side)
- `Platform` adapter interface insulating Wails

**Source files (from current pkg/display/):**

| Current File | Destination | Notes |
|---|---|---|
| `display.go` (~800 LOC window section) | `window.go`, `tiling.go`, `layout.go` | Split by concern |
| `window.go` (90 LOC) | `options.go` | Functional options, own Window type |
| `window_state.go` (261 LOC) | `state.go` | JSON persistence |
| `layout.go` (149 LOC) | `layout.go` | Named arrangements |

### pkg/systray — System tray and panel

Extracted from `tray.go` plus `TrayMenuItem` types from `display.go`.

**Responsibilities:**
- Tray creation, lifecycle
- Icon management (template for macOS, dual-mode for Windows/Linux)
- Tooltip and label
- Dynamic menu builder (`TrayMenuItem` recursive tree)
- Callback registry (`RegisterTrayMenuCallback`) — stored as `Manager` fields, NOT package-level vars
- Attached panel window (hidden, frameless, offset) — accepts a `WindowHandle` interface (see Shared Types)
- `Platform` adapter interface insulating Wails
- **Migration note:** Current `activeTray` and `trayMenuCallbacks` package-level vars become fields on `Manager`

**Source files:**

| Current File | Destination | Notes |
|---|---|---|
| `tray.go` (200 LOC) | `tray.go` | Core tray lifecycle |
| `display.go` (TrayMenuItem types, SetTrayMenu, etc.) | `menu.go`, `types.go` | Dynamic menu + types |

### pkg/menu — Application menus

Extracted from `menu.go`.

**Responsibilities:**
- Menu builder — constructs menu item trees (labels, accelerators, submenus, separators)
- Menu item types and accelerator bindings
- `Platform` adapter interface insulating Wails
- **Click handlers live in `pkg/display`**, not here. `pkg/menu` builds structure only; the orchestrator injects closures for app-specific actions (open file, new window, etc.)

**Source files:**

| Current File | Destination | Notes |
|---|---|---|
| `menu.go` (185 LOC) | `menu.go` | App menu construction |

### pkg/display — Orchestrator / display server contract

What remains — the glue layer and the API surface TypeScript apps bind to.

**Responsibilities:**
- `Service` struct embedding `core.ServiceRuntime`
- Composes `window.Manager`, `systray.Manager`, `menu.Manager`
- `WSEventManager` — WebSocket pub/sub bridge for TS apps (the display server channel)
- Dialog helpers (file open, save, directory select)
- Clipboard (read/write text, HTML, images)
- Notifications (native system, fallback to dialog)
- Theme detection (dark/light)
- IPC action types
- `Register()` factory for core DI
- Shared types consumed by TS apps: `ScreenInfo`, `WorkArea`

**Source files that stay:**

| File | LOC | Notes |
|---|---|---|
| `events.go` | 365 | WebSocket bridge — the display server contract |
| `dialog.go` | 192 | File/directory dialogs |
| `notification.go` | 127 | System notifications |
| `clipboard.go` | 61 | Clipboard operations |
| `theme.go` | 38 | Theme detection |
| `actions.go` | 20 | IPC message types |
| `display.go` (~500 LOC) | Orchestration, startup, service wiring |

### ui/ — Feature demo (top-level)

Moved from `pkg/display/ui/` to top-level `ui/`.

- Reference implementation demonstrating all CoreGUI capabilities
- Sets the standard pattern for downstream apps (BugSETI, LEM, Mining, IDE)
- Existing Angular code in `pkg/display/ui/` moves as-is to top-level `ui/`; `go:embed` directives update to match
- Placeholder README added explaining its purpose as feature demo
- Future: Playwright inside WebView2 for automated testing, errors surfaced to agents

## Dependency Direction

```
pkg/display (orchestrator)
├── imports pkg/window
├── imports pkg/systray
├── imports pkg/menu
└── imports core/go (DI)

pkg/window ──→ core/go, wails/v3 (behind Platform interface)
pkg/systray ──→ core/go, wails/v3 (behind Platform interface)
pkg/menu ──→ core/go, wails/v3 (behind Platform interface)
```

No circular dependencies. Window, systray, and menu are peers — they do not import each other. Display imports all three and wires them together.

Shared types (`ScreenInfo`, `WorkArea`) live in `pkg/display` since that's the contract layer TS apps consume.

### Shared Types

A `WindowHandle` interface lives in `pkg/display` for cross-package use (e.g. systray attaching a panel window without importing `pkg/window`):

```go
// pkg/display/types.go
type WindowHandle interface {
    Name() string
    Show()
    Hide()
    SetPosition(x, y int)
    SetSize(width, height int)
}
```

Both `pkg/window.PlatformWindow` and `pkg/systray.PlatformTray.AttachWindow()` work with this interface — no `any` types, no cross-package peer imports.

## Wails Insulation Pattern

Each sub-package defines a `Platform` interface — the adapter contract. Wails never leaks past this boundary.

```go
// pkg/window/platform.go
type Platform interface {
    CreateWindow(opts PlatformWindowOptions) PlatformWindow
    GetWindows() []PlatformWindow
}

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

    // Events (for WSEventManager insulation)
    OnWindowEvent(handler func(event WindowEvent))
}
```

```go
// pkg/systray/platform.go
type Platform interface {
    NewTray() PlatformTray
}

type PlatformTray interface {
    SetIcon(data []byte)
    SetTemplateIcon(data []byte)
    SetTooltip(text string)
    SetLabel(text string)
    SetMenu(menu PlatformMenu)
    AttachWindow(w display.WindowHandle)
}
```

```go
// pkg/menu/platform.go
type Platform interface {
    NewMenu() PlatformMenu
}

type PlatformMenu interface {
    Add(label string) PlatformMenuItem
    AddSeparator()
}
```

**Key rules:**
- Interfaces only expose what CoreGUI actually uses — no speculative wrapping
- Wails adapter implementations live in each package (e.g. `pkg/window/wails.go`)
- Mock implementations for testing (e.g. `pkg/window/mock_test.go`)
- If Wails changes (v4, breaking API), update 3 adapter files — nothing else changes

### WSEventManager Insulation

`WSEventManager` (stays in `pkg/display`) currently calls `application.Get()` directly and takes `*application.WebviewWindow` in `AttachWindowListeners`. After the split:

- `AttachWindowListeners` accepts `PlatformWindow` (which has `OnWindowEvent`) instead of the Wails concrete type
- The `application.Get()` call moves into the Wails adapter — the event manager receives an `EventSource` interface
- This allows testing the WebSocket bridge without a Wails runtime

### WindowStateManager Insulation

`CaptureState` currently takes `*application.WebviewWindow`. After the split:

- `CaptureState` accepts `PlatformWindow` interface (has `Name()`, `Position()`, `Size()`, `IsMaximised()`)
- `ApplyState` returns CoreGUI's own `Window` struct, not `application.WebviewWindowOptions`

## Testing Strategy

Each package gets its own test suite with mock platform:

| Package | Test File | Mock | Coverage |
|---|---|---|---|
| `pkg/window` | `window_test.go` | `mockPlatform` + `mockWindow` | CRUD, state persistence, tiling, snapping, layouts, presets |
| `pkg/systray` | `tray_test.go` | `mockPlatform` + `mockTray` | Icon, menu, callbacks, panel attachment |
| `pkg/menu` | `menu_test.go` | `mockPlatform` + `mockMenu` | Construction, item types, accelerators |
| `pkg/display` | `display_test.go` | Composes sub-package mocks | Orchestration, events, dialogs, clipboard, notifications |

Existing tests from `display_test.go` (636 LOC) split to follow their code. Each package gets its own `newTestX(t)` helper for creating a service with mock platform.

Test framework: `testify` (assert + require). Naming convention: `_Good`/`_Bad`/`_Ugly` suffix pattern from core/go.

## Reference Patterns

These existing implementations inform the systray and window patterns:

- **LEM** (`/Users/snider/Code/lthn/LEM/cmd/lem-desktop/tray.go`): Multi-window TrayService with dashboard snapshots, platform-specific icons
- **Mining** (`/Users/snider/Code/snider/Mining/`): Angular custom elements (`createCustomElement`), Wails service facade, sparkline SVG charts
- **core/ide** (`/Users/snider/Code/core/ide/main.go`): Simpler systray with tray panel window

## Licence

EUPL-1.2
