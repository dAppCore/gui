package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackground_CloneMap_Good(t *testing.T) {
	source := map[string]any{"alpha": "one", "beta": 2}

	cloned := cloneMap(source)

	require.NotNil(t, cloned)
	assert.Equal(t, source, cloned)

	source["alpha"] = "mutated"
	assert.Equal(t, "one", cloned["alpha"])
}

func TestBackground_CloneMap_Bad(t *testing.T) {
	cloned := cloneMap(nil)

	require.NotNil(t, cloned)
	assert.Empty(t, cloned)
}

func TestBackground_CloneMap_Ugly(t *testing.T) {
	source := map[string]any{"nested": map[string]any{"value": "original"}}

	cloned := cloneMap(source)
	require.NotNil(t, cloned)

	source["nested"] = map[string]any{"value": "changed"}
	assert.Equal(t, map[string]any{"value": "original"}, cloned["nested"])
}

func TestBackground_DecodeMap_Good(t *testing.T) {
	decoded := decodeMap(map[string]any{"scope": "/app"})

	require.NotNil(t, decoded)
	assert.Equal(t, map[string]any{"scope": "/app"}, decoded)
}

func TestBackground_DecodeMap_Bad(t *testing.T) {
	decoded := decodeMap("not-a-map")

	require.NotNil(t, decoded)
	assert.Empty(t, decoded)
}

func TestBackground_DecodeMap_Ugly(t *testing.T) {
	decoded := decodeMap(nil)

	require.NotNil(t, decoded)
	assert.Empty(t, decoded)
}

func TestBackground_RegisterBackgroundActions_Good(t *testing.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.serviceWorker.register").Run(context.Background(), core.NewOptions(
		core.Option{Key: "scriptURL", Value: "https://example.com/sw.js"},
		core.Option{Key: "options", Value: map[string]any{"scope": "/app"}},
	))

	require.True(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "/app", payload["scope"])
	require.Contains(t, payload, "active")
	assert.Equal(t, map[string]any{"scriptURL": "https://example.com/sw.js"}, payload["active"])
}

func TestBackground_RegisterBackgroundActions_Bad(t *testing.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: "   "},
		core.Option{Key: "requests", Value: nil},
		core.Option{Key: "options", Value: nil},
	))

	require.True(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "", payload["id"])
	assert.Equal(t, "registered", payload["state"])
	assert.Nil(t, payload["requests"])
}

func TestBackground_RegisterBackgroundActions_Ugly(t *testing.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.payment.instrument.set").Run(context.Background(), core.NewOptions(
		core.Option{Key: "key", Value: "  card-01  "},
		core.Option{Key: "details", Value: map[string]any{"network": "visa", "last4": "4242"}},
	))

	require.True(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "card-01", payload["key"])
	assert.Equal(t, map[string]any{"network": "visa", "last4": "4242"}, payload["details"])
}

func TestBackground_AddSync_Good(t *testing.T) {
	r := NewBackgroundRegistry()
	source := map[string]any{"tag": "refresh", "kind": "sync"}
	record := r.AddSync(source)

	require.NotNil(t, record)
	assert.Equal(t, "refresh", record["tag"])
	assert.Equal(t, "sync", record["kind"])
	source["tag"] = "mutated"
	assert.Equal(t, "refresh", record["tag"])
}

func TestBackground_AddSync_Bad(t *testing.T) {
	r := NewBackgroundRegistry()
	record := r.AddSync(nil)

	require.NotNil(t, record)
	assert.Empty(t, record)
}

func TestBackground_AddSync_Ugly(t *testing.T) {
	r := NewBackgroundRegistry()
	first := r.AddSync(map[string]any{"tag": "sync-1"})
	second := r.AddSync(map[string]any{"tag": "sync-2"})

	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.Equal(t, 2, r.SyncRegistrationsCount())
}

func TestBackground_AddPush_Good(t *testing.T) {
	r := NewBackgroundRegistry()
	source := map[string]any{"endpoint": "/push/abc", "auth": "core-local"}
	record := r.AddPush(source)

	require.NotNil(t, record)
	assert.Equal(t, "/push/abc", record["endpoint"])
	assert.Equal(t, "core-local", record["auth"])
	source["endpoint"] = "/push/mutated"
	assert.Equal(t, "/push/abc", record["endpoint"])
}

func TestBackground_AddPush_Bad(t *testing.T) {
	r := NewBackgroundRegistry()
	record := r.AddPush(nil)

	require.NotNil(t, record)
	assert.Empty(t, record)
}

func TestBackground_AddPush_Ugly(t *testing.T) {
	r := NewBackgroundRegistry()
	first := r.AddPush(map[string]any{"endpoint": "/push/abc"})
	second := r.AddPush(map[string]any{"endpoint": "/push/def"})

	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.Equal(t, 2, r.PushSubscriptionsCount())
	assert.Equal(t, "/push/abc", first["endpoint"])
	assert.Equal(t, "/push/def", second["endpoint"])
}
