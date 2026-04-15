package display

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	core "dappco.re/go/core"
	gostore "dappco.re/go/store"
)

type StorageEntry struct {
	Origin    string    `json:"origin"`
	Bucket    string    `json:"bucket"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StorageRegistry struct {
	mu      sync.RWMutex
	entries map[string]StorageEntry
	store   *gostore.Store
}

func NewStorageRegistry() *StorageRegistry {
	registry := &StorageRegistry{entries: make(map[string]StorageEntry)}
	registry.store = openStorageStore()
	registry.loadPersistedEntries()
	return registry
}

func openStorageStore() *gostore.Store {
	path := storageDatabasePath()
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if path != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil
		}
	}
	storeInstance, err := gostore.New(path)
	if err != nil {
		return nil
	}
	return storeInstance
}

func storageDatabasePath() string {
	if override := strings.TrimSpace(core.Env("CORE_GUI_STORAGE_PATH")); override != "" {
		return override
	}
	if strings.TrimSpace(core.Env("CORE_GUI_STORAGE_PERSIST")) == "" {
		return ":memory:"
	}
	home := strings.TrimSpace(core.Env("DIR_HOME"))
	if home == "" {
		return ":memory:"
	}
	return core.Path(home, ".core", "state", fmt.Sprintf("gui-storage-%d.db", os.Getpid()))
}

func makeStorageEntryKey(origin, bucket, key string) string {
	return strings.Join([]string{origin, bucket, key}, "\x00")
}

func storageCompositeKey(origin, bucket, key string) string {
	return core.JSONMarshalString([]string{origin, bucket, key})
}

func decodeStorageCompositeKey(value string) (string, string, string, bool) {
	var parts []string
	if result := core.JSONUnmarshalString(value, &parts); !result.OK || len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func (r *StorageRegistry) loadPersistedEntries() {
	if r == nil || r.store == nil {
		return
	}
	for entry, err := range r.store.AllSeq("storage") {
		if err != nil {
			continue
		}
		var stored StorageEntry
		if result := core.JSONUnmarshalString(entry.Value, &stored); !result.OK {
			if origin, bucket, key, ok := decodeStorageCompositeKey(entry.Key); ok {
				stored = StorageEntry{
					Origin:    origin,
					Bucket:    bucket,
					Key:       key,
					Value:     entry.Value,
					UpdatedAt: time.Now(),
				}
			} else {
				continue
			}
		}
		if stored.UpdatedAt.IsZero() {
			stored.UpdatedAt = time.Now()
		}
		r.entries[makeStorageEntryKey(stored.Origin, stored.Bucket, stored.Key)] = stored
	}
}

func (r *StorageRegistry) Set(origin, bucket, key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry := StorageEntry{
		Origin:    origin,
		Bucket:    bucket,
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	composite := makeStorageEntryKey(origin, bucket, key)
	r.entries[composite] = entry
	if r.store != nil {
		_ = r.store.Set("storage", storageCompositeKey(origin, bucket, key), core.JSONMarshalString(entry))
	}
}

func (r *StorageRegistry) Get(origin, bucket, key string) (StorageEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if entry, ok := r.entries[makeStorageEntryKey(origin, bucket, key)]; ok {
		return entry, true
	}

	var latest StorageEntry
	found := false
	for _, entry := range r.entries {
		if bucket != "" && entry.Bucket != bucket {
			continue
		}
		if key != "" && entry.Key != key {
			continue
		}
		if origin != "" && entry.Origin != origin {
			continue
		}
		if !found || entry.UpdatedAt.After(latest.UpdatedAt) {
			latest = entry
			found = true
		}
	}
	return latest, found
}

func (r *StorageRegistry) Search(query string) []StorageEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	needle := strings.ToLower(strings.TrimSpace(query))
	results := make([]StorageEntry, 0)
	for _, entry := range r.entries {
		if needle == "" ||
			strings.Contains(strings.ToLower(entry.Origin), needle) ||
			strings.Contains(strings.ToLower(entry.Bucket), needle) ||
			strings.Contains(strings.ToLower(entry.Key), needle) ||
			strings.Contains(strings.ToLower(entry.Value), needle) {
			results = append(results, entry)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	return results
}

func (r *StorageRegistry) Close() error {
	if r == nil || r.store == nil {
		return nil
	}
	return r.store.Close()
}
