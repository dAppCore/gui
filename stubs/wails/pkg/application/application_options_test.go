package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestApplicationOptions_ChainMiddleware_Good(t *testing.T) {
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

	assert.Equal(t, []string{"first", "second", "handler"}, calls)
}

func TestApplicationOptions_ChainMiddleware_Bad(t *testing.T) {
	calls := 0
	base := handlerFunc(func() {
		calls++
	})

	ChainMiddleware()(base).ServeHTTP(nil, nil)

	assert.Equal(t, 1, calls)
}

func TestApplicationOptions_ChainMiddleware_Ugly(t *testing.T) {
	calls := 0
	base := handlerFunc(func() {
		calls++
	})

	shortCircuit := func(next Handler) Handler {
		_ = next
		return handlerFunc(func() {})
	}

	ChainMiddleware(shortCircuit)(base).ServeHTTP(nil, nil)

	assert.Zero(t, calls)
}

func TestApplicationOptions_NewRGB_Good(t *testing.T) {
	got := NewRGB(0x11, 0x22, 0x33)

	assert.Equal(t, RGBA{Red: 0x11, Green: 0x22, Blue: 0x33, Alpha: 255}, got)
}

func TestApplicationOptions_NewRGB_Bad(t *testing.T) {
	got := NewRGB(0, 0, 0)

	assert.Equal(t, RGBA{Alpha: 255}, got)
}

func TestApplicationOptions_NewRGB_Ugly(t *testing.T) {
	got := NewRGB(255, 255, 255)

	assert.Equal(t, RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255}, got)
}

func TestApplicationOptions_NewRGBPtr_Good(t *testing.T) {
	got := NewRGBPtr(0x11, 0x22, 0x33)

	require.NotNil(t, got)
	assert.Equal(t, uint32(0x00332211), *got)
}

func TestApplicationOptions_NewRGBPtr_Bad(t *testing.T) {
	got := NewRGBPtr(0, 0, 0)

	require.NotNil(t, got)
	assert.Zero(t, *got)
}

func TestApplicationOptions_NewRGBPtr_Ugly(t *testing.T) {
	got := NewRGBPtr(255, 255, 255)

	require.NotNil(t, got)
	assert.Equal(t, uint32(0x00ffffff), *got)
}

func TestApplicationOptions_AssetFileServerFS_Good(t *testing.T) {
	assert.Nil(t, AssetFileServerFS(nil))
}

func TestApplicationOptions_AssetFileServerFS_Bad(t *testing.T) {
	assert.Nil(t, AssetFileServerFS(fakeFS{}))
}

func TestApplicationOptions_AssetFileServerFS_Ugly(t *testing.T) {
	assert.Nil(t, BundledAssetFileServer(fakeFS{}))
}

func TestApplicationOptions_BundledAssetFileServer_Good(t *testing.T) {
	assert.Nil(t, BundledAssetFileServer(nil))
}

func TestApplicationOptions_BundledAssetFileServer_Bad(t *testing.T) {
	assert.Nil(t, BundledAssetFileServer(fakeFS{}))
}

func TestApplicationOptions_BundledAssetFileServer_Ugly(t *testing.T) {
	assert.Nil(t, BundledAssetFileServer(fakeFS{}))
}
