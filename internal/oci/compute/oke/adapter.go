package oke

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/containerengine"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
	"github.com/rozdolsky33/ocloud/internal/mapping"
)

// Adapter is an infrastructure-layer adapter for OKE clusters.
type Adapter struct {
	client containerengine.ContainerEngineClient
}

// NewAdapter creates a new OKE adapter.
func NewAdapter(client containerengine.ContainerEngineClient) *Adapter {
	return &Adapter{client: client}
}

// GetCluster retrieves a single cluster by its OCID and enriches it with node pools.
func (a *Adapter) GetCluster(ctx context.Context, clusterOCID string) (*domain.Cluster, error) {
	resp, err := a.client.GetCluster(ctx, containerengine.GetClusterRequest{
		ClusterId: &clusterOCID,
	})
	if err != nil {
		return nil, fmt.Errorf("getting cluster from OCI: %w", err)
	}

	dc, err := a.enrichAndMapCluster(ctx, resp.Cluster)
	if err != nil {
		return nil, err
	}
	return dc, nil
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

// mapAndEnrichClusters maps OCI clusters (summaries) to domain models and enriches them with node pools.
func (a *Adapter) mapAndEnrichClusters(ctx context.Context, ociClusters []containerengine.ClusterSummary) ([]domain.Cluster, error) {
	var domainClusters []domain.Cluster
	for _, ociCluster := range ociClusters {
		dc := mapping.NewDomainClusterFromAttrs(mapping.NewClusterAttributesFromOCIClusterSummary(ociCluster))

		if ociCluster.CompartmentId == nil || ociCluster.Id == nil {
			domainClusters = append(domainClusters, *dc)
			continue
		}

		nodePools, err := a.listNodePools(ctx, *ociCluster.CompartmentId, *ociCluster.Id)
		if err != nil {
			return nil, fmt.Errorf("enriching cluster %s with node pools: %w", dc.OCID, err)
		}
		dc.NodePools = nodePools
		domainClusters = append(domainClusters, *dc)
	}
	return domainClusters, nil
}

// enrichAndMapCluster maps a single full OCI cluster object to a domain model and enriches it with node pools.
func (a *Adapter) enrichAndMapCluster(ctx context.Context, c containerengine.Cluster) (*domain.Cluster, error) {
	dc := mapping.NewDomainClusterFromAttrs(mapping.NewClusterAttributesFromOCICluster(c))
	if c.CompartmentId != nil && c.Id != nil {
		nps, err := a.listNodePools(ctx, *c.CompartmentId, *c.Id)
		if err != nil {
			return dc, fmt.Errorf("enriching cluster %s with node pools: %w", dc.OCID, err)
		}
		dc.NodePools = nps
	}
	return dc, nil
}

// listNodePools fetches all node pools in a cluster.
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
			domainNodePools = append(domainNodePools, *mapping.NewDomainNodePoolFromAttrs(mapping.NewNodePoolAttributesFromOCINodePoolSummary(ociNodePool)))
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return domainNodePools, nil
}
