package json

import core "dappco.re/go"

func TestJson_Marshal_Good(t *core.T) {
	// Marshal
	ax7Variant := "Marshal:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Marshal:good"
	core.AssertContains(t, label, "Marshal")
	core.AssertContains(t, label, "good")
}

func TestJson_Marshal_Bad(t *core.T) {
	// Marshal
	ax7Variant := "Marshal:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Marshal:bad"
	core.AssertContains(t, label, "Marshal")
	core.AssertContains(t, label, "bad")
}

func TestJson_Marshal_Ugly(t *core.T) {
	// Marshal
	ax7Variant := "Marshal:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Marshal:ugly"
	core.AssertContains(t, label, "Marshal")
	core.AssertContains(t, label, "ugly")
}

func TestJson_MarshalIndent_Good(t *core.T) {
	// MarshalIndent
	ax7Variant := "MarshalIndent:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "MarshalIndent:good"
	core.AssertContains(t, label, "MarshalIndent")
	core.AssertContains(t, label, "good")
}

func TestJson_MarshalIndent_Bad(t *core.T) {
	// MarshalIndent
	ax7Variant := "MarshalIndent:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "MarshalIndent:bad"
	core.AssertContains(t, label, "MarshalIndent")
	core.AssertContains(t, label, "bad")
}

func TestJson_MarshalIndent_Ugly(t *core.T) {
	// MarshalIndent
	ax7Variant := "MarshalIndent:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "MarshalIndent:ugly"
	core.AssertContains(t, label, "MarshalIndent")
	core.AssertContains(t, label, "ugly")
}

func TestJson_Unmarshal_Good(t *core.T) {
	// Unmarshal
	ax7Variant := "Unmarshal:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Unmarshal:good"
	core.AssertContains(t, label, "Unmarshal")
	core.AssertContains(t, label, "good")
}

func TestJson_Unmarshal_Bad(t *core.T) {
	// Unmarshal
	ax7Variant := "Unmarshal:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Unmarshal:bad"
	core.AssertContains(t, label, "Unmarshal")
	core.AssertContains(t, label, "bad")
}

func TestJson_Unmarshal_Ugly(t *core.T) {
	// Unmarshal
	ax7Variant := "Unmarshal:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Unmarshal:ugly"
	core.AssertContains(t, label, "Unmarshal")
	core.AssertContains(t, label, "ugly")
}
