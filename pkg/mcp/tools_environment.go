// pkg/mcp/tools_environment.go
package mcp

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/gui/pkg/environment"
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
		return nil, ThemeGetOutput{}, fmt.Errorf("unexpected result type from theme query")
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
		return nil, ThemeSystemOutput{}, fmt.Errorf("unexpected result type from environment info query")
	}
	return nil, ThemeSystemOutput{Info: info}, nil
}

// --- Registration ---

func (s *Subsystem) registerEnvironmentTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "theme_get", Description: "Get the current application theme"}, s.themeGet)
	mcp.AddTool(server, &mcp.Tool{Name: "theme_system", Description: "Get system environment and theme information"}, s.themeSystem)
}
