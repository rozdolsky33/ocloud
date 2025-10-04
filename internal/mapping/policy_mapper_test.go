package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	domain "github.com/rozdolsky33/ocloud/internal/domain/identity"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestNewPolicyAttributesFromOCIPolicy_AllFields(t *testing.T) {
	name := "test-policy"
	id := "ocid1.policy.oc1..exampleuniqueID"
	desc := "Test policy description"
	statements := []string{
		"Allow group Devs to manage all-resources in tenancy",
		"Allow group QA to read all-resources in tenancy",
	}
	tCreated := time.Date(2024, 7, 21, 10, 11, 12, 0, time.UTC)

	p := identity.Policy{
		Name:        &name,
		Id:          &id,
		Description: &desc,
		Statements:  statements,
		TimeCreated: &common.SDKTime{Time: tCreated},
		FreeformTags: map[string]string{
			"env": "dev",
		},
		DefinedTags: map[string]map[string]interface{}{
			"ns": map[string]interface{}{"k": "v"},
		},
	}

	attrs := mapping.NewPolicyAttributesFromOCIPolicy(p)
	require.NotNil(t, attrs)
	require.NotNil(t, attrs.Name)
	require.Equal(t, name, *attrs.Name)
	require.NotNil(t, attrs.ID)
	require.Equal(t, id, *attrs.ID)
	require.Equal(t, statements, attrs.Statement)
	require.NotNil(t, attrs.Description)
	require.Equal(t, desc, *attrs.Description)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, tCreated.Equal(*attrs.TimeCreated))
	require.Equal(t, map[string]string{"env": "dev"}, attrs.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, attrs.DefinedTags)
}

func TestNewDomainPolicyFromAttrs_AllFields(t *testing.T) {
	name := "policyA"
	id := "ocid1.policy.oc1..abc"
	desc := "desc"
	created := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)

	attrs := &mapping.PolicyAttributes{
		Name:        &name,
		ID:          &id,
		Statement:   []string{"s1", "s2"},
		Description: &desc,
		TimeCreated: &created,
		FreeformTags: map[string]string{
			"team": "core",
		},
		DefinedTags: map[string]map[string]interface{}{
			"ns": {"k": "v"},
		},
	}

	got := mapping.NewDomainPolicyFromAttrs(attrs)
	require.NotNil(t, got)
	require.IsType(t, &domain.Policy{}, got)
	require.Equal(t, name, got.Name)
	require.Equal(t, id, got.ID)
	require.Equal(t, []string{"s1", "s2"}, got.Statement)
	require.Equal(t, desc, got.Description)
	require.True(t, created.Equal(got.TimeCreated))
	require.Equal(t, map[string]string{"team": "core"}, got.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, got.DefinedTags)
}

func TestNewDomainPolicyFromAttrs_NilFields(t *testing.T) {
	attrs := &mapping.PolicyAttributes{
		// Name, ID, Description, TimeCreated are nil
		Statement:    nil,
		FreeformTags: nil,
		DefinedTags:  nil,
	}

	got := mapping.NewDomainPolicyFromAttrs(attrs)
	require.NotNil(t, got)
	require.Equal(t, "", got.Name)
	require.Equal(t, "", got.ID)
	require.Equal(t, "", got.Description)
	require.True(t, got.TimeCreated.IsZero())
	require.Nil(t, got.FreeformTags)
	require.Nil(t, got.DefinedTags)
	require.Nil(t, got.Statement)
}

func TestEndToEnd_OCIPolicy_To_DomainPolicy(t *testing.T) {
	name := "end2end"
	id := "ocid1.policy.oc1..end2end"
	desc := "e2e desc"
	statements := []string{"rule1"}
	created := time.Now().UTC().Truncate(time.Second)

	op := identity.Policy{
		Name:         &name,
		Id:           &id,
		Statements:   statements,
		Description:  &desc,
		TimeCreated:  &common.SDKTime{Time: created},
		FreeformTags: map[string]string{"a": "b"},
		DefinedTags:  map[string]map[string]interface{}{"ns": map[string]interface{}{"k": "v"}},
	}

	attrs := mapping.NewPolicyAttributesFromOCIPolicy(op)
	got := mapping.NewDomainPolicyFromAttrs(attrs)

	require.Equal(t, name, got.Name)
	require.Equal(t, id, got.ID)
	require.Equal(t, desc, got.Description)
	require.Equal(t, statements, got.Statement)
	require.True(t, created.Equal(got.TimeCreated))
	require.Equal(t, map[string]string{"a": "b"}, got.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, got.DefinedTags)
}
