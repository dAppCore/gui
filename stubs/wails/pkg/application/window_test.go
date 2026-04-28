package application

import (
	core "dappco.re/go"
)

func TestWindow_Interface_Good(t *core.T) {
	var _ Window = (*WebviewWindow)(nil)
	var _ Window = (*BrowserWindow)(nil)
	core.AssertNotEmpty(t, core.Sprintf("%T", t))
}

func TestWindow_Interface_Bad(t *core.T) {
	var w Window

	core.AssertNil(t, w)
	core.AssertNotEmpty(t, core.Sprintf("%T", w))
}

func TestWindow_Interface_Ugly(t *core.T) {
	var w Window = (*WebviewWindow)(nil)

	core.AssertFalse(t, w == nil)
	core.AssertNotEmpty(t, core.Sprintf("%T", w == nil))
}
