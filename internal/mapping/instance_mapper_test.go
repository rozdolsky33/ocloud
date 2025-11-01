package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestNewInstanceAttributesFromOCIInstance_And_ToDomain(t *testing.T) {
	id := "ocid1.instance.oc1..aaa"
	name := "web-1"
	shape := "VM.Standard3.Flex"
	image := "ocid1.image.oc1..img"
	region := "eu-frankfurt-1"
	ad := "AD-1"
	fd := "FD-1"
	ocpus := float32(2)
	mem := float32(16)
	created := time.Now().UTC().Truncate(time.Second)

	inst := core.Instance{
		Id:                 &id,
		DisplayName:        &name,
		LifecycleState:     core.InstanceLifecycleStateRunning,
		Shape:              &shape,
		ImageId:            &image,
		TimeCreated:        &common.SDKTime{Time: created},
		Region:             &region,
		AvailabilityDomain: &ad,
		FaultDomain:        &fd,
		ShapeConfig: &core.InstanceShapeConfig{
			Ocpus:       &ocpus,
			MemoryInGBs: &mem,
		},
		FreeformTags: map[string]string{"role": "web"},
		DefinedTags:  map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewInstanceAttributesFromOCIInstance(inst)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, core.InstanceLifecycleStateRunning, attrs.State)
	require.Equal(t, &shape, attrs.Shape)
	require.Equal(t, &image, attrs.ImageId)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, created.Equal(attrs.TimeCreated.Time))
	require.Equal(t, &region, attrs.Region)
	require.Equal(t, &ad, attrs.AvailabilityDomain)
	require.Equal(t, &fd, attrs.FaultDomain)
	require.NotNil(t, attrs.Vcpus)
	require.InDelta(t, 2, float64(*attrs.Vcpus), 0.001)
	require.NotNil(t, attrs.MemoryInGBs)
	require.InDelta(t, 16, float64(*attrs.MemoryInGBs), 0.001)

	dom := mapping.NewDomainInstanceFromAttrs(attrs)
	require.IsType(t, &domain.Instance{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, string(core.InstanceLifecycleStateRunning), dom.State)
	require.Equal(t, shape, dom.Shape)
	require.Equal(t, image, dom.ImageID)
	require.True(t, created.Equal(dom.TimeCreated))
	require.Equal(t, region, dom.Region)
	require.Equal(t, ad, dom.AvailabilityDomain)
	require.Equal(t, fd, dom.FaultDomain)
	require.Equal(t, 2, dom.VCPUs)
	require.InDelta(t, 16, float64(dom.MemoryGB), 0.001)
	require.Equal(t, map[string]string{"role": "web"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, dom.DefinedTags)
}

func TestNewDomainInstanceFromAttrs_NilValues(t *testing.T) {
	attrs := &mapping.InstanceAttributes{}
	dom := mapping.NewDomainInstanceFromAttrs(attrs)
	require.Equal(t, "", dom.OCID)
	require.Equal(t, "", dom.DisplayName)
	require.Equal(t, "", dom.State)
	require.Equal(t, "", dom.Shape)
	require.Equal(t, "", dom.ImageID)
	require.True(t, dom.TimeCreated.IsZero())
	require.Equal(t, "", dom.Region)
	require.Equal(t, "", dom.AvailabilityDomain)
	require.Equal(t, "", dom.FaultDomain)
	require.Equal(t, 0, dom.VCPUs)
	require.InDelta(t, 0, float64(dom.MemoryGB), 0.001)
	require.Nil(t, dom.FreeformTags)
	require.Nil(t, dom.DefinedTags)
}

func TestNewVnicAttributesFromOCIVnic(t *testing.T) {
	ip := "10.0.0.10"
	subnet := "ocid1.subnet.oc1..xyz"
	host := "web-1"
	skip := true
	nsgIds := []string{"ocid1.nsg.oc1..nsg1", "ocid1.nsg.oc1..nsg2"}

	v := core.Vnic{
		PrivateIp:           &ip,
		SubnetId:            &subnet,
		HostnameLabel:       &host,
		SkipSourceDestCheck: &skip,
		NsgIds:              nsgIds,
	}

	got := mapping.NewVnicAttributesFromOCIVnic(v)
	require.Equal(t, &ip, got.PrivateIp)
	require.Equal(t, &subnet, got.SubnetId)
	require.Equal(t, &host, got.HostnameLabel)
	require.Equal(t, &skip, got.SkipSourceDestCheck)
	require.Equal(t, nsgIds, got.NsgIds)
}

func TestNewVnicAttributesFromOCIVnic_EmptyNsgIds(t *testing.T) {
	ip := "10.0.0.10"
	subnet := "ocid1.subnet.oc1..xyz"

	v := core.Vnic{
		PrivateIp: &ip,
		SubnetId:  &subnet,
		NsgIds:    []string{},
	}

	got := mapping.NewVnicAttributesFromOCIVnic(v)
	require.Equal(t, &ip, got.PrivateIp)
	require.Equal(t, &subnet, got.SubnetId)
	require.NotNil(t, got.NsgIds)
	require.Empty(t, got.NsgIds)
}

func TestNewVnicAttributesFromOCIVnic_NilNsgIds(t *testing.T) {
	ip := "10.0.0.10"
	subnet := "ocid1.subnet.oc1..xyz"

	v := core.Vnic{
		PrivateIp: &ip,
		SubnetId:  &subnet,
		NsgIds:    nil,
	}

	got := mapping.NewVnicAttributesFromOCIVnic(v)
	require.Equal(t, &ip, got.PrivateIp)
	require.Equal(t, &subnet, got.SubnetId)
	require.Nil(t, got.NsgIds)
}
