package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialog_OpenFileDialogStruct_Good(t *testing.T) {
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

	require.NoError(t, err)
	assert.Equal(t, "/tmp/a.txt", got)
	paths, err := dialog.PromptForMultipleSelection()
	require.NoError(t, err)
	assert.Equal(t, []string{"/tmp/a.txt", "/tmp/b.txt"}, paths)
}

func TestDialog_OpenFileDialogStruct_Bad(t *testing.T) {
	dialog := newOpenFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	require.NoError(t, err)
	assert.Empty(t, got)
	assert.True(t, dialog.canChooseFiles)
}

func TestDialog_OpenFileDialogStruct_Ugly(t *testing.T) {
	dialog := newOpenFileDialog()
	selected := []string{"/tmp/input.txt"}
	dialog.SetSelectedFiles(selected)
	selected[0] = "/tmp/mutated.txt"

	paths, err := dialog.PromptForMultipleSelection()

	require.NoError(t, err)
	assert.Equal(t, []string{"/tmp/input.txt"}, paths)
}

func TestDialog_SaveFileDialogStruct_Good(t *testing.T) {
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

	require.NoError(t, err)
	assert.Equal(t, "/tmp/report.csv", got)
}

func TestDialog_SaveFileDialogStruct_Bad(t *testing.T) {
	dialog := newSaveFileDialog()
	dialog.SetOptions(nil)

	got, err := dialog.PromptForSingleSelection()

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDialog_SaveFileDialogStruct_Ugly(t *testing.T) {
	dialog := newSaveFileDialog()
	selected := "/tmp/a.csv"
	dialog.SetSelectedPath(selected)
	selected = "/tmp/b.csv"

	got, err := dialog.PromptForSingleSelection()

	require.NoError(t, err)
	assert.Equal(t, "/tmp/a.csv", got)
}

func TestDialog_MessageDialog_Good(t *testing.T) {
	dialog := newMessageDialog(QuestionDialogType)
	dialog.SetTitle("Confirm").
		SetMessage("Proceed?").
		AddButton("Yes").
		AddButton("No").
		SetDefaultButton("Yes").
		SetCancelButton("No")
	dialog.SetButtonClickedForStub("Yes")

	got, err := dialog.Show()

	require.NoError(t, err)
	assert.Equal(t, "Yes", got)
	assert.Equal(t, QuestionDialogType, dialog.dialogType)
	assert.Equal(t, "Confirm", dialog.title)
	assert.Equal(t, "Proceed?", dialog.message)
	require.Len(t, dialog.buttons, 2)
	assert.True(t, dialog.buttons[0].IsDefault)
	assert.True(t, dialog.buttons[1].IsCancel)
}

func TestDialog_MessageDialog_Bad(t *testing.T) {
	dialog := newMessageDialog(InfoDialogType)

	got, err := dialog.Show()

	require.NoError(t, err)
	assert.Empty(t, got)
	assert.Equal(t, InfoDialogType, dialog.dialogType)
}

func TestDialog_MessageDialog_Ugly(t *testing.T) {
	dialog := newMessageDialog(ErrorDialogType)
	dialog.AddButton("Retry").AddButton("Retry")
	dialog.SetDefaultButton("Retry")
	dialog.SetCancelButton("Retry")

	require.Len(t, dialog.buttons, 2)
	assert.True(t, dialog.buttons[0].IsDefault)
	assert.True(t, dialog.buttons[0].IsCancel)
	assert.True(t, dialog.buttons[1].IsDefault)
	assert.True(t, dialog.buttons[1].IsCancel)
}

func TestDialog_DialogManager_Good(t *testing.T) {
	manager := &DialogManager{}

	assert.Equal(t, InfoDialogType, manager.Info().dialogType)
	assert.Equal(t, QuestionDialogType, manager.Question().dialogType)
	assert.Equal(t, WarningDialogType, manager.Warning().dialogType)
	assert.Equal(t, ErrorDialogType, manager.Error().dialogType)
	assert.NotNil(t, manager.OpenFile())
	assert.NotNil(t, manager.SaveFile())
}

func TestDialog_DialogManager_Bad(t *testing.T) {
	manager := &DialogManager{}

	assert.NotNil(t, manager.OpenFileWithOptions(nil))
	assert.NotNil(t, manager.SaveFileWithOptions(nil))
}

func TestDialog_DialogManager_Ugly(t *testing.T) {
	manager := &DialogManager{}

	open := manager.OpenFileWithOptions(&OpenFileDialogOptions{AllowsMultipleSelection: true})
	save := manager.SaveFileWithOptions(&SaveFileDialogOptions{Filename: "out.csv"})

	assert.True(t, open.multipleAllowed)
	assert.Equal(t, "out.csv", save.filename)
}
