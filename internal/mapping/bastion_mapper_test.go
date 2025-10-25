package mapping_test

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestNewBastionAttributesFromOCIBastionSummary_And_ToDomain(t *testing.T) {
	id := "ocid1.bastion.oc1..abcd1234"
	name := "my-bastion"
	bType := "STANDARD"
	compartmentID := "ocid1.compartment.oc1..xyz"
	vcnID := "ocid1.vcn.oc1..vcn123"
	subnetID := "ocid1.subnet.oc1..subnet456"
	state := bastion.BastionLifecycleStateActive

	ociBastion := bastion.BastionSummary{
		Id:             &id,
		Name:           &name,
		BastionType:    &bType,
		CompartmentId:  &compartmentID,
		TargetVcnId:    &vcnID,
		TargetSubnetId: &subnetID,
		LifecycleState: state,
		FreeformTags:   map[string]string{"env": "prod"},
		DefinedTags:    map[string]map[string]interface{}{"ns": {"key": "value"}},
	}

	attrs := mapping.NewBastionAttributesFromOCIBastionSummary(ociBastion)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &bType, attrs.BastionType)
	require.Equal(t, state, attrs.LifecycleState)
	require.Equal(t, &compartmentID, attrs.CompartmentID)
	require.Equal(t, &vcnID, attrs.TargetVcnID)
	require.Equal(t, &subnetID, attrs.TargetSubnetID)
	require.Nil(t, attrs.MaxSessionTTL) // Not available in BastionSummary
	require.Equal(t, map[string]string{"env": "prod"}, attrs.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"key": "value"}}, attrs.DefinedTags)

	dom := mapping.NewDomainBastionFromAttrs(attrs)
	require.IsType(t, &domain.Bastion{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, bType, dom.BastionType)
	require.Equal(t, string(state), dom.LifecycleState)
	require.Equal(t, compartmentID, dom.CompartmentID)
	require.Equal(t, vcnID, dom.TargetVcnID)
	require.Equal(t, subnetID, dom.TargetSubnetID)
	require.Equal(t, 0, dom.MaxSessionTTL) // Nil pointer becomes 0
	require.Equal(t, map[string]string{"env": "prod"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"key": "value"}}, dom.DefinedTags)
}

func TestNewBastionAttributesFromOCIBastion_And_ToDomain(t *testing.T) {
	id := "ocid1.bastion.oc1..abcd1234"
	name := "my-bastion"
	bType := "STANDARD"
	compartmentID := "ocid1.compartment.oc1..xyz"
	vcnID := "ocid1.vcn.oc1..vcn123"
	subnetID := "ocid1.subnet.oc1..subnet456"
	maxTTL := 10800
	state := bastion.BastionLifecycleStateActive

	ociBastion := bastion.Bastion{
		Id:                     &id,
		Name:                   &name,
		BastionType:            &bType,
		CompartmentId:          &compartmentID,
		TargetVcnId:            &vcnID,
		TargetSubnetId:         &subnetID,
		MaxSessionTtlInSeconds: &maxTTL,
		LifecycleState:         state,
		FreeformTags:           map[string]string{"env": "prod"},
		DefinedTags:            map[string]map[string]interface{}{"ns": {"key": "value"}},
	}

	attrs := mapping.NewBastionAttributesFromOCIBastion(ociBastion)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &bType, attrs.BastionType)
	require.Equal(t, state, attrs.LifecycleState)
	require.Equal(t, &compartmentID, attrs.CompartmentID)
	require.Equal(t, &vcnID, attrs.TargetVcnID)
	require.Equal(t, &subnetID, attrs.TargetSubnetID)
	require.Equal(t, &maxTTL, attrs.MaxSessionTTL)

	dom := mapping.NewDomainBastionFromAttrs(attrs)
	require.IsType(t, &domain.Bastion{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, bType, dom.BastionType)
	require.Equal(t, string(state), dom.LifecycleState)
	require.Equal(t, compartmentID, dom.CompartmentID)
	require.Equal(t, vcnID, dom.TargetVcnID)
	require.Equal(t, subnetID, dom.TargetSubnetID)
	require.Equal(t, maxTTL, dom.MaxSessionTTL)
}

func TestNewDomainBastionFromAttrs_NilValues(t *testing.T) {
	// All pointer fields are nil; expect zero values in domain
	attrs := &mapping.BastionAttributes{}

	dom := mapping.NewDomainBastionFromAttrs(attrs)
	require.Equal(t, "", dom.OCID)
	require.Equal(t, "", dom.DisplayName)
	require.Equal(t, "", dom.BastionType)
	require.Equal(t, "", dom.LifecycleState)
	require.Equal(t, "", dom.CompartmentID)
	require.Equal(t, "", dom.TargetVcnID)
	require.Equal(t, "", dom.TargetSubnetID)
	require.Equal(t, 0, dom.MaxSessionTTL)
	require.Nil(t, dom.FreeformTags)
	require.Nil(t, dom.DefinedTags)
}

func TestNewDomainBastionFromAttrs_PartialValues(t *testing.T) {
	// Test with some nil and some non-nil values
	id := "ocid1.bastion.oc1..partial"
	name := "partial-bastion"
	// Other fields remain nil

	attrs := &mapping.BastionAttributes{
		OCID:        &id,
		DisplayName: &name,
	}

	dom := mapping.NewDomainBastionFromAttrs(attrs)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, "", dom.BastionType)
	require.Equal(t, "", dom.LifecycleState)
	require.Equal(t, "", dom.CompartmentID)
	require.Equal(t, "", dom.TargetVcnID)
	require.Equal(t, "", dom.TargetSubnetID)
	require.Equal(t, 0, dom.MaxSessionTTL)
}
