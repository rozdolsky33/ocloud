package compartment

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	bleveQuery "github.com/blevesearch/bleve/v2/search/query"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"strconv"
	"strings"
)

func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	return &Service{
		identityClient: appCtx.IdentityClient,
		logger:         appCtx.Logger,
		TenancyID:      appCtx.TenancyID,
		TenancyName:    appCtx.TenancyName,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]Compartment, error) {
	logger.LogWithLevel(s.logger, 3, "Listing compartments in tenancy", "tenancyName: ", s.TenancyName, "tenancyID: ", s.TenancyID)

	// prepare the base request
	req := identity.ListCompartmentsRequest{
		CompartmentId:          &s.TenancyID,
		AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
		LifecycleState:         identity.CompartmentLifecycleStateActive,
		CompartmentIdInSubtree: common.Bool(true),
	}

	var compartments []Compartment

	// paginate through results; stop when OpcNextPage is nil
	pageToken := ""
	for {
		if pageToken != "" {
			req.Page = common.String(pageToken)
		}

		resp, err := s.identityClient.ListCompartments(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing compartments: %w", err)
		}

		// scan each compartment summary for a name match
		for _, comp := range resp.Items {
			compartment := mapToCompartment(comp)
			compartments = append(compartments, compartment)

		}

		// if there's no next page, we're done searching
		if resp.OpcNextPage == nil {
			break
		}
		pageToken = *resp.OpcNextPage
	}

	return compartments, nil

}

func (s *Service) Find(ctx context.Context, searchPattern string) ([]Compartment, error) {
	logger.LogWithLevel(s.logger, 3, "Finding compartments using Bleve fuzzy search", "pattern", searchPattern)

	// Step 1: Fetch all compartments
	compartments, err := s.fetchAllCompartments(ctx)
	if err != nil {
		return nil, err
	}

	// Step 2: Build index
	index, err := buildCompartmentIndex(compartments)
	if err != nil {
		return nil, err
	}

	// Step 3: Fuzzy search on multiple fields
	fields := []string{"Name", "Description"}
	matchedIdxs, err := fuzzySearchIndex(index, strings.ToLower(searchPattern), fields, 500)
	if err != nil {
		return nil, err
	}

	// Step 4: Return matched compartments
	var results []Compartment
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(compartments) {
			results = append(results, compartments[idx])
		}
	}

	logger.LogWithLevel(s.logger, 2, "Compartment search complete", "matches", len(results))
	return results, nil
}

func (s *Service) fetchAllCompartments(ctx context.Context) ([]Compartment, error) {
	var all []Compartment
	page := ""

	for {
		resp, err := s.identityClient.ListCompartments(ctx, identity.ListCompartmentsRequest{
			CompartmentId:          &s.TenancyID,
			Page:                   &page,
			AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
			LifecycleState:         identity.CompartmentLifecycleStateActive,
			CompartmentIdInSubtree: common.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("listing compartments: %w", err)
		}
		for _, item := range resp.Items {
			all = append(all, mapToCompartment(item))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return all, nil
}

func buildCompartmentIndex(compartments []Compartment) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	for i, c := range compartments {
		err := index.Index(fmt.Sprintf("%d", i), mapToIndexableCompartment(c))
		if err != nil {
			return nil, fmt.Errorf("indexing failed: %w", err)
		}
	}
	return index, nil
}

func fuzzySearchIndex(index bleve.Index, pattern string, fields []string, limit int) ([]int, error) {
	var queries []bleveQuery.Query

	for _, field := range fields {
		// Fuzzy match (Levenshtein distance)
		fq := bleve.NewFuzzyQuery(pattern)
		fq.SetField(field)
		fq.SetFuzziness(2)
		queries = append(queries, fq)

		// Prefix match (useful for dev1, splunkdev1, etc.)
		pq := bleve.NewPrefixQuery(pattern)
		pq.SetField(field)
		queries = append(queries, pq)

		// Wildcard match (matches anywhere in token)
		wq := bleve.NewWildcardQuery("*" + pattern + "*")
		wq.SetField(field)
		queries = append(queries, wq)
	}

	// OR all queries together
	combinedQuery := bleve.NewDisjunctionQuery(queries...)

	search := bleve.NewSearchRequestOptions(combinedQuery, limit, 0, false)

	result, err := index.Search(search)
	if err != nil {
		return nil, err
	}

	var hits []int
	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil {
			continue
		}
		hits = append(hits, idx)
	}

	return hits, nil
}

func mapToCompartment(compartment identity.Compartment) Compartment {
	return Compartment{
		Name:        *compartment.Name,
		ID:          *compartment.Id,
		Description: *compartment.Description,
	}
}

func mapToIndexableCompartment(compartment Compartment) IndexableCompartment {
	return IndexableCompartment{
		ID:          compartment.ID,
		Name:        strings.ToLower(compartment.Name),
		Description: strings.ToLower(compartment.Description),
	}
}
