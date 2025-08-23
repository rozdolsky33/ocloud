package subnet

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// Adapter is an infrastructure-layer adapter for network subnets.
type Adapter struct {
	client core.VirtualNetworkClient
}

// NewAdapter creates a new subnet adapter.
func NewAdapter(client core.VirtualNetworkClient) *Adapter {
	return &Adapter{client: client}
}

// GetSubnet retrieves a single subnet by its OCID.
func (a *Adapter) GetSubnet(ctx context.Context, ocid string) (*domain.Subnet, error) {
	resp, err := a.client.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("getting subnet from OCI: %w", err)
	}

	sub := a.toDomainModel(resp.Subnet)
	return &sub,
		nil
}

// ListSubnets fetches all subnets in a compartment.
func (a *Adapter) ListSubnets(ctx context.Context, compartmentID string) ([]domain.Subnet, error) {
	var subnets []domain.Subnet
	var page *string

	for {
		resp, err := a.client.ListSubnets(ctx, core.ListSubnetsRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing subnets from OCI: %w", err)
		}

		for _, item := range resp.Items {
			subnets = append(subnets, a.toDomainModel(item))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return subnets, nil
}

// toDomainModel converts an OCI SDK subnet object to our application's domain model.
func (a *Adapter) toDomainModel(s core.Subnet) domain.Subnet {
	return domain.Subnet{
		OCID:                    *s.Id,
		DisplayName:             *s.DisplayName,
		CIDRBlock:               *s.CidrBlock,
		VcnOCID:                 *s.VcnId,
		RouteTableOCID:          *s.RouteTableId,
		SecurityListOCIDs:       s.SecurityListIds,
		DhcpOptionsOCID:         *s.DhcpOptionsId,
		ProhibitPublicIPOnVnic:  *s.ProhibitPublicIpOnVnic,
		ProhibitInternetIngress: *s.ProhibitInternetIngress,
		DNSLabel:                *s.DnsLabel,
		SubnetDomainName:        *s.SubnetDomainName,
	}
}
