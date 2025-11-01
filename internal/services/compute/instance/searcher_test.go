package instance

import (
	"testing"

	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/stretchr/testify/require"
)

func TestSearchableInstance_ToIndexable_WithSecurityListsAndNSGs(t *testing.T) {
	instance := compute.Instance{
		DisplayName:        "web-server-1",
		Hostname:           "web1",
		PrimaryIP:          "10.0.1.5",
		ImageName:          "Oracle Linux 8",
		ImageOS:            "Oracle Linux",
		Shape:              "VM.Standard3.Flex",
		OCID:               "ocid1.instance.oc1..aaa",
		FaultDomain:        "FAULT-DOMAIN-1",
		AvailabilityDomain: "AD-1",
		VcnName:            "production-vcn",
		SubnetName:         "web-subnet",
		SecurityListNames:  []string{"default-security-list", "web-security-list"},
		NsgNames:           []string{"app-tier-nsg", "web-tier-nsg"},
		FreeformTags:       map[string]string{"env": "prod"},
		DefinedTags:        map[string]map[string]interface{}{"namespace": {"key": "value"}},
	}

	searchable := SearchableInstance{instance}
	indexed := searchable.ToIndexable()

	// Test basic fields
	require.Equal(t, "web-server-1", indexed["Name"])
	require.Equal(t, "web1", indexed["Hostname"])
	require.Equal(t, "10.0.1.5", indexed["PrimaryIP"])
	require.Equal(t, "oracle linux 8", indexed["ImageName"])
	require.Equal(t, "oracle linux", indexed["ImageOS"])
	require.Equal(t, "vm.standard3.flex", indexed["Shape"])
	require.Equal(t, "ocid1.instance.oc1..aaa", indexed["OCID"])
	require.Equal(t, "fault-domain-1", indexed["FD"])
	require.Equal(t, "ad-1", indexed["AD"])
	require.Equal(t, "production-vcn", indexed["VcnName"])
	require.Equal(t, "web-subnet", indexed["SubnetName"])

	// Test Security Lists
	securityLists, ok := indexed["SecurityLists"].(string)
	require.True(t, ok, "SecurityLists should be a string")
	require.Contains(t, securityLists, "default-security-list")
	require.Contains(t, securityLists, "web-security-list")

	// Test NSGs
	nsgs, ok := indexed["NSGs"].(string)
	require.True(t, ok, "NSGs should be a string")
	require.Contains(t, nsgs, "app-tier-nsg")
	require.Contains(t, nsgs, "web-tier-nsg")

	// Test tags
	require.NotEmpty(t, indexed["TagsKV"])
	require.NotEmpty(t, indexed["TagsVal"])
}

func TestSearchableInstance_ToIndexable_EmptySecurityListsAndNSGs(t *testing.T) {
	instance := compute.Instance{
		DisplayName:       "app-server",
		PrimaryIP:         "10.0.2.10",
		SecurityListNames: []string{},
		NsgNames:          []string{},
	}

	searchable := SearchableInstance{instance}
	indexed := searchable.ToIndexable()

	// Test Security Lists - should be empty string when empty array
	securityLists, ok := indexed["SecurityLists"].(string)
	require.True(t, ok, "SecurityLists should be a string")
	require.Empty(t, securityLists)

	// Test NSGs - should be empty string when empty array
	nsgs, ok := indexed["NSGs"].(string)
	require.True(t, ok, "NSGs should be a string")
	require.Empty(t, nsgs)
}

func TestSearchableInstance_ToIndexable_NilSecurityListsAndNSGs(t *testing.T) {
	instance := compute.Instance{
		DisplayName:       "db-server",
		PrimaryIP:         "10.0.3.20",
		SecurityListNames: nil,
		NsgNames:          nil,
	}

	searchable := SearchableInstance{instance}
	indexed := searchable.ToIndexable()

	// Test Security Lists - should be empty string when nil
	securityLists, ok := indexed["SecurityLists"].(string)
	require.True(t, ok, "SecurityLists should be a string")
	require.Empty(t, securityLists)

	// Test NSGs - should be empty string when nil
	nsgs, ok := indexed["NSGs"].(string)
	require.True(t, ok, "NSGs should be a string")
	require.Empty(t, nsgs)
}

func TestSearchableInstance_ToIndexable_CaseInsensitive(t *testing.T) {
	instance := compute.Instance{
		DisplayName:       "API-Server",
		SecurityListNames: []string{"Production-Security-List"},
		NsgNames:          []string{"API-Gateway-NSG"},
	}

	searchable := SearchableInstance{instance}
	indexed := searchable.ToIndexable()

	// All fields should be lowercase
	require.Equal(t, "api-server", indexed["Name"])

	securityLists, ok := indexed["SecurityLists"].(string)
	require.True(t, ok)
	require.Equal(t, "production-security-list", securityLists)

	nsgs, ok := indexed["NSGs"].(string)
	require.True(t, ok)
	require.Equal(t, "api-gateway-nsg", nsgs)
}

func TestGetSearchableFields_IncludesSecurityListsAndNSGs(t *testing.T) {
	fields := GetSearchableFields()

	require.Contains(t, fields, "SecurityLists", "SearchableFields should include SecurityLists")
	require.Contains(t, fields, "NSGs", "SearchableFields should include NSGs")

	// Also verify other expected fields still exist
	require.Contains(t, fields, "Name")
	require.Contains(t, fields, "Hostname")
	require.Contains(t, fields, "VcnName")
	require.Contains(t, fields, "SubnetName")
}

func TestGetBoostedFields_Unchanged(t *testing.T) {
	boosted := GetBoostedFields()

	// Boosted fields should remain Name and Hostname
	require.Contains(t, boosted, "Name")
	require.Contains(t, boosted, "Hostname")
	require.Len(t, boosted, 2, "Should only have 2 boosted fields")
}

func TestToSearchableInstances(t *testing.T) {
	instances := []Instance{
		{
			DisplayName:       "instance-1",
			SecurityListNames: []string{"sl1"},
			NsgNames:          []string{"nsg1"},
		},
		{
			DisplayName:       "instance-2",
			SecurityListNames: []string{"sl2"},
			NsgNames:          []string{"nsg2"},
		},
	}

	searchable := ToSearchableInstances(instances)

	require.Len(t, searchable, 2)

	// Test first instance
	indexed0 := searchable[0].ToIndexable()
	require.Equal(t, "instance-1", indexed0["Name"])
	require.Contains(t, indexed0["SecurityLists"].(string), "sl1")
	require.Contains(t, indexed0["NSGs"].(string), "nsg1")

	// Test second instance
	indexed1 := searchable[1].ToIndexable()
	require.Equal(t, "instance-2", indexed1["Name"])
	require.Contains(t, indexed1["SecurityLists"].(string), "sl2")
	require.Contains(t, indexed1["NSGs"].(string), "nsg2")
}
