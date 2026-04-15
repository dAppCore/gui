// pkg/browser/service_test.go
package browser

import (
	"context"
	core "dappco.re/go/core"
	"testing"

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
	c := core.New(
		core.WithService(Register(mp)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
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

	r := c.Action("browser.openURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	require.True(t, r.OK)
	assert.Equal(t, "https://example.com", mp.lastURL)
}

func TestTaskOpenURL_Bad_Scheme(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.openURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "javascript:alert(1)"},
	))
	assert.False(t, r.OK)
	assert.Empty(t, mp.lastURL)
}

func TestTaskOpenURL_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{urlErr: core.NewError("browser not found")}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.openURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	assert.False(t, r.OK)
}

func TestTaskOpenFile_Good(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.openFile").Run(context.Background(), core.NewOptions(
		core.Option{Key: "path", Value: "/tmp/readme.txt"},
	))
	require.True(t, r.OK)
	assert.Equal(t, "/tmp/readme.txt", mp.lastPath)
}

func TestTaskOpenFile_Bad_RelativePath(t *testing.T) {
	mp := &mockPlatform{}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.openFile").Run(context.Background(), core.NewOptions(
		core.Option{Key: "path", Value: "relative/readme.txt"},
	))
	assert.False(t, r.OK)
	assert.Empty(t, mp.lastPath)
}

func TestTaskOpenFile_Bad_PlatformError(t *testing.T) {
	mp := &mockPlatform{fileErr: core.NewError("file not found")}
	_, c := newTestBrowserService(t, mp)

	r := c.Action("browser.openFile").Run(context.Background(), core.NewOptions(
		core.Option{Key: "path", Value: "/nonexistent"},
	))
	assert.False(t, r.OK)
}

func TestTaskOpenURL_Bad_NoService(t *testing.T) {
	c := core.New(core.WithServiceLock())
	r := c.Action("browser.openURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "url", Value: "https://example.com"},
	))
	assert.False(t, r.OK)
}
