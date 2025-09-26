package gateway

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"golang.org/x/sync/errgroup"
)

type Adapter struct {
	client core.VirtualNetworkClient
}

// NewAdapter creates an Adapter that wraps the provided core.VirtualNetworkClient.
// The returned Adapter uses the client to perform OCI virtual network operations.
func NewAdapter(client core.VirtualNetworkClient) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) GatewaysSummary(ctx context.Context, compartmentID, vcnID string) (vcn.Gateways, error) {
	var out vcn.Gateways

	// Service ID -> Name cache for SGW services
	svcNames := make(map[string]string)

	eg, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex

	// 1)--------------------------------------------- Internet Gateway ------------------------------------------------
	eg.Go(func() error {
		req := core.ListInternetGatewaysRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}
		resp, err := a.client.ListInternetGateways(ctx, req)
		if err != nil {
			return fmt.Errorf("list IGWs: %w", err)
		}
		name := "—"
		for _, igw := range resp.Items {
			// show the first enabled gateway name
			if igw.IsEnabled != nil && *igw.IsEnabled {
				if igw.DisplayName != nil && *igw.DisplayName != "" {
					name = fmt.Sprintf("%s (present)", *igw.DisplayName)
				} else {
					name = "present"
				}
				break
			}
		}
		mu.Lock()
		out.InternetGateway = name
		mu.Unlock()
		return nil
	})

	// 2)-------------------------------------------- NAT Gateway ------------------------------------------------------
	eg.Go(func() error {
		req := core.ListNatGatewaysRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}
		resp, err := a.client.ListNatGateways(ctx, req)
		if err != nil {
			return fmt.Errorf("list NATs: %w", err)
		}
		name := "—"
		for _, nat := range resp.Items {
			if nat.DisplayName != nil && *nat.DisplayName != "" {
				name = fmt.Sprintf("%s (present)", *nat.DisplayName)
			} else {
				name = "present"
			}
			break
		}
		mu.Lock()
		out.NatGateway = name
		mu.Unlock()
		return nil
	})

	// 3)----------------------------- Service Gateway (+attached services pretty names) -------------------------------
	eg.Go(func() error {
		// first, cache regional services so we can map IDs -> names
		sr, err := a.client.ListServices(ctx, core.ListServicesRequest{})
		if err != nil {
			return fmt.Errorf("list services: %w", err)
		}
		for _, s := range sr.Items {
			if s.Id != nil {
				label := "-"
				if s.Description != nil && *s.Description != "" {
					label = *s.Description
				}
				svcNames[*s.Id] = label
			}
		}

		// now list SGWs on this VCN
		req := core.ListServiceGatewaysRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}
		resp, err := a.client.ListServiceGateways(ctx, req)
		if err != nil {
			return fmt.Errorf("list SGWs: %w", err)
		}
		if len(resp.Items) == 0 {
			mu.Lock()
			out.ServiceGateway = "—"
			mu.Unlock()
			return nil
		}
		sgw := resp.Items[0]
		name := "-"
		if sgw.DisplayName != nil && *sgw.DisplayName != "" {
			name = *sgw.DisplayName
		}
		var attached []string
		for _, e := range sgw.Services {
			if e.ServiceId != nil {
				if n, ok := svcNames[*e.ServiceId]; ok && n != "" {
					attached = append(attached, n)
				}
			}
		}
		s := name
		if len(attached) > 0 {
			s = fmt.Sprintf("%s (%s)", name, strings.Join(attached, ", "))
		}
		mu.Lock()
		out.ServiceGateway = s
		mu.Unlock()
		return nil
	})

	// 4)------------------------------------------ DRG attachment -----------------------------------------------------
	eg.Go(func() error {
		req := core.ListDrgAttachmentsRequest{
			VcnId: &vcnID,
		}
		resp, err := a.client.ListDrgAttachments(ctx, req)
		if err != nil {
			return fmt.Errorf("list DRG attachments: %w", err)
		}
		if len(resp.Items) == 0 {
			mu.Lock()
			out.Drg = "—"
			mu.Unlock()
			return nil
		}
		att := resp.Items[0]
		status := "attached"
		if att.LifecycleState != "" {
			status = strings.ToLower(string(att.LifecycleState))
		}
		name := "drg"
		if att.DrgId != nil {
			drg, err := a.client.GetDrg(ctx, core.GetDrgRequest{DrgId: att.DrgId})
			if err == nil && drg.DisplayName != nil && *drg.DisplayName != "" {
				name = *drg.DisplayName
			}
		}
		mu.Lock()
		out.Drg = fmt.Sprintf("%s (%s)", name, status)
		mu.Unlock()
		return nil
	})

	// 5)------------------------------------ LPG peers: lpg-name → peer-vcn-name --------------------------------------
	eg.Go(func() error {
		req := core.ListLocalPeeringGatewaysRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}
		resp, err := a.client.ListLocalPeeringGateways(ctx, req)
		if err != nil {
			return fmt.Errorf("list LPGs: %w", err)
		}

		var peers []string
		for _, lpg := range resp.Items {
			lpgName := "-"
			if lpg.DisplayName != nil {
				lpgName = *lpg.DisplayName
			}

			// if we have a peer LPG, fetch its VCN name
			if lpg.PeerId != nil {
				peer, err := a.client.GetLocalPeeringGateway(ctx, core.GetLocalPeeringGatewayRequest{LocalPeeringGatewayId: lpg.PeerId})
				if err == nil && peer.LocalPeeringGateway.VcnId != nil {
					vcn, err := a.client.GetVcn(ctx, core.GetVcnRequest{VcnId: peer.LocalPeeringGateway.VcnId})
					if err == nil && vcn.DisplayName != nil {
						peers = append(peers, fmt.Sprintf("%s → %s", lpgName, *vcn.DisplayName))
						continue
					}
				}
				peers = append(peers, fmt.Sprintf("%s → <peer>", lpgName))
			}
		}
		mu.Lock()
		out.LocalPeeringPeers = peers
		mu.Unlock()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return vcn.Gateways{}, err
	}
	return out, nil
}
