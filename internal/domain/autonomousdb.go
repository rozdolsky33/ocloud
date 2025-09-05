package domain

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/database"
)

// AutonomousDatabase represents an autonomous database instance with its attributes and connection details.
type AutonomousDatabase struct {
	Name              string
	ID                string
	PrivateEndpoint   string
	PrivateEndpointIp string
	ConnectionStrings map[string]string
	Profiles          []database.DatabaseConnectionStringProfile
	FreeformTags      map[string]string
	DefinedTags       map[string]map[string]interface{}
}

// AutonomousDatabaseRepository defines the interface for interacting with Autonomous Database data.
type AutonomousDatabaseRepository interface {
	GetAutonomousDatabase(ctx context.Context, ocid string) (*AutonomousDatabase, error)
	ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]AutonomousDatabase, error)
}
