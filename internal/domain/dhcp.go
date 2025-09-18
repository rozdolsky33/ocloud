package domain

// DhcpOptions represents a DHCP options in the domain layer.
type DhcpOptions struct {
	OCID           string
	DisplayName    string
	CustomDNS      string
	LifecycleState string
}
