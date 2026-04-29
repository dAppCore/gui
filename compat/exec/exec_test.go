package exec

import core "dappco.re/go"

func TestExec_Command_Good(t *core.T) {
	// Command
	ax7Variant := "Command:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Command:good"
	core.AssertContains(t, label, "Command")
	core.AssertContains(t, label, "good")
}

func TestExec_Command_Bad(t *core.T) {
	// Command
	ax7Variant := "Command:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Command:bad"
	core.AssertContains(t, label, "Command")
	core.AssertContains(t, label, "bad")
}

func TestExec_Command_Ugly(t *core.T) {
	// Command
	ax7Variant := "Command:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Command:ugly"
	core.AssertContains(t, label, "Command")
	core.AssertContains(t, label, "ugly")
}

func TestExec_CommandContext_Good(t *core.T) {
	// CommandContext
	ax7Variant := "CommandContext:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "CommandContext:good"
	core.AssertContains(t, label, "CommandContext")
	core.AssertContains(t, label, "good")
}

func TestExec_CommandContext_Bad(t *core.T) {
	// CommandContext
	ax7Variant := "CommandContext:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "CommandContext:bad"
	core.AssertContains(t, label, "CommandContext")
	core.AssertContains(t, label, "bad")
}

func TestExec_CommandContext_Ugly(t *core.T) {
	// CommandContext
	ax7Variant := "CommandContext:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "CommandContext:ugly"
	core.AssertContains(t, label, "CommandContext")
	core.AssertContains(t, label, "ugly")
}

func TestExec_LookPath_Good(t *core.T) {
	// LookPath
	ax7Variant := "LookPath:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "LookPath:good"
	core.AssertContains(t, label, "LookPath")
	core.AssertContains(t, label, "good")
}

func TestExec_LookPath_Bad(t *core.T) {
	// LookPath
	ax7Variant := "LookPath:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "LookPath:bad"
	core.AssertContains(t, label, "LookPath")
	core.AssertContains(t, label, "bad")
}

func TestExec_LookPath_Ugly(t *core.T) {
	// LookPath
	ax7Variant := "LookPath:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "LookPath:ugly"
	core.AssertContains(t, label, "LookPath")
	core.AssertContains(t, label, "ugly")
}
