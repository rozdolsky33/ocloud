package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/configuration/info"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"github.com/rozdolsky33/ocloud/scripts"

	"github.com/pkg/errors"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config"
)

// NewService creates a new authentication service.
func NewService() *Service {
	appCtx := &app.ApplicationContext{
		Logger: logger.Logger,
	}
	service := &Service{
		logger:   appCtx.Logger,
		Provider: config.LoadOCIConfig(),
	}
	logger.Logger.V(logger.Info).Info("Created new authentication service")
	return service
}

// Authenticate authenticates with OCI using the specified profile and region.
func (s *Service) Authenticate(profile, region string) (*AuthenticationResult, error) {
	logger.Logger.V(logger.Info).Info("Starting OCI authentication.", "profile", profile, "region", region)

	// Authenticate via OCI CLI
	ociCmd := exec.Command("oci", "session", "authenticate", "--profile-name", profile, "--region", region)
	ociCmd.Stdout = os.Stdout
	ociCmd.Stderr = os.Stderr

	logger.LogWithLevel(s.logger, logger.Trace, "Running OCI CLI command", "command", "oci session authenticate", "profile", profile, "region", region)

	if err := ociCmd.Run(); err != nil {
		return nil, errors.Wrap(err, "failed to run `oci session authenticate`")
	}

	logger.Logger.V(logger.Info).Info("OCI CLI authentication successful.")

	os.Setenv(flags.EnvKeyProfile, profile)
	os.Setenv(flags.EnvKeyRegion, region)

	logger.LogWithLevel(s.logger, logger.Trace, "Set environment variables", flags.EnvKeyProfile, profile, flags.EnvKeyRegion, region)

	tenancyOCID, err := s.Provider.TenancyOCID()
	if err != nil {
		return nil, errors.Wrap(err, "fetching tenancy OCID")
	}

	logger.LogWithLevel(s.logger, logger.Trace, "Fetched tenancy OCID", "tenancyOCID", tenancyOCID)

	result := &AuthenticationResult{
		TenancyID: tenancyOCID,
		Profile:   profile,
		Region:    region,
	}

	// Try to get a tenancy name from a mapping file
	logger.LogWithLevel(s.logger, logger.Trace, "Attempting to get tenancy name from mapping file")
	tenancies, err := config.LoadTenancyMap()
	if err != nil {
		logger.LogWithLevel(s.logger, logger.Trace, "Failed to load tenancy map, continuing without tenancy name", "error", err)
	} else {
		for _, t := range tenancies {
			if t.TenancyID == tenancyOCID {
				logger.LogWithLevel(s.logger, logger.Trace, "Found tenancy name in mapping file", "tenancy", t.Tenancy)
				result.TenancyName = t.Tenancy
				logger.LogWithLevel(s.logger, logger.Trace, "Set compartment name to tenancy name", "compartmentName", t.Tenancy)
				break
			}
		}
		logger.LogWithLevel(s.logger, logger.Trace, "No matching tenancy found in mapping file", "tenancyOCID", tenancyOCID)
	}

	logger.LogWithLevel(s.logger, logger.Debug, "Authentication successful", "profile", profile, "region", region, "tenancyID", tenancyOCID, "tenancyName", result.TenancyName)
	return result, nil
}

func (s *Service) promptForProfile() (string, error) {
	logger.Logger.V(logger.Info).Info("Prompting user for OCI profile selection.")

	useCustom := util.PromptYesNo("Do you want to enter a custom OCI profile name? (Default is 'DEFAULT')")

	profile := "DEFAULT"
	if useCustom {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter profile name: ")
		customProfile, err := reader.ReadString('\n')
		if err != nil {
			return "", errors.Wrap(err, "reading custom profile input")
		}
		profile = strings.TrimSpace(customProfile)

		logger.LogWithLevel(s.logger, logger.Debug, "Using custom profile", "profile", profile)
		fmt.Printf("Using profile: %s\n", profile)
	} else {
		logger.LogWithLevel(s.logger, logger.Debug, "Using DEFAULT profile")
		fmt.Println("Using DEFAULT profile")
	}

	return profile, nil
}

// GetOCIRegions returns a list of all available OCI regions.
func (s *Service) getOCIRegions() []RegionInfo {
	logger.Logger.V(logger.Info).Info("Fetching list of OCI regions.")

	regions := []string{
		"af-johannesburg-1", "ap-batam-1", "ap-chiyoda-1", "ap-chuncheon-1", "ap-chuncheon-2",
		"ap-dcc-canberra-1", "ap-dcc-gazipur-1", "ap-delhi-1", "ap-hyderabad-1", "ap-ibaraki-1",
		"ap-kulai-1", "ap-melbourne-1", "ap-mumbai-1", "ap-osaka-1", "ap-seoul-1",
		"ap-seoul-2", "ap-singapore-1", "ap-singapore-2", "ap-suwon-1", "ap-sydney-1",
		"ap-tokyo-1", "ca-montreal-1", "ca-toronto-1", "eu-amsterdam-1", "eu-crissier-1",
		"eu-dcc-dublin-1", "eu-dcc-dublin-2", "eu-dcc-milan-1", "eu-dcc-milan-2",
		"eu-dcc-rating-1", "eu-dcc-rating-2", "eu-dcc-zurich-1", "eu-frankfurt-1",
		"eu-frankfurt-2", "eu-jovanovac-1", "eu-madrid-1", "eu-madrid-2",
		"eu-marseille-1", "eu-milan-1", "eu-paris-1", "eu-stockholm-1",
		"eu-zurich-1", "il-jerusalem-1", "me-abudhabi-1", "me-abudhabi-2",
		"me-abudhabi-3", "me-abudhabi-4", "me-alain-1", "me-dcc-doha-1",
		"me-dcc-muscat-1", "me-dubai-1", "me-jeddah-1", "me-riyadh-1",
		"mx-monterrey-1", "mx-queretaro-1", "sa-bogota-1", "sa-santiago-1",
		"sa-saopaulo-1", "sa-valparaiso-1", "sa-vinhedo-1", "uk-cardiff-1",
		"uk-gov-cardiff-1", "uk-gov-london-1", "uk-london-1", "us-ashburn-1",
		"us-ashburn-2", "us-chicago-1", "us-gov-ashburn-1", "us-gov-chicago-1",
		"us-gov-phoenix-1", "us-langley-1", "us-luke-1", "us-phoenix-1",
		"us-saltlake-2", "us-sanjose-1", "us-somerset-1", "us-thames-1",
	}

	var regionInfos []RegionInfo
	for i, r := range regions {
		regionInfos = append(regionInfos, RegionInfo{
			ID:   strconv.Itoa(i + 1),
			Name: r,
		})
	}

	logger.LogWithLevel(s.logger, logger.Trace, "Retrieved OCI regions", "count", len(regionInfos))
	return regionInfos
}

// PromptForRegion prompts the user to select an OCI region.
func (s *Service) promptForRegion() (string, error) {
	logger.Logger.V(logger.Info).Info("Prompting user for OCI region selection.")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter region number or name: ")
	input, err := reader.ReadString('\n')

	if err != nil {
		return "", errors.Wrap(err, "reading region input")
	}

	input = strings.TrimSpace(input)
	logger.LogWithLevel(s.logger, logger.Trace, "User entered region input", "input", input)

	regions := s.getOCIRegions()
	var chosen string

	if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= len(regions) {
		chosen = regions[idx-1].Name
		logger.LogWithLevel(s.logger, logger.Trace, "Selected region by index", "index", idx, "region", chosen)
	} else {
		chosen = input
		logger.LogWithLevel(s.logger, logger.Trace, "Selected region by name", "region", chosen)
	}

	return chosen, nil
}

// viewConfigurationWithErrorHandling is a helper function to handle viewing configuration
// and handling common errors like a missing tenancy mapping file.
func (s *Service) viewConfigurationWithErrorHandling(realm string) error {
	err := info.ViewConfiguration(false, realm)
	if err != nil {
		if strings.Contains(err.Error(), "tenancy mapping file not found") {
			logger.LogWithLevel(s.logger, logger.Trace, "Tenancy mapping file not found, continuing without it", "error", err)
			return nil
		}
		return fmt.Errorf("viewing configuration: %w", err)
	}
	return nil
}

// performInteractiveAuthentication handles the interactive authentication process.
// It prompts the user for profile and region selection, authenticates with OCI,
// and returns the result of the authentication process.
func (s *Service) performInteractiveAuthentication(filter, realm string) (*AuthenticationResult, error) {
	logger.Logger.V(logger.Info).Info("Starting interactive authentication process.")
	profile, err := s.promptForProfile()
	if err != nil {
		return nil, fmt.Errorf("selecting profile: %w", err)
	}

	logger.Logger.V(logger.Info).Info("Profile selected successfully.", "profile", profile)

	err = s.viewConfigurationWithErrorHandling(realm)
	if err != nil {
		return nil, fmt.Errorf("viewing configuration: %w", err)
	}

	logger.LogWithLevel(s.logger, logger.Trace, "Getting OCI regions")
	regions := s.getOCIRegions()
	logger.LogWithLevel(s.logger, logger.Trace, "Displaying regions table", "regionCount", len(regions), "filter", filter)

	if err := DisplayRegionsTable(regions, filter); err != nil {
		return nil, fmt.Errorf("displaying regions: %w", err)
	}

	region, err := s.promptForRegion()
	if err != nil {
		return nil, fmt.Errorf("selecting region: %w", err)
	}

	fmt.Printf("Using region: %s\n", region)
	logger.LogWithLevel(s.logger, logger.Trace, "Region selected", "region", region)

	logger.LogWithLevel(s.logger, logger.Trace, "Authenticating with OCI", "profile", profile, "region", region)
	result, err := s.Authenticate(profile, region)

	if err != nil {
		return nil, fmt.Errorf("authenticating with OCI: %w", err)
	}

	logger.Logger.V(logger.Info).Info("OCI authentication successful.")

	err = s.viewConfigurationWithErrorHandling(realm)
	if err != nil {
		return nil, fmt.Errorf("viewing configuration: %w", err)
	}

	// Prompt for custom environment variables
	if util.PromptYesNo("Do you want to set OCI_TENANCY_NAME and OCI_COMPARTMENT?") {
		logger.Logger.V(logger.Info).Info("Prompting for custom environment variables.")
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Enter %s: ", flags.EnvKeyTenancyName)
		tenancy, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(s.logger, logger.Trace, "Error reading tenancy name input", "error", err)
		}

		fmt.Printf("Enter %s: ", flags.EnvKeyCompartment)
		compartment, err := reader.ReadString('\n')
		if err != nil {
			logger.LogWithLevel(s.logger, logger.Trace, "Error reading compartment input", "error", err)
		}

		tenancy = strings.TrimSpace(tenancy)
		compartment = strings.TrimSpace(compartment)

		logger.LogWithLevel(s.logger, logger.Trace, "Custom environment variables entered", "tenancyName", tenancy, "compartment", compartment)

		if tenancy != "" {
			result.TenancyName = tenancy
			logger.LogWithLevel(s.logger, logger.Trace, "Updated tenancy name", "tenancyName", tenancy)
		}

		if compartment != "" {
			result.CompartmentName = compartment
			logger.LogWithLevel(s.logger, logger.Trace, "Updated compartment", "compartment", compartment)
		}
		logger.Logger.V(logger.Info).Info("Custom environment variables set.")

	} else {
		logger.LogWithLevel(s.logger, logger.Trace, "Skipping variable setup")
		fmt.Println("\n Skipping variable setup.")
	}

	logger.LogWithLevel(s.logger, logger.Debug, "Interactive authentication completed successfully", "profile", profile, "region", region)
	return result, nil
}

// RunOCIAuthRefresher runs the OCI auth refresher script for the specified profile.
func (s *Service) runOCIAuthRefresher(profile string) error {
	logger.Logger.V(logger.Info).Info("Starting OCI auth refresher setup.", "profile", profile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	scriptDir := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCloudDefaultDirName, flags.OCloudScriptsDirName)
	if err := os.MkdirAll(scriptDir, 0o755); err != nil {
		return fmt.Errorf("failed to create script directory: %w", err)
	}

	scriptPath := fmt.Sprintf("%s/oci_auth_refresher.sh", scriptDir)

	// Write the embedded script bytes to the disk
	if err := os.WriteFile(scriptPath, scripts.OCIAuthRefresher, 0o700); err != nil {
		return fmt.Errorf("failed to write OCI auth refresher script to file: %w", err)
	}

	// Use a background context so it can run indefinitely
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("NOHUP=1 %s %s", scriptPath, profile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start OCI auth refresher script: %w", err)
	}

	pid := cmd.Process.Pid
	logger.LogWithLevel(logger.Logger, logger.Debug, "OCI auth refresher script started", "profile", profile, "pid", pid)
	// Write refresher PID to a profile session
	profileDir := filepath.Join(homeDir, flags.OCIConfigDirName, flags.OCISessionsDirName, profile)
	pidFile := filepath.Join(profileDir, flags.OCIRefresherPIDFileName)
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0o644); err != nil {
		return fmt.Errorf("failed to write OCI auth refresher script pid to file: %w", err)
	}

	fmt.Printf("\nOCI auth refresher started for profile %s with PID %d\n", profile, pid)
	fmt.Println("You can verify it's running with: pgrep -af oci_auth_refresher.sh")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nPress Enter to continue... ")
	_, _ = reader.ReadString('\n')
	return nil
}
