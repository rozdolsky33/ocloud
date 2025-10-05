package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestAutonomousDB_Attributes_From_OCI_And_Domain(t *testing.T) {
	id := "ocid1.autonomousdatabase.oc1..db"
	name := "adb1"
	comp := "ocid1.compartment.oc1..cpt"
	dbver := "19c"
	workload := database.AutonomousDatabaseDbWorkloadOltp
	lic := database.AutonomousDatabaseLicenseModelLicenseIncluded
	free := true
	whitelisted := []string{"1.1.1.1"}
	pe := "10.0.2.2"
	peip := "10.0.2.2"
	pelabel := "adb1-pe"
	subnet := "ocid1.subnet.oc1..subnet"
	nsgs := []string{"ocid1.nsg.oc1..nsg1", "ocid1.nsg.oc1..nsg2"}
	mtls := true
	computeModel := database.AutonomousDatabaseComputeModelEcpu
	ecpu := float32(4)
	ocpu := float32(2)
	cores := 2
	sizeTB := 1
	sizeGB := 1024
	autoscale := true
	autoscaleStorage := true
	opsStatus := database.AutonomousDatabaseOperationsInsightsStatusEnabled
	mgmtStatus := database.AutonomousDatabaseDatabaseManagementStatusEnabled
	dataSafeStatus := database.AutonomousDatabaseDataSafeStatusRegistered
	dgEnabled := true
	role := database.AutonomousDatabaseRolePrimary
	peers := []string{"ocid1.autonomousdatabase.oc1..peer"}
	backupDays := 7
	lastBackup := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
	latestRestore := time.Now().UTC().Truncate(time.Second)
	patchModel := "RELEASE_UPDATES"
	nextMaint := "ocid1.maint.oc1..id"
	schedType := "TIME_BASED"
	created := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Second)

	oci := database.AutonomousDatabase{
		Id:                             &id,
		DbName:                         &name,
		CompartmentId:                  &comp,
		LifecycleState:                 database.AutonomousDatabaseLifecycleStateAvailable,
		DbVersion:                      &dbver,
		DbWorkload:                     workload,
		LicenseModel:                   lic,
		IsFreeTier:                     &free,
		WhitelistedIps:                 whitelisted,
		PrivateEndpoint:                &pe,
		PrivateEndpointIp:              &peip,
		PrivateEndpointLabel:           &pelabel,
		SubnetId:                       &subnet,
		NsgIds:                         nsgs,
		IsMtlsConnectionRequired:       &mtls,
		ComputeModel:                   computeModel,
		ComputeCount:                   &ecpu,
		OcpuCount:                      &ocpu,
		CpuCoreCount:                   &cores,
		DataStorageSizeInTBs:           &sizeTB,
		DataStorageSizeInGBs:           &sizeGB,
		IsAutoScalingEnabled:           &autoscale,
		IsAutoScalingForStorageEnabled: &autoscaleStorage,
		OperationsInsightsStatus:       opsStatus,
		DatabaseManagementStatus:       mgmtStatus,
		DataSafeStatus:                 dataSafeStatus,
		IsDataGuardEnabled:             &dgEnabled,
		Role:                           role,
		PeerDbIds:                      peers,
		ConnectionStrings:              &database.AutonomousDatabaseConnectionStrings{AllConnectionStrings: map[string]string{"db": "conn"}},
		ConnectionUrls:                 &database.AutonomousDatabaseConnectionUrls{SqlDevWebUrl: &pe},
		FreeformTags:                   map[string]string{"env": "prod"},
		DefinedTags:                    map[string]map[string]interface{}{"ns": {"k": "v"}},
		TimeCreated:                    &common.SDKTime{Time: created},
	}

	attrs := mapping.NewAutonomousDatabaseAttributesFromOCIAutonomousDatabase(oci)
	require.NotNil(t, attrs)
	// spot check a few attributes
	require.Equal(t, &id, attrs.ID)
	require.Equal(t, &name, attrs.Name)
	require.Equal(t, &comp, attrs.CompartmentOCID)
	require.Equal(t, string(database.AutonomousDatabaseLifecycleStateAvailable), attrs.LifecycleState)
	require.Equal(t, &dbver, attrs.DbVersion)
	require.Equal(t, string(workload), attrs.DbWorkload)
	require.Equal(t, string(lic), attrs.LicenseModel)
	require.Equal(t, &pe, attrs.PrivateEndpoint)
	require.Equal(t, &peip, attrs.PrivateEndpointIp)
	require.Equal(t, &pelabel, attrs.PrivateEndpointLabel)
	require.Equal(t, &subnet, attrs.SubnetId)
	require.Equal(t, nsgs, attrs.NsgIds)
	require.Equal(t, &ecpu, attrs.EcpuCount)
	require.Equal(t, &ocpu, attrs.OcpuCount)
	require.Equal(t, &cores, attrs.CpuCoreCount)
	require.Equal(t, &sizeTB, attrs.DataStorageSizeInTBs)
	require.Equal(t, &sizeGB, attrs.DataStorageSizeInGBs)
	require.Equal(t, &autoscale, attrs.IsAutoScalingEnabled)
	require.Equal(t, &autoscaleStorage, attrs.IsStorageAutoScalingEnabled)
	require.Equal(t, string(opsStatus), attrs.OperationsInsightsStatus)
	require.Equal(t, string(mgmtStatus), attrs.DatabaseManagementStatus)
	require.Equal(t, string(dataSafeStatus), attrs.DataSafeStatus)
	require.Equal(t, &dgEnabled, attrs.IsDataGuardEnabled)
	require.Equal(t, (*string)(&role), attrs.Role)
	require.Equal(t, peers, attrs.PeerAutonomousDbIds)
	require.Equal(t, map[string]string{"db": "conn"}, attrs.ConnectionStrings)
	require.Equal(t, map[string]string{"env": "prod"}, attrs.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, attrs.DefinedTags)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, created.Equal(attrs.TimeCreated.Time))

	// Now populate remaining attrs used only by domain conversion to avoid nil deref
	attrs.BackupRetentionDays = &backupDays
	attrs.LastBackupTime = &lastBackup
	attrs.LatestRestoreTime = &latestRestore
	attrs.PatchModel = (*string)(&patchModel)
	attrs.NextMaintenanceRunId = &nextMaint
	attrs.MaintenanceScheduleType = (*string)(&schedType)

	dom := mapping.NewDomainAutonomousDatabaseFromAttrs(attrs)
	require.IsType(t, &domain.AutonomousDatabase{}, dom)
	require.Equal(t, id, dom.ID)
	require.Equal(t, name, dom.Name)
	require.Equal(t, comp, dom.CompartmentOCID)
	require.Equal(t, string(database.AutonomousDatabaseLifecycleStateAvailable), dom.LifecycleState)
	require.Equal(t, dbver, dom.DbVersion)
	require.Equal(t, string(workload), dom.DbWorkload)
	require.Equal(t, string(lic), dom.LicenseModel)
	require.Equal(t, &free, dom.IsFreeTier)
	require.Equal(t, whitelisted, dom.WhitelistedIps)
	require.Equal(t, pe, dom.PrivateEndpoint)
	require.Equal(t, peip, dom.PrivateEndpointIp)
	require.Equal(t, pelabel, dom.PrivateEndpointLabel)
	require.Equal(t, subnet, dom.SubnetId)
	require.Equal(t, nsgs, dom.NsgIds)
	require.Equal(t, &mtls, dom.IsMtlsRequired)
	require.Equal(t, string(computeModel), dom.ComputeModel)
	require.Equal(t, &ecpu, dom.EcpuCount)
	require.Equal(t, &ocpu, dom.OcpuCount)
	require.Equal(t, &cores, dom.CpuCoreCount)
	require.Equal(t, &sizeTB, dom.DataStorageSizeInTBs)
	require.Equal(t, &sizeGB, dom.DataStorageSizeInGBs)
	require.Equal(t, &autoscale, dom.IsAutoScalingEnabled)
	require.Equal(t, &autoscaleStorage, dom.IsStorageAutoScalingEnabled)
	require.Equal(t, string(opsStatus), dom.OperationsInsightsStatus)
	require.Equal(t, string(mgmtStatus), dom.DatabaseManagementStatus)
	require.Equal(t, string(dataSafeStatus), dom.DataSafeStatus)
	require.Equal(t, &dgEnabled, dom.IsDataGuardEnabled)
	require.Equal(t, string(role), dom.Role)
	require.Equal(t, peers, dom.PeerAutonomousDbIds)
	require.Equal(t, &backupDays, dom.BackupRetentionDays)
	require.Equal(t, &lastBackup, dom.LastBackupTime)
	require.Equal(t, &latestRestore, dom.LatestRestoreTime)
	require.Equal(t, string(patchModel), dom.PatchModel)
	require.Equal(t, nextMaint, dom.NextMaintenanceRunId)
	require.Equal(t, string(schedType), dom.MaintenanceScheduleType)
	require.Equal(t, map[string]string{"db": "conn"}, dom.ConnectionStrings)
	require.NotNil(t, dom.ConnectionUrls)
	require.Equal(t, map[string]string{"env": "prod"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, dom.DefinedTags)
	require.NotNil(t, dom.TimeCreated)
	require.True(t, created.Equal(*dom.TimeCreated))
}
