package loadbalancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchableLoadBalancer_ToIndexable(t *testing.T) {
	lb := LoadBalancer{
		Name:            "Prod-LB",
		OCID:            "ocid1.lb.oc1..abc",
		Type:            "public",
		State:           "Active",
		VcnName:         "VCN-Prod",
		Shape:           "flexible",
		IPAddresses:     []string{"10.0.0.5", " 10.0.0.6 "},
		Hostnames:       []string{"www.Example.com", " api.example.com"},
		SSLCertificates: []string{"CertA", "CertB"},
		Subnets:         []string{"subnetA", "subnetB"},
	}

	s := SearchableLoadBalancer{lb}
	doc := s.ToIndexable()

	assert.Equal(t, "prod-lb", doc["Name"])
	assert.Equal(t, "ocid1.lb.oc1..abc", doc["OCID"])
	assert.Equal(t, "public", doc["Type"])
	assert.Equal(t, "active", doc["State"])
	assert.Equal(t, "vcn-prod", doc["VcnName"])
	assert.Equal(t, "flexible", doc["Shape"])
	assert.Contains(t, doc["IPAddresses"], "10.0.0.5")
	assert.Contains(t, doc["IPAddresses"], "10.0.0.6")
	assert.Contains(t, doc["Hostnames"], "www.example.com")
	assert.Contains(t, doc["Hostnames"], "api.example.com")
	assert.Contains(t, doc["SSLCertificates"], "certa")
	assert.Contains(t, doc["Subnets"], "subneta")
}
