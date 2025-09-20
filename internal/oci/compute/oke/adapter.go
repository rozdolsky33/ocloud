package oke

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/containerengine"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
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
	return &dc, nil
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
		dc := a.toDomainCluster(ociCluster)

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
	dc := a.toDomainCluster(c)
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

// toDomainCluster maps either a full containerengine.Cluster (from Get) or a containerengine.ClusterSummary (from List) into the single domain.Cluster type.
func (a *Adapter) toDomainCluster(ociCluster interface{}) domain.Cluster {
	var (
		clusterID         *string
		displayName       *string
		kubernetesVersion *string
		vcnID             *string
		lifecycleState    string
		endpoints         *containerengine.ClusterEndpoints
		metadata          *containerengine.ClusterMetadata
		freeformTags      map[string]string
		definedTags       map[string]map[string]interface{}
	)

	switch src := ociCluster.(type) {
	case containerengine.Cluster:
		clusterID = src.Id
		displayName = src.Name
		kubernetesVersion = src.KubernetesVersion
		vcnID = src.VcnId
		lifecycleState = string(src.LifecycleState)
		endpoints = src.Endpoints
		metadata = src.Metadata
		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags

	case containerengine.ClusterSummary:
		clusterID = src.Id
		displayName = src.Name
		kubernetesVersion = src.KubernetesVersion
		vcnID = src.VcnId
		lifecycleState = string(src.LifecycleState)
		endpoints = src.Endpoints
		metadata = src.Metadata
		freeformTags = src.FreeformTags
		definedTags = src.DefinedTags

	default:
		return domain.Cluster{}
	}

	domainCluster := domain.Cluster{}

	if clusterID != nil {
		domainCluster.OCID = *clusterID
	}
	if displayName != nil {
		domainCluster.DisplayName = *displayName
	}
	if kubernetesVersion != nil {
		domainCluster.KubernetesVersion = *kubernetesVersion
	}
	if vcnID != nil {
		domainCluster.VcnOCID = *vcnID
	}

	domainCluster.State = lifecycleState

	if endpoints != nil {
		if endpoints.PrivateEndpoint != nil {
			domainCluster.PrivateEndpoint = *endpoints.PrivateEndpoint
		}
		if endpoints.Kubernetes != nil {
			domainCluster.PublicEndpoint = *endpoints.Kubernetes
		}
	}

	if metadata != nil && metadata.TimeCreated != nil {
		domainCluster.TimeCreated = metadata.TimeCreated.Time
	}

	domainCluster.FreeformTags = freeformTags
	domainCluster.DefinedTags = definedTags

	return domainCluster
}
