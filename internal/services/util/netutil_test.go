package util

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExtractHostname(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"example.com", "example.com"},
		{"example.com:443", "example.com"},
		{"http://example.com", "example.com"},
		{"https://example.com:8443", "example.com"},
		{"https://example.com:8443/path/to/res", "example.com"},
		{"ssh://host.internal:22/whatever", "host.internal"},
		{"host.internal/path", "host.internal"},
	}
	for _, c := range cases {
		got := ExtractHostname(c.in)
		require.Equal(t, c.out, got, "input=%s", c.in)
	}
}

func TestResolveHostToIP_Localhost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ipStr, err := ResolveHostToIP(ctx, "localhost")
	require.NoError(t, err)
	require.NotEmpty(t, ipStr)

	ip := net.ParseIP(ipStr)
	require.NotNil(t, ip, "should parse as IP: %s", ipStr)
	require.True(t, ip.IsLoopback(), "expected loopback IP, got %s", ip)
}

func TestResolveHostToIP_Error(t *testing.T) {
	// RFC 2606/6761 reserved TLD .invalid should never resolve.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := ResolveHostToIP(ctx, "definitely-not-existing-domain.invalid")
	require.Error(t, err)
}

func TestIsLocalTCPPortInUse(t *testing.T) {
	// Acquire a free port by listening on :0
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port

	// While listening, the port should be reported as in use
	require.True(t, IsLocalTCPPortInUse(port))

	// Close and allow a short grace period for the OS to release the port
	require.NoError(t, ln.Close())
	time.Sleep(50 * time.Millisecond)

	require.False(t, IsLocalTCPPortInUse(port))

	// Also test a separately freed port
	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	p2 := ln2.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln2.Close())
	time.Sleep(20 * time.Millisecond)
	require.False(t, IsLocalTCPPortInUse(p2))
}
