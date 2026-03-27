package networklb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/networkloadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	nlbLogger "github.com/rozdolsky33/ocloud/internal/logger"
)

func (a *Adapter) enrichBackendHealth(ctx context.Context, nlb networkloadbalancer.NetworkLoadBalancer, dm *domain.NetworkLoadBalancer, deep bool) error {
	start := time.Now()
	nlbID, nlbName := "", ""
	if nlb.Id != nil {
		nlbID = *nlb.Id
	}
	if nlb.DisplayName != nil {
		nlbName = *nlb.DisplayName
	}
	defer func() {
		nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.backend_health", "id", nlbID, "name", nlbName, "backend_sets", len(nlb.BackendSets), "duration_ms", time.Since(start).Milliseconds())
	}()
	if nlb.Id == nil {
		return nil
	}
	healthLocal := make(map[string]string)
	jobs := make(chan Work, len(nlb.BackendSets))
	var mu sync.Mutex
	for bsName := range nlb.BackendSets {
		name := bsName
		jobs <- func() error {
			var hResp networkloadbalancer.GetBackendSetHealthResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					hResp, e = a.nlbClient.GetBackendSetHealth(ctx, networkloadbalancer.GetBackendSetHealthRequest{
						NetworkLoadBalancerId: nlb.Id,
						BackendSetName:        &name,
					})
					return e
				})
			})
			if err != nil {
				return err
			}
			status := strings.ToUpper(string(hResp.BackendSetHealth.Status))
			mu.Lock()
			healthLocal[name] = status
			mu.Unlock()
			return nil
		}
	}
	close(jobs)
	_ = runWithWorkers(ctx, a.workerCount, jobs)
	if dm.BackendHealth == nil {
		dm.BackendHealth = map[string]string{}
	}
	for k, v := range healthLocal {
		dm.BackendHealth[k] = v
	}
	return nil
}

func (a *Adapter) enrichBackendMembers(ctx context.Context, nlb networkloadbalancer.NetworkLoadBalancer, dm *domain.NetworkLoadBalancer, deep bool) error {
	start := time.Now()
	nlbID, nlbName := "", ""
	if nlb.Id != nil {
		nlbID = *nlb.Id
	}
	if nlb.DisplayName != nil {
		nlbName = *nlb.DisplayName
	}
	defer func() {
		nlbLogger.LogWithLevel(nlbLogger.CmdLogger, nlbLogger.Debug, "nlb.enrich.backend_members", "id", nlbID, "name", nlbName, "backend_sets", len(nlb.BackendSets), "duration_ms", time.Since(start).Milliseconds())
	}()
	if nlb.Id == nil {
		return nil
	}
	jobs := make(chan Work, len(nlb.BackendSets))
	var mu sync.Mutex
	for bsName := range nlb.BackendSets {
		name := bsName
		jobs <- func() error {
			var bsResp networkloadbalancer.GetBackendSetResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					bsResp, e = a.nlbClient.GetBackendSet(ctx, networkloadbalancer.GetBackendSetRequest{
						NetworkLoadBalancerId: nlb.Id,
						BackendSetName:        &name,
					})
					return e
				})
			})
			if err != nil {
				return err
			}
			backends := make([]domain.Backend, len(bsResp.BackendSet.Backends))
			setStatus := strings.ToUpper(dm.BackendHealth[name])
			needDeep := deep || (setStatus != "" && setStatus != "OK")
			for i, b := range bsResp.BackendSet.Backends {
				ip := ""
				if b.IpAddress != nil {
					ip = *b.IpAddress
				}
				port := 0
				if b.Port != nil {
					port = *b.Port
				}
				status := "UNKNOWN"
				if needDeep && ip != "" && port > 0 {
					backendName := fmt.Sprintf("%s:%d", ip, port)
					var bhResp networkloadbalancer.GetBackendHealthResponse
					_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
						return a.do(ctx, func() error {
							var e error
							bhResp, e = a.nlbClient.GetBackendHealth(ctx, networkloadbalancer.GetBackendHealthRequest{
								NetworkLoadBalancerId: nlb.Id,
								BackendSetName:        &name,
								BackendName:           &backendName,
							})
							return e
						})
					})
					if bhResp.RawResponse != nil {
						status = strings.ToUpper(string(bhResp.BackendHealth.Status))
					}
				}
				backends[i] = domain.Backend{Name: ip, Port: port, Status: status}
			}
			mu.Lock()
			bs := dm.BackendSets[name]
			bs.Backends = backends
			dm.BackendSets[name] = bs
			mu.Unlock()
			return nil
		}
	}
	close(jobs)
	_ = runWithWorkers(ctx, a.workerCount, jobs)
	return nil
}

func (a *Adapter) enrichBackendHealthFromSummary(ctx context.Context, nlbID string, backendSetNames []string, dm *domain.NetworkLoadBalancer) error {
	healthLocal := make(map[string]string)
	jobs := make(chan Work, len(backendSetNames))
	var mu sync.Mutex
	for _, bsName := range backendSetNames {
		name := bsName
		jobs <- func() error {
			var hResp networkloadbalancer.GetBackendSetHealthResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					hResp, e = a.nlbClient.GetBackendSetHealth(ctx, networkloadbalancer.GetBackendSetHealthRequest{
						NetworkLoadBalancerId: &nlbID,
						BackendSetName:        &name,
					})
					return e
				})
			})
			if err != nil {
				return err
			}
			status := strings.ToUpper(string(hResp.BackendSetHealth.Status))
			mu.Lock()
			healthLocal[name] = status
			mu.Unlock()
			return nil
		}
	}
	close(jobs)
	_ = runWithWorkers(ctx, a.workerCount, jobs)
	if dm.BackendHealth == nil {
		dm.BackendHealth = map[string]string{}
	}
	for k, v := range healthLocal {
		dm.BackendHealth[k] = v
	}
	return nil
}
