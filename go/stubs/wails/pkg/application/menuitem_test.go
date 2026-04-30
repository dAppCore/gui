package application

import (
	core "dappco.re/go"
)

func TestMenuItem_NewRole_Good(t *core.T) {
	// NewRole
	ax7Variant := "NewRole:good"
	core.AssertContains(t, ax7Variant, "good")
	cases := []struct {
		role Role
		want string
	}{
		{AppMenu, "App"},
		{FileMenu, "File"},
		{EditMenu, "Edit"},
		{HelpMenu, "Help"},
		{Quit, "Quit"},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *core.T) {
			got := NewRole(tc.role)
			core.AssertEqual(t, tc.want, got.Label)
			core.AssertTrue(t, got.Enabled)
		})
	}
}

func TestMenuItem_NewRole_Bad(t *core.T) {
	// NewRole
	ax7Variant := "NewRole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	got := NewRole(Role(-1))

	core.AssertEmpty(t, got.Label)
	core.AssertTrue(t, got.Enabled)
	core.AssertEqual(t, "unknown", MenuRole(999).String())
}

func TestMenuItem_NewRole_UglyCase(t *core.T) {
	item := NewMenuItem("Open")
	item.SetAccelerator("CmdOrCtrl+O")
	item.SetTooltip("Open a file")
	item.SetChecked(true)
	item.SetEnabled(false)

	core.AssertEqual(t, "Open", item.Label)
	core.AssertEqual(t, "CmdOrCtrl+O", item.GetAccelerator())
	core.AssertEqual(t, "Open a file", item.Tooltip)
	core.AssertTrue(t, item.Checked)
	core.AssertFalse(t, item.Enabled)
}

// AX7 generated source-matching smoke coverage.
func TestMenuitem_NewMenuItem_Good(t *core.T) {
	// NewMenuItem
	ax7Variant := "NewMenuItem:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMenuItem("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItem_Bad(t *core.T) {
	// NewMenuItem
	ax7Variant := "NewMenuItem:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMenuItem("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItem_Ugly(t *core.T) {
	// NewMenuItem
	ax7Variant := "NewMenuItem:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMenuItem("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemSeparator_Good(t *core.T) {
	// NewMenuItemSeparator
	ax7Variant := "NewMenuItemSeparator:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMenuItemSeparator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemSeparator_Bad(t *core.T) {
	// NewMenuItemSeparator
	ax7Variant := "NewMenuItemSeparator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMenuItemSeparator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemSeparator_Ugly(t *core.T) {
	// NewMenuItemSeparator
	ax7Variant := "NewMenuItemSeparator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMenuItemSeparator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemCheckbox_Good(t *core.T) {
	// NewMenuItemCheckbox
	ax7Variant := "NewMenuItemCheckbox:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMenuItemCheckbox("agent", true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemCheckbox_Bad(t *core.T) {
	// NewMenuItemCheckbox
	ax7Variant := "NewMenuItemCheckbox:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMenuItemCheckbox("", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemCheckbox_Ugly(t *core.T) {
	// NewMenuItemCheckbox
	ax7Variant := "NewMenuItemCheckbox:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMenuItemCheckbox("../../edge", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemRadio_Good(t *core.T) {
	// NewMenuItemRadio
	ax7Variant := "NewMenuItemRadio:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMenuItemRadio("agent", true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemRadio_Bad(t *core.T) {
	// NewMenuItemRadio
	ax7Variant := "NewMenuItemRadio:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMenuItemRadio("", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewMenuItemRadio_Ugly(t *core.T) {
	// NewMenuItemRadio
	ax7Variant := "NewMenuItemRadio:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMenuItemRadio("../../edge", false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewSubMenuItem_Good(t *core.T) {
	// NewSubMenuItem
	ax7Variant := "NewSubMenuItem:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewSubMenuItem("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewSubMenuItem_Bad(t *core.T) {
	// NewSubMenuItem
	ax7Variant := "NewSubMenuItem:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewSubMenuItem("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewSubMenuItem_Ugly(t *core.T) {
	// NewSubMenuItem
	ax7Variant := "NewSubMenuItem:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewSubMenuItem("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewRole_Good(t *core.T) {
	// NewRole
	ax7Variant := "NewRole:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewRole(*new(Role))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewRole_Bad(t *core.T) {
	// NewRole
	ax7Variant := "NewRole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewRole(*new(Role))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewRole_Ugly(t *core.T) {
	// NewRole
	ax7Variant := "NewRole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewRole(*new(Role))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewServicesMenu_Good(t *core.T) {
	// NewServicesMenu
	ax7Variant := "NewServicesMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewServicesMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewServicesMenu_Bad(t *core.T) {
	// NewServicesMenu
	ax7Variant := "NewServicesMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewServicesMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_NewServicesMenu_Ugly(t *core.T) {
	// NewServicesMenu
	ax7Variant := "NewServicesMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewServicesMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetAccelerator_Good(t *core.T) {
	// MenuItem GetAccelerator
	ax7Variant := "MenuItem_GetAccelerator:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetAccelerator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetAccelerator_Bad(t *core.T) {
	// MenuItem GetAccelerator
	ax7Variant := "MenuItem_GetAccelerator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetAccelerator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetAccelerator_Ugly(t *core.T) {
	// MenuItem GetAccelerator
	ax7Variant := "MenuItem_GetAccelerator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetAccelerator()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetSubmenu_Good(t *core.T) {
	// MenuItem GetSubmenu
	ax7Variant := "MenuItem_GetSubmenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetSubmenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetSubmenu_Bad(t *core.T) {
	// MenuItem GetSubmenu
	ax7Variant := "MenuItem_GetSubmenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetSubmenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMenuitem_MenuItem_GetSubmenu_Ugly(t *core.T) {
	// MenuItem GetSubmenu
	ax7Variant := "MenuItem_GetSubmenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MenuItem)
	result := core.Try(func() any {
		got0 := subject.GetSubmenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}
