// pkg/window/service_config_test.go
package window

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newConfigTestWindowService(t *testing.T) (*Service, *Manager) {
	t.Helper()

	mgr := &Manager{
		platform: newMockPlatform(),
		state:    NewStateManagerWithDir(t.TempDir()),
		layout:   NewLayoutManagerWithDir(t.TempDir()),
		windows:  make(map[string]PlatformWindow),
	}

	return &Service{manager: mgr}, mgr
}

func TestServiceConfig_applyConfig_Good(t *testing.T) {
	svc, mgr := newConfigTestWindowService(t)

	stateFile := filepath.Join(t.TempDir(), "window-state.json")
	svc.applyConfig(map[string]any{
		"default_width":  1440,
		"default_height": 900,
		"state_file":     stateFile,
	})

	require.Equal(t, stateFile, mgr.State().filePath())

	pw, err := mgr.Open(WithName("main"))
	require.NoError(t, err)

	width, height := pw.Size()
	assert.Equal(t, 1440, width)
	assert.Equal(t, 900, height)

	mgr.State().SetState("main", WindowState{Width: width, Height: height})
	mgr.State().ForceSync()

	content, err := os.ReadFile(stateFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"main"`)
}

func TestServiceConfig_applyConfig_Bad(t *testing.T) {
	svc, mgr := newConfigTestWindowService(t)
	mgr.SetDefaultWidth(1111)
	mgr.SetDefaultHeight(2222)
	initialPath := mgr.State().filePath()

	svc.applyConfig(map[string]any{
		"default_width":  "wide",
		"default_height": true,
		"state_file":     123,
	})

	assert.Equal(t, 1111, mgr.defaultWidth)
	assert.Equal(t, 2222, mgr.defaultHeight)
	assert.Equal(t, initialPath, mgr.State().filePath())
}

func TestServiceConfig_applyConfig_Ugly(t *testing.T) {
	svc, mgr := newConfigTestWindowService(t)
	initialPath := mgr.State().filePath()

	require.NotPanics(t, func() {
		svc.applyConfig(nil)
	})

	assert.Equal(t, initialPath, mgr.State().filePath())
}
