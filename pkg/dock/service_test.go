// pkg/dock/service_test.go
package dock

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Platform ---

type mockPlatform struct {
	visible         bool
	badge           string
	hasBadge        bool
	progress        float64
	bounceID        int
	bounceType      BounceType
	bounceCalled    bool
	stopBounceCalled bool
	showErr         error
	hideErr         error
	badgeErr        error
	removeErr       error
	progressErr     error
	bounceErr       error
	stopBounceErr   error
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
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "dock")
	return svc, c, mock
}

// --- Tests ---

func TestRegister_Good(t *testing.T) {
	svc, _, _ := newTestDockService(t)
	assert.NotNil(t, svc)
}

func TestQueryVisible_Good(t *testing.T) {
	_, c, _ := newTestDockService(t)
	result, handled, err := c.QUERY(QueryVisible{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)
}

func TestQueryVisible_Bad(t *testing.T) {
	// No dock service registered — QUERY returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(QueryVisible{})
	assert.False(t, handled)
}

func TestTaskShowIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = false // Start hidden

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskShowIcon{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.visible)
	require.NotNil(t, received)
	assert.True(t, received.Visible)
}

func TestTaskHideIcon_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.visible = true // Start visible

	var received *ActionVisibilityChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionVisibilityChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskHideIcon{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, mock.visible)
	require.NotNil(t, received)
	assert.False(t, received.Visible)
}

func TestTaskSetBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	_, handled, err := c.PERFORM(TaskSetBadge{Label: "3"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "3", mock.badge)
	assert.True(t, mock.hasBadge)
}

func TestTaskSetBadge_EmptyLabel_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	_, handled, err := c.PERFORM(TaskSetBadge{Label: ""})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", mock.badge)
	assert.True(t, mock.hasBadge) // Empty string = default system badge indicator
}

func TestTaskRemoveBadge_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	// Set a badge first
	_, _, _ = c.PERFORM(TaskSetBadge{Label: "5"})

	_, handled, err := c.PERFORM(TaskRemoveBadge{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", mock.badge)
	assert.False(t, mock.hasBadge)
}

func TestTaskShowIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.showErr = assert.AnError

	_, handled, err := c.PERFORM(TaskShowIcon{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskHideIcon_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.hideErr = assert.AnError

	_, handled, err := c.PERFORM(TaskHideIcon{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetBadge_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.badgeErr = assert.AnError

	_, handled, err := c.PERFORM(TaskSetBadge{Label: "3"})
	assert.True(t, handled)
	assert.Error(t, err)
}

// --- TaskSetProgressBar ---

func TestTaskSetProgressBar_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionProgressChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionProgressChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskSetProgressBar{Progress: 0.5})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 0.5, mock.progress)
	require.NotNil(t, received)
	assert.Equal(t, 0.5, received.Progress)
}

func TestTaskSetProgressBar_Hide_Good(t *testing.T) {
	// Progress -1.0 hides the indicator
	_, c, mock := newTestDockService(t)
	_, handled, err := c.PERFORM(TaskSetProgressBar{Progress: -1.0})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, -1.0, mock.progress)
}

func TestTaskSetProgressBar_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.progressErr = assert.AnError

	_, handled, err := c.PERFORM(TaskSetProgressBar{Progress: 0.5})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskSetProgressBar_Ugly(t *testing.T) {
	// No dock service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskSetProgressBar{Progress: 0.5})
	assert.False(t, handled)
}

// --- TaskBounce ---

func TestTaskBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionBounceStarted
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionBounceStarted); ok {
			received = &a
		}
		return nil
	})

	result, handled, err := c.PERFORM(TaskBounce{BounceType: BounceInformational})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.bounceCalled)
	assert.Equal(t, BounceInformational, mock.bounceType)
	requestID, ok := result.(int)
	require.True(t, ok)
	assert.Equal(t, 1, requestID)
	require.NotNil(t, received)
	assert.Equal(t, BounceInformational, received.BounceType)
}

func TestTaskBounce_Critical_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)
	result, handled, err := c.PERFORM(TaskBounce{BounceType: BounceCritical})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, BounceCritical, mock.bounceType)
	requestID, ok := result.(int)
	require.True(t, ok)
	assert.Equal(t, 1, requestID)
}

func TestTaskBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.bounceErr = assert.AnError

	_, handled, err := c.PERFORM(TaskBounce{BounceType: BounceInformational})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskBounce_Ugly(t *testing.T) {
	// No dock service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskBounce{BounceType: BounceInformational})
	assert.False(t, handled)
}

// --- TaskStopBounce ---

func TestTaskStopBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	// Start a bounce to get a requestID
	result, _, err := c.PERFORM(TaskBounce{BounceType: BounceInformational})
	require.NoError(t, err)
	requestID := result.(int)

	_, handled, err := c.PERFORM(TaskStopBounce{RequestID: requestID})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.stopBounceCalled)
}

func TestTaskStopBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.stopBounceErr = assert.AnError

	_, handled, err := c.PERFORM(TaskStopBounce{RequestID: 1})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskStopBounce_Ugly(t *testing.T) {
	// No dock service — PERFORM returns handled=false
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.PERFORM(TaskStopBounce{RequestID: 99})
	assert.False(t, handled)
}

func TestTaskRemoveBadge_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.removeErr = assert.AnError

	_, handled, err := c.PERFORM(TaskRemoveBadge{})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestQueryVisible_Ugly(t *testing.T) {
	// Dock icon initially hidden
	mock := &mockPlatform{visible: false}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	result, handled, err := c.QUERY(QueryVisible{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, false, result)
}
