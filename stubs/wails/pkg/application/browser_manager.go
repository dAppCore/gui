package application

import "sync"

// BrowserManager manages browser-related operations in-memory.
//
//	manager := &BrowserManager{}
//	manager.OpenURL("https://example.com")
//	last := manager.LastURL // "https://example.com"
type BrowserManager struct {
	mu       sync.RWMutex
	LastURL  string
	LastFile string
}

// OpenURL stores the URL as the last opened URL.
//
//	manager.OpenURL("https://lthn.io")
//	_ = manager.LastURL // "https://lthn.io"
func (bm *BrowserManager) OpenURL(url string) error {
	if bm == nil {
		return nil
	}
	bm.mu.Lock()
	bm.LastURL = url
	bm.mu.Unlock()
	return nil
}

// OpenFile stores the path as the last opened file.
//
//	manager.OpenFile("/home/user/report.pdf")
//	_ = manager.LastFile // "/home/user/report.pdf"
func (bm *BrowserManager) OpenFile(path string) error {
	if bm == nil {
		return nil
	}
	bm.mu.Lock()
	bm.LastFile = path
	bm.mu.Unlock()
	return nil
}
