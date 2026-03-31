// pkg/mcp/tools_browser.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- browser_open_url ---

type BrowserOpenURLInput struct {
	URL string `json:"url"`
}
type BrowserOpenURLOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) browserOpenURL(_ context.Context, _ *mcp.CallToolRequest, input BrowserOpenURLInput) (*mcp.CallToolResult, BrowserOpenURLOutput, error) {
	r := s.core.Action("browser.openURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: input.URL},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, BrowserOpenURLOutput{}, e
		}
		return nil, BrowserOpenURLOutput{}, nil
	}
	return nil, BrowserOpenURLOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerBrowserTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "browser_open_url", Description: "Open a URL in the default system browser"}, s.browserOpenURL)
}
