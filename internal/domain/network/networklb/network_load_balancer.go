package networklb

import (
	"context"
	"time"
)

type NetworkLoadBalancer struct {
	ID            string
	Name          string
	State         string
	Type          string
	IPAddresses   []string
	OCID          string
	Subnets       []string
	NSGs          []string
	Created       *time.Time
	BackendSets   map[string]BackendSet
	BackendHealth map[string]string
	Listeners     map[string]string
	VcnID         string
	VcnName       string
	SubnetID      string
}

type BackendSet struct {
	Policy   string
	Backends []Backend
	Health   string
}

type Backend struct {
	Name   string
	Port   int
	Status string
}

type NetworkLoadBalancerRepository interface {
	GetNetworkLoadBalancer(ctx context.Context, ocid string) (*NetworkLoadBalancer, error)
	ListNetworkLoadBalancers(ctx context.Context, compartmentID string) ([]NetworkLoadBalancer, error)
	GetEnrichedNetworkLoadBalancer(ctx context.Context, ocid string) (*NetworkLoadBalancer, error)
	ListEnrichedNetworkLoadBalancers(ctx context.Context, compartmentID string) ([]NetworkLoadBalancer, error)
}
