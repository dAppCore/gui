// pkg/mcp/tools_dock.go
package mcp

import (
	"context"

	"forge.lthn.ai/core/gui/pkg/dock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- dock_show ---

type DockShowInput struct{}
type DockShowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockShow(_ context.Context, _ *mcp.CallToolRequest, _ DockShowInput) (*mcp.CallToolResult, DockShowOutput, error) {
	_, _, err := s.core.PERFORM(dock.TaskShowIcon{})
	if err != nil {
		return nil, DockShowOutput{}, err
	}
	return nil, DockShowOutput{Success: true}, nil
}

// --- dock_hide ---

type DockHideInput struct{}
type DockHideOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockHide(_ context.Context, _ *mcp.CallToolRequest, _ DockHideInput) (*mcp.CallToolResult, DockHideOutput, error) {
	_, _, err := s.core.PERFORM(dock.TaskHideIcon{})
	if err != nil {
		return nil, DockHideOutput{}, err
	}
	return nil, DockHideOutput{Success: true}, nil
}

// --- dock_badge ---

type DockBadgeInput struct {
	Label string `json:"label"`
}
type DockBadgeOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockBadge(_ context.Context, _ *mcp.CallToolRequest, input DockBadgeInput) (*mcp.CallToolResult, DockBadgeOutput, error) {
	_, _, err := s.core.PERFORM(dock.TaskSetBadge{Label: input.Label})
	if err != nil {
		return nil, DockBadgeOutput{}, err
	}
	return nil, DockBadgeOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerDockTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "dock_show", Description: "Show the dock/taskbar icon"}, s.dockShow)
	mcp.AddTool(server, &mcp.Tool{Name: "dock_hide", Description: "Hide the dock/taskbar icon"}, s.dockHide)
	mcp.AddTool(server, &mcp.Tool{Name: "dock_badge", Description: "Set the dock/taskbar badge label"}, s.dockBadge)
}
