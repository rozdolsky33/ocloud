package compartment

import "github.com/rozdolsky33/ocloud/internal/domain"

// Compartment is an alias to the domain model, ensuring that the service layer
// uses the authoritative model from the domain layer. This provides a consistent
// type across the application while allowing for local extension if needed in the future.
type Compartment = domain.Compartment
