package image

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/services/util"
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
	ImageTags       util.ResourceTags
}

// IndexableImage represents an image model optimized for indexing and searching in the application.
type IndexableImage struct {
	ID              string
	Name            string
	ImageOSVersion  string
	OperatingSystem string
	LunchMode       string
	Tags            string
	TagValues       string // Separate field for tag values only, to make them directly searchable
}
