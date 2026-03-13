// pkg/mcp/tools_contextmenu.go
package mcp

import (
	"context"
	"encoding/json"

	"forge.lthn.ai/core/gui/pkg/contextmenu"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- contextmenu_add ---

// ContextMenuAddInput uses map[string]any for the menu definition because
// contextmenu.ContextMenuDef contains self-referencing MenuItemDef (Items []MenuItemDef)
// which the MCP SDK schema generator cannot handle (cycle detection panic).
type ContextMenuAddInput struct {
	Name string         `json:"name"`
	Menu map[string]any `json:"menu"`
}
type ContextMenuAddOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) contextMenuAdd(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuAddInput) (*mcp.CallToolResult, ContextMenuAddOutput, error) {
	// Convert map[string]any to ContextMenuDef via JSON round-trip
	menuJSON, _ := json.Marshal(input.Menu)
	var menuDef contextmenu.ContextMenuDef
	_ = json.Unmarshal(menuJSON, &menuDef)
	_, _, err := s.core.PERFORM(contextmenu.TaskAdd{Name: input.Name, Menu: menuDef})
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
	Menu map[string]any `json:"menu"`
}

func (s *Subsystem) contextMenuGet(_ context.Context, _ *mcp.CallToolRequest, input ContextMenuGetInput) (*mcp.CallToolResult, ContextMenuGetOutput, error) {
	result, _, err := s.core.QUERY(contextmenu.QueryGet{Name: input.Name})
	if err != nil {
		return nil, ContextMenuGetOutput{}, err
	}
	menu, _ := result.(*contextmenu.ContextMenuDef)
	if menu == nil {
		return nil, ContextMenuGetOutput{}, nil
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menuJSON, _ := json.Marshal(menu)
	var menuMap map[string]any
	_ = json.Unmarshal(menuJSON, &menuMap)
	return nil, ContextMenuGetOutput{Menu: menuMap}, nil
}

// --- contextmenu_list ---

type ContextMenuListInput struct{}
type ContextMenuListOutput struct {
	Menus map[string]any `json:"menus"`
}

func (s *Subsystem) contextMenuList(_ context.Context, _ *mcp.CallToolRequest, _ ContextMenuListInput) (*mcp.CallToolResult, ContextMenuListOutput, error) {
	result, _, err := s.core.QUERY(contextmenu.QueryList{})
	if err != nil {
		return nil, ContextMenuListOutput{}, err
	}
	menus, _ := result.(map[string]contextmenu.ContextMenuDef)
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menusJSON, _ := json.Marshal(menus)
	var menusMap map[string]any
	_ = json.Unmarshal(menusJSON, &menusMap)
	return nil, ContextMenuListOutput{Menus: menusMap}, nil
}

// --- Registration ---

func (s *Subsystem) registerContextMenuTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_add", Description: "Register a context menu"}, s.contextMenuAdd)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_remove", Description: "Unregister a context menu"}, s.contextMenuRemove)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_get", Description: "Get a context menu by name"}, s.contextMenuGet)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_list", Description: "List all registered context menus"}, s.contextMenuList)
}
