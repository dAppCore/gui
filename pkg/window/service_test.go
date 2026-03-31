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
		Options: []WindowOption{WithName("test"), WithURL("/")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test", info.Name)
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
}

func TestTaskRestoreLayout_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskRestoreLayout{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}
