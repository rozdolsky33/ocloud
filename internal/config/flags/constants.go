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
	FlagNameHelp              = "help"
	FlagNameColor             = "color"
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
	FlagShortHelp              = "h"
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
	FlagDescHelp              = "help for ocloud (shorthand: -h)"
	FlagDescColor             = "Enable colored output"
)

// Flag values and defaults
const (
	FlagValueTrue     = "true"
	FlagValueInfo     = "info"
	FlagValueHelpMode = "help-mode"
)

// Flag prefixes and special strings
const (
	FlagPrefixShortHelp  = "-h"
	FlagPrefixLongHelp   = "--help"
	FlagPrefixLogLevel   = "--log-level"
	FlagPrefixLogLevelEq = "--log-level="
	FlagPrefixColor      = "--color"
	CobraAnnotationKey   = "cobra_annotation_flag_set_by_cobra"
)

// Environment variables
const (
	EnvOCITenancy     = "OCI_CLI_TENANCY"
	EnvOCITenancyName = "OCI_TENANCY_NAME"
	EnvOCICompartment = "OCI_COMPARTMENT"
	EnvOCIRegion      = "OCI_REGION"
)
