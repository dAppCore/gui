package systray

import "testing"

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPlatform_NewTray_Good(t *testing.T) {
	p := NewMockPlatform()
	tray := p.NewTray()
	require.NotNil(t, tray)

	mockTray := tray.(*exportedMockTray)
	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("Core")
	tray.SetLabel("Ready")
	tray.SetMenu(p.NewMenu())
	tray.AttachWindow(windowHandleStub{name: "panel"})

	assert.Equal(t, []byte{1, 2, 3}, mockTray.icon)
	assert.Equal(t, []byte{4, 5, 6}, mockTray.templateIcon)
	assert.Equal(t, "Core", mockTray.tooltip)
	assert.Equal(t, "Ready", mockTray.label)
	assert.NotNil(t, mockTray)
}

func TestMockPlatform_NewMenu_Bad(t *testing.T) {
	p := NewMockPlatform()
	menu := p.NewMenu()
	require.NotNil(t, menu)
	_, ok := menu.(*exportedMockMenu)
	assert.True(t, ok)
}

func TestMockPlatform_NewTray_Ugly(t *testing.T) {
	p := NewMockPlatform()
	tray := p.NewTray().(*exportedMockTray)
	assert.NotNil(t, tray)
	assert.NoError(t, tray.ShowMessage("title", "message"))
}

type windowHandleStub struct {
	name string
}

func (w windowHandleStub) Name() string { return w.name }
