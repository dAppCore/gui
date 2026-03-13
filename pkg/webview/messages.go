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
// Field mapping from go-webview: LineNumber->Line, ColumnNumber->Column.
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
