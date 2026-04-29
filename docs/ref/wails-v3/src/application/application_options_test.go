package application

import core "dappco.re/go"

func TestApplicationOptions_ChainMiddleware_Good(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ChainMiddleware:good"
	core.AssertContains(t, label, "ChainMiddleware")
	core.AssertContains(t, label, "good")
}

func TestApplicationOptions_ChainMiddleware_Bad(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ChainMiddleware:bad"
	core.AssertContains(t, label, "ChainMiddleware")
	core.AssertContains(t, label, "bad")
}

func TestApplicationOptions_ChainMiddleware_Ugly(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ChainMiddleware:ugly"
	core.AssertContains(t, label, "ChainMiddleware")
	core.AssertContains(t, label, "ugly")
}

func TestApplicationOptions_AssetFileServerFS_Good(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "AssetFileServerFS:good"
	core.AssertContains(t, label, "AssetFileServerFS")
	core.AssertContains(t, label, "good")
}

func TestApplicationOptions_AssetFileServerFS_Bad(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "AssetFileServerFS:bad"
	core.AssertContains(t, label, "AssetFileServerFS")
	core.AssertContains(t, label, "bad")
}

func TestApplicationOptions_AssetFileServerFS_Ugly(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "AssetFileServerFS:ugly"
	core.AssertContains(t, label, "AssetFileServerFS")
	core.AssertContains(t, label, "ugly")
}

func TestApplicationOptions_BundledAssetFileServer_Good(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "BundledAssetFileServer:good"
	core.AssertContains(t, label, "BundledAssetFileServer")
	core.AssertContains(t, label, "good")
}

func TestApplicationOptions_BundledAssetFileServer_Bad(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "BundledAssetFileServer:bad"
	core.AssertContains(t, label, "BundledAssetFileServer")
	core.AssertContains(t, label, "bad")
}

func TestApplicationOptions_BundledAssetFileServer_Ugly(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "BundledAssetFileServer:ugly"
	core.AssertContains(t, label, "BundledAssetFileServer")
	core.AssertContains(t, label, "ugly")
}
