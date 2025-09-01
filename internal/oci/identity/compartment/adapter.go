package compartment

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// CompartmentAdapter is an infrastructure-layer adapter that implements the domain.CompartmentRepository interface.
type CompartmentAdapter struct {
	client    identity.IdentityClient
	tenancyID string
}

// NewCompartmentAdapter creates a new adapter for interacting with OCI compartments.
func NewCompartmentAdapter(client identity.IdentityClient, tenancyID string) *CompartmentAdapter {
	return &CompartmentAdapter{
		client:    client,
		tenancyID: tenancyID,
	}
}

// GetCompartment retrieves a single compartment by its OCID.
func (a *CompartmentAdapter) GetCompartment(ctx context.Context, ocid string) (*domain.Compartment, error) {
	resp, err := a.client.GetCompartment(ctx, identity.GetCompartmentRequest{
		CompartmentId: &ocid,
	})
	if err != nil {
		return nil, fmt.Errorf("getting compartment from OCI: %w", err)
	}

	comp := a.toDomainModel(resp.Compartment)
	return &comp, nil
}

// ListCompartments retrieves all active compartments under a given parent compartment.
// It handles pagination to fetch all results from OCI.
func (a *CompartmentAdapter) ListCompartments(ctx context.Context, parentCompartmentID string) ([]domain.Compartment, error) {
	var compartments []domain.Compartment
	page := ""

	// If no parent is specified, use the tenancy root.
	if parentCompartmentID == "" {
		parentCompartmentID = a.tenancyID
	}

	for {
		resp, err := a.client.ListCompartments(ctx, identity.ListCompartmentsRequest{
			CompartmentId:          &parentCompartmentID,
			Page:                   &page,
			AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
			LifecycleState:         identity.CompartmentLifecycleStateActive,
			CompartmentIdInSubtree: common.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("listing compartments from OCI: %w", err)
		}

		for _, item := range resp.Items {
			compartments = append(compartments, a.toDomainModel(item))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return compartments, nil
}

// toDomainModel converts an OCI SDK compartment object to our application's domain model.
// This is the anti-corruption layer in action.
func (a *CompartmentAdapter) toDomainModel(c identity.Compartment) domain.Compartment {
	var parentID string
	if c.CompartmentId != nil {
		parentID = *c.CompartmentId
	}

	var state string
	if c.LifecycleState != "" {
		state = string(c.LifecycleState)
	}

	return domain.Compartment{
		OCID:                *c.Id,
		DisplayName:         *c.Name,
		Description:         *c.Description,
		LifecycleState:      state,
		ParentCompartmentID: parentID,
	}
}
