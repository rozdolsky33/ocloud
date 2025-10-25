package bastion

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

// Bastion is an alias to the domain.Bastion type.
// This maintains backward compatibility while using the domain model.
type Bastion = domain.Bastion

// Service represents a service that manages bastion operations.
// It now depends on the BastionRepository interface instead of concrete OCI clients.
type Service struct {
	bastionRepo   domain.BastionRepository
	bastionClient bastion.BastionClient // Temporarily kept for session management
	networkClient core.VirtualNetworkClient
	computeClient core.ComputeClient
	logger        logr.Logger
	compartmentID string
}

// Config represents the configuration for the bastion service.
type Config struct {
	SshPublicKeyFile  string `yaml:"ssh_pub_key_file"`
	SshPrivateKeyFile string `yaml:"ssh_private_key_file"`
	SessionTimeout    *int   `yaml:"session_timeout"`
}
