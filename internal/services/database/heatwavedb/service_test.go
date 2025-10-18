package heatwavedb

import (
	"context"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/mysql"
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

// TestFuzzySearch tests the FuzzySearch function with various patterns
func TestFuzzySearch(t *testing.T) {
	mockRepo := new(MockHeatWaveDatabaseRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	storage := 1024
	clusterSize := 3

	expectedDBs := []database.HeatWaveDatabase{
		{
			DisplayName:    "prod-mysql-db",
			ID:             "ocid1.mysqldbsystem.oc1..prod",
			LifecycleState: "ACTIVE",
			MysqlVersion:   "8.4.6",
			ShapeName:      "MySQL.4",
			VcnName:        "prod-vcn",
			SubnetName:     "database-subnet",
			IpAddress:      "10.0.20.175",
			HeatWaveCluster: &mysql.HeatWaveClusterSummary{
				ClusterSize: &clusterSize,
			},
			DataStorage: &mysql.DataStorage{
				DataStorageSizeInGBs: &storage,
			},
			FreeformTags: map[string]string{"env": "production"},
		},
		{
			DisplayName:    "dev-mysql-db",
			ID:             "ocid1.mysqldbsystem.oc1..dev",
			LifecycleState: "INACTIVE",
			MysqlVersion:   "8.0.35",
			ShapeName:      "MySQL.2",
			VcnName:        "dev-vcn",
			SubnetName:     "dev-subnet",
			IpAddress:      "10.0.10.100",
			FreeformTags:   map[string]string{"env": "development"},
		},
		{
			DisplayName:    "test-heatwave",
			ID:             "ocid1.mysqldbsystem.oc1..test",
			LifecycleState: "ACTIVE",
			MysqlVersion:   "8.4.6",
			ShapeName:      "MySQL.8",
			VcnName:        "test-vcn",
			SubnetName:     "test-subnet",
			IpAddress:      "10.0.30.50",
		},
	}

	mockRepo.On("ListEnrichedHeatWaveDatabases", ctx, appCtx.CompartmentID).Return(expectedDBs, nil)

	// Test search by name
	t.Run("search by name", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "prod")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		// Should find prod-mysql-db
		found := false
		for _, db := range results {
			if db.DisplayName == "prod-mysql-db" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find prod-mysql-db")
	})

	// Test search by version
	t.Run("search by version", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "8.4.6")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2, "should find at least 2 databases with version 8.4.6")
	})

	// Test search by VCN name
	t.Run("search by VCN name", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "dev-vcn")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, db := range results {
			if db.VcnName == "dev-vcn" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find database in dev-vcn")
	})

	// Test search by IP address
	t.Run("search by IP address", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "10.0.20")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, db := range results {
			if db.IpAddress == "10.0.20.175" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find database with IP 10.0.20.175")
	})

	// Test search by shape
	t.Run("search by shape", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "MySQL.4")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, db := range results {
			if db.ShapeName == "MySQL.4" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find MySQL.4 database")
	})

	// Test search by tag value
	t.Run("search by tag value", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "production")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, db := range results {
			if db.DisplayName == "prod-mysql-db" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find database with production tag")
	})

	// Test empty search pattern returns all
	t.Run("empty search returns all", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "")
		assert.NoError(t, err)
		assert.Len(t, results, 3, "empty search should return all databases")
	})

	// Test whitespace-only search pattern returns all
	t.Run("whitespace search returns all", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "   ")
		assert.NoError(t, err)
		assert.Len(t, results, 3, "whitespace-only search should return all databases")
	})

	// Test search with no matches
	t.Run("search with no matches", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "zzz-completely-nonexistent-xyz-999")
		assert.NoError(t, err)
		assert.Len(t, results, 0, "should return empty results for non-matching pattern")
	})

	mockRepo.AssertExpectations(t)
}
