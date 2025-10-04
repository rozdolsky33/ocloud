package objectstorage

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	storage "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/logger"
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

// FetchPaginatedBuckets lists buckets and returns a page plus pagination metadata.
// When showAll is true, it enriches each bucket by calling GetBucket to retrieve full details.
func (s *Service) FetchPaginatedBuckets(ctx context.Context, limit, pageNum int) ([]storage.Bucket, int, string, error) {
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
		full, e := s.osRepo.GetBucket(ctx, name)
		if e != nil {
			continue
		}
		all[i] = *full
	}

	total := len(all)
	if pageNum <= 0 {
		pageNum = 1
	}
	start := (pageNum - 1) * limit
	end := start + limit
	if start >= total {
		return []storage.Bucket{}, total, "", nil
	}
	if end > total {
		end = total
	}
	paged := all[start:end]
	next := ""
	if end < total {
		next = fmt.Sprintf("%d", pageNum+1)
	}
	return paged, total, next, nil
}
