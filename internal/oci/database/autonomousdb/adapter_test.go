package autonomousdb

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"
	domain "github.com/rozdolsky33/ocloud/internal/domain/database"
)

func sPtr(s string) *string     { return &s }
func bPtr(b bool) *bool         { return &b }
func f32Ptr(v float32) *float32 { return &v }
func iPtr(i int) *int           { return &i }

func TestToDomainAutonomousDB_FromSummary(t *testing.T) {
	ad := &Adapter{}
	created := time.Date(2024, 8, 2, 3, 4, 5, 0, time.UTC)
	oci := database.AutonomousDatabaseSummary{
		DbName:                   sPtr("adb1"),
		Id:                       sPtr("ocid1.autonomousdb.oc1..db"),
		CompartmentId:            sPtr("ocid1.compartment.oc1..comp"),
		LifecycleState:           "AVAILABLE",
		DbVersion:                sPtr("19c"),
		DbWorkload:               database.AutonomousDatabaseSummaryDbWorkloadOltp,
		LicenseModel:             database.AutonomousDatabaseSummaryLicenseModelLicenseIncluded,
		WhitelistedIps:           []string{"1.2.3.4/32"},
		PrivateEndpoint:          sPtr("10.0.0.5"),
		PrivateEndpointIp:        sPtr("10.0.0.5"),
		PrivateEndpointLabel:     sPtr("pep"),
		SubnetId:                 sPtr("ocid1.subnet.oc1..subnet"),
		NsgIds:                   []string{"nsg1"},
		IsMtlsConnectionRequired: bPtr(true),
		ComputeModel:             database.AutonomousDatabaseSummaryComputeModelEcpu,
		ComputeCount:             f32Ptr(2),
		OcpuCount:                f32Ptr(1),
		CpuCoreCount:             iPtr(1),
		DataStorageSizeInTBs:     iPtr(1),
		ConnectionStrings:        &database.AutonomousDatabaseConnectionStrings{AllConnectionStrings: map[string]string{"LOW": "conn"}},
		ConnectionUrls:           &database.AutonomousDatabaseConnectionUrls{SqlDevWebUrl: sPtr("https://sqldev")},
		FreeformTags:             map[string]string{"env": "dev"},
		DefinedTags:              map[string]map[string]interface{}{"ns": {"k": "v"}},
		TimeCreated:              &common.SDKTime{Time: created},
	}
	got := ad.toDomainAutonomousDB(oci)
	if got.Name != "adb1" || got.ID != "ocid1.autonomousdb.oc1..db" || got.CompartmentOCID != "ocid1.compartment.oc1..comp" || got.LifecycleState != "AVAILABLE" || got.DbVersion != "19c" || got.DbWorkload == "" || got.LicenseModel == "" || got.PrivateEndpoint != "10.0.0.5" || got.SubnetId != "ocid1.subnet.oc1..subnet" || got.IsMtlsRequired == nil || !*got.IsMtlsRequired || got.ComputeModel == "" || got.EcpuCount == nil || *got.EcpuCount != 2 || got.ConnectionStrings["LOW"] == "" || got.ConnectionUrls == nil || got.TimeCreated == nil || !got.TimeCreated.Equal(created) {
		t.Fatalf("toDomainAutonomousDB(summary) unexpected mapping: %#v", got)
	}
}

func TestToDomainAutonomousDB_FromFull(t *testing.T) {
	ad := &Adapter{}
	oci := database.AutonomousDatabase{
		DbName:               sPtr("full"),
		Id:                   sPtr("ocid1.autonomousdb.oc1..full"),
		LifecycleState:       database.AutonomousDatabaseLifecycleStateTerminated,
		DataStorageSizeInGBs: iPtr(1024),
	}
	got := ad.toDomainAutonomousDB(oci)
	want := domain.AutonomousDatabase{Name: "full", ID: "ocid1.autonomousdb.oc1..full", LifecycleState: string(database.AutonomousDatabaseLifecycleStateTerminated)}
	if got.Name != want.Name || got.ID != want.ID || got.LifecycleState != want.LifecycleState || got.DataStorageSizeInGBs == nil || *got.DataStorageSizeInGBs != 1024 {
		t.Fatalf("toDomainAutonomousDB(full) mismatch: got %#v", got)
	}
}
