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
}

func newTestCoreSchemeHandler(t *testing.T) (RouteSchemeHandler, *schemeDispatchRecorder) {
	t.Helper()

	c := core.New(
		core.WithService(Register(nil)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	recorder := &schemeDispatchRecorder{}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		name, ok := q.(string)
		if !ok {
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

	c.Action("core.agent", func(_ context.Context, _ core.Options) core.Result {
		recorder.actions = append(recorder.actions, "core.agent")
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
