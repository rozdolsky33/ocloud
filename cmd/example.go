package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/config"

	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
	Use:   "list-instances",
	Short: "List compute instances in the configured compartment",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := core.NewComputeClientWithConfigurationProvider(config.ConfigProvider)
		if err != nil {
			return err
		}

		resp, err := client.ListInstances(context.Background(), core.ListInstancesRequest{
			CompartmentId: &config.CompartmentID,
		})
		if err != nil {
			return err
		}

		for _, inst := range resp.Items {
			fmt.Println(*inst.DisplayName, *inst.Id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
