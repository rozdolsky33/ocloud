package domain

import (
	"context"
	"time"
)

// Instance represents a compute instance in the cloud.
type Instance struct {
	OCID               string
	DisplayName        string
	State              string
	Shape              string
	ImageID            string
	TimeCreated        time.Time
	Region             string
	AvailabilityDomain string
	FaultDomain        string
	VCPUs              int
	MemoryGB           float32
	FreeformTags       map[string]string
	DefinedTags        map[string]map[string]interface{}
	// Enriched fields
	PrimaryIP         string
	SubnetID          string
	SubnetName        string
	VcnID             string
	VcnName           string
	ImageName         string
	ImageOS           string
	Hostname          string
	PrivateDNSEnabled bool
	RouteTableName    string
	RouteTableID      string
}

// InstanceRepository defines the port for interacting with instance storage.
// Implementations will handle the complexity of fetching and enriching instance data.
type InstanceRepository interface {
	ListInstances(ctx context.Context, compartmentID string) ([]Instance, error)
}
