package domain

import "errors"
import "fmt"

var (
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("not found")
)

// NewNotFoundError creates a new error indicating that a resource was not found.
func NewNotFoundError(resourceType, resourceName string) error {
	return fmt.Errorf("%s '%s': %w", resourceType, resourceName, ErrNotFound)
}
