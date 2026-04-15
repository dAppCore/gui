// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"forge.lthn.ai/core/gui/pkg/events"
	"forge.lthn.ai/core/gui/pkg/screen"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsystem_Good_Name(t *testing.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	assert.Equal(t, "display", sub.Name())
}

func TestSubsystem_Good_RegisterTools(t *testing.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	// RegisterTools should not panic with a real mcp.Server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	assert.NotPanics(t, func() { sub.RegisterTools(server) })
	assert.NotEmpty(t, sub.Manifest())
	assert.Contains(t, sub.ManifestText(), "layout_suggest")
	assert.Contains(t, sub.ManifestText(), "window_title_set")
	assert.Contains(t, sub.ManifestText(), "focus_set")
	assert.Contains(t, sub.ManifestText(), "dialog_message")
	assert.Contains(t, sub.ManifestText(), "event_info")
	assert.Contains(t, sub.ManifestText(), "screen_work_area")
}

// Integration test: verify the IPC round-trip that MCP tool handlers use.

type mockClipPlatform struct {
	text string
	ok   bool
}

func (m *mockClipPlatform) Text() (string, bool)  { return m.text, m.ok }
func (m *mockClipPlatform) SetText(t string) bool { m.text = t; m.ok = t != ""; return true }

func TestMCP_Good_ClipboardRoundTrip(t *testing.T) {
	c := core.New(
		core.WithService(clipboard.Register(&mockClipPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	// Verify the IPC path that clipboard_read tool handler uses
	r := c.QUERY(clipboard.QueryText{})
	require.True(t, r.OK)
	content, ok := r.Value.(clipboard.ClipboardContent)
	require.True(t, ok, "expected ClipboardContent type")
	assert.Equal(t, "hello", content.Text)
}

func TestMCP_Bad_NoServices(t *testing.T) {
	c := core.New(core.WithServiceLock())
	// Without any services, QUERY should return OK=false
	r := c.QUERY(clipboard.QueryText{})
	assert.False(t, r.OK)
}

type manifestScreenPlatform struct{}

func (manifestScreenPlatform) GetAll() []screen.Screen {
	return []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}}
}

func (p manifestScreenPlatform) GetPrimary() *screen.Screen {
	all := p.GetAll()
	return &all[0]
}

func (p manifestScreenPlatform) GetCurrent() *screen.Screen {
	return p.GetPrimary()
}

func TestSubsystem_Good_CallTool_LayoutSuggest(t *testing.T) {
	c := core.New(
		core.WithService(screen.Register(manifestScreenPlatform{})),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	result, err := sub.CallTool(context.Background(), "layout_suggest", map[string]any{"window_count": 2})
	require.NoError(t, err)
	assert.Contains(t, result, "left-right")
}

type aliasDialogPlatform struct {
	last dialog.MessageDialogOptions
}

func (m *aliasDialogPlatform) OpenFile(_ dialog.OpenFileOptions) ([]string, error) { return nil, nil }
func (m *aliasDialogPlatform) SaveFile(_ dialog.SaveFileOptions) (string, error)   { return "", nil }
func (m *aliasDialogPlatform) OpenDirectory(_ dialog.OpenDirectoryOptions) (string, error) {
	return "", nil
}
func (m *aliasDialogPlatform) MessageDialog(opts dialog.MessageDialogOptions) (string, error) {
	m.last = opts
	return "OK", nil
}

type aliasEventsPlatform struct {
	mu        sync.Mutex
	listeners map[string]int
}

func (m *aliasEventsPlatform) Emit(_ string, _ ...any) bool { return false }
func (m *aliasEventsPlatform) On(name string, _ func(*events.CustomEvent)) func() {
	m.mu.Lock()
	if m.listeners == nil {
		m.listeners = make(map[string]int)
	}
	m.listeners[name]++
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.listeners[name] > 0 {
			m.listeners[name]--
		}
	}
}
func (m *aliasEventsPlatform) Off(name string) {
	m.mu.Lock()
	delete(m.listeners, name)
	m.mu.Unlock()
}
func (m *aliasEventsPlatform) OnMultiple(name string, callback func(*events.CustomEvent), counter int) {
	_ = m.On(name, callback)
}
func (m *aliasEventsPlatform) Reset() {
	m.mu.Lock()
	m.listeners = make(map[string]int)
	m.mu.Unlock()
}

func TestSubsystem_Good_CallTool_RFCAliases(t *testing.T) {
	windowPlatform := window.NewMockPlatform()
	dialogPlatform := &aliasDialogPlatform{}
	eventsPlatform := &aliasEventsPlatform{}

	c := core.New(
		core.WithService(screen.Register(manifestScreenPlatform{})),
		core.WithService(window.Register(windowPlatform)),
		core.WithService(dialog.Register(dialogPlatform)),
		core.WithService(events.Register(eventsPlatform)),
		core.WithServiceLock(),
	)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case events.QueryServerInfo:
			return core.Result{Value: events.ServerInfo{
				ConnectedClients:  2,
				SubscriptionCount: 5,
				BufferLength:      1,
				BufferCapacity:    100,
			}, OK: true}
		default:
			return core.Result{}
		}
	})
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "window_create", map[string]any{
		"name": "main", "title": "Original", "x": 10, "y": 20, "width": 300, "height": 200,
	})
	require.NoError(t, err)

	_, err = sub.CallTool(context.Background(), "window_title_set", map[string]any{
		"name": "main", "title": "Updated",
	})
	require.NoError(t, err)

	titleResult, err := sub.CallTool(context.Background(), "window_title_get", map[string]any{"name": "main"})
	require.NoError(t, err)
	assert.Contains(t, titleResult, "Updated")

	_, err = sub.CallTool(context.Background(), "focus_set", map[string]any{"name": "main"})
	require.NoError(t, err)

	focusedResult, err := sub.CallTool(context.Background(), "window_focused", map[string]any{})
	require.NoError(t, err)
	assert.Contains(t, focusedResult, "main")

	_, err = sub.CallTool(context.Background(), "window_bounds", map[string]any{
		"name": "main", "x": 25, "y": 35, "width": 640, "height": 480,
	})
	require.NoError(t, err)

	windowResult, err := sub.CallTool(context.Background(), "window_get", map[string]any{"name": "main"})
	require.NoError(t, err)
	assert.Contains(t, windowResult, "\"width\":640")
	assert.Contains(t, windowResult, "\"height\":480")

	screenResult, err := sub.CallTool(context.Background(), "screen_work_area", map[string]any{})
	require.NoError(t, err)
	assert.Contains(t, screenResult, "\"width\":2000")

	dialogResult, err := sub.CallTool(context.Background(), "dialog_message", map[string]any{
		"type": "warning", "title": "Heads up", "message": "Check this",
	})
	require.NoError(t, err)
	assert.Equal(t, dialog.DialogWarning, dialogPlatform.last.Type)
	assert.Contains(t, dialogResult, "OK")

	subscribeResult, err := sub.CallTool(context.Background(), "event_subscribe", map[string]any{"name": "theme:changed"})
	require.NoError(t, err)
	assert.Contains(t, subscribeResult, "success")

	eventInfoResult, err := sub.CallTool(context.Background(), "event_info", map[string]any{})
	require.NoError(t, err)
	assert.Contains(t, eventInfoResult, "\"connectedClients\":2")

	_, err = sub.CallTool(context.Background(), "event_unsubscribe", map[string]any{"name": "theme:changed"})
	require.NoError(t, err)
}
