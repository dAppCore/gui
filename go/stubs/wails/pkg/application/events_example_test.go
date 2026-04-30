//go:build compliance

package application

import core "dappco.re/go"

func ExampleApplicationEvent_Context() {
	core.Println("ApplicationEvent_Context")
	// Output:
	// ApplicationEvent_Context
}

func ExampleApplicationEvent_Cancel() {
	core.Println("ApplicationEvent_Cancel")
	// Output:
	// ApplicationEvent_Cancel
}

func ExampleApplicationEvent_IsCancelled() {
	core.Println("ApplicationEvent_IsCancelled")
	// Output:
	// ApplicationEvent_IsCancelled
}

func ExampleCustomEvent_Cancel() {
	core.Println("CustomEvent_Cancel")
	// Output:
	// CustomEvent_Cancel
}

func ExampleCustomEvent_IsCancelled() {
	core.Println("CustomEvent_IsCancelled")
	// Output:
	// CustomEvent_IsCancelled
}

func ExampleEventManager_Emit() {
	core.Println("EventManager_Emit")
	// Output:
	// EventManager_Emit
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

func ExampleEventManager_OnApplicationEvent() {
	core.Println("EventManager_OnApplicationEvent")
	// Output:
	// EventManager_OnApplicationEvent
}
