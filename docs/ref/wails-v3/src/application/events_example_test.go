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

func ExampleCustomEvent_ToJSON() {
	core.Println("CustomEvent_ToJSON")
	// Output:
	// CustomEvent_ToJSON
}

func ExampleNewWailsEventProcessor() {
	core.Println("NewWailsEventProcessor")
	// Output:
	// NewWailsEventProcessor
}

func ExampleEventProcessor_On() {
	core.Println("EventProcessor_On")
	// Output:
	// EventProcessor_On
}

func ExampleEventProcessor_OnMultiple() {
	core.Println("EventProcessor_OnMultiple")
	// Output:
	// EventProcessor_OnMultiple
}

func ExampleEventProcessor_Once() {
	core.Println("EventProcessor_Once")
	// Output:
	// EventProcessor_Once
}

func ExampleEventProcessor_Emit() {
	core.Println("EventProcessor_Emit")
	// Output:
	// EventProcessor_Emit
}

func ExampleEventProcessor_Off() {
	core.Println("EventProcessor_Off")
	// Output:
	// EventProcessor_Off
}

func ExampleEventProcessor_OffAll() {
	core.Println("EventProcessor_OffAll")
	// Output:
	// EventProcessor_OffAll
}

func ExampleEventProcessor_RegisterHook() {
	core.Println("EventProcessor_RegisterHook")
	// Output:
	// EventProcessor_RegisterHook
}

func ExampleRegisterEvent() {
	core.Println("RegisterEvent")
	// Output:
	// RegisterEvent
}
