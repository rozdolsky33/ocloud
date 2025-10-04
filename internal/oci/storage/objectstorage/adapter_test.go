package objectstorage

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

func TestToDomainBucketFromSummary(t *testing.T) {
	a := &Adapter{}
	tm := commonTime(t, 2025, 10, 4)
	name := "b1"
	ns := "ns1"
	sum := objectstorage.BucketSummary{
		Name:        &name,
		Namespace:   &ns,
		TimeCreated: &common.SDKTime{Time: tm},
	}
	db := a.toDomainBucketFromSummary(sum)
	if db.Name != name || db.Namespace != ns || !db.TimeCreated.Equal(tm) {
		t.Fatalf("unexpected mapping: %+v", db)
	}
}

func TestToDomainBucketFromBucket_MappingAndFlags(t *testing.T) {
	a := &Adapter{}
	tm := commonTime(t, 2025, 10, 4)
	name := "bucketX"
	ns := "ns"
	id := "ocid1.bucket.oc1..xyz"
	kms := "ocid1.key.oc1..kms"
	approxCount := int64(42)
	approxSize := int64(1<<20 + 512) // ~1 MiB
	ro := false
	detail := objectstorage.Bucket{
		Name:               &name,
		Id:                 &id,
		Namespace:          &ns,
		TimeCreated:        &common.SDKTime{Time: tm},
		StorageTier:        objectstorage.BucketStorageTierStandard,
		PublicAccessType:   objectstorage.BucketPublicAccessTypeNopublicaccess,
		KmsKeyId:           &kms,
		Versioning:         objectstorage.BucketVersioningDisabled,
		ReplicationEnabled: &ro, // false
		IsReadOnly:         &ro,
		ApproximateCount:   &approxCount,
		ApproximateSize:    &approxSize,
		FreeformTags:       map[string]string{"env": "dev"},
		DefinedTags:        map[string]map[string]interface{}{"ns": {"k": "v"}},
	}
	db := a.toDomainBucketFromBucket(detail)
	if db.Name != name || db.OCID != id || db.Namespace != ns || !db.TimeCreated.Equal(tm) {
		t.Fatalf("basic fields not mapped: %+v", db)
	}
	if db.StorageTier != "Standard" || db.Visibility != "NoPublicAccess" || db.Versioning != "Disabled" {
		t.Fatalf("enum mappings incorrect: %+v", db)
	}
	if db.Encryption != "Customer-managed (KMS)" { // KMS key present
		t.Fatalf("encryption mapping incorrect: %s", db.Encryption)
	}
	if db.ReplicationEnabled != false || db.IsReadOnly != false {
		t.Fatalf("boolean flags incorrect: %+v", db)
	}
	if db.ApproximateCount != int(approxCount) || db.ApproximateSize != approxSize {
		t.Fatalf("approximate metrics incorrect: %+v", db)
	}
	if db.FreeformTags["env"] != "dev" || db.DefinedTags["ns"]["k"] != "v" {
		t.Fatalf("tags not mapped: %+v", db)
	}
}

func commonTime(t *testing.T, y int, m int, d int) time.Time {
	t.Helper()
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}
