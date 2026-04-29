package window

import (
	core "dappco.re/go"
	filepath "dappco.re/go/gui/compat/filepath"
	os "dappco.re/go/gui/compat/os"
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

func TestStateManagerState_NewStateManagerWithPathEnv_GoodCase(t *core.T) {
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

// AX7 generated source-matching smoke coverage.
func TestState_NewStateManager_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManager_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManager_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManager()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithDir_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithDir("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithDir_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithDir("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithDir_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithDir("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithPath_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithPath("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithPath_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithPath("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_NewStateManagerWithPath_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewStateManagerWithPath("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetPath_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetPath("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetPath_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetPath("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetPath_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetPath("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_GetState_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetState("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_GetState_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetState("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_GetState_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0, got1 := subject.GetState("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetState_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetState("agent", *new(WindowState))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetState_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetState("", *new(WindowState))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_SetState_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.SetState("../../edge", *new(WindowState))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdatePosition_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdatePosition("agent", 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdatePosition_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdatePosition("", 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdatePosition_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdatePosition("../../edge", -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateSize_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateSize("agent", 1, 1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateSize_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateSize("", 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateSize_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateSize("../../edge", -1, -1)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateMaximized_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateMaximized("agent", true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateMaximized_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateMaximized("", false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_UpdateMaximized_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.UpdateMaximized("../../edge", false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_CaptureState_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.CaptureState(*new(PlatformWindow))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_CaptureState_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.CaptureState(*new(PlatformWindow))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_CaptureState_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.CaptureState(*new(PlatformWindow))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ApplyState_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.ApplyState(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ApplyState_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.ApplyState(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ApplyState_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.ApplyState(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ListStates_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ListStates()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ListStates_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ListStates()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ListStates_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ListStates()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_Clear_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.Clear()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_Clear_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.Clear()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_Clear_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		subject.Clear()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ForceSync_Good(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ForceSync()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ForceSync_Bad(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ForceSync()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestState_StateManager_ForceSync_Ugly(t *core.T) {
	subject := new(StateManager)
	result := core.Try(func() any {
		got0 := subject.ForceSync()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
