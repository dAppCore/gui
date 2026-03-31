package keybinding

import "errors"

var ErrorAlreadyRegistered = errors.New("keybinding: accelerator already registered")

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

// ActionTriggered is broadcast when a registered key binding fires.
type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
