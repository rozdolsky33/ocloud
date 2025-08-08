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
// It allows users to select a bastion from a list using keyboard navigation
func RunCreateCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Create a new bastion service
	service := bastionSvc.NewService()

	// Get dummy bastions from the service
	bastions := service.GetDummyBastions()

	// Initialize the Bubble Tea TUI model with our bastions
	model := NewBastionModel(bastions)
	p := tea.NewProgram(model)

	// Run the TUI and wait for user selection
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}

	// Process the result after TUI exits
	if m, ok := m.(BastionModel); ok && m.Choice != "" {
		// User made a selection - find the selected bastion by ID
		var selected bastionSvc.Bastion
		for _, b := range bastions {
			if b.ID == m.Choice {
				selected = b
				break
			}
		}

		// Display the result of the bastion creation
		fmt.Printf("\n---\nCreated Bastion: %s\nID: %s\n", selected.Name, selected.ID)
	} else {
		// User cancelled the operation
		fmt.Println("Operation cancelled.")
	}

	return nil
}
