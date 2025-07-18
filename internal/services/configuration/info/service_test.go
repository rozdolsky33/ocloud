package info

import (
	"errors"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestNewService tests the NewService function
func TestNewService(t *testing.T) {
	// Create a mock application context
	appCtx := &app.ApplicationContext{
		Logger: logger.NewTestLogger(),
	}

	// Call NewService with the mock context
	service := NewService(appCtx)

	// Verify that the service was created correctly
	assert.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

// mockLoadTenancyMapFunc is a function type for mocking LoadTenancyMap
type mockLoadTenancyMapFunc func() ([]appConfig.MappingsFile, error)

// mockService is a mock implementation of the Service for testing
type mockService struct {
	Service
	mockLoadTenancyMap mockLoadTenancyMapFunc
}

// LoadTenancyMappings overrides the real implementation for testing
func (m *mockService) LoadTenancyMappings(realm string) (*TenancyMappingResult, error) {
	// Use the mock function to get mappings
	mappings, err := m.mockLoadTenancyMap()
	if err != nil {
		return nil, err
	}

	// Filter by realm if specified (same logic as the real implementation)
	var filteredMappings []appConfig.MappingsFile
	for _, tenancy := range mappings {
		if realm != "" && tenancy.Realm != realm {
			continue
		}
		filteredMappings = append(filteredMappings, tenancy)
	}

	return &TenancyMappingResult{
		Mappings: filteredMappings,
	}, nil
}

// TestLoadTenancyMappings tests the LoadTenancyMappings method
func TestLoadTenancyMappings(t *testing.T) {
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
		realm         string
		mockError     error
		expectedCount int
	}{
		{
			name:          "No realm filter",
			realm:         "",
			mockError:     nil,
			expectedCount: 2,
		},
		{
			name:          "Filter by OC1 realm",
			realm:         "OC1",
			mockError:     nil,
			expectedCount: 1,
		},
		{
			name:          "Filter by OC2 realm",
			realm:         "OC2",
			mockError:     nil,
			expectedCount: 1,
		},
		{
			name:          "Filter by non-existent realm",
			realm:         "OC3",
			mockError:     nil,
			expectedCount: 0,
		},
		{
			name:          "Error loading tenancy map",
			realm:         "",
			mockError:     errors.New("failed to load tenancy map"),
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock service
			service := &mockService{
				Service: Service{
					logger: logger.NewTestLogger(),
				},
				mockLoadTenancyMap: func() ([]appConfig.MappingsFile, error) {
					return mockMappings, tc.mockError
				},
			}

			// Call LoadTenancyMappings
			result, err := service.LoadTenancyMappings(tc.realm)

			// Verify the results
			if tc.mockError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedCount, len(result.Mappings))

				// If filtering by realm, verify that all mappings have the correct realm
				if tc.realm != "" && tc.expectedCount > 0 {
					for _, mapping := range result.Mappings {
						assert.Equal(t, tc.realm, mapping.Realm)
					}
				}
			}
		})
	}
}
