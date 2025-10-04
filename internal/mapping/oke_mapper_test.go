package mapping_test

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/mapping"
	"github.com/stretchr/testify/require"
)

func TestCluster_Attributes_From_OCI_And_Domain(t *testing.T) {
	id := "ocid1.cluster.oc1..abc"
	name := "prod-oke"
	ver := "v1.29.1"
	vcn := "ocid1.vcn.oc1..vcn"
	priv := "10.0.0.5"
	pub := "https://api.k8s.local"
	created := time.Now().UTC().Truncate(time.Second)

	ociCluster := containerengine.Cluster{
		Id:                &id,
		Name:              &name,
		KubernetesVersion: &ver,
		VcnId:             &vcn,
		LifecycleState:    containerengine.ClusterLifecycleStateActive,
		Endpoints: &containerengine.ClusterEndpoints{
			PrivateEndpoint: &priv,
			Kubernetes:      &pub,
		},
		Metadata: &containerengine.ClusterMetadata{
			TimeCreated: &common.SDKTime{Time: created},
		},
		FreeformTags: map[string]string{"team": "platform"},
		DefinedTags:  map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewClusterAttributesFromOCICluster(ociCluster)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &ver, attrs.KubernetesVersion)
	require.Equal(t, &vcn, attrs.VcnOCID)
	require.Equal(t, string(containerengine.ClusterLifecycleStateActive), attrs.State)
	require.Equal(t, &priv, attrs.PrivateEndpoint)
	require.Equal(t, &pub, attrs.PublicEndpoint)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, created.Equal(*attrs.TimeCreated))

	dom := mapping.NewDomainClusterFromAttrs(attrs)
	require.IsType(t, &domain.Cluster{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, ver, dom.KubernetesVersion)
	require.Equal(t, vcn, dom.VcnOCID)
	require.Equal(t, string(containerengine.ClusterLifecycleStateActive), dom.State)
	require.Equal(t, priv, dom.PrivateEndpoint)
	require.Equal(t, pub, dom.PublicEndpoint)
	require.True(t, created.Equal(dom.TimeCreated))
	require.Equal(t, map[string]string{"team": "platform"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, dom.DefinedTags)
}

func TestClusterSummary_Attributes_From_OCI(t *testing.T) {
	id := "ocid1.cluster.oc1..xyz"
	name := "dev-oke"
	ver := "v1.28.2"
	vcn := "ocid1.vcn.oc1..abc"
	priv := "10.0.1.5"
	pub := "https://api.dev.local"
	created := time.Now().UTC().Truncate(time.Second)

	summary := containerengine.ClusterSummary{
		Id:                &id,
		Name:              &name,
		KubernetesVersion: &ver,
		VcnId:             &vcn,
		LifecycleState:    containerengine.ClusterLifecycleStateActive,
		Endpoints: &containerengine.ClusterEndpoints{
			PrivateEndpoint: &priv,
			Kubernetes:      &pub,
		},
		Metadata: &containerengine.ClusterMetadata{
			TimeCreated: &common.SDKTime{Time: created},
		},
		FreeformTags: map[string]string{"env": "dev"},
		DefinedTags:  map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewClusterAttributesFromOCIClusterSummary(summary)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &ver, attrs.KubernetesVersion)
	require.Equal(t, &vcn, attrs.VcnOCID)
	require.Equal(t, string(containerengine.ClusterLifecycleStateActive), attrs.State)
	require.Equal(t, &priv, attrs.PrivateEndpoint)
	require.Equal(t, &pub, attrs.PublicEndpoint)
	require.NotNil(t, attrs.TimeCreated)
	require.True(t, created.Equal(*attrs.TimeCreated))
}

func TestNodePool_Attributes_From_OCI_And_Domain(t *testing.T) {
	id := "ocid1.nodepool.oc1..np"
	name := "np-a"
	ver := "v1.29.1"
	shape := "VM.Standard3.Flex"
	count := 3

	np := containerengine.NodePool{
		Id:                &id,
		Name:              &name,
		KubernetesVersion: &ver,
		NodeShape:         &shape,
		NodeConfigDetails: &containerengine.NodePoolNodeConfigDetails{Size: &count},
		FreeformTags:      map[string]string{"pool": "a"},
		DefinedTags:       map[string]map[string]interface{}{"ns": {"k": "v"}},
	}

	attrs := mapping.NewNodePoolAttributesFromOCINodePool(np)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &ver, attrs.KubernetesVersion)
	require.Equal(t, &shape, attrs.NodeShape)
	require.Equal(t, &count, attrs.NodeCount)

	dom := mapping.NewDomainNodePoolFromAttrs(attrs)
	require.IsType(t, &domain.NodePool{}, dom)
	require.Equal(t, id, dom.OCID)
	require.Equal(t, name, dom.DisplayName)
	require.Equal(t, ver, dom.KubernetesVersion)
	require.Equal(t, shape, dom.NodeShape)
	require.Equal(t, count, dom.NodeCount)
	require.Equal(t, map[string]string{"pool": "a"}, dom.FreeformTags)
	require.Equal(t, map[string]map[string]interface{}{"ns": {"k": "v"}}, dom.DefinedTags)
}

func TestNodePoolSummary_Attributes_From_OCI(t *testing.T) {
	id := "ocid1.nodepool.oc1..np2"
	name := "np-b"
	ver := "v1.29.1"
	shape := "VM.Standard3.Flex"
	count := 1

	nps := containerengine.NodePoolSummary{
		Id:                &id,
		Name:              &name,
		KubernetesVersion: &ver,
		NodeShape:         &shape,
		NodeConfigDetails: &containerengine.NodePoolNodeConfigDetails{Size: &count},
	}

	attrs := mapping.NewNodePoolAttributesFromOCINodePoolSummary(nps)
	require.NotNil(t, attrs)
	require.Equal(t, &id, attrs.OCID)
	require.Equal(t, &name, attrs.DisplayName)
	require.Equal(t, &ver, attrs.KubernetesVersion)
	require.Equal(t, &shape, attrs.NodeShape)
	require.Equal(t, &count, attrs.NodeCount)
}
