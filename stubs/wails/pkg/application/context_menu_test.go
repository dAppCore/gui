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
