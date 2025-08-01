package config

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/logger"

	"gopkg.in/yaml.v3"
	"path/filepath"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/pkg/errors"
	"os"
)

// For testing purposes
var (
	// MockGetTenancyOCID allows tests to override the GetTenancyOCID function
	MockGetTenancyOCID func() (string, error)
	// MockLookupTenancyID allows tests to override the LookupTenancyID function
	MockLookupTenancyID func(tenancyName string) (string, error)
)

// DefaultTenancyMapPath defines the default file path for the OCI tenancy map configuration in the user's home directory.
// If the home directory cannot be determined, it falls back to an empty string.
var DefaultTenancyMapPath = func() string {
	dir, err := getUserHomeDir()
	if err != nil {
		logger.LogWithLevel(logger.Logger, 1, "failed to get user home directory for tenancy map path", "error", err)
		return ""
	}
	return filepath.Join(dir, ".oci", "tenancy-map.yaml")
}()

const (
	defaultProfile = "DEFAULT"
	envProfileKey  = "OCI_CLI_PROFILE"
	configDir      = ".oci"
	configFile     = "config"
	// EnvTenancyMapPath is the environment variable key used to specify the file path for the OCI tenancy map configuration.
	EnvTenancyMapPath = "OCI_TENANCY_MAP_PATH"
)

// LoadOCIConfig picks the profile from env or default, and logs at debug level.
// If there's an error getting the home directory, it falls back to the default provider.
func LoadOCIConfig() common.ConfigurationProvider {
	profile := GetOCIProfile()
	if profile == defaultProfile {
		logger.LogWithLevel(logger.Logger, 3, "using default profile")
		return common.DefaultConfigProvider()
	}

	logger.LogWithLevel(logger.Logger, 3, "using profile", "profile", profile)

	homeDir, err := getUserHomeDir()
	if err != nil {
		logger.Logger.Error(err, "failed to get user home directory for config path, falling back to default provider")
		return common.DefaultConfigProvider()
	}

	path := filepath.Join(homeDir, configDir, configFile)
	return common.CustomProfileConfigProvider(path, profile)
}

// GetOCIProfile returns OCI_CLI_PROFILE or "DEFAULT".
func GetOCIProfile() string {
	if p := os.Getenv(envProfileKey); p != "" {
		return p
	}
	return defaultProfile
}

// GetTenancyOCID fetches the tenancy OCID (error on failure).
func GetTenancyOCID() (string, error) {
	// Use mock function if set (for testing)
	if MockGetTenancyOCID != nil {
		return MockGetTenancyOCID()
	}

	// Normal implementation
	id, err := LoadOCIConfig().TenancyOCID()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve tenancy OCID from OCI config")
	}
	return id, nil
}

// LookupTenancyID locates the OCID for a given tenancy name.
// It returns an error if the map cannot be loaded or if the name isn't found.
func LookupTenancyID(tenancyName string) (string, error) {
	// Use mock function if set (for testing)
	if MockLookupTenancyID != nil {
		return MockLookupTenancyID(tenancyName)
	}

	// Normal implementation
	path := TenancyMapPath()
	logger.LogWithLevel(logger.Logger, 3, "looking up tenancy in map", "tenancy", tenancyName, "path", path)

	tenancies, err := LoadTenancyMap()
	if err != nil {
		return "", err
	}

	for _, env := range tenancies {
		if env.Tenancy == tenancyName {
			logger.LogWithLevel(logger.Logger, 3, "found tenancy", "tenancy", tenancyName, "tenancyID", env.TenancyID)
			return env.TenancyID, nil
		}
	}

	lookupErr := fmt.Errorf("tenancy %q not found in %s", tenancyName, path)
	logger.Logger.Info("tenancy lookup failed", "error", lookupErr)
	return "", errors.Wrap(lookupErr, "tenancy lookup failed - please check that the tenancy name is correct and exists in the mapping file")
}

// LoadTenancyMap loads the tenancy mapping from the disk at TenancyMapPath.
// It logs debug information and returns a slice of OciTenancyEnvironment.
func LoadTenancyMap() ([]MappingsFile, error) {
	path := TenancyMapPath()
	logger.LogWithLevel(logger.Logger, 3, "loading tenancy map", "path", path)

	if err := ensureFile(path); err != nil {
		logger.Logger.Info("tenancy mapping file not found", "error", err)
		return nil, errors.Wrapf(err, "tenancy mapping file not found (%s) - this is normal if you're not using tenancy name lookup. To set up the mapping file, create a YAML file at %s or set the %s environment variable to point to your mapping file. The file should contain entries mapping tenancy names to OCIDs. Example:\n- environment: prod\n  tenancy: mytenancy\n  tenancy_id: ocid1.tenancy.oc1..aaaaaaaabcdefghijklmnopqrstuvwxyz\n  realm: oc1\n  compartments: mycompartment\n  regions: us-ashburn-1", path, DefaultTenancyMapPath, EnvTenancyMapPath)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logger.Logger.Error(err, "failed to read tenancy mapping file", "path", path)
		return nil, errors.Wrapf(err, "failed to read tenancy mapping file (%s)", path)
	}

	var tenancies []MappingsFile
	if err := yaml.Unmarshal(data, &tenancies); err != nil {
		logger.Logger.Error(err, "failed to parse tenancy mapping file", "path", path)
		return nil, errors.Wrapf(err, "failed to parse tenancy mapping file (%s) - please check that the file is valid YAML", path)
	}

	logger.LogWithLevel(logger.Logger, 3, "loaded tenancy mapping entries", "count", len(tenancies))
	return tenancies, nil
}

// ensureFile verifies the given path exists and is not a directory.
func ensureFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("path %s is a directory, expected a file", path)
	}
	return nil
}

// getUserHomeDir returns the path to the current user's home directory or an error if unable to determine it.
func getUserHomeDir() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Logger.Error(err, "failed to get user home directory")
		return "", fmt.Errorf("getting user home directory: %w", err)
	}
	return dir, nil
}

// TenancyMapPath returns either the overridden path or the default.
func TenancyMapPath() string {
	if p := os.Getenv(EnvTenancyMapPath); p != "" {
		logger.LogWithLevel(logger.Logger, 3, "using tenancy map from env", "path", p)
		return p
	}
	return DefaultTenancyMapPath
}
