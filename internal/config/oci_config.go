package config

import (
	"os"
)

type OCIConfig struct {
	Profile       string
	TenantID      string
	CompartmentID string
	Region        string
	User          string
	Key           string
	Fingerprint   string
}

func (c *OCIConfig) Config() error {
	c.Profile = ""
	c.TenantID = ""
	c.CompartmentID = ""
	c.Region = ""
	c.User = ""
	c.Key = ""
	c.Fingerprint = ""
	return nil
}

//func NewOCIConfig() common.ConfigurationProvider {
//	var config common.ConfigurationProvider
//	profile := os.Getenv("OCI_PROFILE")
//
//	if profile == "" {
//		config = common.DefaultConfigProvider()
//	}else{
//		configPath := getHomeDir() + "./oci/config"
//		config = common.CustomProfileConfigProvider(configPath, profile)
//	}
//	return config
//}

func getHomeDir() string {
	return os.Getenv("HOME")
}
