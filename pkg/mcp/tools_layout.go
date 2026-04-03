// pkg/mcp/tools_layout.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/screen"
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

// --- layout_beside_editor ---

type LayoutBesideEditorInput struct {
	Editor string `json:"editor,omitempty"`
	Window string `json:"window,omitempty"`
}
type LayoutBesideEditorOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutBesideEditor(_ context.Context, _ *mcp.CallToolRequest, input LayoutBesideEditorInput) (*mcp.CallToolResult, LayoutBesideEditorOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskBesideEditor{Editor: input.Editor, Window: input.Window})
	if err != nil {
		return nil, LayoutBesideEditorOutput{}, err
	}
	return nil, LayoutBesideEditorOutput{Success: true}, nil
}

// --- layout_suggest ---

type LayoutSuggestInput struct {
	WindowCount  int `json:"windowCount,omitempty"`
	ScreenWidth  int `json:"screenWidth,omitempty"`
	ScreenHeight int `json:"screenHeight,omitempty"`
}
type LayoutSuggestOutput struct {
	Suggestion window.LayoutSuggestion `json:"suggestion"`
}

func (s *Subsystem) layoutSuggest(_ context.Context, _ *mcp.CallToolRequest, input LayoutSuggestInput) (*mcp.CallToolResult, LayoutSuggestOutput, error) {
	windowCount := input.WindowCount
	if windowCount <= 0 {
		result, _, err := s.core.QUERY(window.QueryWindowList{})
		if err != nil {
			return nil, LayoutSuggestOutput{}, err
		}
		windows, ok := result.([]window.WindowInfo)
		if !ok {
			return nil, LayoutSuggestOutput{}, fmt.Errorf("unexpected result type from window list query")
		}
		windowCount = len(windows)
	}
	screenW, screenH := input.ScreenWidth, input.ScreenHeight
	if screenW <= 0 || screenH <= 0 {
		screenW, screenH = primaryScreenSize(s.core)
	}
	result, handled, err := s.core.QUERY(window.QueryLayoutSuggestion{
		WindowCount:  windowCount,
		ScreenWidth:  screenW,
		ScreenHeight: screenH,
	})
	if err != nil {
		return nil, LayoutSuggestOutput{}, err
	}
	if !handled {
		return nil, LayoutSuggestOutput{}, fmt.Errorf("window service not available")
	}
	suggestion, ok := result.(window.LayoutSuggestion)
	if !ok {
		return nil, LayoutSuggestOutput{}, fmt.Errorf("unexpected result type from layout suggestion query")
	}
	return nil, LayoutSuggestOutput{Suggestion: suggestion}, nil
}

// --- screen_find_space ---

type ScreenFindSpaceInput struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}
type ScreenFindSpaceOutput struct {
	Space window.SpaceInfo `json:"space"`
}

func (s *Subsystem) screenFindSpace(_ context.Context, _ *mcp.CallToolRequest, input ScreenFindSpaceInput) (*mcp.CallToolResult, ScreenFindSpaceOutput, error) {
	screenW, screenH := primaryScreenSize(s.core)
	if screenW <= 0 || screenH <= 0 {
		screenW, screenH = 1920, 1080
	}
	result, handled, err := s.core.QUERY(window.QueryFindSpace{
		Width:        input.Width,
		Height:       input.Height,
		ScreenWidth:  screenW,
		ScreenHeight: screenH,
	})
	if err != nil {
		return nil, ScreenFindSpaceOutput{}, err
	}
	if !handled {
		return nil, ScreenFindSpaceOutput{}, fmt.Errorf("window service not available")
	}
	space, ok := result.(window.SpaceInfo)
	if !ok {
		return nil, ScreenFindSpaceOutput{}, fmt.Errorf("unexpected result type from find space query")
	}
	if space.ScreenWidth == 0 {
		space.ScreenWidth = screenW
	}
	if space.ScreenHeight == 0 {
		space.ScreenHeight = screenH
	}
	return nil, ScreenFindSpaceOutput{Space: space}, nil
}

// --- window_arrange_pair ---

type WindowArrangePairInput struct {
	First  string `json:"first"`
	Second string `json:"second"`
}
type WindowArrangePairOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) windowArrangePair(_ context.Context, _ *mcp.CallToolRequest, input WindowArrangePairInput) (*mcp.CallToolResult, WindowArrangePairOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskArrangePair{First: input.First, Second: input.Second})
	if err != nil {
		return nil, WindowArrangePairOutput{}, err
	}
	return nil, WindowArrangePairOutput{Success: true}, nil
}

// --- layout_stack ---

type LayoutStackInput struct {
	Windows []string `json:"windows,omitempty"`
	OffsetX int      `json:"offsetX,omitempty"`
	OffsetY int      `json:"offsetY,omitempty"`
}
type LayoutStackOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) layoutStack(_ context.Context, _ *mcp.CallToolRequest, input LayoutStackInput) (*mcp.CallToolResult, LayoutStackOutput, error) {
	_, _, err := s.core.PERFORM(window.TaskStackWindows{
		Windows: input.Windows,
		OffsetX: input.OffsetX,
		OffsetY: input.OffsetY,
	})
	if err != nil {
		return nil, LayoutStackOutput{}, err
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
	workflow, ok := window.ParseWorkflowLayout(input.Workflow)
	if !ok {
		return nil, LayoutWorkflowOutput{}, fmt.Errorf("unknown workflow: %s", input.Workflow)
	}
	_, _, err := s.core.PERFORM(window.TaskApplyWorkflow{
		Workflow: workflow,
		Windows:  input.Windows,
	})
	if err != nil {
		return nil, LayoutWorkflowOutput{}, err
	}
	return nil, LayoutWorkflowOutput{Success: true}, nil
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
	mcp.AddTool(server, &mcp.Tool{Name: "layout_beside_editor", Description: "Place a window beside a detected editor window"}, s.layoutBesideEditor)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_suggest", Description: "Suggest an optimal layout for the current screen"}, s.layoutSuggest)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_find_space", Description: "Find an empty area for a new window"}, s.screenFindSpace)
	mcp.AddTool(server, &mcp.Tool{Name: "window_arrange_pair", Description: "Arrange two windows side-by-side"}, s.windowArrangePair)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_stack", Description: "Cascade windows with an offset"}, s.layoutStack)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_workflow", Description: "Apply a predefined workflow layout"}, s.layoutWorkflow)
}

func primaryScreenSize(c *core.Core) (int, int) {
	result, handled, err := c.QUERY(screen.QueryPrimary{})
	if err == nil && handled {
		if scr, ok := result.(*screen.Screen); ok && scr != nil {
			if scr.WorkArea.Width > 0 && scr.WorkArea.Height > 0 {
				return scr.WorkArea.Width, scr.WorkArea.Height
			}
			if scr.Bounds.Width > 0 && scr.Bounds.Height > 0 {
				return scr.Bounds.Width, scr.Bounds.Height
			}
			if scr.Size.Width > 0 && scr.Size.Height > 0 {
				return scr.Size.Width, scr.Size.Height
			}
		}
	}
	return 1920, 1080
}
