package loadbalancer

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/search"
)

// Service provides operations for managing load balancers.
type Service struct {
	repo          domain.LoadBalancerRepository
	logger        logr.Logger
	compartmentID string
}

// NewService creates a new load balancer service.
func NewService(repo domain.LoadBalancerRepository, appCtx *app.ApplicationContext) *Service {
	return &Service{
		repo:          repo,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}
}

// GetLoadBalancer retrieves a load balancer by its OCID.
func (s *Service) GetLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	s.logger.V(logger.Debug).Info("getting load balancer", "ocid", ocid)
	lb, err := s.repo.GetLoadBalancer(ctx, ocid)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	return lb, nil
}

// ListLoadBalancers lists all load balancers in the configured compartment.
func (s *Service) ListLoadBalancers(ctx context.Context) ([]LoadBalancer, error) {
	s.logger.V(logger.Debug).Info("listing load balancers", "compartmentID", s.compartmentID)
	lbs, err := s.repo.ListLoadBalancers(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list load balancers: %w", err)
	}
	return lbs, nil
}

// FetchPaginatedLoadBalancers returns a page of load balancers and pagination metadata.
// If showAll is true, it uses the enriched model; otherwise, it uses the basic model for performance.
func (s *Service) FetchPaginatedLoadBalancers(ctx context.Context, limit, pageNum int, showAll bool) ([]LoadBalancer, int, string, error) {
	s.logger.V(logger.Debug).Info("fetching paginated load balancers", "limit", limit, "page", pageNum, "showAll", showAll)
	var (
		all []LoadBalancer
		err error
	)
	if showAll {
		all, err = s.repo.ListEnrichedLoadBalancers(ctx, s.compartmentID)
	} else {
		all, err = s.repo.ListLoadBalancers(ctx, s.compartmentID)
	}
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing load balancers from repository: %w", err)
	}
	total := len(all)
	if pageNum <= 0 {
		pageNum = 1
	}
	start := (pageNum - 1) * limit
	end := start + limit
	if start >= total {
		return []LoadBalancer{}, total, "", nil
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

// GetEnrichedLoadBalancer retrieves and returns the enriched load balancer by OCID.
func (s *Service) GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	s.logger.V(logger.Debug).Info("getting enriched load balancer", "ocid", ocid)
	lb, err := s.repo.GetEnrichedLoadBalancer(ctx, ocid)
	if err != nil {
		return nil, fmt.Errorf("failed to get enriched load balancer: %w", err)
	}
	return lb, nil
}

// FuzzySearch performs a fuzzy search for load balancers based on the provided search pattern.
func (s *Service) FuzzySearch(ctx context.Context, searchPattern string) ([]LoadBalancer, error) {
	all, err := s.repo.ListEnrichedLoadBalancers(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all load balancers for search: %w", err)
	}

	// Build the search index using the common search package and the load balancer searcher adapter.
	indexables := ToSearchableLoadBalancers(all)
	idxMapping := search.NewIndexMapping(GetSearchableFields())
	idx, err := search.BuildIndex(indexables, idxMapping)
	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	matchedIdxs, err := search.FuzzySearch(idx, searchPattern, GetSearchableFields(), GetBoostedFields())
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}

	results := make([]LoadBalancer, 0, len(matchedIdxs))
	for _, i := range matchedIdxs {
		if i >= 0 && i < len(all) {
			results = append(results, all[i])
		}
	}
	return results, nil
}
