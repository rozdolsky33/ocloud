package domain

import (
	"context"
	"time"
)

// Cluster represents an OKE (Oracle Kubernetes Engine) cluster.
type Cluster struct {
	OCID              string
	DisplayName       string
	KubernetesVersion string
	VcnOCID           string
	State             string
	PrivateEndpoint   string
	PublicEndpoint    string
	TimeCreated       time.Time
	FreeformTags      map[string]string
	DefinedTags       map[string]map[string]interface{}
	NodePools         []NodePool
}

// NodePool represents a node pool within an OKE cluster.
type NodePool struct {
	OCID              string
	DisplayName       string
	KubernetesVersion string
	NodeShape         string
	NodeCount         int
	FreeformTags      map[string]string
	DefinedTags       map[string]map[string]interface{}
}

// ClusterRepository defines the port for interacting with OKE cluster storage.
type ClusterRepository interface {
	GetCluster(ctx context.Context, ocid string) (*Cluster, error)
	ListClusters(ctx context.Context, compartmentID string) ([]Cluster, error)
}
