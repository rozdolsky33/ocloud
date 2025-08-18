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

// Service represents a service that interacts with OCI Bastion and Virtual Network APIs to manage resources.
type Service struct {
	bastionClient bastion.BastionClient
	networkClient core.VirtualNetworkClient
	computeClient core.ComputeClient
	logger        logr.Logger
	compartmentID string
	vcnCache      map[string]*core.Vcn
	subnetCache   map[string]*core.Subnet
}

// Config represents the configuration for the bastion service.
type Config struct {
	SshPublicKeyFile  string `yaml:"ssh_pub_key_file"`
	SshPrivateKeyFile string `yaml:"ssh_private_key_file"`
	SessionTimeout    *int   `yaml:"session_timeout"`
}
