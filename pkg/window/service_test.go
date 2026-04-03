package window

import (
	"context"
	"sync"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/screen"
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

type testScreenPlatform struct {
	screens []screen.Screen
}

func (p *testScreenPlatform) GetAll() []screen.Screen { return p.screens }

func (p *testScreenPlatform) GetPrimary() *screen.Screen {
	for i := range p.screens {
		if p.screens[i].IsPrimary {
			return &p.screens[i]
		}
	}
	return nil
}

func newTestWindowServiceWithScreen(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithService(screen.Register(&testScreenPlatform{
			screens: []screen.Screen{{
				ID: "primary", Name: "Primary", IsPrimary: true,
				Size:     screen.Size{Width: 2560, Height: 1440},
				Bounds:   screen.Rect{X: 0, Y: 0, Width: 2560, Height: 1440},
				WorkArea: screen.Rect{X: 0, Y: 0, Width: 2560, Height: 1440},
			}},
		})),
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
		Opts: []WindowOption{WithName("test"), WithURL("/")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestTaskOpenWindowDescriptor_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{
		Window: &Window{Name: "descriptor", Title: "Descriptor", Width: 640, Height: 480},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "descriptor", info.Name)
	assert.Equal(t, "Descriptor", info.Title)
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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("a")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("b")}})
	_, _, _ = c.PERFORM(TaskMinimise{Name: "b"})

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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetSize{Name: "test", W: 800, H: 600})
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
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

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
		Opts: []WindowOption{WithName("drop-test")},
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
