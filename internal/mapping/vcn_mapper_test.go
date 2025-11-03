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

func TestNewSecurityListAttributesFromOCISecurityList(t *testing.T) {
	id := "ocid1.securitylist.oc1..sl1"
	name := "default-security-list"
	state := core.SecurityListLifecycleStateAvailable

	sl := core.SecurityList{
		Id:             &id,
		DisplayName:    &name,
		LifecycleState: state,
	}

	attrs := NewSecurityListAttributesFromOCISecurityList(sl)
	if attrs.OCID == nil || *attrs.OCID != id {
		t.Errorf("OCID mismatch: expected %s, got %v", id, attrs.OCID)
	}
	if attrs.DisplayName == nil || *attrs.DisplayName != name {
		t.Errorf("DisplayName mismatch: expected %s, got %v", name, attrs.DisplayName)
	}
	if attrs.LifecycleState != state {
		t.Errorf("LifecycleState mismatch: expected %v, got %v", state, attrs.LifecycleState)
	}
}

func TestNewDomainSecurityListFromAttrs(t *testing.T) {
	id := "ocid1.securitylist.oc1..sl1"
	name := "web-security-list"
	state := core.SecurityListLifecycleStateAvailable

	attrs := &SecurityListAttributes{
		OCID:           &id,
		DisplayName:    &name,
		LifecycleState: state,
	}

	dm := NewDomainSecurityListFromAttrs(attrs)
	if dm.OCID != id {
		t.Errorf("OCID mismatch: expected %s, got %s", id, dm.OCID)
	}
	if dm.DisplayName != name {
		t.Errorf("DisplayName mismatch: expected %s, got %s", name, dm.DisplayName)
	}
	if dm.LifecycleState != string(state) {
		t.Errorf("LifecycleState mismatch: expected %s, got %s", string(state), dm.LifecycleState)
	}
}

func TestNewDomainSecurityListFromAttrs_NilValues(t *testing.T) {
	attrs := &SecurityListAttributes{}
	dm := NewDomainSecurityListFromAttrs(attrs)

	if dm.OCID != "" {
		t.Errorf("OCID should be empty, got %s", dm.OCID)
	}
	if dm.DisplayName != "" {
		t.Errorf("DisplayName should be empty, got %s", dm.DisplayName)
	}
	if dm.LifecycleState != "" {
		t.Errorf("LifecycleState should be empty, got %s", dm.LifecycleState)
	}
}

func TestNewNSGAttributesFromOCINSG(t *testing.T) {
	id := "ocid1.networksecuritygroup.oc1..nsg1"
	name := "app-tier-nsg"
	state := core.NetworkSecurityGroupLifecycleStateAvailable

	nsg := core.NetworkSecurityGroup{
		Id:             &id,
		DisplayName:    &name,
		LifecycleState: state,
	}

	attrs := NewNSGAttributesFromOCINSG(nsg)
	if attrs.OCID == nil || *attrs.OCID != id {
		t.Errorf("OCID mismatch: expected %s, got %v", id, attrs.OCID)
	}
	if attrs.DisplayName == nil || *attrs.DisplayName != name {
		t.Errorf("DisplayName mismatch: expected %s, got %v", name, attrs.DisplayName)
	}
	if attrs.LifecycleState != state {
		t.Errorf("LifecycleState mismatch: expected %v, got %v", state, attrs.LifecycleState)
	}
}

func TestNewDomainNSGFromAttrs(t *testing.T) {
	id := "ocid1.networksecuritygroup.oc1..nsg1"
	name := "database-nsg"
	state := core.NetworkSecurityGroupLifecycleStateAvailable

	attrs := &NSGAttributes{
		OCID:           &id,
		DisplayName:    &name,
		LifecycleState: state,
	}

	dm := NewDomainNSGFromAttrs(attrs)
	if dm.OCID != id {
		t.Errorf("OCID mismatch: expected %s, got %s", id, dm.OCID)
	}
	if dm.DisplayName != name {
		t.Errorf("DisplayName mismatch: expected %s, got %s", name, dm.DisplayName)
	}
	if dm.LifecycleState != string(state) {
		t.Errorf("LifecycleState mismatch: expected %s, got %s", string(state), dm.LifecycleState)
	}
}

func TestNewDomainNSGFromAttrs_NilValues(t *testing.T) {
	attrs := &NSGAttributes{}
	dm := NewDomainNSGFromAttrs(attrs)

	if dm.OCID != "" {
		t.Errorf("OCID should be empty, got %s", dm.OCID)
	}
	if dm.DisplayName != "" {
		t.Errorf("DisplayName should be empty, got %s", dm.DisplayName)
	}
	if dm.LifecycleState != "" {
		t.Errorf("LifecycleState should be empty, got %s", dm.LifecycleState)
	}
}
