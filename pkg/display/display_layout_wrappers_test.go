package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type layoutResultWrapperCase struct {
	name      string
	action    string
	zero      any
	call      func(*Service) (any, error)
	setupGood func(*testing.T, *core.Core)
	wantGood  func(*testing.T, any)
}

func runLayoutResultWrapperCase(t *testing.T, tc layoutResultWrapperCase) {
	t.Helper()

	t.Run("good", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		if tc.setupGood != nil {
			tc.setupGood(t, c)
		}

		got, err := tc.call(svc)
		require.NoError(t, err)
		tc.wantGood(t, got)
	})

	t.Run("bad", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: assert.AnError, OK: false}
		})

		got, err := tc.call(svc)
		require.Error(t, err)
		assert.Equal(t, tc.zero, got)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("ugly-action", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: false}
		})

		got, err := tc.call(svc)
		require.Error(t, err)
		assert.Equal(t, tc.zero, got)
		assert.Contains(t, err.Error(), tc.action)
	})

	t.Run("ugly-type", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: true}
		})

		got, err := tc.call(svc)
		require.Error(t, err)
		assert.Equal(t, tc.zero, got)
		assert.Contains(t, err.Error(), "unexpected result type")
	})
}

func TestDisplay_LayoutDelegationWrappers(t *testing.T) {
	errorCases := []errorOnlyWrapperCase{
		{
			name:   "DeleteLayout",
			action: "window.deleteLayout",
			call: func(svc *Service) error {
				return svc.DeleteLayout("development")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("window.deleteLayout", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskDeleteLayout)
					assert.Equal(t, "development", task.Name)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "TileWindows",
			action: "window.tileWindows",
			call: func(svc *Service) error {
				return svc.TileWindows(window.TileModeGrid, []string{"editor", "terminal"})
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("window.tileWindows", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskTileWindows)
					assert.Equal(t, window.TileModeGrid.String(), task.Mode)
					assert.Equal(t, []string{"editor", "terminal"}, task.Windows)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SnapWindow",
			action: "window.snapWindow",
			call: func(svc *Service) error {
				return svc.SnapWindow("preview", window.SnapCenter)
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("window.snapWindow", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskSnapWindow)
					assert.Equal(t, "preview", task.Name)
					assert.Equal(t, window.SnapCenter.String(), task.Position)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "StackWindows",
			action: "window.stackWindows",
			call: func(svc *Service) error {
				return svc.StackWindows([]string{"editor", "preview"}, 24, 18)
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("window.stackWindows", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskStackWindows)
					assert.Equal(t, []string{"editor", "preview"}, task.Windows)
					assert.Equal(t, 24, task.OffsetX)
					assert.Equal(t, 18, task.OffsetY)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ApplyWorkflowLayout",
			action: "window.applyWorkflow",
			call: func(svc *Service) error {
				return svc.ApplyWorkflowLayout(window.WorkflowCoding)
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("window.applyWorkflow", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(window.TaskApplyWorkflow)
					assert.Equal(t, window.WorkflowCoding.String(), task.Workflow)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "LayoutBesideEditor",
		action: "window.layoutBesideEditor",
		zero:   window.LayoutBesideEditorResult{},
		call: func(svc *Service) (any, error) {
			return svc.LayoutBesideEditor("preview", "code", "right", 0.62)
		},
		setupGood: func(t *testing.T, c *core.Core) {
			t.Helper()
			c.Action("window.layoutBesideEditor", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskLayoutBesideEditor)
				assert.Equal(t, "preview", task.Name)
				assert.Equal(t, "code", task.Editor)
				assert.Equal(t, "right", task.Side)
				assert.InDelta(t, 0.62, task.Ratio, 0.0001)
				return core.Result{
					Value: window.LayoutBesideEditorResult{
						Editor: "code",
						EditorBounds: window.WindowBounds{
							X: 10, Y: 20, Width: 640, Height: 800,
						},
						WindowBounds: window.WindowBounds{
							X: 650, Y: 20, Width: 640, Height: 800,
						},
						Side:     "right",
						ScreenID: "screen-1",
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *testing.T, got any) {
			t.Helper()
			result := got.(window.LayoutBesideEditorResult)
			assert.Equal(t, "code", result.Editor)
			assert.Equal(t, "right", result.Side)
			assert.Equal(t, "screen-1", result.ScreenID)
		},
	})

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "FindScreenSpace",
		action: "window.findSpace",
		zero:   window.ScreenSpace{},
		call: func(svc *Service) (any, error) {
			return svc.FindScreenSpace("screen-1", 800, 600, 24)
		},
		setupGood: func(t *testing.T, c *core.Core) {
			t.Helper()
			c.Action("window.findSpace", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskScreenFindSpace)
				assert.Equal(t, "screen-1", task.ScreenID)
				assert.Equal(t, 800, task.Width)
				assert.Equal(t, 600, task.Height)
				assert.Equal(t, 24, task.Padding)
				return core.Result{
					Value: window.ScreenSpace{
						ScreenID: "screen-1",
						X:        100,
						Y:        120,
						Width:    800,
						Height:   600,
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *testing.T, got any) {
			t.Helper()
			space := got.(window.ScreenSpace)
			assert.Equal(t, "screen-1", space.ScreenID)
			assert.Equal(t, 100, space.X)
			assert.Equal(t, 120, space.Y)
			assert.Equal(t, 800, space.Width)
			assert.Equal(t, 600, space.Height)
		},
	})

	runLayoutResultWrapperCase(t, layoutResultWrapperCase{
		name:   "ArrangeWindowPair",
		action: "window.arrangePair",
		zero:   window.PairArrangement{},
		call: func(svc *Service) (any, error) {
			return svc.ArrangeWindowPair("editor", "preview", "screen-1", 0.55)
		},
		setupGood: func(t *testing.T, c *core.Core) {
			t.Helper()
			c.Action("window.arrangePair", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(window.TaskWindowArrangePair)
				assert.Equal(t, "editor", task.Primary)
				assert.Equal(t, "preview", task.Secondary)
				assert.Equal(t, "screen-1", task.ScreenID)
				assert.InDelta(t, 0.55, task.Ratio, 0.0001)
				return core.Result{
					Value: window.PairArrangement{
						Primary: window.WindowBounds{
							X: 0, Y: 0, Width: 800, Height: 600,
						},
						Secondary: window.WindowBounds{
							X: 800, Y: 0, Width: 800, Height: 600,
						},
						Orientation: "horizontal",
						ScreenID:    "screen-1",
					},
					OK: true,
				}
			})
		},
		wantGood: func(t *testing.T, got any) {
			t.Helper()
			arrangement := got.(window.PairArrangement)
			assert.Equal(t, "horizontal", arrangement.Orientation)
			assert.Equal(t, "screen-1", arrangement.ScreenID)
			assert.Equal(t, 800, arrangement.Primary.Width)
			assert.Equal(t, 800, arrangement.Secondary.Width)
		},
	})
}
