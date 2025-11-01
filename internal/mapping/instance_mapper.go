package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

type InstanceAttributes struct {
	OCID               *string
	DisplayName        *string
	State              core.InstanceLifecycleStateEnum
	Shape              *string
	ImageId            *string
	TimeCreated        *common.SDKTime
	Region             *string
	AvailabilityDomain *string
	FaultDomain        *string
	Vcpus              *float32
	MemoryInGBs        *float32
	FreeformTags       map[string]string
	DefinedTags        map[string]map[string]interface{}
}

func NewInstanceAttributesFromOCIInstance(i core.Instance) *InstanceAttributes {
	return &InstanceAttributes{
		OCID:               i.Id,
		DisplayName:        i.DisplayName,
		State:              i.LifecycleState,
		Shape:              i.Shape,
		ImageId:            i.ImageId,
		TimeCreated:        i.TimeCreated,
		Region:             i.Region,
		AvailabilityDomain: i.AvailabilityDomain,
		FaultDomain:        i.FaultDomain,
		Vcpus:              i.ShapeConfig.Ocpus,
		MemoryInGBs:        i.ShapeConfig.MemoryInGBs,
		FreeformTags:       i.FreeformTags,
		DefinedTags:        i.DefinedTags,
	}
}

func NewDomainInstanceFromAttrs(i *InstanceAttributes) *domain.Instance {
	var ocid, displayName, state, shape, imageId, region, availabilityDomain, faultDomain string
	var timeCreated time.Time
	var vcpus int
	var memoryGB float32

	if i.OCID != nil {
		ocid = *i.OCID
	}
	if i.DisplayName != nil {
		displayName = *i.DisplayName
	}
	if i.State != "" {
		state = string(i.State)
	}
	if i.Shape != nil {
		shape = *i.Shape
	}
	if i.ImageId != nil {
		imageId = *i.ImageId
	}
	if i.TimeCreated != nil {
		timeCreated = i.TimeCreated.Time
	}
	if i.Region != nil {
		region = *i.Region
	}
	if i.AvailabilityDomain != nil {
		availabilityDomain = *i.AvailabilityDomain
	}
	if i.FaultDomain != nil {
		faultDomain = *i.FaultDomain
	}
	if i.Vcpus != nil {
		vcpus = int(*i.Vcpus)
	}
	if i.MemoryInGBs != nil {
		memoryGB = *i.MemoryInGBs
	}

	return &domain.Instance{
		OCID:               ocid,
		DisplayName:        displayName,
		State:              state,
		Shape:              shape,
		ImageID:            imageId,
		TimeCreated:        timeCreated,
		Region:             region,
		AvailabilityDomain: availabilityDomain,
		FaultDomain:        faultDomain,
		VCPUs:              vcpus,
		MemoryGB:           memoryGB,
		FreeformTags:       i.FreeformTags,
		DefinedTags:        i.DefinedTags,
	}
}

type VnicAttributes struct {
	PrivateIp           *string
	SubnetId            *string
	HostnameLabel       *string
	SkipSourceDestCheck *bool
	NsgIds              []string
}

func NewVnicAttributesFromOCIVnic(v core.Vnic) *VnicAttributes {
	return &VnicAttributes{
		PrivateIp:           v.PrivateIp,
		SubnetId:            v.SubnetId,
		HostnameLabel:       v.HostnameLabel,
		SkipSourceDestCheck: v.SkipSourceDestCheck,
		NsgIds:              v.NsgIds,
	}
}
