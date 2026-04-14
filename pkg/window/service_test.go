package window

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"dappco.re/go/core/gui/pkg/screen"
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

func requireOpenWindow(t *testing.T, c *core.Core, window Window) WindowInfo {
	t.Helper()
	result, handled, err := c.PERFORM(TaskOpenWindow{Window: &window})
	require.NoError(t, err)
	require.True(t, handled)
	return result.(WindowInfo)
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestApplyConfig_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)

	svc.applyConfig(map[string]any{
		"default_width":  1500,
		"default_height": 900,
	})

	pw, err := svc.manager.Open()
	require.NoError(t, err)
	w, h := pw.Size()
	assert.Equal(t, 1500, w)
	assert.Equal(t, 900, h)
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

func TestTaskOpenWindow_Bad_MissingWindow(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{})
	assert.True(t, handled)
	assert.Error(t, err)
	assert.Nil(t, result)
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
	_ = requireOpenWindow(t, c, Window{Name: "a"})
	_ = requireOpenWindow(t, c, Window{Name: "b"})

	result, handled, err := c.QUERY(QueryWindowList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]WindowInfo)
	assert.Len(t, list, 2)

	byName := make(map[string]WindowInfo, len(list))
	for _, info := range list {
		byName[info.Name] = info
	}

	assert.True(t, byName["a"].Visible)
	assert.False(t, byName["a"].Minimized)
	assert.False(t, byName["b"].Visible)
	assert.True(t, byName["b"].Minimized)
}

func TestQueryWindowByName_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	result, handled, err := c.QUERY(QueryWindowByName{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(*WindowInfo)
	assert.Equal(t, "test", info.Name)
	assert.True(t, info.Visible)
	assert.False(t, info.Minimized)
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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskSetSize{Name: "test", Width: 800, Height: 600})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestTaskMinimiseAndVisibility_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskMinimise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.True(t, info.Minimized)
	assert.False(t, info.Visible)

	_, handled, err = c.PERFORM(TaskSetVisibility{Name: "test", Visible: true})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ = c.QUERY(QueryWindowByName{Name: "test"})
	info = result.(*WindowInfo)
	assert.True(t, info.Visible)
}

func TestTaskSetAlwaysOnTop_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetAlwaysOnTop{Name: "test", AlwaysOnTop: true})
	require.NoError(t, err)
	assert.True(t, handled)

	svc := core.MustServiceFor[*Service](c, "window")
	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	assert.True(t, pw.(*mockWindow).alwaysOnTop)
}

func TestTaskSetBackgroundColour_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetBackgroundColour{
		Name: "test", Red: 10, Green: 20, Blue: 30, Alpha: 40,
	})
	require.NoError(t, err)
	assert.True(t, handled)

	svc := core.MustServiceFor[*Service](c, "window")
	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	assert.Equal(t, [4]uint8{10, 20, 30, 40}, pw.(*mockWindow).backgroundColor)
}

func TestTaskTileWindows_UsesPrimaryScreenSize(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("left")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("right")}})

	_, handled, err := c.PERFORM(TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	require.NoError(t, err)
	assert.True(t, handled)

	left, _, _ := c.QUERY(QueryWindowByName{Name: "left"})
	right, _, _ := c.QUERY(QueryWindowByName{Name: "right"})
	leftInfo := left.(*WindowInfo)
	rightInfo := right.(*WindowInfo)

	assert.Equal(t, 1280, leftInfo.Width)
	assert.Equal(t, 1280, rightInfo.Width)
	assert.Equal(t, 0, leftInfo.X)
	assert.Equal(t, 1280, rightInfo.X)
}

func TestTaskTileWindows_ResetsMaximizedState(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("left")}})

	_, _, _ = c.PERFORM(TaskMaximise{Name: "left"})
	_, handled, err := c.PERFORM(TaskTileWindows{Mode: "left-half", Windows: []string{"left"}})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "left"})
	info := result.(*WindowInfo)
	assert.False(t, info.Maximized)
	assert.Equal(t, 0, info.X)
	assert.Equal(t, 1280, info.Width)
}

func TestTaskSetOpacity_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetOpacity{Name: "test", Opacity: 0.65})
	require.NoError(t, err)
	assert.True(t, handled)

	svc := core.MustServiceFor[*Service](c, "window")
	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	assert.InDelta(t, 0.65, pw.(*mockWindow).opacity, 0.0001)
}

func TestTaskSetOpacity_BadRange(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetOpacity{Name: "test", Opacity: 1.5})
	require.Error(t, err)
	assert.True(t, handled)
}

func TestTaskStackWindows_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("one")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("two")}})

	_, handled, err := c.PERFORM(TaskStackWindows{
		Windows: []string{"one", "two"},
		OffsetX: 20,
		OffsetY: 30,
	})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "two"})
	info := result.(*WindowInfo)
	assert.Equal(t, 20, info.X)
	assert.Equal(t, 30, info.Y)
}

func TestTaskApplyWorkflow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("editor")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("assistant")}})

	_, handled, err := c.PERFORM(TaskApplyWorkflow{
		Workflow: WorkflowCoding,
		Windows:  []string{"editor", "assistant"},
	})
	require.NoError(t, err)
	assert.True(t, handled)

	editorResult, _, _ := c.QUERY(QueryWindowByName{Name: "editor"})
	assistantResult, _, _ := c.QUERY(QueryWindowByName{Name: "assistant"})
	editor := editorResult.(*WindowInfo)
	assistant := assistantResult.(*WindowInfo)
	assert.Greater(t, editor.Width, assistant.Width)
	assert.Equal(t, editor.Width, assistant.X)
}

func TestTaskRestoreLayout_ClearsMaximizedState(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("editor")}})
	_, _, _ = c.PERFORM(TaskMaximise{Name: "editor"})

	svc := core.MustServiceFor[*Service](c, "window")
	err := svc.Manager().Layout().SaveLayout("restore", map[string]WindowState{
		"editor": {X: 12, Y: 34, Width: 640, Height: 480, Maximized: false},
	})
	require.NoError(t, err)

	_, handled, err := c.PERFORM(TaskRestoreLayout{Name: "restore"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "editor"})
	info := result.(*WindowInfo)
	assert.False(t, info.Maximized)
	assert.Equal(t, 12, info.X)
	assert.Equal(t, 640, info.Width)
}

func TestTaskMaximise_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	info := requireOpenWindow(t, c, Window{Name: "drop-test"})
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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_ = requireOpenWindow(t, c, Window{Name: "editor", Width: 960, Height: 1080, X: 0, Y: 0})
	_ = requireOpenWindow(t, c, Window{Name: "terminal", Width: 960, Height: 1080, X: 960, Y: 0})

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
	_ = requireOpenWindow(t, c, Window{Name: "editor", Width: 800, Height: 600, X: 0, Y: 0})
	_ = requireOpenWindow(t, c, Window{Name: "terminal", Width: 800, Height: 600, X: 0, Y: 0})

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
	_ = requireOpenWindow(t, c, Window{Name: "s1", Width: 800, Height: 600})
	_ = requireOpenWindow(t, c, Window{Name: "s2", Width: 800, Height: 600})

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
	_ = requireOpenWindow(t, c, Window{Name: "editor", Width: 800, Height: 600})
	_ = requireOpenWindow(t, c, Window{Name: "terminal", Width: 800, Height: 600})

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

// --- TaskSetZoom / QueryWindowZoom ---

func TestTaskSetZoom_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskSetZoom{Name: "test", Factor: 1.5})
	require.NoError(t, err)
	assert.True(t, handled)

	result, handled, err := c.QUERY(QueryWindowZoom{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.InDelta(t, 1.5, result.(float64), 0.001)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.InDelta(t, 1.5, mw.zoom, 0.001)
}

func TestTaskSetZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetZoom{Name: "nonexistent", Factor: 1.2})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestQueryWindowZoom_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.QUERY(QueryWindowZoom{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- QueryWindowBounds ---

func TestQueryWindowBounds_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "bounds-test", Width: 800, Height: 600, X: 50, Y: 75})

	result, handled, err := c.QUERY(QueryWindowBounds{Name: "bounds-test"})
	require.NoError(t, err)
	assert.True(t, handled)
	bounds := result.(*Bounds)
	require.NotNil(t, bounds)
	assert.Equal(t, 50, bounds.X)
	assert.Equal(t, 75, bounds.Y)
	assert.Equal(t, 800, bounds.Width)
	assert.Equal(t, 600, bounds.Height)
}

func TestQueryWindowBounds_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.QUERY(QueryWindowBounds{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetURL ---

func TestTaskSetURL_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskSetURL{Name: "test", URL: "/settings"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "/settings", mw.url)
}

func TestTaskSetURL_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskSetURL{Name: "nonexistent", URL: "/nope"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetHTML ---

func TestTaskSetHTML_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

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
	_, handled, err := c.PERFORM(TaskSetHTML{Name: "nonexistent", HTML: "<p>nope</p>"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskExecJS ---

func TestTaskExecJS_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskExecJS{Name: "test", JS: "document.title = 'Updated'"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, "document.title = 'Updated'", mw.lastJS)
}

func TestTaskExecJS_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskExecJS{Name: "nonexistent", JS: "alert(1)"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskToggleFullscreen ---

func TestTaskToggleFullscreen_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskToggleFullscreen{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.Equal(t, 1, mw.toggleFullscreenCount)
}

func TestTaskToggleFullscreen_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskToggleFullscreen{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskPrint ---

func TestTaskPrint_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskPrint{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.printCalled)
}

func TestTaskPrint_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskPrint{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskFlash ---

func TestTaskFlash_Good(t *testing.T) {
	svc, c := newTestWindowService(t)
	_ = requireOpenWindow(t, c, Window{Name: "test"})

	_, handled, err := c.PERFORM(TaskFlash{Name: "test", Enabled: true})
	require.NoError(t, err)
	assert.True(t, handled)

	pw, ok := svc.Manager().Get("test")
	require.True(t, ok)
	mw := pw.(*mockWindow)
	assert.True(t, mw.flashing)

	_, _, _ = c.PERFORM(TaskFlash{Name: "test", Enabled: false})
	assert.False(t, mw.flashing)
}

func TestTaskFlash_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskFlash{Name: "nonexistent", Enabled: true})
	assert.True(t, handled)
	assert.Error(t, err)
}
