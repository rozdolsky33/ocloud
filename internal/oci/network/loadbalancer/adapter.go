package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
)

// Adapter implements the domain.LoadBalancerRepository interface for OCI.
type Adapter struct {
	lbClient loadbalancer.LoadBalancerClient
	nwClient core.VirtualNetworkClient
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider common.ConfigurationProvider) (*Adapter, error) {
	lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer client: %w", err)
	}
	nwClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}
	return &Adapter{
		lbClient: lbClient,
		nwClient: nwClient,
	}, nil
}

// GetLoadBalancer retrieves a single Load Balancer and maps it to the domain model (enriched).
func (a *Adapter) GetLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
	response, err := a.lbClient.GetLoadBalancer(ctx, loadbalancer.GetLoadBalancerRequest{
		LoadBalancerId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}

	dm, err := a.enrichAndMapLoadBalancer(ctx, response.LoadBalancer)
	if err != nil {
		return nil, err
	}
	return &dm, nil
}

// ListLoadBalancers returns all load balancers in the compartment (paginated) enriched with details
func (a *Adapter) ListLoadBalancers(ctx context.Context, compartmentID string) ([]domain.LoadBalancer, error) {
	result := make([]domain.LoadBalancer, 0)
	var page *string
	for {
		resp, err := a.lbClient.ListLoadBalancers(ctx, loadbalancer.ListLoadBalancersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing load balancers: %w", err)
		}
		for _, item := range resp.Items {
			mapped, err := a.enrichAndMapLoadBalancer(ctx, item)
			if err != nil {
				return nil, err
			}
			result = append(result, mapped)
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return result, nil
}

// toBaseDomainLoadBalancer maps the base fields without enrichment
func toBaseDomainLoadBalancer(lb loadbalancer.LoadBalancer) domain.LoadBalancer {
	var id, name, shape string
	if lb.Id != nil {
		id = *lb.Id
	}
	if lb.DisplayName != nil {
		name = *lb.DisplayName
	}
	if lb.ShapeName != nil {
		shape = *lb.ShapeName
	}
	// Type: Public/Private
	typeStr := "Public"
	if lb.IsPrivate != nil && *lb.IsPrivate {
		typeStr = "Private"
	}
	var createdTime *time.Time
	if lb.TimeCreated != nil {
		t := lb.TimeCreated.Time
		createdTime = &t
	}

	// IP addresses
	ips := make([]string, 0)
	for _, ip := range lb.IpAddresses {
		label := ""
		if ip.IpAddress != nil {
			label = *ip.IpAddress
		}
		// Try to annotate public/private if available in SDK (best-effort)
		if ip.IsPublic != nil {
			if *ip.IsPublic {
				label = fmt.Sprintf("%s (public)", label)
			} else {
				label = fmt.Sprintf("%s (private)", label)
			}
		}
		if label != "" {
			ips = append(ips, label)
		}
	}

	// Listeners: name -> "proto:port → backendset"
	listeners := make(map[string]string)
	for name, l := range lb.Listeners {
		proto := ""
		if l.Protocol != nil {
			proto = strings.ToLower(*l.Protocol)
		}
		port := 0
		if l.Port != nil {
			port = int(*l.Port)
		}
		backend := ""
		if l.DefaultBackendSetName != nil {
			backend = *l.DefaultBackendSetName
		}
		listeners[name] = fmt.Sprintf("%s:%d → %s", proto, port, backend)
	}

	// Backend sets: only policy and health checker basic info; backends left empty initially
	backendSets := make(map[string]domain.BackendSet)
	for name, bs := range lb.BackendSets {
		policy := ""
		if bs.Policy != nil {
			policy = *bs.Policy
		}
		hc := ""
		if bs.HealthChecker != nil {
			p := ""
			if bs.HealthChecker.Protocol != nil {
				p = strings.ToUpper(*bs.HealthChecker.Protocol)
			}
			port := 0
			if bs.HealthChecker.Port != nil {
				port = int(*bs.HealthChecker.Port)
			}
			hc = fmt.Sprintf("%s:%d", p, port)
		}
		backendSets[name] = domain.BackendSet{Policy: policy, Health: hc, Backends: []domain.Backend{}}
	}

	// Subnets and NSGs (IDs for now; will resolve to names during enrichment)
	subnets := append([]string{}, lb.SubnetIds...)
	nsgs := append([]string{}, lb.NetworkSecurityGroupIds...)

	// Certificates: collect names only (expiry mapping omitted for portability)
	certs := make([]string, 0)
	for name := range lb.Certificates {
		certs = append(certs, name)
	}

	return domain.LoadBalancer{
		ID:              id,
		OCID:            id,
		Name:            name,
		State:           string(lb.LifecycleState),
		Type:            typeStr,
		IPAddresses:     ips,
		Shape:           shape,
		Listeners:       listeners,
		BackendHealth:   map[string]string{},
		Subnets:         subnets,
		NSGs:            nsgs,
		Created:         createdTime,
		BackendSets:     backendSets,
		SSLCertificates: certs,
	}
}

// enrichAndMapLoadBalancer builds the domain model and enriches it with names & health
func (a *Adapter) enrichAndMapLoadBalancer(ctx context.Context, lb loadbalancer.LoadBalancer) (domain.LoadBalancer, error) {
	dm := toBaseDomainLoadBalancer(lb)

	// Resolve Subnet names and CIDRs
	resolvedSubnets := make([]string, 0, len(dm.Subnets))
	for _, sid := range dm.Subnets {
		id := sid
		if id == "" {
			continue
		}
		resp, err := a.nwClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &id})
		if err == nil {
			name := ""
			if resp.Subnet.DisplayName != nil {
				name = *resp.Subnet.DisplayName
			}
			cidr := ""
			if resp.Subnet.CidrBlock != nil {
				cidr = *resp.Subnet.CidrBlock
			}
			if name != "" && cidr != "" {
				resolvedSubnets = append(resolvedSubnets, fmt.Sprintf("%s (%s)", name, cidr))
				continue
			}
		}
		// Fallback to ID when resolution fails
		resolvedSubnets = append(resolvedSubnets, sid)
	}
	dm.Subnets = resolvedSubnets

	// Resolve NSG names
	resolvedNSGs := make([]string, 0, len(dm.NSGs))
	for _, nid := range dm.NSGs {
		id := nid
		if id == "" {
			continue
		}
		resp, err := a.nwClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &id})
		if err == nil {
			name := ""
			if resp.NetworkSecurityGroup.DisplayName != nil {
				name = *resp.NetworkSecurityGroup.DisplayName
			}
			if name != "" {
				resolvedNSGs = append(resolvedNSGs, name)
				continue
			}
		}
		resolvedNSGs = append(resolvedNSGs, nid)
	}
	dm.NSGs = resolvedNSGs

	// Backend set health summaries
	if lb.Id != nil {
		for bsName := range lb.BackendSets {
			name := bsName
			hResp, err := a.lbClient.GetBackendSetHealth(ctx, loadbalancer.GetBackendSetHealthRequest{
				LoadBalancerId: lb.Id,
				BackendSetName: &name,
			})
			if err == nil {
				status := string(hResp.BackendSetHealth.Status)
				if dm.BackendHealth == nil {
					dm.BackendHealth = map[string]string{}
				}
				dm.BackendHealth[name] = status
			}
		}
	}

	return dm, nil
}
