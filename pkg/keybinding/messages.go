// pkg/keybinding/messages.go
package keybinding

import "errors"

// ErrAlreadyRegistered is returned when attempting to add a binding
// that already exists. Callers must TaskRemove first to rebind.
var ErrAlreadyRegistered = errors.New("keybinding: accelerator already registered")

// BindingInfo describes a registered keyboard shortcut.
type BindingInfo struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// --- Queries ---

// QueryList returns all registered bindings. Result: []BindingInfo
type QueryList struct{}

// --- Tasks ---

// TaskAdd registers a new keyboard shortcut. Result: nil
// Returns ErrAlreadyRegistered if the accelerator is already bound.
type TaskAdd struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

// TaskRemove unregisters a keyboard shortcut. Result: nil
type TaskRemove struct {
	Accelerator string `json:"accelerator"`
}

// --- Actions ---

// ActionTriggered is broadcast when a registered shortcut is activated.
type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
