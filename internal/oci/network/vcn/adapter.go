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

// GetVcn fetches VCN details by OCID.
func (a *Adapter) GetVcn(ctx context.Context, vcnID string) (*domain.VCN, error) {
	resp, err := a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
	if err != nil {
		return nil, fmt.Errorf("getting VCN from OCI: %w", err)
	}
	return toVCNModel(&resp.Vcn), nil
}

// ListVcns lists all VCNs in a compartment.
func (a *Adapter) ListVcns(ctx context.Context, compartmentID string) ([]*domain.VCN, error) {
	var vcns []*domain.VCN
	req := core.ListVcnsRequest{
		CompartmentId: &compartmentID,
	}

	for {
		resp, err := a.client.ListVcns(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing VCNs from OCI: %w", err)
		}

		for _, v := range resp.Items {
			vcns = append(vcns, toVCNModel(&v))
		}

		if resp.OpcNextPage == nil {
			break
		}
		req.Page = resp.OpcNextPage
	}

	return vcns, nil
}

// GetDhcpOptions fetches DHCP options resource by OCID.
func (a *Adapter) GetDhcpOptions(ctx context.Context, dhcpID string) (*domain.DhcpOptions, error) {
	resp, err := a.client.GetDhcpOptions(ctx, core.GetDhcpOptionsRequest{DhcpId: &dhcpID})
	if err != nil {
		return nil, fmt.Errorf("getting DHCP options from OCI: %w", err)
	}
	return toDHCPOptionsModel(&resp.DhcpOptions), nil
}

func toVCNModel(v *core.Vcn) *domain.VCN {
	m := &domain.VCN{}
	if v == nil {
		return m
	}
	if v.Id != nil {
		m.OCID = *v.Id
	}
	if v.DisplayName != nil {
		m.DisplayName = *v.DisplayName
	}
	if v.LifecycleState != "" {
		m.LifecycleState = string(v.LifecycleState)
	}
	if v.CompartmentId != nil {
		m.CompartmentID = *v.CompartmentId
	}
	if v.DnsLabel != nil {
		m.DnsLabel = *v.DnsLabel
	}
	if v.VcnDomainName != nil {
		m.DomainName = *v.VcnDomainName
	}
	if v.CidrBlocks != nil {
		m.CidrBlocks = make([]string, len(v.CidrBlocks))
		copy(m.CidrBlocks, v.CidrBlocks)
	}
	if v.Ipv6CidrBlocks != nil && len(v.Ipv6CidrBlocks) > 0 {
		m.Ipv6Enabled = true
	}
	if v.DefaultDhcpOptionsId != nil {
		m.DhcpOptionsID = *v.DefaultDhcpOptionsId
	}
	if v.TimeCreated != nil {
		m.TimeCreated = v.TimeCreated.Time
	}
	if v.FreeformTags != nil {
		m.FreeformTags = v.FreeformTags
	}
	return m
}

func toDHCPOptionsModel(d *core.DhcpOptions) *domain.DhcpOptions {
	m := &domain.DhcpOptions{}
	if d == nil {
		return m
	}
	if d.Id != nil {
		m.OCID = *d.Id
	}
	if d.DisplayName != nil {
		m.DisplayName = *d.DisplayName
	}
	if d.LifecycleState != "" {
		m.LifecycleState = string(d.LifecycleState)
	}
	return m
}
