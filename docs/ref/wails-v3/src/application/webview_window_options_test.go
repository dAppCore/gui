package application

import core "dappco.re/go"

func TestWebviewWindowOptions_NewRGBA_Good(t *core.T) {
	// NewRGBA
	ax7Variant := "NewRGBA:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewRGBA:good"
	core.AssertContains(t, label, "NewRGBA")
	core.AssertContains(t, label, "good")
}

func TestWebviewWindowOptions_NewRGBA_Bad(t *core.T) {
	// NewRGBA
	ax7Variant := "NewRGBA:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewRGBA:bad"
	core.AssertContains(t, label, "NewRGBA")
	core.AssertContains(t, label, "bad")
}

func TestWebviewWindowOptions_NewRGBA_Ugly(t *core.T) {
	// NewRGBA
	ax7Variant := "NewRGBA:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewRGBA:ugly"
	core.AssertContains(t, label, "NewRGBA")
	core.AssertContains(t, label, "ugly")
}

func TestWebviewWindowOptions_NewRGB_Good(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewRGB:good"
	core.AssertContains(t, label, "NewRGB")
	core.AssertContains(t, label, "good")
}

func TestWebviewWindowOptions_NewRGB_Bad(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewRGB:bad"
	core.AssertContains(t, label, "NewRGB")
	core.AssertContains(t, label, "bad")
}

func TestWebviewWindowOptions_NewRGB_Ugly(t *core.T) {
	// NewRGB
	ax7Variant := "NewRGB:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewRGB:ugly"
	core.AssertContains(t, label, "NewRGB")
	core.AssertContains(t, label, "ugly")
}

func TestWebviewWindowOptions_NewRGBPtr_Good(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewRGBPtr:good"
	core.AssertContains(t, label, "NewRGBPtr")
	core.AssertContains(t, label, "good")
}

func TestWebviewWindowOptions_NewRGBPtr_Bad(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewRGBPtr:bad"
	core.AssertContains(t, label, "NewRGBPtr")
	core.AssertContains(t, label, "bad")
}

func TestWebviewWindowOptions_NewRGBPtr_Ugly(t *core.T) {
	// NewRGBPtr
	ax7Variant := "NewRGBPtr:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewRGBPtr:ugly"
	core.AssertContains(t, label, "NewRGBPtr")
	core.AssertContains(t, label, "ugly")
}
