package bastion

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
)

// CanReach checks if the provided bastion can reach a target identified by target VCN and/or Subnet.
// The logic is intentionally simple/minimal:
// - If targetSubnetID is provided, we fetch it and compare its VCN ID with bastion.TargetVcnID.
// - Else if targetVcnID is provided, we compare it directly with bastion.TargetVcnID.
// - If neither targetVcnID nor targetSubnetID is provided, we cannot determine reachability.
func (s *Service) CanReach(ctx context.Context, b Bastion, targetVcnID string, targetSubnetID string) (bool, string) {
	if b.TargetVcnID == "" {
		return false, "Selected Bastion is not configured with a target VCN."
	}

	if targetSubnetID != "" {
		subnet, err := s.getSubnetDetails(ctx, targetSubnetID)
		if err != nil {
			return false, fmt.Sprintf("Unable to verify reachability: failed to fetch target subnet: %v", err)
		}
		if vcnMatches(b.TargetVcnID, subnet) {
			return true, "Bastion target VCN matches the target subnet's VCN."
		}
		return false, fmt.Sprintf("Bastion target VCN %s does not match target subnet's VCN %s", b.TargetVcnID, safeVcnID(subnet))
	}

	// Fall back to VCN comparison if available.
	if targetVcnID != "" {
		if b.TargetVcnID == targetVcnID {
			return true, "Bastion target VCN matches the target VCN."
		}
		return false, fmt.Sprintf("Bastion target VCN %s does not match target VCN %s", b.TargetVcnID, targetVcnID)
	}

	return false, "Target network details are unavailable; cannot verify reachability."
}

// getSubnetDetails retrieves subnet details from OCI.
// This is a temporary helper method until session management is refactored.
func (s *Service) getSubnetDetails(ctx context.Context, subnetID string) (*core.Subnet, error) {
	resp, err := s.networkClient.GetSubnet(ctx, core.GetSubnetRequest{
		SubnetId: &subnetID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting subnet details: %w", err)
	}
	return &resp.Subnet, nil
}

// vcnMatches checks if the provided subnet's VCN ID matches the specified bastion VCN ID. Returns true if they match.
func vcnMatches(bastionVcnID string, subnet *core.Subnet) bool {
	if subnet == nil || subnet.VcnId == nil {
		return false
	}
	return bastionVcnID == *subnet.VcnId
}

// safeVcnID returns the VCN ID of the provided subnet, or an empty string if the subnet is nil or has no VCN ID.
func safeVcnID(subnet *core.Subnet) string {
	if subnet == nil || subnet.VcnId == nil {
		return ""
	}
	return *subnet.VcnId
}
