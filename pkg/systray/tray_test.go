// pkg/systray/tray_test.go
package systray

import (
	core "dappco.re/go"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Setup_Good(t *core.T) {
	m, p := newTestManager()
	err := m.Setup("Core", "Core")
	core.RequireNoError(t, err)
	core.AssertTrue(t, m.IsActive())
	core.AssertLen(t, p.trays, 1)
	core.AssertEqual(t, "Core", p.trays[0].tooltip)
	core.AssertEqual(t, "Core", p.trays[0].label)
	core.AssertNotEmpty(t, p.trays[0].templateIcon) // default icon embedded
}

func TestManager_SetIcon_Good(t *core.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	err := m.SetIcon([]byte{1, 2, 3})
	core.RequireNoError(t, err)
	core.AssertEqual(t, []byte{1, 2, 3}, p.trays[0].icon)
}

func TestManager_SetIcon_Bad(t *core.T) {
	m, _ := newTestManager()
	err := m.SetIcon([]byte{1})
	core.AssertError(t, err) // tray not initialised
}

func TestManager_SetTooltip_Good(t *core.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetTooltip("New Tooltip")
	core.AssertEqual(t, "New Tooltip", p.trays[0].tooltip)
}

func TestManager_SetLabel_Good(t *core.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetLabel("New Label")
	core.AssertEqual(t, "New Label", p.trays[0].label)
}

func TestManager_RegisterCallback_Good(t *core.T) {
	m, _ := newTestManager()
	called := false
	m.RegisterCallback("test-action", func() { called = true })
	cb, ok := m.GetCallback("test-action")
	core.AssertTrue(t, ok)
	cb()
	core.AssertTrue(t, called)
}

func TestManager_RegisterCallback_Bad(t *core.T) {
	m, _ := newTestManager()
	_, ok := m.GetCallback("nonexistent")
	core.AssertFalse(t, ok)
}

func TestManager_UnregisterCallback_Good(t *core.T) {
	m, _ := newTestManager()
	m.RegisterCallback("remove-me", func() {})
	m.UnregisterCallback("remove-me")
	_, ok := m.GetCallback("remove-me")
	core.AssertFalse(t, ok)
}

func TestManager_GetInfo_Good(t *core.T) {
	m, _ := newTestManager()
	info := m.GetInfo()
	core.AssertFalse(t, info["active"].(bool))
	_ = m.Setup("Core", "Core")
	info = m.GetInfo()
	core.AssertTrue(t, info["active"].(bool))
}

func TestManager_Build_Submenu_Recursive_Good(t *core.T) {
	m, p := newTestManager()
	core.RequireNoError(t, m.Setup("Core", "Core"))

	items := []TrayMenuItem{
		{
			Label: "Parent",
			Submenu: []TrayMenuItem{
				{Label: "Child 1"},
				{Label: "Child 2"},
			},
		},
	}

	core.RequireNoError(t, m.SetMenu(items))
	core.AssertLen(t, p.menus, 1)

	menu := p.menus[0]
	core.AssertLen(t, menu.items, 1)
	core.AssertEqual(t, "Parent", menu.items[0])
	core.AssertLen(t, menu.subs, 1)
	core.AssertLen(t, menu.subs[0].items, 2)
	core.AssertEqual(t, "Child 1", menu.subs[0].items[0])
	core.AssertEqual(t, "Child 2", menu.subs[0].items[1])
}
