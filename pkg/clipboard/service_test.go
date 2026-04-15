// pkg/clipboard/service_test.go
package clipboard

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	text     string
	ok       bool
	image    []byte
	hasImage bool
}

func (m *mockPlatform) Text() (string, bool) { return m.text, m.ok }
func (m *mockPlatform) SetText(text string) bool {
	m.text = text
	m.ok = text != ""
	return true
}
func (m *mockPlatform) Image() ([]byte, bool) { return m.image, m.hasImage }
func (m *mockPlatform) SetImage(data []byte) bool {
	m.image = append([]byte(nil), data...)
	m.hasImage = len(data) > 0
	return true
}

func newTestService(t *testing.T) (*Service, *core.Core) {
	t.Helper()
	c := core.New(
		core.WithService(Register(&mockPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	svc := core.MustServiceFor[*Service](c, "clipboard")
	return svc, c
}

func TestRegister_Good(t *testing.T) {
	svc, _ := newTestService(t)
	assert.NotNil(t, svc)
}

func TestQueryText_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryText{})
	require.True(t, r.OK)
	content := r.Value.(ClipboardContent)
	assert.Equal(t, "hello", content.Text)
	assert.True(t, content.HasContent)
}

func TestQueryText_Bad(t *testing.T) {
	// No clipboard service registered
	c := core.New(core.WithServiceLock())
	r := c.QUERY(QueryText{})
	assert.False(t, r.OK)
}

func TestTaskSetText_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.setText").Run(context.Background(), core.NewOptions(
		core.Option{Key: "text", Value: "world"},
	))
	require.True(t, r.OK)
	assert.Equal(t, true, r.Value)

	// Verify via query
	qr := c.QUERY(QueryText{})
	assert.Equal(t, "world", qr.Value.(ClipboardContent).Text)
}

func TestTaskClear_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.clear").Run(context.Background(), core.NewOptions())
	require.True(t, r.OK)

	// Verify empty
	qr := c.QUERY(QueryText{})
	assert.Equal(t, "", qr.Value.(ClipboardContent).Text)
	assert.False(t, qr.Value.(ClipboardContent).HasContent)
}

func TestQueryImage_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.Action("clipboard.setImage").Run(context.Background(), core.NewOptions(
		core.Option{Key: "data", Value: []byte{0x89, 0x50, 0x4e, 0x47}},
	))
	require.True(t, r.OK)
	assert.Equal(t, true, r.Value)

	result := c.QUERY(QueryImage{})
	require.True(t, result.OK)
	content := result.Value.(ImageContent)
	assert.True(t, content.HasImage)
	assert.Equal(t, []byte{0x89, 0x50, 0x4e, 0x47}, content.Data)
}
