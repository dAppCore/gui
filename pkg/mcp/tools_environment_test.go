package mcp

import (
	"context"
	"errors"
	"testing"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEnvironmentToolsTestSubsystem(t *testing.T, query func(core.Query) core.Result) *Subsystem {
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

func TestToolsEnvironment_themeGet_Good(t *testing.T) {
	sub := newEnvironmentToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{
				Value: environment.ThemeInfo{
					IsDark: true,
					Theme:  "dark",
				},
				OK: true,
			}
		}
		return core.Result{}
	})

	_, out, err := sub.themeGet(context.Background(), nil, ThemeGetInput{})
	require.NoError(t, err)
	assert.True(t, out.Theme.IsDark)
	assert.Equal(t, "dark", out.Theme.Theme)
}

func TestToolsEnvironment_themeGet_Bad(t *testing.T) {
	sub := newEnvironmentToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{OK: false, Value: "theme backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.themeGet(context.Background(), nil, ThemeGetInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "theme query failed")
}

func TestToolsEnvironment_themeGet_Ugly(t *testing.T) {
	sub := newEnvironmentToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{OK: true, Value: errors.New("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.themeGet(context.Background(), nil, ThemeGetInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected result type")
}
