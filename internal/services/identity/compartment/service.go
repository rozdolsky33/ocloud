package compartment

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"strings"
)

// NewService initializes and returns a new Service instance using the provided ApplicationContext.
// It injects required dependencies such as IdentityClient, Logger, TenancyID, and TenancyName into the Service.
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	return &Service{
		identityClient: appCtx.IdentityClient,
		logger:         appCtx.Logger,
		TenancyID:      appCtx.TenancyID,
		TenancyName:    appCtx.TenancyName,
	}, nil
}

// List retrieves a paginated list of compartments based on the provided limit and page number parameters.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Compartment, int, string, error) {
	var compartments []Compartment
	var nextPageToken string
	var totalCount int
	logger.LogWithLevel(s.logger, 3, "List compartments", "limit", limit, "pageNum", pageNum, "Total Count", totalCount)

	// Prepare the base request
	// to Create a request with a limit parameter to fetch only the required page
	request := identity.ListCompartmentsRequest{
		CompartmentId:          &s.TenancyID,
		AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
		LifecycleState:         identity.CompartmentLifecycleStateActive,
		CompartmentIdInSubtree: common.Bool(true), // ListCompartments on the tenancy (root compartment) default is false
	}

	// Add limit parameters
	if limit > 0 {
		request.Limit = &limit
		logger.LogWithLevel(s.logger, 3, "Setting limit parameter", "limit", limit)
	}

	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

		// paginate through results; stop when OpcNextPage is nil
		page := ""
		currentPage := 1

		for currentPage <= pageNum {
			// Fetch page token, not actual data
			// Use limit to ensure consistent patination
			tokenRequest := identity.ListCompartmentsRequest{
				CompartmentId:          &s.TenancyID,
				AccessLevel:            identity.ListCompartmentsAccessLevelAccessible,
				LifecycleState:         identity.CompartmentLifecycleStateActive,
				CompartmentIdInSubtree: common.Bool(true), // ListCompartments on the tenancy (root compartment) default is false
				Page:                   &page,
			}

			if limit > 0 {
				tokenRequest.Limit = &limit
			}

			resp, err := s.identityClient.ListCompartments(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("error fetching token: %w", err)
			}

			// If there is no next page, we've reached then end
			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []Compartment{}, 0, "", nil
			}
			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}
		// Set the page token for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)

	}

	// Fetch compartments for the request
	response, err := s.identityClient.ListCompartments(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing compartments: %w", err)
	}

	// Set the total count to the number of compartments returned
	// If we have a next page, this is an estimate
	totalCount = len(response.Items)
	// if we have a next page, we know there is more
	if response.OpcNextPage != nil {
		// If we have a next page token, we know there are more compartments
		// We need to estimate the total count more accurately
		// Since we don't know the exact total count, we'll set it to a value
		// that indicates there are more pages (at least one more page worth of compartments)
		totalCount = pageNum*limit + limit
	}

	// Save the next page toke if available
	if response.OpcNextPage != nil {
		nextPageToken = *response.OpcNextPage
		logger.LogWithLevel(s.logger, 3, "Next page token", "token", nextPageToken, "Total Count", totalCount)

	}

	// Process the compartment
	for _, comp := range response.Items {
		compartments = append(compartments, mapToCompartment(comp))

	}
	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(compartments),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return compartments, totalCount, nextPageToken, nil

}

// Find performs a fuzzy search for compartments based on the provided searchPattern and returns matching compartments.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Compartment, error) {
	logger.LogWithLevel(s.logger, 3, "Finding allCompartments using Bleve fuzzy search", "pattern", searchPattern)

	// Step 1: Fetch all allCompartments
	allCompartments, err := s.fetchAllCompartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all compartments: %w", err)
	}

	// Step 2: Build index
	index, err := util.BuildIndex(allCompartments, func(c Compartment) any {
		return mapToIndexableCompartment(c)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build index: %w", err)
	}

	// Step 3: Fuzzy search on multiple fields
	fields := []string{"Name", "Description"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("failed to fuzzy search index: %w", err)
	}

	// Step 4: Return matched allCompartments
	var results []Compartment
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allCompartments) {
			results = append(results, allCompartments[idx])
		}
	}

	logger.LogWithLevel(s.logger, 2, "Compartment search complete", "matches", len(results))
	return results, nil
}

// fetchAllCompartments retrieves all active compartments within a tenancy, including nested compartments.
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

// mapToCompartment maps an identity.Compartment to a Compartment struct, transferring selected field values.
func mapToCompartment(compartment identity.Compartment) Compartment {
	return Compartment{
		Name:        *compartment.Name,
		ID:          *compartment.Id,
		Description: *compartment.Description,
		CompartmentTags: util.ResourceTags{
			FreeformTags: compartment.FreeformTags,
			DefinedTags:  compartment.DefinedTags,
		},
	}
}

// mapToIndexableCompartment converts a Compartment instance to an IndexableCompartment with lowercased fields.
func mapToIndexableCompartment(compartment Compartment) IndexableCompartment {
	return IndexableCompartment{
		Name:        strings.ToLower(compartment.Name),
		Description: strings.ToLower(compartment.Description),
	}
}
