package display

import (
	"context"
	"net/url"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type schemeDispatchRecorder struct {
	queries []string
	actions []string
	params  map[string]url.Values
}

func newTestCoreSchemeHandler(t *testing.T) (RouteSchemeHandler, *schemeDispatchRecorder) {
	t.Helper()

	c := newTestCore(t)

	recorder := &schemeDispatchRecorder{params: make(map[string]url.Values)}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		var name string
		switch typed := q.(type) {
		case CoreRouteQuery:
			name = typed.Target
			recorder.params[name] = cloneURLValues(typed.Params)
		case string:
			name = typed
		default:
			return core.Result{}
		}

		recorder.queries = append(recorder.queries, name)
		switch name {
		case "core.settings":
			return core.Result{Value: "settings-query", OK: true}
		case "core.store":
			return core.Result{Value: "store-query", OK: true}
		case "core.network":
			return core.Result{Value: "network-query", OK: true}
		case "core.models":
			return core.Result{Value: "models-query", OK: true}
		default:
			return core.Result{}
		}
	})

	c.Action("core.agent", func(_ context.Context, opts core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.agent")
		recorder.params["core.agent"] = url.Values{"q": []string{opts.String("q")}}
		return core.Result{Value: "agent-action", OK: true}
	})
	c.Action("core.wallet", func(_ context.Context, _ core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.wallet")
		return core.Result{Value: "wallet-action", OK: true}
	})
	c.Action("core.identity", func(_ context.Context, _ core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.identity")
		return core.Result{Value: "identity-action", OK: true}
	})

	svc := core.MustServiceFor[*Service](c, "display")
	return svc.SchemeHandler(), recorder
}

func TestSchemeHandler_Handle_ForwardsQueryParameters(t *testing.T) {
	handler, recorder := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://store?q=invoice&tag=a&tag=b")
	require.NoError(t, err)

	result := handler.Handle(parsedURL)
	require.True(t, result.OK)
	assert.Equal(t, "store-query", result.Value)
	assert.Equal(t, []string{"invoice"}, recorder.params["core.store"]["q"])
	assert.Equal(t, []string{"a", "b"}, recorder.params["core.store"]["tag"])

	actionURL, err := url.Parse("core://agent?q=launch")
	require.NoError(t, err)

	result = handler.Handle(actionURL)
	require.True(t, result.OK)
	assert.Equal(t, "agent-action", result.Value)
	assert.Equal(t, []string{"launch"}, recorder.params["core.agent"]["q"])
}

func TestSchemeHandler_Handle_Good(t *testing.T) {
	handler, recorder := newTestCoreSchemeHandler(t)

	tests := []struct {
		rawURL string
		value  string
	}{
		{rawURL: "core://settings", value: "settings-query"},
		{rawURL: "core://store", value: "store-query"},
		{rawURL: "core://network", value: "network-query"},
		{rawURL: "core://models", value: "models-query"},
		{rawURL: "core://agent", value: "agent-action"},
		{rawURL: "core://wallet", value: "wallet-action"},
		{rawURL: "core://identity", value: "identity-action"},
	}

	for _, test := range tests {
		parsedURL, err := url.Parse(test.rawURL)
		require.NoError(t, err)

		result := handler.Handle(parsedURL)
		require.True(t, result.OK, test.rawURL)
		assert.Equal(t, test.value, result.Value)
	}

	assert.Equal(t, []string{
		"core.settings",
		"core.store",
		"core.network",
		"core.models",
	}, recorder.queries)
	assert.Equal(t, []string{
		"core.agent",
		"core.wallet",
		"core.identity",
	}, recorder.actions)
}

func TestSchemeHandler_Handle_Bad(t *testing.T) {
	handler, _ := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://missing")
	require.NoError(t, err)

	result := handler.Handle(parsedURL)
	require.False(t, result.OK)
	assert.ErrorContains(t, result.Value.(error), "unknown core route: missing")
}

func TestSchemeHandler_Handle_Ugly(t *testing.T) {
	handler, _ := newTestCoreSchemeHandler(t)

	parsedURL, err := url.Parse("core://settings/profile")
	require.NoError(t, err)

	result := handler.Handle(parsedURL)
	require.False(t, result.OK)
	assert.ErrorContains(t, result.Value.(error), "malformed core URL")
}
