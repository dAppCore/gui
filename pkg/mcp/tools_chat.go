package mcp

import (
	"context"

	"dappco.re/go/core/gui/pkg/display"
	coreerr "forge.lthn.ai/core/go-log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ChatModelsInput struct{}
type ChatModelsOutput struct {
	SelectedModel string               `json:"selected_model"`
	Models        []display.ModelEntry `json:"models"`
}

func (s *Subsystem) chatModels(_ context.Context, _ *mcp.CallToolRequest, _ ChatModelsInput) (*mcp.CallToolResult, ChatModelsOutput, error) {
	result, _, err := s.core.QUERY(display.QueryChatModels{})
	if err != nil {
		return nil, ChatModelsOutput{}, err
	}
	models, ok := result.([]display.ModelEntry)
	if !ok {
		return nil, ChatModelsOutput{}, coreerr.E("mcp.chatModels", "unexpected result type", nil)
	}
	selected := ""
	for _, model := range models {
		if model.Loaded {
			selected = model.Name
			break
		}
	}
	return nil, ChatModelsOutput{SelectedModel: selected, Models: models}, nil
}

type ChatSelectModelInput struct {
	Model string `json:"model"`
}

func (s *Subsystem) chatSelectModel(_ context.Context, _ *mcp.CallToolRequest, input ChatSelectModelInput) (*mcp.CallToolResult, ChatModelsOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskSelectModel{Model: input.Model})
	if err != nil {
		return nil, ChatModelsOutput{}, err
	}
	models, ok := result.([]display.ModelEntry)
	if !ok {
		return nil, ChatModelsOutput{}, coreerr.E("mcp.chatSelectModel", "unexpected result type", nil)
	}
	selected := ""
	for _, model := range models {
		if model.Loaded {
			selected = model.Name
			break
		}
	}
	return nil, ChatModelsOutput{SelectedModel: selected, Models: models}, nil
}

type ChatSettingsGetInput struct{}
type ChatSettingsOutput struct {
	Settings display.ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsGet(_ context.Context, _ *mcp.CallToolRequest, _ ChatSettingsGetInput) (*mcp.CallToolResult, ChatSettingsOutput, error) {
	result, _, err := s.core.QUERY(display.QueryChatSettingsLoad{})
	if err != nil {
		return nil, ChatSettingsOutput{}, err
	}
	settings, ok := result.(display.ChatSettings)
	if !ok {
		return nil, ChatSettingsOutput{}, coreerr.E("mcp.chatSettingsGet", "unexpected result type", nil)
	}
	return nil, ChatSettingsOutput{Settings: settings}, nil
}

type ChatSettingsSetInput struct {
	Settings display.ChatSettings `json:"settings"`
}

func (s *Subsystem) chatSettingsSet(_ context.Context, _ *mcp.CallToolRequest, input ChatSettingsSetInput) (*mcp.CallToolResult, ChatSettingsOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskChatSettingsSave{Settings: input.Settings})
	if err != nil {
		return nil, ChatSettingsOutput{}, err
	}
	settings, ok := result.(display.ChatSettings)
	if !ok {
		return nil, ChatSettingsOutput{}, coreerr.E("mcp.chatSettingsSet", "unexpected result type", nil)
	}
	return nil, ChatSettingsOutput{Settings: settings}, nil
}

type ChatSettingsResetInput struct{}

func (s *Subsystem) chatSettingsReset(_ context.Context, _ *mcp.CallToolRequest, _ ChatSettingsResetInput) (*mcp.CallToolResult, ChatSettingsOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskChatSettingsReset{})
	if err != nil {
		return nil, ChatSettingsOutput{}, err
	}
	settings, ok := result.(display.ChatSettings)
	if !ok {
		return nil, ChatSettingsOutput{}, coreerr.E("mcp.chatSettingsReset", "unexpected result type", nil)
	}
	return nil, ChatSettingsOutput{Settings: settings}, nil
}

type ChatConversationNewInput struct{}
type ChatConversationOutput struct {
	Conversation display.Conversation `json:"conversation"`
}

func (s *Subsystem) chatConversationNew(_ context.Context, _ *mcp.CallToolRequest, _ ChatConversationNewInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskConversationNew{})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatConversationNew", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatConversationSaveInput struct {
	Conversation display.Conversation `json:"conversation"`
}

func (s *Subsystem) chatConversationSave(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationSaveInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskConversationSave{Conversation: input.Conversation})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatConversationSave", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatConversationGetInput struct {
	ID string `json:"id"`
}

func (s *Subsystem) chatConversationGet(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationGetInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.QUERY(display.QueryConversationGet{ID: input.ID})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatConversationGet", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatConversationsListInput struct{}
type ChatConversationsOutput struct {
	Conversations []display.Conversation `json:"conversations"`
}

func (s *Subsystem) chatConversationsList(_ context.Context, _ *mcp.CallToolRequest, _ ChatConversationsListInput) (*mcp.CallToolResult, ChatConversationsOutput, error) {
	result, _, err := s.core.QUERY(display.QueryConversationsList{})
	if err != nil {
		return nil, ChatConversationsOutput{}, err
	}
	conversations, ok := result.([]display.Conversation)
	if !ok {
		return nil, ChatConversationsOutput{}, coreerr.E("mcp.chatConversationsList", "unexpected result type", nil)
	}
	return nil, ChatConversationsOutput{Conversations: conversations}, nil
}

type ChatConversationsSearchInput struct {
	Query string `json:"q"`
}

func (s *Subsystem) chatConversationsSearch(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationsSearchInput) (*mcp.CallToolResult, ChatConversationsOutput, error) {
	result, _, err := s.core.QUERY(display.QueryConversationsSearch{Query: input.Query})
	if err != nil {
		return nil, ChatConversationsOutput{}, err
	}
	conversations, ok := result.([]display.Conversation)
	if !ok {
		return nil, ChatConversationsOutput{}, coreerr.E("mcp.chatConversationsSearch", "unexpected result type", nil)
	}
	return nil, ChatConversationsOutput{Conversations: conversations}, nil
}

type ChatConversationRenameInput struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *Subsystem) chatConversationRename(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationRenameInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskConversationRename{ID: input.ID, Title: input.Title})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatConversationRename", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatConversationDeleteInput struct {
	ID string `json:"id"`
}
type ChatConversationDeleteOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) chatConversationDelete(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationDeleteInput) (*mcp.CallToolResult, ChatConversationDeleteOutput, error) {
	_, _, err := s.core.PERFORM(display.TaskConversationDelete{ID: input.ID})
	if err != nil {
		return nil, ChatConversationDeleteOutput{}, err
	}
	return nil, ChatConversationDeleteOutput{Success: true}, nil
}

type ChatConversationExportInput struct {
	ID string `json:"id"`
}
type ChatConversationExportOutput struct {
	Markdown string `json:"markdown"`
}

func (s *Subsystem) chatConversationExport(_ context.Context, _ *mcp.CallToolRequest, input ChatConversationExportInput) (*mcp.CallToolResult, ChatConversationExportOutput, error) {
	result, _, err := s.core.QUERY(display.QueryConversationExport{ID: input.ID})
	if err != nil {
		return nil, ChatConversationExportOutput{}, err
	}
	markdown, ok := result.(string)
	if !ok {
		return nil, ChatConversationExportOutput{}, coreerr.E("mcp.chatConversationExport", "unexpected result type", nil)
	}
	return nil, ChatConversationExportOutput{Markdown: markdown}, nil
}

type ChatHistoryInput struct {
	ConversationID string `json:"conversation_id"`
}
type ChatHistoryOutput struct {
	Messages []display.ChatMessage `json:"messages"`
}

func (s *Subsystem) chatHistory(_ context.Context, _ *mcp.CallToolRequest, input ChatHistoryInput) (*mcp.CallToolResult, ChatHistoryOutput, error) {
	result, _, err := s.core.QUERY(display.QueryChatHistory{ConversationID: input.ConversationID})
	if err != nil {
		return nil, ChatHistoryOutput{}, err
	}
	messages, ok := result.([]display.ChatMessage)
	if !ok {
		return nil, ChatHistoryOutput{}, coreerr.E("mcp.chatHistory", "unexpected result type", nil)
	}
	return nil, ChatHistoryOutput{Messages: messages}, nil
}

type ChatSendInput struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

func (s *Subsystem) chatSend(_ context.Context, _ *mcp.CallToolRequest, input ChatSendInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskChatSend{
		ConversationID: input.ConversationID,
		Content:        input.Content,
	})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatSend", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatClearInput struct {
	ConversationID string `json:"conversation_id"`
}

func (s *Subsystem) chatClear(_ context.Context, _ *mcp.CallToolRequest, input ChatClearInput) (*mcp.CallToolResult, ChatConversationOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskChatClear{ConversationID: input.ConversationID})
	if err != nil {
		return nil, ChatConversationOutput{}, err
	}
	conv, ok := result.(display.Conversation)
	if !ok {
		return nil, ChatConversationOutput{}, coreerr.E("mcp.chatClear", "unexpected result type", nil)
	}
	return nil, ChatConversationOutput{Conversation: conv}, nil
}

type ChatAttachmentsGetInput struct {
	ConversationID string `json:"conversation_id"`
}
type ChatAttachmentsOutput struct {
	Attachments []display.ImageAttachment `json:"attachments"`
}

func (s *Subsystem) chatAttachmentsGet(_ context.Context, _ *mcp.CallToolRequest, input ChatAttachmentsGetInput) (*mcp.CallToolResult, ChatAttachmentsOutput, error) {
	result, _, err := s.core.QUERY(display.QueryQueuedImages{ConversationID: input.ConversationID})
	if err != nil {
		return nil, ChatAttachmentsOutput{}, err
	}
	attachments, ok := result.([]display.ImageAttachment)
	if !ok {
		return nil, ChatAttachmentsOutput{}, coreerr.E("mcp.chatAttachmentsGet", "unexpected result type", nil)
	}
	return nil, ChatAttachmentsOutput{Attachments: attachments}, nil
}

type ChatAttachImageInput struct {
	ConversationID string                  `json:"conversation_id"`
	Attachment     display.ImageAttachment `json:"attachment"`
}

func (s *Subsystem) chatAttachImage(_ context.Context, _ *mcp.CallToolRequest, input ChatAttachImageInput) (*mcp.CallToolResult, ChatAttachmentsOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskAttachImage{
		ConversationID: input.ConversationID,
		Attachment:     input.Attachment,
	})
	if err != nil {
		return nil, ChatAttachmentsOutput{}, err
	}
	attachments, ok := result.([]display.ImageAttachment)
	if !ok {
		return nil, ChatAttachmentsOutput{}, coreerr.E("mcp.chatAttachImage", "unexpected result type", nil)
	}
	return nil, ChatAttachmentsOutput{Attachments: attachments}, nil
}

type ChatDetachImageInput struct {
	ConversationID string `json:"conversation_id"`
	AttachmentID   string `json:"attachment_id"`
}

func (s *Subsystem) chatDetachImage(_ context.Context, _ *mcp.CallToolRequest, input ChatDetachImageInput) (*mcp.CallToolResult, ChatAttachmentsOutput, error) {
	result, _, err := s.core.PERFORM(display.TaskDetachImage{
		ConversationID: input.ConversationID,
		AttachmentID:   input.AttachmentID,
	})
	if err != nil {
		return nil, ChatAttachmentsOutput{}, err
	}
	attachments, ok := result.([]display.ImageAttachment)
	if !ok {
		return nil, ChatAttachmentsOutput{}, coreerr.E("mcp.chatDetachImage", "unexpected result type", nil)
	}
	return nil, ChatAttachmentsOutput{Attachments: attachments}, nil
}

func (s *Subsystem) registerChatTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "chat_models", Description: "List locally available chat models and the currently selected model"}, s.chatModels)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_select_model", Description: "Select the active local chat model"}, s.chatSelectModel)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_settings_get", Description: "Load persisted chat settings"}, s.chatSettingsGet)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_settings_set", Description: "Persist chat settings defaults"}, s.chatSettingsSet)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_settings_reset", Description: "Reset chat settings to RFC defaults"}, s.chatSettingsReset)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_new", Description: "Create a new chat conversation"}, s.chatConversationNew)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_save", Description: "Save or import a full chat conversation"}, s.chatConversationSave)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_get", Description: "Load a specific chat conversation"}, s.chatConversationGet)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversations_list", Description: "List saved chat conversations"}, s.chatConversationsList)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversations_search", Description: "Search chat conversations by title or content"}, s.chatConversationsSearch)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_rename", Description: "Rename a chat conversation"}, s.chatConversationRename)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_delete", Description: "Delete a chat conversation"}, s.chatConversationDelete)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_conversation_export", Description: "Export a chat conversation as Markdown"}, s.chatConversationExport)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_history", Description: "Return the ordered message history for a conversation"}, s.chatHistory)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_send", Description: "Append a user message and create the paired assistant response placeholder"}, s.chatSend)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_clear", Description: "Clear all messages from a conversation"}, s.chatClear)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_attachments_get", Description: "List queued image attachments for the next message"}, s.chatAttachmentsGet)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_attach_image", Description: "Queue an image attachment for the next chat message"}, s.chatAttachImage)
	mcp.AddTool(server, &mcp.Tool{Name: "chat_detach_image", Description: "Remove a queued image attachment"}, s.chatDetachImage)
}
