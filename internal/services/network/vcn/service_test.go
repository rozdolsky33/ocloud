package vcn

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

// fakeVCNRepo implements domain.VCNRepository for tests
type fakeVCNRepo struct {
	vcns         []VCN
	errList      error
	enrichedByID map[string]VCN
}

func (f *fakeVCNRepo) GetEnrichedVcn(ctx context.Context, ocid string) (VCN, error) {
	if v, ok := f.enrichedByID[ocid]; ok {
		return v, nil
	}
	return VCN{}, assert.AnError
}

func (f *fakeVCNRepo) ListVcns(ctx context.Context, compartmentID string) ([]VCN, error) {
	if f.errList != nil {
		return nil, f.errList
	}
	// return a shallow copy (without related resources)
	out := make([]VCN, len(f.vcns))
	copy(out, f.vcns)
	return out, nil
}

func (f *fakeVCNRepo) ListEnrichedVcns(ctx context.Context, compartmentID string) ([]VCN, error) {
	if f.errList != nil {
		return nil, f.errList
	}
	return f.vcns, nil
}

func makeVCN(i int) VCN {
	return VCN{
		OCID:           "ocid1.vcn.oc1.." + string(rune('A'+i)),
		DisplayName:    "vcn-" + string(rune('a'+i)),
		LifecycleState: "AVAILABLE",
		CidrBlocks:     []string{"10.0." + string(rune('0'+i)) + ".0/24"},
		DnsLabel:       "corp" + string(rune('a'+i)),
		DomainName:     "corp" + string(rune('a'+i)) + ".oraclevcn.com",
		TimeCreated:    time.Now(),
	}
}

func makeService(repo *fakeVCNRepo) (*Service, *app.ApplicationContext) {
	buf := &bytes.Buffer{}
	appCtx := &app.ApplicationContext{Logger: logger.NewTestLogger(), CompartmentID: "ocid1.compartment.oc1..test", Stdout: buf}
	svc := NewService(repo, appCtx.Logger, appCtx.CompartmentID)
	return svc, appCtx
}

func TestService_ListAndPaginateVCNs(t *testing.T) {
	repo := &fakeVCNRepo{vcns: []VCN{makeVCN(0), makeVCN(1), makeVCN(2)}}
	svc, _ := makeService(repo)
	ctx := context.Background()

	// List
	list, err := svc.ListVcns(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 3)

	// Paginate: limit 2, page 1
	page1, total, next, err := svc.FetchPaginatedVCNs(ctx, 2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, "2", next)
	assert.Len(t, page1, 2)

	// Paginate: page 2
	page2, total2, next2, err := svc.FetchPaginatedVCNs(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3, total2)
	assert.Equal(t, "", next2)
	assert.Len(t, page2, 1)
}
