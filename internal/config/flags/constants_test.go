package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagNames(t *testing.T) {
	// Test common flag names
	assert.Equal(t, "log-level", FlagNameLogLevel)
	assert.Equal(t, "debug", FlagNameDebug)
	assert.Equal(t, "tenancy-id", FlagNameTenancyID)
	assert.Equal(t, "tenancy-name", FlagNameTenancyName)
	assert.Equal(t, "compartment", FlagNameCompartment)
	assert.Equal(t, "help", FlagNameHelp)
	assert.Equal(t, "color", FlagNameColor)
	assert.Equal(t, "limit", FlagNameLimit)
	assert.Equal(t, "page", FlagNamePage)
	assert.Equal(t, "json", FlagNameJSON)
	assert.Equal(t, "version", FlagNameVersion)
	assert.Equal(t, "all", FlagNameAll)
	assert.Equal(t, "sort", FlagNameSort)
	assert.Equal(t, "realm", FlagNameRealm)
	assert.Equal(t, "filter", FlagNameFilter)
	assert.Equal(t, "scope", FlagNameScope)
	assert.Equal(t, "tenancy-scope", FlagNameTenancyScope)

	// Test network toggle flag names
	assert.Equal(t, "gateway", FlagNameGateway)
	assert.Equal(t, "subnet", FlagNameSubnet)
	assert.Equal(t, "nsg", FlagNameNsg)
	assert.Equal(t, "route-table", FlagNameRoute)
	assert.Equal(t, "security-list", FlagNameSecurity)
}

func TestFlagShorthands(t *testing.T) {
	// Test common flag shorthands
	assert.Equal(t, "t", FlagShortTenancyID)
	assert.Equal(t, "c", FlagShortCompartment)
	assert.Equal(t, "h", FlagShortHelp)
	assert.Equal(t, "d", FlagShortDebug)
	assert.Equal(t, "m", FlagShortLimit)
	assert.Equal(t, "p", FlagShortPage)
	assert.Equal(t, "j", FlagShortJSON)
	assert.Equal(t, "v", FlagShortVersion)
	assert.Equal(t, "s", FlagShortSort)
	assert.Equal(t, "r", FlagShortRealm)
	assert.Equal(t, "f", FlagShortFilter)
	assert.Equal(t, "A", FlagShortAll)
	assert.Equal(t, "T", FlagShortTenancyScope)

	// Test network toggle flag shorthands
	assert.Equal(t, "G", FlagShortGateway)
	assert.Equal(t, "S", FlagShortSubnet)
	assert.Equal(t, "N", FlagShortNsg)
	assert.Equal(t, "R", FlagShortRoute)
	assert.Equal(t, "L", FlagShortSecurity)
}

func TestFlagDescriptions(t *testing.T) {
	// Test that flag descriptions are non-empty
	assert.NotEmpty(t, FlagDescLogLevel)
	assert.NotEmpty(t, FlagDescDebug)
	assert.NotEmpty(t, FlagDescTenancyID)
	assert.NotEmpty(t, FlagDescTenancyName)
	assert.NotEmpty(t, FlagDescCompartment)
	assert.NotEmpty(t, FlagDescHelp)
	assert.NotEmpty(t, FlagDescLimit)
	assert.NotEmpty(t, FlagDescPage)
	assert.NotEmpty(t, FlagDescJSON)
	assert.NotEmpty(t, FlagDescVersion)
	assert.NotEmpty(t, FlagDescSort)
	assert.NotEmpty(t, FlagDescRealm)
	assert.NotEmpty(t, FlagDescFilter)
	assert.NotEmpty(t, FlagDescAll)
	assert.NotEmpty(t, FlagDescScope)
	assert.NotEmpty(t, FlagDescTenancyScope)

	// Test network flag descriptions
	assert.NotEmpty(t, FlagDescGateway)
	assert.NotEmpty(t, FlagDescSubnet)
	assert.NotEmpty(t, FlagDescNsg)
	assert.NotEmpty(t, FlagDescRoute)
	assert.NotEmpty(t, FlagDescSecurity)
}

func TestFlagValues(t *testing.T) {
	assert.Equal(t, "true", FlagValueTrue)
	assert.Equal(t, "info", FlagValueInfo)
	assert.Equal(t, "help-mode", FlagValueHelpMode)
}

func TestFlagPrefixes(t *testing.T) {
	assert.Equal(t, "-h", FlagPrefixShortHelp)
	assert.Equal(t, "--help", FlagPrefixLongHelp)
	assert.Equal(t, "--color", FlagPrefixColor)
	assert.Equal(t, "--debug", FlagPrefixDebug)
	assert.Equal(t, "-d", FlagPrefixShortDebug)
	assert.Equal(t, "--version", FlagPrefixVersion)
	assert.Equal(t, "-v", FlagPrefixShortVersion)
	assert.Equal(t, "cobra_annotation_flag_set_by_cobra", CobraAnnotationKey)
}

func TestEnvironmentKeys(t *testing.T) {
	assert.Equal(t, "OCI_CLI_PROFILE", EnvKeyProfile)
	assert.Equal(t, "OCI_CLI_TENANCY", EnvKeyCLITenancy)
	assert.Equal(t, "OCI_TENANCY_NAME", EnvKeyTenancyName)
	assert.Equal(t, "OCI_COMPARTMENT", EnvKeyCompartment)
	assert.Equal(t, "OCI_AUTH_AUTO_REFRESHER", EnvKeyAutoRefresher)
	assert.Equal(t, "OCI_REGION", EnvKeyRegion)
	assert.Equal(t, "OCI_TENANCY_MAP_PATH", EnvKeyTenancyMapPath)
}

func TestDefaults(t *testing.T) {
	assert.Equal(t, "DEFAULT", DefaultProfileName)
	assert.Equal(t, ".oci", OCIConfigDirName)
	assert.Equal(t, "config", OCIConfigFileName)
	assert.Equal(t, ".ocloud", OCloudDefaultDirName)
	assert.Equal(t, "scripts", OCloudScriptsDirName)
	assert.Equal(t, "sessions", OCISessionsDirName)
	assert.Equal(t, "tenancy-map.yaml", TenancyMapFileName)
	assert.Equal(t, "refresher.pid", OCIRefresherPIDFile)
}
