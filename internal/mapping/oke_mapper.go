package mapping

import (
	"time"

	"github.com/oracle/oci-go-sdk/v65/containerengine"
	domain "github.com/rozdolsky33/ocloud/internal/domain/compute"
)

type ClusterAttributes struct {
	OCID              *string
	DisplayName       *string
	KubernetesVersion *string
	VcnOCID           *string
	State             string
	PrivateEndpoint   *string
	PublicEndpoint    *string
	TimeCreated       *time.Time
	FreeformTags      map[string]string
	DefinedTags       map[string]map[string]interface{}
}

func NewClusterAttributesFromOCICluster(c containerengine.Cluster) *ClusterAttributes {
	var privateEndpoint, publicEndpoint *string
	if c.Endpoints != nil {
		privateEndpoint = c.Endpoints.PrivateEndpoint
		publicEndpoint = c.Endpoints.Kubernetes
	}
	var timeCreated *time.Time
	if c.Metadata != nil && c.Metadata.TimeCreated != nil {
		t := c.Metadata.TimeCreated.Time
		timeCreated = &t
	}

	return &ClusterAttributes{
		OCID:              c.Id,
		DisplayName:       c.Name,
		KubernetesVersion: c.KubernetesVersion,
		VcnOCID:           c.VcnId,
		State:             string(c.LifecycleState),
		PrivateEndpoint:   privateEndpoint,
		PublicEndpoint:    publicEndpoint,
		TimeCreated:       timeCreated,
		FreeformTags:      c.FreeformTags,
		DefinedTags:       c.DefinedTags,
	}
}

func NewClusterAttributesFromOCIClusterSummary(c containerengine.ClusterSummary) *ClusterAttributes {
	var privateEndpoint, publicEndpoint *string
	if c.Endpoints != nil {
		privateEndpoint = c.Endpoints.PrivateEndpoint
		publicEndpoint = c.Endpoints.Kubernetes
	}
	var timeCreated *time.Time
	if c.Metadata != nil && c.Metadata.TimeCreated != nil {
		t := c.Metadata.TimeCreated.Time
		timeCreated = &t
	}

	return &ClusterAttributes{
		OCID:              c.Id,
		DisplayName:       c.Name,
		KubernetesVersion: c.KubernetesVersion,
		VcnOCID:           c.VcnId,
		State:             string(c.LifecycleState),
		PrivateEndpoint:   privateEndpoint,
		PublicEndpoint:    publicEndpoint,
		TimeCreated:       timeCreated,
		FreeformTags:      c.FreeformTags,
		DefinedTags:       c.DefinedTags,
	}
}

func NewDomainClusterFromAttrs(c *ClusterAttributes) *domain.Cluster {
	var ocid, displayName, kubernetesVersion, vcnOCID, state, privateEndpoint, publicEndpoint string
	var timeCreated time.Time

	if c.OCID != nil {
		ocid = *c.OCID
	}
	if c.DisplayName != nil {
		displayName = *c.DisplayName
	}
	if c.KubernetesVersion != nil {
		kubernetesVersion = *c.KubernetesVersion
	}
	if c.VcnOCID != nil {
		vcnOCID = *c.VcnOCID
	}
	if c.State != "" {
		state = c.State
	}
	if c.PrivateEndpoint != nil {
		privateEndpoint = *c.PrivateEndpoint
	}
	if c.PublicEndpoint != nil {
		publicEndpoint = *c.PublicEndpoint
	}
	if c.TimeCreated != nil {
		timeCreated = *c.TimeCreated
	}

	return &domain.Cluster{
		OCID:              ocid,
		DisplayName:       displayName,
		KubernetesVersion: kubernetesVersion,
		VcnOCID:           vcnOCID,
		State:             state,
		PrivateEndpoint:   privateEndpoint,
		PublicEndpoint:    publicEndpoint,
		TimeCreated:       timeCreated,
		FreeformTags:      c.FreeformTags,
		DefinedTags:       c.DefinedTags,
	}
}

type NodePoolAttributes struct {
	OCID              *string
	DisplayName       *string
	KubernetesVersion *string
	NodeShape         *string
	NodeCount         *int
	FreeformTags      map[string]string
	DefinedTags       map[string]map[string]interface{}
}

func NewNodePoolAttributesFromOCINodePool(np containerengine.NodePool) *NodePoolAttributes {
	var nodeCount *int
	if np.NodeConfigDetails != nil {
		nodeCount = np.NodeConfigDetails.Size
	}
	return &NodePoolAttributes{
		OCID:              np.Id,
		DisplayName:       np.Name,
		KubernetesVersion: np.KubernetesVersion,
		NodeShape:         np.NodeShape,
		NodeCount:         nodeCount,
		FreeformTags:      np.FreeformTags,
		DefinedTags:       np.DefinedTags,
	}
}

func NewNodePoolAttributesFromOCINodePoolSummary(np containerengine.NodePoolSummary) *NodePoolAttributes {
	var nodeCount *int
	if np.NodeConfigDetails != nil {
		nodeCount = np.NodeConfigDetails.Size
	}
	return &NodePoolAttributes{
		OCID:              np.Id,
		DisplayName:       np.Name,
		KubernetesVersion: np.KubernetesVersion,
		NodeShape:         np.NodeShape,
		NodeCount:         nodeCount,
		FreeformTags:      np.FreeformTags,
		DefinedTags:       np.DefinedTags,
	}
}

func NewDomainNodePoolFromAttrs(np *NodePoolAttributes) *domain.NodePool {
	var ocid, displayName, kubernetesVersion, nodeShape string
	var nodeCount int

	if np.OCID != nil {
		ocid = *np.OCID
	}
	if np.DisplayName != nil {
		displayName = *np.DisplayName
	}
	if np.KubernetesVersion != nil {
		kubernetesVersion = *np.KubernetesVersion
	}
	if np.NodeShape != nil {
		nodeShape = *np.NodeShape
	}
	if np.NodeCount != nil {
		nodeCount = *np.NodeCount
	}

	return &domain.NodePool{
		OCID:              ocid,
		DisplayName:       displayName,
		KubernetesVersion: kubernetesVersion,
		NodeShape:         nodeShape,
		NodeCount:         nodeCount,
		FreeformTags:      np.FreeformTags,
		DefinedTags:       np.DefinedTags,
	}
}
