package loadbalancer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/certificatesmanagement"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	lbLogger "github.com/rozdolsky33/ocloud/internal/logger"
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
	dm := toBaseDomainLoadBalancer(response.LoadBalancer)
	_ = a.enrichBackendHealth(ctx, response.LoadBalancer, &dm, false)
	_ = a.resolveSubnets(ctx, &dm)
	return &dm, nil
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
	return &dm, nil
}

// ListLoadBalancers returns all load balancers in the compartment (paginated) mapped to a base domain model,
// lightly enriched with backend health so the default view can display statuses without the full enrichment cost.
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
				dm := toBaseDomainLoadBalancer(lb)
				// Light enrichment: health and subnets (captures VCN as well)
				_ = a.enrichBackendHealth(ctx, lb, &dm, false)
				_ = a.resolveSubnets(ctx, &dm)
				mapped[idx] = dm
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
	out := make([]domain.LoadBalancer, len(items))
	var wg sync.WaitGroup
	errCh := make(chan error, len(items))
	for i := range items {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			mapped, err := a.enrichAndMapLoadBalancer(ctx, items[idx])
			if err != nil {
				errCh <- err
				return
			}
			out[idx] = mapped
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// toEnrichDomainLoadBalancer maps and enriches the OCI LB into the domain model (alias for enrichAndMapLoadBalancer)
func (a *Adapter) toEnrichDomainLoadBalancer(ctx context.Context, lb loadbalancer.LoadBalancer) (domain.LoadBalancer, error) {
	return a.enrichAndMapLoadBalancer(ctx, lb)
}

// enrichAndMapLoadBalancer builds the domain model and enriches it with names, health, members and SSL certificate info
func (a *Adapter) enrichAndMapLoadBalancer(ctx context.Context, lb loadbalancer.LoadBalancer) (domain.LoadBalancer, error) {
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
		lbLogger.LogWithLevel(lbLogger.CmdLogger, lbLogger.Debug, "lb.enrich.total", "id", id, "name", name, "duration_ms", time.Since(startTotal).Seconds())
	}()

	dm := toBaseDomainLoadBalancer(lb)

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
		if err := a.resolveSubnets(ctx, &dm); err != nil {
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
		if err := a.resolveNSGs(ctx, &dm); err != nil {
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
		if err := a.enrichBackendHealth(ctx, lb, &dm, true); err != nil {
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
		if err := a.enrichBackendMembers(ctx, lb, &dm, true); err != nil {
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
		if err := a.enrichCertificates(ctx, lb, &dm); err != nil {
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
		"duration_total_ms", time.Since(startTotal).Seconds(),
		"resolve_subnets_ms", dResolveSubnets,
		"resolve_nsgs_ms", dResolveNSGs,
		"backend_health_ms", dHealth,
		"backend_members_ms", dMembers,
		"certificates_ms", dCerts,
	)
	return dm, nil
}

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
	// Use cached list if available
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

	var wg sync.WaitGroup
	var mu sync.Mutex
	for n := range nameSet {
		name := n
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		}()
	}

	for cid := range idSet {
		id := cid
		wg.Add(1)
		go func() {
			defer wg.Done()
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
				return
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
		}()
	}
	wg.Wait()
	sort.Strings(out)
	dm.SSLCertificates = out
	return nil
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
	useSSL := false
	routingPolicySet := make(map[string]struct{})
	for name, l := range lb.Listeners {
		// Determine port
		port := 0
		if l.Port != nil {
			port = int(*l.Port)
		}
		// Determine backend set name
		backend := ""
		if l.DefaultBackendSetName != nil {
			backend = *l.DefaultBackendSetName
		}
		// Capture SSL usage
		if l.SslConfiguration != nil {
			useSSL = true
		}
		// Capture routing policy referenced by the listener, if any
		if l.RoutingPolicyName != nil && *l.RoutingPolicyName != "" {
			routingPolicySet[*l.RoutingPolicyName] = struct{}{}
		}
		protoLabel := "http"
		protoUpper := ""
		if l.Protocol != nil {
			protoUpper = strings.ToUpper(*l.Protocol)
		}
		if l.SslConfiguration != nil || port == 443 || port == 8443 {
			protoLabel = "https"
		} else if strings.EqualFold(protoUpper, "TCP") {
			protoLabel = "tcp"
		} else if port == 80 {
			protoLabel = "http"
		}
		listeners[name] = fmt.Sprintf("%s:%d → %s", protoLabel, port, backend)
	}
	// If no routing policies captured from listeners, fall back to keys of the load balancer's routing policies map
	if len(routingPolicySet) == 0 {
		for rpName := range lb.RoutingPolicies {
			if rpName != "" {
				routingPolicySet[rpName] = struct{}{}
			}
		}
	}
	routingPolicies := make([]string, 0, len(routingPolicySet))
	for rp := range routingPolicySet {
		routingPolicies = append(routingPolicies, rp)
	}
	sort.Strings(routingPolicies)

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
			switch port {
			case 443, 8443:
				p = "HTTPS"
			case 80:
				p = "HTTP"
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

	// Hostnames: collect FQDN values from LB hostname map
	hostnames := make([]string, 0)
	for n, h := range lb.Hostnames {
		if h.Hostname != nil && *h.Hostname != "" {
			hostnames = append(hostnames, *h.Hostname)
			continue
		}
		// fallback to the map key if value missing
		if strings.TrimSpace(n) != "" {
			hostnames = append(hostnames, n)
		}
	}
	sort.Strings(hostnames)

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
		RoutingPolicies: routingPolicies,
		UseSSL:          useSSL,
		Hostnames:       hostnames,
	}
}
