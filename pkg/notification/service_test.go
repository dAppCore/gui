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
	sendErr              error
	permGranted          bool
	permErr              error
	revokeErr            error
	registerCategoryErr  error
	lastOpts             NotificationOptions
	lastCategory         NotificationCategory
	sendCalled           bool
	revokeCalled         bool
	registerCategoryCalled bool
}

func (m *mockPlatform) Send(opts NotificationOptions) error {
	m.sendCalled = true
	m.lastOpts = opts
	return m.sendErr
}
func (m *mockPlatform) RequestPermission() (bool, error) { return m.permGranted, m.permErr }
func (m *mockPlatform) CheckPermission() (bool, error)   { return m.permGranted, m.permErr }
func (m *mockPlatform) RevokePermission() error {
	m.revokeCalled = true
	return m.revokeErr
}
func (m *mockPlatform) RegisterCategory(category NotificationCategory) error {
	m.registerCategoryCalled = true
	m.lastCategory = category
	return m.registerCategoryErr
}

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
		Options: NotificationOptions{Title: "Test", Message: "Hello"},
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
		Options: NotificationOptions{Title: "Warn", Message: "Oops", Severity: SeverityWarning},
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

func TestTaskRevokePermission_Good(t *testing.T) {
	mock, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskRevokePermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.revokeCalled)
}

func TestTaskRevokePermission_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskRevokePermission{})
	assert.False(t, handled)
}

func TestTaskRegisterCategory_Good(t *testing.T) {
	mock, c := newTestService(t)
	category := NotificationCategory{
		ID: "message",
		Actions: []NotificationAction{
			{ID: "reply", Title: "Reply"},
			{ID: "dismiss", Title: "Dismiss"},
		},
	}

	_, handled, err := c.PERFORM(TaskRegisterCategory{Category: category})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.registerCategoryCalled)
	assert.Equal(t, "message", mock.lastCategory.ID)
	assert.Len(t, mock.lastCategory.Actions, 2)
	assert.Equal(t, "reply", mock.lastCategory.Actions[0].ID)
}

func TestTaskRegisterCategory_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskRegisterCategory{})
	assert.False(t, handled)
}

func TestActionNotificationActionTriggered_Ugly(t *testing.T) {
	// Verify the action structs are distinct types.
	var triggered ActionNotificationActionTriggered
	var dismissed ActionNotificationDismissed
	triggered.NotificationID = "n1"
	triggered.ActionID = "reply"
	dismissed.NotificationID = "n1"
	assert.Equal(t, "n1", triggered.NotificationID)
	assert.Equal(t, "reply", triggered.ActionID)
	assert.Equal(t, "n1", dismissed.NotificationID)
}

func TestActionNotificationDismissed_Good(t *testing.T) {
	_, c := newTestService(t)

	var received *ActionNotificationDismissed
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionNotificationDismissed); ok {
			received = &a
		}
		return nil
	})

	// Broadcast dismissed action directly (as the platform adapter would).
	_ = c.ACTION(ActionNotificationDismissed{NotificationID: "notif-42"})
	require.NotNil(t, received)
	assert.Equal(t, "notif-42", received.NotificationID)
}

func TestActionNotificationActionTriggered_Good(t *testing.T) {
	_, c := newTestService(t)

	var received *ActionNotificationActionTriggered
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionNotificationActionTriggered); ok {
			received = &a
		}
		return nil
	})

	_ = c.ACTION(ActionNotificationActionTriggered{NotificationID: "notif-7", ActionID: "archive"})
	require.NotNil(t, received)
	assert.Equal(t, "notif-7", received.NotificationID)
	assert.Equal(t, "archive", received.ActionID)
}
