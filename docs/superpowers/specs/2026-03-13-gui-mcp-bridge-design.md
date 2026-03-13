# CoreGUI Spec E: MCP Bridge & WebView Service

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Add `pkg/webview` service wrapping `go-webview` CDP client, add `pkg/mcp` MCP display subsystem exposing all GUI packages as MCP tools, update display orchestrator for webview event bridging

## Context

Specs A–D extracted 14 packages covering all 11 Wails v3 Manager APIs through IPC. The old `core-gui` repo (archived) had an MCP bridge (`mcp_bridge.go`, 1136 lines) exposing ~90 tools via HTTP REST. Two gaps remain:

1. **WebView interaction** — JS eval, console capture, DOM queries, screenshots, network inspection — has no core.Service wrapper
2. **MCP tool layer** — no way for Claude Code or other MCP clients to call the GUI IPC packages

This spec closes both gaps by:
- Wrapping `go-webview` (CDP client) as `pkg/webview` — a core.Service with IPC messages
- Creating `pkg/mcp` — an MCP subsystem that translates tool calls to IPC messages across all 15 packages

The old bridge's non-GUI tools (file ops, process mgmt, lang detect) already live in `core/mcp` via `go-io`, `go-process`, and `go-i18n` — they are not ported here.

## Architecture

### Dependency Direction

```
core/gui/pkg/mcp (MCP subsystem)
├── imports core/go (Core, PERFORM, QUERY)
├── imports core/gui/pkg/* (IPC message types only)
├── imports github.com/modelcontextprotocol/go-sdk/mcp (MCP SDK)
└── structurally satisfies core/mcp Subsystem interface (zero import)

core/gui/pkg/webview (new service)
├── imports core/go (ServiceRuntime, IPC)
├── imports forge.lthn.ai/core/go-webview (CDP client)
└── no Platform interface (go-webview IS the abstraction)
```

No circular dependencies. `core/mcp` never imports `core/gui` — structural typing satisfies the `Subsystem` interface. The consumer wires the subsystem at the application level.

### Integration Point

The consumer (e.g. `cmd/core-gui`) creates the subsystem and passes it to `core/mcp`:

```go
guiSub := guimcp.New(coreInstance)
mcpSvc, _ := coremcp.New(coremcp.WithSubsystem(guiSub))
```

---

## Package Design

### 1. pkg/webview

Wraps `go-webview` (Chrome DevTools Protocol client) as a core.Service. No Platform interface — `go-webview` already abstracts Chrome/Chromium via CDP WebSocket. The service manages a map of window name → `*webview.Webview` connections, lazily created on first access.

**Configuration:**

```go
type Options struct {
    DebugURL string // Chrome debug endpoint (default: "http://localhost:9222")
    Timeout  time.Duration // Operation timeout (default: 30s)
    ConsoleLimit int // Max console messages per window (default: 1000)
}
```

**Register factory:**

```go
func Register(opts ...func(*Options)) func(*core.Core) (any, error)
```

No Platform interface — the service creates `go-webview.Webview` instances directly using the configured debug URL. Each window gets its own CDP connection, managed in `connections map[string]*webview.Webview`.

**IPC messages:**

```go
// Queries (read-only)
type QueryURL          struct{ Window string `json:"window"` }
type QueryTitle        struct{ Window string `json:"window"` }
type QueryConsole      struct{ Window string `json:"window"`; Level string `json:"level,omitempty"`; Limit int `json:"limit,omitempty"` }
type QuerySelector     struct{ Window string `json:"window"`; Selector string `json:"selector"` }
type QuerySelectorAll  struct{ Window string `json:"window"`; Selector string `json:"selector"` }
type QueryDOMTree      struct{ Window string `json:"window"`; Selector string `json:"selector,omitempty"` }

// Tasks (side-effects)
type TaskEvaluate      struct{ Window string `json:"window"`; Script string `json:"script"` }
type TaskClick         struct{ Window string `json:"window"`; Selector string `json:"selector"` }
type TaskType          struct{ Window string `json:"window"`; Selector string `json:"selector"`; Text string `json:"text"` }
type TaskNavigate      struct{ Window string `json:"window"`; URL string `json:"url"` }
type TaskScreenshot    struct{ Window string `json:"window"` } // → ScreenshotResult{Base64 string, MimeType string}
type TaskScroll        struct{ Window string `json:"window"`; X int `json:"x"`; Y int `json:"y"` } // X, Y are absolute scroll position (window.scrollTo)
type TaskHover         struct{ Window string `json:"window"`; Selector string `json:"selector"` }
type TaskSelect        struct{ Window string `json:"window"`; Selector string `json:"selector"`; Value string `json:"value"` }
type TaskCheck         struct{ Window string `json:"window"`; Selector string `json:"selector"`; Checked bool `json:"checked"` }
type TaskUploadFile    struct{ Window string `json:"window"`; Selector string `json:"selector"`; Paths []string `json:"paths"` }
type TaskSetViewport   struct{ Window string `json:"window"`; Width int `json:"width"`; Height int `json:"height"` }
type TaskClearConsole  struct{ Window string `json:"window"` }

// Actions (broadcast)
type ActionConsoleMessage struct{ Window string `json:"window"`; Message ConsoleMessage `json:"message"` }
type ActionException      struct{ Window string `json:"window"`; Exception ExceptionInfo `json:"exception"` }
```

**Our own types** (decoupled from go-webview — field names normalised from go-webview's `LineNumber`→`Line`, `ColumnNumber`→`Column`):

```go
type ConsoleMessage struct {
    Type      string    `json:"type"`      // "log", "warn", "error", "info", "debug"
    Text      string    `json:"text"`
    Timestamp time.Time `json:"timestamp"`
    URL       string    `json:"url,omitempty"`
    Line      int       `json:"line,omitempty"`      // go-webview: Line (same)
    Column    int       `json:"column,omitempty"`    // go-webview: Column (same)
}

type ElementInfo struct {
    TagName    string            `json:"tagName"`
    Attributes map[string]string `json:"attributes,omitempty"`
    InnerText  string            `json:"innerText,omitempty"`
    InnerHTML  string            `json:"innerHTML,omitempty"`
    BoundingBox *BoundingBox     `json:"boundingBox,omitempty"`
}

type BoundingBox struct {
    X      float64 `json:"x"`
    Y      float64 `json:"y"`
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
}

type ExceptionInfo struct {
    Text       string    `json:"text"`
    URL        string    `json:"url,omitempty"`
    Line       int       `json:"line,omitempty"`      // go-webview: LineNumber
    Column     int       `json:"column,omitempty"`    // go-webview: ColumnNumber
    StackTrace string    `json:"stackTrace,omitempty"`
    Timestamp  time.Time `json:"timestamp"`
}

// ScreenshotResult wraps the raw PNG bytes as base64 for JSON/MCP transport.
type ScreenshotResult struct {
    Base64   string `json:"base64"`   // base64-encoded PNG
    MimeType string `json:"mimeType"` // always "image/png"
}
```

**Type conversion:** The service converts between go-webview types and our own types. Key mappings:
- `gowebview.ConsoleMessage` → `ConsoleMessage` (field names match)
- `gowebview.ExceptionInfo` → `ExceptionInfo` (`LineNumber`→`Line`, `ColumnNumber`→`Column`)
- `gowebview.ElementInfo` → `ElementInfo` (field names match)
- `Screenshot() ([]byte, error)` → `ScreenshotResult{Base64: base64.StdEncoding.EncodeToString(png), MimeType: "image/png"}`

**Service logic:** The service registers IPC handlers for all Query/Task types. On each call, it looks up (or lazily creates) the `*webview.Webview` connection for the named window, delegates to go-webview, and converts the result to our own types. Console watchers broadcast `ActionConsoleMessage` via `s.Core().ACTION()`. Exception watchers broadcast `ActionException`.

**Window-name → CDP-target mapping:** Each Wails WebviewWindow runs as a separate CDP target. On first access for a window name, the service calls `go-webview.ListTargets(debugURL)` to enumerate all browser targets, matches by page title or URL containing the window name, and creates a `*webview.Webview` connected to that target's WebSocket URL. If no match is found, falls back to the first available page target (single-window mode).

**Prerequisite:** `go-webview.ListTargets()` currently returns an unexported `targetInfo` type. Before implementation, export this as `TargetInfo` in go-webview (one-line rename, non-breaking since the function is newly introduced).

**Connection lifecycle:** Connections are created on first access via the target mapping above and cached in `connections map[string]*webview.Webview`. When a window is closed (detected via `window.ActionWindowClosed` IPC event), the corresponding webview connection is closed and removed from the map. The service implements `Stoppable` to close all connections on shutdown.

**Config:** None persisted — runtime options only.

---

### 2. pkg/mcp

MCP display subsystem that registers GUI tools with an MCP server. Implements `Name() string` and `RegisterTools(server *mcp.Server)` — structurally satisfying `core/mcp`'s `Subsystem` interface without importing it.

**Constructor:**

```go
func New(c *core.Core) *Subsystem
```

**Subsystem struct:**

```go
type Subsystem struct {
    core *core.Core
}

func (s *Subsystem) Name() string { return "display" }

func (s *Subsystem) RegisterTools(server *mcp.Server) {
    s.registerWebviewTools(server)
    s.registerWindowTools(server)
    s.registerLayoutTools(server)
    s.registerScreenTools(server)
    s.registerClipboardTools(server)
    s.registerDialogTools(server)
    s.registerNotificationTools(server)
    s.registerTrayTools(server)
    s.registerEnvironmentTools(server)
    s.registerBrowserTools(server)
    s.registerContextMenuTools(server)
    s.registerKeybindingTools(server)
    s.registerDockTools(server)
    s.registerLifecycleTools(server)
}
```

**Tool files and counts:**

| File | Tool Count | Tools | IPC Target |
|------|-----------|-------|-----------|
| `tools_webview.go` | 18 | `webview_eval`, `webview_click`, `webview_type`, `webview_navigate`, `webview_screenshot` (returns base64 PNG), `webview_scroll`, `webview_hover`, `webview_select`, `webview_check`, `webview_upload`, `webview_viewport`, `webview_console`, `webview_console_clear`, `webview_query`, `webview_query_all`, `webview_dom_tree`, `webview_url`, `webview_title` | `pkg/webview` |
| `tools_window.go` | 15 | `window_list`, `window_get`, `window_focused`, `window_create`, `window_close`, `window_position`, `window_size`, `window_bounds`, `window_maximize`, `window_minimize`, `window_restore`, `window_focus`, `window_title`, `window_visibility`, `window_fullscreen` | `pkg/window` |
| `tools_layout.go` | 7 | `layout_save`, `layout_restore`, `layout_list`, `layout_delete`, `layout_get`, `layout_tile`, `layout_snap` | `pkg/window` |
| `tools_screen.go` | 5 | `screen_list`, `screen_get`, `screen_primary`, `screen_at_point`, `screen_work_areas` | `pkg/screen` |
| `tools_clipboard.go` | 4 | `clipboard_read`, `clipboard_write`, `clipboard_has`, `clipboard_clear` | `pkg/clipboard` |
| `tools_dialog.go` | 5 | `dialog_open_file`, `dialog_save_file`, `dialog_open_directory`, `dialog_confirm`, `dialog_prompt` | `pkg/dialog` |
| `tools_notification.go` | 3 | `notification_show`, `notification_permission_request`, `notification_permission_check` | `pkg/notification` |
| `tools_tray.go` | 4 | `tray_set_icon`, `tray_set_tooltip`, `tray_set_label`, `tray_info` | `pkg/systray` |
| `tools_environment.go` | 2 | `theme_get`, `theme_system` | `pkg/environment` |
| `tools_browser.go` | 1 | `browser_open_url` | `pkg/browser` |
| `tools_contextmenu.go` | 4 | `contextmenu_add`, `contextmenu_remove`, `contextmenu_get`, `contextmenu_list` | `pkg/contextmenu` |
| `tools_keybinding.go` | 2 | `keybinding_add`, `keybinding_remove` | `pkg/keybinding` |
| `tools_dock.go` | 3 | `dock_show`, `dock_hide`, `dock_badge` | `pkg/dock` |
| `tools_lifecycle.go` | 1 | `app_quit` | `pkg/lifecycle` |

**~74 tools total.**

**Tool handler pattern:** Each handler is a thin IPC translation. Typed Input/Output structs with JSON tags, same pattern as `core/mcp`'s file tool handlers:

```go
type WebviewEvalInput struct {
    Window string `json:"window"`
    Script string `json:"script"`
}

type WebviewEvalOutput struct {
    Result any    `json:"result"`
    Window string `json:"window"`
}

func (s *Subsystem) webviewEval(ctx context.Context, req *mcp.CallToolRequest, input WebviewEvalInput) (*mcp.CallToolResult, WebviewEvalOutput, error) {
    result, _, err := s.core.PERFORM(webview.TaskEvaluate{
        Window: input.Window,
        Script: input.Script,
    })
    if err != nil {
        return nil, WebviewEvalOutput{}, err
    }
    return nil, WebviewEvalOutput{Result: result, Window: input.Window}, nil
}
```

```go
type ClipboardReadInput struct{}

type ClipboardReadOutput struct {
    Content string `json:"content"`
}

func (s *Subsystem) clipboardRead(ctx context.Context, req *mcp.CallToolRequest, input ClipboardReadInput) (*mcp.CallToolResult, ClipboardReadOutput, error) {
    result, _, err := s.core.QUERY(clipboard.QueryText{})
    if err != nil {
        return nil, ClipboardReadOutput{}, err
    }
    content, _ := result.(string)
    return nil, ClipboardReadOutput{Content: content}, nil
}
```

**Tool registration** uses `mcp.AddTool()` from the MCP SDK with typed handler functions:

```go
func (s *Subsystem) registerClipboardTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "clipboard_read",
        Description: "Read the current clipboard content",
    }, s.clipboardRead)

    mcp.AddTool(server, &mcp.Tool{
        Name:        "clipboard_write",
        Description: "Write text to the clipboard",
    }, s.clipboardWrite)

    // ... etc
}
```

---

## Display Orchestrator Changes

### New HandleIPCEvents Cases

```go
case webview.ActionConsoleMessage:
    if s.events != nil {
        s.events.Emit(Event{Type: EventWebviewConsole, Window: m.Window,
            Data: map[string]any{"message": m.Message}})
    }
case webview.ActionException:
    if s.events != nil {
        s.events.Emit(Event{Type: EventWebviewException, Window: m.Window,
            Data: map[string]any{"exception": m.Exception}})
    }
```

### New EventType Constants

```go
EventWebviewConsole   EventType = "webview.console"
EventWebviewException EventType = "webview.exception"
```

### New WS→IPC Cases

```go
case "webview:eval":
    window, _ := msg.Data["window"].(string)
    script, _ := msg.Data["script"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskEvaluate{Window: window, Script: script})
case "webview:click":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskClick{Window: window, Selector: selector})
case "webview:type":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    text, _ := msg.Data["text"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskType{Window: window, Selector: selector, Text: text})
case "webview:navigate":
    window, _ := msg.Data["window"].(string)
    url, _ := msg.Data["url"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskNavigate{Window: window, URL: url})
case "webview:screenshot":
    window, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskScreenshot{Window: window})
case "webview:console":
    window, _ := msg.Data["window"].(string)
    level, _ := msg.Data["level"].(string)
    limit := 100
    if l, ok := msg.Data["limit"].(float64); ok {
        limit = int(l)
    }
    result, handled, err = s.Core().QUERY(webview.QueryConsole{Window: window, Level: level, Limit: limit})
case "webview:query":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QuerySelector{Window: window, Selector: selector})
case "webview:url":
    window, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryURL{Window: window})
case "webview:title":
    window, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryTitle{Window: window})
case "webview:viewport":
    window, _ := msg.Data["window"].(string)
    width, _ := msg.Data["width"].(float64)
    height, _ := msg.Data["height"].(float64)
    result, handled, err = s.Core().PERFORM(webview.TaskSetViewport{Window: window, Width: int(width), Height: int(height)})
case "webview:clear-console":
    window, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskClearConsole{Window: window})
case "webview:query-all":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QuerySelectorAll{Window: window, Selector: selector})
case "webview:dom-tree":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryDOMTree{Window: window, Selector: selector})
case "webview:hover":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskHover{Window: window, Selector: selector})
case "webview:select":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    value, _ := msg.Data["value"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskSelect{Window: window, Selector: selector, Value: value})
case "webview:check":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    checked, _ := msg.Data["checked"].(bool)
    result, handled, err = s.Core().PERFORM(webview.TaskCheck{Window: window, Selector: selector, Checked: checked})
case "webview:upload":
    window, _ := msg.Data["window"].(string)
    selector, _ := msg.Data["selector"].(string)
    pathsRaw, _ := msg.Data["paths"].([]any)
    var paths []string
    for _, p := range pathsRaw {
        if s, ok := p.(string); ok {
            paths = append(paths, s)
        }
    }
    result, handled, err = s.Core().PERFORM(webview.TaskUploadFile{Window: window, Selector: selector, Paths: paths})
case "webview:scroll":
    window, _ := msg.Data["window"].(string)
    x, _ := msg.Data["x"].(float64)
    y, _ := msg.Data["y"].(float64)
    result, handled, err = s.Core().PERFORM(webview.TaskScroll{Window: window, X: int(x), Y: int(y)})
```

---

## Testing Strategy

### pkg/webview

Mock the go-webview client behind a thin interface internal to the package for testing. The interface uses our own types (not go-webview types) to keep tests decoupled:

```go
type connector interface {
    Navigate(url string) error
    Click(selector string) error
    Type(selector, text string) error
    Hover(selector string) error
    Select(selector, value string) error
    Check(selector string, checked bool) error
    Evaluate(script string) (any, error)
    Screenshot() ([]byte, error)
    GetURL() (string, error)
    GetTitle() (string, error)
    GetHTML(selector string) (string, error)  // used internally by QueryDOMTree
    QuerySelector(selector string) (*ElementInfo, error)
    QuerySelectorAll(selector string) ([]*ElementInfo, error)
    GetConsole() []ConsoleMessage
    ClearConsole()
    SetViewport(width, height int) error
    UploadFile(selector string, paths []string) error
    Close() error
}
```

The real implementation wraps `*gowebview.Webview`, converting go-webview types to our own types at the boundary. Tests inject a mock connector:

```go
func TestWebview_Evaluate_Good(t *testing.T) {
    mock := &mockConnector{evalResult: "hello"}
    svc := newTestService(mock)
    c, _ := core.New(core.WithService(svc), core.WithServiceLock())
    c.ServiceStartup(context.Background(), nil)

    result, handled, err := c.PERFORM(webview.TaskEvaluate{Window: "main", Script: "1+1"})
    require.NoError(t, err)
    assert.True(t, handled)
    assert.Equal(t, "hello", result)
}

func TestWebview_Console_Good(t *testing.T) {
    mock := &mockConnector{consoleMessages: []webview.ConsoleMessage{
        {Type: "log", Text: "hello"},
        {Type: "error", Text: "oops"},
    }}
    svc := newTestService(mock)
    c, _ := core.New(core.WithService(svc), core.WithServiceLock())
    c.ServiceStartup(context.Background(), nil)

    result, _, err := c.QUERY(webview.QueryConsole{Window: "main", Level: "error", Limit: 10})
    require.NoError(t, err)
    msgs, _ := result.([]webview.ConsoleMessage)
    assert.Len(t, msgs, 1)
    assert.Equal(t, "oops", msgs[0].Text)
}

func TestWebview_ConnectionCleanup_Good(t *testing.T) {
    mock := &mockConnector{}
    svc := newTestService(mock)
    c, _ := core.New(core.WithService(svc), core.WithServiceLock())
    c.ServiceStartup(context.Background(), nil)

    // Access creates connection
    _, _, _ = c.QUERY(webview.QueryURL{Window: "main"})
    assert.True(t, mock.connected)

    // Window close action triggers cleanup
    _ = c.ACTION(window.ActionWindowClosed{Name: "main"})
    assert.True(t, mock.closed)
}
```

### pkg/mcp

Test that tool handlers correctly translate to IPC calls. Since handlers are unexported (passed to `mcp.AddTool`), tests verify at the IPC level — register real sub-services with mock platforms, create the subsystem, and verify the IPC round-trip works:

```go
func TestMCP_ClipboardRead_Good(t *testing.T) {
    mockClip := &mockClipboardPlatform{content: "hello"}
    c, _ := core.New(
        core.WithService(clipboard.Register(mockClip)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    // Verify the IPC path that the MCP handler would use
    result, handled, err := c.QUERY(clipboard.QueryText{})
    require.NoError(t, err)
    assert.True(t, handled)
    content, _ := result.(string)
    assert.Equal(t, "hello", content)
}

func TestMCP_WindowList_Good(t *testing.T) {
    mockWin := &mockWindowPlatform{windows: []window.Info{{Name: "main"}}}
    c, _ := core.New(
        core.WithService(window.Register(mockWin)),
        core.WithServiceLock(),
    )
    c.ServiceStartup(context.Background(), nil)

    // Verify the IPC path that the MCP handler would use
    result, handled, err := c.QUERY(window.QueryList{})
    require.NoError(t, err)
    assert.True(t, handled)
    windows, _ := result.([]window.Info)
    assert.Len(t, windows, 1)
}

func TestMCP_RegisterTools_Good(t *testing.T) {
    c, _ := core.New(core.WithServiceLock())
    sub := guimcp.New(c)
    assert.Equal(t, "display", sub.Name())
    // RegisterTools requires a real mcp.Server — verify it does not panic
    assert.NotPanics(t, func() { sub.RegisterTools(mcp.NewServer(nil, nil)) })
}
```

---

## Deferred Work

- **Performance metrics**: `QueryPerformance` — requires direct CDP `Performance.getMetrics` call not yet exposed by go-webview's public API. Add when go-webview gains this.
- **Resource listing**: `QueryResources` — requires CDP `Page.getResourceTree`, not yet in go-webview. Add when go-webview gains this.
- **Network inspection**: `QueryNetwork` / `TaskClearNetwork` — requires CDP `Network` domain enablement and request tracking, not yet in go-webview. The existing `InjectNetworkInterceptor` is for request modification, not read-only inspection.
- **WebView action sequences**: Chaining multiple webview actions (click → wait → type → screenshot) as a single MCP tool call. Currently requires multiple sequential tool calls.
- **Angular helpers**: `go-webview`'s `AngularHelper` provides SPA-specific testing utilities — useful for the Angular frontend but not needed for initial MCP bridge.
- **Multi-tab CDP**: `go-webview`'s `CDPClient.NewTab()` could map to window creation for headless browser automation scenarios.

## Licence

EUPL-1.2
