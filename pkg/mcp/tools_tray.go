// pkg/mcp/tools_tray.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
	r := s.core.Action("systray.setIcon").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: systray.TaskSetTrayIcon{Data: input.Data}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TraySetIconOutput{}, e
		}
		return nil, TraySetIconOutput{}, nil
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
	r := s.core.QUERY(systray.QueryConfig{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, TrayInfoOutput{}, e
		}
		return nil, TrayInfoOutput{}, nil
	}
	config, ok := r.Value.(map[string]any)
	if !ok {
		return nil, TrayInfoOutput{}, coreerr.E("mcp.trayInfo", "unexpected result type", nil)
	}
	return nil, TrayInfoOutput{Config: config}, nil
}

// --- Registration ---

func (s *Subsystem) registerTrayTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "tray_set_icon", Description: "Set the system tray icon"}, s.traySetIcon)
	addTool(s, server, &mcp.Tool{Name: "tray_set_tooltip", Description: "Set the system tray tooltip"}, s.traySetTooltip)
	addTool(s, server, &mcp.Tool{Name: "tray_set_label", Description: "Set the system tray label"}, s.traySetLabel)
	addTool(s, server, &mcp.Tool{Name: "tray_info", Description: "Get system tray configuration"}, s.trayInfo)
}
