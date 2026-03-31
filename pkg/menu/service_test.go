package menu

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMenuService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "menu")
	return svc, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestMenuService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestTaskSetAppMenu_Good(t *testing.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{
		{Label: "File", Children: []MenuItem{
			{Label: "New"},
			{Type: "separator"},
			{Label: "Quit"},
		}},
	}
	r := taskRun(c, "menu.setAppMenu", TaskSetAppMenu{Items: items})
	require.True(t, r.OK)
}

func TestQueryGetAppMenu_Good(t *testing.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{{Label: "File"}, {Label: "Edit"}}
	taskRun(c, "menu.setAppMenu", TaskSetAppMenu{Items: items})

	r := c.QUERY(QueryGetAppMenu{})
	require.True(t, r.OK)
	menuItems := r.Value.([]MenuItem)
	assert.Len(t, menuItems, 2)
	assert.Equal(t, "File", menuItems[0].Label)
}

func TestTaskSetAppMenu_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("menu.setAppMenu").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}
