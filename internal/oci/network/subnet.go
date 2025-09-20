package network

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
)

// SubnetRepository implements the domain SubnetRepository interface.

func (c *Client) GetSubnet(ctx context.Context, ocid string) (*subnet.Subnet, error) {
	resp, err := c.vnClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &ocid})
	if err != nil {
		return nil, fmt.Errorf("getting subnet: %w", err)
	}

	s := &subnet.Subnet{
		OCID:           *resp.Id,
		DisplayName:    *resp.DisplayName,
		LifecycleState: string(resp.LifecycleState),
		CidrBlock:      *resp.CidrBlock,
		Public:         !*resp.ProhibitPublicIpOnVnic,
		RouteTableID:   *resp.RouteTableId,
	}

	return s, nil
}

func (c *Client) ListSubnets(ctx context.Context, compartmentID string) ([]subnet.Subnet, error) {
	var subnets []subnet.Subnet
	var page *string

	for {
		req := core.ListSubnetsRequest{
			CompartmentId: &compartmentID,
			Limit:         common.Int(100),
			Page:          page,
		}

		resp, err := c.vnClient.ListSubnets(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing subnets: %w", err)
		}

		for _, s := range resp.Items {
			subnets = append(subnets, subnet.Subnet{
				OCID:           *s.Id,
				DisplayName:    *s.DisplayName,
				LifecycleState: string(s.LifecycleState),
				CidrBlock:      *s.CidrBlock,
				Public:         !*s.ProhibitPublicIpOnVnic,
				RouteTableID:   *s.RouteTableId,
			})
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return subnets, nil
}
