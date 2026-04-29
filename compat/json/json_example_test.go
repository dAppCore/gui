//go:build compliance

package json

import core "dappco.re/go"

func ExampleMarshal() {
	core.Println("Marshal")
	// Output:
	// Marshal
}

func ExampleMarshalIndent() {
	core.Println("MarshalIndent")
	// Output:
	// MarshalIndent
}

func ExampleUnmarshal() {
	core.Println("Unmarshal")
	// Output:
	// Unmarshal
}
