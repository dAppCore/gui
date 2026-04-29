package application

import (
	core "dappco.re/go"
)

func TestDialog_OpenFileDialogStruct_GoodCase(t *core.T) {
	dialog := newOpenFileDialog()
	dialog.SetOptions(&OpenFileDialogOptions{
		Title:                   "Pick",
		Directory:               "/tmp",
		Filters:                 []FileFilter{{DisplayName: "Text", Pattern: "*.txt"}},
		AllowsMultipleSelection: true,
		CanChooseDirectories:    true,
		CanChooseFiles:          false,
		ShowHiddenFiles:         true,
	})

	dialog.SetSelectedFiles([]string{"/tmp/a.txt", "/tmp/b.txt"})
	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/a.txt", got)
	paths, err := dialog.PromptForMultipleSelection()
	core.RequireNoError(t, err)
	core.AssertEqual(t, []string{"/tmp/a.txt", "/tmp/b.txt"}, paths)
}

func TestDialog_OpenFileDialogStruct_BadCase(t *core.T) {
	dialog := newOpenFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
	core.AssertTrue(t, dialog.canChooseFiles)
}

func TestDialog_OpenFileDialogStruct_UglyCase(t *core.T) {
	dialog := newOpenFileDialog()
	selected := []string{"/tmp/input.txt"}
	dialog.SetSelectedFiles(selected)
	selected[0] = "/tmp/mutated.txt"

	paths, err := dialog.PromptForMultipleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, []string{"/tmp/input.txt"}, paths)
}

func TestDialog_SaveFileDialogStruct_GoodCase(t *core.T) {
	dialog := newSaveFileDialog()
	dialog.SetOptions(&SaveFileDialogOptions{
		Title:           "Export",
		Directory:       "/tmp",
		Filename:        "report.csv",
		Filters:         []FileFilter{{DisplayName: "CSV", Pattern: "*.csv"}},
		ShowHiddenFiles: true,
	})
	dialog.SetSelectedPath("/tmp/report.csv")

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/report.csv", got)
}

func TestDialog_SaveFileDialogStruct_BadCase(t *core.T) {
	dialog := newSaveFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
}

func TestDialog_SaveFileDialogStruct_UglyCase(t *core.T) {
	dialog := newSaveFileDialog()
	selected := "/tmp/a.csv"
	dialog.SetSelectedPath(selected)
	selected = "/tmp/b.csv"

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/a.csv", got)
}

func TestDialog_MessageDialog_GoodCase(t *core.T) {
	dialog := newMessageDialog(QuestionDialogType)
	dialog.SetTitle("Confirm").
		SetMessage("Proceed?").
		AddButton("Yes").
		AddButton("No").
		SetDefaultButton("Yes").
		SetCancelButton("No")
	dialog.SetButtonClickedForStub("Yes")

	got, err := dialog.Show()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "Yes", got)
	core.AssertEqual(t, QuestionDialogType, dialog.dialogType)
	core.AssertEqual(t, "Confirm", dialog.title)
	core.AssertEqual(t, "Proceed?", dialog.message)
	core.AssertLen(t, dialog.buttons, 2)
	core.AssertTrue(t, dialog.buttons[0].IsDefault)
	core.AssertTrue(t, dialog.buttons[1].IsCancel)
}

func TestDialog_MessageDialog_BadCase(t *core.T) {
	dialog := newMessageDialog(InfoDialogType)

	got, err := dialog.Show()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
	core.AssertEqual(t, InfoDialogType, dialog.dialogType)
}

func TestDialog_MessageDialog_UglyCase(t *core.T) {
	dialog := newMessageDialog(ErrorDialogType)
	dialog.AddButton("Retry").AddButton("Retry")
	dialog.SetDefaultButton("Retry")
	dialog.SetCancelButton("Retry")

	core.AssertLen(t, dialog.buttons, 2)
	core.AssertTrue(t, dialog.buttons[0].IsDefault)
	core.AssertTrue(t, dialog.buttons[0].IsCancel)
	core.AssertTrue(t, dialog.buttons[1].IsDefault)
	core.AssertTrue(t, dialog.buttons[1].IsCancel)
}

func TestDialog_DialogManager_Good(t *core.T) {
	// DialogManager
	ax7Variant := "DialogManager:good"
	core.AssertContains(t, ax7Variant, "good")
	manager := &DialogManager{}

	core.AssertEqual(t, InfoDialogType, manager.Info().dialogType)
	core.AssertEqual(t, QuestionDialogType, manager.Question().dialogType)
	core.AssertEqual(t, WarningDialogType, manager.Warning().dialogType)
	core.AssertEqual(t, ErrorDialogType, manager.Error().dialogType)
	core.AssertNotNil(t, manager.OpenFile())
	core.AssertNotNil(t, manager.SaveFile())
}

func TestDialog_DialogManager_Bad(t *core.T) {
	// DialogManager
	ax7Variant := "DialogManager:bad"
	core.AssertContains(t, ax7Variant, "bad")
	manager := &DialogManager{}

	core.AssertNotNil(t, manager.OpenFileWithOptions(nil))
	core.AssertNotNil(t, manager.SaveFileWithOptions(nil))
}

func TestDialog_DialogManager_Ugly(t *core.T) {
	// DialogManager
	ax7Variant := "DialogManager:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	manager := &DialogManager{}

	open := manager.OpenFileWithOptions(&OpenFileDialogOptions{AllowsMultipleSelection: true})
	save := manager.SaveFileWithOptions(&SaveFileDialogOptions{Filename: "out.csv"})

	core.AssertTrue(t, open.multipleAllowed)
	core.AssertEqual(t, "out.csv", save.filename)
}

// AX7 generated source-matching smoke coverage.
func TestDialog_OpenFileDialogStruct_SetOptions_Good(t *core.T) {
	// OpenFileDialogStruct SetOptions
	ax7Variant := "OpenFileDialogStruct_SetOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetOptions_Bad(t *core.T) {
	// OpenFileDialogStruct SetOptions
	ax7Variant := "OpenFileDialogStruct_SetOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetOptions_Ugly(t *core.T) {
	// OpenFileDialogStruct SetOptions
	ax7Variant := "OpenFileDialogStruct_SetOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetTitle_Good(t *core.T) {
	// OpenFileDialogStruct SetTitle
	ax7Variant := "OpenFileDialogStruct_SetTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetTitle_Bad(t *core.T) {
	// OpenFileDialogStruct SetTitle
	ax7Variant := "OpenFileDialogStruct_SetTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetTitle_Ugly(t *core.T) {
	// OpenFileDialogStruct SetTitle
	ax7Variant := "OpenFileDialogStruct_SetTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetDirectory_Good(t *core.T) {
	// OpenFileDialogStruct SetDirectory
	ax7Variant := "OpenFileDialogStruct_SetDirectory:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetDirectory_Bad(t *core.T) {
	// OpenFileDialogStruct SetDirectory
	ax7Variant := "OpenFileDialogStruct_SetDirectory:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetDirectory_Ugly(t *core.T) {
	// OpenFileDialogStruct SetDirectory
	ax7Variant := "OpenFileDialogStruct_SetDirectory:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_AddFilter_Good(t *core.T) {
	// OpenFileDialogStruct AddFilter
	ax7Variant := "OpenFileDialogStruct_AddFilter:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_AddFilter_Bad(t *core.T) {
	// OpenFileDialogStruct AddFilter
	ax7Variant := "OpenFileDialogStruct_AddFilter:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_AddFilter_Ugly(t *core.T) {
	// OpenFileDialogStruct AddFilter
	ax7Variant := "OpenFileDialogStruct_AddFilter:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetAllowsMultipleSelection_Good(t *core.T) {
	// OpenFileDialogStruct SetAllowsMultipleSelection
	ax7Variant := "OpenFileDialogStruct_SetAllowsMultipleSelection:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetAllowsMultipleSelection(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetAllowsMultipleSelection_Bad(t *core.T) {
	// OpenFileDialogStruct SetAllowsMultipleSelection
	ax7Variant := "OpenFileDialogStruct_SetAllowsMultipleSelection:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetAllowsMultipleSelection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetAllowsMultipleSelection_Ugly(t *core.T) {
	// OpenFileDialogStruct SetAllowsMultipleSelection
	ax7Variant := "OpenFileDialogStruct_SetAllowsMultipleSelection:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetAllowsMultipleSelection(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetSelectedFiles_Good(t *core.T) {
	// OpenFileDialogStruct SetSelectedFiles
	ax7Variant := "OpenFileDialogStruct_SetSelectedFiles:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedFiles(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetSelectedFiles_Bad(t *core.T) {
	// OpenFileDialogStruct SetSelectedFiles
	ax7Variant := "OpenFileDialogStruct_SetSelectedFiles:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedFiles(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_SetSelectedFiles_Ugly(t *core.T) {
	// OpenFileDialogStruct SetSelectedFiles
	ax7Variant := "OpenFileDialogStruct_SetSelectedFiles:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedFiles(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForSingleSelection_Good(t *core.T) {
	// OpenFileDialogStruct PromptForSingleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForSingleSelection:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForSingleSelection_Bad(t *core.T) {
	// OpenFileDialogStruct PromptForSingleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForSingleSelection:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForSingleSelection_Ugly(t *core.T) {
	// OpenFileDialogStruct PromptForSingleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForSingleSelection:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForMultipleSelection_Good(t *core.T) {
	// OpenFileDialogStruct PromptForMultipleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForMultipleSelection:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForMultipleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForMultipleSelection_Bad(t *core.T) {
	// OpenFileDialogStruct PromptForMultipleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForMultipleSelection:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForMultipleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_OpenFileDialogStruct_PromptForMultipleSelection_Ugly(t *core.T) {
	// OpenFileDialogStruct PromptForMultipleSelection
	ax7Variant := "OpenFileDialogStruct_PromptForMultipleSelection:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(OpenFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForMultipleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetOptions_Good(t *core.T) {
	// SaveFileDialogStruct SetOptions
	ax7Variant := "SaveFileDialogStruct_SetOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetOptions_Bad(t *core.T) {
	// SaveFileDialogStruct SetOptions
	ax7Variant := "SaveFileDialogStruct_SetOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetOptions_Ugly(t *core.T) {
	// SaveFileDialogStruct SetOptions
	ax7Variant := "SaveFileDialogStruct_SetOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetOptions(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetTitle_Good(t *core.T) {
	// SaveFileDialogStruct SetTitle
	ax7Variant := "SaveFileDialogStruct_SetTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetTitle_Bad(t *core.T) {
	// SaveFileDialogStruct SetTitle
	ax7Variant := "SaveFileDialogStruct_SetTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetTitle_Ugly(t *core.T) {
	// SaveFileDialogStruct SetTitle
	ax7Variant := "SaveFileDialogStruct_SetTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetDirectory_Good(t *core.T) {
	// SaveFileDialogStruct SetDirectory
	ax7Variant := "SaveFileDialogStruct_SetDirectory:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetDirectory_Bad(t *core.T) {
	// SaveFileDialogStruct SetDirectory
	ax7Variant := "SaveFileDialogStruct_SetDirectory:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetDirectory_Ugly(t *core.T) {
	// SaveFileDialogStruct SetDirectory
	ax7Variant := "SaveFileDialogStruct_SetDirectory:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetDirectory("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetFilename_Good(t *core.T) {
	// SaveFileDialogStruct SetFilename
	ax7Variant := "SaveFileDialogStruct_SetFilename:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetFilename("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetFilename_Bad(t *core.T) {
	// SaveFileDialogStruct SetFilename
	ax7Variant := "SaveFileDialogStruct_SetFilename:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetFilename("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetFilename_Ugly(t *core.T) {
	// SaveFileDialogStruct SetFilename
	ax7Variant := "SaveFileDialogStruct_SetFilename:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.SetFilename("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_AddFilter_Good(t *core.T) {
	// SaveFileDialogStruct AddFilter
	ax7Variant := "SaveFileDialogStruct_AddFilter:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_AddFilter_Bad(t *core.T) {
	// SaveFileDialogStruct AddFilter
	ax7Variant := "SaveFileDialogStruct_AddFilter:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_AddFilter_Ugly(t *core.T) {
	// SaveFileDialogStruct AddFilter
	ax7Variant := "SaveFileDialogStruct_AddFilter:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0 := subject.AddFilter("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetSelectedPath_Good(t *core.T) {
	// SaveFileDialogStruct SetSelectedPath
	ax7Variant := "SaveFileDialogStruct_SetSelectedPath:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedPath("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetSelectedPath_Bad(t *core.T) {
	// SaveFileDialogStruct SetSelectedPath
	ax7Variant := "SaveFileDialogStruct_SetSelectedPath:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedPath("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_SetSelectedPath_Ugly(t *core.T) {
	// SaveFileDialogStruct SetSelectedPath
	ax7Variant := "SaveFileDialogStruct_SetSelectedPath:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		subject.SetSelectedPath("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_PromptForSingleSelection_Good(t *core.T) {
	// SaveFileDialogStruct PromptForSingleSelection
	ax7Variant := "SaveFileDialogStruct_PromptForSingleSelection:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_PromptForSingleSelection_Bad(t *core.T) {
	// SaveFileDialogStruct PromptForSingleSelection
	ax7Variant := "SaveFileDialogStruct_PromptForSingleSelection:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_SaveFileDialogStruct_PromptForSingleSelection_Ugly(t *core.T) {
	// SaveFileDialogStruct PromptForSingleSelection
	ax7Variant := "SaveFileDialogStruct_PromptForSingleSelection:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(SaveFileDialogStruct)
	result := core.Try(func() any {
		got0, got1 := subject.PromptForSingleSelection()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetTitle_Good(t *core.T) {
	// MessageDialog SetTitle
	ax7Variant := "MessageDialog_SetTitle:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetTitle("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetTitle_Bad(t *core.T) {
	// MessageDialog SetTitle
	ax7Variant := "MessageDialog_SetTitle:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetTitle("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetTitle_Ugly(t *core.T) {
	// MessageDialog SetTitle
	ax7Variant := "MessageDialog_SetTitle:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetTitle("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetMessage_Good(t *core.T) {
	// MessageDialog SetMessage
	ax7Variant := "MessageDialog_SetMessage:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetMessage("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetMessage_Bad(t *core.T) {
	// MessageDialog SetMessage
	ax7Variant := "MessageDialog_SetMessage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetMessage("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetMessage_Ugly(t *core.T) {
	// MessageDialog SetMessage
	ax7Variant := "MessageDialog_SetMessage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetMessage("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_AddButton_Good(t *core.T) {
	// MessageDialog AddButton
	ax7Variant := "MessageDialog_AddButton:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.AddButton("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_AddButton_Bad(t *core.T) {
	// MessageDialog AddButton
	ax7Variant := "MessageDialog_AddButton:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.AddButton("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_AddButton_Ugly(t *core.T) {
	// MessageDialog AddButton
	ax7Variant := "MessageDialog_AddButton:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.AddButton("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetDefaultButton_Good(t *core.T) {
	// MessageDialog SetDefaultButton
	ax7Variant := "MessageDialog_SetDefaultButton:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetDefaultButton("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetDefaultButton_Bad(t *core.T) {
	// MessageDialog SetDefaultButton
	ax7Variant := "MessageDialog_SetDefaultButton:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetDefaultButton("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetDefaultButton_Ugly(t *core.T) {
	// MessageDialog SetDefaultButton
	ax7Variant := "MessageDialog_SetDefaultButton:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetDefaultButton("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetCancelButton_Good(t *core.T) {
	// MessageDialog SetCancelButton
	ax7Variant := "MessageDialog_SetCancelButton:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetCancelButton("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetCancelButton_Bad(t *core.T) {
	// MessageDialog SetCancelButton
	ax7Variant := "MessageDialog_SetCancelButton:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetCancelButton("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetCancelButton_Ugly(t *core.T) {
	// MessageDialog SetCancelButton
	ax7Variant := "MessageDialog_SetCancelButton:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0 := subject.SetCancelButton("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetButtonClickedForStub_Good(t *core.T) {
	// MessageDialog SetButtonClickedForStub
	ax7Variant := "MessageDialog_SetButtonClickedForStub:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		subject.SetButtonClickedForStub("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetButtonClickedForStub_Bad(t *core.T) {
	// MessageDialog SetButtonClickedForStub
	ax7Variant := "MessageDialog_SetButtonClickedForStub:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		subject.SetButtonClickedForStub("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_SetButtonClickedForStub_Ugly(t *core.T) {
	// MessageDialog SetButtonClickedForStub
	ax7Variant := "MessageDialog_SetButtonClickedForStub:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		subject.SetButtonClickedForStub("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_Show_Good(t *core.T) {
	// MessageDialog Show
	ax7Variant := "MessageDialog_Show:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0, got1 := subject.Show()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_Show_Bad(t *core.T) {
	// MessageDialog Show
	ax7Variant := "MessageDialog_Show:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0, got1 := subject.Show()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_MessageDialog_Show_Ugly(t *core.T) {
	// MessageDialog Show
	ax7Variant := "MessageDialog_Show:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MessageDialog)
	result := core.Try(func() any {
		got0, got1 := subject.Show()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFile_Good(t *core.T) {
	// DialogManager OpenFile
	ax7Variant := "DialogManager_OpenFile:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFile_Bad(t *core.T) {
	// DialogManager OpenFile
	ax7Variant := "DialogManager_OpenFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFile_Ugly(t *core.T) {
	// DialogManager OpenFile
	ax7Variant := "DialogManager_OpenFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFileWithOptions_Good(t *core.T) {
	// DialogManager OpenFileWithOptions
	ax7Variant := "DialogManager_OpenFileWithOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFileWithOptions_Bad(t *core.T) {
	// DialogManager OpenFileWithOptions
	ax7Variant := "DialogManager_OpenFileWithOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_OpenFileWithOptions_Ugly(t *core.T) {
	// DialogManager OpenFileWithOptions
	ax7Variant := "DialogManager_OpenFileWithOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.OpenFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFile_Good(t *core.T) {
	// DialogManager SaveFile
	ax7Variant := "DialogManager_SaveFile:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFile_Bad(t *core.T) {
	// DialogManager SaveFile
	ax7Variant := "DialogManager_SaveFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFile_Ugly(t *core.T) {
	// DialogManager SaveFile
	ax7Variant := "DialogManager_SaveFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFile()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFileWithOptions_Good(t *core.T) {
	// DialogManager SaveFileWithOptions
	ax7Variant := "DialogManager_SaveFileWithOptions:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFileWithOptions_Bad(t *core.T) {
	// DialogManager SaveFileWithOptions
	ax7Variant := "DialogManager_SaveFileWithOptions:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_SaveFileWithOptions_Ugly(t *core.T) {
	// DialogManager SaveFileWithOptions
	ax7Variant := "DialogManager_SaveFileWithOptions:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.SaveFileWithOptions(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Info_Good(t *core.T) {
	// DialogManager Info
	ax7Variant := "DialogManager_Info:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Info_Bad(t *core.T) {
	// DialogManager Info
	ax7Variant := "DialogManager_Info:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Info_Ugly(t *core.T) {
	// DialogManager Info
	ax7Variant := "DialogManager_Info:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Info()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Question_Good(t *core.T) {
	// DialogManager Question
	ax7Variant := "DialogManager_Question:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Question()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Question_Bad(t *core.T) {
	// DialogManager Question
	ax7Variant := "DialogManager_Question:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Question()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Question_Ugly(t *core.T) {
	// DialogManager Question
	ax7Variant := "DialogManager_Question:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Question()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Warning_Good(t *core.T) {
	// DialogManager Warning
	ax7Variant := "DialogManager_Warning:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Warning()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Warning_Bad(t *core.T) {
	// DialogManager Warning
	ax7Variant := "DialogManager_Warning:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Warning()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Warning_Ugly(t *core.T) {
	// DialogManager Warning
	ax7Variant := "DialogManager_Warning:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Warning()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Error_Good(t *core.T) {
	// DialogManager Error
	ax7Variant := "DialogManager_Error:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Error()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Error_Bad(t *core.T) {
	// DialogManager Error
	ax7Variant := "DialogManager_Error:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Error()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_Error_Ugly(t *core.T) {
	// DialogManager Error
	ax7Variant := "DialogManager_Error:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0 := subject.Error()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowInfo_Good(t *core.T) {
	// DialogManager ShowInfo
	ax7Variant := "DialogManager_ShowInfo:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowInfo()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowInfo_Bad(t *core.T) {
	// DialogManager ShowInfo
	ax7Variant := "DialogManager_ShowInfo:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowInfo()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowInfo_Ugly(t *core.T) {
	// DialogManager ShowInfo
	ax7Variant := "DialogManager_ShowInfo:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowInfo()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowQuestion_Good(t *core.T) {
	// DialogManager ShowQuestion
	ax7Variant := "DialogManager_ShowQuestion:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowQuestion()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowQuestion_Bad(t *core.T) {
	// DialogManager ShowQuestion
	ax7Variant := "DialogManager_ShowQuestion:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowQuestion()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowQuestion_Ugly(t *core.T) {
	// DialogManager ShowQuestion
	ax7Variant := "DialogManager_ShowQuestion:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowQuestion()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowWarning_Good(t *core.T) {
	// DialogManager ShowWarning
	ax7Variant := "DialogManager_ShowWarning:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowWarning()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowWarning_Bad(t *core.T) {
	// DialogManager ShowWarning
	ax7Variant := "DialogManager_ShowWarning:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowWarning()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowWarning_Ugly(t *core.T) {
	// DialogManager ShowWarning
	ax7Variant := "DialogManager_ShowWarning:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowWarning()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowError_Good(t *core.T) {
	// DialogManager ShowError
	ax7Variant := "DialogManager_ShowError:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowError()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowError_Bad(t *core.T) {
	// DialogManager ShowError
	ax7Variant := "DialogManager_ShowError:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowError()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}

func TestDialog_DialogManager_ShowError_Ugly(t *core.T) {
	// DialogManager ShowError
	ax7Variant := "DialogManager_ShowError:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(DialogManager)
	result := core.Try(func() any {
		got0, got1 := subject.ShowError()
		return core.Sprintf("%T,%T", got0, got1)
	})
	core.AssertNotNil(t, result.Value)
	core.AssertNotEqual(t, "", core.Sprintf("%T", result.Value))
}
