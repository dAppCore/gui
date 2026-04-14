// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/window"
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
