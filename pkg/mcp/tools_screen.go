// pkg/mcp/tools_screen.go
package mcp

import (
	"context"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- screen_list ---

type ScreenListInput struct{}
type ScreenListOutput struct {
	Screens []screen.Screen `json:"screens"`
}

func (s *Subsystem) screenList(_ context.Context, _ *mcp.CallToolRequest, _ ScreenListInput) (*mcp.CallToolResult, ScreenListOutput, error) {
	result, _, err := s.core.QUERY(screen.QueryAll{})
	if err != nil {
		return nil, ScreenListOutput{}, err
	}
	screens, ok := result.([]screen.Screen)
	if !ok {
		return nil, ScreenListOutput{}, coreerr.E("mcp.screenList", "unexpected result type", nil)
	}
	return nil, ScreenListOutput{Screens: screens}, nil
}

// --- screen_get ---

type ScreenGetInput struct {
	ID string `json:"id"`
}
type ScreenGetOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenGet(_ context.Context, _ *mcp.CallToolRequest, input ScreenGetInput) (*mcp.CallToolResult, ScreenGetOutput, error) {
	result, _, err := s.core.QUERY(screen.QueryByID{ID: input.ID})
	if err != nil {
		return nil, ScreenGetOutput{}, err
	}
	scr, ok := result.(*screen.Screen)
	if !ok {
		return nil, ScreenGetOutput{}, coreerr.E("mcp.screenGet", "unexpected result type", nil)
	}
	return nil, ScreenGetOutput{Screen: scr}, nil
}

// --- screen_primary ---

type ScreenPrimaryInput struct{}
type ScreenPrimaryOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenPrimary(_ context.Context, _ *mcp.CallToolRequest, _ ScreenPrimaryInput) (*mcp.CallToolResult, ScreenPrimaryOutput, error) {
	result, _, err := s.core.QUERY(screen.QueryPrimary{})
	if err != nil {
		return nil, ScreenPrimaryOutput{}, err
	}
	scr, ok := result.(*screen.Screen)
	if !ok {
		return nil, ScreenPrimaryOutput{}, coreerr.E("mcp.screenPrimary", "unexpected result type", nil)
	}
	return nil, ScreenPrimaryOutput{Screen: scr}, nil
}

// --- screen_at_point ---

type ScreenAtPointInput struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type ScreenAtPointOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenAtPoint(_ context.Context, _ *mcp.CallToolRequest, input ScreenAtPointInput) (*mcp.CallToolResult, ScreenAtPointOutput, error) {
	result, _, err := s.core.QUERY(screen.QueryAtPoint{X: input.X, Y: input.Y})
	if err != nil {
		return nil, ScreenAtPointOutput{}, err
	}
	scr, ok := result.(*screen.Screen)
	if !ok {
		return nil, ScreenAtPointOutput{}, coreerr.E("mcp.screenAtPoint", "unexpected result type", nil)
	}
	return nil, ScreenAtPointOutput{Screen: scr}, nil
}

// --- screen_work_areas ---

type ScreenWorkAreasInput struct{}
type ScreenWorkAreasOutput struct {
	WorkAreas []screen.Rect `json:"workAreas"`
}

func (s *Subsystem) screenWorkAreas(_ context.Context, _ *mcp.CallToolRequest, _ ScreenWorkAreasInput) (*mcp.CallToolResult, ScreenWorkAreasOutput, error) {
	result, _, err := s.core.QUERY(screen.QueryWorkAreas{})
	if err != nil {
		return nil, ScreenWorkAreasOutput{}, err
	}
	areas, ok := result.([]screen.Rect)
	if !ok {
		return nil, ScreenWorkAreasOutput{}, coreerr.E("mcp.screenWorkAreas", "unexpected result type", nil)
	}
	return nil, ScreenWorkAreasOutput{WorkAreas: areas}, nil
}

// --- screen_work_area ---

type ScreenWorkAreaInput struct {
	ID string `json:"id,omitempty"`
}
type ScreenWorkAreaOutput struct {
	WorkArea screen.Rect `json:"workArea"`
}

func (s *Subsystem) screenWorkArea(_ context.Context, _ *mcp.CallToolRequest, input ScreenWorkAreaInput) (*mcp.CallToolResult, ScreenWorkAreaOutput, error) {
	screens, err := s.allScreens()
	if err != nil {
		return nil, ScreenWorkAreaOutput{}, err
	}
	scr := chooseScreenByIDOrPrimary(screens, input.ID)
	if scr == nil {
		return nil, ScreenWorkAreaOutput{}, nil
	}
	return nil, ScreenWorkAreaOutput{WorkArea: workAreaRect(scr)}, nil
}

// --- screen_find_space ---

type ScreenFindSpaceInput struct {
	ScreenID string `json:"screenId,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}
type ScreenFindSpaceOutput struct {
	ScreenID string      `json:"screenId"`
	Bounds   screen.Rect `json:"bounds"`
}

func (s *Subsystem) screenFindSpace(_ context.Context, _ *mcp.CallToolRequest, input ScreenFindSpaceInput) (*mcp.CallToolResult, ScreenFindSpaceOutput, error) {
	screens, err := s.allScreens()
	if err != nil {
		return nil, ScreenFindSpaceOutput{}, err
	}
	windows, err := s.allWindows()
	if err != nil {
		return nil, ScreenFindSpaceOutput{}, err
	}

	orderedScreens := make([]screen.Screen, 0, len(screens))
	if selected := chooseScreenByIDOrPrimary(screens, input.ScreenID); selected != nil {
		orderedScreens = append(orderedScreens, *selected)
		for _, scr := range screens {
			if scr.ID != selected.ID {
				orderedScreens = append(orderedScreens, scr)
			}
		}
	}

	for _, scr := range orderedScreens {
		workArea := workAreaRect(&scr)
		occupied := make([]screen.Rect, 0, len(windows))
		for _, info := range windows {
			if windowScreen := screenForWindowInfo(screens, info); windowScreen != nil && windowScreen.ID == scr.ID {
				occupied = append(occupied, screen.Rect{X: info.X, Y: info.Y, Width: info.Width, Height: info.Height})
			}
		}
		if candidate, ok := findLargestFreeRect(workArea, occupied, input.Width, input.Height); ok {
			return nil, ScreenFindSpaceOutput{ScreenID: scr.ID, Bounds: candidate}, nil
		}
	}

	return nil, ScreenFindSpaceOutput{}, nil
}

// --- screen_for_window ---

type ScreenForWindowInput struct {
	Name string `json:"name"`
}
type ScreenForWindowOutput struct {
	Screen *screen.Screen `json:"screen"`
}

func (s *Subsystem) screenForWindow(_ context.Context, _ *mcp.CallToolRequest, input ScreenForWindowInput) (*mcp.CallToolResult, ScreenForWindowOutput, error) {
	result, _, err := s.core.QUERY(window.QueryWindowByName{Name: input.Name})
	if err != nil {
		return nil, ScreenForWindowOutput{}, err
	}
	info, _ := result.(*window.WindowInfo)
	if info == nil {
		return nil, ScreenForWindowOutput{}, nil
	}
	centerX := info.X + info.Width/2
	centerY := info.Y + info.Height/2
	screenResult, _, err := s.core.QUERY(screen.QueryAtPoint{X: centerX, Y: centerY})
	if err != nil {
		return nil, ScreenForWindowOutput{}, err
	}
	scr, _ := screenResult.(*screen.Screen)
	return nil, ScreenForWindowOutput{Screen: scr}, nil
}

// --- Registration ---

func (s *Subsystem) registerScreenTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "screen_list", Description: "List all connected displays/screens"}, s.screenList)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_get", Description: "Get information about a specific screen"}, s.screenGet)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_primary", Description: "Get the primary screen"}, s.screenPrimary)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_at_point", Description: "Get the screen at a specific point"}, s.screenAtPoint)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_work_area", Description: "Get the work area for a screen"}, s.screenWorkArea)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_work_areas", Description: "Get work areas for all screens"}, s.screenWorkAreas)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_find_space", Description: "Find the largest empty area on a screen that fits the requested size"}, s.screenFindSpace)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_for_window", Description: "Get the screen containing a window"}, s.screenForWindow)
}
