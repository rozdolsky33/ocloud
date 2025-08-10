package bastion

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/core"
)

// CanReach checks if the provided bastion can reach a target identified by target VCN and/or Subnet.
// The logic is intentionally simple/minimal:
// - If targetSubnetID is provided, we fetch it and compare its VCN ID with bastion.TargetVcnId.
// - Else if targetVcnID is provided, we compare it directly with bastion.TargetVcnId.
// - If neither targetVcnID nor targetSubnetID is provided, we cannot determine reachability.
// Returns ok boolean and a user-friendly reason string.
func (s *Service) CanReach(ctx context.Context, b Bastion, targetVcnID string, targetSubnetID string) (bool, string) {
	if b.TargetVcnId == "" {
		return false, "Selected Bastion is not configured with a target VCN."
	}

	// Prefer subnet if provided, to derive VCN and be precise.
	if targetSubnetID != "" {
		subnet, err := s.fetchSubnetDetails(ctx, targetSubnetID)
		if err != nil {
			return false, fmt.Sprintf("Unable to verify reachability: failed to fetch target subnet: %v", err)
		}
		if vcnMatches(b.TargetVcnId, subnet) {
			return true, "Bastion target VCN matches the target subnet's VCN."
		}
		return false, fmt.Sprintf("Bastion target VCN %s does not match target subnet's VCN %s", b.TargetVcnId, safeVcnID(subnet))
	}

	// Fall back to VCN comparison if available.
	if targetVcnID != "" {
		if b.TargetVcnId == targetVcnID {
			return true, "Bastion target VCN matches the target VCN."
		}
		return false, fmt.Sprintf("Bastion target VCN %s does not match target VCN %s", b.TargetVcnId, targetVcnID)
	}

	return false, "Target network details are unavailable; cannot verify reachability."
}

func vcnMatches(bastionVcnID string, subnet *core.Subnet) bool {
	if subnet == nil || subnet.VcnId == nil {
		return false
	}
	return bastionVcnID == *subnet.VcnId
}

func safeVcnID(subnet *core.Subnet) string {
	if subnet == nil || subnet.VcnId == nil {
		return ""
	}
	return *subnet.VcnId
}
