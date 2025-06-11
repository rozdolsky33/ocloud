package config

// OCITenancyEnvironment represents the configuration for an OCI tenancy environment.
// It includes details such as environment, tenancy, tenancy ID, realm, compartments, and regions.
type OCITenancyEnvironment struct {
	Environment  string `yaml:"environment"`
	Tenancy      string `yaml:"tenancy"`
	TenancyID    string `yaml:"tenancy_id"`
	Realm        string `yaml:"realm"`
	Compartments string `yaml:"compartments"`
	Regions      string `yaml:"regions"`
}
