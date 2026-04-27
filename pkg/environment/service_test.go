// pkg/environment/service_test.go
package environment

import (
	"context"
	"path/filepath"
	"sync"
	"testing"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	isDark            bool
	info              EnvironmentInfo
	accentColour      string
	openFMErr         error
	openFMPath        string
	openFMSelect      bool
	focusFollowsMouse bool
	themeHandler      func(isDark bool)
	mu                sync.Mutex
}

func (m *mockPlatform) IsDarkMode() bool           { return m.isDark }
func (m *mockPlatform) Info() EnvironmentInfo      { return m.info }
func (m *mockPlatform) AccentColour() string       { return m.accentColour }
func (m *mockPlatform) HasFocusFollowsMouse() bool { return m.focusFollowsMouse }
func (m *mockPlatform) OpenFileManager(path string, selectFile bool) error {
	m.openFMPath = path
	m.openFMSelect = selectFile
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
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "environment")
	assert.NotNil(t, svc)
}

func TestQueryTheme_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryTheme{})
	require.True(t, r.OK)
	theme := r.Value.(ThemeInfo)
	assert.True(t, theme.IsDark)
	assert.Equal(t, "dark", theme.Theme)
}

func TestQueryInfo_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryInfo{})
	require.True(t, r.OK)
	info := r.Value.(EnvironmentInfo)
	assert.Equal(t, "darwin", info.OS)
	assert.Equal(t, "arm64", info.Arch)
}

func TestQueryAccentColour_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAccentColour{})
	require.True(t, r.OK)
	assert.Equal(t, "rgb(0,122,255)", r.Value)
}

func TestTaskOpenFileManager_Good(t *testing.T) {
	mock, c := newTestService(t)
	r := c.Action("environment.openFileManager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "/tmp", Select: true}},
	))
	require.True(t, r.OK)
	assert.Equal(t, filepath.Clean("/tmp"), mock.openFMPath)
	assert.True(t, mock.openFMSelect)
}

func TestTaskOpenFileManager_Bad_InvalidPath(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("environment.openFileManager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "../tmp", Select: false}},
	))
	assert.False(t, r.OK)
	assert.Contains(t, r.Value.(error).Error(), "path must be absolute")
}

func TestThemeChange_ActionBroadcast_Good(t *testing.T) {
	mock, c := newTestService(t)

	// Register a listener that captures the action
	var received *ActionThemeChanged
	var mu sync.Mutex
	c.RegisterAction(func(_ *core.Core, msg core.Message) core.Result {
		if a, ok := msg.(ActionThemeChanged); ok {
			mu.Lock()
			received = &a
			mu.Unlock()
		}
		return core.Result{OK: true}
	})

	// Simulate theme change
	mock.simulateThemeChange(false)

	mu.Lock()
	r := received
	mu.Unlock()
	require.NotNil(t, r)
	assert.False(t, r.IsDark)
}

func TestTaskSetTheme_Good_OverrideAndReset(t *testing.T) {
	mock, c := newTestService(t)
	mock.isDark = false

	r := c.Action("environment.setTheme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "dark"}},
	))
	require.True(t, r.OK)

	theme := c.QUERY(QueryTheme{})
	require.True(t, theme.OK)
	info := theme.Value.(ThemeInfo)
	assert.True(t, info.IsDark)
	assert.Equal(t, "dark", info.Theme)

	r = c.Action("environment.setTheme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "system"}},
	))
	require.True(t, r.OK)

	theme = c.QUERY(QueryTheme{})
	require.True(t, theme.OK)
	info = theme.Value.(ThemeInfo)
	assert.False(t, info.IsDark)
	assert.Equal(t, "light", info.Theme)
}

func TestTaskSetTheme_Bad_Invalid(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("environment.setTheme").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTheme{Theme: "sepia"}},
	))
	assert.False(t, r.OK)
}

// --- GetAccentColor ---

func TestQueryAccentColour_Bad_Empty(t *testing.T) {
	// accent colour := "" — still returns handled with empty string
	mock := &mockPlatform{
		isDark:       false,
		accentColour: "",
		info:         EnvironmentInfo{OS: "linux", Arch: "amd64"},
	}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.True(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryAccentColour{})
	require.True(t, r.OK)
	assert.Equal(t, "", r.Value)
}

func TestQueryAccentColour_Ugly_NoService(t *testing.T) {
	// No environment service — query is unhandled
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryAccentColour{})
	assert.False(t, r.OK)
}

// --- OpenFileManager ---

func TestTaskOpenFileManager_Bad_Error(t *testing.T) {
	// platform returns an error on open
	openErr := coreerr.E("test", "file manager unavailable", nil)
	mock := &mockPlatform{openFMErr: openErr}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.True(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.Action("environment.openFileManager").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: TaskOpenFileManager{Path: "/missing", Select: false}},
	))
	assert.False(t, r.OK)
	err, _ := r.Value.(error)
	assert.ErrorIs(t, err, openErr)
}

func TestTaskOpenFileManager_Ugly_NoService(t *testing.T) {
	// No environment service — action is not registered
	c := core.New(core.WithServiceLock())
	r := c.Action("environment.openFileManager").Run(context.Background(), core.NewOptions())
	assert.False(t, r.OK)
}

// --- HasFocusFollowsMouse ---

func TestQueryFocusFollowsMouse_Good_True(t *testing.T) {
	// platform reports focus-follows-mouse enabled (Linux/X11 sloppy focus)
	mock := &mockPlatform{focusFollowsMouse: true}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.True(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryFocusFollowsMouse{})
	require.True(t, r.OK)
	assert.Equal(t, true, r.Value)
}

func TestQueryFocusFollowsMouse_Bad_False(t *testing.T) {
	// platform reports focus-follows-mouse disabled (Windows/macOS default)
	mock := &mockPlatform{focusFollowsMouse: false}
	c := core.New(core.WithService(Register(mock)), core.WithServiceLock())
	require.True(t, c.ServiceStartup(t.Context(), nil).OK)

	r := c.QUERY(QueryFocusFollowsMouse{})
	require.True(t, r.OK)
	assert.Equal(t, false, r.Value)
}

func TestQueryFocusFollowsMouse_Ugly_NoService(t *testing.T) {
	// No environment service — query is unhandled
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryFocusFollowsMouse{})
	assert.False(t, r.OK)
}
