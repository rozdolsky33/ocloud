package configs

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

// FlagShorthand defines single-character aliases for flags
const (
	FlagShortLogLevel          = "l"
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
	FlagSortFindByName         = "n"
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
	FlagDescImageDetails      = "Image details"
	FlagDescFindByName        = "Find resources by name pattern search"
	FlagDescFindByStatement   = "Find resources by statement"
	FlagDescIncludeStatements = "Include statements"
)

// EnvironmentVars defines environment variable names
const (
	EnvOCIRegion      = "OCI_CLI_REGION"
	EnvOCITenancy     = "OCI_CLI_TENANCY"
	EnvOCITenancyName = "OCI_TENANCY_NAME"
	EnvOCICompartment = "OCI_COMPARTMENT"
)

// ErrorMessages defines error text
const (
	invalidFlagMessage = "flag arguments"
)

// FlagConfig represents a command flag configuration
type FlagConfig struct {
	Name        string
	Shorthand   string
	Default     interface{}
	Description string
}
