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
	NodePools         []NodePool
}

// NodePool represents a node pool within an OKE cluster.
type NodePool struct {
	OCID              string
	DisplayName       string
	KubernetesVersion string
	NodeShape         string
	NodeCount         int
}

// ClusterRepository defines the port for interacting with OKE cluster storage.
type ClusterRepository interface {
	ListClusters(ctx context.Context, compartmentID string) ([]Cluster, error)
}
