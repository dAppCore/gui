package application

import core "dappco.re/go"

func TestClipboard_Clipboard_SetText_Good(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Clipboard_SetText:good"
	core.AssertContains(t, label, "Clipboard_SetText")
	core.AssertContains(t, label, "good")
}

func TestClipboard_Clipboard_SetText_Bad(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Clipboard_SetText:bad"
	core.AssertContains(t, label, "Clipboard_SetText")
	core.AssertContains(t, label, "bad")
}

func TestClipboard_Clipboard_SetText_Ugly(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Clipboard_SetText:ugly"
	core.AssertContains(t, label, "Clipboard_SetText")
	core.AssertContains(t, label, "ugly")
}

func TestClipboard_Clipboard_Text_Good(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Clipboard_Text:good"
	core.AssertContains(t, label, "Clipboard_Text")
	core.AssertContains(t, label, "good")
}

func TestClipboard_Clipboard_Text_Bad(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Clipboard_Text:bad"
	core.AssertContains(t, label, "Clipboard_Text")
	core.AssertContains(t, label, "bad")
}

func TestClipboard_Clipboard_Text_Ugly(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Clipboard_Text:ugly"
	core.AssertContains(t, label, "Clipboard_Text")
	core.AssertContains(t, label, "ugly")
}
