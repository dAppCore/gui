package chat

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	core "dappco.re/go/core"
	guimcp "forge.lthn.ai/core/gui/pkg/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockToolExecutor struct {
	calls []ToolCall
}

func (m *mockToolExecutor) Manifest() []guimcp.ToolDescriptor {
	return []guimcp.ToolDescriptor{{
		Name:        "layout_suggest",
		Description: "Suggest a layout",
		InputSchema: map[string]any{"type": "object"},
	}}
}

func (m *mockToolExecutor) ManifestText() string {
	return "Available MCP tools:\n- layout_suggest: Suggest a layout"
}

func (m *mockToolExecutor) CallTool(_ context.Context, name string, arguments map[string]any) (string, error) {
	m.calls = append(m.calls, ToolCall{Name: name, Arguments: arguments})
	return `{"mode":"left-right"}`, nil
}

func newChatCore(t *testing.T, handler http.HandlerFunc, toolExecutor ToolExecutor) *core.Core {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c := core.New(
		core.WithService(Register(
			func(o *Options) { o.APIURL = server.URL },
			func(o *Options) { o.StorePath = filepath.Join(t.TempDir(), "chat.db") },
			func(o *Options) { o.ToolExecutor = toolExecutor },
			func(o *Options) { o.Now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() } },
			func(o *Options) { o.ModelRoots = nil },
		)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return c
}

func TestService_Good_SendAndHistory(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"delta\":{\"content\":\" world\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}, &mockToolExecutor{})

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hi"},
	))
	require.True(t, send.OK)

	conv := send.Value.(Conversation)
	require.Len(t, conv.Messages, 2)
	assert.Equal(t, "user", conv.Messages[0].Role)
	assert.Equal(t, "assistant", conv.Messages[1].Role)
	assert.Equal(t, "Hello world", conv.Messages[1].Content)

	history := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, history.OK)
	assert.Equal(t, conv.ID, history.Value.(Conversation).ID)

	queryHistory := c.QUERY(QueryHistory{ConversationID: conv.ID})
	require.True(t, queryHistory.OK)
	assert.Equal(t, conv.ID, queryHistory.Value.(Conversation).ID)
}

func TestService_Good_ToolCallRoundTrip(t *testing.T) {
	toolExecutor := &mockToolExecutor{}
	requests := 0
	c := newChatCore(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "text/event-stream")
		if requests == 1 {
			_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call-1\",\"function\":{\"name\":\"layout_suggest\",\"arguments\":\"{\\\"window_count\\\":2}\"}}]}}]}\n\n")
			_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"finish_reason\":\"tool_calls\"}]}\n\n")
			_, _ = io.WriteString(w, "data: [DONE]\n\n")
			return
		}
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-2\",\"choices\":[{\"delta\":{\"content\":\"Use a left-right split.\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-2\",\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}, toolExecutor)

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Arrange these windows"},
	))
	require.True(t, send.OK)

	conv := send.Value.(Conversation)
	require.GreaterOrEqual(t, len(conv.Messages), 4)
	require.Len(t, toolExecutor.calls, 1)
	assert.Equal(t, "layout_suggest", toolExecutor.calls[0].Name)
	assert.Equal(t, 2.0, toolExecutor.calls[0].Arguments["window_count"])
	assert.Equal(t, "tool", conv.Messages[2].Role)
	assert.True(t, strings.Contains(conv.Messages[len(conv.Messages)-1].Content, "left-right"))
}

func TestService_Good_SelectModelUpdatesConversation(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"id\":\"chatcmpl-1\",\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}, &mockToolExecutor{})

	created := c.Action("gui.chat.conversations.new").Run(context.Background(), core.NewOptions())
	require.True(t, created.OK)
	conv := created.Value.(Conversation)

	selected := c.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "lemma"},
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, selected.OK)

	updated := c.QUERY(QueryConversationGet{ConversationID: conv.ID})
	require.True(t, updated.OK)
	assert.Equal(t, "lemma", updated.Value.(Conversation).Model)
}

func TestService_Good_SettingsDefaults(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}, &mockToolExecutor{})

	result := c.QUERY(QuerySettingsDefaults{})
	require.True(t, result.OK)
	assert.Equal(t, DefaultSettings(), result.Value.(ChatSettings))

	actionResult := c.Action("gui.chat.settings.defaults").Run(context.Background(), core.Options{})
	require.True(t, actionResult.OK)
	assert.Equal(t, DefaultSettings(), actionResult.Value.(ChatSettings))
}
