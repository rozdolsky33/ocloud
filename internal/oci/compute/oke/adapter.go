package oke

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// Adapter is an infrastructure-layer adapter for OKE clusters.
type Adapter struct {
	client containerengine.ContainerEngineClient
}

// NewAdapter creates a new OKE adapter.
func NewAdapter(client containerengine.ContainerEngineClient) *Adapter {
	return &Adapter{client: client}
}

// ListClusters fetches all clusters in a compartment and enriches them with node pools.
func (a *Adapter) ListClusters(ctx context.Context, compartmentID string) ([]domain.Cluster, error) {
	var ociClusters []containerengine.ClusterSummary
	var page *string

	for {
		resp, err := a.client.ListClusters(ctx, containerengine.ListClustersRequest{
			CompartmentId: &compartmentID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing OKE clusters from OCI: %w", err)
		}
		ociClusters = append(ociClusters, resp.Items...)

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return a.mapAndEnrichClusters(ctx, ociClusters)
}

func (a *Adapter) mapAndEnrichClusters(ctx context.Context, ociClusters []containerengine.ClusterSummary) ([]domain.Cluster, error) {
	var domainClusters []domain.Cluster
	for _, ociCluster := range ociClusters {
		dc := domain.Cluster{
			OCID:              *ociCluster.Id,
			DisplayName:       *ociCluster.Name,
			KubernetesVersion: *ociCluster.KubernetesVersion,
			VcnOCID:           *ociCluster.VcnId,
			State:             string(ociCluster.LifecycleState),
		}
		if ociCluster.Endpoints != nil {
			if ociCluster.Endpoints.PrivateEndpoint != nil {
				dc.PrivateEndpoint = *ociCluster.Endpoints.PrivateEndpoint
			}
			if ociCluster.Endpoints.Kubernetes != nil {
				dc.PublicEndpoint = *ociCluster.Endpoints.Kubernetes
			}
		}
		if ociCluster.Metadata != nil && ociCluster.Metadata.TimeCreated != nil {
			dc.TimeCreated = ociCluster.Metadata.TimeCreated.Time
		}

		// Check if CompartmentId and I'd are not nil before dereferencing
		if ociCluster.CompartmentId == nil || ociCluster.Id == nil {
			domainClusters = append(domainClusters, dc)
			continue
		}

		nodePools, err := a.listNodePools(ctx, *ociCluster.CompartmentId, *ociCluster.Id)
		if err != nil {
			// In a real-world scenario, you might want to handle this more gracefully
			// (e.g., log the error and continue), but for now, we'll fail fast.
			return nil, fmt.Errorf("enriching cluster %s with node pools: %w", dc.OCID, err)
		}
		dc.NodePools = nodePools
		domainClusters = append(domainClusters, dc)
	}
	return domainClusters, nil
}

func (a *Adapter) listNodePools(ctx context.Context, compartmentID, clusterID string) ([]domain.NodePool, error) {
	var domainNodePools []domain.NodePool
	var page *string

	for {
		resp, err := a.client.ListNodePools(ctx, containerengine.ListNodePoolsRequest{
			CompartmentId: &compartmentID,
			ClusterId:     &clusterID,
			Page:          page,
		})
		if err != nil {
			return nil, fmt.Errorf("listing node pools from OCI: %w", err)
		}

		for _, ociNodePool := range resp.Items {
			dnp := domain.NodePool{
				OCID:              *ociNodePool.Id,
				DisplayName:       *ociNodePool.Name,
				KubernetesVersion: *ociNodePool.KubernetesVersion,
				NodeShape:         *ociNodePool.NodeShape,
			}
			if ociNodePool.NodeConfigDetails != nil && ociNodePool.NodeConfigDetails.Size != nil {
				dnp.NodeCount = *ociNodePool.NodeConfigDetails.Size
			}
			domainNodePools = append(domainNodePools, dnp)
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return domainNodePools, nil
}
