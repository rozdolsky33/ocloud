package instance

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
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

func (a *Adapter) GetInstance(ctx context.Context, id string) (*domain.Instance, error) {
	resp, err := a.computeClient.GetInstance(ctx, core.GetInstanceRequest{
		InstanceId: &id,
	})

	if err != nil {
		return nil, fmt.Errorf("getting instance from OCI: %w", err)
	}

	inst := a.toDomainModel(resp.Instance)
	return &inst, nil

}

// ListInstances fetches all running instances in a compartment.
func (a *Adapter) ListInstances(ctx context.Context, compartmentID string) ([]domain.Instance, error) {
	var allInstances []domain.Instance
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
		for _, item := range resp.Items {
			allInstances = append(allInstances, a.toDomainModel(item))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return allInstances, nil
}

// ListEnrichedInstances fetches all running instances in a compartment and enriches them with network and image details.
func (a *Adapter) ListEnrichedInstances(ctx context.Context, compartmentID string) ([]domain.Instance, error) {
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

// GetEnrichedInstance fetches a single instance by OCID and enriches it with network and image details.
func (a *Adapter) GetEnrichedInstance(ctx context.Context, instanceID string) (*domain.Instance, error) {
	resp, err := a.computeClient.GetInstance(ctx, core.GetInstanceRequest{InstanceId: &instanceID})
	if err != nil {
		return nil, fmt.Errorf("getting instance from OCI: %w", err)
	}
	// Enrich and map the single instance
	dm, err := a.enrichAndMapInstance(ctx, resp.Instance)
	if err != nil {
		return nil, err
	}
	return &dm, nil
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
				FreeformTags:       ociInstance.FreeformTags,
				DefinedTags:        ociInstance.DefinedTags,
			}

			vnic, err := a.getPrimaryVnic(ctx, *ociInstance.Id, *ociInstance.CompartmentId)
			if err != nil {
				errChan <- fmt.Errorf("enriching instance %s with network: %w", dm.OCID, err)
				return
			}
			if vnic != nil {
				dm.PrimaryIP = *vnic.PrivateIp
				dm.SubnetID = *vnic.SubnetId
				if vnic.HostnameLabel != nil {
					dm.Hostname = *vnic.HostnameLabel
				}
				dm.PrivateDNSEnabled = vnic.SkipSourceDestCheck == nil || !*vnic.SkipSourceDestCheck
				subnet, err := a.getSubnet(ctx, *vnic.SubnetId)
				if err != nil {
					errChan <- fmt.Errorf("enriching instance %s with subnet: %w", dm.OCID, err)
					return
				}
				if subnet != nil {
					dm.SubnetName = *subnet.DisplayName
					dm.VcnID = *subnet.VcnId
					if subnet.RouteTableId != nil {
						dm.RouteTableID = *subnet.RouteTableId
						rt, err := a.getRouteTable(ctx, *subnet.RouteTableId)
						if err != nil {
							errChan <- fmt.Errorf("enriching instance %s with route table: %w", dm.OCID, err)
							return
						}
						if rt != nil {
							dm.RouteTableName = *rt.DisplayName
						}
					}
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

// enrichAndMapInstance converts a single OCI instance to a domain model and enriches it with details.
func (a *Adapter) enrichAndMapInstance(ctx context.Context, ociInstance core.Instance) (domain.Instance, error) {
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
		FreeformTags:       ociInstance.FreeformTags,
		DefinedTags:        ociInstance.DefinedTags,
	}

	vnic, err := a.getPrimaryVnic(ctx, *ociInstance.Id, *ociInstance.CompartmentId)
	if err != nil {
		return dm, fmt.Errorf("enriching instance %s with network: %w", dm.OCID, err)
	}
	if vnic != nil {
		dm.PrimaryIP = *vnic.PrivateIp
		dm.SubnetID = *vnic.SubnetId
		if vnic.HostnameLabel != nil {
			dm.Hostname = *vnic.HostnameLabel
		}
		dm.PrivateDNSEnabled = vnic.SkipSourceDestCheck == nil || !*vnic.SkipSourceDestCheck
		subnet, err := a.getSubnet(ctx, *vnic.SubnetId)
		if err != nil {
			return dm, fmt.Errorf("enriching instance %s with subnet: %w", dm.OCID, err)
		}
		if subnet != nil {
			dm.SubnetName = *subnet.DisplayName
			dm.VcnID = *subnet.VcnId
			if subnet.RouteTableId != nil {
				dm.RouteTableID = *subnet.RouteTableId
				rt, err := a.getRouteTable(ctx, *subnet.RouteTableId)
				if err != nil {
					return dm, fmt.Errorf("enriching instance %s with route table: %w", dm.OCID, err)
				}
				if rt != nil {
					dm.RouteTableName = *rt.DisplayName
				}
			}
			vcn, err := a.getVcn(ctx, *subnet.VcnId)
			if err != nil {
				return dm, fmt.Errorf("enriching instance %s with vcn: %w", dm.OCID, err)
			}
			if vcn != nil {
				dm.VcnName = *vcn.DisplayName
			}
		}
	}

	image, err := a.getImage(ctx, *ociInstance.ImageId)
	if err != nil {
		return dm, fmt.Errorf("enriching instance %s with image: %w", dm.OCID, err)
	}
	if image != nil {
		dm.ImageName = *image.DisplayName
		dm.ImageOS = *image.OperatingSystem
	}

	return dm, nil
}

// getPrimaryVnic finds the primary VNIC for a given instance.
func (a *Adapter) getPrimaryVnic(ctx context.Context, instanceID, compartmentID string) (*core.Vnic, error) {
	var attachments core.ListVnicAttachmentsResponse
	var err error
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second

	// Retry ListVnicAttachments with a unified helper
	err = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
		var e error
		attachments, e = a.computeClient.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
			CompartmentId: &compartmentID,
			InstanceId:    &instanceID,
		})
		return e
	})
	if err != nil {
		return nil, err
	}

	for _, attach := range attachments.Items {
		if attach.VnicId != nil {
			var resp core.GetVnicResponse
			var vnicErr error

			// Retry GetVnic using a unified helper; if it still fails, move on to the next VNIC
			vnicErr = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
				var e error
				resp, e = a.networkClient.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})
				return e
			})
			if vnicErr == nil {
				if resp.Vnic.IsPrimary != nil && *resp.Vnic.IsPrimary {
					return &resp.Vnic, nil
				}
			}
		}
	}
	return nil, nil
}

// getSubnet fetches subnet details.
func (a *Adapter) getSubnet(ctx context.Context, subnetID string) (*core.Subnet, error) {
	var resp core.GetSubnetResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second

	err = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetSubnet(ctx, core.GetSubnetRequest{SubnetId: &subnetID})
		return e
	})
	if err != nil {
		return nil, err
	}
	return &resp.Subnet, nil
}

// getVcn fetches VCN details.
func (a *Adapter) getVcn(ctx context.Context, vcnID string) (*core.Vcn, error) {
	var resp core.GetVcnResponse
	var err error
	// Retry
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second

	err = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
		return e
	})
	if err != nil {
		return nil, err
	}
	return &resp.Vcn, nil
}

// getImage fetches image details.
func (a *Adapter) getImage(ctx context.Context, imageID string) (*core.Image, error) {
	var resp core.GetImageResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second

	err = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
		var e error
		resp, e = a.computeClient.GetImage(ctx, core.GetImageRequest{ImageId: &imageID})
		return e
	})
	if err != nil {
		return nil, err
	}
	return &resp.Image, nil
}

// getRouteTable fetches route table details.
func (a *Adapter) getRouteTable(ctx context.Context, rtID string) (*core.RouteTable, error) {
	var resp core.GetRouteTableResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second

	err = retryOnRateLimit(ctx, maxRetries, initialBackoff, maxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetRouteTable(ctx, core.GetRouteTableRequest{RtId: &rtID})
		return e
	})
	if err != nil {
		return nil, err
	}
	return &resp.RouteTable, nil
}

// retryOnRateLimit retries the provided operation when OCI responds with HTTP 429 rate limited.
// It applies exponential backoff between retries and preserves the original behavior and error messages.
func retryOnRateLimit(ctx context.Context, maxRetries int, initialBackoff, maxBackoff time.Duration, op func() error) error {
	backoff := initialBackoff
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := op()
		if err == nil {
			return nil
		}

		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			if attempt == maxRetries-1 {
				return fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		return err
	}
	return nil
}

// toDomainModel converts an OCI SDK image object to our application's domain model.
func (a *Adapter) toDomainModel(inst core.Instance) domain.Instance {
	return domain.Instance{
		OCID:        *inst.Id,
		DisplayName: *inst.DisplayName,
		TimeCreated: inst.TimeCreated.Time,
		Shape:       *inst.Shape,
		State:       string(inst.LifecycleState),
	}
}
