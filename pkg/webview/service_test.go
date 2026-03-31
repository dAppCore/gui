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

// --- SetURL ---

func TestTaskSetURL_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskSetURL{Window: "main", URL: "https://example.com/page"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "https://example.com/page", mock.lastNavURL)
}

func TestTaskSetURL_Bad_UnknownWindow(t *testing.T) {
	_, c := newTestService(t, &mockConnector{})
	// Inject a connector factory that errors
	svc := core.MustServiceFor[*Service](c, "webview")
	svc.newConn = func(_, _ string) (connector, error) {
		return nil, core.E("test", "no connection", nil)
	}
	_, _, err := c.PERFORM(TaskSetURL{Window: "bad", URL: "https://example.com"})
	assert.Error(t, err)
}

func TestTaskSetURL_Ugly_EmptyURL(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskSetURL{Window: "main", URL: ""})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", mock.lastNavURL)
}

// --- Zoom ---

func TestQueryZoom_Good(t *testing.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	result, handled, err := c.QUERY(QueryZoom{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.InDelta(t, 1.5, result.(float64), 0.001)
}

func TestQueryZoom_Good_DefaultsToOne(t *testing.T) {
	mock := &mockConnector{} // zoom not set → GetZoom returns 1.0
	_, c := newTestService(t, mock)
	result, handled, err := c.QUERY(QueryZoom{Window: "main"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.InDelta(t, 1.0, result.(float64), 0.001)
}

func TestQueryZoom_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryZoom{Window: "main"})
	assert.False(t, handled)
}

func TestTaskSetZoom_Good(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	_, handled, err := c.PERFORM(TaskSetZoom{Window: "main", Zoom: 2.0})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.InDelta(t, 2.0, mock.lastZoomSet, 0.001)
}

func TestTaskSetZoom_Good_Reset(t *testing.T) {
	mock := &mockConnector{zoom: 1.5}
	_, c := newTestService(t, mock)
	_, _, err := c.PERFORM(TaskSetZoom{Window: "main", Zoom: 1.0})
	require.NoError(t, err)
	assert.InDelta(t, 1.0, mock.zoom, 0.001)
}

func TestTaskSetZoom_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSetZoom{Window: "main", Zoom: 1.5})
	assert.False(t, handled)
}

func TestTaskSetZoom_Ugly_ZeroZoom(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	// Zero zoom is technically valid input; the connector accepts it.
	_, handled, err := c.PERFORM(TaskSetZoom{Window: "main", Zoom: 0})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.InDelta(t, 0.0, mock.lastZoomSet, 0.001)
}

// --- Print ---

func TestTaskPrint_Good_Dialog(t *testing.T) {
	mock := &mockConnector{}
	_, c := newTestService(t, mock)
	result, handled, err := c.PERFORM(TaskPrint{Window: "main", ToPDF: false})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
	assert.True(t, mock.printCalled)
	assert.False(t, mock.printToPDF)
}

func TestTaskPrint_Good_PDF(t *testing.T) {
	pdfHeader := []byte{0x25, 0x50, 0x44, 0x46} // %PDF
	mock := &mockConnector{printPDFBytes: pdfHeader}
	_, c := newTestService(t, mock)
	result, handled, err := c.PERFORM(TaskPrint{Window: "main", ToPDF: true})
	require.NoError(t, err)
	assert.True(t, handled)
	pr, ok := result.(PrintResult)
	require.True(t, ok)
	assert.Equal(t, "application/pdf", pr.MimeType)
	assert.NotEmpty(t, pr.Base64)
	assert.True(t, mock.printToPDF)
}

func TestTaskPrint_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskPrint{Window: "main"})
	assert.False(t, handled)
}

func TestTaskPrint_Bad_Error(t *testing.T) {
	mock := &mockConnector{printErr: core.E("test", "print failed", nil)}
	_, c := newTestService(t, mock)
	_, _, err := c.PERFORM(TaskPrint{Window: "main", ToPDF: true})
	assert.Error(t, err)
}

func TestTaskPrint_Ugly_EmptyPDF(t *testing.T) {
	// toPDF=true but connector returns zero bytes — should still wrap as PrintResult
	mock := &mockConnector{printPDFBytes: []byte{}}
	_, c := newTestService(t, mock)
	result, handled, err := c.PERFORM(TaskPrint{Window: "main", ToPDF: true})
	require.NoError(t, err)
	assert.True(t, handled)
	pr, ok := result.(PrintResult)
	require.True(t, ok)
	assert.Equal(t, "application/pdf", pr.MimeType)
	assert.Equal(t, "", pr.Base64) // empty PDF encodes to empty base64
}
