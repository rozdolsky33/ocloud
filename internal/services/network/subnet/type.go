package subnet

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// Service encapsulates the network client, logger, and compartment ID for managing virtual network operations.
type Service struct {
	networkClient core.VirtualNetworkClient
	logger        logr.Logger
	compartmentID string
}

// Subnet represents the metadata and properties of a subnet resource in the system.
type Subnet struct {
	Name                    string
	ID                      string
	CIDR                    string
	VcnID                   string
	RouteTableID            string
	SecurityListID          []string
	DhcpOptionsID           string
	ProhibitPublicIPOnVnic  bool
	ProhibitInternetIngress bool
	ProhibitInternetEgress  bool
	DNSLabel                string
	SubnetDomainName        string
}

// IndexableSubnet represents a subnet with attributes formatted for indexing purposes.
// Name is the lowercase representation of the subnet's name.
// CIDR specifies the subnet's IP address block.
type IndexableSubnet struct {
	Name string
	CIDR string
}
