// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/display"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/notification"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/webview"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsystem_Good_Name(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := New(c)
	assert.Equal(t, "display", sub.Name())
}

func TestSubsystem_Good_RegisterTools(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	sub := New(c)
	// RegisterTools should not panic with a real mcp.Server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	assert.NotPanics(t, func() { sub.RegisterTools(server) })
}

// Integration test: verify the IPC round-trip that MCP tool handlers use.

type mockClipPlatform struct {
	text string
	ok   bool
}

func (m *mockClipPlatform) Text() (string, bool)  { return m.text, m.ok }
func (m *mockClipPlatform) SetText(t string) bool { m.text = t; m.ok = t != ""; return true }

type mockNotificationPlatform struct {
	sendCalled bool
	lastOpts   notification.NotificationOptions
}

func (m *mockNotificationPlatform) Send(opts notification.NotificationOptions) error {
	m.sendCalled = true
	m.lastOpts = opts
	return nil
}
func (m *mockNotificationPlatform) RequestPermission() (bool, error) { return true, nil }
func (m *mockNotificationPlatform) CheckPermission() (bool, error)   { return true, nil }

type mockEnvironmentPlatform struct {
	isDark bool
}

func (m *mockEnvironmentPlatform) IsDarkMode() bool { return m.isDark }
func (m *mockEnvironmentPlatform) Info() environment.EnvironmentInfo {
	return environment.EnvironmentInfo{}
}
func (m *mockEnvironmentPlatform) AccentColour() string { return "" }
func (m *mockEnvironmentPlatform) OpenFileManager(path string, selectFile bool) error {
	return nil
}
func (m *mockEnvironmentPlatform) OnThemeChange(handler func(isDark bool)) func() {
	return func() {}
}
func (m *mockEnvironmentPlatform) SetTheme(isDark bool) error {
	m.isDark = isDark
	return nil
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
	if len(m.screens) == 0 {
		return nil
	}
	return &m.screens[0]
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
		core.WithService(webview.Register()),
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
