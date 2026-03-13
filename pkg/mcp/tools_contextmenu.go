// pkg/mcp/tools_contextmenu.go
package mcp

import (
	"context"

	"forge.lthn.ai/core/gui/pkg/contextmenu"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- contextmenu_add ---

type ContextMenuAddInput struct {
	Name string                    `json:"name"`
	Menu contextmenu.ContextMenuDef `json:"menu"`
}
type ContextMenuAddOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) contextMenuAdd(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuAddInput) (*mcp.CallToolResult, ContextMenuAddOutput, error) {
	_, _, err := s.core.PERFORM(contextmenu.TaskAdd{Name: input.Name, Menu: input.Menu})
	if err != nil {
		return nil, ContextMenuAddOutput{}, err
	}
	return nil, ContextMenuAddOutput{Success: true}, nil
}

// --- contextmenu_remove ---

type ContextMenuRemoveInput struct {
	Name string `json:"name"`
}
type ContextMenuRemoveOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) contextMenuRemove(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuRemoveInput) (*mcp.CallToolResult, ContextMenuRemoveOutput, error) {
	_, _, err := s.core.PERFORM(contextmenu.TaskRemove{Name: input.Name})
	if err != nil {
		return nil, ContextMenuRemoveOutput{}, err
	}
	return nil, ContextMenuRemoveOutput{Success: true}, nil
}

// --- contextmenu_get ---

type ContextMenuGetInput struct {
	Name string `json:"name"`
}
type ContextMenuGetOutput struct {
	Menu *contextmenu.ContextMenuDef `json:"menu"`
}

func (s *Subsystem) contextMenuGet(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuGetInput) (*mcp.CallToolResult, ContextMenuGetOutput, error) {
	result, _, err := s.core.QUERY(contextmenu.QueryGet{Name: input.Name})
	if err != nil {
		return nil, ContextMenuGetOutput{}, err
	}
	menu, _ := result.(*contextmenu.ContextMenuDef)
	return nil, ContextMenuGetOutput{Menu: menu}, nil
}

// --- contextmenu_list ---

type ContextMenuListInput struct{}
type ContextMenuListOutput struct {
	Menus map[string]contextmenu.ContextMenuDef `json:"menus"`
}

func (s *Subsystem) contextMenuList(_ context.Context, _ *mcp.CallToolRequest, _ ContextMenuListInput) (*mcp.CallToolResult, ContextMenuListOutput, error) {
	result, _, err := s.core.QUERY(contextmenu.QueryList{})
	if err != nil {
		return nil, ContextMenuListOutput{}, err
	}
	menus, _ := result.(map[string]contextmenu.ContextMenuDef)
	return nil, ContextMenuListOutput{Menus: menus}, nil
}

// --- Registration ---

func (s *Subsystem) registerContextMenuTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_add", Description: "Register a context menu"}, s.contextMenuAdd)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_remove", Description: "Unregister a context menu"}, s.contextMenuRemove)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_get", Description: "Get a context menu by name"}, s.contextMenuGet)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_list", Description: "List all registered context menus"}, s.contextMenuList)
}
