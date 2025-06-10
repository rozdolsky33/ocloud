package config

import (
	"fmt"
	"os"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/spf13/viper"
)

// ConfigProvider is the OCI SDK auth provider
var ConfigProvider common.ConfigurationProvider

// CompartmentID holds the OCID for the target compartment
var CompartmentID string

// InitAuth sets up tenancy/profile/region/env for the SDK without requiring a compartment
func InitAuth() error {
	if t := viper.GetString("tenancy"); t != "" {
		os.Setenv("OCI_TENANCY", t)
	}
	if p := viper.GetString("profile"); p != "" {
		os.Setenv("OCI_PROFILE", p)
	}
	if r := viper.GetString("region"); r != "" {
		os.Setenv("OCI_REGION", r)
	}
	ConfigProvider = common.DefaultConfigProvider()
	return nil
}

// Init runs full init: auth + compartment check
func Init() error {
	if err := InitAuth(); err != nil {
		return err
	}
	CompartmentID = viper.GetString("compartment")
	if CompartmentID == "" {
		return fmt.Errorf("compartment OCID must be set (flag -c or env OCI_COMPARTMENT)")
	}
	return nil
}
