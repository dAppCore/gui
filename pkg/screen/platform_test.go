package screen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScreenPlatform_Rect_OriginCornerContains_Good(t *testing.T) {
	r := Rect{X: 10, Y: 20, Width: 100, Height: 50}

	assert.Equal(t, Point{X: 10, Y: 20}, r.Origin())
	assert.Equal(t, Point{X: 110, Y: 70}, r.Corner())
	assert.Equal(t, Point{X: 109, Y: 69}, r.InsideCorner())
	assert.True(t, r.Contains(Point{X: 10, Y: 20}))
	assert.True(t, r.Contains(Point{X: 109, Y: 69}))
	assert.False(t, r.Contains(Point{X: 110, Y: 70}))
	assert.Equal(t, Size{Width: 100, Height: 50}, r.RectSize())
}

func TestScreenPlatform_Rect_IsEmpty_Bad(t *testing.T) {
	assert.True(t, Rect{Width: 0, Height: 50}.IsEmpty())
	assert.True(t, Rect{Width: 50, Height: 0}.IsEmpty())
	assert.True(t, Rect{Width: -1, Height: -1}.IsEmpty())
}

func TestScreenPlatform_Rect_Intersect_Ugly(t *testing.T) {
	left := Rect{X: 0, Y: 0, Width: 100, Height: 100}
	right := Rect{X: 50, Y: 25, Width: 100, Height: 100}

	assert.Equal(t, Rect{X: 50, Y: 25, Width: 50, Height: 75}, left.Intersect(right))
	assert.Equal(t, Rect{}, left.Intersect(Rect{X: 200, Y: 200, Width: 10, Height: 10}))
}

func TestScreenPlatform_Placement_Apply_Good(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		offset    int
		reference OffsetReference
		want      Rect
	}{
		{
			name:      "top-from-begin",
			alignment: AlignTop,
			offset:    25,
			reference: OffsetBegin,
			want:      Rect{X: 125, Y: 100, Width: 200, Height: 100},
		},
		{
			name:      "right-from-end",
			alignment: AlignRight,
			offset:    40,
			reference: OffsetEnd,
			want:      Rect{X: 400, Y: 360, Width: 200, Height: 100},
		},
		{
			name:      "bottom-from-begin",
			alignment: AlignBottom,
			offset:    30,
			reference: OffsetBegin,
			want:      Rect{X: 130, Y: 500, Width: 200, Height: 100},
		},
		{
			name:      "left-from-end",
			alignment: AlignLeft,
			offset:    20,
			reference: OffsetEnd,
			want:      Rect{X: -100, Y: 380, Width: 200, Height: 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := &Screen{
				Bounds:   Rect{X: 100, Y: 200, Width: 200, Height: 100},
				WorkArea: Rect{X: 110, Y: 210, Width: 200, Height: 100},
			}
			parent := &Screen{Bounds: Rect{X: 100, Y: 200, Width: 300, Height: 300}}

			NewPlacement(screen, parent, tt.alignment, tt.offset, tt.reference).Apply()

			assert.Equal(t, tt.want, screen.Bounds)
			assert.Equal(t, tt.want.X+10, screen.WorkArea.X)
			assert.Equal(t, tt.want.Y+10, screen.WorkArea.Y)
		})
	}
}

func TestScreenPlatform_Placement_Apply_Ugly(t *testing.T) {
	screen := &Screen{
		Bounds:   Rect{X: 0, Y: 0, Width: 20, Height: 20},
		WorkArea: Rect{X: 1, Y: 1, Width: 20, Height: 20},
	}
	parent := &Screen{Bounds: Rect{X: 0, Y: 0, Width: 10, Height: 10}}

	NewPlacement(screen, parent, AlignBottom, 999, OffsetEnd).Apply()

	assert.Equal(t, Rect{X: -20, Y: 10, Width: 20, Height: 20}, screen.Bounds)
	assert.Equal(t, Rect{X: -19, Y: 11, Width: 20, Height: 20}, screen.WorkArea)
}
