package vcn

import (
	"context"
)

// Subnet represents a virtual network subnet.
type Subnet struct {
	OCID                    string
	DisplayName             string
	CIDRBlock               string
	VcnOCID                 string
	RouteTableOCID          string
	SecurityListOCIDs       []string
	DhcpOptionsOCID         string
	ProhibitPublicIPOnVnic  bool
	ProhibitInternetIngress bool
	DNSLabel                string
	SubnetDomainName        string
}

// SubnetRepository defines the port for interacting with subnet storage.
type SubnetRepository interface {
	GetSubnet(ctx context.Context, ocid string) (*Subnet, error)
	ListSubnets(ctx context.Context, compartmentID string) ([]Subnet, error)
}
