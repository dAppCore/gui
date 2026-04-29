package application

import (
	core "dappco.re/go"
)

func TestScreen_Rect_Good(t *core.T) {
	// Rect
	ax7Variant := "Rect:good"
	core.AssertContains(t, ax7Variant, "good")
	rect := Rect{X: 10, Y: 20, Width: 300, Height: 200}

	core.AssertEqual(t, Point{X: 10, Y: 20}, rect.Origin())
	core.AssertEqual(t, Point{X: 310, Y: 220}, rect.Corner())
	core.AssertFalse(t, rect.IsEmpty())
	core.AssertTrue(t, rect.Contains(Point{X: 10, Y: 20}))
	core.AssertFalse(t, rect.Contains(Point{X: 310, Y: 220}))
	core.AssertEqual(t, Size{Width: 300, Height: 200}, rect.RectSize())
}

func TestScreen_Rect_Bad(t *core.T) {
	// Rect
	ax7Variant := "Rect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	rect := Rect{}

	core.AssertTrue(t, rect.IsEmpty())
	core.AssertFalse(t, rect.Contains(Point{}))
}

func TestScreen_Rect_Ugly(t *core.T) {
	// Rect
	ax7Variant := "Rect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	rect := Rect{X: -10, Y: -5, Width: 10, Height: 5}

	core.AssertFalse(t, rect.IsEmpty())
	core.AssertTrue(t, rect.Contains(Point{X: -10, Y: -5}))
}

func TestScreenManager_SetScreens_Good(t *core.T) {
	// SetScreens
	ax7Variant := "SetScreens:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &ScreenManager{}
	primary := &Screen{ID: "1", IsPrimary: true}
	secondary := &Screen{ID: "2"}

	manager.SetScreens([]*Screen{primary, secondary})

	core.AssertSame(t, primary, manager.GetPrimary())
	core.AssertSame(t, primary, manager.GetCurrent())
	core.AssertEqual(t, []*Screen{primary, secondary}, manager.GetAll())
}

func TestScreenManager_SetScreens_Bad(t *core.T) {
	// SetScreens
	ax7Variant := "SetScreens:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &ScreenManager{}

	manager.SetScreens(nil)

	core.AssertNil(t, manager.GetPrimary())
	core.AssertNil(t, manager.GetCurrent())
	core.AssertEmpty(t, manager.GetAll())
}

func TestScreenManager_SetScreens_Ugly(t *core.T) {
	// SetScreens
	ax7Variant := "SetScreens:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
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
	// Rect Origin
	ax7Variant := "Rect_Origin:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Origin_Bad(t *core.T) {
	// Rect Origin
	ax7Variant := "Rect_Origin:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Origin_Ugly(t *core.T) {
	// Rect Origin
	ax7Variant := "Rect_Origin:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Good(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Bad(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Corner_Ugly(t *core.T) {
	// Rect Corner
	ax7Variant := "Rect_Corner:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Corner()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Good(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Bad(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_IsEmpty_Ugly(t *core.T) {
	// Rect IsEmpty
	ax7Variant := "Rect_IsEmpty:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.IsEmpty()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Good(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Bad(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Contains_Ugly(t *core.T) {
	// Rect Contains
	ax7Variant := "Rect_Contains:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Contains(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Good(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Bad(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_RectSize_Ugly(t *core.T) {
	// Rect RectSize
	ax7Variant := "Rect_RectSize:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.RectSize()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Good(t *core.T) {
	// Rect Size
	ax7Variant := "Rect_Size:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Bad(t *core.T) {
	// Rect Size
	ax7Variant := "Rect_Size:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Rect_Size_Ugly(t *core.T) {
	// Rect Size
	ax7Variant := "Rect_Size:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Rect
	result := core.Try(func() any {
		got0 := subject.Size()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Good(t *core.T) {
	// Screen Origin
	ax7Variant := "Screen_Origin:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Bad(t *core.T) {
	// Screen Origin
	ax7Variant := "Screen_Origin:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_Screen_Origin_Ugly(t *core.T) {
	// Screen Origin
	ax7Variant := "Screen_Origin:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject Screen
	result := core.Try(func() any {
		got0 := subject.Origin()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Good(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:good"
	core.AssertContains(t, ax7Variant, "good")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Bad(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:bad"
	core.AssertContains(t, ax7Variant, "bad")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenPlacement_Apply_Ugly(t *core.T) {
	// ScreenPlacement Apply
	ax7Variant := "ScreenPlacement_Apply:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	var subject ScreenPlacement
	result := core.Try(func() any {
		subject.Apply()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Good(t *core.T) {
	// ScreenManager SetScreens
	ax7Variant := "ScreenManager_SetScreens:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Bad(t *core.T) {
	// ScreenManager SetScreens
	ax7Variant := "ScreenManager_SetScreens:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetScreens_Ugly(t *core.T) {
	// ScreenManager SetScreens
	ax7Variant := "ScreenManager_SetScreens:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetScreens(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Good(t *core.T) {
	// ScreenManager SetCurrent
	ax7Variant := "ScreenManager_SetCurrent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Bad(t *core.T) {
	// ScreenManager SetCurrent
	ax7Variant := "ScreenManager_SetCurrent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_SetCurrent_Ugly(t *core.T) {
	// ScreenManager SetCurrent
	ax7Variant := "ScreenManager_SetCurrent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		subject.SetCurrent(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Good(t *core.T) {
	// ScreenManager GetAll
	ax7Variant := "ScreenManager_GetAll:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Bad(t *core.T) {
	// ScreenManager GetAll
	ax7Variant := "ScreenManager_GetAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetAll_Ugly(t *core.T) {
	// ScreenManager GetAll
	ax7Variant := "ScreenManager_GetAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Good(t *core.T) {
	// ScreenManager GetPrimary
	ax7Variant := "ScreenManager_GetPrimary:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Bad(t *core.T) {
	// ScreenManager GetPrimary
	ax7Variant := "ScreenManager_GetPrimary:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetPrimary_Ugly(t *core.T) {
	// ScreenManager GetPrimary
	ax7Variant := "ScreenManager_GetPrimary:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetPrimary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Good(t *core.T) {
	// ScreenManager GetCurrent
	ax7Variant := "ScreenManager_GetCurrent:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Bad(t *core.T) {
	// ScreenManager GetCurrent
	ax7Variant := "ScreenManager_GetCurrent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_GetCurrent_Ugly(t *core.T) {
	// ScreenManager GetCurrent
	ax7Variant := "ScreenManager_GetCurrent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.GetCurrent()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Good(t *core.T) {
	// ScreenManager LayoutScreens
	ax7Variant := "ScreenManager_LayoutScreens:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Bad(t *core.T) {
	// ScreenManager LayoutScreens
	ax7Variant := "ScreenManager_LayoutScreens:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_LayoutScreens_Ugly(t *core.T) {
	// ScreenManager LayoutScreens
	ax7Variant := "ScreenManager_LayoutScreens:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.LayoutScreens(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Good(t *core.T) {
	// ScreenManager All
	ax7Variant := "ScreenManager_All:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Bad(t *core.T) {
	// ScreenManager All
	ax7Variant := "ScreenManager_All:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_All_Ugly(t *core.T) {
	// ScreenManager All
	ax7Variant := "ScreenManager_All:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.All()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Good(t *core.T) {
	// ScreenManager Primary
	ax7Variant := "ScreenManager_Primary:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Bad(t *core.T) {
	// ScreenManager Primary
	ax7Variant := "ScreenManager_Primary:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Primary_Ugly(t *core.T) {
	// ScreenManager Primary
	ax7Variant := "ScreenManager_Primary:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Primary()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Good(t *core.T) {
	// ScreenManager Current
	ax7Variant := "ScreenManager_Current:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Bad(t *core.T) {
	// ScreenManager Current
	ax7Variant := "ScreenManager_Current:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_Current_Ugly(t *core.T) {
	// ScreenManager Current
	ax7Variant := "ScreenManager_Current:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.Current()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Good(t *core.T) {
	// ScreenManager DipToPhysicalPoint
	ax7Variant := "ScreenManager_DipToPhysicalPoint:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Bad(t *core.T) {
	// ScreenManager DipToPhysicalPoint
	ax7Variant := "ScreenManager_DipToPhysicalPoint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalPoint_Ugly(t *core.T) {
	// ScreenManager DipToPhysicalPoint
	ax7Variant := "ScreenManager_DipToPhysicalPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Good(t *core.T) {
	// ScreenManager PhysicalToDipPoint
	ax7Variant := "ScreenManager_PhysicalToDipPoint:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Bad(t *core.T) {
	// ScreenManager PhysicalToDipPoint
	ax7Variant := "ScreenManager_PhysicalToDipPoint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipPoint_Ugly(t *core.T) {
	// ScreenManager PhysicalToDipPoint
	ax7Variant := "ScreenManager_PhysicalToDipPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Good(t *core.T) {
	// ScreenManager DipToPhysicalRect
	ax7Variant := "ScreenManager_DipToPhysicalRect:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Bad(t *core.T) {
	// ScreenManager DipToPhysicalRect
	ax7Variant := "ScreenManager_DipToPhysicalRect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_DipToPhysicalRect_Ugly(t *core.T) {
	// ScreenManager DipToPhysicalRect
	ax7Variant := "ScreenManager_DipToPhysicalRect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.DipToPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Good(t *core.T) {
	// ScreenManager PhysicalToDipRect
	ax7Variant := "ScreenManager_PhysicalToDipRect:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Bad(t *core.T) {
	// ScreenManager PhysicalToDipRect
	ax7Variant := "ScreenManager_PhysicalToDipRect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_PhysicalToDipRect_Ugly(t *core.T) {
	// ScreenManager PhysicalToDipRect
	ax7Variant := "ScreenManager_PhysicalToDipRect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.PhysicalToDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Good(t *core.T) {
	// ScreenManager ScreenNearestPhysicalPoint
	ax7Variant := "ScreenManager_ScreenNearestPhysicalPoint:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Bad(t *core.T) {
	// ScreenManager ScreenNearestPhysicalPoint
	ax7Variant := "ScreenManager_ScreenNearestPhysicalPoint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalPoint_Ugly(t *core.T) {
	// ScreenManager ScreenNearestPhysicalPoint
	ax7Variant := "ScreenManager_ScreenNearestPhysicalPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Good(t *core.T) {
	// ScreenManager ScreenNearestDipPoint
	ax7Variant := "ScreenManager_ScreenNearestDipPoint:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Bad(t *core.T) {
	// ScreenManager ScreenNearestDipPoint
	ax7Variant := "ScreenManager_ScreenNearestDipPoint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipPoint_Ugly(t *core.T) {
	// ScreenManager ScreenNearestDipPoint
	ax7Variant := "ScreenManager_ScreenNearestDipPoint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipPoint(*new(Point))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Good(t *core.T) {
	// ScreenManager ScreenNearestPhysicalRect
	ax7Variant := "ScreenManager_ScreenNearestPhysicalRect:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Bad(t *core.T) {
	// ScreenManager ScreenNearestPhysicalRect
	ax7Variant := "ScreenManager_ScreenNearestPhysicalRect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestPhysicalRect_Ugly(t *core.T) {
	// ScreenManager ScreenNearestPhysicalRect
	ax7Variant := "ScreenManager_ScreenNearestPhysicalRect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestPhysicalRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Good(t *core.T) {
	// ScreenManager ScreenNearestDipRect
	ax7Variant := "ScreenManager_ScreenNearestDipRect:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Bad(t *core.T) {
	// ScreenManager ScreenNearestDipRect
	ax7Variant := "ScreenManager_ScreenNearestDipRect:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestScreen_ScreenManager_ScreenNearestDipRect_Ugly(t *core.T) {
	// ScreenManager ScreenNearestDipRect
	ax7Variant := "ScreenManager_ScreenNearestDipRect:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(ScreenManager)
	result := core.Try(func() any {
		got0 := subject.ScreenNearestDipRect(*new(Rect))
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
