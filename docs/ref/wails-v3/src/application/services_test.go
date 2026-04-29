package application

import core "dappco.re/go"

func TestServices_NewService_Good(t *core.T) {
	// NewService
	ax7Variant := "NewService:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewService:good"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "good")
}

func TestServices_NewService_Bad(t *core.T) {
	// NewService
	ax7Variant := "NewService:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewService:bad"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "bad")
}

func TestServices_NewService_Ugly(t *core.T) {
	// NewService
	ax7Variant := "NewService:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewService:ugly"
	core.AssertContains(t, label, "NewService")
	core.AssertContains(t, label, "ugly")
}

func TestServices_NewServiceWithOptions_Good(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewServiceWithOptions:good"
	core.AssertContains(t, label, "NewServiceWithOptions")
	core.AssertContains(t, label, "good")
}

func TestServices_NewServiceWithOptions_Bad(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewServiceWithOptions:bad"
	core.AssertContains(t, label, "NewServiceWithOptions")
	core.AssertContains(t, label, "bad")
}

func TestServices_NewServiceWithOptions_Ugly(t *core.T) {
	// NewServiceWithOptions
	ax7Variant := "NewServiceWithOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewServiceWithOptions:ugly"
	core.AssertContains(t, label, "NewServiceWithOptions")
	core.AssertContains(t, label, "ugly")
}

func TestServices_Service_Instance_Good(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Service_Instance:good"
	core.AssertContains(t, label, "Service_Instance")
	core.AssertContains(t, label, "good")
}

func TestServices_Service_Instance_Bad(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Service_Instance:bad"
	core.AssertContains(t, label, "Service_Instance")
	core.AssertContains(t, label, "bad")
}

func TestServices_Service_Instance_Ugly(t *core.T) {
	// Service Instance
	ax7Variant := "Service_Instance:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Service_Instance:ugly"
	core.AssertContains(t, label, "Service_Instance")
	core.AssertContains(t, label, "ugly")
}
