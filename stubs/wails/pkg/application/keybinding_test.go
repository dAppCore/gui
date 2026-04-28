package application

import (
	core "dappco.re/go"
)

func TestKeyBindingManager_Add_Good(t *core.T) {
	manager := &KeyBindingManager{}
	calls := 0

	manager.Add("CmdOrCtrl+K", func(window Window) {
		calls++
		core.AssertNil(t, window)
	})

	handled := manager.Process("CmdOrCtrl+K", nil)

	core.AssertTrue(t, handled)
	core.AssertEqual(t, 1, calls)
	core.AssertLen(t, manager.GetAll(), 1)
}

func TestKeyBindingManager_Add_Bad(t *core.T) {
	manager := &KeyBindingManager{}

	handled := manager.Process("missing", nil)

	core.AssertFalse(t, handled)
	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Add_Ugly(t *core.T) {
	manager := &KeyBindingManager{}
	calls := 0

	manager.Add("CmdOrCtrl+K", func(Window) { calls++ })
	manager.Add("CmdOrCtrl+K", func(Window) { calls += 10 })

	handled := manager.Process("CmdOrCtrl+K", nil)

	core.AssertTrue(t, handled)
	core.AssertEqual(t, 10, calls)
}

func TestKeyBindingManager_Process_RecoversFromPanic(t *core.T) {
	manager := &KeyBindingManager{}

	manager.Add("CmdOrCtrl+K", func(Window) {
		panic("boom")
	})

	core.AssertFalse(t, manager.Process("CmdOrCtrl+K", nil))
}

func TestKeyBindingManager_Remove_Good(t *core.T) {
	manager := &KeyBindingManager{}
	manager.Add("CmdOrCtrl+K", func(Window) {})

	manager.Remove("CmdOrCtrl+K")

	core.AssertFalse(t, manager.Process("CmdOrCtrl+K", nil))
	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Bad(t *core.T) {
	manager := &KeyBindingManager{}

	manager.Remove("missing")

	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Ugly(t *core.T) {
	manager := &KeyBindingManager{}

	manager.Remove("")

	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_NilReceiver_IsSafe(t *core.T) {
	var manager *KeyBindingManager

	core.AssertNotPanics(t, func() {
		manager.Add("CmdOrCtrl+K", func(Window) {})
		manager.Remove("CmdOrCtrl+K")
		core.AssertFalse(t, manager.Process("CmdOrCtrl+K", nil))
		core.AssertNil(t, manager.GetAll())
	})
}
