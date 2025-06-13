// Package flags defines flag types and domain-specific flag collections for the CLI.
package flags

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

// Environment variables
const (
	EnvOCITenancy     = "OCI_CLI_TENANCY"
	EnvOCITenancyName = "OCI_TENANCY_NAME"
	EnvOCICompartment = "OCI_COMPARTMENT"
	EnvOCIRegion      = "OCI_REGION"
)
