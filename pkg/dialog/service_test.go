// pkg/dialog/service_test.go
package dialog

import (
	"context"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/webview"
	"dappco.re/go/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	openFilePaths []string
	saveFilePath  string
	openDirPath   string
	messageButton string
	openFileErr   error
	saveFileErr   error
	openDirErr    error
	messageErr    error
	lastOpenOpts  OpenFileOptions
	lastSaveOpts  SaveFileOptions
	lastDirOpts   OpenDirectoryOptions
	lastMsgOpts   MessageDialogOptions
}

func (m *mockPlatform) OpenFile(opts OpenFileOptions) ([]string, error) {
	m.lastOpenOpts = opts
	return m.openFilePaths, m.openFileErr
}
func (m *mockPlatform) SaveFile(opts SaveFileOptions) (string, error) {
	m.lastSaveOpts = opts
	return m.saveFilePath, m.saveFileErr
}
func (m *mockPlatform) OpenDirectory(opts OpenDirectoryOptions) (string, error) {
	m.lastDirOpts = opts
	return m.openDirPath, m.openDirErr
}
func (m *mockPlatform) MessageDialog(opts MessageDialogOptions) (string, error) {
	m.lastMsgOpts = opts
	return m.messageButton, m.messageErr
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		openFilePaths: []string{"/tmp/file.txt"},
		saveFilePath:  "/tmp/save.txt",
		openDirPath:   "/tmp/dir",
		messageButton: "OK",
	}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

// --- Good path tests ---

func TestService_Register_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")
	assert.NotNil(t, svc)
}

func TestService_TaskOpenFile_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/a.txt", "/b.txt"}

	r := taskRun(c, "dialog.openFile", TaskOpenFile{
		Options: OpenFileOptions{Title: "Pick", AllowMultiple: true},
	})
	require.True(t, r.OK)
	paths := r.Value.([]string)
	assert.Equal(t, []string{"/a.txt", "/b.txt"}, paths)
	assert.Equal(t, "Pick", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.AllowMultiple)
}

func TestService_TaskOpenFile_FileFilters_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/img.png"}

	filters := []FileFilter{{DisplayName: "Images", Pattern: "*.png;*.jpg"}}
	r := taskRun(c, "dialog.openFile", TaskOpenFile{
		Options: OpenFileOptions{
			Title:   "Select image",
			Filters: filters,
		},
	})
	require.True(t, r.OK)
	assert.Equal(t, []string{"/img.png"}, r.Value.([]string))
	require.Len(t, mock.lastOpenOpts.Filters, 1)
	assert.Equal(t, "Images", mock.lastOpenOpts.Filters[0].DisplayName)
	assert.Equal(t, "*.png;*.jpg", mock.lastOpenOpts.Filters[0].Pattern)
}

func TestService_TaskOpenFile_MultipleSelection_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/a.txt", "/b.txt", "/c.txt"}

	r := taskRun(c, "dialog.openFile", TaskOpenFile{
		Options: OpenFileOptions{AllowMultiple: true},
	})
	require.True(t, r.OK)
	assert.Equal(t, []string{"/a.txt", "/b.txt", "/c.txt"}, r.Value.([]string))
	assert.True(t, mock.lastOpenOpts.AllowMultiple)
}

func TestService_TaskOpenFile_CanChooseOptions_Good(t *testing.T) {
	mock, c := newTestService(t)

	r := taskRun(c, "dialog.openFile", TaskOpenFile{
		Options: OpenFileOptions{
			CanChooseFiles:       true,
			CanChooseDirectories: true,
			ShowHiddenFiles:      true,
		},
	})
	require.True(t, r.OK)
	assert.True(t, mock.lastOpenOpts.CanChooseFiles)
	assert.True(t, mock.lastOpenOpts.CanChooseDirectories)
	assert.True(t, mock.lastOpenOpts.ShowHiddenFiles)
}

func TestService_TaskOpenFileWithOptions_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/log.txt"}

	opts := &OpenFileOptions{
		Title:           "Select log",
		AllowMultiple:   false,
		ShowHiddenFiles: true,
	}
	r := taskRun(c, "dialog.openFile", TaskOpenFileWithOptions{Options: opts})
	require.True(t, r.OK)
	assert.Equal(t, []string{"/log.txt"}, r.Value.([]string))
	assert.Equal(t, "Select log", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.ShowHiddenFiles)
}

func TestService_TaskOpenFileWithOptions_NilOptions_Good(t *testing.T) {
	_, c := newTestService(t)

	r := taskRun(c, "dialog.openFile", TaskOpenFileWithOptions{Options: nil})
	require.True(t, r.OK)
	assert.NotNil(t, r.Value)
}

func TestService_TaskSaveFile_Good(t *testing.T) {
	_, c := newTestService(t)
	r := taskRun(c, "dialog.saveFile", TaskSaveFile{
		Options: SaveFileOptions{Filename: "out.txt"},
	})
	require.True(t, r.OK)
	assert.Equal(t, "/tmp/save.txt", r.Value)
}

func TestService_TaskSaveFile_ShowHidden_Good(t *testing.T) {
	mock, c := newTestService(t)

	r := taskRun(c, "dialog.saveFile", TaskSaveFile{
		Options: SaveFileOptions{Filename: "out.txt", ShowHiddenFiles: true},
	})
	require.True(t, r.OK)
	assert.True(t, mock.lastSaveOpts.ShowHiddenFiles)
}

func TestService_TaskSaveFileWithOptions_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.saveFilePath = "/exports/data.json"

	opts := &SaveFileOptions{
		Title:    "Export data",
		Filename: "data.json",
		Filters:  []FileFilter{{DisplayName: "JSON", Pattern: "*.json"}},
	}
	r := taskRun(c, "dialog.saveFile", TaskSaveFileWithOptions{Options: opts})
	require.True(t, r.OK)
	assert.Equal(t, "/exports/data.json", r.Value.(string))
	assert.Equal(t, "Export data", mock.lastSaveOpts.Title)
	require.Len(t, mock.lastSaveOpts.Filters, 1)
	assert.Equal(t, "JSON", mock.lastSaveOpts.Filters[0].DisplayName)
}

func TestService_TaskSaveFileWithOptions_NilOptions_Good(t *testing.T) {
	_, c := newTestService(t)

	r := taskRun(c, "dialog.saveFile", TaskSaveFileWithOptions{Options: nil})
	require.True(t, r.OK)
	assert.Equal(t, "/tmp/save.txt", r.Value)
}

func TestService_TaskOpenDirectory_Good(t *testing.T) {
	mock, c := newTestService(t)

	r := taskRun(c, "dialog.openDirectory", TaskOpenDirectory{
		Options: OpenDirectoryOptions{Title: "Pick Dir", ShowHiddenFiles: true},
	})
	require.True(t, r.OK)
	assert.Equal(t, "/tmp/dir", r.Value)
	assert.Equal(t, "Pick Dir", mock.lastDirOpts.Title)
	assert.True(t, mock.lastDirOpts.ShowHiddenFiles)
}

func TestService_TaskMessageDialog_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	r := taskRun(c, "dialog.message", TaskMessageDialog{
		Options: MessageDialogOptions{
			Type: DialogQuestion, Title: "Confirm",
			Message: "Sure?", Buttons: []string{"Yes", "No"},
		},
	})
	require.True(t, r.OK)
	assert.Equal(t, "Yes", r.Value)
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
}

func TestService_TaskInfo_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	r := taskRun(c, "dialog.info", TaskInfo{
		Title: "Done", Message: "File saved successfully.",
	})
	require.True(t, r.OK)
	assert.Equal(t, "OK", r.Value.(string))
	assert.Equal(t, DialogInfo, mock.lastMsgOpts.Type)
	assert.Equal(t, "Done", mock.lastMsgOpts.Title)
	assert.Equal(t, "File saved successfully.", mock.lastMsgOpts.Message)
}

func TestService_TaskInfo_WithButtons_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Close"

	r := taskRun(c, "dialog.info", TaskInfo{
		Title: "Notice", Message: "Update available.", Buttons: []string{"Close", "Later"},
	})
	require.True(t, r.OK)
	assert.Equal(t, "Close", r.Value.(string))
	assert.Equal(t, []string{"Close", "Later"}, mock.lastMsgOpts.Buttons)
}

func TestService_TaskQuestion_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	r := taskRun(c, "dialog.question", TaskQuestion{
		Title: "Confirm deletion", Message: "Delete file?", Buttons: []string{"Yes", "No"},
	})
	require.True(t, r.OK)
	assert.Equal(t, "Yes", r.Value.(string))
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
	assert.Equal(t, "Confirm deletion", mock.lastMsgOpts.Title)
}

func TestService_TaskWarning_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	r := taskRun(c, "dialog.warning", TaskWarning{
		Title: "Disk full", Message: "Storage is critically low.", Buttons: []string{"OK"},
	})
	require.True(t, r.OK)
	assert.Equal(t, "OK", r.Value.(string))
	assert.Equal(t, DialogWarning, mock.lastMsgOpts.Type)
	assert.Equal(t, "Disk full", mock.lastMsgOpts.Title)
}

func TestService_TaskError_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	r := taskRun(c, "dialog.error", TaskError{
		Title: "Operation failed", Message: "could not write file: permission denied",
	})
	require.True(t, r.OK)
	assert.Equal(t, "OK", r.Value.(string))
	assert.Equal(t, DialogError, mock.lastMsgOpts.Type)
	assert.Equal(t, "Operation failed", mock.lastMsgOpts.Title)
	assert.Equal(t, "could not write file: permission denied", mock.lastMsgOpts.Message)
}

func TestService_TaskPrompt_Good(t *testing.T) {
	_, c := newTestService(t)

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowList:
			return core.Result{Value: []window.WindowInfo{
				{Name: "editor", Focused: true},
			}, OK: true}
		default:
			return core.Result{}
		}
	})

	var script string
	c.Action("gui.webview.eval", func(_ context.Context, opts core.Options) core.Result {
		task := opts.Get("task").Value.(webview.TaskEvaluate)
		script = task.Script
		assert.Equal(t, "editor", task.Window)
		return core.Result{Value: "draft", OK: true}
	})

	r := taskRun(c, "dialog.prompt", TaskPrompt{
		Title:        "Rename",
		Message:      "Enter a new name",
		DefaultValue: "draft",
	})
	require.True(t, r.OK)
	result := r.Value.(PromptResult)
	assert.Equal(t, "draft", result.Value)
	assert.True(t, result.Confirmed)
	assert.Contains(t, script, "window.prompt(")
	assert.Contains(t, script, "Rename")
	assert.Contains(t, script, "Enter a new name")
}

// --- Bad path tests ---

func TestService_TaskOpenFile_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.openFile").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskOpenFileWithOptions_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.openFile").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskSaveFileWithOptions_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.saveFile").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskInfo_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.info").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskQuestion_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.question").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskWarning_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.warning").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_TaskError_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.error").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

// --- Ugly path tests ---

func TestService_TaskOpenFile_Ugly(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = nil

	r := taskRun(c, "dialog.openFile", TaskOpenFile{
		Options: OpenFileOptions{Title: "Pick"},
	})
	require.True(t, r.OK)
	assert.Nil(t, r.Value.([]string))
}

func TestService_TaskOpenFileWithOptions_MultipleFilters_Ugly(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/doc.pdf"}

	opts := &OpenFileOptions{
		Title: "Select document",
		Filters: []FileFilter{
			{DisplayName: "PDF", Pattern: "*.pdf"},
			{DisplayName: "Word", Pattern: "*.docx"},
			{DisplayName: "All files", Pattern: "*.*"},
		},
	}
	r := taskRun(c, "dialog.openFile", TaskOpenFileWithOptions{Options: opts})
	require.True(t, r.OK)
	assert.Equal(t, []string{"/doc.pdf"}, r.Value.([]string))
	assert.Len(t, mock.lastOpenOpts.Filters, 3)
}

func TestService_TaskSaveFileWithOptions_FiltersAndHidden_Ugly(t *testing.T) {
	mock, c := newTestService(t)

	opts := &SaveFileOptions{
		Title:           "Save",
		Filename:        "output.csv",
		ShowHiddenFiles: true,
		Filters:         []FileFilter{{DisplayName: "CSV", Pattern: "*.csv"}},
	}
	r := taskRun(c, "dialog.saveFile", TaskSaveFileWithOptions{Options: opts})
	require.True(t, r.OK)
	assert.True(t, mock.lastSaveOpts.ShowHiddenFiles)
	assert.Equal(t, "output.csv", mock.lastSaveOpts.Filename)
}

func TestService_UnknownTask_Ugly(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("dialog.nonexistent").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestService_promptOptionsFrom_Good(t *testing.T) {
	got, err := promptOptionsFrom(core.NewOptions(
		core.Option{Key: "task", Value: TaskPrompt{
			Title:        "Rename",
			Message:      "Enter a new name",
			DefaultValue: "draft",
		}},
	))
	require.NoError(t, err)
	assert.Equal(t, TaskPrompt{
		Title:        "Rename",
		Message:      "Enter a new name",
		DefaultValue: "draft",
	}, got)
}

func TestService_promptOptionsFrom_Bad(t *testing.T) {
	got, err := promptOptionsFrom(core.NewOptions(
		core.Option{Key: "title", Value: "Rename"},
		core.Option{Key: "message", Value: "Enter a new name"},
		core.Option{Key: "defaultValue", Value: "draft"},
	))
	require.NoError(t, err)
	assert.Equal(t, TaskPrompt{
		Title:        "Rename",
		Message:      "Enter a new name",
		DefaultValue: "draft",
	}, got)
}

func TestService_promptOptionsFrom_Ugly(t *testing.T) {
	got, err := promptOptionsFrom(core.NewOptions())
	require.NoError(t, err)
	assert.Zero(t, got)
}

func TestService_promptWindowName_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowList:
			return core.Result{Value: []window.WindowInfo{
				{Name: "editor"},
				{Name: "preview", Focused: true},
			}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.promptWindowName()
	require.NoError(t, err)
	assert.Equal(t, "preview", got)
}

func TestService_promptWindowName_Bad(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowList:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.promptWindowName()
	require.Error(t, err)
	assert.Empty(t, got)
}

func TestService_promptWindowName_Ugly(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")

	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowList:
			return core.Result{Value: []window.WindowInfo{
				{Name: "editor"},
				{Name: "terminal"},
			}, OK: true}
		default:
			return core.Result{}
		}
	})

	got, err := svc.promptWindowName()
	require.NoError(t, err)
	assert.Equal(t, "editor", got)
}

func TestService_promptScript_Good(t *testing.T) {
	script := promptScript("Rename", "Enter a new name", "draft")
	assert.Contains(t, script, "window.prompt(")
	assert.Contains(t, script, "Rename")
	assert.Contains(t, script, "Enter a new name")
	assert.Contains(t, script, "draft")
}

func TestService_promptScript_Bad(t *testing.T) {
	script := promptScript("Rename", "", "")
	assert.Contains(t, script, "window.prompt(")
	assert.Contains(t, script, "Rename")
	assert.NotContains(t, script, "Enter a new name")
}

func TestService_promptScript_Ugly(t *testing.T) {
	script := promptScript("", "Line 1\nLine 2", "\"quoted\"")
	assert.Contains(t, script, "Line 1")
	assert.Contains(t, script, "Line 2")
	assert.Contains(t, script, "quoted")
	assert.True(t, strings.Contains(script, "window.prompt("))
}
