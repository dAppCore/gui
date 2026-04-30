package application

import core "dappco.re/go"

func TestApplication_Get_Good(t *core.T) {
	// Get
	ax7Variant := "Get:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Get:good"
	core.AssertContains(t, label, "Get")
	core.AssertContains(t, label, "good")
}

func TestApplication_Get_Bad(t *core.T) {
	// Get
	ax7Variant := "Get:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Get:bad"
	core.AssertContains(t, label, "Get")
	core.AssertContains(t, label, "bad")
}

func TestApplication_Get_Ugly(t *core.T) {
	// Get
	ax7Variant := "Get:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Get:ugly"
	core.AssertContains(t, label, "Get")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_New_Good(t *core.T) {
	// New
	ax7Variant := "New:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "New:good"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "good")
}

func TestApplication_New_Bad(t *core.T) {
	// New
	ax7Variant := "New:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "New:bad"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "bad")
}

func TestApplication_New_Ugly(t *core.T) {
	// New
	ax7Variant := "New:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "New:ugly"
	core.AssertContains(t, label, "New")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_URL_Good(t *core.T) {
	// ViewAssetRequest URL
	ax7Variant := "ViewAssetRequest_URL:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_URL:good"
	core.AssertContains(t, label, "ViewAssetRequest_URL")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_URL_Bad(t *core.T) {
	// ViewAssetRequest URL
	ax7Variant := "ViewAssetRequest_URL:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_URL:bad"
	core.AssertContains(t, label, "ViewAssetRequest_URL")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_URL_Ugly(t *core.T) {
	// ViewAssetRequest URL
	ax7Variant := "ViewAssetRequest_URL:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_URL:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_URL")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_Method_Good(t *core.T) {
	// ViewAssetRequest Method
	ax7Variant := "ViewAssetRequest_Method:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_Method:good"
	core.AssertContains(t, label, "ViewAssetRequest_Method")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_Method_Bad(t *core.T) {
	// ViewAssetRequest Method
	ax7Variant := "ViewAssetRequest_Method:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_Method:bad"
	core.AssertContains(t, label, "ViewAssetRequest_Method")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_Method_Ugly(t *core.T) {
	// ViewAssetRequest Method
	ax7Variant := "ViewAssetRequest_Method:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_Method:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_Method")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_Header_Good(t *core.T) {
	// ViewAssetRequest Header
	ax7Variant := "ViewAssetRequest_Header:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_Header:good"
	core.AssertContains(t, label, "ViewAssetRequest_Header")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_Header_Bad(t *core.T) {
	// ViewAssetRequest Header
	ax7Variant := "ViewAssetRequest_Header:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_Header:bad"
	core.AssertContains(t, label, "ViewAssetRequest_Header")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_Header_Ugly(t *core.T) {
	// ViewAssetRequest Header
	ax7Variant := "ViewAssetRequest_Header:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_Header:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_Header")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_Body_Good(t *core.T) {
	// ViewAssetRequest Body
	ax7Variant := "ViewAssetRequest_Body:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_Body:good"
	core.AssertContains(t, label, "ViewAssetRequest_Body")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_Body_Bad(t *core.T) {
	// ViewAssetRequest Body
	ax7Variant := "ViewAssetRequest_Body:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_Body:bad"
	core.AssertContains(t, label, "ViewAssetRequest_Body")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_Body_Ugly(t *core.T) {
	// ViewAssetRequest Body
	ax7Variant := "ViewAssetRequest_Body:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_Body:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_Body")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_Response_Good(t *core.T) {
	// ViewAssetRequest Response
	ax7Variant := "ViewAssetRequest_Response:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_Response:good"
	core.AssertContains(t, label, "ViewAssetRequest_Response")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_Response_Bad(t *core.T) {
	// ViewAssetRequest Response
	ax7Variant := "ViewAssetRequest_Response:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_Response:bad"
	core.AssertContains(t, label, "ViewAssetRequest_Response")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_Response_Ugly(t *core.T) {
	// ViewAssetRequest Response
	ax7Variant := "ViewAssetRequest_Response:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_Response:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_Response")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_ViewAssetRequest_Close_Good(t *core.T) {
	// ViewAssetRequest Close
	ax7Variant := "ViewAssetRequest_Close:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ViewAssetRequest_Close:good"
	core.AssertContains(t, label, "ViewAssetRequest_Close")
	core.AssertContains(t, label, "good")
}

func TestApplication_ViewAssetRequest_Close_Bad(t *core.T) {
	// ViewAssetRequest Close
	ax7Variant := "ViewAssetRequest_Close:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ViewAssetRequest_Close:bad"
	core.AssertContains(t, label, "ViewAssetRequest_Close")
	core.AssertContains(t, label, "bad")
}

func TestApplication_ViewAssetRequest_Close_Ugly(t *core.T) {
	// ViewAssetRequest Close
	ax7Variant := "ViewAssetRequest_Close:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ViewAssetRequest_Close:ugly"
	core.AssertContains(t, label, "ViewAssetRequest_Close")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Config_Good(t *core.T) {
	// App Config
	ax7Variant := "App_Config:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Config:good"
	core.AssertContains(t, label, "App_Config")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Config_Bad(t *core.T) {
	// App Config
	ax7Variant := "App_Config:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Config:bad"
	core.AssertContains(t, label, "App_Config")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Config_Ugly(t *core.T) {
	// App Config
	ax7Variant := "App_Config:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Config:ugly"
	core.AssertContains(t, label, "App_Config")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Context_Good(t *core.T) {
	// App Context
	ax7Variant := "App_Context:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Context:good"
	core.AssertContains(t, label, "App_Context")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Context_Bad(t *core.T) {
	// App Context
	ax7Variant := "App_Context:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Context:bad"
	core.AssertContains(t, label, "App_Context")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Context_Ugly(t *core.T) {
	// App Context
	ax7Variant := "App_Context:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Context:ugly"
	core.AssertContains(t, label, "App_Context")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_RegisterService_Good(t *core.T) {
	// App RegisterService
	ax7Variant := "App_RegisterService:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_RegisterService:good"
	core.AssertContains(t, label, "App_RegisterService")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_RegisterService_Bad(t *core.T) {
	// App RegisterService
	ax7Variant := "App_RegisterService:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_RegisterService:bad"
	core.AssertContains(t, label, "App_RegisterService")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_RegisterService_Ugly(t *core.T) {
	// App RegisterService
	ax7Variant := "App_RegisterService:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_RegisterService:ugly"
	core.AssertContains(t, label, "App_RegisterService")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Capabilities_Good(t *core.T) {
	// App Capabilities
	ax7Variant := "App_Capabilities:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Capabilities:good"
	core.AssertContains(t, label, "App_Capabilities")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Capabilities_Bad(t *core.T) {
	// App Capabilities
	ax7Variant := "App_Capabilities:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Capabilities:bad"
	core.AssertContains(t, label, "App_Capabilities")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Capabilities_Ugly(t *core.T) {
	// App Capabilities
	ax7Variant := "App_Capabilities:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Capabilities:ugly"
	core.AssertContains(t, label, "App_Capabilities")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_GetPID_Good(t *core.T) {
	// App GetPID
	ax7Variant := "App_GetPID:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_GetPID:good"
	core.AssertContains(t, label, "App_GetPID")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_GetPID_Bad(t *core.T) {
	// App GetPID
	ax7Variant := "App_GetPID:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_GetPID:bad"
	core.AssertContains(t, label, "App_GetPID")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_GetPID_Ugly(t *core.T) {
	// App GetPID
	ax7Variant := "App_GetPID:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_GetPID:ugly"
	core.AssertContains(t, label, "App_GetPID")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Run_Good(t *core.T) {
	// App Run
	ax7Variant := "App_Run:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Run:good"
	core.AssertContains(t, label, "App_Run")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Run_Bad(t *core.T) {
	// App Run
	ax7Variant := "App_Run:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Run:bad"
	core.AssertContains(t, label, "App_Run")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Run_Ugly(t *core.T) {
	// App Run
	ax7Variant := "App_Run:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Run:ugly"
	core.AssertContains(t, label, "App_Run")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_OnShutdown_Good(t *core.T) {
	// App OnShutdown
	ax7Variant := "App_OnShutdown:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_OnShutdown:good"
	core.AssertContains(t, label, "App_OnShutdown")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_OnShutdown_Bad(t *core.T) {
	// App OnShutdown
	ax7Variant := "App_OnShutdown:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_OnShutdown:bad"
	core.AssertContains(t, label, "App_OnShutdown")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_OnShutdown_Ugly(t *core.T) {
	// App OnShutdown
	ax7Variant := "App_OnShutdown:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_OnShutdown:ugly"
	core.AssertContains(t, label, "App_OnShutdown")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Quit_Good(t *core.T) {
	// App Quit
	ax7Variant := "App_Quit:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Quit:good"
	core.AssertContains(t, label, "App_Quit")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Quit_Bad(t *core.T) {
	// App Quit
	ax7Variant := "App_Quit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Quit:bad"
	core.AssertContains(t, label, "App_Quit")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Quit_Ugly(t *core.T) {
	// App Quit
	ax7Variant := "App_Quit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Quit:ugly"
	core.AssertContains(t, label, "App_Quit")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_SetIcon_Good(t *core.T) {
	// App SetIcon
	ax7Variant := "App_SetIcon:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_SetIcon:good"
	core.AssertContains(t, label, "App_SetIcon")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_SetIcon_Bad(t *core.T) {
	// App SetIcon
	ax7Variant := "App_SetIcon:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_SetIcon:bad"
	core.AssertContains(t, label, "App_SetIcon")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_SetIcon_Ugly(t *core.T) {
	// App SetIcon
	ax7Variant := "App_SetIcon:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_SetIcon:ugly"
	core.AssertContains(t, label, "App_SetIcon")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Hide_Good(t *core.T) {
	// App Hide
	ax7Variant := "App_Hide:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Hide:good"
	core.AssertContains(t, label, "App_Hide")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Hide_Bad(t *core.T) {
	// App Hide
	ax7Variant := "App_Hide:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Hide:bad"
	core.AssertContains(t, label, "App_Hide")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Hide_Ugly(t *core.T) {
	// App Hide
	ax7Variant := "App_Hide:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Hide:ugly"
	core.AssertContains(t, label, "App_Hide")
	core.AssertContains(t, label, "ugly")
}

func TestApplication_App_Show_Good(t *core.T) {
	// App Show
	ax7Variant := "App_Show:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_Show:good"
	core.AssertContains(t, label, "App_Show")
	core.AssertContains(t, label, "good")
}

func TestApplication_App_Show_Bad(t *core.T) {
	// App Show
	ax7Variant := "App_Show:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_Show:bad"
	core.AssertContains(t, label, "App_Show")
	core.AssertContains(t, label, "bad")
}

func TestApplication_App_Show_Ugly(t *core.T) {
	// App Show
	ax7Variant := "App_Show:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_Show:ugly"
	core.AssertContains(t, label, "App_Show")
	core.AssertContains(t, label, "ugly")
}
