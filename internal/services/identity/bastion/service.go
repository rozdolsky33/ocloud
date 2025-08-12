package bastion

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// NewService creates a new bastion service
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	bc, err := oci.NewBastionClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create bastion client: %w", err)
	}
	nc, err := oci.NewNetworkClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}
	cc, err := oci.NewComputeClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}
	return &Service{
		bastionClient: bc,
		networkClient: nc,
		computeClient: cc,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
		vcnCache:      make(map[string]*core.Vcn),
		subnetCache:   make(map[string]*core.Subnet),
	}, nil
}

// List retrieves and returns all bastion hosts from the given compartment in the OCI account.
func (s *Service) List(ctx context.Context) (bastions []Bastion, err error) {
	logger.LogWithLevel(s.logger, 1, "Listing Bastions in compartment", "compartmentID", s.compartmentID)
	request := bastion.ListBastionsRequest{
		CompartmentId: &s.compartmentID,
	}
	response, err := s.bastionClient.ListBastions(ctx, request)

	if err != nil {
		return nil, fmt.Errorf("failed to list bastions: %w", err)
	}
	var allBastions []Bastion

	for _, b := range response.Items {
		toBastion := mapToBastion(b)
		if b.TargetVcnId != nil && *b.TargetVcnId != "" {
			vcn, err := s.fetchVcnDetails(ctx, *b.TargetVcnId)
			if err != nil {
				logger.LogWithLevel(s.logger, 2, "Failed to fetch VCN details", "vcnID", *b.TargetVcnId, "error", err)
			} else if vcn.DisplayName != nil {
				toBastion.TargetVcnName = *vcn.DisplayName
			}
		}

		// Fetch Subnet details
		if b.TargetSubnetId != nil && *b.TargetSubnetId != "" {
			subnet, err := s.fetchSubnetDetails(ctx, *b.TargetSubnetId)
			if err != nil {
				logger.LogWithLevel(s.logger, 2, "Failed to fetch Subnet details", "subnetID", *b.TargetSubnetId, "error", err)
			} else if subnet.DisplayName != nil {
				toBastion.TargetSubnetName = *subnet.DisplayName
			}
		}

		allBastions = append(allBastions, toBastion)
	}

	return allBastions, nil
}

// fetchVcnDetails retrieves the VCN details for the given VCN ID.
func (s *Service) fetchVcnDetails(ctx context.Context, vcnID string) (*core.Vcn, error) {

	if vcn, ok := s.vcnCache[vcnID]; ok {
		logger.LogWithLevel(s.logger, 3, "VCN cache hit", "vcnID", vcnID)
		return vcn, nil
	}

	logger.LogWithLevel(s.logger, 3, "VCN cache miss", "vcnID", vcnID)
	resp, err := s.networkClient.GetVcn(ctx, core.GetVcnRequest{
		VcnId: &vcnID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting VCN details: %w", err)
	}

	s.vcnCache[vcnID] = &resp.Vcn
	return &resp.Vcn, nil
}

// fetchSubnetDetails retrieves the subnet details for the given subnet ID.
// It uses a cache to avoid making repeated API calls for the same subnet.
func (s *Service) fetchSubnetDetails(ctx context.Context, subnetID string) (*core.Subnet, error) {
	if subnet, ok := s.subnetCache[subnetID]; ok {
		logger.LogWithLevel(s.logger, 3, "subnet cache hit", "subnetID", subnetID)
		return subnet, nil
	}

	// Cache miss, fetch from API
	logger.LogWithLevel(s.logger, 3, "subnet cache miss", "subnetID", subnetID)
	resp, err := s.networkClient.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &subnetID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting subnet details: %w", err)
	}

	s.subnetCache[subnetID] = &resp.Subnet
	return &resp.Subnet, nil
}

// mapToBastion converts a BastionSummary object to a Bastion object with relevant fields populated.
func mapToBastion(bastion bastion.BastionSummary) Bastion {
	return Bastion{
		ID:             *bastion.Id,
		Name:           *bastion.Name,
		BastionType:    *bastion.BastionType,
		LifecycleState: bastion.LifecycleState,
		TargetVcnId:    *bastion.TargetVcnId,
		TargetSubnetId: *bastion.TargetSubnetId,
	}
}
