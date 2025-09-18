package domain

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
	TimeCreated    time.Time
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
}

type VCNRepository interface {
	GetVcn(ctx context.Context, ocid string) (*VCN, error)
	ListVcns(ctx context.Context, compartmentID string) ([]*VCN, error)
}
