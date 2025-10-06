package image

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for image operations.
type Service struct {
	imageRepo     compute.ImageRepository
	logger        logr.Logger
	compartmentID string
}

// NewService initializes a new Service instance.
func NewService(repo compute.ImageRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		imageRepo:     repo,
		logger:        logger,
		compartmentID: compartmentID,
	}
}

// FetchPaginatedImages retrieves a paginated list of images.
func (s *Service) FetchPaginatedImages(ctx context.Context, limit, pageNum int) ([]Image, int, string, error) {
	s.logger.V(logger.Debug).Info("listing images", "limit", limit, "pageNum", pageNum)

	allImages, err := s.imageRepo.ListImages(ctx, s.compartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing images from repository: %w", err)
	}

	pagedResults, totalCount, nextPageToken := util.PaginateSlice(allImages, limit, pageNum)

	s.logger.Info("completed image listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// FuzzySearch performs a fuzzy search for images.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]Image, error) {
	s.logger.V(logger.Debug).Info("finding images with fuzzy search", "pattern", searchPattern)

	allImages, err := s.imageRepo.ListImages(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all images for search: %w", err)
	}

	searchableImages := ToSearchableImages(allImages)
	indexMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(searchableImages, indexMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	s.logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allImages))

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]Image, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(allImages) {
			results = append(results, allImages[i])
		}
	}

	return results, nil
}
