//go:build compliance

package exec

import core "dappco.re/go"

func ExampleCommand() {
	core.Println("Command")
	// Output:
	// Command
}

func ExampleCommandContext() {
	core.Println("CommandContext")
	// Output:
	// CommandContext
}

func ExampleLookPath() {
	core.Println("LookPath")
	// Output:
	// LookPath
}
