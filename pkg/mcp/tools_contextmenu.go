// pkg/mcp/tools_contextmenu.go
package mcp

import (
	"context"
	"encoding/json"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
		return nil, ContextMenuAddOutput{}, coreerr.E("mcp.contextMenuAdd", "failed to marshal menu definition", err)
	}
	var menuDef contextmenu.ContextMenuDef
	if err := json.Unmarshal(menuJSON, &menuDef); err != nil {
		return nil, ContextMenuAddOutput{}, coreerr.E("mcp.contextMenuAdd", "failed to unmarshal menu definition", err)
	}
	r := s.core.Action("contextmenu.add").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: contextmenu.TaskAdd{Name: input.Name, Menu: menuDef}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuAddOutput{}, e
		}
		return nil, ContextMenuAddOutput{}, nil
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
	r := s.core.Action("contextmenu.remove").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: contextmenu.TaskRemove{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuRemoveOutput{}, e
		}
		return nil, ContextMenuRemoveOutput{}, nil
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
	r := s.core.QUERY(contextmenu.QueryGet{Name: input.Name})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuGetOutput{}, e
		}
		return nil, ContextMenuGetOutput{}, nil
	}
	menu, ok := r.Value.(*contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuGetOutput{}, coreerr.E("mcp.contextMenuGet", "unexpected result type", nil)
	}
	if menu == nil {
		return nil, ContextMenuGetOutput{}, nil
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menuJSON, err := json.Marshal(menu)
	if err != nil {
		return nil, ContextMenuGetOutput{}, coreerr.E("mcp.contextMenuGet", "failed to marshal context menu", err)
	}
	var menuMap map[string]any
	if err := json.Unmarshal(menuJSON, &menuMap); err != nil {
		return nil, ContextMenuGetOutput{}, coreerr.E("mcp.contextMenuGet", "failed to unmarshal context menu", err)
	}
	return nil, ContextMenuGetOutput{Menu: menuMap}, nil
}

// --- contextmenu_list ---

type ContextMenuListInput struct{}
type ContextMenuListOutput struct {
	Menus map[string]any `json:"menus"`
}

func (s *Subsystem) contextMenuList(_ context.Context, _ *mcp.CallToolRequest, _ ContextMenuListInput) (*mcp.CallToolResult, ContextMenuListOutput, error) {
	r := s.core.QUERY(contextmenu.QueryList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ContextMenuListOutput{}, e
		}
		return nil, ContextMenuListOutput{}, nil
	}
	menus, ok := r.Value.(map[string]contextmenu.ContextMenuDef)
	if !ok {
		return nil, ContextMenuListOutput{}, coreerr.E("mcp.contextMenuList", "unexpected result type", nil)
	}
	// Convert to map[string]any via JSON round-trip to avoid cyclic type in schema
	menusJSON, err := json.Marshal(menus)
	if err != nil {
		return nil, ContextMenuListOutput{}, coreerr.E("mcp.contextMenuList", "failed to marshal context menus", err)
	}
	var menusMap map[string]any
	if err := json.Unmarshal(menusJSON, &menusMap); err != nil {
		return nil, ContextMenuListOutput{}, coreerr.E("mcp.contextMenuList", "failed to unmarshal context menus", err)
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
