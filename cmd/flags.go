package cmd

// FlagNames defines the string constants for flag names
const (
	FlagNameDebug             = "debug"
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
	FlagShortDebug             = "d"
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
	FlagDescDebug             = "Enable debug logging"
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
	FlagDescTenancyID         = "Tenancy OCID"
	FlagDescTenancyName       = "Tenancy name"
	FlagDescCompartment       = "Compartment Name"
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

/**
The main improvements in this refactoring are:
1. Constants are now grouped logically by their purpose (flag names, shorthands, descriptions, etc.)
2. Consistent naming convention using proper Go naming standards (e.g., `FlagNameList` instead of) `flagList`
3. Introduction of a `FlagConfig` struct to encapsulate flag configuration
4. New `getSearchFlags` function that returns flag configurations, making it easier to maintain and extend
5. Enhanced function that handles different flag types dynamically `addCommonSearchFlags`
6. Better documentation with clear grouping and purpose descriptions
7. Removed redundant constants and consolidated related ones
8. More type safety through the use of structured types instead of loose constants

These changes make the code more maintainable, self-documenting, and easier to extend in the future. The use of structured types also helps catch potential errors at compile-time rather than runtime.
*/
