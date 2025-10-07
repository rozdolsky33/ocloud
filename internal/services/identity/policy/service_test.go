package policy

import (
	"context"
	"fmt"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// TestServiceStruct tests the basic structure of the Service struct
func TestServiceStruct(t *testing.T) {
	// Create a simple service with nil clients
	service := &Service{
		logger:        logger.NewTestLogger(),
		CompartmentID: "test-compartment-id",
	}

	// Test that the service was created correctly
	assert.NotNil(t, service)
	assert.Equal(t, "test-compartment-id", service.CompartmentID)
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
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
	}

	//service, err := NewService(appCtx.Logger, appCtx.CompartmentID)

	// but if we did, we would expect no error and a valid service
	//assert.NoError(t, err)
	//assert.NotNil(t, service)
	fmt.Println(appCtx)
}

// TestList tests the FetchPaginatedClusters function
func TestList(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for FetchPaginatedClusters since it requires the OCI SDK")

	// This is a placeholder test that would normally test the FetchPaginatedClusters function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call FetchPaginatedClusters with different parameters
	// 3. Verify that the returned policies, total count, and next page token are correct

	service := &Service{
		logger:        logger.NewTestLogger(),
		CompartmentID: "test-compartment-id",
	}

	ctx := context.Background()
	policies, _, _, err := service.FetchPaginatedPolies(ctx, 10, 1)

	// but if we did, we would expect no error and a valid list of policies
	assert.NoError(t, err)
	assert.NotNil(t, policies)
}

// TestFind tests the FuzzySearch function
func TestSearch(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for FuzzySearch since it requires the OCI SDK")

	// This is a placeholder test that would normally test the FuzzySearch function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call FuzzySearch with different search patterns
	// 3. Verify that the returned policies match the search pattern

	service := &Service{
		logger:        logger.NewTestLogger(),
		CompartmentID: "test-compartment-id",
	}

	ctx := context.Background()
	policies, err := service.FuzzySearch(ctx, "test")

	// but if we did, we would expect no error and a valid list of policies
	assert.NoError(t, err)
	assert.NotNil(t, policies)
}

// TestFetchAllPolicies tests the fetchAllPolicies function
func TestSearchAllPolicies(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for fetchAllPolicies since it requires the OCI SDK")

	// This is a placeholder test that would normally test the fetchAllPolicies function
	// In a real test, we would:
	// 1. Create a mock service with mock clients
	// 2. Call fetchAllPolicies
	// 3. Verify that all policies are returned

	service := &Service{
		logger:        logger.NewTestLogger(),
		CompartmentID: "test-compartment-id",
	}

	ctx := context.Background()
	policies, err := service.policyRepo.ListPolicies(ctx, "test.ocid.292393")

	// but if we did, we would expect no error and a valid list of policies
	assert.NoError(t, err)
	assert.NotNil(t, policies)
}

// TestMapToIndexablePolicy verifies the SearchablePolicy adapter produces expected fields
func TestMapToIndexablePolicy(t *testing.T) {
	// Create a test policy
	p := identity.Policy{
		Name:        "TestPolicy",
		ID:          "ocid1.policy.oc1.phx.test",
		Description: "Test policy description",
		Statement:   []string{"Allow group Administrators to manage all-resources in tenancy"},
	}

	// Adapt and convert to an indexable map
	sp := SearchablePolicy{Policy: Policy(p)}
	doc := sp.ToIndexable()

	// Validate expected lower-cased fields and combined statements
	assert.Equal(t, "testpolicy", doc["Name"])                     // lowercased
	assert.Equal(t, "test policy description", doc["Description"]) // lowercased
	assert.Equal(t, "ocid1.policy.oc1.phx.test", doc["OCID"])      // lowercased (already)
	assert.Contains(t, doc["Statements"], "allow group administrators")

	// Tags are not set; flattened forms should exist as empty strings
	assert.Equal(t, "", doc["TagsKV"])
	assert.Equal(t, "", doc["TagsVal"])
}

// TestMapToPolicies tests the mapToPolicies function
func TestMapToPolicies(t *testing.T) {
	// Skip this test since it requires the OCI SDK
	t.Skip("Skipping test for mapToPolicies since it requires the OCI SDK")

	// This is a placeholder test that would normally test the mapToPolicies function
	// In a real test, we would:
	// 1. Create a mock OCI policy
	// 2. Call mapToPolicies with the mock policy
	// 3. Verify that the returned Policy has the expected values
}
