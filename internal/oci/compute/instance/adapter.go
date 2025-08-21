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

// getPrimaryVnic finds the primary VNIC for a given instance.
func (a *Adapter) getPrimaryVnic(ctx context.Context, instanceID, compartmentID string) (*core.Vnic, error) {
	var attachments core.ListVnicAttachmentsResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second
	backoff := initialBackoff

	// Retry loop for ListVnicAttachments
	for attempt := 0; attempt < maxRetries; attempt++ {
		attachments, err = a.computeClient.ListVnicAttachments(ctx, core.ListVnicAttachmentsRequest{
			CompartmentId: &compartmentID,
			InstanceId:    &instanceID,
		})

		// If no error or not a rate limit error, break the retry loop
		if err == nil {
			break
		}

		// Check if it's a rate limit error (HTTP 429)
		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			// If this is the last attempt, return the error
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}

			// Sleep with exponential backoff before retrying
			time.Sleep(backoff)

			// Increase backoff for the next attempt, but don't exceed maxBackoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// For non-rate-limit errors, return immediately
		return nil, err
	}

	for _, attach := range attachments.Items {
		if attach.VnicId != nil {
			var resp core.GetVnicResponse
			var vnicErr error

			// Retry for GetVnic with the same parameters as ListVnicAttachments
			retryBackoff := initialBackoff
			for vnicAttempt := 0; vnicAttempt < maxRetries; vnicAttempt++ {
				resp, vnicErr = a.networkClient.GetVnic(ctx, core.GetVnicRequest{VnicId: attach.VnicId})

				// If no error or not a rate limit error, break the retry loop
				if vnicErr == nil {
					if resp.Vnic.IsPrimary != nil && *resp.Vnic.IsPrimary {
						return &resp.Vnic, nil
					}
					break
				}

				// Check if it's a rate limit error (HTTP 429)
				if serviceErr, ok := common.IsServiceError(vnicErr); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
					// If this is the last attempt, continue to the next VNIC
					if vnicAttempt == maxRetries-1 {
						break
					}

					// Sleep with exponential backoff before retrying
					time.Sleep(retryBackoff)

					// Increase backoff for next attempt, but don't exceed maxBackoff
					retryBackoff *= 2
					if retryBackoff > maxBackoff {
						retryBackoff = maxBackoff
					}

					continue
				}

				// For non-rate-limit errors, break and try the next VNIC
				break
			}
		}
	}
	return nil, nil // No primary VNIC found
}

// getSubnet fetches subnet details.
func (a *Adapter) getSubnet(ctx context.Context, subnetID string) (*core.Subnet, error) {
	var resp core.GetSubnetResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second
	backoff := initialBackoff

	// Retry loop for GetSubnet
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = a.networkClient.GetSubnet(ctx, core.GetSubnetRequest{
			SubnetId: &subnetID,
		})

		// If no error or not a rate limit error, break the retry loop
		if err == nil {
			break
		}

		// Check if it's a rate limit error (HTTP 429)
		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			// If this is the last attempt, return the error
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}

			// Sleep with exponential backoff before retrying
			time.Sleep(backoff)

			// Increase backoff for next attempt, but don't exceed maxBackoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// For non-rate-limit errors, return immediately
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &resp.Subnet, nil
}

// getVcn fetches VCN details.
func (a *Adapter) getVcn(ctx context.Context, vcnID string) (*core.Vcn, error) {
	var resp core.GetVcnResponse
	var err error

	// Retry parameters
	maxRetries := 5
	initialBackoff := 1 * time.Second
	maxBackoff := 32 * time.Second
	backoff := initialBackoff

	// Retry loop for GetVcn
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = a.networkClient.GetVcn(ctx, core.GetVcnRequest{
			VcnId: &vcnID,
		})

		// If no error or not a rate limit error, break the retry loop
		if err == nil {
			break
		}

		// Check if it's a rate limit error (HTTP 429)
		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			// If this is the last attempt, return the error
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}

			// Sleep with exponential backoff before retrying
			time.Sleep(backoff)

			// Increase backoff for next attempt, but don't exceed maxBackoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// For non-rate-limit errors, return immediately
		return nil, err
	}

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
	backoff := initialBackoff

	// Retry loop for GetImage
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = a.computeClient.GetImage(ctx, core.GetImageRequest{
			ImageId: &imageID,
		})

		// If no error or not a rate limit error, break the retry loop
		if err == nil {
			break
		}

		// Check if it's a rate limit error (HTTP 429)
		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			// If this is the last attempt, return the error
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}

			// Sleep with exponential backoff before retrying
			time.Sleep(backoff)

			// Increase backoff for the next attempt, but don't exceed maxBackoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// For non-rate-limit errors, return immediately
		return nil, err
	}

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
	backoff := initialBackoff

	// Retry loop for GetRouteTable
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = a.networkClient.GetRouteTable(ctx, core.GetRouteTableRequest{
			RtId: &rtID,
		})

		// If no error or not a rate limit error, break the retry loop
		if err == nil {
			break
		}

		// Check if it's a rate limit error (HTTP 429)
		if serviceErr, ok := common.IsServiceError(err); ok && serviceErr.GetHTTPStatusCode() == http.StatusTooManyRequests {
			// If this is the last attempt, return the error
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("rate limit exceeded after %d retries: %w", maxRetries, err)
			}

			// Sleep with exponential backoff before retrying
			time.Sleep(backoff)

			// Increase backoff for the next attempt, but don't exceed maxBackoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// For non-rate-limit errors, return immediately
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &resp.RouteTable, nil
}
