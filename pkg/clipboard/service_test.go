// pkg/clipboard/service_test.go
package clipboard

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	text  string
	ok    bool
	img   []byte
	imgOk bool
}

func (m *mockPlatform) Text() (string, bool) { return m.text, m.ok }
func (m *mockPlatform) SetText(text string) bool {
	m.text = text
	m.ok = text != ""
	return true
}
func (m *mockPlatform) Image() ([]byte, bool) { return m.img, m.imgOk }
func (m *mockPlatform) SetImage(data []byte) bool {
	m.img = data
	m.imgOk = len(data) > 0
	return true
}

func newTestService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c, err := core.New(
		core.WithService(Register(&mockPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))
	svc := core.MustServiceFor[*Service](c, "clipboard")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestService(t)
	assert.NotNil(t, svc)
}

func TestQueryText_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryText{})
	require.NoError(t, err)
	assert.True(t, handled)
	content := result.(ClipboardContent)
	assert.Equal(t, "hello", content.Text)
	assert.True(t, content.HasContent)
}

func TestQueryText_Bad(t *testing.T) {
	// No clipboard service registered
	c, _ := core.New(core.WithServiceLock())
	_, handled, _ := c.QUERY(QueryText{})
	assert.False(t, handled)
}

func TestTaskSetText_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.PERFORM(TaskSetText{Text: "world"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Equal(t, true, result)

	// Verify via query
	r, _, _ := c.QUERY(QueryText{})
	assert.Equal(t, "world", r.(ClipboardContent).Text)
}

func TestTaskClear_Good(t *testing.T) {
	_, c := newTestService(t)
	_, handled, err := c.PERFORM(TaskClear{})
	require.NoError(t, err)
	assert.True(t, handled)

	// Verify empty
	r, _, _ := c.QUERY(QueryText{})
	assert.Equal(t, "", r.(ClipboardContent).Text)
	assert.False(t, r.(ClipboardContent).HasContent)
}

func TestQueryImage_Good(t *testing.T) {
	mock := &mockPlatform{img: []byte{1, 2, 3}, imgOk: true}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	result, handled, err := c.QUERY(QueryImage{})
	require.NoError(t, err)
	assert.True(t, handled)
	image := result.(ClipboardImageContent)
	assert.True(t, image.HasContent)
}

func TestTaskSetImage_Good(t *testing.T) {
	mock := &mockPlatform{}
	c, err := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.NoError(t, err)
	require.NoError(t, c.ServiceStartup(context.Background(), nil))

	_, handled, err := c.PERFORM(TaskSetImage{Data: []byte{9, 8, 7}})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.True(t, mock.imgOk)
}
