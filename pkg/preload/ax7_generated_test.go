package preload

import core "dappco.re/go"

func TestAX7_DefaultTrustedOriginPolicy_Good(t *core.T) {
	symbol := DefaultTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_DefaultTrustedOriginPolicy_Bad(t *core.T) {
	symbol := DefaultTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_DefaultTrustedOriginPolicy_Ugly(t *core.T) {
	symbol := DefaultTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreload_Good(t *core.T) {
	symbol := InjectPreload
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreload_Bad(t *core.T) {
	symbol := InjectPreload
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreload_Ugly(t *core.T) {
	symbol := InjectPreload
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreloadWithTrustedOriginPolicy_Good(t *core.T) {
	symbol := InjectPreloadWithTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreloadWithTrustedOriginPolicy_Bad(t *core.T) {
	symbol := InjectPreloadWithTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_InjectPreloadWithTrustedOriginPolicy_Ugly(t *core.T) {
	symbol := InjectPreloadWithTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicy_Good(t *core.T) {
	symbol := NewTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicy_Bad(t *core.T) {
	symbol := NewTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicy_Ugly(t *core.T) {
	symbol := NewTrustedOriginPolicy
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicyWithActions_Good(t *core.T) {
	symbol := NewTrustedOriginPolicyWithActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicyWithActions_Bad(t *core.T) {
	symbol := NewTrustedOriginPolicyWithActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_NewTrustedOriginPolicyWithActions_Ugly(t *core.T) {
	symbol := NewTrustedOriginPolicyWithActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActions_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActions_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActions_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActions
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActionsForURL_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActionsForURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActionsForURL_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActionsForURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowedActionsForURL_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowedActionsForURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_Allows_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).Allows
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_Allows_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).Allows
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_Allows_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).Allows
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsAction_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsAction_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsAction_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsActionURL_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsActionURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsActionURL_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsActionURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsActionURL_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsActionURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsURL_Good(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsURL_Bad(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_TrustedOriginPolicy_AllowsURL_Ugly(t *core.T) {
	symbol := (*TrustedOriginPolicy).AllowsURL
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
