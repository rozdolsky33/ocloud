package mapping

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

func TestNewDomainVCNFromAttrs_Basic(t *testing.T) {
	id := "ocid1.vcn.oc1..vcn"
	name := "main-vcn"
	cidr := "10.0.0.0/16"
	state := core.VcnLifecycleStateAvailable
	created := common.SDKTime{Time: time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)}

	oci := core.Vcn{
		Id:             &id,
		DisplayName:    &name,
		CidrBlocks:     []string{cidr},
		LifecycleState: state,
		TimeCreated:    &created,
	}
	attrs := NewVCNAttributesFromOCIVCN(oci)
	dm := NewDomainVCNFromAttrs(attrs)

	if dm.OCID != id {
		t.Errorf("OCID mismatch: %s", dm.OCID)
	}
	if dm.DisplayName != name {
		t.Errorf("DisplayName mismatch: %s", dm.DisplayName)
	}
	if len(dm.CidrBlocks) != 1 || dm.CidrBlocks[0] != cidr {
		t.Errorf("CidrBlocks mismatch: %v", dm.CidrBlocks)
	}
	if dm.LifecycleState != string(state) {
		t.Errorf("LifecycleState mismatch: %s", dm.LifecycleState)
	}
	if dm.TimeCreated.IsZero() {
		t.Errorf("TimeCreated should be set")
	}
}
