package display

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetwork_InterfaceFlags_Good(t *testing.T) {
	flags := interfaceFlags(net.FlagUp | net.FlagLoopback | net.FlagRunning)

	assert.Equal(t, []string{"up", "loopback", "running"}, flags)
}

func TestNetwork_InterfaceFlags_Bad(t *testing.T) {
	assert.Empty(t, interfaceFlags(0))
}

func TestNetwork_InterfaceFlags_Ugly(t *testing.T) {
	assert.Empty(t, interfaceFlags(net.Flags(1<<30)))
}

func TestNetwork_RenderNetworkPage_Good(t *testing.T) {
	svc := &Service{}
	state := NetworkState{
		Hostname:   "core-host",
		ObservedAt: time.Unix(1_700_000_000, 0).UTC(),
		Interfaces: []NetworkInterfaceState{
			{
				Name:      "en0",
				Index:     2,
				MTU:       1500,
				Addresses: []string{"192.168.0.10/24", "fe80::1/64"},
				Up:        true,
			},
		},
		Peers: []NetworkPeerState{
			{ID: "peer-1", Topic: "timeline", Connected: true, SeenAt: time.Unix(1_700_000_100, 0).UTC()},
		},
	}

	body := svc.renderNetworkPage(state)

	assert.Contains(t, body, "core://network")
	assert.Contains(t, body, "core-host")
	assert.Contains(t, body, "en0")
	assert.Contains(t, body, "192.168.0.10/24")
	assert.Contains(t, body, "Registered peers")
	assert.Contains(t, body, "peer-1")
}

func TestNetwork_RenderNetworkPage_Bad(t *testing.T) {
	svc := &Service{}

	body := svc.renderNetworkPage(NetworkState{
		Hostname:   "<host>",
		ObservedAt: time.Unix(1, 0).UTC(),
	})

	assert.Contains(t, body, "No network interfaces were detected.")
	assert.Contains(t, body, "&lt;host&gt;")
}

func TestNetwork_RenderNetworkPage_Ugly(t *testing.T) {
	svc := &Service{}

	body := svc.renderNetworkPage(NetworkState{
		Hostname:   strings.Repeat("x", 128),
		ObservedAt: time.Unix(1, 0).UTC(),
		Interfaces: []NetworkInterfaceState{
			{Name: "\"quoted\"", Index: 99, MTU: 9, Addresses: []string{"<addr>"}, Up: false, Loopback: true},
		},
	})

	assert.Contains(t, body, "&#34;quoted&#34;")
	assert.Contains(t, body, "&lt;addr&gt;")
	assert.Contains(t, body, "loopback")
}

func TestNetwork_RenderNetworkInterfacePage_Good(t *testing.T) {
	svc := &Service{}
	state := NetworkState{
		Hostname:   "core-host",
		ObservedAt: time.Unix(1_700_000_000, 0).UTC(),
		Peers: []NetworkPeerState{
			{ID: "peer-1", Topic: "timeline", Connected: true, SeenAt: time.Unix(1_700_000_100, 0).UTC()},
		},
	}
	iface := NetworkInterfaceState{
		Name:      "en0",
		Index:     2,
		MTU:       1500,
		Addresses: []string{"192.168.0.10/24"},
		Flags:     []string{"up", "running"},
		Up:        true,
	}

	body := svc.renderNetworkInterfacePage(state, iface)

	assert.Contains(t, body, "core://network/en0")
	assert.Contains(t, body, "core-host")
	assert.Contains(t, body, "192.168.0.10/24")
	assert.Contains(t, body, "Flags: up, running")
	assert.Contains(t, body, "Registered peers")
	assert.Contains(t, body, "peer-1")
}

func TestNetwork_RenderNetworkInterfacePage_Bad(t *testing.T) {
	svc := &Service{}
	state := NetworkState{Hostname: "core-host", ObservedAt: time.Unix(1, 0).UTC()}
	iface := NetworkInterfaceState{Name: "en0", Index: 2, MTU: 1500, Up: false}

	body := svc.renderNetworkInterfacePage(state, iface)

	assert.Contains(t, body, "Flags: none")
	assert.NotContains(t, body, "Registered peers")
}

func TestNetwork_RenderNetworkInterfacePage_Ugly(t *testing.T) {
	svc := &Service{}
	state := NetworkState{Hostname: "<host>", ObservedAt: time.Unix(1, 0).UTC()}
	iface := NetworkInterfaceState{
		Name:      "\"quoted\"",
		Index:     99,
		MTU:       9,
		Addresses: []string{"<addr>"},
		Up:        false,
		Loopback:  true,
	}

	body := svc.renderNetworkInterfacePage(state, iface)

	assert.Contains(t, body, "&#34;quoted&#34;")
	assert.Contains(t, body, "&lt;addr&gt;")
	assert.Contains(t, body, "&lt;host&gt;")
	assert.Contains(t, body, "loopback")
}
