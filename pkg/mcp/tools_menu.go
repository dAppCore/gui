package mcp

import (
	"context"
	"strings"

	"dappco.re/go/core/gui/pkg/menu"
	coreerr "forge.lthn.ai/core/go-log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MenuItemSpec struct {
	Label       string         `json:"label,omitempty"`
	Accelerator string         `json:"accelerator,omitempty"`
	Type        string         `json:"type,omitempty"`
	Checked     bool           `json:"checked,omitempty"`
	Disabled    bool           `json:"disabled,omitempty"`
	Tooltip     string         `json:"tooltip,omitempty"`
	Children    []MenuItemSpec `json:"children,omitempty"`
	Role        string         `json:"role,omitempty"`
}

type MenuGetInput struct{}

type MenuOutput struct {
	Items []MenuItemSpec `json:"items"`
}

type MenuSetInput struct {
	Items []MenuItemSpec `json:"items"`
}

func (s *Subsystem) menuGet(_ context.Context, _ *mcp.CallToolRequest, _ MenuGetInput) (*mcp.CallToolResult, MenuOutput, error) {
	items, err := s.queryMenuItems()
	if err != nil {
		return nil, MenuOutput{}, err
	}
	return nil, MenuOutput{Items: items}, nil
}

func (s *Subsystem) menuSet(_ context.Context, _ *mcp.CallToolRequest, input MenuSetInput) (*mcp.CallToolResult, MenuOutput, error) {
	items, err := decodeMenuItems(input.Items)
	if err != nil {
		return nil, MenuOutput{}, err
	}
	if _, _, err := s.core.PERFORM(menu.TaskSetAppMenu{Items: items}); err != nil {
		return nil, MenuOutput{}, err
	}
	snapshot, err := s.queryMenuItems()
	if err != nil {
		return nil, MenuOutput{}, err
	}
	return nil, MenuOutput{Items: snapshot}, nil
}

func (s *Subsystem) queryMenuItems() ([]MenuItemSpec, error) {
	result, _, err := s.core.QUERY(menu.QueryGetAppMenu{})
	if err != nil {
		return nil, err
	}
	items, ok := result.([]menu.MenuItem)
	if !ok {
		return nil, coreerr.E("mcp.menuGet", "unexpected result type", nil)
	}
	return encodeMenuItems(items), nil
}

func encodeMenuItems(items []menu.MenuItem) []MenuItemSpec {
	if len(items) == 0 {
		return nil
	}
	out := make([]MenuItemSpec, 0, len(items))
	for _, item := range items {
		spec := MenuItemSpec{
			Label:       item.Label,
			Accelerator: item.Accelerator,
			Type:        item.Type,
			Checked:     item.Checked,
			Disabled:    item.Disabled,
			Tooltip:     item.Tooltip,
			Children:    encodeMenuItems(item.Children),
		}
		if item.Role != nil {
			spec.Role = encodeMenuRole(*item.Role)
		}
		out = append(out, spec)
	}
	return out
}

func decodeMenuItems(items []MenuItemSpec) ([]menu.MenuItem, error) {
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]menu.MenuItem, 0, len(items))
	for _, item := range items {
		role, err := decodeMenuRole(item.Role)
		if err != nil {
			return nil, err
		}
		children, err := decodeMenuItems(item.Children)
		if err != nil {
			return nil, err
		}
		out = append(out, menu.MenuItem{
			Label:       item.Label,
			Accelerator: item.Accelerator,
			Type:        item.Type,
			Checked:     item.Checked,
			Disabled:    item.Disabled,
			Tooltip:     item.Tooltip,
			Children:    children,
			Role:        role,
		})
	}
	return out, nil
}

func encodeMenuRole(role menu.MenuRole) string {
	switch role {
	case menu.RoleAppMenu:
		return "app"
	case menu.RoleFileMenu:
		return "file"
	case menu.RoleEditMenu:
		return "edit"
	case menu.RoleViewMenu:
		return "view"
	case menu.RoleWindowMenu:
		return "window"
	case menu.RoleHelpMenu:
		return "help"
	default:
		return ""
	}
}

func decodeMenuRole(role string) (*menu.MenuRole, error) {
	switch strings.TrimSpace(strings.ToLower(role)) {
	case "":
		return nil, nil
	case "app":
		value := menu.RoleAppMenu
		return &value, nil
	case "file":
		value := menu.RoleFileMenu
		return &value, nil
	case "edit":
		value := menu.RoleEditMenu
		return &value, nil
	case "view":
		value := menu.RoleViewMenu
		return &value, nil
	case "window":
		value := menu.RoleWindowMenu
		return &value, nil
	case "help":
		value := menu.RoleHelpMenu
		return &value, nil
	default:
		return nil, coreerr.E("mcp.decodeMenuRole", "unknown menu role: "+role, nil)
	}
}

func (s *Subsystem) registerMenuTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "menu_get", Description: "Get the current application menu structure"}, s.menuGet)
	mcp.AddTool(server, &mcp.Tool{Name: "menu_set", Description: "Set the application menu structure"}, s.menuSet)
}
