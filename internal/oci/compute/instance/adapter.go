package instance

import (
	"context"
	"fmt"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// Adapter is an infrastructure-layer adapter for compute instances.
// It implements the domain.InstanceRepository interface.
type Adapter struct {
	computeClient core.ComputeClient
	networkClient core.VirtualNetworkClient
}

// NewAdapter creates a new instance adapter.
func NewAdapter(computeClient core.ComputeClient, networkClient core.VirtualNetworkClient) *Adapter {
	return &Adapter{
		computeClient: computeClient,
		networkClient: networkClient,
	}
}

// ListInstances fetches all running instances in a compartment and enriches them with network and image details.
func (a *Adapter) ListInstances(ctx context.Context, compartmentID string) ([]domain.Instance, error) {
	var allInstances []core.Instance
	var page *string

	for {
		resp, err := a.computeClient.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId:  &compartmentID,
			LifecycleState: core.InstanceLifecycleStateRunning,
			Page:           page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing instances from OCI: %w", err)
		}
		allInstances = append(allInstances, resp.Items...)

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return a.enrichAndMapInstances(ctx, allInstances)
}

// enrichAndMapInstances converts OCI instances to domain models and enriches them with details.
func (a *Adapter) enrichAndMapInstances(ctx context.Context, ociInstances []core.Instance) ([]domain.Instance, error) {
	domainInstances := make([]domain.Instance, len(ociInstances))
	var wg sync.WaitGroup
	errChan := make(chan error, len(ociInstances))

	for i, ociInstance := range ociInstances {
		wg.Add(1)
		go func(i int, ociInstance core.Instance) {
			defer wg.Done()

			dm := domain.Instance{
				OCID:               *ociInstance.Id,
				DisplayName:        *ociInstance.DisplayName,
				State:              string(ociInstance.LifecycleState),
				Shape:              *ociInstance.Shape,
				ImageID:            *ociInstance.ImageId,
				TimeCreated:        ociInstance.TimeCreated.Time,
				Region:             *ociInstance.Region,
				AvailabilityDomain: *ociInstance.AvailabilityDomain,
				FaultDomain:        *ociInstance.FaultDomain,
				VCPUs:              int(*ociInstance.ShapeConfig.Vcpus),
				MemoryGB:           *ociInstance.ShapeConfig.MemoryInGBs,
			}

			vnic, err := a.getPrimaryVnic(ctx, *ociInstance.Id, *ociInstance.CompartmentId)
			if err != nil {
				errChan <- fmt.Errorf("enriching instance %s with network: %w", dm.OCID, err)
				return
			}
			if vnic != nil {
				dm.PrimaryIP = *vnic.PrivateIp
				dm.SubnetID = *vnic.SubnetId
				subnet, err := a.getSubnet(ctx, *vnic.SubnetId)
				if err != nil {
					errChan <- fmt.Errorf("enriching instance %s with subnet: %w", dm.OCID, err)
					return
				}
				if subnet != nil {
					dm.SubnetName = *subnet.DisplayName
					dm.VcnID = *subnet.VcnId
					vcn, err := a.getVcn(ctx, *subnet.VcnId)
					if err != nil {
						errChan <- fmt.Errorf("enriching instance %s with vcn: %w", dm.OCID, err)
						return
					}
					if vcn != nil {
						dm.VcnName = *vcn.DisplayName
					}
				}
			}

			image, err := a.getImage(ctx, *ociInstance.ImageId)
			if err != nil {
				errChan <- fmt.Errorf("enriching instance %s with image: %w", dm.OCID, err)
				return
			}
			if image != nil {
				dm.ImageName = *image.DisplayName
				dm.ImageOS = *image.OperatingSystem
			}

			domainInstances[i] = dm
		}(i, ociInstance)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return nil, err
	}

	return domainInstances, nil
}

// getPrimaryVnic finds the primary VNIC for a given instance.
func (a *Adapter) getPrimaryVnic(ctx context.Context, instanceID, compartmentID string) (*core.Vnic, error) {
	attachments, err := a.computeClient.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
		CompartmentId: &compartmentID,
		InstanceId:    &instanceID,
	})
	if err != nil {
		return nil, err
	}

	for _, attach := range attachments.Items {
		if attach.VnicId != nil {
			resp, err := a.networkClient.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
			if err == nil && resp.Vnic.IsPrimary != nil && *resp.Vnic.IsPrimary {
				return &resp.Vnic, nil
			}
		}
	}
	return nil, nil // No primary VNIC found
}

// getSubnet fetches subnet details.
func (a *Adapter) getSubnet(ctx context.Context, subnetID string) (*core.Subnet, error) {
	resp, err := a.networkClient.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &subnetID,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Subnet, nil
}

// getVcn fetches VCN details.
func (a *Adapter) getVcn(ctx context.Context, vcnID string) (*core.Vcn, error) {
	resp, err := a.networkClient.GetVcn(ctx, core.GetVcnRequest{
		VcnId: &vcnID,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Vcn, nil
}

// getImage fetches image details.
func (a *Adapter) getImage(ctx context.Context, imageID string) (*core.Image, error) {
	resp, err := a.computeClient.GetImage(ctx, core.GetImageRequest{
		ImageId: &imageID,
	})
	if err != nil {
		return nil, err
	}
	return &resp.Image, nil
}
