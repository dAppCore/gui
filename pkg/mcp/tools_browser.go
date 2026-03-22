// pkg/mcp/tools_browser.go
package mcp

import (
	"context"

	"forge.lthn.ai/core/gui/pkg/browser"
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
	_, _, err := s.core.PERFORM(browser.TaskOpenURL{URL: input.URL})
	if err != nil {
		return nil, BrowserOpenURLOutput{}, err
	}
	return nil, BrowserOpenURLOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerBrowserTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "browser_open_url", Description: "Open a URL in the default system browser"}, s.browserOpenURL)
}
