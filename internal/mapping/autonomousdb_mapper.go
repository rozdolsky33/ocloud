package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
)

type AutonomousDatabaseAttributes struct {
	ID                          *string
	Name                        *string
	CompartmentOCID             *string
	LifecycleState              string
	DbVersion                   *string
	DbWorkload                  string
	LicenseModel                string
	IsFreeTier                  *bool
	WhitelistedIps              []string
	PrivateEndpoint             *string
	PrivateEndpointIp           *string
	PrivateEndpointLabel        *string
	SubnetId                    *string
	NsgIds                      []string
	IsMtlsRequired              *bool
	ComputeModel                string
	EcpuCount                   *float32
	OcpuCount                   *float32
	CpuCoreCount                *int
	DataStorageSizeInTBs        *int
	DataStorageSizeInGBs        *int
	IsAutoScalingEnabled        *bool
	IsStorageAutoScalingEnabled *bool
	OperationsInsightsStatus    string
	DatabaseManagementStatus    string
	DataSafeStatus              string
	IsDataGuardEnabled          *bool
	Role                        *string
	PeerAutonomousDbIds         []string
	BackupRetentionDays         *int
	LastBackupTime              *time.Time
	LatestRestoreTime           *time.Time
	PatchModel                  *string
	NextMaintenanceRunId        *string
	MaintenanceScheduleType     *string
	ConnectionStrings           map[string]string
	Profiles                    []database.DatabaseConnectionStringProfile
	ConnectionUrls              *database.AutonomousDatabaseConnectionUrls
	FreeformTags                map[string]string
	DefinedTags                 map[string]map[string]interface{}
	TimeCreated                 *common.SDKTime
}

func NewAutonomousDatabaseAttributesFromOCIAutonomousDatabase(db database.AutonomousDatabase) *AutonomousDatabaseAttributes {
	return &AutonomousDatabaseAttributes{
		ID:                          db.Id,
		Name:                        db.DbName,
		CompartmentOCID:             db.CompartmentId,
		LifecycleState:              string(db.LifecycleState),
		DbVersion:                   db.DbVersion,
		DbWorkload:                  string(db.DbWorkload),
		LicenseModel:                string(db.LicenseModel),
		IsFreeTier:                  db.IsFreeTier,
		WhitelistedIps:              db.WhitelistedIps,
		PrivateEndpoint:             db.PrivateEndpoint,
		PrivateEndpointIp:           db.PrivateEndpointIp,
		PrivateEndpointLabel:        db.PrivateEndpointLabel,
		SubnetId:                    db.SubnetId,
		NsgIds:                      db.NsgIds,
		IsMtlsRequired:              db.IsMtlsConnectionRequired,
		ComputeModel:                string(db.ComputeModel),
		EcpuCount:                   db.ComputeCount,
		OcpuCount:                   db.OcpuCount,
		CpuCoreCount:                db.CpuCoreCount,
		DataStorageSizeInTBs:        db.DataStorageSizeInTBs,
		DataStorageSizeInGBs:        db.DataStorageSizeInGBs,
		IsAutoScalingEnabled:        db.IsAutoScalingEnabled,
		IsStorageAutoScalingEnabled: db.IsAutoScalingForStorageEnabled,
		OperationsInsightsStatus:    string(db.OperationsInsightsStatus),
		DatabaseManagementStatus:    string(db.DatabaseManagementStatus),
		DataSafeStatus:              string(db.DataSafeStatus),
		IsDataGuardEnabled:          db.IsDataGuardEnabled,
		Role:                        (*string)(&db.Role),
		PeerAutonomousDbIds:         db.PeerDbIds,
		ConnectionStrings:           db.ConnectionStrings.AllConnectionStrings,
		Profiles:                    db.ConnectionStrings.Profiles,
		ConnectionUrls:              db.ConnectionUrls,
		FreeformTags:                db.FreeformTags,
		DefinedTags:                 db.DefinedTags,
		TimeCreated:                 db.TimeCreated,
	}
}

func NewAutonomousDatabaseAttributesFromOCIAutonomousDatabaseSummary(db database.AutonomousDatabaseSummary) *AutonomousDatabaseAttributes {
	return &AutonomousDatabaseAttributes{
		ID:                          db.Id,
		Name:                        db.DbName,
		CompartmentOCID:             db.CompartmentId,
		LifecycleState:              string(db.LifecycleState),
		DbVersion:                   db.DbVersion,
		DbWorkload:                  string(db.DbWorkload),
		LicenseModel:                string(db.LicenseModel),
		IsFreeTier:                  db.IsFreeTier,
		WhitelistedIps:              db.WhitelistedIps,
		PrivateEndpoint:             db.PrivateEndpoint,
		PrivateEndpointIp:           db.PrivateEndpointIp,
		PrivateEndpointLabel:        db.PrivateEndpointLabel,
		SubnetId:                    db.SubnetId,
		NsgIds:                      db.NsgIds,
		IsMtlsRequired:              db.IsMtlsConnectionRequired,
		ComputeModel:                string(db.ComputeModel),
		EcpuCount:                   db.ComputeCount,
		OcpuCount:                   db.OcpuCount,
		CpuCoreCount:                db.CpuCoreCount,
		DataStorageSizeInTBs:        db.DataStorageSizeInTBs,
		DataStorageSizeInGBs:        db.DataStorageSizeInGBs,
		IsAutoScalingEnabled:        db.IsAutoScalingEnabled,
		IsStorageAutoScalingEnabled: db.IsAutoScalingForStorageEnabled,
		OperationsInsightsStatus:    string(db.OperationsInsightsStatus),
		DatabaseManagementStatus:    string(db.DatabaseManagementStatus),
		DataSafeStatus:              string(db.DataSafeStatus),
		IsDataGuardEnabled:          db.IsDataGuardEnabled,
		Role:                        (*string)(&db.Role),
		PeerAutonomousDbIds:         db.PeerDbIds,
		ConnectionStrings:           db.ConnectionStrings.AllConnectionStrings,
		Profiles:                    db.ConnectionStrings.Profiles,
		ConnectionUrls:              db.ConnectionUrls,
		FreeformTags:                db.FreeformTags,
		DefinedTags:                 db.DefinedTags,
		TimeCreated:                 db.TimeCreated,
	}
}

func NewDomainAutonomousDatabaseFromAttrs(attrs *AutonomousDatabaseAttributes) *domain.AutonomousDatabase {
	var timeCreated *time.Time
	if attrs.TimeCreated != nil {
		t := attrs.TimeCreated.Time
		timeCreated = &t
	}
	return &domain.AutonomousDatabase{
		ID:                          *attrs.ID,
		Name:                        *attrs.Name,
		CompartmentOCID:             *attrs.CompartmentOCID,
		LifecycleState:              attrs.LifecycleState,
		DbVersion:                   *attrs.DbVersion,
		DbWorkload:                  attrs.DbWorkload,
		LicenseModel:                attrs.LicenseModel,
		IsFreeTier:                  attrs.IsFreeTier,
		WhitelistedIps:              attrs.WhitelistedIps,
		PrivateEndpoint:             *attrs.PrivateEndpoint,
		PrivateEndpointIp:           *attrs.PrivateEndpointIp,
		PrivateEndpointLabel:        *attrs.PrivateEndpointLabel,
		SubnetId:                    *attrs.SubnetId,
		NsgIds:                      attrs.NsgIds,
		IsMtlsRequired:              attrs.IsMtlsRequired,
		ComputeModel:                attrs.ComputeModel,
		EcpuCount:                   attrs.EcpuCount,
		OcpuCount:                   attrs.OcpuCount,
		CpuCoreCount:                attrs.CpuCoreCount,
		DataStorageSizeInTBs:        attrs.DataStorageSizeInTBs,
		DataStorageSizeInGBs:        attrs.DataStorageSizeInGBs,
		IsAutoScalingEnabled:        attrs.IsAutoScalingEnabled,
		IsStorageAutoScalingEnabled: attrs.IsStorageAutoScalingEnabled,
		OperationsInsightsStatus:    attrs.OperationsInsightsStatus,
		DatabaseManagementStatus:    attrs.DatabaseManagementStatus,
		DataSafeStatus:              attrs.DataSafeStatus,
		IsDataGuardEnabled:          attrs.IsDataGuardEnabled,
		Role:                        *attrs.Role,
		PeerAutonomousDbIds:         attrs.PeerAutonomousDbIds,
		BackupRetentionDays:         attrs.BackupRetentionDays,
		LastBackupTime:              attrs.LastBackupTime,
		LatestRestoreTime:           attrs.LatestRestoreTime,
		PatchModel:                  *attrs.PatchModel,
		NextMaintenanceRunId:        *attrs.NextMaintenanceRunId,
		MaintenanceScheduleType:     *attrs.MaintenanceScheduleType,
		ConnectionStrings:           attrs.ConnectionStrings,
		Profiles:                    attrs.Profiles,
		ConnectionUrls:              attrs.ConnectionUrls,
		FreeformTags:                attrs.FreeformTags,
		DefinedTags:                 attrs.DefinedTags,
		TimeCreated:                 timeCreated,
	}
}
