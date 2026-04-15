package display

import (
	"context"
	"testing"
	"time"

	core "dappco.re/go/core"
	"forge.lthn.ai/core/gui/pkg/p2p"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDisplayP2PTestService(t *testing.T) (*Service, *p2p.Service, *core.Core) {
	t.Helper()

	var displaySvc *Service
	var p2pSvc *p2p.Service
	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc, err := New()
			require.NoError(t, err)
			svc.ServiceRuntime = core.NewServiceRuntime(c, Options{})
			displaySvc = svc
			return core.Result{Value: svc, OK: true}
		}),
		core.WithService(func(c *core.Core) core.Result {
			p2pSvc = p2p.NewService(c, p2p.Options{NodeID: "node-1"})
			return core.Result{Value: p2pSvc, OK: true}
		}),
		core.WithServiceLock(),
	)
	require.True(t, c.ServiceStartup(context.Background(), nil).OK)
	require.NotNil(t, displaySvc)
	require.NotNil(t, p2pSvc)
	displaySvc.events = &WSEventManager{
		eventBuffer: make(chan Event, 1),
		clients:     make(map[*websocket.Conn]*clientState),
	}
	displaySvc.attachP2PBridge()
	return displaySvc, p2pSvc, c
}

func TestDisplayP2P_attachP2PBridge_Good(t *testing.T) {
	displaySvc, p2pSvc, _ := newDisplayP2PTestService(t)

	err := p2pSvc.Publish(context.Background(), p2p.Envelope{
		Topic:    "display",
		Route:    "route-1",
		SenderID: "peer-1",
		Payload:  map[string]any{"hello": "world"},
	})
	require.NoError(t, err)

	select {
	case event := <-displaySvc.events.eventBuffer:
		assert.Equal(t, EventCustomEvent, event.Type)
		require.NotNil(t, event.Data)
		assert.Equal(t, "p2p", event.Data["source"])
		assert.Equal(t, "route-1", event.Data["route"])
		assert.Equal(t, "peer-1", event.Data["sender_id"])
		assert.Equal(t, map[string]any{"hello": "world"}, event.Data["payload"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for bridged event")
	}
}

func TestDisplayP2P_attachP2PBridge_Bad(t *testing.T) {
	c := core.New(core.WithServiceLock())
	svc, err := New()
	require.NoError(t, err)
	svc.ServiceRuntime = core.NewServiceRuntime(c, Options{})

	require.NotPanics(t, func() {
		svc.attachP2PBridge()
	})
}

func TestDisplayP2P_attachP2PBridge_Ugly(t *testing.T) {
	displaySvc, p2pSvc, _ := newDisplayP2PTestService(t)
	displaySvc.events = nil

	require.NotPanics(t, func() {
		err := p2pSvc.Publish(context.Background(), p2p.Envelope{
			Topic:    "display",
			SenderID: "peer-2",
		})
		require.NoError(t, err)
	})
}
