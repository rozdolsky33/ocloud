package setup

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rozdolsky33/ocloud/internal/app"
	appConfig "github.com/rozdolsky33/ocloud/internal/config"
	"github.com/rozdolsky33/ocloud/internal/config/flags"
	"github.com/rozdolsky33/ocloud/internal/logger"
	"github.com/rozdolsky33/ocloud/internal/services/util"
	"gopkg.in/yaml.v3"
)

// ErrCancelled is returned when user cancels an operation via Ctrl+C
var ErrCancelled = errors.New("operation cancelled by user")

// NewService initializes a new Service instance with the provided application context.
func NewService() *Service {
	appCtx := &app.ApplicationContext{
		Logger: logger.Logger,
	}
	service := &Service{
		logger: appCtx.Logger,
	}
	logger.Logger.V(logger.Debug).Info("Creating new configuration setup service.")
	return service
}

// ConfigureTenancyFile creates or updates a tenancy mapping configuration file with user-provided inputs.
// The context allows for graceful cancellation via Ctrl+C.
func (s *Service) ConfigureTenancyFile(ctx context.Context) (err error) {
	logger.Logger.V(logger.Debug).Info("Starting tenancy map configuration.")

	// Set up signal handling for Ctrl+C
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\n\nOperation cancelled by user. Exiting...")
			cancel()
		case <-ctx.Done():
		}
	}()
	defer signal.Stop(sigChan)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting user home directory: %w", err)
	}

	configDir := filepath.Join(home, flags.OCIConfigDirName, flags.OCloudDefaultDirName)
	configFile := filepath.Join(configDir, flags.TenancyMapFileName)

	var mappingFile []appConfig.MappingsFile

	logger.LogWithLevel(s.logger, logger.Trace, "Creating config directory", "dir", configDir)

	// Load existing records if a file exists
	if _, err := os.Stat(configFile); err == nil {
		logger.LogWithLevel(s.logger, logger.Trace, "Loading existing tenancy map")
		mappingFile, err = appConfig.LoadTenancyMap()
		if err != nil {
			return fmt.Errorf("loading existing tenancy map: %w", err)
		}
	} else {
		fmt.Println("\nTenancy mapping file not found at:", configFile)
		logger.Logger.V(logger.Debug).Info("Tenancy mapping file not found.", "path", configFile)
		if !util.PromptYesNo("Do you want to create the file and set up tenancy mapping?") {
			fmt.Println("Setup cancelled. Exiting.")
			return nil
		}
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0o755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		logger.Logger.V(logger.Debug).Info("Configuration directory created.", "dir", configDir)

		logger.LogWithLevel(s.logger, logger.Trace, "Creating new tenancy map")
	}

	reader := bufio.NewReader(os.Stdin)

	logger.Logger.V(logger.Debug).Info("Prompting for new tenancy records.")
	for {
		// Check if context was cancelled before starting a new record
		select {
		case <-ctx.Done():
			return ErrCancelled
		default:
		}

		fmt.Println("\t--- Add a new tenancy record ---")

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
			var err error
			if field.isMulti {
				values[field.name], err = promptMulti(ctx, reader, field.promptText)
			} else if field.name == "realm" {
				values[field.name], err = promptWithValidation(ctx, reader, field.promptText, validateRealm)
			} else if field.name == "tenancy_id" {
				values[field.name], err = promptWithValidation(ctx, reader, field.promptText, validateTenancyID)
			} else {
				values[field.name], err = prompt(ctx, reader, field.promptText)
			}
			if err != nil {
				if errors.Is(err, ErrCancelled) {
					return ErrCancelled
				}
				return fmt.Errorf("reading input for %s: %w", field.name, err)
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
		// Display a record before saving it to the file
		fmt.Println("\t--- Record ---")
		out, err := yaml.Marshal(record)
		if err != nil {
			return fmt.Errorf("marshalling tenancy map: %w", err)
		}
		fmt.Println(string(out))
		if util.PromptYesNo("Do you want to add this record to the tenancy map?") {
			mappingFile = append(mappingFile, record)
			logger.Logger.V(logger.Info).Info("Record added to tenancy map.")
		} else {
			fmt.Println("Record discarded")
			logger.Logger.V(logger.Info).Info("Record discarded.")
		}
		if !util.PromptYesNo("Do you want to add another record?") {
			break
		}
	}
	// Write to a file
	logger.LogWithLevel(s.logger, logger.Trace, "Writing tenancy map to file")
	out, err := yaml.Marshal(mappingFile)
	if err != nil {
		return fmt.Errorf("marshalling tenancy map: %w", err)
	}
	err = os.WriteFile(configFile, out, 0644)
	if err != nil {
		return fmt.Errorf("writing tenancy map to file: %w", err)
	}
	logger.Logger.V(logger.Info).Info("Tenancy map written to file successfully.", "file", configFile)

	return nil
}

// prompt reads user input from the provided reader with a label and returns the trimmed input as a string.
// Returns ErrCancelled if the context is cancelled during input.
func prompt(ctx context.Context, reader *bufio.Reader, label string) (string, error) {
	logger.Logger.V(logger.Debug).Info("Prompting for input.", "label", label)

	// Check context before prompting
	select {
	case <-ctx.Done():
		return "", ErrCancelled
	default:
	}

	fmt.Printf("%s: ", label)

	// Read input in a goroutine so we can monitor context cancellation
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		text, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- strings.TrimSpace(text)
	}()

	select {
	case <-ctx.Done():
		return "", ErrCancelled
	case err := <-errChan:
		return "", err
	case text := <-resultChan:
		return text, nil
	}
}

// promptMulti reads a line of input for a given label and returns the input split into a slice of strings.
// Returns ErrCancelled if the context is cancelled during input.
func promptMulti(ctx context.Context, reader *bufio.Reader, label string) ([]string, error) {
	logger.Logger.V(logger.Debug).Info("Prompting for multi-input.", "label", label)

	// Check context before prompting
	select {
	case <-ctx.Done():
		return nil, ErrCancelled
	default:
	}

	fmt.Printf("%s: ", label)

	// Read input in a goroutine so we can monitor context cancellation
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		text, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- strings.TrimSpace(text)
	}()

	select {
	case <-ctx.Done():
		return nil, ErrCancelled
	case err := <-errChan:
		return nil, err
	case text := <-resultChan:
		return strings.Fields(text), nil
	}
}

// validateRealm ensures the realm is properly formatted
func validateRealm(realm string) (string, error) {
	realm = strings.ToUpper(realm)

	if len(realm) > 4 {
		return "", fmt.Errorf("realm must be no more than 4 characters")
	}

	if len(realm) < 2 || realm[:2] != "OC" {
		return "", fmt.Errorf("realm must start with OC")
	}

	return realm, nil
}

// validateTenancyID ensures the tenancy ID contains the word "tenancy"
func validateTenancyID(tenancyID string) (string, error) {
	if !strings.Contains(tenancyID, "tenancy") {
		return "", fmt.Errorf("tenancy ID must contain the word 'tenancy'")
	}
	return tenancyID, nil
}

// promptWithValidation prompts for input and validates it using the provided validation function.
// Returns ErrCancelled if the context is cancelled during input.
func promptWithValidation(ctx context.Context, reader *bufio.Reader, label string, validate func(string) (string, error)) (string, error) {
	logger.Logger.V(logger.Debug).Info("Prompting for input with validation.", "label", label)
	for {
		// Check context before each prompt attempt
		select {
		case <-ctx.Done():
			return "", ErrCancelled
		default:
		}

		fmt.Printf("%s: ", label)

		// Read input in a goroutine so we can monitor context cancellation
		resultChan := make(chan string, 1)
		errChan := make(chan error, 1)

		go func() {
			text, err := reader.ReadString('\n')
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- strings.TrimSpace(text)
		}()

		var input string
		select {
		case <-ctx.Done():
			return "", ErrCancelled
		case err := <-errChan:
			return "", err
		case input = <-resultChan:
		}

		validated, err := validate(input)
		if err != nil {
			fmt.Printf("Error: %s. Please try again.\n", err)
			logger.Logger.V(logger.Debug).Info("Validation failed.", "error", err)
			continue
		}

		logger.Logger.V(logger.Debug).Info("Validation successful.", "value", validated)
		return validated, nil
	}
}
