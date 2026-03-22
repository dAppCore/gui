// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
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
func (m *mockClipPlatform) SetText(t string) bool  { m.text = t; m.ok = t != ""; return true }

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
