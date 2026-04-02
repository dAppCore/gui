// pkg/mcp/tools_tray.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/systray"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- tray_set_icon ---

type TraySetIconInput struct {
	Data []byte `json:"data"`
}
type TraySetIconOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetIcon(_ context.Context, _ *mcp.CallToolRequest, input TraySetIconInput) (*mcp.CallToolResult, TraySetIconOutput, error) {
	_, _, err := s.core.PERFORM(systray.TaskSetTrayIcon{Data: input.Data})
	if err != nil {
		return nil, TraySetIconOutput{}, err
	}
	return nil, TraySetIconOutput{Success: true}, nil
}

// --- tray_set_tooltip ---

type TraySetTooltipInput struct {
	Tooltip string `json:"tooltip"`
}
type TraySetTooltipOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetTooltip(_ context.Context, _ *mcp.CallToolRequest, input TraySetTooltipInput) (*mcp.CallToolResult, TraySetTooltipOutput, error) {
	// Tooltip is set via the tray menu items; for now this is a no-op placeholder
	_ = input.Tooltip
	return nil, TraySetTooltipOutput{Success: true}, nil
}

// --- tray_set_label ---

type TraySetLabelInput struct {
	Label string `json:"label"`
}
type TraySetLabelOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetLabel(_ context.Context, _ *mcp.CallToolRequest, input TraySetLabelInput) (*mcp.CallToolResult, TraySetLabelOutput, error) {
	// Label is part of the tray configuration; placeholder for now
	_ = input.Label
	return nil, TraySetLabelOutput{Success: true}, nil
}

// --- tray_info ---

type TrayInfoInput struct{}
type TrayInfoOutput struct {
	Config map[string]any `json:"config"`
}

func (s *Subsystem) trayInfo(_ context.Context, _ *mcp.CallToolRequest, _ TrayInfoInput) (*mcp.CallToolResult, TrayInfoOutput, error) {
	result, _, err := s.core.QUERY(systray.QueryConfig{})
	if err != nil {
		return nil, TrayInfoOutput{}, err
	}
	config, ok := result.(map[string]any)
	if !ok {
		return nil, TrayInfoOutput{}, fmt.Errorf("unexpected result type from tray config query")
	}
	return nil, TrayInfoOutput{Config: config}, nil
}

// --- tray_show_message ---

type TrayShowMessageInput struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
type TrayShowMessageOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) trayShowMessage(_ context.Context, _ *mcp.CallToolRequest, input TrayShowMessageInput) (*mcp.CallToolResult, TrayShowMessageOutput, error) {
	_, _, err := s.core.PERFORM(systray.TaskShowMessage{Title: input.Title, Message: input.Message})
	if err != nil {
		return nil, TrayShowMessageOutput{}, err
	}
	return nil, TrayShowMessageOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerTrayTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_icon", Description: "Set the system tray icon"}, s.traySetIcon)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_tooltip", Description: "Set the system tray tooltip"}, s.traySetTooltip)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_label", Description: "Set the system tray label"}, s.traySetLabel)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_info", Description: "Get system tray configuration"}, s.trayInfo)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_show_message", Description: "Show a tray message or notification"}, s.trayShowMessage)
}
