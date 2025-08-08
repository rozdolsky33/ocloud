package bastion

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// Bastion represents a bastion host in the system
type Bastion struct {
	ID               string
	Name             string
	BastionType      string
	TargetVcnId      string
	TargetVcnName    string
	TargetSubnetId   string
	TargetSubnetName string
	LifecycleState   bastion.BastionLifecycleStateEnum
}

type Service struct {
	bastionClient bastion.BastionClient
	networkClient core.VirtualNetworkClient
	logger        logr.Logger
	compartmentID string
	vcnCache      map[string]*core.Vcn
	subnetCache   map[string]*core.Subnet
}
