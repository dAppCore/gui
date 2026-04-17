package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScreen_Rect_Good(t *testing.T) {
	rect := Rect{X: 10, Y: 20, Width: 300, Height: 200}

	assert.Equal(t, Point{X: 10, Y: 20}, rect.Origin())
	assert.Equal(t, Point{X: 310, Y: 220}, rect.Corner())
	assert.False(t, rect.IsEmpty())
	assert.True(t, rect.Contains(Point{X: 10, Y: 20}))
	assert.False(t, rect.Contains(Point{X: 310, Y: 220}))
	assert.Equal(t, Size{Width: 300, Height: 200}, rect.RectSize())
}

func TestScreen_Rect_Bad(t *testing.T) {
	rect := Rect{}

	assert.True(t, rect.IsEmpty())
	assert.False(t, rect.Contains(Point{}))
}

func TestScreen_Rect_Ugly(t *testing.T) {
	rect := Rect{X: -10, Y: -5, Width: 10, Height: 5}

	assert.False(t, rect.IsEmpty())
	assert.True(t, rect.Contains(Point{X: -10, Y: -5}))
}

func TestScreenManager_SetScreens_Good(t *testing.T) {
	manager := &ScreenManager{}
	primary := &Screen{ID: "1", IsPrimary: true}
	secondary := &Screen{ID: "2"}

	manager.SetScreens([]*Screen{primary, secondary})

	require.Same(t, primary, manager.GetPrimary())
	require.Same(t, primary, manager.GetCurrent())
	assert.Equal(t, []*Screen{primary, secondary}, manager.GetAll())
}

func TestScreenManager_SetScreens_Bad(t *testing.T) {
	manager := &ScreenManager{}

	manager.SetScreens(nil)

	assert.Nil(t, manager.GetPrimary())
	assert.Nil(t, manager.GetCurrent())
	assert.Empty(t, manager.GetAll())
}

func TestScreenManager_SetScreens_Ugly(t *testing.T) {
	manager := &ScreenManager{}
	primary := &Screen{ID: "1", IsPrimary: true}
	current := &Screen{ID: "current"}

	manager.SetCurrent(current)
	manager.SetScreens([]*Screen{primary})

	require.Same(t, primary, manager.GetPrimary())
	require.Same(t, current, manager.GetCurrent())
}

func TestScreenManager_SetScreens_CopiesInput(t *testing.T) {
	manager := &ScreenManager{}
	screens := []*Screen{
		{ID: "1", IsPrimary: true},
		{ID: "2"},
	}

	manager.SetScreens(screens)
	screens[0] = &Screen{ID: "mutated"}

	require.Len(t, manager.GetAll(), 2)
	require.Equal(t, "1", manager.GetPrimary().ID)
	require.Equal(t, "1", manager.GetAll()[0].ID)
}
