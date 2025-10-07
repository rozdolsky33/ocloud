package vcn

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_GetEnrichedVcn_ViaRepo(t *testing.T) {
	v := makeVCN(0)
	repo := &fakeVCNRepo{vcns: []VCN{v}, enrichedByID: map[string]VCN{v.OCID: v}}
	svc, _ := makeService(repo)

	got, err := svc.vcnRepo.GetEnrichedVcn(context.Background(), v.OCID)
	assert.NoError(t, err)
	assert.Equal(t, v.OCID, got.OCID)
	assert.Equal(t, v.DisplayName, got.DisplayName)
}
