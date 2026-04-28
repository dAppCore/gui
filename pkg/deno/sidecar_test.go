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
	manager := New(Options{Binary: "/usr/local/bin/deno-custom", Args: []string{"fmt"}})

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
