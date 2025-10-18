package heatwavedb

import (
	"bytes"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/mysql"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

func ptrInt(v int) *int              { return &v }
func ptrBool(v bool) *bool           { return &v }
func ptrTime(t time.Time) *time.Time { return &t }
func ptrString(s string) *string     { return &s }

// Test summary view and JSON output for the list printer
func TestPrintHeatWaveDbsInfo_SummaryAndJSON(t *testing.T) {
	// Prepare test data
	db1 := database.HeatWaveDatabase{
		DisplayName:               "hw-prod-db",
		ID:                        "ocid1.mysqldbsystem.oc1..aaaa",
		LifecycleState:            "ACTIVE",
		MysqlVersion:              "8.0.35",
		ShapeName:                 "MySQL.VM.Standard.E4.4.128GB",
		DataStorageSizeInGBs:      ptrInt(100),
		IsHighlyAvailable:         ptrBool(true),
		IsHeatWaveClusterAttached: ptrBool(true),
		HeatWaveCluster: &mysql.HeatWaveClusterSummary{
			ClusterSize: ptrInt(2),
		},
		DatabaseMode: "READ_WRITE",
		AccessMode:   "UNRESTRICTED",
		SubnetId:     "ocid1.subnet.oc1..subnet1",
		SubnetName:   "hw-subnet",
		VcnID:        "ocid1.vcn.oc1..vcn1",
		VcnName:      "hw-vcn",
		NsgIds:       []string{"ocid1.nsg.oc1..nsg1"},
		NsgNames:     []string{"hw-nsg"},
		TimeCreated:  ptrTime(time.Date(2023, 6, 27, 20, 30, 46, 0, time.UTC)),
	}
	db2 := database.HeatWaveDatabase{
		DisplayName:               "hw-dev-db",
		ID:                        "ocid1.mysqldbsystem.oc1..bbbb",
		LifecycleState:            "INACTIVE",
		MysqlVersion:              "8.0.34",
		ShapeName:                 "MySQL.VM.Standard.E3.1.8GB",
		DataStorageSizeInGBs:      ptrInt(50),
		IsHighlyAvailable:         ptrBool(false),
		IsHeatWaveClusterAttached: ptrBool(false),
		DatabaseMode:              "READ_ONLY",
		AccessMode:                "RESTRICTED",
		SubnetId:                  "ocid1.subnet.oc1..subnet2",
		TimeCreated:               ptrTime(time.Date(2023, 5, 15, 10, 15, 30, 0, time.UTC)),
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	// Table output (summary)
	err := PrintHeatWaveDbsInfo([]database.HeatWaveDatabase{db1, db2}, appCtx, nil, false, false)
	assert.NoError(t, err)
	out := buf.String()
	// Validate presence of key fields for both DBs
	assert.Contains(t, out, "hw-prod-db")
	assert.Contains(t, out, "Lifecycle State")
	assert.Contains(t, out, "MySQL Version")
	assert.Contains(t, out, "Shape")
	assert.Contains(t, out, "Storage")
	assert.Contains(t, out, "High Availability")
	assert.Contains(t, out, "HeatWave Cluster")
	assert.Contains(t, out, "Database Mode")
	assert.Contains(t, out, "Access Mode")

	// JSON output with pagination
	buf.Reset()
	pg := &util.PaginationInfo{TotalCount: 2, Limit: 2, CurrentPage: 1, NextPageToken: ""}
	err = PrintHeatWaveDbsInfo([]database.HeatWaveDatabase{db1, db2}, appCtx, pg, true, false)
	assert.NoError(t, err)
	jsonOut := buf.String()
	assert.Contains(t, jsonOut, `"items"`)
	assert.Contains(t, jsonOut, `"pagination"`)
	assert.Contains(t, jsonOut, `"hw-prod-db"`)
	assert.Contains(t, jsonOut, `"hw-dev-db"`)
}

// Test detailed view output
func TestPrintHeatWaveDbsInfo_DetailedView(t *testing.T) {
	clusterSize := 3
	shapeName := "MySQL.HeatWave.VM.Standard.E3"
	db := database.HeatWaveDatabase{
		DisplayName:               "hw-prod-detailed",
		ID:                        "ocid1.mysqldbsystem.oc1..cccc",
		LifecycleState:            "ACTIVE",
		MysqlVersion:              "8.0.35",
		Description:               "Production HeatWave Database",
		ShapeName:                 "MySQL.VM.Standard.E4.8.256GB",
		DataStorageSizeInGBs:      ptrInt(500),
		IsHighlyAvailable:         ptrBool(true),
		IsHeatWaveClusterAttached: ptrBool(true),
		HeatWaveCluster: &mysql.HeatWaveClusterSummary{
			ClusterSize:    &clusterSize,
			ShapeName:      &shapeName,
			LifecycleState: mysql.HeatWaveClusterLifecycleStateActive,
		},
		DatabaseMode:       "READ_WRITE",
		AccessMode:         "UNRESTRICTED",
		SubnetId:           "ocid1.subnet.oc1..subnet1",
		SubnetName:         "hw-subnet",
		VcnID:              "ocid1.vcn.oc1..vcn1",
		VcnName:            "hw-vcn",
		NsgIds:             []string{"ocid1.nsg.oc1..nsg1", "ocid1.nsg.oc1..nsg2"},
		NsgNames:           []string{"hw-nsg-1", "hw-nsg-2"},
		AvailabilityDomain: "AD-1",
		FaultDomain:        "FAULT-DOMAIN-1",
		ConfigurationId:    "ocid1.mysqlconfiguration.oc1..config1",
		BackupPolicy: &mysql.BackupPolicy{
			IsEnabled:       ptrBool(true),
			RetentionInDays: ptrInt(30),
			WindowStartTime: ptrString("02:00"),
			FreeformTags:    map[string]string{"env": "prod"},
			DefinedTags:     map[string]map[string]interface{}{"operations": {"team": "dba"}},
		},
		Endpoints: []mysql.DbSystemEndpoint{
			{
				IpAddress: ptrString("10.0.1.10"),
				Hostname:  ptrString("hw-prod-detailed.mysql.oraclecloud.com"),
			},
		},
		TimeCreated: ptrTime(time.Date(2023, 6, 27, 20, 30, 46, 0, time.UTC)),
		TimeUpdated: ptrTime(time.Date(2023, 8, 15, 14, 20, 10, 0, time.UTC)),
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	// Detailed view (showAll=true)
	err := PrintHeatWaveDbsInfo([]database.HeatWaveDatabase{db}, appCtx, nil, false, true)
	assert.NoError(t, err)
	out := buf.String()

	// Validate detailed fields
	assert.Contains(t, out, "hw-prod-detailed")
	assert.Contains(t, out, "Production HeatWave Database")
	assert.Contains(t, out, "HeatWave Shape")
	assert.Contains(t, out, "HeatWave State")
	assert.Contains(t, out, "Configuration ID")
	assert.Contains(t, out, "Availability Domain")
	assert.Contains(t, out, "Fault Domain")
	assert.Contains(t, out, "Automatic Backups")
	assert.Contains(t, out, "Retention Days")
	assert.Contains(t, out, "Endpoint 1 IP")
	assert.Contains(t, out, "Endpoint 1 Hostname")
	assert.Contains(t, out, "Time Created")
	assert.Contains(t, out, "Time Updated")
}

// Test empty list handling
func TestPrintHeatWaveDbsInfo_EmptyList(t *testing.T) {
	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	err := PrintHeatWaveDbsInfo([]database.HeatWaveDatabase{}, appCtx, nil, false, false)
	assert.NoError(t, err)
	// The ValidateAndReportEmpty should handle empty lists
	out := buf.String()
	assert.NotEmpty(t, out)
}

// Test JSON output for empty list with pagination
func TestPrintHeatWaveDbsInfo_EmptyJSON(t *testing.T) {
	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	err := PrintHeatWaveDbsInfo([]database.HeatWaveDatabase{}, appCtx, nil, true, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "{}")
}

// Test single database output
func TestPrintHeatWaveDbInfo_Single(t *testing.T) {
	db := database.HeatWaveDatabase{
		DisplayName:               "hw-single-db",
		ID:                        "ocid1.mysqldbsystem.oc1..dddd",
		LifecycleState:            "ACTIVE",
		MysqlVersion:              "8.0.35",
		ShapeName:                 "MySQL.VM.Standard.E4.4.128GB",
		DataStorageSizeInGBs:      ptrInt(100),
		IsHighlyAvailable:         ptrBool(true),
		IsHeatWaveClusterAttached: ptrBool(false),
		DatabaseMode:              "READ_WRITE",
		AccessMode:                "UNRESTRICTED",
		SubnetId:                  "ocid1.subnet.oc1..subnet1",
		TimeCreated:               ptrTime(time.Date(2023, 6, 27, 20, 30, 46, 0, time.UTC)),
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	// Summary view
	err := PrintHeatWaveDbInfo(&db, appCtx, false, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "hw-single-db")
	assert.Contains(t, out, "ACTIVE")

	// JSON output
	buf.Reset()
	err = PrintHeatWaveDbInfo(&db, appCtx, true, false)
	assert.NoError(t, err)
	jsonOut := buf.String()
	assert.Contains(t, jsonOut, `"hw-single-db"`)
	assert.Contains(t, jsonOut, `"ACTIVE"`)
}

// Test boolToString helper
func TestBoolToString(t *testing.T) {
	assert.Equal(t, "true", boolToString(ptrBool(true)))
	assert.Equal(t, "false", boolToString(ptrBool(false)))
	assert.Equal(t, "", boolToString(nil))
}
