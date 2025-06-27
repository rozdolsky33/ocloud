package compartment

import (
	"bytes"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPrintCompartmentsInfo tests the PrintCompartmentsInfo function
func TestPrintCompartmentsInfo(t *testing.T) {
	// Create test compartments
	compartments := []Compartment{
		{
			Name:        "TestCompartment1",
			ID:          "ocid1.compartment.oc1.phx.test1",
			Description: "Test compartment 1 description",
		},
		{
			Name:        "TestCompartment2",
			ID:          "ocid1.compartment.oc1.phx.test2",
			Description: "Test compartment 2 description",
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		TenancyName: "TestTenancy",
		TenancyID:   "ocid1.tenancy.oc1.phx.test",
		Logger:      logger.NewTestLogger(),
		Stdout:      &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintCompartmentsInfo(compartments, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestCompartment1")
	assert.Contains(t, output, "ocid1.compartment.oc1.phx.test1")
	assert.Contains(t, output, "Test compartment 1 description")
	assert.Contains(t, output, "TestCompartment2")
	assert.Contains(t, output, "ocid1.compartment.oc1.phx.test2")
	assert.Contains(t, output, "Test compartment 2 description")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintCompartmentsInfo(compartments, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestCompartment1")
	assert.Contains(t, jsonOutput, "ocid1.compartment.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "Test compartment 1 description")
	assert.Contains(t, jsonOutput, "TestCompartment2")
	assert.Contains(t, jsonOutput, "ocid1.compartment.oc1.phx.test2")
	assert.Contains(t, jsonOutput, "Test compartment 2 description")
}

// TestPrintCompartmentsInfoEmpty tests the PrintCompartmentsInfo function with empty compartments
func TestPrintCompartmentsInfoEmpty(t *testing.T) {
	// Create an empty compartments slice
	compartments := []Compartment{}

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		TenancyName: "TestTenancy",
		TenancyID:   "ocid1.tenancy.oc1.phx.test",
		Logger:      logger.NewTestLogger(),
		Stdout:      &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintCompartmentsInfo(compartments, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintCompartmentsInfo(compartments, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty object
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "{}")
}

// TestPrintCompartmentsInfoWithPagination tests the PrintCompartmentsInfo function with pagination
func TestPrintCompartmentsInfoWithPagination(t *testing.T) {
	// Create test compartments
	compartments := []Compartment{
		{
			Name:        "TestCompartment1",
			ID:          "ocid1.compartment.oc1.phx.test1",
			Description: "Test compartment 1 description",
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
		TenancyName: "TestTenancy",
		TenancyID:   "ocid1.tenancy.oc1.phx.test",
		Logger:      logger.NewTestLogger(),
		Stdout:      &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintCompartmentsInfo(compartments, appCtx, pagination, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestCompartment1")
	assert.Contains(t, output, "ocid1.compartment.oc1.phx.test1")
	assert.Contains(t, output, "Test compartment 1 description")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintCompartmentsInfo(compartments, appCtx, pagination, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestCompartment1")
	assert.Contains(t, jsonOutput, "ocid1.compartment.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "Test compartment 1 description")
	assert.Contains(t, jsonOutput, "\"pagination\"")
	assert.Contains(t, jsonOutput, "\"TotalCount\"")
	assert.Contains(t, jsonOutput, "\"Limit\"")
	assert.Contains(t, jsonOutput, "\"CurrentPage\"")
	assert.Contains(t, jsonOutput, "\"NextPageToken\"")
}
