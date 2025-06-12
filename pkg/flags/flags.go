package flags

import (
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// FlagNames defines the string constants for flag names
const (
	FlagNameLogLevel          = "log-level"
	FlagNameList              = "list"
	FlagNameFind              = "find"
	FlagNameTenancyID         = "tenancy-id"
	FlagNameTenancyName       = "tenancy-name"
	FlagNameCompartment       = "compartment"
	FlagNameCreate            = "create"
	FlagNameType              = "type"
	FlagNameBastionID         = "bastion-id"
	FlagNameTargetIP          = "target-ip"
	FlagNameTTL               = "ttl"
	FlagNamePrivateKey        = "private-key"
	FlagNamePublicKey         = "public-key"
	FlagNameInstanceID        = "instance-id"
	FlagNameUser              = "user"
	FlagNameOKEName           = "oke-name"
	FlagNameLocalFwPort       = "local-fw-port"
	FlagNameHostFwPort        = "host-fw-port"
	FlagNameImageDetails      = "image-details"
	FlagNameFindByName        = "find-by-name"
	FlagNameFindByStatement   = "find-by-statements"
	FlagNameIncludeStatements = "include-statements"
)

// FlagShorthands defines single-character aliases for flags
const (
	FlagShortList              = "l"
	FlagShortFind              = "f"
	FlagShortTenancyID         = "t"
	FlagShortCompartment       = "c"
	FlagShortCreate            = "r"
	FlagShortType              = "y"
	FlagShortBastionID         = "b"
	FlagShortTargetIP          = "i"
	FlagShortTTL               = "m"
	FlagShortPrivateKey        = "a"
	FlagShortPublicKey         = "e"
	FlagShortInstanceID        = "o"
	FlagShortUser              = "u"
	FlagShortOKEName           = "k"
	FlagShortLocalFwPort       = "w"
	FlagShortHostFwPort        = "f"
	FlagShortImageDetails      = "i"
	FlagShortFindByName        = "n"
	FlagShortFindByStatement   = "s"
	FlagShortIncludeStatements = "a"
)

// FlagDescriptions contains help text for flags
const (
	FlagDescLogLevel          = "Set the log verbosity. Supported values are: debug, info, warn, and error."
	FlagDescList              = "List all resources"
	FlagDescFind              = "Find resources by name pattern search"
	FlagDescCreate            = "Create a resource"
	FlagDescType              = "Resource type"
	FlagDescBastionID         = "Bastion OCID"
	FlagDescTargetIP          = "Target IP address"
	FlagDescTTL               = "TTL in minutes"
	FlagDescPrivateKey        = "Private key file path"
	FlagDescPublicKey         = "Public key file path"
	FlagDescInstanceID        = "Instance OCID"
	FlagDescUser              = "User name"
	FlagDescTenancyID         = "OCI tenancy OCID"
	FlagDescTenancyName       = "Tenancy name"
	FlagDescCompartment       = "OCI compartment name"
	FlagDescOKEName           = "OKE cluster name"
	FlagDescLocalFwPort       = "Local firewall port"
	FlagDescHostFwPort        = "Host firewall port"
	FlagDescImageDetails      = "Show image details"
	FlagDescFindByName        = "Find resources by name pattern search"
	FlagDescFindByStatement   = "Find resources by statement"
	FlagDescIncludeStatements = "Include statements"
)

// Flag represents a command flag configuration
type Flag struct {
	Name        string
	Shorthand   string
	Default     interface{}
	Description string
}

// AddBoolFlag adds a boolean flag to the command
func (f Flag) AddBoolFlag(cmd *cobra.Command) *bool {
	var value bool
	if defaultVal, ok := f.Default.(bool); ok {
		value = defaultVal
	}
	cmd.Flags().BoolVarP(&value, f.Name, f.Shorthand, value, f.Description)
	return &value
}

// AddStringFlag adds a string flag to the command
func (f Flag) AddStringFlag(cmd *cobra.Command) *string {
	var value string
	if defaultVal, ok := f.Default.(string); ok {
		value = defaultVal
	}
	cmd.Flags().StringVarP(&value, f.Name, f.Shorthand, value, f.Description)
	return &value
}

// AddIntFlag adds an integer flag to the command
func (f Flag) AddIntFlag(cmd *cobra.Command) *int {
	var value int
	if defaultVal, ok := f.Default.(int); ok {
		value = defaultVal
	}
	cmd.Flags().IntVarP(&value, f.Name, f.Shorthand, value, f.Description)
	return &value
}

// Define flag configurations for reuse across commands

// Instance flags
var (
	ListFlag = Flag{
		Name:        FlagNameList,
		Shorthand:   FlagShortList,
		Default:     false,
		Description: FlagDescList,
	}

	FindFlag = Flag{
		Name:        FlagNameFind,
		Shorthand:   FlagShortFind,
		Default:     "",
		Description: FlagDescFind,
	}

	ImageDetailsFlag = Flag{
		Name:        FlagNameImageDetails,
		Shorthand:   FlagShortImageDetails,
		Default:     false,
		Description: FlagDescImageDetails,
	}
)

// InitGlobalFlags initializes global CLI flags and binds them to environment variables for configuration.
func InitGlobalFlags(root *cobra.Command) {
	root.PersistentFlags().StringVarP(&logger.LogLevel, FlagNameLogLevel, "", "info", logger.LogLevelMsg)
	root.PersistentFlags().BoolVar(&logger.ColoredOutput, "color", false, logger.ColoredOutputMsg)
	root.PersistentFlags().StringP(FlagNameTenancyID, FlagShortTenancyID, "", FlagDescTenancyID)
	root.PersistentFlags().StringP(FlagNameTenancyName, "", "", FlagDescTenancyName)
	root.PersistentFlags().StringP(FlagNameCompartment, FlagShortCompartment, "", FlagDescCompartment)

	_ = viper.BindPFlag(FlagNameTenancyID, root.PersistentFlags().Lookup(FlagNameTenancyID))
	_ = viper.BindPFlag(FlagNameTenancyName, root.PersistentFlags().Lookup(FlagNameTenancyName))
	_ = viper.BindPFlag(FlagNameCompartment, root.PersistentFlags().Lookup(FlagNameCompartment))

	// allow ENV overrides, e.g., OCI_CLI_TENANCY, OCI_TENANCY_NAME, OCI_COMPARTMENT
	viper.SetEnvPrefix("OCI")
	viper.AutomaticEnv()
}
