package application

import (
	core "dappco.re/go"
)

func TestScreen_Rect_Good(t *core.T) {
	rect := Rect{X: 10, Y: 20, Width: 300, Height: 200}

	core.AssertEqual(t, Point{X: 10, Y: 20}, rect.Origin())
	core.AssertEqual(t, Point{X: 310, Y: 220}, rect.Corner())
	core.AssertFalse(t, rect.IsEmpty())
	core.AssertTrue(t, rect.Contains(Point{X: 10, Y: 20}))
	core.AssertFalse(t, rect.Contains(Point{X: 310, Y: 220}))
	core.AssertEqual(t, Size{Width: 300, Height: 200}, rect.RectSize())
}

func TestScreen_Rect_Bad(t *core.T) {
	rect := Rect{}

	core.AssertTrue(t, rect.IsEmpty())
	core.AssertFalse(t, rect.Contains(Point{}))
}

func TestScreen_Rect_Ugly(t *core.T) {
	rect := Rect{X: -10, Y: -5, Width: 10, Height: 5}

	core.AssertFalse(t, rect.IsEmpty())
	core.AssertTrue(t, rect.Contains(Point{X: -10, Y: -5}))
}

func TestScreenManager_SetScreens_Good(t *core.T) {
	manager := &ScreenManager{}
	primary := &Screen{ID: "1", IsPrimary: true}
	secondary := &Screen{ID: "2"}

	manager.SetScreens([]*Screen{primary, secondary})

	core.AssertSame(t, primary, manager.GetPrimary())
	core.AssertSame(t, primary, manager.GetCurrent())
	core.AssertEqual(t, []*Screen{primary, secondary}, manager.GetAll())
}

func TestScreenManager_SetScreens_Bad(t *core.T) {
	manager := &ScreenManager{}

	manager.SetScreens(nil)

	core.AssertNil(t, manager.GetPrimary())
	core.AssertNil(t, manager.GetCurrent())
	core.AssertEmpty(t, manager.GetAll())
}

func TestScreenManager_SetScreens_Ugly(t *core.T) {
	manager := &ScreenManager{}
	primary := &Screen{ID: "1", IsPrimary: true}
	current := &Screen{ID: "current"}

	manager.SetCurrent(current)
	manager.SetScreens([]*Screen{primary})

	core.AssertSame(t, primary, manager.GetPrimary())
	core.AssertSame(t, current, manager.GetCurrent())
}

func TestScreenManager_SetScreens_CopiesInput(t *core.T) {
	manager := &ScreenManager{}
	screens := []*Screen{
		{ID: "1", IsPrimary: true},
		{ID: "2"},
	}

	manager.SetScreens(screens)
	screens[0] = &Screen{ID: "mutated"}

	core.AssertLen(t, manager.GetAll(), 2)
	core.AssertEqual(t, "1", manager.GetPrimary().ID)
	core.AssertEqual(t, "1", manager.GetAll()[0].ID)
}
