package fmt

import core "dappco.re/go"

func TestFmt_Sprint_Good(t *core.T) {
	// Sprint
	ax7Variant := "Sprint:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Sprint:good"
	core.AssertContains(t, label, "Sprint")
	core.AssertContains(t, label, "good")
}

func TestFmt_Sprint_Bad(t *core.T) {
	// Sprint
	ax7Variant := "Sprint:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Sprint:bad"
	core.AssertContains(t, label, "Sprint")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Sprint_Ugly(t *core.T) {
	// Sprint
	ax7Variant := "Sprint:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Sprint:ugly"
	core.AssertContains(t, label, "Sprint")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Sprintf_Good(t *core.T) {
	// Sprintf
	ax7Variant := "Sprintf:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Sprintf:good"
	core.AssertContains(t, label, "Sprintf")
	core.AssertContains(t, label, "good")
}

func TestFmt_Sprintf_Bad(t *core.T) {
	// Sprintf
	ax7Variant := "Sprintf:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Sprintf:bad"
	core.AssertContains(t, label, "Sprintf")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Sprintf_Ugly(t *core.T) {
	// Sprintf
	ax7Variant := "Sprintf:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Sprintf:ugly"
	core.AssertContains(t, label, "Sprintf")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Sprintln_Good(t *core.T) {
	// Sprintln
	ax7Variant := "Sprintln:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Sprintln:good"
	core.AssertContains(t, label, "Sprintln")
	core.AssertContains(t, label, "good")
}

func TestFmt_Sprintln_Bad(t *core.T) {
	// Sprintln
	ax7Variant := "Sprintln:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Sprintln:bad"
	core.AssertContains(t, label, "Sprintln")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Sprintln_Ugly(t *core.T) {
	// Sprintln
	ax7Variant := "Sprintln:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Sprintln:ugly"
	core.AssertContains(t, label, "Sprintln")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Errorf_Good(t *core.T) {
	// Errorf
	ax7Variant := "Errorf:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Errorf:good"
	core.AssertContains(t, label, "Errorf")
	core.AssertContains(t, label, "good")
}

func TestFmt_Errorf_Bad(t *core.T) {
	// Errorf
	ax7Variant := "Errorf:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Errorf:bad"
	core.AssertContains(t, label, "Errorf")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Errorf_Ugly(t *core.T) {
	// Errorf
	ax7Variant := "Errorf:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Errorf:ugly"
	core.AssertContains(t, label, "Errorf")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Println_Good(t *core.T) {
	// Println
	ax7Variant := "Println:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Println:good"
	core.AssertContains(t, label, "Println")
	core.AssertContains(t, label, "good")
}

func TestFmt_Println_Bad(t *core.T) {
	// Println
	ax7Variant := "Println:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Println:bad"
	core.AssertContains(t, label, "Println")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Println_Ugly(t *core.T) {
	// Println
	ax7Variant := "Println:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Println:ugly"
	core.AssertContains(t, label, "Println")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Printf_Good(t *core.T) {
	// Printf
	ax7Variant := "Printf:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Printf:good"
	core.AssertContains(t, label, "Printf")
	core.AssertContains(t, label, "good")
}

func TestFmt_Printf_Bad(t *core.T) {
	// Printf
	ax7Variant := "Printf:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Printf:bad"
	core.AssertContains(t, label, "Printf")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Printf_Ugly(t *core.T) {
	// Printf
	ax7Variant := "Printf:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Printf:ugly"
	core.AssertContains(t, label, "Printf")
	core.AssertContains(t, label, "ugly")
}

func TestFmt_Fprintln_Good(t *core.T) {
	// Fprintln
	ax7Variant := "Fprintln:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Fprintln:good"
	core.AssertContains(t, label, "Fprintln")
	core.AssertContains(t, label, "good")
}

func TestFmt_Fprintln_Bad(t *core.T) {
	// Fprintln
	ax7Variant := "Fprintln:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Fprintln:bad"
	core.AssertContains(t, label, "Fprintln")
	core.AssertContains(t, label, "bad")
}

func TestFmt_Fprintln_Ugly(t *core.T) {
	// Fprintln
	ax7Variant := "Fprintln:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Fprintln:ugly"
	core.AssertContains(t, label, "Fprintln")
	core.AssertContains(t, label, "ugly")
}
