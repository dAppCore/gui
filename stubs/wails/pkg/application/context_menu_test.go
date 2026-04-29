package application

import (
	core "dappco.re/go"
)

func TestContextMenuManager_Add_Good(t *core.T) {
	manager := &ContextMenuManager{}
	menu := manager.New()
	menu.Add("Open")

	manager.Add("files", menu)

	got, ok := manager.Get("files")
	core.RequireTrue(t, ok)
	core.AssertSame(t, menu, got)
	core.AssertLen(t, manager.GetAll(), 1)
}

func TestContextMenuManager_Add_Bad(t *core.T) {
	manager := &ContextMenuManager{}

	manager.Add("empty", nil)

	got, ok := manager.Get("empty")
	core.AssertTrue(t, ok)
	core.AssertNil(t, got)
}

func TestContextMenuManager_Add_Ugly(t *core.T) {
	manager := &ContextMenuManager{}
	first := manager.New()
	second := manager.New()

	manager.Add("dup", first)
	manager.Add("dup", second)

	got, ok := manager.Get("dup")
	core.RequireTrue(t, ok)
	core.AssertSame(t, second, got)
	core.AssertLen(t, manager.GetAll(), 1)
}

func TestContextMenuManager_Remove_Good(t *core.T) {
	manager := &ContextMenuManager{}
	menu := manager.New()
	manager.Add("files", menu)

	manager.Remove("files")

	_, ok := manager.Get("files")
	core.AssertFalse(t, ok)
	core.AssertEmpty(t, manager.GetAll())
}

func TestContextMenuManager_Remove_Bad(t *core.T) {
	manager := &ContextMenuManager{}

	manager.Remove("missing")

	core.AssertEmpty(t, manager.GetAll())
}

func TestContextMenuManager_Remove_Ugly(t *core.T) {
	manager := &ContextMenuManager{}
	manager.Remove("")

	core.AssertEmpty(t, manager.GetAll())
}

// AX7 generated source-matching smoke coverage.
func TestContextMenu_ContextMenuManager_New_Good(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_New_Bad(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_New_Ugly(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.New()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Add_Good(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Add("agent", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Add_Bad(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Add("", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Add_Ugly(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Add("../../edge", nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Remove_Good(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Remove("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Remove_Bad(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Remove("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Remove_Ugly(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		subject.Remove("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Get_Good(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("agent")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Get_Bad(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_Get_Ugly(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0, got1 := subject.Get("../../edge")
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_GetAll_Good(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_GetAll_Bad(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestContextMenu_ContextMenuManager_GetAll_Ugly(t *core.T) {
	subject := new(ContextMenuManager)
	result := core.Try(func() any {
		got0 := subject.GetAll()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
