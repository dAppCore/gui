package window

// WindowInfo contains information about a window.
type WindowInfo struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Visible   bool   `json:"visible"`
	Minimized bool   `json:"minimized"`
	Maximized bool   `json:"maximized"`
	Focused   bool   `json:"focused"`
}

// Bounds describes the position and size of a window.
type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// --- Queries (read-only) ---

// QueryWindowList returns all tracked windows. Result: []WindowInfo
type QueryWindowList struct{}

// QueryWindowByName returns a single window by name. Result: *WindowInfo (nil if not found)
type QueryWindowByName struct{ Name string }

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
type QueryConfig struct{}

// QueryWindowBounds returns the current bounds for a window.
type QueryWindowBounds struct{ Name string }

// QueryFindSpace returns a suggested free placement for a new window.
type QueryFindSpace struct {
	Width  int
	Height int
}

// QueryLayoutSuggestion returns a layout recommendation for the current screen.
type QueryLayoutSuggestion struct {
	WindowCount  int
	ScreenWidth  int
	ScreenHeight int
}

// --- Tasks (side-effects) ---

// TaskOpenWindow creates a new window. Result: WindowInfo
type TaskOpenWindow struct {
	Window *Window
	Opts   []WindowOption
}

// TaskCloseWindow closes a window. Handler persists state BEFORE emitting ActionWindowClosed.
type TaskCloseWindow struct{ Name string }

// TaskSetPosition moves a window.
type TaskSetPosition struct {
	Name string
	X, Y int
}

// TaskSetSize resizes a window.
type TaskSetSize struct {
	Name          string
	Width, Height int
	W, H          int
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

// TaskSetAlwaysOnTop pins a window above others.
type TaskSetAlwaysOnTop struct {
	Name        string
	AlwaysOnTop bool
}

// TaskSetBackgroundColour updates the window background colour.
type TaskSetBackgroundColour struct {
	Name  string
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// TaskSetOpacity updates the window opacity as a value between 0 and 1.
type TaskSetOpacity struct {
	Name    string
	Opacity float32
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

// TaskArrangePair places two windows side-by-side in a balanced split.
type TaskArrangePair struct {
	First  string
	Second string
}

// TaskBesideEditor places a target window beside an editor/IDE window.
type TaskBesideEditor struct {
	Editor string
	Window string
}

// TaskStackWindows cascades windows with a shared offset.
type TaskStackWindows struct {
	Windows []string
	OffsetX int
	OffsetY int
}

// TaskApplyWorkflow applies a predefined workflow layout to windows.
type TaskApplyWorkflow struct {
	Workflow WorkflowLayout
	Windows  []string
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
	Name          string
	Width, Height int
	W, H          int
}

type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }

type ActionFilesDropped struct {
	Name     string   `json:"name"` // window name
	Paths    []string `json:"paths"`
	TargetID string   `json:"targetId,omitempty"`
}

// SpaceInfo describes a suggested empty area on the screen.
type SpaceInfo struct {
	X            int    `json:"x"`
	Y            int    `json:"y"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ScreenWidth  int    `json:"screenWidth"`
	ScreenHeight int    `json:"screenHeight"`
	Reason       string `json:"reason,omitempty"`
}

// LayoutSuggestion describes a recommended layout for a screen.
type LayoutSuggestion struct {
	Mode           string `json:"mode"`
	Columns        int    `json:"columns"`
	Rows           int    `json:"rows"`
	PrimaryWidth   int    `json:"primaryWidth"`
	SecondaryWidth int    `json:"secondaryWidth"`
	Description    string `json:"description"`
}
