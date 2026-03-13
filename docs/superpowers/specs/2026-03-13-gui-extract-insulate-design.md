# CoreGUI Spec A: Extract & Insulate

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Extract 5 new core.Service packages from pkg/display, each with Platform interface insulation and full IPC coverage for the TS SDK

## Context

The Service Conclave IPC design (`2026-03-13-service-conclave-ipc-design.md`) established the three-layer pattern: IPC Bus → Service → Platform Interface. It converted window, systray, and menu into core.Services. Five feature areas remain embedded in `pkg/display` with direct `application.Get()` calls, making them untestable and invisible to the WS bridge.

This spec extracts those features into independent packages following the same pattern. Full IPC coverage ensures the TypeScript SDK (via WSEventManager) can access every feature, enabling PWA-only frontend development.

## Architecture

### Extraction Targets

| Package | Source file | Platform methods | IPC pattern |
|---------|------------|-----------------|-------------|
| `pkg/clipboard` | `display/clipboard.go` | `Text()`, `SetText()` | Query read, Task write/clear |
| `pkg/dialog` | `display/dialog.go` | `OpenFile()`, `SaveFile()`, `OpenDirectory()`, `MessageDialog()` | Tasks (show UI, return result) |
| `pkg/notification` | `display/notification.go` | `Send()`, `RequestPermission()`, `CheckPermission()` | Tasks + Query auth |
| `pkg/environment` | `display/theme.go` + `display/interfaces.go` (EventSource) | `IsDarkMode()`, `Info()`, `AccentColour()`, `OpenFileManager()`, `OnThemeChange()` | Query + Action theme change |
| `pkg/screen` | `display/display.go` (GetScreens/GetWorkAreas) | `GetAll()`, `GetPrimary()` | Query screen info |

### Three-Layer Stack (per package)

```
IPC Bus (core/go QUERY/TASK/ACTION)
    ↓
Service (core.ServiceRuntime[Options] + business logic)
    ↓
Platform Interface (Wails v3 adapter, injected via Register closure)
```

### Registration Pattern

Each package exposes `Register(platform Platform) func(*core.Core) (any, error)`:

```go
core.New(
    core.WithService(display.Register(wailsApp)),
    core.WithService(window.Register(windowPlatform)),
    core.WithService(systray.Register(trayPlatform)),
    core.WithService(menu.Register(menuPlatform)),
    core.WithService(clipboard.Register(clipPlatform)),
    core.WithService(dialog.Register(dialogPlatform)),
    core.WithService(notification.Register(notifyPlatform)),
    core.WithService(environment.Register(envPlatform)),
    core.WithService(screen.Register(screenPlatform)),
    core.WithServiceLock(),
)
```

Display registers first (owns config), then all sub-services. Order among sub-services does not matter — they only depend on display's config query.

### WS Bridge Integration

The display orchestrator's `HandleIPCEvents` gains cases for all new Action types. The existing `WSEventManager.Emit()` pattern is unchanged — display translates IPC Actions to `Event` structs for WebSocket clients.

---

## Package Designs

### 1. pkg/clipboard

**Platform interface:**

```go
type Platform interface {
    Text() (string, bool)
    SetText(text string) bool
}
```

**IPC messages:**

```go
type QueryText struct{}                    // → ClipboardContent{Text string, HasContent bool}
type TaskSetText struct{ Text string }     // → bool
type TaskClear struct{}                    // → bool
```

**Wails adapter:** Wraps `app.Clipboard.Text()` and `app.Clipboard.SetText()`.

**Config:** None — stateless.

**Implementation note:** `TaskClear` is implemented by the Service calling `platform.SetText("")` — no separate Platform method needed.

**WS bridge:** No actions. TS apps use `clipboard:text` (query) and `clipboard:set-text` / `clipboard:clear` (tasks) via WS→IPC.

---

### 2. pkg/dialog

**Platform interface:**

```go
type Platform interface {
    OpenFile(opts OpenFileOptions) ([]string, error)
    SaveFile(opts SaveFileOptions) (string, error)
    OpenDirectory(opts OpenDirectoryOptions) (string, error)
    MessageDialog(opts MessageDialogOptions) (string, error)
}
```

**Our own types (no Wails leakage):**

```go
type OpenFileOptions struct {
    Title         string
    Directory     string
    Filename      string
    Filters       []FileFilter
    AllowMultiple bool
}

type SaveFileOptions struct {
    Title     string
    Directory string
    Filename  string
    Filters   []FileFilter
}

type OpenDirectoryOptions struct {
    Title         string
    Directory     string
    AllowMultiple bool
}

type MessageDialogOptions struct {
    Type    DialogType
    Title   string
    Message string
    Buttons []string
}

type DialogType int

const (
    DialogInfo DialogType = iota
    DialogWarning
    DialogError
    DialogQuestion
)

type FileFilter struct {
    DisplayName string
    Pattern     string
    Extensions  []string // forward-compat — unused by Wails v3 currently
}
```

**IPC messages (all Tasks — dialogs are side-effects):**

```go
type TaskOpenFile struct{ Opts OpenFileOptions }           // → []string
type TaskSaveFile struct{ Opts SaveFileOptions }           // → string
type TaskOpenDirectory struct{ Opts OpenDirectoryOptions }  // → string (or []string if AllowMultiple)
type TaskMessageDialog struct{ Opts MessageDialogOptions }  // → string (button clicked)
```

**Wails adapter:** Translates our types to `application.OpenFileDialogStruct`, `application.SaveFileDialogStruct`, `application.MessageDialog`.

**Config:** None — stateless.

**WS bridge:** No actions. TS apps call `dialog:open-file`, `dialog:save-file`, `dialog:open-directory`, `dialog:message` via WS→IPC Tasks. Response includes selected path(s) or button clicked.

---

### 3. pkg/notification

**Platform interface:**

```go
type Platform interface {
    Send(opts NotificationOptions) error
    RequestPermission() (bool, error)
    CheckPermission() (bool, error)
}
```

**Our own types:**

```go
type NotificationSeverity int

const (
    SeverityInfo NotificationSeverity = iota
    SeverityWarning
    SeverityError
)

type NotificationOptions struct {
    ID       string
    Title    string
    Message  string
    Subtitle string
    Severity NotificationSeverity // used by fallback to select DialogType
}

type PermissionStatus struct {
    Granted bool
}
```

**IPC messages:**

```go
// Queries
type QueryPermission struct{}                     // → PermissionStatus

// Tasks
type TaskSend struct{ Opts NotificationOptions }  // → error only
type TaskRequestPermission struct{}                // → bool (granted)

// Actions
type ActionNotificationClicked struct{ ID string } // future — when Wails supports callbacks
```

**Fallback logic (Service layer, not Platform):** If `platform.Send()` returns an error, the Service PERFORMs `dialog.TaskMessageDialog` via IPC as a fallback, mapping `Severity` to the appropriate `DialogType` (Info→DialogInfo, Warning→DialogWarning, Error→DialogError). This keeps the Platform interface clean (just native notifications) while the Service owns the fallback decision.

**Migration note:** The existing `ShowInfoNotification`, `ShowWarningNotification`, and `ShowErrorNotification` convenience methods are replaced by `TaskSend` with the appropriate `Severity` field. The existing `ConfirmDialog(title, msg) bool` in display is replaced by `dialog.TaskMessageDialog` with `Type: DialogQuestion, Buttons: ["Yes", "No"]` — callers migrate from checking `bool` to checking `result == "Yes"`.

**Wails adapter:** Wraps `*notifications.NotificationService`.

**Config:** None currently.

**WS bridge:** `notification.clicked` (future). TS apps call `notification:send`, `notification:request-permission` via WS→IPC Tasks.

---

### 4. pkg/environment

**Platform interface:**

```go
type Platform interface {
    IsDarkMode() bool
    Info() EnvironmentInfo
    AccentColour() string
    OpenFileManager(path string, selectFile bool) error
    OnThemeChange(handler func(isDark bool)) func()
}
```

**Our own types:**

```go
type EnvironmentInfo struct {
    OS       string
    Arch     string
    Debug    bool
    Platform PlatformInfo
}

type PlatformInfo struct {
    Name    string
    Version string
}

type ThemeInfo struct {
    IsDark bool
    Theme  string // "dark" or "light"
}
```

**IPC messages:**

```go
// Queries
type QueryTheme struct{}                                       // → ThemeInfo
type QueryInfo struct{}                                        // → EnvironmentInfo
type QueryAccentColour struct{}                                 // → string

// Tasks
type TaskOpenFileManager struct{ Path string; Select bool }    // → error

// Actions
type ActionThemeChanged struct{ IsDark bool }
```

**Theme change flow:** The Service's `OnStartup` registers `platform.OnThemeChange()` callback. When fired, the Service broadcasts `ActionThemeChanged` via `s.Core().ACTION()`. The display orchestrator's `HandleIPCEvents` converts this to `Event{Type: EventThemeChange}` for WebSocket clients.

**What it replaces:** `display.ThemeInfo`, `display.GetTheme()`, `display.GetSystemTheme()`, the `EventSource` interface in `display/types.go`, and the `wailsEventSource` adapter in `display/interfaces.go`.

**Wails adapter:** Wraps `app.Env.IsDarkMode()`, `app.Env.Info()`, `app.Env.GetAccentColor()`, `app.Env.OpenFileManager()`, and `app.Event.OnApplicationEvent(ThemeChanged, ...)`.

**Config:** None — read-only system state.

**WS bridge:** `theme.change` (existing EventThemeChange). TS apps also call `environment:theme`, `environment:info`, `environment:accent-colour`, `environment:open-file-manager` via WS→IPC.

---

### 5. pkg/screen

**Platform interface:**

```go
type Platform interface {
    GetAll() []Screen
    GetPrimary() *Screen
}
```

Two methods. Computed queries (`GetAtPoint`, `GetWorkAreas`) are handled by the Service from `GetAll()` results.

**Our own types:**

```go
type Screen struct {
    ID          string
    Name        string
    ScaleFactor float64
    Size        Size
    Bounds      Rect   // position + dimensions (Bounds.X, Bounds.Y for origin)
    WorkArea    Rect
    IsPrimary   bool
    Rotation    float64
}

type Rect struct {
    X, Y, Width, Height int
}

type Size struct {
    Width, Height int
}
```

**IPC messages:**

```go
// Queries
type QueryAll struct{}                                 // → []Screen
type QueryPrimary struct{}                             // → *Screen
type QueryByID struct{ ID string }                     // → *Screen (nil if not found)
type QueryAtPoint struct{ X, Y int }                   // → *Screen (nil if none)
type QueryWorkAreas struct{}                           // → []Rect

// Actions
type ActionScreensChanged struct{ Screens []Screen }   // future — when Wails supports display hotplug
```

**Service-level computed queries:** `QueryAtPoint` iterates `platform.GetAll()` and returns the screen whose Bounds contain the point. `QueryWorkAreas` extracts WorkArea from each screen.

**What moves from display:** `GetScreens()`, `GetWorkAreas()`, `GetScreen(id)`, `GetScreenAtPoint()` — all currently in `display.go` calling `application.Get().GetScreens()` directly.

**Wails adapter:** Wraps `app.GetScreens()` and finds primary from the list.

**Config:** None — read-only hardware state.

**WS bridge:** `screen.change` (existing EventScreenChange, currently unwired). TS apps call `screen:all`, `screen:primary`, `screen:at-point`, `screen:work-areas` via WS→IPC Queries.

---

## Display Orchestrator Changes

After extraction, `pkg/display` sheds ~420 lines and retains:

- **WSEventManager**: Bridges all sub-service Actions to WebSocket events (gains cases for new Action types)
- **Config ownership**: go-config, `.core/gui/config.yaml`, `QueryConfig`/`TaskSaveConfig` handlers
- **App-level wiring**: Creates Wails Platform adapters, passes to each sub-service's `Register()`
- **WS→IPC translation**: Inbound WebSocket messages translated to IPC Tasks/Queries

The orchestrator no longer contains any clipboard/dialog/notification/theme/screen logic — only routing.

### Removed Types

The following are removed from `pkg/display` after extraction:

- `EventSource` interface (`types.go`) — replaced by `environment.Platform.OnThemeChange()`
- `wailsEventSource` adapter (`interfaces.go`) — moves into environment's Wails adapter
- `newWailsEventSource()` factory (`interfaces.go`) — no longer needed
- `WSEventManager.SetupWindowEventListeners()` (`events.go`) — theme events now arrive via IPC `ActionThemeChanged`
- `NewWSEventManager(es EventSource)` signature changes to `NewWSEventManager()` — no longer takes `EventSource`
- `ClipboardContentType` enum (`clipboard.go`) — Wails v3 only supports text
- `DialogManager` interface subset in `interfaces.go` — replaced by `dialog.Platform`
- `wailsDialogManager` adapter — moves into dialog's Wails adapter

### New HandleIPCEvents Cases

```go
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
    switch m := msg.(type) {
    // ... existing window/systray/menu cases ...
    case environment.ActionThemeChanged:
        s.events.Emit(Event{Type: EventThemeChange,
            Data: map[string]any{"isDark": m.IsDark, "theme": themeStr(m.IsDark)}})
    case notification.ActionNotificationClicked:
        s.events.Emit(Event{Type: "notification.clicked",
            Data: map[string]any{"id": m.ID}})
    case screen.ActionScreensChanged:
        s.events.Emit(Event{Type: EventScreenChange,
            Data: map[string]any{"screens": m.Screens}})
    }
    return nil
}
```

### New WS→IPC Cases

```go
func (s *Service) handleWSMessage(conn *websocket.Conn, msg WSMessage) {
    switch msg.Type {
    // ... existing window cases ...
    case "clipboard:text":
        result, handled, err = s.Core().QUERY(clipboard.QueryText{})
    case "clipboard:set-text":
        result, handled, err = s.Core().PERFORM(clipboard.TaskSetText{Text: msg.Data["text"].(string)})
    case "dialog:open-file":
        result, handled, err = s.Core().PERFORM(dialog.TaskOpenFile{Opts: decodeOpenFileOpts(msg.Data)})
    case "notification:send":
        result, handled, err = s.Core().PERFORM(notification.TaskSend{Opts: decodeNotifyOpts(msg.Data)})
    case "environment:theme":
        result, handled, err = s.Core().QUERY(environment.QueryTheme{})
    case "screen:all":
        result, handled, err = s.Core().QUERY(screen.QueryAll{})
    // ... etc.
    }
}
```

## Dependency Direction

```
pkg/display (orchestrator)
├── imports pkg/clipboard (message types only)
├── imports pkg/dialog (message types only)
├── imports pkg/notification (message types only)
├── imports pkg/environment (message types only)
├── imports pkg/screen (message types only)
├── imports pkg/window (message types only)
├── imports pkg/systray (message types only)
├── imports pkg/menu (message types only)
├── imports core/go (DI, IPC)
└── imports core/go-config

pkg/clipboard, pkg/dialog, pkg/screen (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface (Wails adapter injected)

pkg/notification (has fallback dependency)
├── imports core/go (DI, IPC)
├── imports pkg/dialog (message types only — for fallback)
└── uses Platform interface (Wails adapter injected)

pkg/environment (independent)
├── imports core/go (DI, IPC)
└── uses Platform interface (Wails adapter injected)
```

No circular dependencies. Sub-packages do not import the orchestrator. Only notification imports dialog (message types only, for the fallback PERFORM).

## Testing Strategy

Each package is independently testable with mock platforms and a real core.Core:

```go
func TestClipboard_QueryText_Good(t *testing.T) {
    mock := &mockPlatform{text: "hello", ok: true}
    c, _ := core.New(
        core.WithService(clipboard.Register(mock)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    result, handled, err := c.QUERY(clipboard.QueryText{})
    require.NoError(t, err)
    assert.True(t, handled)
    content := result.(clipboard.ClipboardContent)
    assert.Equal(t, "hello", content.Text)
    assert.True(t, content.HasContent)
}

func TestNotification_Fallback_Good(t *testing.T) {
    // notification platform fails → falls back to dialog via IPC
    mockNotify := &mockNotifyPlatform{sendErr: errors.New("no permission")}
    mockDialog := &mockDialogPlatform{}
    c, _ := core.New(
        core.WithService(dialog.Register(mockDialog)),
        core.WithService(notification.Register(mockNotify)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    _, handled, err := c.PERFORM(notification.TaskSend{
        Opts: notification.NotificationOptions{Title: "Test", Message: "Hello"},
    })
    assert.True(t, handled)
    assert.NoError(t, err)
    assert.True(t, mockDialog.messageDialogCalled) // fallback fired
}
```

Integration test with full conclave follows the existing `TestServiceConclave_Good` pattern.

## Deferred Work

- **Clipboard content types**: Image/HTML clipboard (Wails v3 only supports text currently)
- **Notification callbacks**: `ActionNotificationClicked` — requires Wails v3 notification callback support
- **Screen hotplug**: `ActionScreensChanged` — requires Wails v3 display change events
- **Dialog sync returns**: Current Wails v3 dialog API is async; sync wrapper is a follow-up
- **File drop integration**: `EnableFileDrop` on windows — separate from dialog, deferred to Spec B
- **Prompt dialog**: Text-input dialog (`PromptDialog`) — not supported natively by Wails v3, existing stub dropped

## Licence

EUPL-1.2
