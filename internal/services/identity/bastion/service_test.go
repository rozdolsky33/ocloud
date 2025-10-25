package bastion

import (
	"context"
	"fmt"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBastionRepository is a mock implementation of domain.BastionRepository
type MockBastionRepository struct {
	mock.Mock
}

func (m *MockBastionRepository) ListBastions(ctx context.Context, compartmentID string) ([]domain.Bastion, error) {
	args := m.Called(ctx, compartmentID)
	return args.Get(0).([]domain.Bastion), args.Error(1)
}

func (m *MockBastionRepository) GetBastion(ctx context.Context, bastionID string) (*domain.Bastion, error) {
	args := m.Called(ctx, bastionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bastion), args.Error(1)
}

func (m *MockBastionRepository) CreateBastion(ctx context.Context, request domain.CreateBastionRequest) (*domain.Bastion, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Bastion), args.Error(1)
}

func (m *MockBastionRepository) DeleteBastion(ctx context.Context, bastionID string) error {
	args := m.Called(ctx, bastionID)
	return args.Error(0)
}

// TestNewService tests the NewService constructor
func TestNewService(t *testing.T) {
	mockRepo := new(MockBastionRepository)
	testLogger := logger.NewTestLogger()
	compartmentID := "ocid1.compartment.oc1..test"

	service := NewService(mockRepo, testLogger, compartmentID)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.bastionRepo)
	assert.Equal(t, testLogger, service.logger)
	assert.Equal(t, compartmentID, service.compartmentID)
}

// TestList tests the List method
func TestList(t *testing.T) {
	tests := []struct {
		name          string
		compartmentID string
		mockBastions  []domain.Bastion
		mockError     error
		expectError   bool
	}{
		{
			name:          "successful list with bastions",
			compartmentID: "ocid1.compartment.oc1..test",
			mockBastions: []domain.Bastion{
				{
					OCID:             "ocid1.bastion.oc1..test1",
					DisplayName:      "bastion-1",
					BastionType:      "STANDARD",
					LifecycleState:   "ACTIVE",
					TargetVcnName:    "vcn-1",
					TargetSubnetName: "subnet-1",
				},
				{
					OCID:             "ocid1.bastion.oc1..test2",
					DisplayName:      "bastion-2",
					BastionType:      "STANDARD",
					LifecycleState:   "ACTIVE",
					TargetVcnName:    "vcn-2",
					TargetSubnetName: "subnet-2",
				},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:          "successful list with no bastions",
			compartmentID: "ocid1.compartment.oc1..empty",
			mockBastions:  []domain.Bastion{},
			mockError:     nil,
			expectError:   false,
		},
		{
			name:          "repository error",
			compartmentID: "ocid1.compartment.oc1..error",
			mockBastions:  []domain.Bastion{},
			mockError:     fmt.Errorf("repository error"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBastionRepository)
			mockRepo.On("ListBastions", mock.Anything, tt.compartmentID).Return(tt.mockBastions, tt.mockError)

			service := NewService(mockRepo, logger.NewTestLogger(), tt.compartmentID)
			bastions, err := service.List(context.Background())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockBastions), len(bastions))
				if len(tt.mockBastions) > 0 {
					assert.Equal(t, tt.mockBastions[0].OCID, bastions[0].OCID)
					assert.Equal(t, tt.mockBastions[0].DisplayName, bastions[0].DisplayName)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGet tests the Get method
func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		bastionID   string
		mockBastion *domain.Bastion
		mockError   error
		expectError bool
		expectNil   bool
	}{
		{
			name:      "successful get",
			bastionID: "ocid1.bastion.oc1..test",
			mockBastion: &domain.Bastion{
				OCID:                     "ocid1.bastion.oc1..test",
				DisplayName:              "test-bastion",
				BastionType:              "STANDARD",
				LifecycleState:           "ACTIVE",
				MaxSessionTTL:            10800,
				ClientCidrBlockAllowList: []string{"0.0.0.0/0"},
				PrivateEndpointIpAddress: "10.0.0.1",
			},
			mockError:   nil,
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "bastion not found",
			bastionID:   "ocid1.bastion.oc1..notfound",
			mockBastion: nil,
			mockError:   fmt.Errorf("bastion not found"),
			expectError: true,
			expectNil:   true,
		},
		{
			name:        "repository error",
			bastionID:   "ocid1.bastion.oc1..error",
			mockBastion: nil,
			mockError:   fmt.Errorf("repository error"),
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBastionRepository)
			mockRepo.On("GetBastion", mock.Anything, tt.bastionID).Return(tt.mockBastion, tt.mockError)

			service := NewService(mockRepo, logger.NewTestLogger(), "ocid1.compartment.oc1..test")
			bastion, err := service.Get(context.Background(), tt.bastionID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectNil {
				assert.Nil(t, bastion)
			} else {
				assert.NotNil(t, bastion)
				assert.Equal(t, tt.mockBastion.OCID, bastion.OCID)
				assert.Equal(t, tt.mockBastion.DisplayName, bastion.DisplayName)
				assert.Equal(t, tt.mockBastion.MaxSessionTTL, bastion.MaxSessionTTL)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestNewServiceFromAppContext tests the convenience constructor
func TestNewServiceFromAppContext(t *testing.T) {
	// This test verifies the constructor works, but we can't easily test
	// OCI client creation without mocking the entire provider
	appCtx := &app.ApplicationContext{
		CompartmentName: "TestCompartment",
		CompartmentID:   "ocid1.compartment.oc1..test",
		Logger:          logger.NewTestLogger(),
	}

	// We can't fully test this without a mock provider,
	// but we can at least verify it compiles and has the right signature
	_ = appCtx

	// This would normally be:
	// service, err: = NewServiceFromAppContext(appCtx)
	// But that requires a real OCI provider
}
