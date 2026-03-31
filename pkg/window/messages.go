package window

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

type TaskSaveLayout struct{ Name string }

type TaskRestoreLayout struct{ Name string }

type TaskDeleteLayout struct{ Name string }

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

type ActionWindowOpened struct{ Name string }
type ActionWindowClosed struct{ Name string }

type ActionWindowMoved struct {
	Name string
	X, Y int
}

type ActionWindowResized struct {
	Name          string
	Width, Height int
}

type ActionWindowFocused struct{ Name string }
type ActionWindowBlurred struct{ Name string }

type ActionFilesDropped struct {
	Name     string   `json:"name"` // window name
	Paths    []string `json:"paths"`
	TargetID string   `json:"targetId,omitempty"`
}
