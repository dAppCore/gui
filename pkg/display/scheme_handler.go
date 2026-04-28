package display

import (
	"context"
	"net/url"
	"strings"

	core "dappco.re/go"
	coreerr "dappco.re/go/log"
)

type routeDispatchKind uint8

const (
	routeDispatchQuery routeDispatchKind = iota
	routeDispatchAction
)

var coreRouteDispatch = map[string]routeDispatchKind{
	"settings": routeDispatchQuery,
	"store":    routeDispatchQuery,
	"network":  routeDispatchQuery,
	"models":   routeDispatchQuery,
	"agent":    routeDispatchAction,
	"wallet":   routeDispatchAction,
	"identity": routeDispatchAction,
}

type coreSchemeHandler struct {
	core *core.Core
}

// CoreRouteQuery carries a core:// query route target and sanitized URL parameters.
type CoreRouteQuery struct {
	Target string
	Params url.Values
}

// NewCoreSchemeHandler returns the RFC route dispatcher for `core://` URLs.
//
//	handler := display.NewCoreSchemeHandler(c)
func NewCoreSchemeHandler(c *core.Core) RouteSchemeHandler {
	return coreSchemeHandler{core: c}
}

// SchemeHandler exposes the RFC route dispatcher for the active display service.
//
//	handler := svc.SchemeHandler()
func (s *Service) SchemeHandler() RouteSchemeHandler {
	if s == nil || s.ServiceRuntime == nil {
		return coreSchemeHandler{}
	}
	return NewCoreSchemeHandler(s.Core())
}

func (h coreSchemeHandler) Handle(rawURL *url.URL) core.Result {
	if h.core == nil {
		return core.Result{
			Value: coreerr.E("display.coreSchemeHandler.Handle", "core runtime unavailable", nil),
			OK:    false,
		}
	}

	route, dispatch, result := resolveCoreSchemeRoute(rawURL)
	if !result.OK {
		return result
	}

	target := "core." + route
	params := cloneURLValues(rawURL.Query())
	switch dispatch {
	case routeDispatchAction:
		return h.core.Action(target).Run(context.Background(), optionsFromURLValues(params))
	case routeDispatchQuery:
		if len(params) > 0 {
			result = h.core.Query(CoreRouteQuery{Target: target, Params: params})
			if result.OK {
				return result
			}
		}
		result = h.core.Query(target)
		if result.OK {
			return result
		}
		return core.Result{
			Value: coreerr.E("display.coreSchemeHandler.Handle", "query not handled: "+target, nil),
			OK:    false,
		}
	default:
		return core.Result{
			Value: coreerr.E("display.coreSchemeHandler.Handle", "unsupported dispatch kind", nil),
			OK:    false,
		}
	}
}

func resolveCoreSchemeRoute(rawURL *url.URL) (string, routeDispatchKind, core.Result) {
	if rawURL == nil {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", "scheme URL is required", nil),
			OK:    false,
		}
	}
	if !strings.EqualFold(strings.TrimSpace(rawURL.Scheme), "core") {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", "unsupported scheme: "+rawURL.Scheme, nil),
			OK:    false,
		}
	}
	if strings.TrimSpace(rawURL.Opaque) != "" {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}
	if rawURL.User != nil || strings.TrimSpace(rawURL.Fragment) != "" || rawURL.Port() != "" {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}
	if path := strings.TrimSpace(rawURL.Path); path != "" && path != "/" {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}

	route := strings.ToLower(strings.TrimSpace(rawURL.Hostname()))
	if route == "" {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", malformedCoreURL, nil),
			OK:    false,
		}
	}

	dispatch, ok := coreRouteDispatch[route]
	if !ok {
		return "", routeDispatchQuery, core.Result{
			Value: coreerr.E("display.resolveCoreSchemeRoute", "unknown core route: "+route, nil),
			OK:    false,
		}
	}

	return route, dispatch, core.Result{OK: true}
}

const malformedCoreURL = "malformed core URL"

func optionsFromURLValues(values url.Values) core.Options {
	opts := core.NewOptions()
	for key, items := range sanitizeCoreQuery(values) {
		if len(items) == 1 {
			opts.Set(key, items[0])
			continue
		}
		opts.Set(key, append([]string(nil), items...))
	}
	return opts
}
