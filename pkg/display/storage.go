package display

import (
	"sort"
	"strings"
	"sync"
	"time"
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
}

func NewStorageRegistry() *StorageRegistry {
	return &StorageRegistry{entries: make(map[string]StorageEntry)}
}

func (r *StorageRegistry) Set(origin, bucket, key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	composite := strings.Join([]string{origin, bucket, key}, "\x00")
	r.entries[composite] = StorageEntry{
		Origin:    origin,
		Bucket:    bucket,
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
}

func (r *StorageRegistry) Get(origin, bucket, key string) (StorageEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if entry, ok := r.entries[strings.Join([]string{origin, bucket, key}, "\x00")]; ok {
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
