package display

import (
	"context"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/p2p"
)

func (s *Service) attachP2PBridge() {
	router, ok := core.ServiceFor[*p2p.Service](s.Core(), "p2p")
	if !ok || router == nil {
		return
	}
	_ = router.Subscribe(context.Background(), "display", func(envelope p2p.Envelope) {
		if s.events == nil {
			return
		}
		s.events.Emit(Event{
			Type: EventCustomEvent,
			Data: map[string]any{
				"source":    "p2p",
				"topic":     envelope.Topic,
				"route":     envelope.Route,
				"sender_id": envelope.SenderID,
				"payload":   envelope.Payload,
			},
		})
	})
}
