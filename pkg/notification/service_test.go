// pkg/notification/service_test.go
package notification

import (
	"context"
	"errors"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"forge.lthn.ai/core/gui/pkg/dialog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	sendErr     error
	permGranted bool
	permErr     error
	lastOpts    NotificationOptions
	sendCalled  bool
}

func (m *mockPlatform) Send(opts NotificationOptions) error {
	m.sendCalled = true
	m.lastOpts = opts
	return m.sendErr
}
func (m *mockPlatform) RequestPermission() (bool, error) { return m.permGranted, m.permErr }
func (m *mockPlatform) CheckPermission() (bool, error)   { return m.permGranted, m.permErr }

// mockDialogPlatform tracks whether MessageDialog was called (for fallback test).
type mockDialogPlatform struct {
	messageCalled bool
	lastMsgOpts   dialog.MessageDialogOptions
}

func (m *mockDialogPlatform) OpenFile(opts dialog.OpenFileOptions) ([]string, error) { return nil, nil }
func (m *mockDialogPlatform) SaveFile(opts dialog.SaveFileOptions) (string, error)   { return "", nil }
func (m *mockDialogPlatform) OpenDirectory(opts dialog.OpenDirectoryOptions) (string, error) {
	return "", nil
}
func (m *mockDialogPlatform) MessageDialog(opts dialog.MessageDialogOptions) (string, error) {
	m.messageCalled = true
	m.lastMsgOpts = opts
	return "OK", nil
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{permGranted: true}
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
	svc := core.MustServiceFor[*Service](c, "notification")
	assert.NotNil(t, svc)
}

func TestTaskSend_Good(t *testing.T) {
	mock, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskSend{
		Opts: NotificationOptions{Title: "Test", Message: "Hello"},
	})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.sendCalled)
	assert.Equal(t, "Test", mock.lastOpts.Title)
}

func TestTaskSend_Fallback_Good(t *testing.T) {
	// Platform fails -> falls back to dialog via IPC
	mockNotify := &mockPlatform{sendErr: errors.New("no permission")}
	mockDlg := &mockDialogPlatform{}
	c, err := core.New(
		core.WithService(dialog.Register(mockDlg)),
		core.WithService(Register(mockNotify)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, err := c.PERFORM(TaskSend{
		Opts: NotificationOptions{Title: "Warn", Message: "Oops", Severity: SeverityWarning},
	})
	assert.True(t, handled)
	assert.NoError(t, err) // fallback succeeds even though platform failed
	assert.True(t, mockDlg.messageCalled)
	assert.Equal(t, dialog.DialogWarning, mockDlg.lastMsgOpts.Type)
}

func TestQueryPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryPermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	status := result.(PermissionStatus)
	assert.True(t, status.Granted)
}

func TestTaskRequestPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskRequestPermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)
}

func TestTaskSend_Bad(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSend{})
	assert.False(t, handled)
}
