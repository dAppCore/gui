package display

import (
	"context"
	"html"
	"net/url"
	"strings"
	"time"

	core "dappco.re/go/core"
	coreerr "dappco.re/go/core/log"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SchemeHandler func(context.Context, string, url.Values) core.Result

type assetHandlerFunc func(w application.ResponseWriter, r *application.Request)

func (f assetHandlerFunc) ServeHTTP(w application.ResponseWriter, r *application.Request) {
	f(w, r)
}

func (s *Service) HandleScheme(scheme string, handler SchemeHandler) {
	if s.schemeHandlers == nil {
		s.schemeHandlers = make(map[string]SchemeHandler)
	}
	s.schemeHandlers[strings.ToLower(strings.TrimSpace(scheme))] = handler
}

func (s *Service) registerDefaultSchemes() {
	s.HandleScheme("core", func(ctx context.Context, route string, query url.Values) core.Result {
		switch route {
		case "settings":
			return s.Core().Action("gui.chat.settings.load").Run(ctx, core.NewOptions())
		case "models":
			state := s.modelState()
			return core.Result{
				Value: map[string]any{
					"content_type": "application/json",
					"body":         core.JSONMarshalString(state),
					"state":        state,
					"models":       s.chatModels(),
				},
				OK: true,
			}
		case "chat":
			if id := coalesce(query.Get("conversation_id"), query.Get("id")); id != "" {
				return s.Core().Action("gui.chat.history").Run(ctx, core.NewOptions(
					core.Option{Key: "conversation_id", Value: id},
				))
			}
			return s.Core().Action("gui.chat.conversations.list").Run(ctx, core.NewOptions())
		case "store":
			results := s.searchAllStorage(query.Get("q"))
			return core.Result{
				Value: map[string]any{
					"content_type": "text/html",
					"body":         s.renderStoreSearchPage(query.Get("q"), results),
					"route":        route,
					"url":          "core://store",
					"query":        query,
					"results":      results,
				},
				OK: true,
			}
		default:
			return core.Result{
				Value: map[string]any{
					"route": route,
					"query": query,
				},
				OK: true,
			}
		}
	})
}

func (s *Service) ResolveScheme(ctx context.Context, rawURL string) core.Result {
	if strings.TrimSpace(rawURL) == "" {
		return core.Result{Value: coreerr.E("display.ResolveScheme", "scheme URL is required", nil), OK: false}
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return core.Result{Value: err, OK: false}
	}
	handler, ok := s.schemeHandlers[strings.ToLower(parsed.Scheme)]
	if !ok {
		return core.Result{Value: coreerr.E("display.ResolveScheme", "no handler registered for scheme "+parsed.Scheme, nil), OK: false}
	}

	route := strings.Trim(strings.TrimPrefix(parsed.Host+parsed.Path, "/"), "/")
	resolved := handler(ctx, route, parsed.Query())
	if !resolved.OK {
		return resolved
	}

	if payload, ok := resolved.Value.(map[string]any); ok {
		if contentType, _ := payload["content_type"].(string); strings.TrimSpace(contentType) != "" {
			if body, ok := payload["body"].(string); ok && strings.TrimSpace(body) != "" {
				return core.Result{Value: payload, OK: true}
			}
			if !strings.EqualFold(contentType, "text/html") {
				payload["body"] = core.JSONMarshalString(payload["state"])
				return core.Result{Value: payload, OK: true}
			}
		}
	}

	body := s.renderSchemeBody(route, resolved.Value)
	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         body,
			"route":        route,
			"url":          rawURL,
		},
		OK: true,
	}
}

func (s *Service) renderSchemeBody(route string, value any) string {
	title := "core://" + route
	pretty := core.JSONMarshalString(value)
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>" +
		html.EscapeString(title) +
		"</title><style>body{font:14px/1.5 ui-monospace,SFMono-Regular,Menlo,monospace;background:#10171f;color:#edf2f7;margin:0}header{padding:16px 20px;border-bottom:1px solid #243447;background:#111827}main{padding:20px}pre{white-space:pre-wrap;word-break:break-word;background:#0b1220;border:1px solid #243447;border-radius:12px;padding:16px}</style></head><body><header><strong>" +
		html.EscapeString(title) +
		"</strong></header><main><pre>" +
		html.EscapeString(pretty) +
		"</pre></main></body></html>"
}

func (s *Service) renderStoreSearchPage(query string, results []StorageEntry) string {
	safeQuery := html.EscapeString(query)
	var items strings.Builder
	if len(results) == 0 && strings.TrimSpace(query) != "" {
		items.WriteString("<p class=\"empty\">No matches found in Core storage.</p>")
	} else if strings.TrimSpace(query) == "" {
		items.WriteString("<p class=\"meta\">Enter a search term to scan Core storage namespaces.</p>")
	} else {
		items.WriteString("<ul>")
		for _, item := range results {
			items.WriteString("<li class=\"result\"><div class=\"origin\">")
			items.WriteString(html.EscapeString(item.Origin))
			items.WriteString("</div><div class=\"bucket\">")
			items.WriteString(html.EscapeString(item.Bucket))
			items.WriteString("</div><div class=\"key\">")
			items.WriteString(html.EscapeString(item.Key))
			items.WriteString("</div><div class=\"value\">")
			items.WriteString(html.EscapeString(item.Value))
			items.WriteString("</div><div class=\"meta\">Updated ")
			items.WriteString(html.EscapeString(item.UpdatedAt.Format(time.RFC3339)))
			items.WriteString("</div></li>")
		}
		items.WriteString("</ul>")
	}
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>core://store</title><style>body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#0f172a;color:#e2e8f0;margin:0}header{padding:20px;border-bottom:1px solid #1e293b;background:linear-gradient(180deg,#111827,#0f172a)}main{padding:20px;display:grid;gap:16px}form{display:flex;gap:8px;flex-wrap:wrap;align-items:center}input{min-width:min(100%,420px);flex:1 1 320px;border-radius:12px;border:1px solid #334155;background:#020617;color:#e2e8f0;padding:12px 14px}button{border:0;border-radius:12px;background:#38bdf8;color:#082f49;padding:12px 16px;font-weight:700;cursor:pointer}section{background:#020617;border:1px solid #1e293b;border-radius:16px;padding:16px}ul{list-style:none;padding:0;margin:0;display:grid;gap:12px}.result{padding:12px;border:1px solid #1e293b;border-radius:12px;background:#0b1220}.meta{color:#94a3b8}.origin{font-weight:700;color:#7dd3fc}.bucket{color:#cbd5e1;font-size:12px;text-transform:uppercase;letter-spacing:.08em}.key{color:#f8fafc;font-weight:600}.value{white-space:pre-wrap;word-break:break-word;color:#e2e8f0}.empty{color:#94a3b8}code{background:#111827;border-radius:8px;padding:2px 6px}</style></head><body><header><strong>core://store</strong><div class=\"meta\">Search the in-memory storage scopes exposed by the preload shim. Query: <code>" +
		safeQuery +
		"</code></div></header><main><section><form method=\"get\" action=\"core://store\"><input name=\"q\" value=\"" +
		safeQuery +
		"\" placeholder=\"Search keys or values\"><button type=\"submit\">Search</button></form></section><section><div id=\"results\">" + items.String() + "</div></section></main></body></html>"
}

func (s *Service) searchAllStorage(query string) []StorageEntry {
	results := s.storage.Search(query)
	if conversations := s.Core().Action("gui.chat.conversations.search").Run(context.Background(), core.NewOptions(core.Option{Key: "q", Value: query})); conversations.OK {
		switch list := conversations.Value.(type) {
		case []any:
			for _, item := range list {
				results = append(results, StorageEntry{
					Origin:    "core://chat",
					Bucket:    "conversation",
					Key:       "summary",
					Value:     core.JSONMarshalString(item),
					UpdatedAt: time.Now(),
				})
			}
		default:
			if payload := core.JSONMarshalString(list); payload != "null" && payload != "" && payload != "[]" {
				results = append(results, StorageEntry{
					Origin:    "core://chat",
					Bucket:    "conversation",
					Key:       "summary",
					Value:     payload,
					UpdatedAt: time.Now(),
				})
			}
		}
	}
	return results
}

func coalesce(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (s *Service) AssetMiddleware() application.Middleware {
	return func(next application.Handler) application.Handler {
		return assetHandlerFunc(func(w application.ResponseWriter, r *application.Request) {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.URL)), "core://") {
				result := s.ResolveScheme(context.Background(), r.URL)
				if !result.OK {
					w.WriteHeader(404)
					_, _ = w.Write([]byte("core route not found"))
					return
				}
				payload, _ := result.Value.(map[string]any)
				body, _ := payload["body"].(string)
				headers := w.Header()
				contentType, _ := payload["content_type"].(string)
				if strings.TrimSpace(contentType) == "" {
					contentType = "text/html"
				}
				headers["Content-Type"] = []string{contentType + "; charset=utf-8"}
				w.WriteHeader(200)
				_, _ = w.Write([]byte(body))
				return
			}
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}
