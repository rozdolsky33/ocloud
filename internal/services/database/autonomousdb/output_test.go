package autonomousdb

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

// TestPrintAutonomousDbInfo tests the PrintAutonomousDbInfo function
func TestPrintAutonomousDbInfo(t *testing.T) {
	// Create test databases
	databases := []domain.AutonomousDatabase{
		{
			Name:              "TestDatabase1",
			ID:                "ocid1.autonomousdatabase.oc1.phx.test1",
			PrivateEndpoint:   "test-endpoint-1.example.com",
			PrivateEndpointIp: "192.168.1.1",
			ConnectionStrings: map[string]string{
				"HIGH":   "high-connection-string-1",
				"MEDIUM": "medium-connection-string-1",
				"LOW":    "low-connection-string-1",
			},
		},
		{
			Name:              "TestDatabase2",
			ID:                "ocid1.autonomousdatabase.oc1.phx.test2",
			PrivateEndpoint:   "test-endpoint-2.example.com",
			PrivateEndpointIp: "192.168.1.2",
			ConnectionStrings: map[string]string{
				"HIGH":   "high-connection-string-2",
				"MEDIUM": "medium-connection-string-2",
				"LOW":    "low-connection-string-2",
			},
		},
	}

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
		Stdout:          &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintAutonomousDbInfo(databases, appCtx, nil, false, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestDatabase1")
	assert.Contains(t, output, "test-endpoint-1.example.com")
	assert.Contains(t, output, "192.168.1.1")
	assert.NotContains(t, output, "high-connection-string-1")
	assert.NotContains(t, output, "medium-connection-string-1")
	assert.NotContains(t, output, "low-connection-string-1")
	assert.Contains(t, output, "TestDatabase2")
	assert.Contains(t, output, "test-endpoint-2.example.com")
	assert.Contains(t, output, "192.168.1.2")
	assert.NotContains(t, output, "high-connection-string-2")
	assert.NotContains(t, output, "medium-connection-string-2")
	assert.NotContains(t, output, "low-connection-string-2")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, nil, true, false)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestDatabase1")
	assert.Contains(t, jsonOutput, "ocid1.autonomousdatabase.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "test-endpoint-1.example.com")
	assert.Contains(t, jsonOutput, "192.168.1.1")
	assert.Contains(t, jsonOutput, "high-connection-string-1")
	assert.Contains(t, jsonOutput, "medium-connection-string-1")
	assert.Contains(t, jsonOutput, "low-connection-string-1")
	assert.Contains(t, jsonOutput, "TestDatabase2")
	assert.Contains(t, jsonOutput, "ocid1.autonomousdatabase.oc1.phx.test2")
	assert.Contains(t, jsonOutput, "test-endpoint-2.example.com")
	assert.Contains(t, jsonOutput, "192.168.1.2")
	assert.Contains(t, jsonOutput, "high-connection-string-2")
	assert.Contains(t, jsonOutput, "medium-connection-string-2")
	assert.Contains(t, jsonOutput, "low-connection-string-2")
}

// TestPrintAutonomousDbInfoEmpty tests the PrintAutonomousDbInfo function with empty databases
func TestPrintAutonomousDbInfoEmpty(t *testing.T) {
	// Create an empty databases slice
	databases := []domain.AutonomousDatabase{}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
		Stdout:          &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintAutonomousDbInfo(databases, appCtx, nil, false, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, nil, true, false)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty object
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "{}")
}

// TestPrintAutonomousDbInfoWithPagination tests the PrintAutonomousDbInfo function with pagination
func TestPrintAutonomousDbInfoWithPagination(t *testing.T) {
	// Create test databases
	databases := []domain.AutonomousDatabase{
		{
			Name:              "TestDatabase1",
			ID:                "ocid1.autonomousdatabase.oc1.phx.test1",
			PrivateEndpoint:   "test-endpoint-1.example.com",
			PrivateEndpointIp: "192.168.1.1",
			ConnectionStrings: map[string]string{
				"HIGH":   "high-connection-string-1",
				"MEDIUM": "medium-connection-string-1",
				"LOW":    "low-connection-string-1",
			},
		},
	}

	// Create pagination info
	pagination := &util.PaginationInfo{
		TotalCount:    10,
		Limit:         1,
		CurrentPage:   1,
		NextPageToken: "next-page-token",
	}

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
		Stdout:          &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintAutonomousDbInfo(databases, appCtx, pagination, false, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestDatabase1")
	assert.Contains(t, output, "test-endpoint-1.example.com")
	assert.Contains(t, output, "192.168.1.1")
	assert.NotContains(t, output, "high-connection-string-1")
	assert.NotContains(t, output, "medium-connection-string-1")
	assert.NotContains(t, output, "low-connection-string-1")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, pagination, true, false)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestDatabase1")
	assert.Contains(t, jsonOutput, "ocid1.autonomousdatabase.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "test-endpoint-1.example.com")
	assert.Contains(t, jsonOutput, "192.168.1.1")
	assert.Contains(t, jsonOutput, "high-connection-string-1")
	assert.Contains(t, jsonOutput, "medium-connection-string-1")
	assert.Contains(t, jsonOutput, "low-connection-string-1")
	assert.Contains(t, jsonOutput, "\"pagination\"")
	assert.Contains(t, jsonOutput, "\"TotalCount\"")
	assert.Contains(t, jsonOutput, "\"Limit\"")
	assert.Contains(t, jsonOutput, "\"CurrentPage\"")
	assert.Contains(t, jsonOutput, "\"NextPageToken\"")
}

// TestPrintAutonomousDbInfoShowAll tests the PrintAutonomousDbInfo function with the showAll flag set to true.
func TestPrintAutonomousDbInfoShowAll(t *testing.T) {
	databases := []domain.AutonomousDatabase{
		{
			Name:                 "TestDatabase1",
			ID:                   "ocid1.autonomousdatabase.oc1.phx.test1",
			LifecycleState:       "AVAILABLE",
			DbVersion:            "19c",
			DbWorkload:           "OLTP",
			LicenseModel:         "BRING_YOUR_OWN_LICENSE",
			ComputeModel:         "ECPU",
			EcpuCount:            float32Ptr(2.0),
			DataStorageSizeInTBs: intPtr(1),
			IsAutoScalingEnabled: boolPtr(true),
			PrivateEndpoint:      "test-endpoint-1.example.com",
			PrivateEndpointIp:    "192.168.1.1",
			SubnetName:           "Subnet-1",
			VcnName:              "VCN-1",
			NsgNames:             []string{"nsg1", "nsg2"},
			IsMtlsRequired:       boolPtr(false),
			ConnectionStrings: map[string]string{
				"HIGH":     "high-connection-string-1",
				"MEDIUM":   "medium-connection-string-1",
				"LOW":      "low-connection-string-1",
				"TP":       "tp-connection-string-1",
				"TPURGENT": "tpurgent-connection-string-1",
			},
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
		Stdout:          &buf,
	}

	// Test with table output (showAll = true)
	err := PrintAutonomousDbInfo(databases, appCtx, nil, false, true)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestDatabase1")
	assert.Contains(t, output, "Lifecycle State")
	assert.Contains(t, output, "AVAILABLE")
	assert.Contains(t, output, "19c")
	assert.Contains(t, output, "OLTP")
	assert.Contains(t, output, "BRING_YOUR_OWN_LICENSE")
	assert.Contains(t, output, "Compute Model")
	assert.Contains(t, output, "ECPU")
	assert.Contains(t, output, "ECPUs")
	assert.Contains(t, output, "2.00")
	assert.Contains(t, output, "Storage")
	assert.Contains(t, output, "1 TB")
	assert.Contains(t, output, "Auto Scaling")
	assert.Contains(t, output, "true")
	assert.Contains(t, output, "Access Type")
	assert.Contains(t, output, "Virtual cloud network")
	assert.Contains(t, output, "Private IP")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "Private Endpoint")
	assert.Contains(t, output, "test-endpoint-1.example.com")
	assert.Contains(t, output, "Subnet")
	assert.Contains(t, output, "Subnet-1")
	assert.Contains(t, output, "VCN")
	assert.Contains(t, output, "VCN-1")
	assert.Contains(t, output, "NSGs")
	assert.Contains(t, output, "[nsg1 nsg2]")
	assert.Contains(t, output, "mTLS Required")
	assert.Contains(t, output, "false")
	assert.Contains(t, output, "High")
	assert.Contains(t, output, "high-connection-string-1")
	assert.Contains(t, output, "Medium")
	assert.Contains(t, output, "medium-connection-string-1")
	assert.Contains(t, output, "Low")
	assert.Contains(t, output, "low-connection-string-1")
	assert.Contains(t, output, "TP")
	assert.Contains(t, output, "tp-connection-string-1")
	assert.Contains(t, output, "TPURGENT")
	assert.Contains(t, output, "tpurgent-connection-string-1")
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func float32Ptr(f float32) *float32 {
	return &f
}
