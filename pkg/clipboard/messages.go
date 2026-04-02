// pkg/clipboard/messages.go
package clipboard

// QueryText reads the clipboard. Result: ClipboardContent
type QueryText struct{}

// TaskSetText writes text to the clipboard. Result: bool (success)
type TaskSetText struct{ Text string }

// TaskClear clears the clipboard. Result: bool (success)
type TaskClear struct{}

// QueryImage reads an image from the clipboard. Result: ClipboardImageContent
type QueryImage struct{}

// TaskSetImage writes image bytes to the clipboard. Result: bool (success)
type TaskSetImage struct{ Data []byte }

// ClipboardImageContent contains clipboard image data encoded for transport.
type ClipboardImageContent struct {
	Base64     string `json:"base64"`
	MimeType   string `json:"mimeType"`
	HasContent bool   `json:"hasContent"`
}
