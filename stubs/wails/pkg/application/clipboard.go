package application

import "sync"

// Clipboard stores and retrieves text in-memory.
//
//	cb := &Clipboard{}
//	cb.SetText("hello")
//	text, ok := cb.Text() // "hello", true
type Clipboard struct {
	mu   sync.RWMutex
	text string
	set  bool
}

// SetText stores the given text in the in-memory clipboard.
//
//	cb.SetText("copied content")
func (c *Clipboard) SetText(text string) bool {
	if c == nil {
		return false
	}
	c.mu.Lock()
	c.text = text
	c.set = true
	c.mu.Unlock()
	return true
}

// Text returns the stored clipboard text and whether any text has been set.
//
//	text, ok := cb.Text()
//	if !ok { text = "" }
func (c *Clipboard) Text() (string, bool) {
	if c == nil {
		return "", false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.text, c.set
}

// ClipboardManager manages clipboard operations via a lazily created Clipboard.
//
//	manager := &ClipboardManager{}
//	manager.SetText("hello")
//	text, ok := manager.Text() // "hello", true
type ClipboardManager struct {
	mu        sync.Mutex
	clipboard *Clipboard
}

// SetText sets text in the clipboard.
//
//	manager.SetText("some text")
func (cm *ClipboardManager) SetText(text string) bool {
	if cm == nil {
		return false
	}
	return cm.getClipboard().SetText(text)
}

// Text gets text from the clipboard.
//
//	text, ok := manager.Text()
func (cm *ClipboardManager) Text() (string, bool) {
	if cm == nil {
		return "", false
	}
	return cm.getClipboard().Text()
}

// getClipboard returns the clipboard instance, creating it if needed.
func (cm *ClipboardManager) getClipboard() *Clipboard {
	if cm == nil {
		return &Clipboard{}
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.clipboard == nil {
		cm.clipboard = &Clipboard{}
	}
	return cm.clipboard
}
