package networklb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	nlbLogger "github.com/rozdolsky33/ocloud/internal/logger"
)

func cachedFetch[T any](ctx context.Context, a *Adapter, id string, mu *sync.RWMutex, cache map[string]T, fetch func(context.Context, string) (T, error)) (T, bool) {
	var resp T
	fromCache := false
	mu.RLock()
	if cached, ok := cache[id]; ok {
		resp = cached
		fromCache = true
	}
	mu.RUnlock()
	if fromCache {
		return resp, true
	}
	_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		return a.do(ctx, func() error {
			var e error
			resp, e = fetch(ctx, id)
			return e
		})
	})
	mu.Lock()
	cache[id] = resp
	mu.Unlock()
	return resp, false
}

func (a *Adapter) resolveSubnets(ctx context.Context, dm *domain.NetworkLoadBalancer) error {
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
		nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.resolve_vcn", "nlb_id", dm.OCID, "nlb_name", dm.Name, "vcn_id", dm.VcnID, "from_cache", vcnFromCache, "duration_ms", time.Since(vcnStart).Milliseconds())
	}
	nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.resolve_subnets", "nlb_id", dm.OCID, "nlb_name", dm.Name, "subnets", origCount, "cache_hits", cacheHits, "cache_misses", origCount-cacheHits, "duration_ms", time.Since(start).Milliseconds())
	return nil
}

func (a *Adapter) resolveNSGs(ctx context.Context, dm *domain.NetworkLoadBalancer) error {
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
	nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.resolve_nsgs", "nlb_id", dm.OCID, "nlb_name", dm.Name, "nsgs", origCount, "cache_hits", cacheHits, "cache_misses", origCount-cacheHits, "duration_ms", time.Since(start).Milliseconds())
	return nil
}
