package auth

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rozdolsky33/ocloud/internal/app"
	"github.com/rozdolsky33/ocloud/internal/config"
)

// NewService creates a new authentication service.
func NewService(appCtx *app.ApplicationContext) *Service {
	return &Service{
		appCtx: appCtx,
		logger: appCtx.Logger,
	}
}

// InitAuth initializes the OCI authentication configuration.
func InitAuth() error {
	ConfigProvider = config.LoadOCIConfig()
	return nil
}

// PromptForProfile prompts the user to select an OCI profile.
func (s *Service) PromptForProfile() (string, error) {
	fmt.Println("Do you want to use the DEFAULT profile or enter a custom profile name?")
	fmt.Println("1: Use DEFAULT profile")
	fmt.Println("2: Enter custom profile name")
	fmt.Print("Enter your choice (1 or 2): ")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "reading profile choice input")
	}
	choice = strings.TrimSpace(choice)

	profile := "DEFAULT"
	if choice == "2" {
		fmt.Print("Enter profile name: ")
		customProfile, err := reader.ReadString('\n')
		if err != nil {
			return "", errors.Wrap(err, "reading custom profile input")
		}
		profile = strings.TrimSpace(customProfile)
		fmt.Printf("Using profile: %s\n", profile)
	} else {
		fmt.Println("Using DEFAULT profile")
	}

	return profile, nil
}

// GetOCIRegions returns a list of all available OCI regions.
func (s *Service) GetOCIRegions() []RegionInfo {
	// Static list of OCI regions
	regions := []string{
		"af-johannesburg-1", "ap-batam-1", "ap-chiyoda-1", "ap-chuncheon-1",
		"ap-chuncheon-2", "ap-dcc-canberra-1", "ap-dcc-gazipur-1", "ap-hyderabad-1",
		"ap-ibaraki-1", "ap-melbourne-1", "ap-mumbai-1", "ap-osaka-1", "ap-seoul-1",
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
		"uk-gov-cardiff-1", "uk-gov-london-1", "uk-london-1", "us-abilene-1",
		"us-ashburn-1", "us-chicago-1", "us-dallas-1", "us-gov-ashburn-1",
		"us-gov-chicago-1", "us-gov-phoenix-1", "us-langley-1", "us-luke-1",
		"us-phoenix-1", "us-saltlake-2", "us-sanjose-1", "us-somerset-1",
		"us-thames-1",
	}

	// Convert to RegionInfo
	var regionInfos []RegionInfo
	for i, r := range regions {
		regionInfos = append(regionInfos, RegionInfo{
			ID:   strconv.Itoa(i + 1),
			Name: r,
		})
	}

	return regionInfos
}

// PromptForRegion prompts the user to select an OCI region.
func (s *Service) PromptForRegion() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter region number or name: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "reading region input")
	}
	input = strings.TrimSpace(input)

	regions := s.GetOCIRegions()
	var chosen string

	if idx, err := strconv.Atoi(input); err == nil && idx >= 1 && idx <= len(regions) {
		chosen = regions[idx-1].Name
	} else {
		// Assume the input is a region name
		chosen = input
	}

	return chosen, nil
}

// Authenticate authenticates with OCI using the specified profile and region.
func (s *Service) Authenticate(profile, region string) (*AuthenticationResult, error) {
	// Authenticate via OCI CLI
	ociCmd := exec.Command("oci", "session", "authenticate", "--profile-name", profile, "--region", region)
	ociCmd.Stdout = os.Stdout
	ociCmd.Stderr = os.Stderr
	if err := ociCmd.Run(); err != nil {
		return nil, errors.Wrap(err, "failed to run `oci session authenticate`")
	}

	// Set environment variables
	os.Setenv("OCI_PROFILE", profile)
	os.Setenv("OCI_REGION", region)

	// Reload provider with the chosen profile / region
	if err := InitAuth(); err != nil {
		return nil, errors.Wrap(err, "reloading config after auth")
	}

	// Fetch root compartment (tenancy) OCID
	tenancyOCID, err := ConfigProvider.TenancyOCID()
	if err != nil {
		return nil, errors.Wrap(err, "fetching tenancy OCID")
	}

	// Create a result
	result := &AuthenticationResult{
		TenancyID: tenancyOCID,
		Profile:   profile,
		Region:    region,
	}

	// Try to get a tenancy name from a mapping file
	tenancies, err := config.LoadTenancyMap()
	if err == nil {
		for _, t := range tenancies {
			if t.TenancyID == tenancyOCID {
				result.TenancyName = t.Tenancy
				break
			}
		}
	}

	return result, nil
}

// GetCurrentEnvironment returns the current OCI environment variables.
func (s *Service) GetCurrentEnvironment() (*AuthenticationResult, error) {
	// Initialize auth
	if err := InitAuth(); err != nil {
		return nil, errors.Wrap(err, "initializing auth")
	}

	// Fetch root compartment (tenancy) OCID
	tenancyOCID, err := ConfigProvider.TenancyOCID()
	if err != nil {
		return nil, errors.Wrap(err, "fetching tenancy OCID")
	}

	// Create a result
	result := &AuthenticationResult{
		TenancyID: tenancyOCID,
		Profile:   config.GetOCIProfile(),
		Region:    os.Getenv("OCI_REGION"),
	}

	// Try to get a tenancy name from a mapping file
	tenancies, err := config.LoadTenancyMap()
	if err == nil {
		for _, t := range tenancies {
			if t.TenancyID == tenancyOCID {
				result.TenancyName = t.Tenancy
				break
			}
		}
	}

	return result, nil
}
