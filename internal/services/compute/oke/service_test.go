package oke

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/stretchr/testify/assert"
)

// mockClusterRepository is a mock implementation of the ClusterRepository for testing.
type mockClusterRepository struct {
	clusters []compute.Cluster
	err      error
}

func (m *mockClusterRepository) GetCluster(ctx context.Context, ocid string) (*Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockClusterRepository) ListClusters(ctx context.Context, compartmentID string) ([]compute.Cluster, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.clusters, nil
}

// GetClusters is required to satisfy the domain.ClusterRepository interface.
func (m *mockClusterRepository) GetClusters(ctx context.Context, ocid string) ([]compute.Cluster, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.clusters, nil
}

func TestService_Find(t *testing.T) {
	mockRepo := &mockClusterRepository{
		clusters: []compute.Cluster{
			{DisplayName: "Test Cluster", KubernetesVersion: "v1.25.4"},
			{DisplayName: "Another Cluster", KubernetesVersion: "v1.24.8"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, err := service.Find(context.Background(), "test")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Test Cluster", results[0].DisplayName)
}

func TestService_List(t *testing.T) {
	mockRepo := &mockClusterRepository{
		clusters: []compute.Cluster{
			{DisplayName: "Test Cluster"},
			{DisplayName: "Another Cluster"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, _, _, err := service.FetchPaginatedClusters(context.Background(), 10, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestService_List_Error(t *testing.T) {
	expectedErr := errors.New("some error")
	mockRepo := &mockClusterRepository{
		err: expectedErr,
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	_, _, _, err := service.FetchPaginatedClusters(context.Background(), 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
