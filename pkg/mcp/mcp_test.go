// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	"dappco.re/go/core/gui/pkg/clipboard"
	"dappco.re/go/core/gui/pkg/display"
	"dappco.re/go/core/gui/pkg/environment"
	"dappco.re/go/core/gui/pkg/menu"
	"dappco.re/go/core/gui/pkg/notification"
	"dappco.re/go/core/gui/pkg/screen"
	"dappco.re/go/core/gui/pkg/webview"
	"dappco.re/go/core/gui/pkg/window"
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsystem_Good_Name(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := NewSubsystem(c)
	assert.Equal(t, "display", sub.Name())
}

func TestSubsystem_Good_RegisterTools(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := NewSubsystem(c)
	// RegisterTools should not panic with a real mcp.Server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	assert.NotPanics(t, func() { sub.RegisterTools(server) })
}

// Integration test: verify the IPC round-trip that MCP tool handlers use.

type mockClipPlatform struct {
	text     string
	ok       bool
	image    []byte
	hasImage bool
}

func (m *mockClipPlatform) Text() (string, bool)  { return m.text, m.ok }
func (m *mockClipPlatform) SetText(t string) bool { m.text = t; m.ok = t != ""; return true }
func (m *mockClipPlatform) Image() ([]byte, bool) { return m.image, m.hasImage }
func (m *mockClipPlatform) SetImage(data []byte) bool {
	m.image = append([]byte(nil), data...)
	m.hasImage = len(data) > 0
	return true
}

func TestMCP_Good_ClipboardRoundTrip(t *testing.T) {
	c, err := core.New(
		core.WithService(clipboard.Register(&mockClipPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	// Verify the IPC path that clipboard_read tool handler uses
	result, handled, err := c.QUERY(clipboard.QueryText{})
	require.NoError(t, err)
	assert.True(t, handled)
	content, ok := result.(clipboard.ClipboardContent)
	require.True(t, ok, "expected ClipboardContent type")
	assert.Equal(t, "hello", content.Text)
}

func TestMCP_Good_DialogMessage(t *testing.T) {
	mock := &mockNotificationPlatform{}
	c, err := core.New(
		core.WithService(notification.Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)
	_, result, err := sub.dialogMessage(context.Background(), nil, DialogMessageInput{
		Title:   "Alias",
		Message: "Hello",
		Kind:    "error",
	})
	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, mock.sendCalled)
	assert.Equal(t, notification.SeverityError, mock.lastOpts.Severity)
}

func TestMCP_Good_ThemeSetString(t *testing.T) {
	mock := &mockEnvironmentPlatform{isDark: true}
	c, err := core.New(
		core.WithService(environment.Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)
	_, result, err := sub.themeSet(context.Background(), nil, ThemeSetInput{Theme: "light"})
	require.NoError(t, err)
	assert.Equal(t, "light", result.Theme.Theme)
	assert.False(t, result.Theme.IsDark)
	assert.False(t, mock.isDark)
}

func TestMCP_Good_WindowTitleSetAlias(t *testing.T) {
	c, err := core.New(
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, err := c.PERFORM(window.TaskOpenWindow{
		Window: &window.Window{Name: "alias-win", Title: "Original", URL: "/"},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	sub := New(c)
	_, result, err := sub.windowTitleSet(context.Background(), nil, WindowTitleInput{
		Name:  "alias-win",
		Title: "Updated",
	})
	require.NoError(t, err)
	assert.True(t, result.Success)

	queried, handled, err := c.QUERY(window.QueryWindowByName{Name: "alias-win"})
	require.NoError(t, err)
	assert.True(t, handled)
	info, ok := queried.(*window.WindowInfo)
	require.True(t, ok)
	require.NotNil(t, info)
	assert.Equal(t, "Updated", info.Title)
}

func TestMCP_Good_ChatRoundTrip(t *testing.T) {
	c, err := core.New(
		core.WithService(display.Register(nil)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)

	_, created, err := sub.chatConversationNew(context.Background(), nil, ChatConversationNewInput{})
	require.NoError(t, err)
	require.NotEmpty(t, created.Conversation.ID)

	_, models, err := sub.chatModels(context.Background(), nil, ChatModelsInput{})
	require.NoError(t, err)
	assert.Equal(t, "lemer", models.SelectedModel)

	_, selected, err := sub.chatSelectModel(context.Background(), nil, ChatSelectModelInput{Model: "lemma"})
	require.NoError(t, err)
	assert.Equal(t, "lemma", selected.SelectedModel)

	_, sent, err := sub.chatSend(context.Background(), nil, ChatSendInput{
		ConversationID: created.Conversation.ID,
		Content:        "Summarise the RFC delta.",
	})
	require.NoError(t, err)
	require.Len(t, sent.Conversation.Messages, 2)
	assert.Equal(t, "lemma", sent.Conversation.Model)

	_, history, err := sub.chatHistory(context.Background(), nil, ChatHistoryInput{
		ConversationID: created.Conversation.ID,
	})
	require.NoError(t, err)
	require.Len(t, history.Messages, 2)

	_, exported, err := sub.chatConversationExport(context.Background(), nil, ChatConversationExportInput{
		ID: created.Conversation.ID,
	})
	require.NoError(t, err)
	assert.Contains(t, exported.Markdown, "Summarise the RFC delta.")
}

func TestMCP_Good_ChatConversationSaveAndAttachments(t *testing.T) {
	c, err := core.New(
		core.WithService(display.Register(nil)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)

	_, saved, err := sub.chatConversationSave(context.Background(), nil, ChatConversationSaveInput{
		Conversation: display.Conversation{
			Messages: []display.ChatMessage{
				{Role: "user", Content: "Imported via MCP."},
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, saved.Conversation.ID)
	assert.Equal(t, "Imported via MCP.", saved.Conversation.Title)

	_, attachments, err := sub.chatAttachImage(context.Background(), nil, ChatAttachImageInput{
		ConversationID: saved.Conversation.ID,
		Attachment: display.ImageAttachment{
			Filename: "diagram.png",
			MimeType: "image/png",
			Data:     "ZmFrZQ==",
		},
	})
	require.NoError(t, err)
	require.Len(t, attachments.Attachments, 1)

	_, listed, err := sub.chatAttachmentsGet(context.Background(), nil, ChatAttachmentsGetInput{
		ConversationID: saved.Conversation.ID,
	})
	require.NoError(t, err)
	require.Len(t, listed.Attachments, 1)
	assert.Equal(t, attachments.Attachments[0].ID, listed.Attachments[0].ID)
}

func TestMCP_Good_ChatStreamingThinkingAndToolCalls(t *testing.T) {
	c, err := core.New(
		core.WithService(display.Register(nil)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)

	_, created, err := sub.chatConversationNew(context.Background(), nil, ChatConversationNewInput{})
	require.NoError(t, err)

	_, sent, err := sub.chatSend(context.Background(), nil, ChatSendInput{
		ConversationID: created.Conversation.ID,
		Content:        "Use a local tool while streaming.",
	})
	require.NoError(t, err)
	require.Len(t, sent.Conversation.Messages, 2)

	_, thinkingStart, err := sub.chatThinkingStart(context.Background(), nil, ChatThinkingStartInput{
		ConversationID: created.Conversation.ID,
	})
	require.NoError(t, err)
	assert.True(t, thinkingStart.Thinking.Active)

	_, thinkingAppend, err := sub.chatThinkingAppend(context.Background(), nil, ChatThinkingAppendInput{
		ConversationID: created.Conversation.ID,
		Content:        "Inspecting local tool availability.",
	})
	require.NoError(t, err)
	assert.Contains(t, thinkingAppend.Thinking.Content, "local tool availability")

	_, streamStart, err := sub.chatStreamStart(context.Background(), nil, ChatStreamStartInput{
		ConversationID: created.Conversation.ID,
	})
	require.NoError(t, err)
	require.Len(t, streamStart.Conversation.Messages, 2)
	assert.True(t, streamStart.Conversation.Messages[1].Streaming)

	_, streamAppend, err := sub.chatStreamAppend(context.Background(), nil, ChatStreamAppendInput{
		ConversationID: created.Conversation.ID,
		Content:        "Tool output ready.",
	})
	require.NoError(t, err)
	assert.Equal(t, "Tool output ready.", streamAppend.Conversation.Messages[1].Content)

	_, thinkingEnd, err := sub.chatThinkingEnd(context.Background(), nil, ChatThinkingEndInput{
		ConversationID: created.Conversation.ID,
	})
	require.NoError(t, err)
	assert.False(t, thinkingEnd.Thinking.Active)

	_, streamFinish, err := sub.chatStreamFinish(context.Background(), nil, ChatStreamFinishInput{
		ConversationID: created.Conversation.ID,
		FinishReason:   "stop",
	})
	require.NoError(t, err)
	assert.False(t, streamFinish.Conversation.Messages[1].Streaming)
	assert.Equal(t, "stop", streamFinish.Conversation.Messages[1].FinishReason)

	_, recorded, err := sub.chatRecordToolCall(context.Background(), nil, ChatRecordToolCallInput{
		ConversationID: created.Conversation.ID,
		Call: display.ToolCall{
			ID:   "tool-1",
			Name: "chat_models",
			Arguments: map[string]any{
				"scope": "local",
			},
		},
		Result: display.ToolResult{
			ToolCallID: "tool-1",
			Content:    "lemer, lemma, lemmy",
		},
	})
	require.NoError(t, err)
	require.Len(t, recorded.Conversation.Messages[1].ToolCalls, 1)
	assert.Equal(t, "chat_models", recorded.Conversation.Messages[1].ToolCalls[0].Call.Name)

	_, history, err := sub.chatHistory(context.Background(), nil, ChatHistoryInput{
		ConversationID: created.Conversation.ID,
	})
	require.NoError(t, err)
	require.Len(t, history.Messages, 2)
	assert.Contains(t, history.Messages[1].Thinking.Content, "local tool availability")
	require.Len(t, history.Messages[1].ToolCalls, 1)
	assert.Equal(t, "lemer, lemma, lemmy", history.Messages[1].ToolCalls[0].Result.Content)
}

func TestMCP_Good_MenuRoundTrip(t *testing.T) {
	c, err := core.New(
		core.WithService(menu.Register(menu.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)

	items := []MenuItemSpec{
		{
			Label: "File",
			Children: []MenuItemSpec{
				{Label: "Export", Accelerator: "CmdOrCtrl+E"},
				{Type: "separator"},
				{Label: "Close", Accelerator: "CmdOrCtrl+W"},
			},
		},
		{
			Role: "help",
		},
	}

	_, updated, err := sub.menuSet(context.Background(), nil, MenuSetInput{Items: items})
	require.NoError(t, err)
	require.Len(t, updated.Items, 2)
	assert.Equal(t, "File", updated.Items[0].Label)
	require.Len(t, updated.Items[0].Children, 3)
	assert.Equal(t, "help", updated.Items[1].Role)

	_, fetched, err := sub.menuGet(context.Background(), nil, MenuGetInput{})
	require.NoError(t, err)
	assert.Equal(t, updated.Items, fetched.Items)
}

func TestMCP_Good_ScreenWorkAreaAlias(t *testing.T) {
	c, err := core.New(
		core.WithService(screen.Register(&mockScreenPlatform{
			screens: []screen.Screen{
				{
					ID:        "1",
					Name:      "Primary",
					IsPrimary: true,
					WorkArea:  screen.Rect{X: 0, Y: 24, Width: 1920, Height: 1056},
					Bounds:    screen.Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
					Size:      screen.Size{Width: 1920, Height: 1080},
				},
			},
		})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := New(c)
	_, plural, err := sub.screenWorkAreas(context.Background(), nil, ScreenWorkAreasInput{})
	require.NoError(t, err)
	_, alias, err := sub.screenWorkArea(context.Background(), nil, ScreenWorkAreasInput{})
	require.NoError(t, err)
	assert.Equal(t, plural, alias)
	assert.Len(t, alias.WorkAreas, 1)
	assert.Equal(t, 24, alias.WorkAreas[0].Y)
}

func TestMCP_Good_ScreenForWindow(t *testing.T) {
	c, err := core.New(
		core.WithService(display.Register(nil)),
		core.WithService(screen.Register(&mockScreenPlatform{
			screens: []screen.Screen{
				{
					ID:        "1",
					Name:      "Primary",
					IsPrimary: true,
					WorkArea:  screen.Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
					Bounds:    screen.Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
					Size:      screen.Size{Width: 1920, Height: 1080},
				},
				{
					ID:       "2",
					Name:     "Secondary",
					WorkArea: screen.Rect{X: 1920, Y: 0, Width: 1280, Height: 1024},
					Bounds:   screen.Rect{X: 1920, Y: 0, Width: 1280, Height: 1024},
					Size:     screen.Size{Width: 1280, Height: 1024},
				},
			},
		})),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, err := c.PERFORM(window.TaskOpenWindow{
		Window: &window.Window{Name: "editor", Title: "Editor", X: 100, Y: 100, Width: 800, Height: 600},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	sub := New(c)
	_, out, err := sub.screenForWindow(context.Background(), nil, ScreenForWindowInput{Window: "editor"})
	require.NoError(t, err)
	require.NotNil(t, out.Screen)
	assert.Equal(t, "Primary", out.Screen.Name)
}

func TestMCP_Good_WebviewErrors(t *testing.T) {
	c, err := core.New(
		core.WithService(webview.Register(webview.Options{})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	require.NoError(t, c.ACTION(webview.ActionException{
		Window: "main",
		Exception: webview.ExceptionInfo{
			Text:       "boom",
			URL:        "https://example.com/app.js",
			Line:       12,
			Column:     4,
			StackTrace: "Error: boom",
		},
	}))

	sub := New(c)
	_, out, err := sub.webviewErrors(context.Background(), nil, WebviewErrorsInput{Window: "main"})
	require.NoError(t, err)
	require.Len(t, out.Errors, 1)
	assert.Equal(t, "boom", out.Errors[0].Text)
}

func TestMCP_Bad_NoServices(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	// Without any services, QUERY should return handled=false
	_, handled, _ := c.QUERY(clipboard.QueryText{})
	assert.False(t, handled)
}

type mockEnvPlatform struct {
	isDark bool
}

func (m *mockEnvPlatform) IsDarkMode() bool                                   { return m.isDark }
func (m *mockEnvPlatform) Info() environment.EnvironmentInfo                  { return environment.EnvironmentInfo{} }
func (m *mockEnvPlatform) AccentColour() string                               { return "" }
func (m *mockEnvPlatform) OpenFileManager(path string, selectFile bool) error { return nil }
func (m *mockEnvPlatform) HasFocusFollowsMouse() bool                         { return false }
func (m *mockEnvPlatform) OnThemeChange(handler func(isDark bool)) func() {
	return func() {}
}

type mockScreenPlatform struct {
	screens []screen.Screen
}

func (m *mockScreenPlatform) GetAll() []screen.Screen { return m.screens }
func (m *mockScreenPlatform) GetPrimary() *screen.Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}
func (m *mockScreenPlatform) GetCurrent() *screen.Screen { return m.GetPrimary() }

func TestMCP_Good_ThemeSetRoundTrip(t *testing.T) {
	c, err := core.New(
		core.WithService(environment.Register(&mockEnvPlatform{isDark: true})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	sub := NewSubsystem(c)
	_, output, err := sub.themeSet(context.Background(), nil, ThemeSetInput{Theme: "light"})
	require.NoError(t, err)
	assert.True(t, output.Success)

	result, handled, err := c.QUERY(environment.QueryTheme{})
	require.NoError(t, err)
	assert.True(t, handled)
	theme := result.(environment.ThemeInfo)
	assert.Equal(t, "light", theme.Theme)
	assert.False(t, theme.IsDark)
}

func TestMCP_Good_ScreenFindSpaceAndArrangePair(t *testing.T) {
	c, err := core.New(
		core.WithService(screen.Register(&mockScreenPlatform{screens: []screen.Screen{
			{
				ID: "1", Name: "Primary", IsPrimary: true,
				Bounds:   screen.Rect{X: 0, Y: 0, Width: 1600, Height: 900},
				WorkArea: screen.Rect{X: 0, Y: 0, Width: 1600, Height: 900},
			},
		}})),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, _, err = c.PERFORM(window.TaskOpenWindow{Window: &window.Window{Name: "editor", X: 0, Y: 0, Width: 800, Height: 900}})
	require.NoError(t, err)
	_, _, err = c.PERFORM(window.TaskOpenWindow{Window: &window.Window{Name: "preview", X: 800, Y: 0, Width: 800, Height: 450}})
	require.NoError(t, err)

	sub := NewSubsystem(c)

	_, free, err := sub.screenFindSpace(context.Background(), nil, ScreenFindSpaceInput{Width: 300, Height: 300})
	require.NoError(t, err)
	assert.Equal(t, "1", free.ScreenID)
	assert.Equal(t, screen.Rect{X: 800, Y: 450, Width: 800, Height: 450}, free.Bounds)

	_, arranged, err := sub.windowArrangePair(context.Background(), nil, WindowArrangePairInput{
		First: "editor", Second: "preview",
	})
	require.NoError(t, err)
	assert.Equal(t, screen.Rect{X: 0, Y: 0, Width: 800, Height: 900}, arranged.FirstBounds)
	assert.Equal(t, screen.Rect{X: 800, Y: 0, Width: 800, Height: 900}, arranged.SecondBounds)
}
