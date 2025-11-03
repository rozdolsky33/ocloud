package instance

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/mapping"
)

const (
	defaultMaxRetries     = 5
	defaultInitialBackoff = 1 * time.Second
	defaultMaxBackoff     = 32 * time.Second
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

// GetEnrichedInstance fetches a single instance by OCID and enriches it with network and image details.
func (a *Adapter) GetEnrichedInstance(ctx context.Context, instanceID string) (*domain.Instance, error) {
	resp, err := a.computeClient.GetInstance(ctx, core.GetInstanceRequest{InstanceId: &instanceID})
	if err != nil {
		return nil, fmt.Errorf("getting instance from OCI: %w", err)
	}
	dm, err := a.enrichAndMapInstance(ctx, resp.Instance)
	if err != nil {
		return nil, err
	}
	return dm, nil
}

// ListInstances fetches all instances in a compartment.
func (a *Adapter) ListInstances(ctx context.Context, compartmentID string) ([]domain.Instance, error) {
	var allInstances []domain.Instance
	var page *string

	for {
		resp, err := a.computeClient.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing instances from OCI: %w", err)
		}
		for _, item := range resp.Items {
			allInstances = append(allInstances, *mapping.NewDomainInstanceFromAttrs(mapping.NewInstanceAttributesFromOCIInstance(item)))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return allInstances, nil
}

// ListEnrichedInstances fetches all instances in a compartment and enriches them with network and image details.
func (a *Adapter) ListEnrichedInstances(ctx context.Context, compartmentID string) ([]domain.Instance, error) {
	var allInstances []core.Instance
	var page *string

	for {
		resp, err := a.computeClient.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId: &compartmentID,
			Page:          page,
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

			dm := mapping.NewDomainInstanceFromAttrs(mapping.NewInstanceAttributesFromOCIInstance(ociInstance))

			if err := a.enrichDomainInstance(ctx, dm, ociInstance); err != nil {
				errChan <- err
				return
			}

			domainInstances[i] = *dm
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
func (a *Adapter) enrichAndMapInstance(ctx context.Context, ociInstance core.Instance) (*domain.Instance, error) {
	dm := mapping.NewDomainInstanceFromAttrs(mapping.NewInstanceAttributesFromOCIInstance(ociInstance))
	if err := a.enrichDomainInstance(ctx, dm, ociInstance); err != nil {
		return dm, err
	}

	return dm, nil
}

// enrichDomainInstance enriches the given domain instance with network, subnet, VCN, route table and image details.
func (a *Adapter) enrichDomainInstance(ctx context.Context, dm *domain.Instance, ociInstance core.Instance) error {
	vnic, err := a.getPrimaryVnic(ctx, *ociInstance.Id, *ociInstance.CompartmentId)
	if err != nil {
		return fmt.Errorf("enriching instance %s with network: %w", dm.OCID, err)
	}
	if vnic != nil {
		vnicAttrs := mapping.NewVnicAttributesFromOCIVnic(*vnic)
		dm.PrimaryIP = *vnicAttrs.PrivateIp
		dm.SubnetID = *vnicAttrs.SubnetId
		if vnicAttrs.HostnameLabel != nil {
			dm.Hostname = *vnicAttrs.HostnameLabel
		}
		dm.PrivateDNSEnabled = vnicAttrs.SkipSourceDestCheck == nil || !*vnicAttrs.SkipSourceDestCheck

		// Extract NSG IDs from VNIC and resolve to names
		dm.NsgIDs = vnicAttrs.NsgIds
		if len(vnicAttrs.NsgIds) > 0 {
			dm.NsgNames = make([]string, 0, len(vnicAttrs.NsgIds))
			for _, nsgID := range vnicAttrs.NsgIds {
				nsgName, err := a.getNsgName(ctx, nsgID)
				if err != nil {
					return fmt.Errorf("enriching instance %s with NSG name for %s: %w", dm.OCID, nsgID, err)
				}
				dm.NsgNames = append(dm.NsgNames, nsgName)
			}
		}

		subnet, err := a.getSubnet(ctx, *vnicAttrs.SubnetId)
		if err != nil {
			return fmt.Errorf("enriching instance %s with subnet: %w", dm.OCID, err)
		}
		if subnet != nil {
			subnetAttrs := mapping.NewSubnetAttributesFromOCISubnet(*subnet)
			dm.SubnetName = *subnetAttrs.DisplayName
			dm.VcnID = *subnet.VcnId

			// Extract Security List IDs from subnet and resolve to names
			dm.SecurityListIDs = subnetAttrs.SecurityListIds
			if len(subnetAttrs.SecurityListIds) > 0 {
				dm.SecurityListNames = make([]string, 0, len(subnetAttrs.SecurityListIds))
				for _, slID := range subnetAttrs.SecurityListIds {
					slName, err := a.getSecurityListName(ctx, slID)
					if err != nil {
						return fmt.Errorf("enriching instance %s with security list name for %s: %w", dm.OCID, slID, err)
					}
					dm.SecurityListNames = append(dm.SecurityListNames, slName)
				}
			}

			if subnetAttrs.RouteTableId != nil {
				dm.RouteTableID = *subnetAttrs.RouteTableId
				rt, err := a.getRouteTable(ctx, *subnetAttrs.RouteTableId)
				if err != nil {
					return fmt.Errorf("enriching instance %s with route table: %w", dm.OCID, err)
				}
				if rt != nil {
					rtAttrs := mapping.NewRouteTableAttributesFromOCIRouteTable(*rt)
					dm.RouteTableName = *rtAttrs.DisplayName
				}
			}
			vcn, err := a.getVcn(ctx, *subnet.VcnId)
			if err != nil {
				return fmt.Errorf("enriching instance %s with vcn: %w", dm.OCID, err)
			}
			if vcn != nil {
				vcnAttrs := mapping.NewVcnAttributesFromOCIVcn(*vcn)
				dm.VcnName = *vcnAttrs.DisplayName
			}
		}
	}

	image, err := a.getImage(ctx, *ociInstance.ImageId)
	if err != nil {
		return fmt.Errorf("enriching instance %s with image: %w", dm.OCID, err)
	}
	if image != nil {
		imageAttrs := mapping.NewImageAttributesFromOCIImage(*image)
		dm.ImageName = *imageAttrs.DisplayName
		dm.ImageOS = *imageAttrs.OperatingSystem
	}
	return nil
}

// getPrimaryVnic finds the primary VNIC for a given instance.
func (a *Adapter) getPrimaryVnic(ctx context.Context, instanceID, compartmentID string) (*core.Vnic, error) {
	var attachments core.ListVnicAttachmentsResponse
	var err error
	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
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

			vnicErr = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
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

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
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

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
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

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
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

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetRouteTable(ctx, core.GetRouteTableRequest{RtId: &rtID})
		return e
	})
	if err != nil {
		return nil, err
	}
	return &resp.RouteTable, nil
}

// getSecurityListName fetches the display name for a security list by OCID.
func (a *Adapter) getSecurityListName(ctx context.Context, securityListID string) (string, error) {
	var resp core.GetSecurityListResponse
	var err error

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetSecurityList(ctx, core.GetSecurityListRequest{SecurityListId: &securityListID})
		return e
	})
	if err != nil {
		return "", err
	}
	if resp.SecurityList.DisplayName != nil {
		return *resp.SecurityList.DisplayName, nil
	}
	return "", nil
}

// getNsgName fetches the display name for a network security group by OCID.
func (a *Adapter) getNsgName(ctx context.Context, nsgID string) (string, error) {
	var resp core.GetNetworkSecurityGroupResponse
	var err error

	err = retryOnRateLimit(ctx, defaultMaxRetries, defaultInitialBackoff, defaultMaxBackoff, func() error {
		var e error
		resp, e = a.networkClient.GetNetworkSecurityGroup(ctx, core.GetNetworkSecurityGroupRequest{NetworkSecurityGroupId: &nsgID})
		return e
	})
	if err != nil {
		return "", err
	}
	if resp.NetworkSecurityGroup.DisplayName != nil {
		return *resp.NetworkSecurityGroup.DisplayName, nil
	}
	return "", nil
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
