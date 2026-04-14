package display

import (
	"context"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"

	coreerr "forge.lthn.ai/core/go-log"
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
	Origin      string    `json:"origin"`
	ConversationID string  `json:"conversation_id"`
	Title       string    `json:"title"`
	Role        string    `json:"role"`
	Snippet     string    `json:"snippet"`
	UpdatedAt   time.Time `json:"updated_at"`
}

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

func (s *Service) handleCoreScheme(_ context.Context, path string, params url.Values) (SchemeResponse, error) {
	switch path {
	case "settings":
		settings := s.chat.Settings()
		models := s.chat.Models()
		return SchemeResponse{
			Scheme:      "core",
			Path:        "settings",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"settings": settings,
				"models":   models,
			},
		}, nil
	case "store":
		query := strings.TrimSpace(params.Get("q"))
		results := s.searchStore(query)
		return SchemeResponse{
			Scheme:      "core",
			Path:        "store",
			ContentType: "application/json",
			StatusCode:  200,
			Data: map[string]any{
				"query":   query,
				"results": results,
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
			if query != "" && !strings.Contains(strings.ToLower(msg.Content), query) && !strings.Contains(strings.ToLower(conv.Title), query) {
				continue
			}
			results = append(results, StoreSearchResult{
				Origin:         "chat://conversations",
				ConversationID: conv.ID,
				Title:          conv.Title,
				Role:           msg.Role,
				Snippet:        trimSnippet(msg.Content),
				UpdatedAt:      conv.UpdatedAt,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	return results
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
