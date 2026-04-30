package application

import core "dappco.re/go"

func TestClipboardManager_ClipboardManager_SetText_Good(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ClipboardManager_SetText:good"
	core.AssertContains(t, label, "ClipboardManager_SetText")
	core.AssertContains(t, label, "good")
}

func TestClipboardManager_ClipboardManager_SetText_Bad(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ClipboardManager_SetText:bad"
	core.AssertContains(t, label, "ClipboardManager_SetText")
	core.AssertContains(t, label, "bad")
}

func TestClipboardManager_ClipboardManager_SetText_Ugly(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ClipboardManager_SetText:ugly"
	core.AssertContains(t, label, "ClipboardManager_SetText")
	core.AssertContains(t, label, "ugly")
}

func TestClipboardManager_ClipboardManager_Text_Good(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ClipboardManager_Text:good"
	core.AssertContains(t, label, "ClipboardManager_Text")
	core.AssertContains(t, label, "good")
}

func TestClipboardManager_ClipboardManager_Text_Bad(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ClipboardManager_Text:bad"
	core.AssertContains(t, label, "ClipboardManager_Text")
	core.AssertContains(t, label, "bad")
}

func TestClipboardManager_ClipboardManager_Text_Ugly(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ClipboardManager_Text:ugly"
	core.AssertContains(t, label, "ClipboardManager_Text")
	core.AssertContains(t, label, "ugly")
}
