package display

import (
	"io"
	"net"
	"net/url"
	"os" // Note: AX-6 — os.Getenv intrinsic, core.Env(...) preferred where reading config
	"path/filepath"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/log"
	"gopkg.in/yaml.v3"
)

const maxViewManifestBytes = 1 << 20

type ViewManifest struct {
	AppID       string                    `yaml:"app_id" json:"app_id"`
	Name        string                    `yaml:"name" json:"name"`
	Preloads    []ManifestPreload         `yaml:"preloads" json:"preloads"`
	Windows     map[string]ManifestWindow `yaml:"windows" json:"windows"`
	HLCRF       []HLCRFTemplate           `yaml:"hlcrf" json:"hlcrf"`
	Permissions []string                  `yaml:"permissions" json:"permissions"`
}

type ManifestPreload struct {
	Path    string `yaml:"path" json:"path"`
	Inline  string `yaml:"inline" json:"inline"`
	Enabled *bool  `yaml:"enabled" json:"enabled,omitempty"`
}

type ManifestWindow struct {
	Title   string `yaml:"title" json:"title"`
	Width   int    `yaml:"width" json:"width"`
	Height  int    `yaml:"height" json:"height"`
	Preload bool   `yaml:"preload" json:"preload"`
}

type HLCRFTemplate struct {
	Name     string `yaml:"name" json:"name"`
	Tag      string `yaml:"tag" json:"tag"`
	Template string `yaml:"template" json:"template"`
}

type loadedManifest struct {
	Path     string
	BaseDir  string
	Manifest ViewManifest
}

func (s *Service) loadManifestForOrigin(pageURL string) (*loadedManifest, error) {
	s.manifestMu.Lock()
	defer s.manifestMu.Unlock()
	if s.manifestCache == nil {
		s.manifestCache = make(map[string]*loadedManifest)
	}
	if cached, ok := s.manifestCache[pageURL]; ok {
		return cached, nil
	}

	path, err := discoverManifestPath(pageURL)
	if err != nil {
		return nil, err
	}
	stream := (&core.Fs{}).NewUnrestricted().Open(path)
	if !stream.OK {
		return nil, coreResultError(stream, "failed to open view manifest")
	}
	reader, ok := stream.Value.(io.Reader)
	if !ok {
		core.CloseStream(stream.Value)
		return nil, coreerr.E("display.loadManifestForOrigin", "view manifest stream is not readable", nil)
	}
	defer core.CloseStream(stream.Value)
	body, err := io.ReadAll(io.LimitReader(reader, maxViewManifestBytes+1))
	if err != nil {
		return nil, coreerr.E("display.loadManifestForOrigin", "failed to read view manifest", err)
	}
	if len(body) > maxViewManifestBytes {
		return nil, coreerr.E("display.loadManifestForOrigin", "view manifest exceeds 1048576 bytes", nil)
	}
	var manifest ViewManifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return nil, coreerr.E("display.loadManifestForOrigin", "failed to parse view manifest", err)
	}
	loaded := &loadedManifest{
		Path:     path,
		BaseDir:  manifestBaseDir(path),
		Manifest: manifest,
	}

	s.manifestCache[pageURL] = loaded
	return loaded, nil
}

func manifestBaseDir(manifestPath string) string {
	baseDir := filepath.Dir(manifestPath)
	if filepath.Base(baseDir) == ".core" {
		return filepath.Dir(baseDir)
	}
	return baseDir
}

func safeManifestPreloadPath(baseDir, preloadPath string) (string, error) {
	return safeManifestRelativePath(baseDir, preloadPath, "preload path")
}

func safeManifestRelativePath(baseDir, relativePath, label string) (string, error) {
	trimmed := core.Trim(relativePath)
	if trimmed == "" {
		return "", coreerr.E("display.safeManifestRelativePath", label+" is empty", nil)
	}
	if filepath.IsAbs(trimmed) {
		return "", coreerr.E("display.safeManifestRelativePath", label+" must be relative", nil)
	}

	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to resolve manifest base directory", err)
	}
	baseResolved, err := filepath.EvalSymlinks(baseAbs)
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to resolve manifest base symlinks", err)
	}
	candidateAbs, err := filepath.Abs(filepath.Join(baseAbs, trimmed))
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to resolve manifest relative path", err)
	}
	rel, err := filepath.Rel(baseAbs, candidateAbs)
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to compare manifest relative path", err)
	}
	if rel == ".." || core.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", coreerr.E("display.safeManifestRelativePath", label+" escapes manifest directory", nil)
	}
	if !(&core.Fs{}).NewUnrestricted().Exists(candidateAbs) {
		return "", coreerr.E("display.safeManifestRelativePath", label+" does not exist", os.ErrNotExist)
	}
	candidateResolved, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to resolve manifest relative path symlinks", err)
	}
	rel, err = filepath.Rel(baseResolved, candidateResolved)
	if err != nil {
		return "", coreerr.E("display.safeManifestRelativePath", "failed to compare resolved manifest relative path", err)
	}
	if rel == ".." || core.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", coreerr.E("display.safeManifestRelativePath", label+" escapes manifest directory", nil)
	}
	return candidateResolved, nil
}

func discoverManifestPath(pageURL string) (string, error) {
	trimmed := core.Trim(pageURL)
	fsys := (&core.Fs{}).NewUnrestricted()
	candidates := make([]string, 0, 4)
	if filepath.VolumeName(trimmed) != "" {
		appendLocalManifestCandidates(&candidates, fsys, trimmed)
	} else {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return "", err
		}
		switch parsed.Scheme {
		case "", "file":
			path := parsed.Path
			if path == "" {
				path = trimmed
			}
			appendLocalManifestCandidates(&candidates, fsys, path)
		default:
			if parsed.Host != "" {
				host, err := manifestHostPathComponent(parsed)
				if err != nil {
					return "", err
				}
				home := core.Trim(os.Getenv("DIR_HOME"))
				if home == "" {
					home = core.Trim(core.Env("DIR_HOME"))
				}
				if home != "" {
					candidates = append(candidates, filepath.Join(home, ".core", "apps", host, ".core", "view.yaml"))
				}
			}
		}
	}
	for _, candidate := range candidates {
		if fsys.Exists(candidate) {
			return candidate, nil
		}
	}
	return "", coreerr.E("display.discoverManifestPath", "view manifest not found", nil)
}

func appendLocalManifestCandidates(candidates *[]string, fsys *core.Fs, path string) {
	if candidates == nil || fsys == nil || !fsys.Exists(path) {
		return
	}
	if fsys.IsDir(path) {
		*candidates = append(*candidates, filepath.Join(path, ".core", "view.yaml"))
		return
	}
	dir := filepath.Dir(path)
	*candidates = append(*candidates, filepath.Join(dir, ".core", "view.yaml"))
	*candidates = append(*candidates, filepath.Join(filepath.Dir(dir), ".core", "view.yaml"))
}

func manifestHostPathComponent(parsed *url.URL) (string, error) {
	host := parsed.Hostname()
	if host == "" {
		return "", coreerr.E("display.manifestHostPathComponent", "manifest host is empty", nil)
	}
	if err := validateManifestHostPathComponent(host); err != nil {
		return "", err
	}
	return host, nil
}

func validateManifestHostPathComponent(host string) error {
	for i := 0; i < len(host); i++ {
		if host[i] < 0x20 || host[i] == 0x7f {
			return coreerr.E("display.validateManifestHostPathComponent", "manifest host contains control character", nil)
		}
	}
	if containsAny(host, "[]") {
		return coreerr.E("display.validateManifestHostPathComponent", "manifest host contains IPv6 brackets", nil)
	}
	if host == "." || host == ".." || core.HasPrefix(host, "../") || core.Contains(host, "/../") {
		return coreerr.E("display.validateManifestHostPathComponent", "manifest host contains relative path segment", nil)
	}
	if containsAny(host, `/\`) {
		return coreerr.E("display.validateManifestHostPathComponent", "manifest host contains path separator", nil)
	}
	if core.Contains(host, ":") && net.ParseIP(host) == nil {
		return coreerr.E("display.validateManifestHostPathComponent", "manifest host contains invalid colon", nil)
	}
	return nil
}

func containsAny(value, chars string) bool {
	for _, char := range chars {
		if core.Contains(value, string(char)) {
			return true
		}
	}
	return false
}

func (s *Service) manifestWindowConfig(pageURL string) map[string]ManifestWindow {
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		return nil
	}
	if len(loaded.Manifest.Windows) == 0 {
		return nil
	}
	windows := make(map[string]ManifestWindow, len(loaded.Manifest.Windows))
	for name, cfg := range loaded.Manifest.Windows {
		windows[name] = cfg
	}
	return windows
}

func (s *Service) readManifestPreload(baseDir, preloadPath string) ([]byte, error) {
	resolvedPath, err := safeManifestPreloadPath(baseDir, preloadPath)
	if err != nil {
		return nil, err
	}
	result := (&core.Fs{}).NewUnrestricted().Read(resolvedPath)
	if !result.OK {
		return nil, coreResultError(result, "failed to read manifest preload")
	}
	content, ok := result.Value.(string)
	if !ok {
		return nil, coreerr.E("display.readManifestPreload", "manifest preload content is not text", nil)
	}
	return []byte(content), nil
}

func coreResultError(result core.Result, fallback string) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return coreerr.E("display.coreResultError", fallback, nil)
}
