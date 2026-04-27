package p2p

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestRouter_Publish_Good(t *testing.T) {
	driver := &fakeDriver{}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates", SenderID: "peer-1"})
	require.NoError(t, err)
	require.Len(t, driver.published, 1)
	assert.Equal(t, "updates", driver.published[0].Topic)
	assert.Equal(t, "peer-1", driver.published[0].SenderID)
	assert.WithinDuration(t, time.Now(), driver.published[0].ReceivedAt, time.Second)
}

func TestRouter_Publish_Bad(t *testing.T) {
	router := New(nil)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	require.NoError(t, err)
}

func TestRouter_Publish_Ugly(t *testing.T) {
	driver := &fakeDriver{publishErr: errors.New("publish failed")}
	router := New(driver)

	err := router.Publish(context.Background(), Envelope{Topic: "updates"})
	require.Error(t, err)
	assert.Equal(t, "publish failed", err.Error())
}

func TestRouter_Subscribe_Good(t *testing.T) {
	driver := &fakeDriver{}
	router := New(driver)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(envelope Envelope) {
		calls++
		assert.Equal(t, "peer-1", envelope.SenderID)
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)

	peers := router.Peers()
	require.Len(t, peers, 1)
	assert.Equal(t, "peer-1", peers[0].ID)
	assert.Equal(t, "timeline", peers[0].Topic)
	assert.True(t, peers[0].Connected)
}

func TestRouter_Subscribe_Bad(t *testing.T) {
	router := New(nil)

	calls := 0
	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {
		calls++
	})
	require.NoError(t, err)
	assert.Zero(t, calls)
	assert.Empty(t, router.Peers())
}

func TestRouter_Subscribe_Ugly(t *testing.T) {
	driver := &fakeDriver{subscribeErr: errors.New("subscribe failed")}
	router := New(driver)

	err := router.Subscribe(context.Background(), "timeline", func(Envelope) {})
	require.Error(t, err)
	assert.Equal(t, "subscribe failed", err.Error())
}
