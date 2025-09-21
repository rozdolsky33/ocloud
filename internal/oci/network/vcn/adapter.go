package vcn

import (
	"context"
	"fmt"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
)

// Adapter provides access to VCN-related OCI APIs.
// It is infra-layer and should be used by the service layer.
type Adapter struct {
	client core.VirtualNetworkClient
}

// NewAdapter creates a new adapter instance.
func NewAdapter(client core.VirtualNetworkClient) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GetEnrichedVcn(ctx context.Context, vcnID string) (domain.VCN, error) {
	resp, err := a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: &vcnID})
	if err != nil {
		return domain.VCN{}, fmt.Errorf("getting VCN from OCI: %w", err)
	}
	m := toDomainVCNModel(resp.Vcn)
	return m, nil
}

// ListVcns lists all VCNs in a given compartment.
func (a *Adapter) ListVcns(ctx context.Context, compartmentID string) ([]domain.VCN, error) {
	req := core.ListVcnsRequest{CompartmentId: &compartmentID}
	var out []domain.VCN
	for {
		resp, err := a.client.ListVcns(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("listing VCNs from OCI: %w", err)
		}
		for _, v := range resp.Items {
			out = append(out, toDomainVCNModel(v))
		}
		if resp.OpcNextPage == nil {
			break
		}
		req.Page = resp.OpcNextPage
	}
	return out, nil
}

// ListEnrichedVcns lists VCNs and enriches them with all related resources in parallel.
func (a *Adapter) ListEnrichedVcns(ctx context.Context, compartmentID string) ([]domain.VCN, error) {
	vcns, err := a.ListVcns(ctx, compartmentID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(vcns)*4)

	for i := range vcns {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Gateways
			vcns[i].Gateways, err = a.ListInternetGateways(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			nats, err := a.ListNatGateways(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].Gateways = append(vcns[i].Gateways, nats...)
			sgws, err := a.ListServiceGateways(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].Gateways = append(vcns[i].Gateways, sgws...)
			lpgs, err := a.ListLocalPeeringGateways(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].Gateways = append(vcns[i].Gateways, lpgs...)
			drg, err := a.ListDrgAttachments(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].Gateways = append(vcns[i].Gateways, drg...)

			// Subnets
			vcns[i].Subnets, err = a.ListSubnets(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}

			// Route tables, Security lists, NSGs
			vcns[i].RouteTables, err = a.ListRouteTables(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].SecurityLists, err = a.ListSecurityLists(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}
			vcns[i].NSGs, err = a.ListNetworkSecurityGroups(ctx, compartmentID, vcns[i].OCID)
			if err != nil {
				errCh <- err
			}

			// DHCP options (default for the VCN)
			if vcns[i].DhcpOptionsID != "" {
				dhcp, derr := a.GetDhcpOptions(ctx, vcns[i].DhcpOptionsID)
				if derr != nil {
					errCh <- derr
				} else {
					vcns[i].DhcpOptions = dhcp
				}
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return vcns, nil
}

func toDomainVCNModel(v core.Vcn) domain.VCN {
	return domain.VCN{
		OCID:           *v.Id,
		DisplayName:    *v.DisplayName,
		LifecycleState: string(v.LifecycleState),
		CompartmentID:  *v.CompartmentId,
		DnsLabel:       *v.DnsLabel,
		DomainName:     *v.VcnDomainName,
		CidrBlocks:     cloneStrings(v.CidrBlocks),
		Ipv6Enabled:    len(v.Ipv6CidrBlocks) > 0,
		DhcpOptionsID:  *v.DefaultDhcpOptionsId,
		TimeCreated:    v.TimeCreated.Time,
		FreeformTags:   v.FreeformTags,
		DefinedTags:    v.DefinedTags,
	}
}

func cloneStrings(in []string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func (a *Adapter) ListInternetGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListInternetGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListInternetGateways(ctx, req)
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Internet"})
	}
	return gateways, nil
}

func (a *Adapter) ListNatGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListNatGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListNatGateways(ctx, req)
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "NAT"})
	}
	return gateways, nil
}

func (a *Adapter) ListServiceGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListServiceGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListServiceGateways(ctx, req)
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Service"})
	}
	return gateways, nil
}

func (a *Adapter) ListLocalPeeringGateways(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListLocalPeeringGatewaysRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListLocalPeeringGateways(ctx, req)
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "Local Peering"})
	}
	return gateways, nil
}

func (a *Adapter) ListDrgAttachments(ctx context.Context, compartmentID, vcnID string) ([]domain.Gateway, error) {
	req := core.ListDrgAttachmentsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListDrgAttachments(ctx, req)
	if err != nil {
		return nil, err
	}
	var gateways []domain.Gateway
	for _, item := range resp.Items {
		gateways = append(gateways, domain.Gateway{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState), Type: "DRG"})
	}
	return gateways, nil
}

func (a *Adapter) ListRouteTables(ctx context.Context, compartmentID, vcnID string) ([]domain.RouteTable, error) {
	req := core.ListRouteTablesRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListRouteTables(ctx, req)
	if err != nil {
		return nil, err
	}
	var rts []domain.RouteTable
	for _, item := range resp.Items {
		rts = append(rts, domain.RouteTable{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return rts, nil
}

func (a *Adapter) ListSecurityLists(ctx context.Context, compartmentID, vcnID string) ([]domain.SecurityList, error) {
	req := core.ListSecurityListsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListSecurityLists(ctx, req)
	if err != nil {
		return nil, err
	}
	var sls []domain.SecurityList
	for _, item := range resp.Items {
		sls = append(sls, domain.SecurityList{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return sls, nil
}

func (a *Adapter) ListNetworkSecurityGroups(ctx context.Context, compartmentID, vcnID string) ([]domain.NSG, error) {
	req := core.ListNetworkSecurityGroupsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListNetworkSecurityGroups(ctx, req)
	if err != nil {
		return nil, err
	}
	var nsgs []domain.NSG
	for _, item := range resp.Items {
		nsgs = append(nsgs, domain.NSG{OCID: *item.Id, DisplayName: *item.DisplayName, LifecycleState: string(item.LifecycleState)})
	}
	return nsgs, nil
}

func (a *Adapter) GetDhcpOptions(ctx context.Context, dhcpID string) (domain.DhcpOptions, error) {
	req := core.GetDhcpOptionsRequest{DhcpId: &dhcpID}
	resp, err := a.client.GetDhcpOptions(ctx, req)
	if err != nil {
		return domain.DhcpOptions{}, err
	}
	return domain.DhcpOptions{OCID: *resp.Id, DisplayName: *resp.DisplayName, LifecycleState: string(resp.LifecycleState), DomainNameType: ""}, nil
}

func (a *Adapter) ListSubnets(ctx context.Context, compartmentID, vcnID string) ([]domain.Subnet, error) {
	req := core.ListSubnetsRequest{CompartmentId: &compartmentID, VcnId: &vcnID}
	resp, err := a.client.ListSubnets(ctx, req)
	if err != nil {
		return nil, err
	}
	var subnets []domain.Subnet
	for _, item := range resp.Items {
		var id, name, cidr, rtID string
		if item.Id != nil {
			id = *item.Id
		}
		if item.DisplayName != nil {
			name = *item.DisplayName
		}
		if item.CidrBlock != nil {
			cidr = *item.CidrBlock
		}
		if item.RouteTableId != nil {
			rtID = *item.RouteTableId
		}
		public := item.ProhibitPublicIpOnVnic == nil || !*item.ProhibitPublicIpOnVnic
		var slIDs []string
		if item.SecurityListIds != nil {
			slIDs = item.SecurityListIds
		}
		subnets = append(subnets, domain.Subnet{OCID: id, DisplayName: name, LifecycleState: string(item.LifecycleState), CidrBlock: cidr, Public: public, RouteTableID: rtID, SecurityListIDs: slIDs})
	}
	return subnets, nil
}
