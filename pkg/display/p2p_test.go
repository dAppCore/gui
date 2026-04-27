package display

import (
	"context"
	"sync"
	"testing"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/p2p"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDisplayP2PTestService(t *testing.T) (*Service, *p2p.Service, *core.Core) {
	t.Helper()

	var p2pSvc *p2p.Service
	driver := newLoopbackP2PDriver()
	displaySvc, c := newServiceWithMockApp(t, func(c *core.Core) core.Result {
		p2pSvc = p2p.NewServiceWithDriver(c, p2p.Options{NodeID: "node-1"}, driver)
		return core.Result{Value: p2pSvc, OK: true}
	})
	require.NotNil(t, displaySvc)
	require.NotNil(t, p2pSvc)
	return displaySvc, p2pSvc, c
}

type loopbackP2PDriver struct {
	mu       sync.Mutex
	handlers map[string]func(p2p.Envelope)
}

func newLoopbackP2PDriver() *loopbackP2PDriver {
	return &loopbackP2PDriver{handlers: make(map[string]func(p2p.Envelope))}
}

func (d *loopbackP2PDriver) Publish(_ context.Context, envelope p2p.Envelope) error {
	d.mu.Lock()
	handler := d.handlers[envelope.Topic]
	d.mu.Unlock()
	if handler != nil {
		handler(envelope)
	}
	return nil
}

func (d *loopbackP2PDriver) Subscribe(_ context.Context, topic string, handler func(p2p.Envelope)) error {
	d.mu.Lock()
	d.handlers[topic] = handler
	d.mu.Unlock()
	return nil
}

type immediateP2PDriver struct {
	envelope p2p.Envelope
}

func (d *immediateP2PDriver) Publish(context.Context, p2p.Envelope) error {
	return nil
}

func (d *immediateP2PDriver) Subscribe(_ context.Context, topic string, handler func(p2p.Envelope)) error {
	if handler != nil {
		handler(p2p.Envelope{
			Topic:    topic,
			Route:    d.envelope.Route,
			SenderID: d.envelope.SenderID,
			Payload:  d.envelope.Payload,
		})
	}
	return nil
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

func TestDisplayP2P_OnStartup_InitializesEventManagerBeforeBridge(t *testing.T) {
	const route = "route-startup"

	driver := &immediateP2PDriver{
		envelope: p2p.Envelope{
			Route:    route,
			SenderID: "peer-startup",
			Payload:  map[string]any{"hello": "world"},
		},
	}

	displaySvc, _ := newServiceWithMockApp(t, func(c *core.Core) core.Result {
		p2pSvc := p2p.NewServiceWithDriver(c, p2p.Options{NodeID: "node-startup"}, driver)
		return core.Result{Value: p2pSvc, OK: true}
	})

	select {
	case event := <-displaySvc.events.eventBuffer:
		assert.Equal(t, EventCustomEvent, event.Type)
		require.NotNil(t, event.Data)
		assert.Equal(t, "p2p", event.Data["source"])
		assert.Equal(t, route, event.Data["route"])
		assert.Equal(t, "peer-startup", event.Data["sender_id"])
		assert.Equal(t, map[string]any{"hello": "world"}, event.Data["payload"])
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for bridged startup event")
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
