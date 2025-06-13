package compute

import (
	"context"
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/core"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// ListInstances lists all instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the AppContext struct.
func ListInstances(appCtx *app.AppContext) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(appCtx.Logger, 1, "ListInstances()")

	client, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	ctx := context.Background()
	instances, err := FetchInstances(ctx, client, appCtx)
	if err != nil {
		return fmt.Errorf("fetching instances: %w", err)
	}

	// Display instance information
	displayInstances(instances)

	return nil
}

// displayInstances prints formatted instance information to the console.
func displayInstances(instances []Instance) {
	fmt.Println("\nInstances:")
	for _, inst := range instances {
		fmt.Println()
		fmt.Println("Name:", inst.Name)
		fmt.Println("ID:", inst.ID)
		fmt.Printf("Private IP: %s AD: %s\tFD: %s\tRegion: %s\n",
			inst.IP, inst.Placement.AvailabilityDomain, inst.Placement.FaultDomain, inst.Placement.Region)
		fmt.Printf("Shape: %s\tMemory: %d GB\tvCPUs: %d\n",
			inst.Shape, int(inst.Resources.MemoryGB), inst.Resources.VCPUs)
		fmt.Println("State:", inst.State)
		fmt.Println("Created:", inst.CreatedAt)
		fmt.Println("Subnet ID: ", inst.SubnetID)
	}
}

// FindInstances searches for instances in the OCI compartment matching the given name pattern.
// It uses the pre-initialized compute and network clients from the AppContext struct.
// Parameters:
// - application: The application with all clients, logger, and resolved IDs.
// - namePattern: The pattern used to match instance names.
// - showImageDetails: A flag indicating whether to include image details in the output.
// Returns an error if the operation fails.
func FindInstances(application *app.AppContext, namePattern string, showImageDetails bool) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(application.Logger, 1, "FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)

	client, err := oci.NewComputeClient(application.Provider)
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	ctx := context.Background()
	instances, err := FetchInstances(ctx, client, application)
	if err != nil {
		return fmt.Errorf("fetching instances: %w", err)
	}

	// Filter instances by name pattern
	var matchedInstances []Instance
	for _, inst := range instances {
		if strings.Contains(strings.ToLower(inst.Name), strings.ToLower(namePattern)) {
			matchedInstances = append(matchedInstances, inst)
		}
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

	displayInstances(matchedInstances)
	return nil
}

// FetchInstances retrieves all running instances from the specified compartment
// and enriches them with network information.
func FetchInstances(ctx context.Context, computeClient core.ComputeClient, appCtx *app.AppContext) ([]Instance, error) {
	// Create a VNets client
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return nil, fmt.Errorf("creating network client: %w", err)
	}

	// Fetch basic instance information
	instances, err := listInstances(ctx, computeClient, appCtx.CompartmentID)
	if err != nil {
		return nil, err
	}

	// Enrich instances with network information
	for i := range instances {
		if err := enrichInstanceWithNetworkInfo(ctx, computeClient, networkClient, &instances[i], appCtx.CompartmentID); err != nil {
			return nil, fmt.Errorf("enriching instance %s with network info: %w", instances[i].ID, err)
		}
	}

	return instances, nil
}

// listInstances retrieves all running instances from the specified compartment.
func listInstances(ctx context.Context, client core.ComputeClient, compartmentID string) ([]Instance, error) {
	var instances []Instance
	page := ""

	for {
		resp, err := client.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &compartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing instances: %w", err)
		}

		for _, oc := range resp.Items {
			inst := mapToInstance(oc)
			instances = append(instances, inst)
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return instances, nil
}

// enrichInstanceWithNetworkInfo adds network-related information to the instance.
func enrichInstanceWithNetworkInfo(
	ctx context.Context,
	computeClient core.ComputeClient,
	networkClient core.VirtualNetworkClient,
	instance *Instance,
	compartmentID string,
) error {
	// Get the VNIC attachments for this instance
	vaResp, err := computeClient.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: &compartmentID,
		InstanceId:    &instance.ID,
	})
	if err != nil {
		return fmt.Errorf("listing VNIC attachments: %w", err)
	}

	// Find the primary VNIC
	for _, attachment := range vaResp.Items {
		vnicResp, err := networkClient.GetVnic(ctx, core.GetVnicRequest{
			VnicId: attachment.VnicId,
		})
		if err != nil {
			return fmt.Errorf("getting VNIC details: %w", err)
		}

		// Use the primary VNIC for the instance's network information
		if vnicResp.IsPrimary != nil && *vnicResp.IsPrimary {
			instance.IP = *vnicResp.PrivateIp
			instance.SubnetID = *vnicResp.SubnetId
			return nil
		}
	}

	return nil
}

// mapToInstance transforms the OCI SDK Instance into our local Instance type.
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
		Shape:           *oc.Shape,
		ImageID:         *oc.ImageId,
		SubnetID:        "", // to be filled later
		State:           oc.LifecycleState,
		CreatedAt:       *oc.TimeCreated,
		OperatingSystem: "", // added per TODO
	}
}
