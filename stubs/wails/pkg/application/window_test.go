package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindow_Interface_Good(t *testing.T) {
	var _ Window = (*WebviewWindow)(nil)
	var _ Window = (*BrowserWindow)(nil)
}

func TestWindow_Interface_Bad(t *testing.T) {
	var w Window

	assert.Nil(t, w)
}

func TestWindow_Interface_Ugly(t *testing.T) {
	var w Window = (*WebviewWindow)(nil)

	assert.False(t, w == nil)
}
