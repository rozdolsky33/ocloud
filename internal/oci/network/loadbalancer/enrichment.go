package loadbalancer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	lbLogger "github.com/rozdolsky33/ocloud/internal/logger"
)

// cachedFetch centralizes cache lookup, rate-limited fetch via Adapter.do, and cache population.
// It returns the response and whether it was served from a cache.
func cachedFetch[T any](ctx context.Context, a *Adapter, id string, mu *sync.RWMutex, cache map[string]T, fetch func(context.Context, string) (T, error)) (T, bool) {
	var resp T
	fromCache := false
	// Fast path: read lock and check cache
	mu.RLock()
	if cached, ok := cache[id]; ok {
		resp = cached
		fromCache = true
	}
	mu.RUnlock()
	if fromCache {
		return resp, true
	}
	// Miss: perform fetch with retry and rate-limited do wrapper
	_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		return a.do(ctx, func() error {
			var e error
			resp, e = fetch(ctx, id)
			return e
		})
	})
	// Store in cache
	mu.Lock()
	cache[id] = resp
	mu.Unlock()
	return resp, false
}

// resolveSubnets resolves subnet IDs on the domain model to "Name (CIDR)" and captures the VCN context
func (a *Adapter) resolveSubnets(ctx context.Context, dm *domain.LoadBalancer) error {
	start := time.Now()
	origCount := len(dm.Subnets)
	cacheHits := 0
	resolved := make([]string, 0, len(dm.Subnets))
	var capturedVcnID string
	for _, sid := range dm.Subnets {
		id := sid
		if id == "" {
			continue
		}
		resp, fromCache := cachedFetch(ctx, a, id, &a.muSubnets, a.subnetCache, func(ctx context.Context, id string) (core.GetSubnetResponse, error) {
			return a.nwClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &id})
		})
		if fromCache {
			cacheHits++
		}

		if resp.Subnet.Id != nil {
			if capturedVcnID == "" && resp.Subnet.VcnId != nil {
				capturedVcnID = *resp.Subnet.VcnId
			}
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
		resolved = append(resolved, sid)
	}
	dm.Subnets = resolved

	if capturedVcnID != "" {
		vcnStart := time.Now()
		vcnResp, vcnFromCache := cachedFetch(ctx, a, capturedVcnID, &a.muVcns, a.vcnCache, func(ctx context.Context, id string) (core.GetVcnResponse, error) {
			return a.nwClient.GetVcn(ctx, core.GetVcnRequest{VcnId: &id})
		})
		dm.VcnID = capturedVcnID
		if vcnResp.Vcn.DisplayName != nil {
			dm.VcnName = *vcnResp.Vcn.DisplayName
		}
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.resolve_vcn", "lb_id", dm.OCID, "lb_name", dm.Name, "vcn_id", dm.VcnID, "from_cache", vcnFromCache, "duration_ms", time.Since(vcnStart).Milliseconds())
	}
	lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.resolve_subnets", "lb_id", dm.OCID, "lb_name", dm.Name, "subnets", origCount, "cache_hits", cacheHits, "cache_misses", origCount-cacheHits, "duration_ms", time.Since(start).Milliseconds())
	return nil
}

// resolveNSGs resolves NSG IDs on the domain model to display names best-effort
func (a *Adapter) resolveNSGs(ctx context.Context, dm *domain.LoadBalancer) error {
	start := time.Now()
	origCount := len(dm.NSGs)
	cacheHits := 0
	resolved := make([]string, 0, len(dm.NSGs))
	for _, nid := range dm.NSGs {
		id := nid
		if id == "" {
			continue
		}
		resp, fromCache := cachedFetch(ctx, a, id, &a.muNsgs, a.nsgCache, func(ctx context.Context, id string) (core.GetNetworkSecurityGroupResponse, error) {
			return a.nwClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &id})
		})
		if fromCache {
			cacheHits++
		}
		if resp.NetworkSecurityGroup.DisplayName != nil && *resp.NetworkSecurityGroup.DisplayName != "" {
			resolved = append(resolved, *resp.NetworkSecurityGroup.DisplayName)
			continue
		}
		resolved = append(resolved, nid)
	}
	dm.NSGs = resolved
	lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.resolve_nsgs", "lb_id", dm.OCID, "lb_name", dm.Name, "nsgs", origCount, "cache_hits", cacheHits, "cache_misses", origCount-cacheHits, "duration_ms", time.Since(start).Milliseconds())
	return nil
}
