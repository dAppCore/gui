// pkg/clipboard/messages.go
package clipboard

// QueryText reads the clipboard. Result: ClipboardContent
type QueryText struct{}

// QueryImage reads image data from the clipboard. Result: ImageContent
type QueryImage struct{}

// TaskSetText writes text to the clipboard. Result: bool (success)
type TaskSetText struct{ Text string }

// TaskSetImage writes image data to the clipboard. Result: bool (success)
type TaskSetImage struct{ Data []byte }

// TaskClear clears the clipboard. Result: bool (success)
type TaskClear struct{}
