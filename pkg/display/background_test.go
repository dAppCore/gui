package display

import (
	"context"

	core "dappco.re/go"
)

func TestBackground_CloneMap_Good(t *core.T) {
	source := map[string]any{"alpha": "one", "beta": 2}

	cloned := cloneMap(source)

	core.AssertNotNil(t, cloned)
	core.AssertEqual(t, source, cloned)

	source["alpha"] = "mutated"
	core.AssertEqual(t, "one", cloned["alpha"])
}

func TestBackground_CloneMap_Bad(t *core.T) {
	cloned := cloneMap(nil)

	core.AssertNotNil(t, cloned)
	core.AssertEmpty(t, cloned)
}

func TestBackground_CloneMap_Ugly(t *core.T) {
	source := map[string]any{"nested": map[string]any{"value": "original"}}

	cloned := cloneMap(source)
	core.AssertNotNil(t, cloned)

	source["nested"].(map[string]any)["value"] = "changed"
	core.AssertEqual(t, map[string]any{"value": "original"}, cloned["nested"])
}

func TestBackground_DecodeMap_Good(t *core.T) {
	source := map[string]any{"scope": "/app"}
	decoded := decodeMap(source)

	core.AssertNotNil(t, decoded)
	core.AssertEqual(t, map[string]any{"scope": "/app"}, decoded)
	source["scope"] = "/mutated"
	if decoded["scope"] != "/app" {
		t.Fatalf("decoded map changed after source mutation: %v", decoded["scope"])
	}
}

func TestBackground_DecodeMap_Bad(t *core.T) {
	decoded := decodeMap("not-a-map")

	core.AssertNotNil(t, decoded)
	core.AssertEmpty(t, decoded)
}

func TestBackground_DecodeMap_Ugly(t *core.T) {
	decoded := decodeMap(nil)

	core.AssertNotNil(t, decoded)
	core.AssertEmpty(t, decoded)
}

func TestBackground_RegisterBackgroundActions_Good(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.serviceWorker.register").Run(context.Background(), core.NewOptions(
		core.Option{Key: "scriptURL", Value: "https://example.com/sw.js"},
		core.Option{Key: "options", Value: map[string]any{"scope": "/app"}},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "/app", payload["scope"])
	core.AssertContains(t, payload, "active")
	core.AssertEqual(t, map[string]any{"scriptURL": "https://example.com/sw.js"}, payload["active"])
}

func TestBackground_RegisterBackgroundActions_Bad(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.background.fetch").Run(context.Background(), core.NewOptions(
		core.Option{Key: "id", Value: "   "},
		core.Option{Key: "requests", Value: nil},
		core.Option{Key: "options", Value: nil},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "", payload["id"])
	core.AssertEqual(t, "registered", payload["state"])
	core.AssertNil(t, payload["requests"])
}

func TestBackground_RegisterBackgroundActions_Ugly(t *core.T) {
	svc, c := newTestDisplayService(t)
	svc.background = NewBackgroundRegistry()

	result := c.Action("core.payment.instrument.set").Run(context.Background(), core.NewOptions(
		core.Option{Key: "key", Value: "  card-01  "},
		core.Option{Key: "details", Value: map[string]any{"network": "visa", "last4": "4242"}},
	))

	core.RequireTrue(t, result.OK)
	payload, ok := result.Value.(map[string]any)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "card-01", payload["key"])
	core.AssertEqual(t, map[string]any{"network": "visa", "last4": "4242"}, payload["details"])
}

func TestBackground_AddSync_Good(t *core.T) {
	r := NewBackgroundRegistry()
	source := map[string]any{"tag": "refresh", "kind": "sync"}
	record := r.AddSync(source)

	core.AssertNotNil(t, record)
	core.AssertEqual(t, "refresh", record["tag"])
	core.AssertEqual(t, "sync", record["kind"])
	source["tag"] = "mutated"
	core.AssertEqual(t, "refresh", record["tag"])
}

func TestBackground_AddSync_Bad(t *core.T) {
	r := NewBackgroundRegistry()
	record := r.AddSync(nil)

	core.AssertNotNil(t, record)
	core.AssertEmpty(t, record)
}

func TestBackground_AddSync_Ugly(t *core.T) {
	r := NewBackgroundRegistry()
	first := r.AddSync(map[string]any{"tag": "sync-1"})
	second := r.AddSync(map[string]any{"tag": "sync-2"})

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
	core.AssertEqual(t, 2, r.SyncRegistrationsCount())
}

func TestBackground_AddPush_Good(t *core.T) {
	r := NewBackgroundRegistry()
	source := map[string]any{"endpoint": "/push/abc", "auth": "core-local"}
	record := r.AddPush(source)

	core.AssertNotNil(t, record)
	core.AssertEqual(t, "/push/abc", record["endpoint"])
	core.AssertEqual(t, "core-local", record["auth"])
	source["endpoint"] = "/push/mutated"
	core.AssertEqual(t, "/push/abc", record["endpoint"])
}

func TestBackground_AddPush_Bad(t *core.T) {
	r := NewBackgroundRegistry()
	record := r.AddPush(nil)

	core.AssertNotNil(t, record)
	core.AssertEmpty(t, record)
}

func TestBackground_AddPush_Ugly(t *core.T) {
	r := NewBackgroundRegistry()
	first := r.AddPush(map[string]any{"endpoint": "/push/abc"})
	second := r.AddPush(map[string]any{"endpoint": "/push/def"})

	core.AssertNotNil(t, first)
	core.AssertNotNil(t, second)
	core.AssertEqual(t, 2, r.PushSubscriptionsCount())
	core.AssertEqual(t, "/push/abc", first["endpoint"])
	core.AssertEqual(t, "/push/def", second["endpoint"])
}
