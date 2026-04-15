package systray

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSystrayService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "systray")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
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
	r := taskRun(c, "systray.setIcon", TaskSetTrayIcon{Data: icon})
	require.True(t, r.OK)
}

func TestTaskSetTrayMenu_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)

	require.NoError(t, svc.manager.Setup("Test", "Test"))

	items := []TrayMenuItem{
		{Label: "Open", ActionID: "open"},
		{Type: "separator"},
		{Label: "Quit", ActionID: "quit"},
	}
	r := taskRun(c, "systray.setMenu", TaskSetTrayMenu{Items: items})
	require.True(t, r.OK)
}

func TestTaskSetTrayIcon_Bad(t *testing.T) {
	// No systray service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("systray.setIcon").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}
