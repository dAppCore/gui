// pkg/clipboard/platform.go
package clipboard

// Platform abstracts the system clipboard backend.
type Platform interface {
	Text() (string, bool)
	SetText(text string) bool
}

// ClipboardContent is the result of QueryText.
type ClipboardContent struct {
	Text       string `json:"text"`
	HasContent bool   `json:"hasContent"`
}
