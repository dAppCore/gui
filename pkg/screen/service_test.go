// pkg/screen/service_test.go
package screen

import (
	"context"
	"testing"

	core "dappco.re/go/core"
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
func (m *mockPlatform) GetCurrent() *Screen {
	if m.current != nil {
		return m.current
	}
	return m.GetPrimary()
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
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	return mock, c
}

func TestRegister_Good(t *testing.T) {
	_, c := newTestService(t)
	svc := core.MustServiceFor[*Service](c, "screen")
	assert.NotNil(t, svc)
}

func TestQueryAll_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAll{})
	require.True(t, r.OK)
	screens := r.Value.([]Screen)
	assert.Len(t, screens, 2)
}

func TestQueryPrimary_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryPrimary{})
	require.True(t, r.OK)
	scr := r.Value.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)
	assert.True(t, scr.IsPrimary)
}

func TestQueryByID_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryByID{ID: "2"})
	require.True(t, r.OK)
	scr := r.Value.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryByID_Bad(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryByID{ID: "99"})
	require.True(t, r.OK)
	assert.Nil(t, r.Value)
}

func TestQueryAtPoint_Good(t *testing.T) {
	_, c := newTestService(t)

	// Point on primary screen
	r := c.QUERY(QueryAtPoint{X: 100, Y: 100})
	require.True(t, r.OK)
	scr := r.Value.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "Built-in", scr.Name)

	// Point on external screen
	r2 := c.QUERY(QueryAtPoint{X: 3000, Y: 500})
	scr = r2.Value.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

func TestQueryAtPoint_Bad(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryAtPoint{X: -1000, Y: -1000})
	require.True(t, r.OK)
	assert.Nil(t, r.Value)
}

func TestQueryWorkAreas_Good(t *testing.T) {
	_, c := newTestService(t)
	r := c.QUERY(QueryWorkAreas{})
	require.True(t, r.OK)
	areas := r.Value.([]Rect)
	assert.Len(t, areas, 2)
	assert.Equal(t, 38, areas[0].Y) // primary has menu bar offset
}

// --- QueryCurrent ---

func TestQueryCurrent_Good(t *testing.T) {
	// current falls back to primary when not explicitly set
	_, c := newTestService(t)
	r := c.QUERY(QueryCurrent{})
	require.True(t, r.OK)
	scr := r.Value.(*Screen)
	require.NotNil(t, scr)
	assert.True(t, scr.IsPrimary)
	assert.Equal(t, "Built-in", scr.Name)
}

func TestQueryCurrent_Bad(t *testing.T) {
	// no screens at all → GetCurrent returns nil
	mock := &mockPlatform{screens: []Screen{}}
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryCurrent{})
	require.True(t, r.OK)
	assert.Nil(t, r.Value)
}

func TestQueryCurrent_Ugly(t *testing.T) {
	// current is explicitly set to the external screen
	mock := &mockPlatform{
		screens: []Screen{
			{ID: "1", Name: "Built-in", IsPrimary: true,
				Bounds: Rect{X: 0, Y: 0, Width: 2560, Height: 1600}},
			{ID: "2", Name: "External",
				Bounds: Rect{X: 2560, Y: 0, Width: 1920, Height: 1080}},
		},
	}
	mock.current = &mock.screens[1]
	c := core.New(
		core.WithService(Register(mock)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	r := c.QUERY(QueryCurrent{})
	scr := r.Value.(*Screen)
	require.NotNil(t, scr)
	assert.Equal(t, "External", scr.Name)
}

// --- Rect geometry helpers ---

func TestRect_Origin_Good(t *testing.T) {
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.Origin()
	assert.Equal(t, Point{X: 10, Y: 20}, pt)
}

func TestRect_Corner_Good(t *testing.T) {
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.Corner()
	assert.Equal(t, Point{X: 110, Y: 70}, pt)
}

func TestRect_InsideCorner_Good(t *testing.T) {
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}
	pt := r.InsideCorner()
	assert.Equal(t, Point{X: 109, Y: 69}, pt)
}

func TestRect_IsEmpty_Good(t *testing.T) {
	assert.False(t, Rect{X: 0, Y: 0, Width: 1, Height: 1}.IsEmpty())
}

func TestRect_IsEmpty_Bad(t *testing.T) {
	assert.True(t, Rect{}.IsEmpty())
	assert.True(t, Rect{Width: 0, Height: 10}.IsEmpty())
	assert.True(t, Rect{Width: 10, Height: -1}.IsEmpty())
}

func TestRect_Contains_Good(t *testing.T) {
	r := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	assert.True(t, r.Contains(Point{X: 0, Y: 0}))
	assert.True(t, r.Contains(Point{X: 50, Y: 50}))
	assert.True(t, r.Contains(Point{X: 99, Y: 99}))
}

func TestRect_Contains_Bad(t *testing.T) {
	r := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	// exclusive right/bottom edge
	assert.False(t, r.Contains(Point{X: 100, Y: 50}))
	assert.False(t, r.Contains(Point{X: 50, Y: 100}))
	assert.False(t, r.Contains(Point{X: -1, Y: 50}))
}

func TestRect_Contains_Ugly(t *testing.T) {
	// zero-size rect never contains anything
	r := Rect{X: 5, Y: 5, Width: 0, Height: 0}
	assert.False(t, r.Contains(Point{X: 5, Y: 5}))
}

func TestRect_RectSize_Good(t *testing.T) {
	r := Rect{X: 100, Y: 200, Width: 1920, Height: 1080}
	sz := r.RectSize()
	assert.Equal(t, Size{Width: 1920, Height: 1080}, sz)
}

func TestRect_Intersect_Good(t *testing.T) {
	a := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	b := Rect{X: 50, Y: 50, Width: 100, Height: 100}
	overlap := a.Intersect(b)
	assert.Equal(t, Rect{X: 50, Y: 50, Width: 50, Height: 50}, overlap)
}

func TestRect_Intersect_Bad(t *testing.T) {
	// no overlap
	a := Rect{X: 0, Y: 0, Width: 50, Height: 50}
	b := Rect{X: 100, Y: 100, Width: 50, Height: 50}
	overlap := a.Intersect(b)
	assert.True(t, overlap.IsEmpty())
}

func TestRect_Intersect_Ugly(t *testing.T) {
	// empty rect intersects nothing
	a := Rect{X: 0, Y: 0, Width: 0, Height: 0}
	b := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	overlap := a.Intersect(b)
	assert.True(t, overlap.IsEmpty())
}

// --- ScreenPlacement ---

func TestScreenPlacement_Apply_Good(t *testing.T) {
	// secondary placed to the RIGHT of primary, no offset
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 2560, Height: 1600},
		WorkArea: Rect{X: 0, Y: 38, Width: 2560, Height: 1562},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 3000, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 3000, Y: 0, Width: 1920, Height: 1080},
	}
	NewPlacement(secondary, primary, AlignRight, 0, OffsetBegin).Apply()
	assert.Equal(t, 2560, secondary.Bounds.X)
	assert.Equal(t, 0, secondary.Bounds.Y)
	assert.Equal(t, 2560, secondary.WorkArea.X)
}

func TestScreenPlacement_Apply_Bad(t *testing.T) {
	// screen placed ABOVE primary: newY = primary.Y - secondary.Height
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 0, Y: -600, Width: 1920, Height: 600},
		WorkArea: Rect{X: 0, Y: -600, Width: 1920, Height: 600},
	}
	NewPlacement(secondary, primary, AlignTop, 0, OffsetBegin).Apply()
	assert.Equal(t, 0, secondary.Bounds.X)
	assert.Equal(t, -600, secondary.Bounds.Y)
}

func TestScreenPlacement_Apply_Ugly(t *testing.T) {
	// END offset reference — places secondary flush to the bottom-right of parent
	primary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
		WorkArea: Rect{X: 0, Y: 0, Width: 1920, Height: 1080},
	}
	secondary := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 800, Height: 600},
		WorkArea: Rect{X: 0, Y: 0, Width: 800, Height: 600},
	}
	// AlignBottom + OffsetEnd + offset=0 → secondary starts at right edge of parent
	NewPlacement(secondary, primary, AlignBottom, 0, OffsetEnd).Apply()
	assert.Equal(t, 1920-800, secondary.Bounds.X) // flush right
	assert.Equal(t, 1080, secondary.Bounds.Y)      // just below parent
}
