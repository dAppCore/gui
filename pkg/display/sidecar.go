package display

import (
	"context"
	"strings"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/deno"
)

func (s *Service) registerSidecarActions() {
	if strings.TrimSpace(core.Env("CORE_DENO_ENABLE")) != "" && s.sidecar == nil {
		manager := deno.New(deno.Options{
			Binary: strings.TrimSpace(core.Env("CORE_DENO_BINARY")),
			Dir:    strings.TrimSpace(core.Env("CORE_DENO_DIR")),
			Args:   splitCommandArgs(core.Env("CORE_DENO_ARGS")),
		})
		s.sidecar = manager
		_, _ = s.sidecar.Start(context.Background())
	}

	s.Core().Action("display.sidecar.start", func(ctx context.Context, _ core.Options) core.Result {
		manager := s.ensureSidecar()
		status, err := manager.Start(ctx)
		return core.Result{}.New(status, err)
	})
	s.Core().Action("display.sidecar.stop", func(ctx context.Context, _ core.Options) core.Result {
		manager := s.ensureSidecar()
		status, err := manager.Stop(ctx)
		return core.Result{}.New(status, err)
	})
	s.Core().Action("display.sidecar.status", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.ensureSidecar().Status(), OK: true}
	})
}

func (s *Service) ensureSidecar() *deno.Manager {
	if s.sidecar == nil {
		s.sidecar = deno.New(deno.Options{
			Binary: strings.TrimSpace(core.Env("CORE_DENO_BINARY")),
			Dir:    strings.TrimSpace(core.Env("CORE_DENO_DIR")),
			Args:   splitCommandArgs(core.Env("CORE_DENO_ARGS")),
		})
	}
	return s.sidecar
}

func splitCommandArgs(value string) []string {
	fields := strings.Fields(strings.TrimSpace(value))
	if len(fields) == 0 {
		return nil
	}
	return fields
}
