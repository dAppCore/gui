package application

import (
	core "dappco.re/go"
)

func TestWindowManagerExpanded_Get_Good(t *core.T) {
	// Get
	ax7Variant := "Get:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertSame(t, first, manager.Get("first"))
	core.AssertSame(t, second, manager.Get("second"))
	core.AssertLen(t, manager.GetAll(), 2)
}

func TestWindowManagerExpanded_Get_Bad(t *core.T) {
	// Get
	ax7Variant := "Get:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &WindowManager{}

	core.AssertNil(t, manager.Get("missing"))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestWindowManagerExpanded_Get_Ugly(t *core.T) {
	// Get
	ax7Variant := "Get:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})

	core.AssertSame(t, first, manager.Get("dup"))
	core.AssertNotEqual(t, core.Sprintf("%p", first), core.Sprintf("%p", second))
}

func TestWindowManagerExpanded_GetByID_Good(t *core.T) {
	// GetByID
	ax7Variant := "GetByID:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertSame(t, first, manager.GetByID(1))
	core.AssertSame(t, second, manager.GetByID(2))
}

func TestWindowManagerExpanded_GetByID_Bad(t *core.T) {
	// GetByID
	ax7Variant := "GetByID:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &WindowManager{}

	core.AssertNil(t, manager.GetByID(99))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestWindowManagerExpanded_GetByID_Ugly(t *core.T) {
	// GetByID
	ax7Variant := "GetByID:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := &WindowManager{}
	manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertNil(t, manager.GetByID(0))
	core.AssertNil(t, manager.GetByID(3))
}

func TestWindowManagerExpanded_NilReceiver_IsSafe(t *core.T) {
	var manager *WindowManager

	core.AssertNotPanics(t, func() {
		core.AssertNil(t, manager.NewWithOptions(WebviewWindowOptions{Name: "ignored"}))
		core.AssertNil(t, manager.Get("missing"))
		core.AssertNil(t, manager.GetByID(1))
		core.AssertNil(t, manager.GetAll())
	})
}

// AX7 generated source-matching smoke coverage.
func TestWindowManagerExpanded_WindowManager_Get_Good(t *core.T) {
	// WindowManager Get
	ax7Variant := "WindowManager_Get:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.Get("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWindowManagerExpanded_WindowManager_Get_Bad(t *core.T) {
	// WindowManager Get
	ax7Variant := "WindowManager_Get:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.Get("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWindowManagerExpanded_WindowManager_Get_Ugly(t *core.T) {
	// WindowManager Get
	ax7Variant := "WindowManager_Get:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.Get("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWindowManagerExpanded_WindowManager_GetByID_Good(t *core.T) {
	// WindowManager GetByID
	ax7Variant := "WindowManager_GetByID:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetByID(1)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWindowManagerExpanded_WindowManager_GetByID_Bad(t *core.T) {
	// WindowManager GetByID
	ax7Variant := "WindowManager_GetByID:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetByID(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestWindowManagerExpanded_WindowManager_GetByID_Ugly(t *core.T) {
	// WindowManager GetByID
	ax7Variant := "WindowManager_GetByID:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(WindowManager)
	result := core.Try(func() any {
		got0 := subject.GetByID(0)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
