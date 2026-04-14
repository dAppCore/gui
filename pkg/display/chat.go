package display

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"

	"forge.lthn.ai/core/config"
	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

const chatConfigSection = "chat"

type ModelEntry struct {
	Name           string `json:"name"`
	Architecture   string `json:"architecture"`
	QuantBits      int    `json:"quant_bits"`
	SizeBytes      int64  `json:"size_bytes"`
	Loaded         bool   `json:"loaded"`
	Backend        string `json:"backend"`
	SupportsVision bool   `json:"supports_vision,omitempty"`
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

type ImageAttachment struct {
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type ToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Content    string `json:"content"`
}

type ToolInvocation struct {
	Call      ToolCall   `json:"call"`
	Result    ToolResult `json:"result"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   time.Time  `json:"ended_at"`
	Error     string     `json:"error,omitempty"`
}

type ThinkingState struct {
	Active     bool      `json:"active"`
	Content    string    `json:"content"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

type ChatMessage struct {
	ID           string            `json:"id"`
	Role         string            `json:"role"`
	Content      string            `json:"content"`
	CreatedAt    time.Time         `json:"created_at"`
	Streaming    bool              `json:"streaming,omitempty"`
	FinishReason string            `json:"finish_reason,omitempty"`
	Attachments  []ImageAttachment `json:"attachments,omitempty"`
	Thinking     *ThinkingState    `json:"thinking,omitempty"`
	ToolCalls    []ToolInvocation  `json:"tool_calls,omitempty"`
}

type Conversation struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Model     string        `json:"model"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Messages  []ChatMessage `json:"messages"`
	Settings  *ChatSettings `json:"settings,omitempty"`
}

type ChatSnapshot struct {
	Settings         ChatSettings                 `json:"settings"`
	SelectedModel    string                       `json:"selected_model"`
	Conversations    map[string]Conversation      `json:"conversations"`
	QueuedImages     map[string][]ImageAttachment `json:"queued_images"`
	Thinking         map[string]ThinkingState     `json:"thinking"`
	StreamingMessage map[string]string            `json:"streaming_message"`
	Models           []ModelEntry                 `json:"models"`
}

type ChatStore struct {
	mu               sync.RWMutex
	settings         ChatSettings
	selectedModel    string
	conversations    map[string]Conversation
	queuedImages     map[string][]ImageAttachment
	thinking         map[string]ThinkingState
	streamingMessage map[string]string
	models           []ModelEntry
	nextID           uint64
}

func defaultChatSettings() ChatSettings {
	return ChatSettings{
		Temperature:   1.0,
		TopP:          0.95,
		TopK:          64,
		MaxTokens:     2048,
		ContextWindow: 8192,
		SystemPrompt:  "You are a helpful assistant.",
		DefaultModel:  "lemer",
	}
}

func defaultModelEntries() []ModelEntry {
	return []ModelEntry{
		{
			Name:           "lemer",
			Architecture:   "gemma3",
			QuantBits:      4,
			SizeBytes:      1_500_000_000,
			Loaded:         true,
			Backend:        "metal",
			SupportsVision: true,
		},
		{
			Name:           "lemma",
			Architecture:   "gemma3",
			QuantBits:      8,
			SizeBytes:      3_200_000_000,
			Loaded:         false,
			Backend:        "metal",
			SupportsVision: true,
		},
		{
			Name:         "lemmy",
			Architecture: "qwen3",
			QuantBits:    4,
			SizeBytes:    1_100_000_000,
			Loaded:       false,
			Backend:      "ollama",
		},
	}
}

func NewChatStore() *ChatStore {
	settings := defaultChatSettings()
	return &ChatStore{
		settings:         settings,
		selectedModel:    settings.DefaultModel,
		conversations:    make(map[string]Conversation),
		queuedImages:     make(map[string][]ImageAttachment),
		thinking:         make(map[string]ThinkingState),
		streamingMessage: make(map[string]string),
		models:           defaultModelEntries(),
	}
}

func (s *ChatStore) Load(cfg *config.Config) {
	if cfg == nil {
		return
	}

	var snapshot ChatSnapshot
	if err := cfg.Get(chatConfigSection, &snapshot); err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if snapshot.Settings == (ChatSettings{}) {
		snapshot.Settings = defaultChatSettings()
	}
	if snapshot.SelectedModel == "" {
		snapshot.SelectedModel = snapshot.Settings.DefaultModel
	}
	if len(snapshot.Models) == 0 {
		snapshot.Models = defaultModelEntries()
	}
	if snapshot.Conversations == nil {
		snapshot.Conversations = make(map[string]Conversation)
	}
	if snapshot.QueuedImages == nil {
		snapshot.QueuedImages = make(map[string][]ImageAttachment)
	}
	if snapshot.Thinking == nil {
		snapshot.Thinking = make(map[string]ThinkingState)
	}
	if snapshot.StreamingMessage == nil {
		snapshot.StreamingMessage = make(map[string]string)
	}

	s.settings = snapshot.Settings
	s.selectedModel = snapshot.SelectedModel
	s.models = snapshot.Models
	s.conversations = snapshot.Conversations
	s.queuedImages = snapshot.QueuedImages
	s.thinking = snapshot.Thinking
	s.streamingMessage = snapshot.StreamingMessage

	for _, conv := range s.conversations {
		if conv.ID == "" {
			continue
		}
		if n := parseCounter(conv.ID); n > s.nextID {
			s.nextID = n
		}
		for _, msg := range conv.Messages {
			if n := parseCounter(msg.ID); n > s.nextID {
				s.nextID = n
			}
			for _, attachment := range msg.Attachments {
				if n := parseCounter(attachment.ID); n > s.nextID {
					s.nextID = n
				}
			}
		}
	}
	for _, attachments := range s.queuedImages {
		for _, attachment := range attachments {
			if n := parseCounter(attachment.ID); n > s.nextID {
				s.nextID = n
			}
		}
	}
	s.ensureAttachmentIdentifiersLocked()
}

func (s *ChatStore) Snapshot() ChatSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshotLocked()
}

func (s *ChatStore) snapshotLocked() ChatSnapshot {
	return ChatSnapshot{
		Settings:         s.settings,
		SelectedModel:    s.selectedModel,
		Conversations:    cloneConversationMap(s.conversations),
		QueuedImages:     cloneAttachmentMap(s.queuedImages),
		Thinking:         cloneThinkingMap(s.thinking),
		StreamingMessage: cloneStringMap(s.streamingMessage),
		Models:           append([]ModelEntry(nil), s.models...),
	}
}

func (s *ChatStore) persist(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}
	if err := cfg.Set(chatConfigSection, s.Snapshot()); err != nil {
		return err
	}
	return cfg.Commit()
}

func (s *ChatStore) Models() []ModelEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]ModelEntry(nil), s.models...)
}

func (s *ChatStore) SelectedModel() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectedModel
}

func (s *ChatStore) SelectModel(name string) ([]ModelEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := -1
	for i := range s.models {
		if s.models[i].Name == name {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, coreerr.E("display.chat.SelectModel", "unknown model: "+name, nil)
	}

	for i := range s.models {
		s.models[i].Loaded = i == index
	}
	s.selectedModel = name
	if s.settings.DefaultModel == "" {
		s.settings.DefaultModel = name
	}
	return append([]ModelEntry(nil), s.models...), nil
}

func (s *ChatStore) Settings() ChatSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

func (s *ChatStore) SaveSettings(settings ChatSettings) ChatSettings {
	s.mu.Lock()
	defer s.mu.Unlock()

	if settings.DefaultModel == "" {
		settings.DefaultModel = s.selectedModel
	}
	s.settings = settings
	if s.selectedModel == "" {
		s.selectedModel = settings.DefaultModel
	}
	return s.settings
}

func (s *ChatStore) ResetSettings() ChatSettings {
	return s.SaveSettings(defaultChatSettings())
}

func (s *ChatStore) SaveConversation(input Conversation) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	conv := cloneConversation(input)
	existing, exists := s.conversations[conv.ID]

	if conv.ID == "" {
		conv.ID = s.nextIdentifier("conv")
	}
	if conv.CreatedAt.IsZero() {
		switch {
		case exists && !existing.CreatedAt.IsZero():
			conv.CreatedAt = existing.CreatedAt
		default:
			conv.CreatedAt = now
		}
	}
	if conv.Model == "" {
		conv.Model = s.nextResponseModelLocked(conv)
	}
	if conv.Settings == nil && exists && existing.Settings != nil {
		copySettings := *existing.Settings
		conv.Settings = &copySettings
	}

	conv.Messages = s.normalizeMessagesLocked(conv.Messages, conv.CreatedAt)
	if strings.TrimSpace(conv.Title) == "" {
		conv.Title = deriveConversationTitle(firstUserMessage(conv.Messages))
	}
	if strings.TrimSpace(conv.Title) == "" {
		conv.Title = "New conversation"
	}
	if conv.UpdatedAt.IsZero() {
		conv.UpdatedAt = conv.CreatedAt
		if len(conv.Messages) > 0 {
			conv.UpdatedAt = conv.Messages[len(conv.Messages)-1].CreatedAt
		}
	}
	if conv.UpdatedAt.Before(conv.CreatedAt) {
		conv.UpdatedAt = conv.CreatedAt
	}

	s.conversations[conv.ID] = conv
	return cloneConversation(conv), nil
}

func (s *ChatStore) NewConversation() Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	id := s.nextIdentifier("conv")
	model := s.selectedModel
	if model == "" {
		model = s.settings.DefaultModel
	}
	conv := Conversation{
		ID:        id,
		Title:     "New conversation",
		Model:     model,
		CreatedAt: now,
		UpdatedAt: now,
		Messages:  []ChatMessage{},
	}
	s.conversations[id] = conv
	return conv
}

func (s *ChatStore) ListConversations() []Conversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Conversation, 0, len(s.conversations))
	for _, conv := range s.conversations {
		items = append(items, cloneConversation(conv))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items
}

func (s *ChatStore) Conversation(id string) (Conversation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.conversations[id]
	if !ok {
		return Conversation{}, false
	}
	return cloneConversation(conv), true
}

func (s *ChatStore) QueuedImages(conversationID string) ([]ImageAttachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.conversations[conversationID]; !ok {
		return nil, coreerr.E("display.chat.QueuedImages", "conversation not found: "+conversationID, nil)
	}
	return append([]ImageAttachment(nil), s.queuedImages[conversationID]...), nil
}

func (s *ChatStore) DeleteConversation(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conversations[id]; !ok {
		return false
	}
	delete(s.conversations, id)
	delete(s.queuedImages, id)
	delete(s.thinking, id)
	delete(s.streamingMessage, id)
	return true
}

func (s *ChatStore) SearchConversations(query string) []Conversation {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return s.ListConversations()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []Conversation
	for _, conv := range s.conversations {
		if conversationMatchesSearchQuery(conv, query) {
			items = append(items, cloneConversation(conv))
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})
	return items
}

func (s *ChatStore) History(conversationID string) ([]ChatMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return nil, coreerr.E("display.chat.History", "conversation not found: "+conversationID, nil)
	}
	return cloneMessages(conv.Messages), nil
}

func (s *ChatStore) ClearConversation(conversationID string) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, coreerr.E("display.chat.ClearConversation", "conversation not found: "+conversationID, nil)
	}
	conv.Messages = nil
	conv.Title = "New conversation"
	conv.UpdatedAt = time.Now().UTC()
	s.conversations[conversationID] = conv
	delete(s.queuedImages, conversationID)
	delete(s.thinking, conversationID)
	delete(s.streamingMessage, conversationID)
	return cloneConversation(conv), nil
}

func (s *ChatStore) QueueImage(conversationID string, attachment ImageAttachment) ([]ImageAttachment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return nil, coreerr.E("display.chat.QueueImage", "conversation not found: "+conversationID, nil)
	}
	if err := validateImageAttachment(attachment); err != nil {
		return nil, coreerr.E("display.chat.QueueImage", err.Error(), nil)
	}
	model := s.nextResponseModelLocked(conv)
	if model != "" && !s.modelSupportsVisionLocked(model) {
		return nil, coreerr.E("display.chat.QueueImage", "selected model does not support vision: "+model, nil)
	}
	if attachment.ID == "" {
		attachment.ID = s.nextIdentifier("img")
	}
	s.queuedImages[conversationID] = append(s.queuedImages[conversationID], attachment)
	return append([]ImageAttachment(nil), s.queuedImages[conversationID]...), nil
}

func (s *ChatStore) RemoveQueuedImage(conversationID, attachmentID string) ([]ImageAttachment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conversations[conversationID]; !ok {
		return nil, coreerr.E("display.chat.RemoveQueuedImage", "conversation not found: "+conversationID, nil)
	}
	if attachmentID == "" {
		return nil, coreerr.E("display.chat.RemoveQueuedImage", "attachment id is required", nil)
	}

	attachments := s.queuedImages[conversationID]
	filtered := attachments[:0]
	removed := false
	for _, attachment := range attachments {
		if attachment.ID == attachmentID {
			removed = true
			continue
		}
		filtered = append(filtered, attachment)
	}
	if !removed {
		return nil, coreerr.E("display.chat.RemoveQueuedImage", "attachment not found: "+attachmentID, nil)
	}
	if len(filtered) == 0 {
		delete(s.queuedImages, conversationID)
		return []ImageAttachment{}, nil
	}
	s.queuedImages[conversationID] = append([]ImageAttachment(nil), filtered...)
	return append([]ImageAttachment(nil), s.queuedImages[conversationID]...), nil
}

func (s *ChatStore) SendMessage(conversationID, content string) (Conversation, ChatMessage, ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, ChatMessage{}, ChatMessage{}, coreerr.E("display.chat.SendMessage", "conversation not found: "+conversationID, nil)
	}
	if strings.TrimSpace(content) == "" {
		return Conversation{}, ChatMessage{}, ChatMessage{}, coreerr.E("display.chat.SendMessage", "message content is required", nil)
	}

	now := time.Now().UTC()
	model := s.nextResponseModelLocked(conv)
	userMessage := ChatMessage{
		ID:          s.nextIdentifier("msg"),
		Role:        "user",
		Content:     content,
		CreatedAt:   now,
		Attachments: append([]ImageAttachment(nil), s.queuedImages[conversationID]...),
	}
	assistantMessage := ChatMessage{
		ID:        s.nextIdentifier("msg"),
		Role:      "assistant",
		Content:   buildAssistantPlaceholder(model, content),
		CreatedAt: now.Add(250 * time.Millisecond),
		Streaming: false,
	}
	if thinking, ok := s.thinking[conversationID]; ok && strings.TrimSpace(thinking.Content) != "" {
		copyThinking := thinking
		copyThinking.Active = false
		copyThinking.FinishedAt = time.Now().UTC()
		assistantMessage.Thinking = &copyThinking
		s.thinking[conversationID] = copyThinking
	}

	conv.Messages = append(conv.Messages, userMessage, assistantMessage)
	if len(conv.Messages) == 2 {
		conv.Title = deriveConversationTitle(content)
	}
	conv.UpdatedAt = assistantMessage.CreatedAt
	conv.Model = model
	s.conversations[conversationID] = conv
	delete(s.queuedImages, conversationID)
	delete(s.streamingMessage, conversationID)

	return cloneConversation(conv), userMessage, assistantMessage, nil
}

func (s *ChatStore) StartStreaming(conversationID string) (Conversation, ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, ChatMessage{}, coreerr.E("display.chat.StartStreaming", "conversation not found: "+conversationID, nil)
	}

	now := time.Now().UTC()
	conv.Model = s.nextResponseModelLocked(conv)
	var assistantMessage *ChatMessage
	if canReuseAssistantForStreaming(conv) {
		assistantMessage = &conv.Messages[len(conv.Messages)-1]
		assistantMessage.Content = ""
		assistantMessage.Streaming = true
		assistantMessage.FinishReason = ""
		if thinking, ok := s.thinking[conversationID]; ok {
			copyThinking := thinking
			assistantMessage.Thinking = &copyThinking
		}
	} else {
		message := ChatMessage{
			ID:        s.nextIdentifier("msg"),
			Role:      "assistant",
			Content:   "",
			CreatedAt: now,
			Streaming: true,
		}
		if thinking, ok := s.thinking[conversationID]; ok {
			copyThinking := thinking
			message.Thinking = &copyThinking
		}
		conv.Messages = append(conv.Messages, message)
		assistantMessage = &conv.Messages[len(conv.Messages)-1]
	}

	s.streamingMessage[conversationID] = assistantMessage.Content
	conv.UpdatedAt = now
	s.conversations[conversationID] = conv
	return cloneConversation(conv), cloneMessage(*assistantMessage), nil
}

func (s *ChatStore) AppendStreaming(conversationID, content string) (Conversation, ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, ChatMessage{}, coreerr.E("display.chat.AppendStreaming", "conversation not found: "+conversationID, nil)
	}

	index := latestStreamingAssistantIndex(conv.Messages)
	if index < 0 {
		return Conversation{}, ChatMessage{}, coreerr.E("display.chat.AppendStreaming", "streaming assistant message not found: "+conversationID, nil)
	}

	conv.Messages[index].Content += content
	if thinking, ok := s.thinking[conversationID]; ok {
		copyThinking := thinking
		conv.Messages[index].Thinking = &copyThinking
	}
	s.streamingMessage[conversationID] = conv.Messages[index].Content
	conv.UpdatedAt = time.Now().UTC()
	s.conversations[conversationID] = conv
	return cloneConversation(conv), cloneMessage(conv.Messages[index]), nil
}

func (s *ChatStore) FinishStreaming(conversationID, finishReason string) (Conversation, ChatMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, ChatMessage{}, coreerr.E("display.chat.FinishStreaming", "conversation not found: "+conversationID, nil)
	}

	index := latestStreamingAssistantIndex(conv.Messages)
	if index < 0 {
		return Conversation{}, ChatMessage{}, coreerr.E("display.chat.FinishStreaming", "streaming assistant message not found: "+conversationID, nil)
	}

	conv.Messages[index].Streaming = false
	conv.Messages[index].FinishReason = strings.TrimSpace(finishReason)
	if thinking, ok := s.thinking[conversationID]; ok {
		copyThinking := thinking
		if copyThinking.Active {
			copyThinking.Active = false
			if copyThinking.FinishedAt.IsZero() {
				copyThinking.FinishedAt = time.Now().UTC()
			}
			s.thinking[conversationID] = copyThinking
		}
		conv.Messages[index].Thinking = &copyThinking
	}
	delete(s.streamingMessage, conversationID)
	conv.UpdatedAt = time.Now().UTC()
	s.conversations[conversationID] = conv
	return cloneConversation(conv), cloneMessage(conv.Messages[index]), nil
}

func (s *ChatStore) StartThinking(conversationID string) (ThinkingState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conversations[conversationID]; !ok {
		return ThinkingState{}, coreerr.E("display.chat.StartThinking", "conversation not found: "+conversationID, nil)
	}
	state := ThinkingState{
		Active:    true,
		StartedAt: time.Now().UTC(),
	}
	s.thinking[conversationID] = state
	s.syncThinkingToStreamingMessageLocked(conversationID, state)
	return state, nil
}

func (s *ChatStore) AppendThinking(conversationID, content string) (ThinkingState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.conversations[conversationID]; !ok {
		return ThinkingState{}, coreerr.E("display.chat.AppendThinking", "conversation not found: "+conversationID, nil)
	}
	state, ok := s.thinking[conversationID]
	if !ok {
		state = ThinkingState{Active: true, StartedAt: time.Now().UTC()}
	}
	state.Active = true
	state.Content += content
	s.thinking[conversationID] = state
	s.syncThinkingToStreamingMessageLocked(conversationID, state)
	return state, nil
}

func (s *ChatStore) EndThinking(conversationID string) (ThinkingState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.thinking[conversationID]
	if !ok {
		return ThinkingState{}, coreerr.E("display.chat.EndThinking", "thinking state not found: "+conversationID, nil)
	}
	state.Active = false
	state.FinishedAt = time.Now().UTC()
	s.thinking[conversationID] = state
	s.syncThinkingToStreamingMessageLocked(conversationID, state)
	return state, nil
}

func (s *ChatStore) RecordToolResult(conversationID string, call ToolCall, result ToolResult, errText string) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok {
		return Conversation{}, coreerr.E("display.chat.RecordToolResult", "conversation not found: "+conversationID, nil)
	}
	index := latestAssistantIndex(conv.Messages)
	if index < 0 {
		return Conversation{}, coreerr.E("display.chat.RecordToolResult", "assistant message not found", nil)
	}
	last := &conv.Messages[index]
	startedAt := time.Now().UTC()
	last.ToolCalls = append(last.ToolCalls, ToolInvocation{
		Call:      call,
		Result:    result,
		StartedAt: startedAt,
		EndedAt:   time.Now().UTC(),
		Error:     errText,
	})
	conv.UpdatedAt = time.Now().UTC()
	s.conversations[conversationID] = conv
	return cloneConversation(conv), nil
}

func (s *ChatStore) RenameConversation(id, title string) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[id]
	if !ok {
		return Conversation{}, coreerr.E("display.chat.RenameConversation", "conversation not found: "+id, nil)
	}

	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		cleanTitle = "Untitled conversation"
	}
	conv.Title = cleanTitle
	conv.UpdatedAt = time.Now().UTC()
	s.conversations[id] = conv
	return cloneConversation(conv), nil
}

func (s *ChatStore) ExportConversationMarkdown(id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.conversations[id]
	if !ok {
		return "", coreerr.E("display.chat.ExportConversationMarkdown", "conversation not found: "+id, nil)
	}

	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(conv.Title)
	builder.WriteString("\n\n")
	builder.WriteString("- Conversation ID: ")
	builder.WriteString(conv.ID)
	builder.WriteString("\n")
	builder.WriteString("- Model: ")
	builder.WriteString(conv.Model)
	builder.WriteString("\n")
	builder.WriteString("- Updated: ")
	builder.WriteString(conv.UpdatedAt.UTC().Format(time.RFC3339))
	builder.WriteString("\n")

	for _, message := range conv.Messages {
		builder.WriteString("\n## ")
		builder.WriteString(roleHeading(message.Role))
		builder.WriteString("\n\n")
		builder.WriteString(message.Content)
		builder.WriteString("\n")

		if len(message.Attachments) > 0 {
			builder.WriteString("\n### Attachments\n")
			for _, attachment := range message.Attachments {
				builder.WriteString("- ")
				builder.WriteString(attachment.Filename)
				builder.WriteString(" (")
				builder.WriteString(attachment.MimeType)
				builder.WriteString(")\n")
			}
		}

		if message.Thinking != nil && strings.TrimSpace(message.Thinking.Content) != "" {
			builder.WriteString("\n### Thinking\n\n")
			builder.WriteString(message.Thinking.Content)
			builder.WriteString("\n")
		}

		if len(message.ToolCalls) > 0 {
			builder.WriteString("\n### Tool Calls\n")
			for _, invocation := range message.ToolCalls {
				builder.WriteString("\n#### ")
				builder.WriteString(invocation.Call.Name)
				builder.WriteString("\n\n")
				if len(invocation.Call.Arguments) > 0 {
					arguments, err := json.MarshalIndent(invocation.Call.Arguments, "", "  ")
					if err == nil {
						builder.WriteString("```json\n")
						builder.Write(arguments)
						builder.WriteString("\n```\n")
					}
				}
				if invocation.Result.Content != "" {
					builder.WriteString("\n")
					builder.WriteString(invocation.Result.Content)
					builder.WriteString("\n")
				}
				if invocation.Error != "" {
					builder.WriteString("\nError: ")
					builder.WriteString(invocation.Error)
					builder.WriteString("\n")
				}
			}
		}
	}

	return builder.String(), nil
}

type QueryChatHistory struct {
	ConversationID string `json:"conversation_id"`
}

type TaskChatSend struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

type TaskChatClear struct {
	ConversationID string `json:"conversation_id"`
}

type QueryChatModels struct{}

type TaskSelectModel struct {
	Model string `json:"model"`
}

type TaskChatSettingsSave struct {
	Settings ChatSettings `json:"settings"`
}

type QueryChatSettingsLoad struct{}

type TaskChatSettingsReset struct{}

type QueryConversationsList struct{}

type QueryConversationGet struct {
	ID string `json:"id"`
}

type QueryConversationsSearch struct {
	Query string `json:"q"`
}

type QueryQueuedImages struct {
	ConversationID string `json:"conversation_id"`
}

type QueryConversationExport struct {
	ID string `json:"id"`
}

type TaskConversationDelete struct {
	ID string `json:"id"`
}

type TaskConversationNew struct{}

type TaskConversationSave struct {
	Conversation Conversation `json:"conversation"`
}

type TaskConversationRename struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskAttachImage struct {
	ConversationID string          `json:"conversation_id"`
	Attachment     ImageAttachment `json:"attachment"`
}

type TaskDetachImage struct {
	ConversationID string `json:"conversation_id"`
	AttachmentID   string `json:"attachment_id"`
}

type TaskThinkingStart struct {
	ConversationID string `json:"conversation_id"`
}

type TaskThinkingAppend struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

type TaskThinkingEnd struct {
	ConversationID string `json:"conversation_id"`
}

type TaskRecordToolCall struct {
	ConversationID string     `json:"conversation_id"`
	Call           ToolCall   `json:"call"`
	Result         ToolResult `json:"result"`
	Error          string     `json:"error,omitempty"`
}

type TaskChatStreamStart struct {
	ConversationID string `json:"conversation_id"`
}

type TaskChatStreamAppend struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

type TaskChatStreamFinish struct {
	ConversationID string `json:"conversation_id"`
	FinishReason   string `json:"finish_reason,omitempty"`
}

func (s *Service) handleChatQuery(_ *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryChatHistory:
		history, err := s.chat.History(q.ConversationID)
		return history, true, err
	case QueryChatModels:
		return s.chat.Models(), true, nil
	case QueryChatSettingsLoad:
		return s.chat.Settings(), true, nil
	case QueryConversationsList:
		return s.chat.ListConversations(), true, nil
	case QueryConversationGet:
		conv, ok := s.chat.Conversation(q.ID)
		if !ok {
			return nil, true, coreerr.E("display.chat.QueryConversationGet", "conversation not found: "+q.ID, nil)
		}
		return conv, true, nil
	case QueryConversationsSearch:
		return s.chat.SearchConversations(q.Query), true, nil
	case QueryQueuedImages:
		attachments, err := s.chat.QueuedImages(q.ConversationID)
		return attachments, true, err
	case QueryConversationExport:
		content, err := s.chat.ExportConversationMarkdown(q.ID)
		return content, true, err
	default:
		return nil, false, nil
	}
}

func (s *Service) handleChatTask(_ *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case TaskConversationNew:
		conv := s.chat.NewConversation()
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.conversation.created",
				Data: map[string]any{
					"conversation": conv,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskConversationSave:
		conv, err := s.chat.SaveConversation(t.Conversation)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.conversation.saved",
				Data: map[string]any{
					"conversation": conv,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskConversationRename:
		conv, err := s.chat.RenameConversation(t.ID, t.Title)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.conversation.renamed",
				Data: map[string]any{
					"conversation_id": conv.ID,
					"title":           conv.Title,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskChatSend:
		conv, userMessage, assistantMessage, err := s.chat.SendMessage(t.ConversationID, t.Content)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.message",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"user":            userMessage,
					"assistant":       assistantMessage,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskChatClear:
		conv, err := s.chat.ClearConversation(t.ConversationID)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.cleared",
				Data: map[string]any{
					"conversation_id": conv.ID,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskSelectModel:
		models, err := s.chat.SelectModel(t.Model)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.model.selected",
				Data: map[string]any{
					"model":  t.Model,
					"models": models,
				},
			})
		}
		return models, true, s.chat.persist(s.configFile)
	case TaskChatSettingsSave:
		settings := s.chat.SaveSettings(t.Settings)
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.settings.updated",
				Data: map[string]any{
					"settings": settings,
				},
			})
		}
		return settings, true, s.chat.persist(s.configFile)
	case TaskChatSettingsReset:
		settings := s.chat.ResetSettings()
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.settings.updated",
				Data: map[string]any{
					"settings": settings,
				},
			})
		}
		return settings, true, s.chat.persist(s.configFile)
	case TaskConversationDelete:
		deleted := s.chat.DeleteConversation(t.ID)
		if !deleted {
			return nil, true, coreerr.E("display.chat.TaskConversationDelete", "conversation not found: "+t.ID, nil)
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.conversation.deleted",
				Data: map[string]any{
					"conversation_id": t.ID,
				},
			})
		}
		return deleted, true, s.chat.persist(s.configFile)
	case TaskAttachImage:
		attachments, err := s.chat.QueueImage(t.ConversationID, t.Attachment)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.attachments.updated",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"attachments":     attachments,
				},
			})
		}
		return attachments, true, s.chat.persist(s.configFile)
	case TaskDetachImage:
		attachments, err := s.chat.RemoveQueuedImage(t.ConversationID, t.AttachmentID)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.attachments.updated",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"attachments":     attachments,
				},
			})
		}
		return attachments, true, s.chat.persist(s.configFile)
	case TaskThinkingStart:
		state, err := s.chat.StartThinking(t.ConversationID)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.thinking.start",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"thinking":        state,
				},
			})
		}
		return state, true, s.chat.persist(s.configFile)
	case TaskThinkingAppend:
		state, err := s.chat.AppendThinking(t.ConversationID, t.Content)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.thinking.delta",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"thinking":        state,
					"delta":           t.Content,
				},
			})
		}
		return state, true, s.chat.persist(s.configFile)
	case TaskThinkingEnd:
		state, err := s.chat.EndThinking(t.ConversationID)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.thinking.end",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"thinking":        state,
				},
			})
		}
		return state, true, s.chat.persist(s.configFile)
	case TaskRecordToolCall:
		conv, err := s.chat.RecordToolResult(t.ConversationID, t.Call, t.Result, t.Error)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.tool.call",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"call":            t.Call,
					"result":          t.Result,
					"error":           t.Error,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskChatStreamStart:
		conv, message, err := s.chat.StartStreaming(t.ConversationID)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.stream.start",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"message":         message,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskChatStreamAppend:
		conv, message, err := s.chat.AppendStreaming(t.ConversationID, t.Content)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.stream.delta",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"message":         message,
					"delta":           t.Content,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	case TaskChatStreamFinish:
		conv, message, err := s.chat.FinishStreaming(t.ConversationID, t.FinishReason)
		if err != nil {
			return nil, true, err
		}
		if s.events != nil {
			s.events.Emit(Event{
				Type: "chat.stream.end",
				Data: map[string]any{
					"conversation_id": t.ConversationID,
					"message":         message,
					"finish_reason":   t.FinishReason,
				},
			})
		}
		return conv, true, s.chat.persist(s.configFile)
	default:
		return nil, false, nil
	}
}

func deriveConversationTitle(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return "New conversation"
	}
	runes := []rune(content)
	if len(runes) <= 50 {
		return content
	}
	return strings.TrimSpace(string(runes[:50])) + "..."
}

func conversationMatchesSearchQuery(conv Conversation, query string) bool {
	if query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(conversationSearchHaystack(conv)), query)
}

func conversationSearchHaystack(conv Conversation) string {
	parts := []string{conv.Title, conv.Model}
	for _, message := range conv.Messages {
		parts = append(parts, messageSearchHaystack(message))
	}
	return strings.Join(parts, " ")
}

func messageSearchHaystack(message ChatMessage) string {
	parts := []string{message.Role, message.Content}
	if message.Thinking != nil {
		parts = append(parts, message.Thinking.Content)
	}
	for _, attachment := range message.Attachments {
		parts = append(parts, attachment.Filename, attachment.MimeType)
	}
	for _, invocation := range message.ToolCalls {
		parts = append(parts, invocation.Call.Name, invocation.Result.Content, invocation.Error)
		if len(invocation.Call.Arguments) > 0 {
			if encoded, err := json.Marshal(invocation.Call.Arguments); err == nil {
				parts = append(parts, string(encoded))
			}
		}
	}
	return strings.Join(parts, " ")
}

func buildAssistantPlaceholder(model, prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "Waiting for the local inference pipeline."
	}
	if model == "" {
		model = "local model"
	}
	return "Local inference is not wired in this workspace yet. " +
		"Captured your prompt for " + model + " and stored it in chat history."
}

func roleHeading(role string) string {
	if role == "" {
		return "Message"
	}
	lower := strings.ToLower(role)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

func canReuseAssistantForStreaming(conv Conversation) bool {
	if len(conv.Messages) == 0 {
		return false
	}
	last := conv.Messages[len(conv.Messages)-1]
	if last.Role != "assistant" {
		return false
	}
	if last.Streaming {
		return true
	}
	if len(conv.Messages) < 2 {
		return false
	}
	previous := conv.Messages[len(conv.Messages)-2]
	if previous.Role != "user" {
		return false
	}
	return last.Content == buildAssistantPlaceholder(conv.Model, previous.Content)
}

func parseCounter(value string) uint64 {
	if idx := strings.LastIndex(value, "-"); idx >= 0 && idx < len(value)-1 {
		value = value[idx+1:]
	}
	var counter uint64
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0
		}
		counter = counter*10 + uint64(ch-'0')
	}
	return counter
}

func (s *ChatStore) nextIdentifier(prefix string) string {
	s.nextID++
	return prefix + "-" + strconvFormatUint(s.nextID)
}

func strconvFormatUint(v uint64) string {
	if v == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for v > 0 {
		i--
		digits[i] = byte('0' + v%10)
		v /= 10
	}
	return string(digits[i:])
}

func cloneConversationMap(src map[string]Conversation) map[string]Conversation {
	dst := make(map[string]Conversation, len(src))
	for key, value := range src {
		dst[key] = cloneConversation(value)
	}
	return dst
}

func cloneConversation(conv Conversation) Conversation {
	clone := conv
	clone.Messages = cloneMessages(conv.Messages)
	return clone
}

func cloneMessages(messages []ChatMessage) []ChatMessage {
	clones := make([]ChatMessage, 0, len(messages))
	for _, message := range messages {
		clones = append(clones, cloneMessage(message))
	}
	return clones
}

func cloneMessage(message ChatMessage) ChatMessage {
	clone := message
	clone.Attachments = append([]ImageAttachment(nil), message.Attachments...)
	if message.Thinking != nil {
		copyThinking := *message.Thinking
		clone.Thinking = &copyThinking
	}
	if len(message.ToolCalls) > 0 {
		clone.ToolCalls = make([]ToolInvocation, 0, len(message.ToolCalls))
		for _, invocation := range message.ToolCalls {
			clone.ToolCalls = append(clone.ToolCalls, ToolInvocation{
				Call: ToolCall{
					ID:        invocation.Call.ID,
					Name:      invocation.Call.Name,
					Arguments: cloneAnyMap(invocation.Call.Arguments),
				},
				Result: ToolResult{
					ToolCallID: invocation.Result.ToolCallID,
					Content:    invocation.Result.Content,
				},
				StartedAt: invocation.StartedAt,
				EndedAt:   invocation.EndedAt,
				Error:     invocation.Error,
			})
		}
	}
	return clone
}

func cloneAttachmentMap(src map[string][]ImageAttachment) map[string][]ImageAttachment {
	dst := make(map[string][]ImageAttachment, len(src))
	for key, value := range src {
		dst[key] = append([]ImageAttachment(nil), value...)
	}
	return dst
}

func cloneThinkingMap(src map[string]ThinkingState) map[string]ThinkingState {
	dst := make(map[string]ThinkingState, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func cloneStringMap(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func cloneAnyMap(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]any, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func latestStreamingAssistantIndex(messages []ChatMessage) int {
	for index := len(messages) - 1; index >= 0; index-- {
		if messages[index].Role == "assistant" && messages[index].Streaming {
			return index
		}
	}
	return -1
}

func latestAssistantIndex(messages []ChatMessage) int {
	for index := len(messages) - 1; index >= 0; index-- {
		if messages[index].Role == "assistant" {
			return index
		}
	}
	return -1
}

func (s *ChatStore) syncThinkingToStreamingMessageLocked(conversationID string, state ThinkingState) {
	conv, ok := s.conversations[conversationID]
	if !ok {
		return
	}
	index := latestStreamingAssistantIndex(conv.Messages)
	if index < 0 {
		return
	}
	copyThinking := state
	conv.Messages[index].Thinking = &copyThinking
	s.conversations[conversationID] = conv
}

func (s *ChatStore) ensureAttachmentIdentifiersLocked() {
	for conversationID, attachments := range s.queuedImages {
		for index := range attachments {
			if attachments[index].ID == "" {
				attachments[index].ID = s.nextIdentifier("img")
			}
		}
		s.queuedImages[conversationID] = attachments
	}
	for conversationID, conv := range s.conversations {
		changed := false
		for messageIndex := range conv.Messages {
			for attachmentIndex := range conv.Messages[messageIndex].Attachments {
				if conv.Messages[messageIndex].Attachments[attachmentIndex].ID != "" {
					continue
				}
				conv.Messages[messageIndex].Attachments[attachmentIndex].ID = s.nextIdentifier("img")
				changed = true
			}
		}
		if changed {
			s.conversations[conversationID] = conv
		}
	}
}

func (s *ChatStore) nextResponseModelLocked(conv Conversation) string {
	if s.selectedModel != "" {
		return s.selectedModel
	}
	if conv.Model != "" {
		return conv.Model
	}
	return s.settings.DefaultModel
}

func (s *ChatStore) modelSupportsVisionLocked(name string) bool {
	for _, model := range s.models {
		if model.Name == name {
			return model.SupportsVision
		}
	}
	return false
}

func (s *ChatStore) normalizeMessagesLocked(messages []ChatMessage, base time.Time) []ChatMessage {
	if base.IsZero() {
		base = time.Now().UTC()
	}
	normalized := make([]ChatMessage, 0, len(messages))
	for index, message := range messages {
		if message.ID == "" {
			message.ID = s.nextIdentifier("msg")
		}
		if message.CreatedAt.IsZero() {
			message.CreatedAt = base.Add(time.Duration(index) * time.Millisecond)
		}
		for attachmentIndex := range message.Attachments {
			if message.Attachments[attachmentIndex].ID == "" {
				message.Attachments[attachmentIndex].ID = s.nextIdentifier("img")
			}
		}
		normalized = append(normalized, message)
	}
	return normalized
}

func firstUserMessage(messages []ChatMessage) string {
	for _, message := range messages {
		if message.Role == "user" && strings.TrimSpace(message.Content) != "" {
			return message.Content
		}
	}
	return ""
}

func validateImageAttachment(attachment ImageAttachment) error {
	if strings.TrimSpace(attachment.Filename) == "" {
		return coreerr.E("display.chat.validateImageAttachment", "attachment filename is required", nil)
	}
	if strings.TrimSpace(attachment.Data) == "" {
		return coreerr.E("display.chat.validateImageAttachment", "attachment data is required", nil)
	}
	switch strings.TrimSpace(strings.ToLower(attachment.MimeType)) {
	case "image/png", "image/jpeg", "image/webp", "image/gif":
		return nil
	default:
		return coreerr.E("display.chat.validateImageAttachment", "unsupported image type: "+attachment.MimeType, nil)
	}
}
