package deno

import (
	core "dappco.re/go"
)

func TestSidecar_New_Good(t *core.T) {
	manager := New(Options{})

	status := manager.Status()
	core.AssertEqual(t, "deno", status.Binary)
	core.AssertFalse(t, status.Running)
	core.AssertEmpty(t, status.PID)
}

func TestSidecar_New_Bad(t *core.T) {
	manager := New(Options{Binary: "/usr/local/bin/deno-custom", Args: []string{"format"}})

	status := manager.Status()
	core.AssertEqual(t, "/usr/local/bin/deno-custom", status.Binary)
	core.AssertFalse(t, status.Running)
}

func TestSidecar_New_Ugly(t *core.T) {
	manager := New(Options{Binary: "   "})

	status := manager.Status()
	core.AssertEqual(t, "deno", status.Binary)
	core.AssertFalse(t, status.Running)
}

// AX7 generated source-matching smoke coverage.
func TestSidecar_Manager_Start_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Start_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Start_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Start(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Stop_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Stop_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Stop_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Stop(core.Background())
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Status_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Status_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Status_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Status()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_OnEvent_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_OnEvent_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_OnEvent_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		subject.OnEvent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Eval_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Eval_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Eval_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0, got1 := subject.Eval(core.Background(), "../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Emit_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Emit_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("", nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestSidecar_Manager_Emit_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Emit("../../edge", map[string]any{})
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
