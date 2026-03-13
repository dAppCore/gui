// pkg/clipboard/messages.go
package clipboard

// QueryText reads the clipboard. Result: ClipboardContent
type QueryText struct{}

// TaskSetText writes text to the clipboard. Result: bool (success)
type TaskSetText struct{ Text string }

// TaskClear clears the clipboard. Result: bool (success)
type TaskClear struct{}
