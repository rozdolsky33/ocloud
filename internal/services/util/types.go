package util

// ResourceTags represents a collection of user-defined and predefined tags associated with a resource.
// FreeformTags is a simple key-value pair map defined by the user for tagging purposes.
// DefinedTags is a nested map containing namespace and keys with associated values for structured tagging.
type ResourceTags struct {
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

// PaginationInfo holds information about the current page and total results
type PaginationInfo struct {
	CurrentPage   int    `json:"CurrentPage"`
	TotalCount    int    `json:"TotalCount"`
	Limit         int    `json:"Limit"`
	NextPageToken string `json:"NextPageToken"`
}

// JSONResponse represents a generic JSON structure containing a list of items and optional pagination information.
type JSONResponse[T any] struct {
	Items      []T             `json:"items"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}
