package mcp

import (
	"context"
	"errors"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/webview"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newWebviewToolsTestSubsystem(t *testing.T, handler func(name string, opts core.Options) core.Result) *Subsystem {
	t.Helper()

	c := core.New(core.WithServiceLock())
	c.Action("webview.devtoolsOpen", func(_ context.Context, opts core.Options) core.Result {
		if handler != nil {
			return handler("webview.devtoolsOpen", opts)
		}
		return core.Result{}
	})
	c.Action("webview.devtoolsClose", func(_ context.Context, opts core.Options) core.Result {
		if handler != nil {
			return handler("webview.devtoolsClose", opts)
		}
		return core.Result{}
	})
	return New(c)
}

func TestToolsWebview_webviewDevTools_Good(t *testing.T) {
	var calls []string

	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		calls = append(calls, name)
		switch task := opts.Get("task").Value.(type) {
		case webview.TaskDevToolsOpen:
			assert.Equal(t, "main", task.Window)
		case webview.TaskDevToolsClose:
			assert.Equal(t, "main", task.Window)
		default:
			t.Fatalf("unexpected task type %T", task)
		}
		return core.Result{OK: true}
	})

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerWebviewTools(server)

	result, err := sub.CallTool(context.Background(), "webview_devtools_open", map[string]any{"window": "main"})
	require.NoError(t, err)
	assert.Contains(t, result, "\"success\":true")

	result, err = sub.CallTool(context.Background(), "webview_devtools_close", map[string]any{"window": "main"})
	require.NoError(t, err)
	assert.Contains(t, result, "\"success\":true")
	assert.Equal(t, []string{"webview.devtoolsOpen", "webview.devtoolsClose"}, calls)
}

func TestToolsWebview_webviewDevToolsOpen_Bad(t *testing.T) {
	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		task, ok := opts.Get("task").Value.(webview.TaskDevToolsOpen)
		require.True(t, ok)
		assert.Equal(t, "main", task.Window)
		assert.Equal(t, "webview.devtoolsOpen", name)
		return core.Result{Value: errors.New("devtools unavailable"), OK: false}
	})

	_, _, err := sub.webviewDevToolsOpen(context.Background(), nil, WebviewDevToolsOpenInput{Window: "main"})
	require.Error(t, err)
	assert.Equal(t, "devtools unavailable", err.Error())
}

func TestToolsWebview_webviewDevToolsClose_Ugly(t *testing.T) {
	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		task, ok := opts.Get("task").Value.(webview.TaskDevToolsClose)
		require.True(t, ok)
		assert.Equal(t, "main", task.Window)
		assert.Equal(t, "webview.devtoolsClose", name)
		return core.Result{Value: "suppressed failure", OK: false}
	})

	_, out, err := sub.webviewDevToolsClose(context.Background(), nil, WebviewDevToolsCloseInput{Window: "main"})
	require.NoError(t, err)
	assert.False(t, out.Success)
}
