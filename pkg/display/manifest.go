package display

import (
	"errors"
	"io"
	"net"
	"net/url"
	"os" // Note: AX-6 — os.Getenv intrinsic, core.Env(...) preferred where reading config
	"path/filepath"
	"sync" // Note: AX-6 — sync.Mutex for manifest cache guard, no core wrapper in pinned core module

	core "dappco.re/go/core"
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
	if s.manifestCache == nil {
		s.manifestCache = make(map[string]*loadedManifest)
	}
	if cached, ok := s.manifestCache[pageURL]; ok {
		s.manifestMu.Unlock()
		return cached, nil
	}
	s.manifestMu.Unlock()

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
		return nil, errors.New("view manifest stream is not readable")
	}
	defer core.CloseStream(stream.Value)
	body, err := io.ReadAll(io.LimitReader(reader, maxViewManifestBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxViewManifestBytes {
		return nil, errors.New("view manifest exceeds 1048576 bytes")
	}
	var manifest ViewManifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return nil, err
	}
	loaded := &loadedManifest{
		Path:     path,
		BaseDir:  manifestBaseDir(path),
		Manifest: manifest,
	}

	s.manifestMu.Lock()
	if s.manifestCache == nil {
		s.manifestCache = make(map[string]*loadedManifest)
	}
	s.manifestCache[pageURL] = loaded
	s.manifestMu.Unlock()
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
		return "", errors.New(label + " is empty")
	}
	if filepath.IsAbs(trimmed) {
		return "", errors.New(label + " must be relative")
	}

	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	baseResolved, err := filepath.EvalSymlinks(baseAbs)
	if err != nil {
		return "", err
	}
	candidateAbs, err := filepath.Abs(filepath.Join(baseAbs, trimmed))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(baseAbs, candidateAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || core.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New(label + " escapes manifest directory")
	}
	if !(&core.Fs{}).NewUnrestricted().Exists(candidateAbs) {
		return "", errors.New(label + " does not exist")
	}
	candidateResolved, err := filepath.EvalSymlinks(candidateAbs)
	if err != nil {
		return "", err
	}
	rel, err = filepath.Rel(baseResolved, candidateResolved)
	if err != nil {
		return "", err
	}
	if rel == ".." || core.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New(label + " escapes manifest directory")
	}
	return candidateResolved, nil
}

func discoverManifestPath(pageURL string) (string, error) {
	parsed, err := url.Parse(core.Trim(pageURL))
	if err != nil {
		return "", err
	}
	fsys := (&core.Fs{}).NewUnrestricted()
	candidates := make([]string, 0, 4)
	switch parsed.Scheme {
	case "", "file":
		path := parsed.Path
		if path == "" {
			path = pageURL
		}
		if fsys.Exists(path) {
			if fsys.IsDir(path) {
				candidates = append(candidates, filepath.Join(path, ".core", "view.yaml"))
			} else {
				dir := filepath.Dir(path)
				candidates = append(candidates, filepath.Join(dir, ".core", "view.yaml"))
				candidates = append(candidates, filepath.Join(filepath.Dir(dir), ".core", "view.yaml"))
			}
		}
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
	for _, candidate := range candidates {
		if fsys.Exists(candidate) {
			return candidate, nil
		}
	}
	return "", errors.New("view manifest not found")
}

func manifestHostPathComponent(parsed *url.URL) (string, error) {
	host := parsed.Hostname()
	if host == "" {
		return "", errors.New("manifest host is empty")
	}
	if err := validateManifestHostPathComponent(host); err != nil {
		return "", err
	}
	return host, nil
}

func validateManifestHostPathComponent(host string) error {
	for i := 0; i < len(host); i++ {
		if host[i] < 0x20 || host[i] == 0x7f {
			return errors.New("manifest host contains control character")
		}
	}
	if containsAny(host, "[]") {
		return errors.New("manifest host contains IPv6 brackets")
	}
	if host == "." || host == ".." || core.HasPrefix(host, "../") || core.Contains(host, "/../") {
		return errors.New("manifest host contains relative path segment")
	}
	if containsAny(host, `/\`) {
		return errors.New("manifest host contains path separator")
	}
	if core.Contains(host, ":") && net.ParseIP(host) == nil {
		return errors.New("manifest host contains invalid colon")
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
		return nil, errors.New("manifest preload content is not text")
	}
	return []byte(content), nil
}

type manifestCacheState struct {
	manifestCache map[string]*loadedManifest
	manifestMu    sync.Mutex
}

func coreResultError(result core.Result, fallback string) error {
	if err, ok := result.Value.(error); ok {
		return err
	}
	return errors.New(fallback)
}
