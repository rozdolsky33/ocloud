package vcn

import (
	"context"
	"testing"

	dn "github.com/rozdolsky33/ocloud/internal/domain/network/vcn"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestService_FuzzySearch_VCN(t *testing.T) {
	v1 := makeVCN(0)
	v1.DisplayName = "prod-vcn"
	v1.FreeformTags = map[string]string{"env": "prod"}
	v1.DnsLabel = "corp"
	v1.DomainName = "corp.oraclevcn.com"
	v1.Gateways = []dn.Gateway{{OCID: "gw1", DisplayName: "igw-prod"}}
	v1.Subnets = []dn.Subnet{{OCID: "sn1", DisplayName: "subnet-a"}}
	v1.NSGs = []dn.NSG{{OCID: "nsg1", DisplayName: "web-nsg"}}
	v1.RouteTables = []dn.RouteTable{{OCID: "rt1", DisplayName: "rt-main"}}
	v1.SecurityLists = []dn.SecurityList{{OCID: "sl1", DisplayName: "seclist-1"}}

	v2 := makeVCN(1)
	v2.DisplayName = "dev-vcn"
	v2.FreeformTags = map[string]string{"env": "dev"}
	v2.DnsLabel = "dev"
	v2.DomainName = "dev.oraclevcn.com"
	v2.Gateways = []dn.Gateway{{OCID: "gw2", DisplayName: "nat-dev"}}

	repo := &fakeVCNRepo{vcns: []VCN{v1, v2}}
	svc := NewService(repo, logger.NewTestLogger(), "ocid1.compartment.oc1..test")

	ctx := context.Background()

	// Search by name
	res, err := svc.FuzzySearch(ctx, "prod")
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "prod-vcn", res[0].DisplayName)
	}

	// Search by tag value
	res, err = svc.FuzzySearch(ctx, "dev")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	// Search by gateway name
	res, err = svc.FuzzySearch(ctx, "igw-prod")
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Equal(t, "prod-vcn", res[0].DisplayName)
	}

	// Search by OCID substring should still match
	res, err = svc.FuzzySearch(ctx, v2.OCID[len(v2.OCID)-3:])
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
