package systray

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWindowHandle struct {
	name       string
	showCalled bool
	hideCalled bool
}

func (w *mockWindowHandle) Name() string              { return w.name }
func (w *mockWindowHandle) Show()                     { w.showCalled = true }
func (w *mockWindowHandle) Hide()                     { w.hideCalled = true }
func (w *mockWindowHandle) SetPosition(x, y int)      {}
func (w *mockWindowHandle) SetSize(width, height int) {}

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

	_, handled, err := c.PERFORM(TaskSetTrayTooltip{Tooltip: "Updated"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskSetLabel_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetTrayLabel{Label: "Updated"})
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

func TestTaskSetTrayTooltip_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetTrayTooltip{Tooltip: "Updated"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Updated", svc.manager.Tray().(*mockTray).tooltip)
}

func TestTaskSetTrayLabel_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskSetTrayLabel{Label: "CoreGUI"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "CoreGUI", svc.manager.Tray().(*mockTray).label)
}

func TestTaskShowMessage_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	_, handled, err := c.PERFORM(TaskShowMessage{Title: "Heads up", Message: "Background work finished"})
	require.NoError(t, err)
	assert.True(t, handled)
	tray := svc.manager.Tray().(*mockTray)
	assert.Equal(t, "Heads up", tray.lastMessageTitle)
	assert.Equal(t, "Background work finished", tray.lastMessageBody)
}

func TestTaskSetTrayIcon_Bad(t *testing.T) {
	// No systray service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetTrayIcon{Data: nil})
	assert.False(t, handled)
}

func TestTaskShowMessage_Smoke(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))
	_, handled, err := c.PERFORM(TaskShowMessage{Title: "Hello", Message: "World"})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestTaskShowHidePanel_Good(t *testing.T) {
	svc, c := newTestSystrayService(t)
	require.NoError(t, svc.manager.Setup("Test", "Test"))

	panel := &mockWindowHandle{name: "panel"}
	require.NoError(t, svc.manager.AttachWindow(panel))

	_, handled, err := c.PERFORM(TaskShowPanel{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, panel.showCalled)

	_, handled, err = c.PERFORM(TaskHidePanel{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, panel.hideCalled)
}
