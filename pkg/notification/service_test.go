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
	sendErr       error
	permGranted   bool
	permErr       error
	revokeErr     error
	revokeCalled  bool
	lastOpts      NotificationOptions
	sendCalled    bool
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

// --- TaskRevokePermission ---

func TestTaskRevokePermission_Good(t *testing.T) {
	mock, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskRevokePermission{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.revokeCalled)
}

func TestTaskRevokePermission_Bad(t *testing.T) {
	mock, c := newTestService(t)
	mock.revokeErr = errors.New("cannot revoke")
	_, handled, err := c.PERFORM(TaskRevokePermission{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskRevokePermission_Ugly(t *testing.T) {
	// No service registered — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskRevokePermission{})
	assert.False(t, handled)
}

// --- TaskRegisterCategory ---

func TestTaskRegisterCategory_Good(t *testing.T) {
	_, c := newTestService(t)
	category := NotificationCategory{
		ID: "message",
		Actions: []NotificationAction{
			{ID: "reply", Title: "Reply"},
			{ID: "delete", Title: "Delete", Destructive: true},
		},
	}
	_, handled, err := c.PERFORM(TaskRegisterCategory{Category: category})
	require.NoError(t, err)
	assert.True(t, handled)

	svc := core.MustServiceFor[*Service](c, "notification")
	stored, ok := svc.categories["message"]
	require.True(t, ok)
	assert.Equal(t, 2, len(stored.Actions))
	assert.Equal(t, "reply", stored.Actions[0].ID)
	assert.True(t, stored.Actions[1].Destructive)
}

func TestTaskRegisterCategory_Bad(t *testing.T) {
	// No service registered — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskRegisterCategory{Category: NotificationCategory{ID: "x"}})
	assert.False(t, handled)
}

func TestTaskRegisterCategory_Ugly(t *testing.T) {
	// Re-registering a category replaces the previous one
	_, c := newTestService(t)
	first := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "a", Title: "A"}}}
	second := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "b", Title: "B"}, {ID: "c", Title: "C"}}}

	_, _, err := c.PERFORM(TaskRegisterCategory{Category: first})
	require.NoError(t, err)
	_, _, err = c.PERFORM(TaskRegisterCategory{Category: second})
	require.NoError(t, err)

	svc := core.MustServiceFor[*Service](c, "notification")
	assert.Equal(t, 2, len(svc.categories["chat"].Actions))
	assert.Equal(t, "b", svc.categories["chat"].Actions[0].ID)
}

// --- NotificationOptions with Actions ---

func TestTaskSend_WithActions_Good(t *testing.T) {
	mock, c := newTestService(t)
	options := NotificationOptions{
		Title:      "Team Chat",
		Message:    "New message from Alice",
		CategoryID: "message",
		Actions: []NotificationAction{
			{ID: "reply", Title: "Reply"},
			{ID: "dismiss", Title: "Dismiss"},
		},
	}
	_, handled, err := c.PERFORM(TaskSend{Options: options})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "message", mock.lastOpts.CategoryID)
	assert.Equal(t, 2, len(mock.lastOpts.Actions))
}

// --- ActionNotificationActionTriggered ---

func TestActionNotificationActionTriggered_Good(t *testing.T) {
	// ActionNotificationActionTriggered is broadcast by external code; confirm it can be received
	_, c := newTestService(t)
	var received *ActionNotificationActionTriggered
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionNotificationActionTriggered); ok {
			received = &a
		}
		return nil
	})
	_ = c.ACTION(ActionNotificationActionTriggered{NotificationID: "n1", ActionID: "reply"})
	require.NotNil(t, received)
	assert.Equal(t, "n1", received.NotificationID)
	assert.Equal(t, "reply", received.ActionID)
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
	_ = c.ACTION(ActionNotificationDismissed{ID: "n2"})
	require.NotNil(t, received)
	assert.Equal(t, "n2", received.ID)
}

func TestQueryPermission_Bad(t *testing.T) {
	// No service — QUERY returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(QueryPermission{})
	assert.False(t, handled)
}

func TestQueryPermission_Ugly(t *testing.T) {
	// Platform returns error — QUERY returns error with handled=true
	mock := &mockPlatform{permErr: errors.New("platform error")}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	_, handled, queryErr := c.QUERY(QueryPermission{})
	assert.True(t, handled)
	assert.Error(t, queryErr)
}
