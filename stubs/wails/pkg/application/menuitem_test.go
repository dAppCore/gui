package application

import (
	core "dappco.re/go"
)

func TestMenuItem_NewRole_Good(t *core.T) {
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
	got := NewRole(Role(-1))

	core.AssertEmpty(t, got.Label)
	core.AssertTrue(t, got.Enabled)
	core.AssertEqual(t, "unknown", MenuRole(999).String())
}

func TestMenuItem_NewRole_Ugly(t *core.T) {
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
