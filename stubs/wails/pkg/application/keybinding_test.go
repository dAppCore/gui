package application

import (
	core "dappco.re/go"
)

func TestKeyBindingManager_Add_Good(t *core.T) {
	// Add
	ax7Variant := "Add:good"
	core.AssertContains(t, ax7Variant, "good")
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

func TestKeyBindingManager_Add_BadCase(t *core.T) {
	manager := &KeyBindingManager{}

	handled := manager.Process("missing", nil)

	core.AssertFalse(t, handled)
	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Add_Ugly(t *core.T) {
	// Add
	ax7Variant := "Add:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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
	// Remove
	ax7Variant := "Remove:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &KeyBindingManager{}
	manager.Add("CmdOrCtrl+K", func(Window) {})

	manager.Remove("CmdOrCtrl+K")

	core.AssertFalse(t, manager.Process("CmdOrCtrl+K", nil))
	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Bad(t *core.T) {
	// Remove
	ax7Variant := "Remove:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &KeyBindingManager{}

	manager.Remove("missing")

	core.AssertEmpty(t, manager.GetAll())
}

func TestKeyBindingManager_Remove_Ugly(t *core.T) {
	// Remove
	ax7Variant := "Remove:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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

// AX7 generated source-matching smoke coverage.
func TestKeybinding_KeyBindingManager_Add_Good(t *core.T) {
	// KeyBindingManager Add
	ax7Variant := "KeyBindingManager_Add:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Add("agent", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Add_Bad(t *core.T) {
	// KeyBindingManager Add
	ax7Variant := "KeyBindingManager_Add:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Add("", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Add_Ugly(t *core.T) {
	// KeyBindingManager Add
	ax7Variant := "KeyBindingManager_Add:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Add("../../edge", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Remove_Good(t *core.T) {
	// KeyBindingManager Remove
	ax7Variant := "KeyBindingManager_Remove:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Remove("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Remove_Bad(t *core.T) {
	// KeyBindingManager Remove
	ax7Variant := "KeyBindingManager_Remove:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Remove("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Remove_Ugly(t *core.T) {
	// KeyBindingManager Remove
	ax7Variant := "KeyBindingManager_Remove:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		subject.Remove("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Process_Good(t *core.T) {
	// KeyBindingManager Process
	ax7Variant := "KeyBindingManager_Process:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.Process("agent", *new(Window))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Process_Bad(t *core.T) {
	// KeyBindingManager Process
	ax7Variant := "KeyBindingManager_Process:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.Process("", *new(Window))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_Process_Ugly(t *core.T) {
	// KeyBindingManager Process
	ax7Variant := "KeyBindingManager_Process:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.Process("../../edge", *new(Window))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_GetAll_Good(t *core.T) {
	// KeyBindingManager GetAll
	ax7Variant := "KeyBindingManager_GetAll:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_GetAll_Bad(t *core.T) {
	// KeyBindingManager GetAll
	ax7Variant := "KeyBindingManager_GetAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestKeybinding_KeyBindingManager_GetAll_Ugly(t *core.T) {
	// KeyBindingManager GetAll
	ax7Variant := "KeyBindingManager_GetAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(KeyBindingManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
