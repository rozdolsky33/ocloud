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

var DefaultTenancyMapPath = filepath.Join(userHomeDir(), ".oci", "tenancy-map.yaml")

const (
	defaultProfile    = "DEFAULT"
	envProfileKey     = "OCI_CLI_PROFILE"
	configDir         = ".oci"
	configFile        = "config"
	EnvTenancyMapPath = "OCI_TENANCY_MAP_PATH"
)

type OciTenancyEnvironment struct {
	Environment  string `yaml:"environment"`
	Tenancy      string `yaml:"tenancy"`
	TenancyId    string `yaml:"tenancy_id"`
	Realm        string `yaml:"realm"`
	Compartments string `yaml:"compartments"`
	Regions      string `yaml:"regions"`
}

// LoadOCIConfig picks the profile from env or default, and logs at debug level.
func LoadOCIConfig() common.ConfigurationProvider {
	profile := GetOCIProfile()
	if profile == defaultProfile {
		logrus.Debug("using default profile")
		return common.DefaultConfigProvider()
	}

	logrus.Debugf("using profile %s", profile)
	path := filepath.Join(userHomeDir(), configDir, configFile)
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
	id, err := LoadOCIConfig().TenancyOCID()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve tenancy OCID from OCI config")
	}
	return id, nil
}

// LookUpTenancyID locates the OCID for a given tenancy name.
// It returns an error if the map cannot be loaded or if the name isn't found.
func LookUpTenancyID(tenancyName string) (string, error) {
	path := tenancyMapPath()
	logrus.Debugf("looking up tenancy %q in map %s", tenancyName, path)

	tenancies, err := LoadTenancyMap()
	if err != nil {
		return "", err
	}

	for _, env := range tenancies {
		if env.Tenancy == tenancyName {
			logrus.Debugf("found tenancy %q => OCID %s", tenancyName, env.TenancyId)
			return env.TenancyId, nil
		}
	}

	lookupErr := fmt.Errorf("tenancy %q not found in %s", tenancyName, path)
	logrus.Warnf("tenancy lookup failed: %v", lookupErr)
	return "", errors.Wrap(lookupErr, "tenancy lookup failed")
}

// LoadTenancyMap loads the tenancy mapping from disk at tenancyMapPath.
// It logs debug information and returns a slice of OciTenancyEnvironment.
func LoadTenancyMap() ([]OciTenancyEnvironment, error) {
	path := tenancyMapPath()
	logrus.Debugf("loading tenancy map from %s", path)

	if err := ensureFile(path); err != nil {
		logrus.Warnf("tenancy mapping file not found: %v", err)
		return nil, errors.Wrapf(err, "tenancy mapping file not found (%s)", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logrus.Errorf("failed to read tenancy mapping file %s: %v", path, err)
		return nil, errors.Wrapf(err, "failed to read tenancy mapping file (%s)", path)
	}

	var tenancies []OciTenancyEnvironment
	if err := yaml.Unmarshal(data, &tenancies); err != nil {
		logrus.Errorf("failed to parse tenancy mapping file %s: %v", path, err)
		return nil, errors.Wrapf(err, "failed to parse tenancy mapping file (%s)", path)
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

func userHomeDir() string {
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
