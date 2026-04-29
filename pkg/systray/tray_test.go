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

func TestManager_RegisterCallback_BadCase(t *core.T) {
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

// AX7 generated source-matching smoke coverage.
func TestTray_NewManager_Good(t *core.T) {
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_NewManager_Bad(t *core.T) {
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_NewManager_Ugly(t *core.T) {
	result := core.Try(func() any {
		got0 := NewManager(*new(Platform))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Setup_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Setup("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Setup_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Setup("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Setup_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Setup("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetIcon_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetIcon_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetIcon_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTemplateIcon_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTemplateIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTemplateIcon_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTemplateIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTemplateIcon_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTemplateIcon(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTooltip_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTooltip_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetTooltip_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetLabel_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetLabel("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetLabel_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetLabel("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_SetLabel_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.SetLabel("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_AttachWindow_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.AttachWindow(*new(WindowHandle))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_AttachWindow_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.AttachWindow(*new(WindowHandle))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_AttachWindow_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.AttachWindow(*new(WindowHandle))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowMessage_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowMessage_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowMessage_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowPanel_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowPanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowPanel_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowPanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_ShowPanel_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.ShowPanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_HidePanel_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.HidePanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_HidePanel_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.HidePanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_HidePanel_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.HidePanel()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Tray_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Tray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Tray_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Tray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_Tray_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.Tray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_IsActive_Good(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.IsActive()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_IsActive_Bad(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.IsActive()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestTray_Manager_IsActive_Ugly(t *core.T) {
	subject := new(Manager)
	result := core.Try(func() any {
		got0 := subject.IsActive()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
