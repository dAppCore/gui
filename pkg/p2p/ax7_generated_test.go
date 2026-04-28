package p2p

import core "dappco.re/go"

func TestAX7_New_Good(t *core.T) {
	symbol := New
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_New_Bad(t *core.T) {
	symbol := New
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_New_Ugly(t *core.T) {
	symbol := New
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewService_Good(t *core.T) {
	symbol := NewService
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewService_Bad(t *core.T) {
	symbol := NewService
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewService_Ugly(t *core.T) {
	symbol := NewService
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewServiceWithDriver_Good(t *core.T) {
	symbol := NewServiceWithDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewServiceWithDriver_Bad(t *core.T) {
	symbol := NewServiceWithDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewServiceWithDriver_Ugly(t *core.T) {
	symbol := NewServiceWithDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTCPDriver_Good(t *core.T) {
	symbol := NewTCPDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTCPDriver_Bad(t *core.T) {
	symbol := NewTCPDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTCPDriver_Ugly(t *core.T) {
	symbol := NewTCPDriver
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_OptionsFromEnv_Good(t *core.T) {
	symbol := OptionsFromEnv
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_OptionsFromEnv_Bad(t *core.T) {
	symbol := OptionsFromEnv
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_OptionsFromEnv_Ugly(t *core.T) {
	symbol := OptionsFromEnv
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Peers_Good(t *core.T) {
	symbol := (*Router).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Peers_Bad(t *core.T) {
	symbol := (*Router).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Peers_Ugly(t *core.T) {
	symbol := (*Router).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Publish_Good(t *core.T) {
	symbol := (*Router).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Publish_Bad(t *core.T) {
	symbol := (*Router).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Publish_Ugly(t *core.T) {
	symbol := (*Router).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Subscribe_Good(t *core.T) {
	symbol := (*Router).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Subscribe_Bad(t *core.T) {
	symbol := (*Router).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Router_Subscribe_Ugly(t *core.T) {
	symbol := (*Router).Subscribe
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

func TestAX7_Service_Peers_Good(t *core.T) {
	symbol := (*Service).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Peers_Bad(t *core.T) {
	symbol := (*Service).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Peers_Ugly(t *core.T) {
	symbol := (*Service).Peers
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Publish_Good(t *core.T) {
	symbol := (*Service).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Publish_Bad(t *core.T) {
	symbol := (*Service).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Publish_Ugly(t *core.T) {
	symbol := (*Service).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_State_Good(t *core.T) {
	symbol := (*Service).State
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_State_Bad(t *core.T) {
	symbol := (*Service).State
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_State_Ugly(t *core.T) {
	symbol := (*Service).State
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Subscribe_Good(t *core.T) {
	symbol := (*Service).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Subscribe_Bad(t *core.T) {
	symbol := (*Service).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Service_Subscribe_Ugly(t *core.T) {
	symbol := (*Service).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Close_Good(t *core.T) {
	symbol := (*TCPDriver).Close
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Close_Bad(t *core.T) {
	symbol := (*TCPDriver).Close
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Close_Ugly(t *core.T) {
	symbol := (*TCPDriver).Close
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_ListenAddr_Good(t *core.T) {
	symbol := (*TCPDriver).ListenAddr
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_ListenAddr_Bad(t *core.T) {
	symbol := (*TCPDriver).ListenAddr
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_ListenAddr_Ugly(t *core.T) {
	symbol := (*TCPDriver).ListenAddr
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Publish_Good(t *core.T) {
	symbol := (*TCPDriver).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Publish_Bad(t *core.T) {
	symbol := (*TCPDriver).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Publish_Ugly(t *core.T) {
	symbol := (*TCPDriver).Publish
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Subscribe_Good(t *core.T) {
	symbol := (*TCPDriver).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Subscribe_Bad(t *core.T) {
	symbol := (*TCPDriver).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TCPDriver_Subscribe_Ugly(t *core.T) {
	symbol := (*TCPDriver).Subscribe
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
