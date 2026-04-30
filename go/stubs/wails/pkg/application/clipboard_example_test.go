//go:build compliance

package application

import core "dappco.re/go"

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

func ExampleClipboardManager_SetText() {
	core.Println("ClipboardManager_SetText")
	// Output:
	// ClipboardManager_SetText
}

func ExampleClipboardManager_Text() {
	core.Println("ClipboardManager_Text")
	// Output:
	// ClipboardManager_Text
}
