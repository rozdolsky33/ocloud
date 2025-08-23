package domain

import "context"

// Compartment represents a logical container for cloud resources.
// This is our application's internal representation, decoupled from the OCI SDK.
type Compartment struct {
	OCID                string
	DisplayName         string
	Description         string
	LifecycleState      string
	ParentCompartmentID string
}

// CompartmentRepository defines the port for interacting with compartment storage.
// The application layer will use this interface to talk to the infrastructure layer.
type CompartmentRepository interface {
	ListCompartments(ctx context.Context, parentCompartmentID string) ([]Compartment, error)
	GetCompartment(ctx context.Context, ocid string) (*Compartment, error)
}
