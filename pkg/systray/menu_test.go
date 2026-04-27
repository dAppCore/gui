package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingTrayPlatform struct {
	tray *recordingTray
	menu *recordingTrayMenu
}

func (p *recordingTrayPlatform) NewTray() PlatformTray {
	p.tray = &recordingTray{}
	return p.tray
}

func (p *recordingTrayPlatform) NewMenu() PlatformMenu {
	p.menu = &recordingTrayMenu{}
	return p.menu
}

type recordingTray struct {
	icon, templateIcon []byte
	tooltip, label     string
	menu               PlatformMenu
	attachedWindow     WindowHandle
}

func (t *recordingTray) SetIcon(data []byte)                     { t.icon = append([]byte(nil), data...) }
func (t *recordingTray) SetTemplateIcon(data []byte)             { t.templateIcon = append([]byte(nil), data...) }
func (t *recordingTray) SetTooltip(text string)                  { t.tooltip = text }
func (t *recordingTray) SetLabel(text string)                    { t.label = text }
func (t *recordingTray) SetMenu(menu PlatformMenu)               { t.menu = menu }
func (t *recordingTray) AttachWindow(w WindowHandle)             { t.attachedWindow = w }
func (t *recordingTray) ShowMessage(title, message string) error { return nil }

type recordingTrayMenu struct {
	items []*recordingTrayMenuItem
	subs  []*recordingTrayMenu
}

func (m *recordingTrayMenu) Add(label string) PlatformMenuItem {
	item := &recordingTrayMenuItem{label: label, enabled: true}
	m.items = append(m.items, item)
	return item
}

func (m *recordingTrayMenu) AddSeparator() {
	m.items = append(m.items, &recordingTrayMenuItem{label: "---"})
}

func (m *recordingTrayMenu) AddSubmenu(label string) PlatformMenu {
	sub := &recordingTrayMenu{}
	m.items = append(m.items, &recordingTrayMenuItem{label: label, submenu: sub})
	m.subs = append(m.subs, sub)
	return sub
}

type recordingTrayMenuItem struct {
	label, tooltip   string
	checked, enabled bool
	submenu          *recordingTrayMenu
	onClick          func()
}

func (i *recordingTrayMenuItem) SetTooltip(text string)  { i.tooltip = text }
func (i *recordingTrayMenuItem) SetChecked(checked bool) { i.checked = checked }
func (i *recordingTrayMenuItem) SetEnabled(enabled bool) { i.enabled = enabled }
func (i *recordingTrayMenuItem) OnClick(fn func())       { i.onClick = fn }

func TestManager_SetMenu_Good(t *testing.T) {
	platform := &recordingTrayPlatform{}
	mgr := NewManager(platform)
	require.NoError(t, mgr.Setup("Core", "Core"))

	clicked := 0
	items := []TrayMenuItem{
		{Label: "Open", Tooltip: "open", ActionID: "open"},
		{Type: "separator"},
		{Label: "More", Submenu: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}}},
		{Label: "Disabled", Disabled: true},
		{Label: "Checked", Checked: true},
	}
	mgr.RegisterCallback("open", func() { clicked++ })
	mgr.RegisterCallback("nested", func() { clicked += 10 })

	require.NoError(t, mgr.SetMenu(items))
	require.NotNil(t, platform.menu)
	require.NotNil(t, platform.tray)

	assert.Equal(t, "Core", platform.tray.tooltip)
	assert.Equal(t, "Core", platform.tray.label)
	assert.Len(t, platform.menu.items, 5)
	assert.Equal(t, "Open", platform.menu.items[0].label)
	assert.Equal(t, "---", platform.menu.items[1].label)
	assert.Equal(t, "More", platform.menu.items[2].label)
	assert.False(t, platform.menu.items[3].enabled)
	assert.True(t, platform.menu.items[4].checked)

	require.NotNil(t, platform.menu.items[0].onClick)
	platform.menu.items[0].onClick()
	require.Len(t, platform.menu.subs, 1)
	require.NotNil(t, platform.menu.subs[0].items[0].onClick)
	platform.menu.subs[0].items[0].onClick()

	assert.Equal(t, 11, clicked)
	assert.Len(t, mgr.GetInfo()["menuItems"].([]TrayMenuItem), 5)
}

func TestManager_SetMenu_Bad(t *testing.T) {
	mgr := NewManager(&recordingTrayPlatform{})
	err := mgr.SetMenu([]TrayMenuItem{{Label: "Quit"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tray not initialised")
}

func TestManager_GetCallback_Ugly(t *testing.T) {
	mgr := NewManager(&recordingTrayPlatform{})
	mgr.RegisterCallback("quit", func() {})
	cb, ok := mgr.GetCallback("quit")
	require.True(t, ok)
	assert.NotNil(t, cb)

	mgr.UnregisterCallback("quit")
	_, ok = mgr.GetCallback("quit")
	assert.False(t, ok)
}
