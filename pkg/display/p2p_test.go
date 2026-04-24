package display

import (
	"context"
	"reflect"
	"testing"
	"time"
	"unsafe"

	core "dappco.re/go/core"
	"dappco.re/go/gui/pkg/p2p"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/pkg/application"
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

func replaceP2PDriver(t *testing.T, svc *p2p.Service, driver p2p.Driver) {
	t.Helper()

	router := p2p.New(driver)
	field := reflect.ValueOf(svc).Elem().FieldByName("router")
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(router))
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

	customEvents := &WSEventManager{
		eventBuffer: make(chan Event, 1),
		clients:     make(map[*websocket.Conn]*clientState),
	}
	driver := &immediateP2PDriver{
		envelope: p2p.Envelope{
			Route:    route,
			SenderID: "peer-startup",
			Payload:  map[string]any{"hello": "world"},
		},
	}

	c := core.New(
		core.WithService(func(c *core.Core) core.Result {
			svc, err := New()
			require.NoError(t, err)
			svc.ServiceRuntime = core.NewServiceRuntime(c, Options{})
			svc.wailsApp = &application.App{Logger: application.Logger{}}
			svc.events = customEvents
			return core.Result{Value: svc, OK: true}
		}),
		core.WithService(func(c *core.Core) core.Result {
			p2pSvc := p2p.NewService(c, p2p.Options{NodeID: "node-startup"})
			replaceP2PDriver(t, p2pSvc, driver)
			return core.Result{Value: p2pSvc, OK: true}
		}),
		core.WithServiceLock(),
	)

	require.True(t, c.ServiceStartup(context.Background(), nil).OK)

	select {
	case event := <-customEvents.eventBuffer:
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
