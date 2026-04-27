// pkg/systray/service_config_test.go
package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newConfigTestSystrayService(t *testing.T) (*Service, *Manager) {
	t.Helper()

	mgr := NewManager(newMockPlatform())
	return &Service{manager: mgr}, mgr
}

func TestServiceConfig_applyConfig_Good(t *testing.T) {
	svc, mgr := newConfigTestSystrayService(t)

	svc.applyConfig(map[string]any{
		"tooltip": "Core Ready",
		"icon":    "assets/tray.png",
	})

	info := mgr.GetInfo()
	assert.Equal(t, "Core Ready", info["tooltip"])
	assert.Equal(t, "Core Ready", info["label"])
	assert.Equal(t, "assets/tray.png", svc.iconPath)
	assert.True(t, mgr.IsActive())
}

func TestServiceConfig_applyConfig_Bad(t *testing.T) {
	svc, mgr := newConfigTestSystrayService(t)

	svc.applyConfig(map[string]any{
		"tooltip": 123,
		"icon":    true,
	})

	info := mgr.GetInfo()
	assert.Equal(t, "Core", info["tooltip"])
	assert.Equal(t, "Core", info["label"])
	assert.Empty(t, svc.iconPath)
}

func TestServiceConfig_applyConfig_Ugly(t *testing.T) {
	svc, mgr := newConfigTestSystrayService(t)

	require.NotPanics(t, func() {
		svc.applyConfig(nil)
	})

	info := mgr.GetInfo()
	assert.Equal(t, "Core", info["tooltip"])
	assert.Equal(t, "Core", info["label"])
	assert.Empty(t, svc.iconPath)
}
