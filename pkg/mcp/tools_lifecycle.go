// pkg/mcp/tools_lifecycle.go
package mcp

import (
	"context"

	"forge.lthn.ai/core/gui/pkg/lifecycle"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- app_quit ---

type AppQuitInput struct{}
type AppQuitOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) appQuit(_ context.Context, _ *mcp.CallToolRequest, _ AppQuitInput) (*mcp.CallToolResult, AppQuitOutput, error) {
	// Broadcast the will-terminate action which triggers application shutdown
	err := s.core.ACTION(lifecycle.ActionWillTerminate{})
	if err != nil {
		return nil, AppQuitOutput{}, err
	}
	return nil, AppQuitOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerLifecycleTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "app_quit", Description: "Quit the application"}, s.appQuit)
}
