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
	RoutingPolicies []string
	UseSSL          bool
	Hostnames       []string
	VcnID           string
	VcnName         string
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
	// Basic getters
	GetLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error)
	ListLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error)
	// Enriched getters (may perform additional API calls)
	GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*LoadBalancer, error)
	ListEnrichedLoadBalancers(ctx context.Context, compartmentID string) ([]LoadBalancer, error)
}
