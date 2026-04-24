package mcp

import (
	"context"
	"errors"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/screen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newScreenToolsTestSubsystem(t *testing.T, query func(core.Query) core.Result) *Subsystem {
	t.Helper()
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if query != nil {
			return query(q)
		}
		return core.Result{}
	})
	return New(c)
}

func TestToolsScreen_screenGet_Good(t *testing.T) {
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryByID:
			return core.Result{
				Value: &screen.Screen{
					ID:   "primary",
					Name: "Main",
					Bounds: screen.Rect{
						X: 0, Y: 0, Width: 1920, Height: 1080,
					},
					WorkArea: screen.Rect{
						X: 0, Y: 0, Width: 1920, Height: 1040,
					},
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	_, out, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "primary"})
	require.NoError(t, err)
	require.NotNil(t, out.Screen)
	assert.Equal(t, "primary", out.Screen.ID)
	assert.Equal(t, "Main", out.Screen.Name)
}

func TestToolsScreen_screenGet_Bad(t *testing.T) {
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(screen.QueryByID); ok {
			return core.Result{OK: false, Value: "screen backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "missing"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "screen query failed")
}

func TestToolsScreen_screenGet_Ugly(t *testing.T) {
	sub := newScreenToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(screen.QueryByID); ok {
			return core.Result{OK: true, Value: errors.New("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.screenGet(context.Background(), nil, ScreenGetInput{ID: "broken"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected result type")
}
