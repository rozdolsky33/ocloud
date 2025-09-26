package identity

import (
	"context"
	"time"
)

// Policy represents a set of rules defining access permissions for resources within a system.
type Policy struct {
	Name         string
	ID           string
	Statement    []string
	Description  string
	TimeCreated  time.Time
	FreeformTags map[string]string
	DefinedTags  map[string]map[string]interface{}
}

// PolicyRepository defines the port for interacting with policy storage.
type PolicyRepository interface {
	GetPolicy(ctx context.Context, ocid string) (*Policy, error)
	ListPolicies(ctx context.Context, compartmentID string) ([]Policy, error)
}
