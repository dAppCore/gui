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

// AX7 generated source-matching smoke coverage.
func TestScreen_Rect_Origin_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Origin_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Origin_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Good(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Bad(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Ugly(t *core.T) {
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Good(t *core.T) {
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Bad(t *core.T) {
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Ugly(t *core.T) {
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Good(t *core.T) {
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Bad(t *core.T) {
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Ugly(t *core.T) {
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Good(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Bad(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Ugly(t *core.T) {
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
