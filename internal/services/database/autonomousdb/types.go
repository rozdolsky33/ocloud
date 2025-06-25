package autonomousdb

import (
	"github.com/go-logr/logr"
	"github.com/oracle/oci-go-sdk/v65/database"
)

// Service provides operations and functionalities related to database management, logging, and compartment handling.
type Service struct {
	dbClient      database.DatabaseClient
	logger        logr.Logger
	compartmentID string
}

// AutonomousDatabase represents an autonomous database instance with its attributes and connection details.
type AutonomousDatabase struct {
	Name              string
	ID                string
	PrivateEndpoint   string
	PrivateEndpointIp string
	ConnectionStrings map[string]string
	Profiles          []database.DatabaseConnectionStringProfile
}

// IndexableAutonomousDatabase represents a simplified autonomous database structure indexed by its name.
type IndexableAutonomousDatabase struct {
	Name string
}
