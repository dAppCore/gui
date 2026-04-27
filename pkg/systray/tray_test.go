// pkg/systray/tray_test.go
package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestManager() (*Manager, *mockPlatform) {
	p := newMockPlatform()
	return NewManager(p), p
}

func TestManager_Setup_Good(t *testing.T) {
	m, p := newTestManager()
	err := m.Setup("Core", "Core")
	require.NoError(t, err)
	assert.True(t, m.IsActive())
	assert.Len(t, p.trays, 1)
	assert.Equal(t, "Core", p.trays[0].tooltip)
	assert.Equal(t, "Core", p.trays[0].label)
	assert.NotEmpty(t, p.trays[0].templateIcon) // default icon embedded
}

func TestManager_SetIcon_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	err := m.SetIcon([]byte{1, 2, 3})
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, p.trays[0].icon)
}

func TestManager_SetIcon_Bad(t *testing.T) {
	m, _ := newTestManager()
	err := m.SetIcon([]byte{1})
	assert.Error(t, err) // tray not initialised
}

func TestManager_SetTooltip_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetTooltip("New Tooltip")
	assert.Equal(t, "New Tooltip", p.trays[0].tooltip)
}

func TestManager_SetLabel_Good(t *testing.T) {
	m, p := newTestManager()
	_ = m.Setup("Core", "Core")
	_ = m.SetLabel("New Label")
	assert.Equal(t, "New Label", p.trays[0].label)
}

func TestManager_RegisterCallback_Good(t *testing.T) {
	m, _ := newTestManager()
	called := false
	m.RegisterCallback("test-action", func() { called = true })
	cb, ok := m.GetCallback("test-action")
	assert.True(t, ok)
	cb()
	assert.True(t, called)
}

func TestManager_RegisterCallback_Bad(t *testing.T) {
	m, _ := newTestManager()
	_, ok := m.GetCallback("nonexistent")
	assert.False(t, ok)
}

func TestManager_UnregisterCallback_Good(t *testing.T) {
	m, _ := newTestManager()
	m.RegisterCallback("remove-me", func() {})
	m.UnregisterCallback("remove-me")
	_, ok := m.GetCallback("remove-me")
	assert.False(t, ok)
}

func TestManager_GetInfo_Good(t *testing.T) {
	m, _ := newTestManager()
	info := m.GetInfo()
	assert.False(t, info["active"].(bool))
	_ = m.Setup("Core", "Core")
	info = m.GetInfo()
	assert.True(t, info["active"].(bool))
}

func TestManager_Build_Submenu_Recursive_Good(t *testing.T) {
	m, p := newTestManager()
	require.NoError(t, m.Setup("Core", "Core"))

	items := []TrayMenuItem{
		{
			Label: "Parent",
			Submenu: []TrayMenuItem{
				{Label: "Child 1"},
				{Label: "Child 2"},
			},
		},
	}

	require.NoError(t, m.SetMenu(items))
	require.Len(t, p.menus, 1)

	menu := p.menus[0]
	require.Len(t, menu.items, 1)
	assert.Equal(t, "Parent", menu.items[0])
	require.Len(t, menu.subs, 1)
	require.Len(t, menu.subs[0].items, 2)
	assert.Equal(t, "Child 1", menu.subs[0].items[0])
	assert.Equal(t, "Child 2", menu.subs[0].items[1])
}
