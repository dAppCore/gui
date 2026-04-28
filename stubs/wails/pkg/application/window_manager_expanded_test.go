package application

import (
	core "dappco.re/go"
)

func TestWindowManagerExpanded_Get_Good(t *core.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertSame(t, first, manager.Get("first"))
	core.AssertSame(t, second, manager.Get("second"))
	core.AssertLen(t, manager.GetAll(), 2)
}

func TestWindowManagerExpanded_Get_Bad(t *core.T) {
	manager := &WindowManager{}

	core.AssertNil(t, manager.Get("missing"))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestWindowManagerExpanded_Get_Ugly(t *core.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "dup"})

	core.AssertSame(t, first, manager.Get("dup"))
	core.AssertNotEqual(t, core.Sprintf("%p", first), core.Sprintf("%p", second))
}

func TestWindowManagerExpanded_GetByID_Good(t *core.T) {
	manager := &WindowManager{}
	first := manager.NewWithOptions(WebviewWindowOptions{Name: "first"})
	second := manager.NewWithOptions(WebviewWindowOptions{Name: "second"})

	core.AssertSame(t, first, manager.GetByID(1))
	core.AssertSame(t, second, manager.GetByID(2))
}

func TestWindowManagerExpanded_GetByID_Bad(t *core.T) {
	manager := &WindowManager{}

	core.AssertNil(t, manager.GetByID(99))
	core.AssertNotEmpty(t, core.Sprintf("%T", manager))
}

func TestWindowManagerExpanded_GetByID_Ugly(t *core.T) {
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
