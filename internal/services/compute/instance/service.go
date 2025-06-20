package instance

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/rozdolsky33/ocloud/internal/util"
	"strconv"
	"strings"
	"sync"

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

		// We need to fetch page tokens until we reach the desired page
		page := ""
		currentPage := 1

		for currentPage < pageNum {
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
			currentPage++
		}

		// Set the page token for the actual request
		request.Page = &page
		logger.LogWithLevel(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)
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

		// Store a reference to the instance in the map for easy lookup
		// We need to get the address of the instance in the slice
		instanceMap[*oc.Id] = &instances[len(instances)-1]
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
func (s *Service) Find(ctx context.Context, searchPattern string, showImageDetails bool) ([]Instance, error) {
	logger.LogWithLevel(s.logger, 1, "finding instances", "pattern", searchPattern)

	var instanceMap = make(map[string]*Instance)
	var allInstances []Instance
	page := ""

	// Fetch all Instances
	for {
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
		page = *resp.OpcNextPage
	}

	// Enrich with VNICs using the same approach as List
	if err := s.enrichInstancesWithVnics(ctx, instanceMap); err != nil {
		logger.LogWithLevel(s.logger, 1, "failed to enrich VNICs", "error", err)
	}

	// Enrich with image details if requested
	if showImageDetails {
		if err := s.enrichInstancesWithImageDetails(ctx, instanceMap); err != nil {
			logger.LogWithLevel(s.logger, 1, "failed to enrich image details", "error", err)
		}
	}

	// Create indexable documents from enriched instances
	var indexableDocs []IndexableInstance
	for _, inst := range instanceMap {
		indexableDocs = append(indexableDocs, toIndexableInstance(*inst))
	}

	//2. Create an in memory Bleve index
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, fmt.Errorf("failed to create Bleve index for instances: %w", err)
	}
	for i, doc := range indexableDocs {
		err := index.Index(fmt.Sprintf("%d", i), doc)
		if err != nil {
			return nil, fmt.Errorf("failed to index instances: %w", err)
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
// This method fetches the primary VNIC per instance either concurrently or sequentially based on configuration.
func (s *Service) enrichInstancesWithVnics(ctx context.Context, instanceMap map[string]*Instance) error {
	if s.enableConcurrency {
		logger.LogWithLevel(s.logger, 1, "processing VNICs in parallel (concurrency enabled)")
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, inst := range instanceMap {
			wg.Add(1)
			go func(inst *Instance) {
				defer wg.Done()
				vnic, err := s.fetchPrimaryVnicForInstance(ctx, inst.ID)
				if err != nil {
					logger.LogWithLevel(s.logger, 1, "error fetching primary VNIC", "instanceID", inst.ID, "error", err)
					return
				}
				if vnic != nil {
					mu.Lock()
					inst.IP = *vnic.PrivateIp
					inst.SubnetID = *vnic.SubnetId
					mu.Unlock()
				}
			}(inst)
		}
		wg.Wait()
	} else {
		logger.LogWithLevel(s.logger, 1, "processing VNICs sequentially (concurrency disabled)")
		for _, inst := range instanceMap {
			vnic, err := s.fetchPrimaryVnicForInstance(ctx, inst.ID)
			if err != nil {
				logger.LogWithLevel(s.logger, 1, "error fetching primary VNIC", "instanceID", inst.ID, "error", err)
				continue
			}
			if vnic != nil {
				inst.IP = *vnic.PrivateIp
				inst.SubnetID = *vnic.SubnetId
			}
		}
	}

	return nil
}

// fetchPrimaryVnicForInstance finds the primary VNIC for a given instance ID.
func (s *Service) fetchPrimaryVnicForInstance(ctx context.Context, instanceID string) (*core.Vnic, error) {
	attachments, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: &s.compartmentID,
		InstanceId:    &instanceID,
	})

	if err != nil {
		logger.LogWithLevel(s.logger, 1, "error listing VNIC attachments", "instanceID", instanceID, "error", err)
		return nil, nil
	}

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
	logger.LogWithLevel(logger.CmdLogger, 1, "Converted tags to indexable string", "<<<<<<<<flattenedTags>>>>>>>>", flattenedTags)
	tagValues, _ := util.ExtractTagValues(instance.InstanceTags.FreeformTags, instance.InstanceTags.DefinedTags)
	logger.LogWithLevel(logger.CmdLogger, 1, "Converted tag values to indexable string", "<--------tagValues-------->", tagValues)

	return IndexableInstance{
		ID:              instance.ID,
		Name:            strings.ToLower(instance.Name),
		OperatingSystem: strings.ToLower(instance.ImageOS),
		CreatedAt:       strings.ToLower(instance.CreatedAt.String()),
		Tags:            flattenedTags,
		TagValues:       tagValues,
	}
}
