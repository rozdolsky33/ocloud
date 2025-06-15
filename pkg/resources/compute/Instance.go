package compute

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/printer"
	"strings"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// NewService constructs a compute Service, wiring up clients once.
func NewService(cfg common.ConfigurationProvider, appCtx *app.AppContext) (*Service, error) {
	cc, err := oci.NewComputeClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute client: %w", err)
	}
	nc, err := oci.NewNetworkClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}

	return &Service{
		compute:            cc,
		network:            nc,
		logger:             appCtx.Logger,
		compartmentID:      appCtx.CompartmentID,
		disableConcurrency: appCtx.DisableConcurrency,
	}, nil
}

// List retrieves instances in the compartment with pagination support.
func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]Instance, int, string, error) {
	// Log input parameters at debug level
	logger.VerboseInfo(s.logger, 3, "List() called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	var all []Instance
	page := ""
	var instanceIDs []string
	instanceMap := make(map[string]*Instance)
	totalCount := 0
	var nextPageToken string

	// Step 1: Fetch all instances
	for {
		resp, err := s.compute.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &s.compartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           &page,
		})
		if err != nil {
			return nil, 0, "", fmt.Errorf("listing instances: %w", err)
		}

		// Add instances to our collection
		for _, oc := range resp.Items {
			inst := mapToInstance(oc)
			all = append(all, inst)

			// Store instance ID for bulk VNIC attachment lookup
			instanceIDs = append(instanceIDs, *oc.Id)

			// Store a reference to the instance in the map for easy lookup
			// We need to get the address of the instance in the slice
			instanceMap[*oc.Id] = &all[len(all)-1]
		}

		// Keep track of total count
		totalCount = len(all)

		// Log the number of instances fetched so far
		logger.VerboseInfo(s.logger, 3, "Fetched instances",
			"count", len(resp.Items),
			"totalSoFar", totalCount)

		// If we've collected enough instances for the requested page, or there are no more pages, break
		if resp.OpcNextPage == nil {
			logger.VerboseInfo(s.logger, 3, "No more pages available")
			break
		}

		// Save the next page token
		nextPageToken = *resp.OpcNextPage
		logger.VerboseInfo(s.logger, 3, "Next page token", "token", nextPageToken)

		// If we're not collecting all instances for pagination, break after the first page
		if limit > 0 {
			logger.VerboseInfo(s.logger, 3, "Stopping after first page due to limit", "limit", limit)
			break
		}

		page = nextPageToken
	}

	// Apply pagination if requested
	paginatedInstances := all
	if limit > 0 && pageNum > 0 {
		logger.VerboseInfo(s.logger, 3, "Applying pagination",
			"pageNum", pageNum,
			"limit", limit,
			"totalInstances", len(all))

		// Calculate start and end indices for the requested page
		startIdx := (pageNum - 1) * limit
		endIdx := startIdx + limit

		logger.VerboseInfo(s.logger, 3, "Calculated pagination indices",
			"startIdx", startIdx,
			"endIdx", endIdx)

		// Adjust indices if they're out of bounds
		if startIdx >= len(all) {
			// If the requested page is beyond available data, return empty result
			logger.VerboseInfo(s.logger, 3, "Requested page is beyond available data",
				"startIdx", startIdx,
				"totalInstances", len(all))
			return []Instance{}, totalCount, nextPageToken, nil
		}
		if endIdx > len(all) {
			endIdx = len(all)
			logger.VerboseInfo(s.logger, 3, "Adjusted end index to match available data",
				"endIdx", endIdx)
		}

		// Extract the requested page
		paginatedInstances = all[startIdx:endIdx]
		logger.VerboseInfo(s.logger, 3, "Extracted page of instances",
			"pageSize", len(paginatedInstances))

		// Update instance map to only include instances in the current page
		newInstanceMap := make(map[string]*Instance)
		for i := range paginatedInstances {
			newInstanceMap[paginatedInstances[i].ID] = &paginatedInstances[i]
		}
		instanceMap = newInstanceMap
		logger.VerboseInfo(s.logger, 3, "Updated instance map for current page",
			"mapSize", len(instanceMap))
	}

	// Step 2: Fetch VNIC attachments for the instances in the current page
	if len(instanceMap) > 0 {
		err := s.enrichInstancesWithVnics(ctx, instanceMap)
		if err != nil {
			logger.VerboseInfo(s.logger, 1, "error enriching instances with VNICs", "error", err)
			// Continue with the instances we have, even if VNIC enrichment failed
		}
	}

	logger.VerboseInfo(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(paginatedInstances),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", nextPageToken != "")
	return paginatedInstances, totalCount, nextPageToken, nil
}

// enrichInstancesWithVnics fetches VNIC attachments for all instances in bulk
// and updates the instances with their primary VNIC information.
func (s *Service) enrichInstancesWithVnics(ctx context.Context, instanceMap map[string]*Instance) error {
	// Create a map to store VNIC attachments by instance ID
	vnicAttachmentsByInstance := make(map[string][]core.VnicAttachment)

	// Fetch all VNIC attachments for the compartment
	page := ""
	for {
		vaResp, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return fmt.Errorf("listing VNIC attachments: %w", err)
		}

		// Group VNIC attachments by instance ID
		for _, attach := range vaResp.Items {
			if attach.InstanceId != nil {
				instanceID := *attach.InstanceId
				vnicAttachmentsByInstance[instanceID] = append(vnicAttachmentsByInstance[instanceID], attach)
			}
		}

		if vaResp.OpcNextPage == nil {
			break
		}
		page = *vaResp.OpcNextPage
	}

	// Check if concurrency is disabled
	if s.disableConcurrency {
		// Process VNIC attachments sequentially
		logger.VerboseInfo(s.logger, 1, "processing VNIC attachments sequentially (concurrency disabled)")

		// For each instance, find its primary VNIC
		for instanceID, attachments := range vnicAttachmentsByInstance {
			// Skip if we don't have this instance in our map
			if _, ok := instanceMap[instanceID]; !ok {
				continue
			}

			// Process each instance's VNIC attachments sequentially
			for _, attach := range attachments {
				// Skip if VNIC ID is nil
				if attach.VnicId == nil {
					continue
				}

				// Get VNIC details
				vnic, err := s.network.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
				if err != nil {
					logger.VerboseInfo(s.logger, 2, "GetVnic error", "error", err, "vnicID", *attach.VnicId)
					continue
				}

				// Check if this is the primary VNIC
				if vnic.IsPrimary != nil && *vnic.IsPrimary {
					if instance, ok := instanceMap[instanceID]; ok {
						instance.IP = *vnic.PrivateIp
						instance.SubnetID = *vnic.SubnetId
					}
					// Found primary VNIC, no need to check other attachments for this instance
					break
				}
			}
		}
	} else {
		// Process VNIC attachments in parallel (default behavior)
		logger.VerboseInfo(s.logger, 1, "processing VNIC attachments in parallel (concurrency enabled)")

		// Use a wait group to process VNICs in parallel
		var wg sync.WaitGroup
		vnicChan := make(chan VnicInfo, len(instanceMap))

		// For each instance, find its primary VNIC
		for instanceID, attachments := range vnicAttachmentsByInstance {
			// Skip if we don't have this instance in our map
			if _, ok := instanceMap[instanceID]; !ok {
				continue
			}

			// Process each instance's VNIC attachments in a separate goroutine
			for _, attach := range attachments {
				wg.Add(1)
				go func(instanceID string, attach core.VnicAttachment) {
					defer wg.Done()

					// Skip if VNIC ID is nil
					if attach.VnicId == nil {
						return
					}

					// Get VNIC details
					vnic, err := s.network.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
					if err != nil {
						logger.VerboseInfo(s.logger, 2, "GetVnic error", "error", err, "vnicID", *attach.VnicId)
						return
					}

					// Check if this is the primary VNIC
					if vnic.IsPrimary != nil && *vnic.IsPrimary {
						vnicChan <- VnicInfo{
							InstanceID: instanceID,
							Ip:         *vnic.PrivateIp,
							SubnetID:   *vnic.SubnetId,
						}
					}
				}(instanceID, attach)
			}
		}

		// Close the channel when all goroutines are done
		go func() {
			wg.Wait()
			close(vnicChan)
		}()

		// Process results from the channel
		for info := range vnicChan {
			if instance, ok := instanceMap[info.InstanceID]; ok {
				instance.IP = info.Ip
				instance.SubnetID = info.SubnetID
			}
		}
	}

	return nil
}

// Find searches instances by name pattern.
func (s *Service) Find(ctx context.Context, pattern string) ([]Instance, error) {
	logger.VerboseInfo(s.logger, 1, "finding instances", "pattern", pattern)

	// Check if the pattern is an exact match for a display name
	// If so, we can use the server-side filtering for better performance
	exactMatchResp, err := s.compute.ListInstances(ctx, core.ListInstancesRequest{
		CompartmentId:  &s.compartmentID,
		LifecycleState: core.InstanceLifecycleStateRunning,
		DisplayName:    &pattern,
	})
	if err != nil {
		return nil, fmt.Errorf("listing instances with exact name match: %w", err)
	}

	// If we found an exact match, return it
	if len(exactMatchResp.Items) > 0 {
		var matched []Instance
		for _, oc := range exactMatchResp.Items {
			inst := mapToInstance(oc)
			// enrich IP & Subnet
			err := s.enrichVnic(ctx, &inst, oc.Id)
			if err != nil {
				logger.VerboseInfo(s.logger, 1, "failed to enrich VNIC info", "instance", *oc.DisplayName, "error", err)
			}
			matched = append(matched, inst)
		}
		logger.VerboseInfo(s.logger, 2, "found exact matching instances", "pattern", pattern, "count", len(matched))
		return matched, nil
	}

	// If no exact match, we need to do a partial match
	// We'll use pagination to avoid loading all instances at once
	var matched []Instance
	page := ""

	for {
		resp, err := s.compute.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &s.compartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing instances for partial name match: %w", err)
		}

		for _, oc := range resp.Items {
			// Check if the display name contains the pattern (case-insensitive)
			if strings.Contains(strings.ToLower(*oc.DisplayName), strings.ToLower(pattern)) {
				inst := mapToInstance(oc)
				// enrich IP & Subnet
				err := s.enrichVnic(ctx, &inst, oc.Id)
				if err != nil {
					logger.VerboseInfo(s.logger, 1, "failed to enrich VNIC info", "instance", *oc.DisplayName, "error", err)
				}
				matched = append(matched, inst)
			}
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	logger.VerboseInfo(s.logger, 2, "found partial matching instances", "pattern", pattern, "count", len(matched))
	return matched, nil
}

// enrichVnic queries attachments and sets IP/Subnet on the model.
// This is a best-effort operation - if no primary VNIC is found, the instance
// will be returned with empty IP/Subnet fields rather than failing.
func (s *Service) enrichVnic(ctx context.Context, inst *Instance, instanceID *string) error {
	vaResp, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: &s.compartmentID,
		InstanceId:    instanceID,
	})
	if err != nil {
		// Log the error but don't fail the entire operation
		logger.VerboseInfo(s.logger, 1, "error listing VNIC attachments", "instanceID", *instanceID, "error", err)
		return nil
	}

	for _, attach := range vaResp.Items {
		vnic, err := s.network.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
		if err != nil {
			logger.VerboseInfo(s.logger, 2, "GetVnic error", "error", err)
			continue
		}
		if vnic.IsPrimary != nil && *vnic.IsPrimary {
			inst.IP = *vnic.PrivateIp
			inst.SubnetID = *vnic.SubnetId
			return nil
		}
	}

	// If no primary VNIC is found, log a warning but don't fail the operation
	logger.VerboseInfo(s.logger, 1, "no primary VNIC found for instance", "instanceID", *instanceID)
	return nil
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
			VCPUs:    *oc.ShapeConfig.Vcpus,
			MemoryGB: *oc.ShapeConfig.MemoryInGBs,
		},
		Shape:     *oc.Shape,
		ImageID:   *oc.ImageId,
		SubnetID:  "", // to be filled later
		State:     oc.LifecycleState,
		CreatedAt: *oc.TimeCreated,
	}
}

// ListInstances lists instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the AppContext struct and supports pagination.
func ListInstances(appCtx *app.AppContext, limit int, page int) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(appCtx.Logger, 1, "ListInstances()", "limit", limit, "page", page)

	service, err := NewService(appCtx.Provider, appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	instances, totalCount, nextPageToken, err := service.List(ctx, limit, page)
	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	// Display instance information with pagination details
	PrintInstancesTable(instances, appCtx, &PaginationInfo{
		CurrentPage:   page,
		TotalCount:    totalCount,
		Limit:         limit,
		NextPageToken: nextPageToken,
	})

	return nil
}

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It uses the pre-initialized compute and network clients from the AppContext struct.
// Parameters:
// - appCtx: The application with all clients, logger, and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include image details in the output.
// Returns an error if the operation fails.
func FindInstances(appCtx *app.AppContext, namePattern string, showImageDetails bool) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(appCtx.Logger, 1, "FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)

	service, err := NewService(appCtx.Provider, appCtx)
	if err != nil {
		return fmt.Errorf("creating compute service: %w", err)
	}

	ctx := context.Background()
	matchedInstances, err := service.Find(ctx, namePattern)
	if err != nil {
		return fmt.Errorf("finding instances: %w", err)
	}

	// Display matched instances
	if len(matchedInstances) == 0 {
		fmt.Printf("No instances found matching pattern: %s\n", namePattern)
		return nil
	}

	// If showImageDetails is true, fetch and display image information
	if showImageDetails {
		// This would be implemented in a future update
		fmt.Println("Image details functionality not yet implemented")
	}

	PrintInstancesTable(matchedInstances, appCtx, nil)
	return nil
}

func PrintInstancesTable(instances []Instance, appCtx *app.AppContext, pagination *PaginationInfo) {
	// Create a table printer with the tenancy name as the title
	tablePrinter := printer.NewTablePrinter(appCtx.TenancyName)

	// Convert instances to a format suitable for the printer
	if len(instances) == 0 {
		fmt.Println("No instances found.")
		if pagination != nil && pagination.TotalCount > 0 {
			fmt.Printf("Page %d is empty. Total records: %d\n", pagination.CurrentPage, pagination.TotalCount)
			if pagination.CurrentPage > 1 {
				fmt.Printf("Try a lower page number (e.g., --page %d)\n", pagination.CurrentPage-1)
			}
		}
		return
	}

	// Print each instance as a key-value table with a title
	for _, instance := range instances {
		// Create a map with the instance data
		instanceData := map[string]string{
			"ID":         instance.ID,
			"AD":         instance.Placement.AvailabilityDomain,
			"FD":         instance.Placement.FaultDomain,
			"Region":     instance.Placement.Region,
			"Shape":      instance.Shape,
			"vCPUs":      fmt.Sprintf("%d", instance.Resources.VCPUs),
			"Created":    instance.CreatedAt.String(),
			"Subnet ID":  instance.SubnetID,
			"Name":       instance.Name,
			"Private IP": instance.IP,
			"Memory":     fmt.Sprintf("%d GB", int(instance.Resources.MemoryGB)),
			"State":      string(instance.State),
		}

		// Define the order of keys to match the example
		orderedKeys := []string{
			"ID",
			"AD",
			"FD",
			"Region",
			"Shape",
			"vCPUs",
			"Created",
			"Subnet ID",
			"Name",
			"Private IP",
			"Memory",
			"State",
		}

		// Print the table with ordered keys and colored title components
		tablePrinter.PrintKeyValueTableWithTitleOrdered(appCtx, instance.Name, instanceData, orderedKeys)
	}

	// Log pagination information if available
	if pagination != nil {
		// Log pagination information at INFO level
		appCtx.Logger.Info("--- Pagination Information ---",
			"page", pagination.CurrentPage,
			"records", fmt.Sprintf("%d/%d", len(instances), pagination.TotalCount),
			"limit", pagination.Limit)

		// Add debug logs for navigation hints
		if pagination.CurrentPage > 1 {
			logger.VerboseInfo(appCtx.Logger, 2, "Pagination navigation",
				"action", "previous page",
				"page", pagination.CurrentPage-1,
				"limit", pagination.Limit)
		}

		if len(instances) == pagination.Limit && len(instances) < pagination.TotalCount {
			logger.VerboseInfo(appCtx.Logger, 2, "Pagination navigation",
				"action", "next page",
				"page", pagination.CurrentPage+1,
				"limit", pagination.Limit)
		}
	}
}
