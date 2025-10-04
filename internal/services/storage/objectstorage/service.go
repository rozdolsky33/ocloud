package objectstorage

import (
	"github.com/go-logr/logr"
	storage "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
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
