package instance

import (
	"reflect"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

func sref(s string) *string { return &s }

func TestToBaseDomainInstanceModel(t *testing.T) {
	ad := &Adapter{}
	tm := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	v := 4
	vcpus := &v
	m := float32(16)
	mem := &m
	oci := core.Instance{
		Id:             sref("ocid1.instance.oc1..inst"),
		DisplayName:    sref("i-1"),
		TimeCreated:    &common.SDKTime{Time: tm},
		Shape:          sref("VM.Standard3.Flex"),
		LifecycleState: "RUNNING",
		FaultDomain:    sref("FAULT-DOMAIN-1"),
		ShapeConfig:    &core.InstanceShapeConfig{Vcpus: vcpus, MemoryInGBs: mem},
	}
	d := ad.toBaseDomainInstanceModel(oci)
	expect := domain.Instance{OCID: "ocid1.instance.oc1..inst", DisplayName: "i-1", TimeCreated: tm, Shape: "VM.Standard3.Flex", State: "RUNNING", VCPUs: 4, MemoryGB: 16, FaultDomain: "FAULT-DOMAIN-1"}
	if !reflect.DeepEqual(d, expect) {
		t.Fatalf("toBaseDomainInstanceModel mismatch: got %#v want %#v", d, expect)
	}
}

func TestToEnrichDomainInstanceModel(t *testing.T) {
	ad := &Adapter{}
	tm := time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)
	free := map[string]string{"env": "dev"}
	def := map[string]map[string]interface{}{"ns": {"k": "v"}}
	oci := core.Instance{
		Id:                 sref("ocid1.instance.oc1..inst2"),
		DisplayName:        sref("i-2"),
		LifecycleState:     "STOPPED",
		Shape:              sref("VM.Standard.E4.Flex"),
		ImageId:            sref("ocid1.image.oc1..img"),
		TimeCreated:        &common.SDKTime{Time: tm},
		Region:             sref("eu-frankfurt-1"),
		AvailabilityDomain: sref("vSDF:EU-FRANKFURT-1-AD-1"),
		FaultDomain:        sref("FAULT-DOMAIN-2"),
		ShapeConfig: func() *core.InstanceShapeConfig {
			vi := 2
			mi := float32(8)
			return &core.InstanceShapeConfig{Vcpus: &vi, MemoryInGBs: &mi}
		}(),
		FreeformTags: free,
		DefinedTags:  def,
	}
	d := ad.toEnrichDomainInstanceModel(oci)
	expect := domain.Instance{
		OCID:               "ocid1.instance.oc1..inst2",
		DisplayName:        "i-2",
		State:              "STOPPED",
		Shape:              "VM.Standard.E4.Flex",
		ImageID:            "ocid1.image.oc1..img",
		TimeCreated:        tm,
		Region:             "eu-frankfurt-1",
		AvailabilityDomain: "vSDF:EU-FRANKFURT-1-AD-1",
		FaultDomain:        "FAULT-DOMAIN-2",
		VCPUs:              2,
		MemoryGB:           8,
		FreeformTags:       free,
		DefinedTags:        def,
	}
	if !reflect.DeepEqual(d, expect) {
		t.Fatalf("toEnrichDomainInstanceModel mismatch: got %#v want %#v", d, expect)
	}
}
