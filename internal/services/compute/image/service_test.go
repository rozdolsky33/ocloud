package image

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/stretchr/testify/assert"
)

// mockImageRepository is a mock implementation of the ImageRepository for testing.
type mockImageRepository struct {
	images []compute.Image
	err    error
}

func (m *mockImageRepository) GetImage(ctx context.Context, ocid string) (*compute.Image, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, i := range m.images {
		if i.OCID == ocid {
			return &i, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockImageRepository) ListImages(ctx context.Context, compartmentID string) ([]compute.Image, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.images, nil
}

func TestService_Find(t *testing.T) {
	mockRepo := &mockImageRepository{
		images: []compute.Image{
			{DisplayName: "Test Image", OperatingSystem: "Linux"},
			{DisplayName: "Another Image", OperatingSystem: "Windows"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, err := service.FuzzySearch(context.Background(), "test")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Test Image", results[0].DisplayName)
}

func TestService_Get(t *testing.T) {
	mockRepo := &mockImageRepository{
		images: []compute.Image{
			{DisplayName: "Test Image"},
			{DisplayName: "Another Image"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, _, _, err := service.FetchPaginatedImages(context.Background(), 10, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestService_Get_Error(t *testing.T) {
	expectedErr := errors.New("some error")
	mockRepo := &mockImageRepository{
		err: expectedErr,
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	_, _, _, err := service.FetchPaginatedImages(context.Background(), 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
