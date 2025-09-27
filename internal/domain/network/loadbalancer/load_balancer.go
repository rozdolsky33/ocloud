package loadbalancer

import (
	"context"
	"time"
)

type LoadBalancer struct {
	ID              string
	Name            string
	State           string
	Type            string
	IPAddresses     []string
	Shape           string
	Listeners       map[string]string
	BackendHealth   map[string]string
	OCID            string
	Subnets         []string
	NSGs            []string
	Created         *time.Time
	BackendSets     map[string]BackendSet
	SSLCertificates []string
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

type LoadBalancerRepository interface {
	GetLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error)
	ListLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error)
}
