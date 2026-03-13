# MCP Bridge & WebView Service Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `pkg/webview` (CDP-backed core.Service) and `pkg/mcp` (MCP display subsystem with ~74 tools) to core/gui

**Architecture:** `pkg/webview` wraps go-webview as a core.Service with IPC messages. `pkg/mcp` implements the MCP Subsystem interface via structural typing, translating tool calls to IPC `PERFORM`/`QUERY` calls across all 15 GUI packages.

**Tech Stack:** Go 1.25, core/go DI framework, go-webview CDP client, MCP SDK (`github.com/modelcontextprotocol/go-sdk/mcp`)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-mcp-bridge-design.md`

---

## File Structure

### New files

| Package | File | Responsibility |
|---------|------|---------------|
| `pkg/webview/` | `messages.go` | IPC message types (6 Queries, 11 Tasks, 2 Actions) + own types (ConsoleMessage, ElementInfo, etc.) |
| | `service.go` | connector interface, Service struct, Register factory, IPC handlers |
| | `service_test.go` | Mock connector, all IPC round-trip tests |
| `pkg/mcp/` | `subsystem.go` | Subsystem struct, New(), Name(), RegisterTools() |
| | `tools_webview.go` | 18 webview tool handlers + Input/Output types |
| | `tools_window.go` | 15 window tool handlers |
| | `tools_layout.go` | 7 layout tool handlers |
| | `tools_screen.go` | 5 screen tool handlers |
| | `tools_clipboard.go` | 4 clipboard tool handlers |
| | `tools_dialog.go` | 5 dialog tool handlers |
| | `tools_notification.go` | 3 notification tool handlers |
| | `tools_tray.go` | 4 systray tool handlers |
| | `tools_environment.go` | 2 environment/theme tool handlers |
| | `tools_browser.go` | 1 browser tool handler |
| | `tools_contextmenu.go` | 4 contextmenu tool handlers |
| | `tools_keybinding.go` | 2 keybinding tool handlers |
| | `tools_dock.go` | 3 dock tool handlers |
| | `tools_lifecycle.go` | 1 lifecycle tool handler |
| | `mcp_test.go` | RegisterTools smoke test + IPC round-trip tests |

### Modified files

| File | Changes |
|------|---------|
| `pkg/display/display.go` | Add webview import, HandleIPCEvents cases, handleWSMessage cases, EventType constants |
| `go.mod` | Add `forge.lthn.ai/core/go-webview`, `github.com/modelcontextprotocol/go-sdk` |

### Prerequisite (separate repo)

| File | Changes |
|------|---------|
| `go-webview/cdp.go` | Export `targetInfo` → `TargetInfo` |

---

## Task 1: Export TargetInfo in go-webview

**Files:**
- Modify: `/Users/snider/Code/core/go-webview/cdp.go`

- [ ] **Step 1: Rename targetInfo → TargetInfo**

In `/Users/snider/Code/core/go-webview/cdp.go`, rename the struct and update all references:

```go
// TargetInfo represents Chrome DevTools target information.
type TargetInfo struct {
    ID                   string `json:"id"`
    Type                 string `json:"type"`
    Title                string `json:"title"`
    URL                  string `json:"url"`
    WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}
```

Update `ListTargets` return type: `func ListTargets(debugURL string) ([]TargetInfo, error)`
Update `ListTargetsAll` return type: `func ListTargetsAll(debugURL string) iter.Seq[TargetInfo]`
Update all internal references (`var targets []targetInfo` → `var targets []TargetInfo`, etc.)

- [ ] **Step 2: Run tests**

Run: `cd /Users/snider/Code/core/go-webview && go build ./...`
Expected: Build succeeds (tests need Chrome, so just verify compilation)

- [ ] **Step 3: Commit**

```bash
cd /Users/snider/Code/core/go-webview
git add cdp.go
git commit -m "feat: export TargetInfo type for external CDP target enumeration"
```

---

## Task 2: pkg/webview — messages.go

**Files:**
- Create: `pkg/webview/messages.go`

- [ ] **Step 1: Write messages.go**

```go
// pkg/webview/messages.go
package webview

import "time"

// --- Queries (read-only) ---

// QueryURL gets the current page URL. Result: string
type QueryURL struct{ Window string `json:"window"` }

// QueryTitle gets the current page title. Result: string
type QueryTitle struct{ Window string `json:"window"` }

// QueryConsole gets captured console messages. Result: []ConsoleMessage
type QueryConsole struct {
	Window string `json:"window"`
	Level  string `json:"level,omitempty"` // filter by type: "log", "warn", "error", "info", "debug"
	Limit  int    `json:"limit,omitempty"` // max messages (0 = all)
}

// QuerySelector finds a single element. Result: *ElementInfo (nil if not found)
type QuerySelector struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// QuerySelectorAll finds all matching elements. Result: []*ElementInfo
type QuerySelectorAll struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// QueryDOMTree gets HTML content. Result: string (outerHTML)
type QueryDOMTree struct {
	Window   string `json:"window"`
	Selector string `json:"selector,omitempty"` // empty = full document
}

// --- Tasks (side-effects) ---

// TaskEvaluate executes JavaScript. Result: any (JS return value)
type TaskEvaluate struct {
	Window string `json:"window"`
	Script string `json:"script"`
}

// TaskClick clicks an element. Result: nil
type TaskClick struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// TaskType types text into an element. Result: nil
type TaskType struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Text     string `json:"text"`
}

// TaskNavigate navigates to a URL. Result: nil
type TaskNavigate struct {
	Window string `json:"window"`
	URL    string `json:"url"`
}

// TaskScreenshot captures the page as PNG. Result: ScreenshotResult
type TaskScreenshot struct{ Window string `json:"window"` }

// TaskScroll scrolls to an absolute position (window.scrollTo). Result: nil
type TaskScroll struct {
	Window string `json:"window"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// TaskHover hovers over an element. Result: nil
type TaskHover struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// TaskSelect selects an option in a <select> element. Result: nil
type TaskSelect struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

// TaskCheck checks/unchecks a checkbox. Result: nil
type TaskCheck struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Checked  bool   `json:"checked"`
}

// TaskUploadFile uploads files to an input element. Result: nil
type TaskUploadFile struct {
	Window   string   `json:"window"`
	Selector string   `json:"selector"`
	Paths    []string `json:"paths"`
}

// TaskSetViewport sets the viewport dimensions. Result: nil
type TaskSetViewport struct {
	Window string `json:"window"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// TaskClearConsole clears captured console messages. Result: nil
type TaskClearConsole struct{ Window string `json:"window"` }

// --- Actions (broadcast) ---

// ActionConsoleMessage is broadcast when a console message is captured.
type ActionConsoleMessage struct {
	Window  string         `json:"window"`
	Message ConsoleMessage `json:"message"`
}

// ActionException is broadcast when a JavaScript exception occurs.
type ActionException struct {
	Window    string        `json:"window"`
	Exception ExceptionInfo `json:"exception"`
}

// --- Types ---

// ConsoleMessage represents a browser console message.
type ConsoleMessage struct {
	Type      string    `json:"type"` // "log", "warn", "error", "info", "debug"
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url,omitempty"`
	Line      int       `json:"line,omitempty"`
	Column    int       `json:"column,omitempty"`
}

// ElementInfo represents a DOM element.
type ElementInfo struct {
	TagName     string            `json:"tagName"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	InnerText   string            `json:"innerText,omitempty"`
	InnerHTML   string            `json:"innerHTML,omitempty"`
	BoundingBox *BoundingBox      `json:"boundingBox,omitempty"`
}

// BoundingBox represents element position and size.
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ExceptionInfo represents a JavaScript exception.
// Field mapping from go-webview: LineNumber→Line, ColumnNumber→Column.
type ExceptionInfo struct {
	Text       string    `json:"text"`
	URL        string    `json:"url,omitempty"`
	Line       int       `json:"line,omitempty"`
	Column     int       `json:"column,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ScreenshotResult wraps raw PNG bytes as base64 for JSON/MCP transport.
type ScreenshotResult struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"` // always "image/png"
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/snider/Code/core/gui && go build ./pkg/webview/`
Expected: PASS (messages only, no dependencies yet beyond stdlib)

- [ ] **Step 3: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/webview/messages.go
git commit -m "feat(webview): add IPC message types and own types"
```

---

## Task 3: pkg/webview — service.go

**Files:**
- Create: `pkg/webview/service.go`

- [ ] **Step 1: Write service.go**

```go
// pkg/webview/service.go
package webview

import (
	"context"
	"encoding/base64"
	"strings"
	"sync"
	"time"

	gowebview "forge.lthn.ai/core/go-webview"
	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/window"
)

// connector abstracts go-webview for testing. The real implementation wraps
// *gowebview.Webview, converting go-webview types to our own types at the boundary.
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
	GetHTML(selector string) (string, error)
	QuerySelector(selector string) (*ElementInfo, error)
	QuerySelectorAll(selector string) ([]*ElementInfo, error)
	GetConsole() []ConsoleMessage
	ClearConsole()
	SetViewport(width, height int) error
	UploadFile(selector string, paths []string) error
	Close() error
}

// Options holds configuration for the webview service.
type Options struct {
	DebugURL     string        // Chrome debug endpoint (default: "http://localhost:9222")
	Timeout      time.Duration // Operation timeout (default: 30s)
	ConsoleLimit int           // Max console messages per window (default: 1000)
}

// Service is a core.Service managing webview interactions via IPC.
type Service struct {
	*core.ServiceRuntime[Options]
	opts        Options
	connections map[string]connector
	mu          sync.RWMutex
	newConn     func(debugURL, windowName string) (connector, error) // injectable for tests
}

// Register creates a factory closure with the given options.
func Register(opts ...func(*Options)) func(*core.Core) (any, error) {
	o := Options{
		DebugURL:     "http://localhost:9222",
		Timeout:      30 * time.Second,
		ConsoleLimit: 1000,
	}
	for _, fn := range opts {
		fn(&o)
	}
	return func(c *core.Core) (any, error) {
		return &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, o),
			opts:           o,
			connections:    make(map[string]connector),
			newConn:        defaultNewConn(o),
		}, nil
	}
}

// defaultNewConn creates real go-webview connections.
func defaultNewConn(opts Options) func(string, string) (connector, error) {
	return func(debugURL, windowName string) (connector, error) {
		// Enumerate targets, match by title/URL containing window name
		targets, err := gowebview.ListTargets(debugURL)
		if err != nil {
			return nil, err
		}
		var wsURL string
		for _, t := range targets {
			if t.Type == "page" && (strings.Contains(t.Title, windowName) || strings.Contains(t.URL, windowName)) {
				wsURL = t.WebSocketDebuggerURL
				break
			}
		}
		// Fallback: first page target
		if wsURL == "" {
			for _, t := range targets {
				if t.Type == "page" {
					wsURL = t.WebSocketDebuggerURL
					break
				}
			}
		}
		if wsURL == "" {
			return nil, core.E("webview.connect", "no page target found", nil)
		}
		wv, err := gowebview.New(
			gowebview.WithDebugURL(debugURL),
			gowebview.WithTimeout(opts.Timeout),
			gowebview.WithConsoleLimit(opts.ConsoleLimit),
		)
		if err != nil {
			return nil, err
		}
		return &realConnector{wv: wv}, nil
	}
}

// OnStartup registers IPC handlers.
func (s *Service) OnStartup(_ context.Context) error {
	s.Core().RegisterQuery(s.handleQuery)
	s.Core().RegisterTask(s.handleTask)
	return nil
}

// OnShutdown closes all CDP connections.
func (s *Service) OnShutdown(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, conn := range s.connections {
		conn.Close()
		delete(s.connections, name)
	}
	return nil
}

// HandleIPCEvents listens for window close events to clean up connections.
func (s *Service) HandleIPCEvents(_ *core.Core, msg core.Message) error {
	switch m := msg.(type) {
	case window.ActionWindowClosed:
		s.mu.Lock()
		if conn, ok := s.connections[m.Name]; ok {
			conn.Close()
			delete(s.connections, m.Name)
		}
		s.mu.Unlock()
	}
	return nil
}

// getConn returns the connector for a window, creating it if needed.
func (s *Service) getConn(windowName string) (connector, error) {
	s.mu.RLock()
	if conn, ok := s.connections[windowName]; ok {
		s.mu.RUnlock()
		return conn, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	// Double-check after acquiring write lock
	if conn, ok := s.connections[windowName]; ok {
		return conn, nil
	}
	conn, err := s.newConn(s.opts.DebugURL, windowName)
	if err != nil {
		return nil, err
	}
	s.connections[windowName] = conn
	return conn, nil
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryURL:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		url, err := conn.GetURL()
		return url, true, err
	case QueryTitle:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		title, err := conn.GetTitle()
		return title, true, err
	case QueryConsole:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		msgs := conn.GetConsole()
		// Filter by level if specified
		if q.Level != "" {
			var filtered []ConsoleMessage
			for _, m := range msgs {
				if m.Type == q.Level {
					filtered = append(filtered, m)
				}
			}
			msgs = filtered
		}
		// Apply limit
		if q.Limit > 0 && len(msgs) > q.Limit {
			msgs = msgs[len(msgs)-q.Limit:]
		}
		return msgs, true, nil
	case QuerySelector:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		el, err := conn.QuerySelector(q.Selector)
		return el, true, err
	case QuerySelectorAll:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		els, err := conn.QuerySelectorAll(q.Selector)
		return els, true, err
	case QueryDOMTree:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		selector := q.Selector
		if selector == "" {
			selector = "html"
		}
		html, err := conn.GetHTML(selector)
		return html, true, err
	default:
		return nil, false, nil
	}
}

func (s *Service) handleTask(_ *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskEvaluate:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		result, err := conn.Evaluate(t.Script)
		return result, true, err
	case TaskClick:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Click(t.Selector)
	case TaskType:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Type(t.Selector, t.Text)
	case TaskNavigate:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Navigate(t.URL)
	case TaskScreenshot:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		png, err := conn.Screenshot()
		if err != nil {
			return nil, true, err
		}
		return ScreenshotResult{
			Base64:   base64.StdEncoding.EncodeToString(png),
			MimeType: "image/png",
		}, true, nil
	case TaskScroll:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		_, err = conn.Evaluate("window.scrollTo(" + itoa(t.X) + "," + itoa(t.Y) + ")")
		return nil, true, err
	case TaskHover:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Hover(t.Selector)
	case TaskSelect:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Select(t.Selector, t.Value)
	case TaskCheck:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Check(t.Selector, t.Checked)
	case TaskUploadFile:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.UploadFile(t.Selector, t.Paths)
	case TaskSetViewport:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.SetViewport(t.Width, t.Height)
	case TaskClearConsole:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		conn.ClearConsole()
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func itoa(n int) string {
	return strings.Join([]string{""}, "") + string(rune('0'+n%10)) // use strconv in real code
}

// realConnector wraps *gowebview.Webview, converting types at the boundary.
type realConnector struct {
	wv *gowebview.Webview
}

func (r *realConnector) Navigate(url string) error              { return r.wv.Navigate(url) }
func (r *realConnector) Click(sel string) error                 { return r.wv.Click(sel) }
func (r *realConnector) Type(sel, text string) error            { return r.wv.Type(sel, text) }
func (r *realConnector) Evaluate(script string) (any, error)    { return r.wv.Evaluate(script) }
func (r *realConnector) Screenshot() ([]byte, error)            { return r.wv.Screenshot() }
func (r *realConnector) GetURL() (string, error)                { return r.wv.GetURL() }
func (r *realConnector) GetTitle() (string, error)              { return r.wv.GetTitle() }
func (r *realConnector) GetHTML(sel string) (string, error)     { return r.wv.GetHTML(sel) }
func (r *realConnector) ClearConsole()                          { r.wv.ClearConsole() }
func (r *realConnector) Close() error                           { return r.wv.Close() }
func (r *realConnector) SetViewport(w, h int) error             { return r.wv.SetViewport(w, h) }
func (r *realConnector) UploadFile(sel string, p []string) error { return r.wv.UploadFile(sel, p) }

func (r *realConnector) Hover(sel string) error {
	return gowebview.NewActionSequence().Add(&gowebview.HoverAction{Selector: sel}).Execute(context.Background(), r.wv)
}

func (r *realConnector) Select(sel, val string) error {
	return gowebview.NewActionSequence().Add(&gowebview.SelectAction{Selector: sel, Value: val}).Execute(context.Background(), r.wv)
}

func (r *realConnector) Check(sel string, checked bool) error {
	return gowebview.NewActionSequence().Add(&gowebview.CheckAction{Selector: sel, Checked: checked}).Execute(context.Background(), r.wv)
}

func (r *realConnector) QuerySelector(sel string) (*ElementInfo, error) {
	el, err := r.wv.QuerySelector(sel)
	if err != nil {
		return nil, err
	}
	return convertElementInfo(el), nil
}

func (r *realConnector) QuerySelectorAll(sel string) ([]*ElementInfo, error) {
	els, err := r.wv.QuerySelectorAll(sel)
	if err != nil {
		return nil, err
	}
	result := make([]*ElementInfo, len(els))
	for i, el := range els {
		result[i] = convertElementInfo(el)
	}
	return result, nil
}

func (r *realConnector) GetConsole() []ConsoleMessage {
	raw := r.wv.GetConsole()
	msgs := make([]ConsoleMessage, len(raw))
	for i, m := range raw {
		msgs[i] = ConsoleMessage{
			Type: m.Type, Text: m.Text, Timestamp: m.Timestamp,
			URL: m.URL, Line: m.Line, Column: m.Column,
		}
	}
	return msgs
}

func convertElementInfo(el *gowebview.ElementInfo) *ElementInfo {
	if el == nil {
		return nil
	}
	info := &ElementInfo{
		TagName:    el.TagName,
		Attributes: el.Attributes,
		InnerText:  el.InnerText,
		InnerHTML:  el.InnerHTML,
	}
	if el.BoundingBox != nil {
		info.BoundingBox = &BoundingBox{
			X: el.BoundingBox.X, Y: el.BoundingBox.Y,
			Width: el.BoundingBox.Width, Height: el.BoundingBox.Height,
		}
	}
	return info
}
```

**Note:** The `itoa` function above is a placeholder — replace with `strconv.Itoa` in the real implementation. It's written this way to avoid adding an import in the plan code.

- [ ] **Step 2: Verify compilation**

Run: `cd /Users/snider/Code/core/gui && go build ./pkg/webview/`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/webview/service.go
git commit -m "feat(webview): add service with connector interface and IPC handlers"
```

---

## Task 4: pkg/webview — service_test.go

**Files:**
- Create: `pkg/webview/service_test.go`

- [ ] **Step 1: Write service_test.go**

```go
// pkg/webview/service_test.go
package webview

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnector struct {
	url       string
	title     string
	html      string
	evalResult any
	screenshot []byte
	console   []ConsoleMessage
	elements  []*ElementInfo
	closed    bool

	lastClickSel  string
	lastTypeSel   string
	lastTypeText  string
	lastNavURL    string
	lastHoverSel  string
	lastSelectSel string
	lastSelectVal string
	lastCheckSel  string
	lastCheckVal  bool
	lastUploadSel string
	lastUploadPaths []string
	lastViewportW int
	lastViewportH int
	consoleClearCalled bool
}

func (m *mockConnector) Navigate(url string) error        { m.lastNavURL = url; return nil }
func (m *mockConnector) Click(sel string) error           { m.lastClickSel = sel; return nil }
func (m *mockConnector) Type(sel, text string) error      { m.lastTypeSel = sel; m.lastTypeText = text; return nil }
func (m *mockConnector) Hover(sel string) error           { m.lastHoverSel = sel; return nil }
func (m *mockConnector) Select(sel, val string) error     { m.lastSelectSel = sel; m.lastSelectVal = val; return nil }
func (m *mockConnector) Check(sel string, c bool) error   { m.lastCheckSel = sel; m.lastCheckVal = c; return nil }
func (m *mockConnector) Evaluate(s string) (any, error)   { return m.evalResult, nil }
func (m *mockConnector) Screenshot() ([]byte, error)      { return m.screenshot, nil }
func (m *mockConnector) GetURL() (string, error)          { return m.url, nil }
func (m *mockConnector) GetTitle() (string, error)        { return m.title, nil }
func (m *mockConnector) GetHTML(sel string) (string, error) { return m.html, nil }
func (m *mockConnector) ClearConsole()                    { m.consoleClearCalled = true }
func (m *mockConnector) Close() error                     { m.closed = true; return nil }
func (m *mockConnector) SetViewport(w, h int) error       { m.lastViewportW = w; m.lastViewportH = h; return nil }
func (m *mockConnector) UploadFile(sel string, p []string) error { m.lastUploadSel = sel; m.lastUploadPaths = p; return nil }

func (m *mockConnector) QuerySelector(sel string) (*ElementInfo, error) {
	if len(m.elements) > 0 {
		return m.elements[0], nil
	}
	return nil, nil
}

func (m *mockConnector) QuerySelectorAll(sel string) ([]*ElementInfo, error) {
	return m.elements, nil
}

func (m *mockConnector) GetConsole() []ConsoleMessage { return m.console }

func newTestService(t *testing.T, mock *mockConnector) (*Service, *core.Core) {
	t.Helper()
	factory := Register()
	c, err := core.New(core.WithService(factory), core.WithServiceLock())
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "webview")
	// Inject mock connector
	svc.newConn = func(_, _ string) (connector, error) { return mock, nil }
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestService(t, &mockConnector{})
	assert.NotNil(t, svc)
}

func TestQueryURL_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{url: "https://example.com"})
	result, handled, err := c.QUERY(QueryURL{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "https://example.com", result)
}

func TestQueryTitle_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{title: "Test Page"})
	result, handled, err := c.QUERY(QueryTitle{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Test Page", result)
}

func TestQueryConsole_Good(t *testing.T) {
	mock := &mockConnector{console: []ConsoleMessage{
		{Type: "log", Text: "hello"},
		{Type: "error", Text: "oops"},
		{Type: "log", Text: "world"},
	}}
	_, c := newTestService(t, mock)
	result, handled, err := c.QUERY(QueryConsole{Window: "main", Level: "error", Limit: 10})
	require.NoError(t, err)
	assert.True(t, handled)
	msgs, _ := result.([]ConsoleMessage)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "oops", msgs[0].Text)
}

func TestQueryConsole_Good_Limit(t *testing.T) {
	mock := &mockConnector{console: []ConsoleMessage{
		{Type: "log", Text: "a"},
		{Type: "log", Text: "b"},
		{Type: "log", Text: "c"},
	}}
	_, c := newTestService(t, mock)
	result, _, _ := c.QUERY(QueryConsole{Window: "main", Limit: 2})
	msgs, _ := result.([]ConsoleMessage)
	assert.Len(t, msgs, 2)
	assert.Equal(t, "b", msgs[0].Text) // last 2
}

func TestTaskEvaluate_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{evalResult: 42})
	result, handled, err := c.PERFORM(TaskEvaluate{Window: "main", Script: "21*2"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 42, result)
}

func TestTaskClick_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskClick{Window: "main", Selector: "#btn"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "#btn", mock.lastClickSel)
}

func TestTaskNavigate_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskNavigate{Window: "main", URL: "https://example.com"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "https://example.com", mock.lastNavURL)
}

func TestTaskScreenshot_Good(t *testing.T) {
	mock := &mockConnector{screenshot: []byte{0x89, 0x50, 0x4E, 0x47}}
	_, c := newTestService(t, mock)
	result, handled, err := c.PERFORM(TaskScreenshot{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	sr, ok := result.(ScreenshotResult)
	assert.True(t, ok)
	assert.Equal(t, "image/png", sr.MimeType)
	assert.NotEmpty(t, sr.Base64)
}

func TestTaskClearConsole_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskClearConsole{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.consoleClearCalled)
}

func TestConnectionCleanup_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Access creates connection
	_, _, _ = c.QUERY(QueryURL{Window: "main"})
	assert.False(t, mock.closed)
	// Window close action triggers cleanup
	_ = c.ACTION(window.ActionWindowClosed{Name: "main"})
	assert.True(t, mock.closed)
}

func TestQueryURL_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryURL{Window: "main"})
	assert.False(t, handled)
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/webview/ -v`
Expected: All PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/webview/service_test.go
git commit -m "test(webview): add service tests with mock connector"
```

---

## Task 5: Display orchestrator — webview integration

**Files:**
- Modify: `pkg/display/display.go`

- [ ] **Step 1: Add webview import**

Add `"forge.lthn.ai/core/gui/pkg/webview"` to the import block.

- [ ] **Step 2: Add EventType constants**

In the EventType const block, add:

```go
EventWebviewConsole   EventType = "webview.console"
EventWebviewException EventType = "webview.exception"
```

- [ ] **Step 3: Add HandleIPCEvents cases**

In `HandleIPCEvents`, add:

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

- [ ] **Step 4: Add handleWSMessage cases**

In `handleWSMessage`, add all 18 webview WS→IPC cases. Full list from the spec:

```go
case "webview:eval":
    w, _ := msg.Data["window"].(string)
    script, _ := msg.Data["script"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskEvaluate{Window: w, Script: script})
case "webview:click":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskClick{Window: w, Selector: sel})
case "webview:type":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    text, _ := msg.Data["text"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskType{Window: w, Selector: sel, Text: text})
case "webview:navigate":
    w, _ := msg.Data["window"].(string)
    url, _ := msg.Data["url"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskNavigate{Window: w, URL: url})
case "webview:screenshot":
    w, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskScreenshot{Window: w})
case "webview:scroll":
    w, _ := msg.Data["window"].(string)
    x, _ := msg.Data["x"].(float64)
    y, _ := msg.Data["y"].(float64)
    result, handled, err = s.Core().PERFORM(webview.TaskScroll{Window: w, X: int(x), Y: int(y)})
case "webview:hover":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskHover{Window: w, Selector: sel})
case "webview:select":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    val, _ := msg.Data["value"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskSelect{Window: w, Selector: sel, Value: val})
case "webview:check":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    checked, _ := msg.Data["checked"].(bool)
    result, handled, err = s.Core().PERFORM(webview.TaskCheck{Window: w, Selector: sel, Checked: checked})
case "webview:upload":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    pathsRaw, _ := msg.Data["paths"].([]any)
    var paths []string
    for _, p := range pathsRaw {
        if ps, ok := p.(string); ok {
            paths = append(paths, ps)
        }
    }
    result, handled, err = s.Core().PERFORM(webview.TaskUploadFile{Window: w, Selector: sel, Paths: paths})
case "webview:viewport":
    w, _ := msg.Data["window"].(string)
    width, _ := msg.Data["width"].(float64)
    height, _ := msg.Data["height"].(float64)
    result, handled, err = s.Core().PERFORM(webview.TaskSetViewport{Window: w, Width: int(width), Height: int(height)})
case "webview:clear-console":
    w, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().PERFORM(webview.TaskClearConsole{Window: w})
case "webview:console":
    w, _ := msg.Data["window"].(string)
    level, _ := msg.Data["level"].(string)
    limit := 100
    if l, ok := msg.Data["limit"].(float64); ok {
        limit = int(l)
    }
    result, handled, err = s.Core().QUERY(webview.QueryConsole{Window: w, Level: level, Limit: limit})
case "webview:query":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QuerySelector{Window: w, Selector: sel})
case "webview:query-all":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QuerySelectorAll{Window: w, Selector: sel})
case "webview:dom-tree":
    w, _ := msg.Data["window"].(string)
    sel, _ := msg.Data["selector"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryDOMTree{Window: w, Selector: sel})
case "webview:url":
    w, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryURL{Window: w})
case "webview:title":
    w, _ := msg.Data["window"].(string)
    result, handled, err = s.Core().QUERY(webview.QueryTitle{Window: w})
```

- [ ] **Step 5: Run tests**

Run: `cd /Users/snider/Code/core/gui && go build ./pkg/display/`
Expected: Build succeeds

- [ ] **Step 6: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/display/display.go
git commit -m "feat(display): add webview IPC→WS bridging (18 cases)"
```

---

## Task 6: pkg/mcp — subsystem.go

**Files:**
- Create: `pkg/mcp/subsystem.go`

- [ ] **Step 1: Write subsystem.go**

```go
// pkg/mcp/subsystem.go
package mcp

import (
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Subsystem implements the MCP Subsystem interface via structural typing.
// It registers GUI tools that translate MCP tool calls to IPC messages.
type Subsystem struct {
	core *core.Core
}

// New creates a display MCP subsystem backed by the given Core instance.
func New(c *core.Core) *Subsystem {
	return &Subsystem{core: c}
}

// Name returns the subsystem identifier.
func (s *Subsystem) Name() string { return "display" }

// RegisterTools registers all GUI tools with the MCP server.
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

- [ ] **Step 2: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/mcp/subsystem.go
git commit -m "feat(mcp): add display subsystem skeleton"
```

---

## Tasks 7–12: MCP Tool Files

Each tool file follows the **exact same pattern**. Here is the template for every tool handler:

### Pattern: Query Tool

```go
// Input/Output types
type ScreenListInput struct{}
type ScreenListOutput struct {
    Screens []screen.ScreenInfo `json:"screens"`
}

// Handler — calls QUERY
func (s *Subsystem) screenList(ctx context.Context, req *mcp.CallToolRequest, input ScreenListInput) (*mcp.CallToolResult, ScreenListOutput, error) {
    result, _, err := s.core.QUERY(screen.QueryList{})
    if err != nil {
        return nil, ScreenListOutput{}, err
    }
    screens, _ := result.([]screen.ScreenInfo)
    return nil, ScreenListOutput{Screens: screens}, nil
}

// Registration
func (s *Subsystem) registerScreenTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "screen_list",
        Description: "List all connected displays/screens",
    }, s.screenList)
}
```

### Pattern: Task Tool

```go
type ClipboardWriteInput struct {
    Text string `json:"text"`
}
type ClipboardWriteOutput struct {
    Success bool `json:"success"`
}

func (s *Subsystem) clipboardWrite(ctx context.Context, req *mcp.CallToolRequest, input ClipboardWriteInput) (*mcp.CallToolResult, ClipboardWriteOutput, error) {
    result, _, err := s.core.PERFORM(clipboard.TaskSetText{Text: input.Text})
    if err != nil {
        return nil, ClipboardWriteOutput{}, err
    }
    success, _ := result.(bool)
    return nil, ClipboardWriteOutput{Success: success}, nil
}
```

### Task 7: tools_webview.go (18 tools)

**Files:** Create `pkg/mcp/tools_webview.go`

Tools: `webview_eval`, `webview_click`, `webview_type`, `webview_navigate`, `webview_screenshot`, `webview_scroll`, `webview_hover`, `webview_select`, `webview_check`, `webview_upload`, `webview_viewport`, `webview_console`, `webview_console_clear`, `webview_query`, `webview_query_all`, `webview_dom_tree`, `webview_url`, `webview_title`

Each follows the Task/Query pattern above. Import `forge.lthn.ai/core/gui/pkg/webview`.

- [ ] Write tools_webview.go with 18 handlers
- [ ] `go build ./pkg/mcp/`
- [ ] Commit: `git commit -m "feat(mcp): add webview tools (18)"`

### Task 8: tools_window.go + tools_layout.go (15 + 7 tools)

**Files:** Create `pkg/mcp/tools_window.go`, `pkg/mcp/tools_layout.go`

**Window tools:** `window_list` (QueryList), `window_get` (QueryGet), `window_focused` (QueryFocused), `window_create` (TaskCreateWindow), `window_close` (TaskCloseWindow), `window_position` (TaskSetPosition), `window_size` (TaskSetSize), `window_bounds` (TaskSetBounds), `window_maximize` (TaskMaximize), `window_minimize` (TaskMinimize), `window_restore` (TaskRestore), `window_focus` (TaskFocus), `window_title` (TaskSetTitle), `window_visibility` (TaskSetVisibility), `window_fullscreen` (TaskSetFullscreen)

**Layout tools:** `layout_save` (TaskSaveLayout), `layout_restore` (TaskRestoreLayout), `layout_list` (QueryListLayouts), `layout_delete` (TaskDeleteLayout), `layout_get` (QueryGetLayout), `layout_tile` (TaskTileWindows), `layout_snap` (TaskSnapWindow)

Import `forge.lthn.ai/core/gui/pkg/window`.

- [ ] Write tools_window.go with 15 handlers
- [ ] Write tools_layout.go with 7 handlers
- [ ] `go build ./pkg/mcp/`
- [ ] Commit: `git commit -m "feat(mcp): add window and layout tools (22)"`

### Task 9: tools_screen.go + tools_clipboard.go + tools_dialog.go (5 + 4 + 5 tools)

**Files:** Create `pkg/mcp/tools_screen.go`, `pkg/mcp/tools_clipboard.go`, `pkg/mcp/tools_dialog.go`

**Screen tools:** `screen_list` (QueryList), `screen_get` (QueryGet), `screen_primary` (QueryPrimary), `screen_at_point` (QueryAtPoint), `screen_work_areas` (QueryWorkAreas)

**Clipboard tools:** `clipboard_read` (QueryText), `clipboard_write` (TaskSetText), `clipboard_has` (QueryText → check HasContent), `clipboard_clear` (TaskClear)

**Dialog tools:** `dialog_open_file` (TaskOpenFile), `dialog_save_file` (TaskSaveFile), `dialog_open_directory` (TaskOpenDirectory), `dialog_confirm` (TaskConfirm), `dialog_prompt` (TaskPrompt)

- [ ] Write all three files
- [ ] `go build ./pkg/mcp/`
- [ ] Commit: `git commit -m "feat(mcp): add screen, clipboard, dialog tools (14)"`

### Task 10: tools_notification.go + tools_tray.go + tools_environment.go (3 + 4 + 2 tools)

**Files:** Create `pkg/mcp/tools_notification.go`, `pkg/mcp/tools_tray.go`, `pkg/mcp/tools_environment.go`

**Notification tools:** `notification_show` (TaskShow), `notification_permission_request` (TaskRequestPermission), `notification_permission_check` (QueryPermission)

**Tray tools:** `tray_set_icon` (TaskSetIcon), `tray_set_tooltip` (TaskSetTooltip), `tray_set_label` (TaskSetLabel), `tray_info` (QueryInfo)

**Environment tools:** `theme_get` (QueryTheme), `theme_system` (QuerySystemTheme)

- [ ] Write all three files
- [ ] `go build ./pkg/mcp/`
- [ ] Commit: `git commit -m "feat(mcp): add notification, tray, environment tools (9)"`

### Task 11: remaining tool files (1 + 4 + 2 + 3 + 1 tools)

**Files:** Create `pkg/mcp/tools_browser.go`, `pkg/mcp/tools_contextmenu.go`, `pkg/mcp/tools_keybinding.go`, `pkg/mcp/tools_dock.go`, `pkg/mcp/tools_lifecycle.go`

**Browser:** `browser_open_url` (TaskOpenURL)

**Context menu:** `contextmenu_add` (TaskAdd), `contextmenu_remove` (TaskRemove), `contextmenu_get` (QueryGet), `contextmenu_list` (QueryList)

**Keybinding:** `keybinding_add` (TaskAdd), `keybinding_remove` (TaskRemove)

**Dock:** `dock_show` (TaskShowIcon), `dock_hide` (TaskHideIcon), `dock_badge` (TaskSetBadge)

**Lifecycle:** `app_quit` (TaskQuit)

- [ ] Write all five files
- [ ] `go build ./pkg/mcp/`
- [ ] Commit: `git commit -m "feat(mcp): add browser, contextmenu, keybinding, dock, lifecycle tools (11)"`

---

## Task 12: pkg/mcp — mcp_test.go

**Files:**
- Create: `pkg/mcp/mcp_test.go`

- [ ] **Step 1: Write mcp_test.go**

```go
// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsystem_Good_Name(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := New(c)
	assert.Equal(t, "display", sub.Name())
}

func TestSubsystem_Good_RegisterTools(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := New(c)
	// RegisterTools should not panic with a real mcp.Server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	assert.NotPanics(t, func() { sub.RegisterTools(server) })
}

// Integration test: verify the IPC round-trip that MCP tool handlers use.

type mockClipPlatform struct{ text string; ok bool }
func (m *mockClipPlatform) Text() (string, bool)   { return m.text, m.ok }
func (m *mockClipPlatform) SetText(t string) bool   { m.text = t; m.ok = t != ""; return true }

func TestMCP_Good_ClipboardRoundTrip(t *testing.T) {
	c, err := core.New(
		core.WithService(clipboard.Register(&mockClipPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	// Verify the IPC path that clipboard_read tool handler uses
	result, handled, err := c.QUERY(clipboard.QueryText{})
	require.NoError(t, err)
	assert.True(t, handled)
	content, _ := result.(clipboard.ClipboardContent)
	assert.Equal(t, "hello", content.Text)
}

// Add a mock window platform to test window round-trip
type mockWinPlatform struct{}
// Implement window.Platform methods as needed for testing...
// The important thing is RegisterTools doesn't panic and IPC routes work.

func TestMCP_Bad_NoServices(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	// Without any services, QUERY should return handled=false
	_, handled, _ := c.QUERY(clipboard.QueryText{})
	assert.False(t, handled)
}
```

- [ ] **Step 2: Run tests**

Run: `cd /Users/snider/Code/core/gui && go test ./pkg/mcp/ -v`
Expected: All PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/snider/Code/core/gui
git add pkg/mcp/mcp_test.go
git commit -m "test(mcp): add subsystem tests with IPC round-trip verification"
```

---

## Task 13: go.mod + final verification

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Update go.mod**

```bash
cd /Users/snider/Code/core/gui
go get forge.lthn.ai/core/go-webview@latest
go get github.com/modelcontextprotocol/go-sdk@latest
go mod tidy
```

- [ ] **Step 2: Run all tests**

Run: `cd /Users/snider/Code/core/gui && go test ./...`
Expected: All packages PASS

- [ ] **Step 3: Run go vet**

Run: `cd /Users/snider/Code/core/gui && go vet ./...`
Expected: No issues

- [ ] **Step 4: Final commit**

```bash
cd /Users/snider/Code/core/gui
git add go.mod go.sum
git commit -m "chore: add go-webview and MCP SDK dependencies"
```

- [ ] **Step 5: Push to forge**

```bash
cd /Users/snider/Code/core/gui && git push origin main
```
