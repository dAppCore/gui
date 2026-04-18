package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClipboard_SetText_Good(t *testing.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("hello")

	assert.True(t, ok)
	text, present := clipboard.Text()
	assert.True(t, present)
	assert.Equal(t, "hello", text)
}

func TestClipboard_SetText_Bad(t *testing.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("")

	assert.True(t, ok)
	text, present := clipboard.Text()
	assert.True(t, present)
	assert.Empty(t, text)
}

func TestClipboard_SetText_Ugly(t *testing.T) {
	clipboard := &Clipboard{}

	ok := clipboard.SetText("line1\nline2")

	assert.True(t, ok)
	text, present := clipboard.Text()
	require.True(t, present)
	assert.Equal(t, "line1\nline2", text)
}

func TestClipboardManager_Text_Good(t *testing.T) {
	manager := &ClipboardManager{}

	ok := manager.SetText("copied")

	assert.True(t, ok)
	text, present := manager.Text()
	assert.True(t, present)
	assert.Equal(t, "copied", text)
}

func TestClipboardManager_Text_Bad(t *testing.T) {
	manager := &ClipboardManager{}

	text, present := manager.Text()

	assert.False(t, present)
	assert.Empty(t, text)
}

func TestClipboardManager_Text_Ugly(t *testing.T) {
	manager := &ClipboardManager{}
	raw := "zero\x00byte"

	ok := manager.SetText(raw)

	assert.True(t, ok)
	text, present := manager.Text()
	assert.True(t, present)
	assert.Equal(t, raw, text)
}

func TestClipboardManager_NilReceiver_IsSafe(t *testing.T) {
	var manager *ClipboardManager

	assert.NotPanics(t, func() {
		assert.False(t, manager.SetText("hello"))
		text, present := manager.Text()
		assert.Empty(t, text)
		assert.False(t, present)
	})
}
