package systray

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSystrayService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "systray")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestSystrayService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestTaskSetTrayIcon_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)

	// Setup tray first (normally done via config in OnStartup)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	icon := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	_, handled, err := c.PERFORM(TaskSetTrayIcon{Data: icon})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetTooltip_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetTooltip{Tooltip: "Updated"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetLabel_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetLabel{Label: "Updated"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetTrayMenu_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)

	require.NoError(t, svc.manager.Setup("Test", "Test"))

	items := []TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}
	_, handled, err := c.PERFORM(TaskSetTrayMenu{Items: items})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetTrayMenu_Submenu_Good(t *testing.T) {
	p := newMockPlatform()
	c, err := core.New(
		core.WithService(Register(p)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	svc := core.MustServiceFor[*Service](c, "systray")
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetTrayMenu{Items: []TrayMenuItem{
		{
			Label: "File",
			Submenu: []TrayMenuItem{
				{Label: "Open", ActionID: "open"},
			},
		},
	}})
	require.NoError(t, err)
	assert.True(t, handled)
	require.Len(t, p.trays, 1)
	require.NotEmpty(t, p.menus)
	require.Len(t, p.menus[0].submenus, 1)
}

func TestTaskSetTrayIcon_Bad(t *testing.T) {
	// No systray service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetTrayIcon{Data: nil})
	assert.False(t, handled)
}

func TestTaskShowMessage_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))
	_, handled, err := c.PERFORM(TaskShowMessage{Title: "Hello", Message: "World"})
	require.NoError(t, err)
	assert.True(t, handled)
}
