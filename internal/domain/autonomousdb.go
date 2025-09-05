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
	DatabaseTags      ResourceTags
}

// AutonomousDatabaseRepository defines the interface for interacting with Autonomous Database data.
type AutonomousDatabaseRepository interface {
	ListAutonomousDatabases(ctx context.Context, compartmentID string) ([]AutonomousDatabase, error)
	FindAutonomousDatabase(ctx context.Context, compartmentID, name string) (*AutonomousDatabase, error)
	GetAutonomousDatabase(ctx context.Context, ocid string) (*AutonomousDatabase, error)
}
