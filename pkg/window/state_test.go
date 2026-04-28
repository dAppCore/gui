package window

import (
	core "dappco.re/go"
	"os"
	"path/filepath"
	"time"
)

func TestStateManagerState_NewStateManagerWithDir_Good(t *core.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)

	core.AssertNotNil(t, sm)
	core.AssertEqual(t, dir, sm.dataDir())
	core.AssertEqual(t, filepath.Join(dir, "window_state.json"), sm.filePath())
	core.AssertEmpty(t, sm.ListStates())
}

func TestStateManagerState_NewStateManagerWithPathEnv_Good(t *core.T) {
	path := filepath.Join(t.TempDir(), "custom", "window_state.json")
	t.Setenv(windowStateFileEnv, path)

	sm := NewStateManager()

	core.AssertNotNil(t, sm)
	core.AssertEqual(t, path, sm.filePath())
	core.AssertEqual(t, filepath.Dir(path), sm.dataDir())
}

func TestStateManagerState_NewStateManagerWithDir_Bad(t *core.T) {
	sm := NewStateManagerWithDir("")

	core.AssertNotNil(t, sm)
	core.AssertEmpty(t, sm.dataDir())
}

func TestStateManagerState_NewStateManagerWithDir_InvalidFile_Good(t *core.T) {
	dir := t.TempDir()
	core.RequireNoError(t, os.WriteFile(filepath.Join(dir, "window_state.json"), []byte("{invalid"), 0o644))

	sm := NewStateManagerWithDir(dir)

	core.AssertNotNil(t, sm)
	core.AssertEmpty(t, sm.ListStates())
}

func TestStateManagerState_SetPath_Good(t *core.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)
	path := filepath.Join(dir, "custom", "window-state.json")

	sm.SetPath(path)
	sm.SetState("main", WindowState{X: 10, Y: 20, Width: 300, Height: 200})
	sm.ForceSync()

	content, err := os.ReadFile(path)
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), `"main"`)
	core.AssertEqual(t, path, sm.filePath())
	core.AssertEqual(t, filepath.Dir(path), sm.dataDir())
}

func TestStateManagerState_SetPath_Ugly(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	initial := sm.filePath()

	sm.SetPath("")

	core.AssertEqual(t, initial, sm.filePath())
}

func TestStateManagerState_SetState_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("main", WindowState{X: 1, Y: 2, Width: 3, Height: 4, Maximized: true})

	got, ok := sm.GetState("main")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 1, got.X)
	core.AssertEqual(t, 2, got.Y)
	core.AssertEqual(t, 3, got.Width)
	core.AssertEqual(t, 4, got.Height)
	core.AssertTrue(t, got.Maximized)
	core.AssertNotEmpty(t, got.UpdatedAt)
}

func TestStateManagerState_UpdatePosition_Bad(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdatePosition("missing", 30, 40)

	got, ok := sm.GetState("missing")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 30, got.X)
	core.AssertEqual(t, 40, got.Y)
	core.AssertEmpty(t, got.Width)
	core.AssertEmpty(t, got.Height)
}

func TestStateManagerState_UpdateSize_Ugly(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdateSize("missing", -800, -600)

	got, ok := sm.GetState("missing")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, -800, got.Width)
	core.AssertEqual(t, -600, got.Height)
}

func TestStateManagerState_UpdateMaximized_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.UpdateMaximized("main", true)

	got, ok := sm.GetState("main")
	core.RequireTrue(t, ok)
	core.AssertTrue(t, got.Maximized)
}

func TestStateManagerState_CaptureState_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.CaptureState(&mockWindow{name: "captured", x: 50, y: 60, width: 800, height: 600, maximised: true})

	got, ok := sm.GetState("captured")
	core.RequireTrue(t, ok)
	core.AssertEqual(t, 50, got.X)
	core.AssertEqual(t, 60, got.Y)
	core.AssertEqual(t, 800, got.Width)
	core.AssertEqual(t, 600, got.Height)
	core.AssertTrue(t, got.Maximized)
}

func TestStateManagerState_ApplyState_Bad(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	w := &Window{Name: "missing", X: 9, Y: 8, Width: 7, Height: 6}

	sm.ApplyState(w)

	core.AssertEqual(t, 9, w.X)
	core.AssertEqual(t, 8, w.Y)
	core.AssertEqual(t, 7, w.Width)
	core.AssertEqual(t, 6, w.Height)
}

func TestStateManagerState_ApplyState_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("main", WindowState{X: 11, Y: 12, Width: 1300, Height: 900})

	w := &Window{Name: "main", X: 1, Y: 2, Width: 10, Height: 20}
	sm.ApplyState(w)

	core.AssertEqual(t, 11, w.X)
	core.AssertEqual(t, 12, w.Y)
	core.AssertEqual(t, 1300, w.Width)
	core.AssertEqual(t, 900, w.Height)
}

func TestStateManagerState_ListStates_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{})
	sm.SetState("beta", WindowState{})

	names := sm.ListStates()

	core.AssertElementsMatch(t, []string{"alpha", "beta"}, names)
}

func TestStateManagerState_Clear_Good(t *core.T) {
	sm := NewStateManagerWithDir(t.TempDir())
	sm.SetState("alpha", WindowState{})
	sm.SetState("beta", WindowState{})

	sm.Clear()

	core.AssertEmpty(t, sm.ListStates())
}

func TestStateManagerState_ForceSync_Good(t *core.T) {
	dir := t.TempDir()
	sm := NewStateManagerWithDir(dir)
	sm.SetState("main", WindowState{Width: 800, Height: 600})
	time.Sleep(10 * time.Millisecond)

	sm.ForceSync()

	content, err := os.ReadFile(filepath.Join(dir, "window_state.json"))
	core.RequireNoError(t, err)
	core.AssertContains(t, string(content), `"main"`)
}
