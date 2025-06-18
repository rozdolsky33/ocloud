package images

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

type Service struct {
	compute       core.ComputeClient
	logger        logr.Logger
	compartmentID string
}

type Image struct {
	Name            string
	ID              string
	CreatedAt       common.SDKTime
	OperatingSystem string
	ImageOSVersion  string
	LunchMode       string
	ImageTags       ImageTags
}

type ImageTags struct {
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
	Images     []Image         `json:"images"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}
