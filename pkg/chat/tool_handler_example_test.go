package chat

import (
	"context"
	"fmt"
	"strings"

	guimcp "dappco.re/go/gui/pkg/mcp"
)

type exampleToolExecutor struct{}

func (exampleToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{"type": "object"},
	}}
}

func (exampleToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (exampleToolExecutor) CallTool(_ context.Context, name string, _ map[string]any) (string, error) {
	if name == "layout_suggest" {
		return `{"mode":"left-right"}`, nil
	}
	return "", nil
}

func ExampleNewToolCallHandler() {
	handler := NewToolCallHandler(exampleToolExecutor{})
	result, err := handler.OnToolCall(context.Background(), ToolCall{
		ID:        "call-1",
		Name:      "layout_suggest",
		Arguments: map[string]any{"window_count": 2},
	})

	fmt.Println(err == nil)
	fmt.Println(result)
	fmt.Println(strings.Contains(handler.BuildToolManifest(), "layout_suggest"))
	// Output:
	// true
	// {"mode":"left-right"}
	// true
}
