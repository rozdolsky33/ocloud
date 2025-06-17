package instance

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/core"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// NewService creates a new Service instance with OCI compute and network clients using the provided AppContext.
// Returns a Service pointer and an error if the initialization fails.
func NewService(appCtx *app.AppContext) (*Service, error) {
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
// Returns instances, total count, next page token, and an error, if any.
func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]Instance, int, string, error) {
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
	logger.LogWithLevel(s.logger, 3, "Fetched instances for page totalCount:", totalCount)
	// If we have a next page, we know there are more instances
	if resp.OpcNextPage != nil {
		// Estimate total count based on current page and items per page
		totalCount = pageNum * limit
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
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

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
func (s *Service) Find(ctx context.Context, pattern string) ([]Instance, error) {
	logger.LogWithLevel(s.logger, 1, "finding instances", "pattern", pattern)

	var instanceMap = make(map[string]*Instance)

	// Try an exact match first using server-side filtering
	exactMatchResp, err := s.compute.ListInstances(ctx, core.ListInstancesRequest{
		CompartmentId:  &s.compartmentID,
		LifecycleState: core.InstanceLifecycleStateRunning,
		DisplayName:    &pattern,
	})
	if err != nil {
		return nil, fmt.Errorf("listing instances with exact name match: %w", err)
	}

	if len(exactMatchResp.Items) > 0 {
		for _, oc := range exactMatchResp.Items {
			inst := mapToInstance(oc)
			instanceMap[inst.ID] = &inst
		}
		logger.LogWithLevel(s.logger, 2, "found exact matching instances", "count", len(instanceMap))
	} else {
		// Fallback to partial match with pagination
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
				if strings.Contains(strings.ToLower(*oc.DisplayName), strings.ToLower(pattern)) {
					inst := mapToInstance(oc)
					instanceMap[inst.ID] = &inst
				}
			}

			if resp.OpcNextPage == nil {
				break
			}
			page = *resp.OpcNextPage
		}
		logger.LogWithLevel(s.logger, 2, "found partial matching instances", "count", len(instanceMap))
	}

	// Enrich with VNICs using the same approach as List
	if err := s.enrichInstancesWithVnics(ctx, instanceMap); err != nil {
		logger.LogWithLevel(s.logger, 1, "failed to enrich VNICs", "error", err)
	}

	// Convert a map to slice for return
	var result []Instance
	for _, inst := range instanceMap {
		result = append(result, *inst)
	}

	return result, nil
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
