package window

import (
	"testing"

	"dappco.re/go/gui/pkg/screen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskLayoutBesideEditor_Good(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("cursor"), WithTitle("Cursor"), WithPosition(0, 0), WithSize(1400, 1000),
	}}).OK)
	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("assistant"), WithPosition(10, 10), WithSize(400, 500),
	}}).OK)

	r := taskRun(c, "window.layoutBesideEditor", TaskLayoutBesideEditor{Name: "assistant"})
	require.True(t, r.OK)
	result := r.Value.(LayoutBesideEditorResult)
	assert.Equal(t, "cursor", result.Editor)
	assert.Equal(t, "right", result.Side)
	assert.Equal(t, 1500, result.WindowBounds.X)
	assert.Equal(t, 500, result.WindowBounds.Width)
}

func TestTaskLayoutSuggest_Good(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	r := taskRun(c, "window.layoutSuggest", TaskLayoutSuggest{WindowCount: 2})
	require.True(t, r.OK)
	suggestion := r.Value.(LayoutSuggestion)
	assert.Equal(t, "left-right", suggestion.Mode)
}

func TestTaskScreenFindSpace_Good(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{
		WithName("left"), WithPosition(0, 0), WithSize(1000, 1000),
	}}).OK)

	r := taskRun(c, "window.findSpace", TaskScreenFindSpace{Width: 400, Height: 400})
	require.True(t, r.OK)
	space := r.Value.(ScreenSpace)
	assert.GreaterOrEqual(t, space.X, 1000)
	assert.GreaterOrEqual(t, space.Width, 400)
}

func TestTaskWindowArrangePair_Good(t *testing.T) {
	_, c := newTestWindowServiceWithScreen(t, []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}})

	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("primary")}}).OK)
	require.True(t, taskRun(c, "window.open", TaskOpenWindow{Options: []WindowOption{WithName("secondary")}}).OK)

	r := taskRun(c, "window.arrangePair", TaskWindowArrangePair{Primary: "primary", Secondary: "secondary", Ratio: 0.6})
	require.True(t, r.OK)
	arrangement := r.Value.(PairArrangement)
	assert.Equal(t, "horizontal", arrangement.Orientation)
	assert.Equal(t, 1200, arrangement.Primary.Width)
	assert.Equal(t, 1200, arrangement.Secondary.X)
}
