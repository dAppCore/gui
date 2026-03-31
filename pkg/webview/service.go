// pkg/webview/service.go
package webview

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	gowebview "forge.lthn.ai/core/go-webview"
	core "dappco.re/go/core"
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
	GetZoom() (float64, error)
	SetZoom(zoom float64) error
	Print(toPDF bool) ([]byte, error)
	Close() error
}

type Options struct {
	DebugURL     string        // Chrome debug endpoint (default: "http://localhost:9222")
	Timeout      time.Duration // Operation timeout (default: 30s)
	ConsoleLimit int           // Max console messages per window (default: 1000)
}

type Service struct {
	*core.ServiceRuntime[Options]
	options      Options
	connections  map[string]connector
	mu           sync.RWMutex
	newConn      func(debugURL, windowName string) (connector, error) // injectable for tests
	watcherSetup func(conn connector, windowName string)              // called after connection creation
}

// Register binds the webview service to a Core instance.
// core.WithService(webview.Register())
// core.WithService(webview.Register(func(o *Options) { o.DebugURL = "http://localhost:9223" }))
func Register(optionFns ...func(*Options)) func(*core.Core) core.Result {
	o := Options{
		DebugURL:     "http://localhost:9222",
		Timeout:      30 * time.Second,
		ConsoleLimit: 1000,
	}
	for _, fn := range optionFns {
		fn(&o)
	}
	return func(c *core.Core) core.Result {
		svc := &Service{
			ServiceRuntime: core.NewServiceRuntime[Options](c, o),
			options:        o,
			connections:    make(map[string]connector),
			newConn:        defaultNewConn(o),
		}
		svc.watcherSetup = svc.defaultWatcherSetup
		return core.Result{Value: svc, OK: true}
	}
}

// defaultNewConn creates real go-webview connections.
func defaultNewConn(options Options) func(string, string) (connector, error) {
	return func(debugURL, windowName string) (connector, error) {
		// Enumerate targets, match by title/URL containing window name
		targets, err := gowebview.ListTargets(debugURL)
		if err != nil {
			return nil, err
		}
		var wsURL string
		for _, t := range targets {
			if t.Type == "page" && (bytes.Contains([]byte(t.Title), []byte(windowName)) || bytes.Contains([]byte(t.URL), []byte(windowName))) {
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
			gowebview.WithTimeout(options.Timeout),
			gowebview.WithConsoleLimit(options.ConsoleLimit),
		)
		if err != nil {
			return nil, err
		}
		return &realConnector{wv: wv, debugURL: debugURL}, nil
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

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().RegisterQuery(s.handleQuery)
	s.registerTaskActions()
	return core.Result{OK: true}
}

// OnShutdown closes all CDP connections.
func (s *Service) OnShutdown(_ context.Context) core.Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, conn := range s.connections {
		conn.Close()
		delete(s.connections, name)
	}
	return core.Result{OK: true}
}

// HandleIPCEvents listens for window close events to clean up connections.
func (s *Service) HandleIPCEvents(_ *core.Core, msg core.Message) core.Result {
	switch m := msg.(type) {
	case window.ActionWindowClosed:
		s.mu.Lock()
		if conn, ok := s.connections[m.Name]; ok {
			conn.Close()
			delete(s.connections, m.Name)
		}
		s.mu.Unlock()
	}
	return core.Result{OK: true}
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
	conn, err := s.newConn(s.options.DebugURL, windowName)
	if err != nil {
		return nil, err
	}
	s.connections[windowName] = conn
	if s.watcherSetup != nil {
		s.watcherSetup(conn, windowName)
	}
	return conn, nil
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch q := q.(type) {
	case QueryURL:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		url, err := conn.GetURL()
		return core.Result{}.New(url, err)
	case QueryTitle:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		title, err := conn.GetTitle()
		return core.Result{}.New(title, err)
	case QueryConsole:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
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
		return core.Result{Value: msgs, OK: true}
	case QuerySelector:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		el, err := conn.QuerySelector(q.Selector)
		return core.Result{}.New(el, err)
	case QuerySelectorAll:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		els, err := conn.QuerySelectorAll(q.Selector)
		return core.Result{}.New(els, err)
	case QueryDOMTree:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		selector := q.Selector
		if selector == "" {
			selector = "html"
		}
		html, err := conn.GetHTML(selector)
		return core.Result{}.New(html, err)
	case QueryZoom:
		conn, err := s.getConn(q.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		zoom, err := conn.GetZoom()
		return core.Result{}.New(zoom, err)
	default:
		return core.Result{}
	}
}

// registerTaskActions registers all webview task handlers as named Core actions.
func (s *Service) registerTaskActions() {
	c := s.Core()
	c.Action("webview.evaluate", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskEvaluate)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		result, err := conn.Evaluate(t.Script)
		return core.Result{}.New(result, err)
	})
	c.Action("webview.click", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClick)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Click(t.Selector))
	})
	c.Action("webview.type", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskType)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Type(t.Selector, t.Text))
	})
	c.Action("webview.navigate", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskNavigate)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Navigate(t.URL))
	})
	c.Action("webview.screenshot", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskScreenshot)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		png, err := conn.Screenshot()
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: ScreenshotResult{
			Base64:   base64.StdEncoding.EncodeToString(png),
			MimeType: "image/png",
		}, OK: true}
	})
	c.Action("webview.scroll", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskScroll)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		_, err = conn.Evaluate("window.scrollTo(" + strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y) + ")")
		return core.Result{Value: nil, OK: true}.New(err)
	})
	c.Action("webview.hover", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskHover)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Hover(t.Selector))
	})
	c.Action("webview.select", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSelect)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Select(t.Selector, t.Value))
	})
	c.Action("webview.check", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskCheck)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Check(t.Selector, t.Checked))
	})
	c.Action("webview.uploadFile", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskUploadFile)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.UploadFile(t.Selector, t.Paths))
	})
	c.Action("webview.setViewport", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetViewport)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.SetViewport(t.Width, t.Height))
	})
	c.Action("webview.clearConsole", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskClearConsole)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conn.ClearConsole()
		return core.Result{OK: true}
	})
	c.Action("webview.setURL", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetURL)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.Navigate(t.URL))
	})
	c.Action("webview.setZoom", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskSetZoom)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		return core.Result{Value: nil, OK: true}.New(conn.SetZoom(t.Zoom))
	})
	c.Action("webview.print", func(_ context.Context, opts core.Options) core.Result {
		t, _ := opts.Get("task").Value.(TaskPrint)
		conn, err := s.getConn(t.Window)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		pdfBytes, err := conn.Print(t.ToPDF)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if !t.ToPDF {
			return core.Result{OK: true}
		}
		return core.Result{Value: PrintResult{
			Base64:   base64.StdEncoding.EncodeToString(pdfBytes),
			MimeType: "application/pdf",
		}, OK: true}
	})
}

// realConnector wraps *gowebview.Webview, converting types at the boundary.
// debugURL is retained so that PDF printing can issue a Page.printToPDF CDP call
// via a fresh CDPClient, since go-webview v0.1.7 does not expose a PrintToPDF helper.
type realConnector struct {
	wv       *gowebview.Webview
	debugURL string // Chrome debug HTTP endpoint (e.g., http://localhost:9222) for direct CDP calls
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
func (r *realConnector) Close() error                            { return r.wv.Close() }
func (r *realConnector) SetViewport(w, h int) error              { return r.wv.SetViewport(w, h) }
func (r *realConnector) UploadFile(sel string, p []string) error { return r.wv.UploadFile(sel, p) }

// GetZoom returns the current CSS zoom level as a float64.
// zoom, _ := conn.GetZoom()  // 1.0 = 100%, 1.5 = 150%
func (r *realConnector) GetZoom() (float64, error) {
	raw, err := r.wv.Evaluate("parseFloat(document.documentElement.style.zoom) || 1.0")
	if err != nil {
		return 0, core.E("realConnector.GetZoom", "failed to get zoom", err)
	}
	switch v := raw.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 1.0, nil
	}
}

// SetZoom sets the CSS zoom level on the document root element.
// conn.SetZoom(1.5)  // 150%
// conn.SetZoom(1.0)  // reset to normal
func (r *realConnector) SetZoom(zoom float64) error {
	script := "document.documentElement.style.zoom = '" + strconv.FormatFloat(zoom, 'g', -1, 64) + "'; undefined"
	_, err := r.wv.Evaluate(script)
	if err != nil {
		return core.E("realConnector.SetZoom", "failed to set zoom", err)
	}
	return nil
}

// Print triggers window.print() or exports to PDF via Page.printToPDF.
// When toPDF is false the browser print dialog is opened (via window.print()) and nil bytes are returned.
// When toPDF is true a fresh CDPClient is opened against the stored WebSocket URL to issue
// Page.printToPDF, which returns raw PDF bytes.
func (r *realConnector) Print(toPDF bool) ([]byte, error) {
	if !toPDF {
		_, err := r.wv.Evaluate("window.print(); undefined")
		if err != nil {
			return nil, core.E("realConnector.Print", "failed to open print dialog", err)
		}
		return nil, nil
	}

	if r.debugURL == "" {
		return nil, core.E("realConnector.Print", "no debug URL stored; cannot issue Page.printToPDF", nil)
	}

	// Open a dedicated CDPClient for the single Page.printToPDF call.
	// NewCDPClient connects to the first page target at the debug endpoint.
	client, err := gowebview.NewCDPClient(r.debugURL)
	if err != nil {
		return nil, core.E("realConnector.Print", "failed to connect for PDF export", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := client.Call(ctx, "Page.printToPDF", map[string]any{
		"printBackground": true,
	})
	if err != nil {
		return nil, core.E("realConnector.Print", "Page.printToPDF failed", err)
	}

	dataStr, ok := result["data"].(string)
	if !ok {
		return nil, core.E("realConnector.Print", "Page.printToPDF returned no data", nil)
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil, core.E("realConnector.Print", "failed to decode PDF data", err)
	}

	return pdfBytes, nil
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
