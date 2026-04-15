package display

import (
	"context"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"

	coreerr "forge.lthn.ai/core/go-log"
	"forge.lthn.ai/core/go/pkg/core"
)

type SchemeResponse struct {
	Scheme      string         `json:"scheme"`
	Path        string         `json:"path"`
	ContentType string         `json:"content_type"`
	StatusCode  int            `json:"status_code"`
	Data        map[string]any `json:"data"`
}

type SchemeHandler func(ctx context.Context, path string, params url.Values) (SchemeResponse, error)

type StoreSearchResult struct {
	Origin         string    `json:"origin"`
	ConversationID string    `json:"conversation_id"`
	Title          string    `json:"title"`
	Role           string    `json:"role"`
	StorageArea    string    `json:"storage_area,omitempty"`
	Key            string    `json:"key,omitempty"`
	Snippet        string    `json:"snippet"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StoreSearchGroup struct {
	Origin    string              `json:"origin"`
	UpdatedAt time.Time           `json:"updated_at"`
	Results   []StoreSearchResult `json:"results"`
}

type QueryRouteResolve struct {
	URL string `json:"url"`
}

type QueryRouteStore struct {
	Query string `json:"q"`
}

type QueryRouteSettings struct{}

type QueryRouteModels struct{}

func (s *Service) registerBuiltinSchemes() {
	s.HandleScheme("core", s.handleCoreScheme)
}

func (s *Service) HandleScheme(scheme string, handler SchemeHandler) {
	if scheme == "" || handler == nil {
		return
	}
	s.schemes[strings.ToLower(scheme)] = handler
}

func (s *Service) ResolveScheme(ctx context.Context, rawURL string) (SchemeResponse, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return SchemeResponse{}, err
	}
	handler, ok := s.schemes[strings.ToLower(parsed.Scheme)]
	if !ok {
		return SchemeResponse{}, coreerr.E("display.ResolveScheme", "scheme not registered: "+parsed.Scheme, nil)
	}

	path := parsed.Host
	if parsed.Path != "" {
		if path == "" {
			path = strings.TrimPrefix(parsed.Path, "/")
		} else {
			path += parsed.Path
		}
	}

	return handler(ctx, strings.Trim(path, "/"), parsed.Query())
}

func (s *Service) handleRouteQuery(_ *core.Core, q core.Query) (any, bool, error) {
	switch q := q.(type) {
	case QueryRouteResolve:
		response, err := s.ResolveScheme(context.Background(), q.URL)
		return response, true, err
	case QueryRouteStore:
		return s.storeRouteData(q.Query), true, nil
	case QueryRouteSettings:
		return s.settingsRouteData(), true, nil
	case QueryRouteModels:
		return s.modelsRouteData(), true, nil
	default:
		return nil, false, nil
	}
}

func (s *Service) handleCoreScheme(_ context.Context, path string, params url.Values) (SchemeResponse, error) {
	switch path {
	case "settings":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "settings",
			ContentType: "application/json",
			StatusCode:  200,
			Data:        s.settingsRouteData(),
		}, nil
	case "models":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "models",
			ContentType: "application/json",
			StatusCode:  200,
			Data:        s.modelsRouteData(),
		}, nil
	case "store":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "store",
			ContentType: "application/json",
			StatusCode:  200,
			Data:        s.storeRouteData(params.Get("q")),
		}, nil
	case "network":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "network",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"available":     false,
				"connections":   []any{},
				"fleet_nodes":   []any{},
				"peer_count":    0,
				"websocketInfo": s.GetEventInfo(),
			},
		}, nil
	case "agent":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "agent",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"available":        false,
				"dispatch_queue":   []any{},
				"sessions":         []any{},
				"workspace_status": "unavailable",
			},
		}, nil
	case "wallet":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "wallet",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"available":    false,
				"balance":      "0",
				"staking":      []any{},
				"transactions": []any{},
			},
		}, nil
	case "identity":
		return SchemeResponse{
			Scheme:      "core",
			Path:        "identity",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"available":     false,
				"certificates":  []any{},
				"consent_state": "unknown",
				"did":           "",
				"keys":          []any{},
			},
		}, nil
	default:
		return SchemeResponse{}, coreerr.E("display.handleCoreScheme", "unknown core route: "+path, nil)
	}
}

func (s *Service) searchStore(query string) []StoreSearchResult {
	query = strings.ToLower(strings.TrimSpace(query))
	var results []StoreSearchResult

	for _, conv := range s.chat.ListConversations() {
		for _, msg := range conv.Messages {
			haystack := strings.ToLower(strings.Join([]string{conv.Title, messageSearchHaystack(msg)}, " "))
			if query != "" && !strings.Contains(haystack, query) {
				continue
			}
			results = append(results, StoreSearchResult{
				Origin:         "chat://conversations",
				ConversationID: conv.ID,
				Title:          conv.Title,
				Role:           msg.Role,
				Snippet:        trimSnippet(messageStoreSnippet(msg)),
				UpdatedAt:      conv.UpdatedAt,
			})
		}
	}

	if s.browserStorage != nil {
		results = append(results, s.browserStorage.Search(query)...)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	return results
}

func (s *Service) settingsRouteData() map[string]any {
	return map[string]any{
		"settings": s.chat.Settings(),
		"models":   s.chat.Models(),
	}
}

func (s *Service) modelsRouteData() map[string]any {
	models := s.chat.Models()
	loadedModels := make([]ModelEntry, 0, len(models))
	backends := make(map[string]map[string]any)
	var totalSizeBytes int64
	visionCount := 0

	for _, model := range models {
		totalSizeBytes += model.SizeBytes
		if model.SupportsVision {
			visionCount++
		}
		if model.Loaded {
			loadedModels = append(loadedModels, model)
		}

		backend := strings.TrimSpace(model.Backend)
		if backend == "" {
			backend = "unknown"
		}
		summary := backends[backend]
		if summary == nil {
			summary = map[string]any{
				"loaded":           0,
				"model_count":      0,
				"supports_vision":  false,
				"total_size_bytes": int64(0),
			}
		}
		summary["model_count"] = summary["model_count"].(int) + 1
		summary["total_size_bytes"] = summary["total_size_bytes"].(int64) + model.SizeBytes
		if model.Loaded {
			summary["loaded"] = summary["loaded"].(int) + 1
		}
		if model.SupportsVision {
			summary["supports_vision"] = true
		}
		backends[backend] = summary
	}

	return map[string]any{
		"selected_model":   s.chat.SelectedModel(),
		"models":           models,
		"loaded_models":    loadedModels,
		"loaded_count":     len(loadedModels),
		"model_count":      len(models),
		"vision_models":    visionCount,
		"total_size_bytes": totalSizeBytes,
		"backends":         backends,
	}
}

func (s *Service) storeRouteData(query string) map[string]any {
	results := s.searchStore(query)
	return map[string]any{
		"query":   strings.TrimSpace(query),
		"results": results,
		"groups":  groupStoreSearchResults(results),
	}
}

func groupStoreSearchResults(results []StoreSearchResult) []StoreSearchGroup {
	byOrigin := make(map[string][]StoreSearchResult)
	for _, result := range results {
		byOrigin[result.Origin] = append(byOrigin[result.Origin], result)
	}

	groups := make([]StoreSearchGroup, 0, len(byOrigin))
	for origin, groupResults := range byOrigin {
		sort.Slice(groupResults, func(i, j int) bool {
			return groupResults[i].UpdatedAt.After(groupResults[j].UpdatedAt)
		})
		updatedAt := time.Time{}
		if len(groupResults) > 0 {
			updatedAt = groupResults[0].UpdatedAt
		}
		groups = append(groups, StoreSearchGroup{
			Origin:    origin,
			UpdatedAt: updatedAt,
			Results:   groupResults,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].UpdatedAt.Equal(groups[j].UpdatedAt) {
			return groups[i].Origin < groups[j].Origin
		}
		return groups[i].UpdatedAt.After(groups[j].UpdatedAt)
	})
	return groups
}

func messageStoreSnippet(message ChatMessage) string {
	if content := strings.TrimSpace(message.Content); content != "" {
		return content
	}
	if message.Thinking != nil {
		if content := strings.TrimSpace(message.Thinking.Content); content != "" {
			return content
		}
	}
	for _, invocation := range message.ToolCalls {
		if content := strings.TrimSpace(invocation.Result.Content); content != "" {
			return content
		}
		if content := strings.TrimSpace(invocation.Error); content != "" {
			return content
		}
	}
	if len(message.Attachments) > 0 {
		names := make([]string, 0, len(message.Attachments))
		for _, attachment := range message.Attachments {
			names = append(names, attachment.Filename)
		}
		return "Attachments: " + strings.Join(names, ", ")
	}
	return ""
}

func trimSnippet(content string) string {
	content = strings.TrimSpace(content)
	runes := []rune(content)
	if len(runes) <= 120 {
		return content
	}
	return string(runes[:120]) + "..."
}

func (r SchemeResponse) JSON() ([]byte, error) {
	return json.Marshal(r)
}
