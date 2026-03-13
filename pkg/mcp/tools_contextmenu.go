// pkg/mcp/tools_contextmenu.go
package mcp

import (
	"context"
	"encoding/json"
	"fmt"

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
	menuJSON, err := json.Marshal(input.Menu)
	if err != nil {
		return nil, ContextMenuAddOutput{}, fmt.Errorf("failed to marshal menu definition: %w", err)
	}
	var menuDef contextmenu.ContextMenuDef
	if err := json.Unmarshal(menuJSON, &menuDef); err != nil {
		return nil, ContextMenuAddOutput{}, fmt.Errorf("failed to unmarshal menu definition: %w", err)
	}
	_, _, err = s.core.PERFORM(contextmenu.TaskAdd{Name: input.Name, Menu: menuDef})
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
	menu, ok := result.(*contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuGetOutput{}, fmt.Errorf("unexpected result type from context menu get query")
	}
	if menu == nil {
		return nil, ContextMenuGetOutput{}, nil
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menuJSON, err := json.Marshal(menu)
	if err != nil {
		return nil, ContextMenuGetOutput{}, fmt.Errorf("failed to marshal context menu: %w", err)
	}
	var menuMap map[string]any
	if err := json.Unmarshal(menuJSON, &menuMap); err != nil {
		return nil, ContextMenuGetOutput{}, fmt.Errorf("failed to unmarshal context menu: %w", err)
	}
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
	menus, ok := result.(map[string]contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuListOutput{}, fmt.Errorf("unexpected result type from context menu list query")
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menusJSON, err := json.Marshal(menus)
	if err != nil {
		return nil, ContextMenuListOutput{}, fmt.Errorf("failed to marshal context menus: %w", err)
	}
	var menusMap map[string]any
	if err := json.Unmarshal(menusJSON, &menusMap); err != nil {
		return nil, ContextMenuListOutput{}, fmt.Errorf("failed to unmarshal context menus: %w", err)
	}
	return nil, ContextMenuListOutput{Menus: menusMap}, nil
}

// --- Registration ---

func (s *Subsystem) registerContextMenuTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_add", Description: "Register a context menu"}, s.contextMenuAdd)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_remove", Description: "Unregister a context menu"}, s.contextMenuRemove)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_get", Description: "Get a context menu by name"}, s.contextMenuGet)
	mcp.AddTool(server, &mcp.Tool{Name: "contextmenu_list", Description: "List all registered context menus"}, s.contextMenuList)
}
