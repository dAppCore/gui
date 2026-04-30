package application

import core "dappco.re/go"

func newBrowserManager() *BrowserManager {
	return &BrowserManager{}
}

func (bm *BrowserManager) Open(target string) error {
	if bm == nil {
		return nil
	}

	if core.Contains(target, "://") || core.HasPrefix(target, "mailto:") {
		return bm.OpenURL(target)
	}

	return bm.OpenFile(target)
}
