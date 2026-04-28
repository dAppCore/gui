package p2p

import (
	"context"
	core "dappco.re/go"
	"errors"
	"time"
)

type fakeDriver struct {
	published    []Envelope
	publishErr   error
	subscribeErr error
}

func (d *fakeDriver) Publish(_ context.Context, envelope Envelope) error {
	d.published = append(d.published, envelope)
	return d.publishErr
}

func (d *fakeDriver) Subscribe(_ context.Context, topic string, handler func(Envelope)) error {
	if d.subscribeErr != nil {
		return d.subscribeErr
	}
	handler(Envelope{
		Topic:    topic,
		Route:    "route",
		SenderID: "peer-1",
		Payload:  map[string]any{"hello": "world"},
	})
	return nil
}

func TestRouter_Publish_Good(t *core.T) {
	driver := &fakeDriver{}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates", SenderID: "peer-1"})
	core.RequireNoError(t, err)
	core.AssertLen(t, driver.published, 1)
	core.AssertEqual(t, "updates", driver.published[0].Topic)
	core.AssertEqual(t, "peer-1", driver.published[0].SenderID)
	core.AssertLessOrEqual(t, time.Since(driver.published[0].ReceivedAt), time.Second)
}

func TestRouter_Publish_Bad(t *core.T) {
	router := New(nil)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	core.RequireNoError(t, err)
}

func TestRouter_Publish_Ugly(t *core.T) {
	driver := &fakeDriver{publishErr: errors.New("publish failed")}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	core.AssertError(t, err)
	core.AssertEqual(t, "publish failed", err.Error())
}

func TestRouter_Subscribe_Good(t *core.T) {
	driver := &fakeDriver{}
	router := New(driver)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(envelope Envelope) {
		calls++
		core.AssertEqual(t, "peer-1", envelope.SenderID)
	})
	core.RequireNoError(t, err)
	core.AssertEqual(t, 1, calls)

	peers := router.Peers()
	core.AssertLen(t, peers, 1)
	core.AssertEqual(t, "peer-1", peers[0].ID)
	core.AssertEqual(t, "timeline", peers[0].Topic)
	core.AssertTrue(t, peers[0].Connected)
}

func TestRouter_Subscribe_Bad(t *core.T) {
	router := New(nil)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {
		calls++
	})
	core.RequireNoError(t, err)
	core.AssertEmpty(t, calls)
	core.AssertEmpty(t, router.Peers())
}

func TestRouter_Subscribe_Ugly(t *core.T) {
	driver := &fakeDriver{subscribeErr: errors.New("subscribe failed")}
	router := New(driver)

	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {})
	core.AssertError(t, err)
	core.AssertEqual(t, "subscribe failed", err.Error())
}
