// pkg/notification/service_test.go
package notification

import (
	"context"
	core "dappco.re/go/core"
	"testing"

	"forge.lthn.ai/core/gui/pkg/dialog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	sendErr      error
	permGranted  bool
	permErr      error
	revokeErr    error
	revokeCalled bool
	clearErr     error
	clearCalled  bool
	clearID      string
	lastOpts     NotificationOptions
	sendCalled   bool
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
func (m *mockPlatform) Clear(id string) error {
	m.clearCalled = true
	m.clearID = id
	return m.clearErr
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

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "notification")
	assert.NotNil(t, svc)
}

func TestTaskSend_Good(t *testing.T) {
	mock, c := newTestService(t)
	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{Title: "Test", Message: "Hello"},
	})
	require.True(t, r.OK)
	assert.True(t, mock.sendCalled)
	assert.Equal(t, "Test", mock.lastOpts.Title)
}

func TestTaskSend_Fallback_Good(t *testing.T) {
	// Platform fails -> falls back to dialog via IPC
	mockNotify := &mockPlatform{sendErr: core.NewError("no permission")}
	mockDlg := &mockDialogPlatform{}
	c := core.New(
		core.WithService(dialog.Register(mockDlg)),
		core.WithService(Register(mockNotify)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{Title: "Warn", Message: "Oops", Severity: SeverityWarning},
	})
	assert.True(t, r.OK) // fallback succeeds even though platform failed
	assert.True(t, mockDlg.messageCalled)
	assert.Equal(t, dialog.DialogWarning, mockDlg.lastMsgOpts.Type)
}

func TestQueryPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryPermission{})
	require.True(t, r.OK)
	status := r.Value.(PermissionStatus)
	assert.True(t, status.Granted)
}

func TestTaskRequestPermission_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("notification.requestPermission").Run(context.Background(), core.NewOptions())
	require.True(t, r.OK)
	assert.Equal(t, true, r.Value)
}

func TestTaskSend_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("notification.send").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

// --- TaskRevokePermission ---

func TestTaskRevokePermission_Good(t *testing.T) {
	mock, c := newTestService(t)
	r := c.Action("notification.revokePermission").Run(context.Background(), core.NewOptions())
	require.True(t, r.OK)
	assert.True(t, mock.revokeCalled)
}

func TestTaskRevokePermission_Bad(t *testing.T) {
	mock, c := newTestService(t)
	mock.revokeErr = core.NewError("cannot revoke")
	r := c.Action("notification.revokePermission").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestTaskRevokePermission_Ugly(t *testing.T) {
	// No service registered — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("notification.revokePermission").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
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
	r := taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: category})
	require.True(t, r.OK)

	svc := core.MustServiceFor[*Service](c, "notification")
	stored, ok := svc.categories["message"]
	require.True(t, ok)
	assert.Equal(t, 2, len(stored.Actions))
	assert.Equal(t, "reply", stored.Actions[0].ID)
	assert.True(t, stored.Actions[1].Destructive)
}

func TestTaskRegisterCategory_Bad(t *testing.T) {
	// No service registered — action is not registered
	c := core.New(core.WithServiceLock())
	r := taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: NotificationCategory{ID: "x"}})
	assert.False(t, r.OK)
}

func TestTaskRegisterCategory_Ugly(t *testing.T) {
	// Re-registering a category replaces the previous one
	_, c := newTestService(t)
	first := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "a", Title: "A"}}}
	second := NotificationCategory{ID: "chat", Actions: []NotificationAction{{ID: "b", Title: "B"}, {ID: "c", Title: "C"}}}

	require.True(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: first}).OK)
	require.True(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{Category: second}).OK)

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
	r := taskRun(c, "notification.send", TaskSend{Options: options})
	require.True(t, r.OK)
	assert.Equal(t, "message", mock.lastOpts.CategoryID)
	assert.Equal(t, 2, len(mock.lastOpts.Actions))
}

func TestTaskSend_RegisteredCategoryActions_Good(t *testing.T) {
	mock, c := newTestService(t)
	require.True(t, taskRun(c, "notification.registerCategory", TaskRegisterCategory{
		Category: NotificationCategory{
			ID: "message",
			Actions: []NotificationAction{
				{ID: "reply", Title: "Reply"},
				{ID: "dismiss", Title: "Dismiss"},
			},
		},
	}).OK)

	r := taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{
			Title:      "Team Chat",
			Message:    "New message from Alice",
			CategoryID: "message",
		},
	})
	require.True(t, r.OK)
	assert.Equal(t, 2, len(mock.lastOpts.Actions))
	assert.Equal(t, "reply", mock.lastOpts.Actions[0].ID)
}

func TestTaskClear_Good_Specific(t *testing.T) {
	mock, c := newTestService(t)
	require.True(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n1", Title: "One", Message: "Hello"},
	}).OK)

	r := taskRun(c, "notification.clear", TaskClear{ID: "n1"})
	require.True(t, r.OK)
	assert.True(t, mock.clearCalled)
	assert.Equal(t, "n1", mock.clearID)
}

func TestTaskClear_Good_All(t *testing.T) {
	mock, c := newTestService(t)
	require.True(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n1", Title: "One", Message: "Hello"},
	}).OK)
	require.True(t, taskRun(c, "notification.send", TaskSend{
		Options: NotificationOptions{ID: "n2", Title: "Two", Message: "World"},
	}).OK)

	r := taskRun(c, "notification.clear", TaskClear{})
	require.True(t, r.OK)
	assert.True(t, mock.clearCalled)
	assert.Equal(t, "", mock.clearID)

	svc := core.MustServiceFor[*Service](c, "notification")
	assert.Empty(t, svc.active)
}

// --- ActionNotificationActionTriggered ---

func TestActionNotificationActionTriggered_Good(t *testing.T) {
	// ActionNotificationActionTriggered is broadcast by external code; confirm it can be received
	_, c := newTestService(t)
	var received *ActionNotificationActionTriggered
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionNotificationActionTriggered); ok {
			received = &a
		}
		return core.Result{OK: true}
	})
	_ = c.ACTION(ActionNotificationActionTriggered{NotificationID: "n1", ActionID: "reply"})
	require.NotNil(t, received)
	assert.Equal(t, "n1", received.NotificationID)
	assert.Equal(t, "reply", received.ActionID)
}

func TestActionNotificationDismissed_Good(t *testing.T) {
	_, c := newTestService(t)
	var received *ActionNotificationDismissed
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionNotificationDismissed); ok {
			received = &a
		}
		return core.Result{OK: true}
	})
	_ = c.ACTION(ActionNotificationDismissed{ID: "n2"})
	require.NotNil(t, received)
	assert.Equal(t, "n2", received.ID)
}

func TestQueryPermission_Bad(t *testing.T) {
	// No service — QUERY returns handled=false
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryPermission{})
	assert.False(t, r.OK)
}

func TestQueryPermission_Ugly(t *testing.T) {
	// Platform returns error — QUERY returns OK=false (framework does not propagate Value for failed queries)
	mock := &mockPlatform{permErr: core.NewError("platform error")}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	r := c.QUERY(QueryPermission{})
	assert.False(t, r.OK)
}
