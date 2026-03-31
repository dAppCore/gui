// pkg/mcp/tools_dock.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/dock"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- dock_show ---

type DockShowInput struct{}
type DockShowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockShow(_ context.Context, _ *mcp.CallToolRequest, _ DockShowInput) (*mcp.CallToolResult, DockShowOutput, error) {
	r := s.core.Action("dock.showIcon").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockShowOutput{}, e
		}
		return nil, DockShowOutput{}, nil
	}
	return nil, DockShowOutput{Success: true}, nil
}

// --- dock_hide ---

type DockHideInput struct{}
type DockHideOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) dockHide(_ context.Context, _ *mcp.CallToolRequest, _ DockHideInput) (*mcp.CallToolResult, DockHideOutput, error) {
	r := s.core.Action("dock.hideIcon").Run(context.Background(), core.NewOptions())
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockHideOutput{}, e
		}
		return nil, DockHideOutput{}, nil
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
	r := s.core.Action("dock.setBadge").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: dock.TaskSetBadge{Label: input.Label}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, DockBadgeOutput{}, e
		}
		return nil, DockBadgeOutput{}, nil
	}
	return nil, DockBadgeOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerDockTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "dock_show", Description: "Show the dock/taskbar icon"}, s.dockShow)
	mcp.AddTool(server, &mcp.Tool{Name: "dock_hide", Description: "Hide the dock/taskbar icon"}, s.dockHide)
	mcp.AddTool(server, &mcp.Tool{Name: "dock_badge", Description: "Set the dock/taskbar badge label"}, s.dockBadge)
}
