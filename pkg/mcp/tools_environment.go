// pkg/mcp/tools_environment.go
package mcp

import (
	"context"

	"dappco.re/go/core/gui/pkg/display"
	"dappco.re/go/core/gui/pkg/environment"
	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- theme_get ---

type ThemeGetInput struct{}
type ThemeGetOutput struct {
	Theme environment.ThemeInfo `json:"theme"`
}

func (s *Subsystem) themeGet(_ context.Context, _ *mcp.CallToolRequest, _ ThemeGetInput) (*mcp.CallToolResult, ThemeGetOutput, error) {
	result, _, err := s.core.QUERY(environment.QueryTheme{})
	if err != nil {
		return nil, ThemeGetOutput{}, err
	}
	theme, ok := result.(environment.ThemeInfo)
	if !ok {
		return nil, ThemeGetOutput{}, coreerr.E("mcp.themeGet", "unexpected result type", nil)
	}
	return nil, ThemeGetOutput{Theme: theme}, nil
}

// --- theme_system ---

type ThemeSystemInput struct{}
type ThemeSystemOutput struct {
	Info environment.EnvironmentInfo `json:"info"`
}

func (s *Subsystem) themeSystem(_ context.Context, _ *mcp.CallToolRequest, _ ThemeSystemInput) (*mcp.CallToolResult, ThemeSystemOutput, error) {
	result, _, err := s.core.QUERY(environment.QueryInfo{})
	if err != nil {
		return nil, ThemeSystemOutput{}, err
	}
	info, ok := result.(environment.EnvironmentInfo)
	if !ok {
		return nil, ThemeSystemOutput{}, coreerr.E("mcp.themeSystem", "unexpected result type", nil)
	}
	return nil, ThemeSystemOutput{Info: info}, nil
}

// --- theme_set ---

type ThemeSetInput struct {
	Theme string `json:"theme"`
}
type ThemeSetOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) themeSet(_ context.Context, _ *mcp.CallToolRequest, input ThemeSetInput) (*mcp.CallToolResult, ThemeSetOutput, error) {
	_, _, err := s.core.PERFORM(environment.TaskSetTheme{Theme: input.Theme})
	if err != nil {
		return nil, ThemeSetOutput{}, err
	}
	return nil, ThemeSetOutput{Success: true}, nil
}

// --- theme_on_change ---

type ThemeOnChangeInput struct{}
type ThemeOnChangeOutput struct {
	Event  string                  `json:"event"`
	Theme  environment.ThemeInfo   `json:"theme"`
	Server display.EventServerInfo `json:"server"`
}

func (s *Subsystem) themeOnChange(_ context.Context, _ *mcp.CallToolRequest, _ ThemeOnChangeInput) (*mcp.CallToolResult, ThemeOnChangeOutput, error) {
	result, _, err := s.core.QUERY(environment.QueryTheme{})
	if err != nil {
		return nil, ThemeOnChangeOutput{}, err
	}
	theme, ok := result.(environment.ThemeInfo)
	if !ok {
		return nil, ThemeOnChangeOutput{}, coreerr.E("mcp.themeOnChange", "unexpected result type", nil)
	}

	output := ThemeOnChangeOutput{
		Event: "theme.change",
		Theme: theme,
	}

	displaySvc, err := core.ServiceFor[*display.Service](s.core, "display")
	if err == nil && displaySvc != nil {
		output.Server = displaySvc.GetEventInfo()
	}

	return nil, output, nil
}

// --- Registration ---

func (s *Subsystem) registerEnvironmentTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "theme_get", Description: "Get the current application theme"}, s.themeGet)
	mcp.AddTool(server, &mcp.Tool{Name: "theme_set", Description: "Override the application theme to light, dark, or system"}, s.themeSet)
	mcp.AddTool(server, &mcp.Tool{Name: "theme_system", Description: "Get system environment and theme information"}, s.themeSystem)
	mcp.AddTool(server, &mcp.Tool{Name: "theme_on_change", Description: "Describe the theme.change event stream exposed on the display event WebSocket"}, s.themeOnChange)
}
