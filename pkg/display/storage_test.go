package display

import (
	core "dappco.re/go"
	"fmt"
	"strings"
	"time"
)

func storageEntryKey(origin, bucket, key string) string {
	return makeStorageEntryKey(origin, bucket, key)
}

func setStorageEntryTime(r *StorageRegistry, origin, bucket, key string, ts time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.entries[storageEntryKey(origin, bucket, key)]
	entry.UpdatedAt = ts
	r.entries[storageEntryKey(origin, bucket, key)] = entry
}

func TestStorageRegistry_Get_Good(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")
	r.Set("origin-b", "local", "theme", "light")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "local", "theme", time.Unix(200, 0).UTC())

	entry, ok := r.Get("origin-a", "local", "theme")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "dark", entry.Value)
	core.AssertEqual(t, "origin-a", entry.Origin)
}

func TestStorageRegistry_Get_Bad(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")

	entry, ok := r.Get("missing", "local", "theme")
	core.AssertFalse(t, ok)
	core.AssertEmpty(t, entry)
}

func TestStorageRegistry_Get_Ugly(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")
	r.Set("origin-b", "local", "theme", "light")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "local", "theme", time.Unix(200, 0).UTC())

	entry, ok := r.Get("", "local", "")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "origin-b", entry.Origin)
	core.AssertEqual(t, "light", entry.Value)
}

func TestStorageRegistry_Search_Good(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")
	r.Set("origin-b", "session", "token", "bravo")
	r.Set("origin-c", "local", "theme", "alpha-beta")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "session", "token", time.Unix(300, 0).UTC())
	setStorageEntryTime(r, "origin-c", "local", "theme", time.Unix(200, 0).UTC())

	results := r.Search("alpha")
	core.AssertLen(t, results, 2)
	core.AssertEqual(t, "origin-c", results[0].Origin)
	core.AssertEqual(t, "origin-a", results[1].Origin)
}

func TestStorageRegistry_Search_Bad(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")
	r.Set("origin-b", "session", "token", "bravo")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "session", "token", time.Unix(200, 0).UTC())

	results := r.Search("")
	core.AssertLen(t, results, 2)
	core.AssertEqual(t, "origin-b", results[0].Origin)
	core.AssertEqual(t, "origin-a", results[1].Origin)
}

func TestStorageRegistry_Search_Ugly(t *core.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")

	results := r.Search("does-not-exist")
	core.AssertEmpty(t, results)
}

func TestStorageRegistry_Snapshot_Good(t *core.T) {
	r := NewStorageRegistry()
	r.Set("core://settings", "localStorage", "theme", "dark")
	r.Set("core://settings", "cookies", "session", `{"value":"abc","path":"/","secure":false}`)
	r.Set("core://other", "localStorage", "theme", "light")

	snapshot := r.Snapshot("core://settings/profile")
	core.AssertContains(t, snapshot, "localStorage")
	core.AssertContains(t, snapshot, "cookies")
	core.AssertEqual(t, "dark", snapshot["localStorage"]["theme"])
	core.AssertEqual(t, `{"value":"abc","path":"/","secure":false}`, snapshot["cookies"]["session"])
	_, otherOriginPresent := snapshot["other"]
	core.AssertFalse(t, otherOriginPresent)
}

func TestStorageRegistry_Set_Bad(t *core.T) {
	r := NewStorageRegistry()

	core.AssertFalse(t, r.Set("", "localStorage", "theme", "dark"))
	core.AssertFalse(t, r.Set("core://settings", "", "theme", "dark"))
	core.AssertFalse(t, r.Set("core://settings", "localStorage", "", "dark"))
	core.AssertFalse(t, r.Set("core://settings", "localStorage", "theme", strings.Repeat("x", maxStorageValueBytes+1)))
}

func TestStorageRegistry_Delete_Good(t *core.T) {
	r := NewStorageRegistry()
	r.Set("core://settings", "localStorage", "theme", "dark")

	core.AssertTrue(t, r.Delete("core://settings", "localStorage", "theme"))
	_, ok := r.Get("core://settings", "localStorage", "theme")
	core.AssertFalse(t, ok)
}

func TestStorageRegistry_Set_RejectsQuotaOverflow(t *core.T) {
	r := NewStorageRegistry()
	for i := 0; i < maxStorageEntriesPerOrigin; i++ {
		core.RequireTrue(t, r.Set("core://settings", "localStorage", fmt.Sprintf("key-%d", i), "v"))
	}
	core.AssertFalse(t, r.Set("core://settings", "localStorage", "overflow", "v"))
}

func TestStorage_StorageOriginForPageURL_Good(t *core.T) {
	core.AssertEqual(t, "https://app.example.com", storageOriginForPageURL("https://app.example.com/path?q=1"))
	core.AssertEqual(t, "core://settings", storageOriginForPageURL("core://settings/view"))
	core.AssertNotEmpty(t, core.Sprintf("%T", storageOriginForPageURL("https://app.example.com/path?q=1")))
}

func TestStorage_StorageOriginForPageURL_Bad(t *core.T) {
	core.AssertEqual(t, "custom://host/path", storageOriginForPageURL("custom://host/path"))
	observedType := core.Sprintf("%T", storageOriginForPageURL("custom://host/path"))
	core.AssertNotEmpty(t, observedType)
}

func TestStorage_StorageOriginForPageURL_Ugly(t *core.T) {
	core.AssertEqual(t, "", storageOriginForPageURL(""))
	core.AssertEqual(t, "", storageOriginForPageURL("   "))
	core.AssertNotEmpty(t, core.Sprintf("%T", storageOriginForPageURL("")))
}

func TestStorage_Snapshot_BlankOriginReturnsEmpty(t *core.T) {
	r := NewStorageRegistry()
	r.Set("core://settings", "localStorage", "theme", "dark")

	snapshot := r.Snapshot("")

	core.AssertEmpty(t, snapshot)
}

func TestStorage_CompositeKey_Good(t *core.T) {
	key := storageCompositeKey("origin", "bucket", "item")

	origin, bucket, item, ok := decodeStorageCompositeKey(key)
	core.RequireTrue(t, ok)
	core.AssertEqual(t, "origin", origin)
	core.AssertEqual(t, "bucket", bucket)
	core.AssertEqual(t, "item", item)
	core.AssertEqual(t, key, makeStorageEntryKey("origin", "bucket", "item"))
}

func TestStorage_CompositeKey_Bad(t *core.T) {
	origin, bucket, item, ok := decodeStorageCompositeKey("not-json")

	core.AssertFalse(t, ok)
	core.AssertEmpty(t, origin)
	core.AssertEmpty(t, bucket)
	core.AssertEmpty(t, item)
}

func TestStorage_CompositeKey_Ugly(t *core.T) {
	origin, bucket, item, ok := decodeStorageCompositeKey(`["one","two"]`)

	core.AssertFalse(t, ok)
	core.AssertEmpty(t, origin)
	core.AssertEmpty(t, bucket)
	core.AssertEmpty(t, item)
}

func TestStorageRegistry_NilReceiverIsSafe(t *core.T) {
	var r *StorageRegistry

	core.AssertFalse(t, r.Set("core://settings", "localStorage", "theme", "dark"))
	core.AssertFalse(t, r.Delete("core://settings", "localStorage", "theme"))

	entry, ok := r.Get("core://settings", "localStorage", "theme")
	core.AssertFalse(t, ok)
	core.AssertEmpty(t, entry)
	core.AssertEmpty(t, r.Search("theme"))
	core.AssertEmpty(t, r.Snapshot("core://settings"))
}
