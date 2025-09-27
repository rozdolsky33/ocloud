package loadbalancer

import (
	"context"
	x509std "crypto/x509"
	pemenc "encoding/pem"
	"fmt"
	"strings"
	"sync"
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

		items := resp.Items
		mappedCh := make(chan domain.LoadBalancer, len(items))
		errCh := make(chan error, len(items))
		var wg sync.WaitGroup
		for i := range items {
			wg.Add(1)
			go func(it loadbalancer.LoadBalancer) {
				defer wg.Done()
				mapped, err := a.enrichAndMapLoadBalancer(ctx, it)
				if err != nil {
					errCh <- err
					return
				}
				mappedCh <- mapped
			}(items[i])
		}
		wg.Wait()
		close(errCh)
		for e := range errCh {
			if e != nil {
				return nil, e
			}
		}
		close(mappedCh)
		for m := range mappedCh {
			result = append(result, m)
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

// enrichAndMapLoadBalancer builds the domain model and enriches it with names, health, and SSL certificate info using concurrency
func (a *Adapter) enrichAndMapLoadBalancer(ctx context.Context, lb loadbalancer.LoadBalancer) (domain.LoadBalancer, error) {
	dm := toBaseDomainLoadBalancer(lb)

	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	var mu sync.Mutex

	// Resolve Subnet names and CIDRs
	wg.Add(1)
	go func() {
		defer wg.Done()
		resolved := make([]string, 0, len(dm.Subnets))
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
					resolved = append(resolved, fmt.Sprintf("%s (%s)", name, cidr))
					continue
				}
			}
			// Fallback to ID when resolution fails
			resolved = append(resolved, sid)
		}
		mu.Lock()
		dm.Subnets = resolved
		mu.Unlock()
	}()

	// Resolve NSG names
	wg.Add(1)
	go func() {
		defer wg.Done()
		resolved := make([]string, 0, len(dm.NSGs))
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
					resolved = append(resolved, name)
					continue
				}
			}
			resolved = append(resolved, nid)
		}
		mu.Lock()
		dm.NSGs = resolved
		mu.Unlock()
	}()

	// Backend set health summaries (concurrent per backend set)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if lb.Id == nil {
			return
		}
		var inner sync.WaitGroup
		healthLocal := make(map[string]string)
		var hmu sync.Mutex
		for bsName := range lb.BackendSets {
			name := bsName
			inner.Add(1)
			go func(n string) {
				defer inner.Done()
				hResp, err := a.lbClient.GetBackendSetHealth(ctx, loadbalancer.GetBackendSetHealthRequest{
					LoadBalancerId: lb.Id,
					BackendSetName: &n,
				})
				if err != nil {
					return
				}
				status := string(hResp.BackendSetHealth.Status)
				hmu.Lock()
				healthLocal[n] = status
				hmu.Unlock()
			}(name)
		}
		inner.Wait()
		mu.Lock()
		if dm.BackendHealth == nil {
			dm.BackendHealth = map[string]string{}
		}
		for k, v := range healthLocal {
			dm.BackendHealth[k] = v
		}
		mu.Unlock()
	}()

	// SSL Certificates: parse PEM and extract expiry date
	wg.Add(1)
	go func() {
		defer wg.Done()
		certs := make([]string, 0, len(lb.Certificates))
		for name, cert := range lb.Certificates {
			expires := ""
			if cert.PublicCertificate != nil && *cert.PublicCertificate != "" {
				if t, ok := parseCertNotAfter(*cert.PublicCertificate); ok {
					expires = t.Format("2006-01-02")
				}
			}
			if expires != "" {
				certs = append(certs, fmt.Sprintf("%s (Expires: %s)", name, expires))
			} else {
				certs = append(certs, name)
			}
		}
		mu.Lock()
		dm.SSLCertificates = certs
		mu.Unlock()
	}()

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return dm, err
		}
	}

	return dm, nil
}

// parseCertNotAfter attempts to parse the first certificate in a PEM bundle and returns NotAfter
func parseCertNotAfter(pemData string) (time.Time, bool) {
	data := []byte(pemData)
	for {
		var block *pemenc.Block
		block, data = pemenc.Decode(data)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509std.ParseCertificate(block.Bytes)
			if err == nil {
				return c.NotAfter, true
			}
		}
	}
	return time.Time{}, false
}
