package p2p

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"sync"
)

type TCPOptions struct {
	ListenAddr string
	PeerAddrs  []string
	NodeID     string
}

type TCPDriver struct {
	options       TCPOptions
	mu            sync.RWMutex
	listener      net.Listener
	subscriptions map[string][]func(Envelope)
}

func NewTCPDriver(options TCPOptions) *TCPDriver {
	return &TCPDriver{
		options: TCPOptions{
			ListenAddr: strings.TrimSpace(options.ListenAddr),
			PeerAddrs:  append([]string(nil), options.PeerAddrs...),
			NodeID:     strings.TrimSpace(options.NodeID),
		},
		subscriptions: make(map[string][]func(Envelope)),
	}
}

func (d *TCPDriver) ListenAddr() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.listener != nil {
		return d.listener.Addr().String()
	}
	return d.options.ListenAddr
}

func (d *TCPDriver) Subscribe(_ context.Context, topic string, handler func(Envelope)) error {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return errors.New("topic is required")
	}
	if handler == nil {
		return errors.New("handler is required")
	}
	d.mu.Lock()
	d.subscriptions[topic] = append(d.subscriptions[topic], handler)
	d.mu.Unlock()
	return d.ensureListener()
}

func (d *TCPDriver) Publish(ctx context.Context, envelope Envelope) error {
	if strings.TrimSpace(envelope.Topic) == "" {
		return errors.New("topic is required")
	}
	if strings.TrimSpace(envelope.SenderID) == "" {
		envelope.SenderID = d.options.NodeID
	}
	d.dispatch(envelope)
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	for _, peer := range d.options.PeerAddrs {
		peer = strings.TrimSpace(peer)
		if peer == "" {
			continue
		}
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", peer)
		if err != nil {
			return err
		}
		if _, err := conn.Write(append(payload, '\n')); err != nil {
			_ = conn.Close()
			return err
		}
		_ = conn.Close()
	}
	return nil
}

func (d *TCPDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.listener == nil {
		return nil
	}
	err := d.listener.Close()
	d.listener = nil
	return err
}

func (d *TCPDriver) ensureListener() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.listener != nil || strings.TrimSpace(d.options.ListenAddr) == "" {
		return nil
	}
	listener, err := net.Listen("tcp", d.options.ListenAddr)
	if err != nil {
		return err
	}
	d.listener = listener
	go d.acceptLoop(listener)
	return nil
}

func (d *TCPDriver) acceptLoop(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go d.readConn(conn)
	}
}

func (d *TCPDriver) readConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var envelope Envelope
		if err := json.Unmarshal(scanner.Bytes(), &envelope); err != nil {
			continue
		}
		d.dispatch(envelope)
	}
}

func (d *TCPDriver) dispatch(envelope Envelope) {
	d.mu.RLock()
	handlers := append([]func(Envelope){}, d.subscriptions[envelope.Topic]...)
	handlers = append(handlers, d.subscriptions["*"]...)
	d.mu.RUnlock()
	for _, handler := range handlers {
		handler(envelope)
	}
}
