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
	// Configuration ID removed - too technical for summary view
	assert.Contains(t, out, "Availability Domain")
	assert.Contains(t, out, "Fault Domain")
	assert.Contains(t, out, "Automatic Backups")
	assert.Contains(t, out, "Retention Days")
	// Endpoint consolidated - removed duplicate endpoint fields
	assert.Contains(t, out, "Endpoint")
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

// Test getMySQLShapeDetails helper
func TestGetMySQLShapeDetails(t *testing.T) {
	tests := []struct {
		name       string
		shapeName  string
		wantECPU   int
		wantMemory int
		wantFound  bool
	}{
		{
			name:       "MySQL.Free shape",
			shapeName:  "MySQL.Free",
			wantECPU:   2,
			wantMemory: 8,
			wantFound:  true,
		},
		{
			name:       "MySQL.2 shape",
			shapeName:  "MySQL.2",
			wantECPU:   2,
			wantMemory: 16,
			wantFound:  true,
		},
		{
			name:       "MySQL.4 shape",
			shapeName:  "MySQL.4",
			wantECPU:   4,
			wantMemory: 32,
			wantFound:  true,
		},
		{
			name:       "MySQL.8 shape",
			shapeName:  "MySQL.8",
			wantECPU:   8,
			wantMemory: 64,
			wantFound:  true,
		},
		{
			name:       "MySQL.16 shape",
			shapeName:  "MySQL.16",
			wantECPU:   16,
			wantMemory: 128,
			wantFound:  true,
		},
		{
			name:       "MySQL.32 shape",
			shapeName:  "MySQL.32",
			wantECPU:   32,
			wantMemory: 256,
			wantFound:  true,
		},
		{
			name:       "MySQL.48 shape",
			shapeName:  "MySQL.48",
			wantECPU:   48,
			wantMemory: 384,
			wantFound:  true,
		},
		{
			name:       "MySQL.64 shape",
			shapeName:  "MySQL.64",
			wantECPU:   64,
			wantMemory: 512,
			wantFound:  true,
		},
		{
			name:       "MySQL.256 shape",
			shapeName:  "MySQL.256",
			wantECPU:   256,
			wantMemory: 1024,
			wantFound:  true,
		},
		{
			name:       "Unknown shape",
			shapeName:  "MySQL.Unknown",
			wantECPU:   0,
			wantMemory: 0,
			wantFound:  false,
		},
		{
			name:       "Legacy OCPU shape",
			shapeName:  "MySQL.VM.Standard.E4.4.128GB",
			wantECPU:   0,
			wantMemory: 0,
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ecpu, memory, found := getMySQLShapeDetails(tt.shapeName)
			assert.Equal(t, tt.wantECPU, ecpu, "ECPU count mismatch")
			assert.Equal(t, tt.wantMemory, memory, "Memory mismatch")
			assert.Equal(t, tt.wantFound, found, "Found flag mismatch")
		})
	}
}

// Test storage size formatting with ECPU shapes based on actual API response
func TestPrintHeatWaveDbInfo_ECPUShapeAndStorage(t *testing.T) {
	autoExpand := false
	db := database.HeatWaveDatabase{
		DisplayName:          "prod-orion-midtier",
		ID:                   "ocid1.mysqldbsystem.oc2.us-luke-1.aaaaaaaa5cm2tpeirggrcvzbep6nicruo3v233emfgguk32n4kaydrv5vyza",
		LifecycleState:       "ACTIVE",
		MysqlVersion:         "8.4.6",
		Description:          "Prod Orion Midtier MySQL Database service",
		ShapeName:            "MySQL.4", // 4 ECPUs, 32 GB memory (per official Oracle table)
		DataStorageSizeInGBs: ptrInt(1024),
		DataStorage: &mysql.DataStorage{
			IsAutoExpandStorageEnabled: &autoExpand,
			MaxStorageSizeInGBs:        nil,
			AllocatedStorageSizeInGBs:  ptrInt(1024),  // 1 TiB allocated
			DataStorageSizeInGBs:       ptrInt(1024),  // 1 TiB used
			DataStorageSizeLimitInGBs:  ptrInt(98304), // 96 TiB maximum limit
		},
		IsHighlyAvailable:         ptrBool(true),
		IsHeatWaveClusterAttached: ptrBool(false),
		DatabaseMode:              "READ_WRITE",
		AccessMode:                "UNRESTRICTED",
		IpAddress:                 "217.142.42.138",
		Port:                      ptrInt(3306),
		PortX:                     ptrInt(33060),
		SubnetId:                  "ocid1.subnet.oc2.us-luke-1.aaaaaaaaciwbptgw5zr6tha4hkwhmbbh5symeeus673oqiilfefltmadkh6q",
		SubnetName:                "database",
		VcnID:                     "ocid1.vcn.oc2.us-luke-1.amaaaaaalkxuqyqaepjbdrld74cdsr4crmazexzv5o3bjbp7cxoidelbxeqq",
		VcnName:                   "vcn-luf-rho-udeprod1-1",
		AvailabilityDomain:        "PQsp:us-luke-1-ad-1",
		FaultDomain:               "FAULT-DOMAIN-3",
		TimeCreated:               ptrTime(time.Date(2025, 10, 8, 19, 49, 26, 0, time.UTC)),
		TimeUpdated:               ptrTime(time.Date(2025, 10, 18, 9, 41, 39, 0, time.UTC)),
	}

	var buf bytes.Buffer
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: &buf}

	// Summary view
	err := PrintHeatWaveDbInfo(&db, appCtx, false, false)
	assert.NoError(t, err)
	out := buf.String()

	// Validate ECPU shape info is displayed correctly per Oracle Table 5-1
	assert.Contains(t, out, "MySQL.4")
	assert.Contains(t, out, "4")     // ECPU count (MySQL.4 has 4 ECPUs per official table)
	assert.Contains(t, out, "32 GB") // Memory (4 ECPUs × 8 GB = 32 GB)

	// Validate storage is shown in TiB (1024 GB = 1.00 TiB)
	assert.Contains(t, out, "1.00 TiB")
	// Validate storage limit is shown (98304 GB = 96 TiB)
	assert.Contains(t, out, "96.00 TiB")
	// Validate auto-expand is shown
	assert.Contains(t, out, "Auto-Expand Storage")
	assert.Contains(t, out, "false")

	// Detailed view
	buf.Reset()
	err = PrintHeatWaveDbInfo(&db, appCtx, false, true)
	assert.NoError(t, err)
	detailedOut := buf.String()

	// Validate detailed view also shows correct info
	assert.Contains(t, detailedOut, "ECPUs")
	assert.Contains(t, detailedOut, "Memory")
	// Storage should be consolidated when used == allocated
	assert.Contains(t, detailedOut, "Storage")
	assert.Contains(t, detailedOut, "Storage Limit")
	assert.Contains(t, detailedOut, "1.00 TiB")
	assert.Contains(t, detailedOut, "96.00 TiB")
	assert.Contains(t, detailedOut, "Prod Orion Midtier MySQL Database service")
	// Endpoint should be consolidated
	assert.Contains(t, detailedOut, "Endpoint")
}
