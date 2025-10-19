package cacheclusterdb

import (
	"context"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCacheClusterRepository is a mock implementation of domain.CacheClusterRepository
type MockCacheClusterRepository struct {
	mock.Mock
}

func (m *MockCacheClusterRepository) GetCacheCluster(ctx context.Context, clusterId string) (*database.CacheCluster, error) {
	args := m.Called(ctx, clusterId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.CacheCluster), args.Error(1)
}

func (m *MockCacheClusterRepository) ListCacheClusters(ctx context.Context, compartmentID string) ([]database.CacheCluster, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.CacheCluster), args.Error(1)
}

func (m *MockCacheClusterRepository) ListEnrichedCacheClusters(ctx context.Context, compartmentID string) ([]database.CacheCluster, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]database.CacheCluster), args.Error(1)
}

// TestServiceStruct tests the basic structure of the Service struct
func TestServiceStruct(t *testing.T) {
	// Create a simple service
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
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1.iad.test",
		Logger:          logger.NewTestLogger(),
	}

	service := NewService(mockRepo, appCtx)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, appCtx.Logger, service.logger)
	assert.Equal(t, appCtx.CompartmentID, service.compartmentID)
}

// TestListCacheClusters tests the ListCacheClusters function
func TestListCacheClusters(t *testing.T) {
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedClusters := []database.CacheCluster{
		{DisplayName: "cluster-1", ID: "ocid1.rediscluster.oc1.iad.aaa"},
		{DisplayName: "cluster-2", ID: "ocid1.rediscluster.oc1.iad.bbb"},
	}

	mockRepo.On("ListCacheClusters", ctx, appCtx.CompartmentID).Return(expectedClusters, nil).Once()

	clusters, err := service.ListCacheClusters(ctx)

	assert.NoError(t, err)
	assert.Len(t, clusters, 2)
	assert.Equal(t, "cluster-1", clusters[0].DisplayName)
	assert.Equal(t, "cluster-2", clusters[1].DisplayName)

	mockRepo.AssertExpectations(t)
}

// TestFetchPaginatedCacheClusters tests the paginated listing function
func TestFetchPaginatedCacheClusters(t *testing.T) {
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedClusters := []database.CacheCluster{
		{DisplayName: "cluster-1", ID: "ocid1.rediscluster.oc1.iad.aaa"},
		{DisplayName: "cluster-2", ID: "ocid1.rediscluster.oc1.iad.bbb"},
		{DisplayName: "cluster-3", ID: "ocid1.rediscluster.oc1.iad.ccc"},
	}

	mockRepo.On("ListEnrichedCacheClusters", ctx, appCtx.CompartmentID).Return(expectedClusters, nil).Once()

	clusters, totalCount, nextPageToken, err := service.FetchPaginatedCacheClusters(ctx, 2, 1)

	assert.NoError(t, err)
	assert.Len(t, clusters, 2)
	assert.Equal(t, 3, totalCount)
	assert.NotEmpty(t, nextPageToken)
	assert.Equal(t, "cluster-1", clusters[0].DisplayName)
	assert.Equal(t, "cluster-2", clusters[1].DisplayName)

	mockRepo.AssertExpectations(t)
}

// TestFetchPaginatedCacheClusters_SecondPage tests pagination on the second page
func TestFetchPaginatedCacheClusters_SecondPage(t *testing.T) {
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedClusters := []database.CacheCluster{
		{DisplayName: "cluster-1", ID: "ocid1.rediscluster.oc1.iad.aaa"},
		{DisplayName: "cluster-2", ID: "ocid1.rediscluster.oc1.iad.bbb"},
		{DisplayName: "cluster-3", ID: "ocid1.rediscluster.oc1.iad.ccc"},
		{DisplayName: "cluster-4", ID: "ocid1.rediscluster.oc1.iad.ddd"},
	}

	mockRepo.On("ListEnrichedCacheClusters", ctx, appCtx.CompartmentID).Return(expectedClusters, nil).Once()

	clusters, totalCount, nextPageToken, err := service.FetchPaginatedCacheClusters(ctx, 2, 2)

	assert.NoError(t, err)
	assert.Len(t, clusters, 2)
	assert.Equal(t, 4, totalCount)
	assert.Empty(t, nextPageToken)
	assert.Equal(t, "cluster-3", clusters[0].DisplayName)
	assert.Equal(t, "cluster-4", clusters[1].DisplayName)

	mockRepo.AssertExpectations(t)
}

// TestFetchPaginatedCacheClusters_EmptyResults tests listing when no clusters are found
func TestFetchPaginatedCacheClusters_EmptyResults(t *testing.T) {
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	emptyClusters := []database.CacheCluster{}

	mockRepo.On("ListEnrichedCacheClusters", ctx, appCtx.CompartmentID).Return(emptyClusters, nil).Once()
	mockRepo.On("ListCacheClusters", ctx, appCtx.CompartmentID).Return(emptyClusters, nil).Once()

	clusters, totalCount, nextPageToken, err := service.FetchPaginatedCacheClusters(ctx, 10, 1)

	assert.NoError(t, err)
	assert.Len(t, clusters, 0)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, nextPageToken)

	mockRepo.AssertExpectations(t)
}

// TestFuzzySearch tests the FuzzySearch function with various patterns
func TestFuzzySearch(t *testing.T) {
	mockRepo := new(MockCacheClusterRepository)
	appCtx := &app.ApplicationContext{
		CompartmentID: "test-compartment-id",
		Logger:        logger.NewTestLogger(),
	}
	service := NewService(mockRepo, appCtx)
	ctx := context.Background()

	expectedClusters := []database.CacheCluster{
		{
			DisplayName:              "prod-cache-cluster",
			ID:                       "ocid1.rediscluster.oc1.iad.prod",
			LifecycleState:           "ACTIVE",
			SoftwareVersion:          "VALKEY_7_2",
			ClusterMode:              "NONSHARDED",
			NodeCount:                3,
			NodeMemoryInGBs:          16,
			VcnName:                  "prod-vcn",
			SubnetName:               "cache-subnet",
			PrimaryFqdn:              "prod-cluster.redis.us-ashburn-1.oci.oraclecloud.com",
			PrimaryEndpointIpAddress: "10.0.20.100",
			FreeformTags:             map[string]string{"env": "production", "tier": "cache"},
		},
		{
			DisplayName:              "dev-redis-cluster",
			ID:                       "ocid1.rediscluster.oc1.iad.dev",
			LifecycleState:           "ACTIVE",
			SoftwareVersion:          "REDIS_7_0",
			ClusterMode:              "SHARDED",
			ShardCount:               2,
			NodeCount:                6,
			NodeMemoryInGBs:          8,
			VcnName:                  "dev-vcn",
			SubnetName:               "dev-subnet",
			PrimaryFqdn:              "dev-cluster.redis.us-ashburn-1.oci.oraclecloud.com",
			PrimaryEndpointIpAddress: "10.0.10.50",
			FreeformTags:             map[string]string{"env": "development"},
		},
		{
			DisplayName:              "test-valkey-cluster",
			ID:                       "ocid1.rediscluster.oc1.iad.test",
			LifecycleState:           "INACTIVE",
			SoftwareVersion:          "VALKEY_7_2",
			ClusterMode:              "NONSHARDED",
			NodeCount:                1,
			NodeMemoryInGBs:          8,
			VcnName:                  "test-vcn",
			SubnetName:               "test-subnet",
			PrimaryFqdn:              "test-cluster.redis.us-ashburn-1.oci.oraclecloud.com",
			PrimaryEndpointIpAddress: "10.0.30.25",
		},
	}

	mockRepo.On("ListEnrichedCacheClusters", ctx, appCtx.CompartmentID).Return(expectedClusters, nil)

	// Test search by name
	t.Run("search by name", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "prod")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.DisplayName == "prod-cache-cluster" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find prod-cache-cluster")
	})

	// Test search by software version
	t.Run("search by software version", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "VALKEY_7_2")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2, "should find at least 2 clusters with VALKEY_7_2")
	})

	// Test search by cluster mode
	t.Run("search by cluster mode", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "SHARDED")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.ClusterMode == "SHARDED" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find sharded cluster")
	})

	// Test search by VCN name
	t.Run("search by VCN name", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "dev-vcn")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.VcnName == "dev-vcn" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find cluster in dev-vcn")
	})

	// Test search by endpoint FQDN
	t.Run("search by endpoint FQDN", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "prod-cluster.redis")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.DisplayName == "prod-cache-cluster" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find cluster by FQDN")
	})

	// Test search by IP address
	t.Run("search by IP address", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "10.0.20")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.PrimaryEndpointIpAddress == "10.0.20.100" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find cluster with IP 10.0.20.100")
	})

	// Test search by tag value
	t.Run("search by tag value", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "production")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
		found := false
		for _, cluster := range results {
			if cluster.DisplayName == "prod-cache-cluster" {
				found = true
				break
			}
		}
		assert.True(t, found, "should find cluster with production tag")
	})

	// Test search by state
	t.Run("search by state", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "ACTIVE")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2, "should find active clusters")
	})

	// Test empty search pattern returns all
	t.Run("empty search returns all", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "")
		assert.NoError(t, err)
		assert.Len(t, results, 3, "empty search should return all clusters")
	})

	// Test whitespace-only search pattern returns all
	t.Run("whitespace search returns all", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "   ")
		assert.NoError(t, err)
		assert.Len(t, results, 3, "whitespace-only search should return all clusters")
	})

	// Test search with no matches
	t.Run("search with no matches", func(t *testing.T) {
		results, err := service.FuzzySearch(ctx, "zzz-completely-nonexistent-xyz-999")
		assert.NoError(t, err)
		assert.Len(t, results, 0, "should return empty results for non-matching pattern")
	})

	mockRepo.AssertExpectations(t)
}
