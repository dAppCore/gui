// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
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
