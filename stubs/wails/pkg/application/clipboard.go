package application

import "sync"

// Clipboard provides direct read/write access to the system clipboard.
//
//	ok := clipboard.SetText("hello")
//	text, ok := clipboard.Text()
type Clipboard struct {
	mu   sync.RWMutex
	text string
}

// SetText writes the given text to the clipboard and returns true on success.
//
//	ok := clipboard.SetText("copied text")
func (c *Clipboard) SetText(text string) bool {
	c.mu.Lock()
	c.text = text
	c.mu.Unlock()
	return true
}

// Text reads the current clipboard text. Returns the text and true on success.
//
//	text, ok := clipboard.Text()
func (c *Clipboard) Text() (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.text, true
}

// ClipboardManager is the application-level clipboard surface.
// Lazily initialises the underlying Clipboard on first use.
//
//	manager.SetText("hello")
//	text, ok := manager.Text()
type ClipboardManager struct {
	mu        sync.Mutex
	clipboard *Clipboard
}

// SetText writes text to the clipboard and returns true on success.
//
//	ok := manager.SetText("hello")
func (cm *ClipboardManager) SetText(text string) bool {
	return cm.instance().SetText(text)
}

// Text reads the current clipboard text. Returns the text and true on success.
//
//	text, ok := manager.Text()
func (cm *ClipboardManager) Text() (string, bool) {
	return cm.instance().Text()
}

// instance returns the Clipboard, creating it if it does not yet exist.
func (cm *ClipboardManager) instance() *Clipboard {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.clipboard == nil {
		cm.clipboard = &Clipboard{}
	}
	return cm.clipboard
}
