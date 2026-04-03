// pkg/webview/service_test.go
package webview

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnector struct {
	url         string
	title       string
	html        string
	evalResult  any
	evalFn      func(script string) (any, error)
	screenshot  []byte
	console     []ConsoleMessage
	elements    []*ElementInfo
	closed      bool
	pdfBytes    []byte
	printCalled bool

	lastClickSel       string
	lastTypeSel        string
	lastTypeText       string
	lastNavURL         string
	lastHoverSel       string
	lastSelectSel      string
	lastSelectVal      string
	lastCheckSel       string
	lastCheckVal       bool
	lastUploadSel      string
	lastUploadPaths    []string
	lastViewportW      int
	lastViewportH      int
	consoleClearCalled bool
	lastEvalScript     string
}

func (m *mockConnector) Navigate(url string) error { m.lastNavURL = url; return nil }
func (m *mockConnector) Click(sel string) error    { m.lastClickSel = sel; return nil }
func (m *mockConnector) Type(sel, text string) error {
	m.lastTypeSel = sel
	m.lastTypeText = text
	return nil
}
func (m *mockConnector) Hover(sel string) error { m.lastHoverSel = sel; return nil }
func (m *mockConnector) Select(sel, val string) error {
	m.lastSelectSel = sel
	m.lastSelectVal = val
	return nil
}
func (m *mockConnector) Check(sel string, c bool) error {
	m.lastCheckSel = sel
	m.lastCheckVal = c
	return nil
}
func (m *mockConnector) Evaluate(s string) (any, error) {
	m.lastEvalScript = s
	if m.evalFn != nil {
		return m.evalFn(s)
	}
	return m.evalResult, nil
}
func (m *mockConnector) Screenshot() ([]byte, error)        { return m.screenshot, nil }
func (m *mockConnector) GetURL() (string, error)            { return m.url, nil }
func (m *mockConnector) GetTitle() (string, error)          { return m.title, nil }
func (m *mockConnector) GetHTML(sel string) (string, error) { return m.html, nil }
func (m *mockConnector) ClearConsole()                      { m.consoleClearCalled = true }
func (m *mockConnector) Print() error                       { m.printCalled = true; return nil }
func (m *mockConnector) Close() error                       { m.closed = true; return nil }
func (m *mockConnector) SetViewport(w, h int) error {
	m.lastViewportW = w
	m.lastViewportH = h
	return nil
}
func (m *mockConnector) PrintToPDF() ([]byte, error) {
	if len(m.pdfBytes) == 0 {
		return []byte("%PDF-1.4\n"), nil
	}
	return m.pdfBytes, nil
}
func (m *mockConnector) UploadFile(sel string, p []string) error {
	m.lastUploadSel = sel
	m.lastUploadPaths = p
	return nil
}

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
	factory := Register(Options{})
	c, err := core.New(
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithService(factory),
		core.WithServiceLock(),
	)
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

func TestQueryExceptions_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{})

	require.NoError(t, c.ACTION(ActionException{
		Window: "main",
		Exception: ExceptionInfo{
			Text:       "boom",
			URL:        "https://example.com/app.js",
			Line:       12,
			Column:     4,
			StackTrace: "Error: boom",
		},
	}))

	result, handled, err := c.QUERY(QueryExceptions{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	exceptions, _ := result.([]ExceptionInfo)
	require.Len(t, exceptions, 1)
	assert.Equal(t, "boom", exceptions[0].Text)
	assert.Equal(t, 12, exceptions[0].Line)
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

func TestTaskScreenshotElement_Good(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x * 40), G: uint8(y * 40), B: 200, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	mock := &mockConnector{
		screenshot: buf.Bytes(),
		evalFn: func(script string) (any, error) {
			return map[string]any{
				"left":             1.0,
				"top":              1.0,
				"width":            2.0,
				"height":           2.0,
				"devicePixelRatio": 1.0,
			}, nil
		},
	}
	_, c := newTestService(t, mock)

	result, handled, err := c.PERFORM(TaskScreenshotElement{Window: "main", Selector: "#card"})
	require.NoError(t, err)
	assert.True(t, handled)
	sr, ok := result.(ScreenshotResult)
	require.True(t, ok)

	raw, err := base64.StdEncoding.DecodeString(sr.Base64)
	require.NoError(t, err)
	decoded, err := png.Decode(bytes.NewReader(raw))
	require.NoError(t, err)
	assert.Equal(t, image.Rect(0, 0, 2, 2), decoded.Bounds())
}

func TestTaskClearConsole_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskClearConsole{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.consoleClearCalled)
}

func TestTaskDevTools_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{})
	_, _, err := c.PERFORM(window.TaskOpenWindow{Opts: []window.WindowOption{window.WithName("main")}})
	require.NoError(t, err)
	_, handled, err := c.PERFORM(TaskOpenDevTools{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	_, handled, err = c.PERFORM(TaskCloseDevTools{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestDiagnosticsQueries_Good(t *testing.T) {
	mock := &mockConnector{
		evalFn: func(script string) (any, error) {
			switch {
			case strings.Contains(script, "getComputedStyle"):
				return map[string]any{"color": "rgb(1, 2, 3)"}, nil
			case strings.Contains(script, "performance.getEntriesByType(\"navigation\")"):
				return map[string]any{
					"navigationStart":      1.0,
					"domContentLoaded":     2.0,
					"loadEventEnd":         3.0,
					"firstPaint":           4.0,
					"firstContentfulPaint": 5.0,
					"usedJSHeapSize":       6.0,
					"totalJSHeapSize":      7.0,
				}, nil
			case strings.Contains(script, "performance.getEntriesByType(\"resource\")"):
				return []any{
					map[string]any{"name": "app.js", "entryType": "resource", "initiatorType": "script"},
				}, nil
			case strings.Contains(script, "window.__coreNetworkLog"):
				return []any{
					map[string]any{"url": "https://example.com", "method": "GET", "status": 200, "resource": "fetch"},
				}, nil
			default:
				return nil, nil
			}
		},
	}
	_, c := newTestService(t, mock)

	style, handled, err := c.QUERY(QueryComputedStyle{Window: "main", Selector: "#app"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "rgb(1, 2, 3)", style.(map[string]string)["color"])

	perf, handled, err := c.QUERY(QueryPerformance{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1.0, perf.(PerformanceMetrics).NavigationStart)

	resources, handled, err := c.QUERY(QueryResources{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Len(t, resources.([]ResourceEntry), 1)

	network, handled, err := c.QUERY(QueryNetwork{Window: "main", Limit: 10})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Len(t, network.([]NetworkEntry), 1)
}

func TestDiagnosticsTasks_Good(t *testing.T) {
	mock := &mockConnector{pdfBytes: []byte("%PDF-1.7")}
	_, c := newTestService(t, mock)

	_, handled, err := c.PERFORM(TaskHighlight{Window: "main", Selector: "#app", Colour: "#00ff00"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Contains(t, mock.lastEvalScript, "outline")

	_, handled, err = c.PERFORM(TaskInjectNetworkLogging{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Contains(t, mock.lastEvalScript, "__coreNetworkLog")

	_, handled, err = c.PERFORM(TaskClearNetworkLog{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)

	_, handled, err = c.PERFORM(TaskPrint{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.printCalled)

	result, handled, err := c.PERFORM(TaskExportPDF{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	pdf, ok := result.(PDFResult)
	require.True(t, ok)
	assert.Equal(t, "application/pdf", pdf.MimeType)
	assert.NotEmpty(t, pdf.Base64)
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
