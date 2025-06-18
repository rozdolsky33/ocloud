package images

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"strings"
)

// NewService initializes a new Service instance with the provided application context.
// Returns a Service pointer and an error if initialization fails.
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

// List retrieves a paginated list of images with given limit and page number parameters.
// It returns the slice of images, total count, next page token, and an error if encountered.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Image, int, string, error) {
	// Log input parameters at debug level
	logger.LogWithLevel(s.logger, 3, "List() called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	var images []Image
	var nextPageToken string
	var totalCount int

	//Create a request with a limit parameter to fetch only the required page
	request := core.ListImagesRequest{
		CompartmentId: &s.compartmentID,
	}

	// Add limit parameters if specified
	if limit > 0 {
		request.Limit = &limit
		logger.LogWithLevel(s.logger, 3, "Setting limit parameter", "limit", limit)
	}
	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

		// We need to fetch page tokens until we reach the desired page
		page := ""
		currentPage := 1

		for currentPage < pageNum {
			// Fetch just the page token, not actual data
			// Use the same limit to ensure consistent pagination

			tokenRequest := core.ListImagesRequest{
				CompartmentId: &s.compartmentID,
				Page:          &page,
			}
			if limit > 0 {
				tokenRequest.Limit = &limit
			}

			resp, err := s.compute.ListImages(ctx, tokenRequest)
			if err != nil {
				return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
			}

			// If there's no next page, we've reached the end
			if resp.OpcNextPage == nil {
				logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return an empty result since the requested page is beyond available data
				return []Image{}, 0, "", nil
			}
			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}
		// Set the page token for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)
	}

	// Fetch images for the requested page
	resp, err := s.compute.ListImages(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing images: %w", err)
	}
	// Set the total count to the number of instances returned
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
		logger.LogWithLevel(s.logger, 3, "Next page token", "token", nextPageToken)
	}

	// Process the images
	for _, oc := range resp.Items {
		image := mapToImage(oc)
		images = append(images, image)
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(images),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)

	return images, totalCount, nextPageToken, nil
}

// Find searches for images matching the provided search pattern within a compartment and returns the matching images.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Image, error) {
	logger.LogWithLevel(s.logger, 1, "finding image", "pattern", searchPattern)

	var images []Image
	page := ""

	// Normalize pattern to lowercase
	pattern := strings.ToLower(searchPattern)

	// Paginate through images
	for {
		resp, err := s.compute.ListImages(ctx, core.ListImagesRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list images: %w", err)
		}

		for _, oc := range resp.Items {
			img := mapToImage(oc)

			if matchAnyField(pattern, img) {
				images = append(images, img)
			}
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return images, nil
}

// matchAnyField checks if the pattern is a substring of any relevant field
func matchAnyField(pattern string, img Image) bool {
	return strings.Contains(strings.ToLower(img.Name), pattern) ||
		strings.Contains(strings.ToLower(img.OperatingSystem), pattern) ||
		strings.Contains(strings.ToLower(img.ImageOSVersion), pattern)
}

// mapToImage converts a core.Image object to an Image struct, extracting specific fields for use in the application.
func mapToImage(oc core.Image) Image {
	return Image{
		ID:              *oc.Id,
		Name:            *oc.DisplayName,
		CreatedAt:       *oc.TimeCreated,
		OperatingSystem: *oc.OperatingSystem,
		ImageOSVersion:  *oc.OperatingSystemVersion,
		LunchMode:       string(oc.LaunchMode),
	}
}
