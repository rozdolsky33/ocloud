package vcn

// VCNDTO represents a Virtual Cloud Network details needed for summary output.
type VCNDTO struct {
	OCID            string            `json:"OCID"`
	DisplayName     string            `json:"DisplayName"`
	LifecycleState  string            `json:"LifecycleState"`
	CompartmentID   string            `json:"CompartmentID"`
	CompartmentName string            `json:"CompartmentName,omitempty"`
	CidrBlocks      []string          `json:"CidrBlocks"`
	Ipv6Enabled     bool              `json:"Ipv6Enabled"`
	DnsLabel        string            `json:"DnsLabel"`
	DomainName      string            `json:"DomainName"`
	DhcpOptionsID   string            `json:"DhcpOptionsID"`
	DhcpOptionsName string            `json:"DhcpOptionsName,omitempty"`
	DhcpCustomDNS   string            `json:"DhcpCustomDNS,omitempty"`
	TimeCreated     string            `json:"TimeCreated"`
	FreeformTags    map[string]string `json:"FreeformTags,omitempty"`
}
