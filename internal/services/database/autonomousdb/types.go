package autonomousdb

import (
	"github.com/rozdolsky33/ocloud/internal/domain"
)

// AutonomousDatabase represents an autonomous database instance with its attributes and connection details.
type AutonomousDatabase = domain.AutonomousDatabase

// IndexableAutonomousDatabase represents a simplified autonomous database structure indexed by its name.
type IndexableAutonomousDatabase struct {
	Name string
}
