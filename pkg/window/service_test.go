package window

import (
	"context"
	"sync"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestWindowService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestTaskOpenWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.open", TaskOpenWindow{
		Window: &Window{Name: "test", URL: "/"},
	})
	require.True(t, r.OK)
	info := r.Value.(WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestTaskOpenWindow_OptionsFallback_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.open", TaskOpenWindow{
		Options: []WindowOption{WithName("test-fallback"), WithURL("/")},
	})
	require.True(t, r.OK)
	info := r.Value.(WindowInfo)
	assert.Equal(t, "test-fallback", info.Name)
}

func TestTaskOpenWindow_Bad(t *testing.T) {
	// No window service registered — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("window.open").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestQueryWindowList_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("a")}})
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("b")}})

	r := c.QUERY(QueryWindowList{})
	require.True(t, r.OK)
	list := r.Value.([]WindowInfo)
	assert.Len(t, list, 2)
}

func TestQueryWindowByName_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := c.QUERY(QueryWindowByName{Name: "test"})
	require.True(t, r.OK)
	info := r.Value.(*WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestQueryWindowByName_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := c.QUERY(QueryWindowByName{Name: "nonexistent"})
	require.True(t, r.OK) // handled=true, result is nil (not found)
	assert.Nil(t, r.Value)
}

func TestTaskCloseWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.close", TaskCloseWindow{Name: "test"})
	require.True(t, r.OK)

	// Verify window is removed
	r2 := c.QUERY(QueryWindowByName{Name: "test"})
	assert.Nil(t, r2.Value)
}

func TestTaskCloseWindow_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.close", TaskCloseWindow{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskSetPosition_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setPosition", TaskSetPosition{Name: "test", X: 100, Y: 200})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "test"})
	info := r2.Value.(*WindowInfo)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
}

func TestTaskSetSize_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setSize", TaskSetSize{Name: "test", Width: 800, Height: 600})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "test"})
	info := r2.Value.(*WindowInfo)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestTaskMaximise_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.maximise", TaskMaximise{Name: "test"})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "test"})
	info := r2.Value.(*WindowInfo)
	assert.True(t, info.Maximized)
}

func TestFileDrop_Good(t *testing.T) {
	_, c := newTestWindowService(t)

	// Open a window
	r := taskRun(c, "window.open", TaskOpenWindow{
		Options: []WindowOption{WithName("drop-test")},
	})
	info := r.Value.(WindowInfo)
	assert.Equal(t, "drop-test", info.Name)

	// Capture broadcast actions
	var dropped ActionFilesDropped
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionFilesDropped); ok {
			mu.Lock()
			dropped = a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	// Get the mock window and simulate file drop
	svc := core.MustServiceFor[*Service](c, "window")
	pw, ok := svc.Manager().Get("drop-test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	mw.emitFileDrop([]string{"/tmp/file1.txt", "/tmp/file2.txt"}, "upload-zone")

	mu.Lock()
	assert.Equal(t, "drop-test", dropped.Name)
	assert.Equal(t, []string{"/tmp/file1.txt", "/tmp/file2.txt"}, dropped.Paths)
	assert.Equal(t, "upload-zone", dropped.TargetID)
	mu.Unlock()
}

// --- TaskMinimise ---

func TestTaskMinimise_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.minimise", TaskMinimise{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.minimised)
}

func TestTaskMinimise_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.minimise", TaskMinimise{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- TaskFocus ---

func TestTaskFocus_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.focus", TaskFocus{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.focused)
}

func TestTaskFocus_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.focus", TaskFocus{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- TaskRestore ---

func TestTaskRestore_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// First maximise, then restore
	taskRun(c, "window.maximise", TaskMaximise{Name: "test"})

	r := taskRun(c, "window.restore", TaskRestore{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.False(t, mw.maximised)

	// Verify state was updated
	state, ok := svc.Manager().State().GetState("test")
	assert.True(t, ok)
	assert.False(t, state.Maximized)
}

func TestTaskRestore_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.restore", TaskRestore{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- TaskSetTitle ---

func TestTaskSetTitle_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setTitle", TaskSetTitle{Name: "test", Title: "New Title"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	assert.Equal(t, "New Title", pw.Title())
}

func TestTaskSetTitle_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setTitle", TaskSetTitle{Name: "nonexistent", Title: "Nope"})
	assert.False(t, r.OK)
}

// --- TaskSetAlwaysOnTop ---

func TestTaskSetAlwaysOnTop_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setAlwaysOnTop", TaskSetAlwaysOnTop{Name: "test", AlwaysOnTop: true})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.alwaysOnTop)
}

func TestTaskSetAlwaysOnTop_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setAlwaysOnTop", TaskSetAlwaysOnTop{Name: "nonexistent", AlwaysOnTop: true})
	assert.False(t, r.OK)
}

// --- TaskSetBackgroundColour ---

func TestTaskSetBackgroundColour_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setBackgroundColour", TaskSetBackgroundColour{
		Name: "test", Red: 10, Green: 20, Blue: 30, Alpha: 40,
	})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, [4]uint8{10, 20, 30, 40}, mw.backgroundColour)
}

func TestTaskSetBackgroundColour_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setBackgroundColour", TaskSetBackgroundColour{Name: "nonexistent", Red: 1, Green: 2, Blue: 3, Alpha: 4})
	assert.False(t, r.OK)
}

// --- TaskSetVisibility ---

func TestTaskSetVisibility_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setVisibility", TaskSetVisibility{Name: "test", Visible: true})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.visible)

	// Now hide it
	r2 := taskRun(c, "window.setVisibility", TaskSetVisibility{Name: "test", Visible: false})
	require.True(t, r2.OK)
	assert.False(t, mw.visible)
}

func TestTaskSetVisibility_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setVisibility", TaskSetVisibility{Name: "nonexistent", Visible: true})
	assert.False(t, r.OK)
}

// --- TaskFullscreen ---

func TestTaskFullscreen_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Enter fullscreen
	r := taskRun(c, "window.fullscreen", TaskFullscreen{Name: "test", Fullscreen: true})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.fullscreened)

	// Exit fullscreen
	r2 := taskRun(c, "window.fullscreen", TaskFullscreen{Name: "test", Fullscreen: false})
	require.True(t, r2.OK)
	assert.False(t, mw.fullscreened)
}

func TestTaskFullscreen_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.fullscreen", TaskFullscreen{Name: "nonexistent", Fullscreen: true})
	assert.False(t, r.OK)
}

// --- TaskSaveLayout ---

func TestTaskSaveLayout_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(960, 1080), WithPosition(0, 0)}})
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(960, 1080), WithPosition(960, 0)}})

	r := taskRun(c, "window.saveLayout", TaskSaveLayout{Name: "coding"})
	require.True(t, r.OK)

	// Verify layout was saved with correct window states
	layout, ok := svc.Manager().Layout().GetLayout("coding")
	assert.True(t, ok)
	assert.Equal(t, "coding", layout.Name)
	assert.Len(t, layout.Windows, 2)

	editorState, ok := layout.Windows["editor"]
	assert.True(t, ok)
	assert.Equal(t, 0, editorState.X)
	assert.Equal(t, 960, editorState.Width)

	termState, ok := layout.Windows["terminal"]
	assert.True(t, ok)
	assert.Equal(t, 960, termState.X)
	assert.Equal(t, 960, termState.Width)
}

func TestTaskSaveLayout_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	// Saving an empty layout with empty name returns an error from LayoutManager
	r := taskRun(c, "window.saveLayout", TaskSaveLayout{Name: ""})
	assert.False(t, r.OK)
}

// --- TaskRestoreLayout ---

func TestTaskRestoreLayout_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	// Open windows
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(800, 600), WithPosition(0, 0)}})
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(800, 600), WithPosition(0, 0)}})

	// Save a layout with specific positions
	taskRun(c, "window.saveLayout", TaskSaveLayout{Name: "coding"})

	// Move the windows to different positions
	taskRun(c, "window.setPosition", TaskSetPosition{Name: "editor", X: 500, Y: 500})
	taskRun(c, "window.setPosition", TaskSetPosition{Name: "terminal", X: 600, Y: 600})

	// Restore the layout
	r := taskRun(c, "window.restoreLayout", TaskRestoreLayout{Name: "coding"})
	require.True(t, r.OK)

	// Verify windows were moved back to saved positions
	pw, ok := svc.Manager().Get("editor")
	require.True(t, ok)
	x, y := pw.Position()
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	pw2, ok := svc.Manager().Get("terminal")
	require.True(t, ok)
	x2, y2 := pw2.Position()
	assert.Equal(t, 0, x2)
	assert.Equal(t, 0, y2)

	editorState, ok := svc.Manager().State().GetState("editor")
	require.True(t, ok)
	assert.Equal(t, 0, editorState.X)
	assert.Equal(t, 0, editorState.Y)

	terminalState, ok := svc.Manager().State().GetState("terminal")
	require.True(t, ok)
	assert.Equal(t, 0, terminalState.X)
	assert.Equal(t, 0, terminalState.Y)
}

func TestTaskRestoreLayout_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.restoreLayout", TaskRestoreLayout{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- TaskStackWindows ---

func TestTaskStackWindows_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("s1"), WithSize(800, 600)}})
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("s2"), WithSize(800, 600)}})

	r := taskRun(c, "window.stackWindows", TaskStackWindows{Windows: []string{"s1", "s2"}, OffsetX: 25, OffsetY: 35})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("s2")
	require.True(t, ok)
	x, y := pw.Position()
	assert.Equal(t, 25, x)
	assert.Equal(t, 35, y)
}

// --- TaskApplyWorkflow ---

func TestTaskApplyWorkflow_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(800, 600)}})
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(800, 600)}})

	r := taskRun(c, "window.applyWorkflow", TaskApplyWorkflow{Workflow: "side-by-side"})
	require.True(t, r.OK)

	editor, ok := svc.Manager().Get("editor")
	require.True(t, ok)
	terminal, ok := svc.Manager().Get("terminal")
	require.True(t, ok)
	editorX, editorY := editor.Position()
	terminalX, terminalY := terminal.Position()

	assert.Equal(t, 0, editorY)
	assert.Equal(t, 0, terminalY)
	assert.ElementsMatch(t, []int{0, 960}, []int{editorX, terminalX})

	// The assignment order is derived from map iteration, so only the
	// geometry matters here.
	assert.Contains(t, []int{editorX, terminalX}, 0)
	assert.Contains(t, []int{editorX, terminalX}, 960)
}

// --- Zoom ---

func TestQueryWindowZoom_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := c.QUERY(QueryWindowZoom{Name: "test"})
	require.True(t, r.OK)
	assert.Equal(t, 1.0, r.Value.(float64))
}

func TestQueryWindowZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := c.QUERY(QueryWindowZoom{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskSetZoom_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setZoom", TaskSetZoom{Name: "test", Magnification: 1.5})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 1.5, mw.zoom)
}

func TestTaskSetZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setZoom", TaskSetZoom{Name: "nonexistent", Magnification: 1.5})
	assert.False(t, r.OK)
}

func TestTaskZoomIn_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.zoomIn", TaskZoomIn{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.InDelta(t, 1.1, mw.zoom, 0.001)
}

func TestTaskZoomIn_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.zoomIn", TaskZoomIn{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskZoomOut_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})
	// Set zoom to 1.5 first so we can decrease it
	taskRun(c, "window.setZoom", TaskSetZoom{Name: "test", Magnification: 1.5})

	r := taskRun(c, "window.zoomOut", TaskZoomOut{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.InDelta(t, 1.4, mw.zoom, 0.001)
}

func TestTaskZoomOut_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.zoomOut", TaskZoomOut{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskZoomReset_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})
	taskRun(c, "window.setZoom", TaskSetZoom{Name: "test", Magnification: 2.0})

	r := taskRun(c, "window.zoomReset", TaskZoomReset{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 1.0, mw.zoom)
}

func TestTaskZoomReset_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.zoomReset", TaskZoomReset{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- Content ---

func TestTaskSetURL_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setURL", TaskSetURL{Name: "test", URL: "https://example.com"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "https://example.com", mw.url)
}

func TestTaskSetURL_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setURL", TaskSetURL{Name: "nonexistent", URL: "https://example.com"})
	assert.False(t, r.OK)
}

func TestTaskSetHTML_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setHTML", TaskSetHTML{Name: "test", HTML: "<h1>Hello</h1>"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "<h1>Hello</h1>", mw.html)
}

func TestTaskSetHTML_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setHTML", TaskSetHTML{Name: "nonexistent", HTML: "<h1>Hello</h1>"})
	assert.False(t, r.OK)
}

func TestTaskExecJS_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.execJS", TaskExecJS{Name: "test", JS: "document.title = 'Ready'"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Contains(t, mw.execJSCalls, "document.title = 'Ready'")
}

func TestTaskExecJS_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.execJS", TaskExecJS{Name: "nonexistent", JS: "alert(1)"})
	assert.False(t, r.OK)
}

// --- State toggles ---

func TestTaskToggleFullscreen_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Toggle on
	r := taskRun(c, "window.toggleFullscreen", TaskToggleFullscreen{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.fullscreened)

	// Toggle off
	taskRun(c, "window.toggleFullscreen", TaskToggleFullscreen{Name: "test"})
	assert.False(t, mw.fullscreened)
}

func TestTaskToggleFullscreen_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.toggleFullscreen", TaskToggleFullscreen{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskToggleMaximise_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Toggle on
	r := taskRun(c, "window.toggleMaximise", TaskToggleMaximise{Name: "test"})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.maximised)

	// Toggle off
	taskRun(c, "window.toggleMaximise", TaskToggleMaximise{Name: "test"})
	assert.False(t, mw.maximised)
}

func TestTaskToggleMaximise_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.toggleMaximise", TaskToggleMaximise{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- Bounds ---

func TestQueryWindowBounds_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("test"), WithSize(800, 600), WithPosition(100, 200),
	}})

	r := c.QUERY(QueryWindowBounds{Name: "test"})
	require.True(t, r.OK)

	bounds := r.Value.(WindowBounds)
	assert.Equal(t, 100, bounds.X)
	assert.Equal(t, 200, bounds.Y)
	assert.Equal(t, 800, bounds.Width)
	assert.Equal(t, 600, bounds.Height)
}

func TestQueryWindowBounds_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := c.QUERY(QueryWindowBounds{Name: "nonexistent"})
	assert.False(t, r.OK)
}

func TestTaskSetBounds_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setBounds", TaskSetBounds{Name: "test", X: 50, Y: 75, Width: 1024, Height: 768})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 50, mw.x)
	assert.Equal(t, 75, mw.y)
	assert.Equal(t, 1024, mw.width)
	assert.Equal(t, 768, mw.height)
}

func TestTaskSetBounds_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setBounds", TaskSetBounds{Name: "nonexistent", X: 0, Y: 0, Width: 800, Height: 600})
	assert.False(t, r.OK)
}

// --- Content protection ---

func TestTaskSetContentProtection_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.setContentProtection", TaskSetContentProtection{Name: "test", Protection: true})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.contentProtection)
}

func TestTaskSetContentProtection_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.setContentProtection", TaskSetContentProtection{Name: "nonexistent", Protection: true})
	assert.False(t, r.OK)
}

// --- Flash ---

func TestTaskFlash_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.flash", TaskFlash{Name: "test", Enabled: true})
	require.True(t, r.OK)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.flashed)
}

func TestTaskFlash_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.flash", TaskFlash{Name: "nonexistent", Enabled: true})
	assert.False(t, r.OK)
}

// --- Print ---

func TestTaskPrint_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	r := taskRun(c, "window.print", TaskPrint{Name: "test"})
	require.True(t, r.OK)
}

func TestTaskPrint_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	r := taskRun(c, "window.print", TaskPrint{Name: "nonexistent"})
	assert.False(t, r.OK)
}

// --- State queries (IsVisible, IsFullscreen, IsMinimised) ---

func TestQueryWindowBounds_Ugly(t *testing.T) {
	// Verify bounds reflect position and size changes
	_, c := newTestWindowService(t)
	taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("test"), WithSize(1280, 800)}})
	taskRun(c, "window.setBounds", TaskSetBounds{Name: "test", X: 10, Y: 20, Width: 640, Height: 480})

	r := c.QUERY(QueryWindowBounds{Name: "test"})
	require.True(t, r.OK)
	bounds := r.Value.(WindowBounds)
	assert.Equal(t, 10, bounds.X)
	assert.Equal(t, 20, bounds.Y)
	assert.Equal(t, 640, bounds.Width)
	assert.Equal(t, 480, bounds.Height)
}
