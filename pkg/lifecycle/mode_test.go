package lifecycle

import (
	"os"
	"testing"
)

func TestMode_DetectMode_Good(t *testing.T) {
	unsetEnv(t, appModeEnv)
	t.Setenv(ciEnv, "")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode, got %q", mode)
	}
}

func TestMode_DetectMode_Bad(t *testing.T) {
	t.Setenv(appModeEnv, "bogus")
	t.Setenv(ciEnv, "")

	mode := DetectMode()
	if mode != ModeManager {
		t.Fatalf("expected manager mode after invalid env value, got %q", mode)
	}
}

func TestMode_DetectMode_Ugly(t *testing.T) {
	t.Setenv(appModeEnv, "")
	t.Setenv(ciEnv, "true")

	mode := DetectMode()
	if mode != ModeWorker {
		t.Fatalf("expected worker mode in CI headless context, got %q", mode)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	value, ok := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}

	t.Cleanup(func() {
		if ok {
			if err := os.Setenv(key, value); err != nil {
				t.Fatalf("restore %s: %v", key, err)
			}
			return
		}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("restore unset %s: %v", key, err)
		}
	})
}
