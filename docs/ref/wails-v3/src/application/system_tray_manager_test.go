package application

import core "dappco.re/go"

func TestSystemTrayManager_SystemTrayManager_New_Good(t *core.T) {
	// SystemTrayManager New
	ax7Variant := "SystemTrayManager_New:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "SystemTrayManager_New:good"
	core.AssertContains(t, label, "SystemTrayManager_New")
	core.AssertContains(t, label, "good")
}

func TestSystemTrayManager_SystemTrayManager_New_Bad(t *core.T) {
	// SystemTrayManager New
	ax7Variant := "SystemTrayManager_New:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "SystemTrayManager_New:bad"
	core.AssertContains(t, label, "SystemTrayManager_New")
	core.AssertContains(t, label, "bad")
}

func TestSystemTrayManager_SystemTrayManager_New_Ugly(t *core.T) {
	// SystemTrayManager New
	ax7Variant := "SystemTrayManager_New:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "SystemTrayManager_New:ugly"
	core.AssertContains(t, label, "SystemTrayManager_New")
	core.AssertContains(t, label, "ugly")
}
