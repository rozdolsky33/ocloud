package util

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"github.com/stretchr/testify/assert"
)

// TestMarshalDataToJSON tests the MarshalDataToJSONResponse function
func TestMarshalDataToJSON(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a printer that writes to the buffer
	p := printer.New(&buf)

	// Test with empty items and nil pagination
	var items []string
	err := MarshalDataToJSONResponse(p, items, nil)
	assert.NoError(t, err, "MarshalDataToJSONResponse should not return an error")

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
	err = MarshalDataToJSONResponse(p, items, nil)
	assert.NoError(t, err, "MarshalDataToJSONResponse should not return an error")

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
	err = MarshalDataToJSONResponse(p, items, pagination)
	assert.NoError(t, err, "MarshalDataToJSONResponse should not return an error")

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

// TestSplitTextByMaxWidth tests the SplitTextByMaxWidth function
func TestSplitTextByMaxWidth(t *testing.T) {
	// Test with an empty string
	result := SplitTextByMaxWidth("")
	assert.Equal(t, []string{""}, result, "Empty string should return a slice with an empty string")

	// Test with a single word
	result = SplitTextByMaxWidth("SingleWord")
	assert.Equal(t, []string{"SingleWord"}, result, "Single word should return a slice with that word")

	// Test with a short string that fits in one line
	result = SplitTextByMaxWidth("This is a short string")
	assert.Equal(t, []string{"This is a short string"}, result, "Short string should return a slice with that string")

	// Test with a long string that needs to be split
	longString := "This is a very long string that should be split into multiple lines because it exceeds the maximum width"
	result = SplitTextByMaxWidth(longString)
	assert.Greater(t, len(result), 1, "Long string should be split into multiple lines")

	// Test with a string that has exactly the max width
	// The max width in the function is 30 characters
	exactWidthString := "This string has exactly thirty"
	result = SplitTextByMaxWidth(exactWidthString)
	assert.Equal(t, 1, len(result), "String with exact max width should not be split")

	// Test with a string that has multiple spaces
	multiSpaceString := "This   string   has   multiple   spaces"
	result = SplitTextByMaxWidth(multiSpaceString)
	// The function uses strings.Fields which normalizes spaces, but also splits by max width
	assert.Equal(t, []string{"This string has multiple", "spaces"}, result, "String with multiple spaces should be normalized and split if needed")
}
