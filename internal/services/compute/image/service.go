package image

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"github.com/rozdolsky33/ocloud/internal/services/util"
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

// List retrieves a paginated list of image with given limit and page number parameters.
// It returns the slice of image, total count, next page token, and an error if encountered.
func (s *Service) List(ctx context.Context, limit, pageNum int) ([]Image, int, string, error) {
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

	// Fetch image for the requested page
	resp, err := s.compute.ListImages(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing image: %w", err)
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

	// Process the image
	for _, oc := range resp.Items {
		images = append(images, mapToImage(oc))
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

// Find performs a fuzzy search for an image using the provided search pattern and context.
// It returns a slice of matching Image objects or an error if the search fails.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Image, error) {
	logger.LogWithLevel(s.logger, 3, "finding image with bleve fuzzy search", "pattern", searchPattern)
	var allImages []Image

	// 1. Fetch all images
	allImages, err := s.fetchAllImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all images: %w", err)
	}

	// 2. Build index
	index, err := util.BuildIndex(allImages, func(img Image) any {
		return mapToIndexableImage(img)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build index: %w", err)
	}

	// 3. Fuzzy search on multiple fields
	fields := []string{"Name", "ImageOSVersion", "OperatingSystem", "LunchMode"}
	matchedIdxs, err := util.FuzzySearchIndex(index, strings.ToLower(searchPattern), fields)
	if err != nil {
		return nil, fmt.Errorf("failed to fuzzy search index: %w", err)
	}

	// Return matched images
	var matchedImages []Image
	for _, idx := range matchedIdxs {
		if idx >= 0 && idx < len(allImages) {
			matchedImages = append(matchedImages, allImages[idx])
		}
	}

	logger.LogWithLevel(s.logger, 2, "found image", "count", len(matchedImages))
	return matchedImages, nil
}

// fetchAllImages retrieves all images from the service by paginating through the available pages.
// It returns a slice of Image objects and an error in case of failure.
func (s *Service) fetchAllImages(ctx context.Context) ([]Image, error) {
	var allImages []Image
	page := ""
	for {
		resp, err := s.compute.ListImages(ctx, core.ListImagesRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list image: %w", err)
		}
		for _, oc := range resp.Items {
			allImages = append(allImages, mapToImage(oc))
		}
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return allImages, nil
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
		ImageTags: util.ResourceTags{
			FreeformTags: oc.FreeformTags,
			DefinedTags:  oc.DefinedTags,
		},
	}
}

// mapToIndexableImage converts an Image object into an IndexableImage structure optimized for indexing and searching.
func mapToIndexableImage(img Image) IndexableImage {
	flattenedTags, _ := util.FlattenTags(img.ImageTags.FreeformTags, img.ImageTags.DefinedTags)
	tagValues, _ := util.ExtractTagValues(img.ImageTags.FreeformTags, img.ImageTags.DefinedTags)
	return IndexableImage{
		Name:            strings.ToLower(img.Name),
		ImageOSVersion:  strings.ToLower(img.ImageOSVersion),
		OperatingSystem: strings.ToLower(img.OperatingSystem),
		LunchMode:       strings.ToLower(img.LunchMode),
		Tags:            flattenedTags,
		TagValues:       tagValues,
	}
}
