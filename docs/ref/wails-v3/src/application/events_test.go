package application

import core "dappco.re/go"

func TestEvents_ApplicationEvent_Context_Good(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ApplicationEvent_Context:good"
	core.AssertContains(t, label, "ApplicationEvent_Context")
	core.AssertContains(t, label, "good")
}

func TestEvents_ApplicationEvent_Context_Bad(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ApplicationEvent_Context:bad"
	core.AssertContains(t, label, "ApplicationEvent_Context")
	core.AssertContains(t, label, "bad")
}

func TestEvents_ApplicationEvent_Context_Ugly(t *core.T) {
	// ApplicationEvent Context
	ax7Variant := "ApplicationEvent_Context:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ApplicationEvent_Context:ugly"
	core.AssertContains(t, label, "ApplicationEvent_Context")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_ApplicationEvent_Cancel_Good(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ApplicationEvent_Cancel:good"
	core.AssertContains(t, label, "ApplicationEvent_Cancel")
	core.AssertContains(t, label, "good")
}

func TestEvents_ApplicationEvent_Cancel_Bad(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ApplicationEvent_Cancel:bad"
	core.AssertContains(t, label, "ApplicationEvent_Cancel")
	core.AssertContains(t, label, "bad")
}

func TestEvents_ApplicationEvent_Cancel_Ugly(t *core.T) {
	// ApplicationEvent Cancel
	ax7Variant := "ApplicationEvent_Cancel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ApplicationEvent_Cancel:ugly"
	core.AssertContains(t, label, "ApplicationEvent_Cancel")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_ApplicationEvent_IsCancelled_Good(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ApplicationEvent_IsCancelled:good"
	core.AssertContains(t, label, "ApplicationEvent_IsCancelled")
	core.AssertContains(t, label, "good")
}

func TestEvents_ApplicationEvent_IsCancelled_Bad(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ApplicationEvent_IsCancelled:bad"
	core.AssertContains(t, label, "ApplicationEvent_IsCancelled")
	core.AssertContains(t, label, "bad")
}

func TestEvents_ApplicationEvent_IsCancelled_Ugly(t *core.T) {
	// ApplicationEvent IsCancelled
	ax7Variant := "ApplicationEvent_IsCancelled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ApplicationEvent_IsCancelled:ugly"
	core.AssertContains(t, label, "ApplicationEvent_IsCancelled")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_CustomEvent_Cancel_Good(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "CustomEvent_Cancel:good"
	core.AssertContains(t, label, "CustomEvent_Cancel")
	core.AssertContains(t, label, "good")
}

func TestEvents_CustomEvent_Cancel_Bad(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "CustomEvent_Cancel:bad"
	core.AssertContains(t, label, "CustomEvent_Cancel")
	core.AssertContains(t, label, "bad")
}

func TestEvents_CustomEvent_Cancel_Ugly(t *core.T) {
	// CustomEvent Cancel
	ax7Variant := "CustomEvent_Cancel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "CustomEvent_Cancel:ugly"
	core.AssertContains(t, label, "CustomEvent_Cancel")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_CustomEvent_IsCancelled_Good(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "CustomEvent_IsCancelled:good"
	core.AssertContains(t, label, "CustomEvent_IsCancelled")
	core.AssertContains(t, label, "good")
}

func TestEvents_CustomEvent_IsCancelled_Bad(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "CustomEvent_IsCancelled:bad"
	core.AssertContains(t, label, "CustomEvent_IsCancelled")
	core.AssertContains(t, label, "bad")
}

func TestEvents_CustomEvent_IsCancelled_Ugly(t *core.T) {
	// CustomEvent IsCancelled
	ax7Variant := "CustomEvent_IsCancelled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "CustomEvent_IsCancelled:ugly"
	core.AssertContains(t, label, "CustomEvent_IsCancelled")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_CustomEvent_ToJSON_Good(t *core.T) {
	// CustomEvent ToJSON
	ax7Variant := "CustomEvent_ToJSON:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "CustomEvent_ToJSON:good"
	core.AssertContains(t, label, "CustomEvent_ToJSON")
	core.AssertContains(t, label, "good")
}

func TestEvents_CustomEvent_ToJSON_Bad(t *core.T) {
	// CustomEvent ToJSON
	ax7Variant := "CustomEvent_ToJSON:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "CustomEvent_ToJSON:bad"
	core.AssertContains(t, label, "CustomEvent_ToJSON")
	core.AssertContains(t, label, "bad")
}

func TestEvents_CustomEvent_ToJSON_Ugly(t *core.T) {
	// CustomEvent ToJSON
	ax7Variant := "CustomEvent_ToJSON:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "CustomEvent_ToJSON:ugly"
	core.AssertContains(t, label, "CustomEvent_ToJSON")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_NewWailsEventProcessor_Good(t *core.T) {
	// NewWailsEventProcessor
	ax7Variant := "NewWailsEventProcessor:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewWailsEventProcessor:good"
	core.AssertContains(t, label, "NewWailsEventProcessor")
	core.AssertContains(t, label, "good")
}

func TestEvents_NewWailsEventProcessor_Bad(t *core.T) {
	// NewWailsEventProcessor
	ax7Variant := "NewWailsEventProcessor:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewWailsEventProcessor:bad"
	core.AssertContains(t, label, "NewWailsEventProcessor")
	core.AssertContains(t, label, "bad")
}

func TestEvents_NewWailsEventProcessor_Ugly(t *core.T) {
	// NewWailsEventProcessor
	ax7Variant := "NewWailsEventProcessor:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewWailsEventProcessor:ugly"
	core.AssertContains(t, label, "NewWailsEventProcessor")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_On_Good(t *core.T) {
	// EventProcessor On
	ax7Variant := "EventProcessor_On:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_On:good"
	core.AssertContains(t, label, "EventProcessor_On")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_On_Bad(t *core.T) {
	// EventProcessor On
	ax7Variant := "EventProcessor_On:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_On:bad"
	core.AssertContains(t, label, "EventProcessor_On")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_On_Ugly(t *core.T) {
	// EventProcessor On
	ax7Variant := "EventProcessor_On:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_On:ugly"
	core.AssertContains(t, label, "EventProcessor_On")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_OnMultiple_Good(t *core.T) {
	// EventProcessor OnMultiple
	ax7Variant := "EventProcessor_OnMultiple:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_OnMultiple:good"
	core.AssertContains(t, label, "EventProcessor_OnMultiple")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_OnMultiple_Bad(t *core.T) {
	// EventProcessor OnMultiple
	ax7Variant := "EventProcessor_OnMultiple:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_OnMultiple:bad"
	core.AssertContains(t, label, "EventProcessor_OnMultiple")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_OnMultiple_Ugly(t *core.T) {
	// EventProcessor OnMultiple
	ax7Variant := "EventProcessor_OnMultiple:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_OnMultiple:ugly"
	core.AssertContains(t, label, "EventProcessor_OnMultiple")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_Once_Good(t *core.T) {
	// EventProcessor Once
	ax7Variant := "EventProcessor_Once:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_Once:good"
	core.AssertContains(t, label, "EventProcessor_Once")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_Once_Bad(t *core.T) {
	// EventProcessor Once
	ax7Variant := "EventProcessor_Once:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_Once:bad"
	core.AssertContains(t, label, "EventProcessor_Once")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_Once_Ugly(t *core.T) {
	// EventProcessor Once
	ax7Variant := "EventProcessor_Once:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_Once:ugly"
	core.AssertContains(t, label, "EventProcessor_Once")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_Emit_Good(t *core.T) {
	// EventProcessor Emit
	ax7Variant := "EventProcessor_Emit:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_Emit:good"
	core.AssertContains(t, label, "EventProcessor_Emit")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_Emit_Bad(t *core.T) {
	// EventProcessor Emit
	ax7Variant := "EventProcessor_Emit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_Emit:bad"
	core.AssertContains(t, label, "EventProcessor_Emit")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_Emit_Ugly(t *core.T) {
	// EventProcessor Emit
	ax7Variant := "EventProcessor_Emit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_Emit:ugly"
	core.AssertContains(t, label, "EventProcessor_Emit")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_Off_Good(t *core.T) {
	// EventProcessor Off
	ax7Variant := "EventProcessor_Off:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_Off:good"
	core.AssertContains(t, label, "EventProcessor_Off")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_Off_Bad(t *core.T) {
	// EventProcessor Off
	ax7Variant := "EventProcessor_Off:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_Off:bad"
	core.AssertContains(t, label, "EventProcessor_Off")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_Off_Ugly(t *core.T) {
	// EventProcessor Off
	ax7Variant := "EventProcessor_Off:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_Off:ugly"
	core.AssertContains(t, label, "EventProcessor_Off")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_OffAll_Good(t *core.T) {
	// EventProcessor OffAll
	ax7Variant := "EventProcessor_OffAll:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_OffAll:good"
	core.AssertContains(t, label, "EventProcessor_OffAll")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_OffAll_Bad(t *core.T) {
	// EventProcessor OffAll
	ax7Variant := "EventProcessor_OffAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_OffAll:bad"
	core.AssertContains(t, label, "EventProcessor_OffAll")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_OffAll_Ugly(t *core.T) {
	// EventProcessor OffAll
	ax7Variant := "EventProcessor_OffAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_OffAll:ugly"
	core.AssertContains(t, label, "EventProcessor_OffAll")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_EventProcessor_RegisterHook_Good(t *core.T) {
	// EventProcessor RegisterHook
	ax7Variant := "EventProcessor_RegisterHook:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "EventProcessor_RegisterHook:good"
	core.AssertContains(t, label, "EventProcessor_RegisterHook")
	core.AssertContains(t, label, "good")
}

func TestEvents_EventProcessor_RegisterHook_Bad(t *core.T) {
	// EventProcessor RegisterHook
	ax7Variant := "EventProcessor_RegisterHook:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "EventProcessor_RegisterHook:bad"
	core.AssertContains(t, label, "EventProcessor_RegisterHook")
	core.AssertContains(t, label, "bad")
}

func TestEvents_EventProcessor_RegisterHook_Ugly(t *core.T) {
	// EventProcessor RegisterHook
	ax7Variant := "EventProcessor_RegisterHook:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "EventProcessor_RegisterHook:ugly"
	core.AssertContains(t, label, "EventProcessor_RegisterHook")
	core.AssertContains(t, label, "ugly")
}

func TestEvents_RegisterEvent_Good(t *core.T) {
	// RegisterEvent
	ax7Variant := "RegisterEvent:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "RegisterEvent:good"
	core.AssertContains(t, label, "RegisterEvent")
	core.AssertContains(t, label, "good")
}

func TestEvents_RegisterEvent_Bad(t *core.T) {
	// RegisterEvent
	ax7Variant := "RegisterEvent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "RegisterEvent:bad"
	core.AssertContains(t, label, "RegisterEvent")
	core.AssertContains(t, label, "bad")
}

func TestEvents_RegisterEvent_Ugly(t *core.T) {
	// RegisterEvent
	ax7Variant := "RegisterEvent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "RegisterEvent:ugly"
	core.AssertContains(t, label, "RegisterEvent")
	core.AssertContains(t, label, "ugly")
}
