package bastion

import (
	"fmt"
	"net"
	"strings"
)

// extractHostname extracts the hostname from a URL or endpoint string
func extractHostname(endpoint string) string {
	if endpoint == "" {
		return ""
	}

	// Remove protocol prefix if present
	if strings.Contains(endpoint, "://") {
		parts := strings.Split(endpoint, "://")
		if len(parts) > 1 {
			endpoint = parts[1]
		}
	}

	// Extract hostname (remove port and path if present)
	hostname := endpoint
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}
	if strings.Contains(hostname, "/") {
		hostname = strings.Split(hostname, "/")[0]
	}

	return hostname
}

// resolveHostToIP resolves a hostname to its IP address
func resolveHostToIP(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to resolve hostname %s: %w", hostname, err)
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for hostname %s", hostname)
	}

	// Return the first IP address
	return ips[0].String(), nil
}
