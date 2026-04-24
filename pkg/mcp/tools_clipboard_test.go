package mcp

import (
	"context"
	"errors"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/clipboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClipboardToolsTestSubsystem(t *testing.T, query func(core.Query) core.Result) *Subsystem {
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

func TestToolsClipboard_clipboardRead_Good(t *testing.T) {
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{
				Value: clipboard.ClipboardContent{
					Text:       "hello",
					HasContent: true,
				},
				OK: true,
			}
		}
		return core.Result{}
	})

	_, out, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	require.NoError(t, err)
	assert.Equal(t, "hello", out.Content)
}

func TestToolsClipboard_clipboardRead_Bad(t *testing.T) {
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{OK: false, Value: "clipboard backend unavailable"}
		}
		return core.Result{}
	})

	_, _, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "clipboard query failed")
}

func TestToolsClipboard_clipboardRead_Ugly(t *testing.T) {
	sub := newClipboardToolsTestSubsystem(t, func(q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryText); ok {
			return core.Result{OK: true, Value: errors.New("unexpected payload")}
		}
		return core.Result{}
	})

	_, _, err := sub.clipboardRead(context.Background(), nil, ClipboardReadInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected result type")
}
