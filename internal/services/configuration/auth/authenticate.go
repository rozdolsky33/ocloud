package auth

import (
	"fmt"

	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
)

// AuthenticateWithOCI handles the authentication process with Oracle Cloud Infrastructure (OCI) using interactive inputs.
// It performs authentication with the provided filter and realm, displays environment variables, and optionally starts the auth refresher.
func AuthenticateWithOCI(filter, realm string) error {
	s := NewService()
	logger.LogWithLevel(s.logger, logger.Debug, "Authenticating with OCI", "filter", filter, "realm", realm)
	result, err := s.performInteractiveAuthentication(filter, realm)
	if err != nil {
		return fmt.Errorf("performing interactive authentication: %w", err)
	}
	logger.LogWithLevel(s.logger, logger.Debug, "Interactive authentication completed", "tenancyID", result.TenancyID, "tenancyName", result.TenancyName)
	logger.LogWithLevel(s.logger, logger.Debug, "Authentication process completed successfully")
	logger.LogWithLevel(s.logger, logger.Debug, "Starting OCI auth refresher for profile", "profile", result.Profile)
	logger.CmdLogger.V(logger.Debug).Info("Prompting for OCI Auth Refresher setup...")

	if util.PromptYesNo("Do you want to set OCI_AUTH_AUTO_REFRESHER") {
		if err := s.runOCIAuthRefresher(result.Profile); err != nil {
			logger.LogWithLevel(s.logger, logger.Debug, "Failed to start OCI auth refresher", "error", err)
		}
		logger.LogWithLevel(s.logger, logger.Debug, "OCI auth refresher enabled")
	} else {
		logger.LogWithLevel(s.logger, logger.Debug, "OCI auth refresher disabled")
	}

	logger.LogWithLevel(s.logger, logger.Trace, "Displaying environment variables")
	if err = PrintExportVariable(result.Profile, result.TenancyName, result.CompartmentName); err != nil {
		return fmt.Errorf("printing export variables: %w", err)
	}

	return nil
}
