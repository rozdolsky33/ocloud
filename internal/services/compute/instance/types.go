package instance

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/services/util"
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
	InstanceTags      util.ResourceTags
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

// VnicInfo Define a struct to hold VNIC information
type VnicInfo struct {
	InstanceID string
	Ip         string
	SubnetID   string
	Err        error
}

// IndexableInstance represents a simplified structure of an Instance for indexing and searchable purposes.
type IndexableInstance struct {
	ID                   string
	Name                 string
	ImageName            string
	ImageOperatingSystem string
	Tags                 string
	TagValues            string // Separate field for tag values only, to make them directly searchable
}
