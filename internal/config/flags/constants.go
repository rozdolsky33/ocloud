// Package flags provide names, shorthands, descriptions, env keys, and defaults for the CLI.
package flags

// ============================================================================
// Flag Names (common)
// ============================================================================
const (
	FlagNameLogLevel     = "log-level"
	FlagNameDebug        = "debug"
	FlagNameTenancyID    = "tenancy-id"
	FlagNameTenancyName  = "tenancy-name"
	FlagNameCompartment  = "compartment"
	FlagNameHelp         = "help"
	FlagNameColor        = "color"
	FlagNameLimit        = "limit"
	FlagNamePage         = "page"
	FlagNameJSON         = "json"
	FlagNameVersion      = "version"
	FlagNameAll          = "all"
	FlagNameSort         = "sort"
	FlagNameRealm        = "realm"
	FlagNameFilter       = "filter"
	FlagNameScope        = "scope"
	FlagNameTenancyScope = "tenancy-scope"
)

// Flag Names (network toggles)
const (
	FlagNameGateway  = "gateway"
	FlagNameSubnet   = "subnet"
	FlagNameNsg      = "nsg"
	FlagNameRoute    = "route-table"
	FlagNameSecurity = "security-list"
)

// ============================================================================
// Flag Shorthands
// ============================================================================
const (
	// Common
	FlagShortTenancyID    = "t"
	FlagShortCompartment  = "c"
	FlagShortHelp         = "h"
	FlagShortDebug        = "d"
	FlagShortLimit        = "m"
	FlagShortPage         = "p"
	FlagShortJSON         = "j"
	FlagShortVersion      = "v"
	FlagShortSort         = "s"
	FlagShortRealm        = "r"
	FlagShortFilter       = "f"
	FlagShortAll          = "A"
	FlagShortTenancyScope = "T"

	// Network toggles (avoid collisions with common flags)
	FlagShortGateway  = "G"
	FlagShortSubnet   = "S"
	FlagShortNsg      = "N"
	FlagShortRoute    = "R"
	FlagShortSecurity = "L"
)

// ============================================================================
// Flag Descriptions
// ============================================================================
const (
	// Common
	FlagDescLogLevel     = "Set the log verbosity (e.g., info, debug)"
	FlagDescDebug        = "Enable debug logging"
	FlagDescTenancyID    = "OCI tenancy OCID"
	FlagDescTenancyName  = "Tenancy name"
	FlagDescCompartment  = "OCI compartment name or OCID"
	FlagDescHelp         = "Show help"
	FlagDescLimit        = "Maximum number of records to display per page"
	FlagDescPage         = "Page number to display"
	FlagDescJSON         = "Output information in JSON format"
	FlagDescVersion      = "Print the ocloud CLI version"
	FlagDescSort         = "Sort results by field (e.g., name, cidr)"
	FlagDescRealm        = "Filter by realm (e.g., OC1, OC2, OC3)"
	FlagDescFilter       = "Filter regions by prefix (e.g., us, eu, ap)"
	FlagDescAll          = "Show all information"
	FlagDescScope        = "Listing scope: compartment or tenancy"
	FlagDescTenancyScope = "Shortcut: list at tenancy level (overrides --scope)"

	// Network
	FlagDescGateway  = "Display gateway information"
	FlagDescSubnet   = "Display subnet information"
	FlagDescNsg      = "Display network security group information"
	FlagDescRoute    = "Display route table information"
	FlagDescSecurity = "Display security list information"
)

// ============================================================================
// Flag Values / Special
// ============================================================================
const (
	FlagValueTrue     = "true"
	FlagValueInfo     = "info"
	FlagValueHelpMode = "help-mode"
)

// ============================================================================
// CLI Prefixes / Annotations (for parsing/help)
// ============================================================================
const (
	FlagPrefixShortHelp    = "-h"
	FlagPrefixLongHelp     = "--help"
	FlagPrefixColor        = "--color"
	FlagPrefixDebug        = "--debug"
	FlagPrefixShortDebug   = "-d"
	FlagPrefixVersion      = "--version"
	FlagPrefixShortVersion = "-v"

	CobraAnnotationKey = "cobra_annotation_flag_set_by_cobra"
)

// ============================================================================
// Environment Keys
// ============================================================================
const (
	EnvKeyProfile        = "OCI_CLI_PROFILE"
	EnvKeyCLITenancy     = "OCI_CLI_TENANCY"
	EnvKeyTenancyName    = "OCI_TENANCY_NAME"
	EnvKeyCompartment    = "OCI_COMPARTMENT"
	EnvKeyAutoRefresher  = "OCI_AUTH_AUTO_REFRESHER"
	EnvKeyRegion         = "OCI_REGION"
	EnvKeyTenancyMapPath = "OCI_TENANCY_MAP_PATH"
)

// ============================================================================
// Filenames / Defaults
// ============================================================================
const (
	DefaultProfileName = "DEFAULT"

	OCIConfigDirName     = ".oci"
	OCIConfigFileName    = "config"
	OCloudDefaultDirName = ".ocloud"
	OCloudScriptsDirName = "scripts"
	OCISessionsDirName   = "sessions"
	TenancyMapFileName   = "tenancy-map.yaml"
	OCIRefresherPIDFile  = "refresher.pid"
)
