package container

import (
	"context"
	"strings"

	core "dappco.re/go/core"
)

type Service struct {
	*core.ServiceRuntime[TIMOptions]
	manager *TIMManager
}

func NewService(c *core.Core, options TIMOptions) *Service {
	return &Service{
		ServiceRuntime: core.NewServiceRuntime(c, options),
		manager:        NewTIMManager(options),
	}
}

func OptionsFromEnv() TIMOptions {
	return TIMOptions{
		Name:    strings.TrimSpace(core.Env("CORE_TIM_NAME")),
		Image:   strings.TrimSpace(core.Env("CORE_TIM_IMAGE")),
		Command: splitCSV(strings.TrimSpace(core.Env("CORE_TIM_COMMAND"))),
		DataDir: strings.TrimSpace(core.Env("CORE_TIM_DATA_DIR")),
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().Action("container.runtime.detect", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: Detect(), OK: true}
	})
	s.Core().Action("tim.start", func(ctx context.Context, _ core.Options) core.Result {
		state, err := s.manager.Start(ctx)
		return core.Result{}.New(state, err)
	})
	s.Core().Action("tim.stop", func(ctx context.Context, _ core.Options) core.Result {
		state, err := s.manager.Stop(ctx)
		return core.Result{}.New(state, err)
	})
	s.Core().Action("tim.status", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.manager.State(), OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) State() TIMState {
	return s.manager.State()
}

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
