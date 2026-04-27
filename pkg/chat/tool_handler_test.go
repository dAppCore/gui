package chat

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	core "dappco.re/go/core"
	guimcp "dappco.re/go/gui/pkg/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type strictToolExecutor struct {
	mu    sync.Mutex
	calls []ToolCall
}

func (m *strictToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"window_count": map[string]any{"type": "integer"},
			},
		},
	}}
}

func (m *strictToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (m *strictToolExecutor) CallTool(_ context.Context, name string, arguments map[string]any) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, ToolCall{Name: name, Arguments: arguments})
	if name != "layout_suggest" {
		return "", core.E("test.tool", "unknown tool: "+name, nil)
	}
	return `{"mode":"left-right"}`, nil
}

func (m *strictToolExecutor) Calls() []ToolCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]ToolCall(nil), m.calls...)
}

type completionRecorder struct {
	mu        sync.Mutex
	requests  []openAIRequest
	responses [][]string
}

func (r *completionRecorder) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	var completion openAIRequest
	if result := core.JSONUnmarshal(body, &completion); !result.OK {
		http.Error(w, renderToolResultContent(result.Value), http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	r.requests = append(r.requests, completion)
	index := len(r.requests) - 1
	r.mu.Unlock()

	if index >= len(r.responses) {
		http.Error(w, "unexpected completion request", http.StatusInternalServerError)
		return
	}
	writeSSE(w, r.responses[index]...)
}

func (r *completionRecorder) Requests() []openAIRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]openAIRequest(nil), r.requests...)
}

func TestToolCallHandler_Good_ServiceDispatchesInlineToolCall(t *testing.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{
		{
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"layout_suggest\",\"arguments\":{\"window_count\":2}}}"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
		{
			`{"id":"chatcmpl-2","choices":[{"delta":{"content":"Layout applied"}}]}`,
			`{"id":"chatcmpl-2","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
	}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Arrange this workspace"},
	))
	require.True(t, send.OK)

	calls := executor.Calls()
	require.Len(t, calls, 1)
	assert.Equal(t, "layout_suggest", calls[0].Name)
	assert.Equal(t, float64(2), calls[0].Arguments["window_count"])

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	require.Len(t, history, 4)
	assert.Equal(t, "assistant", history[1].Role)
	require.Len(t, history[1].ToolCalls, 1)
	assert.Equal(t, "tool", history[2].Role)
	assert.Contains(t, history[2].Content, "left-right")
	assert.Equal(t, "Layout applied", history[3].Content)

	requests := recorder.Requests()
	require.Len(t, requests, 2)
	require.NotEmpty(t, requests[0].Messages)
	systemPrompt, ok := requests[0].Messages[0].Content.(string)
	require.True(t, ok)
	assert.True(t, strings.HasPrefix(systemPrompt, "Available MCP tools:"))
	assert.Contains(t, systemPrompt, "layout_suggest")
	assert.Contains(t, systemPrompt, "You are a helpful assistant.")
}

func TestToolCallHandler_Bad_UnknownToolErrorAppearsInConversation(t *testing.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{
		{
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"missing_tool\",\"arguments\":{}}}"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
		{
			`{"id":"chatcmpl-2","choices":[{"delta":{"content":"Could not run that tool"}}]}`,
			`{"id":"chatcmpl-2","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		},
	}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Use the missing tool"},
	))
	require.True(t, send.OK)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	require.Len(t, history, 4)
	assert.Equal(t, "tool", history[2].Role)
	assert.Contains(t, history[2].Content, "missing_tool")
	assert.Equal(t, "Could not run that tool", history[3].Content)
}

func TestToolCallHandler_Ugly_MalformedInlineToolCallDoesNotDispatch(t *testing.T) {
	executor := &strictToolExecutor{}
	recorder := &completionRecorder{responses: [][]string{{
		`{"id":"chatcmpl-1","choices":[{"delta":{"content":"{\"tool_call\":{\"name\":\"layout_suggest\",\"arguments\":"}}]}`,
		`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
		`[DONE]`,
	}}}
	c := newChatCore(t, recorder.ServeHTTP, executor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Try malformed JSON"},
	))
	require.True(t, send.OK)

	assert.Empty(t, executor.Calls())
	assert.Len(t, recorder.Requests(), 1)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 0)
	require.Len(t, history, 2)
	assert.Equal(t, "assistant", history[1].Role)
	assert.Contains(t, history[1].Content, "tool_call")
	assert.Empty(t, history[1].ToolCalls)
}
