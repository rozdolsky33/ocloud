package mapping

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNewHeatWaveDatabaseAttributesFromOCIDbSystem(t *testing.T) {
	// Setup
	id := "ocid1.mysqldbsystem.test"
	displayName := "test-heatwave-db"
	compartmentID := "ocid1.compartment.test"
	mysqlVersion := "8.0.35"
	subnetID := "ocid1.subnet.test"
	shapeName := "MySQL.VM.Standard.E4.4.128GB"
	storage := 100
	isHA := true
	isHeatWaveAttached := true
	clusterSize := 2
	now := common.SDKTime{Time: time.Now()}

	dbSystem := mysql.DbSystem{
		Id:                        &id,
		DisplayName:               &displayName,
		CompartmentId:             &compartmentID,
		MysqlVersion:              &mysqlVersion,
		SubnetId:                  &subnetID,
		ShapeName:                 &shapeName,
		DataStorageSizeInGBs:      &storage,
		IsHighlyAvailable:         &isHA,
		IsHeatWaveClusterAttached: &isHeatWaveAttached,
		HeatWaveCluster: &mysql.HeatWaveClusterSummary{
			ClusterSize: &clusterSize,
		},
		LifecycleState: mysql.DbSystemLifecycleStateActive,
		DatabaseMode:   mysql.DbSystemDatabaseModeWrite,
		AccessMode:     mysql.DbSystemAccessModeUnrestricted,
		TimeCreated:    &now,
		TimeUpdated:    &now,
	}

	// Execute
	attrs := NewHeatWaveDatabaseAttributesFromOCIDbSystem(dbSystem)

	// Assert
	assert.NotNil(t, attrs)
	assert.Equal(t, &id, attrs.ID)
	assert.Equal(t, &displayName, attrs.DisplayName)
	assert.Equal(t, &compartmentID, attrs.CompartmentOCID)
	assert.Equal(t, string(mysql.DbSystemLifecycleStateActive), attrs.LifecycleState)
	assert.Equal(t, &mysqlVersion, attrs.MysqlVersion)
	assert.Equal(t, &subnetID, attrs.SubnetId)
	assert.Equal(t, &shapeName, attrs.ShapeName)
	assert.Equal(t, &storage, attrs.DataStorageSizeInGBs)
	assert.Equal(t, &isHA, attrs.IsHighlyAvailable)
	assert.Equal(t, &isHeatWaveAttached, attrs.IsHeatWaveClusterAttached)
	assert.NotNil(t, attrs.HeatWaveCluster)
	assert.Equal(t, &clusterSize, attrs.HeatWaveCluster.ClusterSize)
	assert.Equal(t, string(mysql.DbSystemDatabaseModeWrite), attrs.DatabaseMode)
	assert.Equal(t, string(mysql.DbSystemAccessModeUnrestricted), attrs.AccessMode)
}

func TestNewHeatWaveDatabaseAttributesFromOCIDbSystemSummary(t *testing.T) {
	// Setup
	id := "ocid1.mysqldbsystem.test"
	displayName := "test-heatwave-db"
	compartmentID := "ocid1.compartment.test"
	mysqlVersion := "8.0.35"
	isHA := false
	isHeatWaveAttached := false
	now := common.SDKTime{Time: time.Now()}

	summary := mysql.DbSystemSummary{
		Id:                        &id,
		DisplayName:               &displayName,
		CompartmentId:             &compartmentID,
		MysqlVersion:              &mysqlVersion,
		IsHighlyAvailable:         &isHA,
		IsHeatWaveClusterAttached: &isHeatWaveAttached,
		LifecycleState:            mysql.DbSystemLifecycleStateActive,
		DatabaseMode:              mysql.DbSystemDatabaseModeWrite,
		AccessMode:                mysql.DbSystemAccessModeUnrestricted,
		TimeCreated:               &now,
		TimeUpdated:               &now,
	}

	// Execute
	attrs := NewHeatWaveDatabaseAttributesFromOCIDbSystemSummary(summary)

	// Assert
	assert.NotNil(t, attrs)
	assert.Equal(t, &id, attrs.ID)
	assert.Equal(t, &displayName, attrs.DisplayName)
	assert.Equal(t, &compartmentID, attrs.CompartmentOCID)
	assert.Equal(t, string(mysql.DbSystemLifecycleStateActive), attrs.LifecycleState)
	assert.Equal(t, &mysqlVersion, attrs.MysqlVersion)
	assert.Equal(t, &isHA, attrs.IsHighlyAvailable)
	assert.Equal(t, &isHeatWaveAttached, attrs.IsHeatWaveClusterAttached)
}

func TestNewDomainHeatWaveDatabaseFromAttrs(t *testing.T) {
	// Setup
	id := "ocid1.mysqldbsystem.test"
	displayName := "test-heatwave-db"
	compartmentID := "ocid1.compartment.test"
	lifecycleState := "ACTIVE"
	mysqlVersion := "8.0.35"
	subnetID := "ocid1.subnet.test"
	shapeName := "MySQL.VM.Standard.E4.4.128GB"
	storage := 100
	isHA := true
	isHeatWaveAttached := true
	databaseMode := "READ_WRITE"
	accessMode := "UNRESTRICTED"
	now := common.SDKTime{Time: time.Now()}

	attrs := &HeatWaveDatabaseAttributes{
		ID:                        &id,
		DisplayName:               &displayName,
		CompartmentOCID:           &compartmentID,
		LifecycleState:            lifecycleState,
		MysqlVersion:              &mysqlVersion,
		SubnetId:                  &subnetID,
		ShapeName:                 &shapeName,
		DataStorageSizeInGBs:      &storage,
		IsHighlyAvailable:         &isHA,
		IsHeatWaveClusterAttached: &isHeatWaveAttached,
		DatabaseMode:              databaseMode,
		AccessMode:                accessMode,
		TimeCreated:               &now,
		TimeUpdated:               &now,
	}

	// Execute
	db := NewDomainHeatWaveDatabaseFromAttrs(attrs)

	// Assert
	assert.NotNil(t, db)
	assert.Equal(t, id, db.ID)
	assert.Equal(t, displayName, db.DisplayName)
	assert.Equal(t, compartmentID, db.CompartmentOCID)
	assert.Equal(t, lifecycleState, db.LifecycleState)
	assert.Equal(t, mysqlVersion, db.MysqlVersion)
	assert.Equal(t, subnetID, db.SubnetId)
	assert.Equal(t, shapeName, db.ShapeName)
	assert.Equal(t, &storage, db.DataStorageSizeInGBs)
	assert.Equal(t, &isHA, db.IsHighlyAvailable)
	assert.Equal(t, &isHeatWaveAttached, db.IsHeatWaveClusterAttached)
	assert.Equal(t, databaseMode, db.DatabaseMode)
	assert.Equal(t, accessMode, db.AccessMode)
	assert.NotNil(t, db.TimeCreated)
	assert.NotNil(t, db.TimeUpdated)
}

func TestNewDomainHeatWaveDatabaseFromAttrs_NilValues(t *testing.T) {
	// Setup - all nil values
	attrs := &HeatWaveDatabaseAttributes{
		ID:             nil,
		DisplayName:    nil,
		LifecycleState: "ACTIVE",
	}

	// Execute
	db := NewDomainHeatWaveDatabaseFromAttrs(attrs)

	// Assert
	assert.NotNil(t, db)
	assert.Equal(t, "", db.ID)
	assert.Equal(t, "", db.DisplayName)
	assert.Equal(t, "ACTIVE", db.LifecycleState)
	assert.Nil(t, db.TimeCreated)
	assert.Nil(t, db.TimeUpdated)
}
