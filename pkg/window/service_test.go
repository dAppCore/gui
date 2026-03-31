package window

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestWindowService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestTaskOpenWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{
		Window: &Window{Name: "test", URL: "/"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestTaskOpenWindow_OptionsFallback_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{
		Options: []WindowOption{WithName("test-fallback"), WithURL("/")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test-fallback", info.Name)
}

func TestTaskOpenWindow_Bad(t *testing.T) {
	// No window service registered — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskOpenWindow{})
	assert.False(t, handled)
}

func TestQueryWindowList_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("a")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("b")}})

	result, handled, err := c.QUERY(QueryWindowList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]WindowInfo)
	assert.Len(t, list, 2)
}

func TestQueryWindowByName_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	result, handled, err := c.QUERY(QueryWindowByName{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(*WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestQueryWindowByName_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.QUERY(QueryWindowByName{Name: "nonexistent"})
	require.NoError(t, err)
	assert.True(t, handled) // handled=true, result is nil (not found)
	assert.Nil(t, result)
}

func TestTaskCloseWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify window is removed
	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	assert.Nil(t, result)
}

func TestTaskCloseWindow_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetPosition_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetPosition{Name: "test", X: 100, Y: 200})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
}

func TestTaskSetSize_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetSize{Name: "test", Width: 800, Height: 600})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestTaskMaximise_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskMaximise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.True(t, info.Maximized)
}

func TestFileDrop_Good(t *testing.T) {
	_, c := newTestWindowService(t)

	// Open a window
	result, _, _ := c.PERFORM(TaskOpenWindow{
		Options: []WindowOption{WithName("drop-test")},
	})
	info := result.(WindowInfo)
	assert.Equal(t, "drop-test", info.Name)

	// Capture broadcast actions
	var dropped ActionFilesDropped
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionFilesDropped); ok {
			mu.Lock()
			dropped = a
			mu.Unlock()
		}
		return nil
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
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskMinimise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.minimised)
}

func TestTaskMinimise_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskMinimise{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskFocus ---

func TestTaskFocus_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskFocus{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.focused)
}

func TestTaskFocus_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskFocus{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskRestore ---

func TestTaskRestore_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// First maximise, then restore
	_, _, _ = c.PERFORM(TaskMaximise{Name: "test"})

	_, handled, err := c.PERFORM(TaskRestore{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

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
	_, handled, err := c.PERFORM(TaskRestore{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetTitle ---

func TestTaskSetTitle_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetTitle{Name: "test", Title: "New Title"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	assert.Equal(t, "New Title", pw.Title())
}

func TestTaskSetTitle_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetTitle{Name: "nonexistent", Title: "Nope"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetAlwaysOnTop ---

func TestTaskSetAlwaysOnTop_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetAlwaysOnTop{Name: "test", AlwaysOnTop: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.alwaysOnTop)
}

func TestTaskSetAlwaysOnTop_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetAlwaysOnTop{Name: "nonexistent", AlwaysOnTop: true})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetBackgroundColour ---

func TestTaskSetBackgroundColour_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetBackgroundColour{
		Name: "test", Red: 10, Green: 20, Blue: 30, Alpha: 40,
	})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, [4]uint8{10, 20, 30, 40}, mw.backgroundColour)
}

func TestTaskSetBackgroundColour_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetBackgroundColour{Name: "nonexistent", Red: 1, Green: 2, Blue: 3, Alpha: 4})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetVisibility ---

func TestTaskSetVisibility_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetVisibility{Name: "test", Visible: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.visible)

	// Now hide it
	_, handled, err = c.PERFORM(TaskSetVisibility{Name: "test", Visible: false})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, mw.visible)
}

func TestTaskSetVisibility_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetVisibility{Name: "nonexistent", Visible: true})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskFullscreen ---

func TestTaskFullscreen_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Enter fullscreen
	_, handled, err := c.PERFORM(TaskFullscreen{Name: "test", Fullscreen: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.fullscreened)

	// Exit fullscreen
	_, handled, err = c.PERFORM(TaskFullscreen{Name: "test", Fullscreen: false})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, mw.fullscreened)
}

func TestTaskFullscreen_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskFullscreen{Name: "nonexistent", Fullscreen: true})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSaveLayout ---

func TestTaskSaveLayout_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(960, 1080), WithPosition(0, 0)}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(960, 1080), WithPosition(960, 0)}})

	_, handled, err := c.PERFORM(TaskSaveLayout{Name: "coding"})
	require.NoError(t, err)
	assert.True(t, handled)

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
	_, handled, err := c.PERFORM(TaskSaveLayout{Name: ""})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskRestoreLayout ---

func TestTaskRestoreLayout_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	// Open windows
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(800, 600), WithPosition(0, 0)}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(800, 600), WithPosition(0, 0)}})

	// Save a layout with specific positions
	_, _, _ = c.PERFORM(TaskSaveLayout{Name: "coding"})

	// Move the windows to different positions
	_, _, _ = c.PERFORM(TaskSetPosition{Name: "editor", X: 500, Y: 500})
	_, _, _ = c.PERFORM(TaskSetPosition{Name: "terminal", X: 600, Y: 600})

	// Restore the layout
	_, handled, err := c.PERFORM(TaskRestoreLayout{Name: "coding"})
	require.NoError(t, err)
	assert.True(t, handled)

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
	_, handled, err := c.PERFORM(TaskRestoreLayout{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskStackWindows ---

func TestTaskStackWindows_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("s1"), WithSize(800, 600)}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("s2"), WithSize(800, 600)}})

	_, handled, err := c.PERFORM(TaskStackWindows{Windows: []string{"s1", "s2"}, OffsetX: 25, OffsetY: 35})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("s2")
	require.True(t, ok)
	x, y := pw.Position()
	assert.Equal(t, 25, x)
	assert.Equal(t, 35, y)
}

// --- TaskApplyWorkflow ---

func TestTaskApplyWorkflow_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("editor"), WithSize(800, 600)}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("terminal"), WithSize(800, 600)}})

	_, handled, err := c.PERFORM(TaskApplyWorkflow{Workflow: "side-by-side"})
	require.NoError(t, err)
	assert.True(t, handled)

	editor, ok := svc.Manager().Get("editor")
	require.True(t, ok)
	x, y := editor.Position()
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	terminal, ok := svc.Manager().Get("terminal")
	require.True(t, ok)
	x, y = terminal.Position()
	assert.Equal(t, 960, x)
	assert.Equal(t, 0, y)
}

// --- Zoom ---

func TestQueryWindowZoom_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	result, handled, err := c.QUERY(QueryWindowZoom{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1.0, result.(float64))
}

func TestQueryWindowZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.QUERY(QueryWindowZoom{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetZoom_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetZoom{Name: "test", Magnification: 1.5})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 1.5, mw.zoom)
}

func TestTaskSetZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetZoom{Name: "nonexistent", Magnification: 1.5})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskZoomIn_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskZoomIn{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.InDelta(t, 1.1, mw.zoom, 0.001)
}

func TestTaskZoomIn_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskZoomIn{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskZoomOut_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})
	// Set zoom to 1.5 first so we can decrease it
	_, _, _ = c.PERFORM(TaskSetZoom{Name: "test", Magnification: 1.5})

	_, handled, err := c.PERFORM(TaskZoomOut{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.InDelta(t, 1.4, mw.zoom, 0.001)
}

func TestTaskZoomOut_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskZoomOut{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskZoomReset_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})
	_, _, _ = c.PERFORM(TaskSetZoom{Name: "test", Magnification: 2.0})

	_, handled, err := c.PERFORM(TaskZoomReset{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 1.0, mw.zoom)
}

func TestTaskZoomReset_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskZoomReset{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- Content ---

func TestTaskSetURL_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetURL{Name: "test", URL: "https://example.com"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "https://example.com", mw.url)
}

func TestTaskSetURL_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetURL{Name: "nonexistent", URL: "https://example.com"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetHTML_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetHTML{Name: "test", HTML: "<h1>Hello</h1>"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "<h1>Hello</h1>", mw.html)
}

func TestTaskSetHTML_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetHTML{Name: "nonexistent", HTML: "<h1>Hello</h1>"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskExecJS_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskExecJS{Name: "test", JS: "document.title = 'Ready'"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Contains(t, mw.execJSCalls, "document.title = 'Ready'")
}

func TestTaskExecJS_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskExecJS{Name: "nonexistent", JS: "alert(1)"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- State toggles ---

func TestTaskToggleFullscreen_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Toggle on
	_, handled, err := c.PERFORM(TaskToggleFullscreen{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.fullscreened)

	// Toggle off
	_, _, _ = c.PERFORM(TaskToggleFullscreen{Name: "test"})
	assert.False(t, mw.fullscreened)
}

func TestTaskToggleFullscreen_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskToggleFullscreen{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskToggleMaximise_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	// Toggle on
	_, handled, err := c.PERFORM(TaskToggleMaximise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.maximised)

	// Toggle off
	_, _, _ = c.PERFORM(TaskToggleMaximise{Name: "test"})
	assert.False(t, mw.maximised)
}

func TestTaskToggleMaximise_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskToggleMaximise{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- Bounds ---

func TestQueryWindowBounds_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{
		WithName("test"), WithSize(800, 600), WithPosition(100, 200),
	}})

	result, handled, err := c.QUERY(QueryWindowBounds{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	bounds := result.(WindowBounds)
	assert.Equal(t, 100, bounds.X)
	assert.Equal(t, 200, bounds.Y)
	assert.Equal(t, 800, bounds.Width)
	assert.Equal(t, 600, bounds.Height)
}

func TestQueryWindowBounds_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.QUERY(QueryWindowBounds{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetBounds_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetBounds{Name: "test", X: 50, Y: 75, Width: 1024, Height: 768})
	require.NoError(t, err)
	assert.True(t, handled)

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
	_, handled, err := c.PERFORM(TaskSetBounds{Name: "nonexistent", X: 0, Y: 0, Width: 800, Height: 600})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- Content protection ---

func TestTaskSetContentProtection_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetContentProtection{Name: "test", Protection: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.contentProtection)
}

func TestTaskSetContentProtection_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetContentProtection{Name: "nonexistent", Protection: true})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- Flash ---

func TestTaskFlash_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskFlash{Name: "test", Enabled: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.flashed)
}

func TestTaskFlash_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskFlash{Name: "nonexistent", Enabled: true})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- Print ---

func TestTaskPrint_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskPrint{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskPrint_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskPrint{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- State queries (IsVisible, IsFullscreen, IsMinimised) ---

func TestQueryWindowBounds_Ugly(t *testing.T) {
	// Verify bounds reflect position and size changes
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("test"), WithSize(1280, 800)}})
	_, _, _ = c.PERFORM(TaskSetBounds{Name: "test", X: 10, Y: 20, Width: 640, Height: 480})

	result, _, err := c.QUERY(QueryWindowBounds{Name: "test"})
	require.NoError(t, err)
	bounds := result.(WindowBounds)
	assert.Equal(t, 10, bounds.X)
	assert.Equal(t, 20, bounds.Y)
	assert.Equal(t, 640, bounds.Width)
	assert.Equal(t, 480, bounds.Height)
}
