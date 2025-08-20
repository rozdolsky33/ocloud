package policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// NewService initializes a new Service instance with the provided application context.
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	return &Service{
		identityClient: appCtx.IdentityClient,
		logger:         appCtx.Logger,
		CompartmentID:  appCtx.CompartmentID,
	}, nil
}

// List retrieves a paginated list of policies based on the provided limit and page number parameters.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Policy, int, string, error) {
	logger.LogWithLevel(s.logger, 1, "Listing Policies", "limit", limit, "page", pageNum)

	var policies []Policy
	var nextPageToken string
	var totalCount int

	// Prepare the base request
	request := identity.ListPoliciesRequest{
		CompartmentId: &s.CompartmentID,
	}

	// Add limit parameters if specified
	if limit > 0 {
		request.Limit = &limit
		logger.LogWithLevel(s.logger, 1, "Limiting policies to", "limit", limit)
	}
	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 1, "Calculating page token for page", "pageNum", pageNum)

		// We need to fetch page tokens until we reach the desired page
		page := ""
		currentPage := 1

		for currentPage < pageNum {
			// Fetch Just the page token, not actual data
			// Usu the same limit to ensure consistent pagination
			tokenRequest := identity.ListPoliciesRequest{
				CompartmentId: &s.CompartmentID,
				Page:          &page,
			}

			if limit > 0 {
				tokenRequest.Limit = &limit
			}

			resp, err := s.identityClient.ListPolicies(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
			}

			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []Policy{}, 0, "", nil
			}
			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}
		// Set the page token for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 1, "Using page token for page", "pageNum", pageNum, "token", page)
	}

	// Fetch Policies for the request
	resp, err := s.identityClient.ListPolicies(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing policies: %w", err)
	}
	// Set the total count to the number of policies returned
	// If we have a next page, this is an estimate
	totalCount = len(resp.Items)
	// If we have a next page, we know there are more instances
	if resp.OpcNextPage != nil {
		// Estimate total count based on current page and items per rage
		totalCount = pageNum*limit + limit
	}

	// Save the next page token if available
	if resp.OpcNextPage != nil {
		nextPageToken = *resp.OpcNextPage
		logger.LogWithLevel(s.logger, 1, "Next page token", "token", nextPageToken)
	}

	// Process the policies
	for _, p := range resp.Items {
		policies = append(policies, mapToPolicies(p))
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(policies),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return policies, totalCount, nextPageToken, nil
}

// Find performs a fuzzy search for policies based on the provided searchPattern and returns matching policy.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Policy, error) {
	logger.LogWithLevel(s.logger, 1, "Finding Policies", "pattern", searchPattern)
	var allPolicies []Policy

	// 1. Fetch all policies in the compartment
	allPolicies, err := s.fetchAllPolicies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all policies: %w", err)
	}

	// 2. Build index
	index, err := util.BuildIndex(allPolicies, func(p Policy) any {
		return mapToIndexablePolicy(p)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build index: %w", err)
	}

	// 3. Fuzzy search on multiple files
	fields := []string{"Name", "Description", "Statement", "Tags", "TagValues"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("failed to fuzzy search index: %w", err)
	}

	// Return matched policies
	var matchedPolicies []Policy
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allPolicies) {
			matchedPolicies = append(matchedPolicies, allPolicies[idx])
		}
	}

	logger.LogWithLevel(s.logger, 2, "Found policies", "count", len(matchedPolicies))
	return matchedPolicies, nil
}

// fetchAllPolicies retrieves all policies within a specific compartment using pagination. Returns a slice of Policy or an error.
func (s *Service) fetchAllPolicies(ctx context.Context) ([]Policy, error) {
	var allPolicies []Policy
	page := ""
	for {
		resp, err := s.identityClient.ListPolicies(ctx, identity.ListPoliciesRequest{
			CompartmentId: &s.CompartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list policies: %w", err)
		}
		for _, p := range resp.Items {
			allPolicies = append(allPolicies, mapToPolicies(p))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return allPolicies, nil
}

// mapToPolicies converts an identity.Policy object to an shared Policy representation, mapping all fields correspondingly.
func mapToPolicies(policy identity.Policy) Policy {
	return Policy{
		Name:        *policy.Name,
		ID:          *policy.Id,
		Statement:   policy.Statements,
		Description: *policy.Description,
		PolicyTags: util.ResourceTags{
			FreeformTags: policy.FreeformTags,
			DefinedTags:  policy.DefinedTags,
		},
	}
}

// mapToIndexablePolicy transforms a Policy object into an IndexablePolicy object with indexed and searchable fields.
func mapToIndexablePolicy(p Policy) IndexablePolicy {
	flattenedTags, _ := util.FlattenTags(p.PolicyTags.FreeformTags, p.PolicyTags.DefinedTags)
	tagValues, _ := util.ExtractTagValues(p.PolicyTags.FreeformTags, p.PolicyTags.DefinedTags)
	joinedStatements := strings.Join(p.Statement, " ") // Concatenate all the statements into one string for fuzzy search/indexing.
	return IndexablePolicy{
		Name:        p.Name,
		Description: p.Description,
		Statement:   joinedStatements,
		Tags:        flattenedTags,
		TagValues:   tagValues,
	}

}
