package oke

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// Service encapsulates OCI container engine clients and configuration.
// It provides methods to list and find clusters without printing directly.
type Service struct {
	containerEngineClient containerengine.ContainerEngineClient
	logger                logr.Logger
	compartmentID         string
}

// Cluster represents an OKE cluster with its attributes and connection details.
type Cluster struct {
	Name               string
	ID                 string
	CreatedAt          string
	Version            string
	State              containerengine.ClusterLifecycleStateEnum
	PrivateEndpoint    string
	KubernetesEndpoint string
	VcnID              string
	NodePools          []NodePool
	OKETags            util.ResourceTags
}

// NodePool represents an OKE node pool with its attributes.
type NodePool struct {
	Name      string
	ID        string
	Version   string
	State     containerengine.NodePoolLifecycleStateEnum
	NodeShape string
	NodeCount int
	Image     string
	Ocpus     string
	MemoryGB  string
	NodeTags  util.ResourceTags
}

// JSONResponse represents the JSON response from the OKE API.
type JSONResponse struct {
	Clusters []Cluster `json:"clusters"`
}
