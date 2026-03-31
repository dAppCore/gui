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

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "dialog")
	assert.NotNil(t, svc)
}

func TestTaskOpenFile_Good(t *testing.T) {
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

func TestTaskSaveFile_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskSaveFile{
		Options: SaveFileOptions{Filename: "out.txt"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/save.txt", result)
}

func TestTaskOpenDirectory_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskOpenDirectory{
		Options: OpenDirectoryOptions{Title: "Pick Dir"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/dir", result)
}

func TestTaskMessageDialog_Good(t *testing.T) {
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

func TestTaskOpenFile_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFile{})
	assert.False(t, handled)
}

func TestTaskOpenFileWithOptions_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/docs/report.pdf"}

	result, handled, err := c.PERFORM(TaskOpenFileWithOptions{
		Title:           "Select Document",
		AllowMultiple:   false,
		ShowHiddenFiles: true,
		CanChooseFiles:  true,
	})
	require.NoError(t, err)
	assert.True(t, handled)
	paths := result.([]string)
	assert.Equal(t, []string{"/docs/report.pdf"}, paths)
	assert.Equal(t, "Select Document", mock.lastOpenOpts.Title)
	assert.True(t, mock.lastOpenOpts.ShowHiddenFiles)
	assert.True(t, mock.lastOpenOpts.CanChooseFiles)
}

func TestTaskOpenFileWithOptions_CanChooseDirectories_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.openFilePaths = []string{"/home/user/projects"}

	_, handled, err := c.PERFORM(TaskOpenFileWithOptions{
		Title:                "Pick folder",
		CanChooseDirectories: true,
		CanChooseFiles:       false,
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.lastOpenOpts.CanChooseDirectories)
	assert.False(t, mock.lastOpenOpts.CanChooseFiles)
}

func TestTaskOpenFileWithOptions_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFileWithOptions{})
	assert.False(t, handled)
}

func TestTaskSaveFileWithOptions_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.saveFilePath = "/exports/data.csv"

	result, handled, err := c.PERFORM(TaskSaveFileWithOptions{
		Title:    "Export CSV",
		Filename: "data.csv",
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/exports/data.csv", result)
	assert.Equal(t, "Export CSV", mock.lastSaveOpts.Title)
	assert.Equal(t, "data.csv", mock.lastSaveOpts.Filename)
}

func TestTaskSaveFileWithOptions_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSaveFileWithOptions{})
	assert.False(t, handled)
}

func TestTaskInfo_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskInfo{Title: "Done", Message: "Saved successfully."})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result)
	assert.Equal(t, DialogInfo, mock.lastMsgOpts.Type)
	assert.Equal(t, "Done", mock.lastMsgOpts.Title)
}

func TestTaskQuestion_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	result, handled, err := c.PERFORM(TaskQuestion{
		Title:   "Confirm",
		Message: "Delete this file?",
		Buttons: []string{"Yes", "No"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "Yes", result)
	assert.Equal(t, DialogQuestion, mock.lastMsgOpts.Type)
}

func TestTaskWarning_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskWarning{Title: "Warning", Message: "File exists."})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result)
	assert.Equal(t, DialogWarning, mock.lastMsgOpts.Type)
}

func TestTaskError_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "OK"

	result, handled, err := c.PERFORM(TaskError{Title: "Error", Message: "Write failed."})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "OK", result)
	assert.Equal(t, DialogError, mock.lastMsgOpts.Type)
}

func TestTaskInfo_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskInfo{})
	assert.False(t, handled)
}

func TestTaskQuestion_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskQuestion{})
	assert.False(t, handled)
}

func TestTaskWarning_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskWarning{})
	assert.False(t, handled)
}

func TestTaskError_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskError{})
	assert.False(t, handled)
}
