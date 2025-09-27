package loadbalancer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
)

// Adapter implements the domain.LoadBalancerRepository interface for OCI.
type Adapter struct {
	client loadbalancer.LoadBalancerClient
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(provider common.ConfigurationProvider) (*Adapter, error) {
	client, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create load balancer client: %w", err)
	}
	return &Adapter{
		client: client,
	}, nil
}

// GetLoadBalancer retrieves a single Load Balancer and maps it to the domain model.
func (a *Adapter) GetLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
	response, err := a.client.GetLoadBalancer(ctx, loadbalancer.GetLoadBalancerRequest{
		LoadBalancerId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}

	lb := toDomainLoadBalancer(response.LoadBalancer)
	return &lb, nil
}

// ListLoadBalancers returns all load balancers in the compartment (paginated)
func (a *Adapter) ListLoadBalancers(ctx context.Context, compartmentID string) ([]domain.LoadBalancer, error) {
	result := make([]domain.LoadBalancer, 0)
	var page *string
	for {
		resp, err := a.client.ListLoadBalancers(ctx, loadbalancer.ListLoadBalancersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing load balancers: %w", err)
		}
		for _, item := range resp.Items {
			result = append(result, toDomainLoadBalancer(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return result, nil
}

func toDomainLoadBalancer(lb loadbalancer.LoadBalancer) domain.LoadBalancer {
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
		if ip.IpAddress != nil {
			ips = append(ips, *ip.IpAddress)
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

	// Backend sets: only policy and health checker basic info; backends left empty
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

	// Subnets and NSGs
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
