package menu

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMenuService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "menu")
	return svc, c
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
	_, handled, err := c.PERFORM(TaskSetAppMenu{Items: items})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestQueryGetAppMenu_Good(t *testing.T) {
	_, c := newTestMenuService(t)

	items := []MenuItem{{Label: "File"}, {Label: "Edit"}}
	_, _, _ = c.PERFORM(TaskSetAppMenu{Items: items})

	result, handled, err := c.QUERY(QueryGetAppMenu{})
	require.NoError(t, err)
	assert.True(t, handled)
	menuItems := result.([]MenuItem)
	assert.Len(t, menuItems, 2)
	assert.Equal(t, "File", menuItems[0].Label)
}

func TestTaskSetAppMenu_Bad(t *testing.T) {
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetAppMenu{})
	assert.False(t, handled)
}
