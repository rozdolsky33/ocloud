package bastion

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/bastion"
)

// Bastion represents a bastion host in the system
type Bastion struct {
	ID             string
	Name           string
	BastionType    string
	TargetVcnId    string
	TargetSubnetId string
	LifecycleState bastion.BastionLifecycleStateEnum
}

type Service struct {
	bastionClient bastion.BastionClient
	logger        logr.Logger
	compartmentID string
}
