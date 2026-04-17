package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTCPDriver_Publish_ContinuesAfterPeerFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	received := make(chan Envelope, 1)
	acceptErr := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptErr <- err
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		if scanner.Scan() {
			var envelope Envelope
			if err := json.Unmarshal(scanner.Bytes(), &envelope); err != nil {
				acceptErr <- err
				return
			}
			received <- envelope
			return
		}
		if err := scanner.Err(); err != nil {
			acceptErr <- err
			return
		}
		acceptErr <- context.Canceled
	}()

	driver := NewTCPDriver(TCPOptions{
		PeerAddrs: []string{"127.0.0.1:1", listener.Addr().String()},
		NodeID:    "node-1",
	})

	err = driver.Publish(context.Background(), Envelope{
		Topic:   "updates",
		Payload: map[string]any{"hello": "world"},
	})
	require.Error(t, err)

	select {
	case envelope := <-received:
		assert.Equal(t, "updates", envelope.Topic)
		assert.Equal(t, "node-1", envelope.SenderID)
		assert.Equal(t, map[string]any{"hello": "world"}, envelope.Payload)
	case err := <-acceptErr:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for peer delivery")
	}
}
