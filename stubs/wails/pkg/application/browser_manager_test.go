package application

import (
	core "dappco.re/go"
)

func TestBrowserManager_OpenURL_Good(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("https://example.com")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "https://example.com", manager.LastURL)
}

func TestBrowserManager_OpenURL_Bad(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("")

	core.AssertNoError(t, err)
	core.AssertEmpty(t, manager.LastURL)
}

func TestBrowserManager_OpenURL_Ugly(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("file:///tmp/%00")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "file:///tmp/%00", manager.LastURL)
}

func TestBrowserManager_OpenFile_Good(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("/tmp/report.txt")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "/tmp/report.txt", manager.LastFile)
}

func TestBrowserManager_OpenFile_Bad(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("")

	core.AssertNoError(t, err)
	core.AssertEmpty(t, manager.LastFile)
}

func TestBrowserManager_OpenFile_Ugly(t *core.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("/tmp/\x00report.txt")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "/tmp/\x00report.txt", manager.LastFile)
}

func TestBrowserManager_NilReceiver_IsSafe(t *core.T) {
	var manager *BrowserManager

	core.AssertNotPanics(t, func() {
		core.AssertNoError(t, manager.OpenURL("https://example.com"))
		core.AssertNoError(t, manager.OpenFile("/tmp/report.txt"))
	})
}
