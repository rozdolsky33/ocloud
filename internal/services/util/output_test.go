package util

import (
	"bytes"
	"encoding/json"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMarshalDataToJSON tests the MarshalDataToJSON function
func TestMarshalDataToJSON(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a printer that writes to the buffer
	p := printer.New(&buf)

	// Test with empty items and nil pagination
	var items []string
	err := MarshalDataToJSON(p, items, nil)
	assert.NoError(t, err, "MarshalDataToJSON should not return an error")

	// Verify that the output is valid JSON
	var response JSONResponse[string]
	err = json.Unmarshal(buf.Bytes(), &response)
	assert.NoError(t, err, "Output should be valid JSON")
	assert.Empty(t, response.Items, "Items should be empty")
	assert.Nil(t, response.Pagination, "Pagination should be nil")

	// Reset the buffer
	buf.Reset()

	// Test with non-empty items and nil pagination
	items = []string{"item1", "item2"}
	err = MarshalDataToJSON(p, items, nil)
	assert.NoError(t, err, "MarshalDataToJSON should not return an error")

	// Verify that the output is valid JSON
	err = json.Unmarshal(buf.Bytes(), &response)
	assert.NoError(t, err, "Output should be valid JSON")
	assert.Equal(t, items, response.Items, "Items should match")
	assert.Nil(t, response.Pagination, "Pagination should be nil")

	// Reset the buffer
	buf.Reset()

	// Test with non-empty items and pagination
	pagination := &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    10,
		Limit:         5,
		NextPageToken: "next-page-token",
	}
	err = MarshalDataToJSON(p, items, pagination)
	assert.NoError(t, err, "MarshalDataToJSON should not return an error")

	// Verify that the output is valid JSON
	err = json.Unmarshal(buf.Bytes(), &response)
	assert.NoError(t, err, "Output should be valid JSON")
	assert.Equal(t, items, response.Items, "Items should match")
	assert.Equal(t, pagination, response.Pagination, "Pagination should match")
}

// TestFormatColoredTitle tests the FormatColoredTitle function
func TestFormatColoredTitle(t *testing.T) {
	// Create a test application context
	appCtx := &app.ApplicationContext{
		TenancyName:     "TestTenancy",
		CompartmentName: "TestCompartment",
	}

	// Test with a simple name
	title := FormatColoredTitle(appCtx, "TestName")

	// We can't easily test the exact output since it includes ANSI color codes,
	// but we can check that it contains the expected components
	assert.Contains(t, title, "TestTenancy")
	assert.Contains(t, title, "TestCompartment")
	assert.Contains(t, title, "TestName")

	// Test with an empty name
	title = FormatColoredTitle(appCtx, "")
	assert.Contains(t, title, "TestTenancy")
	assert.Contains(t, title, "TestCompartment")

	// Test with empty tenancy and compartment
	emptyAppCtx := &app.ApplicationContext{}
	title = FormatColoredTitle(emptyAppCtx, "TestName")
	assert.Contains(t, title, "TestName")
}
