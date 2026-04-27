// pkg/dock/service_test.go
package dock

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	visible          bool
	badge            string
	hasBadge         bool
	progress         float64
	bounceID         int
	bounceType       BounceType
	bounceCalled     bool
	stopBounceCalled bool
	showErr          error
	hideErr          error
	badgeErr         error
	removeErr        error
	progressErr      error
	bounceErr        error
	stopBounceErr    error
}

func (m *mockPlatform) ShowIcon() error {
	if m.showErr != nil {
		return m.showErr
	}
	m.visible = true
	return nil
}

func (m *mockPlatform) HideIcon() error {
	if m.hideErr != nil {
		return m.hideErr
	}
	m.visible = false
	return nil
}

func (m *mockPlatform) SetBadge(label string) error {
	if m.badgeErr != nil {
		return m.badgeErr
	}
	m.badge = label
	m.hasBadge = true
	return nil
}

func (m *mockPlatform) RemoveBadge() error {
	if m.removeErr != nil {
		return m.removeErr
	}
	m.badge = ""
	m.hasBadge = false
	return nil
}

func (m *mockPlatform) IsVisible() bool { return m.visible }

func (m *mockPlatform) SetProgressBar(progress float64) error {
	if m.progressErr != nil {
		return m.progressErr
	}
	m.progress = progress
	return nil
}

func (m *mockPlatform) Bounce(bounceType BounceType) (int, error) {
	if m.bounceErr != nil {
		return 0, m.bounceErr
	}
	m.bounceCalled = true
	m.bounceType = bounceType
	m.bounceID++
	return m.bounceID, nil
}

func (m *mockPlatform) StopBounce(requestID int) error {
	if m.stopBounceErr != nil {
		return m.stopBounceErr
	}
	m.stopBounceCalled = true
	return nil
}

// --- Test helpers ---

func newTestDockService(t *testing.T) (*Service, *core.Core, *mockPlatform) {
	t.Helper()
	mock := &mockPlatform{visible: true}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "dock")
	return svc, c, mock
}

func taskRun(c *core.Core, name string, task any) core.Result {
	return c.Action(name).Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: task},
	))
}

func setBadge(c *core.Core, label string) core.Result {
	return c.Action("dock.setBadge").Run(context.Background(), core.NewOptions(
		core.Option{Key: "label", Value: label},
	))
}

// --- Tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestDockService(t)
	assert.NotNil(t, svc)
}

func TestQueryVisible_Good(t *testing.T) {
	_, c, _ := newTestDockService(t)
	r := c.QUERY(QueryVisible{})
	require.True(t, r.OK)
	assert.Equal(t, true, r.Value)
}

func TestQueryVisible_Bad(t *testing.T) {
	// No dock service registered — QUERY returns handled=false
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryVisible{})
	assert.False(t, r.OK)
}

func TestTaskShowIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = false // Start hidden

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.showIcon", TaskShowIcon{})
	require.True(t, r.OK)
	assert.True(t, mock.visible)
	require.NotNil(t, received)
	assert.True(t, received.Visible)
}

func TestTaskHideIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = true // Start visible

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.hideIcon", TaskHideIcon{})
	require.True(t, r.OK)
	assert.False(t, mock.visible)
	require.NotNil(t, received)
	assert.False(t, received.Visible)
}

func TestTaskSetBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	r := setBadge(c, "3")
	require.True(t, r.OK)
	assert.Equal(t, "3", mock.badge)
	assert.True(t, mock.hasBadge)
}

func TestTaskSetBadge_EmptyLabel_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	r := setBadge(c, "")
	require.True(t, r.OK)
	assert.Equal(t, "", mock.badge)
	assert.True(t, mock.hasBadge) // Empty string = default system badge indicator
}

func TestTaskRemoveBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	// Set a badge first
	_ = setBadge(c, "5")

	r := taskRun(c, "dock.removeBadge", TaskRemoveBadge{})
	require.True(t, r.OK)
	assert.Equal(t, "", mock.badge)
	assert.False(t, mock.hasBadge)
}

func TestTaskShowIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.showErr = assert.AnError

	r := taskRun(c, "dock.showIcon", TaskShowIcon{})
	assert.False(t, r.OK)
}

func TestTaskHideIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.hideErr = assert.AnError

	r := taskRun(c, "dock.hideIcon", TaskHideIcon{})
	assert.False(t, r.OK)
}

func TestTaskSetBadge_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.badgeErr = assert.AnError

	r := setBadge(c, "3")
	assert.False(t, r.OK)
}

// --- TaskSetProgressBar ---

func TestTaskSetProgressBar_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionProgressChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionProgressChanged); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.setProgressBar", TaskSetProgressBar{Progress: 0.5})
	require.True(t, r.OK)
	assert.Equal(t, 0.5, mock.progress)
	require.NotNil(t, received)
	assert.Equal(t, 0.5, received.Progress)
}

func TestTaskSetProgressBar_Hide_Good(t *testing.T) {
	// Progress -1.0 hides the indicator
	_, c, mock := newTestDockService(t)
	r := taskRun(c, "dock.setProgressBar", TaskSetProgressBar{Progress: -1.0})
	require.True(t, r.OK)
	assert.Equal(t, -1.0, mock.progress)
}

func TestTaskSetProgressBar_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.progressErr = assert.AnError

	r := taskRun(c, "dock.setProgressBar", TaskSetProgressBar{Progress: 0.5})
	assert.False(t, r.OK)
}

func TestTaskSetProgressBar_Ugly(t *testing.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.setProgressBar").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

// --- TaskBounce ---

func TestTaskBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionBounceStarted
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionBounceStarted); ok {
			received = &a
		}
		return core.Result{OK: true}
	})

	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	require.True(t, r.OK)
	assert.True(t, mock.bounceCalled)
	assert.Equal(t, BounceInformational, mock.bounceType)
	requestID, ok := r.Value.(int)
	require.True(t, ok)
	assert.Equal(t, 1, requestID)
	require.NotNil(t, received)
	assert.Equal(t, BounceInformational, received.BounceType)
}

func TestTaskBounce_Critical_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceCritical})
	require.True(t, r.OK)
	assert.Equal(t, BounceCritical, mock.bounceType)
	requestID, ok := r.Value.(int)
	require.True(t, ok)
	assert.Equal(t, 1, requestID)
}

func TestTaskBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.bounceErr = assert.AnError

	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	assert.False(t, r.OK)
}

func TestTaskBounce_Ugly(t *testing.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.bounce").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

// --- TaskStopBounce ---

func TestTaskStopBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	// Start a bounce to get a requestID
	r := taskRun(c, "dock.bounce", TaskBounce{BounceType: BounceInformational})
	require.True(t, r.OK)
	requestID := r.Value.(int)

	r2 := taskRun(c, "dock.stopBounce", TaskStopBounce{RequestID: requestID})
	require.True(t, r2.OK)
	assert.True(t, mock.stopBounceCalled)
}

func TestTaskStopBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.stopBounceErr = assert.AnError

	r := taskRun(c, "dock.stopBounce", TaskStopBounce{RequestID: 1})
	assert.False(t, r.OK)
}

func TestTaskStopBounce_Ugly(t *testing.T) {
	// No dock service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("dock.stopBounce").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

func TestTaskRemoveBadge_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.removeErr = assert.AnError

	r := taskRun(c, "dock.removeBadge", TaskRemoveBadge{})
	assert.False(t, r.OK)
}

func TestQueryVisible_Ugly(t *testing.T) {
	// Dock icon initially hidden
	mock := &mockPlatform{visible: false}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryVisible{})
	require.True(t, r.OK)
	assert.Equal(t, false, r.Value)
}
