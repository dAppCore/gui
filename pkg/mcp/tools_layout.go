// pkg/mcp/tools_layout.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
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
	r := s.core.Action("window.saveLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSaveLayout{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutSaveOutput{}, e
		}
		return nil, LayoutSaveOutput{}, nil
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
	r := s.core.Action("window.restoreLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskRestoreLayout{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutRestoreOutput{}, e
		}
		return nil, LayoutRestoreOutput{}, nil
	}
	return nil, LayoutRestoreOutput{Success: true}, nil
}

// --- layout_list ---

type LayoutListInput struct{}
type LayoutListOutput struct {
	Layouts []window.LayoutInfo `json:"layouts"`
}

func (s *Subsystem) layoutList(_ context.Context, _ *mcp.CallToolRequest, _ LayoutListInput) (*mcp.CallToolResult, LayoutListOutput, error) {
	r := s.core.QUERY(window.QueryLayoutList{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutListOutput{}, e
		}
		return nil, LayoutListOutput{}, nil
	}
	layouts, ok := r.Value.([]window.LayoutInfo)
	if !ok {
		return nil, LayoutListOutput{}, coreerr.E("mcp.layoutList", "unexpected result type", nil)
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
	r := s.core.Action("window.deleteLayout").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskDeleteLayout{Name: input.Name}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutDeleteOutput{}, e
		}
		return nil, LayoutDeleteOutput{}, nil
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
	r := s.core.QUERY(window.QueryLayoutGet{Name: input.Name})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutGetOutput{}, e
		}
		return nil, LayoutGetOutput{}, nil
	}
	layout, ok := r.Value.(*window.Layout)
	if !ok {
		return nil, LayoutGetOutput{}, coreerr.E("mcp.layoutGet", "unexpected result type", nil)
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
	r := s.core.Action("window.tileWindows").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskTileWindows{Mode: input.Mode, Windows: input.Windows}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutTileOutput{}, e
		}
		return nil, LayoutTileOutput{}, nil
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
	r := s.core.Action("window.snapWindow").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSnapWindow{Name: input.Name, Position: input.Position}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutSnapOutput{}, e
		}
		return nil, LayoutSnapOutput{}, nil
	}
	return nil, LayoutSnapOutput{Success: true}, nil
}

// --- layout_stack ---

type LayoutStackInput struct {
	Windows []string `json:"windows,omitempty"`
	OffsetX int      `json:"offsetX"`
	OffsetY int      `json:"offsetY"`
}
type LayoutStackOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutStack(_ context.Context, _ *mcp.CallToolRequest, input LayoutStackInput) (*mcp.CallToolResult, LayoutStackOutput, error) {
	r := s.core.Action("window.stackWindows").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskStackWindows{Windows: input.Windows, OffsetX: input.OffsetX, OffsetY: input.OffsetY}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutStackOutput{}, e
		}
		return nil, LayoutStackOutput{}, nil
	}
	return nil, LayoutStackOutput{Success: true}, nil
}

// --- layout_workflow ---

type LayoutWorkflowInput struct {
	Workflow string   `json:"workflow"`
	Windows  []string `json:"windows,omitempty"`
}
type LayoutWorkflowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutWorkflow(_ context.Context, _ *mcp.CallToolRequest, input LayoutWorkflowInput) (*mcp.CallToolResult, LayoutWorkflowOutput, error) {
	r := s.core.Action("window.applyWorkflow").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskApplyWorkflow{Workflow: input.Workflow, Windows: input.Windows}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutWorkflowOutput{}, e
		}
		return nil, LayoutWorkflowOutput{}, nil
	}
	return nil, LayoutWorkflowOutput{Success: true}, nil
}

// --- layout_beside_editor ---

type LayoutBesideEditorInput struct {
	Name   string  `json:"name"`
	Editor string  `json:"editor,omitempty"`
	Side   string  `json:"side,omitempty"`
	Ratio  float64 `json:"ratio,omitempty"`
}

type LayoutBesideEditorOutput struct {
	Result window.LayoutBesideEditorResult `json:"result"`
}

func (s *Subsystem) layoutBesideEditor(_ context.Context, _ *mcp.CallToolRequest, input LayoutBesideEditorInput) (*mcp.CallToolResult, LayoutBesideEditorOutput, error) {
	r := s.core.Action("window.layoutBesideEditor").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskLayoutBesideEditor{
			Name: input.Name, Editor: input.Editor, Side: input.Side, Ratio: input.Ratio,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutBesideEditorOutput{}, e
		}
		return nil, LayoutBesideEditorOutput{}, nil
	}
	result, ok := r.Value.(window.LayoutBesideEditorResult)
	if !ok {
		return nil, LayoutBesideEditorOutput{}, coreerr.E("mcp.layoutBesideEditor", "unexpected result type", nil)
	}
	return nil, LayoutBesideEditorOutput{Result: result}, nil
}

// --- layout_suggest ---

type LayoutSuggestInput struct {
	ScreenID    string `json:"screen_id,omitempty"`
	WindowCount int    `json:"window_count,omitempty"`
}

type LayoutSuggestOutput struct {
	Suggestion window.LayoutSuggestion `json:"suggestion"`
}

func (s *Subsystem) layoutSuggest(_ context.Context, _ *mcp.CallToolRequest, input LayoutSuggestInput) (*mcp.CallToolResult, LayoutSuggestOutput, error) {
	r := s.core.Action("window.layoutSuggest").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskLayoutSuggest{ScreenID: input.ScreenID, WindowCount: input.WindowCount}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, LayoutSuggestOutput{}, e
		}
		return nil, LayoutSuggestOutput{}, nil
	}
	result, ok := r.Value.(window.LayoutSuggestion)
	if !ok {
		return nil, LayoutSuggestOutput{}, coreerr.E("mcp.layoutSuggest", "unexpected result type", nil)
	}
	return nil, LayoutSuggestOutput{Suggestion: result}, nil
}

// --- screen_find_space ---

type ScreenFindSpaceInput struct {
	ScreenID string `json:"screen_id,omitempty"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Padding  int    `json:"padding,omitempty"`
}

type ScreenFindSpaceOutput struct {
	Space window.ScreenSpace `json:"space"`
}

func (s *Subsystem) screenFindSpace(_ context.Context, _ *mcp.CallToolRequest, input ScreenFindSpaceInput) (*mcp.CallToolResult, ScreenFindSpaceOutput, error) {
	r := s.core.Action("window.findSpace").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskScreenFindSpace{
			ScreenID: input.ScreenID, Width: input.Width, Height: input.Height, Padding: input.Padding,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ScreenFindSpaceOutput{}, e
		}
		return nil, ScreenFindSpaceOutput{}, nil
	}
	result, ok := r.Value.(window.ScreenSpace)
	if !ok {
		return nil, ScreenFindSpaceOutput{}, coreerr.E("mcp.screenFindSpace", "unexpected result type", nil)
	}
	return nil, ScreenFindSpaceOutput{Space: result}, nil
}

// --- window_arrange_pair ---

type WindowArrangePairInput struct {
	Primary   string  `json:"primary"`
	Secondary string  `json:"secondary"`
	ScreenID  string  `json:"screen_id,omitempty"`
	Ratio     float64 `json:"ratio,omitempty"`
}

type WindowArrangePairOutput struct {
	Arrangement window.PairArrangement `json:"arrangement"`
}

func (s *Subsystem) windowArrangePair(_ context.Context, _ *mcp.CallToolRequest, input WindowArrangePairInput) (*mcp.CallToolResult, WindowArrangePairOutput, error) {
	r := s.core.Action("window.arrangePair").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskWindowArrangePair{
			Primary: input.Primary, Secondary: input.Secondary, ScreenID: input.ScreenID, Ratio: input.Ratio,
		}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, WindowArrangePairOutput{}, e
		}
		return nil, WindowArrangePairOutput{}, nil
	}
	result, ok := r.Value.(window.PairArrangement)
	if !ok {
		return nil, WindowArrangePairOutput{}, coreerr.E("mcp.windowArrangePair", "unexpected result type", nil)
	}
	return nil, WindowArrangePairOutput{Arrangement: result}, nil
}

// --- Registration ---

func (s *Subsystem) registerLayoutTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "layout_save", Description: "Save the current window arrangement as a named layout"}, s.layoutSave)
	addTool(s, server, &mcp.Tool{Name: "layout_restore", Description: "Restore a saved window layout"}, s.layoutRestore)
	addTool(s, server, &mcp.Tool{Name: "layout_list", Description: "List all saved layouts"}, s.layoutList)
	addTool(s, server, &mcp.Tool{Name: "layout_delete", Description: "Delete a saved layout"}, s.layoutDelete)
	addTool(s, server, &mcp.Tool{Name: "layout_get", Description: "Get a specific layout by name"}, s.layoutGet)
	addTool(s, server, &mcp.Tool{Name: "layout_tile", Description: "Tile windows in a grid arrangement"}, s.layoutTile)
	addTool(s, server, &mcp.Tool{Name: "layout_snap", Description: "Snap a window to a screen edge or corner"}, s.layoutSnap)
	addTool(s, server, &mcp.Tool{Name: "layout_stack", Description: "Stack windows in a cascade pattern"}, s.layoutStack)
	addTool(s, server, &mcp.Tool{Name: "layout_workflow", Description: "Apply a preset workflow layout"}, s.layoutWorkflow)
	addTool(s, server, &mcp.Tool{Name: "layout_beside_editor", Description: "Position a window beside the detected editor window"}, s.layoutBesideEditor)
	addTool(s, server, &mcp.Tool{Name: "layout_suggest", Description: "Suggest the best layout for the current screen and window count"}, s.layoutSuggest)
	addTool(s, server, &mcp.Tool{Name: "screen_find_space", Description: "Find an empty rectangle on a screen for a new window"}, s.screenFindSpace)
	addTool(s, server, &mcp.Tool{Name: "window_arrange_pair", Description: "Arrange two windows in an optimal split on one screen"}, s.windowArrangePair)
}
