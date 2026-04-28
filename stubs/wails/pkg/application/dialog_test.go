package application

import (
	core "dappco.re/go"
)

func TestDialog_OpenFileDialogStruct_Good(t *core.T) {
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

func TestDialog_OpenFileDialogStruct_Bad(t *core.T) {
	dialog := newOpenFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
	core.AssertTrue(t, dialog.canChooseFiles)
}

func TestDialog_OpenFileDialogStruct_Ugly(t *core.T) {
	dialog := newOpenFileDialog()
	selected := []string{"/tmp/input.txt"}
	dialog.SetSelectedFiles(selected)
	selected[0] = "/tmp/mutated.txt"

	paths, err := dialog.PromptForMultipleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, []string{"/tmp/input.txt"}, paths)
}

func TestDialog_SaveFileDialogStruct_Good(t *core.T) {
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

func TestDialog_SaveFileDialogStruct_Bad(t *core.T) {
	dialog := newSaveFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
}

func TestDialog_SaveFileDialogStruct_Ugly(t *core.T) {
	dialog := newSaveFileDialog()
	selected := "/tmp/a.csv"
	dialog.SetSelectedPath(selected)
	selected = "/tmp/b.csv"

	got, err := dialog.PromptForSingleSelection()

	core.RequireNoError(t, err)
	core.AssertEqual(t, "/tmp/a.csv", got)
}

func TestDialog_MessageDialog_Good(t *core.T) {
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

func TestDialog_MessageDialog_Bad(t *core.T) {
	dialog := newMessageDialog(InfoDialogType)

	got, err := dialog.Show()

	core.RequireNoError(t, err)
	core.AssertEmpty(t, got)
	core.AssertEqual(t, InfoDialogType, dialog.dialogType)
}

func TestDialog_MessageDialog_Ugly(t *core.T) {
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
	manager := &DialogManager{}

	core.AssertEqual(t, InfoDialogType, manager.Info().dialogType)
	core.AssertEqual(t, QuestionDialogType, manager.Question().dialogType)
	core.AssertEqual(t, WarningDialogType, manager.Warning().dialogType)
	core.AssertEqual(t, ErrorDialogType, manager.Error().dialogType)
	core.AssertNotNil(t, manager.OpenFile())
	core.AssertNotNil(t, manager.SaveFile())
}

func TestDialog_DialogManager_Bad(t *core.T) {
	manager := &DialogManager{}

	core.AssertNotNil(t, manager.OpenFileWithOptions(nil))
	core.AssertNotNil(t, manager.SaveFileWithOptions(nil))
}

func TestDialog_DialogManager_Ugly(t *core.T) {
	manager := &DialogManager{}

	open := manager.OpenFileWithOptions(&OpenFileDialogOptions{AllowsMultipleSelection: true})
	save := manager.SaveFileWithOptions(&SaveFileDialogOptions{Filename: "out.csv"})

	core.AssertTrue(t, open.multipleAllowed)
	core.AssertEqual(t, "out.csv", save.filename)
}
