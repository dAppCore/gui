package display

import "strings"

type AppMode string

const (
	ModeManager AppMode = "manager"
	ModeWorker  AppMode = "worker"
)

func DetectMode(args []string, getenv func(string) string) AppMode {
	if getenv != nil {
		if mode, ok := parseMode(getenv("CORE_GUI_MODE")); ok {
			return mode
		}
	}

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--mode=") {
			continue
		}
		if mode, ok := parseMode(strings.TrimPrefix(arg, "--mode=")); ok {
			return mode
		}
	}

	return ModeManager
}

func parseMode(raw string) (AppMode, bool) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case string(ModeManager):
		return ModeManager, true
	case string(ModeWorker):
		return ModeWorker, true
	default:
		return "", false
	}
}
