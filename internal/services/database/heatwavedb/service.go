package heatwavedb

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/domain/database"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service provides operations and functionalities related to HeatWave database management, logging, and compartment handling.
type Service struct {
	repo          database.HeatWaveDatabaseRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance with the provided application context.
func NewService(repo database.HeatWaveDatabaseRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// ListHeatWaveDb retrieves and returns all HeatWave databases from the given compartment in the OCI account.
func (s *Service) ListHeatWaveDb(ctx context.Context) ([]HeatWaveDatabase, error) {
	s.logger.V(logger.Debug).Info("listing HeatWave databases")
	databases, err := s.repo.ListHeatWaveDatabases(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list HeatWave databases: %w", err)
	}
	return databases, nil
}

// FetchPaginatedHeatWaveDb retrieves a paginated list of HeatWave databases with given limit and page number parameters.
// It returns the slice of databases, total count, next page token, and an error if encountered.
func (s *Service) FetchPaginatedHeatWaveDb(ctx context.Context, limit, pageNum int) ([]HeatWaveDatabase, int, string, error) {
	s.logger.V(logger.Debug).Info("listing HeatWave databases", "limit", limit, "pageNum", pageNum)

	allDatabases, err := s.repo.ListEnrichedHeatWaveDatabases(ctx, s.compartmentID)
	if err != nil {
		allDatabases, err = s.repo.ListHeatWaveDatabases(ctx, s.compartmentID)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed to list HeatWave databases: %w", err)
		}
	}

	if len(allDatabases) == 0 {
		var baseErr error
		allDatabases, baseErr = s.repo.ListHeatWaveDatabases(ctx, s.compartmentID)
		if baseErr != nil {
			return nil, 0, "", fmt.Errorf("failed to list HeatWave databases: %w", baseErr)
		}
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allDatabases, limit, pageNum)

	logger.LogWithLevel(s.logger, logger.Info, "completed HeatWave database listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}
