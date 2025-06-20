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
	Name                string
	ID                  string
	Version             string
	State               string
	CreatedAt           string
	PrivateEndpoint     string
	PrivateEndpointIP   string
	PrivateEndpointPort string
}
