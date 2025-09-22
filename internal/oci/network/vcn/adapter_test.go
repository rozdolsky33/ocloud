package vcn

import (
	"reflect"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
)

func strptr(s string) *string { return &s }

func TestCloneStrings(t *testing.T) {
	// nil input -> nil output
	if got := cloneStrings(nil); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}

	in := []string{"a", "b", "c"}
	got := cloneStrings(in)
	if !reflect.DeepEqual(got, in) {
		t.Fatalf("expected clone equal to input, got %v", got)
	}
	// ensure it's a copy, not the same backing array
	got[0] = "z"
	if in[0] != "a" {
		t.Fatalf("expected original slice to remain unchanged, got %v", in)
	}
}

func TestToDomainVCNModel(t *testing.T) {
	now := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	cidrs := []string{"10.0.0.0/16", "10.1.0.0/16"}
	ipv6Cidrs := []string{"2600:1f18:abcd::/56"}
	freeform := map[string]string{"env": "dev"}
	defined := map[string]map[string]interface{}{"ns": {"k": "v"}}

	v := core.Vcn{
		Id:                   strptr("ocid1.vcn.oc1..test"),
		DisplayName:          strptr("test-vcn"),
		LifecycleState:       "AVAILABLE",
		CompartmentId:        strptr("ocid1.compartment.oc1..comp"),
		DnsLabel:             strptr("mynet"),
		VcnDomainName:        strptr("mynet.oraclevcn.com"),
		CidrBlocks:           cidrs,
		Ipv6CidrBlocks:       ipv6Cidrs,
		DefaultDhcpOptionsId: strptr("ocid1.dhcp.oc1..dhcp"),
		TimeCreated:          &common.SDKTime{Time: now},
		FreeformTags:         freeform,
		DefinedTags:          defined,
	}

	d := toDomainVCNModel(v)

	expected := domain.VCN{
		OCID:           "ocid1.vcn.oc1..test",
		DisplayName:    "test-vcn",
		LifecycleState: "AVAILABLE",
		CompartmentID:  "ocid1.compartment.oc1..comp",
		DnsLabel:       "mynet",
		DomainName:     "mynet.oraclevcn.com",
		CidrBlocks:     cidrs,
		Ipv6Enabled:    true,
		DhcpOptionsID:  "ocid1.dhcp.oc1..dhcp",
		TimeCreated:    now,
		FreeformTags:   freeform,
		DefinedTags:    defined,
	}

	// Compare selected fields (slices and maps by value)
	if d.OCID != expected.OCID || d.DisplayName != expected.DisplayName || d.LifecycleState != expected.LifecycleState || d.CompartmentID != expected.CompartmentID || d.DnsLabel != expected.DnsLabel || d.DomainName != expected.DomainName || d.DhcpOptionsID != expected.DhcpOptionsID || !reflect.DeepEqual(d.CidrBlocks, expected.CidrBlocks) || d.Ipv6Enabled != expected.Ipv6Enabled || !d.TimeCreated.Equal(expected.TimeCreated) || !reflect.DeepEqual(d.FreeformTags, expected.FreeformTags) || !reflect.DeepEqual(d.DefinedTags, expected.DefinedTags) {
		t.Fatalf("toDomainVCNModel mismatch:\n got: %#v\nwant: %#v", d, expected)
	}

	// Ensure slices are cloned, not aliased
	if &d.CidrBlocks[0] == &v.CidrBlocks[0] {
		// addresses in Go for string elements may be reused; mutate to verify decoupling
		// Instead, mutate original slice and ensure domain slice unchanged length/content
		v.CidrBlocks[0] = "changed"
		if d.CidrBlocks[0] == "changed" {
			t.Fatalf("expected CidrBlocks to be a clone, not referencing original slice")
		}
	}
}

func TestToDomainVCNModel_NoIPv6(t *testing.T) {
	v := core.Vcn{
		Id:                   strptr("ocid1.vcn.oc1..noipv6"),
		DisplayName:          strptr("noipv6"),
		LifecycleState:       "AVAILABLE",
		CompartmentId:        strptr("ocid1.compartment.oc1..comp"),
		DnsLabel:             strptr("net"),
		VcnDomainName:        strptr("net.oraclevcn.com"),
		CidrBlocks:           []string{"10.0.0.0/16"},
		Ipv6CidrBlocks:       nil,
		DefaultDhcpOptionsId: strptr("ocid1.dhcp.oc1..dhcp"),
		TimeCreated:          &common.SDKTime{Time: time.Unix(0, 0)},
	}
	d := toDomainVCNModel(v)
	if d.Ipv6Enabled {
		t.Fatalf("expected Ipv6Enabled=false when no IPv6 CIDRs, got true")
	}
}
