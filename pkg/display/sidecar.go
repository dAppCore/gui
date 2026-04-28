package display

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	core "dappco.re/go"
	"dappco.re/go/gui/pkg/deno"
	"dappco.re/go/gui/pkg/internal/coreutil"
)

func (s *Service) registerSidecarActions() {
	if strings.TrimSpace(core.Env("CORE_DENO_ENABLE")) != "" && s.sidecar == nil {
		manager, err := s.sidecarForStart()
		if err != nil {
			if s != nil && s.ServiceRuntime != nil && s.Core() != nil {
				s.Core().LogWarn(err, "display.registerSidecarActions", "skipping sidecar auto-start; invalid sidecar environment")
			}
			s.sidecar = nil
		} else if binary := strings.TrimSpace(manager.Status().Binary); binary != "" {
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
		status, err := s.startSidecar(ctx)
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
		status, err := s.startSidecar(ctx)
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
		s.sidecar = s.newSidecar(deno.Options{
			Binary: strings.TrimSpace(core.Env("CORE_DENO_BINARY")),
			Dir:    strings.TrimSpace(core.Env("CORE_DENO_DIR")),
			Args:   splitCommandArgs(core.Env("CORE_DENO_ARGS")),
			Core:   s.coreRef(),
		})
	}
	return s.sidecar
}

func (s *Service) startSidecar(ctx context.Context) (deno.Status, error) {
	manager, err := s.sidecarForStart()
	if err != nil {
		s.sidecar = nil
		return deno.Status{}, err
	}
	return manager.Start(ctx)
}

func (s *Service) sidecarForStart() (*deno.Manager, error) {
	options, err := sidecarLaunchOptions(s.coreRef())
	if err != nil {
		return nil, err
	}
	if s.sidecar == nil || !s.sidecar.Status().Running {
		s.sidecar = s.newSidecar(options)
	}
	return s.sidecar, nil
}

func (s *Service) coreRef() *core.Core {
	if s != nil && s.ServiceRuntime != nil {
		return s.Core()
	}
	return nil
}

func (s *Service) newSidecar(options deno.Options) *deno.Manager {
	manager := deno.New(options)
	manager.OnEvent(func(event deno.Event) {
		if s == nil || s.events == nil {
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
	return manager
}

func sidecarLaunchOptions(coreRef *core.Core) (deno.Options, error) {
	args := splitCommandArgs(core.Env("CORE_DENO_ARGS"))
	if err := validateSidecarArgs(args, coreRef); err != nil {
		return deno.Options{}, err
	}

	binary, err := validateSidecarBinary(core.Env("CORE_DENO_BINARY"))
	if err != nil {
		return deno.Options{}, err
	}
	dir, err := validateSidecarDir(core.Env("CORE_DENO_DIR"))
	if err != nil {
		return deno.Options{}, err
	}

	return deno.Options{
		Binary: binary,
		Dir:    dir,
		Args:   args,
		Core:   coreRef,
	}, nil
}

func validateSidecarArgs(args []string, coreRef *core.Core) error {
	for _, arg := range args {
		flag := denoPermissionFlag(arg)
		if flag == "" {
			continue
		}
		msg := fmt.Sprintf("CORE_DENO_ARGS contains permission flag %s; deno sandbox is being weakened. This is intentional only if you understand the implications.", flag)
		if strings.EqualFold(strings.TrimSpace(core.Env("CORE_DENO_ALLOW_PERMISSIONS")), "true") {
			fmt.Fprintln(os.Stderr, msg)
			if coreRef != nil {
				coreRef.LogWarn(fmt.Errorf("permission flag %s", flag), "display.sidecar.env", msg)
			}
			continue
		}
		return fmt.Errorf("%s Set CORE_DENO_ALLOW_PERMISSIONS=true to allow this intentionally", msg)
	}
	return nil
}

func denoPermissionFlag(arg string) string {
	token := strings.TrimSpace(arg)
	switch {
	case token == "-A" || strings.HasPrefix(token, "-A="):
		return "-A"
	case token == "--allow-all" || strings.HasPrefix(token, "--allow-all="):
		return "--allow-all"
	case strings.HasPrefix(token, "--allow-"):
		if flag, _, ok := strings.Cut(token, "="); ok {
			return flag
		}
		return token
	default:
		return ""
	}
}

func validateSidecarBinary(value string) (string, error) {
	binary := strings.TrimSpace(value)
	if binary == "" {
		return "", nil
	}
	if !filepath.IsAbs(binary) {
		return "", fmt.Errorf("CORE_DENO_BINARY must be an absolute path: %q", binary)
	}
	info, err := os.Stat(binary)
	if err != nil {
		return "", fmt.Errorf("CORE_DENO_BINARY does not exist: %q: %w", binary, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("CORE_DENO_BINARY must point to a file, not a directory: %q", binary)
	}
	if filepath.Base(binary) != "deno" {
		return "", fmt.Errorf("CORE_DENO_BINARY must point to a binary named deno: %q", binary)
	}
	resolved, err := filepath.EvalSymlinks(binary)
	if err != nil {
		return "", fmt.Errorf("CORE_DENO_BINARY could not be resolved: %q: %w", binary, err)
	}
	if filepath.Base(resolved) != "deno" {
		return "", fmt.Errorf("CORE_DENO_BINARY must resolve to a binary named deno: %q", binary)
	}
	return resolved, nil
}

func validateSidecarDir(value string) (string, error) {
	dir := strings.TrimSpace(value)
	if dir == "" {
		return "", nil
	}
	if hasParentPathComponent(dir) {
		return "", fmt.Errorf("CORE_DENO_DIR must not contain .. path components: %q", dir)
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("CORE_DENO_DIR could not be made absolute: %q: %w", dir, err)
	}
	if err := rejectSymlinkPathComponents(abs); err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("CORE_DENO_DIR does not exist: %q: %w", dir, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("CORE_DENO_DIR must be an existing directory: %q", dir)
	}
	return abs, nil
}

func hasParentPathComponent(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func rejectSymlinkPathComponents(path string) error {
	clean := filepath.Clean(path)
	volume := filepath.VolumeName(clean)
	rest := strings.TrimPrefix(clean, volume)
	current := volume
	if strings.HasPrefix(rest, string(filepath.Separator)) {
		current += string(filepath.Separator)
		rest = strings.TrimPrefix(rest, string(filepath.Separator))
	}
	for _, part := range strings.Split(rest, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("CORE_DENO_DIR component does not exist: %q: %w", current, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("CORE_DENO_DIR must not contain symlink component: %q", current)
		}
	}
	return nil
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
	if err := s.sidecar.Emit("core.ipc.message", map[string]any{
		"type": typeName,
		"data": normalizeSidecarValue(msg),
	}); err != nil {
		coreutil.LogWarn(s.Core(), err, "display.emitSidecarIPC", "sidecar emit failed")
	}
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
