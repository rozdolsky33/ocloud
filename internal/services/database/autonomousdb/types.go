package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/domain/database"
)

// AutonomousDatabase represents an autonomous database instance with its attributes and connection details.
type AutonomousDatabase = database.AutonomousDatabase

// IndexableAutonomousDatabase represents a simplified autonomous database structure indexed by its name.
type IndexableAutonomousDatabase struct {
	Name string
}
