package subnet

import "context"

// Subnet represents a subnet in the domain layer.
type Subnet struct {
	OCID            string
	DisplayName     string
	LifecycleState  string
	CidrBlock       string
	Public          bool
	RouteTableID    string
	SecurityListIDs []string
	NSGIDs          []string
}

type SubnetRepository interface {
	GetSubnet(ctx context.Context, ocid string) (*Subnet, error)
	ListSubnets(ctx context.Context, compartmentID string) ([]Subnet, error)
}
