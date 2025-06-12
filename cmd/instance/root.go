package instance

import (
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/configuration"
	"github.com/rozdolsky33/ocloud/internal/helpers"
	"github.com/rozdolsky33/ocloud/internal/resources/compute"
	"github.com/spf13/cobra"
)

var InstanceCmd = &cobra.Command{
	Use:     "instance",
	Short:   "Find and list OCI instances",
	PreRunE: preConfigE,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		app, err := configuration.NewAppContext(ctx, cmd, args)
		if err != nil {
			return err
		}
		fmt.Println("Running instance command")
		fmt.Println("Compartment: ", app.CompartmentName)

		list, _ := cmd.Flags().GetBool("list")
		//find, _ := cmd.Flags().GetString("find")
		//details, _ := cmd.Flags().GetBool("image-details")

		switch {
		case list:
			compute.ListInstances()

		//case find != "":
		//	resources.FindInstances(
		//		app.ComputeCli,
		//		app.VnetCli,
		//		app.Compartment,
		//		find,
		//		details,
		//	)
		default:
			fmt.Println("Error: must pass --list or --find")
		}
		return nil
	},
}

func preConfigE(cmd *cobra.Command, args []string) error {
	if err := helpers.SetLogger(); err != nil {
		return err
	}
	helpers.InitLogger(helpers.CmdLogger)

	return nil
}
func init() {
	InstanceCmd.Flags().BoolP("list", "l", false, "List all instances")
	InstanceCmd.Flags().StringP("find", "f", "", "Search instances by name")
	InstanceCmd.Flags().BoolP("image-details", "i", false, "Show image details")
}
