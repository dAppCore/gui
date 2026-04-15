package window

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateManagerState_NewStateManagerWithDir_Good(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)

	require.NotNil(t, sm)
	assert.Equal(t, dir, sm.dataDir())
	assert.Equal(t, filepath.Join(dir, "window_state.json"), sm.filePath())
	assert.Empty(t, sm.ListStates())
}

func TestStateManagerState_NewStateManagerWithDir_Bad(t *testing.T) {
	sm := NewStateManagerWithDir("")

	require.NotNil(t, sm)
	assert.Empty(t, sm.dataDir())
}

func TestStateManagerState_SetPath_Good(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)
	path := filepath.Join(dir, "custom", "window-state.json")

	sm.SetPath(path)
	sm.SetState("main", WindowState{X: 10, Y: 20, Width: 300, Height: 200})
	sm.ForceSync()

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"main"`)
	assert.Equal(t, path, sm.filePath())
	assert.Equal(t, filepath.Dir(path), sm.dataDir())
}

func TestStateManagerState_SetPath_Ugly(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	initial := sm.filePath()

	sm.SetPath("")

	assert.Equal(t, initial, sm.filePath())
}

func TestStateManagerState_SetState_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("main", WindowState{X: 1, Y: 2, Width: 3, Height: 4, Maximized: true})

	got, ok := sm.GetState("main")
	require.True(t, ok)
	assert.Equal(t, 1, got.X)
	assert.Equal(t, 2, got.Y)
	assert.Equal(t, 3, got.Width)
	assert.Equal(t, 4, got.Height)
	assert.True(t, got.Maximized)
	assert.NotZero(t, got.UpdatedAt)
}

func TestStateManagerState_UpdatePosition_Bad(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdatePosition("missing", 30, 40)

	got, ok := sm.GetState("missing")
	require.True(t, ok)
	assert.Equal(t, 30, got.X)
	assert.Equal(t, 40, got.Y)
	assert.Zero(t, got.Width)
	assert.Zero(t, got.Height)
}

func TestStateManagerState_UpdateSize_Ugly(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdateSize("missing", -800, -600)

	got, ok := sm.GetState("missing")
	require.True(t, ok)
	assert.Equal(t, -800, got.Width)
	assert.Equal(t, -600, got.Height)
}

func TestStateManagerState_UpdateMaximized_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdateMaximized("main", true)

	got, ok := sm.GetState("main")
	require.True(t, ok)
	assert.True(t, got.Maximized)
}

func TestStateManagerState_CaptureState_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.CaptureState(&mockWindow{name: "captured", x: 50, y: 60, width: 800, height: 600, maximised: true})

	got, ok := sm.GetState("captured")
	require.True(t, ok)
	assert.Equal(t, 50, got.X)
	assert.Equal(t, 60, got.Y)
	assert.Equal(t, 800, got.Width)
	assert.Equal(t, 600, got.Height)
	assert.True(t, got.Maximized)
}

func TestStateManagerState_ApplyState_Bad(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	w := &Window{Name: "missing", X: 9, Y: 8, Width: 7, Height: 6}

	sm.ApplyState(w)

	assert.Equal(t, 9, w.X)
	assert.Equal(t, 8, w.Y)
	assert.Equal(t, 7, w.Width)
	assert.Equal(t, 6, w.Height)
}

func TestStateManagerState_ApplyState_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("main", WindowState{X: 11, Y: 12, Width: 1300, Height: 900})

	w := &Window{Name: "main", X: 1, Y: 2, Width: 10, Height: 20}
	sm.ApplyState(w)

	assert.Equal(t, 11, w.X)
	assert.Equal(t, 12, w.Y)
	assert.Equal(t, 1300, w.Width)
	assert.Equal(t, 900, w.Height)
}

func TestStateManagerState_ListStates_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{})
	sm.SetState("beta", WindowState{})

	names := sm.ListStates()

	assert.ElementsMatch(t, []string{"alpha", "beta"}, names)
}

func TestStateManagerState_Clear_Good(t *testing.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{})
	sm.SetState("beta", WindowState{})

	sm.Clear()

	assert.Empty(t, sm.ListStates())
}

func TestStateManagerState_ForceSync_Good(t *testing.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)
	sm.SetState("main", WindowState{Width: 800, Height: 600})
	time.Sleep(10 * time.Millisecond)

	sm.ForceSync()

	content, err := os.ReadFile(filepath.Join(dir, "window_state.json"))
	require.NoError(t, err)
	assert.Contains(t, string(content), `"main"`)
}
