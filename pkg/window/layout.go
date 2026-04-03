// pkg/window/layout.go
package window

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Layout is a named window arrangement.
// Use: layout := window.Layout{Name: "coding"}
type Layout struct {
	Name      string                 `json:"name"`
	Windows   map[string]WindowState `json:"windows"`
	CreatedAt int64                  `json:"createdAt"`
	UpdatedAt int64                  `json:"updatedAt"`
}

// LayoutInfo is a summary of a layout.
// Use: info := window.LayoutInfo{Name: "coding", WindowCount: 2}
type LayoutInfo struct {
	Name        string `json:"name"`
	WindowCount int    `json:"windowCount"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// LayoutManager persists named window arrangements to ~/.config/Core/layouts.json.
// Use: lm := window.NewLayoutManager()
type LayoutManager struct {
	configDir string
	layouts   map[string]Layout
	mu        sync.RWMutex
}

// NewLayoutManager creates a LayoutManager loading from the default config directory.
// Use: lm := window.NewLayoutManager()
func NewLayoutManager() *LayoutManager {
	lm := &LayoutManager{
		layouts: make(map[string]Layout),
	}
	configDir, err := os.UserConfigDir()
	if err == nil {
		lm.configDir = filepath.Join(configDir, "Core")
	}
	lm.loadLayouts()
	return lm
}

// NewLayoutManagerWithDir creates a LayoutManager loading from a custom config directory.
// Useful for testing or when the default config directory is not appropriate.
// Use: lm := window.NewLayoutManagerWithDir(t.TempDir())
func NewLayoutManagerWithDir(configDir string) *LayoutManager {
	lm := &LayoutManager{
		configDir: configDir,
		layouts:   make(map[string]Layout),
	}
	lm.loadLayouts()
	return lm
}

func (lm *LayoutManager) layoutsFilePath() string {
	return filepath.Join(lm.configDir, "layouts.json")
}

func (lm *LayoutManager) loadLayouts() {
	if lm.configDir == "" {
		return
	}
	data, err := os.ReadFile(lm.layoutsFilePath())
	if err != nil {
		return
	}
	lm.mu.Lock()
	defer lm.mu.Unlock()
	_ = json.Unmarshal(data, &lm.layouts)
}

func (lm *LayoutManager) saveLayouts() {
	if lm.configDir == "" {
		return
	}
	lm.mu.RLock()
	data, err := json.MarshalIndent(lm.layouts, "", "  ")
	lm.mu.RUnlock()
	if err != nil {
		return
	}
	_ = os.MkdirAll(lm.configDir, 0o755)
	_ = os.WriteFile(lm.layoutsFilePath(), data, 0o644)
}

// SaveLayout creates or updates a named layout.
// Use: _ = lm.SaveLayout("coding", windowStates)
func (lm *LayoutManager) SaveLayout(name string, windowStates map[string]WindowState) error {
	if name == "" {
		return fmt.Errorf("layout name cannot be empty")
	}
	now := time.Now().UnixMilli()
	lm.mu.Lock()
	existing, exists := lm.layouts[name]
	layout := Layout{
		Name:      name,
		Windows:   windowStates,
		UpdatedAt: now,
	}
	if exists {
		layout.CreatedAt = existing.CreatedAt
	} else {
		layout.CreatedAt = now
	}
	lm.layouts[name] = layout
	lm.mu.Unlock()
	lm.saveLayouts()
	return nil
}

// GetLayout returns a layout by name.
// Use: layout, ok := lm.GetLayout("coding")
func (lm *LayoutManager) GetLayout(name string) (Layout, bool) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	l, ok := lm.layouts[name]
	return l, ok
}

// ListLayouts returns info summaries for all layouts.
// Use: layouts := lm.ListLayouts()
func (lm *LayoutManager) ListLayouts() []LayoutInfo {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	infos := make([]LayoutInfo, 0, len(lm.layouts))
	for _, l := range lm.layouts {
		infos = append(infos, LayoutInfo{
			Name: l.Name, WindowCount: len(l.Windows),
			CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
		})
	}
	return infos
}

// DeleteLayout removes a layout by name.
// Use: lm.DeleteLayout("coding")
func (lm *LayoutManager) DeleteLayout(name string) {
	lm.mu.Lock()
	delete(lm.layouts, name)
	lm.mu.Unlock()
	lm.saveLayouts()
}
