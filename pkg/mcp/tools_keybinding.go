// pkg/mcp/tools_keybinding.go
package mcp

import (
	"context"

	"forge.lthn.ai/core/gui/pkg/keybinding"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- keybinding_add ---

type KeybindingAddInput struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description,omitempty"`
}
type KeybindingAddOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) keybindingAdd(_ context.Context, _ *mcp.CallToolRequest, input KeybindingAddInput) (*mcp.CallToolResult, KeybindingAddOutput, error) {
	_, _, err := s.core.PERFORM(keybinding.TaskAdd{Accelerator: input.Accelerator, Description: input.Description})
	if err != nil {
		return nil, KeybindingAddOutput{}, err
	}
	return nil, KeybindingAddOutput{Success: true}, nil
}

// --- keybinding_remove ---

type KeybindingRemoveInput struct {
	Accelerator string `json:"accelerator"`
}
type KeybindingRemoveOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) keybindingRemove(_ context.Context, _ *mcp.CallToolRequest, input KeybindingRemoveInput) (*mcp.CallToolResult, KeybindingRemoveOutput, error) {
	_, _, err := s.core.PERFORM(keybinding.TaskRemove{Accelerator: input.Accelerator})
	if err != nil {
		return nil, KeybindingRemoveOutput{}, err
	}
	return nil, KeybindingRemoveOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerKeybindingTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "keybinding_add", Description: "Register a keyboard shortcut"}, s.keybindingAdd)
	mcp.AddTool(server, &mcp.Tool{Name: "keybinding_remove", Description: "Unregister a keyboard shortcut"}, s.keybindingRemove)
}
