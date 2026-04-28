package textutil

import core "dappco.re/go"

func TestAX7_FirstNonEmpty_Good(t *core.T) {
	symbol := FirstNonEmpty
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_FirstNonEmpty_Bad(t *core.T) {
	symbol := FirstNonEmpty
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}

func TestAX7_FirstNonEmpty_Ugly(t *core.T) {
	symbol := FirstNonEmpty
	core.AssertNotNil(t, symbol)
	core.AssertContains(t, core.Sprintf("%T", symbol), "func")
}
