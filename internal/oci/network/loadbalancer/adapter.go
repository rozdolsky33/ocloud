package loadbalancer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	lbLogger "github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"
)

// Adapter implements the domain.LoadBalancerRepository interface for OCI.
type Adapter struct {
	lbClient    loadbalancer.LoadBalancerClient
	nwClient    core.VirtualNetworkClient
	certsClient certificatesmanagement.CertificatesManagementClient
	limiter     *rate.Limiter
	sf          singleflight.Group
	workerCount int
	// caches to reduce repeated OCI calls within a command run
	subnetCache   map[string]core.GetSubnetResponse
	vcnCache      map[string]core.GetVcnResponse
	nsgCache      map[string]core.GetNetworkSecurityGroupResponse
	certListCache map[string][]loadbalancer.Certificate // keyed by LB ID
	muSubnets     sync.RWMutex
	muVcns        sync.RWMutex
	muNsgs        sync.RWMutex
	muCertLists   sync.RWMutex
}

// NewAdapter creates a new Adapter instance using pre-created OCI clients.
func NewAdapter(lbClient loadbalancer.LoadBalancerClient, nwClient core.VirtualNetworkClient, certsClient certificatesmanagement.CertificatesManagementClient) *Adapter {
	ad := &Adapter{
		lbClient:      lbClient,
		nwClient:      nwClient,
		certsClient:   certsClient,
		workerCount:   defaultWorkerCount,
		limiter:       rate.NewLimiter(rate.Limit(defaultRatePerSec), defaultRateBurst),
		subnetCache:   make(map[string]core.GetSubnetResponse),
		vcnCache:      make(map[string]core.GetVcnResponse),
		nsgCache:      make(map[string]core.GetNetworkSecurityGroupResponse),
		certListCache: make(map[string][]loadbalancer.Certificate),
	}
	return ad
}

// GetLoadBalancer retrieves a single Load Balancer and maps it to the basic domain model, adding backend health for usability.
func (a *Adapter) GetLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
	response, err := a.lbClient.GetLoadBalancer(ctx, loadbalancer.GetLoadBalancerRequest{
		LoadBalancerId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	dm := mapping.NewDomainLoadBalancerFromAttrs(mapping.NewLoadBalancerAttributesFromOCILoadBalancer(response.LoadBalancer))
	_ = a.enrichBackendHealth(ctx, response.LoadBalancer, dm, false)
	_ = a.resolveSubnets(ctx, dm)
	return dm, nil
}

// GetEnrichedLoadBalancer retrieves a single Load Balancer and returns the enriched domain model.
func (a *Adapter) GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
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
	return dm, nil
}

// ListLoadBalancers returns all load balancers in the compartment (paginated) mapped to a base domain model,
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
		pageItems := resp.Items
		mapped := make([]domain.LoadBalancer, len(pageItems))
		var wg sync.WaitGroup
		for i := range pageItems {
			wg.Add(1)
			idx := i
			go func() {
				defer wg.Done()
				lb := pageItems[idx]
				dm := mapping.NewDomainLoadBalancerFromAttrs(mapping.NewLoadBalancerAttributesFromOCILoadBalancer(lb))
				_ = a.enrichBackendHealth(ctx, lb, dm, false)
				_ = a.resolveSubnets(ctx, dm)
				mapped[idx] = *dm
			}()
		}
		wg.Wait()
		result = append(result, mapped...)

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return result, nil
}

// ListEnrichedLoadBalancers returns all load balancers in the compartment (paginated) with enrichment
func (a *Adapter) ListEnrichedLoadBalancers(ctx context.Context, compartmentID string) ([]domain.LoadBalancer, error) {
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

		mapped, err := a.enrichAndMapLoadBalancers(ctx, resp.Items)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped...)

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return result, nil
}

// enrichAndMapLoadBalancers converts a slice of OCI LBs to domain models with enrichment using concurrency
func (a *Adapter) enrichAndMapLoadBalancers(ctx context.Context, items []loadbalancer.LoadBalancer) ([]domain.LoadBalancer, error) {
	// Page-level prefetch: collect unique Subnet and NSG IDs across items and resolve them once using the worker pool.
	uniqSubnets := make(map[string]struct{})
	uniqNSGs := make(map[string]struct{})
	for _, lb := range items {
		for _, sid := range lb.SubnetIds {
			if sid != "" {
				uniqSubnets[sid] = struct{}{}
			}
		}
		for _, nid := range lb.NetworkSecurityGroupIds {
			if nid != "" {
				uniqNSGs[nid] = struct{}{}
			}
		}
	}

	// Prefetch subnets (and related VCNs)
	if len(uniqSubnets) > 0 {
		jobs := make(chan Work, len(uniqSubnets))
		for sid := range uniqSubnets {
			id := sid
			jobs <- func() error {
				var sResp core.GetSubnetResponse
				err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
					return a.do(ctx, func() error {
						var e error
						sResp, e = a.nwClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &id})
						return e
					})
				})
				if err != nil {
					return err
				}
				a.muSubnets.Lock()
				a.subnetCache[id] = sResp
				a.muSubnets.Unlock()
				// Prefetch VCN by ID if not in cache
				if sResp.Subnet.VcnId != nil {
					vcnID := *sResp.Subnet.VcnId
					a.muVcns.RLock()
					_, ok := a.vcnCache[vcnID]
					a.muVcns.RUnlock()
					if !ok {
						var vcnResp core.GetVcnResponse
						_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
							return a.do(ctx, func() error {
								var e error
								vcnResp, e = a.nwClient.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
								return e
							})
						})
						if vcnResp.RawResponse != nil {
							a.muVcns.Lock()
							a.vcnCache[vcnID] = vcnResp
							a.muVcns.Unlock()
						}
					}
				}
				return nil
			}
		}
		close(jobs)
		_ = runWithWorkers(ctx, a.workerCount, jobs)
	}

	if len(uniqNSGs) > 0 {
		jobs := make(chan Work, len(uniqNSGs))
		for nid := range uniqNSGs {
			id := nid
			jobs <- func() error {
				var nResp core.GetNetworkSecurityGroupResponse
				err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
					return a.do(ctx, func() error {
						var e error
						nResp, e = a.nwClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &id})
						return e
					})
				})
				if err != nil {
					return err
				}
				a.muNsgs.Lock()
				a.nsgCache[id] = nResp
				a.muNsgs.Unlock()
				return nil
			}
		}
		close(jobs)
		_ = runWithWorkers(ctx, a.workerCount, jobs)
	}

	// Process each LB using a bounded worker pool as well
	out := make([]domain.LoadBalancer, len(items))
	jobs := make(chan Work, len(items))
	var mu sync.Mutex
	for i := range items {
		idx := i
		jobs <- func() error {
			mapped, err := a.enrichAndMapLoadBalancer(ctx, items[idx])
			if err != nil {
				return err
			}
			mu.Lock()
			out[idx] = *mapped
			mu.Unlock()
			return nil
		}
	}
	close(jobs)
	if err := runWithWorkers(ctx, a.workerCount, jobs); err != nil {
		return nil, err
	}
	return out, nil
}

// enrichAndMapLoadBalancer builds the domain model and enriches it with names, health, members and SSL certificate info
func (a *Adapter) enrichAndMapLoadBalancer(ctx context.Context, lb loadbalancer.LoadBalancer) (*domain.LoadBalancer, error) {
	startTotal := time.Now()
	id := ""
	name := ""
	if lb.Id != nil {
		id = *lb.Id
	}
	if lb.DisplayName != nil {
		name = *lb.DisplayName
	}
	// Start log for this LB enrichment
	lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.start", "id", id, "name", name)
	defer func() {
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.total", "id", id, "name", name, "duration_ms", time.Since(startTotal).Milliseconds())
	}()

	dm := mapping.NewDomainLoadBalancerFromAttrs(mapping.NewLoadBalancerAttributesFromOCILoadBalancer(lb))

	var (
		wg              sync.WaitGroup
		errCh           = make(chan error, 5)
		dResolveSubnets int64
		dResolveNSGs    int64
		dHealth         int64
		dMembers        int64
		dCerts          int64
		mu              sync.Mutex
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		if err := a.resolveSubnets(ctx, dm); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dResolveSubnets = time.Since(s).Milliseconds()
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		if err := a.resolveNSGs(ctx, dm); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dResolveNSGs = time.Since(s).Milliseconds()
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		if err := a.enrichBackendHealth(ctx, lb, dm, true); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dHealth = time.Since(s).Milliseconds()
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		if err := a.enrichBackendMembers(ctx, lb, dm, true); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dMembers = time.Since(s).Milliseconds()
		mu.Unlock()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		if err := a.enrichCertificates(ctx, lb, dm); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dCerts = time.Since(s).Milliseconds()
		mu.Unlock()
	}()

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return dm, err
		}
	}
	// Final summary with per-step durations
	lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.summary", "id", id, "name", name,
		"duration_total_ms", time.Since(startTotal).Milliseconds(),
		"resolve_subnets_ms", dResolveSubnets,
		"resolve_nsgs_ms", dResolveNSGs,
		"backend_health_ms", dHealth,
		"backend_members_ms", dMembers,
		"certificates_ms", dCerts,
	)
	return dm, nil
}
