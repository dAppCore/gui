// pkg/dialog/service_test.go
package dialog

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
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
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	return mock, c
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

	result, handled, err := c.PERFORM(TaskOpenFile{
		Options: OpenFileOptions{Title: "Pick", AllowMultiple: true},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	paths := result.([]string)
	assert.Equal(t, []string{"/a.txt", "/b.txt"}, paths)
	assert.Equal(t, "Pick", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.AllowMultiple)
}

func TestService_TaskOpenFile_FileFilters_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/img.png"}

	filters := []FileFilter{{DisplayName: "Images", Pattern: "*.png;*.jpg"}}
	result, handled, err := c.PERFORM(TaskOpenFile{
		Options: OpenFileOptions{
			Title:   "Select image",
			Filters: filters,
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, []string{"/img.png"}, result.([]string))
	require.Len(t, mock.lastOpenOpts.Filters, 1)
	assert.Equal(t, "Images", mock.lastOpenOpts.Filters[0].DisplayName)
	assert.Equal(t, "*.png;*.jpg", mock.lastOpenOpts.Filters[0].Pattern)
}

func TestService_TaskOpenFile_MultipleSelection_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/a.txt", "/b.txt", "/c.txt"}

	result, handled, err := c.PERFORM(TaskOpenFile{
		Options: OpenFileOptions{AllowMultiple: true},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, []string{"/a.txt", "/b.txt", "/c.txt"}, result.([]string))
	assert.True(t, mock.lastOpenOpts.AllowMultiple)
}

func TestService_TaskOpenFile_CanChooseOptions_Good(t *testing.T) {
	mock, c := newTestService(t)

	_, _, err := c.PERFORM(TaskOpenFile{
		Options: OpenFileOptions{
			CanChooseFiles:       true,
			CanChooseDirectories: true,
			ShowHiddenFiles:      true,
		},
	})
	require.NoError(t, err)
	assert.True(t, mock.lastOpenOpts.CanChooseFiles)
	assert.True(t, mock.lastOpenOpts.CanChooseDirectories)
	assert.True(t, mock.lastOpenOpts.ShowHiddenFiles)
}

func TestService_TaskOpenFileWithOptions_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/log.txt"}

	opts := &OpenFileOptions{
		Title:         "Select log",
		AllowMultiple: false,
		ShowHiddenFiles: true,
	}
	result, handled, err := c.PERFORM(TaskOpenFileWithOptions{Options: opts})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, []string{"/log.txt"}, result.([]string))
	assert.Equal(t, "Select log", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.ShowHiddenFiles)
}

func TestService_TaskOpenFileWithOptions_NilOptions_Good(t *testing.T) {
	_, c := newTestService(t)

	result, handled, err := c.PERFORM(TaskOpenFileWithOptions{Options: nil})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.NotNil(t, result)
}

func TestService_TaskSaveFile_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskSaveFile{
		Options: SaveFileOptions{Filename: "out.txt"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/save.txt", result)
}

func TestService_TaskSaveFile_ShowHidden_Good(t *testing.T) {
	mock, c := newTestService(t)

	_, _, err := c.PERFORM(TaskSaveFile{
		Options: SaveFileOptions{Filename: "out.txt", ShowHiddenFiles: true},
	})
	require.NoError(t, err)
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
	result, handled, err := c.PERFORM(TaskSaveFileWithOptions{Options: opts})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/exports/data.json", result.(string))
	assert.Equal(t, "Export data", mock.lastSaveOpts.Title)
	require.Len(t, mock.lastSaveOpts.Filters, 1)
	assert.Equal(t, "JSON", mock.lastSaveOpts.Filters[0].DisplayName)
}

func TestService_TaskSaveFileWithOptions_NilOptions_Good(t *testing.T) {
	_, c := newTestService(t)

	result, handled, err := c.PERFORM(TaskSaveFileWithOptions{Options: nil})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/save.txt", result)
}

func TestService_TaskOpenDirectory_Good(t *testing.T) {
	mock, c := newTestService(t)

	result, handled, err := c.PERFORM(TaskOpenDirectory{
		Options: OpenDirectoryOptions{Title: "Pick Dir", ShowHiddenFiles: true},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/dir", result)
	assert.Equal(t, "Pick Dir", mock.lastDirOpts.Title)
	assert.True(t, mock.lastDirOpts.ShowHiddenFiles)
}

func TestService_TaskMessageDialog_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	result, handled, err := c.PERFORM(TaskMessageDialog{
		Options: MessageDialogOptions{
			Type: DialogQuestion, Title: "Confirm",
			Message: "Sure?", Buttons: []string{"Yes", "No"},
		},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Yes", result)
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
}

func TestService_TaskInfo_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskInfo{
		Title: "Done", Message: "File saved successfully.",
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result.(string))
	assert.Equal(t, DialogInfo, mock.lastMsgOpts.Type)
	assert.Equal(t, "Done", mock.lastMsgOpts.Title)
	assert.Equal(t, "File saved successfully.", mock.lastMsgOpts.Message)
}

func TestService_TaskInfo_WithButtons_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Close"

	result, handled, err := c.PERFORM(TaskInfo{
		Title: "Notice", Message: "Update available.", Buttons: []string{"Close", "Later"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Close", result.(string))
	assert.Equal(t, []string{"Close", "Later"}, mock.lastMsgOpts.Buttons)
}

func TestService_TaskQuestion_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	result, handled, err := c.PERFORM(TaskQuestion{
		Title: "Confirm deletion", Message: "Delete file?", Buttons: []string{"Yes", "No"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Yes", result.(string))
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
	assert.Equal(t, "Confirm deletion", mock.lastMsgOpts.Title)
}

func TestService_TaskWarning_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskWarning{
		Title: "Disk full", Message: "Storage is critically low.", Buttons: []string{"OK"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result.(string))
	assert.Equal(t, DialogWarning, mock.lastMsgOpts.Type)
	assert.Equal(t, "Disk full", mock.lastMsgOpts.Title)
}

func TestService_TaskError_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskError{
		Title: "Operation failed", Message: "could not write file: permission denied",
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result.(string))
	assert.Equal(t, DialogError, mock.lastMsgOpts.Type)
	assert.Equal(t, "Operation failed", mock.lastMsgOpts.Title)
	assert.Equal(t, "could not write file: permission denied", mock.lastMsgOpts.Message)
}

// --- Bad path tests ---

func TestService_TaskOpenFile_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFile{})
	assert.False(t, handled)
}

func TestService_TaskOpenFileWithOptions_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFileWithOptions{})
	assert.False(t, handled)
}

func TestService_TaskSaveFileWithOptions_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSaveFileWithOptions{})
	assert.False(t, handled)
}

func TestService_TaskInfo_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskInfo{})
	assert.False(t, handled)
}

func TestService_TaskQuestion_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskQuestion{})
	assert.False(t, handled)
}

func TestService_TaskWarning_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskWarning{})
	assert.False(t, handled)
}

func TestService_TaskError_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskError{})
	assert.False(t, handled)
}

// --- Ugly path tests ---

func TestService_TaskOpenFile_Ugly(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = nil

	result, handled, err := c.PERFORM(TaskOpenFile{
		Options: OpenFileOptions{Title: "Pick"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result.([]string))
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
	result, handled, err := c.PERFORM(TaskOpenFileWithOptions{Options: opts})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, []string{"/doc.pdf"}, result.([]string))
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
	_, _, err := c.PERFORM(TaskSaveFileWithOptions{Options: opts})
	require.NoError(t, err)
	assert.True(t, mock.lastSaveOpts.ShowHiddenFiles)
	assert.Equal(t, "output.csv", mock.lastSaveOpts.Filename)
}

func TestService_UnknownTask_Ugly(t *testing.T) {
	type unknownTask struct{}
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(unknownTask{})
	assert.False(t, handled)
}
