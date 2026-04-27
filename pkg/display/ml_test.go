package display

import (
	"context"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/chat"
	"github.com/stretchr/testify/assert"
)

func TestML_ModelState_Good(t *testing.T) {
	t.Setenv("CORE_ML_API_URL", "https://ml.example.test/api/")
	svc, c := newTestDisplayService(t)
	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{
			Value: []chat.ModelEntry{
				{Name: "alpha", Loaded: true, SizeBytes: 2048, Backend: "vulkan"},
				{Name: "beta", Loaded: false, SizeBytes: 4096},
			},
			OK: true,
		}
	})

	state := svc.modelState()
	assert.Equal(t, "https://ml.example.test/api", state.APIURL)
	assert.Equal(t, 1, len(state.Loaded))
	assert.Equal(t, int64(2048), state.VRAMBytes)
	assert.Equal(t, "vulkan", state.Backend)
	assert.Equal(t, "https://ml.example.test/api/v1/chat/completions", state.InferenceURL)
}

func TestML_ModelState_Bad(t *testing.T) {
	t.Setenv("CORE_ML_API_URL", "")
	svc, _ := newTestDisplayService(t)

	state := svc.modelState()
	assert.Equal(t, "http://localhost:8090", state.APIURL)
	assert.Empty(t, state.Loaded)
	assert.Equal(t, int64(0), state.VRAMBytes)
	assert.Equal(t, "local", state.Backend)
}

func TestML_ModelState_Ugly(t *testing.T) {
	svc, c := newTestDisplayService(t)
	c.Action("gui.chat.models", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{
			Value: []chat.ModelEntry{
				{Name: "unloaded", Loaded: false, SizeBytes: 8192},
				{Name: "blank-backend", Loaded: true, SizeBytes: 0},
			},
			OK: true,
		}
	})

	state := svc.modelState()
	assert.Equal(t, int64(0), state.VRAMBytes)
	assert.Equal(t, "local", state.Backend)
	assert.Equal(t, "\"line\\nquote\\\"slash\\\\\"", quoteJS("line\nquote\"slash\\"))
	assert.Equal(t, int64(0), estimateVRAM(state.Available))
}
