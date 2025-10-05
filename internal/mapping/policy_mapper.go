package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

type PolicyAttributes struct {
	Name         *string
	ID           *string
	Statement    []string
	Description  *string
	TimeCreated  *time.Time
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

func NewPolicyAttributesFromOCIPolicy(p identity.Policy) *PolicyAttributes {
	return &PolicyAttributes{
		Name:         p.Name,
		ID:           p.Id,
		Statement:    p.Statements,
		Description:  p.Description,
		TimeCreated:  &p.TimeCreated.Time,
		FreeformTags: p.FreeformTags,
		DefinedTags:  p.DefinedTags,
	}
}

func NewDomainPolicyFromAttrs(p *PolicyAttributes) *domain.Policy {
	var name, id, description string
	var timeCreated time.Time

	if p.Name != nil {
		name = *p.Name
	}
	if p.ID != nil {
		id = *p.ID
	}
	if p.Description != nil {
		description = *p.Description
	}
	if p.TimeCreated != nil {
		timeCreated = *p.TimeCreated
	}

	return &domain.Policy{
		Name:         name,
		ID:           id,
		Statement:    p.Statement,
		Description:  description,
		TimeCreated:  timeCreated,
		FreeformTags: p.FreeformTags,
		DefinedTags:  p.DefinedTags,
	}
}
