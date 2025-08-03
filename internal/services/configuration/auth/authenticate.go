package auth

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/logger"
)

// AuthenticateWithOCI handles the authentication process with Oracle Cloud Infrastructure (OCI) using interactive inputs.
// It performs authentication with the provided filter and realm, displays environment variables, and optionally starts the auth refresher.
// Returns an error if any step in the process fails.
func AuthenticateWithOCI(filter, realm string) error {

	s := NewService()

	logger.LogWithLevel(s.logger, 1, "Authenticating with OCI", "filter", filter, "realm", realm)

	result, err := s.performInteractiveAuthentication(filter, realm)
	if err != nil {
		return fmt.Errorf("performing interactive authentication: %w", err)
	}

	logger.LogWithLevel(s.logger, 3, "Interactive authentication completed", "tenancyID", result.TenancyID, "tenancyName", result.TenancyName)

	logger.LogWithLevel(s.logger, 3, "Displaying environment variables")
	if err = PrintExportVariable(result.Profile, result.TenancyName, result.CompartmentName); err != nil {
		return fmt.Errorf("printing export variables: %w", err)
	}

	logger.LogWithLevel(s.logger, 1, "Authentication process completed successfully")

	logger.LogWithLevel(s.logger, 1, "Starting OCI auth refresher for profile", "profile", result.Profile)

	if s.promptYesNo("Do you want to set OCI_AUTH_AUTO_REFRESHER") {
		if err := s.runOCIAuthRefresher(result.Profile); err != nil {
			logger.LogWithLevel(s.logger, 1, "Failed to start OCI auth refresher", "error", err)
		}
		logger.LogWithLevel(s.logger, 1, "OCI auth refresher enabled")
	} else {
		logger.LogWithLevel(s.logger, 1, "OCI auth refresher disabled")
	}

	return nil
}
