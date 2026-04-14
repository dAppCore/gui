package contextmenu

import core "dappco.re/go/core"

var ErrorMenuNotFound = core.E("contextmenu", "menu not found", nil)

// QueryGet returns a named context menu definition. Result: *ContextMenuDef (nil if not found)
//
//	result := c.QUERY(contextmenu.QueryGet{Name: "editor"})
type QueryGet struct {
	Name string `json:"name"`
}

// QueryList returns all registered context menus. Result: map[string]ContextMenuDef
//
//	result := c.QUERY(contextmenu.QueryList{})
type QueryList struct{}

// QueryGetAll returns all context menus as a slice. Result: []ContextMenuDef
//
//	result := c.QUERY(contextmenu.QueryGetAll{})
type QueryGetAll struct{}

// TaskAdd registers a named context menu. Replaces if already exists.
//
//	c.PERFORM(contextmenu.TaskAdd{Name: "editor", Menu: def})
type TaskAdd struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskRemove unregisters a context menu by name. Error: ErrorMenuNotFound if missing.
//
//	c.PERFORM(contextmenu.TaskRemove{Name: "editor"})
type TaskRemove struct {
	Name string `json:"name"`
}

// TaskUpdate replaces a context menu definition. Error: ErrorMenuNotFound if missing.
//
//	c.PERFORM(contextmenu.TaskUpdate{Name: "editor", Menu: newDef})
type TaskUpdate struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskDestroy removes a context menu and releases resources.
//
//	c.PERFORM(contextmenu.TaskDestroy{Name: "editor"})
type TaskDestroy struct {
	Name string `json:"name"`
}

// ActionItemClicked is broadcast when a context menu item is clicked.
type ActionItemClicked struct {
	MenuName string `json:"menuName"`
	ActionID string `json:"actionId"`
	Data     string `json:"data,omitempty"`
}
