package mapping

import (
	"testing"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/loadbalancer"
)

func TestNewDomainLoadBalancerFromAttrs_FullMapping(t *testing.T) {
	id := "ocid1.loadbalancer.oc1..lb"
	name := "prod-lb"
	shape := "flexible"
	isPrivate := false
	// IPs
	pubIP := "1.2.3.4"
	privIP := "10.0.0.5"
	pubTrue := true
	pubFalse := false
	created := common.SDKTime{Time: time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC)}

	// Listener
	lPort := 443
	proto := "TCP" // should be normalized to https due to port 443
	backendSetName := "bs1"
	routePolName := "rp1"

	// BackendSet HealthChecker
	hcProto := "http"
	hcPort := 80

	lb := loadbalancer.LoadBalancer{
		Id:             &id,
		DisplayName:    &name,
		LifecycleState: loadbalancer.LoadBalancerLifecycleStateActive,
		IsPrivate:      &isPrivate,
		ShapeName:      &shape,
		TimeCreated:    &created,
		IpAddresses: []loadbalancer.IpAddress{
			{IpAddress: &pubIP, IsPublic: &pubTrue},
			{IpAddress: &privIP, IsPublic: &pubFalse},
		},
		Listeners: map[string]loadbalancer.Listener{
			"https-listener": {
				Port:                  &lPort,
				Protocol:              &proto,
				DefaultBackendSetName: &backendSetName,
				RoutingPolicyName:     &routePolName,
				SslConfiguration:      &loadbalancer.SslConfiguration{},
			},
		},
		RoutingPolicies: map[string]loadbalancer.RoutingPolicy{
			"rp2": {},
		},
		BackendSets: map[string]loadbalancer.BackendSet{
			backendSetName: {
				Policy: &[]string{"ROUND_ROBIN"}[0],
				HealthChecker: &loadbalancer.HealthChecker{
					Protocol: &hcProto,
					Port:     &[]int{hcPort}[0],
				},
			},
		},
		SubnetIds:               []string{"ocid1.subnet.oc1..a", "ocid1.subnet.oc1..b"},
		NetworkSecurityGroupIds: []string{"ocid1.nsg.oc1..x"},
		Certificates:            map[string]loadbalancer.Certificate{"certA": {}},
		Hostnames:               map[string]loadbalancer.Hostname{"hn1": {Hostname: &[]string{"app.example.com"}[0]}},
	}

	attrs := NewLoadBalancerAttributesFromOCILoadBalancer(lb)
	dm := NewDomainLoadBalancerFromAttrs(attrs)

	if dm.OCID != id || dm.ID != id {
		t.Fatalf("id mapping failed: %+v", dm)
	}
	if dm.Name != name {
		t.Errorf("name mapping failed: %s", dm.Name)
	}
	if dm.Type != "Public" {
		t.Errorf("expected Public type, got %s", dm.Type)
	}
	if dm.Shape != shape {
		t.Errorf("shape mismatch: %s", dm.Shape)
	}
	if dm.Created == nil || dm.Created.IsZero() {
		t.Errorf("created time should be set")
	}
	if len(dm.IPAddresses) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(dm.IPAddresses))
	}
	if dm.IPAddresses[0] != "1.2.3.4 (public)" {
		t.Errorf("public ip annotation mismatch: %v", dm.IPAddresses)
	}
	if dm.IPAddresses[1] != "10.0.0.5 (private)" {
		t.Errorf("private ip annotation mismatch: %v", dm.IPAddresses)
	}

	// Listener string
	val, ok := dm.Listeners["https-listener"]
	if !ok {
		t.Fatalf("listener not mapped")
	}
	if val == "" || val[:5] != "https" {
		t.Errorf("listener proto/format unexpected: %s", val)
	}

	// Routing policies: implementation includes policies referenced by listeners (rp1);
	// it falls back to all LB policies only when none are referenced.
	if len(dm.RoutingPolicies) != 1 {
		t.Fatalf("routing policies expected 1, got %d", len(dm.RoutingPolicies))
	}
	if dm.RoutingPolicies[0] != "rp1" {
		t.Errorf("routing policies mismatch: %v", dm.RoutingPolicies)
	}

	// Backend set
	bs, ok := dm.BackendSets[backendSetName]
	if !ok {
		t.Fatalf("backend set missing")
	}
	if bs.Policy != "ROUND_ROBIN" {
		t.Errorf("backend policy mismatch: %s", bs.Policy)
	}
	if bs.Health != "HTTP:80" {
		t.Errorf("health label mismatch: %s", bs.Health)
	}

	// Certs imply UseSSL=true
	if !dm.UseSSL {
		t.Errorf("UseSSL expected true if ssl listener present/certificates exist")
	}
	if len(dm.SSLCertificates) != 1 || dm.SSLCertificates[0] != "certA" {
		t.Errorf("cert names mismatch: %v", dm.SSLCertificates)
	}

	// Hostnames collected and sorted
	if len(dm.Hostnames) != 1 || dm.Hostnames[0] != "app.example.com" {
		t.Errorf("hostnames mismatch: %v", dm.Hostnames)
	}

	// Subnets and NSGs copied
	if len(dm.Subnets) != 2 || dm.Subnets[0] == "" {
		t.Errorf("subnets mapping mismatch: %v", dm.Subnets)
	}
	if len(dm.NSGs) != 1 || dm.NSGs[0] == "" {
		t.Errorf("nsgs mapping mismatch: %v", dm.NSGs)
	}
}

func TestNewDomainLoadBalancerFromAttrs_DefaultsAndEmpties(t *testing.T) {
	attrs := &LoadBalancerAttributes{}
	dm := NewDomainLoadBalancerFromAttrs(attrs)
	if dm.ID != "" || dm.Name != "" {
		t.Errorf("expected empty id/name by default")
	}
	if dm.Type != "Public" {
		t.Errorf("default type should be Public when IsPrivate is nil")
	}
	if dm.Created != nil {
		t.Errorf("created should be nil when not provided")
	}
	if len(dm.Listeners) != 0 || len(dm.BackendSets) != 0 {
		t.Errorf("collections should be initialized (empty)")
	}
}
