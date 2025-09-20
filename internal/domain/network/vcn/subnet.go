package vcn

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
