package subnet

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	nc, err := oci.NewNetworkClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}
	return &Service{
		networkClient: nc,
		logger:        appCtx.Logger,
		compartmentID: appCtx.CompartmentID,
	}, nil
}

func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]Subnet, int, string, error) {

	var subnets []Subnet
	var nextPageToken string
	var totalCount int

	// Prepare the base request
	// Create a request with limited parameters to fetch only the required page
	request := core.ListSubnetsRequest{
		CompartmentId: &s.compartmentID,
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

		for currentPage < pageNum {
			// Fetch page token, not actual data
			// Use limit to ensure consistent pagination
			tokenRequest := core.ListSubnetsRequest{
				CompartmentId: &s.compartmentID,
				Page:          &page,
			}

			if limit > 0 {
				tokenRequest.Limit = &limit
			}

			resp, err := s.networkClient.ListSubnets(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
			}

			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []Subnet{}, 0, "", nil
			}
			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}
		// Set the page token for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 1, "Using page token for page", "pageNum", pageNum, "token", page)
	}

	// Fetch Subnets for the request
	resp, err := s.networkClient.ListSubnets(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing subnets: %w", err)
	}
	// Set the total count to the number of subnets returned
	// If we have a next page, this is an estimate
	totalCount = len(resp.Items)
	// If we have a next page, we know there are more subnets
	if resp.OpcNextPage != nil {
		// Estimate total count based on my current page and items per rage
		totalCount = pageNum*limit + limit
	}

	//Save the next page token if available
	if resp.OpcNextPage != nil {
		nextPageToken = *resp.OpcNextPage
		logger.LogWithLevel(s.logger, 3, "Next page token", "token", nextPageToken)
	}

	//Process the subnets
	for _, oc := range resp.Items {
		subnets = append(subnets, mapToSubnets(oc))

	}
	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount
	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(subnets),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return subnets, totalCount, nextPageToken, nil
}

func (s *Service) Find(ctx context.Context, namePattern string) ([]Subnet, error) {
	return nil, nil
}

func mapToSubnets(s core.Subnet) Subnet {
	return Subnet{
		Name:                    *s.DisplayName,
		ID:                      *s.Id,
		CIDR:                    *s.CidrBlock,
		VcnID:                   *s.VcnId,
		RouteTableID:            *s.RouteTableId,
		SecurityListID:          s.SecurityListIds,
		DhcpOptionsID:           *s.DhcpOptionsId,
		ProhibitPublicIPOnVnic:  *s.ProhibitPublicIpOnVnic,
		ProhibitInternetIngress: *s.ProhibitInternetIngress,
		ProhibitInternetEgress:  *s.ProhibitInternetIngress,
		DNSLabel:                *s.DnsLabel,
		SubnetDomainName:        *s.SubnetDomainName,
	}
}
