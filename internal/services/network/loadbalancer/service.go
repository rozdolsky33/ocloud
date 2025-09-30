package loadbalancer

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
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

func (s *Service) Find(ctx context.Context, searchPattern string) ([]LoadBalancer, error) {

	// 1: Fetch all load balancers
	all, err := s.repo.ListEnrichedLoadBalancers(ctx, s.compartmentID)
	if err != nil {
		return nil, fmt.Errorf("fetching all load balancers for search: %w", err)
	}

	// 2: Filter by search pattern
	index, err := util.BuildIndex(all, func(lb LoadBalancer) any {
		return mapToIndexableLoadBalancer(lb)
	})

	if err != nil {
		return nil, fmt.Errorf("building search index: %w", err)
	}

	fields := []string{"Name"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("performing fuzzy search: %w", err)
	}
	var results []LoadBalancer
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(all) {
			results = append(results, all[idx])
		}
	}

	return results, nil
}

func mapToIndexableLoadBalancer(lb LoadBalancer) any {
	// Normalize strings to lowercase to match lowercased search patterns
	lower := func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToLower(s)
	}
	lowerSlice := func(in []string) []string {
		if len(in) == 0 {
			return nil
		}
		out := make([]string, 0, len(in))
		for _, v := range in {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			out = append(out, strings.ToLower(v))
		}
		if len(out) == 0 {
			return nil
		}
		return out
	}

	return struct {
		Name            string
		Type            string
		VcnName         string
		Hostnames       []string
		SSLCertificates []string
		Subnets         []string
	}{
		Name:            lower(lb.Name),
		Type:            lower(lb.Type),
		VcnName:         lower(lb.VcnName),
		Hostnames:       lowerSlice(lb.Hostnames),
		SSLCertificates: lowerSlice(lb.SSLCertificates),
		Subnets:         lowerSlice(lb.Subnets),
	}
}
