package cmd

import (
	"context"
	"fmt"
	"github.com/rozdolsky33/ocloud/cmd/configuration"
	"github.com/rozdolsky33/ocloud/cmd/instance"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ocloud",
	Short: "Interact with Oracle Cloud Infrastructure",
	Long:  "",
}

func init() {
	configuration.InitGlobalFlags(rootCmd)
	rootCmd.AddCommand(instance.InstanceCmd)
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
