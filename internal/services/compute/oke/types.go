package oke

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/containerengine"
)

type Service struct {
	containerEngineClient containerengine.ContainerEngineClient
	logger                logr.Logger
	compartmentID         string
}

type Cluster struct {
	Name            string
	ID              string
	CreatedAt       string
	Version         string
	State           containerengine.ClusterLifecycleStateEnum
	PrivateEndpoint string
	NodePools       []NodePool
}

type NodePool struct {
	Name              string
	ID                string
	Version           string
	State             containerengine.NodePoolLifecycleStateEnum
	NodeShape         string
	NodeCount         int
	NodeSourceDetails containerengine.NodeSourceDetails
}

type JSONResponse struct {
	Clusters []Cluster `json:"clusters"`
}
