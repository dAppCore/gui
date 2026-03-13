# GUI Config Wiring — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the in-memory `loadConfig()` stub with real `.core/gui/config.yaml` loading via go-config, and wire `handleConfigTask` to persist changes to disk.

**Architecture:** Display orchestrator owns a `*config.Config` instance pointed at `~/.core/gui/config.yaml`. On startup, it loads file contents into the existing `configData` map. On `TaskSaveConfig`, it calls `cfg.Set()` + `cfg.Commit()` to persist. Sub-services remain unchanged — they already QUERY for their section and receive `map[string]any`.

**Tech Stack:** Go 1.26, `forge.lthn.ai/core/go` v0.2.2 (DI/IPC), `forge.lthn.ai/core/go-config` (Viper-backed YAML config), testify (assert/require)

**Spec:** `docs/superpowers/specs/2026-03-13-gui-config-wiring-design.md`

---

## File Structure

### Modified Files

| File | Changes |
|------|---------|
| `go.mod` | Add `forge.lthn.ai/core/go-config` dependency |
| `pkg/display/display.go` | Add `cfg *config.Config` field, replace `loadConfig()` stub, update `handleConfigTask` to persist via `cfg.Set()` + `cfg.Commit()`, add `guiConfigPath()` helper |
| `pkg/display/display_test.go` | Add config loading + persistence tests |
| `pkg/window/service.go` | Flesh out `applyConfig()` stub: read `default_width`, `default_height`, `state_file` |
| `pkg/systray/service.go` | Flesh out `applyConfig()` stub: read `icon` field (tooltip already works) |
| `pkg/menu/service.go` | Flesh out `applyConfig()` stub: read `show_dev_tools`, expose via accessor |

---

## Task 1: Wire go-config into Display Orchestrator

**Files:**
- Modify: `go.mod`
- Modify: `pkg/display/display.go`
- Modify: `pkg/display/display_test.go`

- [ ] **Step 1: Write failing test — config loads from file**

Add to `pkg/display/display_test.go`:

```go
func TestLoadConfig_Good(t *testing.T) {
	// Create temp config file
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".core", "gui", "config.yaml")
	require.NoError(t, os.MkdirAll(filepath.Dir(cfgPath), 0o755))
	require.NoError(t, os.WriteFile(cfgPath, []byte(`
window:
  default_width: 1280
  default_height: 720
systray:
  tooltip: "Test App"
menu:
  show_dev_tools: false
`), 0o644))

	s, _ := New()
	s.loadConfigFrom(cfgPath)

	// Verify configData was populated from file
	assert.Equal(t, 1280, s.configData["window"]["default_width"])
	assert.Equal(t, "Test App", s.configData["systray"]["tooltip"])
	assert.Equal(t, false, s.configData["menu"]["show_dev_tools"])
}

func TestLoadConfig_Bad_MissingFile(t *testing.T) {
	s, _ := New()
	s.loadConfigFrom(filepath.Join(t.TempDir(), "nonexistent.yaml"))

	// Should not panic, configData stays at empty defaults
	assert.Empty(t, s.configData["window"])
	assert.Empty(t, s.configData["systray"])
	assert.Empty(t, s.configData["menu"])
}
```

- [ ] **Step 2: Add go-config dependency**

```bash
cd /path/to/core/gui
go get forge.lthn.ai/core/go-config
```

The Go workspace will resolve it locally from `~/Code/core/go-config`.

- [ ] **Step 3: Implement `loadConfig()` and `loadConfigFrom()`**

In `pkg/display/display.go`, add the `cfg` field and replace the stub:

```go
import (
	"os"
	"path/filepath"

	"forge.lthn.ai/core/go-config"
)

type Service struct {
	*core.ServiceRuntime[Options]
	wailsApp   *application.App
	app        App
	config     Options
	configData map[string]map[string]any
	cfg        *config.Config  // go-config instance for file persistence
	notifier   *notifications.NotificationService
	events     *WSEventManager
}

func guiConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".core", "gui", "config.yaml")
	}
	return filepath.Join(home, ".core", "gui", "config.yaml")
}

func (s *Service) loadConfig() {
	s.loadConfigFrom(guiConfigPath())
}

func (s *Service) loadConfigFrom(path string) {
	cfg, err := config.New(config.WithPath(path))
	if err != nil {
		// Non-critical — continue with empty configData
		return
	}
	s.cfg = cfg

	for _, section := range []string{"window", "systray", "menu"} {
		var data map[string]any
		if err := cfg.Get(section, &data); err == nil && data != nil {
			s.configData[section] = data
		}
	}
}
```

- [ ] **Step 4: Write failing test — config persists on TaskSaveConfig**

```go
func TestHandleConfigTask_Persists_Good(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	s, _ := New()
	s.loadConfigFrom(cfgPath)  // Creates empty config (file doesn't exist yet)

	// Simulate a TaskSaveConfig through the handler
	c, _ := core.New(
		core.WithService(func(c *core.Core) (any, error) {
			s.ServiceRuntime = core.NewServiceRuntime[Options](c, Options{})
			return s, nil
		}),
		core.WithServiceLock(),
	)
	c.ServiceStartup(context.Background(), nil)

	_, handled, err := c.PERFORM(window.TaskSaveConfig{
		Value: map[string]any{"default_width": 1920},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify file was written
	data, err := os.ReadFile(cfgPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "default_width")
}
```

- [ ] **Step 5: Update `handleConfigTask` to persist**

```go
func (s *Service) handleConfigTask(c *core.Core, t core.Task) (any, bool, error) {
	switch t := t.(type) {
	case window.TaskSaveConfig:
		s.configData["window"] = t.Value
		s.persistSection("window", t.Value)
		return nil, true, nil
	case systray.TaskSaveConfig:
		s.configData["systray"] = t.Value
		s.persistSection("systray", t.Value)
		return nil, true, nil
	case menu.TaskSaveConfig:
		s.configData["menu"] = t.Value
		s.persistSection("menu", t.Value)
		return nil, true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) persistSection(key string, value map[string]any) {
	if s.cfg == nil {
		return
	}
	_ = s.cfg.Set(key, value)
	_ = s.cfg.Commit()
}
```

- [ ] **Step 6: Run tests, verify green**

```bash
core go test --run TestLoadConfig
core go test --run TestHandleConfigTask_Persists
```

---

## Task 2: Flesh Out Sub-Service `applyConfig()` Stubs

**Files:**
- Modify: `pkg/window/service.go`
- Modify: `pkg/systray/service.go`
- Modify: `pkg/menu/service.go`

- [ ] **Step 1: Window — apply default dimensions**

In `pkg/window/service.go`, replace the stub:

```go
func (s *Service) applyConfig(cfg map[string]any) {
	if w, ok := cfg["default_width"]; ok {
		if width, ok := w.(int); ok {
			s.manager.SetDefaultWidth(width)
		}
	}
	if h, ok := cfg["default_height"]; ok {
		if height, ok := h.(int); ok {
			s.manager.SetDefaultHeight(height)
		}
	}
	if sf, ok := cfg["state_file"]; ok {
		if stateFile, ok := sf.(string); ok {
			s.manager.State().SetPath(stateFile)
		}
	}
}
```

> **Note:** The Manager's `SetDefaultWidth`, `SetDefaultHeight`, and `State().SetPath()` methods may need to be added if they don't exist yet. If not present, skip those calls and add a `// TODO:` — this task is about wiring config, not extending Manager's API.

- [ ] **Step 2: Systray — add icon path handling**

In `pkg/systray/service.go`, extend the existing `applyConfig`:

```go
func (s *Service) applyConfig(cfg map[string]any) {
	tooltip, _ := cfg["tooltip"].(string)
	if tooltip == "" {
		tooltip = "Core"
	}
	_ = s.manager.Setup(tooltip, tooltip)

	if iconPath, ok := cfg["icon"].(string); ok && iconPath != "" {
		// Icon loading is deferred to when assets are available.
		// Store the path for later use.
		s.iconPath = iconPath
	}
}
```

Add `iconPath string` field to `Service` struct.

- [ ] **Step 3: Menu — read show_dev_tools flag**

In `pkg/menu/service.go`, replace the stub:

```go
func (s *Service) applyConfig(cfg map[string]any) {
	if v, ok := cfg["show_dev_tools"]; ok {
		if show, ok := v.(bool); ok {
			s.showDevTools = show
		}
	}
}

// ShowDevTools returns whether developer tools menu items should be shown.
func (s *Service) ShowDevTools() bool {
	return s.showDevTools
}
```

Add `showDevTools bool` field to `Service` struct (default `false`).

- [ ] **Step 4: Run full test suite**

```bash
core go test
```

---

## Completion Criteria

1. `loadConfig()` reads from `~/.core/gui/config.yaml` via go-config
2. `handleConfigTask` persists changes to disk via `cfg.Set()` + `cfg.Commit()`
3. Missing/malformed config file does not crash the GUI
4. Sub-service `applyConfig()` methods consume real config values
5. All existing tests continue to pass
6. New tests cover load, persist, and missing-file scenarios

## Deferred Work

- **Manifest/slots integration**: go-scm `Manifest.Layout`/`Manifest.Slots` could feed a `layout` config section for user slot preferences. Not needed yet.
- **Manager API extensions**: `SetDefaultWidth()`, `SetDefaultHeight()`, `State().SetPath()` — add when window Manager is extended.
- **Config file watching**: Viper supports `WatchConfig()` for live reload. Not needed for a desktop app where config changes come through IPC.

## Licence

EUPL-1.2
