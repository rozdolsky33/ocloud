package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rozdolsky33/ocloud/cmd/instance"
	"github.com/rozdolsky33/ocloud/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "ocloud",
	Short: "Interact with Oracle Cloud Infrastructure",
	Long:  "",
}

func init() {
	config.InitGlobalFlags(rootCmd)
	rootCmd.AddCommand(instance.InstanceCmd)
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
