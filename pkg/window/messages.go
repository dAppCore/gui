package window

// WindowInfo contains information about a window.
type WindowInfo struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
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

// TaskRestore restores a maximised or minimised window to its normal state.
type TaskRestore struct{ Name string }

// TaskSetTitle changes a window's title.
type TaskSetTitle struct {
	Name  string
	Title string
}

// TaskSetVisibility shows or hides a window.
type TaskSetVisibility struct {
	Name    string
	Visible bool
}

// TaskFullscreen enters or exits fullscreen mode.
type TaskFullscreen struct {
	Name       string
	Fullscreen bool
}

// --- Layout Queries ---

// QueryLayoutList returns summaries of all saved layouts. Result: []LayoutInfo
type QueryLayoutList struct{}

// QueryLayoutGet returns a layout by name. Result: *Layout (nil if not found)
type QueryLayoutGet struct{ Name string }

// --- Layout Tasks ---

// TaskSaveLayout saves the current window arrangement as a named layout. Result: bool
type TaskSaveLayout struct{ Name string }

// TaskRestoreLayout restores a saved layout by name.
type TaskRestoreLayout struct{ Name string }

// TaskDeleteLayout removes a saved layout by name.
type TaskDeleteLayout struct{ Name string }

// TaskTileWindows arranges windows in a tiling mode.
type TaskTileWindows struct {
	Mode    string   // "left-right", "grid", "left-half", "right-half", etc.
	Windows []string // window names; empty = all
}

// TaskSnapWindow snaps a window to a screen edge/corner.
type TaskSnapWindow struct {
	Name     string // window name
	Position string // "left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right", "center"
}

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

type ActionFilesDropped struct {
	Name     string   `json:"name"`     // window name
	Paths    []string `json:"paths"`
	TargetID string   `json:"targetId,omitempty"`
}
