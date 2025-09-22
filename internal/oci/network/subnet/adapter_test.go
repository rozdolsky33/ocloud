package subnet

import (
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/subnet"
)

func sptr(s string) *string { return &s }
func bptr(b bool) *bool     { return &b }

func TestToDomainModel_Subnet(t *testing.T) {
	ad := &Adapter{}
	oci := core.Subnet{
		Id:                     sptr("ocid1.subnet.oc1..abc"),
		DisplayName:            sptr("subnet-a"),
		LifecycleState:         "AVAILABLE",
		CidrBlock:              sptr("10.0.1.0/24"),
		ProhibitPublicIpOnVnic: bptr(false),
		RouteTableId:           sptr("ocid1.routetable.oc1..rt"),
		SecurityListIds:        []string{"sl1", "sl2"},
	}

	d := ad.toDomainModel(oci)

	expect := domain.Subnet{
		OCID:            "ocid1.subnet.oc1..abc",
		DisplayName:     "subnet-a",
		LifecycleState:  "AVAILABLE",
		CidrBlock:       "10.0.1.0/24",
		Public:          true, // inverse of ProhibitPublicIpOnVnic
		RouteTableID:    "ocid1.routetable.oc1..rt",
		SecurityListIDs: []string{"sl1", "sl2"},
		NSGIDs:          nil,
	}

	if d.OCID != expect.OCID || d.DisplayName != expect.DisplayName || d.LifecycleState != expect.LifecycleState || d.CidrBlock != expect.CidrBlock || d.Public != expect.Public || d.RouteTableID != expect.RouteTableID || !reflect.DeepEqual(d.SecurityListIDs, expect.SecurityListIDs) {
		t.Fatalf("toDomainModel mismatch: got %#v want %#v", d, expect)
	}
}

func TestToDomainModel_Subnet_DefaultsWhenNil(t *testing.T) {
	ad := &Adapter{}
	oci := core.Subnet{ // many fields nil
		LifecycleState:         "TERMINATED",
		ProhibitPublicIpOnVnic: nil, // implies public
	}
	d := ad.toDomainModel(oci)
	if !d.Public {
		t.Fatalf("expected Public=true when ProhibitPublicIpOnVnic is nil")
	}
	if d.OCID != "" || d.DisplayName != "" || d.CidrBlock != "" || d.RouteTableID != "" {
		t.Fatalf("expected empty strings for unset pointers, got %#v", d)
	}
}
