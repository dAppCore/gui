package container

import (
	"os"
	"strings"

	"dappco.re/go/config"
)

type AppMode string

const (
	ModeManager AppMode = "manager"
	ModeWorker  AppMode = "worker"
)

type ModeEnvironment struct {
	Args        []string
	LookupEnv   func(string) (string, bool)
	ConfigValue func(string) string
}

// DetectMode resolves the RFC startup mode from CLI flags first, then config/env.
func DetectMode() AppMode {
	cfg, _ := config.New()
	return DetectModeWithEnvironment(ModeEnvironment{
		Args:      os.Args[1:],
		LookupEnv: os.LookupEnv,
		ConfigValue: func(key string) string {
			if cfg == nil {
				return ""
			}
			var value string
			if err := cfg.Get(key, &value); err != nil {
				return ""
			}
			return value
		},
	})
}

func DetectModeWithEnvironment(environment ModeEnvironment) AppMode {
	args := environment.Args
	if args == nil {
		args = os.Args[1:]
	}

	if value, found := modeArgValue(args); found {
		if mode, ok := parseMode(value); ok {
			return mode
		}
	}

	if environment.LookupEnv != nil {
		if value, ok := environment.LookupEnv("CORE_GUI_MODE"); ok {
			if mode, ok := parseMode(value); ok {
				return mode
			}
		}
	}

	if environment.ConfigValue != nil {
		for _, key := range []string{"gui.mode", "display.mode", "mode"} {
			if mode, ok := parseMode(environment.ConfigValue(key)); ok {
				return mode
			}
		}
	}

	return ModeManager
}

func modeArgValue(args []string) (string, bool) {
	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		if arg == "--mode" || arg == "-mode" {
			if i+1 < len(args) {
				return args[i+1], true
			}
			return "", true
		}
		if value, ok := strings.CutPrefix(arg, "--mode="); ok {
			return value, true
		}
		if value, ok := strings.CutPrefix(arg, "-mode="); ok {
			return value, true
		}
	}
	return "", false
}

func parseMode(value string) (AppMode, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ModeWorker):
		return ModeWorker, true
	case string(ModeManager):
		return ModeManager, true
	default:
		return "", false
	}
}
