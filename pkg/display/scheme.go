package display

import (
	"context"
	"html"
	"net/url"
	"strings"

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
			return s.Core().Action("gui.chat.models").Run(ctx, core.NewOptions())
		case "chat":
			if id := coalesce(query.Get("conversation_id"), query.Get("id")); id != "" {
				return s.Core().Action("gui.chat.history").Run(ctx, core.NewOptions(
					core.Option{Key: "conversation_id", Value: id},
				))
			}
			return s.Core().Action("gui.chat.conversations.list").Run(ctx, core.NewOptions())
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
				headers["Content-Type"] = []string{"text/html; charset=utf-8"}
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
