package application

import (
	core "dappco.re/go"
)

func TestClipboard_SetText_Good(t *core.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("hello")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.AssertTrue(t, present)
	core.AssertEqual(t, "hello", text)
}

func TestClipboard_SetText_Bad(t *core.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.AssertTrue(t, present)
	core.AssertEmpty(t, text)
}

func TestClipboard_SetText_Ugly(t *core.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("line1\nline2")

	core.AssertTrue(t, ok)
	text, present := clipboard.Text()
	core.RequireTrue(t, present)
	core.AssertEqual(t, "line1\nline2", text)
}

func TestClipboardManager_Text_Good(t *core.T) {
	manager := &ClipboardManager{}

	ok := manager.SetText("copied")

	core.AssertTrue(t, ok)
	text, present := manager.Text()
	core.AssertTrue(t, present)
	core.AssertEqual(t, "copied", text)
}

func TestClipboardManager_Text_Bad(t *core.T) {
	manager := &ClipboardManager{}

	text, present := manager.Text()

	core.AssertFalse(t, present)
	core.AssertEmpty(t, text)
}

func TestClipboardManager_Text_Ugly(t *core.T) {
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
