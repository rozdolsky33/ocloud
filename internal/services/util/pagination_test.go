package util

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestLogPaginationInfo tests the LogPaginationInfo function
func TestLogPaginationInfo(t *testing.T) {
	// Create a test logger
	testLogger := logger.NewTestLogger()

	// Create a test application context
	appCtx := &app.ApplicationContext{
		Logger: testLogger,
	}

	// Test with nil pagination
	LogPaginationInfo(nil, appCtx)
	// No assertions needed, just making sure it doesn't panic

	// Test with pagination but no next page
	pagination := &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    10,
		Limit:         10,
		NextPageToken: "",
	}
	LogPaginationInfo(pagination, appCtx)
	// No assertions needed, just making sure it doesn't panic

	// Test with pagination and next page
	pagination = &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    20,
		Limit:         10,
		NextPageToken: "next-page-token",
	}
	LogPaginationInfo(pagination, appCtx)
	// No assertions needed, just making sure it doesn't panic

	// Test with pagination and previous page
	pagination = &PaginationInfo{
		CurrentPage:   2,
		TotalCount:    20,
		Limit:         10,
		NextPageToken: "",
	}
	LogPaginationInfo(pagination, appCtx)
	// No assertions needed, just making sure it doesn't panic

	// Test with pagination, previous page, and next page
	pagination = &PaginationInfo{
		CurrentPage:   2,
		TotalCount:    30,
		Limit:         10,
		NextPageToken: "next-page-token",
	}
	LogPaginationInfo(pagination, appCtx)
	// No assertions needed, just making sure it doesn't panic
}

// TestAdjustPaginationInfo tests the AdjustPaginationInfo function
func TestAdjustPaginationInfo(t *testing.T) {
	// Test when total records displayed are less than total count
	pagination := &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    100,
		Limit:         10,
		NextPageToken: "next-page-token",
	}
	AdjustPaginationInfo(pagination)
	assert.Equal(t, 10, pagination.TotalCount, "TotalCount should be adjusted to 10")

	// Test when total records displayed are equal to the total count
	pagination = &PaginationInfo{
		CurrentPage:   10,
		TotalCount:    100,
		Limit:         10,
		NextPageToken: "",
	}
	AdjustPaginationInfo(pagination)
	assert.Equal(t, 100, pagination.TotalCount, "TotalCount should remain 100")

	// Test when total records displayed are greater than total count
	pagination = &PaginationInfo{
		CurrentPage:   11,
		TotalCount:    100,
		Limit:         10,
		NextPageToken: "",
	}
	AdjustPaginationInfo(pagination)
	assert.Equal(t, 100, pagination.TotalCount, "TotalCount should remain 100")
}

// TestValidateAndReportEmpty tests the ValidateAndReportEmpty function
func TestValidateAndReportEmpty(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Test with non-empty items
	items := []string{"item1", "item2"}
	result := ValidateAndReportEmpty(items, nil, &buf)
	assert.False(t, result, "Should return false for non-empty items")
	assert.Equal(t, "", buf.String(), "Should not write anything to output for non-empty items")

	// Reset the buffer
	buf.Reset()

	// Test with empty items and nil pagination
	items = []string{}
	result = ValidateAndReportEmpty(items, nil, &buf)
	assert.True(t, result, "Should return true for empty items")
	assert.Contains(t, buf.String(), "No Items found.", "Should write 'No Items found.' to output")

	// Reset the buffer
	buf.Reset()

	// Test with empty items and pagination with total count = 0
	pagination := &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    0,
		Limit:         10,
		NextPageToken: "",
	}
	result = ValidateAndReportEmpty(items, pagination, &buf)
	assert.True(t, result, "Should return true for empty items")
	assert.Contains(t, buf.String(), "No Items found.", "Should write 'No Items found.' to output")
	assert.NotContains(t, buf.String(), "Page 1 is empty.", "Should not write pagination info when total count is 0")

	// Reset the buffer
	buf.Reset()

	// Test with empty items and pagination with total count > 0
	pagination = &PaginationInfo{
		CurrentPage:   1,
		TotalCount:    10,
		Limit:         10,
		NextPageToken: "",
	}
	result = ValidateAndReportEmpty(items, pagination, &buf)
	assert.True(t, result, "Should return true for empty items")
	assert.Contains(t, buf.String(), "No Items found.", "Should write 'No Items found.' to output")
	assert.Contains(t, buf.String(), "Page 1 is empty. Total records: 10", "Should write pagination info")
	assert.NotContains(t, buf.String(), "Try a lower page number", "Should not suggest a lower page number for page 1")

	// Reset the buffer
	buf.Reset()

	// Test with empty items and pagination with total count > 0 and current page > 1
	pagination = &PaginationInfo{
		CurrentPage:   2,
		TotalCount:    10,
		Limit:         10,
		NextPageToken: "",
	}
	result = ValidateAndReportEmpty(items, pagination, &buf)
	assert.True(t, result, "Should return true for empty items")
	assert.Contains(t, buf.String(), "No Items found.", "Should write 'No Items found.' to output")
	assert.Contains(t, buf.String(), "Page 2 is empty. Total records: 10", "Should write pagination info")
	assert.Contains(t, buf.String(), "Try a lower page number (e.g., --page 1)", "Should suggest a lower page number")
}

// TestPaginateSlice_Basic verifies normal pagination and next token behavior.
func TestPaginateSlice_Basic(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7}
	paged, total, next := PaginateSlice(items, 3, 1)
	assert.Equal(t, 7, total)
	assert.Equal(t, []int{1, 2, 3}, paged)
	assert.Equal(t, "2", next)

	paged2, total2, next2 := PaginateSlice(items, 3, 2)
	assert.Equal(t, 7, total2)
	assert.Equal(t, []int{4, 5, 6}, paged2)
	assert.Equal(t, "3", next2)

	paged3, total3, next3 := PaginateSlice(items, 3, 3)
	assert.Equal(t, 7, total3)
	assert.Equal(t, []int{7}, paged3)
	assert.Equal(t, "", next3)
}

// TestPaginateSlice_PageBeyondRange returns an empty page and no next token.
func TestPaginateSlice_PageBeyondRange(t *testing.T) {
	items := []string{"a", "b", "c"}
	paged, total, next := PaginateSlice(items, 2, 5)
	assert.Equal(t, 3, total)
	assert.Empty(t, paged)
	assert.Equal(t, "", next)
}

// TestPaginateSlice_ZeroOrNegativeDefaults ensures page<=0 and limit<=0 default to 1.
func TestPaginateSlice_ZeroOrNegativeDefaults(t *testing.T) {
	items := []int{10, 20}
	paged, total, next := PaginateSlice(items, 0, 0)
	assert.Equal(t, 2, total)
	assert.Equal(t, []int{10}, paged) // limit defaults to 1
	assert.Equal(t, "2", next)
}

// TestPaginateSlice_EmptyInput returns empty with a total 0.
func TestPaginateSlice_EmptyInput(t *testing.T) {
	var items []struct{}
	paged, total, next := PaginateSlice(items, 10, 1)
	assert.Equal(t, 0, total)
	assert.Empty(t, paged)
	assert.Equal(t, "", next)
}

// TestPaginateSlice_ExactMultiple tests when total is exact multiple of limit.
func TestPaginateSlice_ExactMultiple(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6}
	paged1, total1, next1 := PaginateSlice(items, 3, 1)
	assert.Equal(t, 6, total1)
	assert.Equal(t, []int{1, 2, 3}, paged1)
	assert.Equal(t, "2", next1)

	paged2, total2, next2 := PaginateSlice(items, 3, 2)
	assert.Equal(t, 6, total2)
	assert.Equal(t, []int{4, 5, 6}, paged2)
	assert.Equal(t, "", next2) // exactly ends at total, so no next
}
