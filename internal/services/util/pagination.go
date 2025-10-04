package util

import "fmt"

// PaginateSlice returns a page of items from the full slice, along with the total count and next page token.
// If pageNum <= 0, it is treated as 1. If the start index exceeds the total count, an empty page is returned.
// The next page token is a string page number (e.g., "2") or empty when there is no next page.
func PaginateSlice[T any](all []T, limit, pageNum int) ([]T, int, string) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if limit <= 0 {
		limit = 1
	}
	total := len(all)
	start := (pageNum - 1) * limit
	end := start + limit
	if start >= total {
		return []T{}, total, ""
	}
	if end > total {
		end = total
	}
	paged := all[start:end]
	next := ""
	if end < total {
		next = fmt.Sprintf("%d", pageNum+1)
	}
	return paged, total, next
}
