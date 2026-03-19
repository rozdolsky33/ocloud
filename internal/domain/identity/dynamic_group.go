package identity

import (
	"context"
	"time"
)

// DynamicGroup represents a group of instances that meet certain criteria.
type DynamicGroup struct {
	OCID           string
	Name           string
	Description    string
	MatchingRule   string
	LifecycleState string
	TimeCreated    time.Time
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
}

// DynamicGroupRepository defines the port for interacting with dynamic group storage.
type DynamicGroupRepository interface {
	GetDynamicGroup(ctx context.Context, ocid string) (*DynamicGroup, error)
	ListDynamicGroups(ctx context.Context, compartmentID string) ([]DynamicGroup, error)
}
