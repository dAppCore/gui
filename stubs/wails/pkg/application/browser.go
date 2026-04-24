package application

import "strings"

func newBrowserManager() *BrowserManager {
	return &BrowserManager{}
}

func (bm *BrowserManager) Open(target string) error {
	if bm == nil {
		return nil
	}

	if strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
		return bm.OpenURL(target)
	}

	return bm.OpenFile(target)
}
