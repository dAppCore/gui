//go:build compliance

package application

import core "dappco.re/go"

func ExampleEventManager_Emit() {
	core.Println("EventManager_Emit")
	// Output:
	// EventManager_Emit
}

func ExampleEventManager_EmitEvent() {
	core.Println("EventManager_EmitEvent")
	// Output:
	// EventManager_EmitEvent
}

func ExampleEventManager_On() {
	core.Println("EventManager_On")
	// Output:
	// EventManager_On
}

func ExampleEventManager_Off() {
	core.Println("EventManager_Off")
	// Output:
	// EventManager_Off
}

func ExampleEventManager_OnMultiple() {
	core.Println("EventManager_OnMultiple")
	// Output:
	// EventManager_OnMultiple
}

func ExampleEventManager_Reset() {
	core.Println("EventManager_Reset")
	// Output:
	// EventManager_Reset
}

func ExampleEventManager_OnApplicationEvent() {
	core.Println("EventManager_OnApplicationEvent")
	// Output:
	// EventManager_OnApplicationEvent
}

func ExampleEventManager_RegisterApplicationEventHook() {
	core.Println("EventManager_RegisterApplicationEventHook")
	// Output:
	// EventManager_RegisterApplicationEventHook
}
