package contextmenu

import "errors"

var ErrorMenuNotFound = errors.New("contextmenu: menu not found")

// --- Queries ---

type QueryGet struct {
	Name string `json:"name"`
}

type QueryList struct{}

// --- Tasks ---

type TaskAdd struct {
	Name string         `json:"name"`
	Menu ContextMenuDef `json:"menu"`
}

type TaskRemove struct {
	Name string `json:"name"`
}

type ActionItemClicked struct {
	MenuName string `json:"menuName"`
	ActionID string `json:"actionId"`
	Data     string `json:"data,omitempty"`
}
