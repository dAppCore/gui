package application

import core "dappco.re/go"

func TestBrowserManager_BrowserManager_OpenURL_Good(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "BrowserManager_OpenURL:good"
	core.AssertContains(t, label, "BrowserManager_OpenURL")
	core.AssertContains(t, label, "good")
}

func TestBrowserManager_BrowserManager_OpenURL_Bad(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "BrowserManager_OpenURL:bad"
	core.AssertContains(t, label, "BrowserManager_OpenURL")
	core.AssertContains(t, label, "bad")
}

func TestBrowserManager_BrowserManager_OpenURL_Ugly(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "BrowserManager_OpenURL:ugly"
	core.AssertContains(t, label, "BrowserManager_OpenURL")
	core.AssertContains(t, label, "ugly")
}

func TestBrowserManager_BrowserManager_OpenFile_Good(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "BrowserManager_OpenFile:good"
	core.AssertContains(t, label, "BrowserManager_OpenFile")
	core.AssertContains(t, label, "good")
}

func TestBrowserManager_BrowserManager_OpenFile_Bad(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "BrowserManager_OpenFile:bad"
	core.AssertContains(t, label, "BrowserManager_OpenFile")
	core.AssertContains(t, label, "bad")
}

func TestBrowserManager_BrowserManager_OpenFile_Ugly(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "BrowserManager_OpenFile:ugly"
	core.AssertContains(t, label, "BrowserManager_OpenFile")
	core.AssertContains(t, label, "ugly")
}
