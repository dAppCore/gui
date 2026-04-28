// pkg/systray/service_config_test.go
package systray

import (
	core "dappco.re/go"
)

func newConfigTestSystrayService(t *core.T) (*Service, *Manager) {
	t.Helper()

	mgr := NewManager(newMockPlatform())
	return &Service{manager: mgr}, mgr
}

func TestServiceConfig_applyConfig_Good(t *core.T) {
	svc, mgr := newConfigTestSystrayService(t)

	svc.applyConfig(map[string]any{
		"tooltip": "Core Ready",
		"icon":    "assets/tray.png",
	})

	info := mgr.GetInfo()
	core.AssertEqual(t, "Core Ready", info["tooltip"])
	core.AssertEqual(t, "Core Ready", info["label"])
	core.AssertEqual(t, "assets/tray.png", svc.iconPath)
	core.AssertTrue(t, mgr.IsActive())
}

func TestServiceConfig_applyConfig_Bad(t *core.T) {
	svc, mgr := newConfigTestSystrayService(t)

	svc.applyConfig(map[string]any{
		"tooltip": 123,
		"icon":    true,
	})

	info := mgr.GetInfo()
	core.AssertEqual(t, "Core", info["tooltip"])
	core.AssertEqual(t, "Core", info["label"])
	core.AssertEmpty(t, svc.iconPath)
}

func TestServiceConfig_applyConfig_Ugly(t *core.T) {
	svc, mgr := newConfigTestSystrayService(t)

	core.AssertNotPanics(t, func() {
		svc.applyConfig(nil)
	})

	info := mgr.GetInfo()
	core.AssertEqual(t, "Core", info["tooltip"])
	core.AssertEqual(t, "Core", info["label"])
	core.AssertEmpty(t, svc.iconPath)
}
