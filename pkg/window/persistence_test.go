// pkg/window/persistence_test.go
package window

import (
	"os"
	core "dappco.re/go/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- StateManager Persistence Tests ---

func TestStateManager_SetAndGet_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	state := WindowState{
		X: 150, Y: 250, Width: 1024, Height: 768,
		Maximized: true, Screen: "primary", URL: "/app",
	}
	sm.SetState("editor", state)

	got, ok := sm.GetState("editor")
	require.True(t, ok)
	assert.Equal(t, 150, got.X)
	assert.Equal(t, 250, got.Y)
	assert.Equal(t, 1024, got.Width)
	assert.Equal(t, 768, got.Height)
	assert.True(t, got.Maximized)
	assert.Equal(t, "primary", got.Screen)
	assert.Equal(t, "/app", got.URL)
	assert.NotZero(t, got.UpdatedAt, "UpdatedAt should be set by SetState")
}

func TestStateManager_UpdatePosition_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{X: 0, Y: 0, Width: 800, Height: 600})

	sm.UpdatePosition("win", 300, 400)

	got, ok := sm.GetState("win")
	require.True(t, ok)
	assert.Equal(t, 300, got.X)
	assert.Equal(t, 400, got.Y)
	// Width/Height should remain unchanged
	assert.Equal(t, 800, got.Width)
	assert.Equal(t, 600, got.Height)
}

func TestStateManager_UpdateSize_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{X: 100, Y: 200, Width: 800, Height: 600})

	sm.UpdateSize("win", 1920, 1080)

	got, ok := sm.GetState("win")
	require.True(t, ok)
	assert.Equal(t, 1920, got.Width)
	assert.Equal(t, 1080, got.Height)
	// Position should remain unchanged
	assert.Equal(t, 100, got.X)
	assert.Equal(t, 200, got.Y)
}

func TestStateManager_UpdateMaximized_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("win", WindowState{Width: 800, Height: 600, Maximized: false})

	sm.UpdateMaximized("win", true)

	got, ok := sm.GetState("win")
	require.True(t, ok)
	assert.True(t, got.Maximized)

	sm.UpdateMaximized("win", false)

	got, ok = sm.GetState("win")
	require.True(t, ok)
	assert.False(t, got.Maximized)
}

func TestStateManager_CaptureState_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	pw := &mockWindow{
		name: "captured", x: 75, y: 125,
		width: 1440, height: 900, maximised: true,
	}

	sm.CaptureState(pw)

	got, ok := sm.GetState("captured")
	require.True(t, ok)
	assert.Equal(t, 75, got.X)
	assert.Equal(t, 125, got.Y)
	assert.Equal(t, 1440, got.Width)
	assert.Equal(t, 900, got.Height)
	assert.True(t, got.Maximized)
	assert.NotZero(t, got.UpdatedAt)
}

func TestStateManager_ApplyState_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("target", WindowState{X: 55, Y: 65, Width: 700, Height: 500})

	w := &Window{Name: "target", Width: 1280, Height: 800, X: 0, Y: 0}
	sm.ApplyState(w)

	assert.Equal(t, 55, w.X)
	assert.Equal(t, 65, w.Y)
	assert.Equal(t, 700, w.Width)
	assert.Equal(t, 500, w.Height)
}

func TestStateManager_ApplyState_Good_NoState(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())

	w := &Window{Name: "untouched", Width: 1280, Height: 800, X: 10, Y: 20}
	sm.ApplyState(w)

	// Window should remain unchanged when no state is saved
	assert.Equal(t, 10, w.X)
	assert.Equal(t, 20, w.Y)
	assert.Equal(t, 1280, w.Width)
	assert.Equal(t, 800, w.Height)
}

func TestStateManager_ListStates_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{Width: 100})
	sm.SetState("beta", WindowState{Width: 200})
	sm.SetState("gamma", WindowState{Width: 300})

	names := sm.ListStates()
	assert.Len(t, names, 3)
	assert.Contains(t, names, "alpha")
	assert.Contains(t, names, "beta")
	assert.Contains(t, names, "gamma")
}

func TestStateManager_Clear_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("a", WindowState{Width: 100})
	sm.SetState("b", WindowState{Width: 200})
	sm.SetState("c", WindowState{Width: 300})

	sm.Clear()

	names := sm.ListStates()
	assert.Empty(t, names)

	_, ok := sm.GetState("a")
	assert.False(t, ok)
}

func TestStateManager_Persistence_Good(t *testing.T) {
	dir := t.TempDir()

	// First manager: write state and force sync to disk
	sm1 := NewStateManagerWithDir(dir)
	sm1.SetState("persist-win", WindowState{
		X: 42, Y: 84, Width: 500, Height: 300,
		Maximized: true, Screen: "secondary", URL: "/settings",
	})
	sm1.ForceSync()

	// Second manager: load from the same directory
	sm2 := NewStateManagerWithDir(dir)

	got, ok := sm2.GetState("persist-win")
	require.True(t, ok)
	assert.Equal(t, 42, got.X)
	assert.Equal(t, 84, got.Y)
	assert.Equal(t, 500, got.Width)
	assert.Equal(t, 300, got.Height)
	assert.True(t, got.Maximized)
	assert.Equal(t, "secondary", got.Screen)
	assert.Equal(t, "/settings", got.URL)
	assert.NotZero(t, got.UpdatedAt)
}

func TestStateManager_SetPath_Good(t *testing.T) {
	dir := t.TempDir()
	path := core.JoinPath(dir, "custom", "window-state.json")

	sm := NewStateManagerWithDir(dir)
	sm.SetPath(path)
	sm.SetState("custom", WindowState{Width: 640, Height: 480})
	sm.ForceSync()

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), "custom")
}

// --- LayoutManager Persistence Tests ---

func TestLayoutManager_SaveAndGet_Good(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	windows := map[string]WindowState{
		"editor":   {X: 0, Y: 0, Width: 960, Height: 1080},
		"terminal": {X: 960, Y: 0, Width: 960, Height: 540},
		"browser":  {X: 960, Y: 540, Width: 960, Height: 540},
	}

	err := lm.SaveLayout("coding", windows)
	require.NoError(t, err)

	layout, ok := lm.GetLayout("coding")
	require.True(t, ok)
	assert.Equal(t, "coding", layout.Name)
	assert.Len(t, layout.Windows, 3)
	assert.Equal(t, 960, layout.Windows["editor"].Width)
	assert.Equal(t, 1080, layout.Windows["editor"].Height)
	assert.Equal(t, 960, layout.Windows["terminal"].X)
	assert.NotZero(t, layout.CreatedAt)
	assert.NotZero(t, layout.UpdatedAt)
	assert.Equal(t, layout.CreatedAt, layout.UpdatedAt, "CreatedAt and UpdatedAt should match on first save")
}

func TestLayoutManager_SaveLayout_EmptyName_Bad(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	err := lm.SaveLayout("", map[string]WindowState{
		"win": {Width: 800},
	})
	assert.Error(t, err)
}

func TestLayoutManager_SaveLayout_Update_Good(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())

	// First save
	err := lm.SaveLayout("evolving", map[string]WindowState{
		"win1": {Width: 800, Height: 600},
	})
	require.NoError(t, err)

	first, ok := lm.GetLayout("evolving")
	require.True(t, ok)
	originalCreatedAt := first.CreatedAt
	originalUpdatedAt := first.UpdatedAt

	// Small delay to ensure UpdatedAt differs
	time.Sleep(2 * time.Millisecond)

	// Second save with same name but different windows
	err = lm.SaveLayout("evolving", map[string]WindowState{
		"win1": {Width: 1024, Height: 768},
		"win2": {Width: 640, Height: 480},
	})
	require.NoError(t, err)

	updated, ok := lm.GetLayout("evolving")
	require.True(t, ok)

	// CreatedAt should be preserved from the original save
	assert.Equal(t, originalCreatedAt, updated.CreatedAt, "CreatedAt should be preserved on update")
	// UpdatedAt should be newer
	assert.GreaterOrEqual(t, updated.UpdatedAt, originalUpdatedAt, "UpdatedAt should advance on update")
	// Windows should reflect the second save
	assert.Len(t, updated.Windows, 2)
	assert.Equal(t, 1024, updated.Windows["win1"].Width)
}

func TestLayoutManager_ListLayouts_Good(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	require.NoError(t, lm.SaveLayout("coding", map[string]WindowState{
		"editor": {Width: 960}, "terminal": {Width: 960},
	}))
	require.NoError(t, lm.SaveLayout("presenting", map[string]WindowState{
		"slides": {Width: 1920},
	}))
	require.NoError(t, lm.SaveLayout("debugging", map[string]WindowState{
		"code": {Width: 640}, "debugger": {Width: 640}, "console": {Width: 640},
	}))

	infos := lm.ListLayouts()
	assert.Len(t, infos, 3)

	// Build a lookup map for assertions regardless of order
	byName := make(map[string]LayoutInfo)
	for _, info := range infos {
		byName[info.Name] = info
	}

	assert.Equal(t, 2, byName["coding"].WindowCount)
	assert.Equal(t, 1, byName["presenting"].WindowCount)
	assert.Equal(t, 3, byName["debugging"].WindowCount)
}

func TestLayoutManager_DeleteLayout_Good(t *testing.T) {
	lm := NewLayoutManagerWithDir(t.TempDir())
	require.NoError(t, lm.SaveLayout("temporary", map[string]WindowState{
		"win": {Width: 800},
	}))

	// Verify it exists
	_, ok := lm.GetLayout("temporary")
	require.True(t, ok)

	lm.DeleteLayout("temporary")

	// Verify it is gone
	_, ok = lm.GetLayout("temporary")
	assert.False(t, ok)

	// Verify list is empty
	assert.Empty(t, lm.ListLayouts())
}

func TestLayoutManager_Persistence_Good(t *testing.T) {
	dir := t.TempDir()

	// First manager: save layout to disk
	lm1 := NewLayoutManagerWithDir(dir)
	err := lm1.SaveLayout("persisted", map[string]WindowState{
		"main":    {X: 0, Y: 0, Width: 1280, Height: 800},
		"sidebar": {X: 1280, Y: 0, Width: 640, Height: 800},
	})
	require.NoError(t, err)

	// Second manager: load from the same directory
	lm2 := NewLayoutManagerWithDir(dir)

	layout, ok := lm2.GetLayout("persisted")
	require.True(t, ok)
	assert.Equal(t, "persisted", layout.Name)
	assert.Len(t, layout.Windows, 2)
	assert.Equal(t, 1280, layout.Windows["main"].Width)
	assert.Equal(t, 800, layout.Windows["main"].Height)
	assert.Equal(t, 640, layout.Windows["sidebar"].Width)
	assert.NotZero(t, layout.CreatedAt)
	assert.NotZero(t, layout.UpdatedAt)
}
