//go:build compliance

package application

import core "dappco.re/go"

func ExampleEventManager_Once() {
	core.Println("EventManager_Once")
	// Output:
	// EventManager_Once
}

func ExampleEventManager_EmitEvent() {
	core.Println("EventManager_EmitEvent")
	// Output:
	// EventManager_EmitEvent
}

func ExampleEventManager_Reset() {
	core.Println("EventManager_Reset")
	// Output:
	// EventManager_Reset
}

func ExampleEventManager_RegisterApplicationEventHook() {
	core.Println("EventManager_RegisterApplicationEventHook")
	// Output:
	// EventManager_RegisterApplicationEventHook
}
