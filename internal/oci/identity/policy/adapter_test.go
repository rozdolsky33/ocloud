package policy

import (
	"reflect"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
)

func sp(s string) *string { return &s }

func TestToDomainModel_Policy(t *testing.T) {
	ad := &Adapter{}
	created := time.Date(2022, 3, 4, 5, 6, 7, 0, time.UTC)
	oci := identity.Policy{
		Id:           sp("ocid1.policy.oc1..abc"),
		Name:         sp("pol-1"),
		Description:  sp("desc"),
		Statements:   []string{"Allow group A to read all-resources in tenancy"},
		TimeCreated:  &common.SDKTime{Time: created},
		FreeformTags: map[string]string{"k": "v"},
		DefinedTags:  map[string]map[string]interface{}{"ns": {"a": 1}},
	}
	d := ad.toDomainModel(oci)
	expect := domain.Policy{
		ID:           "ocid1.policy.oc1..abc",
		Name:         "pol-1",
		Description:  "desc",
		Statement:    []string{"Allow group A to read all-resources in tenancy"},
		TimeCreated:  created,
		FreeformTags: map[string]string{"k": "v"},
		DefinedTags:  map[string]map[string]interface{}{"ns": {"a": 1}},
	}
	if !reflect.DeepEqual(d, expect) {
		t.Fatalf("toDomainModel(policy) mismatch: got %#v want %#v", d, expect)
	}
}
