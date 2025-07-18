package config

// MappingsFile represents the configuration for an OCI tenancy environment.
// It includes details such as environment, tenancy, tenancy ID, realm, compartments, and regions.
type MappingsFile struct {
	Environment  string `yaml:"environment" json:"environment"`
	Tenancy      string `yaml:"tenancy" json:"tenancy"`
	TenancyID    string `yaml:"tenancy_id" json:"tenancy_id"`
	Realm        string `yaml:"realm" json:"realm"`
	Compartments string `yaml:"compartments" json:"compartments"`
	Regions      string `yaml:"regions" json:"regions"`
}
