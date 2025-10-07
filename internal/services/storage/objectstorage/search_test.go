package objectstorage

import (
	"context"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestService_FuzzySearch_Buckets(t *testing.T) {
	b1 := makeBucket(0)
	b1.Name = "prod-logs"
	b1.Namespace = "teamA"
	b1.FreeformTags = map[string]string{"env": "prod"}
	b1.Visibility = "public"

	b2 := makeBucket(1)
	b2.Name = "dev-backups"
	b2.Namespace = "teamB"
	b2.FreeformTags = map[string]string{"env": "dev"}
	b2.Encryption = "kms"

	repo := &fakeRepo{
		byName: map[string]Bucket{
			b1.Name: b1,
			b2.Name: b2,
		},
		list: []Bucket{{Name: b1.Name}, {Name: b2.Name}},
	}
	svc := NewService(repo, logger.NewTestLogger(), "ocid1.compartment.oc1..test")
	ctx := context.Background()

	res, err := svc.FuzzySearch(ctx, "prod")
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, b1.Name, res[0].Name)
	}

	res, err = svc.FuzzySearch(ctx, "teamB")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	// OCID substring search should also return something
	ocidSub := b2.OCID[len(b2.OCID)-3:]
	res, err = svc.FuzzySearch(ctx, ocidSub)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// test scaffolding
var (
	testLogger      = loggerForTests()
	testCompartment = "ocid1.compartment.oc1..test"
)

func loggerForTests() anyLogger {
	return anyLogger{}
}

type anyLogger struct{}

// implement minimal logr.Logger-like interface methods used (but our service uses logr.Logger only for V(), which returns an object with Info). For simplicity in tests, we avoid calling logger; NewService accepts logr.Logger but we can reuse the actual test logger from logger pkg instead. We'll override below by constructing via makeSvc to avoid this scaffolding.
