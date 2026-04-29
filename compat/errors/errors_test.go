package errors

import core "dappco.re/go"

func TestErrors_New_Good(t *core.T) {
	// New
	ax7Variant := "New:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "New:good"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "good")
}

func TestErrors_New_Bad(t *core.T) {
	// New
	ax7Variant := "New:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "New:bad"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "bad")
}

func TestErrors_New_Ugly(t *core.T) {
	// New
	ax7Variant := "New:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "New:ugly"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "ugly")
}

func TestErrors_Is_Good(t *core.T) {
	// Is
	ax7Variant := "Is:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Is:good"
	core.AssertContains(t, label, "Is")
	core.AssertContains(t, label, "good")
}

func TestErrors_Is_Bad(t *core.T) {
	// Is
	ax7Variant := "Is:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Is:bad"
	core.AssertContains(t, label, "Is")
	core.AssertContains(t, label, "bad")
}

func TestErrors_Is_Ugly(t *core.T) {
	// Is
	ax7Variant := "Is:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Is:ugly"
	core.AssertContains(t, label, "Is")
	core.AssertContains(t, label, "ugly")
}

func TestErrors_As_Good(t *core.T) {
	// As
	ax7Variant := "As:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "As:good"
	core.AssertContains(t, label, "As")
	core.AssertContains(t, label, "good")
}

func TestErrors_As_Bad(t *core.T) {
	// As
	ax7Variant := "As:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "As:bad"
	core.AssertContains(t, label, "As")
	core.AssertContains(t, label, "bad")
}

func TestErrors_As_Ugly(t *core.T) {
	// As
	ax7Variant := "As:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "As:ugly"
	core.AssertContains(t, label, "As")
	core.AssertContains(t, label, "ugly")
}

func TestErrors_Join_Good(t *core.T) {
	// Join
	ax7Variant := "Join:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Join:good"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "good")
}

func TestErrors_Join_Bad(t *core.T) {
	// Join
	ax7Variant := "Join:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Join:bad"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "bad")
}

func TestErrors_Join_Ugly(t *core.T) {
	// Join
	ax7Variant := "Join:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Join:ugly"
	core.AssertContains(t, label, "Join")
	core.AssertContains(t, label, "ugly")
}
