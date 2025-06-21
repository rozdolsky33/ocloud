package instance

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/core"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// NewService creates a new Service instance with OCI compute and network clients using the provided ApplicationContext.
// Returns a Service pointer and an error if the initialization fails.
func NewService(appCtx *app.ApplicationContext) (*Service, error) {
	cfg := appCtx.Provider
	cc, err := oci.NewComputeClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}
	nc, err := oci.NewNetworkClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}

	return &Service{
		compute:           cc,
		network:           nc,
		logger:            appCtx.Logger,
		compartmentID:     appCtx.CompartmentID,
		enableConcurrency: appCtx.EnableConcurrency,
		subnetCache:       make(map[string]*core.Subnet),
		vcnCache:          make(map[string]*core.Vcn),
		routeTableCache:   make(map[string]*core.RouteTable),
		pageTokenCache:    make(map[string]map[int]string),
	}, nil
}

// List retrieves a paginated list of running VM instances within a specified compartment.
// It supports pagination through the use of a limit and page number.
// If showImageDetails is true, it also enriches instances with image details.
// Returns instances, total count, next page token, and an error, if any.
func (s *Service) List(ctx context.Context, limit int, pageNum int, showImageDetails bool) ([]Instance, int, string, error) {
	// Log input parameters at debug level
	logger.LogWithLevel(s.logger, 3, "List() called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	// Initialize variables
	var instances []Instance
	instanceMap := make(map[string]*Instance)
	var nextPageToken string
	var totalCount int

	// Create a request with a limit parameter to fetch only the required page
	request := core.ListInstancesRequest{
		CompartmentId:  &s.compartmentID,
		LifecycleState: core.InstanceLifecycleStateRunning,
	}

	// Add limit parameter if specified
	if limit > 0 {
		request.Limit = &limit
		logger.LogWithLevel(s.logger, 3, "Setting limit parameter", "limit", limit)
	}

	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.LogWithLevel(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

		// Check if we have a cache for this compartment
		if _, ok := s.pageTokenCache[s.compartmentID]; !ok {
			s.pageTokenCache[s.compartmentID] = make(map[int]string)
		}

		// Check if we have the page token in the cache
		if token, ok := s.pageTokenCache[s.compartmentID][pageNum]; ok {
			logger.LogWithLevel(s.logger, 3, "Using cached page token", "pageNum", pageNum, "token", token)
			request.Page = &token
		} else {
			// We need to fetch page tokens until we reach the desired page
			page := ""
			currentPage := 1

			// Check if we have any cached page tokens that can help us get closer to the desired page
			var startPage int
			var startToken string
			for p := pageNum - 1; p >= 1; p-- {
				if token, ok := s.pageTokenCache[s.compartmentID][p]; ok {
					startPage = p
					startToken = token
					logger.LogWithLevel(s.logger, 3, "Found cached token for earlier page", "startPage", startPage, "token", startToken)
					break
				}
			}

			// If we found a cached token, start from there
			if startToken != "" {
				page = startToken
				currentPage = startPage + 1
			}

			for currentPage <= pageNum {
				// Fetch just the page token, not the actual data
				// Use the same limit to ensure consistent pagination
				tokenRequest := core.ListInstancesRequest{
					CompartmentId:  &s.compartmentID,
					LifecycleState: core.InstanceLifecycleStateRunning,
					Page:           &page,
				}

				if limit > 0 {
					tokenRequest.Limit = &limit
				}

				resp, err := s.compute.ListInstances(ctx, tokenRequest)
				if err != nil {
					return nil, 0, "", fmt.Errorf("fetching page token: %w", err)
				}

				// If there's no next page, we've reached the end
				if resp.OpcNextPage == nil {
					logger.LogWithLevel(s.logger, 3, "Reached end of data while calculating page token",
						"currentPage", currentPage, "targetPage", pageNum)
					// Return an empty result since the requested page is beyond available data
					return []Instance{}, 0, "", nil
				}

				// Move to the next page
				page = *resp.OpcNextPage

				// Cache the token for this page
				s.pageTokenCache[s.compartmentID][currentPage] = page
				logger.LogWithLevel(s.logger, 3, "Cached page token", "page", currentPage, "token", page)

				currentPage++
			}

			// Set the page token for the actual request
			// We use the token for the page before the one we want
			if token, ok := s.pageTokenCache[s.compartmentID][pageNum-1]; ok {
				request.Page = &token
				logger.LogWithLevel(s.logger, 3, "Using calculated page token", "pageNum", pageNum, "token", token)
			}
		}
	}

	// Fetch the instances for the requested page
	resp, err := s.compute.ListInstances(ctx, request)
	if err != nil {
		return nil, 0, "", fmt.Errorf("listing instances: %w", err)
	}

	// Set the total count to the number of instances returned
	// If we have a next page, this is an estimate
	totalCount = len(resp.Items)
	// If we have a next page, we know there are more instances
	if resp.OpcNextPage != nil {
		// If we have a next page token, we know there are more instances
		// We need to estimate the total count more accurately
		// Since we don't know the exact total count, we'll set it to a value
		// that indicates there are more pages (at least one more page worth of instances)
		totalCount = pageNum*limit + limit
	}

	// Save the next page token if available
	if resp.OpcNextPage != nil {
		nextPageToken = *resp.OpcNextPage
		logger.LogWithLevel(s.logger, 3, "Next page token", "token", nextPageToken)
	}

	// Process the instances
	for _, oc := range resp.Items {
		inst := mapToInstance(oc)
		instances = append(instances, inst)

		// Create a copy of the instance and store a pointer to it in the map
		// This ensures the pointer remains valid even if the slice is reallocated
		instanceCopy := inst
		instanceMap[*oc.Id] = &instanceCopy
	}

	logger.LogWithLevel(s.logger, 3, "Fetched instances for page",
		"pageNum", pageNum, "count", len(instances))

	// Step 2: Fetch VNIC attachments for the instances in the current page
	if len(instanceMap) > 0 {
		err := s.enrichInstancesWithVnics(ctx, instanceMap)
		if err != nil {
			logger.LogWithLevel(s.logger, 1, "error enriching instances with VNICs", "error", err)
			// Continue with the instances we have, even if VNIC enrichment failed
		}

		// Step 3: Fetch image details if requested
		if showImageDetails {
			err := s.enrichInstancesWithImageDetails(ctx, instanceMap)
			if err != nil {
				logger.LogWithLevel(s.logger, 1, "error enriching instances with image details", "error", err)
				// Continue with the instances we have, even if image details enrichment failed
			}
		}
	}

	// Update the instance slice with the enriched data from the instanceMap
	// This ensures that the returned instances have the enriched data
	for i := range instances {
		if enriched, ok := instanceMap[instances[i].ID]; ok {
			instances[i] = *enriched
		}
	}

	// Calculate if there are more pages after the current page
	// The most direct way to determine if there are more pages is to check if there's a next page token
	hasNextPage := resp.OpcNextPage != nil

	// Log detailed pagination information at debug level 1 for better visibility
	if hasNextPage {
		logger.LogWithLevel(s.logger, 1, "Pagination information",
			"currentPage", pageNum,
			"recordsOnThisPage", len(instances),
			"estimatedTotalRecords", totalCount,
			"morePages", "true")
	}

	logger.LogWithLevel(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(instances),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)
	return instances, totalCount, nextPageToken, nil
}

// Find searches for cloud instances matching the given pattern within the compartment.
// It attempts an exact name match first, followed by a partial match if necessary.
// Instances are enriched with network data (VNICs) before being returned as a list.
// If showImageDetails is true, it also enriches instances with image details.
// It uses token caching to improve performance for subsequent searches.
func (s *Service) Find(ctx context.Context, searchPattern string, showImageDetails bool) ([]Instance, error) {
	// Start overall performance tracking
	overallStartTime := time.Now()
	logger.LogWithLevel(s.logger, 1, "finding instances", "pattern", searchPattern)

	var instanceMap = make(map[string]*Instance)
	var allInstances []Instance

	// Initialize pagination variables
	var page string
	currentPage := 1

	// Check if we have a cache for this compartment
	if _, ok := s.pageTokenCache[s.compartmentID]; !ok {
		s.pageTokenCache[s.compartmentID] = make(map[int]string)
		logger.LogWithLevel(s.logger, 3, "Created new page token cache for compartment", "compartmentID", s.compartmentID)
	} else {
		logger.LogWithLevel(s.logger, 3, "Using existing page token cache", "compartmentID", s.compartmentID, "cacheSize", len(s.pageTokenCache[s.compartmentID]))
	}

	// Find the highest cached page number to optimize fetching
	var startPage int
	var startToken string
	for p := range s.pageTokenCache[s.compartmentID] {
		if p > startPage {
			startPage = p
			startToken = s.pageTokenCache[s.compartmentID][p]
		}
	}

	if startToken != "" {
		logger.LogWithLevel(s.logger, 3, "Starting from cached token", "page", startPage, "token", startToken)
		page = startToken
		currentPage = startPage + 1
	}

	// Record start time for performance tracking
	startTime := time.Now()
	fetchCount := 0

	// Fetch all Instances with token caching
	for {
		fetchCount++
		logger.LogWithLevel(s.logger, 3, "Fetching instances", "page", currentPage, "token", page, "fetchCount", fetchCount)
		resp, err := s.compute.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &s.compartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           &page,
		})
		if err != nil {
			return nil, fmt.Errorf("failed listing instances: %w", err)
		}

		for _, oc := range resp.Items {
			inst := mapToInstance(oc)
			allInstances = append(allInstances, inst)

			// Add a pointer to the instance to the map for enrichment
			instanceCopy := inst
			instanceMap[*oc.Id] = &instanceCopy
		}

		if resp.OpcNextPage == nil {
			break
		}

		// Cache the token for this page
		page = *resp.OpcNextPage
		s.pageTokenCache[s.compartmentID][currentPage] = page
		logger.LogWithLevel(s.logger, 3, "Cached page token", "page", currentPage, "token", page)
		currentPage++
	}

	// Log performance metrics for fetching
	fetchDuration := time.Since(startTime)
	logger.LogWithLevel(s.logger, 1, "Instance fetching performance",
		"totalPages", fetchCount,
		"totalInstances", len(allInstances),
		"duration", fetchDuration,
		"startedFromCachedPage", startPage > 0,
		"cachedPagesUsed", startPage)

	// Track enrichment performance
	enrichStartTime := time.Now()

	// Enrich with VNICs using the same approach as List
	vnicStartTime := time.Now()
	if err := s.enrichInstancesWithVnics(ctx, instanceMap); err != nil {
		logger.LogWithLevel(s.logger, 1, "failed to enrich VNICs", "error", err)
	}
	vnicDuration := time.Since(vnicStartTime)
	logger.LogWithLevel(s.logger, 1, "VNIC enrichment performance", "duration", vnicDuration, "instanceCount", len(instanceMap))

	// Enrich with image details if requested
	if showImageDetails {
		imageStartTime := time.Now()
		if err := s.enrichInstancesWithImageDetails(ctx, instanceMap); err != nil {
			logger.LogWithLevel(s.logger, 1, "failed to enrich image details", "error", err)
		}
		imageDuration := time.Since(imageStartTime)
		logger.LogWithLevel(s.logger, 1, "Image details enrichment performance", "duration", imageDuration, "instanceCount", len(instanceMap))
	}

	totalEnrichDuration := time.Since(enrichStartTime)
	logger.LogWithLevel(s.logger, 1, "Total enrichment performance", "duration", totalEnrichDuration, "instanceCount", len(instanceMap))

	// Track search performance
	searchStartTime := time.Now()

	// Create indexable documents from enriched instances
	indexingStartTime := time.Now()
	var indexableDocs []IndexableInstance
	for _, inst := range instanceMap {
		indexableDocs = append(indexableDocs, toIndexableInstance(*inst))
	}
	indexingDuration := time.Since(indexingStartTime)
	logger.LogWithLevel(s.logger, 3, "Created indexable documents", "count", len(indexableDocs), "duration", indexingDuration)

	// Create an in memory Bleve index
	indexCreationStartTime := time.Now()
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create Bleve index for instances: %w", err)
	}

	// Index all documents
	for i, doc := range indexableDocs {
		err := index.Index(fmt.Sprintf("%d", i), doc)
		if err != nil {
			return nil, fmt.Errorf("failed to index instances: %w", err)
		}
	}
	indexCreationDuration := time.Since(indexCreationStartTime)
	logger.LogWithLevel(s.logger, 3, "Created search index", "duration", indexCreationDuration)

	// Prepare a fuzzy query with wildcard
	queryStartTime := time.Now()
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
	queryDuration := time.Since(queryStartTime)
	logger.LogWithLevel(s.logger, 3, "Prepared search query", "query", queryString, "duration", queryDuration)

	// Perform search
	searchExecutionStartTime := time.Now()
	result, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	searchExecutionDuration := time.Since(searchExecutionStartTime)
	logger.LogWithLevel(s.logger, 3, "Executed search", "hits", len(result.Hits), "duration", searchExecutionDuration)

	// Collect matched results
	resultCollectionStartTime := time.Now()
	var matched []Instance
	for _, hit := range result.Hits {
		idx, err := strconv.Atoi(hit.ID)
		if err != nil || idx < 0 || idx >= len(indexableDocs) {
			continue
		}

		// Get the instance ID from the indexable document
		instanceID := indexableDocs[idx].ID

		// Get the enriched instance from the map
		if enriched, ok := instanceMap[instanceID]; ok {
			matched = append(matched, *enriched)
		}
	}
	resultCollectionDuration := time.Since(resultCollectionStartTime)
	logger.LogWithLevel(s.logger, 3, "Collected search results", "matchedCount", len(matched), "duration", resultCollectionDuration)

	// Log overall search performance
	totalSearchDuration := time.Since(searchStartTime)
	logger.LogWithLevel(s.logger, 1, "Search performance",
		"totalDuration", totalSearchDuration,
		"indexingDuration", indexingDuration,
		"indexCreationDuration", indexCreationDuration,
		"queryDuration", queryDuration,
		"searchExecutionDuration", searchExecutionDuration,
		"resultCollectionDuration", resultCollectionDuration,
		"matchedCount", len(matched),
		"totalInstancesSearched", len(allInstances))
	// Log overall performance for the entire Find operation
	overallDuration := time.Since(overallStartTime)
	logger.LogWithLevel(s.logger, 1, "Overall Find operation performance",
		"totalDuration", overallDuration,
		"matchedCount", len(matched),
		"totalInstancesSearched", len(allInstances),
		"startedFromCachedPage", startPage > 0,
		"cachedPagesUsed", startPage)

	logger.LogWithLevel(s.logger, 2, "found instances", "count", len(matched))
	return matched, nil
}

// enrichInstancesWithImageDetails enriches each instance in the provided map with its associated image details.
// This method fetches image details for each instance either concurrently or sequentially based on configuration.
func (s *Service) enrichInstancesWithImageDetails(ctx context.Context, instanceMap map[string]*Instance) error {
	if s.enableConcurrency {
		logger.LogWithLevel(s.logger, 1, "processing image details in parallel (concurrency enabled)")
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, inst := range instanceMap {
			wg.Add(1)
			go func(inst *Instance) {
				defer wg.Done()
				imageDetails, err := s.fetchImageDetails(ctx, inst.ImageID)
				if err != nil {
					logger.LogWithLevel(s.logger, 1, "error fetching image details", "imageID", inst.ImageID, "error", err)
					return
				}
				if imageDetails != nil {
					mu.Lock()
					if imageDetails.DisplayName != nil {
						inst.ImageName = *imageDetails.DisplayName
					}
					if imageDetails.OperatingSystem != nil {
						inst.ImageOS = *imageDetails.OperatingSystem
					}
					// Copy free-form tags
					if len(imageDetails.FreeformTags) > 0 {
						inst.InstanceTags.FreeformTags = make(map[string]string)
						for k, v := range imageDetails.FreeformTags {
							inst.InstanceTags.FreeformTags[k] = v
						}
						logger.LogWithLevel(s.logger, 1, "freeform tags", "tags", inst.InstanceTags.FreeformTags)
					}
					// Copy defined tags
					if len(imageDetails.DefinedTags) > 0 {
						inst.InstanceTags.DefinedTags = make(map[string]map[string]interface{})
						for namespace, tags := range imageDetails.DefinedTags {
							inst.InstanceTags.DefinedTags[namespace] = make(map[string]interface{})
							for k, v := range tags {
								inst.InstanceTags.DefinedTags[namespace][k] = v
							}
						}
						logger.LogWithLevel(s.logger, 1, "defined tags", "tags", inst.InstanceTags.DefinedTags)
					}
					mu.Unlock()
				}
			}(inst)
		}
		wg.Wait()
	} else {
		logger.LogWithLevel(s.logger, 1, "processing image details sequentially (concurrency disabled)")
		for _, inst := range instanceMap {
			imageDetails, err := s.fetchImageDetails(ctx, inst.ImageID)
			if err != nil {
				logger.LogWithLevel(s.logger, 1, "error fetching image details", "imageID", inst.ImageID, "error", err)
				continue
			}
			if imageDetails != nil {
				if imageDetails.DisplayName != nil {
					inst.ImageName = *imageDetails.DisplayName
				}
				if imageDetails.OperatingSystem != nil {
					inst.ImageOS = *imageDetails.OperatingSystem
				}
				// Copy free-form tags
				if len(imageDetails.FreeformTags) > 0 {
					inst.InstanceTags.FreeformTags = make(map[string]string)
					for k, v := range imageDetails.FreeformTags {
						inst.InstanceTags.FreeformTags[k] = v
					}
					logger.LogWithLevel(s.logger, 1, "freeform tags", "tags", inst.InstanceTags.FreeformTags)
				}
				// Copy defined tags
				if len(imageDetails.DefinedTags) > 0 {
					inst.InstanceTags.DefinedTags = make(map[string]map[string]interface{})
					for namespace, tags := range imageDetails.DefinedTags {
						inst.InstanceTags.DefinedTags[namespace] = make(map[string]interface{})
						for k, v := range tags {
							inst.InstanceTags.DefinedTags[namespace][k] = v
						}
					}
					logger.LogWithLevel(s.logger, 1, "defined tags", "tags", inst.InstanceTags.DefinedTags)
				}
			}
		}
	}

	return nil
}

// fetchImageDetails fetches the details of an image from OCI
func (s *Service) fetchImageDetails(ctx context.Context, imageID string) (*core.Image, error) {
	if imageID == "" {
		return nil, nil
	}

	// Create a request to get the image details
	request := core.GetImageRequest{
		ImageId: &imageID,
	}

	// Call the OCI API to get the image details
	response, err := s.compute.GetImage(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("getting image details: %w", err)
	}

	return &response.Image, nil
}

// enrichInstancesWithVnics enriches each instance in the provided map with its associated VNIC information.
// This method uses a batch approach to fetch VNIC attachments for all instances at once, reducing API calls.
func (s *Service) enrichInstancesWithVnics(ctx context.Context, instanceMap map[string]*Instance) error {
	// Extract instance IDs
	var instanceIDs []string
	for id := range instanceMap {
		instanceIDs = append(instanceIDs, id)
	}

	// Batch fetches VNIC attachments for all instances
	logger.LogWithLevel(s.logger, 1, "batch fetching VNIC attachments", "instanceCount", len(instanceIDs))
	attachmentsMap, err := s.batchFetchVnicAttachments(ctx, instanceIDs)
	if err != nil {
		logger.LogWithLevel(s.logger, 1, "error batch fetching VNIC attachments", "error", err)
		// Continue with individual fetching as fallback
	}

	if s.enableConcurrency {
		logger.LogWithLevel(s.logger, 1, "processing VNICs in parallel (concurrency enabled)")
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, inst := range instanceMap {
			wg.Add(1)
			go func(inst *Instance) {
				defer wg.Done()

				// Try to get VNIC attachments from the batch result first
				var vnic *core.Vnic
				if attachments, ok := attachmentsMap[inst.ID]; ok && len(attachments) > 0 {
					// Find primary VNIC from attachments
					for _, attach := range attachments {
						primaryVnic, err := s.getPrimaryVnic(ctx, attach)
						if err != nil {
							logger.LogWithLevel(s.logger, 1, "error getting primary VNIC from batch", "instanceID", inst.ID, "error", err)
							continue
						}
						if primaryVnic != nil {
							vnic = primaryVnic
							break
						}
					}
				}

				// Fallback to individual fetch if not found in a batch
				if vnic == nil {
					var err error
					vnic, err = s.fetchPrimaryVnicForInstance(ctx, inst.ID)
					if err != nil {
						logger.LogWithLevel(s.logger, 1, "error fetching primary VNIC", "instanceID", inst.ID, "error", err)
						return
					}
				}

				if vnic != nil {
					mu.Lock()

					// Basic VNIC information
					inst.IP = *vnic.PrivateIp
					inst.SubnetID = *vnic.SubnetId

					// Set hostname if available
					if vnic.HostnameLabel != nil {
						inst.Hostname = *vnic.HostnameLabel
					}

					// Fetch subnet details
					subnetDetails, err := s.fetchSubnetDetails(ctx, *vnic.SubnetId)
					if err != nil {
						logger.LogWithLevel(s.logger, 1, "error fetching subnet details", "subnetID", *vnic.SubnetId, "error", err)
					} else if subnetDetails != nil {
						// Set subnet name
						if subnetDetails.DisplayName != nil {
							inst.SubnetName = *subnetDetails.DisplayName
						}

						// Set VCN ID
						if subnetDetails.VcnId != nil {
							inst.VcnID = *subnetDetails.VcnId

							// Fetch VCN details
							vcnDetails, err := s.fetchVcnDetails(ctx, *subnetDetails.VcnId)
							if err != nil {
								logger.LogWithLevel(s.logger, 1, "error fetching VCN details", "vcnID", *subnetDetails.VcnId, "error", err)
							} else if vcnDetails != nil && vcnDetails.DisplayName != nil {
								inst.VcnName = *vcnDetails.DisplayName
							}
						}

						// Set private DNS enabled flag
						if subnetDetails.DnsLabel != nil && *subnetDetails.DnsLabel != "" {
							inst.PrivateDNSEnabled = true
						}

						// Set a route table ID and name
						if subnetDetails.RouteTableId != nil {
							inst.RouteTableID = *subnetDetails.RouteTableId

							// Fetch route table details
							routeTableDetails, err := s.fetchRouteTableDetails(ctx, *subnetDetails.RouteTableId)
							if err != nil {
								logger.LogWithLevel(s.logger, 1, "error fetching route table details", "routeTableID", *subnetDetails.RouteTableId, "error", err)
							} else if routeTableDetails != nil && routeTableDetails.DisplayName != nil {
								inst.RouteTableName = *routeTableDetails.DisplayName
							}
						}
					}
					mu.Unlock()
				}
			}(inst)
		}
		wg.Wait()
	} else {
		logger.LogWithLevel(s.logger, 1, "processing VNICs sequentially (concurrency disabled)")
		for _, inst := range instanceMap {
			// Try to get VNIC attachments from the batch result first
			var vnic *core.Vnic
			if attachments, ok := attachmentsMap[inst.ID]; ok && len(attachments) > 0 {
				// Find primary VNIC from attachments
				for _, attach := range attachments {
					primaryVnic, err := s.getPrimaryVnic(ctx, attach)
					if err != nil {
						logger.LogWithLevel(s.logger, 1, "error getting primary VNIC from batch", "instanceID", inst.ID, "error", err)
						continue
					}
					if primaryVnic != nil {
						vnic = primaryVnic
						break
					}
				}
			}

			// Fallback to individual fetch if not found in a batch
			if vnic == nil {
				var err error
				vnic, err = s.fetchPrimaryVnicForInstance(ctx, inst.ID)
				if err != nil {
					logger.LogWithLevel(s.logger, 1, "error fetching primary VNIC", "instanceID", inst.ID, "error", err)
					continue
				}
			}

			if vnic != nil {
				// Basic VNIC information
				inst.IP = *vnic.PrivateIp
				inst.SubnetID = *vnic.SubnetId

				// Set hostname if available
				if vnic.HostnameLabel != nil {
					inst.Hostname = *vnic.HostnameLabel
				}

				// Fetch subnet details
				subnetDetails, err := s.fetchSubnetDetails(ctx, *vnic.SubnetId)
				if err != nil {
					logger.LogWithLevel(s.logger, 1, "error fetching subnet details", "subnetID", *vnic.SubnetId, "error", err)
				} else if subnetDetails != nil {
					// Set subnet name
					if subnetDetails.DisplayName != nil {
						inst.SubnetName = *subnetDetails.DisplayName
					}

					// Set VCN ID
					if subnetDetails.VcnId != nil {
						inst.VcnID = *subnetDetails.VcnId

						// Fetch VCN details
						vcnDetails, err := s.fetchVcnDetails(ctx, *subnetDetails.VcnId)
						if err != nil {
							logger.LogWithLevel(s.logger, 1, "error fetching VCN details", "vcnID", *subnetDetails.VcnId, "error", err)
						} else if vcnDetails != nil && vcnDetails.DisplayName != nil {
							inst.VcnName = *vcnDetails.DisplayName
						}
					}

					// Set private DNS enabled flag
					if subnetDetails.DnsLabel != nil && *subnetDetails.DnsLabel != "" {
						inst.PrivateDNSEnabled = true
					}

					// Set a route table ID and name
					if subnetDetails.RouteTableId != nil {
						inst.RouteTableID = *subnetDetails.RouteTableId

						// Fetch route table details
						routeTableDetails, err := s.fetchRouteTableDetails(ctx, *subnetDetails.RouteTableId)
						if err != nil {
							logger.LogWithLevel(s.logger, 1, "error fetching route table details", "routeTableID", *subnetDetails.RouteTableId, "error", err)
						} else if routeTableDetails != nil && routeTableDetails.DisplayName != nil {
							inst.RouteTableName = *routeTableDetails.DisplayName
						}
					}
				}
			}
		}
	}

	return nil
}

// fetchPrimaryVnicForInstance finds the primary VNIC for a given instance ID.
// It uses a batch approach to fetch VNIC attachments for multiple instances at once.
func (s *Service) fetchPrimaryVnicForInstance(ctx context.Context, instanceID string) (*core.Vnic, error) {
	// Fetch all VNIC attachments for the instance
	attachments, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: &s.compartmentID,
		InstanceId:    &instanceID,
	})

	if err != nil {
		logger.LogWithLevel(s.logger, 1, "error listing VNIC attachments", "instanceID", instanceID, "error", err)
		return nil, nil
	}

	// Find the primary VNIC
	for _, attach := range attachments.Items {
		vnic, err := s.getPrimaryVnic(ctx, attach)
		if err != nil {
			logger.LogWithLevel(s.logger, 1, "error getting primary VNIC", "instanceID", instanceID, "error", err)
			continue
		}
		if vnic != nil {
			return vnic, nil
		}
	}
	logger.LogWithLevel(s.logger, 1, "no primary VNIC found for instance", "instanceID", instanceID)

	return nil, nil
}

// batchFetchVnicAttachments fetches VNIC attachments for multiple instances in a single API call.
// It returns a map of instance ID to VNIC attachments.
func (s *Service) batchFetchVnicAttachments(ctx context.Context, instanceIDs []string) (map[string][]core.VnicAttachment, error) {
	result := make(map[string][]core.VnicAttachment)

	// Fetch all VNIC attachments for the compartment
	var page string
	for {
		resp, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing VNIC attachments: %w", err)
		}

		// Filter attachments by instance ID
		for _, attach := range resp.Items {
			if attach.InstanceId == nil {
				continue
			}

			// Check if this attachment belongs to one of our instances
			for _, id := range instanceIDs {
				if *attach.InstanceId == id {
					result[id] = append(result[id], attach)
					break
				}
			}
		}

		// Check if there are more pages
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return result, nil
}

// getPrimaryVnic retrieves the primary VNIC associated with the provided VnicAttachment.
// It returns the VNIC if it is marked as primary, or nil if no primary VNIC is found.
// In case of an error during the VNIC retrieval process, it returns nil.
func (s *Service) getPrimaryVnic(ctx context.Context, attach core.VnicAttachment) (*core.Vnic, error) {
	if attach.VnicId == nil {
		logger.LogWithLevel(s.logger, 2, "VnicAttachment missing VnicId", "attachment", attach)
		return nil, nil
	}
	resp, err := s.network.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
	if err != nil {
		logger.LogWithLevel(s.logger, 2, "GetVnic error", "error", err, "vnicID", *attach.VnicId)
		return nil, nil
	}

	vnic := resp.Vnic
	if vnic.IsPrimary != nil && *vnic.IsPrimary {
		return &vnic, nil
	}
	logger.LogWithLevel(s.logger, 2, "VnicAttachment missing primary Vnic", "attachment", attach)
	return nil, nil
}

// fetchSubnetDetails retrieves the subnet details for the given subnet ID.
// It uses a cache to avoid making repeated API calls for the same subnet.
func (s *Service) fetchSubnetDetails(ctx context.Context, subnetID string) (*core.Subnet, error) {
	// Check cache first
	if subnet, ok := s.subnetCache[subnetID]; ok {
		logger.LogWithLevel(s.logger, 3, "subnet cache hit", "subnetID", subnetID)
		return subnet, nil
	}

	// Cache miss, fetch from API
	logger.LogWithLevel(s.logger, 3, "subnet cache miss", "subnetID", subnetID)
	resp, err := s.network.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &subnetID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting subnet details: %w", err)
	}

	// Store in cache
	s.subnetCache[subnetID] = &resp.Subnet
	return &resp.Subnet, nil
}

// fetchVcnDetails retrieves the VCN details for the given VCN ID.
// It uses a cache to avoid making repeated API calls for the same VCN.
func (s *Service) fetchVcnDetails(ctx context.Context, vcnID string) (*core.Vcn, error) {
	// Check cache first
	if vcn, ok := s.vcnCache[vcnID]; ok {
		logger.LogWithLevel(s.logger, 3, "VCN cache hit", "vcnID", vcnID)
		return vcn, nil
	}

	// Cache miss, fetch from API
	logger.LogWithLevel(s.logger, 3, "VCN cache miss", "vcnID", vcnID)
	resp, err := s.network.GetVcn(ctx, core.GetVcnRequest{
		VcnId: &vcnID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting VCN details: %w", err)
	}

	// Store in cache
	s.vcnCache[vcnID] = &resp.Vcn
	return &resp.Vcn, nil
}

// fetchRouteTableDetails retrieves the route table details for the given route table ID.
// It uses a cache to avoid making repeated API calls for the same route table.
func (s *Service) fetchRouteTableDetails(ctx context.Context, routeTableID string) (*core.RouteTable, error) {
	// Check cache first
	if routeTable, ok := s.routeTableCache[routeTableID]; ok {
		logger.LogWithLevel(s.logger, 3, "route table cache hit", "routeTableID", routeTableID)
		return routeTable, nil
	}

	// Cache miss, fetch from API
	logger.LogWithLevel(s.logger, 3, "route table cache miss", "routeTableID", routeTableID)
	resp, err := s.network.GetRouteTable(ctx, core.GetRouteTableRequest{
		RtId: &routeTableID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting route table details: %w", err)
	}

	// Store in cache
	s.routeTableCache[routeTableID] = &resp.RouteTable
	return &resp.RouteTable, nil
}

// mapToInstance maps SDK Instance to local model.
func mapToInstance(oc core.Instance) Instance {
	return Instance{
		Name: *oc.DisplayName,
		ID:   *oc.Id,
		IP:   "", // to be filled later
		Placement: Placement{
			Region:             *oc.Region,
			AvailabilityDomain: *oc.AvailabilityDomain,
			FaultDomain:        *oc.FaultDomain,
		},
		Resources: Resources{
			VCPUs:    int(*oc.ShapeConfig.Vcpus),
			MemoryGB: *oc.ShapeConfig.MemoryInGBs,
		},
		Shape:     *oc.Shape,
		ImageID:   *oc.ImageId,
		SubnetID:  "", // to be filled later
		State:     oc.LifecycleState,
		CreatedAt: *oc.TimeCreated,
	}
}

// ToIndexableInstance converts an Instance into an IndexableInstance with simplified and normalized fields for indexing.
func toIndexableInstance(instance Instance) IndexableInstance {
	flattenedTags, _ := util.FlattenTags(instance.InstanceTags.FreeformTags, instance.InstanceTags.DefinedTags)
	tagValues, _ := util.ExtractTagValues(instance.InstanceTags.FreeformTags, instance.InstanceTags.DefinedTags)

	return IndexableInstance{
		ID:                   instance.ID,
		Name:                 strings.ToLower(instance.Name),
		ImageName:            strings.ToLower(instance.ImageName),
		ImageOperatingSystem: strings.ToLower(instance.ImageOS),
		Tags:                 flattenedTags,
		TagValues:            tagValues,
	}
}
