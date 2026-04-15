package container

import "testing"

func TestDetectModeWithEnvironment(t *testing.T) {
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
}
