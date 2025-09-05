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
	return &dc, nil
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

// mapAndEnrichClusters maps OCI clusters (summaries) to domain models and enriches them with node pools.
func (a *Adapter) mapAndEnrichClusters(ctx context.Context, ociClusters []containerengine.ClusterSummary) ([]domain.Cluster, error) {
	var domainClusters []domain.Cluster
	for _, ociCluster := range ociClusters {
		dc := a.toDomainModelSummary(ociCluster)

		if ociCluster.CompartmentId == nil || ociCluster.Id == nil {
			domainClusters = append(domainClusters, dc)
			continue
		}

		nodePools, err := a.listNodePools(ctx, *ociCluster.CompartmentId, *ociCluster.Id)
		if err != nil {
			return nil, fmt.Errorf("enriching cluster %s with node pools: %w", dc.OCID, err)
		}
		dc.NodePools = nodePools
		domainClusters = append(domainClusters, dc)
	}
	return domainClusters, nil
}

// enrichAndMapCluster maps a single full OCI cluster object to a domain model and enriches it with node pools.
func (a *Adapter) enrichAndMapCluster(ctx context.Context, c containerengine.Cluster) (domain.Cluster, error) {
	dc := a.toDomainModel(c)
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
			dnp := domain.NodePool{}
			if ociNodePool.Id != nil {
				dnp.OCID = *ociNodePool.Id
			}
			if ociNodePool.Name != nil {
				dnp.DisplayName = *ociNodePool.Name
			}
			if ociNodePool.KubernetesVersion != nil {
				dnp.KubernetesVersion = *ociNodePool.KubernetesVersion
			}
			if ociNodePool.NodeShape != nil {
				dnp.NodeShape = *ociNodePool.NodeShape
			}
			if ociNodePool.NodeConfigDetails != nil && ociNodePool.NodeConfigDetails.Size != nil {
				dnp.NodeCount = *ociNodePool.NodeConfigDetails.Size
			}
			dnp.FreeformTags = ociNodePool.FreeformTags
			dnp.DefinedTags = ociNodePool.DefinedTags
			domainNodePools = append(domainNodePools, dnp)
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}
	return domainNodePools, nil
}

func (a *Adapter) toDomainModel(c containerengine.Cluster) domain.Cluster {
	dc := domain.Cluster{}
	if c.Id != nil {
		dc.OCID = *c.Id
	}
	if c.Name != nil {
		dc.DisplayName = *c.Name
	}
	if c.KubernetesVersion != nil {
		dc.KubernetesVersion = *c.KubernetesVersion
	}
	if c.VcnId != nil {
		dc.VcnOCID = *c.VcnId
	}
	dc.State = string(c.LifecycleState)
	if c.Endpoints != nil {
		if c.Endpoints.PrivateEndpoint != nil {
			dc.PrivateEndpoint = *c.Endpoints.PrivateEndpoint
		}
		if c.Endpoints.Kubernetes != nil {
			dc.PublicEndpoint = *c.Endpoints.Kubernetes
		}
	}
	if c.Metadata != nil && c.Metadata.TimeCreated != nil {
		dc.TimeCreated = c.Metadata.TimeCreated.Time
	}
	dc.FreeformTags = c.FreeformTags
	dc.DefinedTags = c.DefinedTags
	return dc
}

// toDomainModelSummary maps an OCI ClusterSummary to a domain.Cluster (without enrichment).
func (a *Adapter) toDomainModelSummary(s containerengine.ClusterSummary) domain.Cluster {
	dc := domain.Cluster{}
	if s.Id != nil {
		dc.OCID = *s.Id
	}
	if s.Name != nil {
		dc.DisplayName = *s.Name
	}
	if s.KubernetesVersion != nil {
		dc.KubernetesVersion = *s.KubernetesVersion
	}
	if s.VcnId != nil {
		dc.VcnOCID = *s.VcnId
	}
	dc.State = string(s.LifecycleState)
	if s.Endpoints != nil {
		if s.Endpoints.PrivateEndpoint != nil {
			dc.PrivateEndpoint = *s.Endpoints.PrivateEndpoint
		}
		if s.Endpoints.Kubernetes != nil {
			dc.PublicEndpoint = *s.Endpoints.Kubernetes
		}
	}
	if s.Metadata != nil && s.Metadata.TimeCreated != nil {
		dc.TimeCreated = s.Metadata.TimeCreated.Time
	}
	dc.FreeformTags = s.FreeformTags
	dc.DefinedTags = s.DefinedTags
	return dc
}
