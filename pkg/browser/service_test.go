// pkg/browser/service_test.go
package browser

import (
	"context"
	"errors"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	lastURL  string
	lastPath string
	urlErr   error
	fileErr  error
}

func (m *mockPlatform) OpenURL(url string) error {
	m.lastURL = url
	return m.urlErr
}

func (m *mockPlatform) OpenFile(path string) error {
	m.lastPath = path
	return m.fileErr
}

func newTestBrowserService(t *testing.T, mp *mockPlatform) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "browser")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	mp := &mockPlatform{}
	svc, _ := newTestBrowserService(t, mp)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.platform)
}

func TestTaskOpenURL_Good(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "https://example.com", mp.lastURL)
}

func TestTaskOpenURL_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{urlErr: errors.New("browser not found")}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOpenFile_Good(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenFile{Path: "/tmp/readme.txt"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, "/tmp/readme.txt", mp.lastPath)
}

func TestTaskOpenFile_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{fileErr: errors.New("file not found")}
	_, c := newTestBrowserService(t, mp)

	_, handled, err := c.PERFORM(TaskOpenFile{Path: "/nonexistent"})
	assert.True(t, handled)
	assert.Error(t, err)
}

func TestTaskOpenURL_Bad_NoService(t *testing.T) {
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.PERFORM(TaskOpenURL{URL: "https://example.com"})
	assert.False(t, handled)
}
