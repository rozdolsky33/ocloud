package vcn

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
)

const (
	defaultMaxRetries     = 5
	defaultInitialBackoff = 1 * time.Second
	defaultMaxBackoff     = 32 * time.Second
)

// Adapter provides access to VCN-related OCI APIs.
// It is infra-layer and should be used by the service layer.
type Adapter struct {
	client core.VirtualNetworkClient
}

// NewAdapter creates a new adapter instance.
func NewAdapter(client core.VirtualNetworkClient) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GetEnrichedVcn(ctx context.Context, vcnID string) (domain.VCN, error) {
	var resp core.GetVcnResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
		return e
	})
	if err != nil {
		return domain.VCN{}, fmt.Errorf("getting VCN from OCI: %w", err)
	}
	m := toDomainVCNModel(resp.Vcn)
	return m, nil
}

// ListVcns lists all VCNs in a given compartment.
func (a *Adapter) ListVcns(ctx context.Context, compartmentID string) ([]domain.VCN, error) {
	req := core.ListVcnsRequest{CompartmentId: &compartmentID}
	var out []domain.VCN
	for {
		var resp core.ListVcnsResponse
		err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
			var e error
			resp, e = a.client.ListVcns(ctx, req)
			return e
		})
		if err != nil {
			return nil, fmt.Errorf("listing VCNs from OCI: %w", err)
		}
		for _, v := range resp.Items {
			out = append(out, toDomainVCNModel(v))
		}
		if resp.OpcNextPage == nil {
			break
		}
		req.Page = resp.OpcNextPage
	}
	return out, nil
}

// ListEnrichedVcns lists VCNs and enriches them with all related resources in parallel.
func (a *Adapter) ListEnrichedVcns(ctx context.Context, compartmentID string) ([]domain.VCN, error) {
	vcns, err := a.ListVcns(ctx, compartmentID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(vcns))

	for i := range vcns {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := a.enrichVCN(ctx, &vcns[i]); err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return vcns, nil
}

func (a *Adapter) enrichVCN(ctx context.Context, vcn *domain.VCN) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	var mutex sync.Mutex

	wg.Add(10)
	go func() {
		defer wg.Done()
		gateways, err := a.listInternetGateways(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
			return
		}
		mutex.Lock()
		vcn.Gateways = append(vcn.Gateways, gateways...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		nats, err := a.listNatGateways(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
			return
		}
		mutex.Lock()
		vcn.Gateways = append(vcn.Gateways, nats...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		sgws, err := a.listServiceGateways(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
			return
		}
		mutex.Lock()
		vcn.Gateways = append(vcn.Gateways, sgws...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		lpgs, err := a.listLocalPeeringGateways(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
			return
		}
		mutex.Lock()
		vcn.Gateways = append(vcn.Gateways, lpgs...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		drg, err := a.listDrgAttachments(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
			return
		}
		mutex.Lock()
		vcn.Gateways = append(vcn.Gateways, drg...)
		mutex.Unlock()
	}()
	go func() {
		defer wg.Done()
		var err error
		vcn.Subnets, err = a.listSubnets(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		vcn.RouteTables, err = a.listRouteTables(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		vcn.SecurityLists, err = a.listSecurityLists(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		vcn.NSGs, err = a.listNetworkSecurityGroups(ctx, vcn.CompartmentID, vcn.OCID)
		if err != nil {
			errCh <- err
		}
	}()
	go func() {
		defer wg.Done()
		if vcn.DhcpOptionsID != "" {
			dhcp, err := a.GetDhcpOptions(ctx, vcn.DhcpOptionsID)
			if err != nil {
				errCh <- err
			} else {
				vcn.DhcpOptions = dhcp
			}
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) listInternetGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListInternetGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListInternetGatewaysResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListInternetGateways(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Internet"})
	}
	return gateways, nil
}

func (a *Adapter) listNatGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListNatGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListNatGatewaysResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListNatGateways(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "NAT"})
	}
	return gateways, nil
}

func (a *Adapter) listServiceGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListServiceGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListServiceGatewaysResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListServiceGateways(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Service"})
	}
	return gateways, nil
}

func (a *Adapter) listLocalPeeringGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListLocalPeeringGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListLocalPeeringGatewaysResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListLocalPeeringGateways(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Local Peering"})
	}
	return gateways, nil
}

func (a *Adapter) listDrgAttachments(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListDrgAttachmentsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListDrgAttachmentsResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListDrgAttachments(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "DRG"})
	}
	return gateways, nil
}

func (a *Adapter) listRouteTables(ctx context.Context, compartmentID, vcnID string) ([]domain.RouteTable, error) {
	req := core.ListRouteTablesRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListRouteTablesResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListRouteTables(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var rts []domain.RouteTable
	for _, item := range resp.Items {
		rts = append(rts, domain.RouteTable{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return rts, nil
}

func (a *Adapter) listSecurityLists(ctx context.Context, compartmentID, vcnID string) ([]domain.SecurityList, error) {
	req := core.ListSecurityListsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListSecurityListsResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListSecurityLists(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var sls []domain.SecurityList
	for _, item := range resp.Items {
		sls = append(sls, domain.SecurityList{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return sls, nil
}

func (a *Adapter) listNetworkSecurityGroups(ctx context.Context, compartmentID, vcnID string) ([]domain.NSG, error) {
	req := core.ListNetworkSecurityGroupsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListNetworkSecurityGroupsResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListNetworkSecurityGroups(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var nsgs []domain.NSG
	for _, item := range resp.Items {
		nsgs = append(nsgs, domain.NSG{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return nsgs, nil
}

func (a *Adapter) GetDhcpOptions(ctx context.Context, dhcpID string) (domain.DhcpOptions, error) {
	var resp core.GetDhcpOptionsResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.GetDhcpOptions(ctx, core.GetDhcpOptionsRequest{DhcpId: &dhcpID})
		return e
	})
	if err != nil {
		return domain.DhcpOptions{}, err
	}
	return domain.DhcpOptions{OCID: *resp.Id, DisplayName: *resp.DisplayName, LifecycleState: string(resp.LifecycleState), DomainNameType: ""}, nil
}

func (a *Adapter) listSubnets(ctx context.Context, compartmentID, vcnID string) ([]domain.Subnet, error) {
	req := core.ListSubnetsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	var resp core.ListSubnetsResponse
	err := retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.client.ListSubnets(ctx, req)
		return e
	})
	if err != nil {
		return nil, err
	}
	var subnets []domain.Subnet
	for _, item := range resp.Items {
		var id, name, cidr, rtID string
		if item.Id != nil {
			id = *item.Id
		}
		if item.DisplayName != nil {
			name = *item.DisplayName
		}
		if item.CidrBlock != nil {
			cidr = *item.CidrBlock
		}
		if item.RouteTableId != nil {
			rtID = *item.RouteTableId
		}
		public := item.ProhibitPublicIpOnVnic == nil || !*item.ProhibitPublicIpOnVnic
		var slIDs []string
		if item.SecurityListIds != nil {
			slIDs = item.SecurityListIds
		}
		subnets = append(subnets, domain.Subnet{OCID: id, DisplayName: name, LifecycleState: string(item.LifecycleState), CidrBlock: cidr, Public: public, RouteTableID: rtID, SecurityListIDs: slIDs})
	}
	return subnets, nil
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

func toDomainVCNModel(v core.Vcn) domain.VCN {
	return domain.VCN{
		OCID:           *v.Id,
		DisplayName:    *v.DisplayName,
		LifecycleState: string(v.LifecycleState),
		CompartmentID:  *v.CompartmentId,
		DnsLabel:       *v.DnsLabel,
		DomainName:     *v.VcnDomainName,
		CidrBlocks:     cloneStrings(v.CidrBlocks),
		Ipv6Enabled:    len(v.Ipv6CidrBlocks) > 0,
		DhcpOptionsID:  *v.DefaultDhcpOptionsId,
		TimeCreated:    v.TimeCreated.Time,
		FreeformTags:   v.FreeformTags,
		DefinedTags:    v.DefinedTags,
	}
}

func cloneStrings(in []string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
