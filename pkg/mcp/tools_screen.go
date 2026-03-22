// pkg/mcp/tools_screen.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/screen"
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
		return nil, ScreenListOutput{}, fmt.Errorf("unexpected result type from screen list query")
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
		return nil, ScreenGetOutput{}, fmt.Errorf("unexpected result type from screen get query")
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
		return nil, ScreenPrimaryOutput{}, fmt.Errorf("unexpected result type from screen primary query")
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
		return nil, ScreenAtPointOutput{}, fmt.Errorf("unexpected result type from screen at point query")
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
		return nil, ScreenWorkAreasOutput{}, fmt.Errorf("unexpected result type from screen work areas query")
	}
	return nil, ScreenWorkAreasOutput{WorkAreas: areas}, nil
}

// --- Registration ---

func (s *Subsystem) registerScreenTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "screen_list", Description: "List all connected displays/screens"}, s.screenList)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_get", Description: "Get information about a specific screen"}, s.screenGet)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_primary", Description: "Get the primary screen"}, s.screenPrimary)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_at_point", Description: "Get the screen at a specific point"}, s.screenAtPoint)
	mcp.AddTool(server, &mcp.Tool{Name: "screen_work_areas", Description: "Get work areas for all screens"}, s.screenWorkAreas)
}
