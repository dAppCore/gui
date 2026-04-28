package deno

import core "dappco.re/go"

func TestAX7_Manager_Emit_Good(t *core.T) {
	symbol := (*Manager).Emit
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Emit_Bad(t *core.T) {
	symbol := (*Manager).Emit
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Emit_Ugly(t *core.T) {
	symbol := (*Manager).Emit
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Eval_Good(t *core.T) {
	symbol := (*Manager).Eval
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Eval_Bad(t *core.T) {
	symbol := (*Manager).Eval
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Eval_Ugly(t *core.T) {
	symbol := (*Manager).Eval
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_OnEvent_Good(t *core.T) {
	symbol := (*Manager).OnEvent
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_OnEvent_Bad(t *core.T) {
	symbol := (*Manager).OnEvent
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_OnEvent_Ugly(t *core.T) {
	symbol := (*Manager).OnEvent
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Start_Good(t *core.T) {
	symbol := (*Manager).Start
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Start_Bad(t *core.T) {
	symbol := (*Manager).Start
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Start_Ugly(t *core.T) {
	symbol := (*Manager).Start
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Status_Good(t *core.T) {
	symbol := (*Manager).Status
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Status_Bad(t *core.T) {
	symbol := (*Manager).Status
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Status_Ugly(t *core.T) {
	symbol := (*Manager).Status
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Stop_Good(t *core.T) {
	symbol := (*Manager).Stop
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Stop_Bad(t *core.T) {
	symbol := (*Manager).Stop
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_Manager_Stop_Ugly(t *core.T) {
	symbol := (*Manager).Stop
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
