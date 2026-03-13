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
	openFilePaths  []string
	saveFilePath   string
	openDirPath    string
	messageButton  string
	openFileErr    error
	saveFileErr    error
	openDirErr     error
	messageErr     error
	lastOpenOpts   OpenFileOptions
	lastSaveOpts   SaveFileOptions
	lastDirOpts    OpenDirectoryOptions
	lastMsgOpts    MessageDialogOptions
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
		Opts: OpenFileOptions{Title: "Pick", AllowMultiple: true},
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
		Opts: SaveFileOptions{Filename: "out.txt"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/save.txt", result)
}

func TestTaskOpenDirectory_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskOpenDirectory{
		Opts: OpenDirectoryOptions{Title: "Pick Dir"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/dir", result)
}

func TestTaskMessageDialog_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.messageButton = "Yes"

	result, handled, err := c.PERFORM(TaskMessageDialog{
		Opts: MessageDialogOptions{
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
