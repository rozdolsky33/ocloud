package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/identitydomains"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

type DynamicGroupAttributes struct {
	OCID           *string
	Name           *string
	Description    *string
	MatchingRule   *string
	LifecycleState string
	TimeCreated    *time.Time
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
	DomainURL      string
}

func NewDynamicGroupAttributesFromOCI(dg identity.DynamicGroup) *DynamicGroupAttributes {
	return &DynamicGroupAttributes{
		OCID:           dg.Id,
		Name:           dg.Name,
		Description:    dg.Description,
		MatchingRule:   dg.MatchingRule,
		LifecycleState: string(dg.LifecycleState),
		TimeCreated:    &dg.TimeCreated.Time,
		FreeformTags:   dg.FreeformTags,
		DefinedTags:    dg.DefinedTags,
	}
}

func NewDynamicGroupAttributesFromIDCS(dg identitydomains.DynamicResourceGroup, domainURL string) *DynamicGroupAttributes {
	var timeCreated *time.Time
	if dg.Meta != nil && dg.Meta.Created != nil {
		t, err := time.Parse(time.RFC3339, *dg.Meta.Created)
		if err == nil {
			timeCreated = &t
		}
	}

	id := dg.Id
	if dg.Ocid != nil {
		id = dg.Ocid
	}

	return &DynamicGroupAttributes{
		OCID:         id,
		Name:         dg.DisplayName,
		Description:  dg.Description,
		MatchingRule: dg.MatchingRule,
		TimeCreated:  timeCreated,
		DomainURL:    domainURL,
	}
}

func NewDomainDynamicGroupFromAttrs(dg *DynamicGroupAttributes) *domain.DynamicGroup {
	var name, ocid, description, matchingRule string
	var timeCreated time.Time

	if dg.Name != nil {
		name = *dg.Name
	}
	if dg.OCID != nil {
		ocid = *dg.OCID
	}
	if dg.Description != nil {
		description = *dg.Description
	}
	if dg.MatchingRule != nil {
		matchingRule = *dg.MatchingRule
	}
	if dg.TimeCreated != nil {
		timeCreated = *dg.TimeCreated
	}

	return &domain.DynamicGroup{
		OCID:           ocid,
		Name:           name,
		Description:    description,
		MatchingRule:   matchingRule,
		LifecycleState: dg.LifecycleState,
		TimeCreated:    timeCreated,
		FreeformTags:   dg.FreeformTags,
		DefinedTags:    dg.DefinedTags,
	}
}
