package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// MappingsFile represents the configuration for an OCI tenancy environment.
// It includes details such as environment, tenancy, tenancy ID, realm, compartments, and regions.
type MappingsFile struct {
	Environment  string   `yaml:"environment"`
	Tenancy      string   `yaml:"tenancy"`
	TenancyID    string   `yaml:"tenancy_id"`
	Realm        string   `yaml:"realm"`
	Compartments []string `yaml:"compartments"`
	Regions      []string `yaml:"regions"`
}

// UnmarshalYAML implements custom unmarshaling for MappingsFile to handle
// both string and string array values for Compartments and Regions fields.
func (m *MappingsFile) UnmarshalYAML(value *yaml.Node) error {
	// Define a temporary struct with the same fields but using interface{} for arrays
	// to handle both string and []string values
	type TempMappingsFile struct {
		Environment  string      `yaml:"environment"`
		Tenancy      string      `yaml:"tenancy"`
		TenancyID    string      `yaml:"tenancy_id"`
		Realm        string      `yaml:"realm"`
		Compartments interface{} `yaml:"compartments"`
		Regions      interface{} `yaml:"regions"`
	}

	// Unmarshal into the temporary struct
	var temp TempMappingsFile
	if err := value.Decode(&temp); err != nil {
		return err
	}

	// Copy the simple fields
	m.Environment = temp.Environment
	m.Tenancy = temp.Tenancy
	m.TenancyID = temp.TenancyID
	m.Realm = temp.Realm

	// Handle Compartments field which could be a string or []string
	switch v := temp.Compartments.(type) {
	case string:
		// If it's a single string, convert to a slice with one element
		m.Compartments = []string{v}
	case []interface{}:
		// If it's already a slice, convert each element to string
		m.Compartments = make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				m.Compartments[i] = str
			} else {
				return fmt.Errorf("compartments[%d] is not a string: %v", i, item)
			}
		}
	case []string:
		// If it's already a []string, use it directly
		m.Compartments = v
	case nil:
		// If it's nil, use an empty slice
		m.Compartments = []string{}
	default:
		return fmt.Errorf("compartments must be a string or array of strings, got %T", v)
	}

	// Handle Regions field which could be a string or []string
	switch v := temp.Regions.(type) {
	case string:
		// If it's a single string, convert to a slice with one element
		m.Regions = []string{v}
	case []interface{}:
		// If it's already a slice, convert each element to string
		m.Regions = make([]string, len(v))
		for i, item := range v {
			if str, ok := item.(string); ok {
				m.Regions[i] = str
			} else {
				return fmt.Errorf("regions[%d] is not a string: %v", i, item)
			}
		}
	case []string:
		// If it's already a []string, use it directly
		m.Regions = v
	case nil:
		// If it's nil, use an empty slice
		m.Regions = []string{}
	default:
		return fmt.Errorf("regions must be a string or array of strings, got %T", v)
	}

	return nil
}

// MarshalYAML implements custom marshaling for MappingsFile to ensure
// lowercase field names and maintain the specific field order.
func (m MappingsFile) MarshalYAML() (interface{}, error) {
	// Create a mapping node
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}

	// Add fields in the desired order with lowercase keys
	// First add the environment field
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "environment"},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: m.Environment},
	)

	// Add tenancy field
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "tenancy"},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: m.Tenancy},
	)

	// Add tenancy_id field
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "tenancy_id"},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: m.TenancyID},
	)

	// Add realm field
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "realm"},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: m.Realm},
	)

	// Add a compartment field
	compartmentsKey := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "compartments"}
	compartmentsValue := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

	for _, comp := range m.Compartments {
		compartmentsValue.Content = append(compartmentsValue.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: comp},
		)
	}

	node.Content = append(node.Content, compartmentsKey, compartmentsValue)

	// Add a region field
	regionsKey := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "regions"}
	regionsValue := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}

	for _, reg := range m.Regions {
		regionsValue.Content = append(regionsValue.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: reg},
		)
	}

	node.Content = append(node.Content, regionsKey, regionsValue)

	return node, nil
}
