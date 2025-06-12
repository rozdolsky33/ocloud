package resources

import (
	"context"
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/logger"
)

// ListInstances lists all instances in the specified compartment.
// This version is called from cmd/instance/root.go.
func ListInstances(ctx context.Context, compartmentID string) error {
	logger.Logger.V(1).Info("ListInstances()")
	fmt.Println("Inside Instance resources running List Instances")
	return nil
}

// FindInstances finds instances by name in the specified compartment.
// This version is called from cmd/instance/root.go.
func FindInstances(ctx context.Context, compartmentID, namePattern string, showImageDetails bool) error {
	logger.Logger.V(1).Info("FindInstances()", "namePattern", namePattern, "showImageDetails", showImageDetails)
	fmt.Println("Finding instances with name pattern:", namePattern)
	return nil
}
