package display

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	coreerr "forge.lthn.ai/core/go-log"
	"gopkg.in/yaml.v3"
)

type modelCatalogFile struct {
	DefaultModel  string              `yaml:"default_model"`
	SelectedModel string              `yaml:"selected_model"`
	Models        []modelCatalogEntry `yaml:"models"`
}

type modelCatalogMapFile struct {
	DefaultModel  string                       `yaml:"default_model"`
	SelectedModel string                       `yaml:"selected_model"`
	Models        map[string]modelCatalogEntry `yaml:"models"`
}

type modelCatalogEntry struct {
	Name           string `yaml:"name"`
	Architecture   string `yaml:"architecture"`
	QuantBits      int    `yaml:"quant_bits"`
	SizeBytes      int64  `yaml:"size_bytes"`
	Loaded         *bool  `yaml:"loaded"`
	Backend        string `yaml:"backend"`
	SupportsVision *bool  `yaml:"supports_vision"`
}

type modelCatalog struct {
	DefaultModel  string
	SelectedModel string
	Models        []ModelEntry
}

func modelCatalogPath(guiConfigPath string) string {
	base := strings.TrimSpace(guiConfigPath)
	if base != "" {
		configDir := filepath.Dir(base)
		if strings.EqualFold(filepath.Base(configDir), "gui") {
			return filepath.Join(filepath.Dir(configDir), "models.yaml")
		}
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return ""
	}
	return filepath.Join(home, ".core", "models.yaml")
}

func loadModelCatalog(path string) (modelCatalog, error) {
	if strings.TrimSpace(path) == "" {
		return modelCatalog{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return modelCatalog{}, nil
		}
		return modelCatalog{}, coreerr.E("display.loadModelCatalog", "read model catalog", err)
	}

	var (
		listFile modelCatalogFile
		mapFile  modelCatalogMapFile
		models   []ModelEntry
		errList  error
		errMap   error
	)

	errList = yaml.Unmarshal(data, &listFile)
	if errList == nil {
		models = entriesToModels(listFile.Models)
	}

	errMap = yaml.Unmarshal(data, &mapFile)
	if len(models) == 0 && errMap == nil {
		models = mapEntriesToModels(mapFile.Models)
	}
	if len(models) == 0 && errList != nil && errMap != nil {
		return modelCatalog{}, coreerr.E("display.loadModelCatalog", "unmarshal model catalog", errList)
	}

	defaultModel := firstNonEmpty(strings.TrimSpace(listFile.DefaultModel), strings.TrimSpace(mapFile.DefaultModel))
	selectedModel := firstNonEmpty(strings.TrimSpace(listFile.SelectedModel), strings.TrimSpace(mapFile.SelectedModel))

	return modelCatalog{
		DefaultModel:  defaultModel,
		SelectedModel: selectedModel,
		Models:        models,
	}, nil
}

func entriesToModels(entries []modelCatalogEntry) []ModelEntry {
	models := make([]ModelEntry, 0, len(entries))
	for _, entry := range entries {
		name := strings.TrimSpace(entry.Name)
		if name == "" {
			continue
		}
		model := ModelEntry{
			Name:         name,
			Architecture: strings.TrimSpace(entry.Architecture),
			QuantBits:    entry.QuantBits,
			SizeBytes:    entry.SizeBytes,
			Backend:      strings.TrimSpace(entry.Backend),
		}
		if entry.Loaded != nil {
			model.Loaded = *entry.Loaded
		}
		if entry.SupportsVision != nil {
			model.SupportsVision = *entry.SupportsVision
		}
		models = append(models, model)
	}
	return models
}

func mapEntriesToModels(entries map[string]modelCatalogEntry) []ModelEntry {
	if len(entries) == 0 {
		return nil
	}
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	models := make([]ModelEntry, 0, len(keys))
	for _, key := range keys {
		entry := entries[key]
		if strings.TrimSpace(entry.Name) == "" {
			entry.Name = key
		}
		models = append(models, entriesToModels([]modelCatalogEntry{entry})...)
	}
	return models
}

func (s *ChatStore) LoadModelCatalog(path string) error {
	catalog, err := loadModelCatalog(path)
	if err != nil {
		return err
	}
	if len(catalog.Models) == 0 && catalog.DefaultModel == "" && catalog.SelectedModel == "" {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.models = mergeModelEntries(s.models, catalog.Models)
	if len(s.models) == 0 {
		s.models = defaultModelEntries()
	}

	selected := firstNonEmpty(
		catalog.SelectedModel,
		loadedModelName(s.models),
		s.selectedModel,
		catalog.DefaultModel,
		s.settings.DefaultModel,
		firstModelName(s.models),
	)
	if !containsModelName(s.models, selected) {
		selected = firstModelName(s.models)
	}

	if containsModelName(s.models, selected) {
		for index := range s.models {
			s.models[index].Loaded = s.models[index].Name == selected
		}
	}

	s.selectedModel = selected
	if containsModelName(s.models, catalog.DefaultModel) {
		s.settings.DefaultModel = catalog.DefaultModel
	}
	if !containsModelName(s.models, s.settings.DefaultModel) {
		s.settings.DefaultModel = selected
	}
	return nil
}

func mergeModelEntries(base, overrides []ModelEntry) []ModelEntry {
	if len(base) == 0 {
		base = defaultModelEntries()
	}

	merged := make([]ModelEntry, 0, len(base)+len(overrides))
	indexByName := make(map[string]int, len(base)+len(overrides))
	for _, model := range base {
		clone := model
		merged = append(merged, clone)
		indexByName[model.Name] = len(merged) - 1
	}

	for _, model := range overrides {
		name := strings.TrimSpace(model.Name)
		if name == "" {
			continue
		}
		model.Name = name
		if index, ok := indexByName[name]; ok {
			merged[index] = mergeModelEntry(merged[index], model)
			continue
		}
		merged = append(merged, model)
		indexByName[name] = len(merged) - 1
	}

	sort.SliceStable(merged, func(i, j int) bool {
		if merged[i].Loaded != merged[j].Loaded {
			return merged[i].Loaded
		}
		return merged[i].Name < merged[j].Name
	})
	return merged
}

func mergeModelEntry(base, override ModelEntry) ModelEntry {
	merged := base
	if strings.TrimSpace(override.Architecture) != "" {
		merged.Architecture = strings.TrimSpace(override.Architecture)
	}
	if override.QuantBits != 0 {
		merged.QuantBits = override.QuantBits
	}
	if override.SizeBytes != 0 {
		merged.SizeBytes = override.SizeBytes
	}
	if strings.TrimSpace(override.Backend) != "" {
		merged.Backend = strings.TrimSpace(override.Backend)
	}
	if override.Loaded {
		merged.Loaded = true
	}
	if override.SupportsVision {
		merged.SupportsVision = true
	}
	return merged
}

func firstModelName(models []ModelEntry) string {
	if len(models) == 0 {
		return ""
	}
	return models[0].Name
}

func loadedModelName(models []ModelEntry) string {
	for _, model := range models {
		if model.Loaded {
			return model.Name
		}
	}
	return ""
}

func containsModelName(models []ModelEntry, name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}
	return slices.ContainsFunc(models, func(model ModelEntry) bool {
		return model.Name == name
	})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
