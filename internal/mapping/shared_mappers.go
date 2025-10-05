package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/core"
	domain_vcn "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
)

// ImageAttributes is a provider-agnostic representation of an image.
// Adapters should populate this from their SDK types and pass to the builder.
type ImageAttributes struct {
	ID                     *string
	DisplayName            *string
	OperatingSystem        *string
	OperatingSystemVersion *string
	LaunchMode             string
	TimeCreated            *time.Time
}

func NewImageAttributesFromOCIImage(i core.Image) *ImageAttributes {
	return &ImageAttributes{
		ID:                     i.Id,
		DisplayName:            i.DisplayName,
		OperatingSystem:        i.OperatingSystem,
		OperatingSystemVersion: i.OperatingSystemVersion,
		LaunchMode:             string(i.LaunchMode),
		TimeCreated:            &i.TimeCreated.Time,
	}
}

type SubnetAttributes struct {
	OCID                   *string
	DisplayName            *string
	LifecycleState         core.SubnetLifecycleStateEnum
	CidrBlock              *string
	RouteTableId           *string
	ProhibitPublicIpOnVnic *bool
	SecurityListIds        []string
	VcnId                  *string
}

func NewSubnetAttributesFromOCISubnet(s core.Subnet) *SubnetAttributes {
	return &SubnetAttributes{
		OCID:                   s.Id,
		DisplayName:            s.DisplayName,
		LifecycleState:         s.LifecycleState,
		CidrBlock:              s.CidrBlock,
		RouteTableId:           s.RouteTableId,
		ProhibitPublicIpOnVnic: s.ProhibitPublicIpOnVnic,
		SecurityListIds:        s.SecurityListIds,
		VcnId:                  s.VcnId,
	}
}

func NewDomainSubnetFromAttrs(s *SubnetAttributes) *domain_vcn.Subnet {
	var ocid, displayName, lifecycleState, cidrBlock, routeTableId string
	var public bool

	if s.OCID != nil {
		ocid = *s.OCID
	}
	if s.DisplayName != nil {
		displayName = *s.DisplayName
	}
	if s.LifecycleState != "" {
		lifecycleState = string(s.LifecycleState)
	}
	if s.CidrBlock != nil {
		cidrBlock = *s.CidrBlock
	}
	if s.RouteTableId != nil {
		routeTableId = *s.RouteTableId
	}
	if s.ProhibitPublicIpOnVnic != nil {
		public = !*s.ProhibitPublicIpOnVnic
	}

	return &domain_vcn.Subnet{
		OCID:            ocid,
		DisplayName:     displayName,
		LifecycleState:  lifecycleState,
		CidrBlock:       cidrBlock,
		Public:          public,
		RouteTableID:    routeTableId,
		SecurityListIDs: s.SecurityListIds,
	}
}

type RouteTableAttributes struct {
	OCID           *string
	DisplayName    *string
	LifecycleState core.RouteTableLifecycleStateEnum
}

func NewRouteTableAttributesFromOCIRouteTable(rt core.RouteTable) *RouteTableAttributes {
	return &RouteTableAttributes{
		OCID:           rt.Id,
		DisplayName:    rt.DisplayName,
		LifecycleState: rt.LifecycleState,
	}
}

func NewDomainRouteTableFromAttrs(rt *RouteTableAttributes) *domain_vcn.RouteTable {
	var ocid, displayName, lifecycleState string

	if rt.OCID != nil {
		ocid = *rt.OCID
	}
	if rt.DisplayName != nil {
		displayName = *rt.DisplayName
	}
	if rt.LifecycleState != "" {
		lifecycleState = string(rt.LifecycleState)
	}

	return &domain_vcn.RouteTable{
		OCID:           ocid,
		DisplayName:    displayName,
		LifecycleState: lifecycleState,
	}
}

type VcnAttributes struct {
	DisplayName *string
}

func NewVcnAttributesFromOCIVcn(v core.Vcn) *VcnAttributes {
	return &VcnAttributes{
		DisplayName: v.DisplayName,
	}
}
