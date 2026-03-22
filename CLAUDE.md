# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `forge.lthn.ai/core/gui` — a display/windowing module for the Core web3 desktop framework. It provides window management, dialogs, system tray, clipboard, notifications, theming, layouts, and real-time WebSocket events. Built on **Wails v3** (Go backend) with an **Angular 20** custom element frontend.

## Build & Development Commands

### Go backend
```bash
go build ./...              # Build all packages
go test ./...               # Run all tests
go test ./pkg/display/...   # Run display package tests
go test -race ./...         # Run tests with race detection
go test -cover ./...        # Run tests with coverage
go test -run TestNew ./pkg/display/  # Run a single test
```

### Angular frontend (pkg/display/ui/)
```bash
cd pkg/display/ui
npm install                 # Install dependencies
npm run build               # Production build
npm run watch               # Dev watch mode
npm run start               # Dev server (localhost:4200)
npm test                    # Unit tests (Karma/Jasmine)
```

## Architecture

### Service-based design with Core DI

The display `Service` registers with `forge.lthn.ai/core/go`'s service container via `Register(c *core.Core)`. It embeds `core.ServiceRuntime[Options]` for lifecycle management and access to sibling services.

### Interface abstraction for testability

All Wails application APIs are abstracted behind interfaces in `interfaces.go` (`App`, `WindowManager`, `MenuManager`, `DialogManager`, etc.). The `wailsApp` adapter wraps the real Wails app. Tests inject a `mockApp` instead — see `mocks_test.go` and the `newServiceWithMockApp(t)` helper.

### Key files in pkg/display/

| File | Responsibility |
|------|---------------|
| `display.go` | Service struct, lifecycle (`Startup`), window CRUD, screen queries, tiling/snapping/layout, workflow presets |
| `window.go` | `WindowOption` functional options pattern, `Window` type alias for `application.WebviewWindowOptions` |
| `window_state.go` | `WindowStateManager` — persists window position/size across restarts |
| `layout.go` | `LayoutManager` — save/restore named window arrangements |
| `events.go` | `WSEventManager` — WebSocket pub/sub for window/theme/screen events |
| `interfaces.go` | Abstract interfaces + Wails adapter implementations |
| `actions.go` | `ActionOpenWindow` IPC message type |
| `menu.go` | Application menu construction |
| `tray.go` | System tray setup |
| `dialog.go` | File/directory dialogs |
| `clipboard.go` | Clipboard read/write |
| `notification.go` | System notifications |
| `theme.go` | Dark/light mode detection |
| `mocks_test.go` | Mock implementations of all interfaces for testing |

### Patterns used throughout

- **Functional options**: `WindowOption` functions (`WindowName()`, `WindowTitle()`, `WindowWidth()`, etc.) configure `application.WebviewWindowOptions`
- **Type alias**: `Window = application.WebviewWindowOptions` — direct alias, not a wrapper
- **Event broadcasting**: `WSEventManager` uses gorilla/websocket with a buffered channel (`eventBuffer`) and per-client subscription filtering (supports `"*"` wildcard)
- **Window lookup by name**: Most Service methods iterate `s.app.Window().GetAll()` and type-assert to `*application.WebviewWindow`, then match by `Name()`

## Testing

- Framework: `testify` (assert + require)
- Pattern: `newServiceWithMockApp(t)` creates a `Service` with mock Wails app — no real window system needed
- `newTestCore(t)` creates a real `core.Core` instance for integration-style tests
- Some tests use `defer func() { recover() }()` to handle nil panics from mock methods that return nil pointers (e.g., `Dialog().Info()`)

## CI/CD

Forgejo Actions (`.forgejo/workflows/`):
- **test.yml**: Runs `go test` with race detection and coverage on push to main/dev and PRs to main
- **security-scan.yml**: Security scanning on push to main/dev/feat/* and PRs to main

Both use reusable workflows from `core/go-devops`.

## Dependencies

- `forge.lthn.ai/core/go` — Core framework with service container and DI
- `github.com/wailsapp/wails/v3` — Desktop app framework (alpha.74)
- `github.com/gorilla/websocket` — WebSocket for real-time events
- `github.com/stretchr/testify` — Test assertions

## Repository migration note

Import paths were recently migrated from `github.com/Snider/Core` to `forge.lthn.ai/core/*`. The `cmd/` directories visible in git status are deleted artifacts from this migration and prior app scaffolds.
