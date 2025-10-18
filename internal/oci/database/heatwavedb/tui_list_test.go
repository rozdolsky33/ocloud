package heatwavedb

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/mysql"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseListModel(t *testing.T) {
	autoExpand := false
	storage := 1024
	ha := true
	isHeatWave := false
	now := time.Now()

	dbs := []domain.HeatWaveDatabase{
		{
			ID:                   "ocid1.mysqldbsystem.test1",
			DisplayName:          "test-db-1",
			LifecycleState:       "ACTIVE",
			MysqlVersion:         "8.4.6",
			ShapeName:            "MySQL.4",
			DataStorageSizeInGBs: &storage,
			DataStorage: &mysql.DataStorage{
				IsAutoExpandStorageEnabled: &autoExpand,
				AllocatedStorageSizeInGBs:  &storage,
				DataStorageSizeInGBs:       &storage,
			},
			IsHighlyAvailable:         &ha,
			IsHeatWaveClusterAttached: &isHeatWave,
			DatabaseMode:              "READ_WRITE",
			SubnetName:                "test-subnet",
			TimeCreated:               &now,
		},
		{
			ID:                        "ocid1.mysqldbsystem.test2",
			DisplayName:               "test-db-2",
			LifecycleState:            "INACTIVE",
			MysqlVersion:              "8.0.35",
			ShapeName:                 "MySQL.2",
			DataStorageSizeInGBs:      &storage,
			IsHighlyAvailable:         &ha,
			IsHeatWaveClusterAttached: &isHeatWave,
			DatabaseMode:              "READ_ONLY",
			VcnName:                   "test-vcn",
			TimeCreated:               &now,
		},
	}

	// Execute
	model := NewDatabaseListModel(dbs)

	// Assert
	assert.NotNil(t, model)
	// The model should be created successfully with the provided databases
}

func TestDescribeHeatWaveDatabase(t *testing.T) {
	storage := 1024
	ha := true
	isHeatWave := true
	clusterSize := 3
	autoExpand := false
	now := time.Date(2025, 10, 8, 19, 49, 26, 0, time.UTC)

	db := domain.HeatWaveDatabase{
		ID:             "ocid1.mysqldbsystem.test",
		DisplayName:    "prod-db",
		LifecycleState: "ACTIVE",
		MysqlVersion:   "8.4.6",
		ShapeName:      "MySQL.4",
		DataStorage: &mysql.DataStorage{
			IsAutoExpandStorageEnabled: &autoExpand,
			AllocatedStorageSizeInGBs:  &storage,
			DataStorageSizeInGBs:       &storage,
		},
		IsHighlyAvailable:         &ha,
		IsHeatWaveClusterAttached: &isHeatWave,
		HeatWaveCluster: &mysql.HeatWaveClusterSummary{
			ClusterSize: &clusterSize,
		},
		DatabaseMode: "READ_WRITE",
		SubnetName:   "database-subnet",
		VcnName:      "prod-vcn",
		TimeCreated:  &now,
	}

	// Execute
	description := describeHeatWaveDatabase(db)

	// Assert
	assert.NotEmpty(t, description)
	assert.Contains(t, description, "ACTIVE")
	assert.Contains(t, description, "8.4.6")
	assert.Contains(t, description, "MySQL.4")
	assert.Contains(t, description, "1.0TB") // 1024 GB should be displayed as 1.0TB
	assert.Contains(t, description, "HA")
	assert.Contains(t, description, "HeatWave(3)")
	assert.Contains(t, description, "database-subnet")
	assert.Contains(t, description, "READ_WRITE")
	assert.Contains(t, description, "2025-10-08")
}

func TestDescribeHeatWaveDatabase_SmallStorage(t *testing.T) {
	storage := 512
	ha := false
	isHeatWave := false

	db := domain.HeatWaveDatabase{
		ID:                        "ocid1.mysqldbsystem.test",
		DisplayName:               "dev-db",
		LifecycleState:            "INACTIVE",
		MysqlVersion:              "8.0.35",
		ShapeName:                 "MySQL.2",
		DataStorageSizeInGBs:      &storage,
		IsHighlyAvailable:         &ha,
		IsHeatWaveClusterAttached: &isHeatWave,
		DatabaseMode:              "READ_ONLY",
	}

	// Execute
	description := describeHeatWaveDatabase(db)

	// Assert
	assert.NotEmpty(t, description)
	assert.Contains(t, description, "INACTIVE")
	assert.Contains(t, description, "512GB") // Small storage in GB
	assert.NotContains(t, description, "HA")
	assert.NotContains(t, description, "HeatWave")
}

func TestIsTrue(t *testing.T) {
	trueVal := true
	falseVal := false

	assert.True(t, isTrue(&trueVal))
	assert.False(t, isTrue(&falseVal))
	assert.False(t, isTrue(nil))
}

func TestFilterNonEmpty(t *testing.T) {
	result := filterNonEmpty("a", "", "b", "  ", "c")
	assert.Equal(t, []string{"a", "b", "c"}, result)

	result = filterNonEmpty("", "  ", "")
	assert.Equal(t, []string{}, result)

	result = filterNonEmpty("single")
	assert.Equal(t, []string{"single"}, result)
}
