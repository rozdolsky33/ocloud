package database

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/mysql"
)

// HeatWaveDatabase represents a MySQL HeatWave database instance with its attributes and connection details.
type HeatWaveDatabase struct {
	// Identity & lifecycle
	ID              string
	DisplayName     string
	CompartmentOCID string
	LifecycleState  string
	Description     string
	MysqlVersion    string
	TimeCreated     *time.Time
	TimeUpdated     *time.Time

	// Networking
	SubnetId   string
	SubnetName string
	VcnID      string
	VcnName    string
	NsgIds     []string
	NsgNames   []string
	Endpoints  []mysql.DbSystemEndpoint

	// Capacity & storage
	ShapeName            string
	DataStorageSizeInGBs *int
	IsHighlyAvailable    *bool

	// HeatWave cluster
	IsHeatWaveClusterAttached *bool
	HeatWaveCluster           *mysql.HeatWaveClusterSummary

	// Configuration & mode
	DatabaseMode    string
	AccessMode      string
	ConfigurationId string

	// Placement
	AvailabilityDomain string
	FaultDomain        string

	// Backup & maintenance
	BackupPolicy    *mysql.BackupPolicy
	DeletionPolicy  *mysql.DeletionPolicyDetails
	MaintenanceInfo *mysql.MaintenanceDetails

	// REST & security
	RestInfo *mysql.RestDetails

	// Tags
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

// HeatWaveDatabaseRepository defines the interface for interacting with HeatWave Database data.
type HeatWaveDatabaseRepository interface {
	GetHeatWaveDatabase(ctx context.Context, ocid string) (*HeatWaveDatabase, error)
	ListHeatWaveDatabases(ctx context.Context, compartmentID string) ([]HeatWaveDatabase, error)
	ListEnrichedHeatWaveDatabases(ctx context.Context, compartmentID string) ([]HeatWaveDatabase, error)
}
