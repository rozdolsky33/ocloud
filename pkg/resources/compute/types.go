package compute

import (
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

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
