package images

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	cc, err := oci.NewComputeClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}
	return &Service{
		compute:       cc,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}, nil
}

func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Image, int, string, error) {
	return nil, 0, "", nil
}

func (s *Service) Find(ctx context.Context, searchPattern string) ([]Image, error) {
	return nil, nil
}
