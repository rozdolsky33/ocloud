package networklb

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/networklb"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
)

// Service provides operations for managing network load balancers.
type Service struct {
	repo          domain.NetworkLoadBalancerRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates a new network load balancer service.
func NewService(repo domain.NetworkLoadBalancerRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// GetNetworkLoadBalancer retrieves a network load balancer by its OCID.
func (s *Service) GetNetworkLoadBalancer(ctx context.Context, ocid string) (*NetworkLoadBalancer, error) {
	s.logger.V(logger.Debug).Info("getting network load balancer", "ocid", ocid)
	nlb, err := s.repo.GetNetworkLoadBalancer(ctx, ocid)
	if err != nil {
		return nil, fmt.Errorf("failed to get network load balancer: %w", err)
	}
	return nlb, nil
}

// ListNetworkLoadBalancers lists all network load balancers in the configured compartment.
func (s *Service) ListNetworkLoadBalancers(ctx context.Context) ([]NetworkLoadBalancer, error) {
	s.logger.V(logger.Debug).Info("listing network load balancers", "compartmentID", s.compartmentID)
	nlbs, err := s.repo.ListNetworkLoadBalancers(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list network load balancers: %w", err)
	}
	return nlbs, nil
}

// FetchPaginatedNetworkLoadBalancers returns a page of network load balancers and pagination metadata.
func (s *Service) FetchPaginatedNetworkLoadBalancers(ctx context.Context, limit, pageNum int, showAll bool) ([]NetworkLoadBalancer, int, string, error) {
	s.logger.V(logger.Debug).Info("fetching paginated network load balancers", "limit", limit, "page", pageNum, "showAll", showAll)
	var (
		all []NetworkLoadBalancer
		err error
	)
	if showAll {
		all, err = s.repo.ListEnrichedNetworkLoadBalancers(ctx, s.compartmentID)
	} else {
		all, err = s.repo.ListNetworkLoadBalancers(ctx, s.compartmentID)
	}
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing network load balancers from repository: %w", err)
	}
	total := len(all)
	if pageNum <= 0 {
		pageNum = 1
	}
	start := (pageNum - 1) * limit
	end := start + limit
	if start >= total {
		return []NetworkLoadBalancer{}, total, "", nil
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

// GetEnrichedNetworkLoadBalancer retrieves and returns the enriched network load balancer by OCID.
func (s *Service) GetEnrichedNetworkLoadBalancer(ctx context.Context, ocid string) (*NetworkLoadBalancer, error) {
	s.logger.V(logger.Debug).Info("getting enriched network load balancer", "ocid", ocid)
	nlb, err := s.repo.GetEnrichedNetworkLoadBalancer(ctx, ocid)
	if err != nil {
		return nil, fmt.Errorf("failed to get enriched network load balancer: %w", err)
	}
	return nlb, nil
}

// FuzzySearch performs a fuzzy search for network load balancers based on the provided search pattern.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]NetworkLoadBalancer, error) {
	all, err := s.repo.ListEnrichedNetworkLoadBalancers(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all network load balancers for search: %w", err)
	}

	indexables := ToSearchableNetworkLoadBalancers(all)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}

	results := make([]NetworkLoadBalancer, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}
	return results, nil
}
