package display

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	core "dappco.re/go/core"
	"gopkg.in/yaml.v3"
)

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
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
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
	trimmed := strings.TrimSpace(preloadPath)
	if trimmed == "" {
		return "", errors.New("preload path is empty")
	}
	if filepath.IsAbs(trimmed) {
		return "", errors.New("preload path must be relative")
	}

	baseAbs, err := filepath.Abs(baseDir)
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
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("preload path escapes manifest directory")
	}
	return candidateAbs, nil
}

func discoverManifestPath(pageURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(pageURL))
	if err != nil {
		return "", err
	}
	candidates := make([]string, 0, 4)
	switch parsed.Scheme {
	case "", "file":
		path := parsed.Path
		if path == "" {
			path = pageURL
		}
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				candidates = append(candidates, filepath.Join(path, ".core", "view.yaml"))
			} else {
				dir := filepath.Dir(path)
				candidates = append(candidates, filepath.Join(dir, ".core", "view.yaml"))
				candidates = append(candidates, filepath.Join(filepath.Dir(dir), ".core", "view.yaml"))
			}
		}
	default:
		if parsed.Host != "" {
			home := strings.TrimSpace(os.Getenv("DIR_HOME"))
			if home == "" {
				home = strings.TrimSpace(core.Env("DIR_HOME"))
			}
			if home != "" {
				candidates = append(candidates, filepath.Join(home, ".core", "apps", parsed.Host, ".core", "view.yaml"))
			}
		}
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", errors.New("view manifest not found")
}

func (s *Service) manifestWindowConfig(pageURL string) map[string]ManifestWindow {
	s.manifestMu.Lock()
	defer s.manifestMu.Unlock()
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		return nil
	}
	return loaded.Manifest.Windows
}

func (s *Service) readManifestPreload(baseDir, preloadPath string) ([]byte, error) {
	resolvedPath, err := safeManifestPreloadPath(baseDir, preloadPath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(resolvedPath)
}

type manifestCacheState struct {
	manifestCache map[string]*loadedManifest
	manifestMu    sync.Mutex
}
