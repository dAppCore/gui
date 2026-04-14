package display

import (
	"context"
	"path/filepath"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectMode_Good(t *testing.T) {
	mode := DetectMode([]string{"core-gui", "--mode=worker"}, func(string) string { return "" })
	assert.Equal(t, ModeWorker, mode)

	mode = DetectMode(nil, func(string) string { return "manager" })
	assert.Equal(t, ModeManager, mode)

	mode = DetectMode(nil, func(string) string { return "bogus" })
	assert.Equal(t, ModeManager, mode)
}

func TestChatLifecycle_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	conv := convResult.(Conversation)
	require.NotEmpty(t, conv.ID)

	updatedAttachments, handled, err := c.PERFORM(TaskAttachImage{
		ConversationID: conv.ID,
		Attachment: ImageAttachment{
			Filename: "diagram.png",
			MimeType: "image/png",
			Data:     "ZmFrZQ==",
			Width:    640,
			Height:   480,
		},
	})
	require.NoError(t, err)
	require.True(t, handled)
	attachments := updatedAttachments.([]ImageAttachment)
	require.Len(t, attachments, 1)
	require.NotEmpty(t, attachments[0].ID)

	remainingResult, handled, err := c.PERFORM(TaskDetachImage{
		ConversationID: conv.ID,
		AttachmentID:   attachments[0].ID,
	})
	require.NoError(t, err)
	require.True(t, handled)
	assert.Empty(t, remainingResult.([]ImageAttachment))

	updatedAttachments, handled, err = c.PERFORM(TaskAttachImage{
		ConversationID: conv.ID,
		Attachment: ImageAttachment{
			Filename: "diagram.png",
			MimeType: "image/png",
			Data:     "ZmFrZQ==",
			Width:    640,
			Height:   480,
		},
	})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskThinkingStart{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskThinkingAppend{
		ConversationID: conv.ID,
		Content:        "Consider the local context first.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	sendResult, handled, err := c.PERFORM(TaskChatSend{
		ConversationID: conv.ID,
		Content:        "Explain local inference.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	updated := sendResult.(Conversation)
	require.Len(t, updated.Messages, 2)
	assert.Equal(t, "user", updated.Messages[0].Role)
	assert.Equal(t, "assistant", updated.Messages[1].Role)
	assert.Len(t, updated.Messages[0].Attachments, 1)
	if assert.NotNil(t, updated.Messages[1].Thinking) {
		assert.Contains(t, updated.Messages[1].Thinking.Content, "local context")
	}

	historyResult, handled, err := c.QUERY(QueryChatHistory{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	history := historyResult.([]ChatMessage)
	require.Len(t, history, 2)

	searchResult, handled, err := c.QUERY(QueryConversationsSearch{Query: "inference"})
	require.NoError(t, err)
	require.True(t, handled)
	require.Len(t, searchResult.([]Conversation), 1)

	renamedResult, handled, err := c.PERFORM(TaskConversationRename{
		ID:    conv.ID,
		Title: "Local inference notes",
	})
	require.NoError(t, err)
	require.True(t, handled)
	assert.Equal(t, "Local inference notes", renamedResult.(Conversation).Title)

	exportedResult, handled, err := c.QUERY(QueryConversationExport{ID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	exported := exportedResult.(string)
	assert.Contains(t, exported, "# Local inference notes")
	assert.Contains(t, exported, "## User")
	assert.Contains(t, exported, "diagram.png")

	settingsResult, handled, err := c.PERFORM(TaskChatSettingsSave{
		Settings: ChatSettings{
			Temperature:   0.7,
			TopP:          0.9,
			TopK:          40,
			MaxTokens:     1024,
			ContextWindow: 4096,
			SystemPrompt:  "Be concise.",
			DefaultModel:  "lemma",
		},
	})
	require.NoError(t, err)
	require.True(t, handled)
	assert.Equal(t, float32(0.7), settingsResult.(ChatSettings).Temperature)

	modelsResult, handled, err := c.PERFORM(TaskSelectModel{Model: "lemma"})
	require.NoError(t, err)
	require.True(t, handled)
	models := modelsResult.([]ModelEntry)
	assert.True(t, models[1].Loaded)

	require.NoError(t, svc.chat.persist(svc.configFile))
}

func TestChatStreamingLifecycle_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	conv := convResult.(Conversation)

	_, handled, err = c.PERFORM(TaskChatSend{
		ConversationID: conv.ID,
		Content:        "Stream the answer instead.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	startResult, handled, err := c.PERFORM(TaskChatStreamStart{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	started := startResult.(Conversation)
	require.Len(t, started.Messages, 2)
	assert.True(t, started.Messages[1].Streaming)
	assert.Empty(t, started.Messages[1].Content)

	_, handled, err = c.PERFORM(TaskThinkingStart{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskThinkingAppend{
		ConversationID: conv.ID,
		Content:        "Streaming through the local bridge.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	appendResult, handled, err := c.PERFORM(TaskChatStreamAppend{
		ConversationID: conv.ID,
		Content:        "Hello",
	})
	require.NoError(t, err)
	require.True(t, handled)
	appended := appendResult.(Conversation)
	require.Len(t, appended.Messages, 2)
	assert.Equal(t, "Hello", appended.Messages[1].Content)
	if assert.NotNil(t, appended.Messages[1].Thinking) {
		assert.Contains(t, appended.Messages[1].Thinking.Content, "local bridge")
	}

	appendResult, handled, err = c.PERFORM(TaskChatStreamAppend{
		ConversationID: conv.ID,
		Content:        " world",
	})
	require.NoError(t, err)
	require.True(t, handled)
	appended = appendResult.(Conversation)
	assert.Equal(t, "Hello world", appended.Messages[1].Content)

	finishResult, handled, err := c.PERFORM(TaskChatStreamFinish{
		ConversationID: conv.ID,
		FinishReason:   "stop",
	})
	require.NoError(t, err)
	require.True(t, handled)
	finished := finishResult.(Conversation)
	require.Len(t, finished.Messages, 2)
	assert.False(t, finished.Messages[1].Streaming)
	assert.Equal(t, "stop", finished.Messages[1].FinishReason)

	historyResult, handled, err := c.QUERY(QueryChatHistory{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	history := historyResult.([]ChatMessage)
	require.Len(t, history, 2)
	assert.Equal(t, "Hello world", history[1].Content)
	assert.False(t, history[1].Streaming)
}

func TestChatPersistence_Good(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gui.yaml")

	svc, err := NewService()
	require.NoError(t, err)
	svc.loadConfigFrom(path)

	conv := svc.chat.NewConversation()
	_, _, _, err = svc.chat.SendMessage(conv.ID, "Persist this conversation.")
	require.NoError(t, err)
	require.NoError(t, svc.chat.persist(svc.configFile))

	reloaded, err := NewService()
	require.NoError(t, err)
	reloaded.loadConfigFrom(path)

	restored, ok := reloaded.chat.Conversation(conv.ID)
	require.True(t, ok)
	require.Len(t, restored.Messages, 2)
	assert.Equal(t, "Persist this conversation.", restored.Messages[0].Content)
}

func TestResolveScheme_Good(t *testing.T) {
	c, err := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	svc := core.MustServiceFor[*Service](c, "display")
	conv := svc.chat.NewConversation()
	_, _, _, err = svc.chat.SendMessage(conv.ID, "Searchable store entry.")
	require.NoError(t, err)

	storeResponse, err := svc.ResolveScheme(context.Background(), "core://store?q=searchable")
	require.NoError(t, err)
	assert.Equal(t, "store", storeResponse.Path)
	results := storeResponse.Data["results"].([]StoreSearchResult)
	require.Len(t, results, 2)

	settingsResponse, err := svc.ResolveScheme(context.Background(), "core://settings")
	require.NoError(t, err)
	assert.Equal(t, "settings", settingsResponse.Path)
	assert.Contains(t, settingsResponse.Data, "settings")
	assert.Contains(t, settingsResponse.Data, "models")
}
