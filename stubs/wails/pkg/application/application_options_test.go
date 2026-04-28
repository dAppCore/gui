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
	calls := 0
	base := handlerFunc(func() {
		calls++
	})

	ChainMiddleware()(base).ServeHTTP(nil, nil)

	core.AssertEqual(t, 1, calls)
}

func TestApplicationOptions_ChainMiddleware_Ugly(t *core.T) {
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
	got := NewRGB(0x11, 0x22, 0x33)

	core.AssertEqual(t, RGBA{Red: 0x11, Green: 0x22, Blue: 0x33, Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGB_Bad(t *core.T) {
	got := NewRGB(0, 0, 0)

	core.AssertEqual(t, RGBA{Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGB_Ugly(t *core.T) {
	got := NewRGB(255, 255, 255)

	core.AssertEqual(t, RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestApplicationOptions_NewRGBPtr_Good(t *core.T) {
	got := NewRGBPtr(0x11, 0x22, 0x33)

	core.AssertNotNil(t, got)
	core.AssertEqual(t, uint32(0x00332211), *got)
}

func TestApplicationOptions_NewRGBPtr_Bad(t *core.T) {
	got := NewRGBPtr(0, 0, 0)

	core.AssertNotNil(t, got)
	core.AssertEmpty(t, *got)
}

func TestApplicationOptions_NewRGBPtr_Ugly(t *core.T) {
	got := NewRGBPtr(255, 255, 255)

	core.AssertNotNil(t, got)
	core.AssertEqual(t, uint32(0x00ffffff), *got)
}

func TestApplicationOptions_AssetFileServerFS_Good(t *core.T) {
	core.AssertNil(t, AssetFileServerFS(nil))
	observedType := core.Sprintf("%T", AssetFileServerFS(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_AssetFileServerFS_Bad(t *core.T) {
	core.AssertNil(t, AssetFileServerFS(fakeFS{}))
	observedType := core.Sprintf("%T", AssetFileServerFS(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_AssetFileServerFS_Ugly(t *core.T) {
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Good(t *core.T) {
	core.AssertNil(t, BundledAssetFileServer(nil))
	observedType := core.Sprintf("%T", BundledAssetFileServer(nil))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Bad(t *core.T) {
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}

func TestApplicationOptions_BundledAssetFileServer_Ugly(t *core.T) {
	core.AssertNil(t, BundledAssetFileServer(fakeFS{}))
	observedType := core.Sprintf("%T", BundledAssetFileServer(fakeFS{}))
	core.AssertNotEmpty(t, observedType)
}
