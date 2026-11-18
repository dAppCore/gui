//go:build compliance

package application

import core "github.com/dappcore/go"

func ExampleClipboard_SetText() {
	core.Println("Clipboard_SetText")
	// Output:
	// Clipboard_SetText
}

func ExampleClipboard_Text() {
	core.Println("Clipboard_Text")
	// Output:
	// Clipboard_Text
}
