package heatwavedb

import (
	"strings"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/mysql"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/stretchr/testify/assert"
)

func TestSearcher_GetFields(t *testing.T) {
	fields := GetSearchableFields()
	boost := GetBoostedFields()

	// basic expectations
	assert.Contains(t, fields, "Name")
	assert.Contains(t, fields, "OCID")
	assert.Contains(t, fields, "State")
	assert.Contains(t, fields, "MysqlVersion")
	assert.Contains(t, fields, "ShapeName")
	assert.Contains(t, fields, "VcnName")
	assert.Contains(t, fields, "SubnetName")
	assert.Contains(t, fields, "IpAddress")
	assert.Contains(t, fields, "TagsKV")
	assert.Contains(t, fields, "TagsVal")

	// boosted are subset of fields
	for _, b := range boost {
		assert.Contains(t, fields, b, "boosted field %s should be in searchable fields", b)
	}
}

func TestSearchableHeatWaveDatabase_ToIndexable_LowercasesAndMaps(t *testing.T) {
	storage := 1024
	clusterSize := 3
	port := 3306
	portX := 33060
	ha := true
	isHeatWave := true

	db := database.HeatWaveDatabase{
		ID:             "ocid1.mysqldbsystem.oc1.iad.test",
		DisplayName:    "Prod-MySQL-DB",
		LifecycleState: "ACTIVE",
		Description:    "Production MySQL Database",
		MysqlVersion:   "8.4.6",
		ShapeName:      "MySQL.4",
		DataStorage: &mysql.DataStorage{
			DataStorageSizeInGBs: &storage,
		},
		DatabaseMode:              "READ_WRITE",
		AccessMode:                "READ_WRITE",
		VcnName:                   "Prod-VCN",
		SubnetName:                "Database-Subnet",
		IpAddress:                 "10.0.20.175",
		HostnameLabel:             "prod-mysql",
		Port:                      &port,
		PortX:                     &portX,
		IsHighlyAvailable:         &ha,
		IsHeatWaveClusterAttached: &isHeatWave,
		HeatWaveCluster: &mysql.HeatWaveClusterSummary{
			ClusterSize: &clusterSize,
		},
		AvailabilityDomain: "AD-1",
		FaultDomain:        "FD-1",
		CrashRecovery:      "ENABLED",
		FreeformTags:       map[string]string{"Environment": "Production"},
	}

	m := SearchableHeatWaveDatabase{db}.ToIndexable()

	// strings are lowercased
	stringFields := []string{
		"Name", "OCID", "State", "Description", "MysqlVersion", "ShapeName",
		"DatabaseMode", "AccessMode", "VcnName", "SubnetName", "IpAddress",
		"HostnameLabel", "AvailabilityDomain", "FaultDomain", "CrashRecovery",
		"TagsKV", "TagsVal",
	}

	for _, k := range stringFields {
		if v, ok := m[k].(string); ok {
			assert.Equal(t, strings.ToLower(v), v, "%s should be lowercased", k)
		}
	}

	// specific values
	assert.Equal(t, "prod-mysql-db", m["Name"])
	assert.Equal(t, "active", m["State"])
	assert.Equal(t, "8.4.6", m["MysqlVersion"])
	assert.Equal(t, "mysql.4", m["ShapeName"])
	assert.Equal(t, "1024", m["StorageGB"])
	assert.Equal(t, "3", m["ClusterSize"])
	assert.Equal(t, "read_write", m["DatabaseMode"])
	assert.Equal(t, "prod-vcn", m["VcnName"])
	assert.Equal(t, "database-subnet", m["SubnetName"])
	assert.Equal(t, "10.0.20.175", m["IpAddress"])
}

func TestSearchableHeatWaveDatabase_ToIndexable_HandleNilFields(t *testing.T) {
	db := database.HeatWaveDatabase{
		ID:             "ocid1.mysqldbsystem.test",
		DisplayName:    "test-db",
		LifecycleState: "ACTIVE",
		MysqlVersion:   "8.0.35",
		ShapeName:      "MySQL.2",
		// All optional fields are nil
	}

	m := SearchableHeatWaveDatabase{db}.ToIndexable()

	// Should not panic and should have empty strings for nil fields
	assert.Equal(t, "test-db", m["Name"])
	assert.Equal(t, "active", m["State"])
	assert.Equal(t, "", m["StorageGB"])
	assert.Equal(t, "", m["ClusterSize"])
	assert.Equal(t, "", m["Description"])
	assert.Equal(t, "", m["VcnName"])
	assert.Equal(t, "", m["SubnetName"])
}

func TestSearchableHeatWaveDatabase_ToIndexable_StorageFallback(t *testing.T) {
	storage := 512

	// Test with deprecated field
	db1 := database.HeatWaveDatabase{
		ID:                   "ocid1.test1",
		DisplayName:          "db1",
		DataStorageSizeInGBs: &storage,
	}
	m1 := SearchableHeatWaveDatabase{db1}.ToIndexable()
	assert.Equal(t, "512", m1["StorageGB"])

	// Test with new DataStorage object
	db2 := database.HeatWaveDatabase{
		ID:          "ocid1.test2",
		DisplayName: "db2",
		DataStorage: &mysql.DataStorage{
			DataStorageSizeInGBs: &storage,
		},
	}
	m2 := SearchableHeatWaveDatabase{db2}.ToIndexable()
	assert.Equal(t, "512", m2["StorageGB"])
}

func TestSearchableHeatWaveDatabase_ToIndexable_JoinsSlices(t *testing.T) {
	db := database.HeatWaveDatabase{
		ID:          "ocid1.test",
		DisplayName: "test-db",
		NsgNames:    []string{"NSG-1", "NSG-2", "NSG-3"},
		NsgIds:      []string{"ocid1.nsg.1", "ocid1.nsg.2"},
	}

	m := SearchableHeatWaveDatabase{db}.ToIndexable()

	// Slices should be joined with commas and lowercased
	assert.Equal(t, "nsg-1,nsg-2,nsg-3", m["NsgNames"])
	assert.Equal(t, "ocid1.nsg.1,ocid1.nsg.2", m["NsgIds"])
}

func TestToSearchableHeatWaveDbs(t *testing.T) {
	storage := 1024
	dbs := []database.HeatWaveDatabase{
		{
			ID:          "ocid1.db1",
			DisplayName: "db1",
			DataStorage: &mysql.DataStorage{
				DataStorageSizeInGBs: &storage,
			},
		},
		{
			ID:          "ocid1.db2",
			DisplayName: "db2",
			DataStorage: &mysql.DataStorage{
				DataStorageSizeInGBs: &storage,
			},
		},
	}

	searchable := ToSearchableHeatWaveDbs(dbs)

	assert.Len(t, searchable, 2)
	for i, s := range searchable {
		m := s.ToIndexable()
		assert.Contains(t, m["Name"], "db")
		// OCID is lowercased in ToIndexable, so compare lowercase
		assert.Equal(t, strings.ToLower(dbs[i].ID), m["OCID"].(string))
	}
}
