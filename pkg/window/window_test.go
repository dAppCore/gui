// pkg/window/window_test.go
package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowDefaults(t *testing.T) {
	w := &Window{}
	assert.Equal(t, "", w.Name)
	assert.Equal(t, 0, w.Width)
}

func TestWindowOption_Name_Good(t *testing.T) {
	w := &Window{}
	err := WithName("main")(w)
	require.NoError(t, err)
	assert.Equal(t, "main", w.Name)
}

func TestWindowOption_Title_Good(t *testing.T) {
	w := &Window{}
	err := WithTitle("My App")(w)
	require.NoError(t, err)
	assert.Equal(t, "My App", w.Title)
}

func TestWindowOption_URL_Good(t *testing.T) {
	w := &Window{}
	err := WithURL("/dashboard")(w)
	require.NoError(t, err)
	assert.Equal(t, "/dashboard", w.URL)
}

func TestWindowOption_Size_Good(t *testing.T) {
	w := &Window{}
	err := WithSize(1280, 720)(w)
	require.NoError(t, err)
	assert.Equal(t, 1280, w.Width)
	assert.Equal(t, 720, w.Height)
}

func TestWindowOption_Position_Good(t *testing.T) {
	w := &Window{}
	err := WithPosition(100, 200)(w)
	require.NoError(t, err)
	assert.Equal(t, 100, w.X)
	assert.Equal(t, 200, w.Y)
}

func TestApplyOptions_Good(t *testing.T) {
	w, err := ApplyOptions(
		WithName("test"),
		WithTitle("Test Window"),
		WithURL("/test"),
		WithSize(800, 600),
	)
	require.NoError(t, err)
	assert.Equal(t, "test", w.Name)
	assert.Equal(t, "Test Window", w.Title)
	assert.Equal(t, "/test", w.URL)
	assert.Equal(t, 800, w.Width)
	assert.Equal(t, 600, w.Height)
}

func TestApplyOptions_Bad(t *testing.T) {
	_, err := ApplyOptions(func(w *Window) error {
		return assert.AnError
	})
	assert.Error(t, err)
}

func TestApplyOptions_Empty_Good(t *testing.T) {
	w, err := ApplyOptions()
	require.NoError(t, err)
	assert.NotNil(t, w)
}

// newTestManager creates a Manager with a mock platform and clean state for testing.
func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	m := &Manager{
		platform: p,
		state:    &StateManager{states: make(map[string]WindowState)},
		layout:   &LayoutManager{layouts: make(map[string]Layout)},
		windows:  make(map[string]PlatformWindow),
	}
	return m, p
}

func TestManager_Open_Good(t *testing.T) {
	m, p := newTestManager()
	pw, err := m.Open(WithName("test"), WithTitle("Test"), WithURL("/test"), WithSize(800, 600))
	require.NoError(t, err)
	assert.NotNil(t, pw)
	assert.Equal(t, "test", pw.Name())
	assert.Len(t, p.windows, 1)
}

func TestManager_Open_Defaults_Good(t *testing.T) {
	m, _ := newTestManager()
	pw, err := m.Open()
	require.NoError(t, err)
	assert.Equal(t, "main", pw.Name())
	w, h := pw.Size()
	assert.Equal(t, 1280, w)
	assert.Equal(t, 800, h)
}

func TestManager_Open_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, err := m.Open(func(w *Window) error { return assert.AnError })
	assert.Error(t, err)
}

func TestManager_Get_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("findme"))
	pw, ok := m.Get("findme")
	assert.True(t, ok)
	assert.Equal(t, "findme", pw.Name())
}

func TestManager_Get_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, ok := m.Get("nonexistent")
	assert.False(t, ok)
}

func TestManager_List_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"))
	_, _ = m.Open(WithName("b"))
	names := m.List()
	assert.Len(t, names, 2)
	assert.Contains(t, names, "a")
	assert.Contains(t, names, "b")
}

func TestManager_Remove_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("temp"))
	m.Remove("temp")
	_, ok := m.Get("temp")
	assert.False(t, ok)
}

// --- StateManager Tests ---

// newTestStateManager creates a clean StateManager with a temp dir for testing.
func newTestStateManager(t *testing.T) *StateManager {
	return &StateManager{
		configDir: t.TempDir(),
		states:    make(map[string]WindowState),
	}
}

func TestStateManager_SetGet_Good(t *testing.T) {
	sm := newTestStateManager(t)
	state := WindowState{X: 100, Y: 200, Width: 800, Height: 600, Maximized: false}
	sm.SetState("main", state)
	got, ok := sm.GetState("main")
	assert.True(t, ok)
	assert.Equal(t, 100, got.X)
	assert.Equal(t, 800, got.Width)
}

func TestStateManager_SetGet_Bad(t *testing.T) {
	sm := newTestStateManager(t)
	_, ok := sm.GetState("nonexistent")
	assert.False(t, ok)
}

func TestStateManager_CaptureState_Good(t *testing.T) {
	sm := newTestStateManager(t)
	w := &mockWindow{name: "cap", x: 50, y: 60, width: 1024, height: 768, maximised: true}
	sm.CaptureState(w)
	got, ok := sm.GetState("cap")
	assert.True(t, ok)
	assert.Equal(t, 50, got.X)
	assert.Equal(t, 1024, got.Width)
	assert.True(t, got.Maximized)
}

func TestStateManager_ApplyState_Good(t *testing.T) {
	sm := newTestStateManager(t)
	sm.SetState("win", WindowState{X: 10, Y: 20, Width: 640, Height: 480})
	w := &Window{Name: "win", Width: 1280, Height: 800}
	sm.ApplyState(w)
	assert.Equal(t, 10, w.X)
	assert.Equal(t, 20, w.Y)
	assert.Equal(t, 640, w.Width)
	assert.Equal(t, 480, w.Height)
}

func TestStateManager_ListStates_Good(t *testing.T) {
	sm := newTestStateManager(t)
	sm.SetState("a", WindowState{Width: 100})
	sm.SetState("b", WindowState{Width: 200})
	names := sm.ListStates()
	assert.Len(t, names, 2)
}

func TestStateManager_Clear_Good(t *testing.T) {
	sm := newTestStateManager(t)
	sm.SetState("a", WindowState{Width: 100})
	sm.Clear()
	names := sm.ListStates()
	assert.Empty(t, names)
}

func TestStateManager_Persistence_Good(t *testing.T) {
	dir := t.TempDir()
	sm1 := &StateManager{configDir: dir, states: make(map[string]WindowState)}
	sm1.SetState("persist", WindowState{X: 42, Y: 84, Width: 500, Height: 300})
	sm1.ForceSync()

	sm2 := &StateManager{configDir: dir, states: make(map[string]WindowState)}
	sm2.load()
	got, ok := sm2.GetState("persist")
	assert.True(t, ok)
	assert.Equal(t, 42, got.X)
	assert.Equal(t, 500, got.Width)
}

// --- LayoutManager Tests ---

// newTestLayoutManager creates a clean LayoutManager with a temp dir for testing.
func newTestLayoutManager(t *testing.T) *LayoutManager {
	return &LayoutManager{
		configDir: t.TempDir(),
		layouts:   make(map[string]Layout),
	}
}

func TestLayoutManager_SaveGet_Good(t *testing.T) {
	lm := newTestLayoutManager(t)
	states := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 1080},
	}
	err := lm.SaveLayout("coding", states)
	require.NoError(t, err)

	layout, ok := lm.GetLayout("coding")
	assert.True(t, ok)
	assert.Equal(t, "coding", layout.Name)
	assert.Len(t, layout.Windows, 2)
}

func TestLayoutManager_GetLayout_Bad(t *testing.T) {
	lm := newTestLayoutManager(t)
	_, ok := lm.GetLayout("nonexistent")
	assert.False(t, ok)
}

func TestLayoutManager_ListLayouts_Good(t *testing.T) {
	lm := newTestLayoutManager(t)
	_ = lm.SaveLayout("a", map[string]WindowState{})
	_ = lm.SaveLayout("b", map[string]WindowState{})
	layouts := lm.ListLayouts()
	assert.Len(t, layouts, 2)
}

func TestLayoutManager_DeleteLayout_Good(t *testing.T) {
	lm := newTestLayoutManager(t)
	_ = lm.SaveLayout("temp", map[string]WindowState{})
	lm.DeleteLayout("temp")
	_, ok := lm.GetLayout("temp")
	assert.False(t, ok)
}

// --- Tiling Tests ---

func TestTileMode_String_Good(t *testing.T) {
	assert.Equal(t, "left-half", TileModeLeftHalf.String())
	assert.Equal(t, "grid", TileModeGrid.String())
}

func TestManager_TileWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("a"), WithSize(800, 600))
	_, _ = m.Open(WithName("b"), WithSize(800, 600))
	err := m.TileWindows(TileModeLeftRight, []string{"a", "b"}, 1920, 1080)
	require.NoError(t, err)
	a, _ := m.Get("a")
	b, _ := m.Get("b")
	aw, _ := a.Size()
	bw, _ := b.Size()
	assert.Equal(t, 960, aw)
	assert.Equal(t, 960, bw)
}

func TestManager_TileWindows_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.TileWindows(TileModeLeftRight, []string{"nonexistent"}, 1920, 1080)
	assert.Error(t, err)
}

func TestManager_SnapWindow_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("snap"), WithSize(800, 600))
	err := m.SnapWindow("snap", SnapLeft, 1920, 1080)
	require.NoError(t, err)
	w, _ := m.Get("snap")
	x, _ := w.Position()
	assert.Equal(t, 0, x)
	sw, _ := w.Size()
	assert.Equal(t, 960, sw)
}

func TestManager_StackWindows_Good(t *testing.T) {
	m, _ := newTestManager()
	_, _ = m.Open(WithName("s1"), WithSize(800, 600))
	_, _ = m.Open(WithName("s2"), WithSize(800, 600))
	err := m.StackWindows([]string{"s1", "s2"}, 30, 30)
	require.NoError(t, err)
	s2, _ := m.Get("s2")
	x, y := s2.Position()
	assert.Equal(t, 30, x)
	assert.Equal(t, 30, y)
}

func TestWorkflowLayout_Good(t *testing.T) {
	assert.Equal(t, "coding", WorkflowCoding.String())
	assert.Equal(t, "debugging", WorkflowDebugging.String())
}
