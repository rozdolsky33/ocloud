// Package bastion Network-related utilities (hostname extraction and ctx-aware DNS).
package bastion

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// extractHostname removes schema/port/path and returns just the host portion.
func extractHostname(endpoint string) string {
	if endpoint == "" {
		return ""
	}
	if strings.Contains(endpoint, "://") {
		parts := strings.SplitN(endpoint, "://", 2)
		endpoint = parts[1]
	}
	host := endpoint
	if i := strings.IndexByte(host, '/'); i >= 0 {
		host = host[:i]
	}
	if i := strings.IndexByte(host, ':'); i >= 0 {
		host = host[:i]
	}
	return host
}

// resolveHostToIP resolves hostname to the first IP (IPv4/IPv6). It uses ctx so
// cancellation/timeouts propagate.
func resolveHostToIP(ctx context.Context, hostname string) (string, error) {
	var r net.Resolver
	ips, err := r.LookupIP(ctx, "ip", hostname)
	if err != nil {
		return "", fmt.Errorf("resolve hostname %s: %w", hostname, err)
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IPs found for hostname %s", hostname)
	}
	return ips[0].String(), nil
}
