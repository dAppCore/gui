// pkg/webview/service_test.go
package webview

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnector struct {
	url        string
	title      string
	html       string
	evalResult any
	screenshot []byte
	console    []ConsoleMessage
	elements   []*ElementInfo
	closed     bool

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

	zoom          float64
	lastZoomSet   float64
	printToPDF    bool
	printCalled   bool
	printPDFBytes []byte
	printErr      error
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
func (m *mockConnector) Evaluate(s string) (any, error)     { return m.evalResult, nil }
func (m *mockConnector) Screenshot() ([]byte, error)        { return m.screenshot, nil }
func (m *mockConnector) GetURL() (string, error)            { return m.url, nil }
func (m *mockConnector) GetTitle() (string, error)          { return m.title, nil }
func (m *mockConnector) GetHTML(sel string) (string, error) { return m.html, nil }
func (m *mockConnector) ClearConsole()                      { m.consoleClearCalled = true }
func (m *mockConnector) Close() error                       { m.closed = true; return nil }
func (m *mockConnector) SetViewport(w, h int) error {
	m.lastViewportW = w
	m.lastViewportH = h
	return nil
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

func (m *mockConnector) GetZoom() (float64, error) {
	if m.zoom == 0 {
		return 1.0, nil
	}
	return m.zoom, nil
}

func (m *mockConnector) SetZoom(zoom float64) error {
	m.lastZoomSet = zoom
	m.zoom = zoom
	return nil
}

func (m *mockConnector) Print(toPDF bool) ([]byte, error) {
	m.printCalled = true
	m.printToPDF = toPDF
	return m.printPDFBytes, m.printErr
}

func newTestService(t *testing.T, mock *mockConnector) (*Service, *core.Core) {
	t.Helper()
	factory := Register()
	c := core.New(core.WithService(factory), core.WithServiceLock())
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "webview")
	// Inject mock connector
	svc.newConn = func(_, _ string) (connector, error) { return mock, nil }
	return svc, c
}

func newTestServiceWithWindow(t *testing.T, mock *mockConnector) (*Service, *window.MockPlatform, *core.Core) {
	t.Helper()
	windowPlatform := window.NewMockPlatform()
	c := core.New(
		core.WithService(window.Register(windowPlatform)),
		core.WithService(Register()),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	result := taskRun(c, "window.open", window.TaskOpenWindow{
		Window: &window.Window{Name: "main", Title: "Main", Width: 800, Height: 600},
	})
	require.True(t, result.OK)

	svc := core.MustServiceFor[*Service](c, "webview")
	svc.newConn = func(_, _ string) (connector, error) { return mock, nil }
	return svc, windowPlatform, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestService(t, &mockConnector{})
	assert.NotNil(t, svc)
}

func TestQueryURL_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{url: "https://example.com"})
	r := c.QUERY(QueryURL{Window: "main"})
	require.True(t, r.OK)
	assert.Equal(t, "https://example.com", r.Value)
}

func TestQueryTitle_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{title: "Test Page"})
	r := c.QUERY(QueryTitle{Window: "main"})
	require.True(t, r.OK)
	assert.Equal(t, "Test Page", r.Value)
}

func TestQueryConsole_Good(t *testing.T) {
	mock := &mockConnector{console: []ConsoleMessage{
		{Type: "log", Text: "hello"},
		{Type: "error", Text: "oops"},
		{Type: "log", Text: "world"},
	}}
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryConsole{Window: "main", Level: "error", Limit: 10})
	require.True(t, r.OK)
	msgs, _ := r.Value.([]ConsoleMessage)
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
	r := c.QUERY(QueryConsole{Window: "main", Limit: 2})
	msgs, _ := r.Value.([]ConsoleMessage)
	assert.Len(t, msgs, 2)
	assert.Equal(t, "b", msgs[0].Text) // last 2
}

func TestTaskEvaluate_Good(t *testing.T) {
	_, c := newTestService(t, &mockConnector{evalResult: 42})
	r := taskRun(c, "webview.evaluate", TaskEvaluate{Window: "main", Script: "21*2"})
	require.True(t, r.OK)
	assert.Equal(t, 42, r.Value)
}

func TestTaskEvaluate_Bad_EmptyWindow(t *testing.T) {
	_, c := newTestService(t, &mockConnector{evalResult: 42})
	r := taskRun(c, "webview.evaluate", TaskEvaluate{Window: " ", Script: "21*2"})
	assert.False(t, r.OK)
}

func TestExactWindowTargetMatch_Good(t *testing.T) {
	assert.True(t, exactWindowTargetMatch("main", "main"))
	assert.False(t, exactWindowTargetMatch("main - docs", "main"))
}

func TestTaskClick_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.click", TaskClick{Window: "main", Selector: "#btn"})
	require.True(t, r.OK)
	assert.Equal(t, "#btn", mock.lastClickSel)
}

func TestTaskNavigate_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.navigate", TaskNavigate{Window: "main", URL: "https://example.com"})
	require.True(t, r.OK)
	assert.Equal(t, "https://example.com", mock.lastNavURL)
}

func TestTaskScreenshot_Good(t *testing.T) {
	mock := &mockConnector{screenshot: []byte{0x89, 0x50, 0x4E, 0x47}}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.screenshot", TaskScreenshot{Window: "main"})
	require.True(t, r.OK)
	sr, ok := r.Value.(ScreenshotResult)
	assert.True(t, ok)
	assert.Equal(t, "image/png", sr.MimeType)
	assert.NotEmpty(t, sr.Base64)
}

func TestTaskClearConsole_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.clearConsole", TaskClearConsole{Window: "main"})
	require.True(t, r.OK)
	assert.True(t, mock.consoleClearCalled)
}

func TestConnectionCleanup_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Access creates connection
	c.QUERY(QueryURL{Window: "main"})
	assert.False(t, mock.closed)
	// Window close action triggers cleanup
	_ = c.ACTION(window.ActionWindowClosed{Name: "main"})
	assert.True(t, mock.closed)
}

func TestQueryURL_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryURL{Window: "main"})
	assert.False(t, r.OK)
}

// --- SetURL ---

func TestTaskSetURL_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "main", URL: "https://example.com/page"})
	require.True(t, r.OK)
	assert.Equal(t, "https://example.com/page", mock.lastNavURL)
}

func TestTaskSetURL_Bad_UnknownWindow(t *testing.T) {
	_, c := newTestService(t, &mockConnector{})
	// Inject a connector factory that errors
	svc := core.MustServiceFor[*Service](c, "webview")
	svc.newConn = func(_, _ string) (connector, error) {
		return nil, core.E("test", "no connection", nil)
	}
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "bad", URL: "https://example.com"})
	assert.False(t, r.OK)
}

func TestTaskSetURL_Ugly_EmptyURL(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setURL", TaskSetURL{Window: "main", URL: ""})
	require.True(t, r.OK)
	assert.Equal(t, "", mock.lastNavURL)
}

// --- Zoom ---

func TestQueryZoom_Good(t *testing.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryZoom{Window: "main"})
	require.True(t, r.OK)
	assert.InDelta(t, 1.5, r.Value.(float64), 0.001)
}

func TestQueryZoom_Good_DefaultsToOne(t *testing.T) {
	mock := &mockConnector{} // zoom not set → GetZoom returns 1.0
	_, c := newTestService(t, mock)
	r := c.QUERY(QueryZoom{Window: "main"})
	require.True(t, r.OK)
	assert.InDelta(t, 1.0, r.Value.(float64), 0.001)
}

func TestQueryZoom_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryZoom{Window: "main"})
	assert.False(t, r.OK)
}

func TestTaskSetZoom_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setZoom", TaskSetZoom{Window: "main", Zoom: 2.0})
	require.True(t, r.OK)
	assert.InDelta(t, 2.0, mock.lastZoomSet, 0.001)
}

func TestTaskSetZoom_Good_Reset(t *testing.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.setZoom", TaskSetZoom{Window: "main", Zoom: 1.0})
	require.True(t, r.OK)
	assert.InDelta(t, 1.0, mock.zoom, 0.001)
}

func TestTaskSetZoom_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("webview.setZoom").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestTaskSetZoom_Ugly_ZeroZoom(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Zero zoom is technically valid input; the connector accepts it.
	r := taskRun(c, "webview.setZoom", TaskSetZoom{Window: "main", Zoom: 0})
	require.True(t, r.OK)
	assert.InDelta(t, 0.0, mock.lastZoomSet, 0.001)
}

// --- Print ---

func TestTaskPrint_Good_Dialog(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: false})
	require.True(t, r.OK)
	assert.Nil(t, r.Value)
	assert.True(t, mock.printCalled)
	assert.False(t, mock.printToPDF)
}

func TestTaskPrint_Good_PDF(t *testing.T) {
	pdfHeader := []byte{0x25, 0x50, 0x44, 0x46} // %PDF
	mock := &mockConnector{printPDFBytes: pdfHeader}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	require.True(t, r.OK)
	pr, ok := r.Value.(PrintResult)
	require.True(t, ok)
	assert.Equal(t, "application/pdf", pr.MimeType)
	assert.NotEmpty(t, pr.Base64)
	assert.True(t, mock.printToPDF)
}

func TestTaskPrint_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("webview.print").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestTaskPrint_Bad_Error(t *testing.T) {
	mock := &mockConnector{printErr: core.E("test", "print failed", nil)}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	assert.False(t, r.OK)
}

func TestTaskPrint_Ugly_EmptyPDF(t *testing.T) {
	// toPDF=true but connector returns zero bytes — should still wrap as PrintResult
	mock := &mockConnector{printPDFBytes: []byte{}}
	_, c := newTestService(t, mock)
	r := taskRun(c, "webview.print", TaskPrint{Window: "main", ToPDF: true})
	require.True(t, r.OK)
	pr, ok := r.Value.(PrintResult)
	require.True(t, ok)
	assert.Equal(t, "application/pdf", pr.MimeType)
	assert.Equal(t, "", pr.Base64) // empty PDF encodes to empty base64
}

func TestQueryExceptions_Good(t *testing.T) {
	svc, c := newTestService(t, &mockConnector{})
	svc.recordException("main", ExceptionInfo{Text: "boom", Line: 7})

	r := c.QUERY(QueryExceptions{Window: "main"})
	require.True(t, r.OK)
	exceptions := r.Value.([]ExceptionInfo)
	require.Len(t, exceptions, 1)
	assert.Equal(t, "boom", exceptions[0].Text)
	assert.Equal(t, 7, exceptions[0].Line)
}

func TestTaskDevTools_Good(t *testing.T) {
	_, windowPlatform, c := newTestServiceWithWindow(t, &mockConnector{})

	r := taskRun(c, "webview.devtoolsOpen", TaskDevToolsOpen{Window: "main"})
	require.True(t, r.OK)
	assert.True(t, windowPlatform.Windows[0].DevToolsOpen())

	r = taskRun(c, "webview.devtoolsClose", TaskDevToolsClose{Window: "main"})
	require.True(t, r.OK)
	assert.False(t, windowPlatform.Windows[0].DevToolsOpen())
}
