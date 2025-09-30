package loadbalancer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/rozdolsky33/ocloud/internal/app"
	domain "github.com/rozdolsky33/ocloud/internal/domain/network/loadbalancer"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// fakeRepo is a minimal in-memory implementation of LoadBalancerRepository for tests
type fakeRepo struct {
	items []domain.LoadBalancer
}

func (f *fakeRepo) GetLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
	return nil, nil
}
func (f *fakeRepo) ListLoadBalancers(ctx context.Context, compartmentID string) ([]domain.LoadBalancer, error) {
	return f.items, nil
}
func (f *fakeRepo) GetEnrichedLoadBalancer(ctx context.Context, ocid string) (*domain.LoadBalancer, error) {
	return nil, nil
}
func (f *fakeRepo) ListEnrichedLoadBalancers(ctx context.Context, compartmentID string) ([]domain.LoadBalancer, error) {
	return f.items, nil
}

func Test_mapToIndexableLoadBalancer_LowercasesFields(t *testing.T) {
	lb := LoadBalancer{
		Name:            "Prod-LB-01",
		Type:            "Public",
		VcnName:         "Main-VCN",
		Hostnames:       []string{"App.EXAMPLE.com", "API.example.COM"},
		SSLCertificates: []string{"MyCert", "AnotherCert"},
		Subnets:         []string{"Subnet-A (10.0.0.0/24)", "Subnet-B (10.0.1.0/24)"},
	}

	idx := mapToIndexableLoadBalancer(lb)

	// Marshal to JSON and then into a generic map for easy assertions
	b, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("failed to marshal index: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("failed to unmarshal index: %v", err)
	}

	if m["Name"] != "prod-lb-01" {
		t.Errorf("Name not lowercased, got %v", m["Name"])
	}
	if m["Type"] != "public" {
		t.Errorf("Type not lowercased, got %v", m["Type"])
	}
	if m["VcnName"] != "main-vcn" {
		t.Errorf("VcnName not lowercased, got %v", m["VcnName"])
	}
	// Hostnames
	hosts, _ := m["Hostnames"].([]any)
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hostnames, got %d", len(hosts))
	}
	for _, v := range hosts {
		vs, _ := v.(string)
		if vs != "app.example.com" && vs != "api.example.com" {
			t.Errorf("Hostnames not lowercased: %v", hosts)
		}
	}
	// Certificates
	certs, _ := m["SSLCertificates"].([]any)
	for _, v := range certs {
		vs, _ := v.(string)
		if vs != "mycert" && vs != "anothercert" {
			t.Errorf("SSLCertificates not lowercased: %v", certs)
		}
	}
	// Subnets
	subs, _ := m["Subnets"].([]any)
	for _, v := range subs {
		vs, _ := v.(string)
		if vs != "subnet-a (10.0.0.0/24)" && vs != "subnet-b (10.0.1.0/24)" {
			t.Errorf("Subnets not lowercased: %v", subs)
		}
	}
}

func TestService_Find_ByNameSubstring_CaseInsensitive(t *testing.T) {
	items := []domain.LoadBalancer{
		{ID: "1", Name: "prod-web-lb"},
		{ID: "2", Name: "stage-api-lb"},
		{ID: "3", Name: "PROD-db-lb"},
	}

	repo := &fakeRepo{items: items}
	appCtx := &app.ApplicationContext{CompartmentID: "ocid1.compartment.oc1..test", Logger: logger.CmdLogger}
	service := NewService(repo, appCtx)

	res, err := service.Find(context.Background(), "prod")
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}

	// ensure returned items correspond to those with prod in name (case-insensitive)
	ids := map[string]bool{}
	for _, r := range res {
		ids[r.ID] = true
	}
	if !ids["1"] || !ids["3"] {
		t.Errorf("unexpected results IDs: %v", ids)
	}
}
