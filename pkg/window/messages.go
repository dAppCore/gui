// pkg/window/messages.go
package window

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

type QueryWindowList struct{}

type QueryWindowByName struct{ Name string }

type QueryConfig struct{}

// TaskOpenWindow opens a concrete Window descriptor.
// window.TaskOpenWindow{Window: &window.Window{Name: "settings", URL: "/", Width: 800, Height: 600}}
type TaskOpenWindow struct{ Window *Window }

type TaskCloseWindow struct{ Name string }

type TaskSetPosition struct {
	Name string
	X, Y int
}

type TaskSetSize struct {
	Name          string
	Width, Height int
}

type TaskMaximise struct{ Name string }

type TaskMinimise struct{ Name string }

type TaskFocus struct{ Name string }

type TaskRestore struct{ Name string }

type TaskSetTitle struct {
	Name  string
	Title string
}

type TaskSetAlwaysOnTop struct {
	Name        string
	AlwaysOnTop bool
}

type TaskSetBackgroundColour struct {
	Name  string
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type TaskSetVisibility struct {
	Name    string
	Visible bool
}

type TaskFullscreen struct {
	Name       string
	Fullscreen bool
}

type QueryLayoutList struct{}

type QueryLayoutGet struct{ Name string }

type QueryFindSpace struct {
	Width        int `json:"width"`
	Height       int `json:"height"`
	ScreenWidth  int `json:"screenWidth"`
	ScreenHeight int `json:"screenHeight"`
}

type QueryLayoutSuggestion struct {
	WindowCount  int `json:"windowCount"`
	ScreenWidth  int `json:"screenWidth"`
	ScreenHeight int `json:"screenHeight"`
}

type TaskSaveLayout struct{ Name string }

type TaskRestoreLayout struct{ Name string }

type TaskDeleteLayout struct{ Name string }

type TaskArrangePair struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

type TaskTileWindows struct {
	Mode    string   // "left-right", "grid", "left-half", "right-half", etc.
	Windows []string // window names; empty = all
}

type TaskStackWindows struct {
	Windows []string // window names; empty = all
	OffsetX int
	OffsetY int
}

type TaskSnapWindow struct {
	Name     string // window name
	Position string // "left", "right", "top", "bottom", "top-left", "top-right", "bottom-left", "bottom-right", "center"
}

type TaskApplyWorkflow struct {
	Workflow string
	Windows  []string // window names; empty = all
}

type TaskSaveConfig struct{ Config map[string]any }

// QueryWindowZoom queries the current zoom factor for a named window. Result: float64
type QueryWindowZoom struct{ Name string }

// QueryWindowBounds queries the current bounds for a named window. Result: *Bounds
type QueryWindowBounds struct{ Name string }

// TaskSetZoom sets the zoom factor for a named window.
// c.PERFORM(window.TaskSetZoom{Name: "main", Factor: 1.5})
type TaskSetZoom struct {
	Name   string
	Factor float64
}

// TaskSetURL navigates a named window to a new URL.
// c.PERFORM(window.TaskSetURL{Name: "main", URL: "/settings"})
type TaskSetURL struct {
	Name string
	URL  string
}

// TaskSetHTML replaces the content of a named window with HTML.
// c.PERFORM(window.TaskSetHTML{Name: "main", HTML: "<h1>Hello</h1>"})
type TaskSetHTML struct {
	Name string
	HTML string
}

// TaskExecJS evaluates JavaScript in a named window.
// c.PERFORM(window.TaskExecJS{Name: "main", JS: "document.title = 'Updated'"})
type TaskExecJS struct {
	Name string
	JS   string
}

// TaskToggleFullscreen toggles fullscreen on a named window.
type TaskToggleFullscreen struct{ Name string }

// TaskPrint triggers the platform print dialog for a named window.
type TaskPrint struct{ Name string }

// TaskFlash flashes (or stops flashing) the taskbar entry for a named window (Windows).
type TaskFlash struct {
	Name    string
	Enabled bool
}

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
}

// ActionWindowFocused is broadcast when a window gains focus.
// Use: _ = c.ACTION(window.ActionWindowFocused{Name: "editor"})
type ActionWindowFocused struct{ Name string }

// ActionWindowBlurred is broadcast when a window loses focus.
// Use: _ = c.ACTION(window.ActionWindowBlurred{Name: "editor"})
type ActionWindowBlurred struct{ Name string }

// ActionFilesDropped is broadcast when files are dropped onto a window.
// Use: _ = c.ACTION(window.ActionFilesDropped{Name: "editor", Paths: []string{"/tmp/report.pdf"}})
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

type ArrangedPair struct {
	First  Bounds `json:"first"`
	Second Bounds `json:"second"`
}
