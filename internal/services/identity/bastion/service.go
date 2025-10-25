package bastion

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	ocibastion "github.com/rozdolsky33/ocloud/internal/oci/identity/bastion"
)

// NewService creates a new bastion service with a repository pattern.
// This is the new constructor that uses the repository interface.
func NewService(repo domain.BastionRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		bastionRepo:   repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// NewServiceFromAppContext creates a bastion service from ApplicationContext.
// This is a convenience constructor that maintains backward compatibility with existing code
// while using the repository pattern internally. For new code, prefer using NewService()
// with explicit adapter creation (see compartment/get.go for a reference pattern).
func NewServiceFromAppContext(appCtx *app.ApplicationContext) (*Service, error) {
	bc, err := oci.NewBastionClient(appCtx.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create bastion client: %w", err)
	}
	nc, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}
	cc, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}

	// Create the adapter
	adapter := ocibastion.NewBastionAdapter(bc, nc, appCtx.CompartmentID)

	// Create service with repository
	service := &Service{
		bastionRepo:   adapter,
		bastionClient: bc, // Temporarily kept for session management
		networkClient: nc,
		computeClient: cc,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}

	return service, nil
}

// List retrieves and returns all bastion hosts from the given compartment.
// The enrichment with VCN and Subnet names is now handled by the repository adapter.
func (s *Service) List(ctx context.Context) ([]Bastion, error) {
	logger.LogWithLevel(s.logger, logger.Debug, "Listing Bastions in compartment", "compartmentID", s.compartmentID)

	bastions, err := s.bastionRepo.ListBastions(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bastions: %w", err)
	}

	logger.Logger.V(logger.Debug).Info("Successfully listed bastions.", "count", len(bastions))
	return bastions, nil
}

// Get retrieves a specific bastion by ID.
func (s *Service) Get(ctx context.Context, bastionID string) (*Bastion, error) {
	logger.LogWithLevel(s.logger, logger.Debug, "Getting Bastion", "bastionID", bastionID)

	bastion, err := s.bastionRepo.GetBastion(ctx, bastionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bastion: %w", err)
	}

	logger.Logger.V(logger.Debug).Info("Successfully retrieved bastion.", "bastionID", bastionID)
	return bastion, nil
}
