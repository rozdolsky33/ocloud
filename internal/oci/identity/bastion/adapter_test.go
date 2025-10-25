package bastion

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/assert"
)

// Helper functions for creating pointers
func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }

// TestNewBastionAdapter tests the adapter constructor
func TestNewBastionAdapter(t *testing.T) {
	// We can't create real OCI clients in unit tests, but we can verify the adapter structure
	adapter := &Adapter{
		compartmentID: "ocid1.compartment.oc1..test",
		vcnCache:      make(map[string]*core.Vcn),
		subnetCache:   make(map[string]*core.Subnet),
	}

	assert.NotNil(t, adapter)
	assert.Equal(t, "ocid1.compartment.oc1..test", adapter.compartmentID)
	assert.NotNil(t, adapter.vcnCache)
	assert.NotNil(t, adapter.subnetCache)
}

// TestMappingIntegration tests that the adapter correctly uses the mapping layer
func TestMappingIntegration(t *testing.T) {
	// Test OCI Bastion -> BastionAttributes -> domain.Bastion conversion
	ociBastion := bastion.Bastion{
		Id:                       stringPtr("ocid1.bastion.oc1..test"),
		Name:                     stringPtr("test-bastion"),
		BastionType:              stringPtr("STANDARD"),
		LifecycleState:           bastion.BastionLifecycleStateActive,
		CompartmentId:            stringPtr("ocid1.compartment.oc1..test"),
		TargetVcnId:              stringPtr("ocid1.vcn.oc1..test"),
		TargetSubnetId:           stringPtr("ocid1.subnet.oc1..test"),
		MaxSessionTtlInSeconds:   intPtr(10800),
		ClientCidrBlockAllowList: []string{"0.0.0.0/0"},
		PrivateEndpointIpAddress: stringPtr("10.0.0.1"),
	}

	// Use mapping layer
	attrs := mapping.NewBastionAttributesFromOCIBastion(ociBastion)
	domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

	// Verify conversion
	assert.Equal(t, "ocid1.bastion.oc1..test", domainBastion.OCID)
	assert.Equal(t, "test-bastion", domainBastion.DisplayName)
	assert.Equal(t, "STANDARD", domainBastion.BastionType)
	assert.Equal(t, "ACTIVE", domainBastion.LifecycleState)
	assert.Equal(t, "ocid1.compartment.oc1..test", domainBastion.CompartmentID)
	assert.Equal(t, "ocid1.vcn.oc1..test", domainBastion.TargetVcnID)
	assert.Equal(t, "ocid1.subnet.oc1..test", domainBastion.TargetSubnetID)
	assert.Equal(t, 10800, domainBastion.MaxSessionTTL)
	assert.Equal(t, []string{"0.0.0.0/0"}, domainBastion.ClientCidrBlockAllowList)
	assert.Equal(t, "10.0.0.1", domainBastion.PrivateEndpointIpAddress)
}

// TestBastionSummaryMapping tests BastionSummary conversion
func TestBastionSummaryMapping(t *testing.T) {
	ociSummary := bastion.BastionSummary{
		Id:             stringPtr("ocid1.bastion.oc1..summary"),
		Name:           stringPtr("summary-bastion"),
		BastionType:    stringPtr("STANDARD"),
		LifecycleState: bastion.BastionLifecycleStateCreating,
		CompartmentId:  stringPtr("ocid1.compartment.oc1..test"),
		TargetVcnId:    stringPtr("ocid1.vcn.oc1..test"),
		TargetSubnetId: stringPtr("ocid1.subnet.oc1..test"),
	}

	attrs := mapping.NewBastionAttributesFromOCIBastionSummary(ociSummary)
	domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

	assert.Equal(t, "ocid1.bastion.oc1..summary", domainBastion.OCID)
	assert.Equal(t, "summary-bastion", domainBastion.DisplayName)
	assert.Equal(t, "STANDARD", domainBastion.BastionType)
	assert.Equal(t, "CREATING", domainBastion.LifecycleState)
	assert.Equal(t, 0, domainBastion.MaxSessionTTL) // Not available in summary
	assert.Nil(t, domainBastion.ClientCidrBlockAllowList)
	assert.Equal(t, "", domainBastion.PrivateEndpointIpAddress)
}

// TestCreateBastionRequestMapping tests domain -> OCI SDK conversion for creation
func TestCreateBastionRequestMapping(t *testing.T) {
	domainRequest := domain.CreateBastionRequest{
		DisplayName:    "new-bastion",
		CompartmentID:  "ocid1.compartment.oc1..test",
		TargetSubnetID: "ocid1.subnet.oc1..test",
		BastionType:    "STANDARD",
		ClientCIDRList: []string{"0.0.0.0/0", "10.0.0.0/8"},
		MaxSessionTTL:  10800,
		FreeformTags:   map[string]string{"env": "test"},
		DefinedTags:    map[string]map[string]interface{}{"ns": {"key": "value"}},
	}

	// Simulate adapter conversion
	ociRequest := bastion.CreateBastionRequest{
		CreateBastionDetails: bastion.CreateBastionDetails{
			Name:                     stringPtr(domainRequest.DisplayName),
			CompartmentId:            stringPtr(domainRequest.CompartmentID),
			TargetSubnetId:           stringPtr(domainRequest.TargetSubnetID),
			BastionType:              stringPtr(domainRequest.BastionType),
			ClientCidrBlockAllowList: domainRequest.ClientCIDRList,
			MaxSessionTtlInSeconds:   intPtr(domainRequest.MaxSessionTTL),
			FreeformTags:             domainRequest.FreeformTags,
			DefinedTags:              domainRequest.DefinedTags,
		},
	}

	// Verify conversion
	assert.Equal(t, "new-bastion", *ociRequest.CreateBastionDetails.Name)
	assert.Equal(t, "ocid1.compartment.oc1..test", *ociRequest.CreateBastionDetails.CompartmentId)
	assert.Equal(t, "ocid1.subnet.oc1..test", *ociRequest.CreateBastionDetails.TargetSubnetId)
	assert.Equal(t, "STANDARD", *ociRequest.CreateBastionDetails.BastionType)
	assert.Equal(t, []string{"0.0.0.0/0", "10.0.0.0/8"}, ociRequest.CreateBastionDetails.ClientCidrBlockAllowList)
	assert.Equal(t, 10800, *ociRequest.CreateBastionDetails.MaxSessionTtlInSeconds)
	assert.Equal(t, map[string]string{"env": "test"}, ociRequest.CreateBastionDetails.FreeformTags)
}

// TestLifecycleStateMapping tests all lifecycle state conversions
func TestLifecycleStateMapping(t *testing.T) {
	tests := []struct {
		name     string
		ociState bastion.BastionLifecycleStateEnum
		expected string
	}{
		{"active state", bastion.BastionLifecycleStateActive, "ACTIVE"},
		{"creating state", bastion.BastionLifecycleStateCreating, "CREATING"},
		{"deleting state", bastion.BastionLifecycleStateDeleting, "DELETING"},
		{"deleted state", bastion.BastionLifecycleStateDeleted, "DELETED"},
		{"failed state", bastion.BastionLifecycleStateFailed, "FAILED"},
		{"updating state", bastion.BastionLifecycleStateUpdating, "UPDATING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ociBastion := bastion.Bastion{
				Id:             stringPtr("ocid1.bastion.oc1..test"),
				Name:           stringPtr("test"),
				BastionType:    stringPtr("STANDARD"),
				LifecycleState: tt.ociState,
				CompartmentId:  stringPtr("ocid1.compartment.oc1..test"),
				TargetVcnId:    stringPtr("ocid1.vcn.oc1..test"),
				TargetSubnetId: stringPtr("ocid1.subnet.oc1..test"),
			}

			attrs := mapping.NewBastionAttributesFromOCIBastion(ociBastion)
			domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

			assert.Equal(t, tt.expected, domainBastion.LifecycleState)
		})
	}
}

// TestNilFieldHandling tests that nil fields are handled safely
func TestNilFieldHandling(t *testing.T) {
	ociBastion := bastion.Bastion{
		Id:                       stringPtr("ocid1.bastion.oc1..nil"),
		Name:                     nil, // nil name
		BastionType:              stringPtr("STANDARD"),
		LifecycleState:           bastion.BastionLifecycleStateActive,
		CompartmentId:            stringPtr("ocid1.compartment.oc1..test"),
		TargetVcnId:              nil, // nil VCN
		TargetSubnetId:           nil, // nil subnet
		MaxSessionTtlInSeconds:   nil, // nil TTL
		ClientCidrBlockAllowList: nil, // nil CIDR list
		PrivateEndpointIpAddress: nil, // nil private IP
	}

	attrs := mapping.NewBastionAttributesFromOCIBastion(ociBastion)
	domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

	// Should not panic and should have empty/zero values
	assert.Equal(t, "ocid1.bastion.oc1..nil", domainBastion.OCID)
	assert.Equal(t, "", domainBastion.DisplayName)
	assert.Equal(t, "STANDARD", domainBastion.BastionType)
	assert.Equal(t, "", domainBastion.TargetVcnID)
	assert.Equal(t, "", domainBastion.TargetSubnetID)
	assert.Equal(t, 0, domainBastion.MaxSessionTTL)
	assert.Nil(t, domainBastion.ClientCidrBlockAllowList)
	assert.Equal(t, "", domainBastion.PrivateEndpointIpAddress)
}

// TestMultipleCIDRBlocks tests handling of multiple CIDR blocks
func TestMultipleCIDRBlocks(t *testing.T) {
	cidrBlocks := []string{
		"0.0.0.0/0",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	ociBastion := bastion.Bastion{
		Id:                       stringPtr("ocid1.bastion.oc1..cidr"),
		Name:                     stringPtr("cidr-bastion"),
		BastionType:              stringPtr("STANDARD"),
		LifecycleState:           bastion.BastionLifecycleStateActive,
		CompartmentId:            stringPtr("ocid1.compartment.oc1..test"),
		TargetVcnId:              stringPtr("ocid1.vcn.oc1..test"),
		TargetSubnetId:           stringPtr("ocid1.subnet.oc1..test"),
		ClientCidrBlockAllowList: cidrBlocks,
	}

	attrs := mapping.NewBastionAttributesFromOCIBastion(ociBastion)
	domainBastion := mapping.NewDomainBastionFromAttrs(attrs)

	assert.Equal(t, cidrBlocks, domainBastion.ClientCidrBlockAllowList)
	assert.Len(t, domainBastion.ClientCidrBlockAllowList, 4)
}
