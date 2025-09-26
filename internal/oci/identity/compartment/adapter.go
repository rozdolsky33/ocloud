package compartment

import (
	"context"
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

// Adapter is an infrastructure-layer adapter that implements the domain.CompartmentRepository interface.
type Adapter struct {
	client          identity.IdentityClient
	compartmentOCID string
}

// NewCompartmentAdapter creates a new adapter for interacting with OCI compartments.
func NewCompartmentAdapter(client identity.IdentityClient, ocid string) *Adapter {
	return &Adapter{
		client:          client,
		compartmentOCID: ocid,
	}
}

// GetCompartment retrieves a single compartment by its OCID.
func (a *Adapter) GetCompartment(ctx context.Context, ocid string) (*domain.Compartment, error) {
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
func (a *Adapter) ListCompartments(ctx context.Context, ocid string) ([]domain.Compartment, error) {
	var compartments []domain.Compartment
	page := ""

	for {
		includeSubtree := strings.HasPrefix(ocid, "ocid1.tenancy.")
		resp, err := a.client.ListCompartments(ctx, identity.ListCompartmentsRequest{
			CompartmentId:          &ocid,
			Page:                   &page,
			AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
			LifecycleState:         identity.CompartmentLifecycleStateActive,
			CompartmentIdInSubtree: common.Bool(includeSubtree),
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
func (a *Adapter) toDomainModel(c identity.Compartment) domain.Compartment {
	var state string
	if c.LifecycleState != "" {
		state = string(c.LifecycleState)
	}
	return domain.Compartment{
		OCID:           *c.Id,
		DisplayName:    *c.Name,
		Description:    *c.Description,
		LifecycleState: state,
		FreeformTags:   c.FreeformTags,
		DefinedTags:    c.DefinedTags,
	}
}
