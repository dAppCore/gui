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

	_, handled, err = c.PERFORM(TaskAttachImage{
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
