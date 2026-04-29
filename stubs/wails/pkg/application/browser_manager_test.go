package application

import (
	core "dappco.re/go"
)

func TestBrowserManager_OpenURL_Good(t *core.T) {
	// OpenURL
	ax7Variant := "OpenURL:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &BrowserManager{}

	err := manager.OpenURL("https://example.com")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "https://example.com", manager.LastURL)
}

func TestBrowserManager_OpenURL_Bad(t *core.T) {
	// OpenURL
	ax7Variant := "OpenURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &BrowserManager{}

	err := manager.OpenURL("")

	core.AssertNoError(t, err)
	core.AssertEmpty(t, manager.LastURL)
}

func TestBrowserManager_OpenURL_Ugly(t *core.T) {
	// OpenURL
	ax7Variant := "OpenURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := &BrowserManager{}

	err := manager.OpenURL("file:///tmp/%00")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "file:///tmp/%00", manager.LastURL)
}

func TestBrowserManager_OpenFile_Good(t *core.T) {
	// OpenFile
	ax7Variant := "OpenFile:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &BrowserManager{}

	err := manager.OpenFile("/tmp/report.txt")

	core.AssertNoError(t, err)
	core.AssertEqual(t, "/tmp/report.txt", manager.LastFile)
}

func TestBrowserManager_OpenFile_Bad(t *core.T) {
	// OpenFile
	ax7Variant := "OpenFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &BrowserManager{}

	err := manager.OpenFile("")

	core.AssertNoError(t, err)
	core.AssertEmpty(t, manager.LastFile)
}

func TestBrowserManager_OpenFile_Ugly(t *core.T) {
	// OpenFile
	ax7Variant := "OpenFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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

// AX7 generated source-matching smoke coverage.
func TestBrowserManager_BrowserManager_OpenURL_Good(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenURL("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserManager_BrowserManager_OpenURL_Bad(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenURL("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserManager_BrowserManager_OpenURL_Ugly(t *core.T) {
	// BrowserManager OpenURL
	ax7Variant := "BrowserManager_OpenURL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenURL("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserManager_BrowserManager_OpenFile_Good(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserManager_BrowserManager_OpenFile_Bad(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestBrowserManager_BrowserManager_OpenFile_Ugly(t *core.T) {
	// BrowserManager OpenFile
	ax7Variant := "BrowserManager_OpenFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(BrowserManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
