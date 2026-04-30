package application

import core "dappco.re/go"

func TestEvent_EventManager_Once_Good(t *core.T) {
	// EventManager Once
	ax7Variant := "EventManager_Once:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventManager_Once:good"
	core.AssertContains(t, label, "EventManager_Once")
	core.AssertContains(t, label, "good")
}

func TestEvent_EventManager_Once_Bad(t *core.T) {
	// EventManager Once
	ax7Variant := "EventManager_Once:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventManager_Once:bad"
	core.AssertContains(t, label, "EventManager_Once")
	core.AssertContains(t, label, "bad")
}

func TestEvent_EventManager_Once_Ugly(t *core.T) {
	// EventManager Once
	ax7Variant := "EventManager_Once:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventManager_Once:ugly"
	core.AssertContains(t, label, "EventManager_Once")
	core.AssertContains(t, label, "ugly")
}

func TestEvent_EventManager_EmitEvent_Good(t *core.T) {
	// EventManager EmitEvent
	ax7Variant := "EventManager_EmitEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventManager_EmitEvent:good"
	core.AssertContains(t, label, "EventManager_EmitEvent")
	core.AssertContains(t, label, "good")
}

func TestEvent_EventManager_EmitEvent_Bad(t *core.T) {
	// EventManager EmitEvent
	ax7Variant := "EventManager_EmitEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventManager_EmitEvent:bad"
	core.AssertContains(t, label, "EventManager_EmitEvent")
	core.AssertContains(t, label, "bad")
}

func TestEvent_EventManager_EmitEvent_Ugly(t *core.T) {
	// EventManager EmitEvent
	ax7Variant := "EventManager_EmitEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventManager_EmitEvent:ugly"
	core.AssertContains(t, label, "EventManager_EmitEvent")
	core.AssertContains(t, label, "ugly")
}

func TestEvent_EventManager_Reset_Good(t *core.T) {
	// EventManager Reset
	ax7Variant := "EventManager_Reset:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventManager_Reset:good"
	core.AssertContains(t, label, "EventManager_Reset")
	core.AssertContains(t, label, "good")
}

func TestEvent_EventManager_Reset_Bad(t *core.T) {
	// EventManager Reset
	ax7Variant := "EventManager_Reset:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventManager_Reset:bad"
	core.AssertContains(t, label, "EventManager_Reset")
	core.AssertContains(t, label, "bad")
}

func TestEvent_EventManager_Reset_Ugly(t *core.T) {
	// EventManager Reset
	ax7Variant := "EventManager_Reset:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventManager_Reset:ugly"
	core.AssertContains(t, label, "EventManager_Reset")
	core.AssertContains(t, label, "ugly")
}

func TestEvent_EventManager_RegisterApplicationEventHook_Good(t *core.T) {
	// EventManager RegisterApplicationEventHook
	ax7Variant := "EventManager_RegisterApplicationEventHook:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventManager_RegisterApplicationEventHook:good"
	core.AssertContains(t, label, "EventManager_RegisterApplicationEventHook")
	core.AssertContains(t, label, "good")
}

func TestEvent_EventManager_RegisterApplicationEventHook_Bad(t *core.T) {
	// EventManager RegisterApplicationEventHook
	ax7Variant := "EventManager_RegisterApplicationEventHook:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventManager_RegisterApplicationEventHook:bad"
	core.AssertContains(t, label, "EventManager_RegisterApplicationEventHook")
	core.AssertContains(t, label, "bad")
}

func TestEvent_EventManager_RegisterApplicationEventHook_Ugly(t *core.T) {
	// EventManager RegisterApplicationEventHook
	ax7Variant := "EventManager_RegisterApplicationEventHook:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventManager_RegisterApplicationEventHook:ugly"
	core.AssertContains(t, label, "EventManager_RegisterApplicationEventHook")
	core.AssertContains(t, label, "ugly")
}
