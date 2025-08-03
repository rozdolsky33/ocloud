package setup

import (
	"bufio"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

// NewService initializes a new Service instance with the provided application context.
// Returns a Service pointer.
func NewService() *Service {
	appCtx := &app.ApplicationContext{
		Logger: logger.Logger,
	}
	service := &Service{
		logger: appCtx.Logger,
	}
	return service
}

// ConfigureTenancyFile creates or updates a tenancy mapping configuration file with user-provided inputs.
func (s *Service) ConfigureTenancyFile() (err error) {

	logger.LogWithLevel(s.logger, 1, "Configuring tenancy map")
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".oci", ".ocloud")
	configFile := filepath.Join(configDir, "tenancy-map.yaml")

	var mappingFile []appConfig.MappingsFile

	logger.LogWithLevel(s.logger, 3, "Creating config directory", "dir", configDir)

	// Load existing records if a file exists
	if _, err := os.Stat(configFile); err == nil {
		logger.LogWithLevel(s.logger, 3, "Loading existing tenancy map")
		mappingFile, err = appConfig.LoadTenancyMap()
		if err != nil {
			return fmt.Errorf("loading existing tenancy map: %w", err)
		}
	} else {
		// File doesn't exist, prompt user if they want to create it
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\nTenancy mapping file not found at:", configFile)
		fmt.Print("Do you want to create the file and set up tenancy mapping? (y/n): ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Setup cancelled. Exiting.")
			return nil
		}

		//Create a directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		logger.LogWithLevel(s.logger, 3, "Creating new tenancy map")
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n--- Add a new tenancy record ---")

		type PromptField struct {
			name       string
			promptText string
			isMulti    bool
		}

		// Maintain the exact order of prompts as specified
		promptFields := []PromptField{
			{"environment", "Environment", false},
			{"tenancy", "Tenancy Name", false},
			{"tenancy_id", "Tenancy OCID", false},
			{"realm", "Realm", false},
			{"compartments", "Compartments (space-separated)", true},
			{"regions", "Regions (space-separated)", true},
		}

		// Collect values in the specified order
		values := make(map[string]interface{})

		for _, field := range promptFields {
			if field.isMulti {
				values[field.name] = promptMulti(reader, field.promptText)
			} else if field.name == "realm" {
				values[field.name] = promptWithValidation(reader, field.promptText, validateRealm)
			} else if field.name == "tenancy_id" {
				values[field.name] = promptWithValidation(reader, field.promptText, validateTenancyID)
			} else {
				values[field.name] = prompt(reader, field.promptText)
			}
		}

		// Create a record with fields in the same order as prompted
		record := appConfig.MappingsFile{
			Environment:  values["environment"].(string),
			Tenancy:      values["tenancy"].(string),
			TenancyID:    values["tenancy_id"].(string),
			Realm:        values["realm"].(string),
			Compartments: values["compartments"].([]string),
			Regions:      values["regions"].([]string),
		}
		mappingFile = append(mappingFile, record)

		more := strings.ToLower(prompt(reader, "Add another record? (y/n)"))
		if more == "n" || more == "no" {
			break
		}
	}
	// Write to a file
	logger.LogWithLevel(s.logger, 3, "Writing tenancy map to file")
	out, err := yaml.Marshal(mappingFile)
	if err != nil {
		return fmt.Errorf("marshalling tenancy map: %w", err)
	}
	err = os.WriteFile(configFile, out, 0644)
	if err != nil {
		return fmt.Errorf("writing tenancy map to file: %w", err)
	}
	logger.LogWithLevel(s.logger, 3, "Tenancy map written to file successfully", "file", configFile)

	return nil
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptMulti(reader *bufio.Reader, label string) []string {
	fmt.Printf("%s: ", label)
	text, _ := reader.ReadString('\n')
	return strings.Fields(strings.TrimSpace(text))
}

// validateRealm ensures the realm is properly formatted:
// - No more than 4 characters
// - Starts with "OC" (will be converted to uppercase)
// Returns the validated realm (uppercase) or an error if invalid
func validateRealm(realm string) (string, error) {
	// Convert to uppercase
	realm = strings.ToUpper(realm)

	// Check length
	if len(realm) > 4 {
		return "", fmt.Errorf("realm must be no more than 4 characters")
	}

	// Check if it starts with OC
	if len(realm) < 2 || realm[:2] != "OC" {
		return "", fmt.Errorf("realm must start with OC")
	}

	return realm, nil
}

// validateTenancyID ensures the tenancy ID contains the word "tenancy"
// Returns the tenancy ID or an error if invalid
func validateTenancyID(tenancyID string) (string, error) {
	if !strings.Contains(tenancyID, "tenancy") {
		return "", fmt.Errorf("tenancy ID must contain the word 'tenancy'")
	}
	return tenancyID, nil
}

// promptWithValidation prompts for input and validates it using the provided validation function
func promptWithValidation(reader *bufio.Reader, label string, validate func(string) (string, error)) string {
	for {
		fmt.Printf("%s: ", label)
		text, _ := reader.ReadString('\n')
		input := strings.TrimSpace(text)

		validated, err := validate(input)
		if err != nil {
			fmt.Printf("Error: %s. Please try again.\n", err)
			continue
		}

		return validated
	}
}
