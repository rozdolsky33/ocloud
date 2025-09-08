package autonomousdb

import (
	"bytes"
	"testing"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

func ptrInt(v int) *int              { return &v }
func ptrF32(v float32) *float32      { return &v }
func ptrBool(v bool) *bool           { return &v }
func ptrTime(t time.Time) *time.Time { return &t }

// Test summary view and JSON output for the list printer
func TestPrintAutonomousDbsInfo_SummaryAndJSON(t *testing.T) {
	// Prepare test data: one ECPU-based and one OCPU-based DB
	db1 := domain.AutonomousDatabase{
		Name:                 "AuthzSessionManagementDev",
		ID:                   "ocid1.autonomousdatabase.oc1..aaaa",
		LifecycleState:       "AVAILABLE",
		DbVersion:            "19c",
		DbWorkload:           "OLTP",
		LicenseModel:         "BRING_YOUR_OWN_LICENSE",
		ComputeModel:         "ECPU",
		EcpuCount:            ptrF32(2.0),
		DataStorageSizeInTBs: ptrInt(1),
		PrivateEndpointIp:    "10.0.0.10",
		PrivateEndpoint:      "db1.adb.us-ashburn-1.oraclecloud.com",
		SubnetId:             "ocid1.subnet.oc1..subnet1",
		NsgIds:               []string{"ocid1.nsg.oc1..nsg1"},
		ConnectionStrings: map[string]string{
			"HIGH":     "high-conn",
			"MEDIUM":   "medium-conn",
			"LOW":      "low-conn",
			"TP":       "tp-conn",
			"TPURGENT": "tpurgent-conn",
		},
		TimeCreated: ptrTime(time.Date(2023, 6, 27, 20, 30, 46, 0, time.UTC)),
	}
	db2 := domain.AutonomousDatabase{
		Name:                 "BillingDB",
		ID:                   "ocid1.autonomousdatabase.oc1..bbbb",
		LifecycleState:       "STOPPED",
		DbVersion:            "23ai",
		DbWorkload:           "DW",
		LicenseModel:         "LICENSE_INCLUDED",
		ComputeModel:         "OCPU",
		OcpuCount:            ptrF32(1.0),
		CpuCoreCount:         ptrInt(2),
		DataStorageSizeInTBs: ptrInt(2),
		PrivateEndpointIp:    "10.0.0.20",
		PrivateEndpoint:      "db2.adb.us-ashburn-1.oraclecloud.com",
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	// Table output (summary)
	err := PrintAutonomousDbsInfo([]domain.AutonomousDatabase{db1, db2}, appCtx, nil, false, false)
	assert.NoError(t, err)
	out := buf.String()
	// Validate presence of key fields for both DBs
	assert.Contains(t, out, "AuthzSessionManagementDev")
	assert.Contains(t, out, "Lifecycle State")
	assert.Contains(t, out, "Compute Model")
	assert.Contains(t, out, "ECPUs")
	assert.Contains(t, out, "Storage")
	assert.Contains(t, out, "Private IP")
	assert.Contains(t, out, "Private Endpoint")
	// Summary view does not include connection strings; they are shown in detailed view

	// JSON output with pagination
	buf.Reset()
	pg := &util.PaginationInfo{TotalCount: 2, Limit: 2, CurrentPage: 1, NextPageToken: ""}
	err = PrintAutonomousDbsInfo([]domain.AutonomousDatabase{db1, db2}, appCtx, pg, true, false)
	assert.NoError(t, err)
	jsonOut := buf.String()
	assert.Contains(t, jsonOut, "\"items\"")
	assert.Contains(t, jsonOut, "AuthzSessionManagementDev")
	assert.Contains(t, jsonOut, "BillingDB")
	assert.Contains(t, jsonOut, "\"pagination\"")
}

// Test detailed view including access type inference and storage TB/GB selection
func TestPrintAutonomousDbInfo_DetailedAccessTypeStorage(t *testing.T) {
	db := &domain.AutonomousDatabase{
		Name:                 "DetailedDB",
		ID:                   "ocid1.autonomousdatabase.oc1..cccc",
		LifecycleState:       "AVAILABLE",
		DbVersion:            "19c",
		DbWorkload:           "OLTP",
		LicenseModel:         "BRING_YOUR_OWN_LICENSE",
		ComputeModel:         "ECPU",
		EcpuCount:            ptrF32(2.0),
		DataStorageSizeInGBs: ptrInt(100), // force GB path in Capacity section
		PrivateEndpoint:      "detailed.adb.us-ashburn-1.oraclecloud.com",
		PrivateEndpointIp:    "10.0.0.30",
		PrivateEndpointLabel: "DetailedDB",
		IsMtlsRequired:       ptrBool(false),
		ConnectionStrings: map[string]string{
			"HIGH":   "high-conn",
			"MEDIUM": "medium-conn",
			"LOW":    "low-conn",
		},
	}
	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	err := PrintAutonomousDbInfo(db, appCtx, false, true)
	assert.NoError(t, err)
	out := buf.String()

	// Access Type should be inferred as Virtual cloud network when private endpoint exists and public flag is nil
	assert.Contains(t, out, "Access Type")
	assert.Contains(t, out, "Virtual cloud network")

	// Capacity should show ECPUs and storage in GB
	assert.Contains(t, out, "ECPUs")
	assert.Contains(t, out, "2.00")
	assert.Contains(t, out, "Storage")
	assert.Contains(t, out, "100 GB")

	// Connection strings should be present
	assert.Contains(t, out, "high-conn")
	assert.Contains(t, out, "medium-conn")
	assert.Contains(t, out, "low-conn")
}

// Test empty handling for list printer
func TestPrintAutonomousDbsInfo_Empty(t *testing.T) {
	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	err := PrintAutonomousDbsInfo([]domain.AutonomousDatabase{}, appCtx, nil, false, true)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "No Items found.")
}
