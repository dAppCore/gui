// pkg/webview/service.go
package webview

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

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
	Print() error
	PrintToPDF() ([]byte, error)
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
	opts         Options
	connections  map[string]connector
	exceptions   map[string][]ExceptionInfo
	mu           sync.RWMutex
	newConn      func(debugURL, windowName string) (connector, error) // injectable for tests
	watcherSetup func(conn connector, windowName string)              // called after connection creation
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
		svc := &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, o),
			opts:           o,
			connections:    make(map[string]connector),
			exceptions:     make(map[string][]ExceptionInfo),
			newConn:        defaultNewConn(o),
		}
		svc.watcherSetup = svc.defaultWatcherSetup
		return svc, nil
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

// defaultWatcherSetup wires up console/exception watchers on real connectors.
// It broadcasts ActionConsoleMessage and ActionException via the Core IPC bus.
func (s *Service) defaultWatcherSetup(conn connector, windowName string) {
	rc, ok := conn.(*realConnector)
	if !ok {
		return // test mocks don't need watchers
	}

	cw := gowebview.NewConsoleWatcher(rc.wv)
	cw.AddHandler(func(msg gowebview.ConsoleMessage) {
		_ = s.Core().ACTION(ActionConsoleMessage{
			Window: windowName,
			Message: ConsoleMessage{
				Type:      msg.Type,
				Text:      msg.Text,
				Timestamp: msg.Timestamp,
				URL:       msg.URL,
				Line:      msg.Line,
				Column:    msg.Column,
			},
		})
	})

	ew := gowebview.NewExceptionWatcher(rc.wv)
	ew.AddHandler(func(exc gowebview.ExceptionInfo) {
		_ = s.Core().ACTION(ActionException{
			Window: windowName,
			Exception: ExceptionInfo{
				Text:       exc.Text,
				URL:        exc.URL,
				Line:       exc.LineNumber,
				Column:     exc.ColumnNumber,
				StackTrace: exc.StackTrace,
				Timestamp:  exc.Timestamp,
			},
		})
	})
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
		delete(s.exceptions, m.Name)
		s.mu.Unlock()
	case ActionException:
		s.recordException(m.Window, m.Exception)
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
	if s.watcherSetup != nil {
		s.watcherSetup(conn, windowName)
	}
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
	case QueryComputedStyle:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		result, err := conn.Evaluate(computedStyleScript(q.Selector))
		if err != nil {
			return nil, true, err
		}
		style, err := coerceToMapStringString(result)
		if err != nil {
			return nil, true, err
		}
		return style, true, nil
	case QueryPerformance:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		result, err := conn.Evaluate(performanceScript())
		if err != nil {
			return nil, true, err
		}
		metrics, err := coerceToPerformanceMetrics(result)
		if err != nil {
			return nil, true, err
		}
		return metrics, true, nil
	case QueryResources:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		result, err := conn.Evaluate(resourcesScript())
		if err != nil {
			return nil, true, err
		}
		resources, err := coerceToResourceEntries(result)
		if err != nil {
			return nil, true, err
		}
		return resources, true, nil
	case QueryNetwork:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return nil, true, err
		}
		result, err := conn.Evaluate(networkLogScript(q.Limit))
		if err != nil {
			return nil, true, err
		}
		entries, err := coerceToNetworkEntries(result)
		if err != nil {
			return nil, true, err
		}
		return entries, true, nil
	case QueryExceptions:
		return s.queryExceptions(q.Window, q.Limit), true, nil
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
	case TaskScreenshotElement:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		png, err := captureElementScreenshot(conn, t.Selector)
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
		_, err = conn.Evaluate("window.scrollTo(" + strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y) + ")")
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
	case TaskHighlight:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		_, err = conn.Evaluate(highlightScript(t.Selector, t.Colour))
		return nil, true, err
	case TaskOpenDevTools:
		ws, err := core.ServiceFor[*window.Service](s.Core(), "window")
		if err != nil {
			return nil, true, err
		}
		pw, ok := ws.Manager().Get(t.Window)
		if !ok {
			return nil, true, fmt.Errorf("window not found: %s", t.Window)
		}
		pw.OpenDevTools()
		return nil, true, nil
	case TaskCloseDevTools:
		ws, err := core.ServiceFor[*window.Service](s.Core(), "window")
		if err != nil {
			return nil, true, err
		}
		pw, ok := ws.Manager().Get(t.Window)
		if !ok {
			return nil, true, fmt.Errorf("window not found: %s", t.Window)
		}
		pw.CloseDevTools()
		return nil, true, nil
	case TaskInjectNetworkLogging:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		_, err = conn.Evaluate(networkInitScript())
		return nil, true, err
	case TaskClearNetworkLog:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		_, err = conn.Evaluate(networkClearScript())
		return nil, true, err
	case TaskPrint:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		return nil, true, conn.Print()
	case TaskExportPDF:
		conn, err := s.getConn(t.Window)
		if err != nil {
			return nil, true, err
		}
		pdf, err := conn.PrintToPDF()
		if err != nil {
			return nil, true, err
		}
		return PDFResult{
			Base64:   base64.StdEncoding.EncodeToString(pdf),
			MimeType: "application/pdf",
		}, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) recordException(windowName string, exc ExceptionInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	exceptions := append(s.exceptions[windowName], exc)
	if limit := s.opts.ConsoleLimit; limit > 0 && len(exceptions) > limit {
		exceptions = exceptions[len(exceptions)-limit:]
	}
	s.exceptions[windowName] = exceptions
}

func (s *Service) queryExceptions(windowName string, limit int) []ExceptionInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exceptions := append([]ExceptionInfo(nil), s.exceptions[windowName]...)
	if limit > 0 && len(exceptions) > limit {
		exceptions = exceptions[len(exceptions)-limit:]
	}
	return exceptions
}

func coerceJSON[T any](v any) (T, error) {
	var out T
	raw, err := json.Marshal(v)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, err
	}
	return out, nil
}

func coerceToMapStringString(v any) (map[string]string, error) {
	return coerceJSON[map[string]string](v)
}

func coerceToPerformanceMetrics(v any) (PerformanceMetrics, error) {
	return coerceJSON[PerformanceMetrics](v)
}

func coerceToResourceEntries(v any) ([]ResourceEntry, error) {
	return coerceJSON[[]ResourceEntry](v)
}

func coerceToNetworkEntries(v any) ([]NetworkEntry, error) {
	return coerceJSON[[]NetworkEntry](v)
}

type elementScreenshotBounds struct {
	Left             float64 `json:"left"`
	Top              float64 `json:"top"`
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	DevicePixelRatio float64 `json:"devicePixelRatio"`
}

func elementScreenshotScript(selector string) string {
	sel := jsQuote(selector)
	return fmt.Sprintf(`(function(){
  const el = document.querySelector(%s);
  if (!el) return null;
  try { el.scrollIntoView({block: "center", inline: "center"}); } catch (e) {}
  const rect = el.getBoundingClientRect();
  return {
    left: rect.left,
    top: rect.top,
    width: rect.width,
    height: rect.height,
    devicePixelRatio: window.devicePixelRatio || 1
  };
})()`, sel)
}

func captureElementScreenshot(conn connector, selector string) ([]byte, error) {
	result, err := conn.Evaluate(elementScreenshotScript(selector))
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("webview: element not found: %s", selector)
	}
	bounds, err := coerceJSON[elementScreenshotBounds](result)
	if err != nil {
		return nil, err
	}
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return nil, fmt.Errorf("webview: element has no measurable bounds: %s", selector)
	}
	raw, err := conn.Screenshot()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	scale := bounds.DevicePixelRatio
	if scale <= 0 {
		scale = 1
	}
	left := int(math.Floor(bounds.Left * scale))
	top := int(math.Floor(bounds.Top * scale))
	right := int(math.Ceil((bounds.Left + bounds.Width) * scale))
	bottom := int(math.Ceil((bounds.Top + bounds.Height) * scale))

	srcBounds := img.Bounds()
	if left < srcBounds.Min.X {
		left = srcBounds.Min.X
	}
	if top < srcBounds.Min.Y {
		top = srcBounds.Min.Y
	}
	if right > srcBounds.Max.X {
		right = srcBounds.Max.X
	}
	if bottom > srcBounds.Max.Y {
		bottom = srcBounds.Max.Y
	}
	if right <= left || bottom <= top {
		return nil, fmt.Errorf("webview: element is outside the captured screenshot: %s", selector)
	}

	crop := image.NewRGBA(image.Rect(0, 0, right-left, bottom-top))
	draw.Draw(crop, crop.Bounds(), img, image.Point{X: left, Y: top}, draw.Src)

	var buf bytes.Buffer
	if err := png.Encode(&buf, crop); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// realConnector wraps *gowebview.Webview, converting types at the boundary.
type realConnector struct {
	wv *gowebview.Webview
}

func (r *realConnector) Navigate(url string) error               { return r.wv.Navigate(url) }
func (r *realConnector) Click(sel string) error                  { return r.wv.Click(sel) }
func (r *realConnector) Type(sel, text string) error             { return r.wv.Type(sel, text) }
func (r *realConnector) Evaluate(script string) (any, error)     { return r.wv.Evaluate(script) }
func (r *realConnector) Screenshot() ([]byte, error)             { return r.wv.Screenshot() }
func (r *realConnector) GetURL() (string, error)                 { return r.wv.GetURL() }
func (r *realConnector) GetTitle() (string, error)               { return r.wv.GetTitle() }
func (r *realConnector) GetHTML(sel string) (string, error)      { return r.wv.GetHTML(sel) }
func (r *realConnector) ClearConsole()                           { r.wv.ClearConsole() }
func (r *realConnector) Print() error                            { _, err := r.wv.Evaluate("window.print()"); return err }
func (r *realConnector) Close() error                            { return r.wv.Close() }
func (r *realConnector) SetViewport(w, h int) error              { return r.wv.SetViewport(w, h) }
func (r *realConnector) UploadFile(sel string, p []string) error { return r.wv.UploadFile(sel, p) }
func (r *realConnector) PrintToPDF() ([]byte, error) {
	client, err := r.cdpClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := client.Call(ctx, "Page.printToPDF", map[string]any{
		"printBackground":   true,
		"preferCSSPageSize": true,
	})
	if err != nil {
		return nil, err
	}
	data, ok := result["data"].(string)
	if !ok || data == "" {
		return nil, fmt.Errorf("webview: missing PDF data")
	}
	return base64.StdEncoding.DecodeString(data)
}

func (r *realConnector) cdpClient() (*gowebview.CDPClient, error) {
	rv := reflect.ValueOf(r.wv)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, fmt.Errorf("webview: invalid connector")
	}
	elem := rv.Elem()
	field := elem.FieldByName("client")
	if !field.IsValid() || field.IsNil() {
		return nil, fmt.Errorf("webview: CDP client not available")
	}
	ptr := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
	client, ok := ptr.(*gowebview.CDPClient)
	if !ok || client == nil {
		return nil, fmt.Errorf("webview: unexpected CDP client type")
	}
	return client, nil
}

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
