package chat

import (
	"context"
	"sort"
	"strings"

	core "dappco.re/go/core"
	guimcp "dappco.re/go/gui/pkg/mcp"
)

const mcpToolActionPrefix = "mcp.tool."

// ToolExecutor is the chat-facing subset of the GUI MCP subsystem.
type ToolExecutor interface {
	Manifest() []guimcp.ToolDescriptor
	ManifestText() string
	CallTool(ctx context.Context, name string, arguments map[string]any) (string, error)
}

// ToolCallHandler intercepts model-emitted tool calls and renders the tool
// manifest that is injected into the system prompt.
type ToolCallHandler interface {
	OnToolCall(ctx context.Context, call ToolCall) (result any, err error)
	BuildToolManifest() string
}

type mcpToolCallHandler struct {
	executor ToolExecutor
}

func NewToolCallHandler(executor ToolExecutor) ToolCallHandler {
	if executor == nil {
		return nil
	}
	return &mcpToolCallHandler{executor: executor}
}

func (h *mcpToolCallHandler) OnToolCall(ctx context.Context, call ToolCall) (any, error) {
	if h == nil || h.executor == nil {
		return nil, core.E("chat.tool_call", "tool execution unavailable", nil)
	}
	call.Name = strings.TrimSpace(call.Name)
	if call.Name == "" {
		return nil, core.E("chat.tool_call", "tool name is required", nil)
	}
	if call.Arguments == nil {
		call.Arguments = map[string]any{}
	}
	return h.executor.CallTool(ctx, call.Name, call.Arguments)
}

func (h *mcpToolCallHandler) BuildToolManifest() string {
	if h == nil || h.executor == nil {
		return ""
	}

	tools := h.executor.Manifest()
	if len(tools) == 0 {
		return strings.TrimSpace(h.executor.ManifestText())
	}
	tools = append([]guimcp.ToolDescriptor(nil), tools...)
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	var builder strings.Builder
	builder.WriteString("Available MCP tools:\n")
	for _, tool := range tools {
		builder.WriteString("- ")
		builder.WriteString(tool.Name)
		if strings.TrimSpace(tool.Description) != "" {
			builder.WriteString(": ")
			builder.WriteString(strings.TrimSpace(tool.Description))
		}
		schema := tool.InputSchema
		if schema == nil {
			schema = map[string]any{"type": "object"}
		}
		builder.WriteString("\n  input_schema: ")
		builder.WriteString(jsonString(schema))
		builder.WriteString("\n")
	}
	builder.WriteString("\nWhen a tool is needed, emit exactly one JSON object in this shape: ")
	builder.WriteString(`{"tool_call":{"name":"tool_name","arguments":{}}}`)
	builder.WriteString(".")
	return strings.TrimSpace(builder.String())
}

type actionToolExecutor struct {
	core     *core.Core
	fallback ToolExecutor
}

func newActionToolExecutor(c *core.Core, fallback ToolExecutor) ToolExecutor {
	if c == nil || fallback == nil {
		return fallback
	}
	return &actionToolExecutor{core: c, fallback: fallback}
}

func registerMCPToolActions(c *core.Core, executor ToolExecutor) {
	if c == nil || executor == nil {
		return
	}
	for _, tool := range executor.Manifest() {
		name := strings.TrimSpace(tool.Name)
		if name == "" {
			continue
		}
		actionName := mcpToolActionPrefix + name
		if c.Action(actionName).Exists() {
			continue
		}
		c.Action(actionName, func(ctx context.Context, opts core.Options) core.Result {
			content, err := executor.CallTool(ctx, name, toolArgumentsFromOptions(opts))
			return core.Result{}.New(content, err)
		})
	}
}

func (e *actionToolExecutor) Manifest() []guimcp.ToolDescriptor {
	if e == nil || e.fallback == nil {
		return nil
	}
	return e.fallback.Manifest()
}

func (e *actionToolExecutor) ManifestText() string {
	if e == nil || e.fallback == nil {
		return ""
	}
	return e.fallback.ManifestText()
}

func (e *actionToolExecutor) CallTool(ctx context.Context, name string, arguments map[string]any) (string, error) {
	if e == nil || e.fallback == nil {
		return "", core.E("chat.tool_call", "tool execution unavailable", nil)
	}
	if e.core != nil {
		result := e.core.Action(mcpToolActionPrefix+strings.TrimSpace(name)).Run(ctx, core.NewOptions(core.Option{
			Key:   "arguments",
			Value: arguments,
		}))
		if !result.OK {
			return "", resultError(result)
		}
		return renderToolResultContent(result.Value), nil
	}
	return e.fallback.CallTool(ctx, name, arguments)
}

type inlineToolCallEnvelope struct {
	ToolCall *ToolCall `json:"tool_call"`
}

func parseInlineToolCall(content string) (ToolCall, bool, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" || !strings.Contains(trimmed, "tool_call") {
		return ToolCall{}, false, nil
	}

	var envelope inlineToolCallEnvelope
	if result := core.JSONUnmarshal([]byte(trimmed), &envelope); !result.OK {
		return ToolCall{}, false, resultError(result)
	}
	if envelope.ToolCall == nil {
		return ToolCall{}, false, nil
	}
	call := *envelope.ToolCall
	call.Name = strings.TrimSpace(call.Name)
	if call.Arguments == nil {
		call.Arguments = map[string]any{}
	}
	return call, true, nil
}

func toolArgumentsFromOptions(opts core.Options) map[string]any {
	if value := opts.Get("arguments"); value.OK {
		if arguments, ok := value.Value.(map[string]any); ok {
			return cloneArguments(arguments)
		}
		var arguments map[string]any
		if result := core.JSONUnmarshal([]byte(jsonString(value.Value)), &arguments); result.OK {
			return arguments
		}
	}

	arguments := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		arguments[item.Key] = item.Value
	}
	return arguments
}

func cloneArguments(arguments map[string]any) map[string]any {
	if arguments == nil {
		return map[string]any{}
	}
	clone := make(map[string]any, len(arguments))
	for key, value := range arguments {
		clone[key] = value
	}
	return clone
}

func renderToolResultContent(result any) string {
	switch typed := result.(type) {
	case nil:
		return ""
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		return jsonString(typed)
	}
}

func jsonString(value any) string {
	result := core.JSONMarshal(value)
	if !result.OK {
		return "{}"
	}
	if data, ok := result.Value.([]byte); ok {
		return string(data)
	}
	return "{}"
}

func resultError(result core.Result) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return core.E("chat.tool_call", "unexpected result type", nil)
}
