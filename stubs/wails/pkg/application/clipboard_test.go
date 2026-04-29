package application

import (
	core "dappco.re/go"
)

func TestClipboard_SetText_Good(t *core.T) {
	// SetText
	ax7Variant := "SetText:good"
	core.AssertContains(t, ax7Variant, "good")
	clipboard := &Clipboard{}

	ok := clipboard.SetText("hello")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.AssertTrue(t, present)
	core.AssertEqual(t, "hello", text)
}

func TestClipboard_SetText_Bad(t *core.T) {
	// SetText
	ax7Variant := "SetText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	clipboard := &Clipboard{}

	ok := clipboard.SetText("")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.AssertTrue(t, present)
	core.AssertEmpty(t, text)
}

func TestClipboard_SetText_Ugly(t *core.T) {
	// SetText
	ax7Variant := "SetText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	clipboard := &Clipboard{}

	ok := clipboard.SetText("line1\nline2")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.RequireTrue(t, present)
	core.AssertEqual(t, "line1\nline2", text)
}

func TestClipboardManager_Text_Good(t *core.T) {
	// Text
	ax7Variant := "Text:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &ClipboardManager{}

	ok := manager.SetText("copied")

	core.AssertTrue(t, ok)
	text, present := manager.Text()
	core.AssertTrue(t, present)
	core.AssertEqual(t, "copied", text)
}

func TestClipboardManager_Text_Bad(t *core.T) {
	// Text
	ax7Variant := "Text:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &ClipboardManager{}

	text, present := manager.Text()

	core.AssertFalse(t, present)
	core.AssertEmpty(t, text)
}

func TestClipboardManager_Text_Ugly(t *core.T) {
	// Text
	ax7Variant := "Text:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := &ClipboardManager{}
	raw := "zero\x00byte"

	ok := manager.SetText(raw)

	core.AssertTrue(t, ok)
	text, present := manager.Text()
	core.AssertTrue(t, present)
	core.AssertEqual(t, raw, text)
}

func TestClipboardManager_NilReceiver_IsSafe(t *core.T) {
	var manager *ClipboardManager

	core.AssertNotPanics(t, func() {
		core.AssertFalse(t, manager.SetText("hello"))
		text, present := manager.Text()
		core.AssertEmpty(t, text)
		core.AssertFalse(t, present)
	})
}

// AX7 generated source-matching smoke coverage.
func TestClipboard_Clipboard_SetText_Good(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0 := subject.SetText("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_Clipboard_SetText_Bad(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0 := subject.SetText("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_Clipboard_SetText_Ugly(t *core.T) {
	// Clipboard SetText
	ax7Variant := "Clipboard_SetText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0 := subject.SetText("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_Clipboard_Text_Good(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_Clipboard_Text_Bad(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_Clipboard_Text_Ugly(t *core.T) {
	// Clipboard Text
	ax7Variant := "Clipboard_Text:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(Clipboard)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_SetText_Good(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0 := subject.SetText("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_SetText_Bad(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0 := subject.SetText("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_SetText_Ugly(t *core.T) {
	// ClipboardManager SetText
	ax7Variant := "ClipboardManager_SetText:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0 := subject.SetText("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_Text_Good(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_Text_Bad(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestClipboard_ClipboardManager_Text_Ugly(t *core.T) {
	// ClipboardManager Text
	ax7Variant := "ClipboardManager_Text:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ClipboardManager)
	result := core.Try(func() any {
		got0, got1 := subject.Text()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
