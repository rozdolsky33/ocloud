package instance

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// Service encapsulates OCI compute/network clients and config.
// It provides methods to list and find instances without printing directly.
type Service struct {
	compute           core.ComputeClient
	network           core.VirtualNetworkClient
	logger            logr.Logger
	compartmentID     string
	enableConcurrency bool
	// Caches to reduce API calls
	subnetCache     map[string]*core.Subnet
	vcnCache        map[string]*core.Vcn
	routeTableCache map[string]*core.RouteTable
	// Cache for pagination tokens
	pageTokenCache map[string]map[int]string // compartmentID -> page number -> page token
}

// Instance represents a VM instance in the cloud.
type Instance struct {
	Name              string
	ID                string
	IP                string
	ImageID           string
	SubnetID          string
	Shape             string
	State             core.InstanceLifecycleStateEnum
	CreatedAt         common.SDKTime
	Placement         Placement
	Resources         Resources
	ImageName         string
	ImageOS           string
	InstanceTags      ResourceTags
	Hostname          string
	SubnetName        string
	VcnID             string
	VcnName           string
	PrivateDNSEnabled bool
	RouteTableID      string
	RouteTableName    string
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
type ResourceTags struct {
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

// VnicInfo Define a struct to hold VNIC information
type VnicInfo struct {
	InstanceID string
	Ip         string
	SubnetID   string
	Err        error
}

// JSONResponse Create a response structure that includes instances and pagination info
type JSONResponse struct {
	Instances  []Instance      `json:"instances"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
}

// PaginationInfo holds information about the current page and total results
type PaginationInfo struct {
	CurrentPage   int
	TotalCount    int
	Limit         int
	NextPageToken string
}

// IndexableInstance represents a simplified structure of an Instance for indexing and searchable purposes.
type IndexableInstance struct {
	ID                   string
	Name                 string
	ImageName            string
	ImageOperatingSystem string
	IP                   string
	SubnetID             string
	Tags                 string
	TagValues            string // Separate field for tag values only, to make them directly searchable
}
