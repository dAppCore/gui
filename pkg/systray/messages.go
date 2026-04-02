package systray

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// --- Tasks ---

// TaskSetTrayIcon sets the tray icon.
type TaskSetTrayIcon struct{ Data []byte }

// TaskSetTooltip updates the tray tooltip text.
type TaskSetTooltip struct{ Tooltip string }

// TaskSetLabel updates the tray label text.
type TaskSetLabel struct{ Label string }

// TaskSetTrayMenu sets the tray menu items.
type TaskSetTrayMenu struct{ Items []TrayMenuItem }

// TaskShowPanel shows the tray panel window.
type TaskShowPanel struct{}

// TaskHidePanel hides the tray panel window.
type TaskHidePanel struct{}

// TaskShowMessage shows a tray message or notification.
type TaskShowMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions ---

// ActionTrayClicked is broadcast when the tray icon is clicked.
type ActionTrayClicked struct{}

// ActionTrayMenuItemClicked is broadcast when a tray menu item is clicked.
type ActionTrayMenuItemClicked struct{ ActionID string }
