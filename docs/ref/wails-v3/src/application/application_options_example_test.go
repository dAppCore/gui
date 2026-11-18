//go:build compliance

package application

import core "github.com/dappcore/go"

func ExampleChainMiddleware() {
	core.Println("ChainMiddleware")
	// Output:
	// ChainMiddleware
}

func ExampleAssetFileServerFS() {
	core.Println("AssetFileServerFS")
	// Output:
	// AssetFileServerFS
}

func ExampleBundledAssetFileServer() {
	core.Println("BundledAssetFileServer")
	// Output:
	// BundledAssetFileServer
}
