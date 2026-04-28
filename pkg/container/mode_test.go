package container

import (
	core "dappco.re/go"
	"os"
)

func TestDetectModeWithEnvironment(t *core.T) {
	mode := DetectModeWithEnvironment(ModeEnvironment{
		Args: []string{"--mode=worker"},
	})
	if mode != ModeWorker {
		t.Fatalf("expected worker mode, got %q", mode)
	}

	mode = DetectModeWithEnvironment(ModeEnvironment{
		LookupEnv: func(key string) (string, bool) {
			if key == "CORE_GUI_MODE" {
				return "manager", true
			}
			return "", false
		},
	})
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}

	mode = DetectModeWithEnvironment(ModeEnvironment{
		LookupEnv: func(key string) (string, bool) {
			if key == "CORE_GUI_MODE" {
				return "worker", true
			}
			return "", false
		},
	})
	if mode != ModeWorker {
		t.Fatalf("expected worker mode from environment, got %q", mode)
	}
}

func TestMode_DetectMode_Good(t *core.T) {
	oldArgs := os.Args
	os.Args = []string{"core-gui", "--mode=worker"}
	t.Cleanup(func() { os.Args = oldArgs })
	t.Setenv("CORE_GUI_MODE", "manager")

	mode := DetectMode()
	if mode != ModeWorker {
		t.Fatalf("expected worker mode, got %q", mode)
	}
}

func TestMode_DetectMode_Bad(t *core.T) {
	oldArgs := os.Args
	os.Args = []string{"core-gui"}
	t.Cleanup(func() { os.Args = oldArgs })
	t.Setenv("CORE_GUI_MODE", "bogus")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}
}

func TestMode_DetectMode_Ugly(t *core.T) {
	oldArgs := os.Args
	os.Args = []string{"core-gui", "--unexpected=flag"}
	t.Cleanup(func() { os.Args = oldArgs })
	t.Setenv("CORE_GUI_MODE", " \tbogus\n")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode after malformed input, got %q", mode)
	}
}
