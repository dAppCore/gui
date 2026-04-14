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
	visible       bool
	badge         string
	hasBadge      bool
	progress      float64
	hasProgress   bool
	bounceID      int
	bounceCount   int
	stopBounceIDs []int
	showErr       error
	hideErr       error
	badgeErr      error
	removeErr     error
	progressErr   error
	bounceErr     error
	stopBounceErr error
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

func (m *mockPlatform) SetProgressBar(value float64) error {
	if m.progressErr != nil {
		return m.progressErr
	}
	m.progress = value
	m.hasProgress = value >= 0
	return nil
}

func (m *mockPlatform) Bounce(bounceType BounceType) (int, error) {
	if m.bounceErr != nil {
		return 0, m.bounceErr
	}
	m.bounceCount++
	m.bounceID++
	return m.bounceID, nil
}

func (m *mockPlatform) StopBounce(bounceID int) error {
	if m.stopBounceErr != nil {
		return m.stopBounceErr
	}
	m.stopBounceIDs = append(m.stopBounceIDs, bounceID)
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

func TestTaskSetProgressBar_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionProgressChanged
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionProgressChanged); ok {
			received = &a
		}
		return nil
	})

	_, handled, err := c.PERFORM(TaskSetProgressBar{Value: 0.75})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 0.75, mock.progress)
	assert.True(t, mock.hasProgress)
	require.NotNil(t, received)
	assert.Equal(t, 0.75, received.Value)
}

func TestTaskSetProgressBar_Remove_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	_, _, _ = c.PERFORM(TaskSetProgressBar{Value: 0.5})
	_, handled, err := c.PERFORM(TaskSetProgressBar{Value: -1.0})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.False(t, mock.hasProgress)
}

func TestTaskSetProgressBar_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.progressErr = assert.AnError

	_, handled, err := c.PERFORM(TaskSetProgressBar{Value: 0.5})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	var received *ActionBounceStarted
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionBounceStarted); ok {
			received = &a
		}
		return nil
	})

	result, handled, err := c.PERFORM(TaskBounce{Type: BounceInformational})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1, mock.bounceCount)
	bounceID := result.(int)
	assert.Equal(t, 1, bounceID)
	require.NotNil(t, received)
	assert.Equal(t, 1, received.BounceID)
	assert.Equal(t, BounceInformational, received.Type)
}

func TestTaskBounce_Critical_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	result, handled, err := c.PERFORM(TaskBounce{Type: BounceCritical})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, 1, mock.bounceCount)
	assert.Equal(t, 1, result.(int))
}

func TestTaskBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.bounceErr = assert.AnError

	_, handled, err := c.PERFORM(TaskBounce{Type: BounceInformational})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskStopBounce_Good(t *testing.T) {
	_, c, mock := newTestDockService(t)

	result, _, _ := c.PERFORM(TaskBounce{Type: BounceInformational})
	bounceID := result.(int)

	_, handled, err := c.PERFORM(TaskStopBounce{BounceID: bounceID})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Contains(t, mock.stopBounceIDs, bounceID)
}

func TestTaskStopBounce_Bad(t *testing.T) {
	_, c, mock := newTestDockService(t)
	mock.stopBounceErr = assert.AnError

	_, handled, err := c.PERFORM(TaskStopBounce{BounceID: 99})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestBounceType_Ugly(t *testing.T) {
	// BounceInformational and BounceCritical must be distinct constants.
	assert.NotEqual(t, BounceInformational, BounceCritical)
}

func TestTaskSetProgressBar_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskSetProgressBar{})
	assert.False(t, handled)
}

func TestTaskBounce_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskBounce{})
	assert.False(t, handled)
}

func TestTaskStopBounce_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskStopBounce{})
	assert.False(t, handled)
}
