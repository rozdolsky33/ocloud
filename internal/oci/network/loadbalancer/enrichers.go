package loadbalancer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	lbLogger "github.com/rozdolsky33/ocloud/internal/logger"
)

// enrichBackendHealth fetches overall status per backend set and fills dm.BackendHealth.
// When deep is false, it only fetches per-set health; per-backend health is handled in enrichBackendMembers when deep.
func (a *Adapter) enrichBackendHealth(ctx context.Context, lb loadbalancer.LoadBalancer, dm *domain.LoadBalancer, deep bool) error {
	start := time.Now()
	lbID, lbName := "", ""
	if lb.Id != nil {
		lbID = *lb.Id
	}
	if lb.DisplayName != nil {
		lbName = *lb.DisplayName
	}
	defer func() {
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.backend_health", "id", lbID, "name", lbName, "backend_sets", len(lb.BackendSets), "duration_ms", time.Since(start).Milliseconds())
	}()
	if lb.Id == nil {
		return nil
	}
	healthLocal := make(map[string]string)
	jobs := make(chan Work, len(lb.BackendSets))
	var mu sync.Mutex
	for bsName := range lb.BackendSets {
		name := bsName
		jobs <- func() error {
			var hResp loadbalancer.GetBackendSetHealthResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					hResp, e = a.lbClient.GetBackendSetHealth(ctx, loadbalancer.GetBackendSetHealthRequest{LoadBalancerId: lb.Id, BackendSetName: &name})
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

// enrichBackendMembers fetches backend members per backend set and fills dm.BackendSets[...].Backends.
// It avoids per-backend GetBackendHealth unless deep is true or the set is unhealthy.
func (a *Adapter) enrichBackendMembers(ctx context.Context, lb loadbalancer.LoadBalancer, dm *domain.LoadBalancer, deep bool) error {
	start := time.Now()
	lbID, lbName := "", ""
	if lb.Id != nil {
		lbID = *lb.Id
	}
	if lb.DisplayName != nil {
		lbName = *lb.DisplayName
	}
	defer func() {
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.backend_members", "id", lbID, "name", lbName, "backend_sets", len(lb.BackendSets), "duration_ms", time.Since(start).Milliseconds())
	}()
	if lb.Id == nil {
		return nil
	}
	jobs := make(chan Work, len(lb.BackendSets))
	var mu sync.Mutex
	for bsName := range lb.BackendSets {
		name := bsName
		jobs <- func() error {
			var bsResp loadbalancer.GetBackendSetResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				return a.do(ctx, func() error {
					var e error
					bsResp, e = a.lbClient.GetBackendSet(ctx, loadbalancer.GetBackendSetRequest{LoadBalancerId: lb.Id, BackendSetName: &name})
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
					port = int(*b.Port)
				}
				status := "UNKNOWN"
				if needDeep && ip != "" && port > 0 {
					backendName := fmt.Sprintf("%s:%d", ip, port)
					var bhResp loadbalancer.GetBackendHealthResponse
					_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
						return a.do(ctx, func() error {
							var e error
							bhResp, e = a.lbClient.GetBackendHealth(ctx, loadbalancer.GetBackendHealthRequest{LoadBalancerId: lb.Id, BackendSetName: &name, BackendName: &backendName})
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

// enrichCertificates gathers certificate names/ids and resolves expiry where possible, storing formatted strings in dm.SSLCertificates
func (a *Adapter) enrichCertificates(ctx context.Context, lb loadbalancer.LoadBalancer, dm *domain.LoadBalancer) error {
	start := time.Now()
	lbID, lbName := "", ""
	if lb.Id != nil {
		lbID = *lb.Id
	}
	if lb.DisplayName != nil {
		lbName = *lb.DisplayName
	}
	defer func() {
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.certificates", "id", lbID, "name", lbName, "certs_count", len(dm.SSLCertificates), "duration_ms", time.Since(start).Milliseconds())
	}()
	out := make([]string, 0)
	if lb.Id == nil {
		dm.SSLCertificates = out
		return nil
	}

	nameSet := make(map[string]struct{})
	idSet := make(map[string]struct{})
	for _, l := range lb.Listeners {
		if l.SslConfiguration != nil {
			if l.SslConfiguration.CertificateName != nil {
				if n := strings.TrimSpace(*l.SslConfiguration.CertificateName); n != "" {
					nameSet[n] = struct{}{}
				}
			}
			for _, cid := range l.SslConfiguration.CertificateIds {
				if c := strings.TrimSpace(cid); c != "" {
					idSet[c] = struct{}{}
				}
			}
		}
	}

	certsByName := make(map[string]loadbalancer.Certificate)
	var listResp loadbalancer.ListCertificatesResponse
	var listItems []loadbalancer.Certificate
	cacheHit := false
	if lbID != "" {
		a.muCertLists.RLock()
		if cached, ok := a.certListCache[lbID]; ok {
			listItems = cached
			cacheHit = true
		}
		a.muCertLists.RUnlock()
	}
	lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.certificates.list_cache", "id", lbID, "name", lbName, "cache_hit", cacheHit)
	if listItems == nil {
		_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
			return a.do(ctx, func() error {
				var e error
				listResp, e = a.lbClient.ListCertificates(ctx, loadbalancer.ListCertificatesRequest{LoadBalancerId: lb.Id})
				return e
			})
		})
		listItems = listResp.Items
		if lbID != "" {
			a.muCertLists.Lock()
			a.certListCache[lbID] = listItems
			a.muCertLists.Unlock()
		}
	}
	if len(listItems) > 0 {
		for _, c := range listItems {
			if c.CertificateName != nil {
				certsByName[*c.CertificateName] = c
				nameSet[*c.CertificateName] = struct{}{}
			}
		}
	} else {
		for n, c := range lb.Certificates {
			certsByName[n] = c
			nameSet[n] = struct{}{}
		}
	}

	if len(nameSet) == 0 {
		var getResp loadbalancer.GetLoadBalancerResponse
		if err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
			var e error
			getResp, e = a.lbClient.GetLoadBalancer(ctx, loadbalancer.GetLoadBalancerRequest{LoadBalancerId: lb.Id})
			return e
		}); err == nil {
			for n, c := range getResp.LoadBalancer.Certificates {
				certsByName[n] = c
				nameSet[n] = struct{}{}
			}
			for _, l := range getResp.LoadBalancer.Listeners {
				if l.SslConfiguration != nil && l.SslConfiguration.CertificateName != nil {
					if n := strings.TrimSpace(*l.SslConfiguration.CertificateName); n != "" {
						nameSet[n] = struct{}{}
					}
				}
			}
		}
	}

	jobs := make(chan Work, len(nameSet)+len(idSet))
	var mu sync.Mutex

	for n := range nameSet {
		name := n
		jobs <- func() error {
			expires := ""
			if c, ok := certsByName[name]; ok {
				if c.PublicCertificate != nil && *c.PublicCertificate != "" {
					if t, ok := parseCertNotAfter(*c.PublicCertificate); ok {
						expires = t.Format("2006-01-02")
					}
				}
			}
			display := name
			if expires != "" {
				display = fmt.Sprintf("%s (Expires: %s)", name, expires)
			}
			mu.Lock()
			out = append(out, display)
			mu.Unlock()
			return nil
		}
	}

	for cid := range idSet {
		id := cid
		jobs <- func() error {
			var certResp certificatesmanagement.GetCertificateResponse
			err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
				var e error
				certResp, e = a.certsClient.GetCertificate(ctx, certificatesmanagement.GetCertificateRequest{CertificateId: &id})
				return e
			})
			if err != nil {
				mu.Lock()
				out = append(out, id)
				mu.Unlock()
				return nil
			}
			name := id
			if certResp.Certificate.Name != nil && *certResp.Certificate.Name != "" {
				name = *certResp.Certificate.Name
			}
			var expiresStr string
			if certResp.Certificate.CurrentVersion != nil && certResp.Certificate.CurrentVersion.VersionNumber != nil {
				ver := *certResp.Certificate.CurrentVersion.VersionNumber
				var verResp certificatesmanagement.GetCertificateVersionResponse
				_ = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
					var e error
					verResp, e = a.certsClient.GetCertificateVersion(ctx, certificatesmanagement.GetCertificateVersionRequest{CertificateId: &id, CertificateVersionNumber: &ver})
					return e
				})
				if verResp.CertificateVersion.Validity != nil && verResp.CertificateVersion.Validity.TimeOfValidityNotAfter != nil {
					expiresStr = verResp.CertificateVersion.Validity.TimeOfValidityNotAfter.Time.Format("2006-01-02")
				}
			}
			display := name
			if expiresStr != "" {
				display = fmt.Sprintf("%s (Expires: %s)", name, expiresStr)
			}
			mu.Lock()
			out = append(out, display)
			mu.Unlock()
			return nil
		}
	}

	close(jobs)
	_ = runWithWorkers(ctx, a.workerCount, jobs)

	sort.Strings(out)
	dm.SSLCertificates = out
	return nil
}
