// pkg/window/messages.go
package window

// WindowInfo contains information about a window.
// Use: info := window.WindowInfo{Name: "editor", Title: "Core Editor"}
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
// Use: bounds := window.Bounds{X: 10, Y: 10, Width: 1280, Height: 800}
type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// --- Queries (read-only) ---

// QueryWindowList returns all tracked windows. Result: []WindowInfo
// Use: result, _, err := c.QUERY(window.QueryWindowList{})
type QueryWindowList struct{}

// QueryWindowByName returns a single window by name. Result: *WindowInfo (nil if not found)
// Use: result, _, err := c.QUERY(window.QueryWindowByName{Name: "editor"})
type QueryWindowByName struct{ Name string }

// QueryConfig requests this service's config section from the display orchestrator.
// Result: map[string]any
// Use: result, _, err := c.QUERY(window.QueryConfig{})
type QueryConfig struct{}

// QueryWindowBounds returns the current bounds for a window.
// Use: result, _, err := c.QUERY(window.QueryWindowBounds{Name: "editor"})
type QueryWindowBounds struct{ Name string }

// QueryFindSpace returns a suggested free placement for a new window.
// Use: result, _, err := c.QUERY(window.QueryFindSpace{Width: 1280, Height: 800})
type QueryFindSpace struct {
	Width        int
	Height       int
	ScreenWidth  int
	ScreenHeight int
}

// QueryLayoutSuggestion returns a layout recommendation for the current screen.
// Use: result, _, err := c.QUERY(window.QueryLayoutSuggestion{WindowCount: 2})
type QueryLayoutSuggestion struct {
	WindowCount  int
	ScreenWidth  int
	ScreenHeight int
}

// --- Tasks (side-effects) ---

// TaskOpenWindow creates a new window. Result: WindowInfo
// Use: _, _, err := c.PERFORM(window.TaskOpenWindow{Opts: []window.WindowOption{window.WithName("editor")}})
type TaskOpenWindow struct {
	Window *Window
	Opts   []WindowOption
}

// TaskCloseWindow closes a window after persisting state.
// Platform close events emit ActionWindowClosed through the tracked window handler.
// Use: _, _, err := c.PERFORM(window.TaskCloseWindow{Name: "editor"})
type TaskCloseWindow struct{ Name string }

// TaskSetPosition moves a window.
// Use: _, _, err := c.PERFORM(window.TaskSetPosition{Name: "editor", X: 160, Y: 120})
type TaskSetPosition struct {
	Name string
	X, Y int
}

// TaskSetSize resizes a window.
// Use: _, _, err := c.PERFORM(window.TaskSetSize{Name: "editor", Width: 1280, Height: 800})
type TaskSetSize struct {
	Name          string
	Width, Height int
	W, H          int
}

// TaskMaximise maximises a window.
// Use: _, _, err := c.PERFORM(window.TaskMaximise{Name: "editor"})
type TaskMaximise struct{ Name string }

// TaskMinimise minimises a window.
// Use: _, _, err := c.PERFORM(window.TaskMinimise{Name: "editor"})
type TaskMinimise struct{ Name string }

// TaskFocus brings a window to the front.
// Use: _, _, err := c.PERFORM(window.TaskFocus{Name: "editor"})
type TaskFocus struct{ Name string }

// TaskRestore restores a maximised or minimised window to its normal state.
// Use: _, _, err := c.PERFORM(window.TaskRestore{Name: "editor"})
type TaskRestore struct{ Name string }

// TaskSetTitle changes a window's title.
// Use: _, _, err := c.PERFORM(window.TaskSetTitle{Name: "editor", Title: "Core Editor"})
type TaskSetTitle struct {
	Name  string
	Title string
}

// TaskSetAlwaysOnTop pins a window above others.
// Use: _, _, err := c.PERFORM(window.TaskSetAlwaysOnTop{Name: "editor", AlwaysOnTop: true})
type TaskSetAlwaysOnTop struct {
	Name        string
	AlwaysOnTop bool
}

// TaskSetBackgroundColour updates the window background colour.
// Use: _, _, err := c.PERFORM(window.TaskSetBackgroundColour{Name: "editor", Red: 0, Green: 0, Blue: 0, Alpha: 0})
type TaskSetBackgroundColour struct {
	Name  string
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// TaskSetOpacity updates the window opacity as a value between 0 and 1.
// Use: _, _, err := c.PERFORM(window.TaskSetOpacity{Name: "editor", Opacity: 0.85})
type TaskSetOpacity struct {
	Name    string
	Opacity float32
}

// TaskSetVisibility shows or hides a window.
// Use: _, _, err := c.PERFORM(window.TaskSetVisibility{Name: "editor", Visible: false})
type TaskSetVisibility struct {
	Name    string
	Visible bool
}

// TaskFullscreen enters or exits fullscreen mode.
// Use: _, _, err := c.PERFORM(window.TaskFullscreen{Name: "editor", Fullscreen: true})
type TaskFullscreen struct {
	Name       string
	Fullscreen bool
}

// --- Layout Queries ---

// QueryLayoutList returns summaries of all saved layouts. Result: []LayoutInfo
// Use: result, _, err := c.QUERY(window.QueryLayoutList{})
type QueryLayoutList struct{}

// QueryLayoutGet returns a layout by name. Result: *Layout (nil if not found)
// Use: result, _, err := c.QUERY(window.QueryLayoutGet{Name: "coding"})
type QueryLayoutGet struct{ Name string }

// --- Layout Tasks ---

// TaskSaveLayout saves the current window arrangement as a named layout. Result: bool
// Use: _, _, err := c.PERFORM(window.TaskSaveLayout{Name: "coding"})
type TaskSaveLayout struct{ Name string }

// TaskRestoreLayout restores a saved layout by name.
// Use: _, _, err := c.PERFORM(window.TaskRestoreLayout{Name: "coding"})
type TaskRestoreLayout struct{ Name string }

// TaskDeleteLayout removes a saved layout by name.
// Use: _, _, err := c.PERFORM(window.TaskDeleteLayout{Name: "coding"})
type TaskDeleteLayout struct{ Name string }

// TaskTileWindows arranges windows in a tiling mode.
// Use: _, _, err := c.PERFORM(window.TaskTileWindows{Mode: "grid"})
type TaskTileWindows struct {
	Mode    string   // "left-right", "grid", "left-half", "right-half", etc.
	Windows []string // window names; empty = all
}

// TaskSnapWindow snaps a window to a screen edge/corner.
// Use: _, _, err := c.PERFORM(window.TaskSnapWindow{Name: "editor", Position: "left"})
type TaskSnapWindow struct {
	Name     string // window name
	Position string // "left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right", "center"
}

// TaskArrangePair places two windows side-by-side in a balanced split.
// Use: _, _, err := c.PERFORM(window.TaskArrangePair{First: "editor", Second: "terminal"})
type TaskArrangePair struct {
	First  string
	Second string
}

// TaskBesideEditor places a target window beside an editor/IDE window.
// Use: _, _, err := c.PERFORM(window.TaskBesideEditor{Editor: "editor", Window: "terminal"})
type TaskBesideEditor struct {
	Editor string
	Window string
}

// TaskStackWindows cascades windows with a shared offset.
// Use: _, _, err := c.PERFORM(window.TaskStackWindows{Windows: []string{"editor", "terminal"}})
type TaskStackWindows struct {
	Windows []string
	OffsetX int
	OffsetY int
}

// TaskApplyWorkflow applies a predefined workflow layout to windows.
// Use: _, _, err := c.PERFORM(window.TaskApplyWorkflow{Workflow: window.WorkflowCoding})
type TaskApplyWorkflow struct {
	Workflow WorkflowLayout
	Windows  []string
}

// TaskSaveConfig persists this service's config section via the display orchestrator.
// Use: _, _, err := c.PERFORM(window.TaskSaveConfig{Value: map[string]any{"default_width": 1280}})
type TaskSaveConfig struct{ Value map[string]any }

// --- Actions (broadcasts) ---

// ActionWindowOpened is broadcast when a window is created.
// Use: _ = c.ACTION(window.ActionWindowOpened{Name: "editor"})
type ActionWindowOpened struct{ Name string }

// ActionWindowClosed is broadcast when a window is closed.
// Use: _ = c.ACTION(window.ActionWindowClosed{Name: "editor"})
type ActionWindowClosed struct{ Name string }

// ActionWindowMoved is broadcast when a window is moved.
// Use: _ = c.ACTION(window.ActionWindowMoved{Name: "editor", X: 160, Y: 120})
type ActionWindowMoved struct {
	Name string
	X, Y int
}

// ActionWindowResized is broadcast when a window is resized.
// Use: _ = c.ACTION(window.ActionWindowResized{Name: "editor", Width: 1280, Height: 800})
type ActionWindowResized struct {
	Name          string
	Width, Height int
	W, H          int
}

// ActionWindowFocused is broadcast when a window gains focus.
// Use: _ = c.ACTION(window.ActionWindowFocused{Name: "editor"})
type ActionWindowFocused struct{ Name string }

// ActionWindowBlurred is broadcast when a window loses focus.
// Use: _ = c.ACTION(window.ActionWindowBlurred{Name: "editor"})
type ActionWindowBlurred struct{ Name string }

type ActionFilesDropped struct {
	Name     string   `json:"name"` // window name
	Paths    []string `json:"paths"`
	TargetID string   `json:"targetId,omitempty"`
}

// SpaceInfo describes a suggested empty area on the screen.
// Use: info := window.SpaceInfo{X: 160, Y: 120, Width: 1280, Height: 800}
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
// Use: suggestion := window.LayoutSuggestion{Mode: "side-by-side"}
type LayoutSuggestion struct {
	Mode           string `json:"mode"`
	Columns        int    `json:"columns"`
	Rows           int    `json:"rows"`
	PrimaryWidth   int    `json:"primaryWidth"`
	SecondaryWidth int    `json:"secondaryWidth"`
	Description    string `json:"description"`
}
