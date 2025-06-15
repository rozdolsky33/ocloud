package compute

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// Service encapsulates OCI compute/network clients and config.
// It provides methods to list and find instances without printing directly.
type Service struct {
	compute            core.ComputeClient
	network            core.VirtualNetworkClient
	logger             logr.Logger
	compartmentID      string
	disableConcurrency bool
}

// Instance represents a VM instance in the cloud.
type Instance struct {
	Name            string
	ID              string
	IP              string
	ImageID         string
	SubnetID        string
	Shape           string
	State           core.InstanceLifecycleStateEnum
	CreatedAt       common.SDKTime
	OperatingSystem string
	Placement       Placement
	Resources       Resources
}

// Placement groups region/availability/fault‐domain info.
type Placement struct {
	Region             string
	AvailabilityDomain string
	FaultDomain        string
}

// Resources group CPU and memory sizing.
type Resources struct {
	VCPUs    int
	MemoryGB float32
}

// Define a struct to hold VNIC information
type VnicInfo struct {
	InstanceID string
	Ip         string
	SubnetID   string
	Err        error
}

// PaginationInfo holds information about the current page and total results
type PaginationInfo struct {
	CurrentPage   int
	TotalCount    int
	Limit         int
	NextPageToken string
}
