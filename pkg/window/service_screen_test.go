package window

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/screen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockScreenPlatform struct {
	screens []screen.Screen
}

func (m *mockScreenPlatform) GetAll() []screen.Screen { return m.screens }

func (m *mockScreenPlatform) GetPrimary() *screen.Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}

func newTestWindowServiceWithScreen(t *testing.T, screens []screen.Screen) (*Service, *core.Core) {
	t.Helper()

	c, err := core.New(
		core.WithService(screen.Register(&mockScreenPlatform{screens: screens})),
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	svc := core.MustServiceFor[*Service](c, "window")
	return svc, c
}

func TestTaskTileWindows_Good_UsesPrimaryScreenSize(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		},
	})

	_, _, err := c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}})
	require.NoError(t, err)
	_, _, err = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}})
	require.NoError(t, err)

	_, handled, err := c.PERFORM(TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, err := c.QUERY(QueryWindowByName{Name: "left"})
	require.NoError(t, err)
	left := result.(*WindowInfo)
	assert.Equal(t, 0, left.X)
	assert.Equal(t, 1000, left.Width)
	assert.Equal(t, 1000, left.Height)

	result, _, err = c.QUERY(QueryWindowByName{Name: "right"})
	require.NoError(t, err)
	right := result.(*WindowInfo)
	assert.Equal(t, 1000, right.X)
	assert.Equal(t, 1000, right.Width)
	assert.Equal(t, 1000, right.Height)
}

func TestTaskSnapWindow_Good_UsesPrimaryScreenSize(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		},
	})

	_, _, err := c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("snap"), WithSize(400, 300)}})
	require.NoError(t, err)

	_, handled, err := c.PERFORM(TaskSnapWindow{Name: "snap", Position: "left"})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, err := c.QUERY(QueryWindowByName{Name: "snap"})
	require.NoError(t, err)
	info := result.(*WindowInfo)
	assert.Equal(t, 0, info.X)
	assert.Equal(t, 0, info.Y)
	assert.Equal(t, 1000, info.Width)
	assert.Equal(t, 1000, info.Height)
}

func TestTaskTileWindows_Good_UsesPrimaryWorkAreaOrigin(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{
		{
			ID: "1", Name: "Primary", IsPrimary: true,
			Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
			WorkArea: screen.Rect{X: 100, Y: 50, Width: 2000, Height: 1000},
		},
	})

	_, _, err := c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}})
	require.NoError(t, err)
	_, _, err = c.PERFORM(TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}})
	require.NoError(t, err)

	_, handled, err := c.PERFORM(TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	require.NoError(t, err)
	assert.True(t, handled)

	result, _, err := c.QUERY(QueryWindowByName{Name: "left"})
	require.NoError(t, err)
	left := result.(*WindowInfo)
	assert.Equal(t, 100, left.X)
	assert.Equal(t, 50, left.Y)
	assert.Equal(t, 1000, left.Width)
	assert.Equal(t, 1000, left.Height)

	result, _, err = c.QUERY(QueryWindowByName{Name: "right"})
	require.NoError(t, err)
	right := result.(*WindowInfo)
	assert.Equal(t, 1100, right.X)
	assert.Equal(t, 50, right.Y)
	assert.Equal(t, 1000, right.Width)
	assert.Equal(t, 1000, right.Height)
}
