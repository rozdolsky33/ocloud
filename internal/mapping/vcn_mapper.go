package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
)

type VCNAttributes struct {
	OCID                 *string
	DisplayName          *string
	LifecycleState       core.VcnLifecycleStateEnum
	CompartmentID        *string
	CidrBlocks           []string
	Ipv6CidrBlocks       []string
	DnsLabel             *string
	VcnDomainName        *string
	DefaultDhcpOptionsId *string
	TimeCreated          *common.SDKTime
	FreeformTags         map[string]string
	DefinedTags          map[string]map[string]interface{}
}

func NewVCNAttributesFromOCIVCN(v core.Vcn) *VCNAttributes {
	return &VCNAttributes{
		OCID:                 v.Id,
		DisplayName:          v.DisplayName,
		LifecycleState:       v.LifecycleState,
		CompartmentID:        v.CompartmentId,
		CidrBlocks:           v.CidrBlocks,
		Ipv6CidrBlocks:       v.Ipv6CidrBlocks,
		DnsLabel:             v.DnsLabel,
		VcnDomainName:        v.VcnDomainName,
		DefaultDhcpOptionsId: v.DefaultDhcpOptionsId,
		TimeCreated:          v.TimeCreated,
		FreeformTags:         v.FreeformTags,
		DefinedTags:          v.DefinedTags,
	}
}

func NewDomainVCNFromAttrs(v *VCNAttributes) *domain.VCN {
	var ocid, displayName, lifecycleState, compartmentID, dnsLabel, domainName, dhcpOptionsID string
	var timeCreated time.Time

	if v.OCID != nil {
		ocid = *v.OCID
	}
	if v.DisplayName != nil {
		displayName = *v.DisplayName
	}
	if v.LifecycleState != "" {
		lifecycleState = string(v.LifecycleState)
	}
	if v.CompartmentID != nil {
		compartmentID = *v.CompartmentID
	}
	if v.DnsLabel != nil {
		dnsLabel = *v.DnsLabel
	}
	if v.VcnDomainName != nil {
		domainName = *v.VcnDomainName
	}
	if v.DefaultDhcpOptionsId != nil {
		dhcpOptionsID = *v.DefaultDhcpOptionsId
	}
	if v.TimeCreated != nil {
		timeCreated = v.TimeCreated.Time
	}

	return &domain.VCN{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
		CompartmentID:  compartmentID,
		CidrBlocks:     v.CidrBlocks,
		Ipv6Enabled:    len(v.Ipv6CidrBlocks) > 0,
		DnsLabel:       dnsLabel,
		DomainName:     domainName,
		DhcpOptionsID:  dhcpOptionsID,
		TimeCreated:    timeCreated,
		FreeformTags:   v.FreeformTags,
		DefinedTags:    v.DefinedTags,
	}
}

type GatewayAttributes struct {
	OCID           *string
	DisplayName    *string
	LifecycleState string
	Type           string
}

func NewGatewayAttributesFromOCIInternetGateway(ig core.InternetGateway) *GatewayAttributes {
	return &GatewayAttributes{
		OCID:           ig.Id,
		DisplayName:    ig.DisplayName,
		LifecycleState: string(ig.LifecycleState),
		Type:           "Internet",
	}
}

func NewGatewayAttributesFromOCINatGateway(ng core.NatGateway) *GatewayAttributes {
	return &GatewayAttributes{
		OCID:           ng.Id,
		DisplayName:    ng.DisplayName,
		LifecycleState: string(ng.LifecycleState),
		Type:           "NAT",
	}
}

func NewGatewayAttributesFromOCIServiceGateway(sg core.ServiceGateway) *GatewayAttributes {
	return &GatewayAttributes{
		OCID:           sg.Id,
		DisplayName:    sg.DisplayName,
		LifecycleState: string(sg.LifecycleState),
		Type:           "Service",
	}
}

func NewGatewayAttributesFromOCILocalPeeringGateway(lpg core.LocalPeeringGateway) *GatewayAttributes {
	return &GatewayAttributes{
		OCID:           lpg.Id,
		DisplayName:    lpg.DisplayName,
		LifecycleState: string(lpg.LifecycleState),
		Type:           "Local Peering",
	}
}

func NewGatewayAttributesFromOCIDrgAttachment(drg core.DrgAttachment) *GatewayAttributes {
	return &GatewayAttributes{
		OCID:           drg.Id,
		DisplayName:    drg.DisplayName,
		LifecycleState: string(drg.LifecycleState),
		Type:           "DRG",
	}
}

func NewDomainGatewayFromAttrs(g *GatewayAttributes) *domain.Gateway {
	var ocid, displayName, lifecycleState, typeName string

	if g.OCID != nil {
		ocid = *g.OCID
	}
	if g.DisplayName != nil {
		displayName = *g.DisplayName
	}
	if g.LifecycleState != "" {
		lifecycleState = g.LifecycleState
	}
	if g.Type != "" {
		typeName = g.Type
	}

	return &domain.Gateway{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
		Type:           typeName,
	}
}

type SecurityListAttributes struct {
	OCID           *string
	DisplayName    *string
	LifecycleState core.SecurityListLifecycleStateEnum
}

func NewSecurityListAttributesFromOCISecurityList(sl core.SecurityList) *SecurityListAttributes {
	return &SecurityListAttributes{
		OCID:           sl.Id,
		DisplayName:    sl.DisplayName,
		LifecycleState: sl.LifecycleState,
	}
}

func NewDomainSecurityListFromAttrs(sl *SecurityListAttributes) *domain.SecurityList {
	var ocid, displayName, lifecycleState string

	if sl.OCID != nil {
		ocid = *sl.OCID
	}
	if sl.DisplayName != nil {
		displayName = *sl.DisplayName
	}
	if sl.LifecycleState != "" {
		lifecycleState = string(sl.LifecycleState)
	}

	return &domain.SecurityList{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
	}
}

type NSGAttributes struct {
	OCID           *string
	DisplayName    *string
	LifecycleState core.NetworkSecurityGroupLifecycleStateEnum
}

func NewNSGAttributesFromOCINSG(nsg core.NetworkSecurityGroup) *NSGAttributes {
	return &NSGAttributes{
		OCID:           nsg.Id,
		DisplayName:    nsg.DisplayName,
		LifecycleState: nsg.LifecycleState,
	}
}

func NewDomainNSGFromAttrs(nsg *NSGAttributes) *domain.NSG {
	var ocid, displayName, lifecycleState string

	if nsg.OCID != nil {
		ocid = *nsg.OCID
	}
	if nsg.DisplayName != nil {
		displayName = *nsg.DisplayName
	}
	if nsg.LifecycleState != "" {
		lifecycleState = string(nsg.LifecycleState)
	}

	return &domain.NSG{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
	}
}

type DhcpOptionsAttributes struct {
	OCID           *string
	DisplayName    *string
	LifecycleState core.DhcpOptionsLifecycleStateEnum
	DomainNameType string
}

func NewDhcpOptionsAttributesFromOCIDhcpOptions(do core.DhcpOptions) *DhcpOptionsAttributes {
	return &DhcpOptionsAttributes{
		OCID:           do.Id,
		DisplayName:    do.DisplayName,
		LifecycleState: do.LifecycleState,
	}
}

func NewDomainDhcpOptionsFromAttrs(do *DhcpOptionsAttributes) *domain.DhcpOptions {
	var ocid, displayName, lifecycleState string

	if do.OCID != nil {
		ocid = *do.OCID
	}
	if do.DisplayName != nil {
		displayName = *do.DisplayName
	}
	if do.LifecycleState != "" {
		lifecycleState = string(do.LifecycleState)
	}

	return &domain.DhcpOptions{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
		DomainNameType: do.DomainNameType,
	}
}
