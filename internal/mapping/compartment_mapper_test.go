package mapping_test

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestNewCompartmentAttributesFromOCICompartment_And_ToDomain(t *testing.T) {
	name := "Dev"
	id := "ocid1.compartment.oc1..abcd"
	desc := "Development compartment"
	state := identity.CompartmentLifecycleStateActive

	ocic := identity.Compartment{
		Id:             &id,
		Name:           &name,
		Description:    &desc,
		LifecycleState: state,
		FreeformTags:   map[string]string{"env": "dev"},
		DefinedTags:    map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewCompartmentAttributesFromOCICompartment(ocic)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.Name)
	require.Equal(t, &desc, attrs.Description)
	require.Equal(t, state, attrs.LifecycleState)
	require.Equal(t, map[string]string{"env": "dev"}, attrs.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, attrs.DefinedTags)

	dom := mapping.NewDomainCompartmentFromAttrs(attrs)
	require.IsType(t, &domain.Compartment{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, desc, dom.Description)
	require.Equal(t, string(state), dom.LifecycleState)
	require.Equal(t, map[string]string{"env": "dev"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, dom.DefinedTags)
}

func TestNewDomainCompartmentFromAttrs_NilValues(t *testing.T) {
	// All pointer fields are nil; expect zero values in domain
	attrs := &mapping.CompartmentAttributes{}

	dom := mapping.NewDomainCompartmentFromAttrs(attrs)
	require.Equal(t, "", dom.OCID)
	require.Equal(t, "", dom.DisplayName)
	require.Equal(t, "", dom.Description)
	require.Equal(t, "", dom.LifecycleState)
	require.Nil(t, dom.FreeformTags)
	require.Nil(t, dom.DefinedTags)
}
