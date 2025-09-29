package loadbalancer

import (
	"context"
	x509std "crypto/x509"
	pemenc "encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	lbLogger "github.com/rozdolsky33/ocloud/internal/logger"
	"golang.org/x/sync/errgroup"
)

const (
	defaultMaxRetries     = 5
	defaultInitialBackoff = 1 * time.Second
	defaultMaxBackoff     = 32 * time.Second
	defaultWorkerCount    = 12
	defaultRatePerSec     = 10
	defaultRateBurst      = 5
)

// Work represents a unit of work executed by the worker pool
// It returns an error to allow early cancellation via errgroup.
type Work func() error

// runWithWorkers executes jobs from the channel using n workers and stops on first error or context cancel.
func runWithWorkers(ctx context.Context, n int, jobs <-chan Work) error {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < n; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case w, ok := <-jobs:
					if !ok {
						return nil
					}
					if err := w(); err != nil {
						return err
					}
				}
			}
		})
	}
	return g.Wait()
}

// do apply a central rate limit before performing the given operation.
func (a *Adapter) do(ctx context.Context, op func() error) error {
	if a.limiter != nil {
		if err := a.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	return op()
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
		var resp core.GetSubnetResponse
		var fromCache bool
		// try cache first
		a.muSubnets.RLock()
		if cached, ok := a.subnetCache[id]; ok {
			resp = cached
			fromCache = true
		}
		a.muSubnets.RUnlock()
		if !fromCache {
			// fetch and cache
			_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					resp, e = a.nwClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &id})
					return e
				})
			})
			a.muSubnets.Lock()
			a.subnetCache[id] = resp
			a.muSubnets.Unlock()
		} else {
			cacheHits++
		}
		// use response if present
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

	// If we have a VCN ID, resolve its name and set on the domain model
	if capturedVcnID != "" {
		vcnStart := time.Now()
		var vcnResp core.GetVcnResponse
		var vcnFromCache bool
		a.muVcns.RLock()
		if cached, ok := a.vcnCache[capturedVcnID]; ok {
			vcnResp = cached
			vcnFromCache = true
		}
		a.muVcns.RUnlock()
		if !vcnFromCache {
			_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					vcnResp, e = a.nwClient.GetVcn(ctx, core.GetVcnRequest{VcnId: &capturedVcnID})
					return e
				})
			})
			a.muVcns.Lock()
			a.vcnCache[capturedVcnID] = vcnResp
			a.muVcns.Unlock()
		}
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
		var resp core.GetNetworkSecurityGroupResponse
		var fromCache bool
		a.muNsgs.RLock()
		if cached, ok := a.nsgCache[id]; ok {
			resp = cached
			fromCache = true
		}
		a.muNsgs.RUnlock()
		if !fromCache {
			_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					resp, e = a.nwClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &id})
					return e
				})
			})
			a.muNsgs.Lock()
			a.nsgCache[id] = resp
			a.muNsgs.Unlock()
		} else {
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

// retryOnRateLimit retries the provided operation when OCI responds with HTTP 429 rate limited.
// It applies exponential backoff between retries and preserves the original behavior and error messages.
func retryOnRateLimit(ctx context.Context, maxRetries int, initialBackoff, maxBackoff time.Duration, op func() error) error {
	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := op()
		if err == nil {
			return nil
		}

		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			if attempt == maxRetries-1 {
				return fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		return err
	}
	return nil
}
