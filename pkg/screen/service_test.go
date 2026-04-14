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
	current *Screen
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
func (m *mockPlatform) GetCurrent() *Screen { return m.current }

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

// --- QueryCurrent ---

func TestQueryCurrent_Good(t *testing.T) {
	mock, c := newTestService(t)
	mock.current = &mock.screens[1] // set "External" as current

	result, handled, err := c.QUERY(QueryCurrent{})
	require.NoError(t, err)
	assert.True(t, handled)
	scr := result.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryCurrent_Bad_NilWhenNoCurrentScreen(t *testing.T) {
	// current is nil by default
	_, c := newTestService(t)
	result, handled, err := c.QUERY(QueryCurrent{})
	require.NoError(t, err)
	assert.True(t, handled)
	assert.Nil(t, result)
}

func TestQueryCurrent_Ugly_NoServiceRegistered(t *testing.T) {
	c, err := core.New(core.WithServiceLock())
	require.NoError(t, err)
	_, handled, _ := c.QUERY(QueryCurrent{})
	assert.False(t, handled)
}

// --- Rect geometry helpers ---

func TestRect_Contains_Good(t *testing.T) {
	r := Rect{X: 100, Y: 100, Width: 200, Height: 150}
	assert.True(t, r.Contains(100, 100))  // top-left corner (inclusive)
	assert.True(t, r.Contains(200, 175))  // centre
	assert.True(t, r.Contains(299, 249))  // bottom-right corner (exclusive boundary - 1)
}

func TestRect_Contains_Bad(t *testing.T) {
	r := Rect{X: 100, Y: 100, Width: 200, Height: 150}
	assert.False(t, r.Contains(99, 100))  // just left
	assert.False(t, r.Contains(100, 99))  // just above
	assert.False(t, r.Contains(300, 200)) // right boundary (exclusive)
	assert.False(t, r.Contains(200, 250)) // bottom boundary (exclusive)
}

func TestRect_Center_Good(t *testing.T) {
	r := Rect{X: 0, Y: 0, Width: 200, Height: 100}
	cx, cy := r.Center()
	assert.Equal(t, 100, cx)
	assert.Equal(t, 50, cy)
}

func TestRect_Center_Ugly_OddDimensions(t *testing.T) {
	r := Rect{X: 1, Y: 1, Width: 101, Height: 51}
	cx, cy := r.Center()
	assert.Equal(t, 51, cx) // integer division: 1 + 101/2 = 1 + 50 = 51
	assert.Equal(t, 26, cy) // 1 + 51/2 = 1 + 25 = 26
}

func TestRect_Overlaps_Good(t *testing.T) {
	a := Rect{X: 0, Y: 0, Width: 200, Height: 200}
	b := Rect{X: 100, Y: 100, Width: 200, Height: 200}
	assert.True(t, a.Overlaps(b))
	assert.True(t, b.Overlaps(a))
}

func TestRect_Overlaps_Bad(t *testing.T) {
	a := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	b := Rect{X: 200, Y: 200, Width: 100, Height: 100}
	assert.False(t, a.Overlaps(b))
}

func TestRect_Overlaps_Ugly_AdjacentEdge(t *testing.T) {
	// touching at edge — not overlapping (exclusive right/bottom boundary)
	a := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	b := Rect{X: 100, Y: 0, Width: 100, Height: 100}
	assert.False(t, a.Overlaps(b))
}

// --- ScreenPlacement ---

type mockPlacer struct {
	x, y          int
	width, height int
}

func (m *mockPlacer) SetPosition(x, y int)        { m.x = x; m.y = y }
func (m *mockPlacer) SetSize(width, height int)    { m.width = width; m.height = height }

func TestScreenPlacement_Apply_Good(t *testing.T) {
	p := ScreenPlacement{ScreenID: "1", X: 50, Y: 75, Width: 800, Height: 600}
	placer := &mockPlacer{}
	p.Apply(placer)
	assert.Equal(t, 50, placer.x)
	assert.Equal(t, 75, placer.y)
	assert.Equal(t, 800, placer.width)
	assert.Equal(t, 600, placer.height)
}

func TestScreenPlacement_Apply_Bad_ZeroDimensions(t *testing.T) {
	// Zero dimensions should skip SetSize but still call SetPosition
	p := ScreenPlacement{ScreenID: "1", X: 100, Y: 200, Width: 0, Height: 0}
	placer := &mockPlacer{width: 1280, height: 800}
	p.Apply(placer)
	assert.Equal(t, 100, placer.x)
	assert.Equal(t, 200, placer.y)
	// Size should remain unchanged when both dimensions are zero
	assert.Equal(t, 1280, placer.width)
	assert.Equal(t, 800, placer.height)
}

func TestScreenPlacement_Apply_Ugly_NegativeCoords(t *testing.T) {
	// Negative coordinates are valid (multi-monitor setups with negative origin)
	p := ScreenPlacement{ScreenID: "2", X: -1920, Y: 0, Width: 1920, Height: 1080}
	placer := &mockPlacer{}
	p.Apply(placer)
	assert.Equal(t, -1920, placer.x)
	assert.Equal(t, 0, placer.y)
	assert.Equal(t, 1920, placer.width)
	assert.Equal(t, 1080, placer.height)
}
