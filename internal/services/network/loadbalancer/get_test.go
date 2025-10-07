package loadbalancer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// fakeRepo implements domain.LoadBalancerRepository for tests
type fakeRepo2 struct {
	byOCID      map[string]LoadBalancer
	compartment string
	errGet      error
	errList     error
}

func (f *fakeRepo2) GetLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	if f.errGet != nil {
		return nil, f.errGet
	}
	if lb, ok := f.byOCID[ocid]; ok {
		return &lb, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepo2) ListLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error) {
	if f.errList != nil {
		return nil, f.errList
	}
	out := make([]LoadBalancer, 0, len(f.byOCID))
	for _, v := range f.byOCID {
		out = append(out, v)
	}
	return out, nil
}

func (f *fakeRepo2) GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error) {
	return f.GetLoadBalancer(ctx, ocid)
}

func (f *fakeRepo2) ListEnrichedLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error) {
	return f.ListLoadBalancers(ctx, compartmentID)
}

func makeServiceWith(repo *fakeRepo2) (*Service, *app.ApplicationContext) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), CompartmentID: "ocid1.compartment.oc1..test", Stdout: buf}
	return NewService(repo, appCtx), appCtx
}

func sampleLB(i int) LoadBalancer {
	return LoadBalancer{
		OCID:        fmt.Sprintf("ocid1.loadbalancer.oc1..lb%c", 'A'+i),
		ID:          fmt.Sprintf("id-%c", 'A'+i),
		Name:        fmt.Sprintf("lb-%c", 'a'+i),
		State:       "ACTIVE",
		Type:        "public",
		IPAddresses: []string{fmt.Sprintf("10.0.0.%d", 1+i)},
		Shape:       "flex",
	}
}

func TestService_GetLoadBalancer(t *testing.T) {
	repo := &fakeRepo2{byOCID: map[string]LoadBalancer{"ocid1.loadbalancer.oc1..lbA": sampleLB(0)}}
	svc, _ := makeServiceWith(repo)

	ctx := context.Background()
	lb, err := svc.GetLoadBalancer(ctx, "ocid1.loadbalancer.oc1..lbA")
	assert.NoError(t, err)
	if assert.NotNil(t, lb) {
		assert.Equal(t, "lb-a", lb.Name)
	}

	// not found
	_, err = svc.GetLoadBalancer(ctx, "missing")
	assert.Error(t, err)
}

func TestService_FetchPaginatedLoadBalancers(t *testing.T) {
	items := map[string]LoadBalancer{
		"ocid1.loadbalancer.oc1..lbA": sampleLB(0),
		"ocid1.loadbalancer.oc1..lbB": sampleLB(1),
		"ocid1.loadbalancer.oc1..lbC": sampleLB(2),
	}
	repo := &fakeRepo2{byOCID: items}
	svc, _ := makeServiceWith(repo)
	ctx := context.Background()

	// Page 1, limit 2
	page1, total, next, err := svc.FetchPaginatedLoadBalancers(ctx, 2, 1, false)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, "2", next)
	assert.Len(t, page1, 2)

	// Page 2, limit 2
	page2, total2, next2, err := svc.FetchPaginatedLoadBalancers(ctx, 2, 2, true) // showAll toggled
	assert.NoError(t, err)
	assert.Equal(t, 3, total2)
	assert.Equal(t, "", next2)
	assert.Len(t, page2, 1)
}

func TestPrintLoadBalancersInfo_JSONAndTable(t *testing.T) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), Stdout: buf}

	lbs := []LoadBalancer{
		{OCID: "ocid1.lb.oc1..1", Name: "prod-lb", State: "ACTIVE", IPAddresses: []string{"1.1.1.1"}},
		{OCID: "ocid1.lb.oc1..2", Name: "dev-lb", State: "PROVISIONING", IPAddresses: []string{"2.2.2.2"}},
	}

	// Table (default) minimal
	err := PrintLoadBalancersInfo(lbs, appCtx, nil, false, false)
	assert.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "prod-lb")
	assert.Contains(t, out, "dev-lb")
	buf.Reset()

	// JSON
	err = PrintLoadBalancersInfo(lbs, appCtx, nil, true, true)
	assert.NoError(t, err)
	jsonOut := buf.String()
	assert.NotEmpty(t, jsonOut)
	// Accept either array or object JSON, depending on response shape
	first := jsonOut[0]
	if first != '{' && first != '[' {
		t.Fatalf("unexpected JSON start: %q", first)
	}
}
