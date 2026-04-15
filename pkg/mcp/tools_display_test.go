package mcp

import (
	"context"
	"errors"
	"testing"

	core "dappco.re/go/core"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDisplayToolTestSubsystem(t *testing.T, handler func(core.Options) core.Result) *Subsystem {
	t.Helper()
	c := core.New(core.WithServiceLock())
	if handler != nil {
		c.Action("display.resolveScheme", func(_ context.Context, opts core.Options) core.Result {
			return handler(opts)
		})
	}
	return New(c)
}

func TestToolsDisplay_schemeResolve_Good(t *testing.T) {
	sub := newDisplayToolTestSubsystem(t, func(opts core.Options) core.Result {
		return core.Result{
			Value: map[string]any{
				"url":          opts.String("url"),
				"route":        "store",
				"content_type": "text/html",
				"body":         "<html>core://store</html>",
			},
			OK: true,
		}
	})
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerDisplayTools(server)

	result, err := sub.CallTool(context.Background(), "scheme_resolve", map[string]any{"url": "core://store?q=alpha"})
	require.NoError(t, err)
	assert.Contains(t, result, "core://store?q=alpha")
	assert.Contains(t, result, "\"route\":\"store\"")
	assert.Contains(t, result, "\"content_type\":\"text/html\"")
}

func TestToolsDisplay_schemeResolve_Bad(t *testing.T) {
	sub := newDisplayToolTestSubsystem(t, func(core.Options) core.Result {
		return core.Result{Value: errors.New("display offline"), OK: false}
	})

	_, _, err := sub.schemeResolve(context.Background(), nil, SchemeResolveInput{URL: "core://store"})
	require.Error(t, err)
	assert.Equal(t, "display offline", err.Error())
}

func TestToolsDisplay_schemeResolve_Ugly(t *testing.T) {
	sub := newDisplayToolTestSubsystem(t, func(core.Options) core.Result {
		return core.Result{Value: map[string]any{
			"route":        "store",
			"content_type": "text/html",
			"body":         "<html>fallback</html>",
		}, OK: true}
	})

	_, out, err := sub.schemeResolve(context.Background(), nil, SchemeResolveInput{URL: "core://store?q=beta"})
	require.NoError(t, err)
	assert.Equal(t, "core://store?q=beta", out.URL)
	assert.Equal(t, "store", out.Route)
}
