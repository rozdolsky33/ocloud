package instance

import (
	"context"
	"fmt"
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
		compute:           cc,
		network:           nc,
		logger:            appCtx.Logger,
		compartmentID:     appCtx.CompartmentID,
		enableConcurrency: appCtx.EnableConcurrency,
	}, nil
}

// List retrieves instances in the compartment with pagination support.
func (s *Service) List(ctx context.Context, limit int, pageNum int) ([]Instance, int, string, error) {
	// Log input parameters at debug level
	logger.VerboseInfo(s.logger, 3, "List() called with pagination parameters",
		"limit", limit,
		"pageNum", pageNum)

	// Initialize variables
	var instances []Instance
	instanceMap := make(map[string]*Instance)
	var nextPageToken string
	var totalCount int

	// Create a request with limit parameter to fetch only the required page
	request := core.ListInstancesRequest{
		CompartmentId:  &s.compartmentID,
		LifecycleState: core.InstanceLifecycleStateRunning,
	}

	// Add limit parameter if specified
	if limit > 0 {
		request.Limit = &limit
		logger.VerboseInfo(s.logger, 3, "Setting limit parameter", "limit", limit)
	}

	// If pageNum > 1, we need to fetch the appropriate page token
	if pageNum > 1 && limit > 0 {
		logger.VerboseInfo(s.logger, 3, "Calculating page token for page", "pageNum", pageNum)

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
				logger.VerboseInfo(s.logger, 3, "Reached end of data while calculating page token",
					"currentPage", currentPage, "targetPage", pageNum)
				// Return empty result since the requested page is beyond available data
				return []Instance{}, 0, "", nil
			}

			// Move to the next page
			page = *resp.OpcNextPage
			currentPage++
		}

		// Set the page token for the actual request
		request.Page = &page
		logger.VerboseInfo(s.logger, 3, "Using page token for page", "pageNum", pageNum, "token", page)
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
		// Estimate total count based on current page and items per page
		totalCount = pageNum * limit
	}

	// Save the next page token if available
	if resp.OpcNextPage != nil {
		nextPageToken = *resp.OpcNextPage
		logger.VerboseInfo(s.logger, 3, "Next page token", "token", nextPageToken)
	}

	// Process the instances
	for _, oc := range resp.Items {
		inst := mapToInstance(oc)
		instances = append(instances, inst)

		// Store a reference to the instance in the map for easy lookup
		// We need to get the address of the instance in the slice
		instanceMap[*oc.Id] = &instances[len(instances)-1]
	}

	logger.VerboseInfo(s.logger, 3, "Fetched instances for page",
		"pageNum", pageNum, "count", len(instances))

	// Step 2: Fetch VNIC attachments for the instances in the current page
	if len(instanceMap) > 0 {
		err := s.enrichInstancesWithVnics(ctx, instanceMap)
		if err != nil {
			logger.VerboseInfo(s.logger, 1, "error enriching instances with VNICs", "error", err)
			// Continue with the instances we have, even if VNIC enrichment failed
		}
	}

	// Calculate if there are more pages after the current page
	hasNextPage := pageNum*limit < totalCount

	logger.VerboseInfo(s.logger, 2, "Completed instance listing with pagination",
		"returnedCount", len(instances),
		"totalCount", totalCount,
		"page", pageNum,
		"limit", limit,
		"hasNextPage", hasNextPage)
	return instances, totalCount, nextPageToken, nil
}

// enrichInstancesWithVnics enriches each instance in the provided map with its associated VNIC information.
// This method retrieves all VNIC attachments and processes them either concurrently or sequentially based on configuration.
func (s *Service) enrichInstancesWithVnics(ctx context.Context, instanceMap map[string]*Instance) error {
	attachments, err := s.fetchAllVnicAttachments(ctx)
	if err != nil {
		return fmt.Errorf("listing VNIC attachments: %w", err)
	}

	vnicAttachmentsByInstance := groupAttachmentsByInstance(attachments, instanceMap)

	if s.enableConcurrency {
		logger.VerboseInfo(s.logger, 1, "processing VNIC attachments in parallel (concurrency enabled)")
		return s.processVnicsConcurrently(ctx, vnicAttachmentsByInstance, instanceMap)
	}

	logger.VerboseInfo(s.logger, 1, "processing VNIC attachments sequentially (concurrency disabled)")
	return s.processVnicsSequentially(ctx, vnicAttachmentsByInstance, instanceMap)
}

// fetchAllVnicAttachments lists all VNIC attachments within a compartment, supporting pagination to retrieve all results.
func (s *Service) fetchAllVnicAttachments(ctx context.Context) ([]core.VnicAttachment, error) {
	var attachments []core.VnicAttachment
	page := ""
	for {
		resp, err := s.compute.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
			CompartmentId: &s.compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing VNIC attachments: %w", err)
		}
		attachments = append(attachments, resp.Items...)
		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return attachments, nil
}

// groupAttachmentsByInstance groups VNIC attachments by their associated instance IDs.
// Only includes attachments for which the instance exists in the given instance map.
func groupAttachmentsByInstance(attachments []core.VnicAttachment, instanceMap map[string]*Instance) map[string][]core.VnicAttachment {
	result := make(map[string][]core.VnicAttachment)
	for _, attach := range attachments {
		if attach.InstanceId == nil {
			continue
		}
		instanceID := *attach.InstanceId
		if _, ok := instanceMap[instanceID]; ok {
			result[instanceID] = append(result[instanceID], attach)
		}
	}
	return result
}

// getPrimaryVnic retrieves the primary VNIC associated with the provided VnicAttachment.
// It returns the VNIC if it is marked as primary, or nil if no primary VNIC is found.
// In case of an error during the VNIC retrieval process, it returns nil.
func (s *Service) getPrimaryVnic(ctx context.Context, attach core.VnicAttachment) (*core.Vnic, error) {
	if attach.VnicId == nil {
		logger.VerboseInfo(s.logger, 2, "VnicAttachment missing VnicId", "attachment", attach)
		return nil, nil
	}
	resp, err := s.network.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
	if err != nil {
		logger.VerboseInfo(s.logger, 2, "GetVnic error", "error", err, "vnicID", *attach.VnicId)
		return nil, nil
	}

	vnic := resp.Vnic
	if vnic.IsPrimary != nil && *vnic.IsPrimary {
		return &vnic, nil
	}
	logger.VerboseInfo(s.logger, 2, "VnicAttachment missing primary Vnic", "attachment", attach)
	return nil, nil
}

// processVnicsConcurrently processes VNIC attachments concurrently and updates the instance map with VNIC information.
func (s *Service) processVnicsConcurrently(ctx context.Context, byInstance map[string][]core.VnicAttachment, instanceMap map[string]*Instance) error {
	var wg sync.WaitGroup
	vnicChan := make(chan VnicInfo, len(byInstance))

	for instanceID, attachments := range byInstance {
		for _, attach := range attachments {
			wg.Add(1)
			go func(instanceID string, attach core.VnicAttachment) {
				defer wg.Done()
				vnic, err := s.getPrimaryVnic(ctx, attach)
				if err == nil && vnic != nil {
					vnicChan <- VnicInfo{
						InstanceID: instanceID,
						Ip:         *vnic.PrivateIp,
						SubnetID:   *vnic.SubnetId,
					}
				}
			}(instanceID, attach)
		}
	}
	go func() {
		wg.Wait()
		close(vnicChan)
	}()

	for info := range vnicChan {
		if inst, ok := instanceMap[info.InstanceID]; ok {
			inst.IP = info.Ip
			inst.SubnetID = info.SubnetID
		}
	}
	return nil
}

// processVnicsSequentially processes VNIC attachments sequentially and updates instance properties with VNIC details.
// ctx is the context for managing request deadlines, cancellations, and other request-scoped values.
// byInstance is a map where the key is instance ID and the value is a list of VNIC attachments for that instance.
// instanceMap maps instance IDs to their respective Instance struct for updating VNIC details.
// Returns an error if there is an issue processing the VNIC attachments.
func (s *Service) processVnicsSequentially(ctx context.Context, byInstance map[string][]core.VnicAttachment, instanceMap map[string]*Instance) error {
	for instanceID, attachments := range byInstance {
		for _, attach := range attachments {
			vnic, err := s.getPrimaryVnic(ctx, attach)
			if err != nil {
				continue
			}
			if vnic != nil {
				if inst, ok := instanceMap[instanceID]; ok {
					inst.IP = *vnic.PrivateIp
					inst.SubnetID = *vnic.SubnetId
				}
				break // no need to check other VNICs
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
