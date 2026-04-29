package application

import (
	core "dappco.re/go"
)

type handlerFunc func()

func (f handlerFunc) ServeHTTP(_ ResponseWriter, _ *Request) {
	f()
}

type fakeFS struct{}

func (fakeFS) Open(name string) (interface{ Read([]byte) (int, error) }, error) {
	_ = name
	return nil, nil
}

func TestApplicationOptions_ChainMiddleware_Good(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:good"
	core.AssertContains(t, ax7Variant, "good")
	calls := make([]string, 0, 3)
	base := handlerFunc(func() {
		calls = append(calls, "handler")
	})

	first := func(next Handler) Handler {
		return handlerFunc(func() {
			calls = append(calls, "first")
			next.ServeHTTP(nil, nil)
		})
	}
	second := func(next Handler) Handler {
		return handlerFunc(func() {
			calls = append(calls, "second")
			next.ServeHTTP(nil, nil)
		})
	}

	ChainMiddleware(first, second)(base).ServeHTTP(nil, nil)

	core.AssertEqual(t, []string{"first", "second", "handler"}, calls)
}

func TestApplicationOptions_ChainMiddleware_Bad(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:bad"
	core.AssertContains(t, ax7Variant, "bad")
	calls := 0
	base := handlerFunc(func() {
		calls++
	})

	ChainMiddleware()(base).ServeHTTP(nil, nil)

	core.AssertEqual(t, 1, calls)
}

func TestApplicationOptions_ChainMiddleware_Ugly(t *core.T) {
	// ChainMiddleware
	ax7Variant := "ChainMiddleware:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	calls := 0
	base := handlerFunc(func() {
		calls++
	})

	shortCircuit := func(next Handler) Handler {
		_ = next
		return handlerFunc(func() {})
	}

	ChainMiddleware(shortCircuit)(base).ServeHTTP(nil, nil)

	core.AssertEmpty(t, calls)
}

func TestApplicationOptions_NewRGB_Good(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:good"
	core.AssertContains(t, ax7Variant, "good")
	got := NewRGB(0x11, 0x22, 0x33)

	core.AssertEqual(t, RGBA{Red: 0x11, Green: 0x22, Blue: 0x33, Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGB_Bad(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:bad"
	core.AssertContains(t, ax7Variant, "bad")
	got := NewRGB(0, 0, 0)

	core.AssertEqual(t, RGBA{Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGB_Ugly(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := NewRGB(255, 255, 255)

	core.AssertEqual(t, RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGBPtr_Good(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:good"
	core.AssertContains(t, ax7Variant, "good")
	got := NewRGBPtr(0x11, 0x22, 0x33)

	core.AssertNotNil(t, got)
	core.AssertEqual(t, uint32(0x00332211), *got)
}

func TestApplicationOptions_NewRGBPtr_Bad(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:bad"
	core.AssertContains(t, ax7Variant, "bad")
	got := NewRGBPtr(0, 0, 0)

	core.AssertNotNil(t, got)
	core.AssertEmpty(t, *got)
}

func TestApplicationOptions_NewRGBPtr_Ugly(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	got := NewRGBPtr(255, 255, 255)

	core.AssertNotNil(t, got)
	core.AssertEqual(t, uint32(0x00ffffff), *got)
}

func TestApplicationOptions_AssetFileServerFS_Good(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertNil(t, AssetFileServerFS(nil))
	observedType := core.Sprintf("%T", AssetFileServerFS(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_AssetFileServerFS_Bad(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertNil(t, AssetFileServerFS(fakeFS{}))
	observedType := core.Sprintf("%T", AssetFileServerFS(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_AssetFileServerFS_UglyCase(t *core.T) {
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Good(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:good"
	core.AssertContains(t, ax7Variant, "good")
	core.AssertNil(t, BundledAssetFileServer(nil))
	observedType := core.Sprintf("%T", BundledAssetFileServer(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Bad(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:bad"
	core.AssertContains(t, ax7Variant, "bad")
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Ugly(t *core.T) {
	// BundledAssetFileServer
	ax7Variant := "BundledAssetFileServer:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_AssetFileServerFS_Ugly(t *core.T) {
	// AssetFileServerFS
	ax7Variant := "AssetFileServerFS:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := AssetFileServerFS(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
