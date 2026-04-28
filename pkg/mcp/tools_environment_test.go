package mcp

import (
	"context"
	"errors"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/environment"
)

func newEnvironmentToolsTestSubsystem(t *core.T, query func(core.Query) core.Result) *Subsystem {
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

func TestToolsEnvironment_themeGet_Good(t *core.T) {
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
	core.RequireNoError(t, err)
	core.AssertTrue(t, out.Theme.IsDark)
	core.AssertEqual(t, "dark", out.Theme.Theme)
}

func TestToolsEnvironment_themeGet_Bad(t *core.T) {
	sub := newEnvironmentToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{OK: false, Value: "theme backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.themeGet(context.Background(), nil, ThemeGetInput{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "theme query failed")
}

func TestToolsEnvironment_themeGet_Ugly(t *core.T) {
	sub := newEnvironmentToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{OK: true, Value: errors.New("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.themeGet(context.Background(), nil, ThemeGetInput{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "unexpected result type")
}
