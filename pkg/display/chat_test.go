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

func TestChatConversationSave_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	result, handled, err := c.PERFORM(TaskConversationSave{
		Conversation: Conversation{
			Messages: []ChatMessage{
				{Role: "user", Content: "Imported conversation content that should title itself."},
				{Role: "assistant", Content: "Imported response."},
			},
		},
	})
	require.NoError(t, err)
	require.True(t, handled)

	conv := result.(Conversation)
	require.NotEmpty(t, conv.ID)
	require.Len(t, conv.Messages, 2)
	assert.Contains(t, conv.Title, "Imported conversation content")
	assert.True(t, len(conv.Title) >= len("Imported conversation content"))
	assert.NotEmpty(t, conv.Messages[0].ID)
	assert.NotEmpty(t, conv.Messages[1].ID)
	assert.False(t, conv.CreatedAt.IsZero())
	assert.False(t, conv.UpdatedAt.IsZero())

	loadedResult, handled, err := c.QUERY(QueryConversationGet{ID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	loaded := loadedResult.(Conversation)
	assert.Equal(t, conv.ID, loaded.ID)
	assert.Equal(t, conv.Title, loaded.Title)
}

func TestChatSelectedModelAppliesMidConversation_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	conv := convResult.(Conversation)
	assert.Equal(t, "lemer", conv.Model)

	_, handled, err = c.PERFORM(TaskChatSend{
		ConversationID: conv.ID,
		Content:        "Use the default model first.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskSelectModel{Model: "lemma"})
	require.NoError(t, err)
	require.True(t, handled)

	result, handled, err := c.PERFORM(TaskChatSend{
		ConversationID: conv.ID,
		Content:        "Now switch to the new model.",
	})
	require.NoError(t, err)
	require.True(t, handled)

	updated := result.(Conversation)
	require.Len(t, updated.Messages, 4)
	assert.Equal(t, "lemma", updated.Model)
	assert.Contains(t, updated.Messages[3].Content, "lemma")
}

func TestChatAttachmentValidation_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	conv := convResult.(Conversation)

	_, handled, err = c.PERFORM(TaskSelectModel{Model: "lemmy"})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskAttachImage{
		ConversationID: conv.ID,
		Attachment: ImageAttachment{
			Filename: "photo.png",
			MimeType: "image/png",
			Data:     "ZmFrZQ==",
		},
	})
	require.Error(t, err)
	assert.True(t, handled)
	assert.Contains(t, err.Error(), "does not support vision")

	_, handled, err = c.PERFORM(TaskSelectModel{Model: "lemer"})
	require.NoError(t, err)
	require.True(t, handled)

	_, handled, err = c.PERFORM(TaskAttachImage{
		ConversationID: conv.ID,
		Attachment: ImageAttachment{
			Filename: "photo.svg",
			MimeType: "image/svg+xml",
			Data:     "ZmFrZQ==",
		},
	})
	require.Error(t, err)
	assert.True(t, handled)
	assert.Contains(t, err.Error(), "unsupported image type")
}

func TestChatQueuedImagesQuery_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	conv := convResult.(Conversation)

	updatedAttachments, handled, err := c.PERFORM(TaskAttachImage{
		ConversationID: conv.ID,
		Attachment: ImageAttachment{
			Filename: "diagram.webp",
			MimeType: "image/webp",
			Data:     "ZmFrZQ==",
		},
	})
	require.NoError(t, err)
	require.True(t, handled)
	require.Len(t, updatedAttachments.([]ImageAttachment), 1)

	result, handled, err := c.QUERY(QueryQueuedImages{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	attachments := result.([]ImageAttachment)
	require.Len(t, attachments, 1)
	assert.Equal(t, "diagram.webp", attachments[0].Filename)
}

func TestChatTaskEvents_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)
	svc.events = &WSEventManager{eventBuffer: make(chan Event, 8)}

	convResult, handled, err := c.PERFORM(TaskConversationNew{})
	require.NoError(t, err)
	require.True(t, handled)
	createdEvent := <-svc.events.eventBuffer
	assert.Equal(t, EventType("chat.conversation.created"), createdEvent.Type)

	conv := convResult.(Conversation)
	_, handled, err = c.PERFORM(TaskThinkingStart{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	startEvent := <-svc.events.eventBuffer
	assert.Equal(t, EventType("chat.thinking.start"), startEvent.Type)

	_, handled, err = c.PERFORM(TaskThinkingAppend{
		ConversationID: conv.ID,
		Content:        "Inspecting the prompt.",
	})
	require.NoError(t, err)
	require.True(t, handled)
	appendEvent := <-svc.events.eventBuffer
	assert.Equal(t, EventType("chat.thinking.delta"), appendEvent.Type)
	assert.Equal(t, "Inspecting the prompt.", appendEvent.Data["delta"])

	_, handled, err = c.PERFORM(TaskThinkingEnd{ConversationID: conv.ID})
	require.NoError(t, err)
	require.True(t, handled)
	endEvent := <-svc.events.eventBuffer
	assert.Equal(t, EventType("chat.thinking.end"), endEvent.Type)
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

	modelsResponse, err := svc.ResolveScheme(context.Background(), "core://models")
	require.NoError(t, err)
	assert.Equal(t, "models", modelsResponse.Path)
	assert.Contains(t, modelsResponse.Data, "selected_model")
	assert.Contains(t, modelsResponse.Data, "models")
}

func TestConversationSearchIncludesThinkingToolsAndAttachments_Good(t *testing.T) {
	_, c := newTestDisplayService(t)

	result, handled, err := c.PERFORM(TaskConversationSave{
		Conversation: Conversation{
			Title: "Invoice debugger",
			Messages: []ChatMessage{
				{
					Role: "user",
					Attachments: []ImageAttachment{
						{
							Filename: "receipt.png",
							MimeType: "image/png",
							Data:     "ZmFrZQ==",
						},
					},
				},
				{
					Role: "assistant",
					Thinking: &ThinkingState{
						Content: "Tracing invoice 42 through the local store.",
					},
					ToolCalls: []ToolInvocation{
						{
							Call: ToolCall{
								Name: "store.lookup",
								Arguments: map[string]any{
									"q": "invoice-42",
								},
							},
							Result: ToolResult{
								Content: "invoice 42 marked paid",
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.True(t, handled)
	conv := result.(Conversation)
	require.NotEmpty(t, conv.ID)

	for _, query := range []string{
		"receipt.png",
		"image/png",
		"tracing invoice 42",
		"store.lookup",
		"invoice-42",
		"marked paid",
	} {
		searchResult, handled, err := c.QUERY(QueryConversationsSearch{Query: query})
		require.NoError(t, err)
		require.True(t, handled)
		matches := searchResult.([]Conversation)
		require.Len(t, matches, 1, query)
		assert.Equal(t, conv.ID, matches[0].ID, query)
	}
}

func TestRouteQueriesAndStoreGroups_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)

	_, handled, err := c.PERFORM(TaskConversationSave{
		Conversation: Conversation{
			Title: "Invoice notes",
			Messages: []ChatMessage{
				{
					Role: "assistant",
					ToolCalls: []ToolInvocation{
						{
							Call: ToolCall{
								Name: "store.lookup",
								Arguments: map[string]any{
									"q": "invoice",
								},
							},
							Result: ToolResult{
								Content: "invoice status is paid",
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.True(t, handled)

	require.NoError(t, svc.SaveBrowserStorageState("https://app.example.com", OriginStorageState{
		LocalStorage: map[string]string{
			"invoice_status": "paid",
		},
	}))

	settingsResult, handled, err := c.QUERY(QueryRouteSettings{})
	require.NoError(t, err)
	require.True(t, handled)
	settingsData := settingsResult.(map[string]any)
	assert.Contains(t, settingsData, "settings")
	assert.Contains(t, settingsData, "models")

	modelsResult, handled, err := c.QUERY(QueryRouteModels{})
	require.NoError(t, err)
	require.True(t, handled)
	modelsData := modelsResult.(map[string]any)
	assert.Contains(t, modelsData, "selected_model")
	assert.Contains(t, modelsData, "models")

	resolveResult, handled, err := c.QUERY(QueryRouteResolve{URL: "core://models"})
	require.NoError(t, err)
	require.True(t, handled)
	resolved := resolveResult.(SchemeResponse)
	assert.Equal(t, "models", resolved.Path)
	assert.Contains(t, resolved.Data, "selected_model")

	storeResult, handled, err := c.QUERY(QueryRouteStore{Query: "invoice"})
	require.NoError(t, err)
	require.True(t, handled)
	storeData := storeResult.(map[string]any)

	results := storeData["results"].([]StoreSearchResult)
	require.Len(t, results, 2)

	groups := storeData["groups"].([]StoreSearchGroup)
	require.Len(t, groups, 2)

	origins := []string{groups[0].Origin, groups[1].Origin}
	assert.ElementsMatch(t, []string{"chat://conversations", "https://app.example.com"}, origins)

	var snippets []string
	for _, result := range results {
		snippets = append(snippets, result.Snippet)
	}
	assert.Contains(t, snippets, "invoice status is paid")
}
