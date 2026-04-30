package application

import core "dappco.re/go"

func TestKeyBinding_KeyBindingManager_Register_Good(t *core.T) {
	// KeyBindingManager Register
	ax7Variant := "KeyBindingManager_Register:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "KeyBindingManager_Register:good"
	core.AssertContains(t, label, "KeyBindingManager_Register")
	core.AssertContains(t, label, "good")
}

func TestKeyBinding_KeyBindingManager_Register_Bad(t *core.T) {
	// KeyBindingManager Register
	ax7Variant := "KeyBindingManager_Register:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "KeyBindingManager_Register:bad"
	core.AssertContains(t, label, "KeyBindingManager_Register")
	core.AssertContains(t, label, "bad")
}

func TestKeyBinding_KeyBindingManager_Register_Ugly(t *core.T) {
	// KeyBindingManager Register
	ax7Variant := "KeyBindingManager_Register:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "KeyBindingManager_Register:ugly"
	core.AssertContains(t, label, "KeyBindingManager_Register")
	core.AssertContains(t, label, "ugly")
}

func TestKeyBinding_KeyBindingManager_Unregister_Good(t *core.T) {
	// KeyBindingManager Unregister
	ax7Variant := "KeyBindingManager_Unregister:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "KeyBindingManager_Unregister:good"
	core.AssertContains(t, label, "KeyBindingManager_Unregister")
	core.AssertContains(t, label, "good")
}

func TestKeyBinding_KeyBindingManager_Unregister_Bad(t *core.T) {
	// KeyBindingManager Unregister
	ax7Variant := "KeyBindingManager_Unregister:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "KeyBindingManager_Unregister:bad"
	core.AssertContains(t, label, "KeyBindingManager_Unregister")
	core.AssertContains(t, label, "bad")
}

func TestKeyBinding_KeyBindingManager_Unregister_Ugly(t *core.T) {
	// KeyBindingManager Unregister
	ax7Variant := "KeyBindingManager_Unregister:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "KeyBindingManager_Unregister:ugly"
	core.AssertContains(t, label, "KeyBindingManager_Unregister")
	core.AssertContains(t, label, "ugly")
}
