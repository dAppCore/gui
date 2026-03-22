package menu

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// QueryGetAppMenu returns the current app menu item descriptors.
// Result: []MenuItem
type QueryGetAppMenu struct{}

// TaskSetAppMenu sets the application menu. OnClick closures work because
// core/go IPC is in-process (no serialisation boundary).
type TaskSetAppMenu struct{ Items []MenuItem }

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }
