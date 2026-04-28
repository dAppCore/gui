package systray

import (
	core "dappco.re/go"
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

func TestManager_SetMenu_Good(t *core.T) {
	platform := &recordingTrayPlatform{}
	mgr := NewManager(platform)
	core.RequireNoError(t, mgr.Setup("Core", "Core"))

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

	core.RequireNoError(t, mgr.SetMenu(items))
	core.AssertNotNil(t, platform.menu)
	core.AssertNotNil(t, platform.tray)

	core.AssertEqual(t, "Core", platform.tray.tooltip)
	core.AssertEqual(t, "Core", platform.tray.label)
	core.AssertLen(t, platform.menu.items, 5)
	core.AssertEqual(t, "Open", platform.menu.items[0].label)
	core.AssertEqual(t, "---", platform.menu.items[1].label)
	core.AssertEqual(t, "More", platform.menu.items[2].label)
	core.AssertFalse(t, platform.menu.items[3].enabled)
	core.AssertTrue(t, platform.menu.items[4].checked)

	core.AssertNotNil(t, platform.menu.items[0].onClick)
	platform.menu.items[0].onClick()
	core.AssertLen(t, platform.menu.subs, 1)
	core.AssertNotNil(t, platform.menu.subs[0].items[0].onClick)
	platform.menu.subs[0].items[0].onClick()

	core.AssertEqual(t, 11, clicked)
	core.AssertLen(t, mgr.GetInfo()["menuItems"].([]TrayMenuItem), 5)
}

func TestManager_SetMenu_Bad(t *core.T) {
	mgr := NewManager(&recordingTrayPlatform{})
	err := mgr.SetMenu([]TrayMenuItem{{Label: "Quit"}})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "tray not initialised")
}

func TestManager_GetCallback_Ugly(t *core.T) {
	mgr := NewManager(&recordingTrayPlatform{})
	mgr.RegisterCallback("quit", func() {})
	cb, ok := mgr.GetCallback("quit")
	core.RequireTrue(t, ok)
	core.AssertNotNil(t, cb)

	mgr.UnregisterCallback("quit")
	_, ok = mgr.GetCallback("quit")
	core.AssertFalse(t, ok)
}
