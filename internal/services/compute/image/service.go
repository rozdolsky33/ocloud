package image

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
	"github.com/rozdolsky33/ocloud/internal/util"
	"strconv"
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

// Find performs a fuzzy search for an image using the provided search pattern and context.
// It returns a slice of matching Image objects or an error if the search fails.
func (s *Service) Find(ctx context.Context, searchPattern string) ([]Image, error) {
	logger.LogWithLevel(s.logger, 3, "finding image with bleve fuzzy search", "pattern", searchPattern)

	var allImages []Image
	var indexableDocs []IndexableImage
	page := ""

	// 1. Fetch all image
	for {
		resp, err := s.compute.ListImages(ctx, core.ListImagesRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list image: %w", err)
		}
		for _, oc := range resp.Items {
			img := mapToImage(oc)
			allImages = append(allImages, img)
			indexableDocs = append(indexableDocs, toIndexableImage(img))
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
		return nil, fmt.Errorf("failed to create index: %w", err)
	}
	for i, doc := range indexableDocs {
		err := index.Index(fmt.Sprintf("%d", i), doc)
		if err != nil {
			return nil, fmt.Errorf("indexing image failed: %w", err)
		}
	}

	// 3. Prepare a fuzzy query with wildcard
	searchPattern = strings.ToLower(searchPattern)
	if !strings.HasPrefix(searchPattern, "*") {
		searchPattern = "*" + searchPattern
	}
	if !strings.HasSuffix(searchPattern, "*") {
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
	var matched []Image
	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil || idx < 0 || idx >= len(allImages) {
			continue
		}
		matched = append(matched, allImages[idx])
	}

	logger.LogWithLevel(s.logger, 2, "found image", "count", len(matched))
	return matched, nil
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

// ToIndexableImage converts an Image object into an IndexableImage structure optimized for indexing and searching.
func toIndexableImage(img Image) IndexableImage {
	flattenedTags, _ := util.FlattenTags(img.ImageTags.FreeformTags, img.ImageTags.DefinedTags)
	tagValues, _ := util.ExtractTagValues(img.ImageTags.FreeformTags, img.ImageTags.DefinedTags)
	return IndexableImage{
		ID:              img.ID,
		Name:            strings.ToLower(img.Name),
		OperatingSystem: strings.ToLower(img.OperatingSystem),
		ImageOSVersion:  strings.ToLower(img.ImageOSVersion),
		Tags:            flattenedTags,
		TagValues:       tagValues,
	}
}
