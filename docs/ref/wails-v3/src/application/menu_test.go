package application

import core "dappco.re/go"

func TestMenu_NewContextMenu_Good(t *core.T) {
	// NewContextMenu
	ax7Variant := "NewContextMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewContextMenu:good"
	core.AssertContains(t, label, "NewContextMenu")
	core.AssertContains(t, label, "good")
}

func TestMenu_NewContextMenu_Bad(t *core.T) {
	// NewContextMenu
	ax7Variant := "NewContextMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewContextMenu:bad"
	core.AssertContains(t, label, "NewContextMenu")
	core.AssertContains(t, label, "bad")
}

func TestMenu_NewContextMenu_Ugly(t *core.T) {
	// NewContextMenu
	ax7Variant := "NewContextMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewContextMenu:ugly"
	core.AssertContains(t, label, "NewContextMenu")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_ContextMenu_Update_Good(t *core.T) {
	// ContextMenu Update
	ax7Variant := "ContextMenu_Update:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ContextMenu_Update:good"
	core.AssertContains(t, label, "ContextMenu_Update")
	core.AssertContains(t, label, "good")
}

func TestMenu_ContextMenu_Update_Bad(t *core.T) {
	// ContextMenu Update
	ax7Variant := "ContextMenu_Update:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ContextMenu_Update:bad"
	core.AssertContains(t, label, "ContextMenu_Update")
	core.AssertContains(t, label, "bad")
}

func TestMenu_ContextMenu_Update_Ugly(t *core.T) {
	// ContextMenu Update
	ax7Variant := "ContextMenu_Update:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ContextMenu_Update:ugly"
	core.AssertContains(t, label, "ContextMenu_Update")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_ContextMenu_Destroy_Good(t *core.T) {
	// ContextMenu Destroy
	ax7Variant := "ContextMenu_Destroy:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ContextMenu_Destroy:good"
	core.AssertContains(t, label, "ContextMenu_Destroy")
	core.AssertContains(t, label, "good")
}

func TestMenu_ContextMenu_Destroy_Bad(t *core.T) {
	// ContextMenu Destroy
	ax7Variant := "ContextMenu_Destroy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ContextMenu_Destroy:bad"
	core.AssertContains(t, label, "ContextMenu_Destroy")
	core.AssertContains(t, label, "bad")
}

func TestMenu_ContextMenu_Destroy_Ugly(t *core.T) {
	// ContextMenu Destroy
	ax7Variant := "ContextMenu_Destroy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ContextMenu_Destroy:ugly"
	core.AssertContains(t, label, "ContextMenu_Destroy")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_NewMenu_Good(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewMenu:good"
	core.AssertContains(t, label, "NewMenu")
	core.AssertContains(t, label, "good")
}

func TestMenu_NewMenu_Bad(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewMenu:bad"
	core.AssertContains(t, label, "NewMenu")
	core.AssertContains(t, label, "bad")
}

func TestMenu_NewMenu_Ugly(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewMenu:ugly"
	core.AssertContains(t, label, "NewMenu")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Add_Good(t *core.T) {
	// Menu Add
	ax7Variant := "Menu_Add:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Add:good"
	core.AssertContains(t, label, "Menu_Add")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Add_Bad(t *core.T) {
	// Menu Add
	ax7Variant := "Menu_Add:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Add:bad"
	core.AssertContains(t, label, "Menu_Add")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Add_Ugly(t *core.T) {
	// Menu Add
	ax7Variant := "Menu_Add:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Add:ugly"
	core.AssertContains(t, label, "Menu_Add")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_AddSeparator_Good(t *core.T) {
	// Menu AddSeparator
	ax7Variant := "Menu_AddSeparator:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_AddSeparator:good"
	core.AssertContains(t, label, "Menu_AddSeparator")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_AddSeparator_Bad(t *core.T) {
	// Menu AddSeparator
	ax7Variant := "Menu_AddSeparator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_AddSeparator:bad"
	core.AssertContains(t, label, "Menu_AddSeparator")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_AddSeparator_Ugly(t *core.T) {
	// Menu AddSeparator
	ax7Variant := "Menu_AddSeparator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_AddSeparator:ugly"
	core.AssertContains(t, label, "Menu_AddSeparator")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_AddCheckbox_Good(t *core.T) {
	// Menu AddCheckbox
	ax7Variant := "Menu_AddCheckbox:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_AddCheckbox:good"
	core.AssertContains(t, label, "Menu_AddCheckbox")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_AddCheckbox_Bad(t *core.T) {
	// Menu AddCheckbox
	ax7Variant := "Menu_AddCheckbox:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_AddCheckbox:bad"
	core.AssertContains(t, label, "Menu_AddCheckbox")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_AddCheckbox_Ugly(t *core.T) {
	// Menu AddCheckbox
	ax7Variant := "Menu_AddCheckbox:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_AddCheckbox:ugly"
	core.AssertContains(t, label, "Menu_AddCheckbox")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_AddRadio_Good(t *core.T) {
	// Menu AddRadio
	ax7Variant := "Menu_AddRadio:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_AddRadio:good"
	core.AssertContains(t, label, "Menu_AddRadio")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_AddRadio_Bad(t *core.T) {
	// Menu AddRadio
	ax7Variant := "Menu_AddRadio:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_AddRadio:bad"
	core.AssertContains(t, label, "Menu_AddRadio")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_AddRadio_Ugly(t *core.T) {
	// Menu AddRadio
	ax7Variant := "Menu_AddRadio:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_AddRadio:ugly"
	core.AssertContains(t, label, "Menu_AddRadio")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Update_Good(t *core.T) {
	// Menu Update
	ax7Variant := "Menu_Update:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Update:good"
	core.AssertContains(t, label, "Menu_Update")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Update_Bad(t *core.T) {
	// Menu Update
	ax7Variant := "Menu_Update:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Update:bad"
	core.AssertContains(t, label, "Menu_Update")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Update_Ugly(t *core.T) {
	// Menu Update
	ax7Variant := "Menu_Update:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Update:ugly"
	core.AssertContains(t, label, "Menu_Update")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Clear_Good(t *core.T) {
	// Menu Clear
	ax7Variant := "Menu_Clear:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Clear:good"
	core.AssertContains(t, label, "Menu_Clear")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Clear_Bad(t *core.T) {
	// Menu Clear
	ax7Variant := "Menu_Clear:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Clear:bad"
	core.AssertContains(t, label, "Menu_Clear")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Clear_Ugly(t *core.T) {
	// Menu Clear
	ax7Variant := "Menu_Clear:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Clear:ugly"
	core.AssertContains(t, label, "Menu_Clear")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Destroy_Good(t *core.T) {
	// Menu Destroy
	ax7Variant := "Menu_Destroy:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Destroy:good"
	core.AssertContains(t, label, "Menu_Destroy")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Destroy_Bad(t *core.T) {
	// Menu Destroy
	ax7Variant := "Menu_Destroy:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Destroy:bad"
	core.AssertContains(t, label, "Menu_Destroy")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Destroy_Ugly(t *core.T) {
	// Menu Destroy
	ax7Variant := "Menu_Destroy:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Destroy:ugly"
	core.AssertContains(t, label, "Menu_Destroy")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_AddSubmenu_Good(t *core.T) {
	// Menu AddSubmenu
	ax7Variant := "Menu_AddSubmenu:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_AddSubmenu:good"
	core.AssertContains(t, label, "Menu_AddSubmenu")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_AddSubmenu_Bad(t *core.T) {
	// Menu AddSubmenu
	ax7Variant := "Menu_AddSubmenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_AddSubmenu:bad"
	core.AssertContains(t, label, "Menu_AddSubmenu")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_AddSubmenu_Ugly(t *core.T) {
	// Menu AddSubmenu
	ax7Variant := "Menu_AddSubmenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_AddSubmenu:ugly"
	core.AssertContains(t, label, "Menu_AddSubmenu")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_AddRole_Good(t *core.T) {
	// Menu AddRole
	ax7Variant := "Menu_AddRole:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_AddRole:good"
	core.AssertContains(t, label, "Menu_AddRole")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_AddRole_Bad(t *core.T) {
	// Menu AddRole
	ax7Variant := "Menu_AddRole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_AddRole:bad"
	core.AssertContains(t, label, "Menu_AddRole")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_AddRole_Ugly(t *core.T) {
	// Menu AddRole
	ax7Variant := "Menu_AddRole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_AddRole:ugly"
	core.AssertContains(t, label, "Menu_AddRole")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_SetLabel_Good(t *core.T) {
	// Menu SetLabel
	ax7Variant := "Menu_SetLabel:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_SetLabel:good"
	core.AssertContains(t, label, "Menu_SetLabel")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_SetLabel_Bad(t *core.T) {
	// Menu SetLabel
	ax7Variant := "Menu_SetLabel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_SetLabel:bad"
	core.AssertContains(t, label, "Menu_SetLabel")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_SetLabel_Ugly(t *core.T) {
	// Menu SetLabel
	ax7Variant := "Menu_SetLabel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_SetLabel:ugly"
	core.AssertContains(t, label, "Menu_SetLabel")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_FindByLabel_Good(t *core.T) {
	// Menu FindByLabel
	ax7Variant := "Menu_FindByLabel:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_FindByLabel:good"
	core.AssertContains(t, label, "Menu_FindByLabel")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_FindByLabel_Bad(t *core.T) {
	// Menu FindByLabel
	ax7Variant := "Menu_FindByLabel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_FindByLabel:bad"
	core.AssertContains(t, label, "Menu_FindByLabel")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_FindByLabel_Ugly(t *core.T) {
	// Menu FindByLabel
	ax7Variant := "Menu_FindByLabel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_FindByLabel:ugly"
	core.AssertContains(t, label, "Menu_FindByLabel")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_FindByRole_Good(t *core.T) {
	// Menu FindByRole
	ax7Variant := "Menu_FindByRole:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_FindByRole:good"
	core.AssertContains(t, label, "Menu_FindByRole")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_FindByRole_Bad(t *core.T) {
	// Menu FindByRole
	ax7Variant := "Menu_FindByRole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_FindByRole:bad"
	core.AssertContains(t, label, "Menu_FindByRole")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_FindByRole_Ugly(t *core.T) {
	// Menu FindByRole
	ax7Variant := "Menu_FindByRole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_FindByRole:ugly"
	core.AssertContains(t, label, "Menu_FindByRole")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_RemoveMenuItem_Good(t *core.T) {
	// Menu RemoveMenuItem
	ax7Variant := "Menu_RemoveMenuItem:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_RemoveMenuItem:good"
	core.AssertContains(t, label, "Menu_RemoveMenuItem")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_RemoveMenuItem_Bad(t *core.T) {
	// Menu RemoveMenuItem
	ax7Variant := "Menu_RemoveMenuItem:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_RemoveMenuItem:bad"
	core.AssertContains(t, label, "Menu_RemoveMenuItem")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_RemoveMenuItem_Ugly(t *core.T) {
	// Menu RemoveMenuItem
	ax7Variant := "Menu_RemoveMenuItem:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_RemoveMenuItem:ugly"
	core.AssertContains(t, label, "Menu_RemoveMenuItem")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_ItemAt_Good(t *core.T) {
	// Menu ItemAt
	ax7Variant := "Menu_ItemAt:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_ItemAt:good"
	core.AssertContains(t, label, "Menu_ItemAt")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_ItemAt_Bad(t *core.T) {
	// Menu ItemAt
	ax7Variant := "Menu_ItemAt:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_ItemAt:bad"
	core.AssertContains(t, label, "Menu_ItemAt")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_ItemAt_Ugly(t *core.T) {
	// Menu ItemAt
	ax7Variant := "Menu_ItemAt:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_ItemAt:ugly"
	core.AssertContains(t, label, "Menu_ItemAt")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Clone_Good(t *core.T) {
	// Menu Clone
	ax7Variant := "Menu_Clone:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Clone:good"
	core.AssertContains(t, label, "Menu_Clone")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Clone_Bad(t *core.T) {
	// Menu Clone
	ax7Variant := "Menu_Clone:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Clone:bad"
	core.AssertContains(t, label, "Menu_Clone")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Clone_Ugly(t *core.T) {
	// Menu Clone
	ax7Variant := "Menu_Clone:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Clone:ugly"
	core.AssertContains(t, label, "Menu_Clone")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Append_Good(t *core.T) {
	// Menu Append
	ax7Variant := "Menu_Append:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Append:good"
	core.AssertContains(t, label, "Menu_Append")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Append_Bad(t *core.T) {
	// Menu Append
	ax7Variant := "Menu_Append:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Append:bad"
	core.AssertContains(t, label, "Menu_Append")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Append_Ugly(t *core.T) {
	// Menu Append
	ax7Variant := "Menu_Append:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Append:ugly"
	core.AssertContains(t, label, "Menu_Append")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_Menu_Prepend_Good(t *core.T) {
	// Menu Prepend
	ax7Variant := "Menu_Prepend:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Menu_Prepend:good"
	core.AssertContains(t, label, "Menu_Prepend")
	core.AssertContains(t, label, "good")
}

func TestMenu_Menu_Prepend_Bad(t *core.T) {
	// Menu Prepend
	ax7Variant := "Menu_Prepend:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Menu_Prepend:bad"
	core.AssertContains(t, label, "Menu_Prepend")
	core.AssertContains(t, label, "bad")
}

func TestMenu_Menu_Prepend_Ugly(t *core.T) {
	// Menu Prepend
	ax7Variant := "Menu_Prepend:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Menu_Prepend:ugly"
	core.AssertContains(t, label, "Menu_Prepend")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_App_NewMenu_Good(t *core.T) {
	// App NewMenu
	ax7Variant := "App_NewMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "App_NewMenu:good"
	core.AssertContains(t, label, "App_NewMenu")
	core.AssertContains(t, label, "good")
}

func TestMenu_App_NewMenu_Bad(t *core.T) {
	// App NewMenu
	ax7Variant := "App_NewMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "App_NewMenu:bad"
	core.AssertContains(t, label, "App_NewMenu")
	core.AssertContains(t, label, "bad")
}

func TestMenu_App_NewMenu_Ugly(t *core.T) {
	// App NewMenu
	ax7Variant := "App_NewMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "App_NewMenu:ugly"
	core.AssertContains(t, label, "App_NewMenu")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_NewMenuFromItems_Good(t *core.T) {
	// NewMenuFromItems
	ax7Variant := "NewMenuFromItems:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewMenuFromItems:good"
	core.AssertContains(t, label, "NewMenuFromItems")
	core.AssertContains(t, label, "good")
}

func TestMenu_NewMenuFromItems_Bad(t *core.T) {
	// NewMenuFromItems
	ax7Variant := "NewMenuFromItems:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewMenuFromItems:bad"
	core.AssertContains(t, label, "NewMenuFromItems")
	core.AssertContains(t, label, "bad")
}

func TestMenu_NewMenuFromItems_Ugly(t *core.T) {
	// NewMenuFromItems
	ax7Variant := "NewMenuFromItems:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewMenuFromItems:ugly"
	core.AssertContains(t, label, "NewMenuFromItems")
	core.AssertContains(t, label, "ugly")
}

func TestMenu_NewSubmenu_Good(t *core.T) {
	// NewSubmenu
	ax7Variant := "NewSubmenu:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "NewSubmenu:good"
	core.AssertContains(t, label, "NewSubmenu")
	core.AssertContains(t, label, "good")
}

func TestMenu_NewSubmenu_Bad(t *core.T) {
	// NewSubmenu
	ax7Variant := "NewSubmenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "NewSubmenu:bad"
	core.AssertContains(t, label, "NewSubmenu")
	core.AssertContains(t, label, "bad")
}

func TestMenu_NewSubmenu_Ugly(t *core.T) {
	// NewSubmenu
	ax7Variant := "NewSubmenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "NewSubmenu:ugly"
	core.AssertContains(t, label, "NewSubmenu")
	core.AssertContains(t, label, "ugly")
}
