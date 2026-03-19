package dynamicgroup

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/identitydomains"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/rozdolsky33/ocloud/internal/oci"
)

// Adapter is an infrastructure-layer adapter that implements the domain.DynamicGroupRepository interface.
type Adapter struct {
	client   identity.IdentityClient
	provider common.ConfigurationProvider
}

// NewDynamicGroupAdapter creates a new adapter for interacting with OCI dynamic groups.
func NewDynamicGroupAdapter(client identity.IdentityClient, provider common.ConfigurationProvider) *Adapter {
	return &Adapter{
		client:   client,
		provider: provider,
	}
}

// GetDynamicGroup retrieves a single dynamic group by its OCID.
func (a *Adapter) GetDynamicGroup(ctx context.Context, ocid string) (*domain.DynamicGroup, error) {
	// Try traditional first
	resp, err := a.client.GetDynamicGroup(ctx, identity.GetDynamicGroupRequest{
		DynamicGroupId: &ocid,
	})
	if err == nil {
		return mapping.NewDomainDynamicGroupFromAttrs(mapping.NewDynamicGroupAttributesFromOCI(resp.DynamicGroup)), nil
	}

	// Try domains
	tenancyID, _ := a.provider.TenancyOCID()
	domainsResp, err := a.client.ListDomains(ctx, identity.ListDomainsRequest{
		CompartmentId: &tenancyID,
	})
	if err == nil {
		for _, d := range domainsResp.Items {
			if d.Url == nil {
				continue
			}
			idcsClient, err := oci.NewIdentityDomainsClient(a.provider, *d.Url)
			if err != nil {
				continue
			}
			resp, err := idcsClient.GetDynamicResourceGroup(ctx, identitydomains.GetDynamicResourceGroupRequest{
				DynamicResourceGroupId: &ocid,
			})
			if err == nil {
				return mapping.NewDomainDynamicGroupFromAttrs(mapping.NewDynamicGroupAttributesFromIDCS(resp.DynamicResourceGroup, *d.Url)), nil
			}
		}
	}

	return nil, fmt.Errorf("getting dynamic group from OCI: %w", err)
}

// ListDynamicGroups retrieves all dynamic groups under a given compartment (usually tenancy).
func (a *Adapter) ListDynamicGroups(ctx context.Context, compartmentID string) ([]domain.DynamicGroup, error) {
	var dynamicGroups []domain.DynamicGroup

	domainsResp, err := a.client.ListDomains(ctx, identity.ListDomainsRequest{
		CompartmentId: &compartmentID,
	})

	if err == nil && len(domainsResp.Items) > 0 {
		// We have Identity Domains
		for _, d := range domainsResp.Items {
			if d.Url == nil {
				continue
			}
			idcsClient, err := oci.NewIdentityDomainsClient(a.provider, *d.Url)
			if err != nil {
				continue
			}

			page := ""
			for {
				resp, err := idcsClient.ListDynamicResourceGroups(ctx, identitydomains.ListDynamicResourceGroupsRequest{
					Attributes: common.String("id,ocid,displayName,description,matchingRule,meta"),
					Page:       &page,
				})
				if err != nil {
					break
				}

				for _, item := range resp.DynamicResourceGroups.Resources {
					dgAttrs := mapping.NewDynamicGroupAttributesFromIDCS(item, *d.Url)
					// Prepend domain name to display name if there are multiple domains
					if len(domainsResp.Items) > 1 {
						name := fmt.Sprintf("[%s] %s", *d.DisplayName, *dgAttrs.Name)
						dgAttrs.Name = &name
					}
					dynamicGroups = append(dynamicGroups, *mapping.NewDomainDynamicGroupFromAttrs(dgAttrs))
				}

				if resp.OpcNextPage == nil {
					break
				}
				page = *resp.OpcNextPage
			}
		}
		if len(dynamicGroups) > 0 {
			return dynamicGroups, nil
		}
	}

	// 2. Fallback to traditional Dynamic Groups
	page := ""
	for {
		resp, err := a.client.ListDynamicGroups(ctx, identity.ListDynamicGroupsRequest{
			CompartmentId: &compartmentID,
			Page:          &page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing dynamic groups from OCI: %w", err)
		}

		for _, item := range resp.Items {
			dynamicGroups = append(dynamicGroups, *mapping.NewDomainDynamicGroupFromAttrs(mapping.NewDynamicGroupAttributesFromOCI(item)))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = *resp.OpcNextPage
	}

	return dynamicGroups, nil
}
