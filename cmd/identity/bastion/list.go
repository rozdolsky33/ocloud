package bastion

import (
	"github.com/rozdolsky33/ocloud/internal/app"
	bastionSvc "github.com/rozdolsky33/ocloud/internal/services/identity/bastion"
	"github.com/spf13/cobra"
)

func NewListCmd(appCtx *app.ApplicationContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list",
		Aliases:       []string{"l"},
		Short:         "List all Bastions",
		Long:          "",
		Example:       "",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunListCommand(cmd, appCtx)
		},
	}

	return cmd
}

func RunListCommand(cmd *cobra.Command, appCtx *app.ApplicationContext) error {
	// Create a new bastion service
	service := bastionSvc.NewService()

	// Get dummy bastions
	bastions := service.GetDummyBastions()

	// Print bastions
	bastionSvc.PrintBastions(bastions)

	return nil
}
