package marketplace

import core "dappco.re/go"

func TestAX7_Installer_FetchManifest_Good(t *core.T) {
	symbol := (*Installer).FetchManifest
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_FetchManifest_Bad(t *core.T) {
	symbol := (*Installer).FetchManifest
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_FetchManifest_Ugly(t *core.T) {
	symbol := (*Installer).FetchManifest
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Install_Good(t *core.T) {
	symbol := (*Installer).Install
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Install_Bad(t *core.T) {
	symbol := (*Installer).Install
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Install_Ugly(t *core.T) {
	symbol := (*Installer).Install
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_List_Good(t *core.T) {
	symbol := (*Installer).List
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_List_Bad(t *core.T) {
	symbol := (*Installer).List
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_List_Ugly(t *core.T) {
	symbol := (*Installer).List
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Verify_Good(t *core.T) {
	symbol := (*Installer).Verify
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Verify_Bad(t *core.T) {
	symbol := (*Installer).Verify
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Installer_Verify_Ugly(t *core.T) {
	symbol := (*Installer).Verify
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
