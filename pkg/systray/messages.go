package systray

type QueryConfig struct{}

type TaskSetTrayIcon struct{ Data []byte }

type TaskSetTrayMenu struct{ Items []TrayMenuItem }

type TaskShowPanel struct{}

type TaskHidePanel struct{}

type TaskSaveConfig struct{ Config map[string]any }

type ActionTrayClicked struct{}

type ActionTrayMenuItemClicked struct{ ActionID string }
