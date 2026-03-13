// pkg/mcp/tools_webview.go
package mcp

import (
	"context"

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
	sr, _ := result.(webview.ScreenshotResult)
	return nil, WebviewScreenshotOutput{Base64: sr.Base64, MimeType: sr.MimeType}, nil
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
	msgs, _ := result.([]webview.ConsoleMessage)
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
	el, _ := result.(*webview.ElementInfo)
	return nil, WebviewQueryOutput{Element: el}, nil
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
	els, _ := result.([]*webview.ElementInfo)
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
	html, _ := result.(string)
	return nil, WebviewDOMTreeOutput{HTML: html}, nil
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
	url, _ := result.(string)
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
	title, _ := result.(string)
	return nil, WebviewTitleOutput{Title: title}, nil
}

// --- Registration ---

func (s *Subsystem) registerWebviewTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "webview_eval", Description: "Execute JavaScript in a webview"}, s.webviewEval)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_click", Description: "Click an element in a webview"}, s.webviewClick)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_type", Description: "Type text into an element in a webview"}, s.webviewType)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_navigate", Description: "Navigate a webview to a URL"}, s.webviewNavigate)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_screenshot", Description: "Capture a webview screenshot as base64 PNG"}, s.webviewScreenshot)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_scroll", Description: "Scroll a webview to an absolute position"}, s.webviewScroll)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_hover", Description: "Hover over an element in a webview"}, s.webviewHover)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_select", Description: "Select an option in a select element"}, s.webviewSelect)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_check", Description: "Check or uncheck a checkbox"}, s.webviewCheck)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_upload", Description: "Upload files to a file input element"}, s.webviewUpload)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_viewport", Description: "Set the webview viewport dimensions"}, s.webviewViewport)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_console", Description: "Get captured console messages from a webview"}, s.webviewConsole)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_console_clear", Description: "Clear captured console messages"}, s.webviewConsoleClear)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_query", Description: "Find a single DOM element by CSS selector"}, s.webviewQuery)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_query_all", Description: "Find all DOM elements matching a CSS selector"}, s.webviewQueryAll)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_dom_tree", Description: "Get HTML content of a webview"}, s.webviewDOMTree)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_url", Description: "Get the current URL of a webview"}, s.webviewURL)
	mcp.AddTool(server, &mcp.Tool{Name: "webview_title", Description: "Get the current page title of a webview"}, s.webviewTitle)
}
