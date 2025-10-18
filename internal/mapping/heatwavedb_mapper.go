package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/mysql"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
)

// HeatWaveDatabaseAttributes holds intermediate attributes for mapping from OCI SDK to a domain model.
type HeatWaveDatabaseAttributes struct {
	ID                         *string
	DisplayName                *string
	CompartmentOCID            *string
	LifecycleState             string
	Description                *string
	MysqlVersion               *string
	TimeCreated                *common.SDKTime
	TimeUpdated                *common.SDKTime
	SubnetId                   *string
	NsgIds                     []string
	Endpoints                  []mysql.DbSystemEndpoint
	ShapeName                  *string
	DataStorageSizeInGBs       *int
	DataStorage                *mysql.DataStorage
	IsHighlyAvailable          *bool
	HostnameLabel              *string
	IpAddress                  *string
	Port                       *int
	PortX                      *int
	IsHeatWaveClusterAttached  *bool
	HeatWaveCluster            *mysql.HeatWaveClusterSummary
	DatabaseMode               string
	AccessMode                 string
	ConfigurationId            *string
	AvailabilityDomain         *string
	FaultDomain                *string
	BackupPolicy               *mysql.BackupPolicy
	DeletionPolicy             *mysql.DeletionPolicyDetails
	MaintenanceInfo            *mysql.MaintenanceDetails
	CrashRecovery              string
	PointInTimeRecoveryDetails *mysql.PointInTimeRecoveryDetails
	RestInfo                   *mysql.RestDetails
	SecureConnections          *mysql.SecureConnectionDetails
	EncryptData                *mysql.EncryptDataDetails
	ReadEndpoint               *mysql.ReadEndpointDetails
	DatabaseManagement         string
	CustomerContacts           []mysql.CustomerContact
	LifecycleDetails           *string
	FreeformTags               map[string]string
	DefinedTags                map[string]map[string]interface{}
}

// NewHeatWaveDatabaseAttributesFromOCIDbSystem converts a full OCI MySQL DbSystem to attributes.
func NewHeatWaveDatabaseAttributesFromOCIDbSystem(db mysql.DbSystem) *HeatWaveDatabaseAttributes {
	return &HeatWaveDatabaseAttributes{
		ID:                         db.Id,
		DisplayName:                db.DisplayName,
		CompartmentOCID:            db.CompartmentId,
		LifecycleState:             string(db.LifecycleState),
		Description:                db.Description,
		MysqlVersion:               db.MysqlVersion,
		TimeCreated:                db.TimeCreated,
		TimeUpdated:                db.TimeUpdated,
		SubnetId:                   db.SubnetId,
		NsgIds:                     db.NsgIds,
		Endpoints:                  db.Endpoints,
		ShapeName:                  db.ShapeName,
		DataStorageSizeInGBs:       db.DataStorageSizeInGBs,
		DataStorage:                db.DataStorage,
		IsHighlyAvailable:          db.IsHighlyAvailable,
		HostnameLabel:              db.HostnameLabel,
		IpAddress:                  db.IpAddress,
		Port:                       db.Port,
		PortX:                      db.PortX,
		IsHeatWaveClusterAttached:  db.IsHeatWaveClusterAttached,
		HeatWaveCluster:            db.HeatWaveCluster,
		DatabaseMode:               string(db.DatabaseMode),
		AccessMode:                 string(db.AccessMode),
		ConfigurationId:            db.ConfigurationId,
		AvailabilityDomain:         db.AvailabilityDomain,
		FaultDomain:                db.FaultDomain,
		BackupPolicy:               db.BackupPolicy,
		DeletionPolicy:             db.DeletionPolicy,
		MaintenanceInfo:            db.Maintenance,
		CrashRecovery:              string(db.CrashRecovery),
		PointInTimeRecoveryDetails: db.PointInTimeRecoveryDetails,
		RestInfo:                   db.Rest,
		SecureConnections:          db.SecureConnections,
		EncryptData:                db.EncryptData,
		ReadEndpoint:               db.ReadEndpoint,
		DatabaseManagement:         string(db.DatabaseManagement),
		CustomerContacts:           db.CustomerContacts,
		LifecycleDetails:           db.LifecycleDetails,
		FreeformTags:               db.FreeformTags,
		DefinedTags:                db.DefinedTags,
	}
}

// NewHeatWaveDatabaseAttributesFromOCIDbSystemSummary converts an OCI MySQL DbSystemSummary to attributes.
func NewHeatWaveDatabaseAttributesFromOCIDbSystemSummary(db mysql.DbSystemSummary) *HeatWaveDatabaseAttributes {
	return &HeatWaveDatabaseAttributes{
		ID:                        db.Id,
		DisplayName:               db.DisplayName,
		CompartmentOCID:           db.CompartmentId,
		LifecycleState:            string(db.LifecycleState),
		Description:               db.Description,
		MysqlVersion:              db.MysqlVersion,
		TimeCreated:               db.TimeCreated,
		TimeUpdated:               db.TimeUpdated,
		ShapeName:                 db.ShapeName,
		IsHighlyAvailable:         db.IsHighlyAvailable,
		IsHeatWaveClusterAttached: db.IsHeatWaveClusterAttached,
		HeatWaveCluster:           db.HeatWaveCluster,
		DatabaseMode:              string(db.DatabaseMode),
		AccessMode:                string(db.AccessMode),
		AvailabilityDomain:        db.AvailabilityDomain,
		FaultDomain:               db.FaultDomain,
		Endpoints:                 db.Endpoints,
		BackupPolicy:              db.BackupPolicy,
		DeletionPolicy:            db.DeletionPolicy,
		RestInfo:                  db.Rest,
		FreeformTags:              db.FreeformTags,
		DefinedTags:               db.DefinedTags,
	}
}

// NewDomainHeatWaveDatabaseFromAttrs converts HeatWaveDatabaseAttributes to domain.HeatWaveDatabase.
func NewDomainHeatWaveDatabaseFromAttrs(attrs *HeatWaveDatabaseAttributes) *domain.HeatWaveDatabase {
	// Helper to safely dereference string pointers
	val := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	var timeCreated, timeUpdated *time.Time
	if attrs.TimeCreated != nil {
		t := attrs.TimeCreated.Time
		timeCreated = &t
	}
	if attrs.TimeUpdated != nil {
		t := attrs.TimeUpdated.Time
		timeUpdated = &t
	}

	return &domain.HeatWaveDatabase{
		ID:                         val(attrs.ID),
		DisplayName:                val(attrs.DisplayName),
		CompartmentOCID:            val(attrs.CompartmentOCID),
		LifecycleState:             attrs.LifecycleState,
		Description:                val(attrs.Description),
		MysqlVersion:               val(attrs.MysqlVersion),
		TimeCreated:                timeCreated,
		TimeUpdated:                timeUpdated,
		SubnetId:                   val(attrs.SubnetId),
		NsgIds:                     attrs.NsgIds,
		Endpoints:                  attrs.Endpoints,
		ShapeName:                  val(attrs.ShapeName),
		DataStorageSizeInGBs:       attrs.DataStorageSizeInGBs,
		DataStorage:                attrs.DataStorage,
		IsHighlyAvailable:          attrs.IsHighlyAvailable,
		HostnameLabel:              val(attrs.HostnameLabel),
		IpAddress:                  val(attrs.IpAddress),
		Port:                       attrs.Port,
		PortX:                      attrs.PortX,
		IsHeatWaveClusterAttached:  attrs.IsHeatWaveClusterAttached,
		HeatWaveCluster:            attrs.HeatWaveCluster,
		DatabaseMode:               attrs.DatabaseMode,
		AccessMode:                 attrs.AccessMode,
		ConfigurationId:            val(attrs.ConfigurationId),
		AvailabilityDomain:         val(attrs.AvailabilityDomain),
		FaultDomain:                val(attrs.FaultDomain),
		BackupPolicy:               attrs.BackupPolicy,
		DeletionPolicy:             attrs.DeletionPolicy,
		MaintenanceInfo:            attrs.MaintenanceInfo,
		CrashRecovery:              attrs.CrashRecovery,
		PointInTimeRecoveryDetails: attrs.PointInTimeRecoveryDetails,
		RestInfo:                   attrs.RestInfo,
		SecureConnections:          attrs.SecureConnections,
		EncryptData:                attrs.EncryptData,
		ReadEndpoint:               attrs.ReadEndpoint,
		DatabaseManagement:         attrs.DatabaseManagement,
		CustomerContacts:           attrs.CustomerContacts,
		LifecycleDetails:           val(attrs.LifecycleDetails),
		FreeformTags:               attrs.FreeformTags,
		DefinedTags:                attrs.DefinedTags,
	}
}
