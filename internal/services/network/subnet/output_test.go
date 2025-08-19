package subnet

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/stretchr/testify/assert"
)

// TestPrintSubnetTable tests the PrintSubnetTable function
func TestPrintSubnetTable(t *testing.T) {
	// Create test subnets
	subnets := []Subnet{
		{
			Name:                   "TestSubnet1",
			ID:                     "ocid1.subnet.oc1.phx.test1",
			CIDR:                   "10.0.0.0/24",
			ProhibitPublicIPOnVnic: true,
			DNSLabel:               "test1",
			SubnetDomainName:       "test1.vcn.oraclevcn.com",
		},
		{
			Name:                   "TestSubnet2",
			ID:                     "ocid1.subnet.oc1.phx.test2",
			CIDR:                   "10.0.1.0/24",
			ProhibitPublicIPOnVnic: false,
			DNSLabel:               "test2",
			SubnetDomainName:       "test2.vcn.oraclevcn.com",
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
	err := PrintSubnetTable(subnets, appCtx, nil, false, "")
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "No")
	assert.Contains(t, output, "test1")
	assert.Contains(t, output, "test1.vcn...")
	assert.Contains(t, output, "TestSubnet2")
	assert.Contains(t, output, "10.0.1.0/24")
	assert.Contains(t, output, "Yes")
	assert.Contains(t, output, "test2")
	assert.Contains(t, output, "test2.vcn...")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintSubnetTable(subnets, appCtx, nil, true, "")
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestSubnet1")
	assert.Contains(t, jsonOutput, "ocid1.subnet.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "10.0.0.0/24")
	assert.Contains(t, jsonOutput, "TestSubnet2")
	assert.Contains(t, jsonOutput, "ocid1.subnet.oc1.phx.test2")
	assert.Contains(t, jsonOutput, "10.0.1.0/24")
}

// TestPrintSubnetTableEmpty tests the PrintSubnetTable function with empty subnets
func TestPrintSubnetTableEmpty(t *testing.T) {
	// Create an empty subnets slice
	subnets := []Subnet{}

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
	err := PrintSubnetTable(subnets, appCtx, nil, false, "")
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintSubnetTable(subnets, appCtx, nil, true, "")
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty object
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "{}")
}

// TestPrintSubnetTableWithPagination tests the PrintSubnetTable function with pagination
func TestPrintSubnetTableWithPagination(t *testing.T) {
	// Create test subnets
	subnets := []Subnet{
		{
			Name:                   "TestSubnet1",
			ID:                     "ocid1.subnet.oc1.phx.test1",
			CIDR:                   "10.0.0.0/24",
			ProhibitPublicIPOnVnic: true,
			DNSLabel:               "test1",
			SubnetDomainName:       "test1.vcn.oraclevcn.com",
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
	err := PrintSubnetTable(subnets, appCtx, pagination, false, "")
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "No")
	assert.Contains(t, output, "test1")
	assert.Contains(t, output, "test1.vcn...")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintSubnetTable(subnets, appCtx, pagination, true, "")
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestSubnet1")
	assert.Contains(t, jsonOutput, "ocid1.subnet.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "10.0.0.0/24")
	assert.Contains(t, jsonOutput, "\"pagination\"")
	assert.Contains(t, jsonOutput, "\"TotalCount\"")
	assert.Contains(t, jsonOutput, "\"Limit\"")
	assert.Contains(t, jsonOutput, "\"CurrentPage\"")
	assert.Contains(t, jsonOutput, "\"NextPageToken\"")
}

// TestPrintSubnetInfo tests the PrintSubnetInfo function
func TestPrintSubnetInfo(t *testing.T) {
	// Create test subnets
	subnets := []Subnet{
		{
			Name:                   "TestSubnet1",
			ID:                     "ocid1.subnet.oc1.phx.test1",
			CIDR:                   "10.0.0.0/24",
			ProhibitPublicIPOnVnic: true,
			DNSLabel:               "test1",
			SubnetDomainName:       "test1.vcn.oraclevcn.com",
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
	err := PrintSubnetInfo(subnets, appCtx, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "No")
	assert.Contains(t, output, "test1")
	assert.Contains(t, output, "test1.vcn.oraclevcn.com")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintSubnetInfo(subnets, appCtx, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestSubnet1")
	assert.Contains(t, jsonOutput, "10.0.0.0/24")
}

// TestPrintSubnetInfoEmpty tests the PrintSubnetInfo function with empty subnets
func TestPrintSubnetInfoEmpty(t *testing.T) {
	// Create an empty subnets slice
	subnets := []Subnet{}

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
	err := PrintSubnetInfo(subnets, appCtx, false)
	assert.NoError(t, err)

	// Verify that the output indicates no items found
	output := buf.String()
	assert.Contains(t, output, "No Items found.")

	// Reset the buffer
	buf.Reset()

	// Test with JSON output (useJSON = true)
	err = PrintSubnetInfo(subnets, appCtx, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and indicates an empty items array
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "\"items\": []")
}
