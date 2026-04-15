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
