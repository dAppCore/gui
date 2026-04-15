// pkg/mcp/tools_environment.go
package mcp

import (
	"context"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"forge.lthn.ai/core/gui/pkg/environment"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- theme_get ---

type ThemeGetInput struct{}
type ThemeGetOutput struct {
	Theme environment.ThemeInfo `json:"theme"`
}

func (s *Subsystem) themeGet(_ context.Context, _ *mcp.CallToolRequest, _ ThemeGetInput) (*mcp.CallToolResult, ThemeGetOutput, error) {
	r := s.core.QUERY(environment.QueryTheme{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ThemeGetOutput{}, e
		}
		return nil, ThemeGetOutput{}, nil
	}
	theme, ok := r.Value.(environment.ThemeInfo)
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
	r := s.core.QUERY(environment.QueryInfo{})
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ThemeSystemOutput{}, e
		}
		return nil, ThemeSystemOutput{}, nil
	}
	info, ok := r.Value.(environment.EnvironmentInfo)
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
	Theme environment.ThemeInfo `json:"theme"`
}

func (s *Subsystem) themeSet(_ context.Context, _ *mcp.CallToolRequest, input ThemeSetInput) (*mcp.CallToolResult, ThemeSetOutput, error) {
	r := s.core.Action("environment.setTheme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: environment.TaskSetTheme{Theme: input.Theme}},
	))
	if !r.OK {
		if e, ok := r.Value.(error); ok {
			return nil, ThemeSetOutput{}, e
		}
		return nil, ThemeSetOutput{}, nil
	}
	theme, ok := r.Value.(environment.ThemeInfo)
	if !ok {
		return nil, ThemeSetOutput{}, coreerr.E("mcp.themeSet", "unexpected result type", nil)
	}
	return nil, ThemeSetOutput{Theme: theme}, nil
}

// --- Registration ---

func (s *Subsystem) registerEnvironmentTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "theme_get", Description: "Get the current application theme"}, s.themeGet)
	addTool(s, server, &mcp.Tool{Name: "theme_set", Description: "Set the application theme to dark, light, or system"}, s.themeSet)
	addTool(s, server, &mcp.Tool{Name: "theme_system", Description: "Get system environment and theme information"}, s.themeSystem)
}
