package container

import (
	"flag"
	"os"
	"strings"

	"forge.lthn.ai/core/config"
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
			_ = cfg.Get(key, &value)
			return value
		},
	})
}

func DetectModeWithEnvironment(environment ModeEnvironment) AppMode {
	args := environment.Args
	if len(args) == 0 {
		args = os.Args[1:]
	}

	fs := flag.NewFlagSet("core-gui", flag.ContinueOnError)
	fs.SetOutput(nil)
	modeFlag := fs.String("mode", "", "")
	_ = fs.Parse(args)
	if mode := parseMode(*modeFlag); mode != "" {
		return mode
	}

	if environment.LookupEnv != nil {
		if value, ok := environment.LookupEnv("CORE_GUI_MODE"); ok {
			if mode := parseMode(value); mode != "" {
				return mode
			}
		}
	}

	if environment.ConfigValue != nil {
		for _, key := range []string{"gui.mode", "display.mode", "mode"} {
			if mode := parseMode(environment.ConfigValue(key)); mode != "" {
				return mode
			}
		}
	}

	return ModeManager
}

func parseMode(value string) AppMode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ModeWorker):
		return ModeWorker
	case "", string(ModeManager):
		return ModeManager
	default:
		return ModeManager
	}
}
