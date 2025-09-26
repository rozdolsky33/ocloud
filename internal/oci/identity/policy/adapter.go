package policy

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

type Adapter struct {
	identityClient identity.IdentityClient
}

// NewAdapter creates a new adapter instance.
func NewAdapter(identityClient identity.IdentityClient) *Adapter {
	return &Adapter{identityClient: identityClient}
}

// GetPolicy retrieves a single policy by its OCID.
func (a *Adapter) GetPolicy(ctx context.Context, ocid string) (*domain.Policy, error) {
	resp, err := a.identityClient.GetPolicy(ctx, identity.GetPolicyRequest{
		PolicyId: &ocid,
	})
	if err != nil {

	}
	policy := a.toDomainModel(resp.Policy)
	return &policy, nil
}

// ListPolicies retrieves all policies in a given compartment.
func (a *Adapter) ListPolicies(ctx context.Context, compartmentID string) ([]domain.Policy, error) {
	var policies []domain.Policy
	page := ""
	for {
		resp, err := a.identityClient.ListPolicies(ctx, identity.ListPoliciesRequest{
			CompartmentId: &compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list policies: %w", err)
		}
		for _, item := range resp.Items {
			policies = append(policies, a.toDomainModel(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return policies, nil
}

// toDomainModel converts an OCI SDK policy object to our application's domain model.'
func (a *Adapter) toDomainModel(p identity.Policy) domain.Policy {
	return domain.Policy{
		Name:         *p.Name,
		ID:           *p.Id,
		Statement:    p.Statements,
		TimeCreated:  p.TimeCreated.Time,
		Description:  *p.Description,
		FreeformTags: p.FreeformTags,
		DefinedTags:  p.DefinedTags,
	}
}
