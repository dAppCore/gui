//go:build compliance

package application

import core "dappco.re/go"

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
