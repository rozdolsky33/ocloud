package compartment

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
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
	logger.LogWithLevel(s.logger, 3, "finding compartments with bleve fuzzy search", "pattern", searchPattern)

	var allCompartments []Compartment
	var indexableDocs []IndexableCompartment
	page := ""

	// 1. Fetch all Compartments
	for {
		resp, err := s.identityClient.ListCompartments(ctx, identity.ListCompartmentsRequest{
			CompartmentId: &s.TenancyID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing compartments: %w", err)
		}
		for _, comp := range resp.Items {
			compartment := mapToCompartment(comp)
			allCompartments = append(allCompartments, compartment)
			indexableDocs = append(indexableDocs, mapToIndexableCompartment(compartment))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	// 2. Create an in-memory Bleve index
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create index for compartments: %w", err)
	}
	for i, doc := range indexableDocs {
		err := index.Index(fmt.Sprintf("%d", i), doc)
		if err != nil {
			return nil, fmt.Errorf("indexing compartments failed: %w", err)
		}
	}

	// 3. Prepare a fuzzy query with wildcard
	searchPattern = strings.ToLower(searchPattern)
	if !strings.HasPrefix(searchPattern, "*") {
		searchPattern = "*" + searchPattern
	}
	if !strings.HasPrefix(searchPattern, "*") {
		searchPattern = searchPattern + "*"
	}

	// Create a query that searches across all relevant fields
	// The _all field is a special field that searches across all indexed fields
	// We also explicitly search in Tags and TagValues fields to ensure tag searches work correctly
	queryString := fmt.Sprintf("_all:%s OR Tags:%s OR TagValues:%s",
		searchPattern, searchPattern, searchPattern)

	query := bleve.NewQueryStringQuery(queryString)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 1000 // Increase from default of 10
	// 4. Perform search
	result, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// 5. Collect matched results
	var matched []Compartment
	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil || idx < 0 || idx >= len(allCompartments) {
			continue
		}
		matched = append(matched, allCompartments[idx])
	}

	logger.LogWithLevel(s.logger, 2, "found image", "count", len(matched))

	return matched, nil
}

func mapToCompartment(compartment identity.Compartment) Compartment {
	return Compartment{
		Name:        *compartment.Name,
		ID:          *compartment.Id,
		Description: *compartment.Description,
	}
}

func mapToIndexableCompartment(compartment Compartment) IndexableCompartment {
	flattenedTags, _ := util.FlattenTags(compartment.CompartmentTags.FreeformTags, compartment.CompartmentTags.DefinedTags)
	tagValues, _ := util.ExtractTagValues(compartment.CompartmentTags.FreeformTags, compartment.CompartmentTags.DefinedTags)
	return IndexableCompartment{
		ID:          compartment.ID,
		Name:        compartment.Name,
		Description: compartment.Description,
		Tags:        flattenedTags,
		TagValues:   tagValues,
	}
}
