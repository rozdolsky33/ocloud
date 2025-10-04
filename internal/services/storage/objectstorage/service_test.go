package objectstorage

import (
	"context"
	"testing"
	"time"

	storage "github.com/rozdolsky33/ocloud/internal/domain/storage/objectstorage"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

type fakeRepo struct {
	list     []storage.Bucket
	getCalls map[string]int
}

func (f *fakeRepo) GetBucket(ctx context.Context, ocid string) (*storage.Bucket, error) {
	if f.getCalls == nil {
		f.getCalls = map[string]int{}
	}
	f.getCalls[ocid]++
	b := storage.Bucket{
		Name:               ocid,
		OCID:               "ocid1.bucket.oc1.." + ocid,
		Namespace:          "ns",
		TimeCreated:        time.Unix(1000, 0),
		StorageTier:        "Standard",
		Visibility:         "NoPublicAccess",
		Encryption:         "Oracle-managed",
		Versioning:         "Disabled",
		ReplicationEnabled: true,
		IsReadOnly:         false,
		ApproximateSize:    1024,
		ApproximateCount:   7,
	}
	return &b, nil
}

func (f *fakeRepo) ListBuckets(ctx context.Context, compartmentID string) ([]storage.Bucket, error) {
	return f.list, nil
}

func TestService_FetchPaginatedBuckets_PaginationAndEnrichment(t *testing.T) {
	// Arrange: two summary buckets from ListBuckets
	list := []storage.Bucket{
		{Name: "b1", Namespace: "ns", TimeCreated: time.Unix(1, 0)},
		{Name: "b2", Namespace: "ns", TimeCreated: time.Unix(2, 0)},
	}
	fr := &fakeRepo{list: list, getCalls: map[string]int{}}
	svc := NewService(fr, logger.CmdLogger, "compartment")

	// Act: page 1 with limit 1, showAll=false (no enrichment)
	res, total, next, err := svc.FetchPaginatedBuckets(context.Background(), 1, 1, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total=2, got %d", total)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 item on page, got %d", len(res))
	}
	if next == "" {
		t.Fatalf("expected next page token, got empty")
	}
	if len(fr.getCalls) != 0 {
		t.Fatalf("expected no GetBucket calls, got %v", fr.getCalls)
	}

	// Act: showAll=true triggers enrichment for all items before pagination
	res2, total2, _, err := svc.FetchPaginatedBuckets(context.Background(), 5, 1, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total2 != 2 || len(res2) != 2 {
		t.Fatalf("expected 2 items enriched, got total=%d len=%d", total2, len(res2))
	}
	// Assert enrichment replaced summaries with full details
	for _, b := range res2 {
		if b.OCID == "" || b.ApproximateCount == 0 || b.ApproximateSize == 0 {
			t.Fatalf("bucket not enriched: %+v", b)
		}
	}
	// Ensure GetBucket called for each by name
	if fr.getCalls["b1"] != 1 || fr.getCalls["b2"] != 1 {
		t.Fatalf("expected GetBucket called once per bucket, got: %v", fr.getCalls)
	}
}

// noopLogger implements logr.Logger minimal methods used in service
// We avoid importing the real logger in tests for simplicity.
type noopLogger struct{}

func (noopLogger) Enabled() bool                                             { return false }
func (noopLogger) Info(msg string, keysAndValues ...interface{})             {}
func (noopLogger) Error(err error, msg string, keysAndValues ...interface{}) {}
func (noopLogger) V(level int) noopLogger                                    { return noopLogger{} }
func (noopLogger) WithValues(keysAndValues ...interface{}) noopLogger        { return noopLogger{} }
func (noopLogger) WithName(name string) noopLogger                           { return noopLogger{} }
