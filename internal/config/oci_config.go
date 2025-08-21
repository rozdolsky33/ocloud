package config

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"

	"path/filepath"

	"gopkg.in/yaml.v3"

	"os"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/pkg/errors"
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
	dir, err := GetUserHomeDir()
	if err != nil {
		logger.Logger.V(logger.Debug).Info("failed to get user home directory for tenancy map path", "error", err)
		return ""
	}
	return filepath.Join(dir, flags.OCIConfigDirName, flags.OCloudDefaultDirName, flags.TenancyMapFileName)
}()

// LoadOCIConfig picks the profile from env or default, and logs at debug level.
// If there's an error getting the home directory, it falls back to the default provider.
func LoadOCIConfig() common.ConfigurationProvider {
	logger.Logger.V(logger.Info).Info("Loading OCI configuration...")
	profile := GetOCIProfile()
	if profile == flags.DefaultProfileName {
		logger.LogWithLevel(logger.Logger, logger.Trace, "using default profile")
		return common.DefaultConfigProvider()
	}

	logger.LogWithLevel(logger.Logger, logger.Trace, "using profile", "profile", profile)

	homeDir, err := GetUserHomeDir()
	if err != nil {
		logger.Logger.Error(err, "failed to get user home directory for config path, falling back to default provider")
		return common.DefaultConfigProvider()
	}

	path := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCIConfigFileName)
	return common.CustomProfileConfigProvider(path, profile)
}

// GetOCIProfile returns OCI_CLI_PROFILE or "DEFAULT".
func GetOCIProfile() string {
	if p := os.Getenv(flags.EnvKeyProfile); p != "" {
		return p
	}
	return flags.DefaultProfileName
}

// GetTenancyOCID fetches the tenancy OCID (error on failure).
func GetTenancyOCID() (string, error) {
	logger.Logger.V(logger.Debug).Info("Attempting to get tenancy OCID.")
	// Use mock function if set (for testing)
	if MockGetTenancyOCID != nil {
		return MockGetTenancyOCID()
	}

	id, err := LoadOCIConfig().TenancyOCID()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve tenancy OCID from OCI config")
	}
	logger.Logger.V(logger.Debug).Info("Successfully retrieved tenancy OCID.", "tenancyID", id)
	return id, nil
}

// LookupTenancyID locates the OCID for a given tenancy name.
// It returns an error if the map cannot be loaded or if the name isn't found.
func LookupTenancyID(tenancyName string) (string, error) {
	logger.Logger.V(logger.Debug).Info("Attempting to lookup tenancy ID", "tenancyName", tenancyName)
	// Use mock function if set (for testing)
	if MockLookupTenancyID != nil {
		return MockLookupTenancyID(tenancyName)
	}

	path := TenancyMapPath()
	logger.LogWithLevel(logger.Logger, logger.Trace, "looking up tenancy in map", "tenancy", tenancyName, "path", path)

	tenancies, err := LoadTenancyMap()
	if err != nil {
		return "", err
	}

	for _, env := range tenancies {
		if env.Tenancy == tenancyName {
			logger.LogWithLevel(logger.Logger, logger.Trace, "found tenancy", "tenancy", tenancyName, "tenancyID", env.TenancyID)
			logger.Logger.V(logger.Debug).Info("Successfully looked up tenancy ID.", "tenancyName", tenancyName, "tenancyID", env.TenancyID)
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
	logger.Logger.V(logger.Debug).Info("Attempting to load tenancy map.")
	path := TenancyMapPath()
	logger.LogWithLevel(logger.Logger, logger.Trace, "loading tenancy map", "path", path)

	if err := ensureFile(path); err != nil {
		logger.Logger.Info("tenancy mapping file not found", "error", err)
		return nil, errors.Wrapf(err, "tenancy mapping file not found (%s) - this is normal if you're not using tenancy name lookup. To set up the mapping file, create a YAML file at %s or set the %s environment variable to point to your mapping file. The file should contain entries mapping tenancy names to OCIDs. Example:\n- environment: OcluodOps\n  tenancy: cncloudops\n  tenancy_id: ocid1.tenancy.oc1..aaaaaaaasrwe3nsfsidfxzxyzct\n  realm: OC1\n  compartments:\n    - sandbox\n    - uat\n    - prod\n  regions:\n    - us-chicago-1\n    - us-ashburn-1\n", path, DefaultTenancyMapPath, flags.EnvKeyTenancyMapPath)
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

	logger.LogWithLevel(logger.Logger, logger.Trace, "loaded tenancy mapping entries", "count", len(tenancies))
	logger.Logger.V(logger.Debug).Info("Successfully loaded tenancy map.", "count", len(tenancies))
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

// GetUserHomeDir returns the path to the current user's home directory or an error if unable to determine it.
func GetUserHomeDir() (string, error) {
	logger.Logger.V(logger.Debug).Info("Attempting to get user home directory.")
	dir, err := os.UserHomeDir()
	if err != nil {
		logger.Logger.Error(err, "failed to get user home directory")
		return "", fmt.Errorf("getting user home directory: %w", err)
	}
	logger.Logger.V(logger.Debug).Info("Successfully retrieved user home directory.", "directory", dir)
	return dir, nil
}

// TenancyMapPath returns either the overridden path or the default.
func TenancyMapPath() string {
	if p := os.Getenv(flags.EnvKeyTenancyMapPath); p != "" {
		logger.LogWithLevel(logger.Logger, logger.Trace, "using tenancy map from env", "path", p)
		return p
	}
	return DefaultTenancyMapPath
}
