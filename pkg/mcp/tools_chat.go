// pkg/mcp/tools_chat.go
package mcp

import (
	"context"
	"time"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ChatMessage struct {
	ID           string                `json:"id"`
	Role         string                `json:"role"`
	Content      string                `json:"content"`
	CreatedAt    time.Time             `json:"created_at"`
	Model        string                `json:"model,omitempty"`
	Attachments  []ChatImageAttachment `json:"attachments,omitempty"`
	Thinking     *ChatThinkingState    `json:"thinking,omitempty"`
	ToolCalls    []ChatToolCall        `json:"tool_calls,omitempty"`
	ToolResults  []ChatToolResult      `json:"tool_results,omitempty"`
	ToolCallID   string                `json:"tool_call_id,omitempty"`
	FinishReason string                `json:"finish_reason,omitempty"`
}

type ChatModelEntry struct {
	Name           string `json:"name"`
	Architecture   string `json:"architecture"`
	QuantBits      int    `json:"quant_bits"`
	SizeBytes      int64  `json:"size_bytes"`
	Loaded         bool   `json:"loaded"`
	Backend        string `json:"backend"`
	SupportsVision bool   `json:"supports_vision"`
}

type ChatSettings struct {
	Temperature   float32 `json:"temperature"`
	TopP          float32 `json:"top_p"`
	TopK          int     `json:"top_k"`
	MaxTokens     int     `json:"max_tokens"`
	ContextWindow int     `json:"context_window"`
	SystemPrompt  string  `json:"system_prompt"`
	DefaultModel  string  `json:"default_model"`
}

type ChatConversation struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Model     string        `json:"model"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Messages  []ChatMessage `json:"messages"`
	Settings  *ChatSettings `json:"settings,omitempty"`
}

type ChatConversationSummary struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Model        string    `json:"model"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
}

type ChatThinkingState struct {
	Active     bool      `json:"active"`
	Content    string    `json:"content"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	EndedAt    time.Time `json:"ended_at,omitempty"`
	DurationMS int64     `json:"duration_ms,omitempty"`
}

type ChatToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ChatToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

type ChatImageAttachment struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

// --- chat_send ---

type ChatSendInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}

type ChatSendOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatSend(_ context.Context, _ *mcp.CallToolRequest, input ChatSendInput) (*mcp.CallToolResult, ChatSendOutput, error) {
	result := s.core.Action("gui.chat.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "content", Value: input.Content},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSendOutput{}, err
		}
		return nil, ChatSendOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatSendOutput{}, err
	}
	return nil, ChatSendOutput{Conversation: conversation}, nil
}

// --- chat_clear ---

type ChatClearInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatClearOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatClear(_ context.Context, _ *mcp.CallToolRequest, input ChatClearInput) (*mcp.CallToolResult, ChatClearOutput, error) {
	result := s.core.Action("gui.chat.clear").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatClearOutput{}, err
		}
		return nil, ChatClearOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatClearOutput{}, err
	}
	return nil, ChatClearOutput{Conversation: conversation}, nil
}

// --- chat_history ---

type ChatHistoryInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatHistoryOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatHistory(_ context.Context, _ *mcp.CallToolRequest, input ChatHistoryInput) (*mcp.CallToolResult, ChatHistoryOutput, error) {
	result := s.core.Action("gui.chat.history").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatHistoryOutput{}, err
		}
		return nil, ChatHistoryOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatHistoryOutput{}, err
	}
	return nil, ChatHistoryOutput{Conversation: conversation}, nil
}

// --- chat_models ---

type ChatModelsInput struct{}

type ChatModelsOutput struct {
	Models []ChatModelEntry `json:"models"`
}

func (s *Subsystem) chatModels(_ context.Context, _ *mcp.CallToolRequest, _ ChatModelsInput) (*mcp.CallToolResult, ChatModelsOutput, error) {
	result := s.core.Action("gui.chat.models").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatModelsOutput{}, err
		}
		return nil, ChatModelsOutput{}, nil
	}
	models, err := decodeChatValue[[]ChatModelEntry](result.Value)
	if err != nil {
		return nil, ChatModelsOutput{}, err
	}
	return nil, ChatModelsOutput{Models: models}, nil
}

// --- chat_select_model ---

type ChatSelectModelInput struct {
	Model          string `json:"model"`
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatSelectModelOutput struct {
	Settings ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSelectModel(_ context.Context, _ *mcp.CallToolRequest, input ChatSelectModelInput) (*mcp.CallToolResult, ChatSelectModelOutput, error) {
	result := s.core.Action("gui.chat.selectModel").Run(context.Background(), core.NewOptions(
		core.Option{Key: "model", Value: input.Model},
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSelectModelOutput{}, err
		}
		return nil, ChatSelectModelOutput{}, nil
	}
	settings, err := decodeChatValue[ChatSettings](result.Value)
	if err != nil {
		return nil, ChatSelectModelOutput{}, err
	}
	return nil, ChatSelectModelOutput{Settings: settings}, nil
}

// --- chat_settings_save ---

type ChatSettingsSaveInput ChatSettings

type ChatSettingsSaveOutput struct {
	Settings ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsSave(_ context.Context, _ *mcp.CallToolRequest, input ChatSettingsSaveInput) (*mcp.CallToolResult, ChatSettingsSaveOutput, error) {
	result := s.core.Action("gui.chat.settings.save").Run(context.Background(), core.NewOptions(
		core.Option{Key: "temperature", Value: input.Temperature},
		core.Option{Key: "top_p", Value: input.TopP},
		core.Option{Key: "top_k", Value: input.TopK},
		core.Option{Key: "max_tokens", Value: input.MaxTokens},
		core.Option{Key: "context_window", Value: input.ContextWindow},
		core.Option{Key: "system_prompt", Value: input.SystemPrompt},
		core.Option{Key: "default_model", Value: input.DefaultModel},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSettingsSaveOutput{}, err
		}
		return nil, ChatSettingsSaveOutput{}, nil
	}
	settings, err := decodeChatValue[ChatSettings](result.Value)
	if err != nil {
		return nil, ChatSettingsSaveOutput{}, err
	}
	return nil, ChatSettingsSaveOutput{Settings: settings}, nil
}

// --- chat_settings_load ---

type ChatSettingsLoadInput struct{}

type ChatSettingsLoadOutput struct {
	Settings ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsLoad(_ context.Context, _ *mcp.CallToolRequest, _ ChatSettingsLoadInput) (*mcp.CallToolResult, ChatSettingsLoadOutput, error) {
	result := s.core.Action("gui.chat.settings.load").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSettingsLoadOutput{}, err
		}
		return nil, ChatSettingsLoadOutput{}, nil
	}
	settings, err := decodeChatValue[ChatSettings](result.Value)
	if err != nil {
		return nil, ChatSettingsLoadOutput{}, err
	}
	return nil, ChatSettingsLoadOutput{Settings: settings}, nil
}

// --- chat_settings_defaults ---

type ChatSettingsDefaultsInput struct{}

type ChatSettingsDefaultsOutput struct {
	Settings ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsDefaults(_ context.Context, _ *mcp.CallToolRequest, _ ChatSettingsDefaultsInput) (*mcp.CallToolResult, ChatSettingsDefaultsOutput, error) {
	result := s.core.Action("gui.chat.settings.defaults").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSettingsDefaultsOutput{}, err
		}
		return nil, ChatSettingsDefaultsOutput{}, nil
	}
	settings, err := decodeChatValue[ChatSettings](result.Value)
	if err != nil {
		return nil, ChatSettingsDefaultsOutput{}, err
	}
	return nil, ChatSettingsDefaultsOutput{Settings: settings}, nil
}

// --- chat_settings_reset ---

type ChatSettingsResetInput struct{}

type ChatSettingsResetOutput struct {
	Settings ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsReset(_ context.Context, _ *mcp.CallToolRequest, _ ChatSettingsResetInput) (*mcp.CallToolResult, ChatSettingsResetOutput, error) {
	result := s.core.Action("gui.chat.settings.reset").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatSettingsResetOutput{}, err
		}
		return nil, ChatSettingsResetOutput{}, nil
	}
	settings, err := decodeChatValue[ChatSettings](result.Value)
	if err != nil {
		return nil, ChatSettingsResetOutput{}, err
	}
	return nil, ChatSettingsResetOutput{Settings: settings}, nil
}

// --- chat_conversations_list ---

type ChatConversationsListInput struct{}

type ChatConversationsListOutput struct {
	Conversations []ChatConversationSummary `json:"conversations"`
}

func (s *Subsystem) chatConversationsList(_ context.Context, _ *mcp.CallToolRequest, _ ChatConversationsListInput) (*mcp.CallToolResult, ChatConversationsListOutput, error) {
	result := s.core.Action("gui.chat.conversations.list").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsListOutput{}, err
		}
		return nil, ChatConversationsListOutput{}, nil
	}
	conversations, err := decodeChatValue[[]ChatConversationSummary](result.Value)
	if err != nil {
		return nil, ChatConversationsListOutput{}, err
	}
	return nil, ChatConversationsListOutput{Conversations: conversations}, nil
}

// --- chat_conversations_get ---

type ChatConversationsGetInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatConversationsGetOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatConversationsGet(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsGetInput) (*mcp.CallToolResult, ChatConversationsGetOutput, error) {
	result := s.core.Action("gui.chat.conversations.get").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsGetOutput{}, err
		}
		return nil, ChatConversationsGetOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatConversationsGetOutput{}, err
	}
	return nil, ChatConversationsGetOutput{Conversation: conversation}, nil
}

// --- chat_conversations_delete ---

type ChatConversationsDeleteInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatConversationsDeleteOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) chatConversationsDelete(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsDeleteInput) (*mcp.CallToolResult, ChatConversationsDeleteOutput, error) {
	result := s.core.Action("gui.chat.conversations.delete").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsDeleteOutput{}, err
		}
		return nil, ChatConversationsDeleteOutput{}, nil
	}
	return nil, ChatConversationsDeleteOutput{Success: true}, nil
}

// --- chat_conversations_search ---

type ChatConversationsSearchInput struct {
	Query string `json:"q"`
}

type ChatConversationsSearchOutput struct {
	Conversations []ChatConversationSummary `json:"conversations"`
}

func (s *Subsystem) chatConversationsSearch(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsSearchInput) (*mcp.CallToolResult, ChatConversationsSearchOutput, error) {
	result := s.core.Action("gui.chat.conversations.search").Run(context.Background(), core.NewOptions(
		core.Option{Key: "q", Value: input.Query},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsSearchOutput{}, err
		}
		return nil, ChatConversationsSearchOutput{}, nil
	}
	conversations, err := decodeChatValue[[]ChatConversationSummary](result.Value)
	if err != nil {
		return nil, ChatConversationsSearchOutput{}, err
	}
	return nil, ChatConversationsSearchOutput{Conversations: conversations}, nil
}

// --- chat_conversations_new ---

type ChatConversationsNewInput struct{}

type ChatConversationsNewOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatConversationsNew(_ context.Context, _ *mcp.CallToolRequest, _ ChatConversationsNewInput) (*mcp.CallToolResult, ChatConversationsNewOutput, error) {
	result := s.core.Action("gui.chat.conversations.new").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsNewOutput{}, err
		}
		return nil, ChatConversationsNewOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatConversationsNewOutput{}, err
	}
	return nil, ChatConversationsNewOutput{Conversation: conversation}, nil
}

// --- chat_conversation_save ---

type ChatConversationSaveInput ChatConversation

type ChatConversationSaveOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatConversationSave(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationSaveInput) (*mcp.CallToolResult, ChatConversationSaveOutput, error) {
	result := s.core.Action("gui.chat.conversation.save").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: input.ID},
		core.Option{Key: "title", Value: input.Title},
		core.Option{Key: "model", Value: input.Model},
		core.Option{Key: "created_at", Value: input.CreatedAt},
		core.Option{Key: "updated_at", Value: input.UpdatedAt},
		core.Option{Key: "messages", Value: input.Messages},
		core.Option{Key: "settings", Value: input.Settings},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationSaveOutput{}, err
		}
		return nil, ChatConversationSaveOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatConversationSaveOutput{}, err
	}
	return nil, ChatConversationSaveOutput{Conversation: conversation}, nil
}

// --- chat_conversations_rename ---

type ChatConversationsRenameInput struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type ChatConversationsRenameOutput struct {
	Conversation ChatConversation `json:"conversation"`
}

func (s *Subsystem) chatConversationsRename(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsRenameInput) (*mcp.CallToolResult, ChatConversationsRenameOutput, error) {
	result := s.core.Action("gui.chat.conversations.rename").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: input.ID},
		core.Option{Key: "title", Value: input.Title},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsRenameOutput{}, err
		}
		return nil, ChatConversationsRenameOutput{}, nil
	}
	conversation, err := decodeChatValue[ChatConversation](result.Value)
	if err != nil {
		return nil, ChatConversationsRenameOutput{}, err
	}
	return nil, ChatConversationsRenameOutput{Conversation: conversation}, nil
}

// --- chat_conversations_export ---

type ChatConversationsExportInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type ChatConversationsExportOutput struct {
	Markdown string `json:"markdown"`
}

func (s *Subsystem) chatConversationsExport(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsExportInput) (*mcp.CallToolResult, ChatConversationsExportOutput, error) {
	result := s.core.Action("gui.chat.conversations.export").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "id", Value: input.ID},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatConversationsExportOutput{}, err
		}
		return nil, ChatConversationsExportOutput{}, nil
	}
	markdown, ok := result.Value.(string)
	if !ok {
		return nil, ChatConversationsExportOutput{}, coreerr.E("mcp.chatConversationsExport", "unexpected result type", nil)
	}
	return nil, ChatConversationsExportOutput{Markdown: markdown}, nil
}

// --- chat_attach_image ---

type ChatAttachImageInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Filename       string `json:"filename"`
	MimeType       string `json:"mime_type"`
	Data           string `json:"data"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}

type ChatAttachImageOutput struct {
	Attachment ChatImageAttachment `json:"attachment"`
}

func (s *Subsystem) chatAttachImage(_ context.Context, _ *mcp.CallToolRequest, input ChatAttachImageInput) (*mcp.CallToolResult, ChatAttachImageOutput, error) {
	result := s.core.Action("gui.chat.attachImage").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "filename", Value: input.Filename},
		core.Option{Key: "mime_type", Value: input.MimeType},
		core.Option{Key: "data", Value: input.Data},
		core.Option{Key: "width", Value: input.Width},
		core.Option{Key: "height", Value: input.Height},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatAttachImageOutput{}, err
		}
		return nil, ChatAttachImageOutput{}, nil
	}
	attachment, err := decodeChatValue[ChatImageAttachment](result.Value)
	if err != nil {
		return nil, ChatAttachImageOutput{}, err
	}
	return nil, ChatAttachImageOutput{Attachment: attachment}, nil
}

// --- chat_remove_image ---

type ChatRemoveImageInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Index          int    `json:"index"`
}

type ChatRemoveImageOutput struct {
	Attachment ChatImageAttachment `json:"attachment"`
}

func (s *Subsystem) chatRemoveImage(_ context.Context, _ *mcp.CallToolRequest, input ChatRemoveImageInput) (*mcp.CallToolResult, ChatRemoveImageOutput, error) {
	result := s.core.Action("gui.chat.removeImage").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "index", Value: input.Index},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatRemoveImageOutput{}, err
		}
		return nil, ChatRemoveImageOutput{}, nil
	}
	attachment, err := decodeChatValue[ChatImageAttachment](result.Value)
	if err != nil {
		return nil, ChatRemoveImageOutput{}, err
	}
	return nil, ChatRemoveImageOutput{Attachment: attachment}, nil
}

// --- chat_thinking_start ---

type ChatThinkingStartInput struct {
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
}

type ChatThinkingStartOutput struct {
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id,omitempty"`
	Content        string    `json:"content,omitempty"`
	StartedAt      time.Time `json:"started_at"`
	DurationMS     int64     `json:"duration_ms,omitempty"`
}

func (s *Subsystem) chatThinkingStart(_ context.Context, _ *mcp.CallToolRequest, input ChatThinkingStartInput) (*mcp.CallToolResult, ChatThinkingStartOutput, error) {
	result := s.core.Action("gui.chat.thinking.start").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "message_id", Value: input.MessageID},
		core.Option{Key: "started_at", Value: input.StartedAt},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatThinkingStartOutput{}, err
		}
		return nil, ChatThinkingStartOutput{}, nil
	}
	output, err := decodeChatValue[ChatThinkingStartOutput](result.Value)
	if err != nil {
		return nil, ChatThinkingStartOutput{}, err
	}
	return nil, output, nil
}

// --- chat_thinking_append ---

type ChatThinkingAppendInput struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id,omitempty"`
	Content        string `json:"content"`
}

type ChatThinkingAppendOutput struct {
	Content string `json:"content"`
}

func (s *Subsystem) chatThinkingAppend(_ context.Context, _ *mcp.CallToolRequest, input ChatThinkingAppendInput) (*mcp.CallToolResult, ChatThinkingAppendOutput, error) {
	result := s.core.Action("gui.chat.thinking.append").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "message_id", Value: input.MessageID},
		core.Option{Key: "content", Value: input.Content},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatThinkingAppendOutput{}, err
		}
		return nil, ChatThinkingAppendOutput{}, nil
	}
	content, ok := result.Value.(string)
	if !ok {
		return nil, ChatThinkingAppendOutput{}, coreerr.E("mcp.chatThinkingAppend", "unexpected result type", nil)
	}
	return nil, ChatThinkingAppendOutput{Content: content}, nil
}

// --- chat_thinking_end ---

type ChatThinkingEndInput struct {
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id,omitempty"`
	DurationMS     int64     `json:"duration_ms,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
}

type ChatThinkingEndOutput struct {
	DurationMS int64 `json:"duration_ms"`
}

func (s *Subsystem) chatThinkingEnd(_ context.Context, _ *mcp.CallToolRequest, input ChatThinkingEndInput) (*mcp.CallToolResult, ChatThinkingEndOutput, error) {
	result := s.core.Action("gui.chat.thinking.end").Run(context.Background(), core.NewOptions(
		core.Option{Key: "conversation_id", Value: input.ConversationID},
		core.Option{Key: "message_id", Value: input.MessageID},
		core.Option{Key: "duration_ms", Value: input.DurationMS},
		core.Option{Key: "started_at", Value: input.StartedAt},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, ChatThinkingEndOutput{}, err
		}
		return nil, ChatThinkingEndOutput{}, nil
	}
	durationMS, ok := result.Value.(int64)
	if !ok {
		return nil, ChatThinkingEndOutput{}, coreerr.E("mcp.chatThinkingEnd", "unexpected result type", nil)
	}
	return nil, ChatThinkingEndOutput{DurationMS: durationMS}, nil
}

func decodeChatValue[T any](value any) (T, error) {
	var output T
	result := core.JSONUnmarshalString(core.JSONMarshalString(value), &output)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return output, err
		}
		return output, coreerr.E("mcp.decodeChatValue", "failed to decode chat value", nil)
	}
	return output, nil
}

func (s *Subsystem) registerChatTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{
		Name:        "chat_send",
		Description: `Send a chat message and stream a reply through the chat service. Example: {"conversation_id":"conv-1","content":"Hello"}`,
	}, s.chatSend)
	addTool(s, server, &mcp.Tool{
		Name:        "chat_clear",
		Description: `Clear a conversation's message history. Example: {"id":"conv-1"}`,
	}, s.chatClear)
	addTool(s, server, &mcp.Tool{
		Name:        "chat_history",
		Description: `Get a conversation and its message history. Example: {"id":"conv-1"}`,
	}, s.chatHistory)
	addTool(s, server, &mcp.Tool{Name: "chat_models", Description: "List available chat models"}, s.chatModels)
	addTool(s, server, &mcp.Tool{
		Name:        "chat_select_model",
		Description: `Select the active chat model for the current settings or conversation. Example: {"model":"lemer","conversation_id":"conv-1"}`,
	}, s.chatSelectModel)
	addTool(s, server, &mcp.Tool{Name: "chat_settings_save", Description: "Persist chat settings"}, s.chatSettingsSave)
	addTool(s, server, &mcp.Tool{Name: "chat_settings_load", Description: "Load persisted chat settings"}, s.chatSettingsLoad)
	addTool(s, server, &mcp.Tool{Name: "chat_settings_defaults", Description: "Return the default chat settings"}, s.chatSettingsDefaults)
	addTool(s, server, &mcp.Tool{Name: "chat_settings_reset", Description: "Reset chat settings to defaults"}, s.chatSettingsReset)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_list", Description: "List stored conversations"}, s.chatConversationsList)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_get", Description: "Get a stored conversation by id"}, s.chatConversationsGet)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_delete", Description: "Delete a stored conversation"}, s.chatConversationsDelete)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_search", Description: "Search stored conversations by text"}, s.chatConversationsSearch)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_new", Description: "Create a new conversation"}, s.chatConversationsNew)
	addTool(s, server, &mcp.Tool{Name: "chat_conversation_save", Description: "Save a conversation object"}, s.chatConversationSave)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_rename", Description: "Rename a conversation"}, s.chatConversationsRename)
	addTool(s, server, &mcp.Tool{Name: "chat_conversations_export", Description: "Export a conversation as Markdown"}, s.chatConversationsExport)
	addTool(s, server, &mcp.Tool{Name: "chat_attach_image", Description: "Queue an image attachment for the next chat message"}, s.chatAttachImage)
	addTool(s, server, &mcp.Tool{Name: "chat_remove_image", Description: "Remove a queued image attachment"}, s.chatRemoveImage)
	addTool(s, server, &mcp.Tool{Name: "chat_thinking_start", Description: "Begin a thinking state for a conversation"}, s.chatThinkingStart)
	addTool(s, server, &mcp.Tool{Name: "chat_thinking_append", Description: "Append thinking content for a conversation"}, s.chatThinkingAppend)
	addTool(s, server, &mcp.Tool{Name: "chat_thinking_end", Description: "End a thinking state for a conversation"}, s.chatThinkingEnd)
}
