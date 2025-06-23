package compartment

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	return &Service{
		identityClient: appCtx.IdentityClient,
		logger:         appCtx.Logger,
		TenancyID:      appCtx.TenancyID,
		TenancyName:    appCtx.TenancyName,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]Compartment, error) {
	logger.LogWithLevel(s.logger, 3, "Listing compartments in tenancy", "tenancyName: ", s.TenancyName, "tenancyID: ", s.TenancyID)

	// prepare the base request
	req := identity.ListCompartmentsRequest{
		CompartmentId:          &s.TenancyID,
		AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
		LifecycleState:         identity.CompartmentLifecycleStateActive,
		CompartmentIdInSubtree: common.Bool(true),
	}

	var compartments []Compartment
	// paginate through results; stop when OpcNextPage is nil
	pageToken := ""
	for {
		if pageToken != "" {
			req.Page = common.String(pageToken)
		}

		resp, err := s.identityClient.ListCompartments(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing compartments: %w", err)
		}

		// scan each compartment summary for a name match
		for _, comp := range resp.Items {
			compartment := mapToCompartment(comp)
			compartments = append(compartments, compartment)

		}

		// if there's no next page, we're done searching
		if resp.OpcNextPage == nil {
			break
		}
		pageToken = *resp.OpcNextPage
	}

	return compartments, nil

}

func mapToCompartment(compartment identity.Compartment) Compartment {
	return Compartment{
		Name:        *compartment.Name,
		ID:          *compartment.Id,
		Description: *compartment.Description,
	}
}
