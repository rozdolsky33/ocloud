package subnet

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestPrintSubnetTable tests the PrintSubnetTable function
func TestPrintSubnetTable(t *testing.T) {
	// Create test subnets
	subnets := []Subnet{
		{
			DisplayName:            "TestSubnet1",
			OCID:                   "ocid1.subnet.oc1.phx.test1",
			CIDRBlock:              "10.0.0.0/24",
			ProhibitPublicIPOnVnic: false,
			DNSLabel:               "subnet1",
			SubnetDomainName:       "subnet1.vcn1.oraclevcn.com",
		},
		{
			DisplayName:            "TestSubnet2",
			OCID:                   "ocid1.subnet.oc1.phx.test2",
			CIDRBlock:              "10.0.1.0/24",
			ProhibitPublicIPOnVnic: true,
			DNSLabel:               "subnet2",
			SubnetDomainName:       "subnet2.vcn1.oraclevcn.com",
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
		Stdout: &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintSubnetTable(subnets, appCtx, nil, false, "")
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "Yes")
	assert.Contains(t, output, "subnet1")
	assert.Contains(t, output, "subnet1.vcn1.oraclevcn.com")
	assert.Contains(t, output, "TestSubnet2")
	assert.Contains(t, output, "No")
	assert.Contains(t, output, "10.0.1.0/24")
	assert.Contains(t, output, "subnet2")
	assert.Contains(t, output, "subnet2.vcn1.oraclevcn.com")

	// Test sorting by name
	buf.Reset()
	err = PrintSubnetTable(subnets, appCtx, nil, false, "name")
	assert.NoError(t, err)
	output = buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "TestSubnet2")

	// Test sorting by CIDR
	buf.Reset()
	err = PrintSubnetTable(subnets, appCtx, nil, false, "cidr")
	assert.NoError(t, err)
	output = buf.String()
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "10.0.1.0/24")

	// Test with JSON output (useJSON = true)
	buf.Reset()
	err = PrintSubnetTable(subnets, appCtx, nil, true, "")
	assert.NoError(t, err)
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestSubnet1")
	assert.Contains(t, jsonOutput, "ocid1.subnet.oc1.phx.test1")
	assert.Contains(t, jsonOutput, "10.0.0.0/24")
}

// TestPrintSubnetInfo tests the PrintSubnetInfo function
func TestPrintSubnetInfo(t *testing.T) {
	// Create test subnets
	subnets := []Subnet{
		{
			DisplayName:            "TestSubnet1",
			OCID:                   "ocid1.subnet.oc1.phx.test1",
			CIDRBlock:              "10.0.0.0/24",
			ProhibitPublicIPOnVnic: false,
			DNSLabel:               "subnet1",
			SubnetDomainName:       "subnet1.vcn1.oraclevcn.com",
		},
	}

	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create an application context with the buffer as stdout
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
		Stdout: &buf,
	}

	// Test with table output (useJSON = false)
	err := PrintSubnetInfo(subnets, appCtx, false)
	assert.NoError(t, err)

	// Verify that the output contains the expected information
	output := buf.String()
	assert.Contains(t, output, "TestSubnet1")
	assert.Contains(t, output, "Yes")
	assert.Contains(t, output, "10.0.0.0/24")
	assert.Contains(t, output, "subnet1")
	assert.Contains(t, output, "subnet1.vcn1.oraclevcn.com")

	// Test with JSON output (useJSON = true)
	buf.Reset()
	err = PrintSubnetInfo(subnets, appCtx, true)
	assert.NoError(t, err)

	// Verify that the output is valid JSON and contains the expected information
	jsonOutput := buf.String()
	assert.Contains(t, jsonOutput, "TestSubnet1")
	assert.Contains(t, jsonOutput, "ocid1.subnet.oc1.phx.test1")
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
		Logger: logger.NewTestLogger(),
		Stdout: &buf,
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
	assert.Contains(t, jsonOutput, "{\"items\": []}")
}
