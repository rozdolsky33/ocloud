package bastion

import (
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/bastion"
	"github.com/oracle/oci-go-sdk/v65/core"
	conf "github.com/rozdolsky33/ocloud/internal/config"
	cflags "github.com/rozdolsky33/ocloud/internal/config/flags"
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

type Config struct {
	SshPublicKeyFile  string `yaml:"ssh_pub_key_file"`
	SshPrivateKeyFile string `yaml:"ssh_private_key_file"`
	SessionTimeout    *int   `yaml:"session_timeout"`
}

// sshConfig populates default values for SSH key locations based on the active OCI profile.
// It dynamically constructs the path to the OCI session directory using standard constants.
// The generated default points to the per-profile OCI session directory under:
//
//	<home>/.oci/sessions/<PROFILE>
//
// By default, it selects the standard OCI key filenames inside that directory:
//   - oci_api_key_public.pem (public)
//   - oci_api_key.pem (private)
func (config *Config) sshConfig() {
	homeDir, _ := conf.GetUserHomeDir()
	profile := conf.GetOCIProfile()
	// Build session directory: ~/.oci/sessions/<profile>
	sessionDir := filepath.Join(homeDir, cflags.OCIConfigDirName, cflags.OCISessionsDirName, profile)

	if config.SshPublicKeyFile == "" {
		config.SshPublicKeyFile = filepath.Join(sessionDir, "oci_api_key_public.pem")
	}
	if config.SshPrivateKeyFile == "" {
		config.SshPrivateKeyFile = filepath.Join(sessionDir, "oci_api_key.pem")
	}
}
