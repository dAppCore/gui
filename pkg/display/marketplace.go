package display

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/marketplace"
	coreerr "dappco.re/go/log"
)

const marketplaceHTTPTimeout = 30 * time.Second

var (
	marketplaceHTTPClient = &http.Client{Timeout: marketplaceHTTPTimeout}
	marketplaceGitRunner  func(context.Context, string, ...string) ([]byte, error)
)

type marketplaceListInput struct {
	RegistryURL string `json:"url"`
}

type marketplaceFetchInput struct {
	ManifestURL string `json:"url"`
}

type marketplaceInstallInput struct {
	ManifestURL string `json:"url"`
	InstallDir  string `json:"install_dir,omitempty"`
	GitBinary   string `json:"git_binary,omitempty"`
}

func (s *Service) registerMarketplaceActions() {
	s.Core().Action("display.marketplace.list", func(ctx context.Context, opts core.Options) core.Result {
		input := marketplaceListInput{RegistryURL: marketplaceRegistryURL(opts)}
		if strings.TrimSpace(input.RegistryURL) == "" {
			return core.Result{Value: coreerr.E("display.marketplace.list", "registry url is required", nil), OK: false}
		}
		installer := marketplace.Installer{HTTPClient: marketplaceHTTPClient}
		manifests, err := installer.List(ctx, input.RegistryURL)
		if err != nil {
			return core.Result{Value: coreerr.E("display.marketplace.list", "failed to list marketplace manifests", err), OK: false}
		}
		return core.Result{Value: map[string]any{
			"registry_url": input.RegistryURL,
			"manifests":    manifests,
		}, OK: true}
	})

	s.Core().Action("display.marketplace.fetch", func(ctx context.Context, opts core.Options) core.Result {
		input := marketplaceFetchInput{ManifestURL: strings.TrimSpace(opts.String("url"))}
		if input.ManifestURL == "" {
			return core.Result{Value: coreerr.E("display.marketplace.fetch", "manifest url is required", nil), OK: false}
		}
		installer := marketplace.Installer{HTTPClient: marketplaceHTTPClient}
		manifest, err := installer.FetchManifest(ctx, input.ManifestURL)
		if err != nil {
			return core.Result{Value: coreerr.E("display.marketplace.fetch", "failed to fetch marketplace manifest", err), OK: false}
		}
		return core.Result{Value: manifest, OK: true}
	})

	s.Core().Action("display.marketplace.verify", func(ctx context.Context, opts core.Options) core.Result {
		input := marketplaceFetchInput{ManifestURL: strings.TrimSpace(opts.String("url"))}
		if input.ManifestURL == "" {
			return core.Result{Value: coreerr.E("display.marketplace.verify", "manifest url is required", nil), OK: false}
		}
		installer := marketplace.Installer{HTTPClient: marketplaceHTTPClient}
		manifest, err := installer.Verify(ctx, input.ManifestURL)
		if err != nil {
			return core.Result{Value: coreerr.E("display.marketplace.verify", "failed to verify marketplace manifest", err), OK: false}
		}
		return core.Result{Value: map[string]any{
			"manifest": manifest,
			"digest":   marketplace.DigestManifest(manifest),
		}, OK: true}
	})

	s.Core().Action("display.marketplace.install", func(ctx context.Context, opts core.Options) core.Result {
		input := marketplaceInstallInput{
			ManifestURL: strings.TrimSpace(opts.String("url")),
			InstallDir:  strings.TrimSpace(opts.String("install_dir")),
			GitBinary:   strings.TrimSpace(opts.String("git_binary")),
		}
		if input.ManifestURL == "" {
			return core.Result{Value: coreerr.E("display.marketplace.install", "manifest url is required", nil), OK: false}
		}

		installer := marketplace.Installer{
			HTTPClient: marketplaceHTTPClient,
			GitBinary:  input.GitBinary,
			GitRunner:  marketplaceGitRunner,
			InstallDir: marketplaceInstallRoot(input.InstallDir),
		}
		manifest, err := installer.Verify(ctx, input.ManifestURL)
		if err != nil {
			return core.Result{Value: coreerr.E("display.marketplace.install", "failed to verify marketplace manifest", err), OK: false}
		}
		targetDir, err := installer.Install(ctx, manifest)
		if err != nil {
			return core.Result{Value: coreerr.E("display.marketplace.install", "failed to install marketplace manifest", err), OK: false}
		}
		return core.Result{Value: map[string]any{
			"manifest":    manifest,
			"digest":      marketplace.DigestManifest(manifest),
			"target_dir":  targetDir,
			"install_dir": installer.InstallDir,
		}, OK: true}
	})
}

func marketplaceRegistryURL(opts core.Options) string {
	if url := strings.TrimSpace(opts.String("url")); url != "" {
		return url
	}
	return strings.TrimSpace(core.Env("CORE_MARKETPLACE_REGISTRY_URL"))
}

func marketplaceInstallRoot(raw string) string {
	if trimmed := strings.TrimSpace(raw); trimmed != "" {
		return trimmed
	}
	home := strings.TrimSpace(core.Env("DIR_HOME"))
	if home == "" {
		if configDir, err := os.UserConfigDir(); err == nil && strings.TrimSpace(configDir) != "" {
			return filepath.Join(configDir, "core", "apps")
		}
		if userHome, err := os.UserHomeDir(); err == nil && strings.TrimSpace(userHome) != "" {
			return filepath.Join(userHome, ".core", "apps")
		}
		return ""
	}
	return filepath.Join(home, ".core", "apps")
}
