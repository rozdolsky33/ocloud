package vcn

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/domain"

	"github.com/oracle/oci-go-sdk/v65/core"
)

// Adapter provides access to VCN-related OCI APIs.
// It is infra-layer and should be used by the service layer.
type Adapter struct {
	client core.VirtualNetworkClient
}

// NewAdapter creates a new VCN adapter using the provided configuration provider.
func NewAdapter(client core.VirtualNetworkClient) *Adapter {
	return &Adapter{client: client}
}

// GetVcn retrieves a single VCN by its OCID.
func (a *Adapter) GetVcn(ctx context.Context, vcnID string) (*domain.VCN, error) {
	resp, err := a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
	if err != nil {
		return nil, fmt.Errorf("getting VCN from OCI: %w", err)
	}
	//toDomainVCNModel(resp)
	fmt.Println(resp)

	return nil, nil
}

// ListVcns lists all VCNs in a compartment.
func (a *Adapter) ListVcns(ctx context.Context, compartmentID string) ([]domain.VCN, error) {
	var vcns []domain.VCN
	req := core.ListVcnsRequest{
		CompartmentId: &compartmentID,
	}

	for {
		resp, err := a.client.ListVcns(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing VCNs from OCI: %w", err)
		}

		for _, v := range resp.Items {
			vcns = append(vcns, toDomainVCNModel(v))
		}

		if resp.OpcNextPage == nil {
			break
		}
		req.Page = resp.OpcNextPage
	}

	return vcns, nil
}

// GetDhcpOptions fetches DHCP options resource by OCID.
func (a *Adapter) GetDhcpOptions(ctx context.Context, dhcpID string) (domain.DhcpOptions, error) {
	resp, err := a.client.GetDhcpOptions(ctx, core.GetDhcpOptionsRequest{DhcpId: &dhcpID})
	if err != nil {
		return domain.DhcpOptions{}, fmt.Errorf("getting DHCP options from OCI: %w", err)
	}
	return toDomainDHCPOptionsModel(&resp.DhcpOptions), nil
}

// Map a *core.Vcn -> domain.VCN safely.
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

func toDomainDHCPOptionsModel(d *core.DhcpOptions) domain.DhcpOptions {
	return domain.DhcpOptions{
		OCID:           *d.Id,
		DisplayName:    *d.DisplayName,
		LifecycleState: string(d.LifecycleState),
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
