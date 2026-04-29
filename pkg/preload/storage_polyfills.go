package preload

import (
	strings "dappco.re/go/gui/compat/strings"

	core "dappco.re/go"
)

func renderStoragePolyfills(pageURL string, canPersist bool) string {
	meta := map[string]any{
		"pageURL":       pageURL,
		"storageOrigin": storageOriginForPageURL(pageURL),
		"storeGroup":    "gui.preload.storage",
		"canPersist":    canPersist,
	}

	return strings.ReplaceAll(
		storagePolyfillsAsset,
		"__CORE_PRELOAD_META__",
		core.JSONMarshalString(meta),
	)
}
