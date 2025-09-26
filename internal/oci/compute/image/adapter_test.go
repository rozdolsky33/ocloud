package image

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

func strp(s string) *string { return &s }

func TestToDomainModel_Image(t *testing.T) {
	ad := &Adapter{}
	tm := time.Date(2023, 7, 1, 12, 0, 0, 0, time.UTC)
	oci := core.Image{
		Id:                     strp("ocid1.image.oc1..img"),
		DisplayName:            strp("Ubuntu 22.04"),
		OperatingSystem:        strp("Canonical Ubuntu"),
		OperatingSystemVersion: strp("22.04"),
		LaunchMode:             "NATIVE",
		TimeCreated:            &common.SDKTime{Time: tm},
	}
	d := ad.toDomainModel(oci)
	expect := domain.Image{
		OCID:                   "ocid1.image.oc1..img",
		DisplayName:            "Ubuntu 22.04",
		OperatingSystem:        "Canonical Ubuntu",
		OperatingSystemVersion: "22.04",
		LaunchMode:             "NATIVE",
		TimeCreated:            tm,
	}
	if d != expect {
		t.Fatalf("toDomainModel mismatch: got %#v want %#v", d, expect)
	}
}
