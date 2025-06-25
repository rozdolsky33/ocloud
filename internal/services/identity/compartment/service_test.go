package compartment

import (
	"context"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestServiceStruct tests the basic structure of the Service struct
func TestServiceStruct(t *testing.T) {
	// Create a simple service with nil clients
	service := &Service{
		logger:      logger.NewTestLogger(),
		TenancyID:   "test-tenancy-id",
		TenancyName: "test-tenancy-name",
	}

	// Test that the service was created correctly
	assert.NotNil(t, service)
	assert.Equal(t, "test-tenancy-id", service.TenancyID)
	assert.Equal(t, "test-tenancy-name", service.TenancyName)
}

// TestNewService tests the NewService function
func TestNewService(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for NewService since it requires the OCI SDK")

	// This is a placeholder test that would normally test the NewService function
	// In a real test, we would:
	// 1. Create a mock application context
	// 2. Call NewService with the mock context
	// 3. Verify that the returned service has the expected values

	appCtx := &app.ApplicationContext{
		TenancyName: "TestTenancy",
		TenancyID:   "ocid1.tenancy.oc1.phx.test",
		Logger:      logger.NewTestLogger(),
	}

	service, err := NewService(appCtx)

	// but if we did, we would expect no error and a valid service
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

// TestList tests the List function
func TestList(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for List since it requires the OCI SDK")

	// This is a placeholder test that would normally test the List function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call List with different parameters
	// 3. Verify that the returned compartments, total count, and next page token are correct

	service := &Service{
		logger:      logger.NewTestLogger(),
		TenancyID:   "test-tenancy-id",
		TenancyName: "test-tenancy-name",
	}

	ctx := context.Background()
	compartments, _, _, err := service.List(ctx, 10, 1)

	// but if we did, we would expect no error and a valid list of compartments
	assert.NoError(t, err)
	assert.NotNil(t, compartments)
}

// TestFind tests the Find function
func TestFind(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for Find since it requires the OCI SDK")

	// This is a placeholder test that would normally test the Find function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call Find with different search patterns
	// 3. Verify that the returned compartments match the search pattern

	service := &Service{
		logger:      logger.NewTestLogger(),
		TenancyID:   "test-tenancy-id",
		TenancyName: "test-tenancy-name",
	}

	ctx := context.Background()
	compartments, err := service.Find(ctx, "test")

	// but if we did, we would expect no error and a valid list of compartments
	assert.NoError(t, err)
	assert.NotNil(t, compartments)
}

// TestFetchAllCompartments tests the fetchAllCompartments function
func TestFetchAllCompartments(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for fetchAllCompartments since it requires the OCI SDK")

	// This is a placeholder test that would normally test the fetchAllCompartments function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call fetchAllCompartments
	// 3. Verify that all compartments are returned

	service := &Service{
		logger:      logger.NewTestLogger(),
		TenancyID:   "test-tenancy-id",
		TenancyName: "test-tenancy-name",
	}

	ctx := context.Background()
	compartments, err := service.fetchAllCompartments(ctx)

	// but if we did, we would expect no error and a valid list of compartments
	assert.NoError(t, err)
	assert.NotNil(t, compartments)
}

// TestMapToCompartment tests the mapToCompartment function
func TestMapToCompartment(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for mapToCompartment since it requires the OCI SDK")

	// This is a placeholder test that would normally test the mapToCompartment function
	// In a real test, we would:
	// 1. Create a mock OCI compartment
	// 2. Call mapToCompartment with the mock compartment
	// 3. Verify that the returned Compartment has the expected values
}

// TestMapToIndexableCompartment tests the mapToIndexableCompartment function
func TestMapToIndexableCompartment(t *testing.T) {
	// Create a test compartment
	compartment := Compartment{
		Name:        "TestCompartment",
		ID:          "ocid1.compartment.oc1.phx.test",
		Description: "Test compartment description",
	}

	// Call mapToIndexableCompartment
	indexable := mapToIndexableCompartment(compartment)

	// Verify that the indexable compartment has the expected values
	assert.Equal(t, strings.ToLower(compartment.Name), indexable.Name)
	assert.Equal(t, strings.ToLower(compartment.Description), indexable.Description)
}
