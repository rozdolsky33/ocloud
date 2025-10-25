package bastion

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
)

// Adapter is an infrastructure-layer adapter that implements the domain.BastionRepository interface.
type Adapter struct {
	bastionClient bastion.BastionClient
	networkClient core.VirtualNetworkClient
	compartmentID string
	vcnCache      map[string]*core.Vcn
	subnetCache   map[string]*core.Subnet
}

// NewBastionAdapter creates a new adapter for interacting with OCI Bastions.
func NewBastionAdapter(bastionClient bastion.BastionClient, networkClient core.VirtualNetworkClient, compartmentID string) *Adapter {
	return &Adapter{
		bastionClient: bastionClient,
		networkClient: networkClient,
		compartmentID: compartmentID,
		vcnCache:      make(map[string]*core.Vcn),
		subnetCache:   make(map[string]*core.Subnet),
	}
}

// ListBastions retrieves all bastions in the configured compartment.
// It enriches the bastion data with VCN and Subnet names.
func (a *Adapter) ListBastions(ctx context.Context, compartmentID string) ([]domain.Bastion, error) {
	var bastions []domain.Bastion
	var page *string

	for {
		req := bastion.ListBastionsRequest{
			CompartmentId: &compartmentID,
		}
		if page != nil {
			req.Page = page
		}

		resp, err := a.bastionClient.ListBastions(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing bastions from OCI: %w", err)
		}

		for _, item := range resp.Items {
			// Fetch full bastion details to get all fields including MaxSessionTTL,
			// ClientCidrBlockAllowList, and PrivateEndpointIpAddress
			var domainBastion *domain.Bastion
			if item.Id != nil {
				fullBastion, err := a.GetBastion(ctx, *item.Id)
				if err == nil {
					domainBastion = fullBastion
				} else {
					// Fallback to summary if full fetch fails
					attrs := mapping.NewBastionAttributesFromOCIBastionSummary(item)
					domainBastion = mapping.NewDomainBastionFromAttrs(attrs)
				}
			} else {
				// No ID, use summary
				attrs := mapping.NewBastionAttributesFromOCIBastionSummary(item)
				domainBastion = mapping.NewDomainBastionFromAttrs(attrs)
			}

			bastions = append(bastions, *domainBastion)
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return bastions, nil
}

// GetBastion retrieves a specific bastion by ID.
// It enriches the bastion data with VCN and Subnet names.
func (a *Adapter) GetBastion(ctx context.Context, bastionID string) (*domain.Bastion, error) {
	resp, err := a.bastionClient.GetBastion(ctx, bastion.GetBastionRequest{
		BastionId: &bastionID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting bastion from OCI: %w", err)
	}

	// Convert OCI bastion to domain bastion via mapping layer
	attrs := mapping.NewBastionAttributesFromOCIBastion(resp.Bastion)
	domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

	if resp.Bastion.TargetVcnId != nil && *resp.Bastion.TargetVcnId != "" {
		vcnName, err := a.getVcnName(ctx, *resp.Bastion.TargetVcnId)
		if err == nil {
			domainBastion.TargetVcnName = vcnName
		}
	}

	if resp.Bastion.TargetSubnetId != nil && *resp.Bastion.TargetSubnetId != "" {
		subnetName, err := a.getSubnetName(ctx, *resp.Bastion.TargetSubnetId)
		if err == nil {
			domainBastion.TargetSubnetName = subnetName
		}
	}

	return domainBastion, nil
}

// CreateBastion creates a new bastion.
func (a *Adapter) CreateBastion(ctx context.Context, request domain.CreateBastionRequest) (*domain.Bastion, error) {
	createReq := bastion.CreateBastionRequest{
		CreateBastionDetails: bastion.CreateBastionDetails{
			Name:                     &request.DisplayName,
			CompartmentId:            &request.CompartmentID,
			TargetSubnetId:           &request.TargetSubnetID,
			BastionType:              &request.BastionType,
			ClientCidrBlockAllowList: request.ClientCIDRList,
			FreeformTags:             request.FreeformTags,
			DefinedTags:              request.DefinedTags,
		},
	}

	if request.MaxSessionTTL > 0 {
		createReq.CreateBastionDetails.MaxSessionTtlInSeconds = &request.MaxSessionTTL
	}

	resp, err := a.bastionClient.CreateBastion(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("creating bastion in OCI: %w", err)
	}

	attrs := mapping.NewBastionAttributesFromOCIBastion(resp.Bastion)
	return mapping.NewDomainBastionFromAttrs(attrs), nil
}

// DeleteBastion deletes a bastion.
func (a *Adapter) DeleteBastion(ctx context.Context, bastionID string) error {
	_, err := a.bastionClient.DeleteBastion(ctx, bastion.DeleteBastionRequest{
		BastionId: &bastionID,
	})
	if err != nil {
		return fmt.Errorf("deleting bastion from OCI: %w", err)
	}
	return nil
}

// getVcnName retrieves the display name for a VCN, using cache when available.
func (a *Adapter) getVcnName(ctx context.Context, vcnID string) (string, error) {
	if vcn, ok := a.vcnCache[vcnID]; ok {
		if vcn.DisplayName != nil {
			return *vcn.DisplayName, nil
		}
		return "", nil
	}

	resp, err := a.networkClient.GetVcn(ctx, core.GetVcnRequest{
		VcnId: &vcnID,
	})
	if err != nil {
		return "", fmt.Errorf("getting VCN details: %w", err)
	}

	a.vcnCache[vcnID] = &resp.Vcn

	if resp.Vcn.DisplayName != nil {
		return *resp.Vcn.DisplayName, nil
	}
	return "", nil
}

// getSubnetName retrieves the display name for a Subnet, using cache when available.
func (a *Adapter) getSubnetName(ctx context.Context, subnetID string) (string, error) {
	if subnet, ok := a.subnetCache[subnetID]; ok {
		if subnet.DisplayName != nil {
			return *subnet.DisplayName, nil
		}
		return "", nil
	}

	resp, err := a.networkClient.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &subnetID,
	})
	if err != nil {
		return "", fmt.Errorf("getting Subnet details: %w", err)
	}

	a.subnetCache[subnetID] = &resp.Subnet

	if resp.Subnet.DisplayName != nil {
		return *resp.Subnet.DisplayName, nil
	}
	return "", nil
}
