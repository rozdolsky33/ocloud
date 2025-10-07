package objectstorage

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// fakeRepo implements storage.ObjectStorageRepository for tests
type fakeRepo struct {
	byName  map[string]Bucket
	list    []Bucket
	errList error
}

func (f *fakeRepo) GetBucketNameByOCID(ctx context.Context, compartmentID, bucketOCID string) (string, error) {
	for name, b := range f.byName {
		if b.OCID == bucketOCID {
			return name, nil
		}
	}
	return "", assert.AnError
}

func (f *fakeRepo) GetBucketByName(ctx context.Context, name string) (*Bucket, error) {
	if b, ok := f.byName[name]; ok {
		return &b, nil
	}
	return nil, assert.AnError
}

func (f *fakeRepo) ListBuckets(ctx context.Context, compartmentID string) ([]Bucket, error) {
	if f.errList != nil {
		return nil, f.errList
	}
	// return copy
	out := make([]Bucket, len(f.list))
	copy(out, f.list)
	return out, nil
}

func makeBucket(i int) Bucket {
	return Bucket{
		Name:        "bucket-" + string(rune('a'+i)),
		OCID:        "ocid1.bucket.oc1.." + string(rune('A'+i)),
		Namespace:   "myns",
		StorageTier: "Standard",
		Visibility:  "private",
		Encryption:  "SSE",
		Versioning:  "Enabled",
		TimeCreated: time.Now(),
	}
}

func makeSvc(repo *fakeRepo) (*Service, *app.ApplicationContext) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), CompartmentID: "ocid1.compartment.oc1..test", Stdout: buf}
	svc := NewService(repo, appCtx.Logger, appCtx.CompartmentID)
	return svc, appCtx
}

func TestService_ListBuckets_EnrichAndPaginate(t *testing.T) {
	enriched1 := makeBucket(0)
	enriched1.ApproximateCount = 42
	enriched2 := makeBucket(1)
	enriched2.ApproximateCount = 7

	repo := &fakeRepo{
		byName: map[string]Bucket{
			"bucket-a": enriched1,
			"bucket-b": enriched2,
		},
		list: []Bucket{
			{Name: "bucket-a"},
			{Name: "bucket-b"},
			{Name: "bucket-c"}, // will not enrich (not present in byName)
		},
	}
	svc, _ := makeSvc(repo)
	ctx := context.Background()

	// ListBuckets enriches known buckets
	lst, err := svc.ListBuckets(ctx)
	assert.NoError(t, err)
	assert.Len(t, lst, 3)
	assert.Equal(t, 42, lst[0].ApproximateCount)
	assert.Equal(t, 7, lst[1].ApproximateCount)

	// FetchPaginatedBuckets limit 2 page 1
	page1, total, next, err := svc.FetchPaginatedBuckets(ctx, 2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, "2", next)
	assert.Len(t, page1, 2)
}
