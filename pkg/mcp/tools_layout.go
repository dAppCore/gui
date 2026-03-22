// pkg/mcp/tools_layout.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- layout_save ---

type LayoutSaveInput struct {
	Name string `json:"name"`
}
type LayoutSaveOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutSave(_ context.Context, _ *mcp.CallToolRequest, input LayoutSaveInput) (*mcp.CallToolResult, LayoutSaveOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSaveLayout{Name: input.Name})
	if err != nil {
		return nil, LayoutSaveOutput{}, err
	}
	return nil, LayoutSaveOutput{Success: true}, nil
}

// --- layout_restore ---

type LayoutRestoreInput struct {
	Name string `json:"name"`
}
type LayoutRestoreOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutRestore(_ context.Context, _ *mcp.CallToolRequest, input LayoutRestoreInput) (*mcp.CallToolResult, LayoutRestoreOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskRestoreLayout{Name: input.Name})
	if err != nil {
		return nil, LayoutRestoreOutput{}, err
	}
	return nil, LayoutRestoreOutput{Success: true}, nil
}

// --- layout_list ---

type LayoutListInput struct{}
type LayoutListOutput struct {
	Layouts []window.LayoutInfo `json:"layouts"`
}

func (s *Subsystem) layoutList(_ context.Context, _ *mcp.CallToolRequest, _ LayoutListInput) (*mcp.CallToolResult, LayoutListOutput, error) {
	result, _, err := s.core.QUERY(window.QueryLayoutList{})
	if err != nil {
		return nil, LayoutListOutput{}, err
	}
	layouts, ok := result.([]window.LayoutInfo)
	if !ok {
		return nil, LayoutListOutput{}, fmt.Errorf("unexpected result type from layout list query")
	}
	return nil, LayoutListOutput{Layouts: layouts}, nil
}

// --- layout_delete ---

type LayoutDeleteInput struct {
	Name string `json:"name"`
}
type LayoutDeleteOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutDelete(_ context.Context, _ *mcp.CallToolRequest, input LayoutDeleteInput) (*mcp.CallToolResult, LayoutDeleteOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskDeleteLayout{Name: input.Name})
	if err != nil {
		return nil, LayoutDeleteOutput{}, err
	}
	return nil, LayoutDeleteOutput{Success: true}, nil
}

// --- layout_get ---

type LayoutGetInput struct {
	Name string `json:"name"`
}
type LayoutGetOutput struct {
	Layout *window.Layout `json:"layout"`
}

func (s *Subsystem) layoutGet(_ context.Context, _ *mcp.CallToolRequest, input LayoutGetInput) (*mcp.CallToolResult, LayoutGetOutput, error) {
	result, _, err := s.core.QUERY(window.QueryLayoutGet{Name: input.Name})
	if err != nil {
		return nil, LayoutGetOutput{}, err
	}
	layout, ok := result.(*window.Layout)
	if !ok {
		return nil, LayoutGetOutput{}, fmt.Errorf("unexpected result type from layout get query")
	}
	return nil, LayoutGetOutput{Layout: layout}, nil
}

// --- layout_tile ---

type LayoutTileInput struct {
	Mode    string   `json:"mode"`
	Windows []string `json:"windows,omitempty"`
}
type LayoutTileOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutTile(_ context.Context, _ *mcp.CallToolRequest, input LayoutTileInput) (*mcp.CallToolResult, LayoutTileOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskTileWindows{Mode: input.Mode, Windows: input.Windows})
	if err != nil {
		return nil, LayoutTileOutput{}, err
	}
	return nil, LayoutTileOutput{Success: true}, nil
}

// --- layout_snap ---

type LayoutSnapInput struct {
	Name     string `json:"name"`
	Position string `json:"position"`
}
type LayoutSnapOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutSnap(_ context.Context, _ *mcp.CallToolRequest, input LayoutSnapInput) (*mcp.CallToolResult, LayoutSnapOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskSnapWindow{Name: input.Name, Position: input.Position})
	if err != nil {
		return nil, LayoutSnapOutput{}, err
	}
	return nil, LayoutSnapOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerLayoutTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "layout_save", Description: "Save the current window arrangement as a named layout"}, s.layoutSave)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_restore", Description: "Restore a saved window layout"}, s.layoutRestore)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_list", Description: "List all saved layouts"}, s.layoutList)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_delete", Description: "Delete a saved layout"}, s.layoutDelete)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_get", Description: "Get a specific layout by name"}, s.layoutGet)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_tile", Description: "Tile windows in a grid arrangement"}, s.layoutTile)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_snap", Description: "Snap a window to a screen edge or corner"}, s.layoutSnap)
}
