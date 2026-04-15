package display

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/chat"
	"forge.lthn.ai/core/gui/pkg/window"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisplay_Good_WindowOpenIncludesPreload(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	result := c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{
				window.WithName("preload"),
				window.WithURL("https://example.com"),
			},
		}},
	))
	require.True(t, result.OK)
	require.Len(t, platform.Windows, 1)
	assert.NotEmpty(t, platform.Windows[0].ExecJSCalls())
	assert.Contains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.electron")
	assert.Contains(t, platform.Windows[0].ExecJSCalls()[0], "globalThis.core.ml")
}

func TestDisplay_Good_CoreSchemeRoutesThroughBackend(t *testing.T) {
	platform := window.NewMockPlatform()
	c := core.New(
		core.WithService(Register(nil)),
		core.WithService(chat.Register(func(o *chat.Options) { o.StorePath = filepath.Join(t.TempDir(), "chat.db") })),
		core.WithService(window.Register(platform)),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	require.True(t, c.Action("window.open").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskOpenWindow{
			Options: []window.WindowOption{window.WithName("settings")},
		}},
	)).OK)

	require.True(t, c.Action("window.setURL").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: window.TaskSetURL{Name: "settings", URL: "core://settings"}},
	)).OK)

	require.Len(t, platform.Windows, 1)
	assert.True(t, strings.Contains(platform.Windows[0].HTMLContent(), "core://settings"))
}
