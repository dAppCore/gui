// pkg/webview/messages.go
package webview

import "time"

// --- Queries (read-only) ---

// QueryURL gets the current page URL. Result: string
// Use: result, _, err := c.QUERY(webview.QueryURL{Window: "editor"})
type QueryURL struct {
	Window string `json:"window"`
}

// QueryTitle gets the current page title. Result: string
// Use: result, _, err := c.QUERY(webview.QueryTitle{Window: "editor"})
type QueryTitle struct {
	Window string `json:"window"`
}

// QueryConsole gets captured console messages. Result: []ConsoleMessage
// Use: result, _, err := c.QUERY(webview.QueryConsole{Window: "editor", Level: "error", Limit: 20})
type QueryConsole struct {
	Window string `json:"window"`
	Level  string `json:"level,omitempty"` // filter by type: "log", "warn", "error", "info", "debug"
	Limit  int    `json:"limit,omitempty"` // max messages (0 = all)
}

// QuerySelector finds a single element. Result: *ElementInfo (nil if not found)
// Use: result, _, err := c.QUERY(webview.QuerySelector{Window: "editor", Selector: "#submit"})
type QuerySelector struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// QuerySelectorAll finds all matching elements. Result: []*ElementInfo
// Use: result, _, err := c.QUERY(webview.QuerySelectorAll{Window: "editor", Selector: "button"})
type QuerySelectorAll struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// QueryDOMTree gets HTML content. Result: string (outerHTML)
// Use: result, _, err := c.QUERY(webview.QueryDOMTree{Window: "editor", Selector: "main"})
type QueryDOMTree struct {
	Window   string `json:"window"`
	Selector string `json:"selector,omitempty"` // empty = full document
}

// QueryComputedStyle returns the computed CSS properties for an element.
// Use: result, _, err := c.QUERY(webview.QueryComputedStyle{Window: "editor", Selector: "#panel"})
type QueryComputedStyle struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// QueryPerformance returns page performance metrics.
// Use: result, _, err := c.QUERY(webview.QueryPerformance{Window: "editor"})
type QueryPerformance struct {
	Window string `json:"window"`
}

// QueryResources returns the page's loaded resource entries.
// Use: result, _, err := c.QUERY(webview.QueryResources{Window: "editor"})
type QueryResources struct {
	Window string `json:"window"`
}

// QueryNetwork returns the captured network log.
// Use: result, _, err := c.QUERY(webview.QueryNetwork{Window: "editor", Limit: 50})
type QueryNetwork struct {
	Window string `json:"window"`
	Limit  int    `json:"limit,omitempty"`
}

// QueryExceptions returns captured JavaScript exceptions.
// Use: result, _, err := c.QUERY(webview.QueryExceptions{Window: "editor", Limit: 10})
type QueryExceptions struct {
	Window string `json:"window"`
	Limit  int    `json:"limit,omitempty"`
}

// --- Tasks (side-effects) ---

// TaskEvaluate executes JavaScript. Result: any (JS return value)
// Use: _, _, err := c.PERFORM(webview.TaskEvaluate{Window: "editor", Script: "document.title"})
type TaskEvaluate struct {
	Window string `json:"window"`
	Script string `json:"script"`
}

// TaskClick clicks an element. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskClick{Window: "editor", Selector: "#submit"})
type TaskClick struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// TaskType types text into an element. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskType{Window: "editor", Selector: "#search", Text: "core"})
type TaskType struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Text     string `json:"text"`
}

// TaskNavigate navigates to a URL. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskNavigate{Window: "editor", URL: "https://example.com"})
type TaskNavigate struct {
	Window string `json:"window"`
	URL    string `json:"url"`
}

// TaskScreenshot captures the page as PNG. Result: ScreenshotResult
// Use: result, _, err := c.PERFORM(webview.TaskScreenshot{Window: "editor"})
type TaskScreenshot struct {
	Window string `json:"window"`
}

// TaskScreenshotElement captures a specific element as PNG. Result: ScreenshotResult
// Use: result, _, err := c.PERFORM(webview.TaskScreenshotElement{Window: "editor", Selector: "#panel"})
type TaskScreenshotElement struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// TaskScroll scrolls to an absolute position (window.scrollTo). Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskScroll{Window: "editor", X: 0, Y: 600})
type TaskScroll struct {
	Window string `json:"window"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

// TaskHover hovers over an element. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskHover{Window: "editor", Selector: "#help"})
type TaskHover struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

// TaskSelect selects an option in a <select> element. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskSelect{Window: "editor", Selector: "#theme", Value: "dark"})
type TaskSelect struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

// TaskCheck checks/unchecks a checkbox. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskCheck{Window: "editor", Selector: "#accept", Checked: true})
type TaskCheck struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Checked  bool   `json:"checked"`
}

// TaskUploadFile uploads files to an input element. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskUploadFile{Window: "editor", Selector: "input[type=file]", Paths: []string{"/tmp/report.pdf"}})
type TaskUploadFile struct {
	Window   string   `json:"window"`
	Selector string   `json:"selector"`
	Paths    []string `json:"paths"`
}

// TaskSetViewport sets the viewport dimensions. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskSetViewport{Window: "editor", Width: 1280, Height: 800})
type TaskSetViewport struct {
	Window string `json:"window"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// TaskClearConsole clears captured console messages. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskClearConsole{Window: "editor"})
type TaskClearConsole struct {
	Window string `json:"window"`
}

// TaskHighlight visually highlights an element.
// Use: _, _, err := c.PERFORM(webview.TaskHighlight{Window: "editor", Selector: "#submit", Colour: "#ffcc00"})
type TaskHighlight struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Colour   string `json:"colour,omitempty"`
}

// TaskOpenDevTools opens the browser devtools for the target window. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskOpenDevTools{Window: "editor"})
type TaskOpenDevTools struct {
	Window string `json:"window"`
}

// TaskCloseDevTools closes the browser devtools for the target window. Result: nil
// Use: _, _, err := c.PERFORM(webview.TaskCloseDevTools{Window: "editor"})
type TaskCloseDevTools struct {
	Window string `json:"window"`
}

// TaskInjectNetworkLogging injects fetch/XHR interception into the page.
// Use: _, _, err := c.PERFORM(webview.TaskInjectNetworkLogging{Window: "editor"})
type TaskInjectNetworkLogging struct {
	Window string `json:"window"`
}

// TaskClearNetworkLog clears the captured network log.
// Use: _, _, err := c.PERFORM(webview.TaskClearNetworkLog{Window: "editor"})
type TaskClearNetworkLog struct {
	Window string `json:"window"`
}

// TaskPrint prints the current page using the browser's native print flow.
// Use: _, _, err := c.PERFORM(webview.TaskPrint{Window: "editor"})
type TaskPrint struct {
	Window string `json:"window"`
}

// TaskExportPDF exports the page to a PDF document.
// Use: result, _, err := c.PERFORM(webview.TaskExportPDF{Window: "editor"})
type TaskExportPDF struct {
	Window string `json:"window"`
}

// --- Actions (broadcast) ---

// ActionConsoleMessage is broadcast when a console message is captured.
// Use: _ = c.ACTION(webview.ActionConsoleMessage{Window: "editor", Message: webview.ConsoleMessage{Type: "error", Text: "boom"}})
type ActionConsoleMessage struct {
	Window  string         `json:"window"`
	Message ConsoleMessage `json:"message"`
}

// ActionException is broadcast when a JavaScript exception occurs.
// Use: _ = c.ACTION(webview.ActionException{Window: "editor", Exception: webview.ExceptionInfo{Text: "ReferenceError"}})
type ActionException struct {
	Window    string        `json:"window"`
	Exception ExceptionInfo `json:"exception"`
}

// --- Types ---

// ConsoleMessage represents a browser console message.
// Use: msg := webview.ConsoleMessage{Type: "warn", Text: "slow network"}
type ConsoleMessage struct {
	Type      string    `json:"type"` // "log", "warn", "error", "info", "debug"
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url,omitempty"`
	Line      int       `json:"line,omitempty"`
	Column    int       `json:"column,omitempty"`
}

// ElementInfo represents a DOM element.
// Use: el := webview.ElementInfo{TagName: "button", InnerText: "Save"}
type ElementInfo struct {
	TagName     string            `json:"tagName"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	InnerText   string            `json:"innerText,omitempty"`
	InnerHTML   string            `json:"innerHTML,omitempty"`
	BoundingBox *BoundingBox      `json:"boundingBox,omitempty"`
}

// BoundingBox represents element position and size.
// Use: box := webview.BoundingBox{X: 10, Y: 20, Width: 120, Height: 40}
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ExceptionInfo represents a JavaScript exception.
// Field mapping from go-webview: LineNumber->Line, ColumnNumber->Column.
// Use: err := webview.ExceptionInfo{Text: "ReferenceError", URL: "app://editor"}
type ExceptionInfo struct {
	Text       string    `json:"text"`
	URL        string    `json:"url,omitempty"`
	Line       int       `json:"line,omitempty"`
	Column     int       `json:"column,omitempty"`
	StackTrace string    `json:"stackTrace,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ScreenshotResult wraps raw PNG bytes as base64 for JSON/MCP transport.
// Use: shot := webview.ScreenshotResult{Base64: "iVBORw0KGgo=", MimeType: "image/png"}
type ScreenshotResult struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"` // always "image/png"
}

// PerformanceMetrics summarises browser performance timings.
// Use: metrics := webview.PerformanceMetrics{NavigationStart: 1.2, LoadEventEnd: 42.5}
type PerformanceMetrics struct {
	NavigationStart      float64 `json:"navigationStart"`
	DOMContentLoaded     float64 `json:"domContentLoaded"`
	LoadEventEnd         float64 `json:"loadEventEnd"`
	FirstPaint           float64 `json:"firstPaint,omitempty"`
	FirstContentfulPaint float64 `json:"firstContentfulPaint,omitempty"`
	UsedJSHeapSize       float64 `json:"usedJSHeapSize,omitempty"`
	TotalJSHeapSize      float64 `json:"totalJSHeapSize,omitempty"`
}

// ResourceEntry summarises a loaded resource.
// Use: entry := webview.ResourceEntry{Name: "app.js", EntryType: "resource"}
type ResourceEntry struct {
	Name            string  `json:"name"`
	EntryType       string  `json:"entryType"`
	InitiatorType   string  `json:"initiatorType,omitempty"`
	StartTime       float64 `json:"startTime"`
	Duration        float64 `json:"duration"`
	TransferSize    float64 `json:"transferSize,omitempty"`
	EncodedBodySize float64 `json:"encodedBodySize,omitempty"`
	DecodedBodySize float64 `json:"decodedBodySize,omitempty"`
}

// NetworkEntry summarises a captured fetch/XHR request.
// Use: entry := webview.NetworkEntry{URL: "/api/status", Method: "GET", Status: 200}
type NetworkEntry struct {
	URL       string `json:"url"`
	Method    string `json:"method"`
	Status    int    `json:"status,omitempty"`
	Resource  string `json:"resource,omitempty"`
	OK        bool   `json:"ok,omitempty"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// PDFResult contains exported PDF bytes encoded for transport.
// Use: pdf := webview.PDFResult{Base64: "JVBERi0xLjQ=", MimeType: "application/pdf"}
type PDFResult struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"` // always "application/pdf"
}
