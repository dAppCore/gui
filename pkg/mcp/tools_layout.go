// pkg/mcp/tools_layout.go
package mcp

import (
	"context"
	"strings"

	"dappco.re/go/core/gui/pkg/screen"
	"dappco.re/go/core/gui/pkg/window"
	coreerr "forge.lthn.ai/core/go-log"
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
	_, _, err := s.core.PERFORM(window.TaskStackWindows{Windows: input.Windows, OffsetX: input.OffsetX, OffsetY: input.OffsetY})
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
	_, _, err := s.core.PERFORM(window.TaskApplyWorkflow{Workflow: input.Workflow, Windows: input.Windows})
	if err != nil {
		return nil, LayoutWorkflowOutput{}, err
	}
	return nil, LayoutWorkflowOutput{Success: true}, nil
}

// --- layout_suggest ---

type LayoutSuggestInput struct {
	Width       int `json:"width"`
	Height      int `json:"height"`
	WindowCount int `json:"windowCount"`
}
type LayoutSuggestOutput struct {
	Mode       string        `json:"mode"`
	Placements []screen.Rect `json:"placements"`
}

func (s *Subsystem) layoutSuggest(_ context.Context, _ *mcp.CallToolRequest, input LayoutSuggestInput) (*mcp.CallToolResult, LayoutSuggestOutput, error) {
	width := input.Width
	height := input.Height
	if width <= 0 {
		width = 1920
	}
	if height <= 0 {
		height = 1080
	}
	count := input.WindowCount
	if count <= 0 {
		count = 1
	}

	workArea := screen.Rect{X: 0, Y: 0, Width: width, Height: height}
	switch {
	case count == 1:
		return nil, LayoutSuggestOutput{Mode: "full", Placements: []screen.Rect{workArea}}, nil
	case count == 2:
		if width >= height {
			half := width / 2
			return nil, LayoutSuggestOutput{
				Mode: "side-by-side",
				Placements: []screen.Rect{
					{X: 0, Y: 0, Width: half, Height: height},
					{X: half, Y: 0, Width: width - half, Height: height},
				},
			}, nil
		}
		half := height / 2
		return nil, LayoutSuggestOutput{
			Mode: "stacked",
			Placements: []screen.Rect{
				{X: 0, Y: 0, Width: width, Height: half},
				{X: 0, Y: half, Width: width, Height: height - half},
			},
		}, nil
	case count == 3 && width >= height:
		mainWidth := width * 2 / 3
		sideHeight := height / 2
		return nil, LayoutSuggestOutput{
			Mode: "editor-plus-stack",
			Placements: []screen.Rect{
				{X: 0, Y: 0, Width: mainWidth, Height: height},
				{X: mainWidth, Y: 0, Width: width - mainWidth, Height: sideHeight},
				{X: mainWidth, Y: sideHeight, Width: width - mainWidth, Height: height - sideHeight},
			},
		}, nil
	default:
		cols := 2
		if count > 4 {
			cols = 3
		}
		rows := (count + cols - 1) / cols
		cellWidth := width / cols
		cellHeight := height / rows
		placements := make([]screen.Rect, 0, count)
		for i := 0; i < count; i++ {
			row := i / cols
			col := i % cols
			placements = append(placements, screen.Rect{
				X: col * cellWidth, Y: row * cellHeight,
				Width: cellWidth, Height: cellHeight,
			})
		}
		return nil, LayoutSuggestOutput{Mode: "grid", Placements: placements}, nil
	}
}

// --- layout_beside_editor ---

type LayoutBesideEditorInput struct {
	Name        string   `json:"name"`
	EditorNames []string `json:"editorNames,omitempty"`
}
type LayoutBesideEditorOutput struct {
	Editor string      `json:"editor"`
	Bounds screen.Rect `json:"bounds"`
}

func (s *Subsystem) layoutBesideEditor(_ context.Context, _ *mcp.CallToolRequest, input LayoutBesideEditorInput) (*mcp.CallToolResult, LayoutBesideEditorOutput, error) {
	windows, err := s.allWindows()
	if err != nil {
		return nil, LayoutBesideEditorOutput{}, err
	}
	screens, err := s.allScreens()
	if err != nil {
		return nil, LayoutBesideEditorOutput{}, err
	}

	editorHints := map[string]struct{}{}
	for _, name := range input.EditorNames {
		editorHints[strings.ToLower(name)] = struct{}{}
	}
	defaultHints := []string{"code", "cursor", "vscode", "studio", "goland", "intellij", "webstorm", "xcode", "vim", "nvim", "emacs", "editor"}

	var editor *window.WindowInfo
	for i := range windows {
		if windows[i].Name == input.Name {
			continue
		}
		name := strings.ToLower(windows[i].Name)
		title := strings.ToLower(windows[i].Title)
		if _, ok := editorHints[name]; ok {
			editor = &windows[i]
			break
		}
		for _, hint := range defaultHints {
			if strings.Contains(name, hint) || strings.Contains(title, hint) {
				editor = &windows[i]
				break
			}
		}
		if editor != nil {
			break
		}
	}
	if editor == nil {
		return nil, LayoutBesideEditorOutput{}, coreerr.E("mcp.layoutBesideEditor", "no editor window detected", nil)
	}

	editorScreen := screenForWindowInfo(screens, *editor)
	if editorScreen == nil {
		editorScreen = chooseScreenByIDOrPrimary(screens, "")
	}
	workArea := workAreaRect(editorScreen)

	editorRect := screen.Rect{X: editor.X, Y: editor.Y, Width: editor.Width, Height: editor.Height}
	candidates := []screen.Rect{
		{X: workArea.X, Y: workArea.Y, Width: max(0, editorRect.X-workArea.X), Height: workArea.Height},
		{X: editorRect.X + editorRect.Width, Y: workArea.Y, Width: max(0, workArea.X+workArea.Width-(editorRect.X+editorRect.Width)), Height: workArea.Height},
		{X: workArea.X, Y: workArea.Y, Width: workArea.Width, Height: max(0, editorRect.Y-workArea.Y)},
		{X: workArea.X, Y: editorRect.Y + editorRect.Height, Width: workArea.Width, Height: max(0, workArea.Y+workArea.Height-(editorRect.Y+editorRect.Height))},
	}

	best := screen.Rect{}
	bestArea := -1
	for _, candidate := range candidates {
		area := candidate.Width * candidate.Height
		if candidate.Width <= 0 || candidate.Height <= 0 {
			continue
		}
		if area > bestArea {
			bestArea = area
			best = candidate
		}
	}
	if bestArea <= 0 {
		arranged, err := s.arrangePairOnScreen(editor.Name, input.Name, editorScreen, "")
		if err != nil {
			return nil, LayoutBesideEditorOutput{}, err
		}
		return nil, LayoutBesideEditorOutput{Editor: editor.Name, Bounds: arranged.Second}, nil
	}

	if err := applyRect(s.core, input.Name, best); err != nil {
		return nil, LayoutBesideEditorOutput{}, err
	}
	return nil, LayoutBesideEditorOutput{Editor: editor.Name, Bounds: best}, nil
}

// --- Registration ---

func (s *Subsystem) registerLayoutTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "layout_save", Description: "Save the current window arrangement as a named layout"}, s.layoutSave)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_restore", Description: "Restore a saved window layout"}, s.layoutRestore)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_list", Description: "List all saved layouts"}, s.layoutList)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_delete", Description: "Delete a saved layout"}, s.layoutDelete)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_get", Description: "Get a specific layout by name"}, s.layoutGet)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_tile", Description: "Tile windows in a grid arrangement"}, s.layoutTile)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_suggest", Description: "Suggest an optimal arrangement for the given screen size and window count"}, s.layoutSuggest)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_beside_editor", Description: "Place a window beside a detected editor or IDE window"}, s.layoutBesideEditor)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_snap", Description: "Snap a window to a screen edge or corner"}, s.layoutSnap)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_stack", Description: "Stack windows in a cascade pattern"}, s.layoutStack)
	mcp.AddTool(server, &mcp.Tool{Name: "layout_workflow", Description: "Apply a preset workflow layout"}, s.layoutWorkflow)
}
