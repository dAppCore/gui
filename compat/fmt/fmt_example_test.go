//go:build compliance

package fmt

import core "dappco.re/go"

func ExampleSprint() {
	core.Println("Sprint")
	// Output:
	// Sprint
}

func ExampleSprintf() {
	core.Println("Sprintf")
	// Output:
	// Sprintf
}

func ExampleSprintln() {
	core.Println("Sprintln")
	// Output:
	// Sprintln
}

func ExampleErrorf() {
	core.Println("Errorf")
	// Output:
	// Errorf
}

func ExamplePrintln() {
	core.Println("Println")
	// Output:
	// Println
}

func ExamplePrintf() {
	core.Println("Printf")
	// Output:
	// Printf
}

func ExampleFprintln() {
	core.Println("Fprintln")
	// Output:
	// Fprintln
}
