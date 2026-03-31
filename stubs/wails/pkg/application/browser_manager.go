package application

// BrowserManager handles opening URLs and files in the system browser.
//
//	manager.OpenURL("https://lthn.io")
//	manager.OpenFile("/home/user/document.pdf")
type BrowserManager struct{}

// OpenURL opens the given URL in the default browser.
//
//	err := manager.OpenURL("https://lthn.io")
func (bm *BrowserManager) OpenURL(url string) error {
	return nil
}

// OpenFile opens the given file path in the default browser or file handler.
//
//	err := manager.OpenFile("/home/user/report.html")
func (bm *BrowserManager) OpenFile(path string) error {
	return nil
}
