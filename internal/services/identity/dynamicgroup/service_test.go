package dynamicgroup

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/stretchr/testify/assert"
)

// mockDynamicGroupRepository is a mock implementation of the DynamicGroupRepository for testing.
type mockDynamicGroupRepository struct {
	dynamicGroups []identity.DynamicGroup
	err           error
}

func (m *mockDynamicGroupRepository) GetDynamicGroup(ctx context.Context, ocid string) (*identity.DynamicGroup, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, dg := range m.dynamicGroups {
		if dg.OCID == ocid {
			return &dg, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockDynamicGroupRepository) ListDynamicGroups(ctx context.Context, compartmentID string) ([]identity.DynamicGroup, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.dynamicGroups, nil
}

func TestService_FuzzySearch(t *testing.T) {
	mockRepo := &mockDynamicGroupRepository{
		dynamicGroups: []identity.DynamicGroup{
			{Name: "AppServer-DG", Description: "Application Servers"},
			{Name: "DBServer-DG", Description: "Database Servers"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	results, err := service.FuzzySearch(context.Background(), "app")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "AppServer-DG", results[0].Name)
}

func TestService_FetchPaginateDynamicGroups(t *testing.T) {
	mockRepo := &mockDynamicGroupRepository{
		dynamicGroups: []identity.DynamicGroup{
			{Name: "DG1"},
			{Name: "DG2"},
			{Name: "DG3"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	// Page 1, Limit 2
	results, total, next, err := service.FetchPaginateDynamicGroups(context.Background(), 2, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, 3, total)
	assert.NotEmpty(t, next)

	// Page 2, Limit 2
	results, total, next, err = service.FetchPaginateDynamicGroups(context.Background(), 2, 2)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Empty(t, next)
}

func TestService_List_Error(t *testing.T) {
	expectedErr := errors.New("failed to list")
	mockRepo := &mockDynamicGroupRepository{
		err: expectedErr,
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	_, _, _, err := service.FetchPaginateDynamicGroups(context.Background(), 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
