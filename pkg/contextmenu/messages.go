// pkg/contextmenu/messages.go
package contextmenu

import "errors"

// ErrMenuNotFound is returned when attempting to remove or get a menu
// that does not exist in the registry.
var ErrMenuNotFound = errors.New("contextmenu: menu not found")

// --- Queries ---

// QueryGet returns a single context menu by name. Result: *ContextMenuDef (nil if not found)
type QueryGet struct {
	Name string `json:"name"`
}

// QueryList returns all registered context menus. Result: map[string]ContextMenuDef
type QueryList struct{}

// --- Tasks ---

// TaskAdd registers a context menu. Result: nil
// If a menu with the same name already exists it is replaced (remove + re-add).
type TaskAdd struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

// TaskRemove unregisters a context menu. Result: nil
// Returns ErrMenuNotFound if the menu does not exist.
type TaskRemove struct {
	Name string `json:"name"`
}

// --- Actions ---

// ActionItemClicked is broadcast when a context menu item is clicked.
// The Data field is populated from the CSS --custom-contextmenu-data property
// on the element that triggered the context menu.
type ActionItemClicked struct {
	MenuName string `json:"menuName"`
	ActionID string `json:"actionId"`
	Data     string `json:"data,omitempty"`
}
