// pkg/mcp/tools_webview.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
	r := s.core.Action("webview.evaluate").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskEvaluate{Window: input.Window, Script: input.Script}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewEvalOutput{}, e
		}
		return nil, WebviewEvalOutput{}, nil
	}
	return nil, WebviewEvalOutput{Result: r.Value, Window: input.Window}, nil
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
	r := s.core.Action("webview.click").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskClick{Window: input.Window, Selector: input.Selector}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewClickOutput{}, e
		}
		return nil, WebviewClickOutput{}, nil
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
	r := s.core.Action("webview.type").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskType{Window: input.Window, Selector: input.Selector, Text: input.Text}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewTypeOutput{}, e
		}
		return nil, WebviewTypeOutput{}, nil
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
	r := s.core.Action("webview.navigate").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskNavigate{Window: input.Window, URL: input.URL}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewNavigateOutput{}, e
		}
		return nil, WebviewNavigateOutput{}, nil
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
	r := s.core.Action("webview.screenshot").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskScreenshot{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScreenshotOutput{}, e
		}
		return nil, WebviewScreenshotOutput{}, nil
	}
	sr, ok := r.Value.(webview.ScreenshotResult)
	if !ok {
		return nil, WebviewScreenshotOutput{}, coreerr.E("mcp.webviewScreenshot", "unexpected result type", nil)
	}
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
	r := s.core.Action("webview.scroll").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskScroll{Window: input.Window, X: input.X, Y: input.Y}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewScrollOutput{}, e
		}
		return nil, WebviewScrollOutput{}, nil
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
	r := s.core.Action("webview.hover").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskHover{Window: input.Window, Selector: input.Selector}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewHoverOutput{}, e
		}
		return nil, WebviewHoverOutput{}, nil
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
	r := s.core.Action("webview.select").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskSelect{Window: input.Window, Selector: input.Selector, Value: input.Value}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewSelectOutput{}, e
		}
		return nil, WebviewSelectOutput{}, nil
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
	r := s.core.Action("webview.check").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskCheck{Window: input.Window, Selector: input.Selector, Checked: input.Checked}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewCheckOutput{}, e
		}
		return nil, WebviewCheckOutput{}, nil
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
	r := s.core.Action("webview.uploadFile").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskUploadFile{Window: input.Window, Selector: input.Selector, Paths: input.Paths}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewUploadOutput{}, e
		}
		return nil, WebviewUploadOutput{}, nil
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
	r := s.core.Action("webview.setViewport").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskSetViewport{Window: input.Window, Width: input.Width, Height: input.Height}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewViewportOutput{}, e
		}
		return nil, WebviewViewportOutput{}, nil
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
	r := s.core.QUERY(webview.QueryConsole{Window: input.Window, Level: input.Level, Limit: input.Limit})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewConsoleOutput{}, e
		}
		return nil, WebviewConsoleOutput{}, nil
	}
	msgs, ok := r.Value.([]webview.ConsoleMessage)
	if !ok {
		return nil, WebviewConsoleOutput{}, coreerr.E("mcp.webviewConsole", "unexpected result type", nil)
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
	r := s.core.Action("webview.clearConsole").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: webview.TaskClearConsole{Window: input.Window}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewConsoleClearOutput{}, e
		}
		return nil, WebviewConsoleClearOutput{}, nil
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
	r := s.core.QUERY(webview.QuerySelector{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewQueryOutput{}, e
		}
		return nil, WebviewQueryOutput{}, nil
	}
	el, ok := r.Value.(*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryOutput{}, coreerr.E("mcp.webviewQuery", "unexpected result type", nil)
	}
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
	r := s.core.QUERY(webview.QuerySelectorAll{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewQueryAllOutput{}, e
		}
		return nil, WebviewQueryAllOutput{}, nil
	}
	els, ok := r.Value.([]*webview.ElementInfo)
	if !ok {
		return nil, WebviewQueryAllOutput{}, coreerr.E("mcp.webviewQueryAll", "unexpected result type", nil)
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
	r := s.core.QUERY(webview.QueryDOMTree{Window: input.Window, Selector: input.Selector})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewDOMTreeOutput{}, e
		}
		return nil, WebviewDOMTreeOutput{}, nil
	}
	html, ok := r.Value.(string)
	if !ok {
		return nil, WebviewDOMTreeOutput{}, coreerr.E("mcp.webviewDOMTree", "unexpected result type", nil)
	}
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
	r := s.core.QUERY(webview.QueryURL{Window: input.Window})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewURLOutput{}, e
		}
		return nil, WebviewURLOutput{}, nil
	}
	url, ok := r.Value.(string)
	if !ok {
		return nil, WebviewURLOutput{}, coreerr.E("mcp.webviewURL", "unexpected result type", nil)
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
	r := s.core.QUERY(webview.QueryTitle{Window: input.Window})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WebviewTitleOutput{}, e
		}
		return nil, WebviewTitleOutput{}, nil
	}
	title, ok := r.Value.(string)
	if !ok {
		return nil, WebviewTitleOutput{}, coreerr.E("mcp.webviewTitle", "unexpected result type", nil)
	}
	return nil, WebviewTitleOutput{Title: title}, nil
}

// --- Registration ---

func (s *Subsystem) registerWebviewTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "webview_eval", Description: "Execute JavaScript in a webview"}, s.webviewEval)
	addTool(s, server, &mcp.Tool{Name: "webview_click", Description: "Click an element in a webview"}, s.webviewClick)
	addTool(s, server, &mcp.Tool{Name: "webview_type", Description: "Type text into an element in a webview"}, s.webviewType)
	addTool(s, server, &mcp.Tool{Name: "webview_navigate", Description: "Navigate a webview to a URL"}, s.webviewNavigate)
	addTool(s, server, &mcp.Tool{Name: "webview_screenshot", Description: "Capture a webview screenshot as base64 PNG"}, s.webviewScreenshot)
	addTool(s, server, &mcp.Tool{Name: "webview_scroll", Description: "Scroll a webview to an absolute position"}, s.webviewScroll)
	addTool(s, server, &mcp.Tool{Name: "webview_hover", Description: "Hover over an element in a webview"}, s.webviewHover)
	addTool(s, server, &mcp.Tool{Name: "webview_select", Description: "Select an option in a select element"}, s.webviewSelect)
	addTool(s, server, &mcp.Tool{Name: "webview_check", Description: "Check or uncheck a checkbox"}, s.webviewCheck)
	addTool(s, server, &mcp.Tool{Name: "webview_upload", Description: "Upload files to a file input element"}, s.webviewUpload)
	addTool(s, server, &mcp.Tool{Name: "webview_viewport", Description: "Set the webview viewport dimensions"}, s.webviewViewport)
	addTool(s, server, &mcp.Tool{Name: "webview_console", Description: "Get captured console messages from a webview"}, s.webviewConsole)
	addTool(s, server, &mcp.Tool{Name: "webview_console_clear", Description: "Clear captured console messages"}, s.webviewConsoleClear)
	addTool(s, server, &mcp.Tool{Name: "webview_query", Description: "Find a single DOM element by CSS selector"}, s.webviewQuery)
	addTool(s, server, &mcp.Tool{Name: "webview_query_all", Description: "Find all DOM elements matching a CSS selector"}, s.webviewQueryAll)
	addTool(s, server, &mcp.Tool{Name: "webview_dom_tree", Description: "Get HTML content of a webview"}, s.webviewDOMTree)
	addTool(s, server, &mcp.Tool{Name: "webview_url", Description: "Get the current URL of a webview"}, s.webviewURL)
	addTool(s, server, &mcp.Tool{Name: "webview_title", Description: "Get the current page title of a webview"}, s.webviewTitle)
}
