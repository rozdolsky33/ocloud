package subnet

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
	domainsubnet "github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
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
func (a *Adapter) GetSubnet(ctx context.Context, ocid string) (*domainsubnet.Subnet, error) {
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
func (a *Adapter) ListSubnets(ctx context.Context, compartmentID string) ([]domainsubnet.Subnet, error) {
	var subnets []domainsubnet.Subnet
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

// toDomainModel converts an OCI SDK subnet object to our application domain model.
func (a *Adapter) toDomainModel(s core.Subnet) domainsubnet.Subnet {
	var routeTableID string
	if s.RouteTableId != nil {
		routeTableID = *s.RouteTableId
	}
	var cidr string
	if s.CidrBlock != nil {
		cidr = *s.CidrBlock
	}
	var displayName string
	if s.DisplayName != nil {
		displayName = *s.DisplayName
	}
	var ocid string
	if s.Id != nil {
		ocid = *s.Id
	}
	// Public is the inverse of ProhibitPublicIpOnVnic
	public := s.ProhibitPublicIpOnVnic == nil || !*s.ProhibitPublicIpOnVnic

	return domainsubnet.Subnet{
		OCID:            ocid,
		DisplayName:     displayName,
		LifecycleState:  string(s.LifecycleState),
		CidrBlock:       cidr,
		Public:          public,
		RouteTableID:    routeTableID,
		SecurityListIDs: s.SecurityListIds,
		NSGIDs:          nil,
	}
}
