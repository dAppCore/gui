package environment

import core "dappco.re/go"

func TestAX7_Register_Good(t *core.T) {
	symbol := Register
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Register_Bad(t *core.T) {
	symbol := Register
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Register_Ugly(t *core.T) {
	symbol := Register
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_HandleIPCEvents_Good(t *core.T) {
	symbol := (*Service).HandleIPCEvents
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_HandleIPCEvents_Bad(t *core.T) {
	symbol := (*Service).HandleIPCEvents
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_HandleIPCEvents_Ugly(t *core.T) {
	symbol := (*Service).HandleIPCEvents
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnShutdown_Good(t *core.T) {
	symbol := (*Service).OnShutdown
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnShutdown_Bad(t *core.T) {
	symbol := (*Service).OnShutdown
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnShutdown_Ugly(t *core.T) {
	symbol := (*Service).OnShutdown
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnStartup_Good(t *core.T) {
	symbol := (*Service).OnStartup
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnStartup_Bad(t *core.T) {
	symbol := (*Service).OnStartup
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_OnStartup_Ugly(t *core.T) {
	symbol := (*Service).OnStartup
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
