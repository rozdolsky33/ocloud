package vcn

import (
	"context"
	"fmt"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/core"
	vcn2 "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
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

// GetVcn retrieves a single VCN by its OCID.
func (a *Adapter) GetVcn(ctx context.Context, vcnID string) (*vcn2.VCN, error) {
	resp, err := a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
	if err != nil {
		return nil, fmt.Errorf("getting VCN from OCI: %w", err)
	}
	m := toDomainVCNModel(resp.Vcn)
	return &m, nil
}

// ListVcns lists all VCNs in a given compartment.
func (a *Adapter) ListVcns(ctx context.Context, compartmentID string) ([]vcn2.VCN, error) {
	req := core.ListVcnsRequest{CompartmentId: &compartmentID}
	var out []vcn2.VCN
	for {
		resp, err := a.client.ListVcns(ctx, req)
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

// ListEnrichedVcns lists VCNs and enriches them with DHCP options in parallel.
func (a *Adapter) ListEnrichedVcns(ctx context.Context, compartmentID string) ([]vcn2.VCN, error) {
	vcns, err := a.ListVcns(ctx, compartmentID)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]struct{})
	for _, v := range vcns {
		if v.DhcpOptionsID != "" {
			ids[v.DhcpOptionsID] = struct{}{}
		}
	}
	if len(ids) == 0 {
		return vcns, nil
	}
	dhcpMap := make(map[string]vcn2.DhcpOptions, len(ids))
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(ids))
	for id := range ids {
		id := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			obj, err := a.GetDhcpOptions(ctx, id)
			if err != nil {
				errCh <- fmt.Errorf("get dhcp options %s: %w", id, err)
				return
			}
			mu.Lock()
			dhcpMap[id] = obj
			mu.Unlock()
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return vcns, err
		}
	}
	for i := range vcns {
		if obj, ok := dhcpMap[vcns[i].DhcpOptionsID]; ok {
			vcns[i].DhcpOptions = obj
		}
	}
	return vcns, nil
}

// GetDhcpOptions fetches DHCP options resource by OCID.
func (a *Adapter) GetDhcpOptions(ctx context.Context, dhcpID string) (vcn2.DhcpOptions, error) {
	resp, err := a.client.GetDhcpOptions(ctx, core.GetDhcpOptionsRequest{DhcpId: &dhcpID})
	if err != nil {
		return vcn2.DhcpOptions{}, fmt.Errorf("getting DHCP options from OCI: %w", err)
	}
	return toDomainDHCPOptionsModel(resp.DhcpOptions)
}

func toDomainVCNModel(v core.Vcn) vcn2.VCN {
	return vcn2.VCN{
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

func toDomainDHCPOptionsModel(d core.DhcpOptions) (vcn2.DhcpOptions, error) {
	if d.Id == nil {
		return vcn2.DhcpOptions{}, fmt.Errorf("dhcp options missing id")
	}
	return vcn2.DhcpOptions{
		OCID:           *d.Id,
		DisplayName:    *d.DisplayName,
		LifecycleState: string(d.LifecycleState),
		DomainNameType: string(d.DomainNameType),
	}, nil
}

func cloneStrings(in []string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
