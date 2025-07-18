package info

import (
	"bytes"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestPrintMappingsFile tests the PrintMappingsFile function
func TestPrintMappingsFile(t *testing.T) {
	// Create mock data
	mockMappings := []appConfig.MappingsFile{
		{
			Environment:  "prod",
			Tenancy:      "mytenancy1",
			TenancyID:    "ocid1.tenancy.oc1..aaaaaaaabcdefg1",
			Realm:        "OC1",
			Compartments: "comp1 comp2",
			Regions:      "us-ashburn-1 us-phoenix-1",
		},
		{
			Environment:  "dev",
			Tenancy:      "mytenancy2",
			TenancyID:    "ocid1.tenancy.oc1..aaaaaaaabcdefg2",
			Realm:        "OC2",
			Compartments: "comp3 comp4",
			Regions:      "eu-frankfurt-1 uk-london-1",
		},
	}

	// Create test cases
	testCases := []struct {
		name          string
		mappings      []appConfig.MappingsFile
		useJSON       bool
		expectedEmpty bool
	}{
		{
			name:          "Print mappings in table format",
			mappings:      mockMappings,
			useJSON:       false,
			expectedEmpty: false,
		},
		{
			name:          "Print mappings in JSON format",
			mappings:      mockMappings,
			useJSON:       true,
			expectedEmpty: false,
		},
		{
			name:          "Print empty mappings in table format",
			mappings:      []appConfig.MappingsFile{},
			useJSON:       false,
			expectedEmpty: true,
		},
		{
			name:          "Print empty mappings in JSON format",
			mappings:      []appConfig.MappingsFile{},
			useJSON:       true,
			expectedEmpty: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a buffer to capture output
			var buf bytes.Buffer

			// Create a mock application context with the buffer as stdout
			appCtx := &app.ApplicationContext{
				Stdout: &buf,
			}

			// Call PrintMappingsFile
			err := PrintMappingsFile(tc.mappings, appCtx, tc.useJSON)

			// Verify the results
			assert.NoError(t, err)

			// Check the output
			output := buf.String()
			if tc.useJSON {
				assert.Contains(t, output, "{")
				assert.Contains(t, output, "}")
				if !tc.expectedEmpty {
					assert.Contains(t, output, "environment")
					assert.Contains(t, output, "tenancy")
					assert.Contains(t, output, "realm")
				}
			} else {
				if !tc.expectedEmpty {
					assert.Contains(t, output, "ENVIRONMENT")
					assert.Contains(t, output, "TENANCY")
					assert.Contains(t, output, "COMPARTMENTS")
					assert.Contains(t, output, "REGIONS")
				}
			}
		})
	}
}

// TestSplitTextByMaxWidth tests the splitTextByMaxWidth function
func TestSplitTextByMaxWidth(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput []string
	}{
		{
			name:           "Empty string",
			input:          "",
			expectedOutput: []string{""},
		},
		{
			name:           "Single word",
			input:          "word",
			expectedOutput: []string{"word"},
		},
		{
			name:           "Short text",
			input:          "short text",
			expectedOutput: []string{"short text"},
		},
		{
			name:  "Long text",
			input: "this is a very long text that should be split into multiple lines because it exceeds the maximum width",
			expectedOutput: []string{
				"this is a very long text that",
				"should be split into multiple lines",
				"because it exceeds the maximum",
				"width",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := splitTextByMaxWidth(tc.input)
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}

// TestGroupMappingsByRealm tests the groupMappingsByRealm function
func TestGroupMappingsByRealm(t *testing.T) {
	// Create mock data
	mockMappings := []appConfig.MappingsFile{
		{
			Environment:  "prod",
			Tenancy:      "mytenancy1",
			TenancyID:    "ocid1.tenancy.oc1..aaaaaaaabcdefg1",
			Realm:        "OC1",
			Compartments: "comp1 comp2",
			Regions:      "us-ashburn-1 us-phoenix-1",
		},
		{
			Environment:  "dev",
			Tenancy:      "mytenancy2",
			TenancyID:    "ocid1.tenancy.oc1..aaaaaaaabcdefg2",
			Realm:        "OC2",
			Compartments: "comp3 comp4",
			Regions:      "eu-frankfurt-1 uk-london-1",
		},
		{
			Environment:  "test",
			Tenancy:      "mytenancy3",
			TenancyID:    "ocid1.tenancy.oc1..aaaaaaaabcdefg3",
			Realm:        "OC1",
			Compartments: "comp5 comp6",
			Regions:      "us-ashburn-1 us-phoenix-1",
		},
	}

	// Call groupMappingsByRealm
	result := groupMappingsByRealm(mockMappings)

	// Verify the results
	assert.Equal(t, 2, len(result))
	assert.Equal(t, 2, len(result["OC1"]))
	assert.Equal(t, 1, len(result["OC2"]))
	assert.Equal(t, "mytenancy1", result["OC1"][0].Tenancy)
	assert.Equal(t, "mytenancy3", result["OC1"][1].Tenancy)
	assert.Equal(t, "mytenancy2", result["OC2"][0].Tenancy)
}
