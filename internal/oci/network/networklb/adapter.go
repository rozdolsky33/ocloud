package networklb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/networkloadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	nlbLogger "github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"golang.org/x/time/rate"
)

// Adapter implements the domain.NetworkLoadBalancerRepository interface for OCI.
type Adapter struct {
	nlbClient   networkloadbalancer.NetworkLoadBalancerClient
	nwClient    core.VirtualNetworkClient
	limiter     *rate.Limiter
	workerCount int
	subnetCache map[string]core.GetSubnetResponse
	vcnCache    map[string]core.GetVcnResponse
	nsgCache    map[string]core.GetNetworkSecurityGroupResponse
	muSubnets   sync.RWMutex
	muVcns      sync.RWMutex
	muNsgs      sync.RWMutex
}

// NewAdapter creates a new Adapter instance using pre-created OCI clients.
func NewAdapter(nlbClient networkloadbalancer.NetworkLoadBalancerClient, nwClient core.VirtualNetworkClient) *Adapter {
	return &Adapter{
		nlbClient:   nlbClient,
		nwClient:    nwClient,
		workerCount: defaultWorkerCount,
		limiter:     rate.NewLimiter(rate.Limit(defaultRatePerSec), defaultRateBurst),
		subnetCache: make(map[string]core.GetSubnetResponse),
		vcnCache:    make(map[string]core.GetVcnResponse),
		nsgCache:    make(map[string]core.GetNetworkSecurityGroupResponse),
	}
}

// GetNetworkLoadBalancer retrieves a single Network Load Balancer and maps it to the basic domain model.
func (a *Adapter) GetNetworkLoadBalancer(ctx context.Context, ocid string) (*domain.NetworkLoadBalancer, error) {
	response, err := a.nlbClient.GetNetworkLoadBalancer(ctx, networkloadbalancer.GetNetworkLoadBalancerRequest{
		NetworkLoadBalancerId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get network load balancer: %w", err)
	}
	dm := mapping.NewDomainNetworkLoadBalancerFromAttrs(mapping.NewNetworkLoadBalancerAttributesFromOCI(response.NetworkLoadBalancer))
	// Enrich backend health from the full NLB object
	_ = a.enrichBackendHealth(ctx, response.NetworkLoadBalancer, dm, false)
	_ = a.resolveSubnets(ctx, dm)
	return dm, nil
}

// GetEnrichedNetworkLoadBalancer retrieves a single Network Load Balancer and returns the enriched domain model.
func (a *Adapter) GetEnrichedNetworkLoadBalancer(ctx context.Context, ocid string) (*domain.NetworkLoadBalancer, error) {
	response, err := a.nlbClient.GetNetworkLoadBalancer(ctx, networkloadbalancer.GetNetworkLoadBalancerRequest{
		NetworkLoadBalancerId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get network load balancer: %w", err)
	}
	dm, err := a.enrichAndMapNetworkLoadBalancer(ctx, response.NetworkLoadBalancer)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

// ListNetworkLoadBalancers returns all network load balancers in the compartment mapped to a base domain model.
func (a *Adapter) ListNetworkLoadBalancers(ctx context.Context, compartmentID string) ([]domain.NetworkLoadBalancer, error) {
	result := make([]domain.NetworkLoadBalancer, 0)
	var page *string
	for {
		resp, err := a.nlbClient.ListNetworkLoadBalancers(ctx, networkloadbalancer.ListNetworkLoadBalancersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing network load balancers: %w", err)
		}
		pageItems := resp.Items
		mapped := make([]domain.NetworkLoadBalancer, len(pageItems))
		var wg sync.WaitGroup
		for i := range pageItems {
			wg.Add(1)
			idx := i
			go func() {
				defer wg.Done()
				nlb := pageItems[idx]
				dm := mapping.NewDomainNetworkLoadBalancerFromAttrs(mapping.NewNetworkLoadBalancerAttributesFromOCISummary(nlb))
				// Collect backend set names from the summary for health enrichment
				bsNames := make([]string, 0, len(nlb.BackendSets))
				for name := range nlb.BackendSets {
					bsNames = append(bsNames, name)
				}
				if nlb.Id != nil {
					_ = a.enrichBackendHealthFromSummary(ctx, *nlb.Id, bsNames, dm)
				}
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

// ListEnrichedNetworkLoadBalancers returns all network load balancers with full enrichment.
func (a *Adapter) ListEnrichedNetworkLoadBalancers(ctx context.Context, compartmentID string) ([]domain.NetworkLoadBalancer, error) {
	result := make([]domain.NetworkLoadBalancer, 0)
	var page *string
	for {
		resp, err := a.nlbClient.ListNetworkLoadBalancers(ctx, networkloadbalancer.ListNetworkLoadBalancersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing network load balancers: %w", err)
		}

		mapped, err := a.enrichAndMapNetworkLoadBalancerSummaries(ctx, resp.Items)
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

// enrichAndMapNetworkLoadBalancerSummaries fetches the full NLB for each summary, then enriches.
func (a *Adapter) enrichAndMapNetworkLoadBalancerSummaries(ctx context.Context, items []networkloadbalancer.NetworkLoadBalancerSummary) ([]domain.NetworkLoadBalancer, error) {
	// Prefetch subnets and NSGs
	uniqSubnets := make(map[string]struct{})
	uniqNSGs := make(map[string]struct{})
	for _, nlb := range items {
		if nlb.SubnetId != nil && *nlb.SubnetId != "" {
			uniqSubnets[*nlb.SubnetId] = struct{}{}
		}
		for _, nid := range nlb.NetworkSecurityGroupIds {
			if nid != "" {
				uniqNSGs[nid] = struct{}{}
			}
		}
	}

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

	// For enriched mode, fetch full NLB objects and enrich
	out := make([]domain.NetworkLoadBalancer, len(items))
	jobs := make(chan Work, len(items))
	var mu sync.Mutex
	for i := range items {
		idx := i
		jobs <- func() error {
			nlbID := ""
			if items[idx].Id != nil {
				nlbID = *items[idx].Id
			}
			// Fetch full NLB for enrichment
			response, err := a.nlbClient.GetNetworkLoadBalancer(ctx, networkloadbalancer.GetNetworkLoadBalancerRequest{
				NetworkLoadBalancerId: &nlbID,
			})
			if err != nil {
				return fmt.Errorf("getting network load balancer for enrichment: %w", err)
			}
			mapped, err := a.enrichAndMapNetworkLoadBalancer(ctx, response.NetworkLoadBalancer)
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

// enrichAndMapNetworkLoadBalancer builds the domain model and enriches it with names, health, and members.
func (a *Adapter) enrichAndMapNetworkLoadBalancer(ctx context.Context, nlb networkloadbalancer.NetworkLoadBalancer) (*domain.NetworkLoadBalancer, error) {
	startTotal := time.Now()
	id := ""
	name := ""
	if nlb.Id != nil {
		id = *nlb.Id
	}
	if nlb.DisplayName != nil {
		name = *nlb.DisplayName
	}
	nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.start", "id", id, "name", name)
	defer func() {
		nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.total", "id", id, "name", name, "duration_ms", time.Since(startTotal).Milliseconds())
	}()

	dm := mapping.NewDomainNetworkLoadBalancerFromAttrs(mapping.NewNetworkLoadBalancerAttributesFromOCI(nlb))

	var (
		wg              sync.WaitGroup
		errCh           = make(chan error, 4)
		dResolveSubnets int64
		dResolveNSGs    int64
		dHealth         int64
		dMembers        int64
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
		if err := a.enrichBackendHealth(ctx, nlb, dm, true); err != nil {
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
		if err := a.enrichBackendMembers(ctx, nlb, dm, true); err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		dMembers = time.Since(s).Milliseconds()
		mu.Unlock()
	}()

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return dm, err
		}
	}
	nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.summary", "id", id, "name", name,
		"duration_total_ms", time.Since(startTotal).Milliseconds(),
		"resolve_subnets_ms", dResolveSubnets,
		"resolve_nsgs_ms", dResolveNSGs,
		"backend_health_ms", dHealth,
		"backend_members_ms", dMembers,
	)
	return dm, nil
}
