package domain

import "context"

// Compartment represents a logical container for cloud resources.
// This is our application's internal representation, decoupled from the OCI SDK.
type Compartment struct {
	OCID           string
	DisplayName    string
	Description    string
	LifecycleState string
	FreeformTags   map[string]string
	DefinedTags    map[string]map[string]interface{}
}

// CompartmentRepository defines the port for interacting with compartment storage.
// The application layer will use this interface to talk to the infrastructure layer.
type CompartmentRepository interface {
	GetCompartment(ctx context.Context, ocid string) (*Compartment, error)
	ListCompartments(ctx context.Context, ocid string) ([]Compartment, error)
}
