package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrowserManager_OpenURL_Good(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("https://example.com")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", manager.LastURL)
}

func TestBrowserManager_OpenURL_Bad(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("")

	assert.NoError(t, err)
	assert.Empty(t, manager.LastURL)
}

func TestBrowserManager_OpenURL_Ugly(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenURL("file:///tmp/%00")

	assert.NoError(t, err)
	assert.Equal(t, "file:///tmp/%00", manager.LastURL)
}

func TestBrowserManager_OpenFile_Good(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("/tmp/report.txt")

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/report.txt", manager.LastFile)
}

func TestBrowserManager_OpenFile_Bad(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("")

	assert.NoError(t, err)
	assert.Empty(t, manager.LastFile)
}

func TestBrowserManager_OpenFile_Ugly(t *testing.T) {
	manager := &BrowserManager{}

	err := manager.OpenFile("/tmp/\x00report.txt")

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/\x00report.txt", manager.LastFile)
}

func TestBrowserManager_NilReceiver_IsSafe(t *testing.T) {
	var manager *BrowserManager

	assert.NotPanics(t, func() {
		assert.NoError(t, manager.OpenURL("https://example.com"))
		assert.NoError(t, manager.OpenFile("/tmp/report.txt"))
	})
}
