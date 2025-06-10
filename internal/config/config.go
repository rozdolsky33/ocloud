package config

import (
	"path/filepath"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	defaultProfile = "DEFAULT"
	envProfileKey  = "OCI_CLI_PROFILE"
	configDir      = ".oci"
	configFile     = "config"
)

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

func userHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	return dir
}
