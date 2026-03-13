package window

// WindowInfo contains information about a window.
type WindowInfo struct {
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Maximized bool   `json:"maximized"`
	Focused   bool   `json:"focused"`
}

// --- Queries (read-only) ---

// QueryWindowList returns all tracked windows. Result: []WindowInfo
type QueryWindowList struct{}

// QueryWindowByName returns a single window by name. Result: *WindowInfo (nil if not found)
type QueryWindowByName struct{ Name string }

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// --- Tasks (side-effects) ---

// TaskOpenWindow creates a new window. Result: WindowInfo
type TaskOpenWindow struct{ Opts []WindowOption }

// TaskCloseWindow closes a window. Handler persists state BEFORE emitting ActionWindowClosed.
type TaskCloseWindow struct{ Name string }

// TaskSetPosition moves a window.
type TaskSetPosition struct {
	Name string
	X, Y int
}

// TaskSetSize resizes a window.
type TaskSetSize struct {
	Name string
	W, H int
}

// TaskMaximise maximises a window.
type TaskMaximise struct{ Name string }

// TaskMinimise minimises a window.
type TaskMinimise struct{ Name string }

// TaskFocus brings a window to the front.
type TaskFocus struct{ Name string }

// TaskSaveConfig persists this service's config section via the display orchestrator.
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions (broadcasts) ---

type ActionWindowOpened struct{ Name string }
type ActionWindowClosed struct{ Name string }

type ActionWindowMoved struct {
	Name string
	X, Y int
}

type ActionWindowResized struct {
	Name string
	W, H int
}

type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }
