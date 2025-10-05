package objectstorage

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	storage "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/logger"
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
		full, e := s.osRepo.GetBucketByName(ctx, name)
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
		full, e := s.osRepo.GetBucketByName(ctx, name)
		if e != nil {
			continue
		}
		all[i] = *full
	}
	paged, total, next := util.PaginateSlice(all, limit, pageNum)
	return paged, total, next, nil
}
