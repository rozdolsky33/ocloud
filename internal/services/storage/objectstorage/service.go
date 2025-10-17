package objectstorage

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	storage "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

type Service struct {
	osRepo        storage.ObjectStorageRepository
	logger        logr.Logger
	CompartmentID string
}

func NewService(repo storage.ObjectStorageRepository, logger logr.Logger, compartmentID string) *Service {
	return &Service{
		osRepo:        repo,
		logger:        logger,
		CompartmentID: compartmentID,
	}
}

func (s *Service) ListBuckets(ctx context.Context) ([]Bucket, error) {
	s.logger.V(logger.Debug).Info("listing object storage buckets")
	buckets, err := s.osRepo.ListBuckets(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("listing buckets from repository: %w", err)
	}
	for i := range buckets {
		name := buckets[i].Name
		if name == "" {
			continue
		}
		full, e := s.osRepo.GetBucketByName(ctx, s.CompartmentID, name)
		if e != nil {
			continue
		}
		buckets[i] = *full
	}

	return buckets, nil
}

// FetchPaginatedBuckets lists buckets and returns a page plus pagination metadata.
func (s *Service) FetchPaginatedBuckets(ctx context.Context, limit, pageNum int) ([]Bucket, int, string, error) {
	s.logger.V(logger.Debug).Info("listing object storage buckets", "limit", limit, "page", pageNum)
	all, err := s.osRepo.ListBuckets(ctx, s.CompartmentID)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing buckets from repository: %w", err)
	}
	for i := range all {
		name := all[i].Name
		if name == "" {
			continue
		}
		full, e := s.osRepo.GetBucketByName(ctx, s.CompartmentID, name)
		if e != nil {
			continue
		}
		all[i] = *full
	}
	paged, total, next := util.PaginateSlice(all, limit, pageNum)
	return paged, total, next, nil
}

func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]Bucket, error) {
	s.logger.V(logger.Debug).Info("searching object storage buckets", "pattern", searchPattern)
	// List and enrich buckets similar to ListBuckets behavior
	all, err := s.osRepo.ListBuckets(ctx, s.CompartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all buckets for search: %w", err)
	}
	for i := range all {
		name := all[i].Name
		if name == "" {
			continue
		}
		if full, e := s.osRepo.GetBucketByName(ctx, s.CompartmentID, name); e == nil && full != nil {
			all[i] = *full
		}
	}

	// Build the search index using the common search package and the bucket searcher adapter.
	indexables := ToSearchableBuckets(all)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}

	results := make([]Bucket, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}
	return results, nil
}
