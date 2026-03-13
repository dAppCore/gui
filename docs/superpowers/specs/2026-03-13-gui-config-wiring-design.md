# GUI Config Wiring — Design Spec

**Date:** 2026-03-13
**Status:** Approved
**Scope:** Replace the in-memory `loadConfig()` stub in `pkg/display` with real `go-config` file-backed configuration, and wire `handleConfigTask` to persist changes to disk.

## Context

The Service Conclave IPC integration (previous spec) established the pattern: display owns config, sub-services QUERY for their section during `OnStartup`, and saves route through `TaskSaveConfig` back to display. All the IPC plumbing works today — but `loadConfig()` creates empty maps and `handleConfigTask` only updates memory. Nothing reads from or writes to disk.

`go-config` (`forge.lthn.ai/core/go-config`) provides exactly what we need: `config.New(config.WithPath(...))` loads a YAML file via Viper, `Get(key, &out)` reads sections by dot-notation key, `Set(key, v)` updates in memory, and `Commit()` persists to disk. It handles directory creation, file-not-exists (returns empty config), and concurrent access via `sync.RWMutex`.

## Design

### Config File Location

`.core/gui/config.yaml` — resolved relative to `$HOME`:

```
~/.core/gui/config.yaml
```

This follows the go-config convention (`~/.core/config.yaml`) but scoped to the GUI module.

### Config Format

```yaml
window:
  state_file: window_state.json
  default_width: 1024
  default_height: 768
systray:
  icon: apptray.png
  tooltip: "Core GUI"
menu:
  show_dev_tools: true
```

Top-level keys map 1:1 to sub-service names. Each sub-service receives its section as `map[string]any` via the existing `QueryConfig` / `handleConfigQuery` pattern.

### Changes to `pkg/display/display.go`

1. **New field**: Add `cfg *config.Config` to `Service` struct (replaces reliance on `configData` maps for initial load).

2. **`loadConfig()`**: Replace stub with:
   ```go
   func (s *Service) loadConfig() {
       cfg, err := config.New(config.WithPath(guiConfigPath()))
       if err != nil {
           // Log warning, continue with empty config — GUI should not crash
           // if config file is missing or malformed.
           s.cfg = nil
           return
       }
       s.cfg = cfg

       // Populate configData sections from file
       for _, section := range []string{"window", "systray", "menu"} {
           var data map[string]any
           if err := cfg.Get(section, &data); err == nil && data != nil {
               s.configData[section] = data
           }
       }
   }
   ```

3. **`handleConfigTask()`**: After updating `configData` in memory, persist via `cfg.Set()` + `cfg.Commit()`:
   ```go
   case window.TaskSaveConfig:
       s.configData["window"] = t.Value
       if s.cfg != nil {
           _ = s.cfg.Set("window", t.Value)
           _ = s.cfg.Commit()
       }
       return nil, true, nil
   ```

4. **`guiConfigPath()`**: Helper returning `~/.core/gui/config.yaml`. Falls back to `.core/gui/config.yaml` in CWD if `$HOME` is unresolvable (shouldn't happen in practice).

### What Does NOT Change

- **`configData` map**: Remains the in-memory cache. `handleConfigQuery` still returns from it. This avoids Viper's type coercion quirks — sub-services get raw `map[string]any`.
- **Sub-service `applyConfig()` methods**: Already wired via IPC. They receive real values once `loadConfig()` populates `configData` from the file. No changes needed — the stubs in `window.applyConfig`, `systray.applyConfig`, and `menu.applyConfig` just need to use the values they already receive.
- **IPC message types**: `QueryConfig`, `TaskSaveConfig` — unchanged.
- **go.mod**: Add `forge.lthn.ai/core/go-config` as a direct dependency. The Go workspace resolves it locally.

### Sub-Service `applyConfig()` Implementations

These are small enhancements to existing stubs (not new architecture):

- **`window.applyConfig()`**: Read `default_width`, `default_height` from config and set on Manager defaults. Read `state_file` to configure StateManager path.
- **`systray.applyConfig()`**: Already partially implemented (reads `tooltip`). Add `icon` path reading.
- **`menu.applyConfig()`**: Read `show_dev_tools` bool and store it for `buildMenu()` to conditionally include the Developer menu.

### Error Handling

Config is non-critical. If the file is missing, malformed, or unreadable:
- `loadConfig()` logs a warning and continues with empty `configData` maps.
- Sub-services receive empty maps and apply their hardcoded defaults.
- `handleConfigTask` saves still work — `config.New` creates the file on first `Commit()`.

### Future Work: Manifest/Slots Integration

The go-scm manifest (`Manifest.Layout`, `Manifest.Slots`) declares which UI components go where (HLCRF positions to component names). Config could store user preferences for those slots (e.g., "I want the terminal in the bottom panel"). This is deferred — the current config covers operational settings only. When implemented, it would add a `layout` section to the config file that overrides manifest slot defaults.

## Dependency Direction

```
pkg/display (orchestrator)
├── imports forge.lthn.ai/core/go-config  // NEW
├── imports pkg/window (message types)
├── imports pkg/systray (message types)
├── imports pkg/menu (message types)
└── imports core/go (DI, IPC)
```

No new imports in sub-packages. Only the display orchestrator touches go-config.

## Licence

EUPL-1.2
