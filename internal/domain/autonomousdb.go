package domain

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/v65/database"
)

// AutonomousDatabase represents an autonomous database instance with its attributes and connection details.
type AutonomousDatabase struct {
	// Identity & lifecycle
	ID               string
	Name             string
	CompartmentOCID  string
	Region           string
	LifecycleState   string
	LifecycleDetails string
	DbVersion        string
	DbWorkload       string
	LicenseModel     string
	IsFreeTier       *bool

	// Networking
	IsPubliclyAccessible *bool
	WhitelistedIps       []string
	PrivateEndpoint      string
	PrivateEndpointIp    string
	PrivateEndpointLabel string
	SubnetId             string
	NsgIds               []string
	IsMtlsRequired       *bool

	// Capacity & autoscaling
	ComputeModel                string
	EcpuCount                   *float32
	OcpuCount                   *float32
	CpuCoreCount                *int
	DataStorageSizeInTBs        *int
	IsAutoScalingEnabled        *bool
	IsStorageAutoScalingEnabled *bool

	// Security & management integrations
	OperationsInsightsStatus string
	DatabaseManagementStatus string
	DataSafeStatus           string
	FreeformTags             map[string]string
	DefinedTags              map[string]map[string]interface{}

	// Resiliency / Data Guard
	IsDataGuardEnabled  *bool
	Role                string
	PeerAutonomousDbIds []string

	// Backups & recovery
	BackupRetentionDays *int
	LastBackupTime      *time.Time
	LatestRestoreTime   *time.Time

	// Maintenance & patching
	PatchModel              string
	NextMaintenanceRunId    string
	MaintenanceScheduleType string

	// URLs / Tools
	ConnectionStrings map[string]string
	Profiles          []database.DatabaseConnectionStringProfile
	ConnectionUrls    *database.AutonomousDatabaseConnectionUrls

	// Timestamps
	TimeCreated *time.Time
}

// AutonomousDatabaseRepository defines the interface for interacting with Autonomous Database data.
type AutonomousDatabaseRepository interface {
	GetAutonomousDatabase(ctx context.Context, ocid string) (*AutonomousDatabase, error)
	ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]AutonomousDatabase, error)
}
