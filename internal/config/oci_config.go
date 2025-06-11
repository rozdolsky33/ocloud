package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
var DefaultTenancyMapPath = filepath.Join(getUserHomeDir(), ".oci", "tenancy-map.yaml")

const (
	defaultProfile = "DEFAULT"
	envProfileKey  = "OCI_CLI_PROFILE"
	configDir      = ".oci"
	configFile     = "config"

	// EnvTenancyMapPath is the environment variable key used to specify the file path for the OCI tenancy map configuration.
	EnvTenancyMapPath = "OCI_TENANCY_MAP_PATH"
)

// LoadOCIConfig picks the profile from env or default, and logs at debug level.
func LoadOCIConfig() common.ConfigurationProvider {
	profile := GetOCIProfile()
	if profile == defaultProfile {
		logrus.Debug("using default profile")
		return common.DefaultConfigProvider()
	}

	logrus.Debugf("using profile %s", profile)
	path := filepath.Join(getUserHomeDir(), configDir, configFile)
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
	path := tenancyMapPath()
	logrus.Debugf("looking up tenancy %q in map %s", tenancyName, path)

	tenancies, err := LoadTenancyMap()
	if err != nil {
		return "", err
	}

	for _, env := range tenancies {
		if env.Tenancy == tenancyName {
			logrus.Debugf("found tenancy %q => OCID %s", tenancyName, env.TenancyID)
			return env.TenancyID, nil
		}
	}

	lookupErr := fmt.Errorf("tenancy %q not found in %s", tenancyName, path)
	logrus.Warnf("tenancy lookup failed: %v", lookupErr)
	return "", errors.Wrap(lookupErr, "tenancy lookup failed - please check that the tenancy name is correct and exists in the mapping file")
}

// LoadTenancyMap loads the tenancy mapping from disk at tenancyMapPath.
// It logs debug information and returns a slice of OciTenancyEnvironment.
func LoadTenancyMap() ([]OCITenancyEnvironment, error) {
	path := tenancyMapPath()
	logrus.Debugf("loading tenancy map from %s", path)

	if err := ensureFile(path); err != nil {
		logrus.Warnf("tenancy mapping file not found: %v", err)
		return nil, errors.Wrapf(err, "tenancy mapping file not found (%s) - this is normal if you're not using tenancy name lookup", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("failed to read tenancy mapping file %s: %v", path, err)
		return nil, errors.Wrapf(err, "failed to read tenancy mapping file (%s)", path)
	}

	var tenancies []OCITenancyEnvironment
	if err := yaml.Unmarshal(data, &tenancies); err != nil {
		logrus.Errorf("failed to parse tenancy mapping file %s: %v", path, err)
		return nil, errors.Wrapf(err, "failed to parse tenancy mapping file (%s) - please check that the file is valid YAML", path)
	}

	logrus.Debugf("loaded %d tenancy mapping entries", len(tenancies))
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

// getUserHomeDir returns the path to the current user's home directory or exits if unable to determine it.
func getUserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	return dir
}

// tenancyMapPath returns either the overridden path or the default.
func tenancyMapPath() string {
	if p := os.Getenv(EnvTenancyMapPath); p != "" {
		logrus.Debugf("using tenancy map from env: %s", p)
		return p
	}
	return DefaultTenancyMapPath
}
