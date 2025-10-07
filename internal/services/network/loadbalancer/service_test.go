package loadbalancer

import (
	"context"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	plain      []LoadBalancer
	enriched   []LoadBalancer
	getCalls   int
	listCalls  int
	elistCalls int
}

func (f *fakeRepo) GetLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	f.getCalls++
	for _, lb := range f.plain {
		if lb.OCID == ocid {
			cpy := lb
			return &cpy, nil
		}
	}
	return nil, nil
}

func (f *fakeRepo) ListLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error) {
	f.listCalls++
	return append([]LoadBalancer(nil), f.plain...), nil
}

func (f *fakeRepo) GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	for _, lb := range f.enriched {
		if lb.OCID == ocid {
			cpy := lb
			return &cpy, nil
		}
	}
	return nil, nil
}

func (f *fakeRepo) ListEnrichedLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error) {
	f.elistCalls++
	return append([]LoadBalancer(nil), f.enriched...), nil
}

func newServiceWithData(plain, enriched []LoadBalancer) *Service {
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), CompartmentID: "ocid1.compartment.oc1..test"}
	repo := &fakeRepo{plain: plain, enriched: enriched}
	return NewService(repo, appCtx)
}

func TestFuzzySearch_ByVariousFields(t *testing.T) {
	plain := []LoadBalancer{{Name: "prod-web", OCID: "ocid1.lb.oc1..aaaa", Hostnames: []string{"www.example.com"}}}
	enriched := []LoadBalancer{
		{Name: "Prod-Web", OCID: "ocid1.lb.oc1..AAAA", Hostnames: []string{"www.example.com"}, IPAddresses: []string{"10.0.0.12"}, VcnName: "vcn-prod"},
		{Name: "stage-api", OCID: "ocid1.lb.oc1..bbbb", Hostnames: []string{"api.stage.internal"}},
	}
	svc := newServiceWithData(plain, enriched)

	ctx := context.Background()

	// by name (case-insensitive, substring)
	res, err := svc.FuzzySearch(ctx, "prod")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, "Prod-Web", res[0].Name)

	// by hostname
	res, err = svc.FuzzySearch(ctx, "example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	// by OCID (case-insensitive exact/substring path)
	res, err = svc.FuzzySearch(ctx, "ocid1.lb.oc1..aaaa")
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestFetchPaginatedLoadBalancers_RespectsShowAll(t *testing.T) {
	plain := []LoadBalancer{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	enriched := []LoadBalancer{{Name: "A"}, {Name: "B"}, {Name: "C"}}
	repo := &fakeRepo{plain: plain, enriched: enriched}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), CompartmentID: "ocid1.compartment.oc1..test"}
	svc := NewService(repo, appCtx)
	ctx := context.Background()

	// showAll=false uses ListLoadBalancers
	items, total, next, err := svc.FetchPaginatedLoadBalancers(ctx, 2, 1, false)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, "2", next)
	assert.Len(t, items, 2)
	assert.Equal(t, 1, repo.listCalls)
	assert.Equal(t, 0, repo.elistCalls)

	// showAll=true uses ListEnrichedLoadBalancers
	items, total, next, err = svc.FetchPaginatedLoadBalancers(ctx, 2, 2, true)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, "", next)
	assert.Len(t, items, 1)
	assert.Equal(t, 1, repo.elistCalls)
}
