package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/identitydomains"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestNewDynamicGroupAttributesFromOCI_And_ToDomain(t *testing.T) {
	name := "dg-name"
	id := "ocid1.dynamicgroup.oc1..abcd"
	desc := "dg description"
	rule := "Any {instance.id = 'ocid1.instance.oc1..123'}"
	state := identity.DynamicGroupLifecycleStateActive
	now := time.Now()

	ocidg := identity.DynamicGroup{
		Id:             &id,
		Name:           &name,
		Description:    &desc,
		MatchingRule:   &rule,
		LifecycleState: state,
		TimeCreated:    &common.SDKTime{Time: now},
		FreeformTags:   map[string]string{"env": "prod"},
		DefinedTags:    map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewDynamicGroupAttributesFromOCI(ocidg)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.Name)
	require.Equal(t, &desc, attrs.Description)
	require.Equal(t, &rule, attrs.MatchingRule)
	require.Equal(t, string(state), attrs.LifecycleState)
	require.Equal(t, now, *attrs.TimeCreated)

	dom := mapping.NewDomainDynamicGroupFromAttrs(attrs)
	require.IsType(t, &domain.DynamicGroup{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.Name)
	require.Equal(t, desc, dom.Description)
	require.Equal(t, rule, dom.MatchingRule)
	require.Equal(t, string(state), dom.LifecycleState)
	require.WithinDuration(t, now, dom.TimeCreated, time.Second)
}

func TestNewDynamicGroupAttributesFromIDCS_And_ToDomain(t *testing.T) {
	id := "idcs-dg-id"
	ocid := "ocid1.dynamicgroup.oc1..idcs"
	name := "idcs-dg-name"
	desc := "idcs-dg-desc"
	rule := "All {resource.type = 'fnfunc'}"
	domainURL := "https://idcs-xyz.identity.oraclecloud.com"
	nowStr := time.Now().Format(time.RFC3339)

	idcsdg := identitydomains.DynamicResourceGroup{
		Id:           &id,
		Ocid:         &ocid,
		DisplayName:  &name,
		Description:  &desc,
		MatchingRule: &rule,
		Meta: &identitydomains.Meta{
			Created: &nowStr,
		},
	}

	attrs := mapping.NewDynamicGroupAttributesFromIDCS(idcsdg, domainURL)
	require.NotNil(t, attrs)
	require.Equal(t, &ocid, attrs.OCID)
	require.Equal(t, &name, attrs.Name)
	require.Equal(t, &desc, attrs.Description)
	require.Equal(t, &rule, attrs.MatchingRule)
	require.Equal(t, domainURL, attrs.DomainURL)

	dom := mapping.NewDomainDynamicGroupFromAttrs(attrs)
	require.Equal(t, ocid, dom.OCID)
	require.Equal(t, name, dom.Name)
	require.Equal(t, desc, dom.Description)
	require.Equal(t, rule, dom.MatchingRule)
}
