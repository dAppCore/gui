// pkg/mcp/tools_tray.go
package mcp

import (
	"context"

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
	_, _, err := s.core.PERFORM(systray.TaskSetTrayTooltip{Tooltip: input.Tooltip})
	if err != nil {
		return nil, TraySetTooltipOutput{}, err
	}
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
	_, _, err := s.core.PERFORM(systray.TaskSetTrayLabel{Label: input.Label})
	if err != nil {
		return nil, TraySetLabelOutput{}, err
	}
	return nil, TraySetLabelOutput{Success: true}, nil
}

// --- tray_set_menu ---

type TraySetMenuInput struct {
	Items []map[string]any `json:"items"`
}
type TraySetMenuOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) traySetMenu(_ context.Context, _ *mcp.CallToolRequest, input TraySetMenuInput) (*mcp.CallToolResult, TraySetMenuOutput, error) {
	items := make([]systray.TrayMenuItem, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, decodeTrayMenuItem(item))
	}
	_, _, err := s.core.PERFORM(systray.TaskSetTrayMenu{Items: items})
	if err != nil {
		return nil, TraySetMenuOutput{}, err
	}
	return nil, TraySetMenuOutput{Success: true}, nil
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
		config = map[string]any{}
	}
	return nil, TrayInfoOutput{Config: config}, nil
}

// --- Registration ---

func (s *Subsystem) registerTrayTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_icon", Description: "Set the system tray icon"}, s.traySetIcon)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_tooltip", Description: "Set the system tray tooltip"}, s.traySetTooltip)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_label", Description: "Set the system tray label"}, s.traySetLabel)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_set_menu", Description: "Set the system tray menu"}, s.traySetMenu)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_show_message", Description: "Show a tray balloon or tray message"}, s.trayShowMessage)
	mcp.AddTool(server, &mcp.Tool{Name: "tray_info", Description: "Get system tray configuration"}, s.trayInfo)
}

func decodeTrayMenuItem(input map[string]any) systray.TrayMenuItem {
	item := systray.TrayMenuItem{}
	if label, ok := input["label"].(string); ok {
		item.Label = label
	}
	if itemType, ok := input["type"].(string); ok {
		item.Type = itemType
	}
	if checked, ok := input["checked"].(bool); ok {
		item.Checked = checked
	}
	if disabled, ok := input["disabled"].(bool); ok {
		item.Disabled = disabled
	}
	if tooltip, ok := input["tooltip"].(string); ok {
		item.Tooltip = tooltip
	}
	if actionID, ok := input["actionId"].(string); ok {
		item.ActionID = actionID
	}
	if actionID, ok := input["action_id"].(string); ok && item.ActionID == "" {
		item.ActionID = actionID
	}
	if rawSubmenu, ok := input["submenu"].([]any); ok {
		item.Submenu = make([]systray.TrayMenuItem, 0, len(rawSubmenu))
		for _, child := range rawSubmenu {
			if childMap, ok := child.(map[string]any); ok {
				item.Submenu = append(item.Submenu, decodeTrayMenuItem(childMap))
			}
		}
	}
	return item
}
