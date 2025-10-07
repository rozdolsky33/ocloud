package compartment

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/stretchr/testify/assert"
)

// mockCompartmentRepository is a mock implementation of the CompartmentRepository for testing.
type mockCompartmentRepository struct {
	compartments []identity.Compartment
	err          error
}

func (m *mockCompartmentRepository) GetCompartment(ctx context.Context, ocid string) (*identity.Compartment, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, c := range m.compartments {
		if c.OCID == ocid {
			return &c, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockCompartmentRepository) ListCompartments(ctx context.Context, parentCompartmentID string) ([]identity.Compartment, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.compartments, nil
}

func TestService_Find(t *testing.T) {
	mockRepo := &mockCompartmentRepository{
		compartments: []identity.Compartment{
			{DisplayName: "Test Compartment", Description: "A test compartment"},
			{DisplayName: "Another", Description: "Another one"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	results, err := service.FuzzySearch(context.Background(), "test")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Test Compartment", results[0].DisplayName)
}

func TestService_List(t *testing.T) {
	mockRepo := &mockCompartmentRepository{
		compartments: []identity.Compartment{
			{DisplayName: "Test Compartment"},
			{DisplayName: "Another"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	results, _, _, err := service.FetchPaginateCompartments(context.Background(), 10, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestService_List_Error(t *testing.T) {
	expectedErr := errors.New("some error")
	mockRepo := &mockCompartmentRepository{
		err: expectedErr,
	}
	service := NewService(mockRepo, logr.Discard(), "test-tenancy")

	_, _, _, err := service.FetchPaginateCompartments(context.Background(), 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
