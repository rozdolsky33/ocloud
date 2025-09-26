package image

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service is the application-layer service for image operations.
type Service struct {
	imageRepo     compute.ImageRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates a Service that orchestrates image-related operations using the provided image repository, logger, and compartment ID.
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

	totalCount := len(allImages)
	start := (pageNum - 1) * limit
	end := start + limit

	if start >= totalCount {
		return []Image{}, totalCount, "", nil
	}

	if end > totalCount {
		end = totalCount
	}

	pagedResults := allImages[start:end]

	var nextPageToken string
	if end < totalCount {
		nextPageToken = fmt.Sprintf("%d", pageNum+1)
	}

	s.logger.Info("completed image listing", "returnedCount", len(pagedResults), "totalCount", totalCount)
	return pagedResults, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search for images.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Image, error) {
	s.logger.V(logger.Debug).Info("finding images with fuzzy search", "pattern", searchPattern)

	allImages, err := s.imageRepo.ListImages(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all images for search: %w", err)
	}

	index, err := util.BuildIndex(allImages, func(img Image) any {
		return mapToIndexableImage(img)
	})
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}
	s.logger.V(logger.Debug).Info("Search index built successfully.", "numEntries", len(allImages))

	fields := []string{"Name", "OperatingSystem", "OperatingSystemVersion"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	s.logger.V(logger.Debug).Info("Fuzzy search completed.", "numMatches", len(matchedIdxs))
	var results []Image
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allImages) {
			results = append(results, allImages[idx])
		}
	}

	s.logger.Info("image search complete", "matches", len(results))
	return results, nil
}

// mapToIndexableImage converts a compute.Image into an indexable representation.
// The returned value is an anonymous struct with fields `Name`, `OperatingSystem`,
// and `OperatingSystemVersion`, each lowercased for case-insensitive indexing.
func mapToIndexableImage(img compute.Image) any {
	return struct {
		Name                   string
		OperatingSystem        string
		OperatingSystemVersion string
	}{
		Name:                   strings.ToLower(img.DisplayName),
		OperatingSystem:        strings.ToLower(img.OperatingSystem),
		OperatingSystemVersion: strings.ToLower(img.OperatingSystemVersion),
	}
}
