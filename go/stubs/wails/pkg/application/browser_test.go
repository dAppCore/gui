package application

import core "dappco.re/go"

func TestBrowser_BrowserManager_Open_Good(t *core.T) {
	// BrowserManager Open
	ax7Variant := "BrowserManager_Open:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "BrowserManager_Open:good"
	core.AssertContains(t, label, "BrowserManager_Open")
	core.AssertContains(t, label, "good")
}

func TestBrowser_BrowserManager_Open_Bad(t *core.T) {
	// BrowserManager Open
	ax7Variant := "BrowserManager_Open:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "BrowserManager_Open:bad"
	core.AssertContains(t, label, "BrowserManager_Open")
	core.AssertContains(t, label, "bad")
}

func TestBrowser_BrowserManager_Open_Ugly(t *core.T) {
	// BrowserManager Open
	ax7Variant := "BrowserManager_Open:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "BrowserManager_Open:ugly"
	core.AssertContains(t, label, "BrowserManager_Open")
	core.AssertContains(t, label, "ugly")
}
