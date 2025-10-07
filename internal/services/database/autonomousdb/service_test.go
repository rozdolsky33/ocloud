package autonomousdb

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAutonomousDatabaseRepository is a mock implementation of domain.AutonomousDatabaseRepository
type MockAutonomousDatabaseRepository struct {
	mock.Mock
}

func (m *MockAutonomousDatabaseRepository) ListEnrichedAutonomousDatabase(ctx context.Context, compartmentID string) ([]database.AutonomousDatabase, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.AutonomousDatabase), args.Error(1)
}

func (m *MockAutonomousDatabaseRepository) GetAutonomousDatabase(ctx context.Context, ocid string) (*database.AutonomousDatabase, error) {
	args := m.Called(ctx, ocid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.AutonomousDatabase), args.Error(1)
}

// ListAutonomousDatabases mocks the ListAutonomousDatabases method of domain.AutonomousDatabaseRepository
func (m *MockAutonomousDatabaseRepository) ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]database.AutonomousDatabase, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.AutonomousDatabase), args.Error(1)
}

// FindAutonomousDatabase mocks the FindAutonomousDatabase method of domain.AutonomousDatabaseRepository
func (m *MockAutonomousDatabaseRepository) FindAutonomousDatabase(ctx context.Context, compartmentID, name string) (*database.AutonomousDatabase, error) {
	args := m.Called(ctx, compartmentID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.AutonomousDatabase), args.Error(1)
}

// TestServiceStruct tests the basic structure of the Service struct
func TestServiceStruct(t *testing.T) {
	// Create a simple service with nil clients
	service := &Service{
		logger:        logger.NewTestLogger(),
		compartmentID: "test-compartment-id",
	}

	// Test that the service was created correctly
	assert.NotNil(t, service)
	assert.Equal(t, "test-compartment-id", service.compartmentID)
}

// TestNewService tests the NewService function
func TestNewService(t *testing.T) {
	mockRepo := new(MockAutonomousDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.phx.test",
		Logger:          logger.NewTestLogger(),
	}

	service := NewService(mockRepo, appCtx)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, appCtx.Logger, service.logger)
	assert.Equal(t, appCtx.CompartmentID, service.compartmentID)
}

// TestList tests the FetchPaginatedClusters function
func TestList(t *testing.T) {
	mockRepo := new(MockAutonomousDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedDBs := []database.AutonomousDatabase{
		{Name: "db1", ID: "ocid1.autonomousdatabase.oc1..aaaaaaaan"},
		{Name: "db2", ID: "ocid1.autonomousdatabase.oc1..bbbbbbbbn"},
		{Name: "db3", ID: "ocid1.autonomousdatabase.oc1..ccccccccn"},
	}

	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()

	databases, totalCount, nextPageToken, err := service.FetchPaginatedAutonomousDb(ctx, 2, 1)

	assert.NoError(t, err)
	assert.Len(t, databases, 2)
	assert.Equal(t, expectedDBs[0].Name, databases[0].Name)
	assert.Equal(t, expectedDBs[1].Name, databases[1].Name)
	assert.Equal(t, len(expectedDBs), totalCount)
	assert.Equal(t, "2", nextPageToken) // Because there's a next page

	mockRepo.AssertExpectations(t)

	// Test second page
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()
	databases, totalCount, nextPageToken, err = service.FetchPaginatedAutonomousDb(ctx, 2, 2)

	assert.NoError(t, err)
	assert.Len(t, databases, 1)
	assert.Equal(t, expectedDBs[2].Name, databases[0].Name)
	assert.Equal(t, len(expectedDBs), totalCount)
	assert.Equal(t, "", nextPageToken)

	mockRepo.AssertExpectations(t)

	// Test empty result
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return([]database.AutonomousDatabase{}, nil).Once()
	mockRepo.On("ListAutonomousDatabases", ctx, appCtx.CompartmentID).Return([]database.AutonomousDatabase{}, nil).Once()
	databases, totalCount, nextPageToken, err = service.FetchPaginatedAutonomousDb(ctx, 10, 1)

	assert.NoError(t, err)
	assert.Len(t, databases, 0)
	assert.Equal(t, 0, totalCount)
	assert.Equal(t, "", nextPageToken)

	mockRepo.AssertExpectations(t)

	// Test error case
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return([]database.AutonomousDatabase{}, fmt.Errorf("mock error")).Once()
	mockRepo.On("ListAutonomousDatabases", ctx, appCtx.CompartmentID).Return([]database.AutonomousDatabase{}, fmt.Errorf("mock error")).Once()
	databases, totalCount, nextPageToken, err = service.FetchPaginatedAutonomousDb(ctx, 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error")
	assert.Nil(t, databases)
	assert.Equal(t, 0, totalCount)
	assert.Equal(t, "", nextPageToken)

	mockRepo.AssertExpectations(t)
}

// TestFind tests the FuzzySearch function
func TestFind(t *testing.T) {
	mockRepo := new(MockAutonomousDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedDBs := []database.AutonomousDatabase{
		{Name: "prod-db", ID: "ocid1.autonomousdatabase.oc1..aaaaaaaan"},
		{Name: "test-db", ID: "ocid1.autonomousdatabase.oc1..bbbbbbbbn"},
		{Name: "dev-db", ID: "ocid1.autonomousdatabase.oc1..ccccccccn"},
	}

	// Test case: found
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()
	databases, err := service.FuzzySearch(ctx, "test")

	assert.NoError(t, err)
	assert.Len(t, databases, 1)
	assert.Equal(t, "test-db", databases[0].Name)
	mockRepo.AssertExpectations(t)

	// Test case: not found
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()
	databases, err = service.FuzzySearch(ctx, "nonexistent")

	assert.NoError(t, err) // Fuzzy search returns no error if not found, just an empty list
	assert.Len(t, databases, 0)
	mockRepo.AssertExpectations(t)

	// Test case: error from repository
	mockRepo.On("ListEnrichedAutonomousDatabase", ctx, appCtx.CompartmentID).Return([]database.AutonomousDatabase{}, fmt.Errorf("mock error")).Once()
	databases, err = service.FuzzySearch(ctx, "test")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch all databases: mock error")
	assert.Nil(t, databases)
	mockRepo.AssertExpectations(t)
}

// TestMapToIndexableDatabase tests the SearchableAutonomousDatabase.ToIndexable mapping
func TestMapToIndexableDatabase(t *testing.T) {
	// Create a test database
	db := database.AutonomousDatabase{
		Name: "TestDatabase",
		ID:   "ocid1.autonomousdatabase.oc1.phx.test",
	}

	// Adapt to searchable and build indexable map
	s := SearchableAutonomousDatabase{AutonomousDatabase: db}
	idx := s.ToIndexable()

	// Verify that the indexable database has the expected (lowercased) values
	assert.Equal(t, strings.ToLower(db.Name), idx["Name"])
	assert.Equal(t, strings.ToLower(db.ID), idx["OCID"])
}
