// pkg/environment/service_test.go
package environment

import (
	"context"
	"sync"
	"testing"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	isDark              bool
	info                EnvironmentInfo
	accentColour        string
	openFMErr           error
	focusFollowsMouse   bool
	themeHandler        func(isDark bool)
	mu                  sync.Mutex
}

func (m *mockPlatform) IsDarkMode() bool        { return m.isDark }
func (m *mockPlatform) Info() EnvironmentInfo    { return m.info }
func (m *mockPlatform) AccentColour() string     { return m.accentColour }
func (m *mockPlatform) HasFocusFollowsMouse() bool { return m.focusFollowsMouse }
func (m *mockPlatform) OpenFileManager(path string, selectFile bool) error {
	return m.openFMErr
}
func (m *mockPlatform) OnThemeChange(handler func(isDark bool)) func() {
	m.mu.Lock()
	m.themeHandler = handler
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		m.themeHandler = nil
		m.mu.Unlock()
	}
}

// simulateThemeChange triggers the stored handler (test helper).
func (m *mockPlatform) simulateThemeChange(isDark bool) {
	m.mu.Lock()
	h := m.themeHandler
	m.mu.Unlock()
	if h != nil {
		h(isDark)
	}
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		isDark:       true,
		accentColour: "rgb(0,122,255)",
		info: EnvironmentInfo{
			OS: "darwin", Arch: "arm64",
			Platform: PlatformInfo{Name: "macOS", Version: "14.0"},
		},
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
	svc := core.MustServiceFor[*Service](c, "environment")
	assert.NotNil(t, svc)
}

func TestQueryTheme_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryTheme{})
	require.NoError(t, err)
	assert.True(t, handled)
	theme := result.(ThemeInfo)
	assert.True(t, theme.IsDark)
	assert.Equal(t, "dark", theme.Theme)
}

func TestQueryInfo_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryInfo{})
	require.NoError(t, err)
	assert.True(t, handled)
	info := result.(EnvironmentInfo)
	assert.Equal(t, "darwin", info.OS)
	assert.Equal(t, "arm64", info.Arch)
}

func TestQueryAccentColour_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAccentColour{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "rgb(0,122,255)", result)
}

func TestTaskOpenFileManager_Good(t *testing.T) {
	_, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskOpenFileManager{Path: "/tmp", Select: true})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestThemeChange_ActionBroadcast_Good(t *testing.T) {
	mock, c := newTestService(t)

	// Register a listener that captures the action
	var received *ActionThemeChanged
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
		if a, ok := msg.(ActionThemeChanged); ok {
			mu.Lock()
			received = &a
			mu.Unlock()
		}
		return nil
	})

	// Simulate theme change
	mock.simulateThemeChange(false)

	mu.Lock()
	r := received
	mu.Unlock()
	require.NotNil(t, r)
	assert.False(t, r.IsDark)
}

// --- GetAccentColor ---

func TestQueryAccentColour_Bad_Empty(t *testing.T) {
	// accent colour := "" — still returns handled with empty string
	mock := &mockPlatform{
		isDark:       false,
		accentColour: "",
		info:         EnvironmentInfo{OS: "linux", Arch: "amd64"},
	}
	c, err := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(t.Context(), nil))

	result, handled, err := c.QUERY(QueryAccentColour{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "", result)
}

func TestQueryAccentColour_Ugly_NoService(t *testing.T) {
	// No environment service — query is unhandled
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryAccentColour{})
	assert.False(t, handled)
}

// --- OpenFileManager ---

func TestTaskOpenFileManager_Bad_Error(t *testing.T) {
	// platform returns an error on open
	openErr := coreerr.E("test", "file manager unavailable", nil)
	mock := &mockPlatform{openFMErr: openErr}
	c, err := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(t.Context(), nil))

	_, handled, err := c.PERFORM(TaskOpenFileManager{Path: "/missing", Select: false})
	assert.True(t, handled)
	assert.ErrorIs(t, err, openErr)
}

func TestTaskOpenFileManager_Ugly_NoService(t *testing.T) {
	// No environment service — task is unhandled
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenFileManager{Path: "/tmp", Select: false})
	assert.False(t, handled)
}

// --- HasFocusFollowsMouse ---

func TestQueryFocusFollowsMouse_Good_True(t *testing.T) {
	// platform reports focus-follows-mouse enabled (Linux/X11 sloppy focus)
	mock := &mockPlatform{focusFollowsMouse: true}
	c, err := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(t.Context(), nil))

	result, handled, err := c.QUERY(QueryFocusFollowsMouse{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)
}

func TestQueryFocusFollowsMouse_Bad_False(t *testing.T) {
	// platform reports focus-follows-mouse disabled (Windows/macOS default)
	mock := &mockPlatform{focusFollowsMouse: false}
	c, err := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(t.Context(), nil))

	result, handled, err := c.QUERY(QueryFocusFollowsMouse{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, false, result)
}

func TestQueryFocusFollowsMouse_Ugly_NoService(t *testing.T) {
	// No environment service — query is unhandled
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryFocusFollowsMouse{})
	assert.False(t, handled)
}
