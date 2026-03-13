package window

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestWindowService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestWindowService(t)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.manager)
}

func TestTaskOpenWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.PERFORM(TaskOpenWindow{
		Opts: []WindowOption{WithName("test"), WithURL("/")},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestTaskOpenWindow_Bad(t *testing.T) {
	// No window service registered — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskOpenWindow{})
	assert.False(t, handled)
}

func TestQueryWindowList_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("a")}})
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("b")}})

	result, handled, err := c.QUERY(QueryWindowList{})
	require.NoError(t, err)
	assert.True(t, handled)
	list := result.([]WindowInfo)
	assert.Len(t, list, 2)
}

func TestQueryWindowByName_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	result, handled, err := c.QUERY(QueryWindowByName{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(*WindowInfo)
	assert.Equal(t, "test", info.Name)
}

func TestQueryWindowByName_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	result, handled, err := c.QUERY(QueryWindowByName{Name: "nonexistent"})
	require.NoError(t, err)
	assert.True(t, handled) // handled=true, result is nil (not found)
	assert.Nil(t, result)
}

func TestTaskCloseWindow_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify window is removed
	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	assert.Nil(t, result)
}

func TestTaskCloseWindow_Bad(t *testing.T) {
	_, c := newTestWindowService(t)
	_, handled, err := c.PERFORM(TaskCloseWindow{Name: "nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetPosition_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetPosition{Name: "test", X: 100, Y: 200})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 100, info.X)
	assert.Equal(t, 200, info.Y)
}

func TestTaskSetSize_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskSetSize{Name: "test", W: 800, H: 600})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.Equal(t, 800, info.Width)
	assert.Equal(t, 600, info.Height)
}

func TestTaskMaximise_Good(t *testing.T) {
	_, c := newTestWindowService(t)
	_, _, _ = c.PERFORM(TaskOpenWindow{Opts: []WindowOption{WithName("test")}})

	_, handled, err := c.PERFORM(TaskMaximise{Name: "test"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, _ := c.QUERY(QueryWindowByName{Name: "test"})
	info := result.(*WindowInfo)
	assert.True(t, info.Maximized)
}
