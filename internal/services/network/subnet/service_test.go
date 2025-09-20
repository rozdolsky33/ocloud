package subnet

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain"
	"github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/stretchr/testify/assert"
)

// mockSubnetRepository is a mock implementation of the SubnetRepository for testing.
type mockSubnetRepository struct {
	subnets []vcn.Subnet
	err     error
}

func (m *mockSubnetRepository) GetSubnet(ctx context.Context, ocid string) (*vcn.Subnet, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, s := range m.subnets {
		if s.OCID == ocid {
			return &s, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSubnetRepository) ListSubnets(ctx context.Context, compartmentID string) ([]vcn.Subnet, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.subnets, nil
}

func TestService_Find(t *testing.T) {
	mockRepo := &mockSubnetRepository{
		subnets: []vcn.Subnet{
			{DisplayName: "Test Subnet", CIDRBlock: "10.0.0.0/24"},
			{DisplayName: "Another Subnet", CIDRBlock: "10.0.1.0/24"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, err := service.Find(context.Background(), "test")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Test Subnet", results[0].DisplayName)
}

func TestService_List(t *testing.T) {
	mockRepo := &mockSubnetRepository{
		subnets: []vcn.Subnet{
			{DisplayName: "Test Subnet"},
			{DisplayName: "Another Subnet"},
		},
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	results, _, _, err := service.List(context.Background(), 10, 1)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestService_List_Error(t *testing.T) {
	expectedErr := errors.New("some error")
	mockRepo := &mockSubnetRepository{
		err: expectedErr,
	}
	service := NewService(mockRepo, logr.Discard(), "test-compartment")

	_, _, _, err := service.List(context.Background(), 10, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
}
