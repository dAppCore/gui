package preload

import (
	"strings"

	core "dappco.re/go"
)

func renderElectronShim(pageURL string) string {
	meta := map[string]any{
		"allow":   true,
		"pageURL": pageURL,
	}

	return strings.ReplaceAll(
		electronShimAsset,
		"__CORE_PRELOAD_META__",
		core.JSONMarshalString(meta),
	)
}
