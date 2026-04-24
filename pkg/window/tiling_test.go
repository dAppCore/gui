// pkg/window/tiling_test.go
package window

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTiling_SnapPosition_String_Good(t *testing.T) {
	tests := []struct {
		name string
		pos  SnapPosition
		want string
	}{
		{name: "left", pos: SnapLeft, want: "left"},
		{name: "right", pos: SnapRight, want: "right"},
		{name: "center", pos: SnapCenter, want: "center"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.pos.String())
		})
	}
}

func TestTiling_SnapPosition_String_Bad(t *testing.T) {
	assert.Empty(t, SnapPosition(123).String())
}

func TestTiling_SnapPosition_String_Ugly(t *testing.T) {
	assert.Empty(t, SnapPosition(-1).String())
}

func TestBesideEditor_Good(t *testing.T) {
	rect := BesideEditor(
		Rect{X: 200, Y: 50, Width: 900, Height: 800},
		Size{Width: 1600, Height: 1000},
	)

	assert.Equal(t, Rect{X: 1100, Y: 50, Width: 500, Height: 800}, rect)
}

func TestBesideEditor_Bad(t *testing.T) {
	rect := BesideEditor(
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
		Size{Width: 1600, Height: 900},
	)

	assert.Equal(t, Rect{}, rect)
}

func TestBesideEditor_Ugly(t *testing.T) {
	rect := BesideEditor(
		Rect{X: -100, Y: 20, Width: 700, Height: 900},
		Size{Width: 1600, Height: 1000},
	)

	assert.Equal(t, Rect{X: 600, Y: 20, Width: 1000, Height: 900}, rect)
}

func TestSuggestLayout_Good(t *testing.T) {
	placements := SuggestLayout(
		[]Window{{Name: "editor"}, {Name: "preview"}},
		Rect{X: 0, Y: 0, Width: 1600, Height: 900},
	)

	assert.Equal(t, []WindowPlacement{
		{Name: "editor", Bounds: Rect{X: 0, Y: 0, Width: 989, Height: 900}},
		{Name: "preview", Bounds: Rect{X: 989, Y: 0, Width: 611, Height: 900}},
	}, placements)
}

func TestSuggestLayout_Bad(t *testing.T) {
	placements := SuggestLayout(nil, Rect{X: 0, Y: 0, Width: 1600, Height: 900})

	assert.Nil(t, placements)
}

func TestSuggestLayout_Ugly(t *testing.T) {
	placements := SuggestLayout(
		[]Window{{Name: "one"}, {Name: "two"}, {Name: "three"}},
		Rect{X: 100, Y: 50, Width: 1500, Height: 900},
	)

	assert.Equal(t, []WindowPlacement{
		{Name: "one", Bounds: Rect{X: 100, Y: 50, Width: 750, Height: 450}},
		{Name: "two", Bounds: Rect{X: 850, Y: 50, Width: 750, Height: 450}},
		{Name: "three", Bounds: Rect{X: 100, Y: 500, Width: 750, Height: 450}},
	}, placements)
}

func TestFindEmptySpace_Good(t *testing.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1600, Height: 1000},
		[]Window{{Name: "left", X: 0, Y: 0, Width: 800, Height: 1000}},
		Size{Width: 400, Height: 300},
	)

	assert.True(t, ok)
	assert.Equal(t, Rect{X: 800, Y: 0, Width: 800, Height: 1000}, space)
}

func TestFindEmptySpace_Bad(t *testing.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 800, Height: 600},
		nil,
		Size{Width: 1200, Height: 700},
	)

	assert.False(t, ok)
	assert.Equal(t, Rect{}, space)
}

func TestFindEmptySpace_Ugly(t *testing.T) {
	space, ok := FindEmptySpace(
		Rect{X: 0, Y: 0, Width: 1000, Height: 800},
		[]Window{
			{Name: "header", X: 0, Y: 0, Width: 1000, Height: 100},
			{Name: "rail", X: 0, Y: 100, Width: 300, Height: 700},
			{Name: "inspector", X: 700, Y: 300, Width: 300, Height: 200},
		},
		Size{Width: 200, Height: 200},
	)

	assert.True(t, ok)
	assert.Equal(t, Rect{X: 300, Y: 100, Width: 400, Height: 700}, space)
}

func TestArrangePair_Good(t *testing.T) {
	left, right := ArrangePair(
		Window{Name: "editor", Width: 1280, Height: 800},
		Window{Name: "terminal", Width: 1280, Height: 800},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	assert.Equal(t, Rect{X: 0, Y: 0, Width: 1000, Height: 1000}, left)
	assert.Equal(t, Rect{X: 1000, Y: 0, Width: 1000, Height: 1000}, right)
}

func TestArrangePair_Bad(t *testing.T) {
	left, right := ArrangePair(Window{}, Window{}, Rect{})

	assert.Equal(t, Rect{}, left)
	assert.Equal(t, Rect{}, right)
}

func TestArrangePair_Ugly(t *testing.T) {
	left, right := ArrangePair(
		Window{Name: "chat", Width: 800, Height: 1600},
		Window{Name: "editor", Width: 1600, Height: 900},
		Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	)

	assert.Equal(t, Rect{X: 0, Y: 0, Width: 800, Height: 1000}, left)
	assert.Equal(t, Rect{X: 800, Y: 0, Width: 1200, Height: 1000}, right)
}
