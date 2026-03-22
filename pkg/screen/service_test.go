// pkg/screen/service_test.go
package screen

import (
	"context"
	"testing"

	"forge.lthn.ai/core/go/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPlatform struct {
	screens []Screen
}

func (m *mockPlatform) GetAll() []Screen { return m.screens }
func (m *mockPlatform) GetPrimary() *Screen {
	for i := range m.screens {
		if m.screens[i].IsPrimary {
			return &m.screens[i]
		}
	}
	return nil
}

func newTestService(t *testing.T) (*mockPlatform, *core.Core) {
	t.Helper()
	mock := &mockPlatform{
		screens: []Screen{
			{
				ID: "1", Name: "Built-in", IsPrimary: true,
				Bounds:   Rect{X: 0, Y: 0, Width: 2560, Height: 1600},
				WorkArea: Rect{X: 0, Y: 38, Width: 2560, Height: 1562},
				Size:     Size{Width: 2560, Height: 1600},
			},
			{
				ID: "2", Name: "External",
				Bounds:   Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				WorkArea: Rect{X: 2560, Y: 0, Width: 1920, Height: 1080},
				Size:     Size{Width: 1920, Height: 1080},
			},
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
	svc := core.MustServiceFor[*Service](c, "screen")
	assert.NotNil(t, svc)
}

func TestQueryAll_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAll{})
	require.NoError(t, err)
	assert.True(t, handled)
	screens := result.([]Screen)
	assert.Len(t, screens, 2)
}

func TestQueryPrimary_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryPrimary{})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)
	assert.True(t, scr.IsPrimary)
}

func TestQueryByID_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryByID{ID: "2"})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryByID_Bad(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryByID{ID: "99"})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryAtPoint_Good(t *testing.T) {
	_, c := newTestService(t)

	// Point on primary screen
	result, handled, err := c.QUERY(QueryAtPoint{X: 100, Y: 100})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)

	// Point on external screen
	result, _, _ = c.QUERY(QueryAtPoint{X: 3000, Y: 500})
	scr = result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryAtPoint_Bad(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryAtPoint{X: -1000, Y: -1000})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryWorkAreas_Good(t *testing.T) {
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryWorkAreas{})
	require.NoError(t, err)
	assert.True(t, handled)
	areas := result.([]Rect)
	assert.Len(t, areas, 2)
	assert.Equal(t, 38, areas[0].Y) // primary has menu bar offset
}
