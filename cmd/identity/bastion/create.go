package bastion

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rozdolsky33/ocloud/internal/app"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/spf13/cobra"
)

func NewCreateCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create",
		Aliases:       []string{"c"},
		Short:         "Create a new Bastion",
		Long:          "Create a new Bastion by selecting from available options",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCreateCommand(cmd, appCtx)
		},
	}

	return cmd
}

// RunCreateCommand executes the create command with an interactive TUI
// It implements a two-step selection process:
// 1. First, the user selects a type (Bastion or Session)
// 2. Then, the user selects a specific bastion from a list based on the chosen type
func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Create a new bastion service to interact with bastion resources
	service := bastionSvc.NewService()

	// STEP 1: Type Selection
	// Show the type selection UI first (Bastion or Session)
	typeModel := NewTypeSelectionModel()
	typeProgram := tea.NewProgram(typeModel)

	// Run the type selection TUI and wait for user selection
	typeResult, err := typeProgram.Run()
	if err != nil {
		fmt.Println("Error running type selection TUI:", err)
		os.Exit(1)
	}

	// Process the type selection result
	// If the user didn't make a selection or cancelled, exit early
	typeModel, ok := typeResult.(TypeSelectionModel)
	if !ok || typeModel.Choice == "" {
		// User cancelled the operation
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Check if the user selected the Bastion type
	if typeModel.Choice == TypeBastion {
		// Show the "Under Construction" animation
		ShowConstructionAnimation()
		return nil
	}

	// STEP 2: Bastion Selection
	// Gets the appropriate list of bastions based on the selected type
	var bastions []bastionSvc.Bastion
	if typeModel.Choice == TypeSession {
		// For Session type, use the dummy bastions
		// In a real implementation; this might fetch session-specific bastions
		bastions = service.GetDummyBastions()
	} else {
		// For Bastion type, also use dummy bastions for now
		// In a real implementation; this might fetch a different set of bastions
		bastions = service.GetDummyBastions()
	}

	// Initialize and show the bastion selection UI with the appropriate list
	bastionModel := NewBastionModel(bastions)
	bastionProgram := tea.NewProgram(bastionModel)

	// Run the bastion selection TUI and wait for user selection
	bastionResult, err := bastionProgram.Run()
	if err != nil {
		fmt.Println("Error running bastion selection TUI:", err)
		os.Exit(1)
	}

	// Process the bastion selection result
	if bastionModel, ok := bastionResult.(BastionModel); ok && bastionModel.Choice != "" {
		// User made a selection - find the selected bastion by ID from our list
		var selected bastionSvc.Bastion
		for _, b := range bastions {
			if b.ID == bastionModel.Choice {
				selected = b
				break
			}
		}

		// STEP 3: Session Type Selection
		// Initialize and show the session type selection UI
		sessionTypeModel := NewSessionTypeModel(selected.ID)
		sessionTypeProgram := tea.NewProgram(sessionTypeModel)

		// Run the session type selection TUI and wait for user selection
		sessionTypeResult, err := sessionTypeProgram.Run()
		if err != nil {
			fmt.Println("Error running session type selection TUI:", err)
			os.Exit(1)
		}

		// Process the session type selection result
		sessionTypeModel, ok := sessionTypeResult.(SessionTypeModel)
		if !ok || sessionTypeModel.Choice == "" {
			// User cancelled the operation
			fmt.Println("Operation cancelled.")
			return nil
		}

		// Check if the user selected Managed SSH
		if sessionTypeModel.Choice == TypeManagedSSH {
			// Show the "Under Construction" animation for Managed SSH
			ShowConstructionAnimation()
			return nil
		}

		// STEP 4: Target Type Selection (only for Port-Forwarding)
		if sessionTypeModel.Choice == TypePortForwarding {
			// Initialize and show the target type selection UI
			targetTypeModel := NewTargetTypeModel(selected.ID)
			targetTypeProgram := tea.NewProgram(targetTypeModel)

			// Run the target type selection TUI and wait for user selection
			targetTypeResult, err := targetTypeProgram.Run()
			if err != nil {
				fmt.Println("Error running target type selection TUI:", err)
				os.Exit(1)
			}

			// Process the target type selection result
			targetTypeModel, ok := targetTypeResult.(TargetTypeModel)
			if !ok || targetTypeModel.Choice == "" {
				// User cancelled the operation
				fmt.Println("Operation cancelled.")
				return nil
			}

			// Check if the user selected Database
			if targetTypeModel.Choice == TargetDatabase {
				// Show the "Under Construction" animation for Database
				ShowConstructionAnimation()
				return nil
			}

			// For OKE, display the final result
			fmt.Printf("\n---\nCreated %s Session on Bastion: %s\nID: %s\nTarget: %s\n",
				sessionTypeModel.Choice, selected.Name, selected.ID, targetTypeModel.Choice)
		}
	} else {
		// User cancelled the operation during bastion selection
		fmt.Println("Operation cancelled.")
	}

	return nil
}
