package compute

import (
	"context"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/oci"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ListInstances lists all instances in the configured compartment using the provided application.
// It uses the pre-initialized compute client from the AppContext struct.
func ListInstances(appCtx *app.AppContext) error {
	// Use VerboseInfo to ensure debug logs work with shorthand flags
	logger.VerboseInfo(appCtx.Logger, 1, "ListInstances()")
	client, err := oci.NewComputeClient(appCtx.Provider)
	if err != nil {
		return err
	}
	ctx := context.Background()

	instances, err := FetchInstances(ctx, client, appCtx)
	if err != nil {
		return err
	}

	fmt.Println("\nInstances:")
	fmt.Println(instances)

	// Use the pre-initialized compute client from the AppContext struct
	// No need to create a new client
	// Use application.ComputeClient to list instances

	return nil
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

	// Use the pre-initialized compute and network clients from the AppContext struct
	// No need to create new clients

	// Use application.ComputeClient and application.NetworkClient to find instances
	// ...

	return nil
}

func FetchInstances(ctx context.Context, computeClient core.ComputeClient, appCtx *app.AppContext) ([]Instance, error) {
	// create a VNets client
	networkClient, err := oci.NewNetworkClient(appCtx.Provider)
	if err != nil {
		return nil, err
	}
	var instances []Instance
	page := ""
	for {
		resp, err := computeClient.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &appCtx.CompartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           &page,
		})
		if err != nil {
			return nil, err
		}

		for _, oc := range resp.Items {
			inst := mapToInstance(oc)
			fmt.Println()
			instances = append(instances, inst)
			// 1) get the VNIC attachment(s) for this instance
			vaResp, err := computeClient.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
				CompartmentId: &appCtx.CompartmentID,
				InstanceId:    oc.Id,
			})
			if err != nil {
				return nil, err
			}
			// 2) find the primary VNIC and call GetVnic
			// after ListVnicAttachments into vaResp.Items...

			for _, attachment := range vaResp.Items {
				vnicResp, err := networkClient.GetVnic(ctx, core.GetVnicRequest{
					VnicId: attachment.VnicId,
				})
				if err != nil {
					return nil, err
				}
				// Now vnicResp.IsPrimary exists
				if vnicResp.IsPrimary != nil && *vnicResp.IsPrimary {
					inst.IP = *vnicResp.PrivateIp
					inst.SubnetID = *vnicResp.SubnetId
					break
				}
			}
			instances = append(instances, inst)
			fmt.Println("Name:", *oc.DisplayName)
			fmt.Println("ID:", *oc.Id)
			fmt.Println("Private IP:", inst.IP, "AD:", *oc.AvailabilityDomain, "\t", "FD:", *oc.FaultDomain, "\t", "Region:", *oc.Region)
			fmt.Println("Shape:", *oc.Shape, "\t", "Memory:", *oc.ShapeConfig.MemoryInGBs, "GB", "\t vCPUs:", *oc.ShapeConfig.Vcpus)
			fmt.Println("State:", oc.LifecycleState)
			fmt.Println("Created:", oc.TimeCreated)
			fmt.Println("Subnet ID: ", inst.SubnetID)

		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}
	return instances, nil
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
		ImageID:         *oc.ImageId,
		SubnetID:        "", // to be filled later
		State:           oc.LifecycleState,
		CreatedAt:       *oc.TimeCreated,
		OperatingSystem: "", // added per TODO
	}
}
