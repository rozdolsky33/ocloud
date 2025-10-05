package mapping

import (
	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

type CompartmentAttributes struct {
	OCID           *string
	Name           *string
	Description    *string
	LifecycleState identity.CompartmentLifecycleStateEnum
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
}

func NewCompartmentAttributesFromOCICompartment(c identity.Compartment) *CompartmentAttributes {
	return &CompartmentAttributes{
		OCID:           c.Id,
		Name:           c.Name,
		Description:    c.Description,
		LifecycleState: c.LifecycleState,
		FreeformTags:   c.FreeformTags,
		DefinedTags:    c.DefinedTags,
	}
}

func NewDomainCompartmentFromAttrs(c *CompartmentAttributes) *domain.Compartment {
	var ocid, name, description, lifecycleState string
	if c.OCID != nil {
		ocid = *c.OCID
	}
	if c.Name != nil {
		name = *c.Name
	}
	if c.Description != nil {
		description = *c.Description
	}
	if c.LifecycleState != "" {
		lifecycleState = string(c.LifecycleState)
	}

	return &domain.Compartment{
		OCID:           ocid,
		DisplayName:    name,
		Description:    description,
		LifecycleState: lifecycleState,
		FreeformTags:   c.FreeformTags,
		DefinedTags:    c.DefinedTags,
	}
}
