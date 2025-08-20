package oke

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

// TestPrintOKETable tests the PrintOKETable function
func TestPrintOKETable(t *testing.T) {
	// Create test clusters
	clusters := []Cluster{
		{
			Name:            "TestCluster1",
			ID:              "ocid1.cluster.oc1.phx.test1",
			CreatedAt:       "2023-01-01T00:00:00Z",
			Version:         "v1.24.1",
			State:           "ACTIVE",
			PrivateEndpoint: "10.0.0.1",
			VcnID:           "ocid1.vcn.oc1.phx.test1",
			NodePools: []NodePool{
				{
					Name:      "TestNodePool1",
					ID:        "ocid1.nodepool.oc1.phx.test1",
					Version:   "v1.24.1",
					State:     "ACTIVE",
					NodeShape: "VM.Standard.E3.Flex",
					NodeCount: 2,
					Image:     "ocid1.image.oc1.phx.test1",
					Ocpus:     "1.0",
					MemoryGB:  "16",
				},
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
	err := PrintOKETable(clusters, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestCluster1")
	assert.Contains(t, output, "Cluster")
	assert.Contains(t, output, "v1.24.1")
	assert.Contains(t, output, "10.0.0.1")
	assert.Contains(t, output, "ACTIVE")
	assert.Contains(t, output, "TestNod") // Truncated name
	assert.Contains(t, output, "NodePool")
	assert.Contains(t, output, "VM.Stan") // Truncated shape
	assert.Contains(t, output, "2")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintOKETable(clusters, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestCluster1")
	assert.Contains(t, jsonOutput, "ocid1.cluster.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "v1.24.1")
	assert.Contains(t, jsonOutput, "TestNodePool1")
	assert.Contains(t, jsonOutput, "ocid1.nodepool.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "VM.Standard.E3.Flex")
}

// TestPrintOKETableEmpty tests the PrintOKETable function with empty clusters
func TestPrintOKETableEmpty(t *testing.T) {
	// Create an empty clusters slice
	clusters := []Cluster{}

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
	err := PrintOKETable(clusters, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintOKETable(clusters, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty items array
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "\"items\": []")
}

// TestPrintOKETableWithPagination tests the PrintOKETable function with pagination
func TestPrintOKETableWithPagination(t *testing.T) {
	// Create test clusters
	clusters := []Cluster{
		{
			Name:            "TestCluster1",
			ID:              "ocid1.cluster.oc1.phx.test1",
			CreatedAt:       "2023-01-01T00:00:00Z",
			Version:         "v1.24.1",
			State:           "ACTIVE",
			PrivateEndpoint: "10.0.0.1",
			VcnID:           "ocid1.vcn.oc1.phx.test1",
			NodePools: []NodePool{
				{
					Name:      "TestNodePool1",
					ID:        "ocid1.nodepool.oc1.phx.test1",
					Version:   "v1.24.1",
					State:     "ACTIVE",
					NodeShape: "VM.Standard.E3.Flex",
					NodeCount: 2,
					Image:     "ocid1.image.oc1.phx.test1",
					Ocpus:     "1.0",
					MemoryGB:  "16",
				},
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
	err := PrintOKETable(clusters, appCtx, pagination, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestCluster1")
	assert.Contains(t, output, "Cluster")
	assert.Contains(t, output, "v1.24.1")
	assert.Contains(t, output, "10.0.0.1")
	assert.Contains(t, output, "ACTIVE")
	assert.Contains(t, output, "TestNod") // Truncated name
	assert.Contains(t, output, "NodePool")
	assert.Contains(t, output, "VM.Stan") // Truncated shape
	assert.Contains(t, output, "2")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintOKETable(clusters, appCtx, pagination, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestCluster1")
	assert.Contains(t, jsonOutput, "ocid1.cluster.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "v1.24.1")
	assert.Contains(t, jsonOutput, "TestNodePool1")
	assert.Contains(t, jsonOutput, "ocid1.nodepool.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "VM.Standard.E3.Flex")
	assert.Contains(t, jsonOutput, "\"pagination\"")
	assert.Contains(t, jsonOutput, "\"TotalCount\"")
	assert.Contains(t, jsonOutput, "\"Limit\"")
	assert.Contains(t, jsonOutput, "\"CurrentPage\"")
	assert.Contains(t, jsonOutput, "\"NextPageToken\"")
}

// TestPrintOKEInfo tests the PrintOKEInfo function
func TestPrintOKEInfo(t *testing.T) {
	// Create test clusters
	clusters := []Cluster{
		{
			Name:            "TestCluster1",
			ID:              "ocid1.cluster.oc1.phx.test1",
			CreatedAt:       "2023-01-01T00:00:00Z",
			Version:         "v1.24.1",
			State:           "ACTIVE",
			PrivateEndpoint: "10.0.0.1",
			VcnID:           "ocid1.vcn.oc1.phx.test1",
			NodePools: []NodePool{
				{
					Name:      "TestNodePool1",
					ID:        "ocid1.nodepool.oc1.phx.test1",
					Version:   "v1.24.1",
					State:     "ACTIVE",
					NodeShape: "VM.Standard.E3.Flex",
					NodeCount: 2,
					Image:     "ocid1.image.oc1.phx.test1",
					Ocpus:     "1.0",
					MemoryGB:  "16",
				},
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
	err := PrintOKEInfo(clusters, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestCluster1")
	assert.Contains(t, output, "ocid1.cluster.oc1.phx.test1")
	assert.Contains(t, output, "v1.24.1")
	assert.Contains(t, output, "10.0.0.1")
	assert.Contains(t, output, "ACTIVE")
	assert.Contains(t, output, "TestNod") // Truncated name
	assert.Contains(t, output, "VM.Stan") // Truncated shape
	assert.Contains(t, output, "2")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintOKEInfo(clusters, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestCluster1")
	assert.Contains(t, jsonOutput, "ocid1.cluster.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "v1.24.1")
	assert.Contains(t, jsonOutput, "TestNodePool1")
	assert.Contains(t, jsonOutput, "ocid1.nodepool.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "VM.Standard.E3.Flex")
}

// TestPrintOKEInfoEmpty tests the PrintOKEInfo function with empty clusters
func TestPrintOKEInfoEmpty(t *testing.T) {
	// Create an empty clusters slice
	clusters := []Cluster{}

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
	err := PrintOKEInfo(clusters, appCtx, nil, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintOKEInfo(clusters, appCtx, nil, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty items array
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "\"items\": []")
}
