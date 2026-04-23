package lifecycle

import (
	"os"
	"strings"
)

type AppMode string

const (
	ModeManager AppMode = "manager"
	ModeWorker  AppMode = "worker"
)

const (
	appModeEnv = "CORE_APP_MODE"
	ciEnv      = "CI"
)

// DetectMode returns ModeWorker when CORE_APP_MODE=worker.
//
//	mode := DetectMode()
func DetectMode() AppMode {
	if value, ok := os.LookupEnv(appModeEnv); ok {
		if mode, valid := parseAppMode(value); valid {
			return mode
		}
	}

	if value, ok := os.LookupEnv(ciEnv); ok && isTrue(value) {
		return ModeWorker
	}

	return ModeManager
}

func parseAppMode(value string) (AppMode, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ModeManager):
		return ModeManager, true
	case string(ModeWorker):
		return ModeWorker, true
	case "":
		return "", false
	default:
		return ModeManager, true
	}
}

func isTrue(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "true")
}
