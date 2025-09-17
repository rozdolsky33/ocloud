// Package flags define flag types and domain-specific flag collections for the CLI.
package flags

// FlagNames defines the string constants for flag names
const (
	FlagNameLogLevel       = "log-level"
	FlagNameDebug          = "debug"
	FlagNameTenancyID      = "tenancy-id"
	FlagNameTenancyName    = "tenancy-name"
	FlagNameCompartment    = "compartment"
	FlagNameHelp           = "help"
	FlagNameColor          = "color"
	FlagNameLimit          = "limit"
	FlagNamePage           = "page"
	FlagNameJSON           = "json"
	FlagNameVersion        = "version"
	FlagNameAllInformation = "all"
	FlagNameSort           = "sort"
	FlagNameRealm          = "realm"
	FlagNameFilter         = "filter"
	FlagNameScope          = "scope"
	FlagNameTenancyScope   = "tenancy-scope"
)

// FlagShorthands defines single-character aliases for flags
const (
	FlagShortTenancyID      = "t"
	FlagShortCompartment    = "c"
	FlagShortHelp           = "h"
	FlagShortDebug          = "d"
	FlagShortLimit          = "m"
	FlagShortPage           = "p"
	FlagShortJSON           = "j"
	FlagShortVersion        = "v"
	FlagShortSort           = "s"
	FlagShortRealm          = "r"
	FlagShortFilter         = "f"
	FlagShortAllInformation = "A"
	FlagShortTenancyScope   = "T"
)

// FlagDescriptions contains help text for flags
const (
	FlagDescLogLevel       = "Set the log verbosity debug,"
	FlagDescDebug          = "Enable debug logging"
	FlagDescTenancyID      = "OCI tenancy OCID"
	FlagDescTenancyName    = "Tenancy name"
	FlagDescCompartment    = "OCI compartment name"
	FlagDescHelp           = "help for ocloud (shorthand: -h)"
	FlagDescLimit          = "Maximum number of records to display per page"
	FlagDescPage           = "Page number to display"
	FlagDescJSON           = "Output information in JSON format"
	FlagDescVersion        = "Print the version number of ocloud CLI"
	FlagDescSort           = "Sort results by field (name or cidr)"
	FlagDescRealm          = "Filter by realm (e.g., OC1, OC2, OC2)"
	FlagDescFilter         = "Filter regions by prefix (e.g., us, eu, ap)"
	FlagDescAllInformation = "Show all information"
	FlagDescScope          = "Listing scope: compartment or tenancy"
	FlagDescTenancyScope   = "Shortcut: list at tenancy level (overrides --scope)"
)

// Flag values and defaults
const (
	FlagValueTrue     = "true"
	FlagValueInfo     = "info"
	FlagValueHelpMode = "help-mode"
)

// Flag prefixes and special strings
const (
	FlagPrefixShortHelp    = "-h"
	FlagPrefixLongHelp     = "--help"
	FlagPrefixColor        = "--color"
	FlagPrefixDebug        = "--debug"
	FlagPrefixShortDebug   = "-d"
	FlagPrefixVersion      = "--version"
	FlagPrefixShortVersion = "-v"
	CobraAnnotationKey     = "cobra_annotation_flag_set_by_cobra"
)

// Environment variable keys
const (
	EnvKeyProfile        = "OCI_CLI_PROFILE"
	EnvKeyCLITenancy     = "OCI_CLI_TENANCY"
	EnvKeyTenancyName    = "OCI_TENANCY_NAME"
	EnvKeyCompartment    = "OCI_COMPARTMENT"
	EnvKeyAutoRefresher  = "OCI_AUTH_AUTO_REFRESHER"
	EnvKeyRegion         = "OCI_REGION"
	EnvKeyTenancyMapPath = "OCI_TENANCY_MAP_PATH"
)

// File/system names & defaults
const (
	DefaultProfileName = "DEFAULT"

	OCIConfigDirName        = ".oci"
	OCIConfigFileName       = "config"
	OCloudDefaultDirName    = ".ocloud"
	OCloudScriptsDirName    = "scripts"
	OCISessionsDirName      = "sessions"
	TenancyMapFileName      = "tenancy-map.yaml"
	OCIRefresherPIDFileName = "refresher.pid"
)
