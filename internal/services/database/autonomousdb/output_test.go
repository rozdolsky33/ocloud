package autonomousdb

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

// TestPrintAutonomousDbInfo tests the PrintAutonomousDbInfo function
func TestPrintAutonomousDbInfo(t *testing.T) {
	// Create test databases
	databases := []AutonomousDatabase{
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
	err := PrintAutonomousDbInfo(databases, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestDatabase1")
	assert.Contains(t, output, "test-endpoint-1.example.com")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "high-connection-string-1")
	assert.Contains(t, output, "medium-connection-string-1")
	assert.Contains(t, output, "low-connection-string-1")
	assert.Contains(t, output, "TestDatabase2")
	assert.Contains(t, output, "test-endpoint-2.example.com")
	assert.Contains(t, output, "192.168.1.2")
	assert.Contains(t, output, "high-connection-string-2")
	assert.Contains(t, output, "medium-connection-string-2")
	assert.Contains(t, output, "low-connection-string-2")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, nil, true)
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
	databases := []AutonomousDatabase{}

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
	err := PrintAutonomousDbInfo(databases, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty object
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "{}")
}

// TestPrintAutonomousDbInfoWithPagination tests the PrintAutonomousDbInfo function with pagination
func TestPrintAutonomousDbInfoWithPagination(t *testing.T) {
	// Create test databases
	databases := []AutonomousDatabase{
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
	err := PrintAutonomousDbInfo(databases, appCtx, pagination, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestDatabase1")
	assert.Contains(t, output, "test-endpoint-1.example.com")
	assert.Contains(t, output, "192.168.1.1")
	assert.Contains(t, output, "high-connection-string-1")
	assert.Contains(t, output, "medium-connection-string-1")
	assert.Contains(t, output, "low-connection-string-1")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintAutonomousDbInfo(databases, appCtx, pagination, true)
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
