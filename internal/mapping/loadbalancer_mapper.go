package mapping

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
)

type LoadBalancerAttributes struct {
	ID                      *string
	DisplayName             *string
	LifecycleState          loadbalancer.LoadBalancerLifecycleStateEnum
	IsPrivate               *bool
	ShapeName               *string
	TimeCreated             *time.Time
	IpAddresses             []loadbalancer.IpAddress
	Listeners               map[string]loadbalancer.Listener
	RoutingPolicies         map[string]loadbalancer.RoutingPolicy
	BackendSets             map[string]loadbalancer.BackendSet
	SubnetIds               []string
	NetworkSecurityGroupIds []string
	Certificates            map[string]loadbalancer.Certificate
	Hostnames               map[string]loadbalancer.Hostname
}

func NewLoadBalancerAttributesFromOCILoadBalancer(lb loadbalancer.LoadBalancer) *LoadBalancerAttributes {
	return &LoadBalancerAttributes{
		ID:                      lb.Id,
		DisplayName:             lb.DisplayName,
		LifecycleState:          lb.LifecycleState,
		IsPrivate:               lb.IsPrivate,
		ShapeName:               lb.ShapeName,
		TimeCreated:             &lb.TimeCreated.Time,
		IpAddresses:             lb.IpAddresses,
		Listeners:               lb.Listeners,
		RoutingPolicies:         lb.RoutingPolicies,
		BackendSets:             lb.BackendSets,
		SubnetIds:               lb.SubnetIds,
		NetworkSecurityGroupIds: lb.NetworkSecurityGroupIds,
		Certificates:            lb.Certificates,
		Hostnames:               lb.Hostnames,
	}
}

func NewDomainLoadBalancerFromAttrs(lb *LoadBalancerAttributes) *domain.LoadBalancer {
	deref := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}
	derefInt := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	id := deref(lb.ID)
	name := deref(lb.DisplayName)
	shape := deref(lb.ShapeName)

	// Type
	typeStr := "Public"
	if lb.IsPrivate != nil && *lb.IsPrivate {
		typeStr = "Private"
	}

	var createdTime *time.Time
	if lb.TimeCreated != nil {
		t := *lb.TimeCreated
		createdTime = &t
	}

	ips := make([]string, 0, len(lb.IpAddresses))
	for i := range lb.IpAddresses {
		ip := lb.IpAddresses[i]
		addr := deref(ip.IpAddress)
		if addr == "" {
			continue
		}

		if ip.IsPublic != nil {
			if *ip.IsPublic {
				addr += " (public)"
			} else {
				addr += " (private)"
			}
		}
		ips = append(ips, addr)
	}

	// Listeners (name -> "proto:port -> backendset")
	listeners := make(map[string]string, len(lb.Listeners))
	useSSL := false

	// routing policies referenced by listeners; fallback to LB map if none seen
	routingPolicySet := make(map[string]struct{}, len(lb.RoutingPolicies))
	for lname, l := range lb.Listeners {
		port := derefInt(l.Port)
		backend := deref(l.DefaultBackendSetName)

		if l.SslConfiguration != nil {
			useSSL = true
		}
		if rp := deref(l.RoutingPolicyName); rp != "" {
			routingPolicySet[rp] = struct{}{}
		}

		// Cheap protocol detection:
		// If SSL config exists or typical HTTPS ports → https
		// Else if the protocol string is "TCP" (case-insensitive) → tcp
		// Else common HTTP ports → http
		proto := "http"
		if l.SslConfiguration != nil || port == 443 || port == 8443 {
			proto = "https"
		} else if l.Protocol != nil {
			p := *l.Protocol
			if p == "TCP" || p == "tcp" || p == "Tcp" {
				proto = "tcp"
			} else if port == 80 {
				proto = "http"
			}
		} else if port == 80 {
			proto = "http"
		}

		// Build "proto:port → backend" without fmt
		// (fmt is fine, but concat avoids an allocation)
		val := proto + ":" + strconv.Itoa(port) + " → " + backend
		listeners[lname] = val
	}

	if len(routingPolicySet) == 0 {
		for rpName := range lb.RoutingPolicies {
			if rpName != "" {
				routingPolicySet[rpName] = struct{}{}
			}
		}
	}
	routingPolicies := make([]string, 0, len(routingPolicySet))
	for rp := range routingPolicySet {
		routingPolicies = append(routingPolicies, rp)
	}
	sort.Strings(routingPolicies)

	// Backend sets (policy + health-check summary)
	backendSets := make(map[string]domain.BackendSet, len(lb.BackendSets))
	for bsName, bs := range lb.BackendSets {
		policy := deref(bs.Policy)

		hc := ""
		if bs.HealthChecker != nil {
			// Derive health-check label: PROTO:PORT (prefer explicit; normalize common ports)
			var p string
			if bs.HealthChecker.Protocol != nil {
				// normalize once; no need for ToUpper on every path
				switch *bs.HealthChecker.Protocol {
				case "https", "HTTPS", "Https":
					p = "HTTPS"
				case "http", "HTTP", "Http":
					p = "HTTP"
				case "tcp", "TCP", "Tcp":
					p = "TCP"
				default:
					p = strings.ToUpper(*bs.HealthChecker.Protocol)
				}
			}

			port := 0
			if bs.HealthChecker.Port != nil {
				port = int(*bs.HealthChecker.Port)
			}
			// If port hints common protocols, prefer those (matches your earlier behavior)
			switch port {
			case 443, 8443:
				p = "HTTPS"
			case 80:
				if p == "" {
					p = "HTTP"
				}
			}
			if p != "" {
				hc = p + ":" + strconv.Itoa(port)
			}
		}

		backendSets[bsName] = domain.BackendSet{
			Policy:   policy,
			Health:   hc,
			Backends: []domain.Backend{}, // filled during enrichment
		}
	}

	// Subnets/NSGs – copy IDs (pre-size)
	subnets := make([]string, len(lb.SubnetIds))
	copy(subnets, lb.SubnetIds)

	nsgs := make([]string, len(lb.NetworkSecurityGroupIds))
	copy(nsgs, lb.NetworkSecurityGroupIds)

	// Certificates: collect names only (pre-size)
	certs := make([]string, 0, len(lb.Certificates))
	for cname := range lb.Certificates {
		certs = append(certs, cname)
	}

	// Hostnames (prefer value, fallback to key) + sort
	hostnames := make([]string, 0, len(lb.Hostnames))
	for key, h := range lb.Hostnames {
		if h.Hostname != nil && *h.Hostname != "" {
			hostnames = append(hostnames, *h.Hostname)
		} else if s := strings.TrimSpace(key); s != "" {
			hostnames = append(hostnames, s)
		}
	}
	sort.Strings(hostnames)

	return &domain.LoadBalancer{
		ID:              id,
		OCID:            id,
		Name:            name,
		State:           string(lb.LifecycleState),
		Type:            typeStr,
		IPAddresses:     ips,
		Shape:           shape,
		Listeners:       listeners,
		BackendHealth:   make(map[string]string),
		Subnets:         subnets,
		NSGs:            nsgs,
		Created:         createdTime,
		BackendSets:     backendSets,
		SSLCertificates: certs,
		RoutingPolicies: routingPolicies,
		UseSSL:          useSSL,
		Hostnames:       hostnames,
	}
}
