package subnet

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/core"
)

type Service struct {
	networkClient core.VirtualNetworkClient
	logger        logr.Logger
	compartmentID string
}

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

type IndexableSubnet struct {
	Name string
	CIDR string
}
