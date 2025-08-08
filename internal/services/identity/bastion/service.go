package bastion

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/bastion"
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
	return &Service{
		bastionClient: bc,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}, nil
}

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
		allBastions = append(allBastions, mapToBastion(b))
	}

	return allBastions, nil
}

// GetDummyBastions returns a list of dummy bastion options
func (s *Service) GetDummyBastions() []Bastion {
	return []Bastion{
		{ID: "ocid1.bastion.oc1.dummy.bastion1", Name: "bastion_1"},
		{ID: "ocid1.bastion.oc1.dummy.bastion2", Name: "basstion_1"},
		{ID: "ocid1.bastion.oc1.dummy.bastion3", Name: "bestion three"},
	}
}

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
