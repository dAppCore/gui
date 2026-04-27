package mcp

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToolsLifecycle_appQuit_Good(t *testing.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerLifecycleTools(server)

	result, err := sub.CallTool(context.Background(), "app_quit", map[string]any{})
	require.NoError(t, err)
	assert.Contains(t, result, "\"success\":true")
}

func TestToolsLifecycle_appQuit_Bad(t *testing.T) {
	sub := New(core.New(core.WithServiceLock()))

	_, out, err := sub.appQuit(context.Background(), nil, AppQuitInput{})
	require.NoError(t, err)
	assert.True(t, out.Success)
	assert.Nil(t, err)
}

func TestToolsLifecycle_appQuit_Ugly(t *testing.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerLifecycleTools(server)

	_, err := sub.CallTool(context.Background(), "app_quit", nil)
	require.NoError(t, err)
}
