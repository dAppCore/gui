package display

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func storageEntryKey(origin, bucket, key string) string {
	return strings.Join([]string{origin, bucket, key}, "\x00")
}

func setStorageEntryTime(r *StorageRegistry, origin, bucket, key string, ts time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := r.entries[storageEntryKey(origin, bucket, key)]
	entry.UpdatedAt = ts
	r.entries[storageEntryKey(origin, bucket, key)] = entry
}

func TestStorageRegistry_Get_Good(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")
	r.Set("origin-b", "local", "theme", "light")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "local", "theme", time.Unix(200, 0).UTC())

	entry, ok := r.Get("origin-a", "local", "theme")
	require.True(t, ok)
	assert.Equal(t, "dark", entry.Value)
	assert.Equal(t, "origin-a", entry.Origin)
}

func TestStorageRegistry_Get_Bad(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")

	entry, ok := r.Get("missing", "local", "theme")
	assert.False(t, ok)
	assert.Zero(t, entry)
}

func TestStorageRegistry_Get_Ugly(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "dark")
	r.Set("origin-b", "local", "theme", "light")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "local", "theme", time.Unix(200, 0).UTC())

	entry, ok := r.Get("", "local", "")
	require.True(t, ok)
	assert.Equal(t, "origin-b", entry.Origin)
	assert.Equal(t, "light", entry.Value)
}

func TestStorageRegistry_Search_Good(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")
	r.Set("origin-b", "session", "token", "bravo")
	r.Set("origin-c", "local", "theme", "alpha-beta")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "session", "token", time.Unix(300, 0).UTC())
	setStorageEntryTime(r, "origin-c", "local", "theme", time.Unix(200, 0).UTC())

	results := r.Search("alpha")
	require.Len(t, results, 2)
	assert.Equal(t, "origin-c", results[0].Origin)
	assert.Equal(t, "origin-a", results[1].Origin)
}

func TestStorageRegistry_Search_Bad(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")
	r.Set("origin-b", "session", "token", "bravo")
	setStorageEntryTime(r, "origin-a", "local", "theme", time.Unix(100, 0).UTC())
	setStorageEntryTime(r, "origin-b", "session", "token", time.Unix(200, 0).UTC())

	results := r.Search("")
	require.Len(t, results, 2)
	assert.Equal(t, "origin-b", results[0].Origin)
	assert.Equal(t, "origin-a", results[1].Origin)
}

func TestStorageRegistry_Search_Ugly(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("origin-a", "local", "theme", "alpha")

	results := r.Search("does-not-exist")
	require.Empty(t, results)
}

func TestStorageRegistry_Snapshot_Good(t *testing.T) {
	r := NewStorageRegistry()
	r.Set("core://settings", "localStorage", "theme", "dark")
	r.Set("core://settings", "cookies", "session", `{"value":"abc","path":"/","secure":false}`)
	r.Set("core://other", "localStorage", "theme", "light")

	snapshot := r.Snapshot("core://settings/profile")
	require.Contains(t, snapshot, "localStorage")
	require.Contains(t, snapshot, "cookies")
	assert.Equal(t, "dark", snapshot["localStorage"]["theme"])
	assert.Equal(t, `{"value":"abc","path":"/","secure":false}`, snapshot["cookies"]["session"])
	_, otherOriginPresent := snapshot["other"]
	assert.False(t, otherOriginPresent)
}

func TestStorage_StorageOriginForPageURL_Good(t *testing.T) {
	assert.Equal(t, "https://app.example.com", storageOriginForPageURL("https://app.example.com/path?q=1"))
	assert.Equal(t, "core://settings", storageOriginForPageURL("core://settings/view"))
}

func TestStorage_StorageOriginForPageURL_Bad(t *testing.T) {
	assert.Equal(t, "custom://host/path", storageOriginForPageURL("custom://host/path"))
}

func TestStorage_StorageOriginForPageURL_Ugly(t *testing.T) {
	assert.Equal(t, "", storageOriginForPageURL(""))
	assert.Equal(t, "", storageOriginForPageURL("   "))
}

func TestStorage_CompositeKey_Good(t *testing.T) {
	key := storageCompositeKey("origin", "bucket", "item")

	origin, bucket, item, ok := decodeStorageCompositeKey(key)
	require.True(t, ok)
	assert.Equal(t, "origin", origin)
	assert.Equal(t, "bucket", bucket)
	assert.Equal(t, "item", item)
	assert.Equal(t, "origin\x00bucket\x00item", makeStorageEntryKey("origin", "bucket", "item"))
}

func TestStorage_CompositeKey_Bad(t *testing.T) {
	origin, bucket, item, ok := decodeStorageCompositeKey("not-json")

	assert.False(t, ok)
	assert.Empty(t, origin)
	assert.Empty(t, bucket)
	assert.Empty(t, item)
}

func TestStorage_CompositeKey_Ugly(t *testing.T) {
	origin, bucket, item, ok := decodeStorageCompositeKey(`["one","two"]`)

	assert.False(t, ok)
	assert.Empty(t, origin)
	assert.Empty(t, bucket)
	assert.Empty(t, item)
}
