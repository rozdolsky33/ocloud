package mapping

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

func TestNewBucketAttributesFromOCIBucket_And_NewDomainBucketFromAttrs(t *testing.T) {
	// Arrange OCI bucket with most fields populated
	name := "my-bucket"
	id := "ocid1.bucket.oc1..exampleuniqueID"
	ns := "mynamespace"
	kms := "ocid1.key.oc1..kmsKey"
	approxCount := int64(42)
	approxSize := int64(5 * 1024 * 1024)
	ff := map[string]string{"Env": "Dev"}
	def := map[string]map[string]interface{}{"Oracle-Tags": {"CreatedBy": "tester"}}
	created := common.SDKTime{Time: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)}
	public := objectstorage.BucketPublicAccessTypeObjectread
	storTier := objectstorage.BucketStorageTierStandard
	ver := objectstorage.BucketVersioningEnabled
	rep := true
	ro := false

	ociBucket := objectstorage.Bucket{
		Name:               &name,
		Id:                 &id,
		Namespace:          &ns,
		TimeCreated:        &created,
		PublicAccessType:   public,
		StorageTier:        storTier,
		KmsKeyId:           &kms,
		Versioning:         ver,
		ReplicationEnabled: &rep,
		IsReadOnly:         &ro,
		ApproximateCount:   &approxCount,
		ApproximateSize:    &approxSize,
		FreeformTags:       ff,
		DefinedTags:        def,
	}

	attrs := NewBucketAttributesFromOCIBucket(ociBucket)
	if attrs == nil {
		t.Fatalf("attrs is nil")
	}
	// Validate attributes mapping
	if attrs.Name == nil || *attrs.Name != name {
		t.Errorf("Name mismatch: %v", attrs.Name)
	}
	if attrs.ID == nil || *attrs.ID != id {
		t.Errorf("ID mismatch: %v", attrs.ID)
	}
	if attrs.Namespace == nil || *attrs.Namespace != ns {
		t.Errorf("Namespace mismatch: %v", attrs.Namespace)
	}
	if attrs.TimeCreated == nil || !attrs.TimeCreated.Equal(created.Time) {
		t.Errorf("TimeCreated mismatch: %v", attrs.TimeCreated)
	}
	if attrs.StorageTier != string(storTier) {
		t.Errorf("StorageTier mismatch: %s", attrs.StorageTier)
	}
	if attrs.PublicAccessType != string(public) {
		t.Errorf("PublicAccessType mismatch: %s", attrs.PublicAccessType)
	}
	if attrs.KmsKeyID == nil || *attrs.KmsKeyID != kms {
		t.Errorf("KmsKeyID mismatch: %v", attrs.KmsKeyID)
	}
	if attrs.Versioning != string(ver) {
		t.Errorf("Versioning mismatch: %s", attrs.Versioning)
	}
	if attrs.ReplicationEnabled == nil || *attrs.ReplicationEnabled != rep {
		t.Errorf("ReplicationEnabled mismatch: %v", attrs.ReplicationEnabled)
	}
	if attrs.IsReadOnly == nil || *attrs.IsReadOnly != ro {
		t.Errorf("IsReadOnly mismatch: %v", attrs.IsReadOnly)
	}
	if attrs.ApproximateCount == nil || *attrs.ApproximateCount != approxCount {
		t.Errorf("ApproximateCount mismatch: %v", attrs.ApproximateCount)
	}
	if attrs.ApproximateSize == nil || *attrs.ApproximateSize != approxSize {
		t.Errorf("ApproximateSize mismatch: %v", attrs.ApproximateSize)
	}
	if len(attrs.FreeformTags) != 1 || attrs.FreeformTags["Env"] != "Dev" {
		t.Errorf("FreeformTags mismatch: %v", attrs.FreeformTags)
	}
	if len(attrs.DefinedTags) == 0 || attrs.DefinedTags["Oracle-Tags"]["CreatedBy"] != "tester" {
		t.Errorf("DefinedTags mismatch: %v", attrs.DefinedTags)
	}

	// Build domain from attrs
	db := NewDomainBucketFromAttrs(*attrs)
	if db.Name != name {
		t.Errorf("domain Name: %s", db.Name)
	}
	if db.OCID != id {
		t.Errorf("domain OCID: %s", db.OCID)
	}
	if db.Namespace != ns {
		t.Errorf("domain Namespace: %s", db.Namespace)
	}
	if db.StorageTier != string(storTier) {
		t.Errorf("domain StorageTier: %s", db.StorageTier)
	}
	if db.Visibility != string(public) {
		t.Errorf("domain Visibility: %s", db.Visibility)
	}
	if db.Versioning != string(ver) {
		t.Errorf("domain Versioning: %s", db.Versioning)
	}
	if db.Encryption != "Customer-managed (KMS)" {
		t.Errorf("domain Encryption expected KMS, got: %s", db.Encryption)
	}
	if !db.ReplicationEnabled {
		t.Errorf("domain ReplicationEnabled: %v", db.ReplicationEnabled)
	}
	if db.IsReadOnly {
		t.Errorf("domain IsReadOnly: %v", db.IsReadOnly)
	}
	if db.ApproximateCount != int(approxCount) {
		t.Errorf("domain ApproximateCount: %d", db.ApproximateCount)
	}
	if db.ApproximateSize != approxSize {
		t.Errorf("domain ApproximateSize: %d", db.ApproximateSize)
	}
	if db.TimeCreated.IsZero() {
		t.Errorf("domain TimeCreated should be set")
	}
	if db.FreeformTags["Env"] != "Dev" {
		t.Errorf("domain FreeformTags: %v", db.FreeformTags)
	}
	if db.DefinedTags["Oracle-Tags"]["CreatedBy"] != "tester" {
		t.Errorf("domain DefinedTags: %v", db.DefinedTags)
	}
}

func TestNewBucketAttributesFromOCIBucketSummary_Minimum(t *testing.T) {
	name := "sum-bucket"
	ns := "ns"
	created := common.SDKTime{Time: time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)}

	sum := objectstorage.BucketSummary{
		Name:        &name,
		Namespace:   &ns,
		TimeCreated: &created,
	}
	attrs := NewBucketAttributesFromOCIBucketSummary(sum)
	if attrs.Name == nil || *attrs.Name != name {
		t.Errorf("Name mismatch")
	}
	if attrs.Namespace == nil || *attrs.Namespace != ns {
		t.Errorf("Namespace mismatch")
	}
	if attrs.TimeCreated == nil || !attrs.TimeCreated.Equal(created.Time) {
		t.Errorf("TimeCreated mismatch")
	}

	// Domain built from minimal attrs should apply defaults
	db := NewDomainBucketFromAttrs(*attrs)
	if db.Name != name || db.Namespace != ns {
		t.Errorf("domain basic fields mismatch: %+v", db)
	}
	if db.Encryption != "Oracle-managed" {
		t.Errorf("expected default Oracle-managed encryption, got %s", db.Encryption)
	}
	if db.ApproximateCount != 0 || db.ApproximateSize != 0 {
		t.Errorf("expected zero metrics by default")
	}
}

func TestNewDomainBucketFromAttrs_DefaultsOnNil(t *testing.T) {
	// All pointer fields nil to trigger defaults
	attrs := BucketAttributes{}
	db := NewDomainBucketFromAttrs(attrs)
	if db.Name != "" || db.OCID != "" || db.Namespace != "" {
		t.Errorf("expected empty identifiers, got: %+v", db)
	}
	if db.Encryption != "Oracle-managed" {
		t.Errorf("default encryption mismatch: %s", db.Encryption)
	}
	if db.ReplicationEnabled || db.IsReadOnly {
		t.Errorf("bool defaults expected false: %+v", db)
	}
	if db.ApproximateCount != 0 || db.ApproximateSize != 0 {
		t.Errorf("numeric defaults expected zero: %+v", db)
	}
}
