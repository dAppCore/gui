package window

import (
	"context"
	"testing"

	core "dappco.re/go/core"
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

func (m *mockScreenPlatform) GetCurrent() *screen.Screen {
	return m.GetPrimary()
}

func newTestWindowServiceWithScreen(t *testing.T, screens []screen.Screen) (*Service, *core.Core) {
	t.Helper()

	c := core.New(
		core.WithService(screen.Register(&mockScreenPlatform{screens: screens})),
		core.WithService(Register(newMockPlatform())),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

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

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}}).OK)
	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}}).OK)

	r := taskRun(c, "window.tileWindows", TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "left"})
	require.True(t, r2.OK)
	left := r2.Value.(*WindowInfo)
	assert.Equal(t, 0, left.X)
	assert.Equal(t, 1000, left.Width)
	assert.Equal(t, 1000, left.Height)

	r3 := c.QUERY(QueryWindowByName{Name: "right"})
	require.True(t, r3.OK)
	right := r3.Value.(*WindowInfo)
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

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("snap"), WithSize(400, 300)}}).OK)

	r := taskRun(c, "window.snapWindow", TaskSnapWindow{Name: "snap", Position: "left"})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "snap"})
	require.True(t, r2.OK)
	info := r2.Value.(*WindowInfo)
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

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("left"), WithSize(400, 400)}}).OK)
	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("right"), WithSize(400, 400)}}).OK)

	r := taskRun(c, "window.tileWindows", TaskTileWindows{Mode: "left-right", Windows: []string{"left", "right"}})
	require.True(t, r.OK)

	r2 := c.QUERY(QueryWindowByName{Name: "left"})
	require.True(t, r2.OK)
	left := r2.Value.(*WindowInfo)
	assert.Equal(t, 100, left.X)
	assert.Equal(t, 50, left.Y)
	assert.Equal(t, 1000, left.Width)
	assert.Equal(t, 1000, left.Height)

	r3 := c.QUERY(QueryWindowByName{Name: "right"})
	require.True(t, r3.OK)
	right := r3.Value.(*WindowInfo)
	assert.Equal(t, 1100, right.X)
	assert.Equal(t, 50, right.Y)
	assert.Equal(t, 1000, right.Width)
	assert.Equal(t, 1000, right.Height)
}
