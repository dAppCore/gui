package chat

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/core/io/store"
	coreerr "dappco.re/go/core/log"
	guimcp "forge.lthn.ai/core/gui/pkg/mcp"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"
)

const (
	conversationsGroup = "chat_conversations"
	settingsGroup      = "chat_settings"
	settingsKey        = "global"
)

type Options struct {
	APIURL       string
	StorePath    string
	HTTPClient   *http.Client
	ModelRoots   []string
	ToolExecutor ToolExecutor
	Now          func() time.Time
}

type Service struct {
	*core.ServiceRuntime[Options]
	options            Options
	store              *store.KeyValueStore
	httpClient         *http.Client
	toolExecutor       ToolExecutor
	toolHandler        *ToolCallHandler
	pendingAttachments map[string][]ImageAttachment
	mu                 sync.Mutex
}

type sendInput struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Content        string `json:"content"`
}

type conversationInput struct {
	ID             string `json:"id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
}

type searchInput struct {
	Query string `json:"q"`
}

type renameInput struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type thinkingInput struct {
	ConversationID string    `json:"conversation_id,omitempty"`
	MessageID      string    `json:"message_id,omitempty"`
	Content        string    `json:"content,omitempty"`
	StartedAt      time.Time `json:"started_at,omitempty"`
	DurationMS     int64     `json:"duration_ms,omitempty"`
}

type selectModelInput struct {
	Model          string `json:"model"`
	ConversationID string `json:"conversation_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type attachImageInput struct {
	ConversationID  string `json:"conversation_id,omitempty"`
	ImageAttachment `json:",inline"`
}

type openAIRequest struct {
	Model       string           `json:"model"`
	Messages    []openAIMessage  `json:"messages"`
	Temperature float32          `json:"temperature,omitempty"`
	TopP        float32          `json:"top_p,omitempty"`
	TopK        int              `json:"top_k,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Stream      bool             `json:"stream"`
	Tools       []openAIToolSpec `json:"tools,omitempty"`
	ToolChoice  string           `json:"tool_choice,omitempty"`
}

type openAIMessage struct {
	Role       string           `json:"role"`
	Content    any              `json:"content,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type openAIToolSpec struct {
	Type     string             `json:"type"`
	Function openAIFunctionSpec `json:"function"`
}

type openAIFunctionSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type openAIToolCall struct {
	ID       string             `json:"id,omitempty"`
	Type     string             `json:"type"`
	Function openAIFunctionCall `json:"function"`
}

type openAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type configuredModels struct {
	Default      string `yaml:"default"`
	DefaultModel string `yaml:"default_model"`
	Models       []struct {
		Name         string `yaml:"name"`
		Path         string `yaml:"path"`
		Architecture string `yaml:"architecture"`
		Backend      string `yaml:"backend"`
	} `yaml:"models"`
}

func Register(optionFns ...func(*Options)) func(*core.Core) core.Result {
	options := Options{
		APIURL:     "http://localhost:8090",
		StorePath:  filepath.Join(core.Env("DIR_HOME"), ".core", "gui", "chat.db"),
		HTTPClient: &http.Client{Timeout: 5 * time.Minute},
		ModelRoots: defaultModelRoots(),
		Now:        time.Now,
	}
	for _, fn := range optionFns {
		fn(&options)
	}
	return func(c *core.Core) core.Result {
		svc := &Service{
			ServiceRuntime:     core.NewServiceRuntime[Options](c, options),
			options:            options,
			httpClient:         options.HTTPClient,
			pendingAttachments: make(map[string][]ImageAttachment),
		}
		return core.Result{Value: svc, OK: true}
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	if err := os.MkdirAll(filepath.Dir(s.options.StorePath), 0o755); err != nil {
		return core.Result{Value: err, OK: false}
	}

	keyValueStore, err := store.New(store.Options{Path: s.options.StorePath})
	if err != nil {
		return core.Result{Value: err, OK: false}
	}
	s.store = keyValueStore

	s.toolExecutor = s.options.ToolExecutor
	if s.toolExecutor == nil {
		subsystem := guimcp.New(s.Core())
		server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "coregui-chat", Version: "0.1.0"}, nil)
		subsystem.RegisterTools(server)
		s.toolExecutor = subsystem
	}
	s.toolHandler = NewToolCallHandler(s.toolExecutor)
	s.Core().RegisterQuery(s.handleQuery)
	s.registerActions()
	return core.Result{OK: true}
}

func (s *Service) HandleIPCEvents(_ *core.Core, _ core.Message) core.Result {
	return core.Result{OK: true}
}

func (s *Service) handleQuery(_ *core.Core, q core.Query) core.Result {
	switch typed := q.(type) {
	case QueryHistory:
		conv, err := s.getConversation(typed.ID, typed.ConversationID)
		return core.Result{}.New(conv, err)
	case QueryModels:
		return core.Result{Value: s.discoverModels(), OK: true}
	case QuerySettings:
		return core.Result{Value: s.loadSettings(), OK: true}
	case QueryConversationList:
		conversations, err := s.listConversationSummaries()
		return core.Result{}.New(conversations, err)
	case QueryConversationGet:
		conv, err := s.getConversation(typed.ID, typed.ConversationID)
		return core.Result{}.New(conv, err)
	case QueryConversationSearch:
		results, err := s.searchConversationSummaries(typed.Query)
		return core.Result{}.New(results, err)
	default:
		return core.Result{}
	}
}

func (s *Service) registerActions() {
	c := s.Core()
	c.Action("gui.chat.send", func(ctx context.Context, opts core.Options) core.Result {
		input, err := decodeInput[sendInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conv, err := s.send(ctx, input)
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.clear", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[conversationInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conv, err := s.clearConversation(input.ID, input.ConversationID)
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.history", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[conversationInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conv, err := s.getConversation(input.ID, input.ConversationID)
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.discoverModels(), OK: true}
	})
	c.Action("gui.chat.selectModel", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[selectModelInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		settings, err := s.selectModel(input)
		return core.Result{}.New(settings, err)
	})
	c.Action("gui.chat.settings.save", func(_ context.Context, opts core.Options) core.Result {
		settings, err := decodeInput[ChatSettings](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		err = s.saveSettings(settings)
		return core.Result{}.New(settings, err)
	})
	c.Action("gui.chat.settings.load", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.loadSettings(), OK: true}
	})
	c.Action("gui.chat.settings.reset", func(_ context.Context, _ core.Options) core.Result {
		settings := DefaultSettings()
		err := s.saveSettings(settings)
		return core.Result{}.New(settings, err)
	})
	c.Action("gui.chat.conversations.list", func(_ context.Context, _ core.Options) core.Result {
		conversations, err := s.listConversationSummaries()
		return core.Result{}.New(conversations, err)
	})
	c.Action("gui.chat.conversations.get", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[conversationInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conv, err := s.getConversation(input.ID, input.ConversationID)
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.conversations.delete", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[conversationInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		err = s.deleteConversation(coalesce(input.ID, input.ConversationID))
		return core.Result{}.New(nil, err)
	})
	c.Action("gui.chat.conversations.search", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[searchInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		results, err := s.searchConversationSummaries(input.Query)
		return core.Result{}.New(results, err)
	})
	c.Action("gui.chat.conversations.new", func(_ context.Context, _ core.Options) core.Result {
		conv, err := s.createConversation()
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.conversation.save", func(_ context.Context, opts core.Options) core.Result {
		conv, err := decodeInput[Conversation](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if strings.TrimSpace(conv.ID) == "" {
			return core.Result{Value: coreerr.E("chat.conversation.save", "conversation id is required", nil), OK: false}
		}
		saved, err := s.saveConversation(conv)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		s.emit(ActionConversationUpdated{Conversation: saved})
		return core.Result{Value: saved, OK: true}
	})
	c.Action("gui.chat.conversations.rename", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[renameInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		conv, err := s.renameConversation(input.ID, input.Title)
		return core.Result{}.New(conv, err)
	})
	c.Action("gui.chat.conversations.export", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[conversationInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		markdown, err := s.exportConversation(coalesce(input.ID, input.ConversationID))
		return core.Result{}.New(markdown, err)
	})
	c.Action("gui.chat.attachImage", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[attachImageInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		s.queueAttachment(coalesce(input.ConversationID, "draft"), input.ImageAttachment)
		return core.Result{Value: input.ImageAttachment, OK: true}
	})
	c.Action("gui.chat.thinking.start", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[thinkingInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if strings.TrimSpace(input.ConversationID) == "" {
			return core.Result{Value: coreerr.E("chat.thinking.start", "conversation id is required", nil), OK: false}
		}
		s.emit(ActionThinkingStarted{
			ConversationID: input.ConversationID,
			MessageID:      input.MessageID,
			StartedAt:      input.StartedAt,
		})
		return core.Result{Value: input, OK: true}
	})
	c.Action("gui.chat.thinking.append", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[thinkingInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if strings.TrimSpace(input.ConversationID) == "" {
			return core.Result{Value: coreerr.E("chat.thinking.append", "conversation id is required", nil), OK: false}
		}
		s.emit(ActionThinkingAppended{
			ConversationID: input.ConversationID,
			MessageID:      input.MessageID,
			Content:        input.Content,
		})
		return core.Result{Value: input.Content, OK: true}
	})
	c.Action("gui.chat.thinking.end", func(_ context.Context, opts core.Options) core.Result {
		input, err := decodeInput[thinkingInput](opts)
		if err != nil {
			return core.Result{Value: err, OK: false}
		}
		if strings.TrimSpace(input.ConversationID) == "" {
			return core.Result{Value: coreerr.E("chat.thinking.end", "conversation id is required", nil), OK: false}
		}
		duration := input.DurationMS
		if duration == 0 && !input.StartedAt.IsZero() {
			duration = time.Since(input.StartedAt).Milliseconds()
			if duration < 0 {
				duration = 0
			}
		}
		s.emit(ActionThinkingEnded{
			ConversationID: input.ConversationID,
			MessageID:      input.MessageID,
			DurationMS:     duration,
		})
		return core.Result{Value: duration, OK: true}
	})
}

func decodeInput[T any](opts core.Options) (T, error) {
	var input T
	if task := opts.Get("task"); task.OK {
		if typed, ok := task.Value.(T); ok {
			return typed, nil
		}
	}

	items := make(map[string]any, opts.Len())
	for _, item := range opts.Items() {
		items[item.Key] = item.Value
	}
	if len(items) == 0 {
		return input, nil
	}

	result := core.JSONUnmarshalString(core.JSONMarshalString(items), &input)
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return input, err
		}
		return input, coreerr.E("chat.decodeInput", "failed to decode action input", nil)
	}
	return input, nil
}

func defaultModelRoots() []string {
	roots := []string{filepath.Join(core.Env("DIR_HOME"), ".core", "models")}
	if env := strings.TrimSpace(core.Env("CORE_MODELS_DIR")); env != "" {
		roots = append(roots, env)
	}
	return roots
}

func (s *Service) now() time.Time {
	if s.options.Now != nil {
		return s.options.Now()
	}
	return time.Now()
}

func (s *Service) saveSettings(settings ChatSettings) error {
	payload := core.JSONMarshalString(settings)
	return s.store.Set(settingsGroup, settingsKey, payload)
}

func (s *Service) loadSettings() ChatSettings {
	settings := DefaultSettings()
	if s.store == nil {
		return settings
	}
	payload, err := s.store.Get(settingsGroup, settingsKey)
	if err != nil {
		return settings
	}
	_ = core.JSONUnmarshalString(payload, &settings)
	return settings
}

func (s *Service) selectModel(input selectModelInput) (ChatSettings, error) {
	settings := s.loadSettings()
	settings.DefaultModel = input.Model
	if err := s.saveSettings(settings); err != nil {
		return ChatSettings{}, err
	}

	targetConversation := coalesce(input.ConversationID, input.ID)
	if targetConversation == "" {
		return settings, nil
	}

	conv, err := s.loadConversation(targetConversation)
	if err != nil {
		return ChatSettings{}, err
	}
	conv.Model = input.Model
	conv, err = s.saveConversation(conv)
	if err != nil {
		return ChatSettings{}, err
	}
	s.emit(ActionConversationUpdated{Conversation: conv})
	return settings, nil
}

func (s *Service) saveConversation(conv Conversation) (Conversation, error) {
	if conv.CreatedAt.IsZero() {
		conv.CreatedAt = s.now()
	}
	if strings.TrimSpace(conv.Title) == "" {
		if len(conv.Messages) > 0 && strings.TrimSpace(conv.Messages[0].Content) != "" {
			conv.Title = titleFrom(conv.Messages[0].Content)
		} else {
			conv.Title = "New Chat"
		}
	}
	conv.UpdatedAt = s.now()
	payload := core.JSONMarshalString(conv)
	return conv, s.store.Set(conversationsGroup, conv.ID, payload)
}

func (s *Service) loadConversation(id string) (Conversation, error) {
	payload, err := s.store.Get(conversationsGroup, id)
	if err != nil {
		return Conversation{}, err
	}
	var conv Conversation
	result := core.JSONUnmarshalString(payload, &conv)
	if !result.OK {
		if decodeErr, ok := result.Value.(error); ok {
			return Conversation{}, decodeErr
		}
		return Conversation{}, coreerr.E("chat.loadConversation", "failed to decode conversation", nil)
	}
	return conv, nil
}

func (s *Service) listConversationSummaries() ([]ConversationSummary, error) {
	if s.store == nil {
		return nil, nil
	}
	items, err := s.store.GetAll(conversationsGroup)
	if err != nil {
		return nil, err
	}
	summaries := make([]ConversationSummary, 0, len(items))
	for _, payload := range items {
		var conv Conversation
		if result := core.JSONUnmarshalString(payload, &conv); result.OK {
			summaries = append(summaries, conv.Summary())
		}
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
	})
	return summaries, nil
}

func (s *Service) searchConversationSummaries(query string) ([]ConversationSummary, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	summaries, err := s.listConversationSummaries()
	if err != nil || query == "" {
		return summaries, err
	}

	items, err := s.store.GetAll(conversationsGroup)
	if err != nil {
		return nil, err
	}
	matches := make([]ConversationSummary, 0)
	for _, payload := range items {
		var conv Conversation
		if result := core.JSONUnmarshalString(payload, &conv); !result.OK {
			continue
		}
		if strings.Contains(strings.ToLower(conv.Title), query) {
			matches = append(matches, conv.Summary())
			continue
		}
		for _, message := range conv.Messages {
			if strings.Contains(strings.ToLower(message.Content), query) {
				matches = append(matches, conv.Summary())
				break
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].UpdatedAt.After(matches[j].UpdatedAt)
	})
	return matches, nil
}

func (s *Service) createConversation() (Conversation, error) {
	settings := s.loadSettings()
	now := s.now()
	conv := Conversation{
		ID:        "conv-" + strconv.FormatInt(now.UnixNano(), 36),
		Title:     "New Chat",
		Model:     settings.DefaultModel,
		CreatedAt: now,
		UpdatedAt: now,
		Messages:  nil,
	}
	conv, err := s.saveConversation(conv)
	if err != nil {
		return Conversation{}, err
	}
	s.emit(ActionConversationCreated{Conversation: conv})
	return conv, nil
}

func (s *Service) getConversation(id, conversationID string) (Conversation, error) {
	target := coalesce(id, conversationID)
	if target == "" {
		return Conversation{}, coreerr.E("chat.getConversation", "conversation id is required", nil)
	}
	return s.loadConversation(target)
}

func (s *Service) renameConversation(id, title string) (Conversation, error) {
	conv, err := s.loadConversation(id)
	if err != nil {
		return Conversation{}, err
	}
	conv.Title = strings.TrimSpace(title)
	if conv.Title == "" {
		conv.Title = "New Chat"
	}
	conv, err = s.saveConversation(conv)
	if err != nil {
		return Conversation{}, err
	}
	s.emit(ActionConversationUpdated{Conversation: conv})
	return conv, nil
}

func (s *Service) clearConversation(id, conversationID string) (Conversation, error) {
	conv, err := s.getConversation(id, conversationID)
	if err != nil {
		return Conversation{}, err
	}
	conv.Messages = nil
	conv.UpdatedAt = s.now()
	conv, err = s.saveConversation(conv)
	if err != nil {
		return Conversation{}, err
	}
	s.clearQueuedAttachments(conv.ID)
	s.emit(ActionConversationCleared{ConversationID: conv.ID})
	s.emit(ActionConversationUpdated{Conversation: conv})
	return conv, nil
}

func (s *Service) deleteConversation(id string) error {
	if id == "" {
		return coreerr.E("chat.deleteConversation", "conversation id is required", nil)
	}
	if err := s.store.Delete(conversationsGroup, id); err != nil {
		return err
	}
	s.clearQueuedAttachments(id)
	s.emit(ActionConversationDeleted{ConversationID: id})
	return nil
}

func (s *Service) exportConversation(id string) (string, error) {
	conv, err := s.loadConversation(id)
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(conv.Title)
	builder.WriteString("\n\n")
	for _, message := range conv.Messages {
		builder.WriteString("## ")
		builder.WriteString(strings.Title(message.Role))
		builder.WriteString("\n\n")
		if message.Content != "" {
			builder.WriteString(message.Content)
			builder.WriteString("\n\n")
		}
		for _, result := range message.ToolResults {
			builder.WriteString("> Tool Result (`")
			builder.WriteString(result.ToolCallID)
			builder.WriteString("`)\n>\n> ")
			builder.WriteString(strings.ReplaceAll(result.Content, "\n", "\n> "))
			builder.WriteString("\n\n")
		}
	}
	return builder.String(), nil
}

func (s *Service) queueAttachment(conversationID string, attachment ImageAttachment) {
	key := coalesce(conversationID, "draft")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingAttachments[key] = append(s.pendingAttachments[key], attachment)
	s.emit(ActionImageQueued{ConversationID: key, Attachment: attachment})
}

func (s *Service) drainAttachments(conversationID string) []ImageAttachment {
	s.mu.Lock()
	defer s.mu.Unlock()
	var attachments []ImageAttachment
	keys := []string{conversationID}
	if conversationID != "draft" {
		keys = append(keys, "draft")
	}
	for _, key := range keys {
		if key == "" {
			continue
		}
		attachments = append(attachments, s.pendingAttachments[key]...)
		delete(s.pendingAttachments, key)
	}
	return attachments
}

func (s *Service) clearQueuedAttachments(conversationID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pendingAttachments, conversationID)
	if conversationID != "draft" {
		delete(s.pendingAttachments, "draft")
	}
}

func (s *Service) mergedSettings(global ChatSettings, override *ChatSettings) ChatSettings {
	if override == nil {
		return global
	}
	merged := global
	if override.Temperature != 0 {
		merged.Temperature = override.Temperature
	}
	if override.TopP != 0 {
		merged.TopP = override.TopP
	}
	if override.TopK != 0 {
		merged.TopK = override.TopK
	}
	if override.MaxTokens != 0 {
		merged.MaxTokens = override.MaxTokens
	}
	if override.ContextWindow != 0 {
		merged.ContextWindow = override.ContextWindow
	}
	if override.SystemPrompt != "" {
		merged.SystemPrompt = override.SystemPrompt
	}
	if override.DefaultModel != "" {
		merged.DefaultModel = override.DefaultModel
	}
	return merged
}

func (s *Service) send(ctx context.Context, input sendInput) (Conversation, error) {
	if strings.TrimSpace(input.Content) == "" && len(s.pendingAttachments) == 0 {
		return Conversation{}, coreerr.E("chat.send", "message content is required", nil)
	}

	settings := s.loadSettings()
	var (
		conv    Conversation
		err     error
		created bool
	)
	if input.ConversationID != "" {
		conv, err = s.loadConversation(input.ConversationID)
	} else {
		conv, err = s.createConversation()
		created = true
	}
	if err != nil {
		return Conversation{}, err
	}

	attachments := s.drainAttachments(conv.ID)
	if len(attachments) == 0 && created {
		attachments = s.drainAttachments("draft")
	}

	now := s.now()
	conv.Model = s.resolveModel(conv.Model, settings.DefaultModel)
	userMessage := ChatMessage{
		ID:          "msg-" + strconv.FormatInt(now.UnixNano(), 36),
		Role:        "user",
		Content:     input.Content,
		CreatedAt:   now,
		Model:       conv.Model,
		Attachments: attachments,
	}
	conv.Messages = append(conv.Messages, userMessage)
	if conv.Title == "" || conv.Title == "New Chat" {
		conv.Title = titleFrom(input.Content)
	}
	conv, err = s.saveConversation(conv)
	if err != nil {
		return Conversation{}, err
	}
	s.emit(ActionMessageAdded{ConversationID: conv.ID, Message: userMessage})
	s.emit(ActionConversationUpdated{Conversation: conv})

	for toolRound := 0; toolRound < 3; toolRound++ {
		effectiveSettings := s.mergedSettings(settings, conv.Settings)
		conv.Model = s.resolveModel(conv.Model, effectiveSettings.DefaultModel)

		assistantMessage, err := s.streamAssistant(ctx, conv, effectiveSettings)
		if err != nil {
			return conv, err
		}
		if hasRenderableContent(assistantMessage) {
			conv.Messages = append(conv.Messages, assistantMessage)
			conv, err = s.saveConversation(conv)
			if err != nil {
				return conv, err
			}
			s.emit(ActionMessageAdded{ConversationID: conv.ID, Message: assistantMessage})
			s.emit(ActionConversationUpdated{Conversation: conv})
		}

		if len(assistantMessage.ToolCalls) == 0 {
			break
		}

		results := s.toolHandler.ExecuteAll(ctx, assistantMessage.ToolCalls)
		for _, result := range results {
			toolMessage := ChatMessage{
				ID:          "tool-" + strconv.FormatInt(s.now().UnixNano(), 36),
				Role:        "tool",
				Content:     result.Content,
				CreatedAt:   s.now(),
				Model:       conv.Model,
				ToolCallID:  result.ToolCallID,
				ToolResults: []ToolResult{result},
			}
			conv.Messages = append(conv.Messages, toolMessage)
			s.emit(ActionToolResultReady{ConversationID: conv.ID, MessageID: assistantMessage.ID, Result: result})
			s.emit(ActionMessageAdded{ConversationID: conv.ID, Message: toolMessage})
		}
		conv, err = s.saveConversation(conv)
		if err != nil {
			return conv, err
		}
		s.emit(ActionConversationUpdated{Conversation: conv})
	}

	return conv, nil
}

func (s *Service) streamAssistant(ctx context.Context, conv Conversation, settings ChatSettings) (ChatMessage, error) {
	messageID := "msg-" + strconv.FormatInt(s.now().UnixNano(), 36)
	requestBody := s.buildCompletionRequest(conv, settings)
	payload := core.JSONMarshalString(requestBody)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(s.options.APIURL, "/")+"/v1/chat/completions", bytes.NewBufferString(payload))
	if err != nil {
		return ChatMessage{}, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return ChatMessage{}, err
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(response.Body)
		return ChatMessage{}, coreerr.E("chat.streamAssistant", strings.TrimSpace(string(body)), nil)
	}

	renderer := NewStreamRenderer(StreamCallbacks{
		OnStart: func(streamID string) {
			s.emit(ActionStreamStarted{ConversationID: conv.ID, MessageID: messageID, StreamID: streamID})
		},
		OnToken: func(content string) {
			s.emit(ActionTokenAppended{ConversationID: conv.ID, MessageID: messageID, Content: content})
		},
		OnThinkingStart: func(state ThinkingState) {
			s.emit(ActionThinkingStarted{ConversationID: conv.ID, MessageID: messageID, StartedAt: state.StartedAt})
		},
		OnThinkingAppend: func(content string) {
			s.emit(ActionThinkingAppended{ConversationID: conv.ID, MessageID: messageID, Content: content})
		},
		OnThinkingEnd: func(state ThinkingState) {
			s.emit(ActionThinkingEnded{ConversationID: conv.ID, MessageID: messageID, DurationMS: state.DurationMS})
		},
		OnToolCall: func(call ToolCall) {
			s.emit(ActionToolCallStarted{ConversationID: conv.ID, MessageID: messageID, Call: call})
		},
		OnFinish: func(reason string) {
			s.emit(ActionStreamFinished{ConversationID: conv.ID, MessageID: messageID, FinishReason: reason})
		},
	})
	if err := renderer.Render(response.Body); err != nil {
		return ChatMessage{}, err
	}
	return renderer.Message(messageID, conv.Model, s.now()), nil
}

func (s *Service) buildCompletionRequest(conv Conversation, settings ChatSettings) openAIRequest {
	request := openAIRequest{
		Model:       s.resolveModel(conv.Model, settings.DefaultModel),
		Messages:    make([]openAIMessage, 0, len(conv.Messages)+1),
		Temperature: settings.Temperature,
		TopP:        settings.TopP,
		TopK:        settings.TopK,
		MaxTokens:   settings.MaxTokens,
		Stream:      true,
	}

	systemPrompt := strings.TrimSpace(settings.SystemPrompt)
	if s.toolExecutor != nil {
		manifest := s.toolExecutor.ManifestText()
		if manifest != "" {
			if systemPrompt != "" {
				systemPrompt += "\n\n"
			}
			systemPrompt += manifest + "\nUse tools when helpful. When a tool is needed, emit a tool call with valid JSON arguments."
			for _, tool := range s.toolExecutor.Manifest() {
				request.Tools = append(request.Tools, openAIToolSpec{
					Type: "function",
					Function: openAIFunctionSpec{
						Name:        tool.Name,
						Description: tool.Description,
						Parameters:  tool.InputSchema,
					},
				})
			}
			request.ToolChoice = "auto"
		}
	}
	if systemPrompt != "" {
		request.Messages = append(request.Messages, openAIMessage{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	for _, message := range conv.Messages {
		apiMessage := openAIMessage{Role: message.Role}
		switch message.Role {
		case "user":
			apiMessage.Content = renderUserContent(message)
		case "assistant":
			apiMessage.Content = message.Content
			if len(message.ToolCalls) > 0 {
				apiMessage.ToolCalls = make([]openAIToolCall, 0, len(message.ToolCalls))
				for _, call := range message.ToolCalls {
					apiMessage.ToolCalls = append(apiMessage.ToolCalls, openAIToolCall{
						ID:   call.ID,
						Type: "function",
						Function: openAIFunctionCall{
							Name:      call.Name,
							Arguments: core.JSONMarshalString(call.Arguments),
						},
					})
				}
			}
		case "tool":
			apiMessage.Content = message.Content
			apiMessage.ToolCallID = message.ToolCallID
		default:
			apiMessage.Content = message.Content
		}
		request.Messages = append(request.Messages, apiMessage)
	}

	return request
}

func (s *Service) resolveModel(current, configured string) string {
	if current = strings.TrimSpace(current); current != "" {
		return current
	}
	if configured = strings.TrimSpace(configured); configured != "" {
		return configured
	}
	models := s.discoverModels()
	if len(models) > 0 {
		return models[0].Name
	}
	return "default"
}

func renderUserContent(message ChatMessage) any {
	if len(message.Attachments) == 0 {
		return message.Content
	}
	parts := []map[string]any{
		{"type": "text", "text": message.Content},
	}
	for _, attachment := range message.Attachments {
		parts = append(parts, map[string]any{
			"type": "image_url",
			"image_url": map[string]any{
				"url": "data:" + attachment.MimeType + ";base64," + attachment.Data,
			},
		})
	}
	return parts
}

func hasRenderableContent(message ChatMessage) bool {
	return strings.TrimSpace(message.Content) != "" || message.Thinking != nil || len(message.ToolCalls) > 0
}

func titleFrom(content string) string {
	title := strings.TrimSpace(content)
	if title == "" {
		return "New Chat"
	}
	runes := []rune(title)
	if len(runes) > 50 {
		return string(runes[:50])
	}
	return title
}

func coalesce(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (s *Service) emit(message any) {
	if message == nil {
		return
	}
	_ = s.Core().ACTION(message)
}

func (s *Service) discoverModels() []ModelEntry {
	settings := s.loadSettings()
	models := map[string]ModelEntry{}
	for _, root := range s.options.ModelRoots {
		for _, model := range discoverModelsOnDisk(root) {
			models[model.Name] = model
		}
	}

	configPath := filepath.Join(core.Env("DIR_HOME"), ".core", "models.yaml")
	if payload, err := os.ReadFile(configPath); err == nil {
		var configured configuredModels
		if err := yaml.Unmarshal(payload, &configured); err == nil {
			defaultModel := coalesce(configured.DefaultModel, configured.Default, settings.DefaultModel)
			for _, item := range configured.Models {
				name := coalesce(item.Name, filepath.Base(item.Path))
				entry := models[name]
				entry.Name = name
				entry.Architecture = coalesce(item.Architecture, entry.Architecture)
				entry.Backend = coalesce(item.Backend, entry.Backend, "local")
				if entry.SizeBytes == 0 && item.Path != "" {
					entry.SizeBytes = directorySize(item.Path)
				}
				entry.Loaded = name == defaultModel
				models[name] = entry
			}
		}
	}

	names := make([]string, 0, len(models))
	for name := range models {
		names = append(names, name)
	}
	slices.Sort(names)
	results := make([]ModelEntry, 0, len(names))
	for _, name := range names {
		results = append(results, models[name])
	}
	return results
}

func discoverModelsOnDisk(root string) []ModelEntry {
	if strings.TrimSpace(root) == "" {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var results []ModelEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		modelPath := filepath.Join(root, entry.Name())
		configPath := filepath.Join(modelPath, "config.json")
		if _, err := os.Stat(configPath); err != nil {
			continue
		}
		results = append(results, ModelEntry{
			Name:         entry.Name(),
			Architecture: architectureFromConfig(configPath),
			QuantBits:    quantBitsFromName(entry.Name()),
			SizeBytes:    directorySize(modelPath),
			Backend:      "local",
		})
	}
	return results
}

func architectureFromConfig(configPath string) string {
	payload, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	var parsed map[string]any
	if result := core.JSONUnmarshalString(string(payload), &parsed); !result.OK {
		return ""
	}
	if architectures, ok := parsed["architectures"].([]any); ok && len(architectures) > 0 {
		if name, ok := architectures[0].(string); ok {
			return strings.ToLower(name)
		}
	}
	if modelType, ok := parsed["model_type"].(string); ok {
		return strings.ToLower(modelType)
	}
	return ""
}

func quantBitsFromName(name string) int {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "q4"), strings.Contains(lower, "4bit"):
		return 4
	case strings.Contains(lower, "q8"), strings.Contains(lower, "8bit"):
		return 8
	default:
		return 0
	}
}

func directorySize(root string) int64 {
	var total int64
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".safetensors") || strings.HasSuffix(info.Name(), ".gguf") || strings.HasSuffix(info.Name(), ".bin") {
			total += info.Size()
		}
		return nil
	})
	return total
}
