// pkg/mcp/tools_webview.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/webview"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- webview_eval ---

type WebviewEvalInput struct {
	Window string `json:"window"`
	Script string `json:"script"`
}

type WebviewEvalOutput struct {
	Result any    `json:"result"`
	Window string `json:"window"`
}

func (s *Subsystem) webviewEval(_ context.Context, _ *mcp.CallToolRequest, input WebviewEvalInput) (*mcp.CallToolResult, WebviewEvalOutput, error) {
	result, _, err := s.core.PERFORM(webview.TaskEvaluate{Window: input.Window, Script: input.Script})
	if err != nil {
		return nil, WebviewEvalOutput{}, err
	}
	return nil, WebviewEvalOutput{Result: result, Window: input.Window}, nil
}

// --- webview_click ---

type WebviewClickInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewClickOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewClick(_ context.Context, _ *mcp.CallToolRequest, input WebviewClickInput) (*mcp.CallToolResult, WebviewClickOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskClick{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewClickOutput{}, err
	}
	return nil, WebviewClickOutput{Success: true}, nil
}

// --- webview_type ---

type WebviewTypeInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Text     string `json:"text"`
}

type WebviewTypeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewType(_ context.Context, _ *mcp.CallToolRequest, input WebviewTypeInput) (*mcp.CallToolResult, WebviewTypeOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskType{Window: input.Window, Selector: input.Selector, Text: input.Text})
	if err != nil {
		return nil, WebviewTypeOutput{}, err
	}
	return nil, WebviewTypeOutput{Success: true}, nil
}

// --- webview_navigate ---

type WebviewNavigateInput struct {
	Window string `json:"window"`
	URL    string `json:"url"`
}

type WebviewNavigateOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNavigate(_ context.Context, _ *mcp.CallToolRequest, input WebviewNavigateInput) (*mcp.CallToolResult, WebviewNavigateOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskNavigate{Window: input.Window, URL: input.URL})
	if err != nil {
		return nil, WebviewNavigateOutput{}, err
	}
	return nil, WebviewNavigateOutput{Success: true}, nil
}

// --- webview_screenshot ---

type WebviewScreenshotInput struct {
	Window string `json:"window"`
}

type WebviewScreenshotOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewScreenshot(_ context.Context, _ *mcp.CallToolRequest, input WebviewScreenshotInput) (*mcp.CallToolResult, WebviewScreenshotOutput, error) {
	result, _, err := s.core.PERFORM(webview.TaskScreenshot{Window: input.Window})
	if err != nil {
		return nil, WebviewScreenshotOutput{}, err
	}
	sr, ok := result.(webview.ScreenshotResult)
	if !ok {
		return nil, WebviewScreenshotOutput{}, fmt.Errorf("unexpected result type from webview screenshot")
	}
	return nil, WebviewScreenshotOutput{Base64: sr.Base64, MimeType: sr.MimeType}, nil
}

// --- webview_screenshot_element ---

type WebviewScreenshotElementInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewScreenshotElementOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewScreenshotElement(_ context.Context, _ *mcp.CallToolRequest, input WebviewScreenshotElementInput) (*mcp.CallToolResult, WebviewScreenshotElementOutput, error) {
	result, _, err := s.core.PERFORM(webview.TaskScreenshotElement{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewScreenshotElementOutput{}, err
	}
	sr, ok := result.(webview.ScreenshotResult)
	if !ok {
		return nil, WebviewScreenshotElementOutput{}, fmt.Errorf("unexpected result type from webview element screenshot")
	}
	return nil, WebviewScreenshotElementOutput{Base64: sr.Base64, MimeType: sr.MimeType}, nil
}

// --- webview_scroll ---

type WebviewScrollInput struct {
	Window string `json:"window"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type WebviewScrollOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewScroll(_ context.Context, _ *mcp.CallToolRequest, input WebviewScrollInput) (*mcp.CallToolResult, WebviewScrollOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskScroll{Window: input.Window, X: input.X, Y: input.Y})
	if err != nil {
		return nil, WebviewScrollOutput{}, err
	}
	return nil, WebviewScrollOutput{Success: true}, nil
}

// --- webview_hover ---

type WebviewHoverInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewHoverOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewHover(_ context.Context, _ *mcp.CallToolRequest, input WebviewHoverInput) (*mcp.CallToolResult, WebviewHoverOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskHover{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewHoverOutput{}, err
	}
	return nil, WebviewHoverOutput{Success: true}, nil
}

// --- webview_select ---

type WebviewSelectInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type WebviewSelectOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewSelect(_ context.Context, _ *mcp.CallToolRequest, input WebviewSelectInput) (*mcp.CallToolResult, WebviewSelectOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskSelect{Window: input.Window, Selector: input.Selector, Value: input.Value})
	if err != nil {
		return nil, WebviewSelectOutput{}, err
	}
	return nil, WebviewSelectOutput{Success: true}, nil
}

// --- webview_check ---

type WebviewCheckInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Checked  bool   `json:"checked"`
}

type WebviewCheckOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewCheck(_ context.Context, _ *mcp.CallToolRequest, input WebviewCheckInput) (*mcp.CallToolResult, WebviewCheckOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskCheck{Window: input.Window, Selector: input.Selector, Checked: input.Checked})
	if err != nil {
		return nil, WebviewCheckOutput{}, err
	}
	return nil, WebviewCheckOutput{Success: true}, nil
}

// --- webview_upload ---

type WebviewUploadInput struct {
	Window   string   `json:"window"`
	Selector string   `json:"selector"`
	Paths    []string `json:"paths"`
}

type WebviewUploadOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewUpload(_ context.Context, _ *mcp.CallToolRequest, input WebviewUploadInput) (*mcp.CallToolResult, WebviewUploadOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskUploadFile{Window: input.Window, Selector: input.Selector, Paths: input.Paths})
	if err != nil {
		return nil, WebviewUploadOutput{}, err
	}
	return nil, WebviewUploadOutput{Success: true}, nil
}

// --- webview_viewport ---

type WebviewViewportInput struct {
	Window string `json:"window"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type WebviewViewportOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewViewport(_ context.Context, _ *mcp.CallToolRequest, input WebviewViewportInput) (*mcp.CallToolResult, WebviewViewportOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskSetViewport{Window: input.Window, Width: input.Width, Height: input.Height})
	if err != nil {
		return nil, WebviewViewportOutput{}, err
	}
	return nil, WebviewViewportOutput{Success: true}, nil
}

// --- webview_console ---

type WebviewConsoleInput struct {
	Window string `json:"window"`
	Level  string `json:"level,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type WebviewConsoleOutput struct {
	Messages []webview.ConsoleMessage `json:"messages"`
}

func (s *Subsystem) webviewConsole(_ context.Context, _ *mcp.CallToolRequest, input WebviewConsoleInput) (*mcp.CallToolResult, WebviewConsoleOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryConsole{Window: input.Window, Level: input.Level, Limit: input.Limit})
	if err != nil {
		return nil, WebviewConsoleOutput{}, err
	}
	msgs, ok := result.([]webview.ConsoleMessage)
	if !ok {
		return nil, WebviewConsoleOutput{}, fmt.Errorf("unexpected result type from webview console query")
	}
	return nil, WebviewConsoleOutput{Messages: msgs}, nil
}

// --- webview_console_clear ---

type WebviewConsoleClearInput struct {
	Window string `json:"window"`
}

type WebviewConsoleClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewConsoleClear(_ context.Context, _ *mcp.CallToolRequest, input WebviewConsoleClearInput) (*mcp.CallToolResult, WebviewConsoleClearOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskClearConsole{Window: input.Window})
	if err != nil {
		return nil, WebviewConsoleClearOutput{}, err
	}
	return nil, WebviewConsoleClearOutput{Success: true}, nil
}

// --- webview_devtools_open ---

type WebviewDevToolsOpenInput struct {
	Window string `json:"window"`
}
type WebviewDevToolsOpenOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewDevToolsOpen(_ context.Context, _ *mcp.CallToolRequest, input WebviewDevToolsOpenInput) (*mcp.CallToolResult, WebviewDevToolsOpenOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskOpenDevTools{Window: input.Window})
	if err != nil {
		return nil, WebviewDevToolsOpenOutput{}, err
	}
	return nil, WebviewDevToolsOpenOutput{Success: true}, nil
}

// --- webview_devtools_close ---

type WebviewDevToolsCloseInput struct {
	Window string `json:"window"`
}
type WebviewDevToolsCloseOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewDevToolsClose(_ context.Context, _ *mcp.CallToolRequest, input WebviewDevToolsCloseInput) (*mcp.CallToolResult, WebviewDevToolsCloseOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskCloseDevTools{Window: input.Window})
	if err != nil {
		return nil, WebviewDevToolsCloseOutput{}, err
	}
	return nil, WebviewDevToolsCloseOutput{Success: true}, nil
}

// --- webview_query ---

type WebviewQueryInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewQueryOutput struct {
	Element *webview.ElementInfo `json:"element"`
}

func (s *Subsystem) webviewQuery(_ context.Context, _ *mcp.CallToolRequest, input WebviewQueryInput) (*mcp.CallToolResult, WebviewQueryOutput, error) {
	result, _, err := s.core.QUERY(webview.QuerySelector{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewQueryOutput{}, err
	}
	el, ok := result.(*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryOutput{}, fmt.Errorf("unexpected result type from webview query")
	}
	return nil, WebviewQueryOutput{Element: el}, nil
}

// --- webview_element_info ---

func (s *Subsystem) webviewElementInfo(_ context.Context, _ *mcp.CallToolRequest, input WebviewQueryInput) (*mcp.CallToolResult, WebviewQueryOutput, error) {
	return s.webviewQuery(nil, nil, input)
}

// --- webview_query_all ---

type WebviewQueryAllInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewQueryAllOutput struct {
	Elements []*webview.ElementInfo `json:"elements"`
}

func (s *Subsystem) webviewQueryAll(_ context.Context, _ *mcp.CallToolRequest, input WebviewQueryAllInput) (*mcp.CallToolResult, WebviewQueryAllOutput, error) {
	result, _, err := s.core.QUERY(webview.QuerySelectorAll{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewQueryAllOutput{}, err
	}
	els, ok := result.([]*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryAllOutput{}, fmt.Errorf("unexpected result type from webview query all")
	}
	return nil, WebviewQueryAllOutput{Elements: els}, nil
}

// --- webview_dom_tree ---

type WebviewDOMTreeInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector,omitempty"`
}

type WebviewDOMTreeOutput struct {
	HTML string `json:"html"`
}

func (s *Subsystem) webviewDOMTree(_ context.Context, _ *mcp.CallToolRequest, input WebviewDOMTreeInput) (*mcp.CallToolResult, WebviewDOMTreeOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryDOMTree{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewDOMTreeOutput{}, err
	}
	html, ok := result.(string)
	if !ok {
		return nil, WebviewDOMTreeOutput{}, fmt.Errorf("unexpected result type from webview DOM tree query")
	}
	return nil, WebviewDOMTreeOutput{HTML: html}, nil
}

// --- webview_source ---

func (s *Subsystem) webviewSource(_ context.Context, _ *mcp.CallToolRequest, input WebviewDOMTreeInput) (*mcp.CallToolResult, WebviewDOMTreeOutput, error) {
	return s.webviewDOMTree(nil, nil, input)
}

// --- webview_computed_style ---

type WebviewComputedStyleInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
}

type WebviewComputedStyleOutput struct {
	Style map[string]string `json:"style"`
}

func (s *Subsystem) webviewComputedStyle(_ context.Context, _ *mcp.CallToolRequest, input WebviewComputedStyleInput) (*mcp.CallToolResult, WebviewComputedStyleOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryComputedStyle{Window: input.Window, Selector: input.Selector})
	if err != nil {
		return nil, WebviewComputedStyleOutput{}, err
	}
	style, ok := result.(map[string]string)
	if !ok {
		return nil, WebviewComputedStyleOutput{}, fmt.Errorf("unexpected result type from webview computed style query")
	}
	return nil, WebviewComputedStyleOutput{Style: style}, nil
}

// --- webview_performance ---

type WebviewPerformanceInput struct {
	Window string `json:"window"`
}

type WebviewPerformanceOutput struct {
	Metrics webview.PerformanceMetrics `json:"metrics"`
}

func (s *Subsystem) webviewPerformance(_ context.Context, _ *mcp.CallToolRequest, input WebviewPerformanceInput) (*mcp.CallToolResult, WebviewPerformanceOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryPerformance{Window: input.Window})
	if err != nil {
		return nil, WebviewPerformanceOutput{}, err
	}
	metrics, ok := result.(webview.PerformanceMetrics)
	if !ok {
		return nil, WebviewPerformanceOutput{}, fmt.Errorf("unexpected result type from webview performance query")
	}
	return nil, WebviewPerformanceOutput{Metrics: metrics}, nil
}

// --- webview_resources ---

type WebviewResourcesInput struct {
	Window string `json:"window"`
}

type WebviewResourcesOutput struct {
	Resources []webview.ResourceEntry `json:"resources"`
}

func (s *Subsystem) webviewResources(_ context.Context, _ *mcp.CallToolRequest, input WebviewResourcesInput) (*mcp.CallToolResult, WebviewResourcesOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryResources{Window: input.Window})
	if err != nil {
		return nil, WebviewResourcesOutput{}, err
	}
	resources, ok := result.([]webview.ResourceEntry)
	if !ok {
		return nil, WebviewResourcesOutput{}, fmt.Errorf("unexpected result type from webview resources query")
	}
	return nil, WebviewResourcesOutput{Resources: resources}, nil
}

// --- webview_network ---

type WebviewNetworkInput struct {
	Window string `json:"window"`
	Limit  int    `json:"limit,omitempty"`
}

type WebviewNetworkOutput struct {
	Requests []webview.NetworkEntry `json:"requests"`
}

func (s *Subsystem) webviewNetwork(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkInput) (*mcp.CallToolResult, WebviewNetworkOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryNetwork{Window: input.Window, Limit: input.Limit})
	if err != nil {
		return nil, WebviewNetworkOutput{}, err
	}
	requests, ok := result.([]webview.NetworkEntry)
	if !ok {
		return nil, WebviewNetworkOutput{}, fmt.Errorf("unexpected result type from webview network query")
	}
	return nil, WebviewNetworkOutput{Requests: requests}, nil
}

// --- webview_network_inject ---

type WebviewNetworkInjectInput struct {
	Window string `json:"window"`
}

type WebviewNetworkInjectOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNetworkInject(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkInjectInput) (*mcp.CallToolResult, WebviewNetworkInjectOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskInjectNetworkLogging{Window: input.Window})
	if err != nil {
		return nil, WebviewNetworkInjectOutput{}, err
	}
	return nil, WebviewNetworkInjectOutput{Success: true}, nil
}

// --- webview_network_clear ---

type WebviewNetworkClearInput struct {
	Window string `json:"window"`
}

type WebviewNetworkClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewNetworkClear(_ context.Context, _ *mcp.CallToolRequest, input WebviewNetworkClearInput) (*mcp.CallToolResult, WebviewNetworkClearOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskClearNetworkLog{Window: input.Window})
	if err != nil {
		return nil, WebviewNetworkClearOutput{}, err
	}
	return nil, WebviewNetworkClearOutput{Success: true}, nil
}

// --- webview_highlight ---

type WebviewHighlightInput struct {
	Window   string `json:"window"`
	Selector string `json:"selector"`
	Colour   string `json:"colour,omitempty"`
}

type WebviewHighlightOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewHighlight(_ context.Context, _ *mcp.CallToolRequest, input WebviewHighlightInput) (*mcp.CallToolResult, WebviewHighlightOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskHighlight{Window: input.Window, Selector: input.Selector, Colour: input.Colour})
	if err != nil {
		return nil, WebviewHighlightOutput{}, err
	}
	return nil, WebviewHighlightOutput{Success: true}, nil
}

// --- webview_print ---

type WebviewPrintInput struct {
	Window string `json:"window"`
}

type WebviewPrintOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) webviewPrint(_ context.Context, _ *mcp.CallToolRequest, input WebviewPrintInput) (*mcp.CallToolResult, WebviewPrintOutput, error) {
	_, _, err := s.core.PERFORM(webview.TaskPrint{Window: input.Window})
	if err != nil {
		return nil, WebviewPrintOutput{}, err
	}
	return nil, WebviewPrintOutput{Success: true}, nil
}

// --- webview_pdf ---

type WebviewPDFInput struct {
	Window string `json:"window"`
}

type WebviewPDFOutput struct {
	Base64   string `json:"base64"`
	MimeType string `json:"mimeType"`
}

func (s *Subsystem) webviewPDF(_ context.Context, _ *mcp.CallToolRequest, input WebviewPDFInput) (*mcp.CallToolResult, WebviewPDFOutput, error) {
	result, _, err := s.core.PERFORM(webview.TaskExportPDF{Window: input.Window})
	if err != nil {
		return nil, WebviewPDFOutput{}, err
	}
	pdf, ok := result.(webview.PDFResult)
	if !ok {
		return nil, WebviewPDFOutput{}, fmt.Errorf("unexpected result type from webview pdf task")
	}
	return nil, WebviewPDFOutput{Base64: pdf.Base64, MimeType: pdf.MimeType}, nil
}

// --- webview_url ---

type WebviewURLInput struct {
	Window string `json:"window"`
}

type WebviewURLOutput struct {
	URL string `json:"url"`
}

func (s *Subsystem) webviewURL(_ context.Context, _ *mcp.CallToolRequest, input WebviewURLInput) (*mcp.CallToolResult, WebviewURLOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryURL{Window: input.Window})
	if err != nil {
		return nil, WebviewURLOutput{}, err
	}
	url, ok := result.(string)
	if !ok {
		return nil, WebviewURLOutput{}, fmt.Errorf("unexpected result type from webview URL query")
	}
	return nil, WebviewURLOutput{URL: url}, nil
}

// --- webview_title ---

type WebviewTitleInput struct {
	Window string `json:"window"`
}

type WebviewTitleOutput struct {
	Title string `json:"title"`
}

func (s *Subsystem) webviewTitle(_ context.Context, _ *mcp.CallToolRequest, input WebviewTitleInput) (*mcp.CallToolResult, WebviewTitleOutput, error) {
	result, _, err := s.core.QUERY(webview.QueryTitle{Window: input.Window})
	if err != nil {
		return nil, WebviewTitleOutput{}, err
	}
	title, ok := result.(string)
	if !ok {
		return nil, WebviewTitleOutput{}, fmt.Errorf("unexpected result type from webview title query")
	}
	return nil, WebviewTitleOutput{Title: title}, nil
}

// --- Registration ---

func (s *Subsystem) registerWebviewTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "webview_eval", Description: "Execute JavaScript in a webview"}, s.webviewEval)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_click", Description: "Click an element in a webview"}, s.webviewClick)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_type", Description: "Type text into an element in a webview"}, s.webviewType)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_navigate", Description: "Navigate a webview to a URL"}, s.webviewNavigate)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_screenshot", Description: "Capture a webview screenshot as base64 PNG"}, s.webviewScreenshot)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_screenshot_element", Description: "Capture a specific element as base64 PNG"}, s.webviewScreenshotElement)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_scroll", Description: "Scroll a webview to an absolute position"}, s.webviewScroll)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_hover", Description: "Hover over an element in a webview"}, s.webviewHover)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_select", Description: "Select an option in a select element"}, s.webviewSelect)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_check", Description: "Check or uncheck a checkbox"}, s.webviewCheck)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_upload", Description: "Upload files to a file input element"}, s.webviewUpload)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_viewport", Description: "Set the webview viewport dimensions"}, s.webviewViewport)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_console", Description: "Get captured console messages from a webview"}, s.webviewConsole)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_console_clear", Description: "Clear captured console messages"}, s.webviewConsoleClear)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_clear_console", Description: "Alias for webview_console_clear"}, s.webviewConsoleClear)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_query", Description: "Find a single DOM element by CSS selector"}, s.webviewQuery)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_element_info", Description: "Get detailed information about a DOM element"}, s.webviewElementInfo)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_query_all", Description: "Find all DOM elements matching a CSS selector"}, s.webviewQueryAll)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_dom_tree", Description: "Get HTML content of a webview"}, s.webviewDOMTree)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_source", Description: "Get page HTML source"}, s.webviewSource)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_computed_style", Description: "Get computed styles for an element"}, s.webviewComputedStyle)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_performance", Description: "Get page performance metrics"}, s.webviewPerformance)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_resources", Description: "List loaded page resources"}, s.webviewResources)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_network", Description: "Get captured network requests"}, s.webviewNetwork)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_network_inject", Description: "Inject fetch/XHR network logging"}, s.webviewNetworkInject)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_network_clear", Description: "Clear captured network requests"}, s.webviewNetworkClear)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_highlight", Description: "Visually highlight an element"}, s.webviewHighlight)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_print", Description: "Open the browser print dialog"}, s.webviewPrint)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_pdf", Description: "Export the current page as a PDF"}, s.webviewPDF)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_url", Description: "Get the current URL of a webview"}, s.webviewURL)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_title", Description: "Get the current page title of a webview"}, s.webviewTitle)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_devtools_open", Description: "Open devtools for a webview window"}, s.webviewDevToolsOpen)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_devtools_close", Description: "Close devtools for a webview window"}, s.webviewDevToolsClose)
}
