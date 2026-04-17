package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/clipboard"
	"forge.lthn.ai/core/gui/pkg/environment"
	"forge.lthn.ai/core/gui/pkg/notification"
	"forge.lthn.ai/core/gui/pkg/systray"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorOnlyWrapperCase struct {
	name      string
	action    string
	call      func(*Service) error
	setupGood func(*testing.T, *core.Core)
}

func runErrorOnlyWrapperCase(t *testing.T, tc errorOnlyWrapperCase) {
	t.Helper()

	t.Run("good", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		if tc.setupGood != nil {
			tc.setupGood(t, c)
		}
		require.NoError(t, tc.call(svc))
	})

	t.Run("bad", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: assert.AnError, OK: false}
		})

		err := tc.call(svc)
		require.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("ugly", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.Action(tc.action, func(_ context.Context, _ core.Options) core.Result {
			return core.Result{Value: "unexpected", OK: false}
		})

		err := tc.call(svc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), tc.action)
	})
}

func TestDisplayAPI_TrayWrappers(t *testing.T) {
	cases := []errorOnlyWrapperCase{
		{
			name:   "SetTrayTooltip",
			action: "systray.setTooltip",
			call: func(svc *Service) error {
				return svc.SetTrayTooltip("Helper tooltip")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("systray.setTooltip", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayTooltip)
					assert.Equal(t, "Helper tooltip", task.Tooltip)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SetTrayLabel",
			action: "systray.setLabel",
			call: func(svc *Service) error {
				return svc.SetTrayLabel("Launcher")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("systray.setLabel", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayLabel)
					assert.Equal(t, "Launcher", task.Label)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "SetTrayMenu",
			action: "systray.setMenu",
			call: func(svc *Service) error {
				return svc.SetTrayMenu([]TrayMenuItem{
					{Label: "Open", ActionID: "open"},
					{
						Label:    "More",
						ActionID: "more",
						Children: []TrayMenuItem{{Label: "Nested", ActionID: "nested"}},
					},
				})
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("systray.setMenu", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskSetTrayMenu)
					require.Len(t, task.Items, 2)
					assert.Equal(t, "Open", task.Items[0].Label)
					require.Len(t, task.Items[1].Submenu, 1)
					assert.Equal(t, "nested", task.Items[1].Submenu[0].ActionID)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowTrayMessage",
			action: "systray.showMessage",
			call: func(svc *Service) error {
				return svc.ShowTrayMessage("Status", "Task complete")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("systray.showMessage", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(systray.TaskShowMessage)
					assert.Equal(t, "Status", task.Title)
					assert.Equal(t, "Task complete", task.Message)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}
}

func TestDisplayAPI_ClipboardWrappers(t *testing.T) {
	t.Run("ClearClipboard", func(t *testing.T) {
		runErrorOnlyWrapperCase(t, errorOnlyWrapperCase{
			name:   "ClearClipboard",
			action: "clipboard.clear",
			call: func(svc *Service) error {
				return svc.ClearClipboard()
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("clipboard.clear", func(_ context.Context, opts core.Options) core.Result {
					assert.Equal(t, 0, opts.Len())
					return core.Result{OK: true}
				})
			},
		})
	})

	t.Run("HasClipboard", func(t *testing.T) {
		svc, c := newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{
					Value: clipboard.ClipboardContent{
						Text:       "present",
						HasContent: true,
					},
					OK: true,
				}
			default:
				return core.Result{}
			}
		})
		assert.True(t, svc.HasClipboard())

		svc, c = newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{
					Value: clipboard.ClipboardContent{
						Text:       "",
						HasContent: false,
					},
					OK: true,
				}
			default:
				return core.Result{}
			}
		})
		assert.False(t, svc.HasClipboard())

		svc, c = newTestDisplayAPIService(t)
		c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
			switch q.(type) {
			case clipboard.QueryText:
				return core.Result{OK: false}
			default:
				return core.Result{}
			}
		})
		assert.False(t, svc.HasClipboard())
	})
}

func TestDisplayAPI_NotificationWrappers(t *testing.T) {
	cases := []errorOnlyWrapperCase{
		{
			name:   "ShowNotification",
			action: "notification.send",
			call: func(svc *Service) error {
				return svc.ShowNotification(NotificationOptions{
					ID:       "alert-1",
					Title:    "Deploy",
					Message:  "Done",
					Subtitle: "CI",
				})
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					assert.Equal(t, notification.NotificationOptions{
						ID:       "alert-1",
						Title:    "Deploy",
						Message:  "Done",
						Subtitle: "CI",
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowInfoNotification",
			action: "notification.send",
			call: func(svc *Service) error {
				return svc.ShowInfoNotification("Info", "Ready")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					assert.Equal(t, notification.NotificationOptions{
						Title:   "Info",
						Message: "Ready",
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowWarningNotification",
			action: "notification.send",
			call: func(svc *Service) error {
				return svc.ShowWarningNotification("Warn", "Careful")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					assert.Equal(t, notification.NotificationOptions{
						Title:    "Warn",
						Message:  "Careful",
						Severity: notification.SeverityWarning,
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ShowErrorNotification",
			action: "notification.send",
			call: func(svc *Service) error {
				return svc.ShowErrorNotification("Error", "Failed")
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("notification.send", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskSend)
					assert.Equal(t, notification.NotificationOptions{
						Title:    "Error",
						Message:  "Failed",
						Severity: notification.SeverityError,
					}, task.Options)
					return core.Result{OK: true}
				})
			},
		},
		{
			name:   "ClearNotifications",
			action: "notification.clear",
			call: func(svc *Service) error {
				return svc.ClearNotifications()
			},
			setupGood: func(t *testing.T, c *core.Core) {
				t.Helper()
				c.Action("notification.clear", func(_ context.Context, opts core.Options) core.Result {
					task := opts.Get("task").Value.(notification.TaskClear)
					assert.Empty(t, task.ID)
					return core.Result{OK: true}
				})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runErrorOnlyWrapperCase(t, tc)
		})
	}
}

func TestDisplayAPI_ThemeWrapper(t *testing.T) {
	runErrorOnlyWrapperCase(t, errorOnlyWrapperCase{
		name:   "SetTheme",
		action: "environment.setTheme",
		call: func(svc *Service) error {
			return svc.SetTheme("system")
		},
		setupGood: func(t *testing.T, c *core.Core) {
			t.Helper()
			c.Action("environment.setTheme", func(_ context.Context, opts core.Options) core.Result {
				task := opts.Get("task").Value.(environment.TaskSetTheme)
				assert.Equal(t, "system", task.Theme)
				return core.Result{OK: true}
			})
		},
	})
}

func TestDisplayAPI_GetTrayInfo(t *testing.T) {
	svc, c := newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{
				Value: map[string]any{
					"tooltip": "Ready",
				},
				OK: true,
			}
		default:
			return core.Result{}
		}
	})

	info := svc.GetTrayInfo()
	require.NotNil(t, info)
	assert.Equal(t, "Ready", info["tooltip"])

	svc, c = newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{Value: "unexpected", OK: true}
		default:
			return core.Result{}
		}
	})
	assert.Nil(t, svc.GetTrayInfo())

	svc, c = newTestDisplayAPIService(t)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case systray.QueryInfo:
			return core.Result{OK: false}
		default:
			return core.Result{}
		}
	})
	assert.Nil(t, svc.GetTrayInfo())
}
