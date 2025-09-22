package oke

import (
	"reflect"
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

func strptrOK(s string) *string { return &s }

func TestToDomainCluster_FromSummary(t *testing.T) {
	ad := &Adapter{}
	created := time.Date(2024, 5, 1, 10, 0, 0, 0, time.UTC)
	in := containerengine.ClusterSummary{
		Id:                strptrOK("ocid1.cluster.oc1..clu"),
		Name:              strptrOK("prod-oke"),
		KubernetesVersion: strptrOK("v1.29.0"),
		LifecycleState:    "ACTIVE",
		VcnId:             strptrOK("ocid1.vcn.oc1..vcn"),
		Endpoints:         &containerengine.ClusterEndpoints{PrivateEndpoint: strptrOK("https://priv"), Kubernetes: strptrOK("https://pub")},
		Metadata:          &containerengine.ClusterMetadata{TimeCreated: &common.SDKTime{Time: created}},
	}
	got := ad.toDomainCluster(in)
	want := domain.Cluster{
		OCID:              "ocid1.cluster.oc1..clu",
		DisplayName:       "prod-oke",
		KubernetesVersion: "v1.29.0",
		VcnOCID:           "ocid1.vcn.oc1..vcn",
		State:             "ACTIVE",
		PrivateEndpoint:   "https://priv",
		PublicEndpoint:    "https://pub",
		TimeCreated:       created,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("toDomainCluster(summary) mismatch: got %#v want %#v", got, want)
	}
}

func TestToDomainCluster_FromFullCluster(t *testing.T) {
	ad := &Adapter{}
	in := containerengine.Cluster{
		Id:                strptrOK("ocid1.cluster.oc1..full"),
		Name:              strptrOK("dev-oke"),
		KubernetesVersion: strptrOK("v1.28.1"),
		LifecycleState:    "DELETING",
		VcnId:             strptrOK("ocid1.vcn.oc1..abc"),
	}
	got := ad.toDomainCluster(in)
	if got.OCID != "ocid1.cluster.oc1..full" || got.DisplayName != "dev-oke" || got.KubernetesVersion != "v1.28.1" || got.State != "DELETING" || got.VcnOCID != "ocid1.vcn.oc1..abc" {
		t.Fatalf("toDomainCluster(cluster) unexpected result: %#v", got)
	}
}
