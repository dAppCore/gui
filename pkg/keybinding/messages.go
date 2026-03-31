package keybinding

import "errors"

var ErrorAlreadyRegistered = errors.New("keybinding: accelerator already registered")

type BindingInfo struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

type QueryList struct{}

type TaskAdd struct {
	Accelerator string `json:"accelerator"`
	Description string `json:"description"`
}

type TaskRemove struct {
	Accelerator string `json:"accelerator"`
}

type ActionTriggered struct {
	Accelerator string `json:"accelerator"`
}
