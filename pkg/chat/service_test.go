package chat

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	core "dappco.re/go/core"
	guimcp "dappco.re/go/gui/pkg/mcp"
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

func newChatCore(t *testing.T, handler http.HandlerFunc, toolExecutor ToolExecutor, optionFns ...func(*Options)) *core.Core {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	options := []func(*Options){
		func(o *Options) { o.APIURL = server.URL },
		func(o *Options) { o.StorePath = filepath.Join(t.TempDir(), "chat.db") },
		func(o *Options) { o.ToolExecutor = toolExecutor },
		func(o *Options) { o.Now = func() time.Time { return time.Unix(1_700_000_000, 0).UTC() } },
		func(o *Options) { o.ModelRoots = nil },
	}
	options = append(options, optionFns...)

	c := core.New(
		core.WithService(Register(options...)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return c
}

func createDiscoveredModelRoot(t *testing.T, name, architecture string) string {
	t.Helper()
	root := t.TempDir()
	modelDir := filepath.Join(root, name)
	require.NoError(t, os.MkdirAll(modelDir, 0o755))
	configJSON := `{"model_type":"` + architecture + `","quantization":{"bits":4,"group_size":32}}`
	require.NoError(t, os.WriteFile(filepath.Join(modelDir, "config.json"), []byte(configJSON), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(modelDir, "weights.safetensors"), []byte("fake"), 0o644))
	return root
}

func sequencedNow(start time.Time) func() time.Time {
	current := start.Add(-time.Second)
	return func() time.Time {
		current = current.Add(time.Second)
		return current
	}
}

func writeSSE(w http.ResponseWriter, payloads ...string) {
	w.Header().Set("Content-Type", "text/event-stream")
	for _, payload := range payloads {
		_, _ = io.WriteString(w, "data: "+payload+"\n\n")
	}
}

func latestConversation(t *testing.T, c *core.Core) Conversation {
	t.Helper()
	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)
	conversations, ok := result.Value.([]Conversation)
	require.True(t, ok)
	require.NotEmpty(t, conversations)
	return conversations[0]
}

func historyMessages(t *testing.T, c *core.Core, conversationID string, limit int) []Message {
	t.Helper()
	options := []core.Option{{
		Key:   "conversation_id",
		Value: conversationID,
	}}
	if limit > 0 {
		options = append(options, core.Option{Key: "limit", Value: limit})
	}
	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(options...))
	require.True(t, result.OK)
	messages, ok := result.Value.([]Message)
	require.True(t, ok)
	return messages
}

func TestActionSend_Good_ReturnsAssistantMessageID(t *testing.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Hello"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":" world"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hi"},
		core.Option{Key: "model", Value: "lemma"},
	))
	require.True(t, send.OK)

	messageID, ok := send.Value.(string)
	require.True(t, ok)
	require.NotEmpty(t, messageID)

	conv := latestConversation(t, c)
	require.Len(t, conv.Messages, 2)
	assert.Equal(t, messageID, conv.Messages[1].ID)
	assert.Equal(t, "lemma", conv.Model)
	assert.Equal(t, "Hello world", conv.Messages[1].Content)
}

func TestActionSend_Bad_RejectsEmptyMessage(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "message content is required")
}

func TestActionSend_Ugly_PropagatesUpstreamFailure(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "model unavailable", http.StatusBadGateway)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hi"},
	))
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "model unavailable")
}

func TestActionHistory_Good_HonoursLimit(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Alpha"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	send := c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "One"},
	))
	require.True(t, send.OK)

	conv := latestConversation(t, c)
	history := historyMessages(t, c, conv.ID, 1)
	require.Len(t, history, 1)
	assert.Equal(t, "assistant", history[0].Role)
	assert.Equal(t, "Alpha", history[0].Content)
}

func TestActionHistory_Bad_RequiresConversationID(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "conversation id is required")
}

func TestActionHistory_Ugly_UnknownConversationFails(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
}

func TestActionModels_Good_ReportsSizeAndStatus(t *testing.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	require.True(t, ok)
	require.Len(t, models, 1)
	assert.Equal(t, "lemma", models[0].Name)
	assert.Equal(t, int64(4), models[0].Size)
	assert.Equal(t, "active", models[0].Status)
}

func TestActionModels_Bad_ReturnsFallbackWhenNothingDiscovered(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	require.True(t, ok)
	require.Len(t, models, 1)
	assert.Equal(t, "default", models[0].Name)
	assert.Equal(t, "active", models[0].Status)
}

func TestActionModels_Ugly_ReflectsSelectedModelStatus(t *testing.T) {
	rootA := createDiscoveredModelRoot(t, "alpha", "gemma3")
	rootB := createDiscoveredModelRoot(t, "beta", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{rootA, rootB} })

	selected := c.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "beta"},
	))
	require.True(t, selected.OK)

	result := c.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)

	models, ok := result.Value.([]ModelEntry)
	require.True(t, ok)
	require.Len(t, models, 2)
	statusByName := map[string]string{}
	for _, model := range models {
		statusByName[model.Name] = model.Status
	}
	assert.Equal(t, "available", statusByName["alpha"])
	assert.Equal(t, "active", statusByName["beta"])
}

func TestActionSelectModel_Good_UpdatesConversationAndSettings(t *testing.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	created := c.Action("gui.chat.conversations.new").Run(context.Background(), core.NewOptions())
	require.True(t, created.OK)
	conv := created.Value.(Conversation)

	selected := c.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "lemma"},
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, selected.OK)

	settings := selected.Value.(ChatSettings)
	assert.Equal(t, "lemma", settings.DefaultModel)

	loaded := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, loaded.OK)
	assert.Equal(t, "lemma", loaded.Value.(Conversation).Model)
}

func TestActionSelectModel_Bad_RequiresModelName(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "model is required")
}

func TestActionSelectModel_Ugly_RejectsUnknownDiscoveredModel(t *testing.T) {
	modelRoot := createDiscoveredModelRoot(t, "lemma", "gemma3")
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.ModelRoots = []string{modelRoot} })

	result := c.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: "missing"},
	))
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "model is not available")
}

func TestActionConversationsList_Good_ReturnsNewestFirst(t *testing.T) {
	now := sequencedNow(time.Unix(1_700_000_000, 0).UTC())
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Ack"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{}, func(o *Options) { o.Now = now })

	require.True(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "First"},
	)).OK)
	require.True(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Second"},
	)).OK)

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)
	conversations := result.Value.([]Conversation)
	require.Len(t, conversations, 2)
	assert.Equal(t, "Second", conversations[0].Messages[0].Content)
	assert.Equal(t, "First", conversations[1].Messages[0].Content)
}

func TestActionConversationsList_Bad_EmptyStoreReturnsEmptySlice(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)
	conversations, ok := result.Value.([]Conversation)
	require.True(t, ok)
	assert.Empty(t, conversations)
}

func TestActionConversationsList_Ugly_IgnoresCorruptRows(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Ack"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	require.True(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Good"},
	)).OK)

	svc := core.MustServiceFor[*Service](c, "chat")
	require.NoError(t, svc.store.Set(conversationsGroup, "broken", "{"))

	result := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	require.True(t, result.OK)
	conversations := result.Value.([]Conversation)
	require.Len(t, conversations, 1)
	assert.Equal(t, "Good", conversations[0].Messages[0].Content)
}

func TestActionConversationsLoad_Good_ReturnsConversation(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Reply"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	require.True(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	)).OK)
	conv := latestConversation(t, c)

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, result.OK)
	loaded := result.Value.(Conversation)
	require.Len(t, loaded.Messages, 2)
	assert.Equal(t, "Reply", loaded.Messages[1].Content)
}

func TestActionConversationsLoad_Bad_RequiresConversationID(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "conversation id is required")
}

func TestActionConversationsLoad_Ugly_UnknownConversationFails(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.load").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
}

func TestActionConversationsDelete_Good_RemovesConversation(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w,
			`{"id":"chatcmpl-1","choices":[{"delta":{"content":"Reply"}}]}`,
			`{"id":"chatcmpl-1","choices":[{"finish_reason":"stop"}]}`,
			`[DONE]`,
		)
	}, &mockToolExecutor{})

	require.True(t, c.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "content", Value: "Hello"},
	)).OK)
	conv := latestConversation(t, c)

	deleted := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: conv.ID},
	))
	require.True(t, deleted.OK)
	assert.Equal(t, true, deleted.Value)

	listed := c.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	require.True(t, listed.OK)
	assert.Empty(t, listed.Value.([]Conversation))
}

func TestActionConversationsDelete_Bad_RequiresConversationID(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "conversation id is required")
}

func TestActionConversationsDelete_Ugly_IsIdempotentForMissingConversation(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "missing"},
	))
	require.True(t, result.OK)
	assert.Equal(t, true, result.Value)
}

func TestActionThinkingStart_Good_ReturnsActiveState(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	require.True(t, result.OK)
	state := result.Value.(ThinkingState)
	assert.True(t, state.Active)
	assert.False(t, state.StartedAt.IsZero())
}

func TestActionThinkingStart_Bad_RequiresConversationID(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "conversation id is required")
}

func TestActionThinkingStart_Ugly_RestartReplacesExistingState(t *testing.T) {
	now := sequencedNow(time.Unix(1_700_000_000, 0).UTC())
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{}, func(o *Options) { o.Now = now })

	first := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	second := c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	require.True(t, first.OK)
	require.True(t, second.OK)
	assert.True(t, second.Value.(ThinkingState).StartedAt.After(first.Value.(ThinkingState).StartedAt))
}

func TestActionThinkingStop_Good_ClearsThinkingState(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	require.True(t, c.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
		core.Option{Key: "started_at", Value: time.Unix(1_700_000_000, 0).UTC()},
	)).OK)

	stopped := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
		core.Option{Key: "duration_ms", Value: int64(25)},
	))
	require.True(t, stopped.OK)
	state := stopped.Value.(ThinkingState)
	assert.False(t, state.Active)
	assert.Equal(t, int64(25), state.DurationMS)
}

func TestActionThinkingStop_Bad_RequiresConversationID(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions())
	require.False(t, result.OK)
	require.Error(t, result.Value.(error))
	assert.Contains(t, result.Value.(error).Error(), "conversation id is required")
}

func TestActionThinkingStop_Ugly_AllowsStopWithoutStart(t *testing.T) {
	c := newChatCore(t, func(w http.ResponseWriter, _ *http.Request) {
		writeSSE(w, `[DONE]`)
	}, &mockToolExecutor{})

	result := c.Action("gui.chat.thinking.stop").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: "conv-1"},
	))
	require.True(t, result.OK)
	state := result.Value.(ThinkingState)
	assert.False(t, state.Active)
	assert.True(t, state.DurationMS >= 0)
}
