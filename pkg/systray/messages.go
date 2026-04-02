package systray

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
// Use: result, _, err := c.QUERY(systray.QueryConfig{})
type QueryConfig struct{}

// --- Tasks ---

// TaskSetTrayIcon sets the tray icon.
// Use: _, _, err := c.PERFORM(systray.TaskSetTrayIcon{Data: iconBytes})
type TaskSetTrayIcon struct{ Data []byte }

// TaskSetTooltip updates the tray tooltip text.
// Use: _, _, err := c.PERFORM(systray.TaskSetTooltip{Tooltip: "Core is ready"})
type TaskSetTooltip struct{ Tooltip string }

// TaskSetLabel updates the tray label text.
// Use: _, _, err := c.PERFORM(systray.TaskSetLabel{Label: "Core"})
type TaskSetLabel struct{ Label string }

// TaskSetTrayMenu sets the tray menu items.
// Use: _, _, err := c.PERFORM(systray.TaskSetTrayMenu{Items: items})
type TaskSetTrayMenu struct{ Items []TrayMenuItem }

// TaskShowPanel shows the tray panel window.
// Use: _, _, err := c.PERFORM(systray.TaskShowPanel{})
type TaskShowPanel struct{}

// TaskHidePanel hides the tray panel window.
// Use: _, _, err := c.PERFORM(systray.TaskHidePanel{})
type TaskHidePanel struct{}

// TaskShowMessage shows a tray message or notification.
// Use: _, _, err := c.PERFORM(systray.TaskShowMessage{Title: "Core", Message: "Sync complete"})
type TaskShowMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaskSaveConfig persists this service's config section via the display orchestrator.
// Use: _, _, err := c.PERFORM(systray.TaskSaveConfig{Value: map[string]any{"tooltip": "Core"}})
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions ---

// ActionTrayClicked is broadcast when the tray icon is clicked.
// Use: _ = c.ACTION(systray.ActionTrayClicked{})
type ActionTrayClicked struct{}

// ActionTrayMenuItemClicked is broadcast when a tray menu item is clicked.
// Use: _ = c.ACTION(systray.ActionTrayMenuItemClicked{ActionID: "quit"})
type ActionTrayMenuItemClicked struct{ ActionID string }
