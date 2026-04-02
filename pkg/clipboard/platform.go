// pkg/clipboard/platform.go
package clipboard

import "encoding/base64"

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

// imageReader is an optional clipboard capability for image reads.
type imageReader interface {
	Image() ([]byte, bool)
}

// imageWriter is an optional clipboard capability for image writes.
type imageWriter interface {
	SetImage(data []byte) bool
}

// encodeImageContent converts raw bytes to transport-safe clipboard image content.
func encodeImageContent(data []byte) ClipboardImageContent {
	return ClipboardImageContent{
		Base64:     base64.StdEncoding.EncodeToString(data),
		MimeType:   "image/png",
		HasContent: len(data) > 0,
	}
}
