package p2p

import (
	"context"
	"strings"

	core "dappco.re/go/core"
)

type Options struct {
	ListenAddr string
	PeerAddrs  []string
	NodeID     string
}

type Service struct {
	*core.ServiceRuntime[Options]
	router *Router
	driver *TCPDriver
}

type State struct {
	NodeID     string `json:"node_id"`
	ListenAddr string `json:"listen_addr,omitempty"`
	Peers      []Peer `json:"peers,omitempty"`
}

func NewService(c *core.Core, options Options) *Service {
	driver := NewTCPDriver(TCPOptions{
		ListenAddr: options.ListenAddr,
		PeerAddrs:  options.PeerAddrs,
		NodeID:     options.NodeID,
	})
	return &Service{
		ServiceRuntime: core.NewServiceRuntime(c, options),
		router:         New(driver),
		driver:         driver,
	}
}

func OptionsFromEnv() Options {
	peers := strings.Split(strings.TrimSpace(core.Env("CORE_P2P_PEERS")), ",")
	filtered := make([]string, 0, len(peers))
	for _, peer := range peers {
		peer = strings.TrimSpace(peer)
		if peer != "" {
			filtered = append(filtered, peer)
		}
	}
	return Options{
		ListenAddr: strings.TrimSpace(core.Env("CORE_P2P_ADDR")),
		PeerAddrs:  filtered,
		NodeID:     strings.TrimSpace(core.Env("CORE_P2P_NODE_ID")),
	}
}

func (s *Service) OnStartup(_ context.Context) core.Result {
	s.Core().Action("p2p.publish", func(ctx context.Context, opts core.Options) core.Result {
		payload := mapValue(opts, "payload")
		envelope := Envelope{
			Topic:    opts.String("topic"),
			Route:    opts.String("route"),
			SenderID: coalesce(opts.String("sender_id"), s.Options().NodeID),
			Payload:  payload,
		}
		return core.Result{}.New(nil, s.Publish(ctx, envelope))
	})
	s.Core().Action("p2p.state", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: s.State(), OK: true}
	})
	return core.Result{OK: true}
}

func (s *Service) OnShutdown(_ context.Context) core.Result {
	return core.Result{}.New(nil, s.driver.Close())
}

func (s *Service) Publish(ctx context.Context, envelope Envelope) error {
	return s.router.Publish(ctx, envelope)
}

func (s *Service) Subscribe(ctx context.Context, topic string, handler func(Envelope)) error {
	return s.router.Subscribe(ctx, topic, handler)
}

func (s *Service) Peers() []Peer {
	return s.router.Peers()
}

func (s *Service) State() State {
	return State{
		NodeID:     s.Options().NodeID,
		ListenAddr: s.driver.ListenAddr(),
		Peers:      s.Peers(),
	}
}

func mapValue(opts core.Options, key string) map[string]any {
	result := opts.Get(key)
	if !result.OK {
		return nil
	}
	value := result.Value
	switch typed := value.(type) {
	case map[string]any:
		return typed
	default:
		var normalized map[string]any
		if result := core.JSONUnmarshalString(core.JSONMarshalString(typed), &normalized); result.OK {
			return normalized
		}
		return map[string]any{"value": typed}
	}
}

func coalesce(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
