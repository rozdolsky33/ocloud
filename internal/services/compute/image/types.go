package image

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// Service represents a structure that encapsulates the Compute client, logger, and compartment information.
type Service struct {
	compute       core.ComputeClient
	logger        logr.Logger
	compartmentID string
}

// Image represents the metadata and properties of an image resource in the system.
type Image struct {
	Name            string
	ID              string
	CreatedAt       common.SDKTime
	OperatingSystem string
	ImageOSVersion  string
	LunchMode       string
	ImageTags       ResourceTags
}

// ResourceTags represents a collection of user-defined and predefined tags associated with a resource.
// FreeformTags is a simple key-value pair map defined by the user for tagging purposes.
// DefinedTags is a nested map containing namespace and keys with associated values for structured tagging.
type ResourceTags struct {
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

// PaginationInfo holds information about the current page and total results
type PaginationInfo struct {
	CurrentPage   int
	TotalCount    int
	Limit         int
	NextPageToken string
}

// JSONResponse Create a response structure that includes instances and pagination info
type JSONResponse struct {
	Images     []Image         `json:"image"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// IndexableImage represents an image model optimized for indexing and searching in the application.
type IndexableImage struct {
	ID              string
	Name            string
	OperatingSystem string
	ImageOSVersion  string
	Tags            string
	TagValues       string // Separate field for tag values only, to make them directly searchable
}
