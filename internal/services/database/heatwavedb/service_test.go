package heatwavedb

import (
	"context"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHeatWaveDatabaseRepository is a mock implementation of domain.HeatWaveDatabaseRepository
type MockHeatWaveDatabaseRepository struct {
	mock.Mock
}

func (m *MockHeatWaveDatabaseRepository) ListEnrichedHeatWaveDatabases(ctx context.Context, compartmentID string) ([]database.HeatWaveDatabase, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.HeatWaveDatabase), args.Error(1)
}

func (m *MockHeatWaveDatabaseRepository) GetHeatWaveDatabase(ctx context.Context, ocid string) (*database.HeatWaveDatabase, error) {
	args := m.Called(ctx, ocid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.HeatWaveDatabase), args.Error(1)
}

func (m *MockHeatWaveDatabaseRepository) ListHeatWaveDatabases(ctx context.Context, compartmentID string) ([]database.HeatWaveDatabase, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.HeatWaveDatabase), args.Error(1)
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
	mockRepo := new(MockHeatWaveDatabaseRepository)
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

// TestList tests the FetchPaginatedHeatWaveDb function
func TestList(t *testing.T) {
	mockRepo := new(MockHeatWaveDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedDBs := []database.HeatWaveDatabase{
		{DisplayName: "hw-db1", ID: "ocid1.mysqldbsystem.oc1..aaaaaaaan"},
		{DisplayName: "hw-db2", ID: "ocid1.mysqldbsystem.oc1..bbbbbbbbn"},
		{DisplayName: "hw-db3", ID: "ocid1.mysqldbsystem.oc1..ccccccccn"},
	}

	mockRepo.On("ListEnrichedHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()

	databases, totalCount, nextPageToken, err := service.FetchPaginatedHeatWaveDb(ctx, 2, 1)

	assert.NoError(t, err)
	assert.Len(t, databases, 2)
	assert.Equal(t, 3, totalCount)
	assert.NotEmpty(t, nextPageToken)
	assert.Equal(t, "hw-db1", databases[0].DisplayName)
	assert.Equal(t, "hw-db2", databases[1].DisplayName)

	mockRepo.AssertExpectations(t)
}

// TestListEmptyResults tests listing when no databases are found
func TestListEmptyResults(t *testing.T) {
	mockRepo := new(MockHeatWaveDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	emptyDBs := []database.HeatWaveDatabase{}

	mockRepo.On("ListEnrichedHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(emptyDBs, nil).Once()
	mockRepo.On("ListHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(emptyDBs, nil).Once()

	databases, totalCount, nextPageToken, err := service.FetchPaginatedHeatWaveDb(ctx, 10, 1)

	assert.NoError(t, err)
	assert.Len(t, databases, 0)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, nextPageToken)

	mockRepo.AssertExpectations(t)
}

// TestListHeatWaveDb tests the ListHeatWaveDb function
func TestListHeatWaveDb(t *testing.T) {
	mockRepo := new(MockHeatWaveDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedDBs := []database.HeatWaveDatabase{
		{DisplayName: "hw-db1", ID: "ocid1.mysqldbsystem.oc1..aaaaaaaan"},
		{DisplayName: "hw-db2", ID: "ocid1.mysqldbsystem.oc1..bbbbbbbbn"},
	}

	mockRepo.On("ListHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()

	databases, err := service.ListHeatWaveDb(ctx)

	assert.NoError(t, err)
	assert.Len(t, databases, 2)
	assert.Equal(t, "hw-db1", databases[0].DisplayName)
	assert.Equal(t, "hw-db2", databases[1].DisplayName)

	mockRepo.AssertExpectations(t)
}

// TestFetchPaginatedHeatWaveDb_SecondPage tests pagination on the second page
func TestFetchPaginatedHeatWaveDb_SecondPage(t *testing.T) {
	mockRepo := new(MockHeatWaveDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedDBs := []database.HeatWaveDatabase{
		{DisplayName: "hw-db1", ID: "ocid1.mysqldbsystem.oc1..aaaaaaaan"},
		{DisplayName: "hw-db2", ID: "ocid1.mysqldbsystem.oc1..bbbbbbbbn"},
		{DisplayName: "hw-db3", ID: "ocid1.mysqldbsystem.oc1..ccccccccn"},
		{DisplayName: "hw-db4", ID: "ocid1.mysqldbsystem.oc1..ddddddddn"},
	}

	mockRepo.On("ListEnrichedHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(expectedDBs, nil).Once()

	databases, totalCount, nextPageToken, err := service.FetchPaginatedHeatWaveDb(ctx, 2, 2)

	assert.NoError(t, err)
	assert.Len(t, databases, 2)
	assert.Equal(t, 4, totalCount)
	assert.Empty(t, nextPageToken)
	assert.Equal(t, "hw-db3", databases[0].DisplayName)
	assert.Equal(t, "hw-db4", databases[1].DisplayName)

	mockRepo.AssertExpectations(t)
}
