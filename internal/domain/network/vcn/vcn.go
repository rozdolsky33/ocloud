package vcn

import (
	"context"
	"time"
)

// VCN represents a Virtual Cloud Network in the domain layer.
type VCN struct {
	OCID           string
	DisplayName    string
	LifecycleState string
	CompartmentID  string
	CidrBlocks     []string
	Ipv6Enabled    bool
	DnsLabel       string
	DomainName     string
	DhcpOptionsID  string
	DhcpOptions    DhcpOptions
	TimeCreated    time.Time
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
	Gateways       []Gateway
	Subnets        []Subnet
	RouteTables    []RouteTable
	SecurityLists  []SecurityList
	NSGs           []NSG
}

type VCNRepository interface {
	GetVcn(ctx context.Context, ocid string) (*VCN, error)
	ListVcns(ctx context.Context, compartmentID string) ([]VCN, error)
	ListEnrichedVcns(ctx context.Context, compartmentID string) ([]VCN, error)
}
