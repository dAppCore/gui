package keybinding

import coreerr "forge.lthn.ai/core/go-log"

// ErrorAlreadyRegistered is returned by TaskAdd when the accelerator is already bound.
var ErrorAlreadyRegistered = coreerr.NewError("keybinding: accelerator already registered")

// ErrorNotRegistered is returned by TaskRemove and TaskProcess when the accelerator is unknown.
var ErrorNotRegistered = coreerr.NewError("keybinding: accelerator not registered")

// BindingInfo describes a registered global key binding.
type BindingInfo struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// QueryList returns all registered key bindings. Result: []BindingInfo
type QueryList struct{}

// TaskAdd registers a global key binding. Error: ErrorAlreadyRegistered if accelerator taken.
type TaskAdd struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// TaskRemove unregisters a global key binding by accelerator.
type TaskRemove struct {
	Accelerator string `json:"accelerator"`
}

// TaskProcess programmatically triggers a registered key binding as if the user pressed it.
// Error: ErrorNotRegistered if the accelerator has not been registered.
// _, _, err := c.PERFORM(TaskProcess{Accelerator: "Ctrl+S"})
type TaskProcess struct {
	Accelerator string `json:"accelerator"`
}

// ActionTriggered is broadcast when a registered key binding fires.
type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
