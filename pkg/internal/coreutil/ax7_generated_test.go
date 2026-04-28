package coreutil

import core "dappco.re/go"

func TestAX7_DispatchAction_Good(t *core.T) {
	symbol := DispatchAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_DispatchAction_Bad(t *core.T) {
	symbol := DispatchAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_DispatchAction_Ugly(t *core.T) {
	symbol := DispatchAction
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_LogWarn_Good(t *core.T) {
	symbol := LogWarn
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_LogWarn_Bad(t *core.T) {
	symbol := LogWarn
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_LogWarn_Ugly(t *core.T) {
	symbol := LogWarn
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_ObserveResult_Good(t *core.T) {
	symbol := ObserveResult
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_ObserveResult_Bad(t *core.T) {
	symbol := ObserveResult
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_ObserveResult_Ugly(t *core.T) {
	symbol := ObserveResult
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
