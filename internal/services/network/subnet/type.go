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

// Subnet represents a logical subdivision of a VCN containing resources with shared CIDR and networking configurations.
// Name is the user-defined subnet name.
// ID is the unique identifier of the subnet.
// CIDR defines the IP address block for the subnet.
// VcnID is the unique identifier of the associated Virtual Cloud Network (VCN).
// RouteTableID is the identifier of the associated route table.
// SecurityListID is a list of IDs for associated security lists.
// DhcpOptionsID is the identifier of the associated DHCP options.
// ProhibitPublicIPOnVnic indicates whether public IP addresses are prohibited in the subnet.
// ProhibitInternetIngress specifies if the subnet disallows incoming internet traffic.
// ProhibitInternetEgress specifies if the subnet disallows outgoing internet traffic.
// DNSLabel is the label assigned for DNS configuration within the VCN.
// SubnetDomainName is the fully qualified domain name of the subnet.
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
