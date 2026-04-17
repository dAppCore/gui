package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenuItem_NewRole_Good(t *testing.T) {
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
		t.Run(tc.want, func(t *testing.T) {
			got := NewRole(tc.role)
			assert.Equal(t, tc.want, got.Label)
			assert.True(t, got.Enabled)
		})
	}
}

func TestMenuItem_NewRole_Bad(t *testing.T) {
	got := NewRole(Role(-1))

	assert.Empty(t, got.Label)
	assert.True(t, got.Enabled)
	assert.Equal(t, "unknown", MenuRole(999).String())
}

func TestMenuItem_NewRole_Ugly(t *testing.T) {
	item := NewMenuItem("Open")
	item.SetAccelerator("CmdOrCtrl+O")
	item.SetTooltip("Open a file")
	item.SetChecked(true)
	item.SetEnabled(false)

	assert.Equal(t, "Open", item.Label)
	assert.Equal(t, "CmdOrCtrl+O", item.GetAccelerator())
	assert.Equal(t, "Open a file", item.Tooltip)
	assert.True(t, item.Checked)
	assert.False(t, item.Enabled)
}
