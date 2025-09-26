package policy

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

// TestPrintPolicyInfo tests the PrintPolicyInfo function
func TestPrintPolicyInfo(t *testing.T) {
	// Create test policies
	policies := []identity.Policy{
		{
			Name:        "TestPolicy1",
			ID:          "ocid1.policy.oc1.phx.test1",
			Description: "Test policy 1 description",
			Statement:   []string{"Allow group Administrators to manage all-resources in tenancy"},
			DefinedTags: map[string]map[string]interface{}{
				"Operations": {
					"Environment": "Production",
				},
			},
			FreeformTags: map[string]string{
				"Department": "IT",
			},
		},
		{
			Name:        "TestPolicy2",
			ID:          "ocid1.policy.oc1.phx.test2",
			Description: "Test policy 2 description",
			Statement:   []string{"Allow group NetworkAdmins to manage virtual-network-family in tenancy"},
			DefinedTags: map[string]map[string]interface{}{
				"Operations": {
					"Environment": "Production",
				},
			},
			FreeformTags: map[string]string{
				"Department": "IT",
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

	// Test with table output (useJSON = false)
	err := PrintPolicyInfo(policies, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestPolicy1")
	assert.Contains(t, output, "ocid1.policy.oc1.phx.test1")
	assert.Contains(t, output, "Test policy 1 description")
	assert.Contains(t, output, "TestPolicy2")
	assert.Contains(t, output, "ocid1.policy.oc1.phx.test2")
	assert.Contains(t, output, "Test policy 2 description")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintPolicyInfo(policies, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestPolicy1")
	assert.Contains(t, jsonOutput, "ocid1.policy.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "Test policy 1 description")
	assert.Contains(t, jsonOutput, "TestPolicy2")
	assert.Contains(t, jsonOutput, "ocid1.policy.oc1.phx.test2")
	assert.Contains(t, jsonOutput, "Test policy 2 description")
}

// TestPrintPolicyInfoEmpty tests the PrintPolicyInfo function with empty policies
func TestPrintPolicyInfoEmpty(t *testing.T) {
	// Create an empty policies slice
	policies := []identity.Policy{}

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
	err := PrintPolicyInfo(policies, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintPolicyInfo(policies, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates no items
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "[]")
}

// TestPrintPolicyInfoWithPagination tests the PrintPolicyInfo function with pagination
func TestPrintPolicyInfoWithPagination(t *testing.T) {
	// Create test policies
	policies := []identity.Policy{
		{
			Name:        "TestPolicy1",
			ID:          "ocid1.policy.oc1.phx.test1",
			Description: "Test policy 1 description",
			Statement:   []string{"Allow group Administrators to manage all-resources in tenancy"},
		},
		{
			Name:        "TestPolicy2",
			ID:          "ocid1.policy.oc1.phx.test2",
			Description: "Test policy 2 description",
			Statement:   []string{"Allow group NetworkAdmins to manage virtual-network-family in tenancy"},
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

	// Create pagination info
	pagination := &util.PaginationInfo{
		CurrentPage:   2,
		TotalCount:    10,
		Limit:         5,
		NextPageToken: "next-page-token",
	}

	// Test with table output (useJSON = false)
	err := PrintPolicyInfo(policies, appCtx, pagination, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestPolicy1")
	assert.Contains(t, output, "ocid1.policy.oc1.phx.test1")
	assert.Contains(t, output, "Test policy 1 description")
	assert.Contains(t, output, "TestPolicy2")
	assert.Contains(t, output, "ocid1.policy.oc1.phx.test2")
	assert.Contains(t, output, "Test policy 2 description")
	assert.Contains(t, output, "Page 2")
	assert.Contains(t, output, "Total: 10")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintPolicyInfo(policies, appCtx, pagination, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestPolicy1")
	assert.Contains(t, jsonOutput, "ocid1.policy.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "Test policy 1 description")
	assert.Contains(t, jsonOutput, "TestPolicy2")
	assert.Contains(t, jsonOutput, "ocid1.policy.oc1.phx.test2")
	assert.Contains(t, jsonOutput, "Test policy 2 description")
	assert.Contains(t, jsonOutput, "\"pagination\":")
	assert.Contains(t, jsonOutput, "\"CurrentPage\"")
	assert.Contains(t, jsonOutput, "2")
	assert.Contains(t, jsonOutput, "\"TotalCount\"")
	assert.Contains(t, jsonOutput, "10")
	assert.Contains(t, jsonOutput, "\"Limit\"")
	assert.Contains(t, jsonOutput, "5")
	assert.Contains(t, jsonOutput, "\"NextPageToken\"")
	assert.Contains(t, jsonOutput, "\"next-page-token\"")
}
