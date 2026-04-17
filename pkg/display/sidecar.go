package display

import (
	"context"
	"os/exec"
	"reflect"
	"strings"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/deno"
)

func (s *Service) registerSidecarActions() {
	if strings.TrimSpace(core.Env("CORE_DENO_ENABLE")) != "" && s.sidecar == nil {
		manager := s.ensureSidecar()
		if binary := strings.TrimSpace(manager.Status().Binary); binary != "" {
			if _, err := exec.LookPath(binary); err != nil {
				if s != nil && s.ServiceRuntime != nil && s.Core() != nil {
					s.Core().LogWarn(err, "display.registerSidecarActions", "skipping sidecar auto-start; binary unavailable")
				}
				s.sidecar = nil
			} else if _, err := manager.Start(context.Background()); err != nil {
				if s != nil && s.ServiceRuntime != nil && s.Core() != nil {
					s.Core().LogError(err, "display.registerSidecarActions", "failed to start enabled sidecar")
				}
				s.sidecar = nil
			}
		}
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
	s.Core().Action("core.deno.sidecar.start", func(ctx context.Context, _ core.Options) core.Result {
		status, err := s.ensureSidecar().Start(ctx)
		return core.Result{}.New(status, err)
	})
	s.Core().Action("core.deno.sidecar.stop", func(ctx context.Context, _ core.Options) core.Result {
		status, err := s.ensureSidecar().Stop(ctx)
		return core.Result{}.New(status, err)
	})
	s.Core().Action("core.deno.sidecar.status", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.ensureSidecar().Status(), OK: true}
	})
	s.Core().Action("display.sidecar.eval", func(ctx context.Context, opts core.Options) core.Result {
		result, err := s.ensureSidecar().Eval(ctx, opts.String("code"))
		return core.Result{}.New(result, err)
	})
	s.Core().Action("core.deno.sidecar.eval", func(ctx context.Context, opts core.Options) core.Result {
		result, err := s.ensureSidecar().Eval(ctx, opts.String("code"))
		return core.Result{}.New(result, err)
	})
}

func (s *Service) ensureSidecar() *deno.Manager {
	if s.sidecar == nil {
		var coreRef *core.Core
		if s != nil && s.ServiceRuntime != nil {
			coreRef = s.Core()
		}
		s.sidecar = deno.New(deno.Options{
			Binary: strings.TrimSpace(core.Env("CORE_DENO_BINARY")),
			Dir:    strings.TrimSpace(core.Env("CORE_DENO_DIR")),
			Args:   splitCommandArgs(core.Env("CORE_DENO_ARGS")),
			Core:   coreRef,
		})
		s.sidecar.OnEvent(func(event deno.Event) {
			if s.events == nil {
				return
			}
			s.events.Emit(Event{
				Type: EventCustomEvent,
				Data: map[string]any{
					"source": "deno",
					"name":   event.Name,
					"data":   event.Data,
				},
			})
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

func (s *Service) forwardIPCToSidecar(msg core.Message) {
	if s == nil || s.sidecar == nil {
		return
	}
	status := s.sidecar.Status()
	if !status.Running || !status.Connected {
		return
	}
	typeName := ""
	if t := reflect.TypeOf(msg); t != nil {
		typeName = t.String()
	}
	_ = s.sidecar.Emit("core.ipc.message", map[string]any{
		"type": typeName,
		"data": normalizeSidecarValue(msg),
	})
}

func normalizeSidecarValue(value any) any {
	if value == nil {
		return nil
	}
	var normalized any
	if result := core.JSONUnmarshalString(core.JSONMarshalString(value), &normalized); result.OK {
		return normalized
	}
	return map[string]any{"value": core.JSONMarshalString(value)}
}
